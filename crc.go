// SPDX-License-Identifier: MIT-0
// SPDX-FileCopyrightText:  2024 Istvan Pasztor

// Package crc is an arbitrary-precision CRC calculator that can calculate CRCs
// of any bit width (between CRC-1 and CRC-64) and can process input of any bit
// length (the end of the input data doesn't have to be on a byte boundary).
//
// Whole bytes of the input data are processed with the help of a precalculated
// 256-element accelerator table. If the end of input isn't byte-aligned then
// the remaining (7 or fewer) bits are calculated into the CRC by a tableless
// bit-by-bit method.
//
// This package provides presets for and has been tested against
// the 100+ CRC algorithms listed in Greg Cook's CRC catalogue:
// https://reveng.sourceforge.io/crc-catalogue/all.htm
package crc

import "errors"

// UInt specifies the integer types that can be used for CRC calculations.
// The bit width of the chosen integer type has to be greater than or equal to
// the bit width used by the CRC algorithm.
// For example a CRC-17 algorithm requires uint32 or uint64.
type UInt interface {
	uint8 | uint16 | uint32 | uint64
}

// A CRC instance is a lightweight "throw-away" object that can calculate the
// CRC of your chunked data with zero or more Update() calls.
type CRC[T UInt] interface {
	Update(data []byte)
	UpdateBits(data []byte, bitLen int)
	Final() T   // Final returns the final CRC value
	Residue() T // Residue returns the final CRC value without the xorout step
}

// Algo is a parametrized CRC algorithm. It can be shared and reused by goroutines
// to save on the resources spent on creating the related accelerator table.
type Algo[T UInt] interface {
	NewCRC() CRC[T]                     // Calculate the CRC of chunked data
	Calc(data []byte) T                 // Calculate the CRC of a single chunk of data
	CalcBits(data []byte, bitLen int) T // Calculate the CRC of a single chunk of data
}

// NewAlgo creates a parametrized CRC algorithm instance - this involves the
// calculation of an accelerator table with 256 entries of type T. Ideally you
// create and share one Algo instance per CRC algorithm during the lifespan of
// the process. Width can be between 1...64 (inclusive) - it mustn't exceed the
// bit width of T. Poly and init are always in (unreflected) MSB-first format.
func NewAlgo[T UInt](width int, poly, init, xorout T, refin, refout bool) (Algo[T], error) {
	if err := checkParams(width, poly, init, xorout); err != nil {
		return nil, err
	}
	a := &algo[T]{width, reflect(poly, width), reflect(init, width), xorout, refin, refout, [256]T{}}
	for i := 1; i < 256; i++ {
		a.table[i] = a.bbbUpd(T(i), 0, 8)
	}
	return a, nil
}

func checkParams[T UInt](width int, poly, init, xorout T) error {
	if width <= 0 || (T(1)<<(width-1)) == 0 {
		return errors.New("width must be greater than zero and less than or equal to the bit width of T")
	}
	m := (T(1) << width) - 1
	if poly > m || init > m || xorout > m {
		return errors.New("poly, init or xorout is outside of the range allowed by width")
	}
	return nil
}

type algo[T UInt] struct {
	width   int // width>0 && width<=bitWidth(T)
	refPoly T   // reflected poly
	refInit T   // reflected init
	xorout  T
	refin   bool
	refout  bool
	table   [256]T
}

func (a *algo[T]) NewCRC() CRC[T] {
	return &crc[T]{a, a.refInit}
}

func (a *algo[T]) Calc(data []byte) T {
	return a.CalcBits(data, -1)
}

func (a *algo[T]) CalcBits(data []byte, bitLen int) T {
	c := a.NewCRC()
	c.UpdateBits(data, bitLen)
	return c.Final()
}

func (a *algo[T]) tblUpd(reg T, data []byte, bitLen int) (newReg T) {
	var n, bitsLeft int
	if bitLen < 0 {
		n, bitsLeft = len(data), 0
	} else if bitLen > (len(data) << 3) {
		panic("bitLen is greater than the number of bits in the input data")
	} else {
		n, bitsLeft = bitLen>>3, bitLen&7
	}

	for _, b := range data[:n] {
		if !a.refin {
			b = reflectedBytes[b]
		}
		reg = a.table[byte(reg)^b] ^ (reg >> 8)
	}

	if bitsLeft > 0 { // 7 or less input data bits remaining
		return a.bbbUpd(reg, data[n], bitsLeft)
	}
	return reg
}

// bbbUpd performs a bit-by-bit (tableless) update.
func (a *algo[T]) bbbUpd(reg T, b byte, bitLen int) (newReg T) {
	if !a.refin {
		b = reflectedBytes[b]
	}
	b &= (1 << bitLen) - 1 // zeroing the unused bits
	reg ^= T(b)

	for i := 0; i < bitLen; i++ {
		if (reg & 1) != 0 {
			reg = (reg >> 1) ^ a.refPoly
		} else {
			reg >>= 1
		}
	}
	return reg
}

type crc[T UInt] struct {
	a   *algo[T]
	reg T // reflected (LSB-first) CRC shift register
}

func (c *crc[T]) Update(data []byte) {
	c.reg = c.a.tblUpd(c.reg, data, -1)
}

func (c *crc[T]) UpdateBits(data []byte, bitLen int) {
	c.reg = c.a.tblUpd(c.reg, data, bitLen)
}

func (c *crc[T]) Final() T {
	return c.Residue() ^ c.a.xorout
}

func (c *crc[T]) Residue() T {
	if c.a.refout {
		return c.reg
	}
	return reflect(c.reg, c.a.width)
}

func reflect[T UInt](val T, numBits int) T {
	x := val & 1
	for i := 1; i < numBits; i++ {
		val >>= 1
		x <<= 1
		x |= val & 1
	}
	return x
}

var reflectedBytes [256]byte

func init() {
	for i := byte(1); i != 0; i++ {
		reflectedBytes[i] = reflect(i, 8)
	}
}
