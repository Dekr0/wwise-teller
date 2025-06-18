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

var StagingCollision = errors.New("Staging folder's name collision. Please remove all content in the .cache folder and try again.")

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

func CreateConversionList(
	ctx context.Context,
	in []string,
	out []string,
	conversion string,
	dry bool,
) (string, error) {
	if !dry && runtime.GOOS != "windows" {
		return "", fmt.Errorf("Wwise External Sources Conversion is only available on Windows")
	}
	var err error
	duplicate := make(map[string]struct{}, len(in))
	safeIn := make([]string, 0, len(in))
	for _, i := range in {
		_, err = os.Lstat(i)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("Wave file %s does not exist.", i))
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information about wave file %s", i), "error", err)
			}
		} else {
			if !filepath.IsAbs(i) {
				slog.Error(fmt.Sprintf("%s is not absolute path", i))
				continue
			}
			if _, in := duplicate[i]; in {
				return "", fmt.Errorf("The provided wave files contain duplicates! %s is duplicated", i)
			} else {
				duplicate[i] = struct{}{}
			}
			safeIn = append(safeIn, i)
		}
	}

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
		Sources: make([]ExternalSource, len(safeIn)),
	}

	suffixing := make(map[string]uint8)
	basename := ""
	for i := range list.Sources {
		basename = strings.Split(filepath.Base(safeIn[i]), ".")[0]
		if suffix, in := suffixing[basename]; !in {
			suffixing[basename] = 0
			basename = fmt.Sprintf("%s.wem", basename)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: safeIn[i],
				Conversion: &conversion,
				Destination: basename,
			}
			out[i] = dest
		} else {
			basename = fmt.Sprintf("%s_%d.wem", basename, suffix)
			dest := filepath.Join(destF, basename)
			list.Sources[i] = ExternalSource{
				Path: safeIn[i],
				Conversion: &conversion,
				Destination: basename,
			}
			out[i] = dest
			suffixing[basename] += 1
		}
	}

	x, err := xml.MarshalIndent(&list, "", "    ")
	if err != nil {
		return "", err
	}
	fmt.Println(string(x))

	wsource := filepath.Join(staging, "external_sources.wsources")
	if err := os.WriteFile(wsource, []byte(xml.Header + string(x)), 0777); err != nil {
		return "", err
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
