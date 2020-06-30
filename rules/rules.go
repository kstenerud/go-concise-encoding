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
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/version"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type RuleOptions struct {
	// Concise encoding spec version to adhere to
	ConciseEncodingVersion uint64

	// Limits before the ruleset artificially terminates with an error.
	MaxBytesLength      uint64
	MaxStringLength     uint64
	MaxURILength        uint64
	MaxIDLength         uint64
	MaxMarkupNameLength uint64
	MaxContainerDepth   uint64
	MaxObjectCount      uint64
	MaxReferenceCount   uint64
	// Max bytes total for all array types
}

func NewDefaultRuleOptions() *RuleOptions {
	var options RuleOptions
	options = defaultRuleOptions
	return &options
}

// Rules constrains the order in which builder commands may be sent, such that
// they form a valid and complete Concise Encoding document.
type Rules struct {
	options           RuleOptions
	charValidator     UTF8Validator
	maxDepth          int
	stateStack        []ruleState
	arrayType         ruleEvent
	arrayData         []byte
	chunkByteCount    uint64
	chunkBytesWritten uint64
	arrayBytesWritten uint64
	isFinalChunk      bool
	objectCount       uint64
	unassignedIDs     []interface{}
	assignedIDs       map[interface{}]ruleEvent
	// TODO: Keep track of marked data stats (sizes etc)
	nextReceiver events.DataEventReceiver
}

func NewRules(options *RuleOptions, nextReceiver events.DataEventReceiver) *Rules {
	_this := new(Rules)
	_this.Init(options, nextReceiver)
	return _this
}

func (_this *Rules) Init(options *RuleOptions, nextReceiver events.DataEventReceiver) {
	_this.options = applyDefaultRuleOptions(options)
	_this.stateStack = make([]ruleState, 0, _this.options.MaxContainerDepth)
	_this.nextReceiver = nextReceiver

	_this.Reset()
}

func (_this *Rules) Reset() {
	_this.stateStack = _this.stateStack[:0]
	_this.stackState(stateAwaitingEndDocument)
	_this.stackState(stateAwaitingVersion)
	_this.unassignedIDs = _this.unassignedIDs[:0]
	_this.assignedIDs = make(map[interface{}]ruleEvent)

	_this.arrayType = eventTypeNothing
	_this.arrayData = _this.arrayData[:0]
	_this.chunkByteCount = 0
	_this.chunkBytesWritten = 0
	_this.arrayBytesWritten = 0
	_this.isFinalChunk = false
	_this.objectCount = 0
}

func (_this *Rules) OnVersion(version uint64) {
	_this.assertCurrentStateAllowsType(eventTypeVersion)
	if version != _this.options.ConciseEncodingVersion {
		panic(fmt.Errorf("Expected version %v but got version %v", _this.options.ConciseEncodingVersion, version))
	}
	_this.changeState(stateAwaitingTLO)
	_this.nextReceiver.OnVersion(version)
}

func (_this *Rules) OnPadding(count int) {
	_this.assertCurrentStateAllowsType(eventTypePadding)
	_this.nextReceiver.OnPadding(count)
}

func (_this *Rules) OnNil() {
	_this.addScalar(eventTypeNil)
	_this.nextReceiver.OnNil()
}

func (_this *Rules) OnBool(value bool) {
	_this.addScalar(eventTypeBool)
	_this.nextReceiver.OnBool(value)
}

func (_this *Rules) OnTrue() {
	_this.addScalar(eventTypeBool)
	_this.nextReceiver.OnTrue()
}

func (_this *Rules) OnFalse() {
	_this.addScalar(eventTypeBool)
	_this.nextReceiver.OnFalse()
}

func (_this *Rules) OnPositiveInt(value uint64) {
	_this.onPositiveInt(value)
	_this.nextReceiver.OnPositiveInt(value)
}

func (_this *Rules) OnNegativeInt(value uint64) {
	_this.onNegativeInt()
	_this.nextReceiver.OnNegativeInt(value)
}

func (_this *Rules) OnInt(value int64) {
	if value >= 0 {
		_this.onPositiveInt(uint64(value))
	} else {
		_this.onNegativeInt()
	}
	_this.nextReceiver.OnInt(value)
}

func (_this *Rules) OnBigInt(value *big.Int) {
	if value.IsInt64() {
		_this.OnInt(value.Int64())
		return
	}

	zero := &big.Int{}
	if value.Cmp(zero) < 0 {
		_this.onNegativeInt()
	} else {
		if _this.isAwaitingID() {
			panic(fmt.Errorf("ID values must not be larger than 64 bits"))
		}
		// If we're not waiting for an ID, then the argument to onPositiveInt
		// isn't used.
		unusedValue := uint64(0)
		_this.onPositiveInt(unusedValue)
	}
	_this.nextReceiver.OnBigInt(value)
}

func (_this *Rules) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}
	_this.addScalar(eventTypeFloat)
	_this.nextReceiver.OnFloat(value)
}

func (_this *Rules) OnBigFloat(value *big.Float) {
	_this.addScalar(eventTypeFloat)
	_this.nextReceiver.OnBigFloat(value)
}

func (_this *Rules) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			_this.OnNan(true)
			return
		}
		_this.OnNan(false)
		return
	}

	_this.addScalar(eventTypeFloat)
	_this.nextReceiver.OnDecimalFloat(value)
}

func (_this *Rules) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		_this.OnNan(false)
		return
	case apd.NaNSignaling:
		_this.OnNan(true)
		return
	}

	_this.addScalar(eventTypeFloat)
	_this.nextReceiver.OnBigDecimalFloat(value)
}

func (_this *Rules) OnComplex(value complex128) {
	_this.addScalar(eventTypeCustom)
	_this.nextReceiver.OnComplex(value)
}

func (_this *Rules) OnNan(signaling bool) {
	_this.addScalar(eventTypeNan)
	_this.nextReceiver.OnNan(signaling)
}

func (_this *Rules) OnUUID(value []byte) {
	_this.addScalar(eventTypeUUID)
	_this.nextReceiver.OnUUID(value)
}

func (_this *Rules) OnTime(value time.Time) {
	_this.addScalar(eventTypeTime)
	_this.nextReceiver.OnTime(value)
}

func (_this *Rules) OnCompactTime(value *compact_time.Time) {
	_this.addScalar(eventTypeTime)
	_this.nextReceiver.OnCompactTime(value)
}

func (_this *Rules) OnBytes(value []byte) {
	_this.onBytesBegin()
	_this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		_this.onArrayData([]byte(value))
	}
	_this.nextReceiver.OnBytes(value)
}

func (_this *Rules) OnString(value string) {
	_this.onStringBegin()
	_this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		_this.onArrayData([]byte(value))
	}
	_this.nextReceiver.OnString(value)
}

func (_this *Rules) OnURI(value string) {
	_this.onURIBegin()
	_this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		_this.onArrayData([]byte(value))
	}
	_this.nextReceiver.OnURI(value)
}

func (_this *Rules) OnCustom(value []byte) {
	_this.onCustomBegin()
	_this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		_this.onArrayData([]byte(value))
	}
	_this.nextReceiver.OnCustom(value)
}

func (_this *Rules) OnBytesBegin() {
	_this.onBytesBegin()
	_this.nextReceiver.OnBytesBegin()
}

func (_this *Rules) OnStringBegin() {
	_this.onStringBegin()
	_this.nextReceiver.OnStringBegin()
}

func (_this *Rules) OnURIBegin() {
	_this.onURIBegin()
	_this.nextReceiver.OnURIBegin()
}

func (_this *Rules) OnCustomBegin() {
	_this.onCustomBegin()
	_this.nextReceiver.OnCustomBegin()
}

func (_this *Rules) OnArrayChunk(length uint64, isFinalChunk bool) {
	_this.onArrayChunk(length, isFinalChunk)
	_this.nextReceiver.OnArrayChunk(length, isFinalChunk)
}

func (_this *Rules) OnArrayData(data []byte) {
	_this.onArrayData(data)
	_this.nextReceiver.OnArrayData(data)
}

func (_this *Rules) OnList() {
	_this.beginContainer(eventTypeList, stateAwaitingListItem)
	_this.nextReceiver.OnList()
}

func (_this *Rules) OnMap() {
	_this.beginContainer(eventTypeMap, stateAwaitingMapKey)
	_this.nextReceiver.OnMap()
}

func (_this *Rules) OnMarkup() {
	_this.beginContainer(eventTypeMarkup, stateAwaitingMarkupName)
	_this.nextReceiver.OnMarkup()
}

func (_this *Rules) OnMetadata() {
	_this.beginContainer(eventTypeMetadata, stateAwaitingMetadataKey)
	_this.nextReceiver.OnMetadata()
}

func (_this *Rules) OnComment() {
	_this.beginContainer(eventTypeComment, stateAwaitingCommentItem)
	_this.nextReceiver.OnComment()
}

func (_this *Rules) OnEnd() {
	_this.assertCurrentStateAllowsType(eventTypeEndContainer)

	switch _this.getCurrentStateID() {
	case stateIDAwaitingListItem:
		_this.unstackState()
		_this.onChildEnded(eventTypeList)
	case stateIDAwaitingMapKey:
		_this.unstackState()
		_this.onChildEnded(eventTypeMap)
	case stateIDAwaitingMarkupKey:
		_this.changeState(stateAwaitingMarkupContents)
	case stateIDAwaitingMarkupContents:
		_this.unstackState()
		_this.onChildEnded(eventTypeMarkup)
	case stateIDAwaitingMetadataKey:
		_this.changeState(stateAwaitingMetadataObject)
		_this.incrementObjectCount()
	case stateIDAwaitingCommentItem:
		_this.unstackState()
		_this.incrementObjectCount()
	default:
		panic(fmt.Errorf("BUG: EndContainer() in state %x (%v) failed to trigger", _this.getCurrentState(), _this.getCurrentState()))
	}
	_this.nextReceiver.OnEnd()
}

func (_this *Rules) OnMarker() {
	if uint64(len(_this.assignedIDs)) >= _this.options.MaxReferenceCount {
		panic(fmt.Errorf("Max number of marker IDs (%v) exceeded", _this.options.MaxReferenceCount))
	}
	_this.beginContainer(eventTypeMarker, stateAwaitingMarkerID)
	_this.nextReceiver.OnMarker()
}

func (_this *Rules) OnReference() {
	_this.beginContainer(eventTypeReference, stateAwaitingReferenceID)
	_this.nextReceiver.OnReference()
}

func (_this *Rules) OnEndDocument() {
	_this.assertCurrentStateAllowsType(eventTypeEndDocument)
	_this.nextReceiver.OnEndDocument()
}

func (_this *Rules) onPositiveInt(value uint64) {
	if _this.isAwaitingID() {
		_this.stackID(value)
	}
	_this.addScalar(eventTypePInt)
}

func (_this *Rules) onNegativeInt() {
	_this.addScalar(eventTypeNInt)
}

func (_this *Rules) onBytesBegin() {
	_this.beginArray(eventTypeBytes)
}

func (_this *Rules) onStringBegin() {
	_this.beginArray(eventTypeString)
}

func (_this *Rules) onURIBegin() {
	_this.beginArray(eventTypeURI)
}

func (_this *Rules) onCustomBegin() {
	_this.beginArray(eventTypeCustom)
}

func (_this *Rules) onArrayChunk(length uint64, isFinalChunk bool) {
	_this.assertCurrentStateAllowsType(eventTypeAChunk)

	_this.chunkByteCount = length
	_this.chunkBytesWritten = 0
	_this.isFinalChunk = isFinalChunk
	_this.changeState(stateAwaitingArrayData)

	if length == 0 {
		_this.onArrayChunkEnded()
	}
}

func (_this *Rules) onArrayData(data []byte) {
	_this.assertCurrentStateAllowsType(eventTypeAData)

	dataLength := uint64(len(data))
	if _this.chunkBytesWritten+dataLength > _this.chunkByteCount {
		panic(fmt.Errorf("Chunk length %v exceeded by %v bytes",
			_this.chunkByteCount, _this.chunkBytesWritten+dataLength-_this.chunkByteCount))
	}

	switch _this.arrayType {
	case eventTypeBytes:
		if _this.arrayBytesWritten+dataLength > _this.options.MaxBytesLength {
			panic(fmt.Errorf("Max byte array length (%v) exceeded", _this.options.MaxBytesLength))
		}
	case eventTypeString:
		if _this.arrayBytesWritten+dataLength > _this.options.MaxStringLength {
			panic(fmt.Errorf("Max string length (%v) exceeded", _this.options.MaxStringLength))
		}
		if _this.isStringInsideComment() {
			_this.validateCommentContents(data)
		} else {
			_this.validateStringContents(data)
		}
		if _this.isAwaitingID() {
			_this.arrayData = append(_this.arrayData, data...)
		}
	case eventTypeURI:
		if _this.arrayBytesWritten+dataLength > _this.options.MaxURILength {
			panic(fmt.Errorf("Max URI length (%v) exceeded", _this.options.MaxURILength))
		}
		if _this.isAwaitingID() {
			_this.arrayData = append(_this.arrayData, data...)
		}
		// Note: URI validation happens when the array is complete
	}

	_this.arrayBytesWritten += dataLength
	_this.chunkBytesWritten += dataLength
	if _this.chunkBytesWritten == _this.chunkByteCount {
		_this.onArrayChunkEnded()
	}
}

func (_this *Rules) getCurrentState() ruleState {
	return _this.stateStack[len(_this.stateStack)-1]
}

func (_this *Rules) getCurrentStateID() ruleState {
	return _this.getCurrentState() & ruleState(ruleIDFieldMask)
}

func (_this *Rules) getParentState() ruleState {
	return _this.stateStack[len(_this.stateStack)-2]
}

func (_this *Rules) hasParentState() bool {
	return len(_this.stateStack) > 1
}

func (_this *Rules) changeState(st ruleState) {
	_this.stateStack[len(_this.stateStack)-1] = st
}

func (_this *Rules) stackState(st ruleState) {
	if uint64(len(_this.stateStack)) >= _this.options.MaxContainerDepth {
		panic(fmt.Errorf("Max depth of %v exceeded", _this.options.MaxContainerDepth-rulesMaxDepthAdjust))
	}
	_this.stateStack = append(_this.stateStack, st)
}

func (_this *Rules) unstackState() {
	_this.stateStack = _this.stateStack[:len(_this.stateStack)-1]
}

func (_this *Rules) isAwaitingID() bool {
	if _this.getCurrentState()&ruleState(eventArrayChunk|eventArrayData) != 0 {
		return _this.getParentState()&ruleFlagAwaitingID != 0
	}
	return _this.getCurrentState()&ruleFlagAwaitingID != 0
}

func (_this *Rules) isAwaitingMarkupName() bool {
	return _this.getCurrentState() == stateAwaitingMarkupName
}

func (_this *Rules) stackID(id interface{}) {
	_this.unassignedIDs = append(_this.unassignedIDs, id)
}

func (_this *Rules) unstackID() (id interface{}) {
	id = _this.unassignedIDs[len(_this.unassignedIDs)-1]
	_this.unassignedIDs = _this.unassignedIDs[:len(_this.unassignedIDs)-1]
	return
}

func (_this *Rules) isStringInsideComment() bool {
	return _this.hasParentState() &&
		_this.getParentState()&ruleState(ruleIDFieldMask) == stateIDAwaitingCommentItem
}

func (_this *Rules) validateStringContents(data []byte) {
	for _, ch := range data {
		_this.charValidator.AddByte(int(ch))
	}
}

func (_this *Rules) validateCommentContents(data []byte) {
	for _, ch := range data {
		_this.charValidator.AddByte(int(ch))
		if _this.charValidator.IsCompleteCharacter() {
			validateRulesCommentCharacter(_this.charValidator.GetCharacter())
		}
	}
}

func (_this *Rules) getFirstRealContainer() ruleState {
	for i := len(_this.stateStack) - 1; i >= 0; i-- {
		currentState := _this.stateStack[i]
		if currentState&ruleFlagRealContainer != 0 {
			return currentState
		}
	}
	panic(fmt.Errorf("BUG: Could not find real container in state stack"))
}

func assertStateAllowsType(currentState ruleState, objectType ruleEvent) {
	allowedEventMask := ruleEvent(currentState) & ruleEventsMask
	if objectType&allowedEventMask == 0 {
		panic(fmt.Errorf("%v not allowed while awaiting %v", objectType, currentState))
	}
}

func (_this *Rules) assertCurrentStateAllowsType(objectType ruleEvent) {
	assertStateAllowsType(_this.getCurrentState(), objectType)
}

func (_this *Rules) beginArray(arrayType ruleEvent) {
	_this.assertCurrentStateAllowsType(arrayType)

	_this.arrayType = arrayType
	_this.arrayData = _this.arrayData[:0]
	_this.chunkByteCount = 0
	_this.chunkBytesWritten = 0
	_this.arrayBytesWritten = 0
	_this.isFinalChunk = false

	_this.stackState(stateAwaitingArrayChunk)
}

func (_this *Rules) onArrayChunkEnded() {
	if !_this.isFinalChunk {
		_this.changeState(stateAwaitingArrayChunk)
		return
	}

	_this.unstackState()

	switch _this.arrayType {
	case eventTypeString:
		if _this.isAwaitingMarkupName() {

			if _this.arrayBytesWritten == 0 {
				panic(fmt.Errorf("Markup name cannot be length 0"))
			}
			if _this.arrayBytesWritten > _this.options.MaxMarkupNameLength {
				panic(fmt.Errorf("Markup name length %v exceeds max of %v", _this.arrayBytesWritten, _this.options.MaxMarkupNameLength))
			}
		}
		if _this.isAwaitingID() {
			if _this.arrayBytesWritten == 0 {
				panic(fmt.Errorf("An ID cannot be length 0"))
			}
			if _this.arrayBytesWritten > _this.options.MaxIDLength {
				panic(fmt.Errorf("ID length %v exceeds max of %v", _this.arrayBytesWritten, _this.options.MaxIDLength))
			}
			_this.stackID(string(_this.arrayData))
		}
	case eventTypeURI:
		if _this.arrayBytesWritten < 2 {
			panic(fmt.Errorf("URI length must allow at least a scheme and colon (2 chars)"))
		}
		if _this.isAwaitingID() {
			url, err := url.Parse(string(_this.arrayData))
			if err != nil {
				panic(fmt.Errorf("%v", err))
			}
			_this.stackID(url)
		}
	case eventTypeBytes:
		// Nothing to do
	}

	arrayType := _this.arrayType
	_this.arrayType = eventTypeNothing
	_this.onChildEnded(arrayType)
}

func (_this *Rules) incrementObjectCount() {
	_this.objectCount++
	if _this.objectCount > _this.options.MaxObjectCount {
		panic(fmt.Errorf("Max object count of %v exceeded", _this.options.MaxObjectCount))
	}
}

func (_this *Rules) onChildEnded(childType ruleEvent) {
	_this.incrementObjectCount()

	switch _this.getCurrentStateID() {
	case stateIDAwaitingMetadataObject:
		container := _this.getFirstRealContainer()
		assertStateAllowsType(container, childType)
		_this.unstackState()
		_this.onChildEnded(childType)
	case stateIDAwaitingMarkerObject:
		container := _this.getFirstRealContainer()
		assertStateAllowsType(container, childType)
		markerID := _this.unstackID()
		if _, exists := _this.assignedIDs[markerID]; exists {
			panic(fmt.Errorf("%v: Marker ID already defined", markerID))
		}
		_this.assignedIDs[markerID] = childType
		_this.unstackState()
		_this.onChildEnded(childType)
	case stateIDAwaitingReferenceID:
		container := _this.getFirstRealContainer()
		markerID := _this.unstackID()

		_, ok := markerID.(*url.URL)
		if ok {
			// We have no way to verify what the URL points to, so call it "anything".
			_this.unstackState()
			_this.onChildEnded(eventTypeAny)
			return
		}

		referencedType, ok := _this.assignedIDs[markerID]
		if !ok {
			panic(fmt.Errorf("Referenced ID [%v] not found", markerID))
		}
		assertStateAllowsType(container, referencedType)
		_this.unstackState()
		_this.onChildEnded(referencedType)
	default:
		_this.changeState(childEndRuleStateChanges[_this.getCurrentStateID()])
	}
}

func (_this *Rules) addScalar(scalarType ruleEvent) {
	_this.assertCurrentStateAllowsType(scalarType)
	_this.onChildEnded(scalarType)
}

func (_this *Rules) beginContainer(containerType ruleEvent, newState ruleState) {
	_this.assertCurrentStateAllowsType(containerType)
	_this.stackState(newState)
}

var defaultRuleOptions = RuleOptions{
	ConciseEncodingVersion: version.ConciseEncodingVersion,
	MaxBytesLength:         1000000000,
	MaxStringLength:        100000000,
	MaxURILength:           10000,
	MaxIDLength:            100,
	MaxMarkupNameLength:    100,
	MaxContainerDepth:      1000,
	MaxObjectCount:         10000000,
	MaxReferenceCount:      100000,
	// TODO: References need to check for amplification attacks. Keep count of referenced things and their object counts
}

// The initial rule state comes pre-stacked. This value accounts for it in calculations.
const rulesMaxDepthAdjust = 2

func applyDefaultRuleOptions(original *RuleOptions) RuleOptions {
	var options RuleOptions
	if original == nil {
		options = defaultRuleOptions
	} else {
		options = *original
		if options.ConciseEncodingVersion < 1 {
			options.ConciseEncodingVersion = defaultRuleOptions.ConciseEncodingVersion
		}
		if options.MaxBytesLength < 1 {
			options.MaxBytesLength = defaultRuleOptions.MaxBytesLength
		}
		if options.MaxStringLength < 1 {
			options.MaxStringLength = defaultRuleOptions.MaxStringLength
		}
		if options.MaxURILength < 1 {
			options.MaxURILength = defaultRuleOptions.MaxURILength
		}
		if options.MaxIDLength < 1 {
			options.MaxIDLength = defaultRuleOptions.MaxIDLength
		}
		if options.MaxMarkupNameLength < 1 {
			options.MaxMarkupNameLength = defaultRuleOptions.MaxMarkupNameLength
		}
		if options.MaxContainerDepth < 1 {
			options.MaxContainerDepth = defaultRuleOptions.MaxContainerDepth
		}
		if options.MaxObjectCount < 1 {
			options.MaxObjectCount = defaultRuleOptions.MaxObjectCount
		}
		if options.MaxReferenceCount < 1 {
			options.MaxReferenceCount = defaultRuleOptions.MaxReferenceCount
		}
	}

	options.MaxContainerDepth += rulesMaxDepthAdjust

	return options
}

type ruleEvent int

const (
	eventIDNothing ruleEvent = iota
	eventIDVersion
	eventIDPadding
	eventIDNil
	eventIDBool
	eventIDPInt
	eventIDNInt
	eventIDFloat
	eventIDNan
	eventIDUUID
	eventIDTime
	eventIDList
	eventIDMap
	eventIDMarkup
	eventIDMetadata
	eventIDComment
	eventIDMarker
	eventIDReference
	eventIDEndContainer
	eventIDBytes
	eventIDString
	eventIDURI
	eventIDCustom
	eventIDAChunk
	eventIDAData
	eventIDEndDocument
)

var ruleEventNames = [...]string{
	"nothing",
	"version",
	"padding",
	"nil",
	"bool",
	"positive int",
	"negative int",
	"float",
	"nan",
	"UUID",
	"time",
	"list",
	"map",
	"markup",
	"metadata",
	"comment",
	"marker",
	"reference",
	"end container",
	"bytes",
	"string",
	"URI",
	"Custom",
	"array chunk",
	"array data",
	"end document",
}

func (_this ruleEvent) String() string {
	return ruleEventNames[_this&ruleEvent(ruleIDFieldMask)]
}

type ruleState int

const (
	stateIDAwaitingNothing ruleState = iota
	stateIDAwaitingVersion
	stateIDAwaitingTLO
	stateIDAwaitingListItem
	stateIDAwaitingCommentItem
	stateIDAwaitingMapKey
	stateIDAwaitingMapValue
	stateIDAwaitingMetadataKey
	stateIDAwaitingMetadataValue
	stateIDAwaitingMetadataObject
	stateIDAwaitingMarkupName
	stateIDAwaitingMarkupKey
	stateIDAwaitingMarkupValue
	stateIDAwaitingMarkupContents
	stateIDAwaitingMarkerID
	stateIDAwaitingMarkerObject
	stateIDAwaitingReferenceID
	stateIDAwaitingArrayChunk
	stateIDAwaitingArrayData
	stateIDAwaitingEndDocument
)

var ruleStateNames = [...]string{
	"nothing",
	"version",
	"top-level object",
	"list item",
	"comment contents",
	"map key",
	"map value",
	"metadata key",
	"metadata value",
	"metadata object",
	"markup name",
	"markup attribute key",
	"markup attribute value",
	"markup contents",
	"marker ID",
	"marker object",
	"reference id",
	"array chunk",
	"array data",
	"end document",
}

func (_this ruleState) String() string {
	return ruleStateNames[_this&ruleState(ruleIDFieldMask)]
}

const (
	ruleIDFieldEnd  ruleEvent = 1 << 5
	ruleIDFieldMask           = ruleIDFieldEnd - 1
)

const (
	eventVersion = ruleEvent(ruleIDFieldEnd) << iota
	eventPadding
	eventScalar
	eventPositiveInt
	eventNil
	eventNan
	eventBeginList
	eventBeginMap
	eventBeginMarkup
	eventBeginMetadata
	eventBeginComment
	eventBeginMarker
	eventBeginReference
	eventEndContainer
	eventBeginBytes
	eventBeginString
	eventBeginURI
	eventBeginCustom
	eventArrayChunk
	eventArrayData
	eventEndDocument
	ruleEventsEnd
	ruleEventsMask = (ruleEventsEnd - 1) - (ruleIDFieldEnd - 1)
)

const (
	ruleFlagRealContainer = ruleState(ruleEventsEnd) << iota
	ruleFlagAwaitingID
	ruleFlagsEnd
	ruleFlagsMask = (ruleFlagsEnd - 1) - (ruleState(ruleEventsEnd) - 1)
)

const (
	eventTypeNothing      = eventIDNothing
	eventTypeVersion      = eventIDVersion | eventVersion
	eventTypePadding      = eventIDPadding | eventPadding
	eventTypeNil          = eventIDNil | eventNil
	eventTypeBool         = eventIDBool | eventScalar
	eventTypePInt         = eventIDPInt | eventPositiveInt
	eventTypeNInt         = eventIDNInt | eventScalar
	eventTypeFloat        = eventIDFloat | eventScalar
	eventTypeNan          = eventIDNan | eventNan
	eventTypeUUID         = eventIDUUID | eventScalar
	eventTypeTime         = eventIDTime | eventScalar
	eventTypeList         = eventIDList | eventBeginList
	eventTypeMap          = eventIDMap | eventBeginMap
	eventTypeMarkup       = eventIDMarkup | eventBeginMarkup
	eventTypeMetadata     = eventIDMetadata | eventBeginMetadata
	eventTypeComment      = eventIDComment | eventBeginComment
	eventTypeMarker       = eventIDMarker | eventBeginMarker
	eventTypeReference    = eventIDReference | eventBeginReference
	eventTypeEndContainer = eventIDEndContainer | eventEndContainer
	eventTypeBytes        = eventIDBytes | eventBeginBytes
	eventTypeString       = eventIDString | eventBeginString
	eventTypeURI          = eventIDURI | eventBeginURI
	eventTypeCustom       = eventIDCustom | eventBeginCustom
	eventTypeAChunk       = eventIDAChunk | eventArrayChunk
	eventTypeAData        = eventIDAData | eventArrayData
	eventTypeEndDocument  = eventIDEndDocument | eventEndDocument
	eventTypeAny          = ruleEventsMask
)

// Primary rules
const (
	eventsArray         = eventBeginBytes | eventBeginString | eventBeginURI | eventBeginCustom
	eventsInvisible     = eventPadding | eventBeginComment | eventBeginMetadata
	eventsKeyableObject = eventsInvisible | eventScalar | eventPositiveInt | eventsArray | eventBeginMarker | eventBeginReference
	eventsAnyObject     = eventsKeyableObject | eventNil | eventNan | eventBeginList | eventBeginMap | eventBeginMarkup
	allowAny            = ruleState(eventsAnyObject)
	allowTLO            = allowAny | ruleState(eventEndDocument)
	allowListItem       = allowAny | ruleState(eventEndContainer)
	allowMapKey         = ruleState(eventsKeyableObject | eventEndContainer)
	allowMapValue       = allowAny
	allowCommentItem    = ruleState(eventBeginString | eventBeginComment | eventEndContainer | eventPadding)
	allowMarkupName     = ruleState(eventPositiveInt | eventBeginString | eventPadding)
	allowMarkupContents = ruleState(eventBeginString | eventBeginComment | eventBeginMarkup | eventEndContainer | eventPadding)
	allowMarkerID       = ruleState(eventPositiveInt | eventBeginString | eventPadding)
	allowReferenceID    = ruleState(eventPositiveInt | eventBeginString | eventBeginURI | eventPadding)
	allowArrayChunk     = ruleState(eventArrayChunk)
	allowArrayData      = ruleState(eventArrayData)
	allowVersion        = ruleState(eventVersion)
	allowEndDocument    = ruleState(eventEndDocument | eventBeginComment | eventPadding)

	stateAwaitingNothing        = stateIDAwaitingNothing
	stateAwaitingVersion        = stateIDAwaitingVersion | allowVersion
	stateAwaitingTLO            = stateIDAwaitingTLO | allowTLO | ruleFlagRealContainer
	stateAwaitingEndDocument    = stateIDAwaitingEndDocument | allowEndDocument
	stateAwaitingListItem       = stateIDAwaitingListItem | allowListItem | ruleFlagRealContainer
	stateAwaitingMapKey         = stateIDAwaitingMapKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMapValue       = stateIDAwaitingMapValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMarkupName     = stateIDAwaitingMarkupName | allowMarkupName | ruleFlagRealContainer
	stateAwaitingMarkupKey      = stateIDAwaitingMarkupKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMarkupValue    = stateIDAwaitingMarkupValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMarkupContents = stateIDAwaitingMarkupContents | allowMarkupContents | ruleFlagRealContainer
	stateAwaitingMarkerID       = stateIDAwaitingMarkerID | allowMarkerID | ruleFlagAwaitingID
	stateAwaitingMarkerObject   = stateIDAwaitingMarkerObject | allowAny
	stateAwaitingReferenceID    = stateIDAwaitingReferenceID | allowReferenceID | ruleFlagAwaitingID
	stateAwaitingCommentItem    = stateIDAwaitingCommentItem | allowCommentItem /* Not a "real" container */
	stateAwaitingMetadataKey    = stateIDAwaitingMetadataKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMetadataValue  = stateIDAwaitingMetadataValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMetadataObject = stateIDAwaitingMetadataObject | allowAny
	stateAwaitingArrayChunk     = stateIDAwaitingArrayChunk | allowArrayChunk
	stateAwaitingArrayData      = stateIDAwaitingArrayData | allowArrayData
)

var childEndRuleStateChanges = [...]ruleState{
	/* stateIDAwaitingNothing                */ stateAwaitingNothing,
	/* stateIDAwaitingVersion              > */ stateAwaitingTLO,
	/* stateIDAwaitingTLO                  > */ stateAwaitingEndDocument,
	/* stateIDAwaitingListItem               */ stateAwaitingListItem,
	/* stateIDAwaitingCommentItem            */ stateAwaitingCommentItem,
	/* stateIDAwaitingMapKey               > */ stateAwaitingMapValue,
	/* stateIDAwaitingMapValue             > */ stateAwaitingMapKey,
	/* stateIDAwaitingMetadataKey          > */ stateAwaitingMetadataValue,
	/* stateIDAwaitingMetadataValue        > */ stateAwaitingMetadataKey,
	/* stateIDAwaitingMetadataObject         */ stateIDAwaitingMetadataObject,
	/* stateIDAwaitingMarkupName           > */ stateAwaitingMarkupKey,
	/* stateIDAwaitingMarkupAttributeKey   > */ stateAwaitingMarkupValue,
	/* stateIDAwaitingMarkupAttributeValue > */ stateAwaitingMarkupKey,
	/* stateIDAwaitingMarkupContents         */ stateAwaitingMarkupContents,
	/* stateIDAwaitingMarkerID             > */ stateAwaitingMarkerObject,
	/* stateIDAwaitingMarkerObject           */ stateAwaitingMarkerObject,
	/* stateIDAwaitingReferenceID            */ stateAwaitingReferenceID,
	/* stateIDAwaitingArrayChunk             */ stateAwaitingArrayChunk,
	/* stateIDAwaitingArrayData              */ stateAwaitingArrayData,
	/* stateIDAwaitingEndDocument          > */ stateAwaitingNothing,
}

func validateRulesCommentCharacter(ch int) {
	switch ch {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08 /*, 0x09, 0x0a*/, 0x0b, 0x0c /*, 0x0d*/, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x7f,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0x2028, 0x2029:
		panic(fmt.Errorf("0x%04x: Invalid comment character", ch))
	}
}
