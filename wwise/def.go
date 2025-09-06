package wwise

type LerpType = u8

const (
	LerpLog3      LerpType = 0
	LerpSine      LerpType = 1
	LerpLog1      LerpType = 2
	LerpInvSCurve LerpType = 3
	LerpLinear    LerpType = 4
	LerpSCurve    LerpType = 5
	LerpExp1      LerpType = 6
	LerpnvSine    LerpType = 7
	LerpExp3      LerpType = 8
	LerpConst     LerpType = 9
	LerpCount     LerpType = 10
)

type HircType = uint8

const (
	HircTypeState  HircType = 0x01
	HircTypeSound  HircType = 0x02
	HircTypeAction HircType = 0x03
	HircTypeEvent  HircType = 0x04
)

type ChunkName = string
const (
	ChunkNameBKHD ChunkName = "BKHD"
	ChunkNameDIDX ChunkName = "DIDX"
	ChunkNameDATA ChunkName = "DATA"
	ChunkNameHIRC ChunkName = "HIRC"
)
