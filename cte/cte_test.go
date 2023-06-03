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

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/nullevent"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// TODO
// func TestCurrent(t *testing.T) {
// 	assertDecode(t, `c0
// 	    |.text/plain|

// 	`)
// }

// TODO: Generates incorrect line/column when reporting an error
// func TestBrokenLineColReporting(t *testing.T) {
// 	v := 0.0
// 	v = -v
// 	assertUnmarshal(t, v, `c0
// 0b1@`)
// }

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

// ===========================================================================

func assertMarshal(t *testing.T, value interface{}, expectedDocument string) (successful bool) {
	document, err := NewMarshaler(nil).MarshalToDocument(value)
	if err != nil {
		t.Errorf("Error [%v] while marshaling object %v", err, describe.D(value))
		return
	}
	actualDocument := string(document)
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected marshal of %v to produce document [%v] but got [%v]", describe.D(value), expectedDocument, actualDocument)
		return
	}
	successful = true
	return
}

func assertUnmarshal(t *testing.T, expectedValue interface{}, document string) (successful bool) {
	config := configuration.DefaultCEUnmarshalerConfiguration()
	config.DebugPanics = true
	actualValue, err := NewUnmarshaler(&config).UnmarshalFromDocument([]byte(document), expectedValue)
	if err != nil {
		t.Errorf("Error [%v] while unmarshaling document [%v]", err, document)
		return
	}

	if !equivalence.IsEquivalent(actualValue, expectedValue) {
		t.Errorf("Expected document [%v] to unmarshal to [%v] but got [%v]", document, describe.D(expectedValue), describe.D(actualValue))
		return
	}
	successful = true
	return
}

func assertDecode(t *testing.T, document string) (successful bool) {
	config := configuration.DefaultCEDecoderConfiguration()
	config.DebugPanics = true
	var receiver events.DataEventReceiver
	receiver = nullevent.NewNullEventReceiver()
	receiver = test.NewEventPrinter(receiver)
	receiver = rules.NewRules(receiver, nil)
	err := NewDecoder(&config).DecodeDocument([]byte(document), receiver)

	if err != nil {
		t.Errorf("Error [%v] while decoding document [%v]", err, document)
		return
	}
	successful = true
	return
}

func assertMarshalUnmarshal(t *testing.T, expectedValue interface{}, expectedDocument string) (successful bool) {
	if !assertMarshal(t, expectedValue, expectedDocument) {
		return
	}
	return assertUnmarshal(t, expectedValue, expectedDocument)
}
