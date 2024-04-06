Arbitrary-precision CRC calculator in golang
============================================

Can calculate CRCs of any bit width (between CRC-1 and CRC-64) and can process
input of any bit length. Automatically creates 256-entry accelerator tables for
the used CRC algorithms. Provides presets for and has been tested against the
[100+ CRC algorithms listed in Greg Cook's CRC catalogue](https://reveng.sourceforge.io/crc-catalogue/all.htm).

```go
import "github.com/pasztorpisti/go-crc"

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
```

[Here is the godoc](https://pkg.go.dev/github.com/pasztorpisti/go-crc)
that you probably don't need.
