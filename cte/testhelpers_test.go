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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

var EvV = test.EvV

func NewBigInt(str string) *big.Int {
	return test.NewBigInt(str)
}

func NewBigFloat(str string) *big.Float {
	return test.NewBigFloat(str)
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

func AB(v []bool) test.Event           { return test.AB(v) }
func ACL(l uint64) test.Event          { return test.ACL(l) }
func ACM(l uint64) test.Event          { return test.ACM(l) }
func ADB(v []bool) test.Event          { return test.ADB(v) }
func ADF16(v []float32) test.Event     { return test.ADF16(v) }
func ADF32(v []float32) test.Event     { return test.ADF32(v) }
func ADF64(v []float64) test.Event     { return test.ADF64(v) }
func ADI16(v []int16) test.Event       { return test.ADI16(v) }
func ADI32(v []int32) test.Event       { return test.ADI32(v) }
func ADI64(v []int64) test.Event       { return test.ADI64(v) }
func ADI8(v []int8) test.Event         { return test.ADI8(v) }
func ADT(v string) test.Event          { return test.ADT(v) }
func ADU(v [][]byte) test.Event        { return test.ADU(v) }
func ADU16(v []uint16) test.Event      { return test.ADU16(v) }
func ADU32(v []uint32) test.Event      { return test.ADU32(v) }
func ADU64(v []uint64) test.Event      { return test.ADU64(v) }
func ADU8(v []uint8) test.Event        { return test.ADU8(v) }
func AF16(v []float32) test.Event      { return test.AF16(v) }
func AF32(v []float32) test.Event      { return test.AF32(v) }
func AF64(v []float64) test.Event      { return test.AF64(v) }
func AI16(v []int16) test.Event        { return test.AI16(v) }
func AI32(v []int32) test.Event        { return test.AI32(v) }
func AI64(v []int64) test.Event        { return test.AI64(v) }
func AI8(v []int8) test.Event          { return test.AI8(v) }
func AU(v [][]byte) test.Event         { return test.AU(v) }
func AU16(v []uint16) test.Event       { return test.AU16(v) }
func AU32(v []uint32) test.Event       { return test.AU32(v) }
func AU64(v []uint64) test.Event       { return test.AU64(v) }
func AU8(v []byte) test.Event          { return test.AU8(v) }
func B(v bool) test.Event              { return test.B(v) }
func BAB() test.Event                  { return test.BAB() }
func BAF16() test.Event                { return test.BAF16() }
func BAF32() test.Event                { return test.BAF32() }
func BAF64() test.Event                { return test.BAF64() }
func BAI16() test.Event                { return test.BAI16() }
func BAI32() test.Event                { return test.BAI32() }
func BAI64() test.Event                { return test.BAI64() }
func BAI8() test.Event                 { return test.BAI8() }
func BAU() test.Event                  { return test.BAU() }
func BAU16() test.Event                { return test.BAU16() }
func BAU32() test.Event                { return test.BAU32() }
func BAU64() test.Event                { return test.BAU64() }
func BAU8() test.Event                 { return test.BAU8() }
func BCB() test.Event                  { return test.BCB() }
func BCT() test.Event                  { return test.BCT() }
func BMEDIA() test.Event               { return test.BMEDIA() }
func BREFR() test.Event                { return test.BREFR() }
func BRID() test.Event                 { return test.BRID() }
func BS() test.Event                   { return test.BS() }
func CB(v []byte) test.Event           { return test.CB(v) }
func CM(v string) test.Event           { return test.CM(v) }
func CS(v string) test.Event           { return test.CS(v) }
func CT(v string) test.Event           { return test.CT(v) }
func E() test.Event                    { return test.E() }
func EDGE() test.Event                 { return test.EDGE() }
func L() test.Event                    { return test.L() }
func M() test.Event                    { return test.M() }
func MARK(id string) test.Event        { return test.MARK(id) }
func N(v interface{}) test.Event       { return test.N(v) }
func NAN() test.Event                  { return test.NAN() }
func NODE() test.Event                 { return test.NODE() }
func NULL() test.Event                 { return test.NULL() }
func PAD() test.Event                  { return test.PAD() }
func REFL(id string) test.Event        { return test.REFL(id) }
func REFR(v string) test.Event         { return test.REFR(v) }
func RID(v string) test.Event          { return test.RID(v) }
func S(v string) test.Event            { return test.S(v) }
func SI(id string) test.Event          { return test.SI(id) }
func SNAN() test.Event                 { return test.SNAN() }
func ST(id string) test.Event          { return test.ST(id) }
func T(v compact_time.Time) test.Event { return test.T(v) }
func UID(v []byte) test.Event          { return test.UID(v) }
func V(v uint64) test.Event            { return test.V(v) }

var DebugPrintEvents = false

func decodeToEvents(opts *options.CEDecoderOptions, document []byte) (evts test.Events, err error) {
	var receiver events.DataEventReceiver
	var events *test.EventCollection
	receiver, events = test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	if DebugPrintEvents {
		receiver = test.NewEventPrinter(receiver)
	}
	err = NewDecoder(opts).Decode(bytes.NewBuffer(document), receiver)
	evts = events.Events
	return
}

func encodeEvents(opts *options.CTEEncoderOptions, events ...test.Event) []byte {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(opts)
	encoder.PrepareToEncode(buffer)
	r := rules.NewRules(encoder, nil)
	test.InvokeEventsAsCompleteDocument(r, events...)
	return buffer.Bytes()
}

func assertDecode(t *testing.T, opts *options.CEDecoderOptions, document string, expectedEvents ...test.Event) (successful bool, events test.Events) {
	actualEvents, err := decodeToEvents(opts, []byte(document))
	if err != nil {
		t.Errorf("Error [%v] while decoding document [%v]", err, document)
		return
	}

	if len(expectedEvents) > 0 {
		if !test.AreEventsEquivalent(actualEvents, expectedEvents) {
			t.Errorf("Expected document [%v] to decode to events [%v] but got [%v]", document, test.Events(expectedEvents), actualEvents)
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
		t.Errorf("Expected decode of document [%v] to fail, but got events [%v]", document, events)
	}
}

func assertEncode(t *testing.T, opts *options.CTEEncoderOptions, expectedDocument string, events ...test.Event) (successful bool) {
	actualDocument := string(encodeEvents(opts, events...))
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected events [%v] to encode to document [%v] but got [%v]", events, expectedDocument, actualDocument)
		return
	}
	successful = true
	return
}

func assertDecodeEncode(t *testing.T, decodeOpts *options.CEDecoderOptions,
	encodeOpts *options.CTEEncoderOptions, document string,
	expectedEvents ...test.Event) (successful bool) {
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

func generateCasePermutations(str string) (result []string) {
	casePermutationDFS([]byte(str), 0, &result)
	return
}

func casePermutationDFS(str []byte, index int, result *[]string) {
	if index == len(str) {
		*result = append(*result, string(str))
		return
	}

	casePermutationDFS(str, index+1, result)

	if isASCIIAlpha(str[index]) {
		str[index] ^= (1 << 5)
		casePermutationDFS(str, index+1, result)
	}
}

func isASCIIAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
