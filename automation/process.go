// TODO: nodes can be run in parallel since some nodes are completely isolated
package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/Dekr0/wwise-teller/integration"
	"github.com/Dekr0/wwise-teller/wwise"
)

const ProcessSpecVersion = 0

type ProcessNodeType uint8

const (
	RewireWithNewSourcesType       ProcessNodeType = 0
	RewireSoundsWithOldSourcesType ProcessNodeType = 1
	BasePropModifiersType          ProcessNodeType = 2
)

type ProcessSpec struct {
	Veriosn     uint8           `json:"version"`
	Pipelines []ProcessPipeline `json:"pipelines"`
}

type ProcessPipeline struct {
	Banks       []string                      `json:"banks"`
	Nodes       []ProcessNode                 `json:"nodes"`
	Integration   integration.IntegrationType `json:"integration"`
	Output        string                      `json:"output"`
}

type ProcessNode struct {
	Type ProcessNodeType `json:"type"`
	Spec string          `json:"spec"`
}

func ParseProcessSpec(fspec string) (*ProcessSpec, error) {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return nil, err
	}
	var spec ProcessSpec
	err = json.Unmarshal(blob, &spec)
	if err != nil {
		return nil, err
	}
	if spec.Veriosn != ProcessSpecVersion {
		return nil, fmt.Errorf("Version spec should be %d", ProcessSpecVersion)
	}
	return &spec, nil
}

func Process(ctx context.Context, fspec string) error {
	spec, err := ParseProcessSpec(fspec)
	if err != nil {
		return err
	}
	sem := make(chan struct{}, 3)
	var w sync.WaitGroup
	for _, p := range spec.Pipelines {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case sem <- struct{}{}:
			w.Add(1)
			go func() {
				if err := RunProcessPipeline(ctx, &p); err != nil {
					slog.Error("Failed to run process pipline", "error", err)
				}
				w.Done()
				<- sem
			}()
		default:
			if err := RunProcessPipeline(ctx, &p); err != nil {
				slog.Error("Failed to run process pipeline", "error", err)
			}
		}
	}
	return nil
}

func RunProcessPipeline(ctx context.Context, p *ProcessPipeline) error {
	return nil
}

func ProcessActiveBank(ctx context.Context, bnk *wwise.Bank, fspec string) error {
	spec, err := ParseProcessSpec(fspec)
	if err != nil {
		return err
	}
	if len(spec.Pipelines) == 0 {
		return nil
	}

	return ProcessPipelineActiveBank(ctx, bnk, &spec.Pipelines[0])
}

func ProcessPipelineActiveBank(ctx context.Context, bnk *wwise.Bank, p *ProcessPipeline) error {
	for _, node := range p.Nodes {
		switch node.Type {
		case RewireWithNewSourcesType:
			if err := RewireWithNewSources(ctx, bnk, node.Spec, false); err != nil {
				return err
			}
		case BasePropModifiersType:
			if err := ProcessBaseProps(bnk, node.Spec); err != nil {
				return err
			}
		default:
			slog.Warn(fmt.Sprintf("Unsupport process procedure %d", node.Type))
		}
	}
	return nil
}
