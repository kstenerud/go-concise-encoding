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

package cte

import (
	"bytes"
	"math"
	"testing"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-compact-time"
)

func TestCTEVersion(t *testing.T) {
	// Valid
	assertDecodeEncode(t, nil, nil, "c1 ", BD(), V(1), ED())
	assertDecode(t, nil, "\r\n\t c1 ", BD(), V(1), ED())
	assertDecode(t, nil, "c1     \r\n\t\t\t", BD(), V(1), ED())

	// Missing whitespace
	assertDecodeFails(t, "c1")
	assertDecodeFails(t, "c1{}")

	// Too big
	assertDecodeFails(t, "c100000000000000000000000000000000 ")

	// Not numeric
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		document := string([]byte{'c', byte(i)})
		assertDecodeFails(t, document)
	}

	// Disallowed version numbers
	for i := 0; i < 0x100; i++ {
		switch i {
		case 'c', 'C', ' ', '\n', '\r', '\t':
			continue
		default:
			document := string([]byte{byte(i)})
			assertDecodeFails(t, document)
		}
	}
}

func TestCTENil(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 @nil", BD(), V(1), N(), ED())
	assertDecodeFails(t, "c1 @nill")
	assertDecodeFails(t, "c1 -@nil")
}

func TestCTEBool(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 @true", BD(), V(1), TT(), ED())
	assertDecodeEncode(t, nil, nil, "c1 @false", BD(), V(1), FF(), ED())

	assertEncode(t, nil, "c1 @false", BD(), V(1), B(false), ED())
	assertEncode(t, nil, "c1 @true", BD(), V(1), B(true), ED())

	assertDecodeFails(t, "c1 @truer")
	assertDecodeFails(t, "c1 @falser")
	assertDecodeFails(t, "c1 -@true")
	assertDecodeFails(t, "c1 -@false")
}

func TestCTEDecimalInt(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 0", BD(), V(1), PI(0), ED())
	assertDecodeEncode(t, nil, nil, "c1 123", BD(), V(1), PI(123), ED())
	assertDecodeEncode(t, nil, nil, "c1 9412504234235366", BD(), V(1), PI(9412504234235366), ED())
	assertDecodeEncode(t, nil, nil, "c1 -49523", BD(), V(1), NI(49523), ED())
	assertDecodeEncode(t, nil, nil, "c1 10000000000000000000000000000", BD(), V(1), BI(NewBigInt("10000000000000000000000000000", 10)), ED())
	assertDecodeEncode(t, nil, nil, "c1 -10000000000000000000000000000", BD(), V(1), BI(NewBigInt("-10000000000000000000000000000", 10)), ED())
	assertDecode(t, nil, "c1 100_00_0_00000000000_00000000_000_0", BD(), V(1), BI(NewBigInt("10000000000000000000000000000", 10)), ED())
	assertDecode(t, nil, "c1 -4_9_5__2___3", BD(), V(1), NI(49523), ED())
	assertEncode(t, nil, "c1 100", BD(), V(1), I(100), ED())
	assertEncode(t, nil, "c1 -100", BD(), V(1), I(-100), ED())
	assertDecode(t, nil, "c1 100", BD(), V(1), PI(100), ED())
	assertDecode(t, nil, "c1 -100", BD(), V(1), NI(100), ED())

	assertDecodeFails(t, "c1 1f")
	assertDecodeFails(t, "c1 -1f")
}

func TestCTEBinaryInt(t *testing.T) {
	assertDecode(t, nil, "c1 0b0", BD(), V(1), PI(0), ED())
	assertDecode(t, nil, "c1 0b1", BD(), V(1), PI(1), ED())
	assertDecode(t, nil, "c1 0b101", BD(), V(1), PI(5), ED())
	assertDecode(t, nil, "c1 0b0010100", BD(), V(1), PI(20), ED())
	assertDecode(t, nil, "c1 -0b100", BD(), V(1), NI(4), ED())
	assertDecode(t, nil, "c1 -0b_1_0_0", BD(), V(1), NI(4), ED())

	assertDecode(t, nil, "c1 0b10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 2)), ED())
	assertDecode(t, nil, "c1 -0b10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("-10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 2)), ED())
	assertDecode(t, nil, "c1 0b100000000000000_000000000000_000000000000000000000000_000000000000000000000_0000000000000000000000000000000000000000_0",
		BD(), V(1), BI(NewBigInt("10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 2)), ED())

	assertDecodeFails(t, "c1 0b2")
}

func TestCTEOctalInt(t *testing.T) {
	assertDecode(t, nil, "c1 0o0", BD(), V(1), PI(0), ED())
	assertDecode(t, nil, "c1 0o1", BD(), V(1), PI(1), ED())
	assertDecode(t, nil, "c1 0o7", BD(), V(1), PI(7), ED())
	assertDecode(t, nil, "c1 0o71", BD(), V(1), PI(57), ED())
	assertDecode(t, nil, "c1 0o644", BD(), V(1), PI(420), ED())
	assertDecode(t, nil, "c1 -0o777", BD(), V(1), NI(511), ED())
	assertDecode(t, nil, "c1 -0o_7__7___7", BD(), V(1), NI(511), ED())

	assertDecode(t, nil, "c1 0o1000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("1000000000000000000000000000000000000000000000", 8)), ED())
	assertDecode(t, nil, "c1 -0o1000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("-1000000000000000000000000000000000000000000000", 8)), ED())
	assertDecode(t, nil, "c1 0o1_0000000000000_0000000000000_000000000000000000_0",
		BD(), V(1), BI(NewBigInt("1000000000000000000000000000000000000000000000", 8)), ED())

	assertDecodeFails(t, "c1 0o9")
}

func TestCTEHexInt(t *testing.T) {
	assertDecode(t, nil, "c1 0x0", BD(), V(1), PI(0), ED())
	assertDecode(t, nil, "c1 0x1", BD(), V(1), PI(1), ED())
	assertDecode(t, nil, "c1 0xf", BD(), V(1), PI(0xf), ED())
	assertDecode(t, nil, "c1 0xfedcba9876543210", BD(), V(1), PI(0xfedcba9876543210), ED())
	assertDecode(t, nil, "c1 0xFEDCBA9876543210", BD(), V(1), PI(0xfedcba9876543210), ED())
	assertDecode(t, nil, "c1 -0x88", BD(), V(1), NI(0x88), ED())
	assertDecode(t, nil, "c1 -0x_8_8__5_a_f__d", BD(), V(1), NI(0x885afd), ED())

	assertDecode(t, nil, "c1 0x1000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("1000000000000000000000000000000000000000000000", 16)), ED())
	assertDecode(t, nil, "c1 -0x1000000000000000000000000000000000000000000000",
		BD(), V(1), BI(NewBigInt("-1000000000000000000000000000000000000000000000", 16)), ED())
	assertDecode(t, nil, "c1 0x1_00000000000_00000000000_0000000000000000000000_0",
		BD(), V(1), BI(NewBigInt("1000000000000000000000000000000000000000000000", 16)), ED())

	assertDecodeFails(t, "c1 0xg")
}

func TestCTEFloat(t *testing.T) {
	assertDecode(t, nil, "c1 0.0", BD(), V(1), DF(NewDFloat("0")), ED())
	assertDecode(t, nil, "c1 -0.0", BD(), V(1), DF(NewDFloat("-0")), ED())

	assertDecodeEncode(t, nil, nil, "c1 1.5", BD(), V(1), DF(NewDFloat("1.5")), ED())
	assertDecodeEncode(t, nil, nil, "c1 1.125", BD(), V(1), DF(NewDFloat("1.125")), ED())
	assertDecodeEncode(t, nil, nil, "c1 1.125e+10", BD(), V(1), DF(NewDFloat("1.125e+10")), ED())
	assertDecodeEncode(t, nil, nil, "c1 1.125e-10", BD(), V(1), DF(NewDFloat("1.125e-10")), ED())
	assertDecode(t, nil, "c1 1.125e10", BD(), V(1), DF(NewDFloat("1.125e+10")), ED())

	assertDecodeEncode(t, nil, nil, "c1 -1.5", BD(), V(1), DF(NewDFloat("-1.5")), ED())
	assertDecodeEncode(t, nil, nil, "c1 -1.125", BD(), V(1), DF(NewDFloat("-1.125")), ED())
	assertDecodeEncode(t, nil, nil, "c1 -1.125e+10", BD(), V(1), DF(NewDFloat("-1.125e+10")), ED())
	assertDecodeEncode(t, nil, nil, "c1 -1.125e-10", BD(), V(1), DF(NewDFloat("-1.125e-10")), ED())
	assertDecode(t, nil, "c1 -1.125e10", BD(), V(1), DF(NewDFloat("-1.125e10")), ED())
	assertDecodeEncode(t, nil, nil, "c1 1.0000000000000000001", BD(), V(1), BDF(NewBDF("1.0000000000000000001")), ED())
	assertDecodeEncode(t, nil, nil, "c1 -1.0000000000000000001", BD(), V(1), BDF(NewBDF("-1.0000000000000000001")), ED())

	assertDecodeEncode(t, nil, nil, "c1 0.5", BD(), V(1), DF(NewDFloat("0.5")), ED())
	assertDecodeEncode(t, nil, nil, "c1 0.125", BD(), V(1), DF(NewDFloat("0.125")), ED())
	assertDecode(t, nil, "c1 0.125e+10", BD(), V(1), DF(NewDFloat("0.125e+10")), ED())
	assertDecode(t, nil, "c1 0.125e-10", BD(), V(1), DF(NewDFloat("0.125e-10")), ED())
	assertDecode(t, nil, "c1 0.125e10", BD(), V(1), DF(NewDFloat("0.125e10")), ED())

	assertDecode(t, nil, "c1 -0.5", BD(), V(1), DF(NewDFloat("-0.5")), ED())
	assertDecode(t, nil, "c1 -0.125", BD(), V(1), DF(NewDFloat("-0.125")), ED())
	assertDecode(t, nil, "c1 -0.125e+10", BD(), V(1), DF(NewDFloat("-0.125e+10")), ED())
	assertDecode(t, nil, "c1 -0.125e-10", BD(), V(1), DF(NewDFloat("-0.125e-10")), ED())
	assertDecode(t, nil, "c1 -0.125e10", BD(), V(1), DF(NewDFloat("-0.125e10")), ED())
	assertDecode(t, nil, "c1 -0.125E+10", BD(), V(1), DF(NewDFloat("-0.125e+10")), ED())
	assertDecode(t, nil, "c1 -0.125E-10", BD(), V(1), DF(NewDFloat("-0.125e-10")), ED())
	assertDecode(t, nil, "c1 -0.125E10", BD(), V(1), DF(NewDFloat("-0.125e10")), ED())

	assertDecode(t, nil, "c1 -1.50000000000000000000000001E10000", BD(), V(1), BDF(NewBDF("-1.50000000000000000000000001E10000")), ED())
	assertDecode(t, nil, "c1 1.50000000000000000000000001E10000", BD(), V(1), BDF(NewBDF("1.50000000000000000000000001E10000")), ED())

	assertDecode(t, nil, "c1 1_._1_2_5_e+1_0", BD(), V(1), DF(NewDFloat("1.125e+10")), ED())

	assertDecode(t, nil, "c1 0.1500000000000000000000000000000000000000000000000001e+10000",
		BD(), V(1), BDF(NewBDF("0.1500000000000000000000000000000000000000000000000001e+10000")), ED())
	assertDecode(t, nil, "c1 -0.1500000000000000000000000000000000000000000000000001e+10000",
		BD(), V(1), BDF(NewBDF("-0.1500000000000000000000000000000000000000000000000001e+10000")), ED())
	assertDecode(t, nil, "c1 0.1_50000000000_00000000000_000000000000_0000000000000000_1e+100_0_0",
		BD(), V(1), BDF(NewBDF("0.1500000000000000000000000000000000000000000000000001e+10000")), ED())

	assertEncode(t, nil, "c1 @nan", BD(), V(1), F(common.QuietNan), ED())
	assertEncode(t, nil, "c1 @snan", BD(), V(1), F(common.SignalingNan), ED())

	assertEncode(t, nil, "c1 1.1", BD(), V(1), BF(NewBigFloat("1.1", 10, 2)), ED())

	assertDecodeFails(t, "c1 -0.5.4")
	assertDecodeFails(t, "c1 -0,5.4")
	assertDecodeFails(t, "c1 0.5.4")
	assertDecodeFails(t, "c1 0,5.4")
	assertDecodeFails(t, "c1 -@blah")
	assertDecodeFails(t, "c1 1.1.1")
	assertDecodeFails(t, "c1 1,1")
	assertDecodeFails(t, "c1 1.1e4e5")
	assertDecodeFails(t, "c1 0.a")
	assertDecodeFails(t, "c1 0.5et")
	assertDecodeFails(t, "c1 0.5e99999999999999999999999")
}

func TestCTEHexFloat(t *testing.T) {
	assertDecode(t, nil, "c1 0x0.0", BD(), V(1), F(0x0.0p0), ED())
	assertDecode(t, nil, "c1 0x0.1", BD(), V(1), F(0x0.1p0), ED())
	assertDecode(t, nil, "c1 0x0.1p+10", BD(), V(1), F(0x0.1p+10), ED())
	assertDecode(t, nil, "c1 0x0.1p-10", BD(), V(1), F(0x0.1p-10), ED())
	assertDecode(t, nil, "c1 0x0.1p10", BD(), V(1), F(0x0.1p10), ED())

	assertDecode(t, nil, "c1 0x1.0", BD(), V(1), F(0x1.0p0), ED())
	assertDecode(t, nil, "c1 0x1.1", BD(), V(1), F(0x1.1p0), ED())
	assertDecode(t, nil, "c1 0xf.1p+10", BD(), V(1), F(0xf.1p+10), ED())
	assertDecode(t, nil, "c1 0xf.1p-10", BD(), V(1), F(0xf.1p-10), ED())
	assertDecode(t, nil, "c1 0xf.1p10", BD(), V(1), F(0xf.1p10), ED())

	assertDecode(t, nil, "c1 -0x1.0", BD(), V(1), F(-0x1.0p0), ED())
	assertDecode(t, nil, "c1 -0x1.1", BD(), V(1), F(-0x1.1p0), ED())
	assertDecode(t, nil, "c1 -0xf.1p+10", BD(), V(1), F(-0xf.1p+10), ED())
	assertDecode(t, nil, "c1 -0xf.1p-10", BD(), V(1), F(-0xf.1p-10), ED())
	assertDecode(t, nil, "c1 -0xf.1p10", BD(), V(1), F(-0xf.1p10), ED())

	assertDecode(t, nil, "c1 -0x0.0", BD(), V(1), F(-0x0.0p0), ED())
	assertDecode(t, nil, "c1 -0x0.1", BD(), V(1), F(-0x0.1p0), ED())
	assertDecode(t, nil, "c1 -0x0.1p+10", BD(), V(1), F(-0x0.1p+10), ED())
	assertDecode(t, nil, "c1 -0x0.1p-10", BD(), V(1), F(-0x0.1p-10), ED())
	assertDecode(t, nil, "c1 -0x0.1p10", BD(), V(1), F(-0x0.1p10), ED())

	// Everything too big for float64
	bigExpected := NewBigFloat("-1.54fffe2ac00592375b427ap100000", 16, 22)
	assertDecode(t, nil, "c1 -0x1.54fffe2ac00592375b427ap100000", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.54fffe2ac00592375b427ap100000", BD(), V(1), BF(bigExpected), ED())

	// Coefficient too big for float64
	bigExpected = NewBigFloat("-1.54fffe2ac00592375b427ap100", 16, 22)
	assertDecode(t, nil, "c1 -0x1.54fffe2ac00592375b427ap100", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.54fffe2ac00592375b427ap100", BD(), V(1), BF(bigExpected), ED())
	assertDecode(t, nil, "c1 0x1.5_4fffe2ac_0059237_5b42_7ap1_00", BD(), V(1), BF(bigExpected), ED())
	assertDecode(t, nil, "c1 0x1.5_4FFFE2AC_0059237_5B42_7AP1_00", BD(), V(1), BF(bigExpected), ED())

	// Exponent too big for float64
	bigExpected = NewBigFloat("-1.8p100000", 16, 16)
	assertDecode(t, nil, "c1 -0x1.8p100000", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.8p100000", BD(), V(1), BF(bigExpected), ED())

	assertDecode(t, nil, "c1 -0x_0_._1_p_1_0", BD(), V(1), F(-0x0.1p10), ED())

	bigExpected = NewBigFloat("8.000000000000001p100", 16, 16)
	assertDecode(t, nil, "c1 0x8.000000000000001p100", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 -0x8.000000000000001p100", BD(), V(1), BF(bigExpected), ED())

	assertDecodeFails(t, "c1 -0x0.5.4")
	assertDecodeFails(t, "c1 -0x0,5.4")
	assertDecodeFails(t, "c1 0x0.5.4")
	assertDecodeFails(t, "c1 0x0,5.4")
	assertDecodeFails(t, "c1 -0x@blah")
	assertDecodeFails(t, "c1 0x1.1.1")
	assertDecodeFails(t, "c1 0x1,1")
	assertDecodeFails(t, "c1 0x1.1p4p5")
	assertDecodeFails(t, "c1 -0x0.l")
	assertDecodeFails(t, "c1 -0x0.5pj")
	assertDecodeFails(t, "c1 -0x0.5p1000000000000000000000000000")
}

func TestCTEUUID(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbb`, BD(), V(1),
		UUID([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertDecode(t, nil, `c1 @FEDCBA98-7654-3210-AAAA-BBBBBBBBBBBB`, BD(), V(1),
		UUID([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertDecodeEncode(t, nil, nil, `c1 @00000000-0000-0000-0000-000000000000`, BD(), V(1),
		UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), ED())

	assertDecodeFails(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbb`)
	assertDecodeFails(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbbb`)
	assertEncodeFails(t, nil, BD(), V(1), UUID([]byte{0xfe, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertEncodeFails(t, nil, BD(), V(1), UUID([]byte{0xfe, 0xdc, 0xff, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
}

func TestCTEDate(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 2000-01-01", BD(), V(1), CT(compact_time.NewDate(2000, 1, 1)), ED())
	assertDecodeEncode(t, nil, nil, "c1 -2000-12-31", BD(), V(1), CT(compact_time.NewDate(-2000, 12, 31)), ED())

	assertDecodeFails(t, "c1 0-01-01")
	assertDecodeFails(t, "c1 --2000-01-01")
	assertDecodeFails(t, "c1 0a-01-01")

	assertDecodeFails(t, "c1 2000-013-01")
	assertDecodeFails(t, "c1 2000-30-1")
	assertDecodeFails(t, "c1 2000-0-10")
	assertDecodeFails(t, "c1 2000-1a-10")
	assertDecodeFails(t, "c1 2000-0a-10")
	assertDecodeFails(t, "c1 2000-a-10")

	assertDecodeFails(t, "c1 2000-01-011")
	assertDecodeFails(t, "c1 2000-01-99")
	assertDecodeFails(t, "c1 2000-10-0")
	assertDecodeFails(t, "c1 2000-10-1a")
	assertDecodeFails(t, "c1 2000-10-0a")
	assertDecodeFails(t, "c1 2000-10-a")
}

func TestCTETime(t *testing.T) {
	assertDecode(t, nil, "c1 1:45:00", BD(), V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecode(t, nil, "c1 01:45:00", BD(), V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecodeEncode(t, nil, nil, "c1 23:59:59.101", BD(), V(1), CT(compact_time.NewTime(23, 59, 59, 101000000, "")), ED())
	assertDecodeEncode(t, nil, nil, "c1 10:00:01.93/America/Los_Angeles", BD(), V(1), CT(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertDecodeEncode(t, nil, nil, "c1 10:00:01.93/89.92/1.10", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/89.90/1.1", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8990, 110)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/89.9/1.10", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8990, 110)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/0/0", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/1/1", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ED())

	assertDecodeFails(t, "c1 001:45:00")
	assertDecodeFails(t, "c1 30:45:10")
	assertDecodeFails(t, "c1 1a:45:10")
	assertDecodeFails(t, "c1 0a:45:10")

	assertDecodeFails(t, "c1 1:045:00")
	assertDecodeFails(t, "c1 1:99:10")
	assertDecodeFails(t, "c1 1:1a:10")
	assertDecodeFails(t, "c1 1:0a:10")
	assertDecodeFails(t, "c1 1:a:10")

	assertDecodeFails(t, "c1 1:45:001")
	assertDecodeFails(t, "c1 1:45:99")
	assertDecodeFails(t, "c1 1:45:1e")
	assertDecodeFails(t, "c1 1:45:0e")
	assertDecodeFails(t, "c1 1:45:e")

	assertDecodeFails(t, "c1 1:45:00.3012544133")
	assertDecodeFails(t, "c1 1:45:00.301254f")
	assertDecodeFails(t, "c1 1:45:00.1f")
	assertDecodeFails(t, "c1 1:45:00.0f")
	assertDecodeFails(t, "c1 1:45:00.f")

	assertDecodeFails(t, "c1 10:00:01.93/89.92/1.10a")
	assertDecodeFails(t, "c1 10:00:01.93/89.92a/1.10")
	assertDecodeFails(t, "c1 10:00:01.93/89.92/a.10")
	assertDecodeFails(t, "c1 10:00:01.93/w89.92/1.10")
	assertDecodeFails(t, "c1 10:00:01.93/89.92/1.10a")
	assertDecodeFails(t, "c1 10:00:01.93/89.92/")
	assertDecodeFails(t, "c1 10:00:01.93/89.92/1.101")
	assertDecodeFails(t, "c1 10:00:01.93//1.10")
	assertDecodeFails(t, "c1 10:00:01.93/89.925/1.10")
}

func TestCTETimestamp(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 2000-01-01/19:31:44.901554/Z", BD(), V(1), CT(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "Z")), ED())
	assertDecodeEncode(t, nil, nil, "c1 2020-01-15/13:41:00.000599", BD(), V(1), CT(compact_time.NewTimestamp(2020, 1, 15, 13, 41, 0, 599000, "")), ED())
	assertDecode(t, nil, "c1 2020-01-15/13:41:00.000599", BD(), V(1), CT(compact_time.NewTimestamp(2020, 1, 15, 13, 41, 0, 599000, "")), ED())
	assertDecodeEncode(t, nil, nil, "c1 2020-01-15/10:00:01.93/89.92/1.10", BD(), V(1), CT(compact_time.NewTimestampLatLong(2020, 1, 15, 10, 0, 1, 930000000, 8992, 110)), ED())

	assertDecodeFails(t, "c1 0-01-01/19:31:44.901554")
	assertDecodeFails(t, "c1 1a-01-01/19:31:44.901554")

	assertDecodeFails(t, "c1 2000-0-01/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-001-01/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-100-01/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-1a-01/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-0a-01/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-a-01/19:31:44.901554")

	assertDecodeFails(t, "c1 2000-01-0/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-001/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-100/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-1a/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-0a/19:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-a/19:31:44.901554")

	assertDecodeFails(t, "c1 2000-01-01/019:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/25:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/1a:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/0a:31:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/a:31:44.901554")

	assertDecodeFails(t, "c1 2000-01-01/19:031:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:310:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:1a:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:0a:44.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:a:44.901554")

	assertDecodeFails(t, "c1 2000-01-01/19:31:044.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:31:440.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:31:1a.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:31:0a.901554")
	assertDecodeFails(t, "c1 2000-01-01/19:31:a.901554")

	assertDecodeFails(t, "c1 2000-01-01/19:31:44.9015544348")
	assertDecodeFails(t, "c1 2000-01-01/19:31:44.1a")
	assertDecodeFails(t, "c1 2000-01-01/19:31:44.0a")
	assertDecodeFails(t, "c1 2000-01-01/19:31:44.a")

	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/1.10a")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92a/1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/a1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/a89.92/1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/a")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/a/1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93//1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/1999.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/8965.92/1.10")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.92/1.a")
	assertDecodeFails(t, "c1 2020-01-15/10:00:01.93/89.a/1.10")

	gotime, err := NewTS(2020, 1, 15, 10, 0, 1, 930000000, "").AsGoTime()
	if err != nil {
		panic(err)
	}
	assertEncode(t, nil, "c1 2020-01-15/10:00:01.93", BD(), V(1), GT(gotime), ED())
}

func TestCTEQuotedString(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 "test string"`, BD(), V(1), S("test string"), ED())
	assertDecode(t, nil, `c1 "test\nstring"`, BD(), V(1), S("test\nstring"), ED())
	assertDecode(t, nil, `c1 "test\rstring"`, BD(), V(1), S("test\rstring"), ED())
	assertDecode(t, nil, `c1 "test\tstring"`, BD(), V(1), S("test\tstring"), ED())
	assertDecodeEncode(t, nil, nil, `c1 "test\"string"`, BD(), V(1), S("test\"string"), ED())
	assertDecode(t, nil, `c1 "test\*string"`, BD(), V(1), S("test*string"), ED())
	assertDecode(t, nil, `c1 "test\/string"`, BD(), V(1), S("test/string"), ED())
	assertDecodeEncode(t, nil, nil, `c1 "test\\string"`, BD(), V(1), S("test\\string"), ED())
	assertDecodeEncode(t, nil, nil, `c1 "test\11string"`, BD(), V(1), S("test\u0001string"), ED())
	assertDecodeEncode(t, nil, nil, `c1 "test\4206dstring"`, BD(), V(1), S("test\u206dstring"), ED())
	assertDecode(t, nil, `c1 "test\4206Dstring"`, BD(), V(1), S("test\u206dstring"), ED())
	assertDecode(t, nil, `c1 "test\
string"`, BD(), V(1), S("teststring"), ED())
	assertDecode(t, nil, "c1 \"test\\\r\nstring\"", BD(), V(1), S("teststring"), ED())

	assertDecodeFails(t, `c1 "test\x"`)
	assertDecodeFails(t, `c1 "\1g"`)
}

func TestCTECustomBinary(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 b"12345678"`, BD(), V(1), CUB([]byte{0x12, 0x34, 0x56, 0x78}), ED())
	assertDecodeEncode(t, nil, nil, `c1 b"abcd"`, BD(), V(1), CUB([]byte{0xab, 0xcd}), ED())
	assertDecode(t, nil, `c1 b"ABCD"`, BD(), V(1), CUB([]byte{0xab, 0xcd}), ED())
	assertDecodeFails(t, `c1 b"qwer"`)
}

func TestCTECustomText(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 t"something(123)"`, BD(), V(1), CUT("something(123)"), ED())
	assertDecodeEncode(t, nil, nil, `c1 t"some\\thing(\"123\")"`, BD(), V(1), CUT("some\\thing(\"123\")"), ED())
	assertDecodeEncode(t, nil, nil, `c1 t"some\nthing\11(123)"`, BD(), V(1), CUT("some\nthing\u0001(123)"), ED())
	assertDecodeEncode(t, nil, nil, `c1 t"something('123\r\n\t')"`, BD(), V(1), CUT("something('123\r\n\t')"), ED())

	assertDecodeFails(t, `c1 t"something('123\r\n\t\x')"`)
}

func TestCTEUnquotedString(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c1 a", BD(), V(1), S("a"), ED())
	assertDecodeEncode(t, nil, nil, "c1 abcd", BD(), V(1), S("abcd"), ED())
	assertDecodeEncode(t, nil, nil, "c1 _-.123aF", BD(), V(1), S("_-.123aF"), ED())
	assertDecodeEncode(t, nil, nil, "c1 新しい", BD(), V(1), S("新しい"), ED())
}

func TestCTEInvalidString(t *testing.T) {
	assertDecodeFails(t, "c1 a|b")
	assertDecodeFails(t, "c1 a*b")
}

func TestCTEVerbatimString(t *testing.T) {
	assertDecodeFails(t, "c1 `")
	assertDecodeFails(t, "c1 `A")
	assertDecodeFails(t, "c1 `A ")
	assertDecodeFails(t, "c1 `A xyz")
	assertDecodeFails(t, "c1 `A xyzAx")
	assertDecode(t, nil, "c1 `A \n\n\n\n\n\n\n\n\n\nA", BD(), V(1), VS("\n\n\n\n\n\n\n\n\n\n"), ED())
	assertDecode(t, nil, "c1 `A aA", BD(), V(1), VS("a"), ED())
	assertDecode(t, nil, "c1 `A\taA", BD(), V(1), VS("a"), ED())
	assertDecode(t, nil, "c1 `A\naA", BD(), V(1), VS("a"), ED())
	assertDecode(t, nil, "c1 `A\r\naA", BD(), V(1), VS("a"), ED())
	assertDecode(t, nil, "c1 `#ENDOFSTRING a test\nwith `stuff`#ENDOFSTRING ", BD(), V(1), VS("a test\nwith `stuff`"), ED())
}

func TestCTEURI(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 u"http://example.com"`, BD(), V(1), URI("http://example.com"), ED())
	assertEncode(t, nil, `c1 u"http://x.com/%22quoted%22"`, BD(), V(1), URI(`http://x.com/"quoted"`), ED())
}

func TestCTEArrayUintX(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 |u8x f1 93|`, BD(), V(1), AU8([]byte{0xf1, 0x93}), ED())
	assertDecode(t, nil, `c1 |u8x f 93 |`, BD(), V(1), AU8([]byte{0xf, 0x93}), ED())
	assertDecodeFails(t, `c1 |u8x f14 93|`)
	assertDecodeFails(t, `c1 |u8x f1o 93|`)

	assertDecodeEncode(t, nil, nil, `c1 |u16x f122 9385|`, BD(), V(1), AU16([]uint16{0xf122, 0x9385}), ED())
	assertDecode(t, nil, `c1 |u16x f12 95|`, BD(), V(1), AU16([]uint16{0xf12, 0x95}), ED())
	assertDecodeFails(t, `c1 |u16x f129e 95|`)
	assertDecodeFails(t, `c1 |u16x f12j 95|`)

	assertDecodeEncode(t, nil, nil, `c1 |u32x 7ddf8134 93cd7aac|`, BD(), V(1), AU32([]uint32{0x7ddf8134, 0x93cd7aac}), ED())
	assertDecode(t, nil, `c1 |u32x 7ddf834 93aac|`, BD(), V(1), AU32([]uint32{0x7ddf834, 0x93aac}), ED())
	assertDecodeFails(t, `c1 |u32x 7ddf8134e 93cd7aac|`)
	assertDecodeFails(t, `c1 |u32x 7ddf81x 93cd7aac|`)

	assertDecodeEncode(t, nil, nil, `c1 |u64x 83ff9ac2445aace7 94ff7ac3219465c1|`, BD(), V(1), AU64([]uint64{0x83ff9ac2445aace7, 0x94ff7ac3219465c1}), ED())
	assertDecode(t, nil, `c1 |u64x 83ff9ac245aace7 94ff79465c1|`, BD(), V(1), AU64([]uint64{0x83ff9ac245aace7, 0x94ff79465c1}), ED())
	assertDecodeFails(t, `c1 |u64x 83ff9ac2445aace72 94ff7ac3219465c1|`)
	assertDecodeFails(t, `c1 |u64x 83ff9ac2l 94ff7ac3219465c1|`)
}

func TestCTEArrayUintB(t *testing.T) {
	// TODO: TestCTEArrayUintB
}

func TestCTEArrayUintO(t *testing.T) {
	// TODO: TestCTEArrayUintO
}

func TestCTEArrayInt8(t *testing.T) {
	eOpts := options.DefaultCTEEncoderOptions()

	eOpts.DefaultArrayEncodingBases.Int8 = 2
	assertDecodeEncode(t, nil, eOpts, `c1 |i8b 0 1 -10 101 1111111 -10000000|`, BD(), V(1), AI8([]int8{0, 1, -2, 5, 0x7f, -0x80}), ED())

	eOpts.DefaultArrayEncodingBases.Int8 = 8
	assertDecodeEncode(t, nil, eOpts, `c1 |i8o 0 -10 50 -127|`, BD(), V(1), AI8([]int8{0o0, -0o10, 0o50, -0o127}), ED())

	eOpts.DefaultArrayEncodingBases.Int8 = 10
	assertDecodeEncode(t, nil, eOpts, `c1 |i8 0 10 -50 127 -128|`, BD(), V(1), AI8([]int8{0, 10, -50, 127, -128}), ED())

	eOpts.DefaultArrayEncodingBases.Int8 = 16
	assertDecodeEncode(t, nil, eOpts, `c1 |i8x 0 1 -50 7f -80|`, BD(), V(1), AI8([]int8{0x00, 0x01, -0x50, 0x7f, -0x80}), ED())

	assertDecode(t, nil, `c1 |i8 00 01 -01 0b101 -0b110 0B101 -0B110 0o10 -0o11 0O10 -0O11 0x7f -0x80 0X7f -0X80|`,
		BD(), V(1), AI8([]int8{0, 1, -1, 5, -6, 5, -6, 8, -9, 8, -9, 127, -128, 127, -128}), ED())

	assertDecodeFails(t, "c1 |i8b 10000000|")
	assertDecodeFails(t, "c1 |i8b -10000001|")
	assertDecodeFails(t, "c1 |i8o 178|")
	assertDecodeFails(t, "c1 |i8o -179|")
	assertDecodeFails(t, "c1 |i8 128|")
	assertDecodeFails(t, "c1 |i8 -129|")
	assertDecodeFails(t, "c1 |i8x 80|")
	assertDecodeFails(t, "c1 |i8x -81|")
}

func TestCTEArrayUint8(t *testing.T) {
	eOpts := options.DefaultCTEEncoderOptions()

	eOpts.DefaultArrayEncodingBases.Uint8 = 2
	assertDecodeEncode(t, nil, eOpts, `c1 |u8b 0 1 10 101 1111111 10000000 11111111|`, BD(), V(1), AU8([]uint8{0, 1, 2, 5, 0x7f, 0x80, 0xff}), ED())

	eOpts.DefaultArrayEncodingBases.Uint8 = 8
	assertDecodeEncode(t, nil, eOpts, `c1 |u8o 0 10 50 127 254 377|`, BD(), V(1), AU8([]uint8{0o0, 0o10, 0o50, 0o127, 0o254, 0o377}), ED())

	eOpts.DefaultArrayEncodingBases.Uint8 = 10
	assertDecodeEncode(t, nil, eOpts, `c1 |u8 0 10 50 128 254 255|`, BD(), V(1), AU8([]uint8{0, 10, 50, 128, 254, 255}), ED())

	eOpts.DefaultArrayEncodingBases.Uint8 = 16
	assertDecodeEncode(t, nil, eOpts, `c1 |u8x 00 01 50 7f 80 ff|`, BD(), V(1), AU8([]uint8{0x00, 0x01, 0x50, 0x7f, 0x80, 0xff}), ED())

	assertDecode(t, nil, `c1 |u8 00 01 01 0b101 0b110 0B101 0B110 0o10 0o11 0O10 0O11 0x7f 0x80 0X7f 0X80 0xff 0Xff|`,
		BD(), V(1), AU8([]uint8{0, 1, 1, 5, 6, 5, 6, 8, 9, 8, 9, 127, 128, 127, 128, 255, 255}), ED())

	assertDecodeFails(t, "c1 |u8b 100000000|")
	assertDecodeFails(t, "c1 |u8o 400|")
	assertDecodeFails(t, "c1 |u8 256|")
	assertDecodeFails(t, "c1 |u8x 100|")
}

func TestCTEArrayInt16(t *testing.T) {
	// TODO: TestCTEArrayInt16
}

func TestCTEArrayUint16(t *testing.T) {
	eOpts := options.DefaultCTEEncoderOptions()

	eOpts.DefaultArrayEncodingBases.Uint16 = 2
	assertDecodeEncode(t, nil, eOpts, `c1 |u16b 0 1 10 101 111111111111111 1000000000000000 1111111111111111|`,
		BD(), V(1), AU16([]uint16{0, 1, 2, 5, 0x7fff, 0x8000, 0xffff}), ED())

	eOpts.DefaultArrayEncodingBases.Uint16 = 8
	assertDecodeEncode(t, nil, eOpts, `c1 |u16o 0 10 50 127 254 377 177777|`,
		BD(), V(1), AU16([]uint16{0o0, 0o10, 0o50, 0o127, 0o254, 0o377, 0o177777}), ED())

	eOpts.DefaultArrayEncodingBases.Uint16 = 10
	assertDecodeEncode(t, nil, eOpts, `c1 |u16 0 10 50 128 254 255 65535|`,
		BD(), V(1), AU16([]uint16{0, 10, 50, 128, 254, 255, 65535}), ED())

	eOpts.DefaultArrayEncodingBases.Uint16 = 16
	assertDecodeEncode(t, nil, eOpts, `c1 |u16x 0000 0001 0050 007f 0080 00ff ffff|`,
		BD(), V(1), AU16([]uint16{0x00, 0x01, 0x50, 0x7f, 0x80, 0xff, 0xffff}), ED())

	assertDecode(t, nil, `c1 |u16 00 01 01 0b101 0b110 0B101 0B110 0o10 0o11 0O10 0O11 0x7f 0x80 0X7f 0X80 0xff 0Xff|`,
		BD(), V(1), AU16([]uint16{0, 1, 1, 5, 6, 5, 6, 8, 9, 8, 9, 127, 128, 127, 128, 255, 255}), ED())

	assertDecodeFails(t, "c1 |u16b 10000000000000000|")
	assertDecodeFails(t, "c1 |u16o 200000|")
	assertDecodeFails(t, "c1 |u16 65536|")
	assertDecodeFails(t, "c1 |u16x 10000|")
}

func TestCTEArrayInt32(t *testing.T) {
	// TODO: TestCTEArrayInt32
}

func TestCTEArrayUint32(t *testing.T) {
	// TODO: TestCTEArrayUint32
}

func TestCTEArrayInt64(t *testing.T) {
	// TODO: TestCTEArrayInt64
}

func TestCTEArrayUint64(t *testing.T) {
	// TODO: TestCTEArrayUint64
}

func TestCTEArrayFloat16(t *testing.T) {
	// TODO: TestCTEArrayFloat16
}

func TestCTEArrayFloat32(t *testing.T) {
	// TODO: TestCTEArrayFloat32
}

func TestCTEArrayFloat64(t *testing.T) {
	// TODO: TestCTEArrayFloat64
}

func TestCTEArrayUUID(t *testing.T) {
	// TODO: TestCTEArrayUUID
}

func TestCTEArrayBool(t *testing.T) {
	// TODO: TestCTEArrayBool
}

func TestCTEBadArrayType(t *testing.T) {
	assertDecodeFails(t, `c1 x"01"`)
}

func TestCTEChunked(t *testing.T) {
	assertChunkedStringlike := func(encoded string, startEvent *test.TEvent) {
		assertEncode(t, nil, encoded, BD(), V(1), startEvent, AC(8, false), AD([]byte("abcdefgh")), ED())
		assertEncode(t, nil, encoded, BD(), V(1), startEvent,
			AC(1, true), AD([]byte("a")),
			AC(2, true), AD([]byte("bc")),
			AC(3, true), AD([]byte("def")),
			AC(2, false), AD([]byte("gh")),
			ED())

		assertEncode(t, nil, encoded, BD(), V(1), startEvent,
			AC(1, true), AD([]byte("a")),
			AC(2, true), AD([]byte("bc")),
			AC(3, true), AD([]byte("def")),
			AC(2, true), AD([]byte("gh")),
			AC(0, false), ED())
	}

	assertChunkedByteslike := func(encoded string, startEvent *test.TEvent) {
		assertEncode(t, nil, encoded, BD(), V(1), startEvent, AC(5, false), AD([]byte{0x12, 0x34, 0x56, 0x78, 0x9a}), ED())
		assertEncode(t, nil, encoded, BD(), V(1), startEvent,
			AC(1, true), AD([]byte{0x12}),
			AC(2, true), AD([]byte{0x34, 0x56}),
			AC(2, false), AD([]byte{0x78, 0x9a}),
			ED())
	}

	assertChunkedStringlike("c1 abcdefgh", SB())
	assertChunkedStringlike("c1 `#abcdefgh#", VB())
	assertChunkedStringlike(`c1 u"abcdefgh"`, UB())
	assertChunkedStringlike(`c1 t"abcdefgh"`, CTB())
	assertChunkedByteslike(`c1 b"123456789a"`, CBB())
	assertChunkedByteslike(`c1 |u8x 12 34 56 78 9a|`, AU8B())
}

func TestCTEList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 []`, BD(), V(1), L(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [123]`, BD(), V(1), L(), PI(123), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [test]`, BD(), V(1), L(), S("test"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [-1 a 2 test -3]`, BD(), V(1), L(), NI(1), S("a"), PI(2), S("test"), NI(3), E(), ED())
}

func TestCTEDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertMarshalUnmarshal(t, v, "c1 [[] [] []]")
}

func TestCTEMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 {}`, BD(), V(1), M(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1=2}`, BD(), V(1), M(), PI(1), PI(2), E(), ED())
	assertDecode(t, nil, "c1 {  1 = 2 3=4 \t}", BD(), V(1), M(), PI(1), PI(2), PI(3), PI(4), E(), ED())
	assertDecodeEncode(t, nil, nil, "c1 {nil=@nil 1.5=1000}")

	assertDecode(t, nil, `c1 {email = u"mailto:me@somewhere.com" 1.5 = "a string"}`, BD(), V(1), M(),
		S("email"), URI("mailto:me@somewhere.com"),
		DF(NewDFloat("1.5")), S("a string"),
		E(), ED())

	assertDecodeEncode(t, nil, nil, `c1 {a=@inf b=1}`)
	assertDecodeEncode(t, nil, nil, `c1 {a=-@inf b=1}`)
}

func TestCTEMapBadKVSeparator(t *testing.T) {
	assertDecodeFails(t, "c1 {a:b}")
}

func TestCTEListList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 [[]]`, BD(), V(1), L(), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 []]`, BD(), V(1), L(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 [] 1]`, BD(), V(1), L(), PI(1), L(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 [2] 1]`, BD(), V(1), L(), PI(1), L(), PI(2), E(), PI(1), E(), ED())
}

func TestCTEListMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 [{}]`, BD(), V(1), L(), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 {}]`, BD(), V(1), L(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 {} 1]`, BD(), V(1), L(), PI(1), M(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 [1 {2=3} 1]`, BD(), V(1), L(), PI(1), M(), PI(2), PI(3), E(), PI(1), E(), ED())
}

func TestCTEMapList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 {1=[]}`, BD(), V(1), M(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1=[2] test=[1 2 3]}`, BD(), V(1), M(), PI(1), L(), PI(2), E(), S("test"), L(), PI(1), PI(2), PI(3), E(), E(), ED())
}

func TestCTEMapMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 {1={}}`, BD(), V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1={a=b} test={}}`, BD(), V(1), M(), PI(1), M(), S("a"), S("b"), E(), S("test"), M(), E(), E(), ED())
}

func TestCTEMetadata(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 ()`, BD(), V(1), META(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 (1=2)`, BD(), V(1), META(), PI(1), PI(2), E(), ED())
	assertDecode(t, nil, "c1 (  1 = 2 3=4 \t)", BD(), V(1), META(), PI(1), PI(2), PI(3), PI(4), E(), ED())
}

func TestCTEMarkup(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 <a>`, BD(), V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a 1=2 3=4>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), PI(3), PI(4), E(), E(), ED())
	assertDecode(t, nil, `c1 <a;>`, BD(), V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;a>`, BD(), V(1), MUP(), S("a"), E(), S("a"), E(), ED())
	assertDecode(t, nil, `c1 <a;a string >`, BD(), V(1), MUP(), S("a"), E(), S("a string"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;a<a>>`, BD(), V(1), MUP(), S("a"), E(), S("a"), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecode(t, nil, `c1 <a 1=2 ;>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a 1=2;a>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a 1=2;<a>>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a 1=2;a <a>>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a "), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;***>`, BD(), V(1), MUP(), S("a"), E(), S("***"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;/x>`, BD(), V(1), MUP(), S("a"), E(), S("/x"), E(), ED())

	assertDecodeEncode(t, nil, nil, `c1 <a;\\>`, BD(), V(1), MUP(), S("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;\210>`, BD(), V(1), MUP(), S("a"), E(), S("\u0010"), E(), ED())

	assertDecodeEncode(t, nil, nil, `c1 <a;\\>`, BD(), V(1), MUP(), S("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;\<>`, BD(), V(1), MUP(), S("a"), E(), S("<"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 <a;\>>`, BD(), V(1), MUP(), S("a"), E(), S(">"), E(), ED())
	assertDecodeEncode(t, nil, nil, "c1 <a;\\`>", BD(), V(1), MUP(), S("a"), E(), S("`"), E(), ED())
	assertDecode(t, nil, `c1 <a;\r>`, BD(), V(1), MUP(), S("a"), E(), S("\r"), E(), ED())
	assertDecode(t, nil, `c1 <a;\n>`, BD(), V(1), MUP(), S("a"), E(), S("\n"), E(), ED())
	assertDecode(t, nil, `c1 <a;\t>`, BD(), V(1), MUP(), S("a"), E(), S("\t"), E(), ED())
	assertDecode(t, nil, `c1 <a;\*>`, BD(), V(1), MUP(), S("a"), E(), S("*"), E(), ED())
	assertDecode(t, nil, `c1 <a;\/>`, BD(), V(1), MUP(), S("a"), E(), S("/"), E(), ED())

	assertDecodeFails(t, `c1 <a;\y>`)
}

func TestCTEMarkupVerbatimString(t *testing.T) {
	assertDecode(t, nil, "c1 <s; `## <d></d>##>")
	assertDecode(t, nil, "c1 <s; `## /d##>")
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
}

func TestCTEMarkupComment(t *testing.T) {
	assertDecode(t, nil, "c1 <a;//blah\n>", BD(), V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, nil, "c1 <a;//blah\n a>", BD(), V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S("a"), E(), ED())
	assertDecode(t, nil, "c1 <a;a//blah\n a>", BD(), V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), S("a"), E(), ED())

	assertDecode(t, nil, "c1 <a;/*blah*/>", BD(), V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, nil, "c1 <a;a/*blah*/>", BD(), V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, nil, "c1 <a;/*blah*/a>", BD(), V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S("a"), E(), ED())

	assertDecode(t, nil, "c1 <a;/*/*blah*/*/>", BD(), V(1), MUP(), S("a"), E(), CMT(), CMT(), S("blah"), E(), E(), E(), ED())
	assertDecode(t, nil, "c1 <a;a/*/*blah*/*/>", BD(), V(1), MUP(), S("a"), E(), S("a"), CMT(), CMT(), S("blah"), E(), E(), E(), ED())
	assertDecode(t, nil, "c1 <a;/*/*blah*/*/a>", BD(), V(1), MUP(), S("a"), E(), CMT(), CMT(), S("blah"), E(), E(), S("a"), E(), ED())

	// TODO: Should it be picking up the extra space between the x and comment?
	assertDecode(t, nil, "c1 <a;x /*blah*/ x>", BD(), V(1), MUP(), S("a"), E(), S("x "), CMT(), S("blah"), E(), S("x"), E(), ED())
}

func TestCTEMapMetadata(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 [1 ()a]`, BD(), V(1), L(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1=()a}`, BD(), V(1), M(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1={}}`, BD(), V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {1=(){}}`, BD(), V(1), M(), PI(1), META(), E(), M(), E(), E(), ED())

	assertDecodeEncode(t, nil, nil, `c1 {()()1=()()a}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {()()1=()(){}}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c1 {()()1=()()[]}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), L(), E(), E(), ED())

	assertDecodeEncode(t, nil, nil, `c1 (x=y){(x=y)1=(x=y)(x=y){a=b}}`, BD(), V(1),
		META(), S("x"), S("y"), E(), M(),
		META(), S("x"), S("y"), E(), PI(1),
		META(), S("x"), S("y"), E(),
		META(), S("x"), S("y"), E(),
		M(), S("a"), S("b"), E(), E(), ED())
}

func TestCTENamed(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 @nil`, BD(), V(1), N(), ED())
	assertDecodeEncode(t, nil, nil, `c1 @nan`, BD(), V(1), NAN(), ED())
	assertDecodeEncode(t, nil, nil, `c1 @snan`, BD(), V(1), SNAN(), ED())
	assertDecodeEncode(t, nil, nil, `c1 @inf`, BD(), V(1), F(math.Inf(1)), ED())
	assertDecodeEncode(t, nil, nil, `c1 -@inf`, BD(), V(1), F(math.Inf(-1)), ED())
	assertDecodeEncode(t, nil, nil, `c1 @false`, BD(), V(1), FF(), ED())
	assertDecodeEncode(t, nil, nil, `c1 @true`, BD(), V(1), TT(), ED())
}

func TestCTEMarker(t *testing.T) {
	assertDecodeFails(t, `c1 &2`)
	assertDecode(t, nil, `c1 &1:string`, BD(), V(1), MARK(), PI(1), S("string"), ED())
	assertDecode(t, nil, `c1 &a:string`, BD(), V(1), MARK(), S("a"), S("string"), ED())
	assertDecodeFails(t, `c1 & 1:string`)
	assertDecodeFails(t, `c1 &1 string`)
	assertDecodeFails(t, `c1 &1string`)
	assertDecodeFails(t, `c1 &rgnsekfrnsekrgfnskergnslekrgnslergselrgblserfbserfbvsekrskfrvbskerfbksefbskerbfserbfrbksuerfbsekjrfbdjfgbsdjfgbsdfgbsdjkhfg`)
	assertDecodeFails(t, `c1 &100000000000000000000000000000000000000000000000`)
}

func TestCTEReference(t *testing.T) {
	assertDecode(t, nil, `c1 $2`, BD(), V(1), REF(), PI(2), ED())
	assertDecode(t, nil, `c1 $a`, BD(), V(1), REF(), S("a"), ED())
	assertDecodeFails(t, `c1 $ 1`)

	assertDecode(t, nil, `c1
{
    outside_ref      = $u"https://"
    // The markup type is good for presentation data
}
`)
}

func TestCTEMarkerReference(t *testing.T) {
	assertDecode(t, nil, `c1 [&2:testing $2]`, BD(), V(1), L(), MARK(), PI(2), S("testing"), REF(), PI(2), E(), ED())
	assertDecodeEncode(t, nil, nil, "c1 {first=&1:1000 second=$1}")
}

func TestCTEComment(t *testing.T) {
	// TODO: Better comment formatting
	assertDecodeEncode(t, nil, nil, `c1 {a=@inf /*test*/b=1}`)
}

func TestCTECommentSingleLine(t *testing.T) {
	assertDecodeFails(t, "c1 //")
	assertDecode(t, nil, "c1 //\n", BD(), V(1), CMT(), E(), ED())
	assertDecode(t, nil, "c1 //\r\n", BD(), V(1), CMT(), E(), ED())
	assertDecodeFails(t, "c1 // ")
	assertDecode(t, nil, "c1 // \n", BD(), V(1), CMT(), S(" "), E(), ED())
	assertDecode(t, nil, "c1 // \r\n", BD(), V(1), CMT(), S(" "), E(), ED())
	assertDecodeFails(t, "c1 //a")
	assertDecode(t, nil, "c1 //a\n", BD(), V(1), CMT(), S("a"), E(), ED())
	assertDecode(t, nil, "c1 //a\r\n", BD(), V(1), CMT(), S("a"), E(), ED())
	assertDecode(t, nil, "c1 // This is a comment\n", BD(), V(1), CMT(), S(" This is a comment"), E(), ED())
	assertDecodeFails(t, "c1 /-\n")
}

func TestCTECommentMultiline(t *testing.T) {
	assertDecode(t, nil, "c1 /**/", BD(), V(1), CMT(), E(), ED())
	assertDecode(t, nil, "c1 /* */", BD(), V(1), CMT(), S(" "), E(), ED())
	assertDecode(t, nil, "c1 /* This is a comment */", BD(), V(1), CMT(), S(" This is a comment "), E(), ED())
	assertDecode(t, nil, "c1 /*This is a comment*/", BD(), V(1), CMT(), S("This is a comment"), E(), ED())
}

func TestCTECommentMultilineNested(t *testing.T) {
	assertDecode(t, nil, "c1 /*/**/*/", BD(), V(1), CMT(), CMT(), E(), E(), ED())
	assertDecode(t, nil, "c1 /*/* */*/", BD(), V(1), CMT(), CMT(), S(" "), E(), E(), ED())
	assertDecode(t, nil, "c1 /* /* */ */", BD(), V(1), CMT(), S(" "), CMT(), S(" "), E(), S(" "), E(), ED())
	assertDecode(t, nil, "c1  /* before/* mid */ after*/  ", BD(), V(1), CMT(), S(" before"), CMT(), S(" mid "), E(), S(" after"), E(), ED())
}

func TestCTECommentAfterValue(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c1 [a /**/]`, BD(), V(1), L(), S("a"), CMT(), E(), E(), ED())
}

func TestCTEComplexComment(t *testing.T) {
	// DebugPrintEvents = true

	document := []byte(`c1
/**/ ( /**/ a= /**/ b /**/ ) /**/
<a;
    /**/
    <b>
>`)

	expected := `c1
/**/
(
    /**/
    a = /**/b
    /**/
)
/**/
<a;
    /**/
    <b>
>`

	encoded := &bytes.Buffer{}
	encOpts := options.DefaultCTEEncoderOptions()
	encOpts.Indent = "    "
	encoder := NewEncoder(encOpts)
	encoder.PrepareToEncode(encoded)
	decoder := NewDecoder(nil)
	err := decoder.Decode(bytes.NewBuffer(document), encoder)
	if err != nil {
		t.Error(err)
		return
	}

	actual := string(encoded.Bytes())
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestCTEBufferEdge(t *testing.T) {
	assertDecode(t, nil, `c1
{
     1  = <a;
            <b;
               <c; `+"`"+`##                       ##>
                         >
                       >
}
`)
}

func TestCTEBufferEdge2(t *testing.T) {
	assertDecode(t, nil, `c1
{
    x  = <a;
                     <b;
                             <c; `+"`"+`##                     ##>
                           >
                       >
}
`)
}

func TestCTEComplexExample(t *testing.T) {
	// DebugPrintEvents = true
	assertDecodeWithRules(t, `c1
// Metadata: _ct is the creation time
(_ct = 2019-9-1/22:14:01)
{
    /* /* Nested comments are allowed */ */
    // There are no commas in maps and lists
    (info = "something interesting about a_list")
    a_list           = [1 2 "a string"]
    map              = {2=two 3=3000 1=one}
    string           = "A string value"
    boolean          = @true
    "binary int"     = -0b10001011
    "octal int"      = 0o644
    "regular int"    = -10000000
    "hex int"        = 0xfffe0001
    "decimal float"  = -14.125
    "hex float"      = 0x5.1ec4p20
    uuid             = @f1ce4567-e89b-12d3-a456-426655440000
    date             = 2019-7-1
    time             = 18:04:00.940231541/E/Prague
    timestamp        = 2010-7-15/13:28:15.415942344/Z
    nil              = @nil
    bytes            = |u8x 10 ff 38 9a dd 00 4f 4f 91|
    url              = u"https://example.com/"
    email            = u"mailto:me@somewhere.com"
    1.5              = "Keys don't have to be strings"
    long-string      = `+"`"+`ZZZ
A backtick induces verbatim processing, which in this case will continue
until three Z characters are encountered, similar to how here documents in
bash work.
You can put anything in here, including double-quote ("), or even more
backticks (`+"`"+`). Verbatim processing stops at the end sequence, which in this
case is three Z characters, specified earlier as a sentinel.ZZZ
    marked_object    = &tag1:{
                                description = "This map will be referenced later using $tag1"
                                value = -@inf
                                child_elements = @nil
                                recursive = $tag1
                            }
    ref1             = $tag1
    ref2             = $tag1
    outside_ref      = $u"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    html_compatible  = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en ;
                         <body;
                           Please choose from the following widgets:
                           <div id=parent style=normal ref-id=1 ;
                             /* Here we use a backtick to induce verbatim processing.
                              * In this case, "##" is chosen as the ending sequence
                              */
                             <script; `+"`"+`##
                               document.getElementById('parent').insertAdjacentHTML('beforeend',
                                  '<div id="idChild"> content </div>');
                             ##>
                           >
                         >
                       >
}
`)
}

func TestCTEEncodeDecodeExample(t *testing.T) {
	// DebugPrintEvents = true

	document := []byte(`c1
// _ct is the creation time, in this case referring to the entire document
(
    _ct = 2019-9-1/22:14:01
)
{
    /* Comments look very C-like, except:
       /* Nested comments are allowed! */
    */
    // Notice that there are no commas in maps and lists
    (
        info = "something interesting about a_list"
    )
    a_list = [
        1
        2
        "a string"
    ]
    map = {
        2=two
        3=3000
        1=one
    }
    string = "A string value"
    boolean = @true
    "binary int" = -0b10001011
    "octal int" = 0o644
    "regular int" = -10000000
    "hex int" = 0xfffe0001
    "decimal float" = -14.125
    "hex float" = 0x5.1ec4p20
    uuid = @f1ce4567-e89b-12d3-a456-426655440000
    date = 2019-7-1
    time = 18:04:00.940231541/E/Prague
    timestamp = 2010-7-15/13:28:15.415942344/Z
    nil = @nil
    bytes = |u8x 10 ff 38 9a dd 00 4f 4f 91|
    url = u"https://example.com/"
    email = u"mailto:me@somewhere.com"
    1.5 = "Keys don't have to be strings"
    long-string = ` + "`" + `ZZZ
A backtick induces verbatim processing, which in this case will continue
until three Z characters are encountered, similar to how here documents in
bash work.
You can put anything in here, including double-quote ("), or even more
backticks (` + "`" + `). Verbatim processing stops at the end sequence, which in this
case is three Z characters, specified earlier as a sentinel.ZZZ
    marked_object = &tag1:{
        description = "This map will be referenced later using $tag1"
        value = -@inf
        child_elements = @nil
        recursive = $tag1
    }
    ref1 = $tag1
    ref2 = $tag1
    outside_ref = $u"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    html_compatible  = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en ;
        <body;
            Please choose from the following widgets:
            <div id=parent style=normal ref-id=1 ;
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "##" is chosen as the ending sequence
                 */
                <script; ` + "`" + `##
                    document.getElementById('parent').insertAdjacentHTML('beforeend',
                        '<div id="idChild"> content </div>');
                ##>
            >
        >
    >
}
`)

	expected := `c1
/* _ct is the creation time, in this case referring to the entire document*/
(
    _ct = 2019-09-01/22:14:01
)
{
    /* Comments look very C-like, except:
        /* Nested comments are allowed! */
    */
    /* Notice that there are no commas in maps and lists*/
    (
        info = "something interesting about a_list"
    )
    a_list = [
        1
        2
        "a string"
    ]
    map = {
        2 = two
        3 = 3000
        1 = one
    }
    string = "A string value"
    boolean = @true
    "binary int" = -139
    "octal int" = 420
    "regular int" = -10000000
    "hex int" = 4294836225
    "decimal float" = -14.125
    "hex float" = 5.368896e+06
    uuid = @f1ce4567-e89b-12d3-a456-426655440000
    date = 2019-07-01
    time = 18:04:00.940231541/E/Prague
    timestamp = 2010-07-15/13:28:15.415942344/Z
    nil = @nil
    bytes = |u8x 10 ff 38 9a dd 00 4f 4f 91|
    url = u"https://example.com/"
    email = u"mailto:me@somewhere.com"
    1.5 = "Keys don't have to be strings"
    long-string = ` + "`" + `#A backtick induces verbatim processing, which in this case will continue
until three Z characters are encountered, similar to how here documents in
bash work.
You can put anything in here, including double-quote ("), or even more
backticks (` + "`" + `). Verbatim processing stops at the end sequence, which in this
case is three Z characters, specified earlier as a sentinel.#
    marked_object = &tag1:{
        description = "This map will be referenced later using $tag1"
        value = -@inf
        child_elements = @nil
        recursive = $tag1
    }
    ref1 = $tag1
    ref2 = $tag1
    outside_ref = $u"https://somewhere.else.com/path/to/document.cte#some_tag"
    /* The markup type is good for presentation data*/
    html_compatible = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en;
        <body;
            Please choose from the following widgets: <div id=parent style=normal ref-id=1;
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "##" is chosen as the ending sequence
                 */
                <script;
                    ` + "`" + `#                    document.getElementById('parent').insertAdjacentHTML('beforeend',
                        '<div id="idChild"> content </div>');
                #
                >
            >
        >
    >
}`

	encoded := &bytes.Buffer{}
	encOpts := options.DefaultCTEEncoderOptions()
	encOpts.Indent = "    "
	encoder := NewEncoder(encOpts)
	encoder.PrepareToEncode(encoded)
	decoder := NewDecoder(nil)
	err := decoder.Decode(bytes.NewBuffer(document), encoder)
	if err != nil {
		t.Error(err)
		return
	}

	actual := string(encoded.Bytes())
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestCTEEncodeImpliedVersion(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureVersion
	assertEncode(t, opts, "[1 2]", BD(), V(1), L(), PI(1), PI(2), E(), ED())
}

func TestCTEDecodeImpliedVersion(t *testing.T) {
	opts := options.DefaultCTEDecoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureVersion
	assertDecode(t, opts, "[1 2]", BD(), V(1), L(), PI(1), PI(2), E(), ED())
}

func TestCTEEncodeImpliedList(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureList
	assertEncode(t, opts, "1 2", BD(), V(1), L(), PI(1), PI(2), E(), ED())
}

func TestCTEDecodeImpliedList(t *testing.T) {
	opts := options.DefaultCTEDecoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureList
	assertDecode(t, opts, "1 2", BD(), V(1), L(), PI(1), PI(2), E(), ED())
}

func TestCTEEncodeImpliedMap(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureMap
	assertEncode(t, opts, "1=2", BD(), V(1), M(), PI(1), PI(2), E(), ED())
}

func TestCTEDecodeImpliedMap(t *testing.T) {
	opts := options.DefaultCTEDecoderOptions()
	opts.ImpliedStructure = options.ImpliedStructureMap
	assertDecode(t, opts, "1=2", BD(), V(1), M(), PI(1), PI(2), E(), ED())
}
