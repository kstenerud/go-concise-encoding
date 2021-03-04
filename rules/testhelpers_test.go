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
	"reflect"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/version"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

func NewBytes(count int, startIndex int) []byte {
	return test.GenerateBytes(count, startIndex)
}

func NewString(count int, startIndex int) []byte {
	return []byte(test.GenerateString(count, startIndex))
}

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

var (
	EvBD     = test.EvBD
	EvED     = test.EvED
	EvV      = test.EvV
	EvPAD    = test.EvPAD
	EvNA     = test.EvNA
	EvB      = test.EvB
	EvTT     = test.EvTT
	EvFF     = test.EvFF
	EvPI     = test.EvPI
	EvNI     = test.EvNI
	EvI      = test.EvI
	EvBI     = test.EvBI
	EvBINil  = test.EvBINil
	EvF      = test.EvF
	EvFNAN   = test.EvFNAN
	EvBF     = test.EvBF
	EvBFNil  = test.EvBFNil
	EvDF     = test.EvDF
	EvDFNAN  = test.EvDFNAN
	EvBDF    = test.EvBDF
	EvBDFNil = test.EvBDFNil
	EvBDFNAN = test.EvBDFNAN
	EvNAN    = test.EvNAN
	EvUUID   = test.EvUUID
	EvGT     = test.EvGT
	EvCT     = test.EvCT
	EvL      = test.EvL
	EvM      = test.EvM
	EvMUP    = test.EvMUP
	EvMETA   = test.EvMETA
	EvCMT    = test.EvCMT
	EvE      = test.EvE
	EvMARK   = test.EvMARK
	EvREF    = test.EvREF
	EvAC     = test.EvAC
	EvAD     = test.EvAD
	EvS      = test.EvS
	EvSB     = test.EvSB
	EvRID    = test.EvRID
	EvRB     = test.EvRB
	EvCUB    = test.EvCUB
	EvCBB    = test.EvCBB
	EvCUT    = test.EvCUT
	EvCTB    = test.EvCTB
	EvAB     = test.EvAB
	EvABB    = test.EvABB
	EvAU8    = test.EvAU8
	EvAU8B   = test.EvAU8B
	EvAU16   = test.EvAU16
	EvAU16B  = test.EvAU16B
	EvAU32   = test.EvAU32
	EvAU32B  = test.EvAU32B
	EvAU64   = test.EvAU64
	EvAU64B  = test.EvAU64B
	EvAI8    = test.EvAI8
	EvAI8B   = test.EvAI8B
	EvAI16   = test.EvAI16
	EvAI16B  = test.EvAI16B
	EvAI32   = test.EvAI32
	EvAI32B  = test.EvAI32B
	EvAI64   = test.EvAI64
	EvAI64B  = test.EvAI64B
	EvAF16   = test.EvAF16
	EvAF16B  = test.EvAF16B
	EvAF32   = test.EvAF32
	EvAF32B  = test.EvAF32B
	EvAF64   = test.EvAF64
	EvAF64B  = test.EvAF64B
	EvAUU    = test.EvAUU
	EvAUUB   = test.EvAUUB
)

var (
	MarkerIDTypes                  = test.MarkerIDTypes
	InvalidMarkerIDTypes           = test.InvalidMarkerIDTypes
	ReferenceIDTypes               = test.ReferenceIDTypes
	InvalidReferenceIDTypes        = test.InvalidReferenceIDTypes
	KeyableReferenceIDTypes        = test.KeyableReferenceIDTypes
	InvalidKeyableReferenceIDTypes = test.InvalidKeyableReferenceIDTypes
	ArrayBeginTypes                = test.ArrayBeginTypes
	ValidTLOValues                 = test.ValidTLOValues
	InvalidTLOValues               = test.InvalidTLOValues
	ValidMapKeys                   = test.ValidMapKeys
	InvalidMapKeys                 = test.InvalidMapKeys
	ValidMapValues                 = test.ValidMapValues
	InvalidMapValues               = test.InvalidMapValues
	ValidListValues                = test.ValidListValues
	InvalidListValues              = test.InvalidListValues
	ValidCommentValues             = test.ValidCommentValues
	InvalidCommentValues           = test.InvalidCommentValues
	ValidMarkupNames               = test.ValidMarkupNames
	InvalidMarkupNames             = test.InvalidMarkupNames
	ValidMarkupContents            = test.ValidMarkupContents
	InvalidMarkupContents          = test.InvalidMarkupContents
	ValidAfterArrayBegin           = test.ValidAfterArrayBegin
	InvalidAfterArrayBegin         = test.InvalidAfterArrayBegin
	ValidAfterArrayChunk           = test.ValidAfterArrayChunk
	InvalidAfterArrayChunk         = test.InvalidAfterArrayChunk
	ValidMarkerIDs                 = test.ValidMarkerIDs
	InvalidMarkerIDs               = test.InvalidMarkerIDs
	ValidMarkerValues              = test.ValidMarkerValues
	InvalidMarkerValues            = test.InvalidMarkerValues
	ValidReferenceIDs              = test.ValidReferenceIDs
	InvalidReferenceIDs            = test.InvalidReferenceIDs
)

func ComplementaryEvents(events []*test.TEvent) []*test.TEvent {
	return test.ComplementaryEvents(events)
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
func AUU(v []byte) *test.TEvent              { return test.AUU(v) }
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
func CONST(n string, e bool) *test.TEvent    { return test.CONST(n, e) }
func BD() *test.TEvent                       { return test.BD() }
func ED() *test.TEvent                       { return test.ED() }

func InvokeEvents(receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.InvokeEvents(receiver, events...)
}

// ============================================================================

const rulesCodecVersion = 1

var uint8Type = reflect.TypeOf(uint8(0))

func assertEventsSucceed(t *testing.T, receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.AssertNoPanic(t, events, func() {
		InvokeEvents(receiver, events...)
	})
}

func assertEachEventSucceeds(t *testing.T, prefix func() events.DataEventReceiver, events ...*test.TEvent) {
	for _, event := range events {
		assertEventsSucceed(t, prefix(), event)
	}
}

func assertEventsFail(t *testing.T, receiver events.DataEventReceiver, events ...*test.TEvent) {
	test.AssertPanics(t, events, func() {
		InvokeEvents(receiver, events...)
	})
}

func assertEachEventFails(t *testing.T, prefix func() events.DataEventReceiver, events ...*test.TEvent) {
	for _, event := range events {
		assertEventsFail(t, prefix(), event)
	}
}

func assertEventsMaxDepth(t *testing.T, maxDepth int, events ...*test.TEvent) {
	rules := newRulesWithMaxDepth(maxDepth)
	assertEventsSucceed(t, rules, events...)
}

func assertRulesOnString(t *testing.T, rules *RulesEventReceiver, value string) {
	length := len(value)
	test.AssertNoPanic(t, value, func() { rules.OnArrayBegin(events.ArrayTypeString) })
	test.AssertNoPanic(t, value, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, value, func() { rules.OnArrayData([]byte(value)) })
	}
}

func assertRulesAddBytes(t *testing.T, rules *RulesEventReceiver, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, value, func() { rules.OnArrayBegin(events.ArrayTypeUint8) })
	test.AssertNoPanic(t, value, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, value, func() { rules.OnArrayData(value) })
	}
}

func assertRulesAddRID(t *testing.T, rules *RulesEventReceiver, ResourceID string) {
	length := len(ResourceID)
	test.AssertNoPanic(t, ResourceID, func() { rules.OnArrayBegin(events.ArrayTypeResourceID) })
	test.AssertNoPanic(t, ResourceID, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, ResourceID, func() { rules.OnArrayData([]byte(ResourceID)) })
	}
}

func assertRulesAddCustomBinary(t *testing.T, rules *RulesEventReceiver, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, value, func() { rules.OnArrayBegin(events.ArrayTypeCustomBinary) })
	test.AssertNoPanic(t, value, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, value, func() { rules.OnArrayData(value) })
	}
}

func assertRulesAddCustomText(t *testing.T, rules *RulesEventReceiver, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, string(value), func() { rules.OnArrayBegin(events.ArrayTypeCustomText) })
	test.AssertNoPanic(t, string(value), func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, value, func() { rules.OnArrayData(value) })
	}
}

func newRulesAfterVersion(opts *options.RuleOptions) *RulesEventReceiver {
	rules := NewRules(events.NewNullEventReceiver(), opts)
	rules.OnBeginDocument()
	rules.OnVersion(version.ConciseEncodingVersion)
	return rules
}

func newRulesWithMaxDepth(maxDepth int) *RulesEventReceiver {
	opts := options.DefaultRuleOptions()
	opts.MaxContainerDepth = uint64(maxDepth)
	return newRulesAfterVersion(opts)
}
