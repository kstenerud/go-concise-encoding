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

package concise_encoding

import (
	"math"
	"reflect"
	"testing"

	"github.com/kstenerud/go-compact-time"
)

// TODO: cteDecode function with recover()

func assertCTEDecode(t *testing.T, document string, expectDecoded ...*tevent) {
	actualDecoded, err := cteDecode([]byte(document))
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expectDecoded, actualDecoded) {
		t.Errorf("Expected decode to %v but got %v", expectDecoded, actualDecoded)
	}
}

func assertCTEEncodeDecode(t *testing.T, document string, expectDecoded ...*tevent) {
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
	assertCTEEncodeDecode(t, "c1 ", v(1), ed())
	assertCTEDecode(t, "\r\n\t c1 ", v(1), ed())
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
	assertCTEDecode(t, "c1     \r\n\t\t\t", v(1), ed())
}

func TestCTEUnquotedString(t *testing.T) {
	// TODO: Encode as well
	assertCTEEncodeDecode(t, "c1 a", v(1), s("a"), ed())
	assertCTEEncodeDecode(t, "c1 abcd", v(1), s("abcd"), ed())
	assertCTEEncodeDecode(t, "c1 _-+.:/123aF", v(1), s("_-+.:/123aF"), ed())
	assertCTEEncodeDecode(t, "c1 新しい", v(1), s("新しい"), ed())
}

func TestCTEVerbatimString(t *testing.T) {
	// TODO: Encode as well
	assertCTEDecodeFails(t, "c1 `")
	assertCTEDecodeFails(t, "c1 `A")
	assertCTEDecodeFails(t, "c1 `A ")
	assertCTEDecodeFails(t, "c1 `A xyz")
	assertCTEDecodeFails(t, "c1 `A xyzAx")
	assertCTEDecode(t, "c1 `A aA", v(1), s("a"), ed())
	assertCTEDecode(t, "c1 `A\taA", v(1), s("a"), ed())
	assertCTEDecode(t, "c1 `A\naA", v(1), s("a"), ed())
	assertCTEDecode(t, "c1 `A\r\naA", v(1), s("a"), ed())
	assertCTEDecode(t, "c1 `#ENDOFSTRING a test\nwith `stuff`#ENDOFSTRING ", v(1), s("a test\nwith `stuff`"), ed())
}

func TestCTEDecimalInt(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 0", v(1), pi(0), ed())
	assertCTEEncodeDecode(t, "c1 123", v(1), pi(123), ed())
	assertCTEEncodeDecode(t, "c1 9412504234235366", v(1), pi(9412504234235366), ed())
	assertCTEEncodeDecode(t, "c1 -49523", v(1), ni(49523), ed())
}

func TestCTEBinaryInt(t *testing.T) {
	assertCTEDecode(t, "c1 0b0", v(1), pi(0), ed())
	assertCTEDecode(t, "c1 0b1", v(1), pi(1), ed())
	assertCTEDecode(t, "c1 0b101", v(1), pi(5), ed())
	assertCTEDecode(t, "c1 0b0010100", v(1), pi(20), ed())
	assertCTEDecode(t, "c1 -0b100", v(1), ni(4), ed())
}

func TestCTEOctalInt(t *testing.T) {
	assertCTEDecode(t, "c1 0o0", v(1), pi(0), ed())
	assertCTEDecode(t, "c1 0o1", v(1), pi(1), ed())
	assertCTEDecode(t, "c1 0o7", v(1), pi(7), ed())
	assertCTEDecode(t, "c1 0o71", v(1), pi(57), ed())
	assertCTEDecode(t, "c1 0o644", v(1), pi(420), ed())
	assertCTEDecode(t, "c1 -0o777", v(1), ni(511), ed())
}

func TestCTEHexInt(t *testing.T) {
	assertCTEDecode(t, "c1 0x0", v(1), pi(0), ed())
	assertCTEDecode(t, "c1 0x1", v(1), pi(1), ed())
	assertCTEDecode(t, "c1 0xf", v(1), pi(0xf), ed())
	assertCTEDecode(t, "c1 0xfedcba9876543210", v(1), pi(0xfedcba9876543210), ed())
	assertCTEDecode(t, "c1 0xFEDCBA9876543210", v(1), pi(0xfedcba9876543210), ed())
	assertCTEDecode(t, "c1 -0x88", v(1), ni(0x88), ed())
}

func TestCTEFloat(t *testing.T) {
	assertCTEDecode(t, "c1 0.0", v(1), df("0"), ed())
	assertCTEDecode(t, "c1 -0.0", v(1), df("-0"), ed())

	assertCTEEncodeDecode(t, "c1 1.5", v(1), df("1.5"), ed())
	assertCTEEncodeDecode(t, "c1 1.125", v(1), df("1.125"), ed())
	assertCTEEncodeDecode(t, "c1 1.125e+10", v(1), df("1.125e+10"), ed())
	assertCTEEncodeDecode(t, "c1 1.125e-10", v(1), df("1.125e-10"), ed())
	assertCTEDecode(t, "c1 1.125e10", v(1), df("1.125e+10"), ed())

	assertCTEEncodeDecode(t, "c1 -1.5", v(1), df("-1.5"), ed())
	assertCTEEncodeDecode(t, "c1 -1.125", v(1), df("-1.125"), ed())
	assertCTEEncodeDecode(t, "c1 -1.125e+10", v(1), df("-1.125e+10"), ed())
	assertCTEEncodeDecode(t, "c1 -1.125e-10", v(1), df("-1.125e-10"), ed())
	assertCTEDecode(t, "c1 -1.125e10", v(1), df("-1.125e10"), ed())

	assertCTEEncodeDecode(t, "c1 0.5", v(1), df("0.5"), ed())
	assertCTEEncodeDecode(t, "c1 0.125", v(1), df("0.125"), ed())
	assertCTEDecode(t, "c1 0.125e+10", v(1), df("0.125e+10"), ed())
	assertCTEDecode(t, "c1 0.125e-10", v(1), df("0.125e-10"), ed())
	assertCTEDecode(t, "c1 0.125e10", v(1), df("0.125e10"), ed())

	assertCTEDecode(t, "c1 -0.5", v(1), df("-0.5"), ed())
	assertCTEDecode(t, "c1 -0.125", v(1), df("-0.125"), ed())
	assertCTEDecode(t, "c1 -0.125e+10", v(1), df("-0.125e+10"), ed())
	assertCTEDecode(t, "c1 -0.125e-10", v(1), df("-0.125e-10"), ed())
	assertCTEDecode(t, "c1 -0.125e10", v(1), df("-0.125e10"), ed())
}

func TestCTEHexFloat(t *testing.T) {
	assertCTEDecode(t, "c1 0x0.0", v(1), f(0x0.0p0), ed())
	assertCTEDecode(t, "c1 0x0.1", v(1), f(0x0.1p0), ed())
	assertCTEDecode(t, "c1 0x0.1p+10", v(1), f(0x0.1p+10), ed())
	assertCTEDecode(t, "c1 0x0.1p-10", v(1), f(0x0.1p-10), ed())
	assertCTEDecode(t, "c1 0x0.1p10", v(1), f(0x0.1p10), ed())

	assertCTEDecode(t, "c1 0x1.0", v(1), f(0x1.0p0), ed())
	assertCTEDecode(t, "c1 0x1.1", v(1), f(0x1.1p0), ed())
	assertCTEDecode(t, "c1 0xf.1p+10", v(1), f(0xf.1p+10), ed())
	assertCTEDecode(t, "c1 0xf.1p-10", v(1), f(0xf.1p-10), ed())
	assertCTEDecode(t, "c1 0xf.1p10", v(1), f(0xf.1p10), ed())

	assertCTEDecode(t, "c1 -0x1.0", v(1), f(-0x1.0p0), ed())
	assertCTEDecode(t, "c1 -0x1.1", v(1), f(-0x1.1p0), ed())
	assertCTEDecode(t, "c1 -0xf.1p+10", v(1), f(-0xf.1p+10), ed())
	assertCTEDecode(t, "c1 -0xf.1p-10", v(1), f(-0xf.1p-10), ed())
	assertCTEDecode(t, "c1 -0xf.1p10", v(1), f(-0xf.1p10), ed())

	assertCTEDecode(t, "c1 -0x0.0", v(1), f(-0x0.0p0), ed())
	assertCTEDecode(t, "c1 -0x0.1", v(1), f(-0x0.1p0), ed())
	assertCTEDecode(t, "c1 -0x0.1p+10", v(1), f(-0x0.1p+10), ed())
	assertCTEDecode(t, "c1 -0x0.1p-10", v(1), f(-0x0.1p-10), ed())
	assertCTEDecode(t, "c1 -0x0.1p10", v(1), f(-0x0.1p10), ed())
}

func TestCTEDate(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 2000-01-01", v(1), ct(compact_time.NewDate(2000, 1, 1)), ed())
	assertCTEEncodeDecode(t, "c1 -2000-12-31", v(1), ct(compact_time.NewDate(-2000, 12, 31)), ed())
}

func TestCTETime(t *testing.T) {
	assertCTEDecode(t, "c1 1:45:00", v(1), ct(compact_time.NewTime(1, 45, 0, 0, "")), ed())
	assertCTEEncodeDecode(t, "c1 23:59:59.101", v(1), ct(compact_time.NewTime(23, 59, 59, 101000000, "")), ed())
	assertCTEEncodeDecode(t, "c1 10:00:01.93/America/Los_Angeles", v(1), ct(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ed())
	assertCTEEncodeDecode(t, "c1 10:00:01.93/89.92/1.10", v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ed())
	assertCTEDecode(t, "c1 10:00:01.93/0/0", v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ed())
	assertCTEDecode(t, "c1 10:00:01.93/1/1", v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ed())
}

func TestCTETimestamp(t *testing.T) {
	assertCTEEncodeDecode(t, "c1 2000-01-01/19:31:44.901554/Z", v(1), ct(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "Z")), ed())
}

func TestCTEURI(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 u"http://example.com"`, v(1), uri("http://example.com"), ed())
}

func TestCTEQuotedString(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 "test string"`, v(1), s("test string"), ed())
}

func TestCTEList(t *testing.T) {
	// TODO: Unquoted string
	assertCTEEncodeDecode(t, `c1 []`, v(1), l(), e(), ed())
	assertCTEEncodeDecode(t, `c1 [123]`, v(1), l(), pi(123), e(), ed())
	assertCTEEncodeDecode(t, `c1 [test]`, v(1), l(), s("test"), e(), ed())
	assertCTEEncodeDecode(t, `c1 [-1 a 2 test -3]`, v(1), l(), ni(1), s("a"), pi(2), s("test"), ni(3), e(), ed())
}

func TestCTEMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {}`, v(1), m(), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1=2}`, v(1), m(), pi(1), pi(2), e(), ed())
	assertCTEDecode(t, "c1 {  1 = 2 3=4 \t}", v(1), m(), pi(1), pi(2), pi(3), pi(4), e(), ed())
}

func TestCTEListList(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [[]]`, v(1), l(), l(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 []]`, v(1), l(), pi(1), l(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 [] 1]`, v(1), l(), pi(1), l(), e(), pi(1), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 [2] 1]`, v(1), l(), pi(1), l(), pi(2), e(), pi(1), e(), ed())
}

func TestCTEListMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [{}]`, v(1), l(), m(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 {}]`, v(1), l(), pi(1), m(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 {} 1]`, v(1), l(), pi(1), m(), e(), pi(1), e(), ed())
	assertCTEEncodeDecode(t, `c1 [1 {2=3} 1]`, v(1), l(), pi(1), m(), pi(2), pi(3), e(), pi(1), e(), ed())
}

func TestCTEMapList(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {1=[]}`, v(1), m(), pi(1), l(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1=[2] test=[1 2 3]}`, v(1), m(), pi(1), l(), pi(2), e(), s("test"), l(), pi(1), pi(2), pi(3), e(), e(), ed())
}

func TestCTEMapMap(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 {1={}}`, v(1), m(), pi(1), m(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1={a=b} test={}}`, v(1), m(), pi(1), m(), s("a"), s("b"), e(), s("test"), m(), e(), e(), ed())
}

func TestCTEMetadata(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 ()`, v(1), meta(), e(), ed())
	assertCTEEncodeDecode(t, `c1 (1=2)`, v(1), meta(), pi(1), pi(2), e(), ed())
	assertCTEDecode(t, "c1 (  1 = 2 3=4 \t)", v(1), meta(), pi(1), pi(2), pi(3), pi(4), e(), ed())
}

func TestCTEMarkup(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 <a>`, v(1), mup(), s("a"), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a 1=2 3=4>`, v(1), mup(), s("a"), pi(1), pi(2), pi(3), pi(4), e(), e(), ed())
	assertCTEDecode(t, `c1 <a|>`, v(1), mup(), s("a"), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a|a>`, v(1), mup(), s("a"), e(), s("a"), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a|a string >`, v(1), mup(), s("a"), e(), s("a string "), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a|a<a>>`, v(1), mup(), s("a"), e(), s("a"), mup(), s("a"), e(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a|<a>>`, v(1), mup(), s("a"), e(), mup(), s("a"), e(), e(), e(), ed())
	assertCTEDecode(t, `c1 <a 1=2 |>`, v(1), mup(), s("a"), pi(1), pi(2), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a 1=2|a>`, v(1), mup(), s("a"), pi(1), pi(2), e(), s("a"), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a 1=2|<a>>`, v(1), mup(), s("a"), pi(1), pi(2), e(), mup(), s("a"), e(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 <a 1=2|a <a>>`, v(1), mup(), s("a"), pi(1), pi(2), e(), s("a "), mup(), s("a"), e(), e(), e(), ed())
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 <a|<a>>`, v(1), mup(), s("a"), e(), mup(), s("a"), e(), e(), e(), ed())
}

func TestCTEMarkupComment(t *testing.T) {
	assertCTEDecode(t, "c1 <a|//blah\n>", v(1), mup(), s("a"), e(), cmt(), s("blah"), e(), e(), ed())
	assertCTEDecode(t, "c1 <a|//blah\n a>", v(1), mup(), s("a"), e(), cmt(), s("blah"), e(), s(" a"), e(), ed())
	assertCTEDecode(t, "c1 <a|a//blah\n a>", v(1), mup(), s("a"), e(), s("a"), cmt(), s("blah"), e(), s(" a"), e(), ed())

	assertCTEDecode(t, "c1 <a|/*blah*/>", v(1), mup(), s("a"), e(), cmt(), s("blah"), e(), e(), ed())
	assertCTEDecode(t, "c1 <a|a/*blah*/>", v(1), mup(), s("a"), e(), s("a"), cmt(), s("blah"), e(), e(), ed())
	assertCTEDecode(t, "c1 <a|/*blah*/a>", v(1), mup(), s("a"), e(), cmt(), s("blah"), e(), s("a"), e(), ed())
}

func TestCTEMapMetadata(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 [1 ()a]`, v(1), l(), pi(1), meta(), e(), s("a"), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1=()a}`, v(1), m(), pi(1), meta(), e(), s("a"), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1={}}`, v(1), m(), pi(1), m(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 {1=(){}}`, v(1), m(), pi(1), meta(), e(), m(), e(), e(), ed())

	assertCTEEncodeDecode(t, `c1 {()()1=()()a}`, v(1), m(), meta(), e(), meta(), e(), pi(1), meta(), e(), meta(), e(), s("a"), e(), ed())
	assertCTEEncodeDecode(t, `c1 {()()1=()(){}}`, v(1), m(), meta(), e(), meta(), e(), pi(1), meta(), e(), meta(), e(), m(), e(), e(), ed())
	assertCTEEncodeDecode(t, `c1 {()()1=()()[]}`, v(1), m(), meta(), e(), meta(), e(), pi(1), meta(), e(), meta(), e(), l(), e(), e(), ed())

	assertCTEEncodeDecode(t, `c1 (x=y){(x=y)1=(x=y)(x=y){a=b}}`, v(1),
		meta(), s("x"), s("y"), e(), m(),
		meta(), s("x"), s("y"), e(), pi(1),
		meta(), s("x"), s("y"), e(),
		meta(), s("x"), s("y"), e(),
		m(), s("a"), s("b"), e(), e(), ed())
}

func TestCTEBytes(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 b"01"`, v(1), bin([]byte{0x01}), ed())
	assertCTEEncodeDecode(t, `c1 b"01ff"`, v(1), bin([]byte{0x01, 0xff}), ed())
	assertCTEDecode(t, `c1 b" f 9 C24 F1 0  "`, v(1), bin([]byte{0xf9, 0xc2, 0x4f, 0x10}), ed())
}

func TestCTECustom(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 c"a9"`, v(1), cust([]byte{0xa9}), ed())
	assertCTEDecode(t, "c1 c\"\n8 f 1 9a 7c\td  \"", v(1), cust([]byte{0x8f, 0x19, 0xa7, 0xcd}), ed())
}

func TestCTENamed(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 @nil`, v(1), n(), ed())
	assertCTEEncodeDecode(t, `c1 @nan`, v(1), nan(), ed())
	assertCTEEncodeDecode(t, `c1 @inf`, v(1), f(math.Inf(1)), ed())
	assertCTEEncodeDecode(t, `c1 -@inf`, v(1), f(math.Inf(-1)), ed())
	assertCTEEncodeDecode(t, `c1 @false`, v(1), ff(), ed())
	assertCTEEncodeDecode(t, `c1 @true`, v(1), tt(), ed())
}

func TestCTEUUID(t *testing.T) {
	assertCTEEncodeDecode(t, `c1 @fedcba98-7654-3210-aaaa-bbbbbbbbbbbb`, v(1),
		uuid([]byte{0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10, 0xaa, 0xaa, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}), ed())
	assertCTEEncodeDecode(t, `c1 @00000000-0000-0000-0000-000000000000`, v(1),
		uuid([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), ed())
}

func TestCTEMarker(t *testing.T) {
	assertCTEDecode(t, `c1 &`, v(1), mark(), ed())
	assertCTEDecode(t, `c1 &1 string`, v(1), mark(), pi(1), s("string"), ed())
	assertCTEDecode(t, `c1 &a string`, v(1), mark(), s("a"), s("string"), ed())
	assertCTEDecodeFails(t, `c1 & 1 string`)
}

func TestCTEReference(t *testing.T) {
	assertCTEDecode(t, `c1 #`, v(1), ref(), ed())
	assertCTEDecode(t, `c1 #1 string`, v(1), ref(), pi(1), s("string"), ed())
	assertCTEDecode(t, `c1 #a string`, v(1), ref(), s("a"), s("string"), ed())
	assertCTEDecodeFails(t, `c1 # 1 string`)
}

func TestCTECommentSingleLine(t *testing.T) {
	assertCTEDecode(t, "c1 //", v(1), cmt(), e(), ed())
	assertCTEDecode(t, "c1 //\n", v(1), cmt(), e(), ed())
	assertCTEDecode(t, "c1 //\r\n", v(1), cmt(), e(), ed())
	assertCTEDecode(t, "c1 // ", v(1), cmt(), s(" "), e(), ed())
	assertCTEDecode(t, "c1 // \n", v(1), cmt(), s(" "), e(), ed())
	assertCTEDecode(t, "c1 // \r\n", v(1), cmt(), s(" "), e(), ed())
	assertCTEDecode(t, "c1 //a", v(1), cmt(), s("a"), e(), ed())
	assertCTEDecode(t, "c1 // This is a comment\n", v(1), cmt(), s(" This is a comment"), e(), ed())
}

func TestCTECommentMultiline(t *testing.T) {
	assertCTEDecode(t, "c1 /**/", v(1), cmt(), e(), ed())
	assertCTEDecode(t, "c1 /* */", v(1), cmt(), s(" "), e(), ed())
	assertCTEDecode(t, "c1 /* This is a comment */", v(1), cmt(), s(" This is a comment "), e(), ed())
	assertCTEDecode(t, "c1 /*This is a comment*/", v(1), cmt(), s("This is a comment"), e(), ed())
}

func TestCTECommentMultilineNested(t *testing.T) {
	assertCTEDecode(t, "c1 /*/**/*/", v(1), cmt(), cmt(), e(), e(), ed())
	assertCTEDecode(t, "c1 /*/* */*/", v(1), cmt(), cmt(), s(" "), e(), e(), ed())
	assertCTEDecode(t, "c1 /* /* */ */", v(1), cmt(), s(" "), cmt(), s(" "), e(), s(" "), e(), ed())
	assertCTEDecode(t, "c1  /* before/* mid */ after*/  ", v(1), cmt(), s(" before"), cmt(), s(" mid "), e(), s(" after"), e(), ed())
}
