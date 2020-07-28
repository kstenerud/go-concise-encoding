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
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func assertCBEMarshalUnmarshal(t *testing.T, expected interface{}) {
	marshalOptions := options.DefaultCBEMarshalerOptions()
	unmarshalOptions := options.DefaultCBEUnmarshalerOptions()
	assertCBEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
	marshalOptions.Iterator.RecursionSupport = true
	assertCBEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
}

func assertCBEMarshalUnmarshalWithOptions(t *testing.T,
	marshalOptions *options.CBEMarshalerOptions,
	unmarshalOptions *options.CBEUnmarshalerOptions,
	expected interface{}) {

	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
	buffer := &bytes.Buffer{}
	err := cbe.Marshal(expected, buffer, marshalOptions)
	if err != nil {
		t.Errorf("CBE Marshal error: %v", err)
		return
	}
	document := buffer.Bytes()

	var actual interface{}
	actual, err = cbe.Unmarshal(buffer, expected, unmarshalOptions)
	if err != nil {
		t.Errorf("CBE Unmarshal error: %v\n- While unmarshaling %v", err, describe.D(document))
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CBE Unmarshal: Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertCTEMarshalUnmarshal(t *testing.T, expected interface{}) {
	marshalOptions := options.DefaultCTEMarshalerOptions()
	unmarshalOptions := options.DefaultCTEUnmarshalerOptions()
	assertCTEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
	marshalOptions.Iterator.RecursionSupport = true
	assertCTEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
}

func assertCTEMarshalUnmarshalWithOptions(t *testing.T,
	marshalOptions *options.CTEMarshalerOptions,
	unmarshalOptions *options.CTEUnmarshalerOptions,
	expected interface{}) {

	debug.DebugOptions.PassThroughPanics = true
	defer func() {
		debug.DebugOptions.PassThroughPanics = false
	}()
	buffer := &bytes.Buffer{}
	err := cte.Marshal(expected, buffer, marshalOptions)
	if err != nil {
		t.Errorf("CTE Marshal error: %v", err)
		return
	}

	var actual interface{}
	actual, err = cte.Unmarshal(buffer, expected, unmarshalOptions)
	if err != nil {
		t.Errorf("CTE Unmarshal error: %v\n- While unmarshaling %v", err, string(buffer.Bytes()))
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CTE Unmarshal: Expected %v but got %v\n- While unmarshaling %v", describe.D(expected), describe.D(actual), string(buffer.Bytes()))
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}) {
	assertCBEMarshalUnmarshal(t, expected)
	assertCTEMarshalUnmarshal(t, expected)
}

func assertMarshalUnmarshalWithBufferSize(t *testing.T, bufferSize int, expected interface{}) {
	cbeMarshalOptions := options.DefaultCBEMarshalerOptions()
	cbeMarshalOptions.Encoder.BufferSize = bufferSize
	cbeUnmarshalOptions := options.DefaultCBEUnmarshalerOptions()
	cbeUnmarshalOptions.Decoder.BufferSize = bufferSize
	assertCBEMarshalUnmarshalWithOptions(t, cbeMarshalOptions, cbeUnmarshalOptions, expected)

	cteMarshalOptions := options.DefaultCTEMarshalerOptions()
	cteMarshalOptions.Encoder.BufferSize = bufferSize
	cteUnmarshalOptions := options.DefaultCTEUnmarshalerOptions()
	cteUnmarshalOptions.Decoder.BufferSize = bufferSize
	assertCTEMarshalUnmarshalWithOptions(t, cteMarshalOptions, cteUnmarshalOptions, expected)
}

// ============================================================================

// Tests

func TestMarshalUnmarshal(t *testing.T) {
	assertMarshalUnmarshal(t, 101)
	assertMarshalUnmarshal(t, *test.NewTestingOuterStruct(1))
	assertMarshalUnmarshal(t, *test.NewBlankTestingOuterStruct())
}

func TestMarshalUnmarshalSmallBuffer(t *testing.T) {
	assertMarshalUnmarshalWithBufferSize(t, 100, *test.NewBlankTestingOuterStruct())
}

func demonstrateCTEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := cte.Marshal(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CTE: %v\n", string(buffer.Bytes()))
	// Prints: Marshaled CTE: c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}
}

func demonstrateCTEUnmarshal() {
	data := bytes.NewBuffer([]byte(`c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}`))

	value, err := cte.Unmarshal(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CTE: %v\n", value)
	// Prints: Unmarshaled CTE: map[a key:2.5 900:2020-01-15/13:41:00.599]
}

func demonstrateCBEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := cbe.Marshal(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CBE: %v\n", buffer.Bytes())
	// Prints: Marshaled CBE: [1 121 133 97 32 107 101 121 112 0 0 32 64 106 132 3 155 188 18 0 32 109 47 80 0 123]
}

func demonstrateCBEUnmarshal() {
	data := bytes.NewBuffer([]byte{0x01, 0x79, 0x85, 0x61, 0x20, 0x6b, 0x65,
		0x79, 0x70, 0x00, 0x00, 0x20, 0x40, 0x6a, 0x84, 0x03, 0x9b, 0xbc, 0x12,
		0x00, 0x20, 0x6d, 0x2f, 0x50, 0x00, 0x7b})

	value, err := cbe.Unmarshal(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CBE: %v\n", value)
	// Prints: Unmarshaled CBE: map[a key:2.5 900:2020-01-15/13:41:00.599]
}

func TestDemonstrate(t *testing.T) {
	demonstrateCTEMarshal()
	demonstrateCTEUnmarshal()
	demonstrateCBEMarshal()
	demonstrateCBEUnmarshal()
}

type SomeStruct struct {
	A int
	B string
	C *SomeStruct
}

func TestDemonstrateRecursiveStructInMap(t *testing.T) {
	document := "c1 {my-value = &1:{A=100 B=test C=#1}}"
	template := map[string]*SomeStruct{}
	result, err := cte.UnmarshalFromBytes([]byte(document), template, nil)
	if err != nil {
		t.Error(err)
	}
	v := result.(map[string]*SomeStruct)
	s := v["my-value"]
	// Can't naively print a recursive structure in go (Printf will stack overflow), so we print each piece manually.
	fmt.Printf("A: %v, B: %v, Ptr to C: %p, ptr to s: %p\n", s.A, s.B, s.C, s)
	// Prints: A: 100, B: test, Ptr to C: 0xc0001f4600, ptr to s: 0xc0001f4600

	encodedDocument, err := cte.MarshalToBytes(v, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Re-encoded CTE: %v\n", string(encodedDocument))
	// Prints: Re-encoded CTE: c1 {my-value=&0:{A=100 B=test C=#0}}

	encodedDocument, err = cbe.MarshalToBytes(v, nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Re-encoded CBE: %v\n", encodedDocument)
	// Prints: Re-encoded CBE: [1 121 136 109 121 45 118 97 108 117 101 151 0 121 129 65 100 129 66 132 116 101 115 116 129 67 152 0 123 123]
}

func TestEmptyListWithIndents(t *testing.T) {
	v := []interface{}{}
	opts := options.DefaultCTEMarshalerOptions()
	opts.Encoder.Indent = "    "
	encodedDocument, err := cte.MarshalToBytes(v, opts)
	if err != nil {
		t.Error(err)
	}
	expected := "c1 []"
	actual := string(encodedDocument)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
