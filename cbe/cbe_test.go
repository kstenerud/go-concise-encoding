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
	"math"
	"reflect"
	"testing"
	"time"
)

// TODO: Remove this when releasing V1
func TestCBEVersion1(t *testing.T) {
	assertDecode(t, nil, []byte{header, 1, 1}, BD(), EvV, I(1), ED())
}

func TestCBEVersion(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, 1}, BD(), EvV, I(1), ED())
}

func TestCBEPadding(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePadding}, BD(), EvV, PAD(1), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePadding, typePadding, typePadding}, BD(), EvV, PAD(1), PAD(1), PAD(1), ED())
}

func TestCBENA(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeNA}, BD(), EvV, NA(), ED())
}

func TestCBENACat(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNACat, 1}, BD(), EvV, NACat(), I(1), ED())
}

func TestCBEBool(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeTrue}, BD(), EvV, TT(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeFalse}, BD(), EvV, FF(), ED())
}

func TestCBEIntEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePosInt8})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt8})
	assertDecodeFails(t, []byte{header, ceVer, typePosInt16, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt16, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typePosInt32, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt32, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typePosInt, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt, 0x05})
	assertDecodeFails(t, []byte{header, ceVer, typeNegInt, 0xff})
}

func TestCBEInt(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	encoder.OnVersion(1)
	encoder.OnList()
	encoder.OnInt(100)
	encoder.OnInt(-100)
	encoder.OnEnd()
	encoder.OnEndDocument()

	expected := []byte{header, 1, typeList, 100, 0x9c, typeEndContainer}
	if !reflect.DeepEqual(buffer.Bytes(), expected) {
		t.Errorf("Expected first buffer %v but got %v", expected, buffer.Bytes())
	}
}

func TestCBEPositiveInt(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, 0}, BD(), EvV, I(0), ED())
	assertDecodeEncode(t, []byte{header, ceVer, 100}, BD(), EvV, I(100), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt8, 101}, BD(), EvV, PI(101), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt8, 0xff}, BD(), EvV, PI(255), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt16, 0x00, 0x01}, BD(), EvV, PI(0x100), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt16, 0xff, 0xff}, BD(), EvV, PI(0xffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt32, 0x00, 0x00, 0x01, 0x00}, BD(), EvV, PI(0x10000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt32, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, PI(0xffffffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt, 0x05, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, PI(0x100000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, PI(0x10000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, BD(), EvV, PI(0x1000000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, PI(0x100000000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, PI(0xffffffffffffffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, BI(NewBigInt("18446744073709551616", 10)), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePosInt, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, BI(NewBigInt("4722366482869645213696", 10)), ED())
	assertEncode(t, nil, []byte{header, ceVer, typePosInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, BI(NewBigInt("18446744073709551615", 10)), ED())

	assertEncode(t, nil, []byte{header, ceVer, typeNA, typeNA}, BD(), EvV, BI(nil), NA(), ED())
}

func TestCBENegativeInt(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, 0xff}, BD(), EvV, I(-1), ED())
	assertDecodeEncode(t, []byte{header, ceVer, 0x9c}, BD(), EvV, I(-100), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt8, 101}, BD(), EvV, NI(101), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt8, 0xff}, BD(), EvV, NI(255), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt16, 0x00, 0x01}, BD(), EvV, NI(0x100), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt16, 0xff, 0xff}, BD(), EvV, NI(0xffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt32, 0x00, 0x00, 0x01, 0x00}, BD(), EvV, NI(0x10000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt32, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, NI(0xffffffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt, 0x05, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, NI(0x100000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, NI(0x10000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00}, BD(), EvV, NI(0x1000000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, NI(0x100000000000000), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, NI(0xffffffffffffffff), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, BI(NewBigInt("-18446744073709551616", 10)), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNegInt, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, BD(), EvV, BI(NewBigInt("-4722366482869645213696", 10)), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeNegInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, BD(), EvV, BI(NewBigInt("-9223372036854775807", 10)), ED())
}

func TestCBEBinaryFloatEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeFloat16, 0xd1})
	assertDecodeFails(t, []byte{header, ceVer, typeFloat32, 0xd1, 0x00, 0x00})
	assertDecodeFails(t, []byte{header, ceVer, typeFloat64, 0xd1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestCBEBinaryFloat(t *testing.T) {
	nanBits := math.Float64bits(math.NaN())
	quietNan := math.Float64frombits(nanBits | uint64(1<<50))
	signalingNan := math.Float64frombits(nanBits & ^uint64(1<<50))
	zero := float64(0)
	negZero := -zero
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x80, 0x00}, BD(), EvV, F(quietNan), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x81, 0x00}, BD(), EvV, F(signalingNan), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x82, 0x00}, BD(), EvV, F(math.Inf(1)), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x83, 0x00}, BD(), EvV, F(math.Inf(-1)), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x02}, BD(), EvV, F(zero), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x03}, BD(), EvV, F(negZero), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x80, 0x00}, BD(), EvV, NAN(), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x81, 0x00}, BD(), EvV, SNAN(), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typeFloat16, 0xd1, 0x17}, BD(), EvV, F(0x1.a2p-80), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeFloat32, 0x80, 0xf4, 0xa7, 0x71}, BD(), EvV, F(0x1.4fe9p100), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeFloat64, 0x00, 0x00, 0xc2, 0x99, 0x91, 0xfe, 0xb4, 0x20}, BD(), EvV, F(0x1.4fe9199c2p-500), ED())
}

func TestCBEDecimalFloatEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeDecimal, 0x04})
}

func TestCBEDecimalFloat(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeDecimal, 0x06, 0x0c}, BD(), EvV, DF(NewDFloat("1.2")), ED())
}

func TestCBEBigFloat(t *testing.T) {
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x88, 0x9c, 0x01, 0xa3, 0xbf, 0xc0, 0x04}, BD(), EvV, BF(NewBigFloat("9.445283e+5000", 10, 7)), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typeDecimal,
		0xcf, 0x9d, 0x01, 0xd1, 0x8e, 0xa2, 0xe6, 0x83, 0x8a, 0xbf, 0xc1, 0xbb,
		0xe1, 0xf3, 0xdf, 0xfc, 0xee, 0xac, 0xe5, 0xfe, 0xe1, 0x8f, 0xe2, 0x43},
		BD(), EvV, BDF(NewBDF("-9.4452837206285466345998345667683453466347345e-5000")), ED())

	assertEncode(t, nil, []byte{header, ceVer, typeNA, typeNA}, BD(), EvV, BF(nil), NA(), ED())
}

func TestCBEBigDecimalFloat(t *testing.T) {
	assertEncode(t, nil, []byte{header, ceVer, typeDecimal, 0x88, 0x9c, 0x01, 0xa3, 0xbf, 0xc0, 0x04}, BD(), EvV, BDF(NewBDF("9.445283e+5000")), ED())
	assertEncode(t, nil, []byte{header, ceVer, typeNA, typeNA}, BD(), EvV, BDF(nil), NA(), ED())
}

func TestCBEUUIDEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeUUID, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
}

func TestCBEUUID(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeUUID, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		BD(), EvV, UUID([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}), ED())
}

func TestCBETimeEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeDate, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeTime, 0x01})
	assertDecodeFails(t, []byte{header, ceVer, typeTimestamp, 0x01})
}

func TestCBETime(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeDate, 0x95, 0x7f, 0x3e}, BD(), EvV, CT(NewDate(-2000, 12, 21)), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeTime, 0xfe, 0x4f, 0xd6, 0xdc, 0x8b, 0x14, 0xfd}, BD(), EvV, CT(NewTime(8, 41, 05, 999999999, "")), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeTimestamp, 0x01, 0x00, 0x10, 0x02, 00, 0x10, 'E', '/', 'B', 'e', 'r', 'l', 'i', 'n'}, BD(), EvV, CT(NewTS(2000, 1, 1, 0, 0, 0, 0, "Europe/Berlin")), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeTimestamp, 0x8d, 0x1c, 0xb0, 0xd7, 0x06, 0x1f, 0x99, 0x12, 0xd5, 0x2e, 0x2f, 0x04}, BD(), EvV, CT(NewTSLL(3190, 8, 31, 0, 54, 47, 394129000, 5994, 1071)), ED())

	tz, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	assertEncode(t, nil, []byte{header, ceVer, typeTimestamp, 0x01, 0x00, 0x10, 0x02, 00, 0x10, 'E', '/', 'B', 'e', 'r', 'l', 'i', 'n'}, BD(), EvV, GT(time.Date(2000, 1, 1, 0, 0, 0, 0, tz)), ED())
}

func TestCBEConstant(t *testing.T) {
	assertEncode(t, nil, []byte{header, ceVer, typeString1, 'a'}, BD(), EvV, CONST("x", true), S("a"), ED())
	assertEncodeFails(t, nil, BD(), EvV, CONST("x", false), S("a"), ED())
	assertEncodeFails(t, nil, BD(), EvV, CONST("x", false), ED())
}

func TestCBEShortStringEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeString1})
	assertDecodeFails(t, []byte{header, ceVer, typeString2, 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString3, 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString4, 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString5, 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString6, 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString7, 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString8, 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString9, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString10, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString11, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString12, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString13, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString14, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
	assertDecodeFails(t, []byte{header, ceVer, typeString15, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'})
}

func TestCBEShortString(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeString0}, BD(), EvV, S(""), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString1, 'a'}, BD(), EvV, S("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString2, 'a', 'a'}, BD(), EvV, S("aa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString3, 'a', 'a', 'a'}, BD(), EvV, S("aaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString4, 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString5, 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString6, 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString7, 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString8, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString9, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString10, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString11, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString12, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString13, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString14, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaaaaaa"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString15, 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}, BD(), EvV, S("aaaaaaaaaaaaaaa"), ED())
}

func TestCBEStringEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeString, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeString, 0x04, 'a'})
}

func TestCBEString(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typeString, 0x00}, BD(), EvV, SB(), AC(0, false), ED())
	assertDecode(t, nil, []byte{header, ceVer, typeString, 0x02, 'a'}, BD(), EvV, S("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeString, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, S("00000000001111111111"), ED())
}

func TestCBERIDEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, TypeRID, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, TypeRID, 0x04, 'a'})

	assertDecodeFails(t, []byte{header, ceVer, TypeRIDCat, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, TypeRIDCat, 0x04, 'a'})
}

func TestCBERID(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, TypeRID, 0x00}, BD(), EvV, RB(), AC(0, false), ED())
	assertDecodeEncode(t, []byte{header, ceVer, TypeRID, 0x02, 'a'}, BD(), EvV, RID("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, TypeRID, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, RID("00000000001111111111"), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, TypeRIDCat, 0x00, 1}, BD(), EvV, RBCat(), AC(0, false), I(1), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, TypeRIDCat, 0x02, 'a', 1}, BD(), EvV, RIDCat("a"), I(1), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, TypeRIDCat, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1', 1}, BD(), EvV, RIDCat("00000000001111111111"), I(1), ED())
}

func TestCBECustomBinaryEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeCustomBinary, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeCustomBinary, 0x04, 'a'})
}

func TestCBECustomBinary(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x00}, BD(), EvV, CBB(), AC(0, false), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x02, 'a'}, BD(), EvV, CUB([]byte("a")), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomBinary, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, CUB([]byte("00000000001111111111")), ED())
}

func TestCBECustomTextEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeCustomText, 0x02})
	assertDecodeFails(t, []byte{header, ceVer, typeCustomText, 0x04, 'a'})
}

func TestCBECustomText(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x00}, BD(), EvV, CTB(), AC(0, false), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x02, 'a'}, BD(), EvV, CUT("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeCustomText, 0x28, '0', '0', '0', '0', '0', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1', '1', '1', '1', '1', '1'}, BD(), EvV, CUT("00000000001111111111"), ED())
}

func TestCBEArrayUint8EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt8, 0x04, 0xfa})
}

func TestCBEArrayUint8(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt8, 0x02, 0x01}, BD(), EvV, AU8([]byte{1}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt8, 0x04, 0xfa, 0x11}, BD(), EvV, AU8([]byte{0xfa, 0x11}), ED())
}

func TestCBEArrayUint16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt16, 0x02, 0xfa})
}

func TestCBEArrayUint16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt16, 0x02, 0x01, 0x02}, BD(), EvV, AU16([]uint16{0x0201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt16, 0x04, 0xfa, 0x11, 0x01, 0x02}, BD(), EvV, AU16([]uint16{0x11fa, 0x0201}), ED())
}

func TestCBEArrayUint32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayUint32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x02, 0x01, 0x02, 0x03, 0x04}, BD(), EvV, AU32([]uint32{0x04030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt32, 0x04, 1, 2, 3, 4, 5, 6, 7, 8}, BD(), EvV, AU32([]uint32{0x04030201, 0x08070605}), ED())
}

func TestCBEArrayUint64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayUint64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x02, 1, 2, 3, 4, 5, 6, 7, 8}, BD(), EvV, AU64([]uint64{0x0807060504030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typePosInt64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}, BD(), EvV, AU64([]uint64{0x0807060504030201, 0x0605040302010009}), ED())
}

func TestCBEArrayInt8EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt8, 0x04, 0xfa})
}

func TestCBEArrayInt8(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt8, 0x02, 0x01}, BD(), EvV, AI8([]int8{1}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt8, 0x04, 0xfa, 0x11}, BD(), EvV, AI8([]int8{-6, 0x11}), ED())
}

func TestCBEArrayInt16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt16, 0x02, 0xfa})
}

func TestCBEArrayInt16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt16, 0x02, 0x01, 0x02}, BD(), EvV, AI16([]int16{0x0201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt16, 0x04, 0xfa, 0x11, 0x9c, 0xff}, BD(), EvV, AI16([]int16{0x11fa, -100}), ED())
}

func TestCBEArrayInt32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayInt32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x02, 0x01, 0x02, 0x03, 0x04}, BD(), EvV, AI32([]int32{0x04030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt32, 0x04, 1, 2, 3, 4, 0x9c, 0xff, 0xff, 0xff}, BD(), EvV, AI32([]int32{0x04030201, -100}), ED())
}

func TestCBEArrayInt64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayInt64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x02, 1, 2, 3, 4, 5, 6, 7, 7}, BD(), EvV, AI64([]int64{0x0707060504030201}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeNegInt64, 0x04, 1, 2, 3, 4, 5, 6, 7, 7, 0x9c, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, BD(), EvV, AI64([]int64{0x0707060504030201, -100}), ED())
}

func TestCBEArrayFloat16EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat16, 0x02, 0xfa})
}

func TestCBEArrayFloat16(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat16, 0x02, 0x01, 0x02}, BD(), EvV, AF16([]byte{0x01, 0x02}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat16, 0x04, 0xfa, 0x11, 0x9c, 0xff}, BD(), EvV, AF16([]byte{0xfa, 0x11, 0x9c, 0xff}), ED())
}

func TestCBEArrayFloat32EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x02, 1, 2, 3})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x04, 1, 2, 3, 4, 5})
}

func TestCBEArrayFloat32(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x02, 0x00, 0x58, 0x93, 0x54}, BD(), EvV, AF32([]float32{0x4.9acp40}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat32, 0x04, 0x00, 0xc8, 0xfc, 0xb5, 000, 0xb1, 0x48, 0xd0}, BD(), EvV, AF32([]float32{-0x1.f99p-20, -0xc.8b1p30}), ED())
}

func TestCBEArrayFloat64EOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1, 2})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1, 2, 3, 4})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1, 2, 3, 4, 5})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1, 2, 3, 4, 5, 6})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 1, 2, 3, 4, 5, 6, 7})
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x04, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCBEArrayFloat64(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x02, 0x66, 0x46, 0x74, 0x3e, 0x33, 0x16, 0x09, 0xc2}, BD(), EvV, AF64([]float64{-0xc.8b199f3a2333p30}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeFloat64, 0x04, 0x00, 0x00, 0xcc, 0xea, 0xf1, 0x7f, 0x32, 0xbe, 0x00, 0x10, 0x90, 0xea, 0xfc, 0x87, 0x18, 0xc3}, BD(), EvV, AF64([]float64{-0x4.9ffc7ab3p-30, -0x1.887fcea901p50}), ED())
}

func TestCBEArrayBooleanEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typePlane2, typeTrue, 0x12, 0x00})
}

func TestCBEArrayBoolean(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeTrue, 0x02, 0x01}, BD(), EvV, AB(1, []byte{0x01}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeTrue, 0x08, 0x03}, BD(), EvV, AB(4, []byte{0x03}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typePlane2, typeTrue, 0x26, 0xfe, 0xc1, 0x03}, BD(), EvV, AB(19, []byte{0xfe, 0xc1, 0x03}), ED())
}

func TestCBEArrayUUID(t *testing.T) {
	// TODO: TestCBEArrayUUID
}

func TestCBEMarker(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeMarker, 1, typeString1, 'a'}, BD(), EvV, MARK(), I(1), S("a"), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMarker, typeString1, 'a', typeString4, 't', 'e', 's', 't'}, BD(), EvV, MARK(), S("a"), S("test"), ED())
}

func TestCBEReference(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeReference, 1}, BD(), EvV, REF(), I(1), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeReference, typeString1, 'a'}, BD(), EvV, REF(), S("a"), ED())
}

func TestCBEContainers(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1, typeEndContainer}, BD(), EvV, L(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMap, 1, typeEndContainer}, BD(), EvV, M(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMetadata, 1, typeEndContainer}, BD(), EvV, META(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeComment, 1, typeEndContainer}, BD(), EvV, CMT(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMarkup, 1, typeEndContainer, typeEndContainer}, BD(), EvV, MUP(), I(1), E(), E(), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1,
		typeList, typeString1, 'a', typeEndContainer,
		typeMap, typeString1, 'a', 100, typeEndContainer,
		typeMetadata, typeString1, 'a', 100, typeEndContainer,
		typeComment, typeString1, 'a', typeEndContainer,
		typeMarkup, typeString1, 'a', typeString1, 'a', 50, typeEndContainer, typeString1, 'a', typeEndContainer,
		typeEndContainer,
	},
		BD(), EvV, L(), I(1),
		L(), S("a"), E(),
		M(), S("a"), I(100), E(),
		META(), S("a"), I(100), E(),
		CMT(), S("a"), E(),
		MUP(), S("a"), S("a"), I(50), E(), S("a"), E(),
		E(), ED())
}

func TestCBEMultipartArray(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typeString, 0x03, 'a', 0x02, 'b'}, BD(), EvV, SB(), AC(1, true), AD([]byte{'a'}), AC(1, false), AD([]byte{'b'}), ED())
	assertDecode(t, nil, []byte{header, ceVer, typePlane2, typePosInt16, 0x03, 0x01, 0x02, 0x02, 0x03, 0x04}, BD(), EvV, AU16B(), AC(1, true), AD([]byte{0x01, 0x02}), AC(1, false), AD([]byte{0x03, 0x04}), ED())
}

func TestCBEChunkedArray(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	InvokeEvents(encoder, EvV, AU16B(), AC(2, true), AD([]byte{1, 0, 2, 0}), AC(2, false), AD([]byte{3, 0, 4, 0}), ED())

	expected := []byte{header, ceVer, typePlane2, typePosInt16, 0x05, 1, 0, 2, 0, 0x04, 3, 0, 4, 0}
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
