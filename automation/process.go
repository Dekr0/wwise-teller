// TODO: nodes can be run in parallel since some nodes are completely isolated
package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/Dekr0/wwise-teller/integration"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

const ProcessSpecVersion = 0

type ProcessScriptType uint8

// TODO: this will cause migration problem on the user side since they use 
// number to tag type
// No need to add placeholder creation because this can be achive by rewiring 
// with new sources, and then use type 1 integration
const (
	TypeRewireWithNewSources ProcessScriptType = iota
	TypeRewireWithOldSources
	TypeBasePropModifiers 
	TypeImportAsRanSeqCntr
	TypeReplaceAudioSources 
	TypeRanSeqModifiers
	TypeNewSoundToRanSeqCntr // Add new sound objects to a random / sequence container
	TypeBulkProcessBaseProp  // Apply one base property to all listed hierarchies
	TypeStreamTypeModifiers  // Change stream type of a sound object
	ProcessScriptTypeCount
)

type IProcessScript interface {
	Encode(string) error
}

type Processor struct {
	Veriosn     uint8           `json:"version"`
	Pipelines []ProcessPipeline `json:"pipelines"`
}

func (p *Processor) NewPipeline() {
	p.Pipelines = append(p.Pipelines, ProcessPipeline{
		BanksWorkspace: "",
		ScriptsWorkspace: "",
		Banks: make([]string, 0, 1),
		Scripts: make([]ProcessScript, 0, 4),
		Integration: integration.IntegrationTypeDefault,
		Output: "",
	})
}

type ProcessPipeline struct {
	BanksWorkspace   string                      `json:"banksWorkspace"`
	ScriptsWorkspace string                      `json:"scriptsWorkspace"`
	Banks          []string                      `json:"banks"`
	Scripts        []ProcessScript               `json:"scripts"`
	Integration      integration.IntegrationType `json:"integration"`
	Output           string                      `json:"output"`
}

func (p *ProcessPipeline) Script(name string) bool {
	return slices.ContainsFunc(p.Scripts, func(script ProcessScript) bool {
		return script.Script == name
	})
}

func (p *ProcessPipeline) NewScript(name string, t ProcessScriptType) {
	if p.Script(name) {
		return
	}
	p.Scripts = append(p.Scripts, ProcessScript{
		Script: name,
		Type: t,
	})
}

type ProcessScript struct {
	Type ProcessScriptType `json:"type"`
	Script string          `json:"script"`
}

// Perform validation and filtering as well
func ParseProcessor(fspec string) (*Processor, error) {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return nil, err
	}
	var spec Processor
	err = json.Unmarshal(blob, &spec)
	if err != nil {
		return nil, err
	}
	if spec.Veriosn != ProcessSpecVersion {
		return nil, fmt.Errorf("Version spec should be %d.", ProcessSpecVersion)
	}
	spec.Pipelines = slices.DeleteFunc(spec.Pipelines, func(p ProcessPipeline) bool {
		if !filepath.IsAbs(p.BanksWorkspace) {
			slog.Error(fmt.Sprintf("%s is not an absolute path.", p.BanksWorkspace))
			return true
		}
		stat, err := os.Lstat(p.BanksWorkspace)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("%s does not exist.", p.ScriptsWorkspace))
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information about %s", p.BanksWorkspace), "error", err)
			}
			return true
		}
		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("Bank workspace %s is not a directory.", p.BanksWorkspace))
			return true
		}

		if !filepath.IsAbs(p.ScriptsWorkspace) {
			slog.Error(fmt.Sprintf("%s is not an absolute path.", p.ScriptsWorkspace))
			return true
		}
		stat, err = os.Lstat(p.ScriptsWorkspace)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("%s does not exist.", p.ScriptsWorkspace))
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information about %s", p.ScriptsWorkspace), "error", err)
			}
			return true
		}
		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("Script workspace %s is not a directory.", p.ScriptsWorkspace))
			return true
		}

		if len(p.Banks) <= 0 {
			return true
		}
		if len(p.Scripts) <= 0 {
			return true
		}
		if p.Integration == integration.IntegrationTypeNone {
			return true
		}

		stat, err = os.Lstat(p.Output)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(p.Output, 0777)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to create output directory %s", p.Output), "error", err)
					return true
				}
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information about output %s", p.Output))
				return true
			}
		} else {
			if !stat.IsDir() {
				slog.Error(fmt.Sprintf("Output %s is not a directory.", p.Output))
			}
		}

		return false
	})

	for i := range spec.Pipelines {
		banksWorkspace := spec.Pipelines[i].BanksWorkspace
		nodesWorkspace := spec.Pipelines[i].ScriptsWorkspace
		spec.Pipelines[i].Banks = slices.DeleteFunc(spec.Pipelines[i].Banks, func(b string) bool {
			if !filepath.IsAbs(b) {
				b = filepath.Join(banksWorkspace, b)
			}
			stat, err := os.Lstat(b)
			if err != nil {
				if os.IsNotExist(err) {
					slog.Error(fmt.Sprintf("Bank %s does not exist.", b))
				} else {
					slog.Error(fmt.Sprintf("Fail to obtain information about bank %s", b))
				}
				return true
			}
			if stat.IsDir() {
				slog.Error(fmt.Sprintf("%s is a directory.", b))
				return true
			}
			return false
		})
		spec.Pipelines[i].Scripts = slices.DeleteFunc(spec.Pipelines[i].Scripts, func(n ProcessScript) bool {
			if !filepath.IsAbs(n.Script) {
				n.Script = filepath.Join(nodesWorkspace, n.Script)
			}
			stat, err := os.Lstat(n.Script)
			if err != nil {
				if os.IsNotExist(err) {
					slog.Error(fmt.Sprintf("Process script %s does not exist.", n.Script))
				} else {
					slog.Error(fmt.Sprintf("Fail to obtain information about process script %s", n.Script))
				}
				return true
			}
			if stat.IsDir() {
				slog.Error(fmt.Sprintf("%s is a directory.", n.Script))
				return true
			}
			if n.Type < TypeRewireWithNewSources || n.Type >= ProcessScriptTypeCount {
				slog.Error(fmt.Sprintf("Unsupported process script type %d", n.Type))
				return true
			}
			return false
		})
	}
	spec.Pipelines = slices.DeleteFunc(spec.Pipelines, func(p ProcessPipeline) bool {
		if len(p.Banks) <= 0 {
			return true
		}
		if len(p.Scripts) <= 0 {
			return true
		}
		return false
	})
	return &spec, nil
}

func Process(ctx context.Context, fspec string) {
	spec, err := ParseProcessor(fspec)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse processor %s", fspec), "error", err)
		return
	}
	sems := make(chan struct{}, 8)
	var w sync.WaitGroup
	for i := range spec.Pipelines {
		p := &spec.Pipelines[i]
		select {
		case <- ctx.Done():
			slog.Error("Process is cancelled", "error", ctx.Err().Error())
			w.Wait()
			return
		case sems <- struct{}{}:
			w.Add(1)
			go func(j int) {
				defer func() {
					<- sems
					w.Done()
				}()
				slog.Info(fmt.Sprintf("(Routine) Running process pipeline %d", j))
				RunProcessPipeline(ctx, p, sems)
				slog.Info(fmt.Sprintf("(Routine) Finishing process pipeline %d", j))
			}(i)
		default:
			slog.Info(fmt.Sprintf("Running process pipeline %d", i))
			RunProcessPipeline(ctx, p, sems)
			slog.Info(fmt.Sprintf("Finishing process pipeline %d", i))
		}
	}
	w.Wait()
}

func RunProcessPipeline(ctx context.Context, p *ProcessPipeline, sems chan struct{}) {
	var w sync.WaitGroup
	bnks := make([]*wwise.Bank, len(p.Banks), len(p.Banks))
	for i := range p.Banks {
		select {
		case <- ctx.Done():
			slog.Error("Process pipeline is cancelled", "error", ctx.Err().Error())
			w.Wait()
			return
		case sems <- struct{}{}:
			w.Add(1)
			go func(j int) {
				defer func() {
					<- sems
					w.Done()
				}()
				slog.Info(fmt.Sprintf("(Routine) Running process script on bank %s", p.Banks[j]))
				bnks[i] = RunProcessScripts(ctx, p.Banks[j], p)
				slog.Info(fmt.Sprintf("(Routine) Processed bank %s", p.Banks[j]))
			}(i)
		default:
			slog.Info(fmt.Sprintf("Running process script on bank %s", p.Banks[i]))
			bnks[i] = RunProcessScripts(ctx, p.Banks[i], p)
			slog.Info(fmt.Sprintf("Processed bank %s", p.Banks[i]))
		}
	}
	w.Wait()

	switch p.Integration {
	case integration.IntegrationTypeNone:
		panic("Panic Trap")
	case integration.IntegrationTypeDefault:
		 for i, bnk := range bnks {
			if bnk == nil {
				continue
			}
			ctx, cancel := context.WithTimeout(ctx, time.Second * 8)
			defer cancel()
			bnkData, err := bnk.Encode(ctx, false, false)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to encode sound bank %s", p.Banks[i]), "error", err)
				slog.Warn(fmt.Sprintf("Skipping sound bank %s", p.Banks[i]))
				continue
			}
			basename := filepath.Base(p.Banks[i])
			err = os.WriteFile(filepath.Join(p.Output, basename), bnkData, 0777)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to save sound bank %s to %s", basename, p.Output), "error", err)
			}
		}
	case integration.IntegrationTypeHelldivers2:
		slog.Info("Using Helldivers 2 integration")
		bnksData := [][]byte{}
		metasData := [][]byte{}
		for i, bnk := range bnks {
			if bnk == nil {
				continue
			}
			ctx, cancel := context.WithTimeout(ctx, time.Second * 8)
			defer cancel()
			bnkData, err := bnk.Encode(ctx, true, false)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to encode sound bank %s", p.Banks[i]), "error", err)
				slog.Warn(fmt.Sprintf("Skipping sound bank %s", p.Banks[i]))
				continue
			}
			meta := bnk.META()
			if meta == nil {
				slog.Error(fmt.Sprintf("Sound bank %s is missing META data for integration", p.Banks[i]))
				slog.Warn(fmt.Sprintf("Skipping sound bank %s", p.Banks[i]))
				continue
			}
			bnksData = append(bnksData, bnkData)
			metasData = append(metasData, meta.B)
		}
		if len(bnksData) <= 0  {
			slog.Warn("No sound bank data available for integration.")
			return
		}
		err := helldivers.GenHelldiversPatchStableMulti(bnksData, metasData, p.Output)
		if err != nil {
			slog.Error("Failed to run Helldivers 2 integration", "error", err)
		} else {
			slog.Info("Generated Helldivers 2 patch")
		}
	default:
		panic("Panic Trap")
	}
}

func RunProcessScripts(ctx context.Context, bank string, p *ProcessPipeline) *wwise.Bank {
	if !filepath.IsAbs(bank) {
		bank = filepath.Join(p.BanksWorkspace, bank)
	}
	bnk, err := parser.ParseBank(bank, ctx, false)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse sound bank %s", bank), "error", err)
		return nil
	}
	for _, script := range p.Scripts {
		select {
		case <- ctx.Done():
			slog.Error("Process scripts execution is cancelled", "error", ctx.Err().Error())
			return nil
		default:
		}
		if !filepath.IsAbs(script.Script) {
			script.Script = filepath.Join(p.ScriptsWorkspace, script.Script)
		}
		slog.Info(fmt.Sprintf("Running processing script %s", script.Script))
		switch script.Type {
		case TypeRewireWithNewSources:
			err = RewireWithNewSources(ctx, bnk, script.Script, false)
		case TypeBasePropModifiers:
			err = ProcessBaseProps(bnk, script.Script)
		case TypeImportAsRanSeqCntr:
			err = ImportAsRanSeqCntr(ctx, bnk, script.Script)
		case TypeReplaceAudioSources:
			err = ReplaceAudioSources(ctx, bnk, script.Script, false)
		case TypeRanSeqModifiers:
			err = ProcessRanSeq(bnk, script.Script)
		case TypeNewSoundToRanSeqCntr:
			err = NewSoundToRanSeqCntr(ctx, bnk, script.Script)
		case TypeBulkProcessBaseProp:
			err = BulkProcessBaseProp(bnk, script.Script)
		case TypeStreamTypeModifiers:
			err = ToStreamTypeBnk(bnk, script.Script)
		default:
			panic(fmt.Sprintf("Unsupport process script type %d.", script.Type))
		}
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to run process script %s", script.Script), "error", err)
		}
	}
	return bnk
}

func ProcessActiveBank(ctx context.Context, bnk *wwise.Bank, fspec string) error {
	spec, err := ParseProcessor(fspec)
	if err != nil {
		return err
	}
	if len(spec.Pipelines) == 0 {
		return nil
	}

	return ProcessPipelineActiveBank(ctx, bnk, &spec.Pipelines[0])
}

func ProcessPipelineActiveBank(ctx context.Context, bnk *wwise.Bank, p *ProcessPipeline) error {
	for _, node := range p.Scripts {
		switch node.Type {
		case TypeRewireWithNewSources:
			if err := RewireWithNewSources(ctx, bnk, node.Script, false); err != nil {
				return err
			}
		case TypeBasePropModifiers:
			if err := ProcessBaseProps(bnk, node.Script); err != nil {
				return err
			}
		default:
			slog.Warn(fmt.Sprintf("Unsupport process procedure %d", node.Type))
		}
	}
	return nil
}
