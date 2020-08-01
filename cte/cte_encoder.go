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
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Receives data events, constructing a CTE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Encoder struct {
	buff               buffer.StreamingWriteBuffer
	chunkBuffer        []byte
	remainingChunkSize uint64
	containerState     []cteEncoderState
	containerItemCount []int
	currentState       cteEncoderState
	currentItemCount   int
	options            options.CTEEncoderOptions
	nextPrefix         string
	prefixGenerators   []cteEncoderPrefixGenerator
}

// Create a new CTE encoder, which will receive data events and write a document
// to writer. If options is nil, default options will be used.
func NewEncoder(writer io.Writer, options *options.CTEEncoderOptions) *Encoder {
	_this := &Encoder{}
	_this.Init(writer, options)
	return _this
}

// Initialize this encoder, which will receive data events and write a document
// to writer. If options is nil, default options will be used.
func (_this *Encoder) Init(writer io.Writer, options *options.CTEEncoderOptions) {
	_this.options = *options.WithDefaultsApplied()
	_this.buff.Init(writer, _this.options.BufferSize)
	_this.prefixGenerators = cteEncoderPrefixGenerators[:]
	if len(_this.options.Indent) > 0 {
		_this.prefixGenerators = cteEncoderPrettyPrefixHandlers[:]
	}
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnPadding(count int) {
	// Nothing to do
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.addFmt("c%d", version)
	_this.setState(cteEncoderStateAwaitTLO)
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
		_this.addFmt("%v#%v", _this.nextPrefix, value)
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
	// TODO: Hex float?
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
		panic(fmt.Errorf("Expected UUID length 16 but got %v", len(v)))
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
			panic(fmt.Errorf("Unknown compact time timezone type %v", value.TimezoneIs))
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
		panic(fmt.Errorf("Unknown compact time type %v", value.TimeIs))
	}
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnBytes(value []byte) {
	_this.addPrefix()
	_this.encodeHex('b', value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnURI(value []byte) {
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
		_this.addFmt(`%v#u"%v"`, _this.nextPrefix, asString)
		_this.currentItemCount++
		_this.transitionState()
	default:
		_this.addPrefix()
		_this.addFmt(`u"%v"`, asString)
		_this.currentItemCount++
		_this.transitionState()
	}
}

func (_this *Encoder) handleStringMarkerID(value []byte) {
	_this.addPrefix()
	_this.unstackState()
	_this.nextPrefix = fmt.Sprintf("%v&%v:", _this.nextPrefix, string(value))
}

func (_this *Encoder) handleStringReferenceID(value []byte) {
	_this.addPrefix()
	_this.unstackState()
	_this.addFmt("%v#%v", _this.nextPrefix, string(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) handleStringMarkupItem(value []byte) {
	_this.addPrefix()
	_this.addString(asMarkupContent(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) handleStringCommentItem(value []byte) {
	_this.addPrefix()
	_this.addString(string(value))
	_this.transitionState()
}

func (_this *Encoder) handleStringNormal(value []byte) {
	_this.addPrefix()
	_this.addString(asPotentialQuotedString(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnString(value []byte) {
	switch _this.currentState &^ cteEncoderStateWithInvisibleItem {
	case cteEncoderStateAwaitMarkerID:
		_this.handleStringMarkerID(value)
	case cteEncoderStateAwaitReferenceID:
		_this.handleStringReferenceID(value)
	case cteEncoderStateAwaitMarkupItem, cteEncoderStateAwaitMarkupFirstItem:
		_this.handleStringMarkupItem(value)
	case cteEncoderStateAwaitCommentItem:
		_this.handleStringCommentItem(value)
	default:
		_this.handleStringNormal(value)
	}
}

func (_this *Encoder) OnVerbatimString(value []byte) {
	_this.addPrefix()
	_this.addString(asVerbatimString(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnCustomBinary(value []byte) {
	_this.addPrefix()
	_this.encodeHex('c', value)
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnCustomText(value []byte) {
	_this.addPrefix()
	_this.addString(asCustomText(value))
	_this.currentItemCount++
	_this.transitionState()
}

func (_this *Encoder) OnBytesBegin() {
	_this.stackState(cteEncoderStateAwaitBytes, ``)
}

func (_this *Encoder) OnStringBegin() {
	_this.stackState(cteEncoderStateAwaitQuotedString, ``)
}

func (_this *Encoder) OnVerbatimStringBegin() {
	_this.stackState(cteEncoderStateAwaitVerbatimString, ``)
}

func (_this *Encoder) OnURIBegin() {
	_this.stackState(cteEncoderStateAwaitURI, ``)
}

func (_this *Encoder) OnCustomBinaryBegin() {
	_this.stackState(cteEncoderStateAwaitCustomBinary, ``)
}

func (_this *Encoder) OnCustomTextBegin() {
	_this.stackState(cteEncoderStateAwaitCustomText, ``)
}

func (_this *Encoder) finalizeArray() {
	oldState := _this.currentState
	_this.unstackState()
	switch oldState {
	case cteEncoderStateAwaitBytes:
		_this.OnBytes(_this.chunkBuffer)
	case cteEncoderStateAwaitQuotedString:
		_this.OnString(_this.chunkBuffer)
	case cteEncoderStateAwaitVerbatimString:
		_this.OnVerbatimString(_this.chunkBuffer)
	case cteEncoderStateAwaitURI:
		_this.OnURI(_this.chunkBuffer)
	case cteEncoderStateAwaitCustomBinary:
		_this.OnCustomBinary(_this.chunkBuffer)
	case cteEncoderStateAwaitCustomText:
		_this.OnCustomText(_this.chunkBuffer)
	}
	_this.chunkBuffer = _this.chunkBuffer[:0]
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
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitListFirstItem, "[")
}

func (_this *Encoder) OnMap() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMapFirstKey, "{")
}

func (_this *Encoder) OnMarkup() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMarkupFirstItem, "")
	_this.stackState(cteEncoderStateAwaitMarkupName, "<")
}

func (_this *Encoder) OnMetadata() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMetaFirstKey, "(")
}

func (_this *Encoder) OnComment() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitCommentItem, "/*")
}

func (_this *Encoder) OnEnd() {
	if _this.currentState == cteEncoderStateAwaitMarkupKey {
		_this.unstackState()
		return
	}

	oldState := _this.currentState &^ cteEncoderStateWithInvisibleItem
	if _this.currentItemCount > 0 && oldState != cteEncoderStateAwaitCommentItem {
		_this.applyIndentation(-1)
	}
	// TODO: Make this nicer
	isInvisible := oldState == cteEncoderStateAwaitMetaKey ||
		oldState == cteEncoderStateAwaitMetaFirstKey ||
		oldState == cteEncoderStateAwaitCommentItem
	_this.unstackState()
	if isInvisible {
		_this.currentState |= cteEncoderStateWithInvisibleItem
		_this.currentItemCount++

		if len(_this.options.Indent) > 0 {
			if _this.currentState&cteEncoderStateAwaitCommentItem != 0 {
				return
			}
			_this.nextPrefix = _this.generateIndentation(0)
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

// Internal

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
	_this.setState(cteEncoderStateTransitions[_this.currentState])
}

func (_this *Encoder) setState(newState cteEncoderState) {
	_this.currentState = newState
	_this.nextPrefix = _this.prefixGenerators[_this.currentState](_this)
}

func (_this *Encoder) addPrefix() {
	if len(_this.nextPrefix) > 0 {
		_this.addString(_this.nextPrefix)
		_this.nextPrefix = ""
	}
}

func (_this *Encoder) addString(str string) {
	// TODO: String continuation
	// TODO: Max column option?
	dst := _this.buff.Allocate(len(str))
	copy(dst, str)
}

func (_this *Encoder) addFmt(format string, args ...interface{}) {
	// TODO: Make something more efficient
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

func (_this *Encoder) applyIndentation(levelOffset int) {
	if len(_this.options.Indent) > 0 {
		level := len(_this.containerState) + levelOffset
		indent := strings.Repeat(_this.options.Indent, level)
		_this.addString("\n" + indent)
	}
}

func (_this *Encoder) generateIndentation(levelOffset int) string {
	level := len(_this.containerState) + levelOffset
	return "\n" + strings.Repeat(_this.options.Indent, level)
}

func (_this *Encoder) generateNoPrefix() string {
	return ""
}

func (_this *Encoder) generateSpacePrefix() string {
	return " "
}

func (_this *Encoder) generateIndentPrefix() string {
	return _this.generateIndentation(0)
}

func (_this *Encoder) generatePipePrefix() string {
	return "|"
}

func (_this *Encoder) generatePipeIndentPrefix() string {
	return "|" + _this.generateIndentation(0)
}

func (_this *Encoder) generateEqualsPrefix() string {
	return "="
}

func (_this *Encoder) generateSpaceEqualsPrefix() string {
	return " = "
}

type cteEncoderState int64

const (
	cteEncoderStateAwaitTLO cteEncoderState = iota * 2
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
	cteEncoderStateAwaitQuotedStringLast
	cteEncoderStateAwaitVerbatimString
	cteEncoderStateAwaitBytes
	cteEncoderStateAwaitBytesLast
	cteEncoderStateAwaitURI
	cteEncoderStateAwaitURILast
	cteEncoderStateAwaitCustomBinary
	cteEncoderStateAwaitCustomText
	cteEncoderStateAwaitCustomLast
	cteEncoderStateCount

	cteEncoderStateWithInvisibleItem cteEncoderState = 1
)

var cteEncoderStateNames = []string{
	"TLO",
	"TLO",
	"ListFirstItem",
	"ListFirstItem",
	"ListItem",
	"ListItem",
	"MapFirstKey",
	"MapFirstKey",
	"MapKey",
	"MapKey",
	"MapValue",
	"MapValue",
	"MetaFirstKey",
	"MetaFirstKey",
	"MetaKey",
	"MetaKey",
	"MetaValue",
	"MetaValue",
	"MarkupName",
	"MarkupName",
	"MarkupKey",
	"MarkupKey",
	"MarkupValue",
	"MarkupValue",
	"MarkupFirstItem",
	"MarkupFirstItem",
	"MarkupItem",
	"MarkupItem",
	"CommentItem",
	"CommentItem",
	"MarkerID",
	"MarkerID",
	"ReferenceID",
	"ReferenceID",
	"QuotedString",
	"QuotedString",
	"QuotedStringLast",
	"QuotedStringLast",
	"VerbatimString",
	"VerbatimString",
	"Bytes",
	"Bytes",
	"BytesLast",
	"BytesLast",
	"URI",
	"URI",
	"URILast",
	"URILast",
	"CustomBinary",
	"CustomBinary",
	"CustomText",
	"CustomText",
	"CustomLast",
	"CustomLast",
}

func (_this cteEncoderState) String() string {
	return cteEncoderStateNames[_this]
}

type cteEncoderPrefixGenerator func(*Encoder) string

var cteEncoderPrefixGenerators [cteEncoderStateCount]cteEncoderPrefixGenerator

func init() {
	for i := 0; i < int(cteEncoderStateCount); i++ {
		cteEncoderPrefixGenerators[i] = (*Encoder).generateNoPrefix
	}
	cteEncoderPrefixGenerators[cteEncoderStateAwaitTLO] = (*Encoder).generateSpacePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitListFirstItem] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitListItem] = (*Encoder).generateSpacePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMapFirstKey] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMapKey] = (*Encoder).generateSpacePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMapValue] = (*Encoder).generateEqualsPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMetaFirstKey] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMetaKey] = (*Encoder).generateSpacePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMetaValue] = (*Encoder).generateEqualsPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkupName] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkupKey] = (*Encoder).generateSpacePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkupValue] = (*Encoder).generateEqualsPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkupFirstItem] = (*Encoder).generatePipePrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkupItem] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitCommentItem] = (*Encoder).generateNoPrefix
	cteEncoderPrefixGenerators[cteEncoderStateAwaitMarkerID] = (*Encoder).generateNoPrefix
	for i := 0; i < int(cteEncoderStateCount); i += 2 {
		cteEncoderPrefixGenerators[i+1] = cteEncoderPrefixGenerators[i]
	}
}

var cteEncoderPrettyPrefixHandlers [cteEncoderStateCount]cteEncoderPrefixGenerator

func init() {
	for i := 0; i < int(cteEncoderStateCount); i++ {
		cteEncoderPrettyPrefixHandlers[i] = (*Encoder).generateNoPrefix
	}
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitTLO] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitListFirstItem] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitListItem] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMapFirstKey] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMapKey] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMapValue] = (*Encoder).generateSpaceEqualsPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMetaFirstKey] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMetaKey] = (*Encoder).generateIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMetaValue] = (*Encoder).generateSpaceEqualsPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkupName] = (*Encoder).generateNoPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkupKey] = (*Encoder).generateSpacePrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkupValue] = (*Encoder).generateEqualsPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkupFirstItem] = (*Encoder).generatePipeIndentPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkupItem] = (*Encoder).generateNoPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitCommentItem] = (*Encoder).generateNoPrefix
	cteEncoderPrettyPrefixHandlers[cteEncoderStateAwaitMarkerID] = (*Encoder).generateNoPrefix
	for i := 0; i < int(cteEncoderStateCount); i += 2 {
		cteEncoderPrettyPrefixHandlers[i+1] = cteEncoderPrettyPrefixHandlers[i]
	}
}

var cteEncoderStateTransitions [cteEncoderStateCount]cteEncoderState

func init() {
	cteEncoderStateTransitions[cteEncoderStateAwaitTLO] = cteEncoderStateAwaitTLO
	cteEncoderStateTransitions[cteEncoderStateAwaitListFirstItem] = cteEncoderStateAwaitListItem
	cteEncoderStateTransitions[cteEncoderStateAwaitListItem] = cteEncoderStateAwaitListItem
	cteEncoderStateTransitions[cteEncoderStateAwaitMapFirstKey] = cteEncoderStateAwaitMapValue
	cteEncoderStateTransitions[cteEncoderStateAwaitMapKey] = cteEncoderStateAwaitMapValue
	cteEncoderStateTransitions[cteEncoderStateAwaitMapValue] = cteEncoderStateAwaitMapKey
	cteEncoderStateTransitions[cteEncoderStateAwaitMetaFirstKey] = cteEncoderStateAwaitMetaValue
	cteEncoderStateTransitions[cteEncoderStateAwaitMetaKey] = cteEncoderStateAwaitMetaValue
	cteEncoderStateTransitions[cteEncoderStateAwaitMetaValue] = cteEncoderStateAwaitMetaKey
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupName] = cteEncoderStateAwaitMarkupKey
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupKey] = cteEncoderStateAwaitMarkupValue
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupValue] = cteEncoderStateAwaitMarkupKey
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupFirstItem] = cteEncoderStateAwaitMarkupItem
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupItem] = cteEncoderStateAwaitMarkupItem
	cteEncoderStateTransitions[cteEncoderStateAwaitCommentItem] = cteEncoderStateAwaitCommentItem

	// cteEncoderStateTransitions[cteEncoderStateAwaitQuotedString] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitQuotedStringLast] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitBytes] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitBytesLast] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitURI] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitURILast] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitCustom] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitCustomLast] = cteEncoderStateAwait

	for i := 0; i < int(cteEncoderStateCount); i += 2 {
		cteEncoderStateTransitions[i+1] = cteEncoderStateTransitions[i]
	}
	// for i := cteEncoderState(0); i < cteEncoderStateCount; i += 2 {
	// 	cteEncoderStateTransitions[i+1] = i
	// }
}

var cteEncoderTerminators [cteEncoderStateCount]string

func init() {
	// cteEncoderTerminators[cteEncoderStateAwaitTLO] = ""
	cteEncoderTerminators[cteEncoderStateAwaitListFirstItem] = "]"
	cteEncoderTerminators[cteEncoderStateAwaitListItem] = "]"
	cteEncoderTerminators[cteEncoderStateAwaitMapFirstKey] = "}"
	cteEncoderTerminators[cteEncoderStateAwaitMapKey] = "}"
	// cteEncoderTerminators[cteEncoderStateAwaitMapValue] = ""
	cteEncoderTerminators[cteEncoderStateAwaitMetaFirstKey] = ")"
	cteEncoderTerminators[cteEncoderStateAwaitMetaKey] = ")"
	// cteEncoderTerminators[cteEncoderStateAwaitMetaValue] = ""
	// cteEncoderTerminators[cteEncoderStateAwaitMarkupName] = ""
	cteEncoderTerminators[cteEncoderStateAwaitMarkupKey] = ""
	// cteEncoderTerminators[cteEncoderStateAwaitMarkupValue] = ""
	cteEncoderTerminators[cteEncoderStateAwaitMarkupFirstItem] = ">"
	cteEncoderTerminators[cteEncoderStateAwaitMarkupItem] = ">"
	cteEncoderTerminators[cteEncoderStateAwaitCommentItem] = "*/"
	// cteEncoderTerminators[cteEncoderStateAwaitMarkerID] = ""
	// cteEncoderTerminators[cteEncoderStateAwaitMarkerItem] = ""
	// cteEncoderTerminators[cteEncoderStateAwaitReferenceID] = ""
	// cteEncoderTerminators[cteEncoderStateAwaitQuotedString] = ""
	cteEncoderTerminators[cteEncoderStateAwaitQuotedStringLast] = `"`
	// cteEncoderTerminators[cteEncoderStateAwaitBytes] = ""
	cteEncoderTerminators[cteEncoderStateAwaitBytesLast] = `"`
	// cteEncoderTerminators[cteEncoderStateAwaitURI] = ""
	cteEncoderTerminators[cteEncoderStateAwaitURILast] = `"`
	// cteEncoderTerminators[cteEncoderStateAwaitCustom] = ""
	cteEncoderTerminators[cteEncoderStateAwaitCustomLast] = `"`

	for i := 0; i < int(cteEncoderStateCount); i += 2 {
		cteEncoderTerminators[i+1] = cteEncoderTerminators[i]
	}
}
