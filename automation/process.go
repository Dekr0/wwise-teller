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

	"github.com/Dekr0/wwise-teller/integration"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

const ProcessSpecVersion = 0

type ProcessScriptType uint8

const (
	RewireWithNewSourcesType ProcessScriptType = 0
	RewireWithOldSourcesType ProcessScriptType = 1
	BasePropModifiersType    ProcessScriptType = 2
	ImportAsRanSeqCntrType   ProcessScriptType = 3
)

type Processor struct {
	Veriosn     uint8           `json:"version"`
	Pipelines []ProcessPipeline `json:"pipelines"`
}

type ProcessPipeline struct {
	BanksWorkspace   string                      `json:"banksWorkspace"`
	ScriptsWorkspace string                      `json:"scriptsWorkspace"`
	Banks          []string                      `json:"banks"`
	Scripts        []ProcessScript               `json:"scripts"`
	Integration      integration.IntegrationType `json:"integration"`
	Output           string                      `json:"output"`
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
			switch n.Type {
			case RewireWithNewSourcesType, RewireWithOldSourcesType, BasePropModifiersType, ImportAsRanSeqCntrType:
			default:
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
	for _, p := range spec.Pipelines {
		select {
		case <- ctx.Done():
			slog.Error(ctx.Err().Error())
			w.Wait()
			return
		case sems <- struct{}{}:
			w.Add(1)
			go RunProcessPipeline(ctx, &p, &w, sems)
		default:
			RunProcessPipeline(ctx, &p, nil, nil)
		}
	}
	w.Wait()
}

func RunProcessPipeline(
	ctx context.Context,
	p *ProcessPipeline,
	w *sync.WaitGroup,
	sems chan struct{},
) {
	if w != nil { defer w.Done() }
	var _w sync.WaitGroup
	if sems != nil { defer func(){ <- sems }() }
	for _, b := range p.Banks {
		select {
		case <- ctx.Done():
			slog.Error(ctx.Err().Error())
			_w.Wait()
			return
		case sems <- struct{}{}:
			_w.Add(1)
			go RunProcessScripts(ctx, b, p, &_w, sems)
		default:
			RunProcessScripts(ctx, b, p, nil, nil)
		}
	}
	_w.Wait()
}

func RunProcessScripts(
	ctx context.Context,
	bank string,
	p *ProcessPipeline,
	w *sync.WaitGroup,
	sems chan struct{},
) {
	if w != nil { defer w.Done() }
	if sems != nil { defer func() { <- sems }() }
	if !filepath.IsAbs(bank) {
		bank = filepath.Join(p.BanksWorkspace, bank)
	}
	bnk, err := parser.ParseBank(bank, ctx, false)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to parse sound bank %s", bank), "error", err)
		return
	}
	basename := filepath.Base(bank)
	for _, script := range p.Scripts {
		select {
		case <- ctx.Done():
			slog.Error(ctx.Err().Error())
			return
		default:
		}
		if !filepath.IsAbs(script.Script) {
			script.Script = filepath.Join(p.ScriptsWorkspace, script.Script)
		}
		switch script.Type {
		case RewireWithNewSourcesType:
			err = RewireWithNewSources(ctx, bnk, script.Script, false)
		case BasePropModifiersType:
			err = ProcessBaseProps(bnk, script.Script)
		case ImportAsRanSeqCntrType:
			err = ImportAsRanSeqCntr(ctx, bnk, script.Script)
		default:
			panic("Panic Trap")
		}
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to run process script %s", script.Script), "error", err)
		}
	}
	switch p.Integration {
	case integration.IntegrationTypeNone:
		panic("Panic Trap")
	case integration.IntegrationTypeDefault:
		bnkData, err := bnk.Encode(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to encode sound bank %s", basename), "error", err)
			return
		}
		err = os.WriteFile(filepath.Join(p.Output, basename), bnkData, 0777)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to save sound bank %s to %s", basename, p.Output), "error", err)
			return
		}
	case integration.IntegrationTypeHelldivers2:
		bnkData, err := bnk.Encode(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to encode sound bank %s", basename), "error", err)
			return
		}
		meta := bnk.META()
		if meta == nil {
			slog.Error(fmt.Sprintf("Sound bank %s is missing META chunk.", basename))
			return
		}
		if err := helldivers.GenHelldiversPatchStable(bnkData, meta.B, p.Output); err != nil {
			slog.Error(fmt.Sprintf("Failed to pack sound bank %s into Helldivers 2 game init archive", basename))
			return
		}
	default:
		panic("Panic Trap")
	}
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
		case RewireWithNewSourcesType:
			if err := RewireWithNewSources(ctx, bnk, node.Script, false); err != nil {
				return err
			}
		case BasePropModifiersType:
			if err := ProcessBaseProps(bnk, node.Script); err != nil {
				return err
			}
		default:
			slog.Warn(fmt.Sprintf("Unsupport process procedure %d", node.Type))
		}
	}
	return nil
}
