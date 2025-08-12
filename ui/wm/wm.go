package wm

type Layout uint8

type DockTag uint8

type WM struct {
	Opens       []bool
	ControlBits   uint8
	Layout        uint8
}

func AllocateWM(w *WM) {
	w.Opens = make([]bool, DockTagCount)
	w.ControlBits = 0b0000_0011
}
