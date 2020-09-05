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

package cte

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Receives data events, constructing a CTE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling Encoder's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Encoder struct {
	buff                   buffer.StreamingWriteBuffer
	chunkBuffer            []byte
	remainingChunkSize     uint64
	containerState         []cteEncoderState
	containerItemCount     []int
	currentState           cteEncoderState
	currentItemCount       int
	nextPrefix             string
	prefixGenerators       []cteEncoderPrefixGenerator
	opts                   options.CTEEncoderOptions
	skipFirstMap           bool
	skipFirstList          bool
	containerDepth         int
	transitionFromSuffixes []string
}

// Create a new CTE encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func NewEncoder(opts *options.CTEEncoderOptions) *Encoder {
	_this := &Encoder{}
	_this.Init(opts)
	return _this
}

// Initialize this encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func (_this *Encoder) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.buff.Init(_this.opts.BufferSize)
	_this.prefixGenerators = cteEncoderPrefixGenerators[:]
	_this.transitionFromSuffixes = cteTransitionFromSuffixes[:]
	if len(_this.opts.Indent) > 0 {
		_this.prefixGenerators = cteEncoderPrettyPrefixHandlers[:]
		_this.transitionFromSuffixes = cteTransitionFromSuffixesPretty[:]
	}
	_this.skipFirstList = _this.opts.ImpliedStructure == options.ImpliedStructureList
	_this.skipFirstMap = _this.opts.ImpliedStructure == options.ImpliedStructureMap
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.buff.SetWriter(writer)
}

func (_this *Encoder) Reset() {
	_this.buff.Reset()
	_this.chunkBuffer = _this.chunkBuffer[:0]
	_this.remainingChunkSize = 0
	_this.containerState = _this.containerState[:0]
	_this.containerItemCount = _this.containerItemCount[:0]
	_this.currentState = 0
	_this.currentItemCount = 0
	_this.nextPrefix = ""
	_this.skipFirstList = _this.opts.ImpliedStructure == options.ImpliedStructureList
	_this.skipFirstMap = _this.opts.ImpliedStructure == options.ImpliedStructureMap
	_this.containerDepth = 0
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnBeginDocument() {
	switch _this.opts.ImpliedStructure {
	case options.ImpliedStructureList:
		_this.stackState(cteEncoderStateAwaitListFirstItem, "")
	case options.ImpliedStructureMap:
		_this.stackState(cteEncoderStateAwaitMapFirstKey, "")
	default:
		_this.setState(cteEncoderStateAwaitVersion)
	}
}

func (_this *Encoder) OnPadding(_ int) {
	// Nothing to do
}

func (_this *Encoder) OnVersion(version uint64) {
	if _this.opts.ImpliedStructure != options.ImpliedStructureNone {
		return
	}
	_this.addFmt("c%d", version)
	_this.transitionState()
}

func (_this *Encoder) OnNil() {
	_this.addPrefix()
	_this.addString("@nil")
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *Encoder) OnTrue() {
	_this.addPrefix()
	_this.addString("@true")
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnFalse() {
	_this.addPrefix()
	_this.addString("@false")
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	_this.addPrefix()
	_this.addFmt("%v", value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnPositiveInt(value uint64) {
	switch _this.currentState {
	case cteEncoderStateAwaitMarkerID:
		_this.unstackState()
		_this.nextPrefix = fmt.Sprintf("%v&%v:", _this.nextPrefix, value)
	case cteEncoderStateAwaitReferenceID:
		_this.unstackState()
		_this.addFmt("%v$%v", _this.nextPrefix, value)
		_this.currentItemCount++
		_this.transitionState()
	default:
		_this.addPrefix()
		_this.addFmt("%d", value)
		_this.currentItemCount++
		_this.transitionState()
	}
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	_this.addPrefix()
	_this.addFmt("-%d", value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}
	_this.addPrefix()
	if math.IsInf(value, 0) {
		if value < 0 {
			_this.addString("-@inf")
		} else {
			_this.addString("@inf")
		}
		_this.currentItemCount++
		_this.transitionState()
		return
	}
	_this.addFmt("%g", value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnBigFloat(value *big.Float) {
	_this.addPrefix()
	_this.addString(conversions.BigFloatToString(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.addPrefix()
	_this.addString(value.Text('g'))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.addPrefix()
	_this.addString(value.Text('g'))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.addPrefix()
	if signaling {
		_this.addString("@snan")
	} else {
		_this.addString("@nan")
	}
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("expected UUID length 16 but got %v", len(v)))
	}
	_this.addPrefix()
	_this.addFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnTime(value time.Time) {
	_this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (_this *Encoder) OnCompactTime(value *compact_time.Time) {
	tz := func(v *compact_time.Time) string {
		switch v.TimezoneIs {
		case compact_time.TypeUTC:
			return ""
		case compact_time.TypeAreaLocation:
			return fmt.Sprintf("/%v", v.AreaLocation)
		case compact_time.TypeLatitudeLongitude:
			return fmt.Sprintf("/%.2f/%.2f", float64(v.LatitudeHundredths)/100, float64(v.LongitudeHundredths)/100)
		default:
			panic(fmt.Errorf("unknown compact time timezone type %v", value.TimezoneIs))
		}
	}
	subsec := func(v *compact_time.Time) string {
		if v.Nanosecond == 0 {
			return ""
		}

		str := fmt.Sprintf("%.9f", float64(v.Nanosecond)/float64(1000000000))
		for str[len(str)-1] == '0' {
			str = str[:len(str)-1]
		}
		return str[1:]
	}
	_this.addPrefix()
	switch value.TimeIs {
	case compact_time.TypeDate:
		_this.addFmt("%d-%02d-%02d", value.Year, value.Month, value.Day)
	case compact_time.TypeTime:
		_this.addFmt("%02d:%02d:%02d%v%v", value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	case compact_time.TypeTimestamp:
		_this.addFmt("%d-%02d-%02d/%02d:%02d:%02d%v%v",
			value.Year, value.Month, value.Day, value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	default:
		panic(fmt.Errorf("unknown compact time type %v", value.TimeIs))
	}
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.addPrefix()
	switch arrayType {
	case events.ArrayTypeString:
		_this.handleString(value)
		return
	case events.ArrayTypeVerbatimString:
		_this.handleVerbatimString(value)
		return
	case events.ArrayTypeURI:
		_this.handleURI(value)
		return
	case events.ArrayTypeCustomBinary:
		_this.handleCustomBinary(value)
		return
	case events.ArrayTypeCustomText:
		_this.handleCustomText(value)
		return
	case events.ArrayTypeUint8:
		_this.encodeUint8Array(value)
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", arrayType))
	}
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.stackState(cteEncoderStateAwaitQuotedString, ``)
	case events.ArrayTypeVerbatimString:
		_this.stackState(cteEncoderStateAwaitVerbatimString, ``)
	case events.ArrayTypeURI:
		_this.stackState(cteEncoderStateAwaitURI, ``)
	case events.ArrayTypeCustomBinary:
		_this.stackState(cteEncoderStateAwaitCustomBinary, ``)
	case events.ArrayTypeCustomText:
		_this.stackState(cteEncoderStateAwaitCustomText, ``)
	case events.ArrayTypeUint8:
		_this.stackState(cteEncoderStateAwaitArrayU8, ``)
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", arrayType))
	}
}

func (_this *Encoder) OnArrayChunk(length uint64, moreChunksFollow bool) {
	if moreChunksFollow {
		return
	}

	if length == 0 {
		_this.finalizeArray()
		return
	}

	_this.remainingChunkSize = length
}

func (_this *Encoder) OnArrayData(data []byte) {
	// TODO: In future, don't buffer up the entire array; write it out as it arrives
	_this.chunkBuffer = append(_this.chunkBuffer, data...)
	_this.remainingChunkSize -= uint64(len(data))
	if _this.remainingChunkSize == 0 {
		_this.finalizeArray()
	}
}

func (_this *Encoder) OnList() {
	if _this.skipFirstList {
		_this.skipFirstList = false
		return
	}

	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitListFirstItem, "[")
	_this.containerDepth++
}

func (_this *Encoder) OnMap() {
	if _this.skipFirstMap {
		_this.skipFirstMap = false
		return
	}

	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMapFirstKey, "{")
	_this.containerDepth++
}

func (_this *Encoder) OnMarkup() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMarkupFirstItem, "")
	_this.stackState(cteEncoderStateAwaitMarkupName, "<")
	_this.containerDepth += 2
}

func (_this *Encoder) OnMetadata() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMetaFirstKey, "(")
	_this.containerDepth++
}

func (_this *Encoder) OnComment() {
	if _this.nextPrefix == "" &&
		_this.currentState != cteEncoderStateAwaitMapValue &&
		_this.currentState != cteEncoderStateAwaitMetaValue &&
		_this.currentState != cteEncoderStateAwaitMarkupValue {
		_this.nextPrefix = " "
	}
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitCommentItem, "/*")
	_this.containerDepth++
}

func (_this *Encoder) OnEnd() {
	if _this.containerDepth <= 0 {
		return
	}
	_this.containerDepth--

	if _this.currentState == cteEncoderStateAwaitMarkupKey {
		_this.unstackState()
		return
	}

	if _this.currentState == cteEncoderStateAwaitMarkupFirstItem {
		// When markup has no contents, we don't need a semicolon. This code
		// removes the automatic semicolon suffix after the attrs in that case.
		_this.buff.EraseLastByte()
	}

	oldState := _this.currentState
	if _this.currentItemCount > 0 && oldState != cteEncoderStateAwaitCommentItem {
		_this.applyIndentation(-1)
	}
	// TODO: Make this nicer
	isInvisible := oldState == cteEncoderStateAwaitMetaKey ||
		oldState == cteEncoderStateAwaitMetaFirstKey ||
		oldState == cteEncoderStateAwaitCommentItem
	_this.unstackState()
	if isInvisible {
		_this.currentItemCount++

		if len(_this.opts.Indent) > 0 {
			if _this.currentState&cteEncoderStateAwaitCommentItem != 0 {
				return
			}
			_this.nextPrefix = _this.generateIndentPrefix()
		} else {
			_this.nextPrefix = ""
		}
	} else {
		_this.currentItemCount++
		_this.transitionState()
	}
}

func (_this *Encoder) OnMarker() {
	_this.stackState(cteEncoderStateAwaitMarkerID, "")
}

func (_this *Encoder) OnReference() {
	_this.stackState(cteEncoderStateAwaitReferenceID, "")
}

func (_this *Encoder) OnEndDocument() {
	if _this.currentItemCount == 0 {
		_this.addString(" ")
	}
	_this.buff.Flush()
}

// ============================================================================

// Array Handlers

func (_this *Encoder) handleURI(value []byte) {
	asString := string(value)
	for _, ch := range value {
		if ch == '"' {
			asString = strings.ReplaceAll(asString, "\"", "%22")
			break
		}
	}

	switch _this.currentState {
	case cteEncoderStateAwaitReferenceID:
		_this.unstackState()
		_this.addFmt(`%v$u"%v"`, _this.nextPrefix, asString)
		_this.currentItemCount++
		_this.transitionState()
	default:
		_this.addPrefix()
		_this.addFmt(`u"%v"`, asString)
		_this.currentItemCount++
		_this.transitionState()
	}
}

func (_this *Encoder) handleStringNormal(value []byte) {
	_this.addPrefix()
	_this.addString(asPotentialQuotedString(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) handleString(value []byte) {
	switch _this.currentState {
	case cteEncoderStateAwaitMarkerID:
		_this.addPrefix()
		_this.unstackState()
		_this.nextPrefix = fmt.Sprintf("%v&%v:", _this.nextPrefix, string(value))
	case cteEncoderStateAwaitReferenceID:
		_this.addPrefix()
		_this.unstackState()
		_this.addFmt("%v$%v", _this.nextPrefix, string(value))
		_this.currentItemCount++
		_this.transitionState()
	case cteEncoderStateAwaitMarkupItem, cteEncoderStateAwaitMarkupFirstItem:
		_this.addPrefix()
		_this.addString(asMarkupContent(value))
		_this.currentItemCount++
		_this.transitionState()
	case cteEncoderStateAwaitCommentItem:
		_this.addPrefix()
		_this.addString(string(value))
		_this.transitionState()
	default:
		_this.addPrefix()
		_this.addString(asPotentialQuotedString(value))
		_this.currentItemCount++
		_this.transitionState()
	}
}

func (_this *Encoder) handleVerbatimString(value []byte) {
	_this.addPrefix()
	_this.addString(asVerbatimString(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) handleCustomBinary(value []byte) {
	_this.addPrefix()
	_this.encodeHex('b', value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) handleCustomText(value []byte) {
	_this.addPrefix()
	_this.addString(asCustomText(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) finalizeArray() {
	oldState := _this.currentState
	_this.unstackState()
	switch oldState {
	case cteEncoderStateAwaitArrayU8:
		_this.OnArray(events.ArrayTypeUint8, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	case cteEncoderStateAwaitQuotedString:
		_this.OnArray(events.ArrayTypeString, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	case cteEncoderStateAwaitVerbatimString:
		_this.OnArray(events.ArrayTypeVerbatimString, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	case cteEncoderStateAwaitURI:
		_this.OnArray(events.ArrayTypeURI, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	case cteEncoderStateAwaitCustomBinary:
		_this.OnArray(events.ArrayTypeCustomBinary, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	case cteEncoderStateAwaitCustomText:
		_this.OnArray(events.ArrayTypeCustomText, uint64(len(_this.chunkBuffer)), _this.chunkBuffer)
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", oldState))
	}
	_this.chunkBuffer = _this.chunkBuffer[:0]
}

// ============================================================================

// State

func (_this *Encoder) stackState(newState cteEncoderState, prefix string) {
	_this.addString(prefix)
	_this.containerState = append(_this.containerState, _this.currentState)
	_this.containerItemCount = append(_this.containerItemCount, 0)
	_this.currentItemCount = 0
	_this.setState(newState)
}

func (_this *Encoder) unstackState() {
	_this.addString(cteEncoderTerminators[_this.currentState])
	prevState := _this.containerState[len(_this.containerState)-1]
	_this.containerState = _this.containerState[:len(_this.containerState)-1]
	_this.currentItemCount = _this.containerItemCount[len(_this.containerItemCount)-1]
	_this.containerItemCount = _this.containerItemCount[:len(_this.containerItemCount)-1]
	_this.setState(prevState)
}

func (_this *Encoder) transitionState() {
	_this.addString(_this.transitionFromSuffixes[_this.currentState])
	_this.setState(cteEncoderStateTransitions[_this.currentState])
}

func (_this *Encoder) setState(newState cteEncoderState) {
	_this.currentState = newState
	_this.nextPrefix = _this.prefixGenerators[_this.currentState](_this)
}

// ============================================================================

// Encoding

func (_this *Encoder) addPrefix() {
	if len(_this.nextPrefix) > 0 {
		_this.addString(_this.nextPrefix)
		_this.nextPrefix = ""
	}
}

func (_this *Encoder) addString(str string) {
	// TODO: String continuation
	dst := _this.buff.Allocate(len(str))
	copy(dst, str)
}

func (_this *Encoder) addFmt(format string, args ...interface{}) {
	_this.addString(fmt.Sprintf(format, args...))
}

var hexToChar = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
}

func (_this *Encoder) encodeHex(prefix byte, value []byte) {
	dst := _this.buff.Allocate(len(value)*2 + 3)
	dst[0] = prefix
	dst[1] = '"'
	dst[len(dst)-1] = '"'
	dst = dst[2 : len(dst)-1]
	for i := 0; i < len(value); i++ {
		b := value[i]
		dst[i*2] = hexToChar[b>>4]
		dst[i*2+1] = hexToChar[b&15]
	}
}

func (_this *Encoder) encodeUint8Array(value []uint8) {
	header := []byte("|u8x")
	dst := _this.buff.Allocate(len(value)*3 + len(header) + 1)
	copy(dst, header)
	dst = dst[len(header):]
	base := 0
	for _, b := range value {
		dst[base] = ' '
		dst[base+1] = hexToChar[b>>4]
		dst[base+2] = hexToChar[b&15]
		base += 3
	}
	dst[base] = '|'
}

func (_this *Encoder) applyIndentation(levelOffset int) {
	if len(_this.opts.Indent) > 0 {
		level := len(_this.containerState) + levelOffset
		indent := strings.Repeat(_this.opts.Indent, level)
		_this.addString("\n" + indent)
	}
}

func (_this *Encoder) generateNoPrefix() string {
	return ""
}

func (_this *Encoder) generateSpacePrefix() string {
	return " "
}

func (_this *Encoder) generateIndentPrefix() string {
	level := len(_this.containerState)
	return "\n" + strings.Repeat(_this.opts.Indent, level)
}

// ============================================================================

// Data

type cteEncoderState int64

const (
	cteEncoderStateAwaitVersion cteEncoderState = iota
	cteEncoderStateAwaitTLO
	cteEncoderStateAwaitListFirstItem
	cteEncoderStateAwaitListItem
	cteEncoderStateAwaitMapFirstKey
	cteEncoderStateAwaitMapKey
	cteEncoderStateAwaitMapValue
	cteEncoderStateAwaitMetaFirstKey
	cteEncoderStateAwaitMetaKey
	cteEncoderStateAwaitMetaValue
	cteEncoderStateAwaitMarkupName
	cteEncoderStateAwaitMarkupKey
	cteEncoderStateAwaitMarkupValue
	cteEncoderStateAwaitMarkupFirstItem
	cteEncoderStateAwaitMarkupItem
	cteEncoderStateAwaitCommentItem
	cteEncoderStateAwaitMarkerID
	cteEncoderStateAwaitReferenceID
	cteEncoderStateAwaitQuotedString
	cteEncoderStateAwaitVerbatimString
	cteEncoderStateAwaitURI
	cteEncoderStateAwaitCustomBinary
	cteEncoderStateAwaitCustomText
	cteEncoderStateAwaitArrayBool
	cteEncoderStateAwaitArrayU8
	cteEncoderStateAwaitArrayU16
	cteEncoderStateAwaitArrayU32
	cteEncoderStateAwaitArrayU64
	cteEncoderStateAwaitArrayI8
	cteEncoderStateAwaitArrayI16
	cteEncoderStateAwaitArrayI32
	cteEncoderStateAwaitArrayI64
	cteEncoderStateAwaitArrayF16
	cteEncoderStateAwaitArrayF32
	cteEncoderStateAwaitArrayF64
	cteEncoderStateAwaitArrayUUID
	cteEncoderStateCount
)

var cteEncoderStateNames = [cteEncoderStateCount]string{
	cteEncoderStateAwaitVersion:         "Version",
	cteEncoderStateAwaitTLO:             "TLO",
	cteEncoderStateAwaitListFirstItem:   "ListFirstItem",
	cteEncoderStateAwaitListItem:        "ListItem",
	cteEncoderStateAwaitMapFirstKey:     "MapFirstKey",
	cteEncoderStateAwaitMapKey:          "MapKey",
	cteEncoderStateAwaitMapValue:        "MapValue",
	cteEncoderStateAwaitMetaFirstKey:    "MetaFirstKey",
	cteEncoderStateAwaitMetaKey:         "MetaKey",
	cteEncoderStateAwaitMetaValue:       "MetaValue",
	cteEncoderStateAwaitMarkupName:      "MarkupName",
	cteEncoderStateAwaitMarkupKey:       "MarkupKey",
	cteEncoderStateAwaitMarkupValue:     "MarkupValue",
	cteEncoderStateAwaitMarkupFirstItem: "MarkupFirstItem",
	cteEncoderStateAwaitMarkupItem:      "MarkupItem",
	cteEncoderStateAwaitCommentItem:     "CommentItem",
	cteEncoderStateAwaitMarkerID:        "MarkerID",
	cteEncoderStateAwaitReferenceID:     "ReferenceID",
	cteEncoderStateAwaitQuotedString:    "QuotedString",
	cteEncoderStateAwaitVerbatimString:  "VerbatimString",
	cteEncoderStateAwaitURI:             "URI",
	cteEncoderStateAwaitCustomBinary:    "CustomBinary",
	cteEncoderStateAwaitCustomText:      "CustomText",
	cteEncoderStateAwaitArrayBool:       "ArrayBool",
	cteEncoderStateAwaitArrayU8:         "ArrayU8",
	cteEncoderStateAwaitArrayU16:        "ArrayU16",
	cteEncoderStateAwaitArrayU32:        "ArrayU32",
	cteEncoderStateAwaitArrayU64:        "ArrayU64",
	cteEncoderStateAwaitArrayI8:         "ArrayI8",
	cteEncoderStateAwaitArrayI16:        "ArrayI16",
	cteEncoderStateAwaitArrayI32:        "ArrayI32",
	cteEncoderStateAwaitArrayI64:        "ArrayI64",
	cteEncoderStateAwaitArrayF16:        "ArrayF16",
	cteEncoderStateAwaitArrayF32:        "ArrayF32",
	cteEncoderStateAwaitArrayF64:        "ArrayF64",
	cteEncoderStateAwaitArrayUUID:       "ArrayUUID",
}

func (_this cteEncoderState) String() string {
	return cteEncoderStateNames[_this]
}

type cteEncoderPrefixGenerator func(*Encoder) string

var cteEncoderPrefixGenerators = [cteEncoderStateCount]cteEncoderPrefixGenerator{
	cteEncoderStateAwaitVersion:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitTLO:             (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitListFirstItem:   (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitListItem:        (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitMapFirstKey:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMapKey:          (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitMapValue:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMetaFirstKey:    (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMetaKey:         (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitMetaValue:       (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupName:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupKey:       (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitMarkupValue:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupFirstItem: (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupItem:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCommentItem:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkerID:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitReferenceID:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitQuotedString:    (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitVerbatimString:  (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitURI:             (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCustomBinary:    (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCustomText:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayBool:       (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU8:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI8:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayUUID:       (*Encoder).generateNoPrefix,
}

var cteEncoderPrettyPrefixHandlers = [cteEncoderStateCount]cteEncoderPrefixGenerator{
	cteEncoderStateAwaitVersion:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitTLO:             (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitListFirstItem:   (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitListItem:        (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMapFirstKey:     (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMapKey:          (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMapValue:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMetaFirstKey:    (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMetaKey:         (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMetaValue:       (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupName:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupKey:       (*Encoder).generateSpacePrefix,
	cteEncoderStateAwaitMarkupValue:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkupFirstItem: (*Encoder).generateIndentPrefix,
	cteEncoderStateAwaitMarkupItem:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCommentItem:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitMarkerID:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitReferenceID:     (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitQuotedString:    (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitVerbatimString:  (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitURI:             (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCustomBinary:    (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitCustomText:      (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayBool:       (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU8:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayU64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI8:         (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayI64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF16:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF32:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayF64:        (*Encoder).generateNoPrefix,
	cteEncoderStateAwaitArrayUUID:       (*Encoder).generateNoPrefix,
}

var cteEncoderStateTransitions = [cteEncoderStateCount]cteEncoderState{
	cteEncoderStateAwaitVersion:         cteEncoderStateAwaitTLO,
	cteEncoderStateAwaitTLO:             cteEncoderStateAwaitTLO,
	cteEncoderStateAwaitListFirstItem:   cteEncoderStateAwaitListItem,
	cteEncoderStateAwaitListItem:        cteEncoderStateAwaitListItem,
	cteEncoderStateAwaitMapFirstKey:     cteEncoderStateAwaitMapValue,
	cteEncoderStateAwaitMapKey:          cteEncoderStateAwaitMapValue,
	cteEncoderStateAwaitMapValue:        cteEncoderStateAwaitMapKey,
	cteEncoderStateAwaitMetaFirstKey:    cteEncoderStateAwaitMetaValue,
	cteEncoderStateAwaitMetaKey:         cteEncoderStateAwaitMetaValue,
	cteEncoderStateAwaitMetaValue:       cteEncoderStateAwaitMetaKey,
	cteEncoderStateAwaitMarkupName:      cteEncoderStateAwaitMarkupKey,
	cteEncoderStateAwaitMarkupKey:       cteEncoderStateAwaitMarkupValue,
	cteEncoderStateAwaitMarkupValue:     cteEncoderStateAwaitMarkupKey,
	cteEncoderStateAwaitMarkupFirstItem: cteEncoderStateAwaitMarkupItem,
	cteEncoderStateAwaitMarkupItem:      cteEncoderStateAwaitMarkupItem,
	cteEncoderStateAwaitCommentItem:     cteEncoderStateAwaitCommentItem,
}

var cteEncoderTerminators = [cteEncoderStateCount]string{
	cteEncoderStateAwaitListFirstItem:   "]",
	cteEncoderStateAwaitListItem:        "]",
	cteEncoderStateAwaitMapFirstKey:     "}",
	cteEncoderStateAwaitMapKey:          "}",
	cteEncoderStateAwaitMetaFirstKey:    ")",
	cteEncoderStateAwaitMetaKey:         ")",
	cteEncoderStateAwaitMarkupKey:       ";",
	cteEncoderStateAwaitMarkupFirstItem: ">",
	cteEncoderStateAwaitMarkupItem:      ">",
	cteEncoderStateAwaitCommentItem:     "*/",
}

var cteTransitionFromSuffixes = [cteEncoderStateCount]string{
	cteEncoderStateAwaitMapFirstKey:  "=",
	cteEncoderStateAwaitMapKey:       "=",
	cteEncoderStateAwaitMetaFirstKey: "=",
	cteEncoderStateAwaitMetaKey:      "=",
	cteEncoderStateAwaitMarkupKey:    "=",
}

var cteTransitionFromSuffixesPretty = [cteEncoderStateCount]string{
	cteEncoderStateAwaitMapFirstKey:  " = ",
	cteEncoderStateAwaitMapKey:       " = ",
	cteEncoderStateAwaitMetaFirstKey: " = ",
	cteEncoderStateAwaitMetaKey:      " = ",
	cteEncoderStateAwaitMarkupKey:    "=", // No spaces for markup k-v pairs
}
