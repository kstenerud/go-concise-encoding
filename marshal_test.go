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
	"fmt"
	"testing"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
)

func TestMarshalUnmarshal(t *testing.T) {
	// assertMarshalUnmarshal(t, 101)
	assertMarshalUnmarshal(t, *test.NewTestingOuterStruct(1))
	// assertMarshalUnmarshal(t, *test.NewBlankTestingOuterStruct())
}

func TestMarshalUnmarshalSmallBuffer(t *testing.T) {
	// assertMarshalUnmarshalWithBufferSize(t, 100, *test.NewBlankTestingOuterStruct())
}

type SomeStruct struct {
	A int
	B string
	C *SomeStruct
}

func TestDemonstrateRecursiveStructInMap(t *testing.T) {
	document := "c0 {my-value = &1:{a=100 b=test c=$1}}"
	template := map[string]*SomeStruct{}
	result, err := ce.UnmarshalCTEFromDocument([]byte(document), template, nil)
	if err != nil {
		t.Error(err)
	}
	v := result.(map[string]*SomeStruct)
	s := v["my-value"]
	// Can't naively print a recursive structure in go (Printf will stack overflow), so we print each piece manually.
	fmt.Printf("A: %v, B: %v, Ptr to C: %p, ptr to s: %p\n", s.A, s.B, s.C, s)
	// Prints: A: 100, B: test, Ptr to C: 0xc0001f4600, ptr to s: 0xc0001f4600

	encodedDocument, err := ce.MarshalCTEToDocument(v, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Re-encoded CTE: %v\n", string(encodedDocument))
	// Prints: Re-encoded CTE: c0 {my-value=&0:{A=100 B=test C=$0}}

	encodedDocument, err = ce.MarshalCBEToDocument(v, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Re-encoded CBE: %v\n", encodedDocument)
	// Prints: Re-encoded CBE: [3 0 121 136 109 121 45 118 97 108 117 101 151 0 121 129 97 100 129 98 132 116 101 115 116 129 99 152 0 123 123]
}

func TestEmptyListWithIndents(t *testing.T) {
	v := []interface{}{}
	opts := options.DefaultCTEMarshalerOptions()
	opts.Encoder.Indent = "    "
	encodedDocument, err := ce.MarshalCTEToDocument(v, opts)
	if err != nil {
		t.Error(err)
	}
	expected := "c0\n[]"
	actual := string(encodedDocument)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
