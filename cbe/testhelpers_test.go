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
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

var EvV = test.EvV

const (
	header = 0x83
	ceVer  = version.ConciseEncodingVersion

	typeDecimal      = 0x65
	typePosInt       = 0x66
	typeNegInt       = 0x67
	typePosInt8      = 0x68
	typeNegInt8      = 0x69
	typePosInt16     = 0x6a
	typeNegInt16     = 0x6b
	typePosInt32     = 0x6c
	typeNegInt32     = 0x6d
	typePosInt64     = 0x6e
	typeNegInt64     = 0x6f
	typeFloat16      = 0x70
	typeFloat32      = 0x71
	typeFloat64      = 0x72
	typeUID          = 0x73
	typeReserved74   = 0x74
	typeReserved75   = 0x75
	typeRelationship = 0x76
	typeComment      = 0x77
	typeMarkup       = 0x78
	typeMap          = 0x79
	typeList         = 0x7a
	typeEndContainer = 0x7b
	typeFalse        = 0x7c
	typeTrue         = 0x7d
	typeNil          = 0x7e
	typePadding      = 0x7f
	typeString0      = 0x80
	typeString1      = 0x81
	typeString2      = 0x82
	typeString3      = 0x83
	typeString4      = 0x84
	typeString5      = 0x85
	typeString6      = 0x86
	typeString7      = 0x87
	typeString8      = 0x88
	typeString9      = 0x89
	typeString10     = 0x8a
	typeString11     = 0x8b
	typeString12     = 0x8c
	typeString13     = 0x8d
	typeString14     = 0x8e
	typeString15     = 0x8f
	typeString       = 0x90
	typeRID          = 0x91
	typeCustomBinary = 0x92
	typeCustomText   = 0x93
	typePlane2       = 0x94
	typeArrayUint8   = 0x95
	typeArrayBit     = 0x96
	typeMarker       = 0x97
	typeReference    = 0x98
	typeDate         = 0x99
	typeTime         = 0x9a
	typeTimestamp    = 0x9b

	typeShortArrayInt8    = 0x00
	typeShortArrayUint16  = 0x10
	typeShortArrayInt16   = 0x20
	typeShortArrayUint32  = 0x30
	typeShortArrayInt32   = 0x40
	typeShortArrayUint64  = 0x50
	typeShortArrayInt64   = 0x60
	typeShortArrayFloat16 = 0x70
	typeShortArrayFloat32 = 0x80
	typeShortArrayFloat64 = 0x90
	typeShortArrayUID     = 0xa0

	typeNA           = 0xe0
	typeRIDCat       = 0xe1
	typeRIDRef       = 0xe2
	typeMedia        = 0xe3
	typeArrayInt8    = 0xff
	typeArrayUint16  = 0xfe
	typeArrayInt16   = 0xfd
	typeArrayUint32  = 0xfc
	typeArrayInt32   = 0xfb
	typeArrayUint64  = 0xfa
	typeArrayInt64   = 0xf9
	typeArrayFloat16 = 0xf8
	typeArrayFloat32 = 0xf7
	typeArrayFloat64 = 0xf6
	typeArrayUID     = 0xf5
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
func N() *test.TEvent                        { return test.N() }
func NA() *test.TEvent                       { return test.NA() }
func PAD(v int) *test.TEvent                 { return test.PAD(v) }
func B(v bool) *test.TEvent                  { return test.B(v) }
func PI(v uint64) *test.TEvent               { return test.PI(v) }
func NI(v uint64) *test.TEvent               { return test.NI(v) }
func BI(v *big.Int) *test.TEvent             { return test.BI(v) }
func NAN() *test.TEvent                      { return test.NAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UID(v []byte) *test.TEvent              { return test.UID(v) }
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
func MB() *test.TEvent                       { return test.MB() }
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

func InvokeEvents(receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.InvokeEvents(receiver, events...)
}

var DebugPrintEvents = false

func decodeToEvents(opts *options.CBEDecoderOptions, document []byte) (evts []*test.TEvent, err error) {
	var topLevelReceiver events.DataEventReceiver
	ter := test.NewTEventStore()
	topLevelReceiver = rules.NewRules(ter, nil)
	if DebugPrintEvents {
		topLevelReceiver = test.NewStdoutTEventPrinter(topLevelReceiver)
	}
	err = NewDecoder(opts).Decode(bytes.NewBuffer(document), topLevelReceiver)
	evts = ter.Events
	return
}

func encodeEvents(opts *options.CBEEncoderOptions, events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := NewEncoder(opts)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	test.InvokeEvents(receiver, events...)
	return buffer.Bytes()
}

func assertDecode(t *testing.T, opts *options.CBEDecoderOptions, document []byte, expectedEvents ...*test.TEvent) (successful bool, events []*test.TEvent) {
	actualEvents, err := decodeToEvents(opts, document)
	if err != nil {
		t.Errorf("Error [%v] while decoding %v", err, describe.D(document))
		return
	}

	if len(expectedEvents) > 0 {
		if !equivalence.IsEquivalent(actualEvents, expectedEvents) {
			t.Errorf("Expected document %v to decode to events %v but got %v", describe.D(document), expectedEvents, actualEvents)
			return
		}
	}
	events = actualEvents
	successful = true
	return
}

func assertDecodeFails(t *testing.T, document []byte) {
	_, err := decodeToEvents(nil, document)
	if err == nil {
		t.Errorf("Expected decode to fail in document %v", describe.D(document))
	}
}

func assertEncode(t *testing.T, opts *options.CBEEncoderOptions, expectedDocument []byte, events ...*test.TEvent) (successful bool) {
	actualDocument := encodeEvents(opts, events...)
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		t.Errorf("Expected encoded document %v but got %v after sending events %v", describe.D(expectedDocument), describe.D(actualDocument), events)
		return
	}
	successful = true
	return
}

func assertEncodeFails(t *testing.T, opts *options.CBEEncoderOptions, events ...*test.TEvent) (successful bool) {
	successful = test.AssertPanics(t, "encode", func() {
		encodeEvents(opts, events...)
	})
	return
}

func assertDecodeEncode(t *testing.T, document []byte, expectedEvents ...*test.TEvent) (successful bool) {
	successful, actualEvents := assertDecode(t, nil, document, expectedEvents...)
	if !successful {
		return
	}
	return assertEncode(t, nil, document, actualEvents...)
}

func assertMarshal(t *testing.T, value interface{}, expectedDocument []byte) (successful bool) {
	document, err := NewMarshaler(nil).MarshalToDocument(value)
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
	actualValue, err := NewUnmarshaler(nil).UnmarshalFromDocument([]byte(document), expectedValue)
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
