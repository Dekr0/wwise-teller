// TODO:
// - Switch to "Events", "Attenuations", "Game Sync", "FX", ...
// - File saving
// - Enable "Modifying Everything"
package ui

import (
	"slices"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type CmdPalette struct {
	name   string
	action func()
}

type RankCmdPalette struct {
	rank int
	cmd  *CmdPalette
}

type CmdPaletteMngr struct {
	query      string
	cmdPalette []*CmdPalette
	filtered   []*RankCmdPalette
	selected   int
}

func NewCmdPaletteMngr(
	dockMngr *DockManager,
	conf *config.Config,
	loop *async.EventLoop,
	modalQ *ModalQ,
) *CmdPaletteMngr {
	mngr := &CmdPaletteMngr{
		"",
		make([]*CmdPalette, 0, 16),
		[]*RankCmdPalette{},
		0,
	}
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"config",
		func() { pushConfigModalFunc(modalQ, conf) },
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"focus file explorer",
		func() { imgui.SetWindowFocusStr("File Explorer") },
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"focus bank explorer",
		func() { imgui.SetWindowFocusStr("Bank Explorer") },
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"focus log",
		func() { imgui.SetWindowFocusStr("Log") },
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"focus object editor",
		func() { imgui.SetWindowFocusStr("Object Editor") },
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"Use docking layout 1",
		func() {
			dockMngr.layout = Layout01
			dockMngr.rebuild = true
		},
	})
	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"Use docking layout 2",
		func() {
			dockMngr.layout = Layout02
			dockMngr.rebuild = true
		},
	})

	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"integration: extract sound banks from Helldivers 2 game archives",
		func() { pushSelectGameArchiveModal(modalQ, loop, conf) },
	})

	mngr.cmdPalette = append(mngr.cmdPalette, &CmdPalette{
		"Disable All Guard Rails",
		func() { ModifiyEverything = !ModifiyEverything },
	})

	mngr.filtered = make([]*RankCmdPalette, len(mngr.cmdPalette))
	for i, c := range mngr.cmdPalette {
		mngr.filtered[i] = &RankCmdPalette{-1, c}
	}
	
	return mngr
}

func (c *CmdPaletteMngr) filter() {
	i := 0
	old := len(c.filtered)
	for _, cmd := range c.cmdPalette {
		rank := fuzzy.RankMatch(c.query, cmd.name)
		if rank == -1 {
			continue
		}
		if i < len(c.filtered) {
			c.filtered[i].rank = rank
			c.filtered[i].cmd = cmd
		} else {
			c.filtered = append(c.filtered, &RankCmdPalette{rank, cmd})
		}
		i += 1
	}
	if i < old {
		c.filtered = slices.Delete(c.filtered, i, old)
	}
	slices.SortFunc(c.filtered, func(a *RankCmdPalette, b *RankCmdPalette) int {
		if a.rank < b.rank {
			return -1
		}
		if a.rank == b.rank {
			return 0
		}
		return 1
	})
	c.selected = 0
}

func (c *CmdPaletteMngr) SetNext(delta int8) {
	if c.selected + int(delta) >= 0 && c.selected + int(delta) < len(c.filtered) {
		c.selected += int(delta)
	}
}
