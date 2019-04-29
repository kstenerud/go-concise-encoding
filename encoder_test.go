package cbe

import "bytes"
import "testing"


func assertEncoded(t* testing.T, encoder *Encoder, expected []byte) {
	actual := encoder.Encoded()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func TestPadding(t *testing.T) {
	assertEncoded(t, New(9).Padding(1),[]byte{127})
	assertEncoded(t, New(9).Padding(2),[]byte{127, 127})
}

func TestNil(t *testing.T) {
	assertEncoded(t, New(9).Nil(),[]byte{126})
}

func TestIntSmall(t *testing.T) {
	assertEncoded(t, New(9).Int(0), []byte{0})
	assertEncoded(t, New(9).Int(1), []byte{1})
	assertEncoded(t, New(9).Int(105), []byte{105})
	assertEncoded(t, New(9).Int(106), []byte{106})
	assertEncoded(t, New(9).Int(-1), []byte{255})
	assertEncoded(t, New(9).Int(-105), []byte{151})
	assertEncoded(t, New(9).Int(-106), []byte{150})
}

func TestInt8(t *testing.T) {
	assertEncoded(t, New(9).Int(107), []byte{0x70, 0x6b})
	assertEncoded(t, New(9).Int(-107), []byte{0x78, 0x6b})
	assertEncoded(t, New(9).Int(255), []byte{0x70, 0xff})
	assertEncoded(t, New(9).Int(-255), []byte{0x78, 0xff})

}

func TestInt16(t *testing.T) {
	assertEncoded(t, New(9).Int(0x100), []byte{0x71, 0x00, 0x01})
	assertEncoded(t, New(9).Int(0x7fff), []byte{0x71, 0xff, 0x7f})
	assertEncoded(t, New(9).Int(-0x7fff), []byte{0x79, 0xff, 0x7f})
	assertEncoded(t, New(9).Int(0xffff), []byte{0x71, 0xff, 0xff})
	assertEncoded(t, New(9).Int(-0xffff), []byte{0x79, 0xff, 0xff})
}

func TestInt32(t *testing.T) {
	assertEncoded(t, New(9).Int(0x10000), []byte{0x72, 0x00, 0x00, 0x01, 0x00})
	assertEncoded(t, New(9).Int(0x7fffffff), []byte{0x72, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, New(9).Int(-0x7fffffff), []byte{0x7a, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, New(9).Int(0xffffffff), []byte{0x72, 0xff, 0xff, 0xff, 0xff})
	assertEncoded(t, New(9).Int(-0xffffffff), []byte{0x7a, 0xff, 0xff, 0xff, 0xff})
}

func TestInt64(t *testing.T) {
	assertEncoded(t, New(9).Int(0x100000000), []byte{0x73, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00})
	assertEncoded(t, New(9).Int(0x7fffffffffffffff), []byte{0x73, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	assertEncoded(t, New(9).Int(-0x7fffffffffffffff), []byte{0x7b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
}

func TestFloat64(t *testing.T) {
	assertEncoded(t, New(9).Float(1.0123), []byte{0x6e, 0x51, 0xda, 0x1b, 0x7c, 0x61, 0x32, 0xf0, 0x3f})
}

func TestList(t *testing.T) {
	assertEncoded(t, New(9).ListBegin(), []byte{0x93})
	assertEncoded(t, New(9).ListBegin().ListEnd(), []byte{0x93, 0x95})
	// TODO: More
}

func TestMap(t *testing.T) {
	assertEncoded(t, New(9).MapBegin(), []byte{0x94})
	assertEncoded(t, New(9).MapBegin().ListEnd(), []byte{0x94, 0x95})
	// TODO: More
}

func TestBinary(t *testing.T) {
	assertEncoded(t, New(9).Binary([]byte{}), []byte{0x91, 0x00})
	assertEncoded(t, New(9).Binary([]byte{1}), []byte{0x91, 0x04, 0x01})
	assertEncoded(t, New(9).Binary([]byte{1, 2}), []byte{0x91, 0x08, 0x01, 0x02})
	assertEncoded(t, New(9).Binary([]byte{1, 2, 3}), []byte{0x91, 0x0c, 0x01, 0x02, 0x03})
	// TODO: Longer than 64 bytes
}

func TestString(t *testing.T) {
	assertEncoded(t, New(9).String(""), []byte{0x80})
	assertEncoded(t, New(9).String("0"), []byte{0x81, 0x30})
	assertEncoded(t, New(9).String("01"), []byte{0x82, 0x30, 0x31})
	assertEncoded(t, New(9).String("0123456789012345"), []byte{
		0x90, 0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
	})
	// TODO: Longer than 64 bytes
}
