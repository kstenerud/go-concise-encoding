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

	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func TestCBEEncoderMultiUse(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(configuration.New())
	encoder.PrepareToEncode(buffer)
	test.InvokeEventsAsCompleteDocument(encoder, test.EvV, test.M(), test.E())

	buffer2 := &bytes.Buffer{}
	encoder.PrepareToEncode(buffer2)
	test.InvokeEventsAsCompleteDocument(encoder, test.EvV, test.M(), test.E())

	expected := []byte{0x81, 0x00, 0x99, 0x9b}
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
	assertMarshalUnmarshal(t, v, []byte{0x81, 0x00, 0x9a, 0x9a, 0x9b, 0x9a, 0x9b, 0x9a, 0x9b, 0x9b})
}

func TestCBEUnmarshalPartialBadMap(t *testing.T) {
	// {1=2, 3=}
	document := []byte{0x81, 0x00, 0x99, 0x01, 0x02, 0x03, 0x9b}
	expectedValue := map[interface{}]interface{}{
		1: 2,
	}
	marshaler := NewUnmarshaler(configuration.New())
	actualValue, err := marshaler.UnmarshalFromDocument(document, nil)
	if err == nil {
		t.Errorf("Expected an error")
	}
	if !equivalence.IsEquivalent(expectedValue, actualValue) {
		t.Errorf("Expected %v but got %v", describe.D(expectedValue), describe.D(actualValue))
	}
}

func TestCBEUnmarshalMapTruncatedDocument(t *testing.T) {
	// {1=2, 3=<EOF>
	document := []byte{0x81, 0x00, 0x99, 0x01, 0x02, 0x03}
	expectedValue := map[interface{}]interface{}{
		1: 2,
	}
	marshaler := NewUnmarshaler(configuration.New())
	actualValue, err := marshaler.UnmarshalFromDocument(document, nil)
	if err == nil {
		t.Errorf("Expected an error")
	}
	if !equivalence.IsEquivalent(expectedValue, actualValue) {
		t.Errorf("Expected %v but got %v", describe.D(expectedValue), describe.D(actualValue))
	}
}

func TestCBEUnmarshalListTruncatedDocument(t *testing.T) {
	// [1 2 3 <EOF>
	document := []byte{0x81, 0x00, 0x9a, 0x01, 0x02, 0x03}
	expectedValue := []interface{}{
		1, 2, 3,
	}
	marshaler := NewUnmarshaler(configuration.New())
	actualValue, err := marshaler.UnmarshalFromDocument(document, nil)
	if err == nil {
		t.Errorf("Expected an error")
	}
	if !equivalence.IsEquivalent(expectedValue, actualValue) {
		t.Errorf("Expected %v but got %v", describe.D(expectedValue), describe.D(actualValue))
	}
}

func TestCBEUnmarshalDeepTruncatedDocument(t *testing.T) {
	// [1 2 3 {4=5 6=[7<EOF>
	document := []byte{0x81, 0x00, 0x9a, 0x01, 0x02, 0x03, 0x99, 0x04, 0x05, 0x06, 0x9a, 0x07}
	expectedValue := []interface{}{
		1, 2, 3,
		map[interface{}]interface{}{
			4: 5,
			6: []interface{}{
				7,
			},
		},
	}
	marshaler := NewUnmarshaler(configuration.New())
	actualValue, err := marshaler.UnmarshalFromDocument(document, nil)
	if err == nil {
		t.Errorf("Expected an error")
	}
	if !equivalence.IsEquivalent(expectedValue, actualValue) {
		t.Errorf("Expected %v but got %v", describe.D(expectedValue), describe.D(actualValue))
	}
}

// ===========================================================================

func assertMarshal(t *testing.T, value interface{}, expectedDocument []byte) (successful bool) {
	document, err := NewMarshaler(configuration.New()).MarshalToDocument(value)
	if err != nil {
		t.Errorf("Error [%v] while marshaling %v", err, describe.D(value))
		return
	}
	if !equivalence.IsEquivalent(document, expectedDocument) {
		t.Errorf("Expected encoded document %v but got %v while marshaling %v", describe.D(expectedDocument), describe.D(document), describe.D(value))
		return
	}
	successful = true
	return
}

func assertUnmarshal(t *testing.T, expectedValue interface{}, document []byte) (successful bool) {
	actualValue, err := NewUnmarshaler(configuration.New()).UnmarshalFromDocument([]byte(document), expectedValue)
	if err != nil {
		t.Errorf("Error [%v] while unmarshaling %v", err, describe.D(document))
		return
	}

	if !equivalence.IsEquivalent(actualValue, expectedValue) {
		t.Errorf("Expected unmarshaled [%v] but got [%v] while unmarshaling %v", describe.D(expectedValue), describe.D(actualValue), describe.D(document))
		return
	}
	successful = true
	return
}

func assertMarshalUnmarshal(t *testing.T, expectedValue interface{}, expectedDocument []byte) (successful bool) {
	if !assertMarshal(t, expectedValue, expectedDocument) {
		return
	}
	return assertUnmarshal(t, expectedValue, expectedDocument)
}
