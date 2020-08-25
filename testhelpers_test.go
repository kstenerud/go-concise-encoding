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
	"math/big"
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

func cbeDecode(opts *options.CBEDecoderOptions, document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
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
	receiver := test.NewTER()
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

	document := cteEncode(encodeOpts, events...)
	return cteDecode(decodeOpts, document)
}

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
func NAN() *test.TEvent                      { return test.NAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UUID(v []byte) *test.TEvent             { return test.UUID(v) }
func GT(v time.Time) *test.TEvent            { return test.GT(v) }
func CT(v *compact_time.Time) *test.TEvent   { return test.CT(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func VS(v string) *test.TEvent               { return test.VS(v) }
func URI(v string) *test.TEvent              { return test.URI(v) }
func CUB(v []byte) *test.TEvent              { return test.CUB(v) }
func CUT(v string) *test.TEvent              { return test.CUT(v) }
func AU8(v []byte) *test.TEvent              { return test.AU8(v) }
func SB() *test.TEvent                       { return test.SB() }
func VB() *test.TEvent                       { return test.VB() }
func UB() *test.TEvent                       { return test.UB() }
func CBB() *test.TEvent                      { return test.CBB() }
func CTB() *test.TEvent                      { return test.CTB() }
func AU8B() *test.TEvent                     { return test.AU8B() }
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
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

// ============================================================================
// Encode/Decode

func assertEncodeDecodeCBEOpts(t *testing.T,
	encodeOpts *options.CBEEncoderOptions,
	decodeOpts *options.CBEDecoderOptions,
	expected ...*test.TEvent) {

	actual, err := cbeEncodeDecode(encodeOpts, decodeOpts, expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CBE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecodeCBE(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeCBEOpts(t, nil, nil, expected...)
}

func assertEncodeDecodeCTEOpts(t *testing.T,
	encodeOpts *options.CTEEncoderOptions,
	decodeOpts *options.CTEDecoderOptions,
	expected ...*test.TEvent) {

	actual, err := cteEncodeDecode(encodeOpts, decodeOpts, expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CTE: Expected %v but got %v", expected, actual)
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
	assertEncodeDecodeCTEOpts(t, cteEncodeOpts, cteDecodeOpts, expected...)
}

func assertEncodeDecode(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeOpts(t, nil, nil, nil, nil, expected...)
}

func assertEncodeDecodeImpliedStructure(t *testing.T,
	impliedStruct options.ImpliedStructure,
	version uint64,
	cteDocument string,
	expected ...*test.TEvent) {

	cteEncodeOpts := options.DefaultCTEEncoderOptions()
	cteEncodeOpts.ConciseEncodingVersion = version
	cteEncodeOpts.ImpliedStructure = impliedStruct
	cteDecodeOpts := options.DefaultCTEDecoderOptions()
	cteDecodeOpts.ConciseEncodingVersion = version
	cteDecodeOpts.ImpliedStructure = impliedStruct
	cbeEncodeOpts := options.DefaultCBEEncoderOptions()
	cbeEncodeOpts.ConciseEncodingVersion = version
	cbeEncodeOpts.ImpliedStructure = impliedStruct
	cbeDecodeOpts := options.DefaultCBEDecoderOptions()
	cbeDecodeOpts.ConciseEncodingVersion = version
	cbeDecodeOpts.ImpliedStructure = impliedStruct

	events, err := cteDecode(cteDecodeOpts, []byte(cteDocument))
	if err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expected, events) {
		t.Errorf("CTE: Expected %v but got %v", expected, events)
		return
	}
	actualCTEDoc := string(cteEncode(cteEncodeOpts, expected...))
	if actualCTEDoc != cteDocument {
		t.Errorf("CTE: Expected doc [%v] but got [%v]", cteDocument, actualCTEDoc)
	}

	actualCBEDoc := cbeEncode(cbeEncodeOpts, expected...)
	events, err = cbeDecode(cbeDecodeOpts, actualCBEDoc)
	if err != nil {
		t.Error(err)
		return
	}
	if !equivalence.IsEquivalent(expected, events) {
		t.Errorf("CBE: Expected %v but got %v", expected, events)
		return
	}
}

// ============================================================================
// Marshal/Unmarshal

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
