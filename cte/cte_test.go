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
