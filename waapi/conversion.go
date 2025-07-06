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

	"github.com/Dekr0/wwise-teller/utils"
	"github.com/google/uuid"
)

type ConversionFormatType uint8

const (
	ConversionFormatTypePCM      ConversionFormatType = 0
	ConversionFormatTypeADPCM    ConversionFormatType = 1
	ConversionFormatTypeVORBIS   ConversionFormatType = 2
	ConversionFormatTypeWEMOpus  ConversionFormatType = 3
)

var StagingCollision = errors.New("Staging folder's name collision.")
var NoConversionProj = errors.New("WWISETELLER_WPROJ enviromental variable (used for locating Wwise Project to automate conversion) is not set.")

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

var WEMCache = ""

func InitWEMCache() error {
	var err error
	WEMCache, err = os.MkdirTemp("", "wwise-teller-wem-cache")
	if err != nil {
		return err
	}
	return nil
}

func CleanWEMCache() {
	os.Remove(WEMCache)
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

	if utils.Tmp == "" {
		if err := utils.InitTmp(); err != nil {
			return "", err
		}
	}

	staging := filepath.Join(utils.Tmp, uuid.New().String())
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

	if utils.Tmp == "" {
		if err := utils.InitTmp(); err != nil {
			return "", err
		}
	}

	staging := filepath.Join(utils.Tmp, uuid.New().String())
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
	cli, err := GetWwiseCLI()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(
		ctx,
		cli,
		"convert-external-source", project,
		"--platform", "Windows",
		"--source-file", wsource,
		"--output", output,
	)
	res, err := cmd.CombinedOutput()
	for line := range bytes.SplitSeq(res, []byte{'\n'}) {
		if bytes.EqualFold(line, []byte{'\r'}) {
			continue
		}
		slog.Info(string(line))
	}
	return err
}

func GetWwiseCLI() (string, error) {
	wwiseRoot := os.Getenv("WWISEROOT")
	if !filepath.IsAbs(wwiseRoot) {
		return "", fmt.Errorf("Path specified in environmental variable WWISEROOT is not in absolute path.")
	}
	stat, err := os.Lstat(wwiseRoot)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("Path specified in environmental variable WWISEROOT is not a directory.")
	}
	cli := filepath.Join(wwiseRoot, "Authoring/x64/Release/bin/WwiseConsole.exe")
	stat, err = os.Lstat(cli)
	if err != nil {
		return "", err
	}
	return cli, nil
}

type WEMInfo struct {
	SampleRate  int
	NumChannels int
	ChannelMask int
	NumSamples  int
	Encoding    string
}

func GetVGMStream() string {
	if runtime.GOOS == "Windows" {
		return "vgmstream-cli.exe"
	} else {
		return "vgmstream-cli"
	}
}

func ExportWEMByte(ctx context.Context, wem []byte, wave bool) (string, error) {
	if WEMCache == "" {
		if err := InitWEMCache(); err != nil {
			return "", err
		}
	}
	wemFile := filepath.Join(WEMCache, uuid.NewString() + ".wem")
	err := os.WriteFile(wemFile, wem, 0777)
	if err != nil {
		return "", err
	}
	defer os.Remove(wemFile)
	return ExportWEMFile(wemFile, ctx, wave)
}

func ExportWEMFile(path string, ctx context.Context, wave bool) (string, error) {
	if WEMCache == "" {
		if err := InitWEMCache(); err != nil {
			return "", err
		}
	}
	cmdName := GetVGMStream()
	tmpWAV := filepath.Join(WEMCache, uuid.NewString() + ".wav")
	cmd := exec.CommandContext(ctx, cmdName, "-o", tmpWAV, path)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return tmpWAV, nil
}

func GetWEMInfoByte(wem []byte, ctx context.Context) (float32, float32, error) {
	tmpWEM := filepath.Join(utils.Tmp, uuid.NewString() + ".wem")
	err := os.WriteFile(tmpWEM, wem, 0777)
	if err != nil {
		return 0.0, 0.0, err
	}
	defer os.Remove(tmpWEM)
	return GetWEMInfoFile(tmpWEM, ctx)
}

func GetWEMInfoFile(wem string, ctx context.Context) (float32, float32, error) {
	cmdName := GetVGMStream() 
	cmd := exec.CommandContext(ctx, cmdName, "-m", wem)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return 0.0, 0.0, err
	}
	return 0.0, 0.0, nil
}
