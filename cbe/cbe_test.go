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
	"testing"
)

func TestCBEVersion(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01}, BD(), V(1), ED())
}

func TestCBEPadding(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, typePadding}, BD(), V(1), PAD(1), ED())
	assertDecodeEncode(t, []byte{0x01, typePadding, typePadding, typePadding}, BD(), V(1), PAD(1), PAD(1), PAD(1), ED())
}

func TestCBENil(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, typeNil}, BD(), V(1), N(), ED())
}

func TestCBEBool(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, typeTrue}, BD(), V(1), B(true), ED())
	assertDecodeEncode(t, []byte{0x01, typeFalse}, BD(), V(1), B(false), ED())
}

func TestCBEIntEOF(t *testing.T) {
	assertDecodeFails(t, []byte{0x01, typePosInt8})
	assertDecodeFails(t, []byte{0x01, typeNegInt8})
	assertDecodeFails(t, []byte{0x01, typePosInt16, 0x01})
	assertDecodeFails(t, []byte{0x01, typeNegInt16, 0x01})
	assertDecodeFails(t, []byte{0x01, typePosInt32, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{0x01, typeNegInt32, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{0x01, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{0x01, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{0x01, typePosInt, 0x01})
	assertDecodeFails(t, []byte{0x01, typeNegInt, 0x05})
	assertDecodeFails(t, []byte{0x01, typeNegInt, 0xff})
}

func TestCBEPositiveInt(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, 0}, BD(), V(1), PI(0), ED())
	assertDecodeEncode(t, []byte{0x01, 100}, BD(), V(1), PI(100), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt8, 101}, BD(), V(1), PI(101), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt8, 0xff}, BD(), V(1), PI(255), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt16, 0x00, 0x01}, BD(), V(1), PI(0x100), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt16, 0xff, 0xff}, BD(), V(1), PI(0xffff), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt32, 0x00, 0x00, 0x01, 0x00}, BD(), V(1), PI(0x10000), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt32, 0xff, 0xff, 0xff, 0xff}, BD(), V(1), PI(0xffffffff), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt, 0x05, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), PI(0x100000000), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), PI(0x10000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, BD(), V(1), PI(0x1000000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), PI(0x100000000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), V(1), PI(0xffffffffffffffff), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), BI(NewBigInt("18446744073709551616")), ED())
	assertDecodeEncode(t, []byte{0x01, typePosInt, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), BI(NewBigInt("4722366482869645213696")), ED())
}

func TestCBENegativeInt(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, 0xff}, BD(), V(1), NI(1), ED())
	assertDecodeEncode(t, []byte{0x01, 0x9c}, BD(), V(1), NI(100), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt8, 101}, BD(), V(1), NI(101), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt8, 0xff}, BD(), V(1), NI(255), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt16, 0x00, 0x01}, BD(), V(1), NI(0x100), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt16, 0xff, 0xff}, BD(), V(1), NI(0xffff), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt32, 0x00, 0x00, 0x01, 0x00}, BD(), V(1), NI(0x10000), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt32, 0xff, 0xff, 0xff, 0xff}, BD(), V(1), NI(0xffffffff), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt, 0x05, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), NI(0x100000000), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), NI(0x10000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, BD(), V(1), NI(0x1000000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), NI(0x100000000000000), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), V(1), NI(0xffffffffffffffff), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), BI(NewBigInt("-18446744073709551616")), ED())
	assertDecodeEncode(t, []byte{0x01, typeNegInt, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), V(1), BI(NewBigInt("-4722366482869645213696")), ED())
}

func TestCBEBinaryFloatEOF(t *testing.T) {
	assertDecodeFails(t, []byte{0x01, typeFloat16, 0xd1})
	assertDecodeFails(t, []byte{0x01, typeFloat32, 0xd1, 0x00, 0x00})
	assertDecodeFails(t, []byte{0x01, typeFloat64, 0xd1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestCBEBinaryFloat(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, typeFloat16, 0xd1, 0x17}, BD(), V(1), F(0x1.a2p-80), ED())
	assertDecodeEncode(t, []byte{0x01, typeFloat32, 0x80, 0xf4, 0xa7, 0x71}, BD(), V(1), F(0x1.4fe9p100), ED())
	assertDecodeEncode(t, []byte{0x01, typeFloat64, 0x00, 0x00, 0xc2, 0x99, 0x91, 0xfe, 0xb4, 0x20}, BD(), V(1), F(0x1.4fe9199c2p-500), ED())
}

// func TestCBEDecimalFloatEOF(t *testing.T) {
// 	assertDecodeFails(t, []byte{0x01, typeDecimal, 0x04})
// }

// func TestCBEDecimalFloat(t *testing.T) {
// 	assertDecodeEncode(t, []byte{0x01, typeDecimal, 0x00}, BD(), V(1), DF(NewDFloat("1.2")), ED())
// }

func TestCBEUUIDEOF(t *testing.T) {
	assertDecodeFails(t, []byte{0x01, typeUUID, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
}

func TestCBEUUID(t *testing.T) {
	assertDecodeEncode(t, []byte{0x01, typeUUID, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		BD(), V(1), UUID([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}), ED())
}
