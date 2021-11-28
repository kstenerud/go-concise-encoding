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
	"testing"

	"github.com/kstenerud/go-concise-encoding/options"
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
	// TODO: AF16() is broken and needs to take a float32 arg instead
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x nan|", BD(), EvV, AF16([]uint8{0xc1, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x snan|", BD(), EvV, AF16([]uint8{0x81, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x inf|", BD(), EvV, AF16([]uint8{0x80, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16x -inf|", BD(), EvV, AF16([]uint8{0x80, 0xff}), ED())
	eOpts.DefaultFormats.Array.Float16 = options.CTEEncodingFormatUnset
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 nan|", BD(), EvV, AF16([]uint8{0xc1, 0x7f}), ED())
	// assertDecodeEncode(t, nil, eOpts, "c0\n|f16 snan|", BD(), EvV, AF16([]uint8{0x81, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 inf|", BD(), EvV, AF16([]uint8{0x80, 0x7f}), ED())
	assertDecodeEncode(t, nil, eOpts, "c0\n|f16 -inf|", BD(), EvV, AF16([]uint8{0x80, 0xff}), ED())

	assertDecodeFails(t, "c0 |f16 0x1.fep+128|")
	assertDecodeFails(t, "c0 |f16 0x1.fep-127|")
	assertDecodeFails(t, "c0 |f16 0x1.fffffffffffffffffffffffff|")
	assertDecodeFails(t, "c0 |f16 -0x1.fffffffffffffffffffffffff|")
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
