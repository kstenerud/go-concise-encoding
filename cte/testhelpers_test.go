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
	"math/big"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-equivalence"
)

func TT() *test.TEvent                       { return test.TT() }
func FF() *test.TEvent                       { return test.FF() }
func I(v int64) *test.TEvent                 { return test.I(v) }
func F(v float64) *test.TEvent               { return test.F(v) }
func BF(v *big.Float) *test.TEvent           { return test.BF(v) }
func DF(v compact_float.DFloat) *test.TEvent { return test.DF(v) }
func BDF(v *apd.Decimal) *test.TEvent        { return test.BDF(v) }
func V(v uint64) *test.TEvent                { return test.V(v) }
func N() *test.TEvent                        { return test.N() }
func PAD(v int) *test.TEvent                 { return test.PAD(v) }
func B(v bool) *test.TEvent                  { return test.B(v) }
func PI(v uint64) *test.TEvent               { return test.PI(v) }
func NI(v uint64) *test.TEvent               { return test.NI(v) }
func BI(v *big.Int) *test.TEvent             { return test.BI(v) }
func CPLX(v complex128) *test.TEvent         { return test.CPLX(v) }
func NAN() *test.TEvent                      { return test.NAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UUID(v []byte) *test.TEvent             { return test.UUID(v) }
func GT(v time.Time) *test.TEvent            { return test.GT(v) }
func CT(v *compact_time.Time) *test.TEvent   { return test.CT(v) }
func BIN(v []byte) *test.TEvent              { return test.BIN(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func URI(v string) *test.TEvent              { return test.URI(v) }
func CUST(v []byte) *test.TEvent             { return test.CUST(v) }
func BB() *test.TEvent                       { return test.BB() }
func SB() *test.TEvent                       { return test.SB() }
func UB() *test.TEvent                       { return test.UB() }
func CB() *test.TEvent                       { return test.CB() }
func AC(l uint64, term bool) *test.TEvent    { return test.AC(l, term) }
func AD(v []byte) *test.TEvent               { return test.AD(v) }
func L() *test.TEvent                        { return test.L() }
func M() *test.TEvent                        { return test.M() }
func MUP() *test.TEvent                      { return test.MUP() }
func META() *test.TEvent                     { return test.META() }
func CMT() *test.TEvent                      { return test.CMT() }
func E() *test.TEvent                        { return test.E() }
func MARK() *test.TEvent                     { return test.MARK() }
func REF() *test.TEvent                      { return test.REF() }
func ED() *test.TEvent                       { return test.ED() }

func decodeToEvents(document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
	err = Decode(bytes.NewBuffer(document), receiver, nil)
	events = receiver.Events
	return
}

func encodeEvents(events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(buffer, nil)
	test.InvokeEvents(encoder, events...)
	return buffer.Bytes()
}

func assertDecode(t *testing.T, document string, expectedEvents ...*test.TEvent) (successful bool, events []*test.TEvent) {
	actualEvents, err := decodeToEvents([]byte(document))
	if err != nil {
		t.Error(err)
		return
	}

	if len(expectedEvents) > 0 {
		if !equivalence.IsEquivalent(actualEvents, expectedEvents) {
			t.Errorf("Expected events %v but got %v", expectedEvents, actualEvents)
			return
		}
	}
	events = actualEvents
	successful = true
	return
}

func assertDecodeFails(t *testing.T, document string) {
	_, err := decodeToEvents([]byte(document))
	if err == nil {
		t.Errorf("Expected decode to fail")
	}
}

func assertEncode(t *testing.T, expectedDocument string, events ...*test.TEvent) (successful bool) {
	actualDocument := string(encodeEvents(events...))
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected document [%v] but got [%v]", expectedDocument, actualDocument)
		return
	}
	successful = true
	return
}

func assertDecodeEncode(t *testing.T, document string, expectedEvents ...*test.TEvent) (successful bool) {
	successful, actualEvents := assertDecode(t, document, expectedEvents...)
	if !successful {
		return
	}
	return assertEncode(t, document, actualEvents...)
}

func assertMarshal(t *testing.T, value interface{}, expectedDocument string) (successful bool) {
	document, err := MarshalToBytes(value, nil)
	if err != nil {
		t.Error(err)
		return
	}
	actualDocument := string(document)
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected document [%v] but got [%v]", expectedDocument, actualDocument)
		return
	}
	successful = true
	return
}

func assertUnmarshal(t *testing.T, expectedValue interface{}, document string) (successful bool) {
	actualValue, err := UnmarshalFromBytes([]byte(document), expectedValue, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actualValue, expectedValue) {
		t.Errorf("Expected unmarshaled [%v] but got [%v]", expectedValue, actualValue)
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
