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
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type IndentType int

const (
	IndentTypeSpace = iota
	IndentTypeTab
)

type BracePosition int

const (
	BracePositionAdjacent = iota
	BracePositionNextLine
)

// TODO: Actually use these options
type CTEEncoderOptions struct {
	IndentCount   int
	IndentType    IndentType
	BracePosition BracePosition
}

type CTEEncoder struct {
	buff           buffer
	containerState []cteEncoderState
	currentState   cteEncoderState
	options        *CTEEncoderOptions
}

func NewCTEEncoder(options *CTEEncoderOptions) *CTEEncoder {
	if options == nil {
		options = &CTEEncoderOptions{}
	}
	return &CTEEncoder{
		options: options,
	}
}

func (this *CTEEncoder) Document() []byte {
	return this.buff.bytes
}

func (this *CTEEncoder) OnPadding(count int) {
	// Nothing to do
}

func (this *CTEEncoder) OnVersion(version uint64) {
	this.addFmt("c%d ", version)
}

func (this *CTEEncoder) OnNil() {
	this.addPrefix()
	this.addString("@nil")
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnBool(value bool) {
	if value {
		this.OnTrue()
	} else {
		this.OnFalse()
	}
}

func (this *CTEEncoder) OnTrue() {
	this.addPrefix()
	this.addString("@true")
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnFalse() {
	this.addPrefix()
	this.addString("@false")
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnInt(value int64) {
	if value >= 0 {
		this.OnPositiveInt(uint64(value))
	} else {
		this.OnNegativeInt(uint64(-value))
	}
}

func (this *CTEEncoder) OnBigInt(value *big.Int) {
	this.addFmt("%v", value)
}

func (this *CTEEncoder) OnPositiveInt(value uint64) {
	this.addPrefix()
	this.addFmt("%d", value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnNegativeInt(value uint64) {
	this.addPrefix()
	this.addFmt("-%d", value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnBinaryFloat(value float64) {
	if math.IsNaN(value) {
		this.OnNan(isSignalingNan(value))
		return
	}
	this.addPrefix()
	if math.IsInf(value, 0) {
		if value < 0 {
			this.addString("-@inf")
		} else {
			this.addString("@inf")
		}
		return
	}
	// TODO: Hex float?
	this.addFmt("%g", value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnDecimalFloat(value compact_float.DFloat) {
	this.addString(value.Text('g'))
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnBigDecimalFloat(value *apd.Decimal) {
	this.addString(value.Text('g'))
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnNan(signaling bool) {
	this.addPrefix()
	if signaling {
		this.addString("@snan")
	} else {
		this.addString("@nan")
	}
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("Expected UUID length 16 but got %v", len(v)))
	}
	this.addPrefix()
	this.addFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnComplex(value complex128) {
	this.addPrefix()
	panic("TODO: CTEEncoder.OnComplex")
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnTime(value time.Time) {
	this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (this *CTEEncoder) OnCompactTime(value *compact_time.Time) {
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
	this.addPrefix()
	switch value.TimeIs {
	case compact_time.TypeDate:
		this.addFmt("%d-%02d-%02d", value.Year, value.Month, value.Day)
	case compact_time.TypeTime:
		this.addFmt("%02d:%02d:%02d%v%v", value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	case compact_time.TypeTimestamp:
		this.addFmt("%d-%02d-%02d/%02d:%02d:%02d%v%v",
			value.Year, value.Month, value.Day, value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	default:
		panic(fmt.Errorf("Unknown compact time type %v", value.TimeIs))
	}
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnBytes(value []byte) {
	this.addPrefix()
	this.encodeHex('b', value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnURI(value string) {
	this.addPrefix()
	// TODO: URL escaping
	this.addFmt(`u"%v"`, value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnString(value string) {
	this.addPrefix()
	if this.currentState == cteEncoderStateAwaitMarkupItem ||
		this.currentState == cteEncoderStateAwaitMarkupFirstItemPre ||
		this.currentState == cteEncoderStateAwaitMarkupFirstItemPost ||
		isUnquotedString(value) {
		this.addString(value)
	} else {
		this.addFmt(`"%v"`, value)
	}
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnCustom(value []byte) {
	this.addPrefix()
	this.encodeHex('c', value)
	this.addSuffix()
	this.transitionState()
}

func (this *CTEEncoder) OnBytesBegin() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitBytes, `b"`)
}

func (this *CTEEncoder) OnStringBegin() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitQuotedString, `"`)
}

func (this *CTEEncoder) OnURIBegin() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitURI, `u"`)
}

func (this *CTEEncoder) OnCustomBegin() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitCustom, `c"`)
}

func (this *CTEEncoder) OnArrayChunk(length uint64, isFinalChunk bool) {
	panic("TODO: CTEEncoder.OnArrayChunk")
}

func (this *CTEEncoder) OnArrayData(data []byte) {
	panic("TODO: CTEEncoder.OnArrayData")
	dst := this.buff.Allocate(len(data))
	copy(dst, data)
}

func (this *CTEEncoder) OnList() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitListFirstItem, "[")
}

func (this *CTEEncoder) OnMap() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitMapFirstKey, "{")
}

func (this *CTEEncoder) OnMarkup() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitMarkupFirstItemPre, "")
	this.stackState(cteEncoderStateAwaitMarkupName, "<")
}

func (this *CTEEncoder) OnMetadata() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitMetaFirstKey, "(")
}

func (this *CTEEncoder) OnComment() {
	this.addPrefix()
	this.stackState(cteEncoderStateAwaitCommentItem, "/*")
}

func (this *CTEEncoder) OnEnd() {
	// TODO: Make this nicer
	isInvisible := this.currentState == cteEncoderStateAwaitMetaKey ||
		this.currentState == cteEncoderStateAwaitMetaFirstKey
	this.unstackState()
	if isInvisible {
		this.currentState |= cteEncoderStateWithInvisibleItem
	} else {
		this.addSuffix()
		this.transitionState()
	}
}

func (this *CTEEncoder) OnMarker() {
	this.addPrefix()
	panic("TODO: CTEEncoder.OnMarker")
	this.addSuffix()
}

func (this *CTEEncoder) OnReference() {
	this.addPrefix()
	panic("TODO: CTEEncoder.OnReference")
	this.addSuffix()
}

func (this *CTEEncoder) OnEndDocument() {
}

// ============================================================================

func (this *CTEEncoder) stackState(newState cteEncoderState, prefix string) {
	this.containerState = append(this.containerState, this.currentState)
	this.currentState = newState
	this.addString(prefix)
}

func (this *CTEEncoder) unstackState() {
	this.addString(cteEncoderTerminators[this.currentState])
	this.currentState = this.containerState[len(this.containerState)-1]
	this.containerState = this.containerState[:len(this.containerState)-1]
}

func (this *CTEEncoder) transitionState() {
	this.currentState = cteEncoderStateTransitions[this.currentState]
}

func (this *CTEEncoder) addPrefix() {
	cteEncoderPrefixHandlers[this.currentState](this)
}

func (this *CTEEncoder) addSuffix() {
	cteEncoderSuffixHandlers[this.currentState](this)
}

func (this *CTEEncoder) addString(str string) {
	dst := this.buff.Allocate(len(str))
	copy(dst, str)
}

func (this *CTEEncoder) addFmt(format string, args ...interface{}) {
	// TODO: Make something more efficient
	this.addString(fmt.Sprintf(format, args...))
}

var hexToChar = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
}

func (this *CTEEncoder) encodeHex(prefix byte, value []byte) {
	dst := this.buff.Allocate(len(value)*2 + 3)
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

func (this *CTEEncoder) suffixNone() {
}

func (this *CTEEncoder) suffixEquals() {
	this.addString("=")
}

func (this *CTEEncoder) prefixNone() {
}

func (this *CTEEncoder) prefixIndent() {
}

func (this *CTEEncoder) prefixSpacer() {
	this.addString(" ")
}

func (this *CTEEncoder) prefixPipe() {
	this.addString("|")
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
	cteEncoderPrefixHandlers[cteEncoderStateAwaitListItem] = (*CTEEncoder).prefixSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapFirstKey] = (*CTEEncoder).prefixIndent
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapKey] = (*CTEEncoder).prefixSpacer
	cteEncoderPrefixHandlers[cteEncoderStateAwaitMapValue] = (*CTEEncoder).prefixIndent
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
