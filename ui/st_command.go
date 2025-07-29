// TODO:
// - File saving
// - Enable "Modifying Everything"
package ui

import (
	"fmt"
	"slices"

	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type CmdPalette struct {
	Name   string
	Action func()
}

type RankCmdPalette struct {
	Rank int
	Cmd  *CmdPalette
}

type CmdPaletteMngr struct {
	Query      string
	CmdPalette []CmdPalette
	Filtered   []RankCmdPalette
	Selected   int
}

func NewCmdPaletteMngrP(c *CmdPaletteMngr, dockMngr *dockmanager.DockManager) {
	c.Query = ""
	c.CmdPalette = make([]CmdPalette, 0, 16)
	c.Filtered = make([]RankCmdPalette, 0, 16)
	c.Selected = 0
	c.CmdPalette = append(c.CmdPalette, CmdPalette{
		"config",
		func() { pushConfigModalFunc() },
	})
	for _, name := range dockmanager.DockWindowNames {
		cmdName := name
		c.CmdPalette = append(c.CmdPalette, CmdPalette{
			fmt.Sprintf("focus %s", name),
			func() { dockMngr.SetFocus(cmdName) },
		})
	}
	for i := range dockmanager.LayoutCount {
		li := i
		c.CmdPalette = append(c.CmdPalette, CmdPalette{
			dockmanager.LayoutName[i],
			func() { dockMngr.SetLayout(li) },
		})
	}
	c.CmdPalette = append(c.CmdPalette, CmdPalette{
		"integration: extract sound banks from Helldivers 2 game archives",
		func() { pushSelectGameArchiveModal() },
	})

	c.CmdPalette = append(c.CmdPalette, CmdPalette{
		"Disable All Guard Rails",
		func() { ModifiyEverything = !ModifiyEverything },
	})
	c.Filtered = make([]RankCmdPalette, len(c.CmdPalette))
	for i := range c.CmdPalette {
		c.Filtered[i] = RankCmdPalette{-1, &c.CmdPalette[i]}
	}
}

func (c *CmdPaletteMngr) Filter() {
	c.Filtered = slices.Delete(c.Filtered, 0, len(c.Filtered))
	for i := range c.CmdPalette {
		cmd := &c.CmdPalette[i]
		rank := fuzzy.RankMatch(c.Query, c.CmdPalette[i].Name)
		if rank == -1 {
			continue
		}
		c.Filtered = append(c.Filtered, RankCmdPalette{rank, cmd})
	}
	slices.SortFunc(c.Filtered, func(a RankCmdPalette, b RankCmdPalette) int {
		if a.Rank < b.Rank {
			return -1
		}
		if a.Rank == b.Rank {
			return 0
		}
		return 1
	})
	c.Selected = 0
}

func (c *CmdPaletteMngr) SetNext(delta int8) {
	if c.Selected + int(delta) >= 0 && c.Selected + int(delta) < len(c.Filtered) {
		c.Selected += int(delta)
	}
}
