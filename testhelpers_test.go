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
	"time"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

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

func NewTS(year, month, day, hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	return test.NewTS(year, month, day, hour, minute, second, nanosecond, areaLocation)
}

func NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	return test.NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

func TT() *test.TEvent                       { return test.TT() }
func FF() *test.TEvent                       { return test.FF() }
func I(v int64) *test.TEvent                 { return test.I(v) }
func F(v float64) *test.TEvent               { return test.F(v) }
func BF(v *big.Float) *test.TEvent           { return test.BF(v) }
func DF(v compact_float.DFloat) *test.TEvent { return test.DF(v) }
func BDF(v *apd.Decimal) *test.TEvent        { return test.BDF(v) }
func V(v uint64) *test.TEvent                { return test.V(v) }
func NA() *test.TEvent                       { return test.NA() }
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
func MUP() *test.TEvent                      { return test.MUP() }
func META() *test.TEvent                     { return test.META() }
func CMT() *test.TEvent                      { return test.CMT() }
func E() *test.TEvent                        { return test.E() }
func MARK() *test.TEvent                     { return test.MARK() }
func REF() *test.TEvent                      { return test.REF() }
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

func cbeDecode(opts *options.CBEDecoderOptions, document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTEventStore()
	r := rules.NewRules(receiver, nil)
	err = ce.NewCBEDecoder(opts).Decode(bytes.NewBuffer(document), r)
	events = receiver.Events
	return
}

func cbeEncode(encodeOpts *options.CBEEncoderOptions, events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCBEEncoder(encodeOpts)
	r := rules.NewRules(encoder, nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEvents(r, events...)
	return buffer.Bytes()
}

func cbeEncodeDecode(encodeOpts *options.CBEEncoderOptions,
	decodeOpts *options.CBEDecoderOptions,
	expected ...*test.TEvent) (events []*test.TEvent, err error) {

	return cbeDecode(decodeOpts, cbeEncode(encodeOpts, expected...))
}

func cteDecode(opts *options.CTEDecoderOptions, document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTEventStore()
	r := rules.NewRules(receiver, nil)
	err = ce.NewCTEDecoder(opts).Decode(bytes.NewBuffer(document), r)
	events = receiver.Events
	return
}

func cteEncode(encodeOpts *options.CTEEncoderOptions, events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCTEEncoder(encodeOpts)
	r := rules.NewRules(encoder, nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEvents(r, events...)
	return buffer.Bytes()
}

func cteEncodeDecode(encodeOpts *options.CTEEncoderOptions,
	decodeOpts *options.CTEDecoderOptions,
	events ...*test.TEvent) (decodedEvents []*test.TEvent, err error) {

	events = test.FilterEventsForCTE(events)
	document := cteEncode(encodeOpts, events...)
	return cteDecode(decodeOpts, document)
}

// ============================================================================
// Encode/Decode

func assertEncodeDecodeCBEOpts(t *testing.T,
	encodeOpts *options.CBEEncoderOptions,
	decodeOpts *options.CBEDecoderOptions,
	expectedEvents ...*test.TEvent) {

	var document []byte
	var actualEvents []*test.TEvent
	var err error

	test.AssertNoPanic(t, fmt.Sprintf("CBE Encode %v", expectedEvents), func() {
		document = cbeEncode(encodeOpts, expectedEvents...)
	})

	test.AssertNoPanic(t, fmt.Sprintf("CBE Decode %v", describe.D(document)), func() {
		actualEvents, err = cbeDecode(decodeOpts, document)
	})
	if err != nil {
		t.Errorf("%v while decoding CBE document: %v", err, describe.D(document))
		return
	}

	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		t.Errorf("CBE: Expected %v but got %v while decoding %v", expectedEvents, actualEvents, describe.D(document))
	}
}

func assertEncodeDecodeCBE(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeCBEOpts(t, nil, nil, expected...)
}

func assertEncodeDecodeCTEOpts(t *testing.T,
	encodeOpts *options.CTEEncoderOptions,
	decodeOpts *options.CTEDecoderOptions,
	expectedEvents ...*test.TEvent) {

	var document []byte
	var actualEvents []*test.TEvent
	var err error

	test.AssertNoPanic(t, fmt.Sprintf("CTE Encode %v", expectedEvents), func() {
		document = cteEncode(encodeOpts, expectedEvents...)
	})

	test.AssertNoPanic(t, fmt.Sprintf("CTE Decode %v", string(document)), func() {
		actualEvents, err = cteDecode(decodeOpts, document)
	})
	if err != nil {
		t.Errorf("%v while decoding CTE document: %v", err, string(document))
		return
	}

	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		t.Errorf("CTE: Expected %v but got %v while decoding %v", expectedEvents, actualEvents, string(document))
	}
}

func assertEncodeDecodeCTE(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeCTEOpts(t, nil, nil, expected...)
}

func assertEncodeDecodeOpts(t *testing.T,
	cbeEncodeOpts *options.CBEEncoderOptions,
	cbeDecodeOpts *options.CBEDecoderOptions,
	cteEncodeOpts *options.CTEEncoderOptions,
	cteDecodeOpts *options.CTEDecoderOptions,
	expected ...*test.TEvent) {

	assertEncodeDecodeCBEOpts(t, cbeEncodeOpts, cbeDecodeOpts, expected...)
	assertEncodeDecodeCTEOpts(t, cteEncodeOpts, cteDecodeOpts, test.FilterEventsForCTE(expected)...)
}

func assertEncodeDecode(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeOpts(t, nil, nil, nil, nil, expected...)
}

func assertEncodeDecodeSet(t *testing.T, prefix []*test.TEvent, suffix []*test.TEvent, events []*test.TEvent) {
	for _, event := range events {
		allEvents := []*test.TEvent{}
		allEvents = append(allEvents, prefix...)
		allEvents = append(allEvents, event)
		allEvents = append(allEvents, test.Completions[event]...)
		allEvents = append(allEvents, suffix...)
		assertEncodeDecode(t, allEvents...)
	}
}

func assertEncodeDecodeSetTLO(t *testing.T, prefix []*test.TEvent, suffix []*test.TEvent, events []*test.TEvent) {
	for _, event := range events {
		allEvents := []*test.TEvent{}
		allEvents = append(allEvents, prefix...)
		allEvents = append(allEvents, event)
		allEvents = append(allEvents, test.Completions[event]...)
		allEvents = append(allEvents, suffix...)
		allEvents = test.FilterEventsForTLO(allEvents)
		assertEncodeDecode(t, allEvents...)
	}
}

func assertDecodeCBECTE(t *testing.T,
	cteEncodeOpts *options.CTEEncoderOptions,
	cteDecodeOpts *options.CTEDecoderOptions,
	cbeEncodeOpts *options.CBEEncoderOptions,
	cbeDecodeOpts *options.CBEDecoderOptions,
	cteExpectedDocument string,
	cbeExpectedDocument []byte,
	expectedEvents ...*test.TEvent) {

	var actualEvents *test.TEventStore

	textDecoder := ce.NewCTEDecoder(cteDecodeOpts)
	textEncoder := ce.NewCTEEncoder(cteEncodeOpts)
	textRules := ce.NewRules(textEncoder, nil)
	var cteActualDocument *bytes.Buffer

	binDecoder := ce.NewCBEDecoder(cbeDecodeOpts)
	binEncoder := ce.NewCBEEncoder(cbeEncodeOpts)
	var cbeActualDocument *bytes.Buffer

	cbeActualDocument = &bytes.Buffer{}
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(binEncoder, nil)); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(cbeExpectedDocument, cbeActualDocument.Bytes()) {
		t.Errorf("Expected %v but got %v", describe.D(cbeExpectedDocument), describe.D(cbeActualDocument.Bytes()))
	}

	actualEvents = test.NewTEventStore()
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(actualEvents, nil)); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expectedEvents, actualEvents.Events) {
		t.Errorf("Expected %v but got %v", expectedEvents, actualEvents.Events)
	}

	cteActualDocument = &bytes.Buffer{}
	textEncoder.PrepareToEncode(cteActualDocument)
	if err := binDecoder.DecodeDocument(cbeExpectedDocument, textRules); err != nil {
		t.Error(err)
		return
	}
	// Don't check the text document for equality since it won't be exactly the same.

	actualEvents = test.NewTEventStore()
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(actualEvents, nil)); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expectedEvents, actualEvents.Events) {
		t.Errorf("Expected %v but got %v", expectedEvents, actualEvents.Events)
	}
}

func assertDecodeEncode(t *testing.T,
	cteEncodeOpts *options.CTEEncoderOptions,
	cteDecodeOpts *options.CTEDecoderOptions,
	cbeEncodeOpts *options.CBEEncoderOptions,
	cbeDecodeOpts *options.CBEDecoderOptions,
	cteExpectedDocument string,
	cbeExpectedDocument []byte,
	expectedEvents ...*test.TEvent) {

	var actualEvents *test.TEventStore

	textDecoder := ce.NewCTEDecoder(cteDecodeOpts)
	textEncoder := ce.NewCTEEncoder(cteEncodeOpts)
	textRules := ce.NewRules(textEncoder, nil)
	var cteActualDocument *bytes.Buffer

	binDecoder := ce.NewCBEDecoder(cbeDecodeOpts)
	binEncoder := ce.NewCBEEncoder(cbeEncodeOpts)
	binRules := ce.NewRules(binEncoder, nil)
	var cbeActualDocument *bytes.Buffer

	cbeActualDocument = &bytes.Buffer{}
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), binRules); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(cbeExpectedDocument, cbeActualDocument.Bytes()) {
		t.Errorf("Expected %v but got %v", describe.D(cbeExpectedDocument), describe.D(cbeActualDocument.Bytes()))
	}

	actualEvents = test.NewTEventStore()
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(actualEvents, nil)); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expectedEvents, actualEvents.Events) {
		t.Errorf("Expected %v but got %v", expectedEvents, actualEvents.Events)
	}

	cteActualDocument = &bytes.Buffer{}
	textEncoder.PrepareToEncode(cteActualDocument)
	if err := binDecoder.DecodeDocument(cbeExpectedDocument, textRules); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(cteExpectedDocument, string(cteActualDocument.Bytes())) {
		t.Errorf("Expected [%v] but got [%v]", cteExpectedDocument, string(cteActualDocument.Bytes()))
	}

	actualEvents = test.NewTEventStore()
	binEncoder.PrepareToEncode(cbeActualDocument)
	if err := textDecoder.DecodeDocument([]byte(cteExpectedDocument), ce.NewRules(actualEvents, nil)); err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expectedEvents, actualEvents.Events) {
		t.Errorf("Expected %v but got %v", expectedEvents, actualEvents.Events)
	}
}

// ============================================================================
// Marshal/Unmarshal

func assertCBEMarshalUnmarshal(t *testing.T, expected interface{}) {
	marshalOptions := options.DefaultCBEMarshalerOptions()
	unmarshalOptions := options.DefaultCBEUnmarshalerOptions()
	assertCBEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
	// marshalOptions.Iterator.RecursionSupport = true
	// assertCBEMarshalUnmarshalWithOptions(t, marshalOptions, unmarshalOptions, expected)
}

func assertCBEMarshalUnmarshalWithOptions(t *testing.T,
	marshalOptions *options.CBEMarshalerOptions,
	unmarshalOptions *options.CBEUnmarshalerOptions,
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

	buffer := &bytes.Buffer{}
	err := ce.MarshalCTE(expected, buffer, marshalOptions)
	if err != nil {
		t.Errorf("CTE Marshal error: %v", err)
		return
	}

	var actual interface{}
	actual, err = ce.UnmarshalCTE(buffer, expected, unmarshalOptions)
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
	// assertCTEMarshalUnmarshal(t, expected)
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
