package script_editor

import (
	"github.com/Dekr0/wwise-teller/automation"
)

type ProcessorEditor struct {
	Path         []string
	Processes    []automation.Processor
	ScriptState  map[string]any
}

func (p *ProcessorEditor) LoadProcessor(fspec string) error {
	spec, err := automation.ParseProcessor(fspec)
	if err != nil {
		return err
	}
	for _, pipeline := range spec.Pipelines {
		for _, script := pipeline.Scripts {

		}
	}
}
