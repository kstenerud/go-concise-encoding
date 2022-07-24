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
	"math/big"
	"net/url"
	"testing"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

var EvV = test.EvV

const (
	ceVer = version.ConciseEncodingVersion
)

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

func cbeDecode(opts *options.CEDecoderOptions, document []byte) (evts test.Events, err error) {
	var receiver events.DataEventReceiver
	var events *test.EventCollection
	receiver, events = test.NewEventCollector(nil)
	if DebugPrintEvents {
		receiver = test.NewEventPrinter(receiver)
	}
	r := rules.NewRules(receiver, nil)
	err = ce.NewCBEDecoder(opts).Decode(bytes.NewBuffer(document), r)
	evts = events.Events
	return
}

func cbeEncode(encodeOpts *options.CBEEncoderOptions, evts ...test.Event) []byte {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCBEEncoder(encodeOpts)
	r := rules.NewRules(encoder, nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEventsAsCompleteDocument(r, evts...)
	return buffer.Bytes()
}

func cteDecode(opts *options.CEDecoderOptions, document []byte) (evts test.Events, err error) {
	var receiver events.DataEventReceiver
	var events *test.EventCollection
	receiver, events = test.NewEventCollector(nil)
	if DebugPrintEvents {
		receiver = test.NewEventPrinter(receiver)
	}
	r := rules.NewRules(receiver, nil)
	err = ce.NewCTEDecoder(opts).Decode(bytes.NewBuffer(document), r)
	evts = events.Events
	return
}

func cteEncode(encodeOpts *options.CTEEncoderOptions, events ...test.Event) []byte {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCTEEncoder(encodeOpts)
	r := rules.NewRules(encoder, nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEventsAsCompleteDocument(r, events...)
	return buffer.Bytes()
}

// ============================================================================
// Encode/Decode

func assertEncodeDecodeCBEOpts(t *testing.T,
	encodeOpts *options.CBEEncoderOptions,
	decodeOpts *options.CEDecoderOptions,
	expectedEvents ...test.Event) {

	var document []byte
	var actualEvents test.Events
	var err error

	test.AssertNoPanic(t, fmt.Sprintf("CBE Encode [%v]", test.Events(expectedEvents)), func() {
		document = cbeEncode(encodeOpts, expectedEvents...)
	})

	test.AssertNoPanic(t, fmt.Sprintf("CBE Decode %v", describe.D(document)), func() {
		actualEvents, err = cbeDecode(decodeOpts, document)
	})
	if err != nil {
		t.Errorf("Error [%v] while decoding CBE document %v", err, describe.D(document))
		return
	}

	if !test.AreEventsEquivalent(expectedEvents, actualEvents) {
		t.Errorf("CBE: Expected [%v] but got [%v] while decoding %v", test.Events(expectedEvents), actualEvents, describe.D(document))
	}
}

func assertEncodeDecodeCBE(t *testing.T, expected ...test.Event) {
	assertEncodeDecodeCBEOpts(t, nil, nil, expected...)
}

func assertEncodeDecodeCTEOpts(t *testing.T,
	encodeOpts *options.CTEEncoderOptions,
	decodeOpts *options.CEDecoderOptions,
	expectedEvents ...test.Event) {

	var document []byte
	var actualEvents test.Events
	var err error

	test.AssertNoPanic(t, fmt.Sprintf("CTE Encode [%v]", test.Events(expectedEvents)), func() {
		document = cteEncode(encodeOpts, expectedEvents...)
	})

	test.AssertNoPanic(t, fmt.Sprintf("CTE Decode %v", string(document)), func() {
		actualEvents, err = cteDecode(decodeOpts, document)
	})
	if err != nil {
		t.Errorf("Error [%v] while decoding CTE document: %v", err, string(document))
		return
	}

	if !test.AreEventsEquivalent(expectedEvents, actualEvents) {
		t.Errorf("CTE: Expected [%v] but got [%v] while decoding %v", test.Events(expectedEvents), actualEvents, string(document))
	}
}

func assertEncodeDecodeCTE(t *testing.T, expected ...test.Event) {
	assertEncodeDecodeCTEOpts(t, nil, nil, expected...)
}

func assertEncodeDecodeOpts(t *testing.T,
	cbeEncodeOpts *options.CBEEncoderOptions,
	cbeDecodeOpts *options.CEDecoderOptions,
	cteEncodeOpts *options.CTEEncoderOptions,
	cteDecodeOpts *options.CEDecoderOptions,
	expected ...test.Event) {

	assertEncodeDecodeCBEOpts(t, cbeEncodeOpts, cbeDecodeOpts, expected...)
	assertEncodeDecodeCTEOpts(t, cteEncodeOpts, cteDecodeOpts, test.RemoveEvents(expected, test.Padding...)...)
}

func assertEncodeDecode(t *testing.T, expected ...test.Event) {
	assertEncodeDecodeOpts(t, nil, nil, nil, nil, expected...)
}

func assertDecodeEncode(t *testing.T,
	cteEncodeOpts *options.CTEEncoderOptions,
	cteDecodeOpts *options.CEDecoderOptions,
	cbeEncodeOpts *options.CBEEncoderOptions,
	cbeDecodeOpts *options.CEDecoderOptions,
	cteExpectedDocument string,
	cbeExpectedDocument []byte,
	expectedEvents ...test.Event) {

	var receiver events.DataEventReceiver
	var events *test.EventCollection

	textDecoder := ce.NewCTEDecoder(cteDecodeOpts)
	textEncoder := ce.NewCTEEncoder(cteEncodeOpts)
	textRules := ce.NewRules(textEncoder, nil)
	var cteActualDocument *bytes.Buffer

	binDecoder := ce.NewCBEDecoder(cbeDecodeOpts)
	binEncoder := ce.NewCBEEncoder(cbeEncodeOpts)
	binRules := ce.NewRules(binEncoder, nil)

	cbeActualDocument := &bytes.Buffer{}
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), binRules); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(cbeExpectedDocument, cbeActualDocument.Bytes()) {
		t.Errorf("Expected CBE document %v but got %v", describe.D(cbeExpectedDocument), describe.D(cbeActualDocument.Bytes()))
	}

	receiver, events = test.NewEventCollector(nil)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(receiver, nil)); err != nil {
		t.Error(err)
		return
	}
	if !events.IsEquivalentTo(expectedEvents) {
		t.Errorf("Expected CBE document %v to produce events [%v] but got [%v]", describe.D(cbeActualDocument.Bytes()), test.Events(expectedEvents), events.Events)
	}

	cteActualDocument = &bytes.Buffer{}
	textEncoder.PrepareToEncode(cteActualDocument)
	if err := binDecoder.DecodeDocument(cbeExpectedDocument, textRules); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(cteExpectedDocument, cteActualDocument.String()) {
		t.Errorf("Expected CTE document [%v] but got [%v]", cteExpectedDocument, cteActualDocument.String())
	}

	receiver, events = test.NewEventCollector(nil)
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(receiver, nil)); err != nil {
		t.Error(err)
		return
	}
	if !events.IsEquivalentTo(expectedEvents) {
		t.Errorf("Expected CTE document %v to produce events [%v] but got [%v]", cbeActualDocument, test.Events(expectedEvents), events.Events)
	}
}

// ============================================================================
// Marshal/Unmarshal

func assertCBEMarshalUnmarshal(t *testing.T, expected interface{}) {
	marshalOptions := options.DefaultCBEMarshalerOptions()
	unmarshalOptions := options.DefaultCEUnmarshalerOptions()
	assertCBEMarshalUnmarshalWithOptions(t, &marshalOptions, &unmarshalOptions, expected)
	// marshalOptions.Iterator.RecursionSupport = true
	// assertCBEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
}

func assertCBEMarshalUnmarshalWithOptions(t *testing.T,
	marshalOptions *options.CBEMarshalerOptions,
	unmarshalOptions *options.CEUnmarshalerOptions,
	expected interface{}) {

	buffer := &bytes.Buffer{}
	err := ce.MarshalCBE(expected, buffer, marshalOptions)
	if err != nil {
		t.Errorf("CBE Marshal error: %v", err)
		return
	}
	document := buffer.Bytes()

	var actual interface{}
	actual, err = ce.UnmarshalCBE(buffer, expected, unmarshalOptions)
	if err != nil {
		t.Errorf("CBE Unmarshal error: %v\n- While unmarshaling %v", err, describe.D(document))
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CBE Unmarshal: Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}) {
	assertCBEMarshalUnmarshal(t, expected)
	// assertCTEMarshalUnmarshal(t, expected)
}

func assertEncodeDecodeEventStreams(t *testing.T, eventStreams []test.Events) {
	for _, stream := range eventStreams {
		assertEncodeDecode(t, stream...)
	}
}
