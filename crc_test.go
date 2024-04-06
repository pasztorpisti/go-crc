// SPDX-License-Identifier: MIT-0
// SPDX-FileCopyrightText:  2024 Istvan Pasztor

package crc_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/pasztorpisti/go-crc"
)

// This example demonstrates the use of the CRC-5/USB preset and a custom
// CRC-16 algorithm.
func Example() {
	// Using CRC-5/USB to calculate the CRC of a byte slice:
	fmt.Printf("usb1: %#x\n", crc.CRC5USB.Calc([]byte("123456789")))

	// Calculating the CRC when the data arrives in chunks:
	c := crc.CRC5USB.NewCRC()
	c.UpdateBits([]byte("12345"), 8*4+2)
	// The previous call consumed four bytes and the two least significant bits
	// of the last byte. '5' is 0b00110101 in binary so that update would have
	// had the same effect with inputs like "1234\x01" and "1234\xfd".
	// The call below provides the 6 most significant bits of the '5'.
	c.UpdateBits([]byte{0b001101}, 6)
	c.Update([]byte("6789"))
	fmt.Printf("usb2: %#x\n", c.Final())

	// Custom polynomial:
	// 0xa2eb was picked from the CRC Polynomial Zoo:
	// https://users.ece.cmu.edu/~koopman/crc/crc16.html
	a, err := crc.NewAlgo[uint16](16, 0xa2eb, 0xffff, 0xffff, true, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("zoo/a2eb: %#x\n", a.Calc([]byte("123456789")))

	// Output:
	// usb1: 0x19
	// usb2: 0x19
	// zoo/a2eb: 0x4e4c
}

type crc64[T crc.UInt] struct {
	crc.CRC[T]
}

func (c *crc64[T]) Final() uint64 {
	return uint64(c.CRC.Final())
}

func (c *crc64[T]) Residue() uint64 {
	return uint64(c.CRC.Residue())
}

type algo64[T crc.UInt] struct {
	algo crc.Algo[T]
}

func (a *algo64[T]) NewCRC() crc.CRC[uint64] {
	return &crc64[T]{a.algo.NewCRC()}
}

func (a *algo64[T]) Calc(data []byte) uint64 {
	return uint64(a.algo.Calc(data))
}

func (a *algo64[T]) CalcBits(data []byte, bitLen int) uint64 {
	return uint64(a.algo.CalcBits(data, bitLen))
}

var presets = []struct {
	name           string
	preset         crc.Algo[uint64]
	check          uint64
	residue        uint64
	codeWord       string
	codeWordBitLen int
}{
	{"CRC3GSM", &algo64[uint8]{crc.CRC3GSM}, 0x4, 0x2, "CRC3GSM\xe0", 59},
	{"CRC3ROHC", &algo64[uint8]{crc.CRC3ROHC}, 0x6, 0x0, "CRC3ROHC\x06", 67},

	{"CRC4INTERLAKEN", &algo64[uint8]{crc.CRC4INTERLAKEN}, 0xb, 0x2, "CRC4INTERLAKEN\x40", 116},
	{"CRC4G704", &algo64[uint8]{crc.CRC4G704}, 0x7, 0x0, "CRC4G704\x09", 68},

	{"CRC5USB", &algo64[uint8]{crc.CRC5USB}, 0x19, 0x06, "CRC5USB\x0d", 61},
	{"CRC5EPCC1G2", &algo64[uint8]{crc.CRC5EPCC1G2}, 0x00, 0x00, "CRC5EPCC1G2\xc0", 93},
	{"CRC5G704", &algo64[uint8]{crc.CRC5G704}, 0x07, 0x00, "CRC5G704\x02", 69},

	{"CRC6G704", &algo64[uint8]{crc.CRC6G704}, 0x06, 0x00, "CRC6G704\x0b", 70},
	{"CRC6CDMA2000B", &algo64[uint8]{crc.CRC6CDMA2000B}, 0x3b, 0x00, "CRC6CDMA2000B\xec", 110},
	{"CRC6DARC", &algo64[uint8]{crc.CRC6DARC}, 0x26, 0x00, "CRC6DARC\x02", 70},
	{"CRC6CDMA2000A", &algo64[uint8]{crc.CRC6CDMA2000A}, 0x0d, 0x00, "CRC6CDMA2000A\x7c", 110},
	{"CRC6GSM", &algo64[uint8]{crc.CRC6GSM}, 0x13, 0x3a, "CRC6GSMT", 62},

	{"CRC7MMC", &algo64[uint8]{crc.CRC7MMC}, 0x75, 0x00, "CRC7MMC\xae", 63},
	{"CRC7UMTS", &algo64[uint8]{crc.CRC7UMTS}, 0x61, 0x00, "CRC7UMTS\x94", 71},
	{"CRC7ROHC", &algo64[uint8]{crc.CRC7ROHC}, 0x53, 0x00, "CRC7ROHC\x1e", 71},

	{"CRC8SMBUS", &algo64[uint8]{crc.CRC8SMBUS}, 0xf4, 0x00, "CRC8SMBUS\x0f", 80},
	{"CRC8I4321", &algo64[uint8]{crc.CRC8I4321}, 0xa1, 0xac, "CRC8I4321\x9a", 80},
	{"CRC8ROHC", &algo64[uint8]{crc.CRC8ROHC}, 0xd0, 0x00, "CRC8ROHC\x26", 72},
	{"CRC8GSMA", &algo64[uint8]{crc.CRC8GSMA}, 0x37, 0x00, "CRC8GSMA\xeb", 72},
	{"CRC8MIFAREMAD", &algo64[uint8]{crc.CRC8MIFAREMAD}, 0x99, 0x00, "CRC8MIFAREMAD\xed", 112},
	{"CRC8ICODE", &algo64[uint8]{crc.CRC8ICODE}, 0x7e, 0x00, "CRC8ICODE\x1f", 80},
	{"CRC8HITAG", &algo64[uint8]{crc.CRC8HITAG}, 0xb4, 0x00, "CRC8HITAG\xc7", 80},
	{"CRC8SAEJ1850", &algo64[uint8]{crc.CRC8SAEJ1850}, 0x4b, 0xc4, "CRC8SAEJ1850z", 104},
	{"CRC8TECH3250", &algo64[uint8]{crc.CRC8TECH3250}, 0x97, 0x00, "CRC8TECH3250A", 104},
	{"CRC8OPENSAFETY", &algo64[uint8]{crc.CRC8OPENSAFETY}, 0x3e, 0x00, "CRC8OPENSAFETYn", 120},
	{"CRC8AUTOSAR", &algo64[uint8]{crc.CRC8AUTOSAR}, 0xdf, 0x42, "CRC8AUTOSAR\xa7", 96},
	{"CRC8NRSC5", &algo64[uint8]{crc.CRC8NRSC5}, 0xf7, 0x00, "CRC8NRSC5\x06", 80},
	{"CRC8MAXIMDOW", &algo64[uint8]{crc.CRC8MAXIMDOW}, 0xa1, 0x00, "CRC8MAXIMDOW\x99", 104},
	{"CRC8DARC", &algo64[uint8]{crc.CRC8DARC}, 0x15, 0x00, "CRC8DARCw", 72},
	{"CRC8GSMB", &algo64[uint8]{crc.CRC8GSMB}, 0x94, 0x53, "CRC8GSMB\x93", 72},
	{"CRC8LTE", &algo64[uint8]{crc.CRC8LTE}, 0xea, 0x00, "CRC8LTE\xe3", 64},
	{"CRC8CDMA2000", &algo64[uint8]{crc.CRC8CDMA2000}, 0xda, 0x00, "CRC8CDMA2000\xbd", 104},
	{"CRC8WCDMA", &algo64[uint8]{crc.CRC8WCDMA}, 0x25, 0x00, "CRC8WCDMA\xb1", 80},
	{"CRC8BLUETOOTH", &algo64[uint8]{crc.CRC8BLUETOOTH}, 0x26, 0x00, "CRC8BLUETOOTHD", 112},
	{"CRC8DVBS2", &algo64[uint8]{crc.CRC8DVBS2}, 0xbc, 0x00, "CRC8DVBS2\x92", 80},

	{"CRC10GSM", &algo64[uint16]{crc.CRC10GSM}, 0x12a, 0x0c6, "CRC10GSM\xb7\x40", 74},
	{"CRC10ATM", &algo64[uint16]{crc.CRC10ATM}, 0x199, 0x000, "CRC10ATM\xdd\x80", 74},
	{"CRC10CDMA2000", &algo64[uint16]{crc.CRC10CDMA2000}, 0x233, 0x000, "CRC10CDMA2000\xe7\xc0", 114},

	{"CRC11UMTS", &algo64[uint16]{crc.CRC11UMTS}, 0x061, 0x000, "CRC11UMTS\x8d\xc0", 83},
	{"CRC11FLEXRAY", &algo64[uint16]{crc.CRC11FLEXRAY}, 0x5a3, 0x000, "CRC11FLEXRAY\xc3\x20", 107},

	{"CRC12DECT", &algo64[uint16]{crc.CRC12DECT}, 0xf5b, 0x000, "CRC12DECT\xd4\x90", 84},
	{"CRC12UMTS", &algo64[uint16]{crc.CRC12UMTS}, 0xdaf, 0x000, "CRC12UMTS\x10\xd0", 84},
	{"CRC12GSM", &algo64[uint16]{crc.CRC12GSM}, 0xb34, 0x178, "CRC12GSM\xcd\x00", 76},
	{"CRC12CDMA2000", &algo64[uint16]{crc.CRC12CDMA2000}, 0xd4d, 0x000, "CRC12CDMA2000\x89\xf0", 116},

	{"CRC13BBC", &algo64[uint16]{crc.CRC13BBC}, 0x04fa, 0x0000, "CRC13BBC\x17h", 77},

	{"CRC14DARC", &algo64[uint16]{crc.CRC14DARC}, 0x082d, 0x0000, "CRC14DARC\x1c\x3f", 86},
	{"CRC14GSM", &algo64[uint16]{crc.CRC14GSM}, 0x30ae, 0x031e, "CRC14GSM\xd4T", 78},

	{"CRC15CAN", &algo64[uint16]{crc.CRC15CAN}, 0x059e, 0x0000, "CRC15CANC\xf0", 79},
	{"CRC15MPT1327", &algo64[uint16]{crc.CRC15MPT1327}, 0x2566, 0x6815, "CRC15MPT1327\x07\xa0", 111},

	{"CRC16DECTX", &algo64[uint16]{crc.CRC16DECTX}, 0x007f, 0x0000, "CRC16DECTXm\xa1", 96},
	{"CRC16DECTR", &algo64[uint16]{crc.CRC16DECTR}, 0x007e, 0x0589, "CRC16DECTRJ\xfa", 96},
	{"CRC16NRSC5", &algo64[uint16]{crc.CRC16NRSC5}, 0xa066, 0x0000, "CRC16NRSC5\x27\x25", 96},
	{"CRC16XMODEM", &algo64[uint16]{crc.CRC16XMODEM}, 0x31c3, 0x0000, "CRC16XMODEM\xd2\x98", 104},
	{"CRC16GSM", &algo64[uint16]{crc.CRC16GSM}, 0xce3c, 0x1d0f, "CRC16GSM\x18\xa9", 80},
	{"CRC16SPIFUJITSU", &algo64[uint16]{crc.CRC16SPIFUJITSU}, 0xe5cc, 0x0000, "CRC16SPIFUJITSUvw", 136},
	{"CRC16IBM3740", &algo64[uint16]{crc.CRC16IBM3740}, 0x29b1, 0x0000, "CRC16IBM3740\xd8\xfe", 112},
	{"CRC16GENIBUS", &algo64[uint16]{crc.CRC16GENIBUS}, 0xd64e, 0x1d0f, "CRC16GENIBUSN\xe2", 112},
	{"CRC16KERMIT", &algo64[uint16]{crc.CRC16KERMIT}, 0x2189, 0x0000, "CRC16KERMIT1b", 104},
	{"CRC16TMS37157", &algo64[uint16]{crc.CRC16TMS37157}, 0x26b1, 0x0000, "CRC16TMS37157\xd6\xcb", 120},
	{"CRC16RIELLO", &algo64[uint16]{crc.CRC16RIELLO}, 0x63d0, 0x0000, "CRC16RIELLO\x8d\x09", 104},
	{"CRC16ISOIEC144433A", &algo64[uint16]{crc.CRC16ISOIEC144433A}, 0xbf05, 0x0000, "CRC16ISOIEC144433A\x07\xf5", 160},
	{"CRC16MCRF4XX", &algo64[uint16]{crc.CRC16MCRF4XX}, 0x6f91, 0x0000, "CRC16MCRF4XX\xa17", 112},
	{"CRC16IBMSDLC", &algo64[uint16]{crc.CRC16IBMSDLC}, 0x906e, 0xf0b8, "CRC16IBMSDLC2\x0c", 112},
	{"CRC16PROFIBUS", &algo64[uint16]{crc.CRC16PROFIBUS}, 0xa819, 0xe394, "CRC16PROFIBUS\xf6\xe2", 120},
	{"CRC16EN13757", &algo64[uint16]{crc.CRC16EN13757}, 0xc2b7, 0xa366, "CRC16EN13757\xf9K", 112},
	{"CRC16DNP", &algo64[uint16]{crc.CRC16DNP}, 0xea82, 0x66c5, "CRC16DNPj\x2e", 80},
	{"CRC16OPENSAFETYA", &algo64[uint16]{crc.CRC16OPENSAFETYA}, 0x5d38, 0x0000, "CRC16OPENSAFETYA\xd7\x7b", 144},
	{"CRC16M17", &algo64[uint16]{crc.CRC16M17}, 0x772b, 0x0000, "CRC16M17\x10\xfd", 80},
	{"CRC16LJ1200", &algo64[uint16]{crc.CRC16LJ1200}, 0xbdf4, 0x0000, "CRC16LJ1200x\x9a", 104},
	{"CRC16OPENSAFETYB", &algo64[uint16]{crc.CRC16OPENSAFETYB}, 0x20fe, 0x0000, "CRC16OPENSAFETYB\x9c\xa9", 144},
	{"CRC16UMTS", &algo64[uint16]{crc.CRC16UMTS}, 0xfee8, 0x0000, "CRC16UMTS\xfd\xd4", 88},
	{"CRC16DDS110", &algo64[uint16]{crc.CRC16DDS110}, 0x9ecf, 0x0000, "CRC16DDS110\xfa\x81", 104},
	{"CRC16CMS", &algo64[uint16]{crc.CRC16CMS}, 0xaee7, 0x0000, "CRC16CMS\xf6\x04", 80},
	{"CRC16ARC", &algo64[uint16]{crc.CRC16ARC}, 0xbb3d, 0x0000, "CRC16ARCg\xda", 80},
	{"CRC16MAXIMDOW", &algo64[uint16]{crc.CRC16MAXIMDOW}, 0x44c2, 0xb001, "CRC16MAXIMDOW\x2f\x29", 120},
	{"CRC16MODBUS", &algo64[uint16]{crc.CRC16MODBUS}, 0x4b37, 0x0000, "CRC16MODBUS\xde\x98", 104},
	{"CRC16USB", &algo64[uint16]{crc.CRC16USB}, 0xb4c8, 0xb001, "CRC16USBXz", 80},
	{"CRC16T10DIF", &algo64[uint16]{crc.CRC16T10DIF}, 0xd0db, 0x0000, "CRC16T10DIF\xef\xdb", 104},
	{"CRC16TELEDISK", &algo64[uint16]{crc.CRC16TELEDISK}, 0x0fb3, 0x0000, "CRC16TELEDISK\xaeG", 120},
	{"CRC16CDMA2000", &algo64[uint16]{crc.CRC16CDMA2000}, 0x4c06, 0x0000, "CRC16CDMA2000\x0a\xd4", 120},

	{"CRC17CANFD", &algo64[uint32]{crc.CRC17CANFD}, 0x04f03, 0x00000, "CRC17CANFD\xdc2\x80", 97},

	{"CRC21CANFD", &algo64[uint32]{crc.CRC21CANFD}, 0x0ed841, 0x000000, "CRC21CANFD\xa1\x2e\xb8", 101},

	{"CRC24BLE", &algo64[uint32]{crc.CRC24BLE}, 0xc25a56, 0x000000, "CRC24BLE\x0f\xaas", 88},
	{"CRC24INTERLAKEN", &algo64[uint32]{crc.CRC24INTERLAKEN}, 0xb4f3e6, 0x144e63, "CRC24INTERLAKEN\xbc\xba\xb3", 144},
	{"CRC24FLEXRAYB", &algo64[uint32]{crc.CRC24FLEXRAYB}, 0x1f23b8, 0x000000, "CRC24FLEXRAYBX\x60\xee", 128},
	{"CRC24FLEXRAYA", &algo64[uint32]{crc.CRC24FLEXRAYA}, 0x7979bd, 0x000000, "CRC24FLEXRAYA\xd1\xc3\x86", 128},
	{"CRC24LTEB", &algo64[uint32]{crc.CRC24LTEB}, 0x23ef52, 0x000000, "CRC24LTEBz\xe3\x84", 96},
	{"CRC24OS9", &algo64[uint32]{crc.CRC24OS9}, 0x200fa5, 0x800fe3, "CRC24OS9\x7c\xa8\xfa", 88},
	{"CRC24LTEA", &algo64[uint32]{crc.CRC24LTEA}, 0xcde703, 0x000000, "CRC24LTEA\x7d\xd6\xab", 96},
	{"CRC24OPENPGP", &algo64[uint32]{crc.CRC24OPENPGP}, 0x21cf02, 0x000000, "CRC24OPENPGP\xf3\x27\x1c", 120},

	{"CRC30CDMA", &algo64[uint32]{crc.CRC30CDMA}, 0x04c34abf, 0x34efa55a, "CRC30CDMA\x90\x22h\x40", 102},

	{"CRC31PHILIPS", &algo64[uint32]{crc.CRC31PHILIPS}, 0x0ce9e46c, 0x4eaf26f1, "CRC31PHILIPSoL\x18\x12", 127},

	{"CRC32XFER", &algo64[uint32]{crc.CRC32XFER}, 0xbd0be338, 0x00000000, "CRC32XFER\x05\x9f\x1fZ", 104},
	{"CRC32CKSUM", &algo64[uint32]{crc.CRC32CKSUM}, 0x765e7680, 0xc704dd7b, "CRC32CKSUM\x25\x11Y\x8e", 112},
	{"CRC32MPEG2", &algo64[uint32]{crc.CRC32MPEG2}, 0x0376e6e7, 0x00000000, "CRC32MPEG2\xa7\x88\xc25", 112},
	{"CRC32BZIP2", &algo64[uint32]{crc.CRC32BZIP2}, 0xfc891918, 0xc704dd7b, "CRC32BZIP2\x89\xb4\x92F", 112},
	{"CRC32JAMCRC", &algo64[uint32]{crc.CRC32JAMCRC}, 0x340bc6d9, 0x00000000, "CRC32JAMCRC\xd9\x7c8\x02", 120},
	{"CRC32ISOHDLC", &algo64[uint32]{crc.CRC32ISOHDLC}, 0xcbf43926, 0xdebb20e3, "CRC32ISOHDLC\xb8\x13\x23\xa2", 128},
	{"CRC32ISCSI", &algo64[uint32]{crc.CRC32ISCSI}, 0xe3069283, 0xb798b438, "CRC32ISCSI\x0ay\xd9\x83", 112},
	{"CRC32MEF", &algo64[uint32]{crc.CRC32MEF}, 0xd2c22f51, 0x00000000, "CRC32MEFq\xdf\xd8\x1a", 96},
	{"CRC32CDROMEDC", &algo64[uint32]{crc.CRC32CDROMEDC}, 0x6ec2edc4, 0x00000000, "CRC32CDROMEDCjZY\x08", 136},
	{"CRC32AIXM", &algo64[uint32]{crc.CRC32AIXM}, 0x3010bf7f, 0x00000000, "CRC32AIXM\x1ae\x05\xe9", 104},
	{"CRC32BASE91D", &algo64[uint32]{crc.CRC32BASE91D}, 0x87315576, 0x45270551, "CRC32BASE91D\x03\xa4\x11\x22", 128},
	{"CRC32AUTOSAR", &algo64[uint32]{crc.CRC32AUTOSAR}, 0x1697d06a, 0x904cddbf, "CRC32AUTOSARj\xbaq\xe2", 128},

	{"CRC40GSM", &algo64[uint64]{crc.CRC40GSM}, 0xd4164fc646, 0xc4ff8071ff, "CRC40GSM\xf9\xaf6\xf3\x87", 104},

	{"CRC64GOISO", &algo64[uint64]{crc.CRC64GOISO}, 0xb90956c775a41001, 0x5300000000000000, "CRC64GOISO1\x17\xc4\x07\xaa\x93\xd2r", 144},
	{"CRC64MS", &algo64[uint64]{crc.CRC64MS}, 0x75d4b74f024eceea, 0x0000000000000000, "CRC64MS\x21\x1d\x84\x0eC\x7d\xb9\xe9", 120},
	{"CRC64ECMA182", &algo64[uint64]{crc.CRC64ECMA182}, 0x6c40df5f0b497347, 0x0000000000000000, "CRC64ECMA1821\xec\x21\x1f\x0f\x40E6", 160},
	{"CRC64WE", &algo64[uint64]{crc.CRC64WE}, 0x62ec59e3f1a4f00a, 0xfcacbebd5931a992, "CRC64WE\x9d\x02\xc9\x5c\xfb\xfcpG", 120},
	{"CRC64XZ", &algo64[uint64]{crc.CRC64XZ}, 0x995dc9bbdf1939fa, 0x49958c9abd7d353f, "CRC64XZ\x40\x8a\xa6\xc4\x0fFz\xd8", 120},
	{"CRC64REDIS", &algo64[uint64]{crc.CRC64REDIS}, 0xe9c6d914c4b8d9ca, 0x0000000000000000, "CRC64REDIS\xd0DOjw\x01\xbe\xa2", 144},
}

func TestCRC(t *testing.T) {
	for _, p := range presets {
		t.Run(p.name, func(t *testing.T) {
			c := p.preset.CalcBits([]byte("123456789"), -1)
			if c != p.check {
				t.Errorf("check=%x, want %x", c, p.check)
			}
		})
	}
}

func TestResidue(t *testing.T) {
	for _, p := range presets {
		t.Run(p.name, func(t *testing.T) {
			c := p.preset.NewCRC()
			c.UpdateBits([]byte(p.codeWord), p.codeWordBitLen)
			r := c.Residue()
			if r != p.residue {
				t.Errorf("residue=%x, want %x", r, p.residue)
			}
		})
	}
}

func Benchmark_CRC8_Calc_100MB(b *testing.B) {
	data := make([]byte, 100*1024*1024)
	rand.New(rand.NewSource(42)).Read(data)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crc.CRC8.Calc(data)
	}
}

func Benchmark_CRC16_Calc_100MB(b *testing.B) {
	data := make([]byte, 100*1024*1024)
	rand.New(rand.NewSource(42)).Read(data)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crc.CRC16.Calc(data)
	}
}

func Benchmark_CRC32_Calc_100MB(b *testing.B) {
	data := make([]byte, 100*1024*1024)
	rand.New(rand.NewSource(42)).Read(data)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crc.CRC32.Calc(data)
	}
}

func Benchmark_CRC64_Calc_100MB(b *testing.B) {
	data := make([]byte, 100*1024*1024)
	rand.New(rand.NewSource(42)).Read(data)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crc.CRC64.Calc(data)
	}
}

// Getting rid of some unhelpful "Unused variable" warnings.
var _ = unused(
	crc.CRC32C, crc.CRC32D, crc.CRC32Q, crc.A, crc.B, crc.X25, crc.CRC16X25,
	crc.XMODEM, crc.KERMIT, crc.CRC16CCITT, crc.CRC16CCITTFALSE, crc.CRC16AUGCCITT,
	crc.V41LSB, crc.V41MSB, crc.PKZIP, crc.V42, crc.XZ, crc.POSIX, crc.CASTAGNOLI,
)

func unused(_ ...any) int { return 0 }
