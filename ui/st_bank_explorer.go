package ui

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type BankManager struct {
	banks sync.Map
	writeLock *atomic.Bool
}

type bankTab struct {
	bank *wwise.Bank
	idQuery string
	typeQuery int32
	filtered []wwise.HircObj
	lSelStorage imgui.SelectionBasicStorage
	// Sync
	writeLock *atomic.Bool
}

func (b *bankTab) filter() {
	if b.bank.HIRC() == nil {
		return
	}
	if !utils.IsDigit(b.idQuery) {
		return
	}
	hirc := b.bank.HIRC()
	i := 0
	old := len(b.filtered)
	for _, d := range hirc.HircObjs {
		if b.typeQuery > 0 && b.typeQuery != int32(d.HircType()) {
			continue
		}
		id, err := d.HircID()
		if err != nil {
			if i < len(b.filtered) {
				b.filtered[i] = d
			} else {
				b.filtered = append(b.filtered, d)
			}
			i += 1
			continue
		}
		if !fuzzy.Match(b.idQuery, strconv.FormatUint(uint64(id), 10)) {
			continue
		}
		if i < len(b.filtered) {
			b.filtered[i] = d
		} else {
			b.filtered = append(b.filtered, d)
		}
		i += 1
	}
	if i < old {
		b.filtered = slices.Delete(b.filtered, i, old)
	}
}

func (b *bankTab) encode(ctx context.Context) ([]byte, error) {
	b.writeLock.Store(true)
	defer b.writeLock.Store(false)
	type result struct {
		data []byte
		err  error
	}
	c := make(chan *result)
	go func() {
		data, err := b.bank.Encode(ctx)
		c <- &result{data, err}
	}()

	select {
	case <- ctx.Done():
		<- c
		return nil, ctx.Err()
	case r := <- c:
		return r.data, r.err
	}
}

func (b *BankManager) openBank(ctx context.Context, path string) error {
	type result struct {
		bank *wwise.Bank
		err error
	}

	c := make(chan *result, 1)

	if _, in := b.banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}
	
	go func() {
		bank, err := parser.ParseBank(path, ctx)
		c <- &result{bank, err}
	}()

	var bank *wwise.Bank
	select {
	case res := <- c:
		if res.bank == nil {
			return res.err
		}
		bank = res.bank
	case <- ctx.Done():
		return ctx.Err()
	}

	if _, in := b.banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}


	filtered := []wwise.HircObj{}
	if bank.HIRC() != nil {
		hirc := bank.HIRC()
		filtered = make([]wwise.HircObj, len(hirc.HircObjs))
		for i, o := range hirc.HircObjs {
			filtered[i] = o
		}
	}

	t := &bankTab{
		writeLock: &atomic.Bool{},
		bank: bank,
		idQuery: "",
		typeQuery: 0,
		filtered: filtered,
		lSelStorage: *imgui.NewSelectionBasicStorage(),
	}
	t.writeLock.Store(false)

	b.banks.Store(path, t)

	return nil
}

func (b *BankManager) closeBank(del string) {
	b.banks.Delete(del)
}
