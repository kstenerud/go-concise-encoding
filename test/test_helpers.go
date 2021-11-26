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

// Test helper code.
package test

import (
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"
	"github.com/kstenerud/go-concise-encoding/version"
)

// ----------------------------------------------------------------------------
// Pass through panics
// ----------------------------------------------------------------------------

// Causes the library to pass through all panics for the duration of the current
// function instead of converting them to error objects. This can be useful for
// tracking down the ultimate cause.
//
// Usage: defer test.PassThroughPanics(true)()
func PassThroughPanics(shouldPassThrough bool) func() {
	oldValue := debug.DebugOptions.PassThroughPanics
	debug.DebugOptions.PassThroughPanics = shouldPassThrough
	return func() {
		debug.DebugOptions.PassThroughPanics = oldValue
	}
}

// ----------------------------------------------------------------------------
// Constructors for common data types
// ----------------------------------------------------------------------------

func NewBigInt(str string) *big.Int {
	bi := new(big.Int)
	_, success := bi.SetString(str, 0)
	if !success {
		panic(fmt.Errorf("malformed unit test: Cannot convert [%v] to big.Int", str))
	}
	return bi
}

func countDecimalSigDigits(str string) int {
	for len(str) > 1 && str[0] == '0' {
		str = str[1:]
	}
	sigDigits := 0
	for _, b := range str {
		switch {
		case b >= '0' && b <= '9':
			sigDigits++
		case b == '.':
			// Ignore
		default:
			return sigDigits
		}
	}
	return sigDigits
}

var bfHexMatcher = regexp.MustCompile(`0[xX]([0-9a-fA-F]+)(\.[0-9a-fA-F]+)?[pP]?.*`)

func NewBigFloat(str string) *big.Float {
	match := bfHexMatcher.FindStringSubmatch(str)
	if len(match) > 1 {
		sigDigits := len(match[1])
		if len(match) > 2 {
			sigDigits += len(match[2]) - 1
		}
		bf := &big.Float{}
		bf.SetPrec(uint(common.HexDigitsToBits(sigDigits)))
		if _, success := bf.SetString(str); success {
			return bf
		} else {
			panic(fmt.Errorf("could not convert [%v] to big float", str))
		}
	}

	sigDigits := countDecimalSigDigits(str)
	f, _, err := big.ParseFloat(str, 10, uint(common.DecimalDigitsToBits(sigDigits)), big.ToNearestEven)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Cannot convert [%v] to big.Float: %w", str, err))
	}
	return f
}

func NewDFloat(str string) compact_float.DFloat {
	df, err := compact_float.DFloatFromString(str)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Cannot convert [%v] to compact_float.DFloat: %w", str, err))
	}
	return df
}

func NewBDF(str string) *apd.Decimal {
	v, _, err := apd.NewFromString(str)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Cannot convert [%v] to apd.Decimal: %w", str, err))
	}
	return v
}

func NewRID(RIDString string) *url.URL {
	rid, err := url.Parse(RIDString)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Bad URL (%v): %w", RIDString, err))
	}
	return rid
}

func NewNode(value interface{}, children []interface{}) *types.Node {
	return &types.Node{
		Value:    value,
		Children: children,
	}
}

func NewEdge(source interface{}, description interface{}, destination interface{}) *types.Edge {
	return &types.Edge{
		Source:      source,
		Description: description,
		Destination: destination,
	}
}

func NewDate(year, month, day int) compact_time.Time {
	t := compact_time.NewDate(year, month, day)
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTime(hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	t := compact_time.NewTime(hour, minute, second, nanosecond, compact_time.TZAtAreaLocation(areaLocation))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTimeLL(hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	t := compact_time.NewTime(hour, minute, second, nanosecond, compact_time.TZAtLatLong(latitudeHundredths, longitudeHundredths))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTimeOff(hour, minute, second, nanosecond, minutesOffset int) compact_time.Time {
	t := compact_time.NewTime(hour, minute, second, nanosecond, compact_time.TZWithMiutesOffsetFromUTC(minutesOffset))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTS(year, month, day, hour, minute, second, nanosecond int, areaLocation string) compact_time.Time {
	t := compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, compact_time.TZAtAreaLocation(areaLocation))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTSLL(year, month, day, hour, minute, second, nanosecond, latitudeHundredths, longitudeHundredths int) compact_time.Time {
	t := compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, compact_time.TZAtLatLong(latitudeHundredths, longitudeHundredths))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func NewTSOff(year, month, day, hour, minute, second, nanosecond, minutesOffset int) compact_time.Time {
	t := compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, compact_time.TZWithMiutesOffsetFromUTC(minutesOffset))
	err := t.Validate()
	if err != nil {
		panic(err)
	}
	return t
}

func AsGoTime(t compact_time.Time) time.Time {
	gt, err := t.AsGoTime()
	if err != nil {
		panic(err)
	}
	return gt
}

func AsCompactTime(t time.Time) compact_time.Time {
	ct := compact_time.AsCompactTime(t)
	err := ct.Validate()
	if err != nil {
		panic(err)
	}
	return ct
}

// ----------------------------------------------------------------------------
// Panics
// ----------------------------------------------------------------------------

func ReportPanic(function func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	function()
	return
}

func AssertNoPanic(t *testing.T, name interface{}, function func()) {
	if debug.DebugOptions.PassThroughPanics {
		function()
	} else {
		if err := ReportPanic(function); err != nil {
			t.Errorf("Error [%v] in %v", err, name)
		}
	}
}

func AssertPanics(t *testing.T, name interface{}, function func()) bool {
	if err := ReportPanic(function); err == nil {
		t.Errorf("Expected an error in %v", name)
		return false
	}
	return true
}

// ----------------------------------------------------------------------------
// Utilities
// ----------------------------------------------------------------------------

func BytesToString(bytes []byte) string {
	var builder strings.Builder
	builder.WriteByte('[')
	for i, c := range bytes {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteByte(hexChar(c >> 4))
		builder.WriteByte(hexChar(c & 15))
	}
	builder.WriteByte(']')
	return builder.String()
}

// ----------------------------------------------------------------------------
// Generators
// ----------------------------------------------------------------------------

func GenerateString(charCount int, startIndex int) string {
	charRange := int('z' - 'a')
	var object strings.Builder
	for i := 0; i < charCount; i++ {
		ch := 'a' + (i+charCount+startIndex)%charRange
		object.WriteByte(byte(ch))
	}
	return object.String()
}

func GenerateBytes(length int, startIndex int) []byte {
	return []byte(GenerateString(length, startIndex))
}

func GenerateArrayElements(elemSize, length int) []byte {
	var result []byte
	for i := 0; i < length; i++ {
		for e := 0; e < elemSize; e++ {
			result = append(result, 1)
		}
	}
	return result
}

func InvokeEvents(receiver events.DataEventReceiver, events ...*TEvent) {
	for _, event := range events {
		event.Invoke(receiver)
	}
}

func CloneBytes(bytes []byte) []byte {
	bytesCopy := make([]byte, len(bytes))
	copy(bytesCopy, bytes)
	return bytesCopy
}

// ----------------------------------------------------------------------------
// Events
// ----------------------------------------------------------------------------

var (
	EvBD      = BD()
	EvED      = ED()
	EvV       = V(version.ConciseEncodingVersion)
	EvPAD     = PAD(1)
	EvCOM     = COM(false, "a")
	EvN       = NULL()
	EvB       = B(true)
	EvTT      = TT()
	EvFF      = FF()
	EvPI      = PI(1)
	EvNI      = NI(1)
	EvI       = I(0)
	EvBI      = BI(NewBigInt("1"))
	EvBINull  = BI(nil)
	EvF       = BF(0x1.0p-1)
	EvFNAN    = BF(math.NaN())
	EvBF      = BBF(NewBigFloat("0x0.1"))
	EvBFNull  = BBF(nil)
	EvDF      = DF(NewDFloat("0.1"))
	EvDFNAN   = DF(NewDFloat("nan"))
	EvBDF     = BDF(NewBDF("0.1"))
	EvBDFNull = BDF(nil)
	EvBDFNAN  = BDF(NewBDF("nan"))
	EvNAN     = NAN()
	EvUID     = UID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	EvGT      = GT(time.Date(2020, time.Month(1), 1, 1, 1, 1, 1, time.UTC))
	EvCT      = T(NewDate(2020, 1, 1))
	EvL       = L()
	EvM       = M()
	EvMUP     = MU("a")
	EvNODE    = NODE()
	EvEDGE    = EDGE()
	EvE       = E()
	EvMARK    = MARK("a")
	EvREF     = REF("a")
	EvRREF    = RREF("a")
	EvAC      = AC(1, false)
	EvAD      = AD([]byte{1})
	EvS       = S("a")
	EvSB      = SB()
	EvRID     = RID("http://z.com")
	EvRB      = RB()
	EvCUB     = CB([]byte{1})
	EvCBB     = CBB()
	EvCUT     = CT("a")
	EvCTB     = CTB()
	EvAB      = AB(1, []byte{1})
	EvABB     = ABB()
	EvAU8     = AU8([]uint8{1})
	EvAU8B    = AU8B()
	EvAU16    = AU16([]uint16{1})
	EvAU16B   = AU16B()
	EvAU32    = AU32([]uint32{1})
	EvAU32B   = AU32B()
	EvAU64    = AU64([]uint64{1})
	EvAU64B   = AU64B()
	EvAI8     = AI8([]int8{1})
	EvAI8B    = AI8B()
	EvAI16    = AI16([]int16{1})
	EvAI16B   = AI16B()
	EvAI32    = AI32([]int32{1})
	EvAI32B   = AI32B()
	EvAI64    = AI64([]int64{1})
	EvAI64B   = AI64B()
	EvAF16    = AF16([]byte{1, 2})
	EvAF16B   = AF16B()
	EvAF32    = AF32([]float32{1})
	EvAF32B   = AF32B()
	EvAF64    = AF64([]float64{1})
	EvAF64B   = AF64B()
	EvAUU     = AU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvAUUB    = AUB()
	EvMB      = MB()
)

var allEvents = []*TEvent{
	EvBD, EvED, EvV, EvPAD, EvCOM, EvN, EvB, EvTT, EvFF, EvPI, EvNI, EvI,
	EvBI, EvBINull, EvF, EvFNAN, EvBF, EvBFNull, EvDF, EvDFNAN, EvBDF, EvBDFNull,
	EvBDFNAN, EvNAN, EvUID, EvGT, EvCT, EvL, EvM, EvMUP, EvNODE, EvEDGE, EvE,
	EvMARK, EvREF, EvRREF, EvAC, EvAD, EvS, EvSB, EvRID, EvRB,
	EvCUB, EvCBB, EvCUT, EvCTB, EvAB, EvABB, EvAU8, EvAU8B, EvAU16,
	EvAU16B, EvAU32, EvAU32B, EvAU64, EvAU64B, EvAI8, EvAI8B, EvAI16, EvAI16B,
	EvAI32, EvAI32B, EvAI64, EvAI64B, EvAF16, EvAF16B, EvAF32, EvAF32B, EvAF64,
	EvAF64B, EvAUU, EvAUUB, EvMB,
}

func FilterAllEvents(events []*TEvent, filter func(*TEvent) []*TEvent) []*TEvent {
	filtered := []*TEvent{}
	for _, event := range events {
		filtered = append(filtered, filter(event)...)
	}
	return filtered
}

func FilterCTE(event *TEvent) []*TEvent {
	switch event.Type {
	case TEventPadding:
		return []*TEvent{}
	case TEventArrayBooleanBegin, TEventArrayFloat16Begin,
		TEventArrayFloat32Begin, TEventArrayFloat64Begin,
		TEventArrayInt8Begin, TEventArrayInt16Begin, TEventArrayInt32Begin,
		TEventArrayInt64Begin, TEventArrayUint8Begin, TEventArrayUint16Begin,
		TEventArrayUint32Begin, TEventArrayUint64Begin, TEventArrayUIDBegin,
		TEventCustomBinaryBegin, TEventCustomTextBegin, TEventResourceIDBegin,
		TEventRemoteRefBegin, TEventStringBegin:
		return []*TEvent{S("x")}
	case TEventArrayChunk, TEventArrayData:
		return []*TEvent{}
	default:
		return []*TEvent{event}
	}
}

func FilterContainer(event *TEvent) []*TEvent {
	switch event.Type {
	case TEventEnd:
		return []*TEvent{}
	default:
		return []*TEvent{event}
	}
}

func FilterKey(event *TEvent) []*TEvent {
	switch event.Type {
	case TEventEnd, TEventReference:
		return []*TEvent{}
	default:
		return []*TEvent{event}
	}
}

func FilterMarker(event *TEvent) []*TEvent {
	switch event.Type {
	case TEventComment, TEventMarker, TEventReference:
		return []*TEvent{}
	default:
		return []*TEvent{event}
	}
}

func FilterEventsForCTE(events []*TEvent) []*TEvent {
	return FilterAllEvents(events, FilterCTE)
}

func FilterEventsForMarker(events []*TEvent) []*TEvent {
	return FilterAllEvents(events, FilterMarker)
}

func FilterEventsForContainer(events []*TEvent) []*TEvent {
	return FilterAllEvents(events, FilterContainer)
}

func FilterEventsForKey(events []*TEvent) []*TEvent {
	return FilterAllEvents(events, FilterKey)
}

func ComplementaryEvents(events []*TEvent) []*TEvent {
	complementary := make([]*TEvent, 0, len(allEvents)/2)
	for _, event := range allEvents {
		for _, compareEvent := range events {
			if event == compareEvent {
				goto Skip
			}
		}
		complementary = append(complementary, event)
	Skip:
	}
	return complementary
}

var (
	ArrayBeginTypes = []*TEvent{
		EvSB, EvRB, EvCBB, EvCTB, EvABB, EvAU8B, EvAU16B, EvAU32B, EvAU64B,
		EvAI8B, EvAI16B, EvAI32B, EvAI64B, EvAF16B, EvAF32B, EvAF64B, EvAUUB, EvMB,
	}

	StringArrayBeginTypes = []*TEvent{
		EvSB, EvRB, EvCTB,
	}

	NonStringArrayBeginTypes = []*TEvent{
		EvCBB, EvABB, EvAU8B, EvAU16B, EvAU32B, EvAU64B,
		EvAI8B, EvAI16B, EvAI32B, EvAI64B, EvAF16B, EvAF32B, EvAF64B, EvAUUB, EvMB,
	}

	ValidTLOValues   = ComplementaryEvents(InvalidTLOValues)
	InvalidTLOValues = []*TEvent{EvBD, EvED, EvV, EvE, EvAC, EvAD, EvREF}

	ValidMapKeys = []*TEvent{
		EvPAD, EvCOM, EvB, EvTT, EvFF, EvPI, EvNI, EvI, EvBI, EvF, EvBF, EvDF,
		EvBDF, EvUID, EvGT, EvCT, EvMARK, EvS, EvSB, EvRID, EvRB,
		EvREF, EvE,
	}
	InvalidMapKeys = ComplementaryEvents(ValidMapKeys)

	ValidMapValues   = ComplementaryEvents(InvalidMapValues)
	InvalidMapValues = []*TEvent{EvBD, EvED, EvV, EvE, EvAC, EvAD}

	ValidListValues   = ComplementaryEvents(InvalidListValues)
	InvalidListValues = []*TEvent{EvBD, EvED, EvV, EvAC, EvAD}

	ValidMarkupContents   = []*TEvent{EvPAD, EvCOM, EvS, EvSB, EvMUP, EvE}
	InvalidMarkupContents = ComplementaryEvents(ValidMarkupContents)

	ValidAfterNonStringArrayBegin   = []*TEvent{EvAC, EvCOM}
	InvalidAfterNonStringArrayBegin = ComplementaryEvents(ValidAfterNonStringArrayBegin)

	ValidAfterStringArrayBegin   = []*TEvent{EvAC}
	InvalidAfterStringArrayBegin = ComplementaryEvents(ValidAfterStringArrayBegin)

	ValidAfterArrayChunk   = []*TEvent{EvAD}
	InvalidAfterArrayChunk = ComplementaryEvents(ValidAfterArrayChunk)

	ValidMarkerValues   = ComplementaryEvents(InvalidMarkerValues)
	InvalidMarkerValues = []*TEvent{EvBD, EvED, EvV, EvE, EvAC, EvAD, EvMARK, EvREF, EvRREF}

	Padding                     = []*TEvent{EvPAD}
	CommentsPaddingRefEnd       = []*TEvent{EvPAD, EvCOM, EvREF, EvE}
	CommentsPaddingMarkerRefEnd = []*TEvent{EvPAD, EvCOM, EvMARK, EvREF, EvE}

	ValidEdgeSources   = ComplementaryEvents(InvalidEdgeSources)
	InvalidEdgeSources = []*TEvent{EvBD, EvED, EvV, EvAC, EvAD, EvN, EvBDFNull, EvBFNull, EvBINull}

	ValidEdgeDescriptions   = ValidListValues
	InvalidEdgeDescriptions = InvalidListValues

	ValidEdgeDestinations    = ValidEdgeSources
	InvalidOEdgeDestinations = InvalidEdgeSources

	ValidNodeValues   = ComplementaryEvents(InvalidNodeValues)
	InvalidNodeValues = []*TEvent{EvBD, EvED, EvV, EvAC, EvAD}
)

func containsEvent(events []*TEvent, event *TEvent) bool {
	for _, e := range events {
		if e.Type == event.Type {
			return true
		}
	}
	return false
}

func RemoveEvents(srcEvents []*TEvent, disallowedEvents []*TEvent) (events []*TEvent) {
	for _, event := range srcEvents {
		if !containsEvent(disallowedEvents, event) {
			events = append(events, event)
		}
	}
	return
}

func copyEvents(events []*TEvent) []*TEvent {
	newEvents := make([]*TEvent, len(events))
	copy(newEvents, events)
	return newEvents
}

var basicCompletions = map[TEventType][]*TEvent{
	TEventList:              {E()},
	TEventMap:               {E()},
	TEventMarkup:            {E(), E()},
	TEventNode:              {I(1), E()},
	TEventEdge:              {RID("a"), RID("b"), I(1)},
	TEventStringBegin:       {AC(1, false), AD([]byte{'a'})},
	TEventResourceIDBegin:   {AC(1, false), AD([]byte{'a'})},
	TEventRemoteRefBegin:    {AC(1, false), AD([]byte{'a'})},
	TEventCustomBinaryBegin: {AC(1, false), AD([]byte{1})},
	TEventCustomTextBegin:   {AC(1, false), AD([]byte{'a'})},
	TEventArrayBooleanBegin: {AC(1, false), AD([]byte{1})},
	TEventArrayUint8Begin:   {AC(1, false), AD(GenerateArrayElements(1, 1))},
	TEventArrayUint16Begin:  {AC(1, false), AD(GenerateArrayElements(2, 1))},
	TEventArrayUint32Begin:  {AC(1, false), AD(GenerateArrayElements(4, 1))},
	TEventArrayUint64Begin:  {AC(1, false), AD(GenerateArrayElements(8, 1))},
	TEventArrayInt8Begin:    {AC(1, false), AD(GenerateArrayElements(1, 1))},
	TEventArrayInt16Begin:   {AC(1, false), AD(GenerateArrayElements(2, 1))},
	TEventArrayInt32Begin:   {AC(1, false), AD(GenerateArrayElements(4, 1))},
	TEventArrayInt64Begin:   {AC(1, false), AD(GenerateArrayElements(8, 1))},
	TEventArrayFloat16Begin: {AC(1, false), AD(GenerateArrayElements(2, 1))},
	TEventArrayFloat32Begin: {AC(1, false), AD(GenerateArrayElements(4, 1))},
	TEventArrayFloat64Begin: {AC(1, false), AD(GenerateArrayElements(8, 1))},
	TEventArrayUIDBegin:     {AC(1, false), AD(GenerateArrayElements(16, 1))},
	TEventMediaBegin:        {AC(1, false), AD([]byte{'a'}), AC(0, false)},
}

func getBasicCompletion(stream []*TEvent) []*TEvent {
	if len(stream) == 0 {
		return []*TEvent{}
	}
	lastEvent := stream[len(stream)-1]
	return basicCompletions[lastEvent.Type]
}

func allPossibleEventStreams(
	docBegin []*TEvent,
	prefix []*TEvent,
	suffix []*TEvent,
	docEnd []*TEvent,
	event *TEvent,
	possibleFollowups []*TEvent) (allEvents [][]*TEvent) {

	switch event.Type {
	case TEventMarker:
		for _, following := range RemoveEvents(RemoveEvents(possibleFollowups, InvalidMarkerValues), CommentsPaddingMarkerRefEnd) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, event)
			newStream = append(newStream, following)
			newStream = append(newStream, getBasicCompletion(newStream)...)
			newStream = append(newStream, suffix...)
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}
	case TEventReference:
		for _, following := range RemoveEvents(RemoveEvents(possibleFollowups, InvalidMarkerValues), CommentsPaddingMarkerRefEnd) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, L(), MARK("a"))
			newStream = append(newStream, following)
			newStream = append(newStream, getBasicCompletion(newStream)...)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, REF("a"))
			newStream = append(newStream, suffix...)
			newStream = append(newStream, E())
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}

	case TEventPadding:
		for _, following := range RemoveEvents(possibleFollowups, CommentsPaddingMarkerRefEnd) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, event)
			newStream = append(newStream, following)
			newStream = append(newStream, getBasicCompletion(newStream)...)
			newStream = append(newStream, suffix...)
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}
	case TEventArrayBooleanBegin, TEventArrayFloat16Begin,
		TEventArrayFloat32Begin, TEventArrayFloat64Begin,
		TEventArrayInt8Begin, TEventArrayInt16Begin, TEventArrayInt32Begin,
		TEventArrayInt64Begin, TEventArrayUint8Begin, TEventArrayUint16Begin,
		TEventArrayUint32Begin, TEventArrayUint64Begin, TEventArrayUIDBegin,
		TEventCustomBinaryBegin, TEventCustomTextBegin, TEventResourceIDBegin,
		TEventRemoteRefBegin, TEventStringBegin:
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, AC(0, false))
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	case TEventArrayChunk, TEventArrayData:
		// TODO: Implement this better somehow?
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	case TEventEnd, TEventComment:
		// Skip
	default:
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, getBasicCompletion(newStream)...)
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	}
	return
}

func GenerateAllVariants(
	docBegin []*TEvent,
	prefix []*TEvent,
	suffix []*TEvent,
	docEnd []*TEvent,
	possibleFollowups []*TEvent) (events [][]*TEvent) {

	for _, event := range possibleFollowups {
		events = append(events, allPossibleEventStreams(docBegin, prefix, suffix, docEnd, event, possibleFollowups)...)
	}

	return
}

func EventOrNull(eventType TEventType, value interface{}) *TEvent {
	if value == nil {
		eventType = TEventNull
	}
	return NewTEvent(eventType, value, nil)
}

// ----------------------------------------------------------------------------
// Stored event convenience constructors
// ----------------------------------------------------------------------------

func TT() *TEvent                       { return NewTEvent(TEventTrue, nil, nil) }
func FF() *TEvent                       { return NewTEvent(TEventFalse, nil, nil) }
func I(v int64) *TEvent                 { return NewTEvent(TEventInt, v, nil) }
func BF(v float64) *TEvent              { return NewTEvent(TEventFloat, v, nil) }
func BBF(v *big.Float) *TEvent          { return EventOrNull(TEventBigFloat, v) }
func DF(v compact_float.DFloat) *TEvent { return NewTEvent(TEventDecimalFloat, v, nil) }
func BDF(v *apd.Decimal) *TEvent        { return EventOrNull(TEventBigDecimalFloat, v) }
func V(v uint64) *TEvent                { return NewTEvent(TEventVersion, v, nil) }
func NULL() *TEvent                     { return NewTEvent(TEventNull, nil, nil) }
func PAD(v int) *TEvent                 { return NewTEvent(TEventPadding, v, nil) }
func COM(m bool, v string) *TEvent      { return NewTEvent(TEventComment, m, v) }
func B(v bool) *TEvent                  { return NewTEvent(TEventBool, v, nil) }
func PI(v uint64) *TEvent               { return NewTEvent(TEventPInt, v, nil) }
func NI(v uint64) *TEvent               { return NewTEvent(TEventNInt, v, nil) }
func BI(v *big.Int) *TEvent             { return EventOrNull(TEventBigInt, v) }
func NAN() *TEvent                      { return NewTEvent(TEventNan, nil, nil) }
func SNAN() *TEvent                     { return NewTEvent(TEventSNan, nil, nil) }
func UID(v []byte) *TEvent              { return NewTEvent(TEventUID, v, nil) }
func GT(v time.Time) *TEvent            { return NewTEvent(TEventTime, v, nil) }
func T(v compact_time.Time) *TEvent     { return EventOrNull(TEventCompactTime, v) }
func S(v string) *TEvent                { return NewTEvent(TEventString, v, nil) }
func RID(v string) *TEvent              { return NewTEvent(TEventResourceID, v, nil) }
func RREF(v string) *TEvent             { return NewTEvent(TEventRemoteRef, v, nil) }
func CB(v []byte) *TEvent               { return NewTEvent(TEventCustomBinary, v, nil) }
func CT(v string) *TEvent               { return NewTEvent(TEventCustomText, v, nil) }
func AB(l uint64, v []byte) *TEvent     { return NewTEvent(TEventArrayBoolean, l, v) }
func AI8(v []int8) *TEvent              { return NewTEvent(TEventArrayInt8, v, nil) }
func AI16(v []int16) *TEvent            { return NewTEvent(TEventArrayInt16, v, nil) }
func AI32(v []int32) *TEvent            { return NewTEvent(TEventArrayInt32, v, nil) }
func AI64(v []int64) *TEvent            { return NewTEvent(TEventArrayInt64, v, nil) }
func AU8(v []byte) *TEvent              { return NewTEvent(TEventArrayUint8, v, nil) }
func AU16(v []uint16) *TEvent           { return NewTEvent(TEventArrayUint16, v, nil) }
func AU32(v []uint32) *TEvent           { return NewTEvent(TEventArrayUint32, v, nil) }
func AU64(v []uint64) *TEvent           { return NewTEvent(TEventArrayUint64, v, nil) }
func AF16(v []byte) *TEvent             { return NewTEvent(TEventArrayFloat16, v, nil) }
func AF32(v []float32) *TEvent          { return NewTEvent(TEventArrayFloat32, v, nil) }
func AF64(v []float64) *TEvent          { return NewTEvent(TEventArrayFloat64, v, nil) }
func AU(v [][]byte) *TEvent             { return NewTEvent(TEventArrayUID, v, nil) }
func SB() *TEvent                       { return NewTEvent(TEventStringBegin, nil, nil) }
func RB() *TEvent                       { return NewTEvent(TEventResourceIDBegin, nil, nil) }
func RRB() *TEvent                      { return NewTEvent(TEventRemoteRefBegin, nil, nil) }
func CBB() *TEvent                      { return NewTEvent(TEventCustomBinaryBegin, nil, nil) }
func CTB() *TEvent                      { return NewTEvent(TEventCustomTextBegin, nil, nil) }
func ABB() *TEvent                      { return NewTEvent(TEventArrayBooleanBegin, nil, nil) }
func AI8B() *TEvent                     { return NewTEvent(TEventArrayInt8Begin, nil, nil) }
func AI16B() *TEvent                    { return NewTEvent(TEventArrayInt16Begin, nil, nil) }
func AI32B() *TEvent                    { return NewTEvent(TEventArrayInt32Begin, nil, nil) }
func AI64B() *TEvent                    { return NewTEvent(TEventArrayInt64Begin, nil, nil) }
func AU8B() *TEvent                     { return NewTEvent(TEventArrayUint8Begin, nil, nil) }
func AU16B() *TEvent                    { return NewTEvent(TEventArrayUint16Begin, nil, nil) }
func AU32B() *TEvent                    { return NewTEvent(TEventArrayUint32Begin, nil, nil) }
func AU64B() *TEvent                    { return NewTEvent(TEventArrayUint64Begin, nil, nil) }
func AF16B() *TEvent                    { return NewTEvent(TEventArrayFloat16Begin, nil, nil) }
func AF32B() *TEvent                    { return NewTEvent(TEventArrayFloat32Begin, nil, nil) }
func AF64B() *TEvent                    { return NewTEvent(TEventArrayFloat64Begin, nil, nil) }
func AUB() *TEvent                      { return NewTEvent(TEventArrayUIDBegin, nil, nil) }
func MB() *TEvent                       { return NewTEvent(TEventMediaBegin, nil, nil) }
func AC(l uint64, more bool) *TEvent    { return NewTEvent(TEventArrayChunk, l, more) }
func AD(v []byte) *TEvent               { return NewTEvent(TEventArrayData, v, nil) }
func L() *TEvent                        { return NewTEvent(TEventList, nil, nil) }
func M() *TEvent                        { return NewTEvent(TEventMap, nil, nil) }
func MU(id string) *TEvent              { return NewTEvent(TEventMarkup, id, nil) }
func NODE() *TEvent                     { return NewTEvent(TEventNode, nil, nil) }
func EDGE() *TEvent                     { return NewTEvent(TEventEdge, nil, nil) }
func E() *TEvent                        { return NewTEvent(TEventEnd, nil, nil) }
func MARK(id string) *TEvent            { return NewTEvent(TEventMarker, id, nil) }
func REF(id string) *TEvent             { return NewTEvent(TEventReference, id, nil) }
func CONST(n string) *TEvent            { return NewTEvent(TEventConstant, n, nil) }
func BD() *TEvent                       { return NewTEvent(TEventBeginDocument, nil, nil) }
func ED() *TEvent                       { return NewTEvent(TEventEndDocument, nil, nil) }

// Converts a go value into a stored event
func EventForValue(value interface{}) *TEvent {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return NULL()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return B(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return I(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return PI(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return BF(rv.Float())
	case reflect.String:
		return S(rv.String())
	case reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.Uint8:
			return AU8(rv.Bytes())
		}
	case reflect.Ptr:
		if rv.IsNil() {
			return NULL()
		}
		switch rv.Type() {
		case common.TypePBigDecimalFloat:
			return BDF(rv.Interface().(*apd.Decimal))
		case common.TypePBigFloat:
			return BBF(rv.Interface().(*big.Float))
		case common.TypePBigInt:
			return BI(rv.Interface().(*big.Int))
		case common.TypePURL:
			return RID(rv.Interface().(*url.URL).String())
		}
		return EventForValue(rv.Elem().Interface())
	case reflect.Struct:
		switch rv.Type() {
		case common.TypeBigDecimalFloat:
			v := rv.Interface().(apd.Decimal)
			return BDF(&v)
		case common.TypeBigFloat:
			v := rv.Interface().(big.Float)
			return BBF(&v)
		case common.TypeBigInt:
			v := rv.Interface().(big.Int)
			return BI(&v)
		case common.TypeCompactTime:
			v := rv.Interface().(compact_time.Time)
			return T(v)
		case common.TypeDFloat:
			v := rv.Interface().(compact_float.DFloat)
			return DF(v)
		case common.TypeTime:
			v := rv.Interface().(time.Time)
			return GT(v)
		case common.TypeURL:
			v := rv.Interface().(url.URL)
			return RID(v.String())
		}
	case reflect.Array:
		if rv.Type() == common.TypeUID {
			v := value.(types.UID)
			return UID(v[:])
		}
	}
	panic(fmt.Errorf("TEST CODE BUG: Unhandled kind: %v", rv.Kind()))
}

// ----------------------------------------------------------------------------
// Testing structures
// ----------------------------------------------------------------------------

// These are just complex structures used by a lot of the subsystem tests.

type TestingInnerStruct struct {
	Inner int
}

type TestingOuterStruct struct {
	Bo     bool
	PBo    *bool
	By     byte
	PBy    *byte
	I      int
	PI     *int
	I8     int8
	PI8    *int8
	I16    int16
	PI16   *int16
	I32    int32
	PI32   *int32
	I64    int64
	PI64   *int64
	U      uint
	PU     *uint
	U8     uint8
	PU8    *uint8
	U16    uint16
	PU16   *uint16
	U32    uint32
	PU32   *uint32
	U64    uint64
	PU64   *uint64
	BI     big.Int
	PBI    *big.Int
	F32    float32
	PF32   *float32
	F64    float64
	PF64   *float64
	BF     big.Float
	PBF    *big.Float
	DF     compact_float.DFloat
	BDF    apd.Decimal
	PBDF   *apd.Decimal
	St     string
	Au8    [4]byte
	Su8    []byte
	Sl     []interface{}
	M      map[interface{}]interface{}
	IS     TestingInnerStruct
	PIS    *TestingInnerStruct
	Time   time.Time
	PTime  *time.Time
	CTime  compact_time.Time
	PCTime compact_time.Time
	PURL   *url.URL
	URL    url.URL
	UID    types.UID
}

func (_this *TestingOuterStruct) GetRepresentativeEvents(includeFakes bool) (events []*TEvent) {
	ade := func(e ...*TEvent) {
		events = append(events, e...)
	}
	adv := func(value interface{}) {
		ade(EventForValue(value))
	}
	anv := func(name string, value interface{}) {
		adv(name)
		adv(value)
	}
	ane := func(name string, e ...*TEvent) {
		adv(name)
		ade(e...)
	}

	ade(M())

	anv("Bo", _this.Bo)
	anv("PBo", _this.PBo)
	anv("By", _this.By)
	anv("PBy", _this.PBy)
	anv("I", _this.I)
	anv("PI", _this.PI)
	anv("I8", _this.I8)
	anv("PI8", _this.PI8)
	anv("I16", _this.I16)
	anv("PI16", _this.PI16)
	anv("I32", _this.I32)
	anv("PI32", _this.PI32)
	anv("I64", _this.I64)
	anv("PI64", _this.PI64)
	anv("U", _this.U)
	anv("PU", _this.PU)
	anv("U8", _this.U8)
	anv("PU8", _this.PU8)
	anv("U16", _this.U16)
	anv("PU16", _this.PU16)
	anv("U32", _this.U32)
	anv("PU32", _this.PU32)
	anv("U64", _this.U64)
	anv("PU64", _this.PU64)
	anv("BI", _this.BI)
	anv("PBI", _this.PBI)
	anv("F32", _this.F32)
	anv("PF32", _this.PF32)
	anv("F64", _this.F64)
	anv("PF64", _this.PF64)
	anv("BF", _this.BF)
	anv("PBF", _this.PBF)
	anv("DF", _this.DF)
	anv("BDF", _this.BDF)
	anv("PBDF", _this.PBDF)
	anv("St", _this.St)
	ane("Au8", AU8(_this.Au8[:]))
	anv("Su8", _this.Su8)

	ane("Sl", L())
	for _, v := range _this.Sl {
		adv(v)
	}
	ade(E())

	ane("M", M())
	for k, v := range _this.M {
		adv(k)
		adv(v)
	}
	ade(E())

	ane("IS", M())
	anv("Inner", _this.IS.Inner)
	ade(E())

	if _this.PIS != nil {
		ane("PIS", M())
		anv("Inner", _this.PIS.Inner)
		ade(E())
	}

	anv("Time", _this.Time)
	anv("PTime", _this.PTime)
	anv("CTime", _this.CTime)
	anv("PCTime", _this.PCTime)
	anv("PURL", _this.PURL)
	anv("UID", _this.UID)

	if includeFakes {
		ane("F1", B(true))
		ane("F2", B(false))
		ane("F3", I(1))
		ane("F4", I(-1))
		ane("F5", BF(1.1))
		ane("F6", BBF(NewBigFloat("1.1")))
		ane("F7", DF(NewDFloat("1.1")))
		ane("F8", BDF(NewBDF("1.1")))
		ane("F9", NULL())
		ane("F10", BI(NewBigInt("1000")))
		ane("F12", NAN())
		ane("F13", SNAN())
		ane("F14", UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
		ane("F15", GT(_this.Time))
		ane("F16", T(_this.CTime))
		ane("F17", AU8([]byte{1}))
		ane("F18", S("xyz"))
		ane("F19", RID("http://example.com"))
		// ane("F20", cust([]byte{1}))
		ane("FakeList", L(), E())
		ane("FakeMap", M(), E())
		ane("FakeDeep", L(), M(), S("A"), L(),
			B(true),
			B(false),
			I(1),
			I(-1),
			BF(1.1),
			BBF(NewBigFloat("1.1")),
			DF(NewDFloat("1.1")),
			BDF(NewBDF("1.1")),
			NULL(),
			BI(NewBigInt("1000")),
			NAN(),
			SNAN(),
			UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
			GT(_this.Time),
			T(_this.CTime),
			AU8([]byte{1}),
			S("xyz"),
			RID("http://example.com"),
			// cust([]byte{1}),
			E(), E(), E())
	}

	ade(E())
	return
}

func NewTestingOuterStruct(baseValue int) *TestingOuterStruct {
	_this := new(TestingOuterStruct)
	_this.Init(baseValue)
	return _this
}

func NewBlankTestingOuterStruct() *TestingOuterStruct {
	_this := new(TestingOuterStruct)
	_this.CTime.Year = 1
	_this.CTime.Month = 1
	_this.CTime.Day = 1
	return _this
}

func (_this *TestingOuterStruct) Init(baseValue int) {
	_this.Bo = baseValue&1 == 1
	_this.PBo = &_this.Bo
	_this.By = byte(baseValue + int(unsafe.Offsetof(_this.By)))
	_this.PBy = &_this.By
	_this.I = 100000 + baseValue + int(unsafe.Offsetof(_this.I))
	_this.PI = &_this.I
	_this.I8 = int8(baseValue + int(unsafe.Offsetof(_this.I8)))
	_this.PI8 = &_this.I8
	_this.I16 = int16(-1000 - baseValue - int(unsafe.Offsetof(_this.I16)))
	_this.PI16 = &_this.I16
	_this.I32 = int32(1000000000 + baseValue + int(unsafe.Offsetof(_this.I32)))
	_this.PI32 = &_this.I32
	_this.I64 = int64(1000000000000000) + int64(baseValue+int(unsafe.Offsetof(_this.I64)))
	_this.PI64 = &_this.I64
	_this.U = uint(1000000 + baseValue + int(unsafe.Offsetof(_this.U)))
	_this.PU = &_this.U
	_this.U8 = uint8(baseValue + int(unsafe.Offsetof(_this.U8)))
	_this.PU8 = &_this.U8
	_this.U16 = uint16(10000 + baseValue + int(unsafe.Offsetof(_this.U16)))
	_this.PU16 = &_this.U16
	_this.U32 = uint32(100000000 + baseValue + int(unsafe.Offsetof(_this.U32)))
	_this.PU32 = &_this.U32
	_this.U64 = uint64(1000000000000) + uint64(baseValue+int(unsafe.Offsetof(_this.U64)))
	_this.PU64 = &_this.U64
	_this.PBI = NewBigInt(fmt.Sprintf("-10000000000000000000000000000000000000%v", unsafe.Offsetof(_this.PBI)))
	_this.BI = *_this.PBI
	_this.F32 = float32(1000000+baseValue+int(unsafe.Offsetof(_this.F32))) + 0.5
	_this.PF32 = &_this.F32
	_this.F64 = float64(1000000000000) + float64(baseValue+int(unsafe.Offsetof(_this.F64))) + 0.5
	_this.PF64 = &_this.F64
	_this.PBF = NewBigFloat("12345678901234567890123.1234567")
	_this.BF = *_this.PBF
	_this.DF = NewDFloat(fmt.Sprintf("-100000000000000%ve-1000000", unsafe.Offsetof(_this.DF)))
	_this.PBDF = NewBDF("-1.234567890123456789777777777777777777771234e-10000")
	_this.BDF = *_this.PBDF
	_this.St = GenerateString(baseValue+5, baseValue)
	_this.Au8[0] = byte(baseValue + int(unsafe.Offsetof(_this.Au8)))
	_this.Au8[1] = byte(baseValue + int(unsafe.Offsetof(_this.Au8)+1))
	_this.Au8[2] = byte(baseValue + int(unsafe.Offsetof(_this.Au8)+2))
	_this.Au8[3] = byte(baseValue + int(unsafe.Offsetof(_this.Au8)+3))
	_this.Su8 = GenerateBytes(baseValue+1, baseValue)
	_this.M = make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		_this.Sl = append(_this.Sl, i)
		_this.M[fmt.Sprintf("key%v", i)] = i
	}
	_this.IS.Inner = baseValue + 15
	_this.PIS = new(TestingInnerStruct)
	_this.PIS.Inner = baseValue + 16
	testTime := time.Date(30000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	_this.PTime = &testTime
	_this.CTime = NewTS(-1000, 1, 1, 1, 1, 1, 1, "Europe/Berlin")
	_this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
	_this.UID = types.UID{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x90, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0, 0xff}
}
