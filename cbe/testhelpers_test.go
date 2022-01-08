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
	"math/big"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

var EvV = test.EvV

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
func RREF(v string) *test.TEvent             { return test.RREF(v) }
func RRB() *test.TEvent                      { return test.RRB() }
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

func InvokeEvents(receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.InvokeEvents(receiver, events...)
}

var DebugPrintEvents = false

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
