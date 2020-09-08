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

// Imposes the structural rules that enforce a well-formed concise encoding
// document.
package rules

import (
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/internal/unicode"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Rules is a DataEventsReceiver passthrough object that constrains the order
// and contents of events to ensure that they form a valid and complete Concise
// Encoding document.
//
// Put this right after your event generator in the event receiver chain to
// enforce correctly formed documents.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Rules struct {
	opts                  options.RuleOptions
	charValidator         UTF8Validator
	maxDepth              int
	stateStack            []ruleState
	arrayType             events.ArrayType
	arrayElementBitCount  int
	arrayData             []byte
	chunkByteCount        uint64
	chunkBytesWritten     uint64
	arrayBytesWritten     uint64
	moreChunksFollow      bool
	objectCount           uint64
	unassignedIDs         []interface{}
	assignedIDs           map[interface{}]ruleEvent
	unmatchedIDs          map[interface{}]bool
	realMaxContainerDepth uint64
	// TODO: Keep track of marked data stats (sizes etc)
	nextReceiver events.DataEventReceiver
}

// The initial rule state comes pre-stacked. This value accounts for it in calculations.
const rulesMaxDepthAdjust = 2

// Create a new rules set.
// If opts = nil, defaults are used.
func NewRules(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) *Rules {
	_this := new(Rules)
	_this.Init(nextReceiver, opts)
	return _this
}

// Initialize a rules set.
// If opts = nil, defaults are used.
func (_this *Rules) Init(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.realMaxContainerDepth = _this.opts.MaxContainerDepth + rulesMaxDepthAdjust
	_this.stateStack = make([]ruleState, 0, _this.realMaxContainerDepth)
	_this.nextReceiver = nextReceiver

	_this.Reset()
}

// Reset the rules set back to its initial state.
func (_this *Rules) Reset() {
	_this.stateStack = _this.stateStack[:0]
	_this.stackState(stateAwaitingEndDocument)
	_this.stackState(stateAwaitingBeginDocument)
	_this.unassignedIDs = _this.unassignedIDs[:0]
	_this.assignedIDs = make(map[interface{}]ruleEvent)
	_this.unmatchedIDs = make(map[interface{}]bool)

	_this.arrayType = events.ArrayTypeInvalid
	_this.arrayData = _this.arrayData[:0]
	_this.chunkByteCount = 0
	_this.chunkBytesWritten = 0
	_this.arrayBytesWritten = 0
	_this.moreChunksFollow = false
	_this.objectCount = 0
}

// ============================================================================

// DataEventReceiver

func (_this *Rules) OnBeginDocument() {
	_this.assertCurrentStateAllowsType(eventTypeBeginDocument)
	_this.changeState(stateAwaitingVersion)
	_this.nextReceiver.OnBeginDocument()
}

func (_this *Rules) OnVersion(version uint64) {
	_this.assertCurrentStateAllowsType(eventTypeVersion)
	if version != _this.opts.ConciseEncodingVersion {
		panic(fmt.Errorf("expected version %v but got version %v", _this.opts.ConciseEncodingVersion, version))
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
	if value == nil {
		_this.OnNil()
		return
	}

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
	if value == nil {
		_this.OnNil()
		return
	}

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
	if value == nil {
		_this.OnNil()
		return
	}

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
	if value == nil {
		_this.OnNil()
		return
	}

	_this.addScalar(eventTypeTime)
	_this.nextReceiver.OnCompactTime(value)
}

func (_this *Rules) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	// TODO: Valid map keys: make this nicer
	switch arrayType {
	case events.ArrayTypeString, events.ArrayTypeVerbatimString:
		// TODO: Isn't this done in onArrayData?
		_this.validateString(value)
	case events.ArrayTypeURI, events.ArrayTypeCustomBinary, events.ArrayTypeCustomText:
	// OK
	default:
		switch _this.getCurrentStateID() {
		case stateIDAwaitingMapKey, stateIDAwaitingMarkupKey, stateIDAwaitingMetadataKey:
			panic(fmt.Errorf("Array type %v not allowed as map key", arrayType))
		}
	}
	_this.beginArray(arrayType)
	_this.onArrayChunk(elementCount, false)
	if len(value) > 0 {
		_this.onArrayData(value)
	}
	_this.nextReceiver.OnArray(arrayType, elementCount, value)
}

func (_this *Rules) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString, events.ArrayTypeVerbatimString:
	// OK
	case events.ArrayTypeURI, events.ArrayTypeCustomBinary, events.ArrayTypeCustomText:
	// OK
	default:
		switch _this.getCurrentStateID() {
		case stateIDAwaitingMapKey, stateIDAwaitingMarkupKey, stateIDAwaitingMetadataKey:
			panic(fmt.Errorf("Array type %v not allowed as map key", arrayType))
		}
	}
	_this.beginArray(arrayType)
	_this.nextReceiver.OnArrayBegin(arrayType)
}

func (_this *Rules) OnArrayChunk(length uint64, moreChunksFollow bool) {
	_this.onArrayChunk(length, moreChunksFollow)
	_this.nextReceiver.OnArrayChunk(length, moreChunksFollow)
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
	if uint64(len(_this.assignedIDs)) >= _this.opts.MaxReferenceCount {
		panic(fmt.Errorf("Max number of marker IDs (%v) exceeded", _this.opts.MaxReferenceCount))
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
	for markerID := range _this.unmatchedIDs {
		_, ok := _this.assignedIDs[markerID]
		if !ok {
			panic(fmt.Errorf("unmatched reference to marker ID %v", markerID))
		}
	}
	_this.nextReceiver.OnEndDocument()
	_this.Reset()
}

// ============================================================================

// Internal

func (_this *Rules) onPositiveInt(value uint64) {
	if _this.isAwaitingID() {
		_this.stackID(value)
	}
	_this.addScalar(eventTypePInt)
}

func (_this *Rules) onNegativeInt() {
	_this.addScalar(eventTypeNInt)
}

func (_this *Rules) onArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.assertCurrentStateAllowsType(eventTypeAChunk)
	_this.chunkByteCount = common.ElementCountToByteCount(_this.arrayElementBitCount, elementCount)
	_this.chunkBytesWritten = 0
	_this.moreChunksFollow = moreChunksFollow
	_this.changeState(stateAwaitingArrayData)

	if elementCount == 0 {
		_this.onArrayChunkEnded()
	}
}

func (_this *Rules) validateString(data []byte) {
	if _this.isStringInsideComment() {
		_this.validateCommentContents(data)
	} else if _this.isAwaitingID() {
		_this.validateIDContents(data)
	} else {
		_this.validateStringContents(data)
	}
}

func (_this *Rules) onArrayData(data []byte) {
	_this.assertCurrentStateAllowsType(eventTypeAData)

	dataLength := uint64(len(data))
	if _this.chunkBytesWritten+dataLength > _this.chunkByteCount {
		panic(fmt.Errorf("chunk length %v exceeded by %v bytes",
			_this.chunkByteCount, _this.chunkBytesWritten+dataLength-_this.chunkByteCount))
	}

	switch _this.arrayType {
	case events.ArrayTypeString, events.ArrayTypeVerbatimString:
		if _this.arrayBytesWritten+dataLength > _this.opts.MaxStringLength {
			panic(fmt.Errorf("max string length (%v) exceeded", _this.opts.MaxStringLength))
		}
		_this.validateString(data)
		if _this.isAwaitingID() {
			_this.arrayData = append(_this.arrayData, data...)
		}
	case events.ArrayTypeURI:
		if _this.arrayBytesWritten+dataLength > _this.opts.MaxURILength {
			panic(fmt.Errorf("max URI length (%v) exceeded", _this.opts.MaxURILength))
		}
		if _this.isAwaitingID() {
			_this.arrayData = append(_this.arrayData, data...)
		}
		// Note: URI validation happens when the array is complete
	default:
		if _this.arrayBytesWritten+dataLength > _this.opts.MaxBytesLength {
			panic(fmt.Errorf("max byte array length (%v) exceeded", _this.opts.MaxBytesLength))
		}
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
	if uint64(len(_this.stateStack)) >= _this.realMaxContainerDepth {
		panic(fmt.Errorf("max depth of %v exceeded", _this.realMaxContainerDepth-rulesMaxDepthAdjust))
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
		_this.charValidator.AddByte(ch)
	}
}

func (_this *Rules) validateCommentContents(data []byte) {
	for _, ch := range data {
		_this.charValidator.AddByte(ch)
		if _this.charValidator.IsCompleteCharacter() {
			validateRulesCommentCharacter(_this.charValidator.GetCharacter())
		}
	}
}

func (_this *Rules) validateIDContents(data []byte) {
	for _, ch := range data {
		_this.charValidator.AddByte(ch)
		if _this.charValidator.IsCompleteCharacter() {
			validateRulesIDCharacter(_this.charValidator.GetCharacter())
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

func (_this *Rules) beginArray(arrayType events.ArrayType) {
	_this.assertCurrentStateAllowsType(arrayTypeToRuleEvent[arrayType])

	_this.arrayType = arrayType
	_this.arrayElementBitCount = arrayType.ElementSize()
	_this.arrayData = _this.arrayData[:0]
	_this.chunkByteCount = 0
	_this.chunkBytesWritten = 0
	_this.arrayBytesWritten = 0
	_this.moreChunksFollow = true

	_this.stackState(stateAwaitingArrayChunk)
}

func (_this *Rules) onArrayChunkEnded() {
	if _this.moreChunksFollow {
		_this.changeState(stateAwaitingArrayChunk)
		return
	}

	_this.unstackState()

	switch _this.arrayType {
	case events.ArrayTypeString, events.ArrayTypeVerbatimString:
		if _this.isAwaitingMarkupName() {

			if _this.arrayBytesWritten == 0 {
				panic(fmt.Errorf("markup name cannot be length 0"))
			}
			if _this.arrayBytesWritten > _this.opts.MaxMarkupNameLength {
				panic(fmt.Errorf("markup name length %v exceeds max of %v", _this.arrayBytesWritten, _this.opts.MaxMarkupNameLength))
			}
		}
		if _this.isAwaitingID() {
			if _this.arrayBytesWritten < minMarkerIDLength {
				panic(fmt.Errorf("an ID cannot be less than %v characters", minMarkerIDLength))
			}
			if _this.arrayBytesWritten > maxMarkerIDLength {
				panic(fmt.Errorf("ID length %v exceeds max of %v", _this.arrayBytesWritten, maxMarkerIDLength))
			}
			if _this.arrayData[0] >= '0' && _this.arrayData[0] <= '9' {
				panic(fmt.Errorf("ID first character cannot be a digit"))
			}
			_this.stackID(string(_this.arrayData))
		}
	case events.ArrayTypeURI:
		if _this.isAwaitingID() {
			uri, err := url.Parse(string(_this.arrayData))
			if err != nil {
				panic(fmt.Errorf("%v", err))
			}
			_this.stackID(uri)
		}
	}

	arrayType := _this.arrayType
	_this.arrayType = events.ArrayTypeInvalid
	_this.onChildEnded(arrayTypeToRuleEvent[arrayType])
}

func (_this *Rules) incrementObjectCount() {
	_this.objectCount++
	if _this.objectCount > _this.opts.MaxObjectCount {
		panic(fmt.Errorf("Max object count of %v exceeded", _this.opts.MaxObjectCount))
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
			_this.unmatchedIDs[markerID] = true
			// We have no way to verify what the unmatched ref points to, so call it "anything".
			_this.unstackState()
			_this.onChildEnded(eventTypeAny)
			return
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

type ruleEvent int

const (
	eventIDNothing ruleEvent = iota
	eventIDBeginDocument
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
	eventIDArray
	eventIDString
	eventIDVerbatimString
	eventIDURI
	eventIDAChunk
	eventIDAData
	eventIDEndDocument
)

var ruleEventNames = [...]string{
	eventIDNothing:        "nothing",
	eventIDBeginDocument:  "begin document",
	eventIDVersion:        "version",
	eventIDPadding:        "padding",
	eventIDNil:            "nil",
	eventIDBool:           "bool",
	eventIDPInt:           "positive int",
	eventIDNInt:           "negative int",
	eventIDFloat:          "float",
	eventIDNan:            "nan",
	eventIDUUID:           "UUID",
	eventIDTime:           "time",
	eventIDList:           "list",
	eventIDMap:            "map",
	eventIDMarkup:         "markup",
	eventIDMetadata:       "metadata",
	eventIDComment:        "comment",
	eventIDMarker:         "marker",
	eventIDReference:      "reference",
	eventIDEndContainer:   "end container",
	eventIDArray:          "array",
	eventIDString:         "string",
	eventIDVerbatimString: "verbatim string",
	eventIDURI:            "URI",
	eventIDAChunk:         "array chunk",
	eventIDAData:          "array data",
	eventIDEndDocument:    "end document",
}

func (_this ruleEvent) String() string {
	return ruleEventNames[_this&ruleIDFieldMask]
}

type ruleState int

const (
	stateIDAwaitingNothing ruleState = iota
	stateIDAwaitingBeginDocument
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
	stateIDAwaitingNothing:        "nothing",
	stateIDAwaitingBeginDocument:  "begin document",
	stateIDAwaitingVersion:        "version",
	stateIDAwaitingTLO:            "top-level object",
	stateIDAwaitingListItem:       "list item",
	stateIDAwaitingCommentItem:    "comment contents",
	stateIDAwaitingMapKey:         "map key",
	stateIDAwaitingMapValue:       "map value",
	stateIDAwaitingMetadataKey:    "metadata key",
	stateIDAwaitingMetadataValue:  "metadata value",
	stateIDAwaitingMetadataObject: "metadata object",
	stateIDAwaitingMarkupName:     "markup name",
	stateIDAwaitingMarkupKey:      "markup attribute key",
	stateIDAwaitingMarkupValue:    "markup attribute value",
	stateIDAwaitingMarkupContents: "markup contents",
	stateIDAwaitingMarkerID:       "marker ID",
	stateIDAwaitingMarkerObject:   "marker object",
	stateIDAwaitingReferenceID:    "reference id",
	stateIDAwaitingArrayChunk:     "array chunk",
	stateIDAwaitingArrayData:      "array data",
	stateIDAwaitingEndDocument:    "end document",
}

func (_this ruleState) String() string {
	return ruleStateNames[_this&ruleState(ruleIDFieldMask)]
}

const minMarkerIDLength = 1
const maxMarkerIDLength = 30

const (
	ruleIDFieldEnd  ruleEvent = 1 << 5
	ruleIDFieldMask           = ruleIDFieldEnd - 1
)

const (
	eventBeginDocument = ruleIDFieldEnd << iota
	eventVersion
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
	eventBeginArray
	eventBeginString
	eventBeginVerbatimString
	eventBeginURI
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
	eventTypeNothing        = eventIDNothing
	eventTypeBeginDocument  = eventIDBeginDocument | eventBeginDocument
	eventTypeVersion        = eventIDVersion | eventVersion
	eventTypePadding        = eventIDPadding | eventPadding
	eventTypeNil            = eventIDNil | eventNil
	eventTypeBool           = eventIDBool | eventScalar
	eventTypePInt           = eventIDPInt | eventPositiveInt
	eventTypeNInt           = eventIDNInt | eventScalar
	eventTypeFloat          = eventIDFloat | eventScalar
	eventTypeNan            = eventIDNan | eventNan
	eventTypeUUID           = eventIDUUID | eventScalar
	eventTypeTime           = eventIDTime | eventScalar
	eventTypeList           = eventIDList | eventBeginList
	eventTypeMap            = eventIDMap | eventBeginMap
	eventTypeMarkup         = eventIDMarkup | eventBeginMarkup
	eventTypeMetadata       = eventIDMetadata | eventBeginMetadata
	eventTypeComment        = eventIDComment | eventBeginComment
	eventTypeMarker         = eventIDMarker | eventBeginMarker
	eventTypeReference      = eventIDReference | eventBeginReference
	eventTypeEndContainer   = eventIDEndContainer | eventEndContainer
	eventTypeArray          = eventIDArray | eventBeginArray
	eventTypeString         = eventIDString | eventBeginString
	eventTypeVerbatimString = eventIDVerbatimString | eventBeginVerbatimString
	eventTypeURI            = eventIDURI | eventBeginURI
	eventTypeAChunk         = eventIDAChunk | eventArrayChunk
	eventTypeAData          = eventIDAData | eventArrayData
	eventTypeEndDocument    = eventIDEndDocument | eventEndDocument
	eventTypeAny            = ruleEventsMask
)

// Primary rules
const (
	eventsArray         = eventBeginArray | eventBeginString | eventBeginVerbatimString | eventBeginURI
	eventsInvisible     = eventPadding | eventBeginComment | eventBeginMetadata
	eventsKeyableObject = eventsInvisible | eventScalar | eventPositiveInt | eventsArray | eventBeginMarker
	eventsAnyObject     = eventsKeyableObject | eventNil | eventNan | eventBeginList | eventBeginMap | eventBeginMarkup | eventBeginReference
	allowAny            = ruleState(eventsAnyObject)
	allowTLO            = allowAny | ruleState(eventEndDocument)
	allowListItem       = allowAny | ruleState(eventEndContainer)
	allowMapKey         = ruleState(eventsKeyableObject | eventEndContainer)
	allowMapValue       = allowAny
	allowCommentItem    = ruleState(eventBeginString | eventBeginComment | eventEndContainer | eventPadding)
	allowMarkupName     = ruleState(eventPositiveInt | eventBeginString | eventPadding)
	allowMarkupContents = ruleState(eventBeginString | eventBeginVerbatimString | eventBeginComment | eventBeginMarkup | eventEndContainer | eventPadding)
	allowMarkerID       = ruleState(eventPositiveInt | eventBeginString | eventPadding)
	allowReferenceID    = ruleState(eventPositiveInt | eventBeginString | eventBeginURI | eventPadding)
	allowArrayChunk     = ruleState(eventArrayChunk)
	allowArrayData      = ruleState(eventArrayData)
	allowBeginDocument  = ruleState(eventBeginDocument)
	allowVersion        = ruleState(eventVersion)
	allowEndDocument    = ruleState(eventEndDocument | eventBeginComment | eventPadding)

	stateAwaitingNothing        = stateIDAwaitingNothing
	stateAwaitingBeginDocument  = stateIDAwaitingBeginDocument | allowBeginDocument
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
	stateIDAwaitingNothing:        stateAwaitingNothing,
	stateIDAwaitingBeginDocument:  stateAwaitingVersion,
	stateIDAwaitingVersion:        stateAwaitingTLO,
	stateIDAwaitingTLO:            stateAwaitingEndDocument,
	stateIDAwaitingListItem:       stateAwaitingListItem,
	stateIDAwaitingCommentItem:    stateAwaitingCommentItem,
	stateIDAwaitingMapKey:         stateAwaitingMapValue,
	stateIDAwaitingMapValue:       stateAwaitingMapKey,
	stateIDAwaitingMetadataKey:    stateAwaitingMetadataValue,
	stateIDAwaitingMetadataValue:  stateAwaitingMetadataKey,
	stateIDAwaitingMetadataObject: stateIDAwaitingMetadataObject,
	stateIDAwaitingMarkupName:     stateAwaitingMarkupKey,
	stateIDAwaitingMarkupKey:      stateAwaitingMarkupValue,
	stateIDAwaitingMarkupValue:    stateAwaitingMarkupKey,
	stateIDAwaitingMarkupContents: stateAwaitingMarkupContents,
	stateIDAwaitingMarkerID:       stateAwaitingMarkerObject,
	stateIDAwaitingMarkerObject:   stateAwaitingMarkerObject,
	stateIDAwaitingReferenceID:    stateAwaitingReferenceID,
	stateIDAwaitingArrayChunk:     stateAwaitingArrayChunk,
	stateIDAwaitingArrayData:      stateAwaitingArrayData,
	stateIDAwaitingEndDocument:    stateAwaitingNothing,
}

var arrayTypeToRuleEvent = [...]ruleEvent{
	events.ArrayTypeBoolean:        eventTypeArray,
	events.ArrayTypeUint8:          eventTypeArray,
	events.ArrayTypeUint16:         eventTypeArray,
	events.ArrayTypeUint32:         eventTypeArray,
	events.ArrayTypeUint64:         eventTypeArray,
	events.ArrayTypeInt8:           eventTypeArray,
	events.ArrayTypeInt16:          eventTypeArray,
	events.ArrayTypeInt32:          eventTypeArray,
	events.ArrayTypeInt64:          eventTypeArray,
	events.ArrayTypeFloat16:        eventTypeArray,
	events.ArrayTypeFloat32:        eventTypeArray,
	events.ArrayTypeFloat64:        eventTypeArray,
	events.ArrayTypeUUID:           eventTypeArray,
	events.ArrayTypeString:         eventTypeString,
	events.ArrayTypeVerbatimString: eventTypeString,
	events.ArrayTypeURI:            eventTypeURI,
	events.ArrayTypeCustomBinary:   eventTypeArray,
	events.ArrayTypeCustomText:     eventTypeArray,
}

func validateRulesCommentCharacter(ch rune) {
	const commentUnsafe = unicode.Control
	if unicode.CharHasProperty(ch, commentUnsafe) {
		panic(fmt.Errorf("0x%04x: Invalid comment character", ch))
	}
}

func validateRulesIDCharacter(ch rune) {
	if !unicode.CharHasProperty(ch, unicode.MarkerIDSafe) {
		panic(fmt.Errorf("0x%04x: Invalid ID character", ch))
	}
}
