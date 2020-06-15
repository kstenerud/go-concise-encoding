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
	"net/url"
	"time"

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
}

func DefaultRuleOptions() *RuleOptions {
	var options RuleOptions
	options = defaultRuleOptions
	return &options
}

type Rules struct {
	options           RuleOptions
	charValidator     Utf8Validator
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
	nextReceiver      DataEventReceiver
}

func NewRules(options *RuleOptions, nextReceiver DataEventReceiver) *Rules {
	this := new(Rules)
	this.Init(options, nextReceiver)
	return this
}

func (this *Rules) Init(options *RuleOptions, nextReceiver DataEventReceiver) {
	this.options = applyDefaultRuleOptions(options)
	this.stateStack = make([]ruleState, 0, this.options.MaxContainerDepth)
	this.nextReceiver = nextReceiver

	this.Reset()
}

func (this *Rules) Reset() {
	this.stateStack = this.stateStack[:0]
	this.stackState(stateAwaitingEndDocument)
	this.stackState(stateAwaitingVersion)
	this.unassignedIDs = this.unassignedIDs[:0]
	this.assignedIDs = make(map[interface{}]ruleEvent)

	this.arrayType = eventTypeNothing
	this.arrayData = this.arrayData[:0]
	this.chunkByteCount = 0
	this.chunkBytesWritten = 0
	this.arrayBytesWritten = 0
	this.isFinalChunk = false
	this.objectCount = 0
}

func (this *Rules) OnVersion(version uint64) {
	this.assertCurrentStateAllowsType(eventTypeVersion)
	if version != this.options.ConciseEncodingVersion {
		panic(fmt.Errorf("Expected version %v but got version %v", this.options.ConciseEncodingVersion, version))
	}
	this.changeState(stateAwaitingTLO)
	this.nextReceiver.OnVersion(version)
}

func (this *Rules) OnPadding(count int) {
	this.assertCurrentStateAllowsType(eventTypePadding)
	this.nextReceiver.OnPadding(count)
}

func (this *Rules) OnNil() {
	this.addScalar(eventTypeNil)
	this.nextReceiver.OnNil()
}

func (this *Rules) OnBool(value bool) {
	this.addScalar(eventTypeBool)
	this.nextReceiver.OnBool(value)
}

func (this *Rules) OnTrue() {
	this.addScalar(eventTypeBool)
	this.nextReceiver.OnTrue()
}

func (this *Rules) OnFalse() {
	this.addScalar(eventTypeBool)
	this.nextReceiver.OnFalse()
}

func (this *Rules) OnPositiveInt(value uint64) {
	this.onPositiveInt(value)
	this.nextReceiver.OnPositiveInt(value)
}

func (this *Rules) OnNegativeInt(value uint64) {
	this.onNegativeInt()
	this.nextReceiver.OnNegativeInt(value)
}

func (this *Rules) OnInt(value int64) {
	if value >= 0 {
		this.onPositiveInt(uint64(value))
	} else {
		this.onNegativeInt()
	}
	this.nextReceiver.OnInt(value)
}

func (this *Rules) OnBigInt(value *big.Int) {
	if value.IsInt64() {
		this.OnInt(value.Int64())
		return
	}

	zero := &big.Int{}
	if value.Cmp(zero) < 0 {
		this.onNegativeInt()
	} else {
		if this.isAwaitingID() {
			panic(fmt.Errorf("ID values must not be larger than 64 bits"))
		}
		// If we're not waiting for an ID, then the argument to onPositiveInt
		// isn't used.
		unusedValue := uint64(0)
		this.onPositiveInt(unusedValue)
	}
	this.nextReceiver.OnBigInt(value)
}

func (this *Rules) OnFloat(value float64) {
	if math.IsNaN(value) {
		this.OnNan(isSignalingNan(value))
		return
	}
	this.addScalar(eventTypeFloat)
	this.nextReceiver.OnFloat(value)
}

func (this *Rules) OnBigFloat(value *big.Float) {
	this.addScalar(eventTypeFloat)
	this.nextReceiver.OnBigFloat(value)
}

func (this *Rules) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			this.OnNan(true)
			return
		}
		this.OnNan(false)
		return
	}

	this.addScalar(eventTypeFloat)
	this.nextReceiver.OnDecimalFloat(value)
}

func (this *Rules) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		this.OnNan(false)
		return
	case apd.NaNSignaling:
		this.OnNan(true)
		return
	}

	this.addScalar(eventTypeFloat)
	this.nextReceiver.OnBigDecimalFloat(value)
}

func (this *Rules) OnComplex(value complex128) {
	this.addScalar(eventTypeCustom)
	this.nextReceiver.OnComplex(value)
}

func (this *Rules) OnNan(signaling bool) {
	this.addScalar(eventTypeNan)
	this.nextReceiver.OnNan(signaling)
}

func (this *Rules) OnUUID(value []byte) {
	this.addScalar(eventTypeUUID)
	this.nextReceiver.OnUUID(value)
}

func (this *Rules) OnTime(value time.Time) {
	this.addScalar(eventTypeTime)
	this.nextReceiver.OnTime(value)
}

func (this *Rules) OnCompactTime(value *compact_time.Time) {
	this.addScalar(eventTypeTime)
	this.nextReceiver.OnCompactTime(value)
}

func (this *Rules) OnBytes(value []byte) {
	this.onBytesBegin()
	this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		this.onArrayData([]byte(value))
	}
	this.nextReceiver.OnBytes(value)
}

func (this *Rules) OnString(value string) {
	this.onStringBegin()
	this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		this.onArrayData([]byte(value))
	}
	this.nextReceiver.OnString(value)
}

func (this *Rules) OnURI(value string) {
	this.onURIBegin()
	this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		this.onArrayData([]byte(value))
	}
	this.nextReceiver.OnURI(value)
}

func (this *Rules) OnCustom(value []byte) {
	this.onCustomBegin()
	this.onArrayChunk(uint64(len(value)), true)
	if len(value) > 0 {
		this.onArrayData([]byte(value))
	}
	this.nextReceiver.OnCustom(value)
}

func (this *Rules) OnBytesBegin() {
	this.onBytesBegin()
	this.nextReceiver.OnBytesBegin()
}

func (this *Rules) OnStringBegin() {
	this.onStringBegin()
	this.nextReceiver.OnStringBegin()
}

func (this *Rules) OnURIBegin() {
	this.onURIBegin()
	this.nextReceiver.OnURIBegin()
}

func (this *Rules) OnCustomBegin() {
	this.onCustomBegin()
	this.nextReceiver.OnCustomBegin()
}

func (this *Rules) OnArrayChunk(length uint64, isFinalChunk bool) {
	this.onArrayChunk(length, isFinalChunk)
	this.nextReceiver.OnArrayChunk(length, isFinalChunk)
}

func (this *Rules) OnArrayData(data []byte) {
	this.onArrayData(data)
	this.nextReceiver.OnArrayData(data)
}

func (this *Rules) OnList() {
	this.beginContainer(eventTypeList, stateAwaitingListItem)
	this.nextReceiver.OnList()
}

func (this *Rules) OnMap() {
	this.beginContainer(eventTypeMap, stateAwaitingMapKey)
	this.nextReceiver.OnMap()
}

func (this *Rules) OnMarkup() {
	this.beginContainer(eventTypeMarkup, stateAwaitingMarkupName)
	this.nextReceiver.OnMarkup()
}

func (this *Rules) OnMetadata() {
	this.beginContainer(eventTypeMetadata, stateAwaitingMetadataKey)
	this.nextReceiver.OnMetadata()
}

func (this *Rules) OnComment() {
	this.beginContainer(eventTypeComment, stateAwaitingCommentItem)
	this.nextReceiver.OnComment()
}

func (this *Rules) OnEnd() {
	this.assertCurrentStateAllowsType(eventTypeEndContainer)

	switch this.getCurrentStateId() {
	case stateIdAwaitingListItem:
		this.unstackState()
		this.onChildEnded(eventTypeList)
	case stateIdAwaitingMapKey:
		this.unstackState()
		this.onChildEnded(eventTypeMap)
	case stateIdAwaitingMarkupKey:
		this.changeState(stateAwaitingMarkupContents)
	case stateIdAwaitingMarkupContents:
		this.unstackState()
		this.onChildEnded(eventTypeMarkup)
	case stateIdAwaitingMetadataKey:
		this.changeState(stateAwaitingMetadataObject)
		this.incrementObjectCount()
	case stateIdAwaitingCommentItem:
		this.unstackState()
		this.incrementObjectCount()
	default:
		panic(fmt.Errorf("BUG: EndContainer() in state %x (%v) failed to trigger", this.getCurrentState(), this.getCurrentState()))
	}
	this.nextReceiver.OnEnd()
}

func (this *Rules) OnMarker() {
	if uint64(len(this.assignedIDs)) >= this.options.MaxReferenceCount {
		panic(fmt.Errorf("Max number of marker IDs (%v) exceeded", this.options.MaxReferenceCount))
	}
	this.beginContainer(eventTypeMarker, stateAwaitingMarkerID)
	this.nextReceiver.OnMarker()
}

func (this *Rules) OnReference() {
	this.beginContainer(eventTypeReference, stateAwaitingReferenceID)
	this.nextReceiver.OnReference()
}

func (this *Rules) OnEndDocument() {
	this.assertCurrentStateAllowsType(eventTypeEndDocument)
	this.nextReceiver.OnEndDocument()
}

func (this *Rules) onPositiveInt(value uint64) {
	if this.isAwaitingID() {
		this.stackId(value)
	}
	this.addScalar(eventTypePInt)
}

func (this *Rules) onNegativeInt() {
	this.addScalar(eventTypeNInt)
}

func (this *Rules) onBytesBegin() {
	this.beginArray(eventTypeBytes)
}

func (this *Rules) onStringBegin() {
	this.beginArray(eventTypeString)
}

func (this *Rules) onURIBegin() {
	this.beginArray(eventTypeURI)
}

func (this *Rules) onCustomBegin() {
	this.beginArray(eventTypeCustom)
}

func (this *Rules) onArrayChunk(length uint64, isFinalChunk bool) {
	this.assertCurrentStateAllowsType(eventTypeAChunk)

	this.chunkByteCount = length
	this.chunkBytesWritten = 0
	this.isFinalChunk = isFinalChunk
	this.changeState(stateAwaitingArrayData)

	if length == 0 {
		this.onArrayChunkEnded()
	}
}

func (this *Rules) onArrayData(data []byte) {
	this.assertCurrentStateAllowsType(eventTypeAData)

	dataLength := uint64(len(data))
	if this.chunkBytesWritten+dataLength > this.chunkByteCount {
		panic(fmt.Errorf("Chunk length %v exceeded by %v bytes",
			this.chunkByteCount, this.chunkBytesWritten+dataLength-this.chunkByteCount))
	}

	switch this.arrayType {
	case eventTypeBytes:
		if this.arrayBytesWritten+dataLength > this.options.MaxBytesLength {
			panic(fmt.Errorf("Max byte array length (%v) exceeded", this.options.MaxBytesLength))
		}
	case eventTypeString:
		if this.arrayBytesWritten+dataLength > this.options.MaxStringLength {
			panic(fmt.Errorf("Max string length (%v) exceeded", this.options.MaxStringLength))
		}
		if this.isStringInsideComment() {
			this.validateCommentContents(data)
		} else {
			this.validateStringContents(data)
		}
		if this.isAwaitingID() {
			this.arrayData = append(this.arrayData, data...)
		}
	case eventTypeURI:
		if this.arrayBytesWritten+dataLength > this.options.MaxURILength {
			panic(fmt.Errorf("Max URI length (%v) exceeded", this.options.MaxURILength))
		}
		if this.isAwaitingID() {
			this.arrayData = append(this.arrayData, data...)
		}
		// Note: URI validation happens when the array is complete
	}

	this.arrayBytesWritten += dataLength
	this.chunkBytesWritten += dataLength
	if this.chunkBytesWritten == this.chunkByteCount {
		this.onArrayChunkEnded()
	}
}

func (this *Rules) getCurrentState() ruleState {
	return this.stateStack[len(this.stateStack)-1]
}

func (this *Rules) getCurrentStateId() ruleState {
	return this.getCurrentState() & ruleState(ruleIDFieldMask)
}

func (this *Rules) getParentState() ruleState {
	return this.stateStack[len(this.stateStack)-2]
}

func (this *Rules) hasParentState() bool {
	return len(this.stateStack) > 1
}

func (this *Rules) changeState(st ruleState) {
	this.stateStack[len(this.stateStack)-1] = st
}

func (this *Rules) stackState(st ruleState) {
	if uint64(len(this.stateStack)) >= this.options.MaxContainerDepth {
		panic(fmt.Errorf("Max depth of %v exceeded", this.options.MaxContainerDepth-rulesMaxDepthAdjust))
	}
	this.stateStack = append(this.stateStack, st)
}

func (this *Rules) unstackState() {
	this.stateStack = this.stateStack[:len(this.stateStack)-1]
}

func (this *Rules) isAwaitingID() bool {
	if this.getCurrentState()&ruleState(eventArrayChunk|eventArrayData) != 0 {
		return this.getParentState()&ruleFlagAwaitingID != 0
	}
	return this.getCurrentState()&ruleFlagAwaitingID != 0
}

func (this *Rules) isAwaitingMarkupName() bool {
	return this.getCurrentState() == stateAwaitingMarkupName
}

func (this *Rules) stackId(id interface{}) {
	this.unassignedIDs = append(this.unassignedIDs, id)
}

func (this *Rules) unstackId() (id interface{}) {
	id = this.unassignedIDs[len(this.unassignedIDs)-1]
	this.unassignedIDs = this.unassignedIDs[:len(this.unassignedIDs)-1]
	return
}

func (this *Rules) isStringInsideComment() bool {
	return this.hasParentState() &&
		this.getParentState()&ruleState(ruleIDFieldMask) == stateIdAwaitingCommentItem
}

func (this *Rules) validateStringContents(data []byte) {
	for _, ch := range data {
		this.charValidator.AddByte(int(ch))
	}
}

func (this *Rules) validateCommentContents(data []byte) {
	for _, ch := range data {
		this.charValidator.AddByte(int(ch))
		if this.charValidator.IsCompleteCharacter() {
			validateRulesCommentCharacter(this.charValidator.Character())
		}
	}
}

func (this *Rules) getFirstRealContainer() ruleState {
	for i := len(this.stateStack) - 1; i >= 0; i-- {
		currentState := this.stateStack[i]
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

func (this *Rules) assertCurrentStateAllowsType(objectType ruleEvent) {
	assertStateAllowsType(this.getCurrentState(), objectType)
}

func (this *Rules) beginArray(arrayType ruleEvent) {
	this.assertCurrentStateAllowsType(arrayType)

	this.arrayType = arrayType
	this.arrayData = this.arrayData[:0]
	this.chunkByteCount = 0
	this.chunkBytesWritten = 0
	this.arrayBytesWritten = 0
	this.isFinalChunk = false

	this.stackState(stateAwaitingArrayChunk)
}

func (this *Rules) onArrayChunkEnded() {
	if !this.isFinalChunk {
		this.changeState(stateAwaitingArrayChunk)
		return
	}

	this.unstackState()

	switch this.arrayType {
	case eventTypeString:
		if this.isAwaitingMarkupName() {

			if this.arrayBytesWritten == 0 {
				panic(fmt.Errorf("Markup name cannot be length 0"))
			}
			if this.arrayBytesWritten > this.options.MaxMarkupNameLength {
				panic(fmt.Errorf("Markup name length %v exceeds max of %v", this.arrayBytesWritten, this.options.MaxMarkupNameLength))
			}
		}
		if this.isAwaitingID() {
			if this.arrayBytesWritten == 0 {
				panic(fmt.Errorf("An ID cannot be length 0"))
			}
			if this.arrayBytesWritten > this.options.MaxIDLength {
				panic(fmt.Errorf("ID length %v exceeds max of %v", this.arrayBytesWritten, this.options.MaxIDLength))
			}
			this.stackId(string(this.arrayData))
		}
	case eventTypeURI:
		if this.arrayBytesWritten < 2 {
			panic(fmt.Errorf("URI length must allow at least a scheme and colon (2 chars)"))
		}
		if this.isAwaitingID() {
			url, err := url.Parse(string(this.arrayData))
			if err != nil {
				panic(fmt.Errorf("%v", err))
			}
			this.stackId(url)
		}
	case eventTypeBytes:
		// Nothing to do
	}

	arrayType := this.arrayType
	this.arrayType = eventTypeNothing
	this.onChildEnded(arrayType)
}

func (this *Rules) incrementObjectCount() {
	this.objectCount++
	if this.objectCount > this.options.MaxObjectCount {
		panic(fmt.Errorf("Max object count of %v exceeded", this.options.MaxObjectCount))
	}
}

func (this *Rules) onChildEnded(childType ruleEvent) {
	this.incrementObjectCount()

	switch this.getCurrentStateId() {
	case stateIdAwaitingMetadataObject:
		container := this.getFirstRealContainer()
		assertStateAllowsType(container, childType)
		this.unstackState()
		this.onChildEnded(childType)
	case stateIdAwaitingMarkerObject:
		container := this.getFirstRealContainer()
		assertStateAllowsType(container, childType)
		markerID := this.unstackId()
		if _, exists := this.assignedIDs[markerID]; exists {
			panic(fmt.Errorf("%v: Marker ID already defined", markerID))
		}
		this.assignedIDs[markerID] = childType
		this.unstackState()
		this.onChildEnded(childType)
	case stateIdAwaitingReferenceID:
		container := this.getFirstRealContainer()
		markerID := this.unstackId()

		_, ok := markerID.(*url.URL)
		if ok {
			// We have no way to verify what the URL points to, so call it "anything".
			this.unstackState()
			this.onChildEnded(eventTypeAny)
			return
		}

		referencedType, ok := this.assignedIDs[markerID]
		if !ok {
			panic(fmt.Errorf("Referenced ID [%v] not found", markerID))
		}
		assertStateAllowsType(container, referencedType)
		this.unstackState()
		this.onChildEnded(referencedType)
	default:
		this.changeState(childEndRuleStateChanges[this.getCurrentStateId()])
	}
}

func (this *Rules) addScalar(scalarType ruleEvent) {
	this.assertCurrentStateAllowsType(scalarType)
	this.onChildEnded(scalarType)
}

func (this *Rules) beginContainer(containerType ruleEvent, newState ruleState) {
	this.assertCurrentStateAllowsType(containerType)
	this.stackState(newState)
}

var defaultRuleOptions = RuleOptions{
	ConciseEncodingVersion: ConciseEncodingVersion,
	MaxBytesLength:         1000000000,
	MaxStringLength:        100000000,
	MaxURILength:           10000,
	MaxIDLength:            100,
	MaxMarkupNameLength:    100,
	MaxContainerDepth:      1000,
	MaxObjectCount:         10000000,
	MaxReferenceCount:      100000,
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
	eventIdNothing ruleEvent = iota
	eventIdVersion
	eventIdPadding
	eventIdNil
	eventIdBool
	eventIdPInt
	eventIdNInt
	eventIdFloat
	eventIdNan
	eventIdUUID
	eventIdTime
	eventIdList
	eventIdMap
	eventIdMarkup
	eventIdMetadata
	eventIdComment
	eventIdMarker
	eventIdReference
	eventIdEndContainer
	eventIdBytes
	eventIdString
	eventIdURI
	eventIdCustom
	eventIdAChunk
	eventIdAData
	eventIdEndDocument
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

func (this ruleEvent) String() string {
	return ruleEventNames[this&ruleEvent(ruleIDFieldMask)]
}

type ruleState int

const (
	stateIdAwaitingNothing ruleState = iota
	stateIdAwaitingVersion
	stateIdAwaitingTLO
	stateIdAwaitingListItem
	stateIdAwaitingCommentItem
	stateIdAwaitingMapKey
	stateIdAwaitingMapValue
	stateIdAwaitingMetadataKey
	stateIdAwaitingMetadataValue
	stateIdAwaitingMetadataObject
	stateIdAwaitingMarkupName
	stateIdAwaitingMarkupKey
	stateIdAwaitingMarkupValue
	stateIdAwaitingMarkupContents
	stateIdAwaitingMarkerID
	stateIdAwaitingMarkerObject
	stateIdAwaitingReferenceID
	stateIdAwaitingArrayChunk
	stateIdAwaitingArrayData
	stateIdAwaitingEndDocument
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

func (this ruleState) String() string {
	return ruleStateNames[this&ruleState(ruleIDFieldMask)]
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
	eventTypeNothing      = eventIdNothing
	eventTypeVersion      = eventIdVersion | eventVersion
	eventTypePadding      = eventIdPadding | eventPadding
	eventTypeNil          = eventIdNil | eventNil
	eventTypeBool         = eventIdBool | eventScalar
	eventTypePInt         = eventIdPInt | eventPositiveInt
	eventTypeNInt         = eventIdNInt | eventScalar
	eventTypeFloat        = eventIdFloat | eventScalar
	eventTypeNan          = eventIdNan | eventNan
	eventTypeUUID         = eventIdUUID | eventScalar
	eventTypeTime         = eventIdTime | eventScalar
	eventTypeList         = eventIdList | eventBeginList
	eventTypeMap          = eventIdMap | eventBeginMap
	eventTypeMarkup       = eventIdMarkup | eventBeginMarkup
	eventTypeMetadata     = eventIdMetadata | eventBeginMetadata
	eventTypeComment      = eventIdComment | eventBeginComment
	eventTypeMarker       = eventIdMarker | eventBeginMarker
	eventTypeReference    = eventIdReference | eventBeginReference
	eventTypeEndContainer = eventIdEndContainer | eventEndContainer
	eventTypeBytes        = eventIdBytes | eventBeginBytes
	eventTypeString       = eventIdString | eventBeginString
	eventTypeURI          = eventIdURI | eventBeginURI
	eventTypeCustom       = eventIdCustom | eventBeginCustom
	eventTypeAChunk       = eventIdAChunk | eventArrayChunk
	eventTypeAData        = eventIdAData | eventArrayData
	eventTypeEndDocument  = eventIdEndDocument | eventEndDocument
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

	stateAwaitingNothing        = stateIdAwaitingNothing
	stateAwaitingVersion        = stateIdAwaitingVersion | allowVersion
	stateAwaitingTLO            = stateIdAwaitingTLO | allowTLO | ruleFlagRealContainer
	stateAwaitingEndDocument    = stateIdAwaitingEndDocument | allowEndDocument
	stateAwaitingListItem       = stateIdAwaitingListItem | allowListItem | ruleFlagRealContainer
	stateAwaitingMapKey         = stateIdAwaitingMapKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMapValue       = stateIdAwaitingMapValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMarkupName     = stateIdAwaitingMarkupName | allowMarkupName | ruleFlagRealContainer
	stateAwaitingMarkupKey      = stateIdAwaitingMarkupKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMarkupValue    = stateIdAwaitingMarkupValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMarkupContents = stateIdAwaitingMarkupContents | allowMarkupContents | ruleFlagRealContainer
	stateAwaitingMarkerID       = stateIdAwaitingMarkerID | allowMarkerID | ruleFlagAwaitingID
	stateAwaitingMarkerObject   = stateIdAwaitingMarkerObject | allowAny
	stateAwaitingReferenceID    = stateIdAwaitingReferenceID | allowReferenceID | ruleFlagAwaitingID
	stateAwaitingCommentItem    = stateIdAwaitingCommentItem | allowCommentItem /* Not a "real" container */
	stateAwaitingMetadataKey    = stateIdAwaitingMetadataKey | allowMapKey | ruleFlagRealContainer
	stateAwaitingMetadataValue  = stateIdAwaitingMetadataValue | allowMapValue | ruleFlagRealContainer
	stateAwaitingMetadataObject = stateIdAwaitingMetadataObject | allowAny
	stateAwaitingArrayChunk     = stateIdAwaitingArrayChunk | allowArrayChunk
	stateAwaitingArrayData      = stateIdAwaitingArrayData | allowArrayData
)

var childEndRuleStateChanges = [...]ruleState{
	/* stateIdAwaitingNothing                */ stateAwaitingNothing,
	/* stateIdAwaitingVersion              > */ stateAwaitingTLO,
	/* stateIdAwaitingTLO                  > */ stateAwaitingEndDocument,
	/* stateIdAwaitingListItem               */ stateAwaitingListItem,
	/* stateIdAwaitingCommentItem            */ stateAwaitingCommentItem,
	/* stateIdAwaitingMapKey               > */ stateAwaitingMapValue,
	/* stateIdAwaitingMapValue             > */ stateAwaitingMapKey,
	/* stateIdAwaitingMetadataKey          > */ stateAwaitingMetadataValue,
	/* stateIdAwaitingMetadataValue        > */ stateAwaitingMetadataKey,
	/* stateIdAwaitingMetadataObject         */ stateIdAwaitingMetadataObject,
	/* stateIdAwaitingMarkupName           > */ stateAwaitingMarkupKey,
	/* stateIdAwaitingMarkupAttributeKey   > */ stateAwaitingMarkupValue,
	/* stateIdAwaitingMarkupAttributeValue > */ stateAwaitingMarkupKey,
	/* stateIdAwaitingMarkupContents         */ stateAwaitingMarkupContents,
	/* stateIdAwaitingMarkerID             > */ stateAwaitingMarkerObject,
	/* stateIdAwaitingMarkerObject           */ stateAwaitingMarkerObject,
	/* stateIdAwaitingReferenceID            */ stateAwaitingReferenceID,
	/* stateIdAwaitingArrayChunk             */ stateAwaitingArrayChunk,
	/* stateIdAwaitingArrayData              */ stateAwaitingArrayData,
	/* stateIdAwaitingEndDocument          > */ stateAwaitingNothing,
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
