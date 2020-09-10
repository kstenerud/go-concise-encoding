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
	"math/big"
	"testing"

	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-compact-time"
)

func TestCTEVersion(t *testing.T) {
	assertDecodeEncode(t, "c1 ", BD(), V(1), ED())
	assertDecode(t, nil, "\r\n\t c1 ", BD(), V(1), ED())
	assertDecode(t, nil, "c1     \r\n\t\t\t", BD(), V(1), ED())
}

func TestCTEVersionNotNumeric(t *testing.T) {
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		document := string([]byte{'c', byte(i)})
		assertDecodeFails(t, document)
	}
}

func TestCTEVersionMissingWhitespace(t *testing.T) {
	assertDecodeFails(t, "c1")
}

func TestCTEVersionTooBig(t *testing.T) {
	assertDecodeFails(t, "c100000000000000000000000000000000 ")
}

func TestCTEBadVersion(t *testing.T) {
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

func TestCTENoWSAfterVersion(t *testing.T) {
	assertDecodeFails(t, "c1{}")
}

func TestCTENil(t *testing.T) {
	assertDecodeEncode(t, "c1 @nil", BD(), V(1), N(), ED())
}

func TestCTEBool(t *testing.T) {
	assertDecodeEncode(t, "c1 @true", BD(), V(1), TT(), ED())
	assertDecodeEncode(t, "c1 @false", BD(), V(1), FF(), ED())
}

func TestCTEDecimalInt(t *testing.T) {
	assertDecodeEncode(t, "c1 0", BD(), V(1), PI(0), ED())
	assertDecodeEncode(t, "c1 123", BD(), V(1), PI(123), ED())
	assertDecodeEncode(t, "c1 9412504234235366", BD(), V(1), PI(9412504234235366), ED())
	assertDecodeEncode(t, "c1 -49523", BD(), V(1), NI(49523), ED())
	assertDecodeEncode(t, "c1 10000000000000000000000000000", BD(), V(1), BI(NewBigInt("10000000000000000000000000000", 10)), ED())
	assertDecodeEncode(t, "c1 -10000000000000000000000000000", BD(), V(1), BI(NewBigInt("-10000000000000000000000000000", 10)), ED())
	assertDecode(t, nil, "c1 -4_9_5__2___3", BD(), V(1), NI(49523), ED())
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
}

func TestCTEFloat(t *testing.T) {
	assertDecode(t, nil, "c1 0.0", BD(), V(1), DF(NewDFloat("0")), ED())
	assertDecode(t, nil, "c1 -0.0", BD(), V(1), DF(NewDFloat("-0")), ED())

	assertDecodeEncode(t, "c1 1.5", BD(), V(1), DF(NewDFloat("1.5")), ED())
	assertDecodeEncode(t, "c1 1.125", BD(), V(1), DF(NewDFloat("1.125")), ED())
	assertDecodeEncode(t, "c1 1.125e+10", BD(), V(1), DF(NewDFloat("1.125e+10")), ED())
	assertDecodeEncode(t, "c1 1.125e-10", BD(), V(1), DF(NewDFloat("1.125e-10")), ED())
	assertDecode(t, nil, "c1 1.125e10", BD(), V(1), DF(NewDFloat("1.125e+10")), ED())

	assertDecodeEncode(t, "c1 -1.5", BD(), V(1), DF(NewDFloat("-1.5")), ED())
	assertDecodeEncode(t, "c1 -1.125", BD(), V(1), DF(NewDFloat("-1.125")), ED())
	assertDecodeEncode(t, "c1 -1.125e+10", BD(), V(1), DF(NewDFloat("-1.125e+10")), ED())
	assertDecodeEncode(t, "c1 -1.125e-10", BD(), V(1), DF(NewDFloat("-1.125e-10")), ED())
	assertDecode(t, nil, "c1 -1.125e10", BD(), V(1), DF(NewDFloat("-1.125e10")), ED())

	assertDecodeEncode(t, "c1 0.5", BD(), V(1), DF(NewDFloat("0.5")), ED())
	assertDecodeEncode(t, "c1 0.125", BD(), V(1), DF(NewDFloat("0.125")), ED())
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

	assertDecodeFails(t, "c1 -0.5.4")
	assertDecodeFails(t, "c1 -0,5.4")
	assertDecodeFails(t, "c1 0.5.4")
	assertDecodeFails(t, "c1 0,5.4")
	assertDecodeFails(t, "c1 -@blah")
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
	bigExpected, _, err := big.ParseFloat("-1.54fffe2ac00592375b427ap100000", 16, 90, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	assertDecode(t, nil, "c1 -0x1.54fffe2ac00592375b427ap100000", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.54fffe2ac00592375b427ap100000", BD(), V(1), BF(bigExpected), ED())

	// Coefficient too big for float64
	bigExpected, _, err = big.ParseFloat("-1.54fffe2ac00592375b427ap100", 16, 90, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	assertDecode(t, nil, "c1 -0x1.54fffe2ac00592375b427ap100", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.54fffe2ac00592375b427ap100", BD(), V(1), BF(bigExpected), ED())

	// Exponent too big for float64
	bigExpected, _, err = big.ParseFloat("-1.8p100000", 16, 64, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	assertDecode(t, nil, "c1 -0x1.8p100000", BD(), V(1), BF(bigExpected), ED())
	bigExpected = bigExpected.Neg(bigExpected)
	assertDecode(t, nil, "c1 0x1.8p100000", BD(), V(1), BF(bigExpected), ED())

	assertDecode(t, nil, "c1 -0x_0_._1_p_1_0", BD(), V(1), F(-0x0.1p10), ED())
}

func TestCTEUUID(t *testing.T) {
	assertDecodeEncode(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbb`, BD(), V(1),
		UUID([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertDecodeEncode(t, `c1 @00000000-0000-0000-0000-000000000000`, BD(), V(1),
		UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), ED())
}

func TestCTEDate(t *testing.T) {
	assertDecodeEncode(t, "c1 2000-01-01", BD(), V(1), CT(compact_time.NewDate(2000, 1, 1)), ED())
	assertDecodeEncode(t, "c1 -2000-12-31", BD(), V(1), CT(compact_time.NewDate(-2000, 12, 31)), ED())
}

func TestCTETime(t *testing.T) {
	assertDecode(t, nil, "c1 1:45:00", BD(), V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecode(t, nil, "c1 01:45:00", BD(), V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecodeEncode(t, "c1 23:59:59.101", BD(), V(1), CT(compact_time.NewTime(23, 59, 59, 101000000, "")), ED())
	assertDecodeEncode(t, "c1 10:00:01.93/America/Los_Angeles", BD(), V(1), CT(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertDecodeEncode(t, "c1 10:00:01.93/89.92/1.10", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/0/0", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ED())
	assertDecode(t, nil, "c1 10:00:01.93/1/1", BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ED())
}

func TestCTETimestamp(t *testing.T) {
	assertDecodeEncode(t, "c1 2000-01-01/19:31:44.901554/Z", BD(), V(1), CT(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "Z")), ED())
	assertDecodeEncode(t, "c1 2020-01-15/13:41:00.000599", BD(), V(1), CT(compact_time.NewTimestamp(2020, 1, 15, 13, 41, 0, 599000, "")), ED())
	assertDecode(t, nil, "c1 2020-01-15/13:41:00.000599", BD(), V(1), CT(compact_time.NewTimestamp(2020, 1, 15, 13, 41, 0, 599000, "")), ED())
}

func TestCTEQuotedString(t *testing.T) {
	assertDecodeEncode(t, `c1 "test string"`, BD(), V(1), S("test string"), ED())
	assertDecode(t, nil, `c1 "test\nstring"`, BD(), V(1), S("test\nstring"), ED())
	assertDecode(t, nil, `c1 "test\rstring"`, BD(), V(1), S("test\rstring"), ED())
	assertDecode(t, nil, `c1 "test\tstring"`, BD(), V(1), S("test\tstring"), ED())
	assertDecodeEncode(t, `c1 "test\"string"`, BD(), V(1), S("test\"string"), ED())
	assertDecode(t, nil, `c1 "test\*string"`, BD(), V(1), S("test*string"), ED())
	assertDecode(t, nil, `c1 "test\/string"`, BD(), V(1), S("test/string"), ED())
	assertDecodeEncode(t, `c1 "test\\string"`, BD(), V(1), S("test\\string"), ED())
	assertDecodeEncode(t, `c1 "test\11string"`, BD(), V(1), S("test\u0001string"), ED())
	assertDecodeEncode(t, `c1 "test\4206dstring"`, BD(), V(1), S("test\u206dstring"), ED())
	assertDecode(t, nil, `c1 "test\
string"`, BD(), V(1), S("teststring"), ED())
	assertDecode(t, nil, "c1 \"test\\\r\nstring\"", BD(), V(1), S("teststring"), ED())
}

func TestCTECustomBinary(t *testing.T) {
	assertDecodeEncode(t, `c1 b"12345678"`, BD(), V(1), CUB([]byte{0x12, 0x34, 0x56, 0x78}), ED())
}

func TestCTECustomText(t *testing.T) {
	assertDecodeEncode(t, `c1 t"something(123)"`, BD(), V(1), CUT("something(123)"), ED())
	assertDecodeEncode(t, `c1 t"some\\thing(\"123\")"`, BD(), V(1), CUT("some\\thing(\"123\")"), ED())
	assertDecodeEncode(t, `c1 t"some\nthing\11(123)"`, BD(), V(1), CUT("some\nthing\u0001(123)"), ED())
}

func TestCTEUnquotedString(t *testing.T) {
	assertDecodeEncode(t, "c1 a", BD(), V(1), S("a"), ED())
	assertDecodeEncode(t, "c1 abcd", BD(), V(1), S("abcd"), ED())
	assertDecodeEncode(t, "c1 _-.123aF", BD(), V(1), S("_-.123aF"), ED())
	assertDecodeEncode(t, "c1 新しい", BD(), V(1), S("新しい"), ED())
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
	assertDecodeEncode(t, `c1 u"http://example.com"`, BD(), V(1), URI("http://example.com"), ED())
	assertEncode(t, nil, `c1 u"http://x.com/%22quoted%22"`, BD(), V(1), URI(`http://x.com/"quoted"`), ED())
}

// TODO: Other array types

func TestCTEUintXArray(t *testing.T) {
	assertDecodeEncode(t, `c1 |u8x f1 93|`, BD(), V(1), AU8([]byte{0xf1, 0x93}), ED())
	assertDecode(t, nil, `c1 |u8x f 93 |`, BD(), V(1), AU8([]byte{0xf, 0x93}), ED())
	assertDecodeEncode(t, `c1 |u16x f122 9385|`, BD(), V(1), AU16([]uint16{0xf122, 0x9385}), ED())
	assertDecode(t, nil, `c1 |u16x f12 95|`, BD(), V(1), AU16([]uint16{0xf12, 0x95}), ED())
	assertDecodeEncode(t, `c1 |u32x 7ddf8134 93cd7aac|`, BD(), V(1), AU32([]uint32{0x7ddf8134, 0x93cd7aac}), ED())
	assertDecode(t, nil, `c1 |u32x 7ddf834 93aac|`, BD(), V(1), AU32([]uint32{0x7ddf834, 0x93aac}), ED())
	assertDecodeEncode(t, `c1 |u64x 83ff9ac2445aace7 94ff7ac3219465c1|`, BD(), V(1), AU64([]uint64{0x83ff9ac2445aace7, 0x94ff7ac3219465c1}), ED())
	assertDecode(t, nil, `c1 |u64x 83ff9ac245aace7 94ff79465c1|`, BD(), V(1), AU64([]uint64{0x83ff9ac245aace7, 0x94ff79465c1}), ED())
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
	assertDecodeEncode(t, `c1 []`, BD(), V(1), L(), E(), ED())
	assertDecodeEncode(t, `c1 [123]`, BD(), V(1), L(), PI(123), E(), ED())
	assertDecodeEncode(t, `c1 [test]`, BD(), V(1), L(), S("test"), E(), ED())
	assertDecodeEncode(t, `c1 [-1 a 2 test -3]`, BD(), V(1), L(), NI(1), S("a"), PI(2), S("test"), NI(3), E(), ED())
}

func TestCTEDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertMarshalUnmarshal(t, v, "c1 [[] [] []]")
}

func TestCTEMap(t *testing.T) {
	assertDecodeEncode(t, `c1 {}`, BD(), V(1), M(), E(), ED())
	assertDecodeEncode(t, `c1 {1=2}`, BD(), V(1), M(), PI(1), PI(2), E(), ED())
	assertDecode(t, nil, "c1 {  1 = 2 3=4 \t}", BD(), V(1), M(), PI(1), PI(2), PI(3), PI(4), E(), ED())
	assertDecodeEncode(t, "c1 {nil=@nil 1.5=1000}")

	assertDecode(t, nil, `c1 {email = u"mailto:me@somewhere.com" 1.5 = "a string"}`, BD(), V(1), M(),
		S("email"), URI("mailto:me@somewhere.com"),
		DF(NewDFloat("1.5")), S("a string"),
		E(), ED())

	assertDecodeEncode(t, `c1 {a=@inf b=1}`)
	assertDecodeEncode(t, `c1 {a=-@inf b=1}`)
}

func TestCTEMapBadKVSeparator(t *testing.T) {
	assertDecodeFails(t, "c1 {a:b}")
}

func TestCTEListList(t *testing.T) {
	assertDecodeEncode(t, `c1 [[]]`, BD(), V(1), L(), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 []]`, BD(), V(1), L(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 [] 1]`, BD(), V(1), L(), PI(1), L(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, `c1 [1 [2] 1]`, BD(), V(1), L(), PI(1), L(), PI(2), E(), PI(1), E(), ED())
}

func TestCTEListMap(t *testing.T) {
	assertDecodeEncode(t, `c1 [{}]`, BD(), V(1), L(), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 {}]`, BD(), V(1), L(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 {} 1]`, BD(), V(1), L(), PI(1), M(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, `c1 [1 {2=3} 1]`, BD(), V(1), L(), PI(1), M(), PI(2), PI(3), E(), PI(1), E(), ED())
}

func TestCTEMapList(t *testing.T) {
	assertDecodeEncode(t, `c1 {1=[]}`, BD(), V(1), M(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1=[2] test=[1 2 3]}`, BD(), V(1), M(), PI(1), L(), PI(2), E(), S("test"), L(), PI(1), PI(2), PI(3), E(), E(), ED())
}

func TestCTEMapMap(t *testing.T) {
	assertDecodeEncode(t, `c1 {1={}}`, BD(), V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1={a=b} test={}}`, BD(), V(1), M(), PI(1), M(), S("a"), S("b"), E(), S("test"), M(), E(), E(), ED())
}

func TestCTEMetadata(t *testing.T) {
	assertDecodeEncode(t, `c1 ()`, BD(), V(1), META(), E(), ED())
	assertDecodeEncode(t, `c1 (1=2)`, BD(), V(1), META(), PI(1), PI(2), E(), ED())
	assertDecode(t, nil, "c1 (  1 = 2 3=4 \t)", BD(), V(1), META(), PI(1), PI(2), PI(3), PI(4), E(), ED())
}

func TestCTEMarkup(t *testing.T) {
	assertDecodeEncode(t, `c1 <a>`, BD(), V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2 3=4>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), PI(3), PI(4), E(), E(), ED())
	assertDecode(t, nil, `c1 <a;>`, BD(), V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a;a>`, BD(), V(1), MUP(), S("a"), E(), S("a"), E(), ED())
	assertDecode(t, nil, `c1 <a;a string >`, BD(), V(1), MUP(), S("a"), E(), S("a string"), E(), ED())
	assertDecodeEncode(t, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a;a<a>>`, BD(), V(1), MUP(), S("a"), E(), S("a"), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecode(t, nil, `c1 <a 1=2 ;>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2;a>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2;<a>>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2;a <a>>`, BD(), V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a "), MUP(), S("a"), E(), E(), E(), ED())

	assertDecodeEncode(t, `c1 <a;\\>`, BD(), V(1), MUP(), S("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, `c1 <a;\210>`, BD(), V(1), MUP(), S("a"), E(), S("\u0010"), E(), ED())
}

func TestCTEMarkupVerbatimString(t *testing.T) {
	assertDecode(t, nil, "c1 <s; `## <d></d>##>")
	assertDecode(t, nil, "c1 <s; `## /d##>")
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertDecodeEncode(t, `c1 <a;<a>>`, BD(), V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
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
}

func TestCTEMapMetadata(t *testing.T) {
	assertDecodeEncode(t, `c1 [1 ()a]`, BD(), V(1), L(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {1=()a}`, BD(), V(1), M(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {1={}}`, BD(), V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1=(){}}`, BD(), V(1), M(), PI(1), META(), E(), M(), E(), E(), ED())

	assertDecodeEncode(t, `c1 {()()1=()()a}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {()()1=()(){}}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {()()1=()()[]}`, BD(), V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), L(), E(), E(), ED())

	assertDecodeEncode(t, `c1 (x=y){(x=y)1=(x=y)(x=y){a=b}}`, BD(), V(1),
		META(), S("x"), S("y"), E(), M(),
		META(), S("x"), S("y"), E(), PI(1),
		META(), S("x"), S("y"), E(),
		META(), S("x"), S("y"), E(),
		M(), S("a"), S("b"), E(), E(), ED())
}

func TestCTENamed(t *testing.T) {
	assertDecodeEncode(t, `c1 @nil`, BD(), V(1), N(), ED())
	assertDecodeEncode(t, `c1 @nan`, BD(), V(1), NAN(), ED())
	assertDecodeEncode(t, `c1 @snan`, BD(), V(1), SNAN(), ED())
	assertDecodeEncode(t, `c1 @inf`, BD(), V(1), F(math.Inf(1)), ED())
	assertDecodeEncode(t, `c1 -@inf`, BD(), V(1), F(math.Inf(-1)), ED())
	assertDecodeEncode(t, `c1 @false`, BD(), V(1), FF(), ED())
	assertDecodeEncode(t, `c1 @true`, BD(), V(1), TT(), ED())
}

func TestCTEMarker(t *testing.T) {
	assertDecodeFails(t, `c1 &2`)
	assertDecode(t, nil, `c1 &1:string`, BD(), V(1), MARK(), PI(1), S("string"), ED())
	assertDecode(t, nil, `c1 &a:string`, BD(), V(1), MARK(), S("a"), S("string"), ED())
	assertDecodeFails(t, `c1 & 1:string`)
	assertDecodeFails(t, `c1 &1 string`)
	assertDecodeFails(t, `c1 &1string`)
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
	assertDecodeEncode(t, "c1 {first=&1:1000 second=$1}")
}

func TestCTEComment(t *testing.T) {
	// TODO: Better comment formatting
	assertDecodeEncode(t, `c1 {a=@inf /*test*/b=1}`)
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
	assertDecodeEncode(t, `c1 [a /**/]`, BD(), V(1), L(), S("a"), CMT(), E(), E(), ED())
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
