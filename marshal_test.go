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
	"testing"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func assertCBEMarshalUnmarshal(t *testing.T, expected interface{}) {
	options := &CBEMarshalerOptions{
		Iterator: IteratorOptions{
			UseReferences: true,
		},
	}
	document, err := MarshalCBE(expected, options)
	if err != nil {
		t.Errorf("CBE Marshal error: %v", err)
		return
	}

	var actual interface{}
	actual, err = UnmarshalCBE(document, expected, nil)
	if err != nil {
		t.Errorf("CBE Unmarshal error: %v\n- While unmarshaling %v", err, describe.D(document))
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CBE Unmarshal: Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertCTEMarshalUnmarshal(t *testing.T, expected interface{}) {
	options := &CTEMarshalerOptions{
		Iterator: IteratorOptions{
			UseReferences: true,
		},
	}
	document, err := MarshalCTE(expected, options)
	if err != nil {
		t.Errorf("CTE Marshal error: %v", err)
		return
	}

	var actual interface{}
	actual, err = UnmarshalCTE(document, expected, nil)
	if err != nil {
		t.Errorf("CTE Unmarshal error: %v\n- While unmarshaling %v", err, string(document))
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CTE Unmarshal: Expected %v but got %v\n- While unmarshaling %v", describe.D(expected), describe.D(actual), string(document))
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}) {
	assertCBEMarshalUnmarshal(t, expected)
	assertCTEMarshalUnmarshal(t, expected)
}

func TestMarshalUnmarshal(t *testing.T) {
	assertMarshalUnmarshal(t, 101)
	assertMarshalUnmarshal(t, *newTestingOuterStruct(1))
	assertMarshalUnmarshal(t, *newBlankTestingOuterStruct())
}
