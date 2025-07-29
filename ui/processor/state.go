package processor

import (
	"log/slog"
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
	spec, err := automation.ParseProcessor(fspec)
	if err != nil {
		return err
	}

	for _, pipeline := range spec.Pipelines {
		for _, script := range pipeline.Scripts {
			scriptPath := filepath.Join(pipeline.BanksWorkspace, script.Script)
			if _, in := p.ProcessScripts[scriptPath]; in {
				continue
			}
			var err error
			switch script.Type {
			case automation.TypeRewireWithNewSources, automation.TypeReplaceAudioSources:
				s := automation.WavSoundMap{}
				if err = automation.DecodeReplaceAudioSourcesScript(&s, scriptPath); err != nil {
					slog.Error("Failed to decode replace audio ")

				}
				p.ProcessScripts[scriptPath] = &s
			case automation.TypeBasePropModifiers:
			case automation.TypeImportAsRanSeqCntr:
			case automation.TypeRanSeqModifiers:
			case automation.TypeNewSoundToRanSeqCntr:
			case automation.TypeBulkProcessBaseProp:
			}
		}
	}
	return nil
}
