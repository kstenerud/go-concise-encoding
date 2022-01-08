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
	"reflect"
	"testing"
)

func TestCBEArrayBitEOF(t *testing.T) {
	assertDecodeFails(t, []byte{header, ceVer, typeArrayBit, 0x12, 0x00})
}

func TestCBEArrayBit(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x02, 0x01}, BD(), EvV, AB(1, []byte{0x01}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x08, 0x03}, BD(), EvV, AB(4, []byte{0x03}), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeArrayBit, 0x26, 0xfe, 0xc1, 0x03}, BD(), EvV, AB(19, []byte{0xfe, 0xc1, 0x03}), ED())
}

func TestCBEMarker(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeMarker, 1, 'x', typeString1, 'a'}, BD(), EvV, MARK("x"), S("a"), ED())
}

func TestCBEReference(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeList, typeMarker, 1, 'x', typeString1, 'a', typeReference, 1, 'x', typeEndContainer},
		BD(), EvV, L(), MARK("x"), S("a"), REF("x"), E(), ED())
}

func TestCBEContainers(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1, typeEndContainer}, BD(), EvV, L(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMap, 1, 1, typeEndContainer}, BD(), EvV, M(), I(1), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeMarkup, 1, 'a', typeEndContainer, typeEndContainer}, BD(), EvV, MUP("a"), E(), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeNode, typeTrue, 1, typeEndContainer}, BD(), EvV, NODE(), TT(), I(1), E(), ED())
	assertDecodeEncode(t, []byte{header, ceVer, typeEdge, 1, 2, 3}, BD(), EvV, EDGE(), I(1), I(2), I(3), ED())

	assertDecodeEncode(t, []byte{header, ceVer, typeList, 1,
		typeList, typeString1, 'a', typeEndContainer,
		typeMap, typeString1, 'a', 100, typeEndContainer,
		typeMarkup, 1, 'a', typeString1, 'a', 50, typeEndContainer, typeString1, 'a', typeEndContainer,
		typeNode, typeTrue, 1, typeEndContainer,
		typeEdge, 1, 2, 3,
		typeEndContainer,
	},
		BD(), EvV, L(), I(1),
		L(), S("a"), E(),
		M(), S("a"), I(100), E(),
		MUP("a"), S("a"), I(50), E(), S("a"), E(),
		NODE(), TT(), I(1), E(),
		EDGE(), I(1), I(2), I(3),
		E(), ED())
}

func TestCBEMultipartArray(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typeString, 0x03, 'a', 0x02, 'b'}, BD(), EvV, SB(), AC(1, true), AD([]byte{'a'}), AC(1, false), AD([]byte{'b'}), ED())
	assertDecode(t, nil, []byte{header, ceVer, typePlane2, typeArrayUint16, 0x03, 0x01, 0x02, 0x02, 0x03, 0x04}, BD(), EvV, AU16B(), AC(1, true), AD([]byte{0x01, 0x02}), AC(1, false), AD([]byte{0x03, 0x04}), ED())
}

func TestCBEChunkedArray(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	InvokeEvents(encoder, EvV, AU16B(), AC(2, true), AD([]byte{1, 0, 2, 0}), AC(2, false), AD([]byte{3, 0, 4, 0}), ED())

	expected := []byte{header, ceVer, typePlane2, typeArrayUint16, 0x05, 1, 0, 2, 0, 0x04, 3, 0, 4, 0}
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

func TestRemoteReference(t *testing.T) {
	assertDecode(t, nil, []byte{header, ceVer, typePlane2, typeRemoteRef, 0x02, 'a'}, BD(), EvV, RREF("a"), ED())
}

func TestMultichunk(t *testing.T) {
	assertDecodeEncode(t, []byte{header, ceVer, typeString, 0x03, 'a', 0}, BD(), EvV, SB(), AC(1, true), AD([]byte{'a'}), AC(0, false), ED())
}

func TestEdge(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer, typeEdge, typeRID, 0x02, 'a', typeRID, 0x02, 'b', 1},
		BD(), EvV, EDGE(), RID("a"), RID("b"), I(1), ED())
}

func TestNode(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer,
			typeNode, typeNull, typeString1, 'a', typeRID, 0x02, 'b',
			typeNode, typeNull, typeEndContainer, typeEndContainer},
		BD(), EvV, NODE(), NULL(), S("a"), RID("b"), NODE(), NULL(), E(), E(), ED())
}

func TestMedia(t *testing.T) {
	assertDecodeEncode(t,
		[]byte{header, ceVer, typePlane2, typeMedia, 0x02, 'a', 0x00},
		BD(), EvV, MB(), AC(1, false), AD([]byte("a")), AC(0, false), ED())
}
