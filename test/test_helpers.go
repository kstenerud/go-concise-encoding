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
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"
	"github.com/kstenerud/go-concise-encoding/version"
)

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
	if err := ReportPanic(function); err != nil {
		t.Errorf("Error [%v] in %v", err, name)
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

func InvokeEventsAsCompleteDocument(receiver events.DataEventReceiver, events ...Event) {
	BD().Invoke(receiver)
	for _, event := range events {
		event.Invoke(receiver)
	}
	ED().Invoke(receiver)
}

// ----------------------------------------------------------------------------
// Events
// ----------------------------------------------------------------------------

var (
	EvAB      = AB([]bool{true})
	EvAC      = ACL(1)
	EvAD      = ADU8([]byte{1})
	EvAF16    = AF16([]float32{1})
	EvAF32    = AF32([]float32{1})
	EvAF64    = AF64([]float64{1})
	EvAI16    = AI16([]int16{1})
	EvAI32    = AI32([]int32{1})
	EvAI64    = AI64([]int64{1})
	EvAI8     = AI8([]int8{1})
	EvAU16    = AU16([]uint16{1})
	EvAU32    = AU32([]uint32{1})
	EvAU64    = AU64([]uint64{1})
	EvAU8     = AU8([]uint8{1})
	EvAU      = AU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvB       = B(true)
	EvBAB     = BAB()
	EvBAF16   = BAF16()
	EvBAF32   = BAF32()
	EvBAF64   = BAF64()
	EvBAI16   = BAI16()
	EvBAI32   = BAI32()
	EvBAI64   = BAI64()
	EvBAI8    = BAI8()
	EvBAU16   = BAU16()
	EvBAU32   = BAU32()
	EvBAU64   = BAU64()
	EvBAU8    = BAU8()
	EvBAU     = BAU()
	EvBBF     = N(NewBigFloat("0x0.1"))
	EvBBFNull = N(nil)
	EvBDF     = N(NewBDF("0.1"))
	EvBDFNAN  = N(NewBDF("nan"))
	EvBDFNull = N(nil)
	EvBF      = N(0x1.0p-1)
	EvBI      = N(NewBigInt("1"))
	EvBINull  = N(nil)
	EvBCB     = BCB(1)
	EvBCT     = BCT(1)
	EvCS      = CS("a")
	EvCM      = CM("a")
	EvTC      = T(NewDate(2020, 1, 1))
	EvCB      = CB(1, []byte{1})
	EvCT      = CT(1, "a")
	EvDF      = N(NewDFloat("0.1"))
	EvDFNAN   = N(compact_float.SignalingNaN())
	EvE       = E()
	EvEDGE    = EDGE()
	EvFNAN    = N(math.NaN())
	EvT       = T(compact_time.AsCompactTime(time.Date(2020, time.Month(1), 1, 1, 1, 1, 1, time.UTC)))
	EvI       = N(0)
	EvL       = L()
	EvM       = M()
	EvMARK    = MARK("a")
	EvMEDIA   = MEDIA("a/b", []byte{0})
	EvBMEDIA  = BMEDIA("a/b")
	EvNAN     = N(compact_float.QuietNaN())
	EvSNAN    = N(compact_float.SignalingNaN())
	EvNI      = N(-1)
	EvNODE    = NODE()
	EvNULL    = NULL()
	EvPAD     = PAD()
	EvPI      = N(1)
	EvBRID    = BRID()
	EvREFL    = REFL("a")
	EvREFR    = REFR("a")
	EvRID     = RID("http://z.com")
	EvS       = S("a")
	EvSB      = BS()
	EvREC     = REC("a")
	EvRT      = RT("a")
	EvUID     = UID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	EvV       = V(version.ConciseEncodingVersion)
)

var allEvents = Events{
	EvV, EvPAD, EvCS, EvCM, EvNULL, EvB, EvPI, EvNI, EvI,
	EvBI, EvBINull, EvBF, EvFNAN, EvBBF, EvBBFNull, EvDF, EvDFNAN, EvBDF, EvBDFNull,
	EvBDFNAN, EvNAN, EvSNAN, EvUID, EvT, EvTC, EvL, EvM, EvREC,
	EvNODE, EvEDGE, EvE,
	EvMARK, EvREFL, EvREFR, EvAC, EvAD, EvS, EvSB, EvRID, EvBRID,
	EvCB, EvBCB, EvCT, EvBCT, EvAB, EvBAB, EvAU8, EvBAU8, EvAU16,
	EvBAU16, EvAU32, EvBAU32, EvAU64, EvBAU64, EvAI8, EvBAI8, EvAI16, EvBAI16,
	EvAI32, EvBAI32, EvAI64, EvBAI64, EvAF16, EvBAF16, EvAF32, EvBAF32, EvAF64,
	EvBAF64, EvAU, EvBAU, EvBMEDIA, EvMEDIA,
}

func ComplementaryEvents(events Events) Events {
	complementary := make(Events, 0, len(allEvents)/2)
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
	ArrayBeginTypes = Events{
		EvSB, EvBRID, EvBCB, EvBCT, EvBAB, EvBAU8, EvBAU16, EvBAU32, EvBAU64,
		EvBAI8, EvBAI16, EvBAI32, EvBAI64, EvBAF16, EvBAF32, EvBAF64, EvBAU, EvBMEDIA,
	}

	StringArrayBeginTypes = Events{
		EvSB, EvBRID, EvBCT,
	}

	NonStringArrayBeginTypes = Events{
		EvBCB, EvBAB, EvBAU8, EvBAU16, EvBAU32, EvBAU64,
		EvBAI8, EvBAI16, EvBAI32, EvBAI64, EvBAF16, EvBAF32, EvBAF64, EvBAU, EvBMEDIA,
	}

	ValidTLOValues   = ComplementaryEvents(InvalidTLOValues)
	InvalidTLOValues = Events{EvV, EvE, EvAC, EvAD, EvREFL, EvRT}

	ValidMapKeys = Events{
		EvPAD, EvCS, EvCM, EvB, EvPI, EvNI, EvI, EvBI,
		EvUID, EvT, EvTC, EvMARK, EvS, EvSB, EvRID, EvBRID,
		EvREFL, EvE,
	}
	InvalidMapKeys = ComplementaryEvents(ValidMapKeys)

	ValidMapValues   = ComplementaryEvents(InvalidMapValues)
	InvalidMapValues = Events{EvV, EvE, EvAC, EvAD, EvRT}

	ValidListValues   = ComplementaryEvents(InvalidListValues)
	InvalidListValues = Events{EvV, EvAC, EvAD, EvRT}

	ValidRecordTypeValues   = ComplementaryEvents(InvalidRecordTypeValues)
	InvalidRecordTypeValues = Events{EvV, EvAC, EvAD, EvMARK, EvREFL, EvRT}

	ValidRecordValues   = ComplementaryEvents(InvalidRecordValues)
	InvalidRecordValues = Events{EvV, EvAC, EvAD, EvRT}

	ValidAfterNonStringArrayBegin   = Events{EvAC, EvCS, EvCM}
	InvalidAfterNonStringArrayBegin = ComplementaryEvents(ValidAfterNonStringArrayBegin)

	ValidAfterStringArrayBegin   = Events{EvAC}
	InvalidAfterStringArrayBegin = ComplementaryEvents(ValidAfterStringArrayBegin)

	ValidAfterArrayChunk   = Events{EvAD}
	InvalidAfterArrayChunk = ComplementaryEvents(ValidAfterArrayChunk)

	ValidMarkerValues   = ComplementaryEvents(InvalidMarkerValues)
	InvalidMarkerValues = Events{EvV, EvE, EvAC, EvAD, EvMARK, EvREFL, EvREFR, EvRT}

	Padding                     = Events{EvPAD}
	CommentsPaddingRefEnd       = Events{EvPAD, EvCS, EvCM, EvREFL, EvE}
	CommentsPaddingMarkerRefEnd = Events{EvPAD, EvCS, EvCM, EvMARK, EvREFL, EvE}

	ValidEdgeSources   = ComplementaryEvents(InvalidEdgeSources)
	InvalidEdgeSources = Events{EvV, EvAC, EvAD, EvNULL, EvBDFNull, EvBBFNull, EvBINull, EvRT}

	ValidEdgeDescriptions   = ValidListValues
	InvalidEdgeDescriptions = InvalidListValues

	ValidEdgeDestinations    = ValidEdgeSources
	InvalidOEdgeDestinations = InvalidEdgeSources

	ValidNodeValues   = ComplementaryEvents(InvalidNodeValues)
	InvalidNodeValues = Events{EvV, EvAC, EvAD, EvRT}
)

func containsEventType(events Events, event Event) bool {
	for _, e := range events {
		if reflect.TypeOf(e) == reflect.TypeOf(event) {
			return true
		}
	}
	return false
}

func RemoveEvents(srcEvents Events, disallowedEvents ...Event) (events Events) {
	for _, event := range srcEvents {
		if !containsEventType(disallowedEvents, event) {
			events = append(events, event)
		}
	}
	return
}

func copyEvents(events Events) Events {
	newEvents := make(Events, len(events))
	copy(newEvents, events)
	return newEvents
}

var basicCompletions = map[reflect.Type]Events{
	reflect.TypeOf(L()): {E()},
	reflect.TypeOf(M()): {E()},
	// reflect.TypeOf(RT("a")):       {E(), N(1)},
	reflect.TypeOf(REC("a")):      {S("a"), E()},
	reflect.TypeOf(NODE()):        {N(1), E()},
	reflect.TypeOf(EDGE()):        {RID("a"), RID("b"), N(1), E()},
	reflect.TypeOf(BS()):          {ACL(1), ADT("a")},
	reflect.TypeOf(BRID()):        {ACL(1), ADT("a")},
	reflect.TypeOf(BREFR()):       {ACL(1), ADT("a")},
	reflect.TypeOf(BCB(0)):        {ACL(1), ADU8([]byte{1})},
	reflect.TypeOf(BCT(0)):        {ACL(1), ADT("a")},
	reflect.TypeOf(BAB()):         {ACL(1), ADB([]bool{true})},
	reflect.TypeOf(BAU8()):        {ACL(1), ADU8([]uint8{0})},
	reflect.TypeOf(BAU16()):       {ACL(1), ADU16([]uint16{0})},
	reflect.TypeOf(BAU32()):       {ACL(1), ADU32([]uint32{0})},
	reflect.TypeOf(BAU64()):       {ACL(1), ADU64([]uint64{0})},
	reflect.TypeOf(BAI8()):        {ACL(1), ADI8([]int8{0})},
	reflect.TypeOf(BAI16()):       {ACL(1), ADI16([]int16{0})},
	reflect.TypeOf(BAI32()):       {ACL(1), ADI32([]int32{0})},
	reflect.TypeOf(BAI64()):       {ACL(1), ADI64([]int64{0})},
	reflect.TypeOf(BAF16()):       {ACL(1), ADF16([]float32{0})},
	reflect.TypeOf(BAF32()):       {ACL(1), ADF32([]float32{0})},
	reflect.TypeOf(BAF64()):       {ACL(1), ADF64([]float64{0})},
	reflect.TypeOf(BAU()):         {ACL(1), ADU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})},
	reflect.TypeOf(BMEDIA("e/b")): {ACL(0)},
}

func GetBasicCompletion(stream Events) Events {
	if len(stream) == 0 {
		return Events{}
	}
	lastEvent := stream[len(stream)-1]
	return basicCompletions[reflect.TypeOf(lastEvent)]
}

func allPossibleEventStreams(
	docBegin Events,
	prefix Events,
	suffix Events,
	docEnd Events,
	event Event,
	possibleFollowups Events) (allEvents []Events) {

	switch event.(type) {
	case *EventMarker:
		for _, following := range RemoveEvents(RemoveEvents(possibleFollowups, InvalidMarkerValues...), CommentsPaddingMarkerRefEnd...) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, event)
			newStream = append(newStream, following)
			newStream = append(newStream, GetBasicCompletion(newStream)...)
			newStream = append(newStream, suffix...)
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}
	case *EventReferenceLocal:
		for _, following := range RemoveEvents(RemoveEvents(possibleFollowups, InvalidMarkerValues...), CommentsPaddingMarkerRefEnd...) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, L(), MARK("a"))
			newStream = append(newStream, following)
			newStream = append(newStream, GetBasicCompletion(newStream)...)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, REFL("a"))
			newStream = append(newStream, suffix...)
			newStream = append(newStream, E())
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}

	case *EventPadding:
		for _, following := range RemoveEvents(possibleFollowups, CommentsPaddingMarkerRefEnd...) {
			newStream := copyEvents(docBegin)
			newStream = append(newStream, prefix...)
			newStream = append(newStream, event)
			newStream = append(newStream, following)
			newStream = append(newStream, GetBasicCompletion(newStream)...)
			newStream = append(newStream, suffix...)
			newStream = append(newStream, docEnd...)
			allEvents = append(allEvents, newStream)
		}
	case *EventBeginArrayBit, *EventBeginArrayFloat16,
		*EventBeginArrayFloat32, *EventBeginArrayFloat64,
		*EventBeginArrayInt8, *EventBeginArrayInt16, *EventBeginArrayInt32,
		*EventBeginArrayInt64, *EventBeginArrayUint8, *EventBeginArrayUint16,
		*EventBeginArrayUint32, *EventBeginArrayUint64, *EventBeginArrayUID,
		*EventBeginCustomBinary, *EventBeginCustomText, *EventBeginResourceID,
		*EventBeginReferenceRemote, *EventBeginString:
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, ACL(0))
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	case *EventArrayChunkMore, *EventArrayChunkLast,
		*EventArrayDataBit, *EventArrayDataFloat16,
		*EventArrayDataFloat32, *EventArrayDataFloat64,
		*EventArrayDataInt8, *EventArrayDataInt16, *EventArrayDataInt32,
		*EventArrayDataInt64, *EventArrayDataUint8, *EventArrayDataUint16,
		*EventArrayDataUint32, *EventArrayDataUint64, *EventArrayDataUID,
		*EventArrayDataText:
		// TODO: Implement this better somehow?
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	case *EventEndContainer, *EventCommentSingleLine, *EventCommentMultiline:
		// Skip
	default:
		newStream := copyEvents(docBegin)
		newStream = append(newStream, prefix...)
		newStream = append(newStream, event)
		newStream = append(newStream, GetBasicCompletion(newStream)...)
		newStream = append(newStream, suffix...)
		newStream = append(newStream, docEnd...)
		allEvents = append(allEvents, newStream)
	}
	return
}

func GenerateAllVariants(
	docBegin Events,
	prefix Events,
	suffix Events,
	docEnd Events,
	possibleFollowups Events) (events []Events) {

	if containsRecords(docBegin, prefix, suffix, docEnd, possibleFollowups) {
		// Crude implementation: Add a basic record type to the top if necessary
		docBegin = append(docBegin, EvRT, EvS, EvE)
	}

	for _, event := range possibleFollowups {
		events = append(events, allPossibleEventStreams(docBegin, prefix, suffix, docEnd, event, possibleFollowups)...)
	}

	return
}

func containsRecords(eventSets ...Events) bool {
	for _, eventSet := range eventSets {
		for _, event := range eventSet {
			if event.IsEquivalentTo(EvREC) {
				return true
			}
		}
	}
	return false
}

// Converts a go value into a stored event
func EventForValue(value interface{}) Event {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return NULL()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return B(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return N(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return N(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return N(rv.Float())
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
			return N(rv.Interface().(*apd.Decimal))
		case common.TypePBigFloat:
			return N(rv.Interface().(*big.Float))
		case common.TypePBigInt:
			return N(rv.Interface().(*big.Int))
		case common.TypePURL:
			return RID(rv.Interface().(*url.URL).String())
		}
		return EventForValue(rv.Elem().Interface())
	case reflect.Struct:
		switch rv.Type() {
		case common.TypeBigDecimalFloat:
			v := rv.Interface().(apd.Decimal)
			return N(&v)
		case common.TypeBigFloat:
			v := rv.Interface().(big.Float)
			return N(&v)
		case common.TypeBigInt:
			v := rv.Interface().(big.Int)
			return N(&v)
		case common.TypeCompactTime:
			v := rv.Interface().(compact_time.Time)
			return T(v)
		case common.TypeDFloat:
			v := rv.Interface().(compact_float.DFloat)
			return N(v)
		case common.TypeTime:
			v := rv.Interface().(time.Time)
			return T(compact_time.AsCompactTime(v))
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

func (_this *TestingOuterStruct) GetRepresentativeEvents(includeFakes bool) (events Events) {
	ade := func(e ...Event) {
		events = append(events, e...)
	}
	adv := func(value interface{}) {
		ade(EventForValue(value))
	}
	anv := func(name string, value interface{}) {
		adv(name)
		adv(value)
	}
	ane := func(name string, e ...Event) {
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
		ane("F3", N(1))
		ane("F4", N(-1))
		ane("F5", N(1.1))
		ane("F6", N(NewBigFloat("1.1")))
		ane("F7", N(NewDFloat("1.1")))
		ane("F8", N(NewBDF("1.1")))
		ane("F9", NULL())
		ane("F10", N(NewBigInt("1000")))
		ane("F12", N(compact_float.QuietNaN()))
		ane("F13", N(compact_float.SignalingNaN()))
		ane("F14", UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
		ane("F15", T(compact_time.AsCompactTime(_this.Time)))
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
			N(1),
			N(-1),
			N(1.1),
			N(NewBigFloat("1.1")),
			N(NewDFloat("1.1")),
			N(NewBDF("1.1")),
			NULL(),
			N(NewBigInt("1000")),
			N(compact_float.QuietNaN()),
			N(compact_float.SignalingNaN()),
			UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
			T(compact_time.AsCompactTime(_this.Time)),
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
