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
	"math"
	"reflect"
	"testing"

	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-compact-time"
)

// TODO: cteDecode function with recover()

func assertCTEDecode(t *testing.T, document string, expectDecoded ...*test.TEvent) {
	actualDecoded, err := cteDecode([]byte(document))
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectDecoded, actualDecoded) {
		t.Errorf("Expected decode to %v but got %v", expectDecoded, actualDecoded)
	}
}

func assertCTEEncodeDecode(t *testing.T, document string, expectDecoded ...*test.TEvent) {
	actualDecoded, err := cteDecode([]byte(document))
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectDecoded, actualDecoded) {
		t.Errorf("Expected decode to %v but got %v", expectDecoded, actualDecoded)
	}

	expectEncoded := document
	actualEncoded := string(cteEncode(expectDecoded...))
	if actualEncoded != expectEncoded {
		t.Errorf("Expected encode to [%v] but got [%v]", expectEncoded, actualEncoded)
	}
}

func assertCTEDecodeFails(t *testing.T, document string) {
	_, err := cteDecode([]byte(document))
	if err == nil {
		t.Errorf("Expected CTE decode to fail")
	}
}

// Tests

func TestCTEVersion(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 ", V(1), ED())
	assertCTEDecode(t, "\r\n\t c1 ", V(1), ED())
}

func TestCTEVersionNotNumeric(t *testing.T) {
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		document := string([]byte{'c', byte(i)})
		assertCTEDecodeFails(t, document)
	}
}

func TestCTEVersionMissingWhitespace(t *testing.T) {
	assertCTEDecodeFails(t, "c1")
}

func TestCTEBadVersion(t *testing.T) {
	for i := 0; i < 0x100; i++ {
		switch i {
		case 'c', 'C', ' ', '\n', '\r', '\t':
			continue
		default:
			document := string([]byte{byte(i)})
			assertCTEDecodeFails(t, document)
		}
	}
}

func TestCTEWhitespace(t *testing.T) {
	assertCTEDecode(t, "c1     \r\n\t\t\t", V(1), ED())
}

func TestCTEUnquotedString(t *testing.T) {
	// TODO: Encode as well
	assertCTEEncodeDecode(t, "c1 a", V(1), S("a"), ED())
	assertCTEEncodeDecode(t, "c1 abcd", V(1), S("abcd"), ED())
	assertCTEEncodeDecode(t, "c1 _-+.:/123aF", V(1), S("_-+.:/123aF"), ED())
	assertCTEEncodeDecode(t, "c1 新しい", V(1), S("新しい"), ED())
}

func TestCTEVerbatimString(t *testing.T) {
	// TODO: Encode as well
	assertCTEDecodeFails(t, "c1 `")
	assertCTEDecodeFails(t, "c1 `A")
	assertCTEDecodeFails(t, "c1 `A ")
	assertCTEDecodeFails(t, "c1 `A xyz")
	assertCTEDecodeFails(t, "c1 `A xyzAx")
	assertCTEDecode(t, "c1 `A aA", V(1), S("a"), ED())
	assertCTEDecode(t, "c1 `A\taA", V(1), S("a"), ED())
	assertCTEDecode(t, "c1 `A\naA", V(1), S("a"), ED())
	assertCTEDecode(t, "c1 `A\r\naA", V(1), S("a"), ED())
	assertCTEDecode(t, "c1 `#ENDOFSTRING a test\nwith `stuff`#ENDOFSTRING ", V(1), S("a test\nwith `stuff`"), ED())
}

func TestCTEDecimalInt(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 0", V(1), PI(0), ED())
	assertCTEEncodeDecode(t, "c1 123", V(1), PI(123), ED())
	assertCTEEncodeDecode(t, "c1 9412504234235366", V(1), PI(9412504234235366), ED())
	assertCTEEncodeDecode(t, "c1 -49523", V(1), NI(49523), ED())
}

func TestCTEBinaryInt(t *testing.T) {
	assertCTEDecode(t, "c1 0b0", V(1), PI(0), ED())
	assertCTEDecode(t, "c1 0b1", V(1), PI(1), ED())
	assertCTEDecode(t, "c1 0b101", V(1), PI(5), ED())
	assertCTEDecode(t, "c1 0b0010100", V(1), PI(20), ED())
	assertCTEDecode(t, "c1 -0b100", V(1), NI(4), ED())
}

func TestCTEOctalInt(t *testing.T) {
	assertCTEDecode(t, "c1 0o0", V(1), PI(0), ED())
	assertCTEDecode(t, "c1 0o1", V(1), PI(1), ED())
	assertCTEDecode(t, "c1 0o7", V(1), PI(7), ED())
	assertCTEDecode(t, "c1 0o71", V(1), PI(57), ED())
	assertCTEDecode(t, "c1 0o644", V(1), PI(420), ED())
	assertCTEDecode(t, "c1 -0o777", V(1), NI(511), ED())
}

func TestCTEHexInt(t *testing.T) {
	assertCTEDecode(t, "c1 0x0", V(1), PI(0), ED())
	assertCTEDecode(t, "c1 0x1", V(1), PI(1), ED())
	assertCTEDecode(t, "c1 0xf", V(1), PI(0xf), ED())
	assertCTEDecode(t, "c1 0xfedcba9876543210", V(1), PI(0xfedcba9876543210), ED())
	assertCTEDecode(t, "c1 0xFEDCBA9876543210", V(1), PI(0xfedcba9876543210), ED())
	assertCTEDecode(t, "c1 -0x88", V(1), NI(0x88), ED())
}

func TestCTEFloat(t *testing.T) {
	assertCTEDecode(t, "c1 0.0", V(1), DF(test.NewDFloat("0")), ED())
	assertCTEDecode(t, "c1 -0.0", V(1), DF(test.NewDFloat("-0")), ED())

	assertCTEEncodeDecode(t, "c1 1.5", V(1), DF(test.NewDFloat("1.5")), ED())
	assertCTEEncodeDecode(t, "c1 1.125", V(1), DF(test.NewDFloat("1.125")), ED())
	assertCTEEncodeDecode(t, "c1 1.125e+10", V(1), DF(test.NewDFloat("1.125e+10")), ED())
	assertCTEEncodeDecode(t, "c1 1.125e-10", V(1), DF(test.NewDFloat("1.125e-10")), ED())
	assertCTEDecode(t, "c1 1.125e10", V(1), DF(test.NewDFloat("1.125e+10")), ED())

	assertCTEEncodeDecode(t, "c1 -1.5", V(1), DF(test.NewDFloat("-1.5")), ED())
	assertCTEEncodeDecode(t, "c1 -1.125", V(1), DF(test.NewDFloat("-1.125")), ED())
	assertCTEEncodeDecode(t, "c1 -1.125e+10", V(1), DF(test.NewDFloat("-1.125e+10")), ED())
	assertCTEEncodeDecode(t, "c1 -1.125e-10", V(1), DF(test.NewDFloat("-1.125e-10")), ED())
	assertCTEDecode(t, "c1 -1.125e10", V(1), DF(test.NewDFloat("-1.125e10")), ED())

	assertCTEEncodeDecode(t, "c1 0.5", V(1), DF(test.NewDFloat("0.5")), ED())
	assertCTEEncodeDecode(t, "c1 0.125", V(1), DF(test.NewDFloat("0.125")), ED())
	assertCTEDecode(t, "c1 0.125e+10", V(1), DF(test.NewDFloat("0.125e+10")), ED())
	assertCTEDecode(t, "c1 0.125e-10", V(1), DF(test.NewDFloat("0.125e-10")), ED())
	assertCTEDecode(t, "c1 0.125e10", V(1), DF(test.NewDFloat("0.125e10")), ED())

	assertCTEDecode(t, "c1 -0.5", V(1), DF(test.NewDFloat("-0.5")), ED())
	assertCTEDecode(t, "c1 -0.125", V(1), DF(test.NewDFloat("-0.125")), ED())
	assertCTEDecode(t, "c1 -0.125e+10", V(1), DF(test.NewDFloat("-0.125e+10")), ED())
	assertCTEDecode(t, "c1 -0.125e-10", V(1), DF(test.NewDFloat("-0.125e-10")), ED())
	assertCTEDecode(t, "c1 -0.125e10", V(1), DF(test.NewDFloat("-0.125e10")), ED())
}

func TestCTEHexFloat(t *testing.T) {
	assertCTEDecode(t, "c1 0x0.0", V(1), F(0x0.0p0), ED())
	assertCTEDecode(t, "c1 0x0.1", V(1), F(0x0.1p0), ED())
	assertCTEDecode(t, "c1 0x0.1p+10", V(1), F(0x0.1p+10), ED())
	assertCTEDecode(t, "c1 0x0.1p-10", V(1), F(0x0.1p-10), ED())
	assertCTEDecode(t, "c1 0x0.1p10", V(1), F(0x0.1p10), ED())

	assertCTEDecode(t, "c1 0x1.0", V(1), F(0x1.0p0), ED())
	assertCTEDecode(t, "c1 0x1.1", V(1), F(0x1.1p0), ED())
	assertCTEDecode(t, "c1 0xf.1p+10", V(1), F(0xf.1p+10), ED())
	assertCTEDecode(t, "c1 0xf.1p-10", V(1), F(0xf.1p-10), ED())
	assertCTEDecode(t, "c1 0xf.1p10", V(1), F(0xf.1p10), ED())

	assertCTEDecode(t, "c1 -0x1.0", V(1), F(-0x1.0p0), ED())
	assertCTEDecode(t, "c1 -0x1.1", V(1), F(-0x1.1p0), ED())
	assertCTEDecode(t, "c1 -0xf.1p+10", V(1), F(-0xf.1p+10), ED())
	assertCTEDecode(t, "c1 -0xf.1p-10", V(1), F(-0xf.1p-10), ED())
	assertCTEDecode(t, "c1 -0xf.1p10", V(1), F(-0xf.1p10), ED())

	assertCTEDecode(t, "c1 -0x0.0", V(1), F(-0x0.0p0), ED())
	assertCTEDecode(t, "c1 -0x0.1", V(1), F(-0x0.1p0), ED())
	assertCTEDecode(t, "c1 -0x0.1p+10", V(1), F(-0x0.1p+10), ED())
	assertCTEDecode(t, "c1 -0x0.1p-10", V(1), F(-0x0.1p-10), ED())
	assertCTEDecode(t, "c1 -0x0.1p10", V(1), F(-0x0.1p10), ED())
}

func TestCTEDate(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 2000-01-01", V(1), CT(compact_time.NewDate(2000, 1, 1)), ED())
	assertCTEEncodeDecode(t, "c1 -2000-12-31", V(1), CT(compact_time.NewDate(-2000, 12, 31)), ED())
}

func TestCTETime(t *testing.T) {
	assertCTEDecode(t, "c1 1:45:00", V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertCTEEncodeDecode(t, "c1 23:59:59.101", V(1), CT(compact_time.NewTime(23, 59, 59, 101000000, "")), ED())
	assertCTEEncodeDecode(t, "c1 10:00:01.93/America/Los_Angeles", V(1), CT(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertCTEEncodeDecode(t, "c1 10:00:01.93/89.92/1.10", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ED())
	assertCTEDecode(t, "c1 10:00:01.93/0/0", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ED())
	assertCTEDecode(t, "c1 10:00:01.93/1/1", V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ED())
}

func TestCTETimestamp(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 2000-01-01/19:31:44.901554/Z", V(1), CT(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "Z")), ED())
}

func TestCTEURI(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 u"http://example.com"`, V(1), URI("http://example.com"), ED())
}

func TestCTEQuotedString(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 "test string"`, V(1), S("test string"), ED())
}

func TestCTEList(t *testing.T) {
	// TODO: Unquoted string
	assertCTEEncodeDecode(t, `c1 []`, V(1), L(), E(), ED())
	assertCTEEncodeDecode(t, `c1 [123]`, V(1), L(), PI(123), E(), ED())
	assertCTEEncodeDecode(t, `c1 [test]`, V(1), L(), S("test"), E(), ED())
	assertCTEEncodeDecode(t, `c1 [-1 a 2 test -3]`, V(1), L(), NI(1), S("a"), PI(2), S("test"), NI(3), E(), ED())
}

func TestCTEMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {}`, V(1), M(), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1=2}`, V(1), M(), PI(1), PI(2), E(), ED())
	assertCTEDecode(t, "c1 {  1 = 2 3=4 \t}", V(1), M(), PI(1), PI(2), PI(3), PI(4), E(), ED())

	assertCTEDecode(t, `c1 {email = u"mailto:me@somewhere.com" 1.5 = "a string"}`, V(1), M(),
	 S("email"), URI("mailto:me@somewhere.com"),
	DF(test.NewDFloat("1.5")), S("a string"),
	 E(), ED())
}

func TestCTEListList(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [[]]`, V(1), L(), L(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 []]`, V(1), L(), PI(1), L(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 [] 1]`, V(1), L(), PI(1), L(), E(), PI(1), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 [2] 1]`, V(1), L(), PI(1), L(), PI(2), E(), PI(1), E(), ED())
}

func TestCTEListMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [{}]`, V(1), L(), M(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 {}]`, V(1), L(), PI(1), M(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 {} 1]`, V(1), L(), PI(1), M(), E(), PI(1), E(), ED())
	assertCTEEncodeDecode(t, `c1 [1 {2=3} 1]`, V(1), L(), PI(1), M(), PI(2), PI(3), E(), PI(1), E(), ED())
}

func TestCTEMapList(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {1=[]}`, V(1), M(), PI(1), L(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1=[2] test=[1 2 3]}`, V(1), M(), PI(1), L(), PI(2), E(), S("test"), L(), PI(1), PI(2), PI(3), E(), E(), ED())
}

func TestCTEMapMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {1={}}`, V(1), M(), PI(1), M(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1={a=b} test={}}`, V(1), M(), PI(1), M(), S("a"), S("b"), E(), S("test"), M(), E(), E(), ED())
}

func TestCTEMetadata(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 ()`, V(1), META(), E(), ED())
	assertCTEEncodeDecode(t, `c1 (1=2)`, V(1), META(), PI(1), PI(2), E(), ED())
	assertCTEDecode(t, "c1 (  1 = 2 3=4 \t)", V(1), META(), PI(1), PI(2), PI(3), PI(4), E(), ED())
}

func TestCTEMarkup(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 <a>`, V(1), MUP(), S("a"), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a 1=2 3=4>`, V(1), MUP(), S("a"), PI(1), PI(2), PI(3), PI(4), E(), E(), ED())
	assertCTEDecode(t, `c1 <a|>`, V(1), MUP(), S("a"), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a|a>`, V(1), MUP(), S("a"), E(), S("a"), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a|a string >`, V(1), MUP(), S("a"), E(), S("a string "), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a|a<a>>`, V(1), MUP(), S("a"), E(), S("a"), MUP(), S("a"), E(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertCTEDecode(t, `c1 <a 1=2 |>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a 1=2|a>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a"), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a 1=2|<a>>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), MUP(), S("a"), E(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 <a 1=2|a <a>>`, V(1), MUP(), S("a"), PI(1), PI(2), E(), S("a "), MUP(), S("a"), E(), E(), E(), ED())
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 <a|<a>>`, V(1), MUP(), S("a"), E(), MUP(), S("a"), E(), E(), E(), ED())
}

func TestCTEMarkupComment(t *testing.T) {
	assertCTEDecode(t, "c1 <a|//blah\n>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertCTEDecode(t, "c1 <a|//blah\n a>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S(" a"), E(), ED())
	assertCTEDecode(t, "c1 <a|a//blah\n a>", V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), S(" a"), E(), ED())

	assertCTEDecode(t, "c1 <a|/*blah*/>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), E(), ED())
	assertCTEDecode(t, "c1 <a|a/*blah*/>", V(1), MUP(), S("a"), E(), S("a"), CMT(), S("blah"), E(), E(), ED())
	assertCTEDecode(t, "c1 <a|/*blah*/a>", V(1), MUP(), S("a"), E(), CMT(), S("blah"), E(), S("a"), E(), ED())
}

func TestCTEMapMetadata(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [1 ()a]`, V(1), L(), PI(1), META(), E(), S("a"), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1=()a}`, V(1), M(), PI(1), META(), E(), S("a"), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1={}}`, V(1), M(), PI(1), M(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 {1=(){}}`, V(1), M(), PI(1), META(), E(), M(), E(), E(), ED())

	assertCTEEncodeDecode(t, `c1 {()()1=()()a}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), S("a"), E(), ED())
	assertCTEEncodeDecode(t, `c1 {()()1=()(){}}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), M(), E(), E(), ED())
	assertCTEEncodeDecode(t, `c1 {()()1=()()[]}`, V(1), M(), META(), E(), META(), E(), PI(1), META(), E(), META(), E(), L(), E(), E(), ED())

	assertCTEEncodeDecode(t, `c1 (x=y){(x=y)1=(x=y)(x=y){a=b}}`, V(1),
		META(), S("x"), S("y"), E(), M(),
		META(), S("x"), S("y"), E(), PI(1),
		META(), S("x"), S("y"), E(),
		META(), S("x"), S("y"), E(),
		M(), S("a"), S("b"), E(), E(), ED())
}

func TestCTEBytes(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 b"01"`, V(1), BIN([]byte{0x01}), ED())
	assertCTEEncodeDecode(t, `c1 b"01ff"`, V(1), BIN([]byte{0x01, 0xff}), ED())
	assertCTEDecode(t, `c1 b" f 9 C24 F1 0  "`, V(1), BIN([]byte{0xf9, 0xc2, 0x4f, 0x10}), ED())
}

func TestCTECustom(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 c"a9"`, V(1), CUST([]byte{0xa9}), ED())
	assertCTEDecode(t, "c1 c\"\n8 f 1 9a 7c\td  \"", V(1), CUST([]byte{0x8f, 0x19, 0xa7, 0xcd}), ED())
}

func TestCTENamed(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 @nil`, V(1), N(), ED())
	assertCTEEncodeDecode(t, `c1 @nan`, V(1), NAN(), ED())
	assertCTEEncodeDecode(t, `c1 @inf`, V(1), F(math.Inf(1)), ED())
	assertCTEEncodeDecode(t, `c1 -@inf`, V(1), F(math.Inf(-1)), ED())
	assertCTEEncodeDecode(t, `c1 @false`, V(1), FF(), ED())
	assertCTEEncodeDecode(t, `c1 @true`, V(1), TT(), ED())
}

func TestCTEUUID(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbb`, V(1),
		UUID([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ED())
	assertCTEEncodeDecode(t, `c1 @00000000-0000-0000-0000-000000000000`, V(1),
		UUID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), ED())
}

func TestCTEMarker(t *testing.T) {
	assertCTEDecode(t, `c1 &`, V(1), MARK(), ED())
	assertCTEDecode(t, `c1 &1 string`, V(1), MARK(), PI(1), S("string"), ED())
	assertCTEDecode(t, `c1 &a string`, V(1), MARK(), S("a"), S("string"), ED())
	assertCTEDecodeFails(t, `c1 & 1 string`)
}

func TestCTEReference(t *testing.T) {
	assertCTEDecode(t, `c1 #`, V(1), REF(), ED())
	assertCTEDecode(t, `c1 #1 string`, V(1), REF(), PI(1), S("string"), ED())
	assertCTEDecode(t, `c1 #a string`, V(1), REF(), S("a"), S("string"), ED())
	assertCTEDecodeFails(t, `c1 # 1 string`)
}

func TestCTECommentSingleLine(t *testing.T) {
	assertCTEDecode(t, "c1 //", V(1), CMT(), E(), ED())
	assertCTEDecode(t, "c1 //\n", V(1), CMT(), E(), ED())
	assertCTEDecode(t, "c1 //\r\n", V(1), CMT(), E(), ED())
	assertCTEDecode(t, "c1 // ", V(1), CMT(), S(" "), E(), ED())
	assertCTEDecode(t, "c1 // \n", V(1), CMT(), S(" "), E(), ED())
	assertCTEDecode(t, "c1 // \r\n", V(1), CMT(), S(" "), E(), ED())
	assertCTEDecode(t, "c1 //a", V(1), CMT(), S("a"), E(), ED())
	assertCTEDecode(t, "c1 // This is a comment\n", V(1), CMT(), S(" This is a comment"), E(), ED())
}

func TestCTECommentMultiline(t *testing.T) {
	assertCTEDecode(t, "c1 /**/", V(1), CMT(), E(), ED())
	assertCTEDecode(t, "c1 /* */", V(1), CMT(), S(" "), E(), ED())
	assertCTEDecode(t, "c1 /* This is a comment */", V(1), CMT(), S(" This is a comment "), E(), ED())
	assertCTEDecode(t, "c1 /*This is a comment*/", V(1), CMT(), S("This is a comment"), E(), ED())
}

func TestCTECommentMultilineNested(t *testing.T) {
	assertCTEDecode(t, "c1 /*/**/*/", V(1), CMT(), CMT(), E(), E(), ED())
	assertCTEDecode(t, "c1 /*/* */*/", V(1), CMT(), CMT(), S(" "), E(), E(), ED())
	assertCTEDecode(t, "c1 /* /* */ */", V(1), CMT(), S(" "), CMT(), S(" "), E(), S(" "), E(), ED())
	assertCTEDecode(t, "c1  /* before/* mid */ after*/  ", V(1), CMT(), S(" before"), CMT(), S(" mid "), E(), S(" after"), E(), ED())
}
