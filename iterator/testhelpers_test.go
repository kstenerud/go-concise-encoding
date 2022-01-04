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
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/types"
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

func TT() *test.TEvent                       { return test.TT() }
func FF() *test.TEvent                       { return test.FF() }
func I(v int64) *test.TEvent                 { return test.I(v) }
func BF(v float64) *test.TEvent              { return test.BF(v) }
func BBF(v *big.Float) *test.TEvent          { return test.BBF(v) }
func DF(v compact_float.DFloat) *test.TEvent { return test.DF(v) }
func BDF(v *apd.Decimal) *test.TEvent        { return test.BDF(v) }
func NULL() *test.TEvent                     { return test.NULL() }
func PAD(v int) *test.TEvent                 { return test.PAD(v) }
func B(v bool) *test.TEvent                  { return test.B(v) }
func PI(v uint64) *test.TEvent               { return test.PI(v) }
func NI(v uint64) *test.TEvent               { return test.NI(v) }
func BI(v *big.Int) *test.TEvent             { return test.BI(v) }
func NAN() *test.TEvent                      { return test.QNAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UID(v []byte) *test.TEvent              { return test.UID(v) }
func GT(v time.Time) *test.TEvent            { return test.GT(v) }
func CT(v compact_time.Time) *test.TEvent    { return test.T(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func RID(v string) *test.TEvent              { return test.RID(v) }
func CUB(v []byte) *test.TEvent              { return test.CB(v) }
func CUT(v string) *test.TEvent              { return test.CT(v) }
func AB(l uint64, v []byte) *test.TEvent     { return test.AB(l, v) }
func AU8(v []byte) *test.TEvent              { return test.AU8(v) }
func AU16(v []uint16) *test.TEvent           { return test.AU16(v) }
func AU32(v []uint32) *test.TEvent           { return test.AU32(v) }
func AU64(v []uint64) *test.TEvent           { return test.AU64(v) }
func AI8(v []int8) *test.TEvent              { return test.AI8(v) }
func AI16(v []int16) *test.TEvent            { return test.AI16(v) }
func AI32(v []int32) *test.TEvent            { return test.AI32(v) }
func AI64(v []int64) *test.TEvent            { return test.AI64(v) }
func AF16(v []float32) *test.TEvent          { return test.AF16(v) }
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
func AUUB() *test.TEvent                     { return test.AUB() }
func MB() *test.TEvent                       { return test.MB() }
func AC(l uint64, more bool) *test.TEvent    { return test.AC(l, more) }
func AD(v []byte) *test.TEvent               { return test.AD(v) }
func L() *test.TEvent                        { return test.L() }
func M() *test.TEvent                        { return test.M() }
func MUP(id string) *test.TEvent             { return test.MU(id) }
func NODE() *test.TEvent                     { return test.NODE() }
func EDGE() *test.TEvent                     { return test.EDGE() }
func E() *test.TEvent                        { return test.E() }
func MARK(id string) *test.TEvent            { return test.MARK(id) }
func REF(id string) *test.TEvent             { return test.REF(id) }
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

func InvokeEvents(receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.InvokeEvents(receiver, events...)
}

func iterateObject(object interface{},
	eventReceiver events.DataEventReceiver,
	sessionOptions *options.IteratorSessionOptions,
	iteratorOptions *options.IteratorOptions) {

	session := NewSession(nil, sessionOptions)
	iter := session.NewIterator(eventReceiver, iteratorOptions)
	iter.Iterate(object)
}

func assertIterateWithOptions(t *testing.T,
	sessionOptions *options.IteratorSessionOptions,
	iteratorOptions *options.IteratorOptions,
	obj interface{},
	evts ...*test.TEvent) {

	expected := append([]*test.TEvent{BD(), EvV}, evts...)
	expected = append(expected, ED())
	receiver := test.NewTEventStore(events.NewNullEventReceiver())
	iterateObject(obj, receiver, sessionOptions, iteratorOptions)

	if !equivalence.IsEquivalent(expected, receiver.Events) {
		t.Errorf("Expected %v but got %v", expected, receiver.Events)
	}
}

func assertIterate(t *testing.T, obj interface{}, events ...*test.TEvent) {
	assertIterateWithOptions(t,
		options.DefaultIteratorSessionOptions(),
		options.DefaultIteratorOptions(),
		obj,
		events...)
}
