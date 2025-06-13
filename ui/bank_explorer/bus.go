package bank_explorer 

import "github.com/Dekr0/wwise-teller/wwise"

type BusFilter struct {
}

type BusViewer struct {
	BusFilter  BusFilter
	ActiveBus  wwise.HircObj
}
