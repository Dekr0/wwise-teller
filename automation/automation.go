package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dekr0/wwise-teller/integration"
	"github.com/Dekr0/wwise-teller/wwise"
)

const AutomationSpecVersion = 0

type AutomationSpec struct {
	Veriosn   uint8                `json:"version"`
	Pipelines []AutomationPipeline `json:"pipelines"`
}

type AutomationPipeline struct {
	Banks       []string                      `json:"banks"`
	Nodes       []AutomationNode              `json:"nodes"`
	Integration   integration.IntegrationType `json:"integration"`
	Output        string                      `json:"output"`
}

type AutomationNodeType uint8

func AutomateActiveBank(ctx context.Context, bnk *wwise.Bank, fspec string) error {
	spec, err := ParseAutomationSpec(fspec)
	if err != nil {
		return err
	}
	for _, pipeline := range spec.Pipelines {
		for _, node := range pipeline.Nodes {
			switch node.Type {
			case RewireWithNewSourcesType:
				if err := RewireWithNewSources(ctx, bnk, node.Spec, false); err != nil {
					return err
				}
			case BasePropModifiersType:
				if err := AutomateBaseProp(bnk, node.Spec); err != nil {
					return err
				}
			default:
				slog.Warn(fmt.Sprintf("Unsupport automation procedure %d", node.Type))
			}
		}
	}
	return nil
}

func ParseAutomationSpec(fspec string) (*AutomationSpec, error) {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return nil, err
	}
	var spec AutomationSpec
	err = json.Unmarshal(blob, &spec)
	if err != nil {
		return nil, err
	}
	if spec.Veriosn != AutomationSpecVersion {
		return nil, fmt.Errorf("Version spec should be %d", AutomationSpecVersion)
	}
	return &spec, nil
}

const (
	RewireWithNewSourcesType  AutomationNodeType = 0
	RewireSoundsWithOldSourcesType  AutomationNodeType = 1
	BasePropModifiersType AutomationNodeType = 2
)

type AutomationNode struct {
	Type AutomationNodeType `json:"type"`
	Spec string             `json:"spec"`
}
