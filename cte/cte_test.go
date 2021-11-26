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
	"fmt"
	"testing"

	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
)

// TODO: Remove this when releasing V1
func TestCTEVersion1(t *testing.T) {
	assertDecode(t, nil, "c1 1", BD(), EvV, PI(1), ED())
}

func TestCTEDocumentBegin(t *testing.T) {
	// Disallowed version numbers
	for i := 0; i < 0x100; i++ {
		switch i {
		case 'c', 'C', ' ', '\n', '\r', '\t':
			continue
		default:
			document := string([]byte{byte(i), '1', ' ', '1'})
			assertDecodeFails(t, document)
		}
	}
}

func TestCTEVersion(t *testing.T) {
	// Not numeric
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		document := string([]byte{'c', byte(i), ' ', '1'})
		assertDecodeFails(t, document)
	}
}

func TestCTEStringWithNul(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
"test\0string"`, BD(), EvV, S("test\u0000string"), ED())
}

func TestCTEArrayFloat16(t *testing.T) {
	eOpts := options.DefaultCTEEncoderOptions()

	eOpts.DefaultFormats.Array.Float16 = options.CTEEncodingFormatHexadecimal
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x 1.fep+10 -1.3p-40 1.18p+127 1.18p-126|",
		BD(), EvV, AF16([]uint8{0xff, 0x44, 0x98, 0xab, 0x0c, 0x7f, 0x8c, 0x00}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x|", BD(), EvV, AF16([]uint8{}), ED())
	assertDecodeFails(t, "c0\n|f16x -|")

	eOpts.DefaultFormats.Array.Float16 = options.CTEEncodingFormatUnset
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 250 -0.25|",
		BD(), EvV, AF16([]uint8{0x7a, 0x43, 0x80, 0xbe}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16|", BD(), EvV, AF16([]uint8{}), ED())
	assertDecodeFails(t, "c0\n|f16 -|")

	assertDecode(t, nil, "c0 |f16 0.25 0x4.dp-30|",
		BD(), EvV, AF16([]uint8{0x80, 0x3e, 0x9a, 0x31}), ED())

	eOpts.DefaultFormats.Array.Float16 = options.CTEEncodingFormatHexadecimal
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x nan|", BD(), EvV, AF16([]uint8{0xc1, 0xff}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x snan|", BD(), EvV, AF16([]uint8{0x81, 0xff}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x inf|", BD(), EvV, AF16([]uint8{0x80, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x -inf|", BD(), EvV, AF16([]uint8{0x80, 0xff}), ED())
	eOpts.DefaultFormats.Array.Float16 = options.CTEEncodingFormatUnset
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 nan|", BD(), EvV, AF16([]uint8{0xc1, 0xff}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 snan|", BD(), EvV, AF16([]uint8{0x81, 0xff}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 inf|", BD(), EvV, AF16([]uint8{0x80, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 -inf|", BD(), EvV, AF16([]uint8{0x80, 0xff}), ED())

	assertDecodeFails(t, "c0 |f16 0x1.fep+128|")
	assertDecodeFails(t, "c0 |f16 0x1.fep-127|")
	assertDecodeFails(t, "c0 |f16 0x1.fffffffffffffffffffffffff|")
	assertDecodeFails(t, "c0 |f16 -0x1.fffffffffffffffffffffffff|")
}

func TestCTEArrayUID(t *testing.T) {
	// TODO: TestCTEArrayUID
}

func TestCTEArrayBool(t *testing.T) {
	// TODO: TestCTEArrayBool
}

func TestCTEChunked(t *testing.T) {
	assertChunkedStringlike := func(encoded string, startEvent *test.TEvent) {
		assertEncode(t, nil, encoded, BD(), EvV, startEvent, AC(8, false), AD([]byte("abcdefgh")), ED())
		assertEncode(t, nil, encoded, BD(), EvV, startEvent,
			AC(1, true), AD([]byte("a")),
			AC(2, true), AD([]byte("bc")),
			AC(3, true), AD([]byte("def")),
			AC(2, false), AD([]byte("gh")),
			ED())

		assertEncode(t, nil, encoded, BD(), EvV, startEvent,
			AC(1, true), AD([]byte("a")),
			AC(2, true), AD([]byte("bc")),
			AC(3, true), AD([]byte("def")),
			AC(2, true), AD([]byte("gh")),
			AC(0, false), ED())
	}

	assertChunkedByteslike := func(encoded string, startEvent *test.TEvent) {
		assertEncode(t, nil, encoded, BD(), EvV, startEvent, AC(5, false), AD([]byte{0x12, 0x34, 0x56, 0x78, 0x9a}), ED())
		assertEncode(t, nil, encoded, BD(), EvV, startEvent,
			AC(1, true), AD([]byte{0x12}),
			AC(2, true), AD([]byte{0x34, 0x56}),
			AC(2, false), AD([]byte{0x78, 0x9a}),
			ED())
	}

	assertChunkedStringlike(`c0
"abcdefgh"`, SB())
	//TODO: assertChunkedStringlike("c0 `# abcdefgh#", VB())
	assertChunkedStringlike("c0\n@\"abcdefgh\"", RB())
	assertChunkedStringlike("c0\n|ct abcdefgh|", CTB())
	assertChunkedByteslike("c0\n|cb 12 34 56 78 9a|", CBB())
	assertChunkedByteslike("c0\n|u8x 12 34 56 78 9a|", AU8B())
}

func TestCTEList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[]`, BD(), EvV, L(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    123
]`, BD(), EvV, L(), PI(123), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    "test"
]`, BD(), EvV, L(), S("test"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    -1
    "a"
    2
    "test"
    -3
]`, BD(), EvV, L(), NI(1), S("a"), PI(2), S("test"), NI(3), E(), ED())
}

func TestCTEDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertMarshalUnmarshal(t, v, `c0
[
    []
    []
    []
]`)
}

func TestCTEMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
{}`, BD(), EvV, M(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
{
    1 = 2
}`, BD(), EvV, M(), PI(1), PI(2), E(), ED())
	assertDecode(t, nil, "c0 {  1 = 2 3=4 \t}", BD(), EvV, M(), PI(1), PI(2), PI(3), PI(4), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
{
    "null" = null
    1.5 = 1000
}`)

	assertDecode(t, nil, `c0 {"email" = @"mailto:me@somewhere.com" 1.5 = "a string"}`, BD(), EvV, M(),
		S("email"), RID("mailto:me@somewhere.com"),
		DF(NewDFloat("1.5")), S("a string"),
		E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
{
    "a" = inf
    "b" = 1
}`)
	assertDecodeEncode(t, nil, nil, `c0
{
    "a" = -inf
    "b" = 1
}`)
}

func TestCTEMapBadKVSeparator(t *testing.T) {
	assertDecodeFails(t, "c0 {a:b}")
}

func TestCTEListList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    []
]`, BD(), EvV, L(), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    []
]`, BD(), EvV, L(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    []
    1
]`, BD(), EvV, L(), PI(1), L(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    [
        2
    ]
    1
]`, BD(), EvV, L(), PI(1), L(), PI(2), E(), PI(1), E(), ED())
}

func TestCTEListMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    {}
]`, BD(), EvV, L(), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    {}
]`, BD(), EvV, L(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    {}
    1
]`, BD(), EvV, L(), PI(1), M(), E(), PI(1), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    1
    {
        2 = 3
    }
    1
]`, BD(), EvV, L(), PI(1), M(), PI(2), PI(3), E(), PI(1), E(), ED())
}

func TestCTEMapList(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
{
    1 = []
}`, BD(), EvV, M(), PI(1), L(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
{
    1 = [
        2
    ]
    "test" = [
        1
        2
        3
    ]
}`, BD(), EvV, M(), PI(1), L(), PI(2), E(), S("test"), L(), PI(1), PI(2), PI(3), E(), E(), ED())
}

func TestCTEMapMap(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
{
    1 = {}
}`, BD(), EvV, M(), PI(1), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
{
    1 = {
        "a" = "b"
    }
    "test" = {}
}`, BD(), EvV, M(), PI(1), M(), S("a"), S("b"), E(), S("test"), M(), E(), E(), ED())
}

func TestCTEMarkup(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
<a>`, BD(), EvV, MUP("a"), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a 1=2 3=4>`, BD(), EvV, MUP("a"), PI(1), PI(2), PI(3), PI(4), E(), E(), ED())
	assertDecode(t, nil, `c0
<a;>`, BD(), EvV, MUP("a"), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    a
>`, BD(), EvV, MUP("a"), E(), S("a"), E(), ED())
	assertDecode(t, nil, `c0 <a;a string >`, BD(), EvV, MUP("a"), E(), S("a string"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    <a>
>`, BD(), EvV, MUP("a"), E(), MUP("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    <a>
>`, BD(), EvV, MUP("a"), E(), MUP("a"), E(), E(), E(), ED())
	assertDecode(t, nil, `c0 <a 1=2 ;>`, BD(), EvV, MUP("a"), PI(1), PI(2), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a 1=2;
    a
>`, BD(), EvV, MUP("a"), PI(1), PI(2), E(), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a 1=2;
    <a>
>`, BD(), EvV, MUP("a"), PI(1), PI(2), E(), MUP("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a 1=2;
    a 
    <a>
>`, BD(), EvV, MUP("a"), PI(1), PI(2), E(), S("a "), MUP("a"), E(), E(), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    ***
>`, BD(), EvV, MUP("a"), E(), S("***"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    /x
>`, BD(), EvV, MUP("a"), E(), S("/x"), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
<a;
    \\
>`, BD(), EvV, MUP("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    \210
>`, BD(), EvV, MUP("a"), E(), S("\u0010"), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
<a;
    \\
>`, BD(), EvV, MUP("a"), E(), S("\\"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    \<
>`, BD(), EvV, MUP("a"), E(), S("<"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
<a;
    \>
>`, BD(), EvV, MUP("a"), E(), S(">"), E(), ED())
	assertDecode(t, nil, `c0 <a;\*>`, BD(), EvV, MUP("a"), E(), S("*"), E(), ED())
	assertDecode(t, nil, `c0 <a;\/>`, BD(), EvV, MUP("a"), E(), S("/"), E(), ED())

	assertDecodeFails(t, `c0 <a;\y>`)
}

func TestCTEMarkupVerbatimString(t *testing.T) {
	assertDecode(t, nil, `c0 <s; \.## <d></d>##>`)
	assertDecode(t, nil, `c0 <s; \.## /d##>`)
}

func TestCTEMarkupMarkup(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
<a;
    <a>
>`, BD(), EvV, MUP("a"), E(), MUP("a"), E(), E(), E(), ED())
}

func TestCTEMarkupComment(t *testing.T) {
	assertDecode(t, nil, "c0 <a;//blah\n>", BD(), EvV, MUP("a"), E(), COM(false, "blah"), E(), ED())
	assertDecode(t, nil, "c0 <a;//blah\n a>", BD(), EvV, MUP("a"), E(), COM(false, "blah"), S("a"), E(), ED())
	assertDecode(t, nil, "c0 <a;a//blah\n a>", BD(), EvV, MUP("a"), E(), S("a"), COM(false, "blah"), S("a"), E(), ED())

	assertDecode(t, nil, "c0 <a;/*blah*/>", BD(), EvV, MUP("a"), E(), COM(true, "blah"), E(), ED())
	assertDecode(t, nil, "c0 <a;a/*blah*/>", BD(), EvV, MUP("a"), E(), S("a"), COM(true, "blah"), E(), ED())
	assertDecode(t, nil, "c0 <a;/*blah*/a>", BD(), EvV, MUP("a"), E(), COM(true, "blah"), S("a"), E(), ED())

	assertDecode(t, nil, "c0 <a;/*/*blah*/*/>", BD(), EvV, MUP("a"), E(), COM(true, "/*blah*/"), E(), ED())
	assertDecode(t, nil, "c0 <a;a/*/*blah*/*/>", BD(), EvV, MUP("a"), E(), S("a"), COM(true, "/*blah*/"), E(), ED())
	assertDecode(t, nil, "c0 <a;/*/*blah*/*/a>", BD(), EvV, MUP("a"), E(), COM(true, "/*blah*/"), S("a"), E(), ED())

	// TODO: Should it be picking up the extra space between the x and comment?
	assertDecode(t, nil, "c0 <a;x /*blah*/ x>", BD(), EvV, MUP("a"), E(), S("x "), COM(true, "blah"), S("x"), E(), ED())
}

func TestCTENamed(t *testing.T) {
	assertDecodeEncode(t, nil, nil, "c0\nnull", BD(), EvV, NULL(), ED())
	assertDecodeEncode(t, nil, nil, "c0\nnan", BD(), EvV, NAN(), ED())
	assertDecodeEncode(t, nil, nil, "c0\nsnan", BD(), EvV, SNAN(), ED())
	assertDecodeEncode(t, nil, nil, "c0\ninf", BD(), EvV, DF(compact_float.Infinity()), ED())
	assertDecodeEncode(t, nil, nil, "c0\n-inf", BD(), EvV, DF(compact_float.NegativeInfinity()), ED())
	assertDecodeEncode(t, nil, nil, "c0\nfalse", BD(), EvV, FF(), ED())
	assertDecodeEncode(t, nil, nil, "c0\ntrue", BD(), EvV, TT(), ED())
}

func TestCTEMarker(t *testing.T) {
	assertDecodeFails(t, `c0 &2`)
	assertDecode(t, nil, `c0 &1:"string"`, BD(), EvV, MARK("1"), S("string"), ED())
	assertDecode(t, nil, `c0 &a:"string"`, BD(), EvV, MARK("a"), S("string"), ED())
	assertDecodeFails(t, `c0 & 1:"string"`)
	assertDecodeFails(t, `c0 &1 "string"`)
	assertDecodeFails(t, `c0 &1"string"`)
	assertDecodeFails(t, `c0 &rgnsekfrnsekrgfnskergnslekrgnslergselrgblserfbserfbvsekrskfrvbskerfbksefbskerbfserbfrbksuerfbsekjrfbdjfgbsdjfgbsdfgbsdjkhfg`)
}

func TestCTEReference(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    &2:"aaaaa"
    $2
]`, BD(), EvV, L(), MARK("2"), S("aaaaa"), REF("2"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    &a:"aaaaa"
    $a
]`, BD(), EvV, L(), MARK("a"), S("aaaaa"), REF("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
[
    &a:"aaaaa"
    "a"
]`, BD(), EvV, L(), MARK("a"), S("aaaaa"), S("a"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
$"http://x.y"`, BD(), EvV, RREF("http://x.y"), ED())
	assertDecodeFails(t, `c0 $ 1`)
}

func TestCTEMarkerReference(t *testing.T) {
	assertDecode(t, nil, `c0 [&2:"testing" $2]`, BD(), EvV, L(), MARK("2"), S("testing"), REF("2"), E(), ED())
	assertDecodeEncode(t, nil, nil, `c0
{
    "first" = &1:1000
    "second" = $1
}`)
}

func TestCTEMarkerReference2(t *testing.T) {
	assertDecode(t, nil, `c0 {"keys" = &1:"foo" $1 = 1}`,
		BD(), EvV,
		M(),
		S("keys"),
		MARK("1"), S("foo"),
		REF("1"), I(1),
		E(),
		ED())
}

func TestCTEComment(t *testing.T) {
	// TODO: Better comment formatting
	assertDecodeEncode(t, nil, nil, `c0
{
    "a" = inf
    /* test */
    "b" = 1
}`)
}

func TestCTECommentSingleLine(t *testing.T) {
	assertDecodeFails(t, "c0 //")
	assertDecode(t, nil, "c0 //\n1", BD(), EvV, COM(false, ""), PI(1), ED())
	assertDecode(t, nil, "c0 //\r\n1", BD(), EvV, COM(false, ""), PI(1), ED())
	assertDecodeFails(t, "c0 // ")
	assertDecode(t, nil, "c0 // \n1", BD(), EvV, COM(false, " "), PI(1), ED())
	assertDecode(t, nil, "c0 // \r\n1", BD(), EvV, COM(false, " "), PI(1), ED())
	assertDecodeFails(t, "c0 //a")
	assertDecode(t, nil, "c0 //a\n1", BD(), EvV, COM(false, "a"), PI(1), ED())
	assertDecode(t, nil, "c0 //a\r\n1", BD(), EvV, COM(false, "a"), PI(1), ED())
	assertDecode(t, nil, "c0 // This is a comment\n1", BD(), EvV, COM(false, " This is a comment"), PI(1), ED())
	assertDecodeFails(t, "c0 /-\n")
}

func TestCTECommentMultiline(t *testing.T) {
	assertDecode(t, nil, "c0 /**/ 1", BD(), EvV, COM(true, ""), PI(1), ED())
	assertDecode(t, nil, "c0 /**/ 1", BD(), EvV, COM(true, ""), PI(1), ED())
	assertDecode(t, nil, "c0 /* This is a comment */ 1", BD(), EvV, COM(true, " This is a comment "), PI(1), ED())
	assertDecode(t, nil, "c0 /*This is a comment*/ 1", BD(), EvV, COM(true, "This is a comment"), PI(1), ED())
}

func TestCTECommentMultilineNested(t *testing.T) {
	assertDecode(t, nil, "c0 /*/**/*/ 1", BD(), EvV, COM(true, "/**/"), PI(1), ED())
	assertDecode(t, nil, "c0 /*/**/ */ 1", BD(), EvV, COM(true, "/**/ "), PI(1), ED())
	assertDecode(t, nil, "c0 /* /**/ */ 1", BD(), EvV, COM(true, " /**/ "), PI(1), ED())
	assertDecode(t, nil, "c0  /* before/* mid */ after*/ 1  ", BD(), EvV, COM(true, " before/* mid */ after"), PI(1), ED())
	assertDecode(t, nil, "c0 /* x /* y */ 10 */ 5", BD(), EvV, COM(true, " x /* y */ 10 "), PI(5), ED())
	assertDecode(t, nil, "c0 /* x /* y */ na */ 5", BD(), EvV, COM(true, " x /* y */ na "), PI(5), ED())
}

func TestCTECommentAfterValue(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    "a"
    /**/
]`, BD(), EvV, L(), S("a"), COM(true, ""), E(), ED())
}

func TestCTEComplexComment(t *testing.T) {
	document := []byte(`c0
/**/ { /**/ "a"= /**/ "b" /**/ "c"= /**/
<a;
    /**/
    <b>
>}`)

	expected := `c0
/**/
{
    /**/
    "a" = /**/
    "b"
    /**/
    "c" = /**/
    <a;
        /**/
        <b>
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

	actual := encoded.String()
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestCTECommentFollowing(t *testing.T) {
	assertDecode(t, nil, `c0 {"a"="b" /**/}`, BD(), EvV, M(), S("a"), S("b"), COM(true, ""), E(), ED())
	assertDecode(t, nil, `c0 {"a"=2 /**/}`, BD(), EvV, M(), S("a"), PI(2), COM(true, ""), E(), ED())
	assertDecode(t, nil, `c0 {"a"=-2 /**/}`, BD(), EvV, M(), S("a"), NI(2), COM(true, ""), E(), ED())
	// TODO: All other bare values: float, date/time, etc
	assertDecode(t, nil, `c0 {"a"=1.5 /**/}`, BD(), EvV, M(), S("a"), DF(NewDFloat("1.5")), COM(true, ""), E(), ED())
	// TODO: Also test for //
}

func TestCTECommentPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
{
    "a" = "b"
    /**/
}`, BD(), EvV, M(), S("a"), S("b"), COM(true, ""), E(), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
/**/
1`, BD(), EvV, COM(true, ""), PI(1), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
/* a */
1`, BD(), EvV, COM(true, " a "), PI(1), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
/* /**/ */
1`, BD(), EvV, COM(true, " /**/ "), PI(1), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
/* /* a */ */
1`, BD(), EvV, COM(true, " /* a */ "), PI(1), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
/**/
"a"`, BD(), EvV, COM(true, ""), S("a"), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    /* xyz */
    "a"
]`, BD(), EvV, L(), COM(true, " xyz "), S("a"), E(), ED())
}

func TestCTEMarkupPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
<a>`, BD(), EvV, MUP("a"), E(), E(), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
<a "x"=1>`, BD(), EvV, MUP("a"), S("x"), PI(1), E(), E(), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
<a;
    aaa
>`, BD(), EvV, MUP("a"), E(), S("aaa"), E(), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
<a "x"="y";
    aaa
>`, BD(), EvV, MUP("a"), S("x"), S("y"), E(), S("aaa"), E(), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
<a "x"="y" "z"=1;
    aaa
>`, BD(), EvV, MUP("a"), S("x"), S("y"), S("z"), PI(1), E(), S("aaa"), E(), ED())
}

func TestCTEPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, "c0\n1", BD(), EvV, PI(1), ED())
}

func TestCTEListPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	// Empty 1 level
	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[]`, BD(), EvV, L(), E(), ED())

	// Empty 2 level
	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    []
]`, BD(), EvV, L(), L(), E(), E(), ED())
	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    []
    []
]`, BD(), EvV, L(), L(), E(), L(), E(), E(), ED())

	// 1 level
	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    1
    2
]`, BD(), EvV, L(), PI(1), PI(2), E(), ED())

	// 2 level
	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    [
        1
        2
    ]
    [
        3
        4
    ]
]`, BD(), EvV, L(), L(), PI(1), PI(2), E(), L(), PI(3), PI(4), E(), E(), ED())
}

func TestCTEMapPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	// // Empty 1 level
	assertDecodeEncode(t, nil, opts, `c0
{}`, BD(), EvV, M(), E(), ED())

	// Empty 2 level
	assertDecodeEncode(t, nil, opts, `c0
{
    "a" = {}
}`, BD(), EvV, M(), S("a"), M(), E(), E(), ED())
	assertDecodeEncode(t, nil, opts, `c0
{
    "a" = {}
    "b" = {}
}`, BD(), EvV, M(), S("a"), M(), E(), S("b"), M(), E(), E(), ED())

	// 1 level
	assertDecodeEncode(t, nil, opts, `c0
{
    1 = 2
}`, BD(), EvV, M(), PI(1), PI(2), E(), ED())

	// 2 level
	assertDecodeEncode(t, nil, opts, `c0
{
    "a" = {
        1 = 2
        3 = 4
    }
    "b" = {
        5 = 6
        7 = 8
    }
}`, BD(), EvV, M(), S("a"), M(), PI(1), PI(2), PI(3), PI(4), E(), S("b"), M(), PI(5), PI(6), PI(7), PI(8), E(), E(), ED())
}

func TestCTEArrayPretty(t *testing.T) {
	opts := options.DefaultCTEEncoderOptions()

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
|u8x 22 33|`, BD(), EvV, AU8([]uint8{0x22, 0x33}), ED())

	opts.Indent = "    "
	assertDecodeEncode(t, nil, opts, `c0
[
    |u8x 22 33|
    |u8x 66 77|
]`, BD(), EvV, L(), AU8([]uint8{0x22, 0x33}), AU8([]uint8{0x66, 0x77}), E(), ED())
}

func TestCTEMarkupVerbatimPretty(t *testing.T) {
	assertDecode(t, nil, `c0 <blah; \.# aaa #>`,
		BD(), EvV, MUP("blah"), E(), S("aaa"), E(), ED())
}

func TestCTEBufferEdge(t *testing.T) {
	assertDecode(t, nil, `c0
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
	assertDecode(t, nil, `c0
{
    "x"  = <a;
                     <b;
                             <c; `+"`"+`##                     ##>
                           >
                       >
}
`)
}

func TestCTEEncodeDecodeExample(t *testing.T) {
	document := `c0
{
    /* Comments look very C-like, except: /* Nested comments are allowed! */ */
    /* Notice that there are no commas in maps and lists */
    "a_list" = [
        1
        2
        "a string"
    ]
    "map" = {
        2 = "two"
        3 = 3000
        1 = "one"
    }
    "string" = "A string value"
    "boolean" = true
    "regular int" = -10000000
    "decimal float" = -14.125
    "uid" = f1ce4567-e89b-12d3-a456-426655440000
    "date" = 2019-07-01
    "time" = 18:04:00.940231541/E/Prague
    "timestamp" = 2010-07-15/13:28:15.415942344/Z
    "null" = null
    "bytes" = |u8x 10 ff 38 9a dd 00 4f 4f 91|
    "url" = @"https://example.com/"
    "email" = @"mailto:me@somewhere.com"
    1.5 = "Keys don't have to be strings"
    "marked_object" = &tag1:{
        "description" = "This map will be referenced later using $tag1"
        "value" = -inf
        "child_elements" = null
        "recursive" = $tag1
    }
    "ref1" = $tag1
    "ref2" = $tag1
    "outside_ref" = $"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    "html_compatible" = <html "xmlns"=@"http://www.w3.org/1999/xhtml" "xml:lang"="en";
        <body;
            Please choose from the following widgets:
            <div "id"="parent" "style"="normal" "ref-id"=1;
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "#" is chosen as the ending sequence */
            >
        >
    >
}`

	expected := `c0
{
    /* Comments look very C-like, except: /* Nested comments are allowed! */ */
    /* Notice that there are no commas in maps and lists */
    "a_list" = [
        1
        2
        "a string"
    ]
    "map" = {
        2 = "two"
        3 = 3000
        1 = "one"
    }
    "string" = "A string value"
    "boolean" = true
    "regular int" = -10000000
    "decimal float" = -14.125
    "uid" = f1ce4567-e89b-12d3-a456-426655440000
    "date" = 2019-07-01
    "time" = 18:04:00.940231541/Europe/Prague
    "timestamp" = 2010-07-15/13:28:15.415942344
    "null" = null
    "bytes" = |u8x 10 ff 38 9a dd 00 4f 4f 91|
    "url" = @"https://example.com/"
    "email" = @"mailto:me@somewhere.com"
    1.5 = "Keys don't have to be strings"
    "marked_object" = &tag1:{
        "description" = "This map will be referenced later using $tag1"
        "value" = -inf
        "child_elements" = null
        "recursive" = $tag1
    }
    "ref1" = $tag1
    "ref2" = $tag1
    "outside_ref" = $"https://somewhere.else.com/path/to/document.cte#some_tag"
    // The markup type is good for presentation data
    "html_compatible" = <html "xmlns"=@"http://www.w3.org/1999/xhtml" "xml:lang"="en";
        <body;
            Please choose from the following widgets: 
            <div "id"="parent" "style"="normal" "ref-id"=1;
                /* Here we use a backtick to induce verbatim processing.
                 * In this case, "#" is chosen as the ending sequence */
            >
        >
    >
}`
	// TODO: "Please choose from the following widgets:" is getting a space appended to it
	encoded := &bytes.Buffer{}
	encOpts := options.DefaultCTEEncoderOptions()
	encOpts.Indent = "    "
	encoder := NewEncoder(encOpts)
	encoder.PrepareToEncode(encoded)
	decoder := NewDecoder(nil)
	err := decoder.Decode(bytes.NewBuffer([]byte(document)), encoder)
	if err != nil {
		t.Errorf("Error [%v] while decoding %v", err, document)
		return
	}

	actual := encoded.String()
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestMapValueComment(t *testing.T) {
	assertEncode(t, nil, `c0
{
    1 = /**/
    1
}`, BD(), EvV, M(), PI(1), COM(true, ""), PI(1), E(), ED())
}

func TestEmptyDocument(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
null`, BD(), EvV, NULL(), ED())
}

func TestNestedComment(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    /* a /* nested */ comment */
    1
]`, BD(), EvV, L(), COM(true, " a /* nested */ comment "), PI(1), E(), ED())
}

func TestMarkupComment(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
<a;
    /* comment */
    1
>`, BD(), EvV, MUP("a"), E(), COM(true, " comment "), S("1"), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
<a;
    /* comment */
>`, BD(), EvV, MUP("a"), E(), COM(true, " comment "), E(), ED())
}

func TestIdentifier(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
&1a:1`, BD(), EvV, MARK("1a"), I(1), ED())
	assertDecodeEncode(t, nil, nil, `c0
&人気:1`, BD(), EvV, MARK("人気"), I(1), ED())

	assertDecodeFails(t, "c0 &~:1")
	assertDecodeFails(t, "c0 &12345|78:1")
	assertDecodeFails(t, "c0 &12345\u000178:1")
}

func TestCTEEdge(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
@(@"a" @"b" 1)`, BD(), EvV, EDGE(), RID("a"), RID("b"), I(1), ED())

	assertDecodeEncode(t, nil, nil, `c0
{
    true = @(@"a" @"b" 1)
}`, BD(), EvV, M(), TT(), EDGE(), RID("a"), RID("b"), I(1), E(), ED())
}

func TestCTEMedia(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
|application/x-sh 23 21 2f 62 69 6e 2f 73 68 0a 0a 65 63 68 6f 20 68 65 6c 6c 6f 20 77 6f 72 6c 64 0a|`,
		BD(), EvV, MB(), AC(16, false), AD([]byte("application/x-sh")), AC(28, false), AD([]byte(`#!/bin/sh

echo hello world
`)), ED())

	assertDecodeEncode(t, nil, nil, `c0
|a|`, BD(), EvV, MB(), AC(1, false), AD([]byte("a")), AC(0, false), ED())
}

func testDecodeCasePermutations(t *testing.T, name string, events ...*test.TEvent) {
	for _, v := range generateCasePermutations(name) {
		ev := []*test.TEvent{BD(), EvV}
		ev = append(ev, events...)
		ev = append(ev, ED())
		assertDecode(t, nil, fmt.Sprintf("C0 %v", v), ev...)
	}
}

func TestMixedCase(t *testing.T) {
	testDecodeCasePermutations(t, "null", NULL())
	testDecodeCasePermutations(t, "nan", NAN())
	testDecodeCasePermutations(t, "snan", SNAN())
	testDecodeCasePermutations(t, "inf", DF(compact_float.Infinity()))
	testDecodeCasePermutations(t, "-inf", DF(compact_float.NegativeInfinity()))
	testDecodeCasePermutations(t, "false", FF())
	testDecodeCasePermutations(t, "true", TT())
}

func TestSpacing(t *testing.T) {
	assertDecodeFails(t, `c0[]`)
	assertDecodeFails(t, `c0 ["a""b"]`)
	assertDecodeFails(t, `c0 ["a"[]]`)
	assertDecodeFails(t, `c0 [[]"a"]`)
	assertDecodeFails(t, `c0 [[][]]`)
	assertDecodeFails(t, `c0 [{}"a"]`)
	assertDecodeFails(t, `c0 [{}{}]`)
	assertDecodeFails(t, `c0 [<a>"a"]`)
	assertDecodeFails(t, `c0 [<a><a>]`)
	assertDecodeFails(t, `c0 [(@"a" @"a" 1)"a"]`)
	assertDecodeFails(t, `c0 [(@"a" @"a" 1)(@"a" @"a" 1)]`)

	assertDecode(t, nil, `c0 ["a" /* comment */ "b"]`, BD(), EvV, L(), S("a"), COM(true, " comment "), S("b"), E(), ED())

	// TODO: This should not fail
	assertDecodeFails(t, `c0 ["a"/* comment */ "b"]`)
}

func TestMismatchedContainerEnd(t *testing.T) {
	assertDecode(t, nil, `c0 []`, BD(), EvV, L(), E(), ED())
	assertDecodeFails(t, `c0 [}`)
	assertDecodeFails(t, `c0 [>`)
	assertDecodeFails(t, `c0 [)`)

	assertDecode(t, nil, `c0 {}`, BD(), EvV, M(), E(), ED())
	assertDecodeFails(t, `c0 {]`)
	assertDecodeFails(t, `c0 {>`)
	assertDecodeFails(t, `c0 {)`)

	assertDecode(t, nil, `c0 <a>`, BD(), EvV, MUP("a"), E(), E(), ED())
	assertDecodeFails(t, `c0 <a}`)
	assertDecodeFails(t, `c0 <a]`)
	assertDecodeFails(t, `c0 <a)`)

	assertDecode(t, nil, `c0 <a 1=2>`, BD(), EvV, MUP("a"), I(1), I(2), E(), E(), ED())
	assertDecodeFails(t, `c0 <a 1=2}`)
	assertDecodeFails(t, `c0 <a 1=2]`)
	assertDecodeFails(t, `c0 <a 1=2)`)

	assertDecode(t, nil, `c0 <a;a>`, BD(), EvV, MUP("a"), E(), S("a"), E(), ED())
	assertDecodeFails(t, `c0 <a;a}`)
	assertDecodeFails(t, `c0 <a;a]`)
	assertDecodeFails(t, `c0 <a;a)`)

	assertDecode(t, nil, `c0 <a 1=2;a>`, BD(), EvV, MUP("a"), I(1), I(2), E(), S("a"), E(), ED())
	assertDecodeFails(t, `c0 <a 1=2;a}`)
	assertDecodeFails(t, `c0 <a 1=2;a]`)
	assertDecodeFails(t, `c0 <a 1=2;a)`)

	assertDecode(t, nil, `c0 @(@"a" @"a" 1)`, BD(), EvV, EDGE(), RID("a"), RID("a"), I(1), ED())
	assertDecodeFails(t, `c0 @(@"a" @"a" 1]`)
	assertDecodeFails(t, `c0 @(@"a" @"a" 1>`)
	assertDecodeFails(t, `c0 @(@"a" @"a" 1}`)
}

func TestSingleLineCommentAndObject(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
[
    // a comment
    1
]`, BD(), EvV, L(), COM(false, " a comment"), PI(1), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
{
    // a comment
    1 = 2
}`, BD(), EvV, M(), COM(false, " a comment"), PI(1), PI(2), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
<x;
    // a comment
    blah
>`, BD(), EvV, MUP("x"), E(), COM(false, " a comment"), S("blah"), E(), ED())
}

func TestNode(t *testing.T) {
	assertDecodeEncode(t, nil, nil, `c0
("a"
    1
    2
)`, BD(), EvV, NODE(), S("a"), PI(1), PI(2), E(), ED())

	assertDecodeEncode(t, nil, nil, `c0
(null
    1
    (2
    )
    (3
        4
        5
    )
)`, BD(), EvV, NODE(), NULL(),
		PI(1),
		NODE(), PI(2), E(),
		NODE(), PI(3),
		PI(4),
		PI(5),
		E(),
		E(), ED())
}
