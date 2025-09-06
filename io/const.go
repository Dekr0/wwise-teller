package io

import (
	"errors"
)

const Size8  = 1
const Size16 = 2
const Size32 = 4
const Size64 = 8

const MaxStzSize = 255

var ExceedStzSize = errors.New("Zero-terminated string reach maximum length of 255")

var ExceedV128Size = errors.New("Exceed size of little endian 128 bits")
