package io

import "fmt"

// Panic when pos is > 8
func Bit(v u8, pos u8) bool {
	if pos >= 8 {
		panic(fmt.Sprintf("%d exceeds length of 8 bits vector", pos))
	}
	return (v >> pos) & 1 > 0
}

// Panic when pos is > 8
func SetBit(v u8, pos u8, set bool) u8 {
	if pos >= 8 {
		panic(fmt.Sprintf("%d exceeds length of 8 bits vector", pos))
	}
	if !set {
		return v & (^(1 << pos))
	}
	return v | (1 << pos)
}
