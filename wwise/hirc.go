package wwise

import (
	"fmt"
	"runtime"
)

type PropBundle struct {
	Ps []uint8
	Vs [][]byte
}

type PropBundleHandle struct {
	Id         uint32
	PropBundle PropBundle
}

type PropBundleSystem struct {
	ActionPropBundle    []PropBundleHandle
	AuxBusPropBundle    []PropBundleHandle
	BaseParamPropBundle []PropBundleHandle
	BusPropBundle       []PropBundleHandle
	ModulatorPropBundle []PropBundleHandle
}

func AllocatePropBundleSystem(s *PropBundleSystem, busBank bool) {
	s.ActionPropBundle = make([]PropBundleHandle, 0, 64)
	if busBank {
		s.AuxBusPropBundle = make([]PropBundleHandle, 0, 64)
		s.BusPropBundle = make([]PropBundleHandle, 0, 64)
		s.BaseParamPropBundle = make([]PropBundleHandle, 0, 4)
	} else {
		s.AuxBusPropBundle = make([]PropBundleHandle, 0, 4)
		s.BusPropBundle = make([]PropBundleHandle, 0, 4)
		s.BaseParamPropBundle = make([]PropBundleHandle, 0, 128)
	}
	s.ModulatorPropBundle = make([]PropBundleHandle, 0, 4)
}
