// SPDX-License-Identifier: MIT-0
// SPDX-FileCopyrightText:  2024 Istvan Pasztor

package crc

import "sync"

// Preset is a proxy to a wrapped Algo instance that is created only when a
// method of the Preset is called. Unused presets never create Algo instances
// so their tables aren't allocated and calculated.
//
// The methods of a Preset instance use synchronization to create the Algo
// instance in a goroutine-safe way. The synchronization overhead is
// insignificant if the program doesn't call the Preset methods very often.
// If it becomes significant it can be avoided by using the underlying Algo
// instance directly after querying it with the Preset.Algo method.
// The methods of the CRC instances returned by Preset.NewCRC aren't affected
// by the previously mentioned synchronization.
type Preset[T UInt] interface {
	Algo[T]
	Algo() Algo[T]
}

// NewPreset creates a preset or returns an error in case of invalid parameters.
// Poly and init are always in (unreflected) MSB-first format.
func newPreset[T UInt](width int, poly, init, xorout T, refin, refout bool) (Preset[T], error) {
	if err := checkParams(width, poly, init, xorout); err != nil {
		return nil, err
	}
	return &preset[T]{width: width, poly: poly, init: init, xorout: xorout,
		refin: refin, refout: refout}, nil
}

// MustNewPreset creates a preset or panics in case of invalid parameters.
// Poly and init are always in (unreflected) MSB-first format.
func mustNewPreset[T UInt](width int, poly, init, xorout T, refin, refout bool) Preset[T] {
	p, err := newPreset(width, poly, init, xorout, refin, refout)
	if err != nil {
		panic(err)
	}
	return p
}

type preset[T UInt] struct {
	width    int
	poly     T
	init     T
	xorout   T
	refin    bool
	refout   bool
	algo     Algo[T]
	algoOnce sync.Once
}

func (p *preset[T]) NewCRC() CRC[T] {
	return p.Algo().NewCRC()
}

func (p *preset[T]) Calc(data []byte) T {
	return p.Algo().Calc(data)
}

func (p *preset[T]) CalcBits(data []byte, bitLen int) T {
	return p.Algo().CalcBits(data, bitLen)
}

func (p *preset[T]) Algo() Algo[T] {
	p.algoOnce.Do(func() {
		a, err := NewAlgo(p.width, p.poly, p.init, p.xorout, p.refin, p.refout)
		if err != nil {
			panic("invalid CRC preset")
		}
		p.algo = a
	})
	return p.algo
}

// These presets provide quick access to well documented CRC algorithms.
// Use their NewCRC and Calc methods to perform CRC calculations.
//
// Source: https://reveng.sourceforge.io/crc-catalogue/all.htm
var (
	CRC8  = CRC8SMBUS
	CRC16 = CRC16ARC
	CRC32 = CRC32ISOHDLC
	CRC64 = CRC64ECMA182

	CRC32C = CRC32ISCSI
	CRC32D = CRC32BASE91D
	CRC32Q = CRC32AIXM

	A = CRC16ISOIEC144433A
	B = CRC16IBMSDLC

	X25             = CRC16IBMSDLC
	CRC16X25        = CRC16IBMSDLC
	XMODEM          = CRC16XMODEM
	KERMIT          = CRC16KERMIT
	CRC16CCITT      = CRC16KERMIT
	CRC16CCITTFALSE = CRC16IBM3740 // commonly misidentified as CRC-16/CCITT
	CRC16AUGCCITT   = CRC16SPIFUJITSU
	V41LSB          = CRC16KERMIT
	V41MSB          = CRC16XMODEM

	PKZIP      = CRC32ISOHDLC
	V42        = CRC32ISOHDLC
	XZ         = CRC32ISOHDLC
	POSIX      = CRC32CKSUM
	CASTAGNOLI = CRC32ISCSI

	CRC3GSM  = mustNewPreset[uint8](3, 0x3, 0x0, 0x7, false, false) // CRC-3/GSM
	CRC3ROHC = mustNewPreset[uint8](3, 0x3, 0x7, 0x0, true, true)   // CRC-3/ROHC

	CRC4INTERLAKEN = mustNewPreset[uint8](4, 0x3, 0xf, 0xf, false, false) // CRC-4/INTERLAKEN
	CRC4G704       = mustNewPreset[uint8](4, 0x3, 0x0, 0x0, true, true)   // CRC-4/G-704      Alias: CRC-4/ITU

	CRC5USB     = mustNewPreset[uint8](5, 0x05, 0x1f, 0x1f, true, true)   // CRC-5/USB
	CRC5EPCC1G2 = mustNewPreset[uint8](5, 0x09, 0x09, 0x00, false, false) // CRC-5/EPC-C1G2   Alias: CRC-5/EPC
	CRC5G704    = mustNewPreset[uint8](5, 0x15, 0x00, 0x00, true, true)   // CRC-5/G-704      Alias: CRC-5/ITU

	CRC6G704      = mustNewPreset[uint8](6, 0x03, 0x00, 0x00, true, true)   // CRC-6/G-704      Alias: CRC-6/ITU
	CRC6CDMA2000B = mustNewPreset[uint8](6, 0x07, 0x3f, 0x00, false, false) // CRC-6/CDMA2000-B
	CRC6DARC      = mustNewPreset[uint8](6, 0x19, 0x00, 0x00, true, true)   // CRC-6/DARC
	CRC6CDMA2000A = mustNewPreset[uint8](6, 0x27, 0x3f, 0x00, false, false) // CRC-6/CDMA2000-A
	CRC6GSM       = mustNewPreset[uint8](6, 0x2f, 0x00, 0x3f, false, false) // CRC-6/GSM

	CRC7MMC  = mustNewPreset[uint8](7, 0x09, 0x00, 0x00, false, false) // CRC-7/MMC        Alias: CRC-7
	CRC7UMTS = mustNewPreset[uint8](7, 0x45, 0x00, 0x00, false, false) // CRC-7/UMTS
	CRC7ROHC = mustNewPreset[uint8](7, 0x4f, 0x7f, 0x00, true, true)   // CRC-7/ROHC

	CRC8SMBUS      = mustNewPreset[uint8](8, 0x07, 0x00, 0x00, false, false) // CRC-8/SMBUS      Alias: CRC-8
	CRC8I4321      = mustNewPreset[uint8](8, 0x07, 0x00, 0x55, false, false) // CRC-8/I-432-1
	CRC8ROHC       = mustNewPreset[uint8](8, 0x07, 0xff, 0x00, true, true)   // CRC-8/ROHC
	CRC8GSMA       = mustNewPreset[uint8](8, 0x1d, 0x00, 0x00, false, false) // CRC-8/GSM-A
	CRC8MIFAREMAD  = mustNewPreset[uint8](8, 0x1d, 0xc7, 0x00, false, false) // CRC-8/MIFARE-MAD
	CRC8ICODE      = mustNewPreset[uint8](8, 0x1d, 0xfd, 0x00, false, false) // CRC-8/I-CODE
	CRC8HITAG      = mustNewPreset[uint8](8, 0x1d, 0xff, 0x00, false, false) // CRC-8/HITAG
	CRC8SAEJ1850   = mustNewPreset[uint8](8, 0x1d, 0xff, 0xff, false, false) // CRC-8/SAE-J1850
	CRC8TECH3250   = mustNewPreset[uint8](8, 0x1d, 0xff, 0x00, true, true)   // CRC-8/TECH-3250  Alias: CRC-8/AES, CRC-8/EBU
	CRC8OPENSAFETY = mustNewPreset[uint8](8, 0x2f, 0x00, 0x00, false, false) // CRC-8/OPENSAFETY
	CRC8AUTOSAR    = mustNewPreset[uint8](8, 0x2f, 0xff, 0xff, false, false) // CRC-8/AUTOSAR
	CRC8NRSC5      = mustNewPreset[uint8](8, 0x31, 0xff, 0x00, false, false) // CRC-8/NRSC-5
	CRC8MAXIMDOW   = mustNewPreset[uint8](8, 0x31, 0x00, 0x00, true, true)   // CRC-8/MAXIM-DOW  Alias: CRC-8/MAXIM, DOW-CRC
	CRC8DARC       = mustNewPreset[uint8](8, 0x39, 0x00, 0x00, true, true)   // CRC-8/DARC
	CRC8GSMB       = mustNewPreset[uint8](8, 0x49, 0x00, 0xff, false, false) // CRC-8/GSM-B
	CRC8LTE        = mustNewPreset[uint8](8, 0x9b, 0x00, 0x00, false, false) // CRC-8/LTE
	CRC8CDMA2000   = mustNewPreset[uint8](8, 0x9b, 0xff, 0x00, false, false) // CRC-8/CDMA2000
	CRC8WCDMA      = mustNewPreset[uint8](8, 0x9b, 0x00, 0x00, true, true)   // CRC-8/WCDMA
	CRC8BLUETOOTH  = mustNewPreset[uint8](8, 0xa7, 0x00, 0x00, true, true)   // CRC-8/BLUETOOTH
	CRC8DVBS2      = mustNewPreset[uint8](8, 0xd5, 0x00, 0x00, false, false) // CRC-8/DVB-S2

	CRC10GSM      = mustNewPreset[uint16](10, 0x175, 0x000, 0x3ff, false, false) // CRC-10/GSM
	CRC10ATM      = mustNewPreset[uint16](10, 0x233, 0x000, 0x000, false, false) // CRC-10/ATM       Alias: CRC-10, CRC-10/I-610
	CRC10CDMA2000 = mustNewPreset[uint16](10, 0x3d9, 0x3ff, 0x000, false, false) // CRC-10/CDMA2000

	CRC11UMTS    = mustNewPreset[uint16](11, 0x307, 0x000, 0x000, false, false) // CRC-11/UMTS
	CRC11FLEXRAY = mustNewPreset[uint16](11, 0x385, 0x01a, 0x000, false, false) // CRC-11/FLEXRAY

	CRC12DECT     = mustNewPreset[uint16](12, 0x80f, 0x000, 0x000, false, false) // CRC-12/DECT      Alias: X-CRC-12
	CRC12UMTS     = mustNewPreset[uint16](12, 0x80f, 0x000, 0x000, false, true)  // CRC-12/UMTS      Alias: CRC-12/3GPP
	CRC12GSM      = mustNewPreset[uint16](12, 0xd31, 0x000, 0xfff, false, false) // CRC-12/GSM
	CRC12CDMA2000 = mustNewPreset[uint16](12, 0xf13, 0xfff, 0x000, false, false) // CRC-12/CDMA2000

	CRC13BBC = mustNewPreset[uint16](13, 0x1cf5, 0x0000, 0x0000, false, false) // CRC-13/BBC

	CRC14DARC = mustNewPreset[uint16](14, 0x0805, 0x0000, 0x0000, true, true)   // CRC-14/DARC
	CRC14GSM  = mustNewPreset[uint16](14, 0x202d, 0x0000, 0x3fff, false, false) // CRC-14/GSM

	CRC15CAN     = mustNewPreset[uint16](15, 0x4599, 0x0000, 0x0000, false, false) // CRC-15/CAN       Alias: CRC-15
	CRC15MPT1327 = mustNewPreset[uint16](15, 0x6815, 0x0000, 0x0001, false, false) // CRC-15/MPT1327

	CRC16DECTX         = mustNewPreset[uint16](16, 0x0589, 0x0000, 0x0000, false, false) // CRC-16/DECT-X    Alias: X-CRC-16
	CRC16DECTR         = mustNewPreset[uint16](16, 0x0589, 0x0000, 0x0001, false, false) // CRC-16/DECT-R    Alias: R-CRC-16
	CRC16NRSC5         = mustNewPreset[uint16](16, 0x080b, 0xffff, 0x0000, true, true)   // CRC-16/NRSC-5
	CRC16XMODEM        = mustNewPreset[uint16](16, 0x1021, 0x0000, 0x0000, false, false) // CRC-16/XMODEM    Alias: CRC-16/ACORN, CRC-16/LTE, CRC-16/V-41-MSB, XMODEM, ZMODEM
	CRC16GSM           = mustNewPreset[uint16](16, 0x1021, 0x0000, 0xffff, false, false) // CRC-16/GSM
	CRC16SPIFUJITSU    = mustNewPreset[uint16](16, 0x1021, 0x1d0f, 0x0000, false, false) // CRC-16/SPI-FUJITSU Alias: CRC-16/AUG-CCITT
	CRC16IBM3740       = mustNewPreset[uint16](16, 0x1021, 0xffff, 0x0000, false, false) // CRC-16/IBM-3740  Alias: CRC-16/AUTOSAR, CRC-16/CCITT-FALSE
	CRC16GENIBUS       = mustNewPreset[uint16](16, 0x1021, 0xffff, 0xffff, false, false) // CRC-16/GENIBUS   Alias: CRC-16/DARC, CRC-16/EPC, CRC-16/EPC-C1G2, CRC-16/I-CODE
	CRC16KERMIT        = mustNewPreset[uint16](16, 0x1021, 0x0000, 0x0000, true, true)   // CRC-16/KERMIT    Alias: CRC-16/BLUETOOTH, CRC-16/CCITT, CRC-16/CCITT-TRUE, CRC-16/V-41-LSB, CRC-CCITT, KERMIT
	CRC16TMS37157      = mustNewPreset[uint16](16, 0x1021, 0x89ec, 0x0000, true, true)   // CRC-16/TMS37157
	CRC16RIELLO        = mustNewPreset[uint16](16, 0x1021, 0xb2aa, 0x0000, true, true)   // CRC-16/RIELLO
	CRC16ISOIEC144433A = mustNewPreset[uint16](16, 0x1021, 0xc6c6, 0x0000, true, true)   // CRC-16/ISO-IEC-14443-3-A Alias: CRC-A
	CRC16MCRF4XX       = mustNewPreset[uint16](16, 0x1021, 0xffff, 0x0000, true, true)   // CRC-16/MCRF4XX
	CRC16IBMSDLC       = mustNewPreset[uint16](16, 0x1021, 0xffff, 0xffff, true, true)   // CRC-16/IBM-SDLC  Alias: CRC-16/ISO-HDLC, CRC-16/ISO-IEC-14443-3-B, CRC-16/X-25, CRC-B, X-25
	CRC16PROFIBUS      = mustNewPreset[uint16](16, 0x1dcf, 0xffff, 0xffff, false, false) // CRC-16/PROFIBUS  Alias: CRC-16/IEC-61158-2
	CRC16EN13757       = mustNewPreset[uint16](16, 0x3d65, 0x0000, 0xffff, false, false) // CRC-16/EN-13757
	CRC16DNP           = mustNewPreset[uint16](16, 0x3d65, 0x0000, 0xffff, true, true)   // CRC-16/DNP
	CRC16OPENSAFETYA   = mustNewPreset[uint16](16, 0x5935, 0x0000, 0x0000, false, false) // CRC-16/OPENSAFETY-A
	CRC16M17           = mustNewPreset[uint16](16, 0x5935, 0xffff, 0x0000, false, false) // CRC-16/M17
	CRC16LJ1200        = mustNewPreset[uint16](16, 0x6f63, 0x0000, 0x0000, false, false) // CRC-16/LJ1200
	CRC16OPENSAFETYB   = mustNewPreset[uint16](16, 0x755b, 0x0000, 0x0000, false, false) // CRC-16/OPENSAFETY-B
	CRC16UMTS          = mustNewPreset[uint16](16, 0x8005, 0x0000, 0x0000, false, false) // CRC-16/UMTS      Alias: CRC-16/BUYPASS, CRC-16/VERIFONE
	CRC16DDS110        = mustNewPreset[uint16](16, 0x8005, 0x800d, 0x0000, false, false) // CRC-16/DDS-110
	CRC16CMS           = mustNewPreset[uint16](16, 0x8005, 0xffff, 0x0000, false, false) // CRC-16/CMS
	CRC16ARC           = mustNewPreset[uint16](16, 0x8005, 0x0000, 0x0000, true, true)   // CRC-16/ARC       Alias: ARC, CRC-16, CRC-16/LHA, CRC-IBM
	CRC16MAXIMDOW      = mustNewPreset[uint16](16, 0x8005, 0x0000, 0xffff, true, true)   // CRC-16/MAXIM-DOW Alias: CRC-16/MAXIM
	CRC16MODBUS        = mustNewPreset[uint16](16, 0x8005, 0xffff, 0x0000, true, true)   // CRC-16/MODBUS    Alias: MODBUS
	CRC16USB           = mustNewPreset[uint16](16, 0x8005, 0xffff, 0xffff, true, true)   // CRC-16/USB
	CRC16T10DIF        = mustNewPreset[uint16](16, 0x8bb7, 0x0000, 0x0000, false, false) // CRC-16/T10-DIF
	CRC16TELEDISK      = mustNewPreset[uint16](16, 0xa097, 0x0000, 0x0000, false, false) // CRC-16/TELEDISK
	CRC16CDMA2000      = mustNewPreset[uint16](16, 0xc867, 0xffff, 0x0000, false, false) // CRC-16/CDMA2000

	CRC17CANFD = mustNewPreset[uint32](17, 0x1685b, 0x00000, 0x00000, false, false) // CRC-17/CAN-FD

	CRC21CANFD = mustNewPreset[uint32](21, 0x102899, 0x000000, 0x000000, false, false) // CRC-21/CAN-FD

	CRC24BLE        = mustNewPreset[uint32](24, 0x00065b, 0x555555, 0x000000, true, true)   // CRC-24/BLE
	CRC24INTERLAKEN = mustNewPreset[uint32](24, 0x328b63, 0xffffff, 0xffffff, false, false) // CRC-24/INTERLAKEN
	CRC24FLEXRAYB   = mustNewPreset[uint32](24, 0x5d6dcb, 0xabcdef, 0x000000, false, false) // CRC-24/FLEXRAY-B
	CRC24FLEXRAYA   = mustNewPreset[uint32](24, 0x5d6dcb, 0xfedcba, 0x000000, false, false) // CRC-24/FLEXRAY-A
	CRC24LTEB       = mustNewPreset[uint32](24, 0x800063, 0x000000, 0x000000, false, false) // CRC-24/LTE-B
	CRC24OS9        = mustNewPreset[uint32](24, 0x800063, 0xffffff, 0xffffff, false, false) // CRC-24/OS-9
	CRC24LTEA       = mustNewPreset[uint32](24, 0x864cfb, 0x000000, 0x000000, false, false) // CRC-24/LTE-A
	CRC24OPENPGP    = mustNewPreset[uint32](24, 0x864cfb, 0xb704ce, 0x000000, false, false) // CRC-24/OPENPGP   Alias: CRC-24

	CRC30CDMA = mustNewPreset[uint32](30, 0x2030b9c7, 0x3fffffff, 0x3fffffff, false, false) // CRC-30/CDMA

	CRC31PHILIPS = mustNewPreset[uint32](31, 0x04c11db7, 0x7fffffff, 0x7fffffff, false, false) // CRC-31/PHILIPS

	CRC32XFER     = mustNewPreset[uint32](32, 0x000000af, 0x00000000, 0x00000000, false, false) // CRC-32/XFER
	CRC32CKSUM    = mustNewPreset[uint32](32, 0x04c11db7, 0x00000000, 0xffffffff, false, false) // CRC-32/CKSUM     Alias: CKSUM, CRC-32/POSIX
	CRC32MPEG2    = mustNewPreset[uint32](32, 0x04c11db7, 0xffffffff, 0x00000000, false, false) // CRC-32/MPEG-2
	CRC32BZIP2    = mustNewPreset[uint32](32, 0x04c11db7, 0xffffffff, 0xffffffff, false, false) // CRC-32/BZIP2     Alias: CRC-32/AAL5, CRC-32/DECT-B, B-CRC-32
	CRC32JAMCRC   = mustNewPreset[uint32](32, 0x04c11db7, 0xffffffff, 0x00000000, true, true)   // CRC-32/JAMCRC    Alias: JAMCRC
	CRC32ISOHDLC  = mustNewPreset[uint32](32, 0x04c11db7, 0xffffffff, 0xffffffff, true, true)   // CRC-32/ISO-HDLC  Alias: CRC-32, CRC-32/ADCCP, CRC-32/V-42, CRC-32/XZ, PKZIP
	CRC32ISCSI    = mustNewPreset[uint32](32, 0x1edc6f41, 0xffffffff, 0xffffffff, true, true)   // CRC-32/ISCSI     Alias: CRC-32/BASE91-C, CRC-32/CASTAGNOLI, CRC-32/INTERLAKEN, CRC-32C
	CRC32MEF      = mustNewPreset[uint32](32, 0x741b8cd7, 0xffffffff, 0x00000000, true, true)   // CRC-32/MEF       Note: this algorithm uses Koopman's polynomial
	CRC32CDROMEDC = mustNewPreset[uint32](32, 0x8001801b, 0x00000000, 0x00000000, true, true)   // CRC-32/CD-ROM-EDC
	CRC32AIXM     = mustNewPreset[uint32](32, 0x814141ab, 0x00000000, 0x00000000, false, false) // CRC-32/AIXM      Alias: CRC-32Q
	CRC32BASE91D  = mustNewPreset[uint32](32, 0xa833982b, 0xffffffff, 0xffffffff, true, true)   // CRC-32/BASE91-D  Alias: CRC-32D
	CRC32AUTOSAR  = mustNewPreset[uint32](32, 0xf4acfb13, 0xffffffff, 0xffffffff, true, true)   // CRC-32/AUTOSAR

	CRC40GSM = mustNewPreset[uint64](40, 0x0004820009, 0x0000000000, 0xffffffffff, false, false) // CRC-40/GSM

	CRC64GOISO   = mustNewPreset[uint64](64, 0x000000000000001b, 0xffffffffffffffff, 0xffffffffffffffff, true, true)   // CRC-64/GO-ISO
	CRC64MS      = mustNewPreset[uint64](64, 0x259c84cba6426349, 0xffffffffffffffff, 0x0000000000000000, true, true)   // CRC-64/MS
	CRC64ECMA182 = mustNewPreset[uint64](64, 0x42f0e1eba9ea3693, 0x0000000000000000, 0x0000000000000000, false, false) // CRC-64/ECMA-182  Alias: CRC-64
	CRC64WE      = mustNewPreset[uint64](64, 0x42f0e1eba9ea3693, 0xffffffffffffffff, 0xffffffffffffffff, false, false) // CRC-64/WE
	CRC64XZ      = mustNewPreset[uint64](64, 0x42f0e1eba9ea3693, 0xffffffffffffffff, 0xffffffffffffffff, true, true)   // CRC-64/XZ        Alias: CRC-64/GO-ECMA
	CRC64REDIS   = mustNewPreset[uint64](64, 0xad93d23594c935a9, 0x0000000000000000, 0x0000000000000000, true, true)   // CRC-64/REDIS
)
