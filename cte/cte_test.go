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

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-compact-time"
)

// debug.DebugOptions.PassThroughPanics = true
// defer func() { debug.DebugOptions.PassThroughPanics = false }()

func TestCTEVersion(t *testing.T) {
	assertDecodeEncode(t, "c1 ", V(1), ED())
	assertDecode(t, "\r\n\t c1 ", V(1), ED())
	assertDecode(t, "c1     \r\n\t\t\t", V(1), ED())
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

func TestCTEUnquotedString(t *testing.T) {
	assertDecodeEncode(t, "c1 a", V(1), S("a"), ED())
	assertDecodeEncode(t, "c1 abcd", V(1), S("abcd"), ED())
	assertDecodeEncode(t, "c1 _-.:123aF", V(1), S("_-.:123aF"), ED())
	assertDecodeEncode(t, "c1 新しい", V(1), S("新しい"), ED())
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
	assertDecode(t, "c1 `A \n\n\n\n\n\n\n\n\n\nA", V(1), VS("\n\n\n\n\n\n\n\n\n\n"), ED())
	assertDecode(t, "c1 `A aA", V(1), VS("a"), ED())
	assertDecode(t, "c1 `A\taA", V(1), VS("a"), ED())
	assertDecode(t, "c1 `A\naA", V(1), VS("a"), ED())
	assertDecode(t, "c1 `A\r\naA", V(1), VS("a"), ED())
	assertDecode(t, "c1 `#ENDOFSTRING a test\nwith `stuff`#ENDOFSTRING ", V(1), VS("a test\nwith `stuff`"), ED())
}

func TestCTEChunkedString(t *testing.T) {
	assertEncode(t, "c1 abcdefgh", V(1), SB(), AC(8, false), AD([]byte("abcdefgh")), ED())

	assertEncode(t, "c1 abcdefgh", V(1), SB(),
		AC(1, true), AD([]byte("a")),
		AC(2, true), AD([]byte("bc")),
		AC(3, true), AD([]byte("def")),
		AC(2, false), AD([]byte("gh")),
		ED())

	assertEncode(t, "c1 abcdefgh", V(1), SB(),
		AC(1, true), AD([]byte("a")),
		AC(2, true), AD([]byte("bc")),
		AC(3, true), AD([]byte("def")),
		AC(2, true), AD([]byte("gh")),
		AC(0, false), ED())
}

func TestCTEDecimalInt(t *testing.T) {
	assertDecodeEncode(t, "c1 0", V(1), PI(0), ED())
	assertDecodeEncode(t, "c1 123", V(1), PI(123), ED())
	assertDecodeEncode(t, "c1 9412504234235366", V(1), PI(9412504234235366), ED())
	assertDecodeEncode(t, "c1 -49523", V(1), NI(49523), ED())
	assertDecodeEncode(t, "c1 10000000000000000000000000000", V(1), BI(test.NewBigInt("10000000000000000000000000000")), ED())
	assertDecodeEncode(t, "c1 -10000000000000000000000000000", V(1), BI(test.NewBigInt("-10000000000000000000000000000")), ED())
	assertDecode(t, "c1 -4_9_5__2___3", V(1), NI(49523), ED())
}

func TestCTEBinaryInt(t *testing.T) {
	assertDecode(t, "c1 0b0", V(1), PI(0), ED())
	assertDecode(t, "c1 0b1", V(1), PI(1), ED())
	assertDecode(t, "c1 0b101", V(1), PI(5), ED())
	assertDecode(t, "c1 0b0010100", V(1), PI(20), ED())
	assertDecode(t, "c1 -0b100", V(1), NI(4), ED())
	assertDecode(t, "c1 -0b_1_0_0", V(1), NI(4), ED())
}

func TestCTEOctalInt(t *testing.T) {
	assertDecode(t, "c1 0o0", V(1), PI(0), ED())
	assertDecode(t, "c1 0o1", V(1), PI(1), ED())
	assertDecode(t, "c1 0o7", V(1), PI(7), ED())
	assertDecode(t, "c1 0o71", V(1), PI(57), ED())
	assertDecode(t, "c1 0o644", V(1), PI(420), ED())
	assertDecode(t, "c1 -0o777", V(1), NI(511), ED())
	assertDecode(t, "c1 -0o_7__7___7", V(1), NI(511), ED())
}

func TestCTEHexInt(t *testing.T) {
	assertDecode(t, "c1 0x0", V(1), PI(0), ED())
	assertDecode(t, "c1 0x1", V(1), PI(1), ED())
	assertDecode(t, "c1 0xf", V(1), PI(0xf), ED())
	assertDecode(t, "c1 0xfedcba9876543210", V(1), PI(0xfedcba9876543210), ED())
	assertDecode(t, "c1 0xFEDCBA9876543210", V(1), PI(0xfedcba9876543210), ED())
	assertDecode(t, "c1 -0x88", V(1), NI(0x88), ED())
	assertDecode(t, "c1 -0x_8_8__5_a_f__d", V(1), NI(0x885afd), ED())
}

func TestCTEFloat(t *testing.T) {
	assertDecode(t, "c1 0.0", V(1), DF(test.NewDFloat("0")), ED())
	assertDecode(t, "c1 -0.0", V(1), DF(test.NewDFloat("-0")), ED())

	assertDecodeEncode(t, "c1 1.5", V(1), DF(test.NewDFloat("1.5")), ED())
	assertDecodeEncode(t, "c1 1.125", V(1), DF(test.NewDFloat("1.125")), ED())
	assertDecodeEncode(t, "c1 1.125e+10", V(1), DF(test.NewDFloat("1.125e+10")), ED())
	assertDecodeEncode(t, "c1 1.125e-10", V(1), DF(test.NewDFloat("1.125e-10")), ED())
	assertDecode(t, "c1 1.125e10", V(1), DF(test.NewDFloat("1.125e+10")), ED())

	assertDecodeEncode(t, "c1 -1.5", V(1), DF(test.NewDFloat("-1.5")), ED())
	assertDecodeEncode(t, "c1 -1.125", V(1), DF(test.NewDFloat("-1.125")), ED())
	assertDecodeEncode(t, "c1 -1.125e+10", V(1), DF(test.NewDFloat("-1.125e+10")), ED())
	assertDecodeEncode(t, "c1 -1.125e-10", V(1), DF(test.NewDFloat("-1.125e-10")), ED())
	assertDecode(t, "c1 -1.125e10", V(1), DF(test.NewDFloat("-1.125e10")), ED())

	assertDecodeEncode(t, "c1 0.5", V(1), DF(test.NewDFloat("0.5")), ED())
	assertDecodeEncode(t, "c1 0.125", V(1), DF(test.NewDFloat("0.125")), ED())
	assertDecode(t, "c1 0.125e+10", V(1), DF(test.NewDFloat("0.125e+10")), ED())
	assertDecode(t, "c1 0.125e-10", V(1), DF(test.NewDFloat("0.125e-10")), ED())
	assertDecode(t, "c1 0.125e10", V(1), DF(test.NewDFloat("0.125e10")), ED())

	assertDecode(t, "c1 -0.5", V(1), DF(test.NewDFloat("-0.5")), ED())
	assertDecode(t, "c1 -0.125", V(1), DF(test.NewDFloat("-0.125")), ED())
	assertDecode(t, "c1 -0.125e+10", V(1), DF(test.NewDFloat("-0.125e+10")), ED())
	assertDecode(t, "c1 -0.125e-10", V(1), DF(test.NewDFloat("-0.125e-10")), ED())
	assertDecode(t, "c1 -0.125e10", V(1), DF(test.NewDFloat("-0.125e10")), ED())
	assertDecode(t, "c1 -0.125E+10", V(1), DF(test.NewDFloat("-0.125e+10")), ED())
	assertDecode(t, "c1 -0.125E-10", V(1), DF(test.NewDFloat("-0.125e-10")), ED())
	assertDecode(t, "c1 -0.125E10", V(1), DF(test.NewDFloat("-0.125e10")), ED())

	assertDecode(t, "c1 -1.50000000000000000000000001E10000", V(1), BDF(test.NewBDF("-1.50000000000000000000000001E10000")), ED())
	assertDecode(t, "c1 1.50000000000000000000000001E10000", V(1), BDF(test.NewBDF("1.50000000000000000000000001E10000")), ED())

	assertDecode(t, "c1 1_._1_2_5_e+1_0", V(1), DF(test.NewDFloat("1.125e+10")), ED())

	assertDecodeFails(t, "c1 -0.5.4")
	assertDecodeFails(t, "c1 -0,5.4")
	assertDecodeFails(t, "c1 0.5.4")
	assertDecodeFails(t, "c1 0,5.4")
	assertDecodeFails(t, "c1 -@blah")
}

func TestCTEHexFloat(t *testing.T) {
	assertDecode(t, "c1 0x0.0", V(1), F(0x0.0p0), ED())
	assertDecode(t, "c1 0x0.1", V(1), F(0x0.1p0), ED())
	assertDecode(t, "c1 0x0.1p+10", V(1), F(0x0.1p+10), ED())
	assertDecode(t, "c1 0x0.1p-10", V(1), F(0x0.1p-10), ED())
	assertDecode(t, "c1 0x0.1p10", V(1), F(0x0.1p10), ED())

	assertDecode(t, "c1 0x1.0", V(1), F(0x1.0p0), ED())
	assertDecode(t, "c1 0x1.1", V(1), F(0x1.1p0), ED())
	assertDecode(t, "c1 0xf.1p+10", V(1), F(0xf.1p+10), ED())
	assertDecode(t, "c1 0xf.1p-10", V(1), F(0xf.1p-10), ED())
	assertDecode(t, "c1 0xf.1p10", V(1), F(0xf.1p10), ED())

	assertDecode(t, "c1 -0x1.0", V(1), F(-0x1.0p0), ED())
	assertDecode(t, "c1 -0x1.1", V(1), F(-0x1.1p0), ED())
	assertDecode(t, "c1 -0xf.1p+10", V(1), F(-0xf.1p+10), ED())
	assertDecode(t, "c1 -0xf.1p-10", V(1), F(-0xf.1p-10), ED())
	assertDecode(t, "c1 -0xf.1p10", V(1), F(-0xf.1p10), ED())

	assertDecode(t, "c1 -0x0.0", V(1), F(-0x0.0p0), ED())
	assertDecode(t, "c1 -0x0.1", V(1), F(-0x0.1p0), ED())
	assertDecode(t, "c1 -0x0.1p+10", V(1), F(-0x0.1p+10), ED())
	assertDecode(t, "c1 -0x0.1p-10", V(1), F(-0x0.1p-10), ED())
	assertDecode(t, "c1 -0x0.1p10", V(1), F(-0x0.1p10), ED())

	assertDecode(t, "c1 -0x_0_._1_p_1_0", V(1), F(-0x0.1p10), ED())
}

func TestCTEDate(t *testing.T) {
	assertDecodeEncode(t, "c1 2000-01-01", V(1), CT(compact_time.NewDate(2000, 1, 1)), ED())
	assertDecodeEncode(t, "c1 -2000-12-31", V(1), CT(compact_time.NewDate(-2000, 12, 31)), ED())
}

func TestCTETime(t *testing.T) {
	assertDecode(t, "c1 1:45:00", V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecode(t, "c1 01:45:00", V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertDecodeEncode(t, "c1 23:59:59.101", V(1), CT(compact_time.NewTime(23, 59, 59, 101000000, "")), ED())
	assertDecodeEncode(t, "c1 10:00:01.93/America/Los_Angeles", V(1), CT(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertDecodeEncode(t, "c1 10:00:01.93/89.92/1.10", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ED())
	assertDecode(t, "c1 10:00:01.93/0/0", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ED())
	assertDecode(t, "c1 10:00:01.93/1/1", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ED())
}

func TestCTETimestamp(t *testing.T) {
	assertDecodeEncode(t, "c1 2000-01-01/19:31:44.901554/Z", V(1), CT(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "Z")), ED())
}

func TestCTEURI(t *testing.T) {
	assertDecodeEncode(t, `c1 u"http://example.com"`, V(1), URI("http://example.com"), ED())
	assertEncode(t, `c1 u"http://x.com/%22quoted%22"`, V(1), URI(`http://x.com/"quoted"`), ED())
}

func TestCTEQuotedString(t *testing.T) {
	assertDecodeEncode(t, `c1 "test string"`, V(1), S("test string"), ED())
	assertDecode(t, `c1 "test\nstring"`, V(1), S("test\nstring"), ED())
	assertDecode(t, `c1 "test\rstring"`, V(1), S("test\rstring"), ED())
	assertDecode(t, `c1 "test\tstring"`, V(1), S("test\tstring"), ED())
	assertDecodeEncode(t, `c1 "test\"string"`, V(1), S("test\"string"), ED())
	assertDecode(t, `c1 "test\*string"`, V(1), S("test*string"), ED())
	assertDecode(t, `c1 "test\/string"`, V(1), S("test/string"), ED())
	assertDecodeEncode(t, `c1 "test\\string"`, V(1), S("test\\string"), ED())
	assertDecodeEncode(t, `c1 "test\11string"`, V(1), S("test\u0001string"), ED())
	assertDecodeEncode(t, `c1 "test\4206dstring"`, V(1), S("test\u206dstring"), ED())
	assertDecode(t, `c1 "test\
string"`, V(1), S("teststring"), ED())
	assertDecode(t, "c1 \"test\\\r\nstring\"", V(1), S("teststring"), ED())
}

func TestCTECustomText(t *testing.T) {
	assertDecodeEncode(t, `c1 t"something(123)"`, V(1), CUT("something(123)"), ED())
	assertDecodeEncode(t, `c1 t"some\\thing(\"123\")"`, V(1), CUT("some\\thing(\"123\")"), ED())
	assertDecodeEncode(t, `c1 t"some\nthing\11(123)"`, V(1), CUT("some\nthing\u0001(123)"), ED())
}

func TestCTEList(t *testing.T) {
	assertDecodeEncode(t, `c1 []`, V(1), L(), E(), ED())
	assertDecodeEncode(t, `c1 [123]`, V(1), L(), PI(123), E(), ED())
	assertDecodeEncode(t, `c1 [test]`, V(1), L(), S("test"), E(), ED())
	assertDecodeEncode(t, `c1 [-1 a 2 test -3]`, V(1), L(), NI(1), S("a"), PI(2), S("test"), NI(3), E(), ED())
}

func TestCTEMap(t *testing.T) {
	assertDecodeEncode(t, `c1 {}`, V(1), M(), E(), ED())
	assertDecodeEncode(t, `c1 {1=2}`, V(1), M(), PI(1), PI(2), E(), ED())
	assertDecode(t, "c1 {  1 = 2 3=4 \t}", V(1), M(), PI(1), PI(2), PI(3), PI(4), E(), ED())

	assertDecode(t, `c1 {email = u"mailto:me@somewhere.com" 1.5 = "a string"}`, V(1), M(),
		S("email"), URI("mailto:me@somewhere.com"),
		DF(test.NewDFloat("1.5")), S("a string"),
		E(), ED())
}

func TestCTEMapBadKVSeparator(t *testing.T) {
	assertDecodeFails(t, "c1 {a:b}")
}

func TestCTEListList(t *testing.T) {
	assertDecodeEncode(t, `c1 [[]]`, V(1), L(), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 []]`, V(1), L(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 [] 1]`, V(1), L(), PI(1), L(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, `c1 [1 [2] 1]`, V(1), L(), PI(1), L(), PI(2), E(), PI(1), E(), ED())
}

func TestCTEListMap(t *testing.T) {
	assertDecodeEncode(t, `c1 [{}]`, V(1), L(), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 {}]`, V(1), L(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 [1 {} 1]`, V(1), L(), PI(1), M(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, `c1 [1 {2=3} 1]`, V(1), L(), PI(1), M(), PI(2), PI(3), E(), PI(1), E(), ED())
}

func TestCTEMapList(t *testing.T) {
	assertDecodeEncode(t, `c1 {1=[]}`, V(1), M(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1=[2] test=[1 2 3]}`, V(1), M(), PI(1), L(), PI(2), E(), S("test"), L(), PI(1), PI(2), PI(3), E(), E(), ED())
}

func TestCTEMapMap(t *testing.T) {
	assertDecodeEncode(t, `c1 {1={}}`, V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1={a=b} test={}}`, V(1), M(), PI(1), M(), S("a"), S("b"), E(), S("test"), M(), E(), E(), ED())
}

func TestCTEMetadata(t *testing.T) {
	assertDecodeEncode(t, `c1 ()`, V(1), META(), E(), ED())
	assertDecodeEncode(t, `c1 (1=2)`, V(1), META(), PI(1), PI(2), E(), ED())
	assertDecode(t, "c1 (  1 = 2 3=4 \t)", V(1), META(), PI(1), PI(2), PI(3), PI(4), E(), ED())
}

func TestCTEMarkup(t *testing.T) {
	assertDecodeEncode(t, `c1 <a>`, V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2 3=4>`, V(1), MUP(), S("a"), PI(1), PI(2), PI(3), PI(4), E(), E(), ED())
	assertDecode(t, `c1 <a|>`, V(1), MUP(), S("a"), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a|a>`, V(1), MUP(), S("a"), E(), S("a"), E(), ED())
	assertDecode(t, `c1 <a|a string >`, V(1), MUP(), S("a"), E(), S("a string"), E(), ED())
	assertDecodeEncode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a|a<a>>`, V(1), MUP(), S("a"), E(), S("a"), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecode(t, `c1 <a 1=2 |>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2|a>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2|<a>>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, `c1 <a 1=2|a <a>>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a "), MUP(), S("a"), E(), E(), E(), ED())

	assertDecodeEncode(t, `c1 <a|\\>`, V(1), MUP(), S("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, `c1 <a|\210>`, V(1), MUP(), S("a"), E(), S("\u0010"), E(), ED())
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertDecodeEncode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
}

func TestCTEMarkupComment(t *testing.T) {
	assertDecode(t, "c1 <a|//blah\n>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, "c1 <a|//blah\n a>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S("a"), E(), ED())
	assertDecode(t, "c1 <a|a//blah\n a>", V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), S("a"), E(), ED())

	assertDecode(t, "c1 <a|/*blah*/>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, "c1 <a|a/*blah*/>", V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), E(), ED())
	assertDecode(t, "c1 <a|/*blah*/a>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S("a"), E(), ED())

	assertDecode(t, "c1 <a|/*/*blah*/*/>", V(1), MUP(), S("a"), E(), CMT(), CMT(), S("blah"), E(), E(), E(), ED())
	assertDecode(t, "c1 <a|a/*/*blah*/*/>", V(1), MUP(), S("a"), E(), S("a"), CMT(), CMT(), S("blah"), E(), E(), E(), ED())
	assertDecode(t, "c1 <a|/*/*blah*/*/a>", V(1), MUP(), S("a"), E(), CMT(), CMT(), S("blah"), E(), E(), S("a"), E(), ED())
}

func TestCTEMapMetadata(t *testing.T) {
	assertDecodeEncode(t, `c1 [1 ()a]`, V(1), L(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {1=()a}`, V(1), M(), PI(1), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {1={}}`, V(1), M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {1=(){}}`, V(1), M(), PI(1), META(), E(), M(), E(), E(), ED())

	assertDecodeEncode(t, `c1 {()()1=()()a}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), S("a"), E(), ED())
	assertDecodeEncode(t, `c1 {()()1=()(){}}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), M(), E(), E(), ED())
	assertDecodeEncode(t, `c1 {()()1=()()[]}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), L(), E(), E(), ED())

	assertDecodeEncode(t, `c1 (x=y){(x=y)1=(x=y)(x=y){a=b}}`, V(1),
		META(), S("x"), S("y"), E(), M(),
		META(), S("x"), S("y"), E(), PI(1),
		META(), S("x"), S("y"), E(),
		META(), S("x"), S("y"), E(),
		M(), S("a"), S("b"), E(), E(), ED())
}

func TestCTEBytes(t *testing.T) {
	assertDecodeEncode(t, `c1 b"01"`, V(1), BIN([]byte{0x01}), ED())
	assertDecodeEncode(t, `c1 b"01ff"`, V(1), BIN([]byte{0x01, 0xff}), ED())
	assertDecode(t, `c1 b" f 9 C24 F1 0  "`, V(1), BIN([]byte{0xf9, 0xc2, 0x4f, 0x10}), ED())
}

func TestCTECustom(t *testing.T) {
	assertDecodeEncode(t, `c1 c"a9"`, V(1), CUB([]byte{0xa9}), ED())
	assertDecode(t, "c1 c\"\n8 f 1 9a 7c\td  \"", V(1), CUB([]byte{0x8f, 0x19, 0xa7, 0xcd}), ED())
}

func TestCTEBadArrayType(t *testing.T) {
	assertDecodeFails(t, `c1 x"01"`)
}

func TestCTENamed(t *testing.T) {
	assertDecodeEncode(t, `c1 @nil`, V(1), N(), ED())
	assertDecodeEncode(t, `c1 @nan`, V(1), NAN(), ED())
	assertDecodeEncode(t, `c1 @snan`, V(1), SNAN(), ED())
	assertDecodeEncode(t, `c1 @inf`, V(1), F(math.Inf(1)), ED())
	assertDecodeEncode(t, `c1 -@inf`, V(1), F(math.Inf(-1)), ED())
	assertDecodeEncode(t, `c1 @false`, V(1), FF(), ED())
	assertDecodeEncode(t, `c1 @true`, V(1), TT(), ED())
}

func TestCTEUUID(t *testing.T) {
	assertDecodeEncode(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbb`, V(1),
		UUID([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertDecodeEncode(t, `c1 @00000000-0000-0000-0000-000000000000`, V(1),
		UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), ED())
}

func TestCTEMarker(t *testing.T) {
	assertDecodeFails(t, `c1 &2`)
	assertDecode(t, `c1 &1:string`, V(1), MARK(), PI(1), S("string"), ED())
	assertDecode(t, `c1 &a:string`, V(1), MARK(), S("a"), S("string"), ED())
	assertDecodeFails(t, `c1 & 1:string`)
	assertDecodeFails(t, `c1 &1 string`)
	assertDecodeFails(t, `c1 &1string`)
}

func TestCTEReference(t *testing.T) {
	assertDecode(t, `c1 #2`, V(1), REF(), PI(2), ED())
	assertDecode(t, `c1 #a`, V(1), REF(), S("a"), ED())
	assertDecodeFails(t, `c1 # 1`)
}

func TestCTEMarkerReference(t *testing.T) {
	assertDecode(t, `c1 [&2:testing #2]`, V(1), L(), MARK(), PI(2), S("testing"), REF(), PI(2), E(), ED())
}

func TestCTECommentSingleLine(t *testing.T) {
	assertDecodeFails(t, "c1 //")
	assertDecode(t, "c1 //\n", V(1), CMT(), E(), ED())
	assertDecode(t, "c1 //\r\n", V(1), CMT(), E(), ED())
	assertDecodeFails(t, "c1 // ")
	assertDecode(t, "c1 // \n", V(1), CMT(), S(" "), E(), ED())
	assertDecode(t, "c1 // \r\n", V(1), CMT(), S(" "), E(), ED())
	assertDecodeFails(t, "c1 //a")
	assertDecode(t, "c1 //a\n", V(1), CMT(), S("a"), E(), ED())
	assertDecode(t, "c1 //a\r\n", V(1), CMT(), S("a"), E(), ED())
	assertDecode(t, "c1 // This is a comment\n", V(1), CMT(), S(" This is a comment"), E(), ED())
	assertDecodeFails(t, "c1 /-\n")
}

func TestCTECommentMultiline(t *testing.T) {
	assertDecode(t, "c1 /**/", V(1), CMT(), E(), ED())
	assertDecode(t, "c1 /* */", V(1), CMT(), S(" "), E(), ED())
	assertDecode(t, "c1 /* This is a comment */", V(1), CMT(), S(" This is a comment "), E(), ED())
	assertDecode(t, "c1 /*This is a comment*/", V(1), CMT(), S("This is a comment"), E(), ED())
}

func TestCTECommentMultilineNested(t *testing.T) {
	assertDecode(t, "c1 /*/**/*/", V(1), CMT(), CMT(), E(), E(), ED())
	assertDecode(t, "c1 /*/* */*/", V(1), CMT(), CMT(), S(" "), E(), E(), ED())
	assertDecode(t, "c1 /* /* */ */", V(1), CMT(), S(" "), CMT(), S(" "), E(), S(" "), E(), ED())
	assertDecode(t, "c1  /* before/* mid */ after*/  ", V(1), CMT(), S(" before"), CMT(), S(" mid "), E(), S(" after"), E(), ED())
}

//

func TestMapFloatKey(t *testing.T) {
	assertDecodeEncode(t, "c1 {nil=@nil 1.5=1000}")
}

func TestMarkerReference(t *testing.T) {
	assertDecodeEncode(t, "c1 {first=&1:1000 second=#1}")
}

func TestDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertMarshalUnmarshal(t, v, "c1 [[] [] []]")
}

func TestInf(t *testing.T) {
	assertDecodeEncode(t, `c1 {a=@inf b=1}`)
	assertDecodeEncode(t, `c1 {a=-@inf b=1}`)
}

func TestComment(t *testing.T) {
	// TODO: Better comment formatting
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
	assertDecodeEncode(t, `c1 {a=@inf /*test*/b=1}`)
}

func TestURIReference(t *testing.T) {
	assertDecode(t, `c1
{
    outside_ref      = #u"https://"
    // The markup type is good for presentation data
}
`)
}

func TestMarkupVerbatimString(t *testing.T) {
	assertDecode(t, "c1 <s| `## <d></d>##>")
	assertDecode(t, "c1 <s| `## /d##>")
}

func TestBufferEdge(t *testing.T) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
	assertDecode(t, `c1
{
     1  = <a|
            <b|
               <c| `+"`"+`##                       ##>
                         >
                       >
}
`)
}

func TestBufferEdge2(t *testing.T) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
	assertDecode(t, `c1
{
    x  = <a|
                     <b|
                             <c| `+"`"+`##                     ##>
                           >
                       >
}
`)
}

func TestComplexExample(t *testing.T) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
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
    bytes            = b"10ff389add004f4f91"
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
                                description = "This map will be referenced later using #tag1"
                                value = -@inf
                                child_elements = @nil
                                recursive = #tag1
                            }
    ref1             = #tag1
    ref2             = #tag1
    outside_ref      = #u"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    html_compatible  = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en |
                         <body|
                           Please choose from the following widgets:
                           <div id=parent style=normal ref-id=1 |
                             /* Here we use a backtick to induce verbatim processing.
                              * In this case, "##" is chosen as the ending sequence
                              */
                             <script| `+"`"+`##
                               document.getElementById('parent').insertAdjacentHTML('beforeend',
                                  '<div id="idChild"> content </div>');
                             ##>
                           >
                         >
                       >
}
`)
}

func TestEncodeDecodeExample(t *testing.T) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
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
    bytes = b"10ff389add004f4f91"
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
        description = "This map will be referenced later using #tag1"
        value = -@inf
        child_elements = @nil
        recursive = #tag1
    }
    ref1 = #tag1
    ref2 = #tag1
    outside_ref = #u"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    html_compatible  = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en |
        <body|
            Please choose from the following widgets:
            <div id=parent style=normal ref-id=1 |
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "##" is chosen as the ending sequence
                 */
                <script| ` + "`" + `##
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
    bytes = b"10ff389add004f4f91"
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
        description = "This map will be referenced later using #tag1"
        value = -@inf
        child_elements = @nil
        recursive = #tag1
    }
    ref1 = #tag1
    ref2 = #tag1
    outside_ref = #u"https://somewhere.else.com/path/to/document.cte#some_tag"
    /* The markup type is good for presentation data*/
    html_compatible = <html xmlns=u"http://www.w3.org/1999/xhtml" xml:lang=en|
        <body|
            Please choose from the following widgets: <div id=parent style=normal ref-id=1|
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "##" is chosen as the ending sequence
                 */|
                <script|
                                        document.getElementById('parent').insertAdjacentHTML('beforeend',
                        '\<div id="idChild"\> content \</div\>');
                
                >
            >
        >
    >
}`

	encoded := &bytes.Buffer{}
	encOpts := options.DefaultCTEEncoderOptions()
	encOpts.Indent = "    "
	encoder := NewEncoder(encoded, encOpts)
	decoder := NewDecoder(bytes.NewBuffer(document), encoder, nil)
	err := decoder.Decode()
	if err != nil {
		t.Error(err)
		return
	}

	actual := string(encoded.Bytes())
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
