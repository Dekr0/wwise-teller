package processor

import (
	"path/filepath"

	"github.com/Dekr0/wwise-teller/automation"
	"github.com/google/uuid"
)

type ProcessorEditor struct {
	Path            []string
	InMemory        []bool
	Processes       []automation.Processor
	ProcessScripts  map[string]automation.IProcessScript
}

func NewP(p *ProcessorEditor) {
	p.Path = make([]string, 0, 4)
	p.Processes = make([]automation.Processor, 0, 4)
	p.ProcessScripts = make(map[string]automation.IProcessScript, 4)
}

func New() (p ProcessorEditor) {
	p.Path = make([]string, 0, 4)
	p.Processes = make([]automation.Processor, 0, 4)
	p.ProcessScripts = make(map[string]automation.IProcessScript, 4)
	return p
}

func (p *ProcessorEditor) NewProcessor(name string) {
	p.Path = append(p.Path, filepath.Join(uuid.NewString(), name))
	p.InMemory = append(p.InMemory, true)
}

func (p *ProcessorEditor) LoadProcessor(fspec string) error {
	return nil
}
