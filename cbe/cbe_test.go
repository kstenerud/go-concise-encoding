// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package cbe

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCBERIDEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeRID, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeRID, 0x04, 'a'})
}

func TestCBERID(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeRID, 0x00}, BD(), EvV, RID(""), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeRID, 0x02, 'a'}, BD(), EvV, RID("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeRID, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, RID("00000000001111111111"), ED())
}

func TestCBECustomBinaryEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeCustomBinary, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeCustomBinary, 0x04, 'a'})
}

func TestCBECustomBinary(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x00}, BD(), EvV, CUB([]byte{}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x02, 'a'}, BD(), EvV, CUB([]byte("a")), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, CUB([]byte("00000000001111111111")), ED())
}

func TestCBECustomTextEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeCustomText, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeCustomText, 0x04, 'a'})
}

func TestCBECustomText(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x00}, BD(), EvV, CUT(""), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x02, 'a'}, BD(), EvV, CUT("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, CUT("00000000001111111111"), ED())
}

func TestCBEArrayUint8EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeArrayUint8, 0x04, 0xfa})
}

func TestCBEArrayUint8(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayUint8, 0x02, 0x01}, BD(), EvV, AU8([]byte{1}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayUint8, 0x04, 0xfa, 0x11}, BD(), EvV, AU8([]byte{0xfa, 0x11}), ED())
}

func TestCBEArrayUint16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint16, 0x02, 0xfa})
}

func TestCBEArrayUint16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint16 | 1, 0x01, 0x02}, BD(), EvV, AU16([]uint16{0x0201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint16 | 2, 0xfa, 0x11, 0x01, 0x02}, BD(), EvV, AU16([]uint16{0x11fa, 0x0201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayUint16, 0x020,
		0xfa, 0x11, 0x01, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, 0x00},
		BD(), EvV, AU16([]uint16{0x11fa, 0x0201, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000,
			0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0001}), ED())
}

func TestCBEArrayUint32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayUint32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint32 | 1, 0x01, 0x02, 0x03, 0x04}, BD(), EvV, AU32([]uint32{0x04030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint32 | 2, 1, 2, 3, 4, 5, 6, 7, 8}, BD(), EvV, AU32([]uint32{0x04030201, 0x08070605}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayUint32, 0x20,
		1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		BD(), EvV, AU32([]uint32{0x04030201, 0x08070605, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}), ED())
}

func TestCBEArrayUint64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayUint64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint64 | 1, 1, 2, 3, 4, 5, 6, 7, 8}, BD(), EvV, AU64([]uint64{0x0807060504030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayUint64 | 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}, BD(), EvV, AU64([]uint64{0x0807060504030201, 0x0605040302010009}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayUint64, 0x20,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0},
		BD(), EvV, AU64([]uint64{0x0807060504030201, 0x0605040302010009, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}), ED())
}

func TestCBEArrayInt8EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt8, 0x04, 0xfa})
}

func TestCBEArrayInt8(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt8 | 1, 0x01}, BD(), EvV, AI8([]int8{1}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt8 | 2, 0xfa, 0x11}, BD(), EvV, AI8([]int8{-6, 0x11}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayInt8, 0x20, 0xfa, 0x11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff},
		BD(), EvV, AI8([]int8{-6, 0x11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1}), ED())
}

func TestCBEArrayInt16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt16, 0x02, 0xfa})
}

func TestCBEArrayInt16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt16 | 1, 0x01, 0x02}, BD(), EvV, AI16([]int16{0x0201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt16 | 2, 0xfa, 0x11, 0x9c, 0xff}, BD(), EvV, AI16([]int16{0x11fa, -100}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayInt16, 0x20,
		0xfa, 0x11, 0x9c, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xfe, 0xff},
		BD(), EvV, AI16([]int16{0x11fa, -100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -2}), ED())
}

func TestCBEArrayInt32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayInt32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt32 | 1, 0x01, 0x02, 0x03, 0x04}, BD(), EvV, AI32([]int32{0x04030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt32 | 2, 1, 2, 3, 4, 0x9c, 0xff, 0xff, 0xff}, BD(), EvV, AI32([]int32{0x04030201, -100}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayInt32, 0x20,
		1, 2, 3, 4, 0x9c, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xfe, 0xff, 0xff, 0xff},
		BD(), EvV, AI32([]int32{0x04030201, -100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -2}), ED())
}

func TestCBEArrayInt64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayInt64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt64 | 1, 1, 2, 3, 4, 5, 6, 7, 7}, BD(), EvV, AI64([]int64{0x0707060504030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayInt64 | 2, 1, 2, 3, 4, 5, 6, 7, 7, 0x9c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, AI64([]int64{0x0707060504030201, -100}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayInt64, 0x20,
		1, 2, 3, 4, 5, 6, 7, 7, 0x9c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		BD(), EvV, AI64([]int64{0x0707060504030201, -100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -2}), ED())
}

func TestCBEArrayFloat16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat16, 0x02, 0xfa})
}

func TestCBEArrayFloat16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat16 | 1, 0x01, 0x02}, BD(), EvV, AF16([]byte{0x01, 0x02}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat16 | 2, 0xfa, 0x11, 0x9c, 0xff}, BD(), EvV, AF16([]byte{0xfa, 0x11, 0x9c, 0xff}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayFloat16, 0x20,
		0xfa, 0x11, 0x9c, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BD(), EvV, AF16([]byte{0xfa, 0x11, 0x9c, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), ED())
}

func TestCBEArrayFloat32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayFloat32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat32 | 1, 0x00, 0x58, 0x93, 0x54}, BD(), EvV, AF32([]float32{0x4.9acp40}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat32 | 2, 0x00, 0xc8, 0xfc, 0xb5, 000, 0xb1, 0x48, 0xd0}, BD(), EvV, AF32([]float32{-0x1.f99p-20, -0xc.8b1p30}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayFloat32, 0x20,
		0x00, 0xc8, 0xfc, 0xb5, 000, 0xb1, 0x48, 0xd0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BD(), EvV, AF32([]float32{-0x1.f99p-20, -0xc.8b1p30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), ED())
}

func TestCBEArrayFloat64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayFloat64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat64 | 1, 0x66, 0x46, 0x74, 0x3e, 0x33, 0x16, 0x09, 0xc2}, BD(), EvV, AF64([]float64{-0xc.8b199f3a2333p30}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeShortArrayFloat64 | 2, 0x00, 0x00, 0xcc, 0xea, 0xf1, 0x7f, 0x32, 0xbe, 0x00, 0x10, 0x90, 0xea, 0xfc, 0x87, 0x18, 0xc3}, BD(), EvV, AF64([]float64{-0x4.9ffc7ab3p-30, -0x1.887fcea901p50}), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeArrayFloat64, 0x20,
		0x00, 0x00, 0xcc, 0xea, 0xf1, 0x7f, 0x32, 0xbe, 0x00, 0x10, 0x90, 0xea,
		0xfc, 0x87, 0x18, 0xc3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BD(), EvV, AF64([]float64{-0x4.9ffc7ab3p-30, -0x1.887fcea901p50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), ED())
}

func TestCBEArrayBitEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeArrayBit, 0x12, 0x00})
}

func TestCBEArrayBit(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x02, 0x01}, BD(), EvV, AB(1, []byte{0x01}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x08, 0x03}, BD(), EvV, AB(4, []byte{0x03}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x26, 0xfe, 0xc1, 0x03}, BD(), EvV, AB(19, []byte{0xfe, 0xc1, 0x03}), ED())
}

func TestCBEArrayUID(t *testing.T) {
	// TODO: TestCBEArrayUID
}

func TestCBEMarker(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeMarker, 1, 'x', typeString1, 'a'}, BD(), EvV, MARK("x"), S("a"), ED())
}

func TestCBEReference(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeList, typeMarker, 1, 'x', typeString1, 'a', typeReference, 1, 'x', typeEndContainer},
		BD(), EvV, L(), MARK("x"), S("a"), REF("x"), E(), ED())
}

func TestCBEContainers(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1, typeEndContainer}, BD(), EvV, L(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMap, 1, 1, typeEndContainer}, BD(), EvV, M(), I(1), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMarkup, 1, 'a', typeEndContainer, typeEndContainer}, BD(), EvV, MUP("a"), E(), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNode, typeTrue, 1, typeEndContainer}, BD(), EvV, NODE(), TT(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeEdge, 1, 2, 3}, BD(), EvV, EDGE(), I(1), I(2), I(3), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1,
		typeList, typeString1, 'a', typeEndContainer,
		typeMap, typeString1, 'a', 100, typeEndContainer,
		typeMarkup, 1, 'a', typeString1, 'a', 50, typeEndContainer, typeString1, 'a', typeEndContainer,
		typeNode, typeTrue, 1, typeEndContainer,
		typeEdge, 1, 2, 3,
		typeEndContainer,
	},
		BD(), EvV, L(), I(1),
		L(), S("a"), E(),
		M(), S("a"), I(100), E(),
		MUP("a"), S("a"), I(50), E(), S("a"), E(),
		NODE(), TT(), I(1), E(),
		EDGE(), I(1), I(2), I(3),
		E(), ED())
}

func TestCBEMultipartArray(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typeString, 0x03, 'a', 0x02, 'b'}, BD(), EvV, SB(), AC(1, true), AD([]byte{'a'}), AC(1, false), AD([]byte{'b'}), ED())
	assertDecode(t, nil, []byte{header, ceVer, typePlane2, typeArrayUint16, 0x03, 0x01, 0x02, 0x02, 0x03, 0x04}, BD(), EvV, AU16B(), AC(1, true), AD([]byte{0x01, 0x02}), AC(1, false), AD([]byte{0x03, 0x04}), ED())
}

func TestCBEChunkedArray(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	InvokeEvents(encoder, EvV, AU16B(), AC(2, true), AD([]byte{1, 0, 2, 0}), AC(2, false), AD([]byte{3, 0, 4, 0}), ED())

	expected := []byte{header, ceVer, typePlane2, typeArrayUint16, 0x05, 1, 0, 2, 0, 0x04, 3, 0, 4, 0}
	if !reflect.DeepEqual(buffer.Bytes(), expected) {
		t.Errorf("Expected first buffer %v but got %v", expected, buffer.Bytes())
	}
}

func TestCBEEncoderMultiUse(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	InvokeEvents(encoder, EvV, M(), E(), ED())

	buffer2 := &bytes.Buffer{}
	encoder.PrepareToEncode(buffer2)
	InvokeEvents(encoder, EvV, M(), E(), ED())

	expected := []byte{header, ceVer, typeMap, typeEndContainer}
	if !reflect.DeepEqual(buffer.Bytes(), expected) {
		t.Errorf("Expected first buffer %v but got %v", expected, buffer.Bytes())
	}
	if !reflect.DeepEqual(buffer2.Bytes(), expected) {
		t.Errorf("Expected second buffer %v but got %v", expected, buffer2.Bytes())
	}
}

func TestCBEDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertMarshalUnmarshal(t, v, []byte{header, ceVer, 0x7a, 0x7a, 0x7b, 0x7a, 0x7b, 0x7a, 0x7b, 0x7b})
}

func TestRemoteReference(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typePlane2, typeRemoteRef, 0x02, 'a'}, BD(), EvV, RREF("a"), ED())
}

func TestMultichunk(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeString, 0x03, 'a', 0}, BD(), EvV, SB(), AC(1, true), AD([]byte{'a'}), AC(0, false), ED())
}

func TestEdge(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer, typeEdge, typeRID, 0x02, 'a', typeRID, 0x02, 'b', 1},
		BD(), EvV, EDGE(), RID("a"), RID("b"), I(1), ED())
}

func TestNode(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer,
			typeNode, typeNull, typeString1, 'a', typeRID, 0x02, 'b',
			typeNode, typeNull, typeEndContainer, typeEndContainer},
		BD(), EvV, NODE(), NULL(), S("a"), RID("b"), NODE(), NULL(), E(), E(), ED())
}

func TestMedia(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer, typePlane2, typeMedia, 0x02, 'a', 0x00},
		BD(), EvV, MB(), AC(1, false), AD([]byte("a")), AC(0, false), ED())
}
