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
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func cbeDecode(document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
	err = ce.NewCBEDecoder(nil).Decode(bytes.NewBuffer(document), receiver)
	events = receiver.Events
	return
}

func cbeEncodeDecode(expected ...*test.TEvent) (events []*test.TEvent, err error) {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCBEEncoder(nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEvents(encoder, expected...)
	document := buffer.Bytes()

	return cbeDecode(document)
}

func cteDecode(document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
	err = ce.NewCTEDecoder(nil).Decode(bytes.NewBuffer(document), receiver)
	events = receiver.Events
	return
}

func cteEncode(events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := ce.NewCTEEncoder(nil)
	encoder.PrepareToEncode(buffer)
	test.InvokeEvents(encoder, events...)
	return buffer.Bytes()
}

func cteEncodeDecode(events ...*test.TEvent) (decodedEvents []*test.TEvent, err error) {
	return cteDecode(cteEncode(events...))
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
func BIN(v []byte) *test.TEvent              { return test.BIN(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func VS(v string) *test.TEvent               { return test.VS(v) }
func URI(v string) *test.TEvent              { return test.URI(v) }
func CUB(v []byte) *test.TEvent              { return test.CUB(v) }
func CUT(v string) *test.TEvent              { return test.CUT(v) }
func BB() *test.TEvent                       { return test.BB() }
func SB() *test.TEvent                       { return test.SB() }
func VB() *test.TEvent                       { return test.VB() }
func UB() *test.TEvent                       { return test.UB() }
func CBB() *test.TEvent                      { return test.CBB() }
func CTB() *test.TEvent                      { return test.CTB() }
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

func assertEncodeDecodeCBE(t *testing.T, expected ...*test.TEvent) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() {
		debug.DebugOptions.PassThroughPanics = false
	}()
	actual, err := cbeEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CBE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecodeCTE(t *testing.T, expected ...*test.TEvent) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	actual, err := cteEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("CTE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecode(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeCBE(t, expected...)
	assertEncodeDecodeCTE(t, expected...)
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

	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()
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

	debug.DebugOptions.PassThroughPanics = true
	defer func() {
		debug.DebugOptions.PassThroughPanics = false
	}()
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
