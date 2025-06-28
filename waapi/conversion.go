package waapi

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

type ConversionFormatType uint8

const (
	ConversionFormatTypePCM      ConversionFormatType = 0
	ConversionFormatTypeADPCM    ConversionFormatType = 1
	ConversionFormatTypeVORBIS   ConversionFormatType = 2
	ConversionFormatTypeWEMOpus  ConversionFormatType = 3
)

var StagingCollision = errors.New("Staging folder's name collision. Please remove all content in the .cache folder and try again.")
var NoConversionProj = errors.New("WWISETELLER_WPROJ enviromental variable (used for locating Wwise Project to automate conversion) is not set.")

var Tmp string = ""

type ExternalSourcesList struct {
	XMLName       xml.Name       `xml:"ExternalSourcesList"`
	SchemaVersion uint8          `xml:"SchemaVersion,attr"`
	Root          string         `xml:"Root,attr"`
	Sources     []ExternalSource 
}

type ExternalSource struct {
	XMLName       xml.Name `xml:"Source"`
	Path          string   `xml:"Path,attr"`
	Conversion   *string   `xml:"Conversion,attr,omitempty"`
	Destination   string   `xml:"Destination,attr"`
	AnalysisType *uint8    `xml:"AnalysisTypes,attr,omitempty"`
}

func InitTmp() error {
	var err error
	Tmp, err = os.MkdirTemp("", "wwise-teller-")
	if err != nil {
		return err
	}
	return nil
}

func CleanTmp() error {
	return os.RemoveAll(Tmp)
}

func GetProject() (string, error) {
	proj := os.Getenv("WWISETELLER_WPROJ")
	if proj == "" {
		return "", NoConversionProj
	}
	if !filepath.IsAbs(proj) {
		return "", fmt.Errorf("File path of Wwise project %s is not in aboslute path.", proj)
	}
	stat, err := os.Lstat(proj)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", fmt.Errorf("File path of Wwise project %s is a directory.", proj)
	}
	return proj, nil
}

// Assume there's no duplicate wave files. 
// Assume all wave files exist.
// Assume all wave files are in full path.
func CreateConversionList[T any](
	ctx context.Context,
	wavsMap map[string]T,
	wemsMap map[string]T,
	conversion string,
	dry bool,
) (string, error) {
	if !dry && runtime.GOOS != "windows" {
		return "", fmt.Errorf("Wwise External Sources Conversion is only available on Windows")
	}
	var err error

	staging := filepath.Join(Tmp, uuid.New().String())
	_, err = os.Lstat(staging)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err := os.Mkdir(staging, 0777)
		if err != nil {
			return "", err
		}
	} else {
		return "", StagingCollision
	}
	destF := filepath.Join(staging, "Windows")

	list := ExternalSourcesList{
		SchemaVersion: 1,
		Root: staging,
		Sources: make([]ExternalSource, len(wavsMap)),
	}

	suffixing := make(map[string]uint8)
	basename := ""
	i := 0
	for wav, v := range wavsMap {
		basename = strings.Split(filepath.Base(wav), ".")[0]
		if suffix, in := suffixing[basename]; !in {
			suffixing[basename] = 0
			basename = fmt.Sprintf("%s.wem", basename)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: wav,
				Conversion: &conversion,
				Destination: basename,
			}
			if _, in := wemsMap[dest]; in {
				panic("Panic Trap")
			}
			wemsMap[dest] = v
		} else {
			basename = fmt.Sprintf("%s_%d.wem", basename, suffix)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: wav,
				Conversion: &conversion,
				Destination: basename,
			}
			if _, in := wemsMap[dest]; in {
				panic("Panic Trap")
			}
			wemsMap[dest] = v
			suffixing[basename] += 1
		}
		i += 1
	}

	x, err := xml.MarshalIndent(&list, "", "    ")
	if err != nil {
		return "", err
	}

	wsource := filepath.Join(staging, "external_sources.wsources")
	if !dry {
		if err := os.WriteFile(wsource, []byte(xml.Header + string(x)), 0777); err != nil {
			return "", err
		}
	} else {
		fmt.Println(xml.Header + string(x))
	}

	return wsource, nil
}

func CreateConversionListInOrder(
	ctx context.Context,
	wavsMap map[string]uint8,
	wems []string,
	conversion string,
	dry bool,
) (string, error) {
	if !dry && runtime.GOOS != "windows" {
		return "", fmt.Errorf("Wwise External Sources Conversion is only available on Windows")
	}
	var err error

	staging := filepath.Join(Tmp, uuid.New().String())
	_, err = os.Lstat(staging)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err := os.Mkdir(staging, 0777)
		if err != nil {
			return "", err
		}
	} else {
		return "", StagingCollision
	}
	destF := filepath.Join(staging, "Windows")

	list := ExternalSourcesList{
		SchemaVersion: 1,
		Root: staging,
		Sources: make([]ExternalSource, len(wavsMap)),
	}

	suffixing := make(map[string]uint8)
	basename := ""
	i := 0
	for wav, idx := range wavsMap {
		basename = strings.Split(filepath.Base(wav), ".")[0]
		if suffix, in := suffixing[basename]; !in {
			suffixing[basename] = 0
			basename = fmt.Sprintf("%s.wem", basename)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: wav,
				Conversion: &conversion,
				Destination: basename,
			}
			wems[idx] = dest
		} else {
			basename = fmt.Sprintf("%s_%d.wem", basename, suffix)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: wav,
				Conversion: &conversion,
				Destination: basename,
			}
			wems[idx] = dest
			suffixing[basename] += 1
		}
		i += 1
	}

	x, err := xml.MarshalIndent(&list, "", "    ")
	if err != nil {
		return "", err
	}

	wsource := filepath.Join(staging, "external_sources.wsources")
	if !dry {
		if err := os.WriteFile(wsource, []byte(xml.Header + string(x)), 0777); err != nil {
			return "", err
		}
	} else {
		fmt.Println(xml.Header + string(x))
	}

	return wsource, nil
}

func WwiseConversion(ctx context.Context, wsource string, project string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Wwise External Sources Conversion is only available on Windows")
	}
	if !filepath.IsAbs(wsource) {
		return fmt.Errorf("%s is not in the form of absolute path.", wsource)
	}
	if !filepath.IsAbs(project) {
		return fmt.Errorf("%s is not in the form of absolute path.", project)
	}

	_, err := os.Lstat(wsource)
	if err != nil {
		return err
	}
	_, err = os.Lstat(project)
	if err != nil {
		return err
	}
	output := filepath.Dir(wsource)
	cmd := exec.CommandContext(
		ctx,
		"WwiseConsole.exe",
		"convert-external-source", project,
		"--platform", "Windows",
		"--source-file", wsource,
		"--output", output,
	)
	res, err := cmd.CombinedOutput()
	for line := range bytes.SplitSeq(res, []byte{'\n'}) {
		slog.Info(string(line))
	}
	return err
}
