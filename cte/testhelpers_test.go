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
	"net/url"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

var EvV = test.EvV

func NewBigInt(str string, base int) *big.Int {
	return test.NewBigInt(str, base)
}

func NewBigFloat(str string, base int, significantDigits int) *big.Float {
	return test.NewBigFloat(str, base, significantDigits)
}

func NewDFloat(str string) compact_float.DFloat {
	return test.NewDFloat(str)
}

func NewBDF(str string) *apd.Decimal {
	return test.NewBDF(str)
}

func NewRID(RIDString string) *url.URL {
	return test.NewRID(RIDString)
}

func NewDate(year, month, day int) compact_time.Time {
	return test.NewDate(year, month, day)
}

func NewTime(hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	return test.NewTime(hour, minute, second, nanosecond, areaLocation)
}

func NewTimeLL(hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	return test.NewTimeLL(hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

func NewTimeOff(hour, minute, second, nanosecond, minutesOffset int) compact_time.Time {
	return test.NewTimeOff(hour, minute, second, nanosecond, minutesOffset)
}
func NewTS(year, month, day, hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	return test.NewTS(year, month, day, hour, minute, second, nanosecond, areaLocation)
}

func NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	return test.NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

func NewTSOff(year, month, day, hour, minute, second, nanosecond, minutesOffset int) compact_time.Time {
	return test.NewTSOff(year, month, day, hour, minute, second, nanosecond, minutesOffset)
}

func TT() *test.TEvent                       { return test.TT() }
func FF() *test.TEvent                       { return test.FF() }
func I(v int64) *test.TEvent                 { return test.I(v) }
func F(v float64) *test.TEvent               { return test.F(v) }
func BF(v *big.Float) *test.TEvent           { return test.BF(v) }
func DF(v compact_float.DFloat) *test.TEvent { return test.DF(v) }
func BDF(v *apd.Decimal) *test.TEvent        { return test.BDF(v) }
func NA() *test.TEvent                       { return test.NA() }
func N() *test.TEvent                        { return test.N() }
func PAD(v int) *test.TEvent                 { return test.PAD(v) }
func B(v bool) *test.TEvent                  { return test.B(v) }
func PI(v uint64) *test.TEvent               { return test.PI(v) }
func NI(v uint64) *test.TEvent               { return test.NI(v) }
func BI(v *big.Int) *test.TEvent             { return test.BI(v) }
func NAN() *test.TEvent                      { return test.NAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UUID(v []byte) *test.TEvent             { return test.UUID(v) }
func GT(v time.Time) *test.TEvent            { return test.GT(v) }
func CT(v compact_time.Time) *test.TEvent    { return test.CT(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func RID(v string) *test.TEvent              { return test.RID(v) }
func CUB(v []byte) *test.TEvent              { return test.CUB(v) }
func CUT(v string) *test.TEvent              { return test.CUT(v) }
func AB(l uint64, v []byte) *test.TEvent     { return test.AB(l, v) }
func AU8(v []byte) *test.TEvent              { return test.AU8(v) }
func AU16(v []uint16) *test.TEvent           { return test.AU16(v) }
func AU32(v []uint32) *test.TEvent           { return test.AU32(v) }
func AU64(v []uint64) *test.TEvent           { return test.AU64(v) }
func AI8(v []int8) *test.TEvent              { return test.AI8(v) }
func AI16(v []int16) *test.TEvent            { return test.AI16(v) }
func AI32(v []int32) *test.TEvent            { return test.AI32(v) }
func AI64(v []int64) *test.TEvent            { return test.AI64(v) }
func AF16(v []byte) *test.TEvent             { return test.AF16(v) }
func AF32(v []float32) *test.TEvent          { return test.AF32(v) }
func AF64(v []float64) *test.TEvent          { return test.AF64(v) }
func SB() *test.TEvent                       { return test.SB() }
func RB() *test.TEvent                       { return test.RB() }
func RBCat() *test.TEvent                    { return test.RBCat() }
func CBB() *test.TEvent                      { return test.CBB() }
func CTB() *test.TEvent                      { return test.CTB() }
func ABB() *test.TEvent                      { return test.ABB() }
func AU8B() *test.TEvent                     { return test.AU8B() }
func AU16B() *test.TEvent                    { return test.AU16B() }
func AU32B() *test.TEvent                    { return test.AU32B() }
func AU64B() *test.TEvent                    { return test.AU64B() }
func AI8B() *test.TEvent                     { return test.AI8B() }
func AI16B() *test.TEvent                    { return test.AI16B() }
func AI32B() *test.TEvent                    { return test.AI32B() }
func AI64B() *test.TEvent                    { return test.AI64B() }
func AF16B() *test.TEvent                    { return test.AF16B() }
func AF32B() *test.TEvent                    { return test.AF32B() }
func AF64B() *test.TEvent                    { return test.AF64B() }
func AUUB() *test.TEvent                     { return test.AUUB() }
func AC(l uint64, more bool) *test.TEvent    { return test.AC(l, more) }
func AD(v []byte) *test.TEvent               { return test.AD(v) }
func L() *test.TEvent                        { return test.L() }
func M() *test.TEvent                        { return test.M() }
func MUP(id string) *test.TEvent             { return test.MUP(id) }
func CMT() *test.TEvent                      { return test.CMT() }
func REL() *test.TEvent                      { return test.REL() }
func E() *test.TEvent                        { return test.E() }
func MARK(id string) *test.TEvent            { return test.MARK(id) }
func REF(id string) *test.TEvent             { return test.REF(id) }
func RIDREF() *test.TEvent                   { return test.RIDREF() }
func CONST(n string) *test.TEvent            { return test.CONST(n) }
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

var DebugPrintEvents = false

func decodeToEvents(opts *options.CTEDecoderOptions, document []byte) (evts []*test.TEvent, err error) {
	var receiver events.DataEventReceiver
	ter := test.NewTEventStore()
	receiver = ter
	receiver = rules.NewRules(receiver, nil)
	if DebugPrintEvents {
		receiver = test.NewStdoutTEventPrinter(receiver)
	}
	err = NewDecoder(opts).Decode(bytes.NewBuffer(document), receiver)
	evts = ter.Events
	return
}

func encodeEvents(opts *options.CTEEncoderOptions, events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(opts)
	encoder.PrepareToEncode(buffer)
	r := rules.NewRules(encoder, nil)
	test.InvokeEvents(r, events...)
	return buffer.Bytes()
}

func assertDecode(t *testing.T, opts *options.CTEDecoderOptions, document string, expectedEvents ...*test.TEvent) (successful bool, events []*test.TEvent) {
	actualEvents, err := decodeToEvents(opts, []byte(document))
	if err != nil {
		t.Errorf("Error [%v] while decoding document [%v]", err, document)
		return
	}

	if len(expectedEvents) > 0 {
		if !test.AreAllEventsEqual(actualEvents, expectedEvents) {
			t.Errorf("Expected document [%v] to decode to events %v but got %v", document, expectedEvents, actualEvents)
			return
		}
	}
	events = actualEvents
	successful = true
	return
}

func assertDecodeFails(t *testing.T, document string) {
	events, err := decodeToEvents(nil, []byte(document))
	if err == nil {
		t.Errorf("Expected decode of document [%v] to fail, but got events %v", document, events)
	}
}

func assertEncode(t *testing.T, opts *options.CTEEncoderOptions, expectedDocument string, events ...*test.TEvent) (successful bool) {
	actualDocument := string(encodeEvents(opts, events...))
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected events %v to encode to document [%v] but got [%v]", events, expectedDocument, actualDocument)
		return
	}
	successful = true
	return
}

func assertEncodeFails(t *testing.T, opts *options.CTEEncoderOptions, events ...*test.TEvent) {
	test.AssertPanics(t, "encode", func() {
		encodeEvents(opts, events...)
	})
}

func assertDecodeEncode(t *testing.T, decodeOpts *options.CTEDecoderOptions,
	encodeOpts *options.CTEEncoderOptions, document string,
	expectedEvents ...*test.TEvent) (successful bool) {
	successful, actualEvents := assertDecode(t, decodeOpts, document, expectedEvents...)
	if !successful {
		return
	}
	return assertEncode(t, encodeOpts, document, actualEvents...)
}

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
	actualValue, err := NewUnmarshaler(nil).UnmarshalFromDocument([]byte(document), expectedValue)
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

func assertMarshalUnmarshal(t *testing.T, expectedValue interface{}, expectedDocument string) (successful bool) {
	if !assertMarshal(t, expectedValue, expectedDocument) {
		return
	}
	return assertUnmarshal(t, expectedValue, expectedDocument)
}
