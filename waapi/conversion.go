package waapi

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var StagingCollision = errors.New("Staging folder's name collision. Please remove all content in the .cache folder and try again.")

const CACHE = ".cache"

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

func WavToWem(
	ctx context.Context,
	in []string,
	out []string,
	project string,
	conversion string,
) (string, error) {
	_, err := os.Lstat(project)
	if err != nil {
		return "", err
	}

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
			if _, in := duplicate[i]; in {
				return "", fmt.Errorf("The provided wave files contain duplicates! %s is duplicated", i)
			} else {
				duplicate[i] = struct{}{}
			}
			safeIn = append(safeIn, i)
		}
	}

	namespace := uuid.New()
	info, err := os.Lstat(CACHE)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err := os.Mkdir(CACHE, 0777)
		if err != nil {
			return "", err
		}
	}
	if !info.IsDir() {
		return "", fmt.Errorf(".cache exists but it's a file")
	}

	staging := filepath.Join(CACHE, namespace.String())
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
		Root: CACHE,
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

	wsource := filepath.Join(staging, "external_sources.wsources")
	if err := os.WriteFile(wsource, []byte(xml.Header + string(x)), 0777); err != nil {
		return "", err
	}

	return wsource, nil
}
