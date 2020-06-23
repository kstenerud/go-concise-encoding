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

package concise_encoding

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type BinaryFloatEncodeAs int

const (
	BinaryFloatEncodeAsDecimal = iota
	BinaryFloatEncodeAsBinary
)

type BracePosition int

const (
	BracePositionAdjacent = iota
	BracePositionNextLine
)

type CTEEncoderOptions struct {
	Indent string
	// TODO: BracePosition option
	BracePosition BracePosition
	// TODO: BinaryFloatEncoding option
	BinaryFloatEncoding BinaryFloatEncodeAs
}

type CTEEncoder struct {
	buff           buffer
	containerState []cteEncoderState
	currentState   cteEncoderState
	options        CTEEncoderOptions
}

func NewCTEEncoder(options *CTEEncoderOptions) *CTEEncoder {
	_this := &CTEEncoder{}
	_this.Init(options)
	return _this
}

func (_this *CTEEncoder) Init(options *CTEEncoderOptions) {
	if options != nil {
		_this.options = *options
	}
}

func (_this *CTEEncoder) Document() []byte {
	return _this.buff.bytes
}

func (_this *CTEEncoder) OnPadding(count int) {
	// Nothing to do
}

func (_this *CTEEncoder) OnVersion(version uint64) {
	_this.addFmt("c%d ", version)
}

func (_this *CTEEncoder) OnNil() {
	_this.addPrefix()
	_this.addString("@nil")
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *CTEEncoder) OnTrue() {
	_this.addPrefix()
	_this.addString("@true")
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnFalse() {
	_this.addPrefix()
	_this.addString("@false")
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *CTEEncoder) OnBigInt(value *big.Int) {
	_this.addFmt("%v", value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnPositiveInt(value uint64) {
	_this.addPrefix()
	_this.addFmt("%d", value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnNegativeInt(value uint64) {
	_this.addPrefix()
	_this.addFmt("-%d", value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(isSignalingNan(value))
		return
	}
	_this.addPrefix()
	if math.IsInf(value, 0) {
		if value < 0 {
			_this.addString("-@inf")
		} else {
			_this.addString("@inf")
		}
		return
	}
	// TODO: Hex float?
	_this.addFmt("%g", value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnBigFloat(value *big.Float) {
	_this.addPrefix()
	_this.addString(bigFloatToString(value))
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnDecimalFloat(value compact_float.DFloat) {
	_this.addPrefix()
	_this.addString(value.Text('g'))
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnBigDecimalFloat(value *apd.Decimal) {
	_this.addPrefix()
	_this.addString(value.Text('g'))
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnNan(signaling bool) {
	_this.addPrefix()
	if signaling {
		_this.addString("@snan")
	} else {
		_this.addString("@nan")
	}
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("Expected UUID length 16 but got %v", len(v)))
	}
	_this.addPrefix()
	_this.addFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnComplex(value complex128) {
	_this.addPrefix()
	panic("TODO: CTEEncoder.OnComplex")
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnTime(value time.Time) {
	_this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (_this *CTEEncoder) OnCompactTime(value *compact_time.Time) {
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
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnBytes(value []byte) {
	_this.addPrefix()
	_this.encodeHex('b', value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnURI(value string) {
	_this.addPrefix()
	// TODO: URL escaping
	_this.addFmt(`u"%v"`, value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnString(value string) {
	_this.addPrefix()
	if _this.currentState == cteEncoderStateAwaitMarkupItem ||
		_this.currentState == cteEncoderStateAwaitMarkupFirstItemPre ||
		_this.currentState == cteEncoderStateAwaitMarkupFirstItemPost ||
		isUnquotedString(value) {
		_this.addString(value)
	} else {
		_this.addFmt(`"%v"`, value)
	}
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnCustom(value []byte) {
	_this.addPrefix()
	_this.encodeHex('c', value)
	_this.addSuffix()
	_this.transitionState()
}

func (_this *CTEEncoder) OnBytesBegin() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitBytes, `b"`)
}

func (_this *CTEEncoder) OnStringBegin() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitQuotedString, `"`)
}

func (_this *CTEEncoder) OnURIBegin() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitURI, `u"`)
}

func (_this *CTEEncoder) OnCustomBegin() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitCustom, `c"`)
}

func (_this *CTEEncoder) OnArrayChunk(length uint64, isFinalChunk bool) {
	panic("TODO: CTEEncoder.OnArrayChunk")
}

func (_this *CTEEncoder) OnArrayData(data []byte) {
	panic("TODO: CTEEncoder.OnArrayData")
	dst := _this.buff.Allocate(len(data))
	copy(dst, data)
}

func (_this *CTEEncoder) OnList() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitListFirstItem, "[")
}

func (_this *CTEEncoder) OnMap() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMapFirstKey, "{")
}

func (_this *CTEEncoder) OnMarkup() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMarkupFirstItemPre, "")
	_this.stackState(cteEncoderStateAwaitMarkupName, "<")
}

func (_this *CTEEncoder) OnMetadata() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitMetaFirstKey, "(")
}

func (_this *CTEEncoder) OnComment() {
	_this.addPrefix()
	_this.stackState(cteEncoderStateAwaitCommentItem, "/*")
}

func (_this *CTEEncoder) OnEnd() {
	// TODO: Make this nicer
	isInvisible := _this.currentState == cteEncoderStateAwaitMetaKey ||
		_this.currentState == cteEncoderStateAwaitMetaFirstKey
	_this.unstackState()
	if isInvisible {
		_this.currentState |= cteEncoderStateWithInvisibleItem
	} else {
		_this.addSuffix()
		_this.transitionState()
	}
}

func (_this *CTEEncoder) OnMarker() {
	_this.addPrefix()
	panic("TODO: CTEEncoder.OnMarker")
	_this.addSuffix()
}

func (_this *CTEEncoder) OnReference() {
	_this.addPrefix()
	panic("TODO: CTEEncoder.OnReference")
	_this.addSuffix()
}

func (_this *CTEEncoder) OnEndDocument() {
}

// ============================================================================

func (_this *CTEEncoder) stackState(newState cteEncoderState, prefix string) {
	_this.containerState = append(_this.containerState, _this.currentState)
	_this.currentState = newState
	_this.addString(prefix)
}

func (_this *CTEEncoder) unstackState() {
	_this.addString(cteEncoderTerminators[_this.currentState])
	_this.currentState = _this.containerState[len(_this.containerState)-1]
	_this.containerState = _this.containerState[:len(_this.containerState)-1]
}

func (_this *CTEEncoder) transitionState() {
	_this.currentState = cteEncoderStateTransitions[_this.currentState]
}

func (_this *CTEEncoder) addPrefix() {
	cteEncoderPrefixHandlers[_this.currentState](_this)
}

func (_this *CTEEncoder) addSuffix() {
	cteEncoderSuffixHandlers[_this.currentState](_this)
}

func (_this *CTEEncoder) addString(str string) {
	dst := _this.buff.Allocate(len(str))
	copy(dst, str)
}

func (_this *CTEEncoder) addFmt(format string, args ...interface{}) {
	// TODO: Make something more efficient
	_this.addString(fmt.Sprintf(format, args...))
}

var hexToChar = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
}

func (_this *CTEEncoder) encodeHex(prefix byte, value []byte) {
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

func (_this *CTEEncoder) suffixNone() {
}

func (_this *CTEEncoder) suffixEquals() {
	_this.addString("=")
}

func (_this *CTEEncoder) prefixNone() {
}

func (_this *CTEEncoder) prefixIndent() {
	if len(_this.options.Indent) > 0 {
		level := len(_this.containerState)
		indent := strings.Repeat(_this.options.Indent, level)
		_this.addString("\n" + indent)
	}
}

func (_this *CTEEncoder) prefixSpacer() {
	_this.addString(" ")
}

func (_this *CTEEncoder) prefixIndentOrSpacer() {
	if len(_this.options.Indent) > 0 {
		_this.prefixIndent()
	} else {
		_this.addString(" ")
	}
}

func (_this *CTEEncoder) prefixPipe() {
	_this.addString("|")
}

type cteEncoderState int

const (
	/*  0 */ cteEncoderStateAwaitTLO cteEncoderState = iota * 2
	/*  2 */ cteEncoderStateAwaitListFirstItem
	/*  4 */ cteEncoderStateAwaitListItem
	/*  6 */ cteEncoderStateAwaitMapFirstKey
	/*  8 */ cteEncoderStateAwaitMapKey
	/* 10 */ cteEncoderStateAwaitMapValue
	/* 12 */ cteEncoderStateAwaitMetaFirstKey
	/* 14 */ cteEncoderStateAwaitMetaKey
	/* 16 */ cteEncoderStateAwaitMetaValue
	/* 18 */ cteEncoderStateAwaitMarkupName
	/* 20 */ cteEncoderStateAwaitMarkupKey
	/* 22 */ cteEncoderStateAwaitMarkupValue
	/* 24 */ cteEncoderStateAwaitMarkupFirstItemPre
	/* 26 */ cteEncoderStateAwaitMarkupFirstItemPost
	/* 28 */ cteEncoderStateAwaitMarkupItem
	cteEncoderStateAwaitCommentItem
	cteEncoderStateAwaitMarkerID
	cteEncoderStateAwaitMarkerItem
	cteEncoderStateAwaitReferenceID
	cteEncoderStateAwaitQuotedString
	cteEncoderStateAwaitQuotedStringLast
	cteEncoderStateAwaitBytes
	cteEncoderStateAwaitBytesLast
	cteEncoderStateAwaitURI
	cteEncoderStateAwaitURILast
	cteEncoderStateAwaitCustom
	cteEncoderStateAwaitCustomLast
	cteEncoderStateCount

	cteEncoderStateWithInvisibleItem cteEncoderState = 1
)

type cteEncoderPrefixFunction func(*CTEEncoder)

var cteEncoderPrefixHandlers [cteEncoderStateCount]cteEncoderPrefixFunction

func init() {
	for i := 0; i < int(cteEncoderStateCount); i++ {
		cteEncoderPrefixHandlers[i] = (*CTEEncoder).prefixNone
	}
	cteEncoderPrefixHandlers[cteEncoderStateAwaitTLO] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitListFirstItem] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitListItem] = (*CTEEncoder).prefixIndentOrSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapFirstKey] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapKey] = (*CTEEncoder).prefixIndentOrSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapValue] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMetaFirstKey] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMetaKey] = (*CTEEncoder).prefixSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMetaValue] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupName] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupKey] = (*CTEEncoder).prefixSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupValue] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupFirstItemPre] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupFirstItemPost] = (*CTEEncoder).prefixPipe
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkupItem] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitCommentItem] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkerID] = (*CTEEncoder).prefixNone
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMarkerItem] = (*CTEEncoder).prefixSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitReferenceID] = (*CTEEncoder).prefixNone
}

var cteEncoderSuffixHandlers [cteEncoderStateCount]cteEncoderPrefixFunction

func init() {
	for i := 0; i < int(cteEncoderStateCount); i++ {
		cteEncoderSuffixHandlers[i] = (*CTEEncoder).suffixNone
	}

	cteEncoderSuffixHandlers[cteEncoderStateAwaitMapFirstKey] = (*CTEEncoder).suffixEquals
	cteEncoderSuffixHandlers[cteEncoderStateAwaitMapKey] = (*CTEEncoder).suffixEquals
	cteEncoderSuffixHandlers[cteEncoderStateAwaitMetaFirstKey] = (*CTEEncoder).suffixEquals
	cteEncoderSuffixHandlers[cteEncoderStateAwaitMetaKey] = (*CTEEncoder).suffixEquals
	cteEncoderSuffixHandlers[cteEncoderStateAwaitMarkupKey] = (*CTEEncoder).suffixEquals

	for i := 0; i < int(cteEncoderStateCount); i += 2 {
		cteEncoderSuffixHandlers[i+1] = cteEncoderSuffixHandlers[i]
	}
}

var cteEncoderStateTransitions [cteEncoderStateCount]cteEncoderState

func init() {
	// cteEncoderStateTransitions[cteEncoderStateAwaitTLO] = cteEncoderStateAwait
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
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupFirstItemPre] = cteEncoderStateAwaitMarkupFirstItemPost
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupFirstItemPost] = cteEncoderStateAwaitMarkupItem
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkupItem] = cteEncoderStateAwaitMarkupItem
	cteEncoderStateTransitions[cteEncoderStateAwaitCommentItem] = cteEncoderStateAwaitCommentItem
	cteEncoderStateTransitions[cteEncoderStateAwaitMarkerID] = cteEncoderStateAwaitMarkerItem
	// cteEncoderStateTransitions[cteEncoderStateAwaitMarkerItem] = cteEncoderStateAwait
	// cteEncoderStateTransitions[cteEncoderStateAwaitReferenceID] = cteEncoderStateAwait
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
	cteEncoderTerminators[cteEncoderStateAwaitMarkupFirstItemPre] = ">"
	cteEncoderTerminators[cteEncoderStateAwaitMarkupFirstItemPost] = ">"
	cteEncoderTerminators[cteEncoderStateAwaitMarkupItem] = ">"
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

func isUnquotedString(str string) bool {
	bytes := []byte(str)

	if len(bytes) == 0 {
		return false
	}

	if !hasProperty(bytes[0], ctePropertyUnquotedStart) {
		return false
	}

	for i := 1; i < len(bytes); i++ {
		if !hasProperty(bytes[i], ctePropertyUnquotedMid) {
			return false
		}
	}

	return true
}
