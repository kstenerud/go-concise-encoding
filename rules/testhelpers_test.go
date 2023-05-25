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

package rules

import (
	"math/big"
	"net/url"
	"testing"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/nullevent"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
)

const (
	ceVer = version.ConciseEncodingVersion
)

func NewBytes(count int, startIndex int) []byte {
	return test.GenerateBytes(count, startIndex)
}

func NewString(count int, startIndex int) string {
	return test.GenerateString(count, startIndex)
}

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

var (
	EvV = test.EvV
)

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
func BD() test.Event                      { return test.BD() }
func BS() test.Event                      { return test.BS() }
func CB(t uint64, v []byte) test.Event    { return test.CB(t, v) }
func CM(v string) test.Event              { return test.CM(v) }
func CS(v string) test.Event              { return test.CS(v) }
func CT(t uint64, v string) test.Event    { return test.CT(t, v) }
func E() test.Event                       { return test.E() }
func ED() test.Event                      { return test.ED() }
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
func REC(id string) test.Event            { return test.REC(id) }
func REFL(id string) test.Event           { return test.REFL(id) }
func REFR(v string) test.Event            { return test.REFR(v) }
func RID(v string) test.Event             { return test.RID(v) }
func RT(id string) test.Event             { return test.RT(id) }
func S(v string) test.Event               { return test.S(v) }
func SNAN() test.Event                    { return test.SNAN() }
func T(v compact_time.Time) test.Event    { return test.T(v) }
func UID(v []byte) test.Event             { return test.UID(v) }
func V(v uint64) test.Event               { return test.V(v) }

var DebugPrintEvents = false

func invokeEvents(receiver events.DataEventReceiver, events ...test.Event) {
	if DebugPrintEvents {
		receiver = test.NewEventPrinter(receiver)
	}

	for _, event := range events {
		event.Invoke(receiver)
	}
}

// ============================================================================

func assertEvents(t *testing.T, allEvents ...test.Event) {
	assertEventsSucceed(t, NewRules(nullevent.NewNullEventReceiver(), nil), allEvents...)
}

func assertEventsSucceed(t *testing.T, receiver events.DataEventReceiver, events ...test.Event) {
	test.AssertNoPanic(t, events, func() {
		invokeEvents(receiver, events...)
	})
}

func assertEventStreamsSucceed(t *testing.T, eventStreams []test.Events) {
	for _, stream := range eventStreams {
		rules := NewRules(nullevent.NewNullEventReceiver(), nil)
		test.AssertNoPanic(t, stream, func() {
			invokeEvents(rules, stream...)
		})
	}
}

func assertEventStreamsFail(t *testing.T, eventStreams []test.Events) {
	for _, stream := range eventStreams {
		rules := NewRules(nullevent.NewNullEventReceiver(), nil)
		test.AssertPanics(t, stream, func() {
			invokeEvents(rules, stream...)
		})
	}
}

func assertEventsFail(t *testing.T, receiver events.DataEventReceiver, events ...test.Event) {
	test.AssertPanics(t, events, func() {
		invokeEvents(receiver, events...)
	})
}

func assertEachEventFails(t *testing.T, prefix func() events.DataEventReceiver, events ...test.Event) {
	for _, event := range events {
		assertEventsFail(t, prefix(), event)
	}
}

func assertEventsMaxDepth(t *testing.T, maxDepth int, events ...test.Event) {
	rules := newRulesWithMaxDepth(maxDepth)
	assertEventsSucceed(t, rules, events...)
}

func newRulesAfterVersion(config *configuration.RuleConfiguration) *RulesEventReceiver {
	rules := NewRules(nullevent.NewNullEventReceiver(), config)
	rules.OnBeginDocument()
	rules.OnVersion(ceVer)
	return rules
}

func newRulesWithMaxDepth(maxDepth int) *RulesEventReceiver {
	config := configuration.DefaultRuleConfiguration()
	config.MaxContainerDepth = uint64(maxDepth)
	return newRulesAfterVersion(&config)
}
