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

package iterator

import (
	"math/big"
	"net/url"
	"testing"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/types"
	"github.com/kstenerud/go-describe"
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

func NewTS(year, month, day, hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	return test.NewTS(year, month, day, hour, minute, second, nanosecond, areaLocation)
}

func NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	return test.NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths)
}

func NewTestingOuterStruct(baseValue int) *test.TestingOuterStruct {
	return test.NewTestingOuterStruct(baseValue)
}

func NewNode(value interface{}, children []interface{}) *types.Node {
	return test.NewNode(value, children)
}

func NewEdge(source interface{}, description interface{}, destination interface{}) *types.Edge {
	return test.NewEdge(source, description, destination)
}

func AB(v []bool) test.Event              { return test.AB(v) }
func ACL(l uint64) test.Event             { return test.ACL(l) }
func ACM(l uint64) test.Event             { return test.ACM(l) }
func ADB(v []bool) test.Event             { return test.ADB(v) }
func ADF16(v []float32) test.Event        { return test.ADF16(v) }
func ADF32(v []float32) test.Event        { return test.ADF32(v) }
func ADF64(v []float64) test.Event        { return test.ADF64(v) }
func ADI16(v []int16) test.Event          { return test.ADI16(v) }
func ADI32(v []int32) test.Event          { return test.ADI32(v) }
func ADI64(v []int64) test.Event          { return test.ADI64(v) }
func ADI8(v []int8) test.Event            { return test.ADI8(v) }
func ADT(v string) test.Event             { return test.ADT(v) }
func ADU(v [][]byte) test.Event           { return test.ADU(v) }
func ADU16(v []uint16) test.Event         { return test.ADU16(v) }
func ADU32(v []uint32) test.Event         { return test.ADU32(v) }
func ADU64(v []uint64) test.Event         { return test.ADU64(v) }
func ADU8(v []uint8) test.Event           { return test.ADU8(v) }
func AF16(v []float32) test.Event         { return test.AF16(v) }
func AF32(v []float32) test.Event         { return test.AF32(v) }
func AF64(v []float64) test.Event         { return test.AF64(v) }
func AI16(v []int16) test.Event           { return test.AI16(v) }
func AI32(v []int32) test.Event           { return test.AI32(v) }
func AI64(v []int64) test.Event           { return test.AI64(v) }
func AI8(v []int8) test.Event             { return test.AI8(v) }
func AU(v [][]byte) test.Event            { return test.AU(v) }
func AU16(v []uint16) test.Event          { return test.AU16(v) }
func AU32(v []uint32) test.Event          { return test.AU32(v) }
func AU64(v []uint64) test.Event          { return test.AU64(v) }
func AU8(v []byte) test.Event             { return test.AU8(v) }
func B(v bool) test.Event                 { return test.B(v) }
func BAB() test.Event                     { return test.BAB() }
func BAF16() test.Event                   { return test.BAF16() }
func BAF32() test.Event                   { return test.BAF32() }
func BAF64() test.Event                   { return test.BAF64() }
func BAI16() test.Event                   { return test.BAI16() }
func BAI32() test.Event                   { return test.BAI32() }
func BAI64() test.Event                   { return test.BAI64() }
func BAI8() test.Event                    { return test.BAI8() }
func BAU() test.Event                     { return test.BAU() }
func BAU16() test.Event                   { return test.BAU16() }
func BAU32() test.Event                   { return test.BAU32() }
func BAU64() test.Event                   { return test.BAU64() }
func BAU8() test.Event                    { return test.BAU8() }
func BCB(v uint64) test.Event             { return test.BCB(v) }
func BCT(v uint64) test.Event             { return test.BCT(v) }
func BMEDIA(v string) test.Event          { return test.BMEDIA(v) }
func BREFR() test.Event                   { return test.BREFR() }
func BRID() test.Event                    { return test.BRID() }
func BS() test.Event                      { return test.BS() }
func CB(t uint64, v []byte) test.Event    { return test.CB(t, v) }
func CM(v string) test.Event              { return test.CM(v) }
func CS(v string) test.Event              { return test.CS(v) }
func CT(t uint64, v string) test.Event    { return test.CT(t, v) }
func E() test.Event                       { return test.E() }
func EDGE() test.Event                    { return test.EDGE() }
func L() test.Event                       { return test.L() }
func M() test.Event                       { return test.M() }
func MARK(id string) test.Event           { return test.MARK(id) }
func MEDIA(t string, v []byte) test.Event { return test.MEDIA(t, v) }
func N(v interface{}) test.Event          { return test.N(v) }
func NAN() test.Event                     { return test.NAN() }
func NODE() test.Event                    { return test.NODE() }
func NULL() test.Event                    { return test.NULL() }
func PAD() test.Event                     { return test.PAD() }
func REFL(id string) test.Event           { return test.REFL(id) }
func REFR(v string) test.Event            { return test.REFR(v) }
func RID(v string) test.Event             { return test.RID(v) }
func S(v string) test.Event               { return test.S(v) }
func SI(id string) test.Event             { return test.SI(id) }
func SNAN() test.Event                    { return test.SNAN() }
func ST(id string) test.Event             { return test.ST(id) }
func T(v compact_time.Time) test.Event    { return test.T(v) }
func UID(v []byte) test.Event             { return test.UID(v) }
func V(v uint64) test.Event               { return test.V(v) }

func iterateObject(object interface{},
	eventReceiver events.DataEventReceiver,
	iteratorConfiguration *configuration.IteratorConfiguration) {

	session := NewSession(nil, iteratorConfiguration)
	iter := session.NewIterator(eventReceiver)
	iter.Iterate(object)
}

func assertIterateWithConfiguration(t *testing.T,
	iteratorConfiguration *configuration.IteratorConfiguration,
	obj interface{},
	evts ...test.Event) {

	expected := test.Events(append(test.Events{EvV}, evts...))
	receiver, container := test.NewEventCollector(nil)
	iterateObject(obj, receiver, iteratorConfiguration)

	if !container.IsEquivalentTo(expected) {
		t.Errorf("Expected %v to iterate to events [%v] but got [%v]", describe.D(obj), expected, container.Events)
	}
}

func assertIterate(t *testing.T, obj interface{}, events ...test.Event) {
	config := configuration.DefaultIteratorConfiguration()
	assertIterateWithConfiguration(t, &config, obj, events...)
}
