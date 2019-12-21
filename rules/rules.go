// Rules encompasses every event that can occur within a CBE or CTE document
// during the encoding or decoding process.
package rules

import (
	"fmt"
	"math"
	"net/url"
)

// Limits before the ruleset artificially terminates with an error.
type Limits struct {
	MaxBytesLength      uint64
	MaxStringLength     uint64
	MaxURILength        uint64
	MaxIDLength         uint64
	MaxMarkupNameLength uint64
	MaxContainerDepth   uint64
	MaxObjectCount      uint64
	MaxReferenceCount   uint64
}

func DefaultLimits() *Limits {
	return &Limits{
		MaxBytesLength:      1000000000,
		MaxStringLength:     100000000,
		MaxURILength:        10000,
		MaxIDLength:         100,
		MaxMarkupNameLength: 100,
		MaxContainerDepth:   1000,
		MaxObjectCount:      1000000,
		MaxReferenceCount:   100000,
	}
}

func validateLimits(limits *Limits) {
	if limits.MaxBytesLength < 1 {
		panic(fmt.Errorf("MaxByteLength must be greater than 0"))
	}
	if limits.MaxStringLength < 1 {
		panic(fmt.Errorf("MaxStringLength must be greater than 0"))
	}
	if limits.MaxURILength < 1 {
		panic(fmt.Errorf("MaxURILength must be greater than 0"))
	}
	if limits.MaxIDLength < 1 {
		panic(fmt.Errorf("MaxIDLength must be greater than 0"))
	}
	if limits.MaxMarkupNameLength < 1 {
		panic(fmt.Errorf("MaxMarkupNameLength must be greater than 0"))
	}
	if limits.MaxContainerDepth < 1 {
		panic(fmt.Errorf("MaxContainerDepth must be greater than 0"))
	}
	if limits.MaxObjectCount < 1 {
		panic(fmt.Errorf("MaxObjectCount must be greater than 0"))
	}
}

// The initial rule state comes pre-stacked. This value accounts for it in calculations.
const maxDepthAdjust = 2

type event int

const (
	eventIdNothing event = iota
	eventIdVersion
	eventIdPadding
	eventIdNil
	eventIdBool
	eventIdPositiveInt
	eventIdNegativeInt
	eventIdFloat
	eventIdNan
	eventIdTime
	eventIdList
	eventIdMap
	eventIdMarkup
	eventIdMetadata
	eventIdComment
	eventIdMarker
	eventIdReference
	eventIdEndContainer
	eventIdEndList
	eventIdEndMap
	eventIdEndMarkupAttributes
	eventIdEndMarkup
	eventIdEndMetadata
	eventIdEndComment
	eventIdBytes
	eventIdString
	eventIdURI
	eventIdArrayChunk
	eventIdArrayData
	eventIdEndDocument
)

var eventNames = [...]string{
	"nothing",
	"version",
	"padding",
	"nil",
	"bool",
	"positive int",
	"negative int",
	"float",
	"nan",
	"time",
	"list",
	"map",
	"markup",
	"metadata",
	"comment",
	"marker",
	"reference",
	"end container",
	"end list",
	"end map",
	"end markup attributes",
	"end markup",
	"end metadata",
	"end comment",
	"bytes",
	"string",
	"URI",
	"array chunk",
	"array data",
	"end document",
}

func getEventName(index event) string {
	return eventNames[index&event(idMask)]
}

type state int

const (
	stateIdAwaitingNothing state = iota
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
	stateIdAwaitingMarkupAttributeKey
	stateIdAwaitingMarkupAttributeValue
	stateIdAwaitingMarkupContents
	stateIdAwaitingMarkerID
	stateIdAwaitingMarkerObject
	stateIdAwaitingReferenceID
	stateIdAwaitingArrayChunk
	stateIdAwaitingArrayData
	stateIdAwaitingEndDocument
)

var stateNames = [...]string{
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

func getStateName(index state) string {
	return stateNames[index&state(idMask)]
}

const (
	idEnd  event = 1 << 5
	idMask       = idEnd - 1
)

const (
	eventVersion = event(idEnd) << iota
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
	eventEndList
	eventEndMap
	eventEndMarkupAttributes
	eventEndMarkup
	eventEndMetadata
	eventEndComment
	eventBeginBytes
	eventBeginString
	eventBeginURI
	eventArrayChunk
	eventArrayData
	eventEndDocument
	eventsEnd
	eventsMask = (eventsEnd - 1) - (idEnd - 1)
)

const (
	flagRealContainer = state(eventsEnd) << iota
	flagAwaitingID
	flagsEnd
	flagsMask = (flagsEnd - 1) - (state(eventsEnd) - 1)
)

const (
	typeNothing             = eventIdNothing
	typeVersion             = eventIdVersion | eventVersion
	typePadding             = eventIdPadding | eventPadding
	typeNil                 = eventIdNil | eventNil
	typeBool                = eventIdBool | eventScalar
	typePositiveInt         = eventIdPositiveInt | eventPositiveInt
	typeNegativeInt         = eventIdNegativeInt | eventScalar
	typeFloat               = eventIdFloat | eventScalar
	typeNan                 = eventIdNan | eventNan
	typeTime                = eventIdTime | eventScalar
	typeList                = eventIdList | eventBeginList
	typeMap                 = eventIdMap | eventBeginMap
	typeMarkup              = eventIdMarkup | eventBeginMarkup
	typeMetadata            = eventIdMetadata | eventBeginMetadata
	typeComment             = eventIdComment | eventBeginComment
	typeMarker              = eventIdMarker | eventBeginMarker
	typeReference           = eventIdReference | eventBeginReference
	typeEndContainer        = eventIdEndContainer | eventEndList | eventEndMap | eventEndMarkupAttributes | eventEndMarkup | eventEndComment | eventEndMetadata
	typeEndList             = eventIdEndList | eventEndList
	typeEndMap              = eventIdEndMap | eventEndMap
	typeEndMarkupAttributes = eventIdEndMarkupAttributes | eventEndMarkupAttributes
	typeEndMarkup           = eventIdEndMarkup | eventEndMarkup
	typeEndMetadata         = eventIdEndMetadata | eventEndMetadata
	typeEndComment          = eventIdEndComment | eventEndComment
	typeBytes               = eventIdBytes | eventBeginBytes
	typeString              = eventIdString | eventBeginString
	typeURI                 = eventIdURI | eventBeginURI
	typeArrayChunk          = eventIdArrayChunk | eventArrayChunk
	typeArrayData           = eventIdArrayData | eventArrayData
	typeEndDocument         = eventIdEndDocument | eventEndDocument
	typeAny                 = eventsMask
)

// Primary rules
const (
	eventsArray              = eventBeginBytes | eventBeginString | eventBeginURI
	eventsInvisible          = eventPadding | eventBeginComment | eventBeginMetadata
	eventsKeyableObject      = eventsInvisible | eventScalar | eventPositiveInt | eventsArray | eventBeginMarker | eventBeginReference
	eventsAnyObject          = eventsKeyableObject | eventNil | eventNan | eventBeginList | eventBeginMap | eventBeginMarkup
	allowAny                 = state(eventsAnyObject)
	allowTLO                 = allowAny | state(eventEndDocument)
	allowListItem            = allowAny | state(eventEndList)
	allowMapKey              = state(eventsKeyableObject | eventEndMap)
	allowMetadataKey         = state(eventsKeyableObject | eventEndMetadata)
	allowMarkupAttributesKey = state(eventsKeyableObject | eventEndMarkupAttributes)
	allowMapValue            = allowAny
	allowCommentItem         = state(eventBeginString | eventBeginComment | eventEndComment | eventPadding)
	allowMarkupName          = state(eventPositiveInt | eventBeginString | eventPadding)
	allowMarkupContents      = state(eventBeginString | eventBeginComment | eventBeginMarkup | eventEndMarkup | eventPadding)
	allowMarkerID            = state(eventPositiveInt | eventBeginString | eventPadding)
	allowReferenceID         = state(eventPositiveInt | eventBeginString | eventBeginURI | eventPadding)
	allowArrayChunk          = state(eventArrayChunk)
	allowArrayData           = state(eventArrayData)
	allowVersion             = state(eventVersion)
	allowEndDocument         = state(eventEndDocument | eventBeginComment | eventPadding)

	stateAwaitingNothing        = stateIdAwaitingNothing
	stateAwaitingVersion        = stateIdAwaitingVersion | allowVersion
	stateAwaitingTLO            = stateIdAwaitingTLO | allowTLO | flagRealContainer
	stateAwaitingEndDocument    = stateIdAwaitingEndDocument | allowEndDocument
	stateAwaitingListItem       = stateIdAwaitingListItem | allowListItem | flagRealContainer
	stateAwaitingMapKey         = stateIdAwaitingMapKey | allowMapKey | flagRealContainer
	stateAwaitingMapValue       = stateIdAwaitingMapValue | allowMapValue | flagRealContainer
	stateAwaitingMarkupName     = stateIdAwaitingMarkupName | allowMarkupName | flagRealContainer
	stateAwaitingMarkupKey      = stateIdAwaitingMarkupAttributeKey | allowMarkupAttributesKey | flagRealContainer
	stateAwaitingMarkupValue    = stateIdAwaitingMarkupAttributeValue | allowMapValue | flagRealContainer
	stateAwaitingMarkupContents = stateIdAwaitingMarkupContents | allowMarkupContents | flagRealContainer
	stateAwaitingMarkerID       = stateIdAwaitingMarkerID | allowMarkerID | flagAwaitingID
	stateAwaitingMarkerObject   = stateIdAwaitingMarkerObject | allowAny
	stateAwaitingReferenceID    = stateIdAwaitingReferenceID | allowReferenceID | flagAwaitingID
	stateAwaitingCommentItem    = stateIdAwaitingCommentItem | allowCommentItem /* Not a "real" container */
	stateAwaitingMetadataKey    = stateIdAwaitingMetadataKey | allowMetadataKey | flagRealContainer
	stateAwaitingMetadataValue  = stateIdAwaitingMetadataValue | allowMapValue | flagRealContainer
	stateAwaitingMetadataObject = stateIdAwaitingMetadataObject | allowAny
	stateAwaitingArrayChunk     = stateIdAwaitingArrayChunk | allowArrayChunk
	stateAwaitingArrayData      = stateIdAwaitingArrayData | allowArrayData

	DecoderStateAwaitingVersion     = stateIdAwaitingVersion
	DecoderStateAwaitingArrayChunk  = stateIdAwaitingArrayChunk
	DecoderStateAwaitingArrayData   = stateIdAwaitingArrayData
	DecoderStateAwaitingEndDocument = stateIdAwaitingEndDocument
	DecoderStateEnded               = stateIdAwaitingNothing
)

var childEndStateChanges = [...]state{
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

func reportError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func validateCommentCharacter(ch int) error {
	switch ch {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08 /*, 0x09, 0x0a*/, 0x0b, 0x0c /*, 0x0d*/, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x7f,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0x2028, 0x2029:
		return reportError("0x%04x: Invalid comment character", ch)
	default:
		return nil
	}
	return nil
}

type Rules struct {
	version           uint64
	charValidator     Utf8Validator
	limits            *Limits
	maxDepth          int
	stateStack        []state
	arrayType         event
	arrayData         []byte
	chunkByteCount    uint64
	chunkBytesWritten uint64
	arrayBytesWritten uint64
	isFinalChunk      bool
	objectCount       uint64
	unassignedIDs     []interface{}
	assignedIDs       map[interface{}]event
}

func (this *Rules) getCurrentState() state {
	return this.stateStack[len(this.stateStack)-1]
}

func (this *Rules) getCurrentStateId() state {
	return this.getCurrentState() & state(idMask)
}

func (this *Rules) getParentState() state {
	return this.stateStack[len(this.stateStack)-2]
}

func (this *Rules) changeState(st state) {
	this.stateStack[len(this.stateStack)-1] = st
}

func (this *Rules) stackState(st state) error {
	if uint64(len(this.stateStack)) >= this.limits.MaxContainerDepth {
		return reportError("Max depth of %v exceeded", this.limits.MaxContainerDepth-maxDepthAdjust)
	}
	this.stateStack = append(this.stateStack, st)
	return nil
}

func (this *Rules) unstackState() {
	this.stateStack = this.stateStack[:len(this.stateStack)-1]
}

func (this *Rules) isAwaitingID() bool {
	if this.getCurrentState()&state(eventArrayChunk|eventArrayData) != 0 {
		return this.getParentState()&flagAwaitingID != 0
	}
	return this.getCurrentState()&flagAwaitingID != 0
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
	return len(this.stateStack) >= 2 && this.stateStack[len(this.stateStack)-2]&stateAwaitingCommentItem != 0
}

func (this *Rules) validateStringContents(data []byte) error {
	for _, ch := range data {
		if err := this.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
	}
	return nil
}

func (this *Rules) validateCommentContents(data []byte) error {
	for _, ch := range data {
		if err := this.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
		if this.charValidator.IsCompleteCharacter() {
			if err := validateCommentCharacter(this.charValidator.Character()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Rules) getFirstRealContainer() (state, error) {
	for i := len(this.stateStack) - 1; i >= 0; i-- {
		currentState := this.stateStack[i]
		if currentState&flagRealContainer != 0 {
			return currentState, nil
		}
	}
	return stateAwaitingNothing, reportError("BUG: Could not find real container in state stack")
}

func assertStateAllowsType(currentState state, objectType event) error {
	allowedEventMask := event(currentState) & eventsMask
	if objectType&allowedEventMask == 0 {
		return reportError("%v not allowed while awaiting %v", getEventName(objectType), getStateName(currentState))
	}
	return nil
}

func (this *Rules) assertCurrentStateAllowsType(objectType event) error {
	return assertStateAllowsType(this.getCurrentState(), objectType)
}

func (this *Rules) beginArray(arrayType event) error {
	if err := this.assertCurrentStateAllowsType(arrayType); err != nil {
		return err
	}

	this.arrayType = arrayType
	this.arrayData = this.arrayData[:0]
	this.chunkByteCount = 0
	this.chunkBytesWritten = 0
	this.arrayBytesWritten = 0
	this.isFinalChunk = false

	return this.stackState(stateAwaitingArrayChunk)
}

func (this *Rules) onArrayChunkEnded() error {
	if !this.isFinalChunk {
		this.changeState(stateAwaitingArrayChunk)
		return nil
	}

	this.unstackState()

	switch this.arrayType {
	case typeString:
		if this.isAwaitingMarkupName() {

			if this.arrayBytesWritten == 0 {
				return reportError("Markup name cannot be length 0")
			}
			if this.arrayBytesWritten > this.limits.MaxMarkupNameLength {
				return reportError("Markup name length %v exceeds max of %v", this.arrayBytesWritten, this.limits.MaxMarkupNameLength)
			}
		}
		if this.isAwaitingID() {
			if this.arrayBytesWritten == 0 {
				return reportError("An ID cannot be length 0")
			}
			if this.arrayBytesWritten > this.limits.MaxIDLength {
				return reportError("ID length %v exceeds max of %v", this.arrayBytesWritten, this.limits.MaxIDLength)
			}
			this.stackId(string(this.arrayData))
		}
	case typeURI:
		if this.arrayBytesWritten < 2 {
			return reportError("URI length must allow at least a scheme and colon (2 chars)")
		}
		if this.isAwaitingID() {
			url, err := url.Parse(string(this.arrayData))
			if err != nil {
				return reportError("%v", err)
			}
			this.stackId(url)
		}
	case typeBytes:
		// Nothing to do
	}

	arrayType := this.arrayType
	this.arrayType = typeNothing
	return this.onChildEnded(arrayType)
}

func (this *Rules) incrementObjectCount() error {
	this.objectCount++
	if this.objectCount > this.limits.MaxObjectCount {
		return reportError("Max object count of %v exceeded", this.limits.MaxObjectCount)
	}
	return nil
}

func (this *Rules) onChildEnded(childType event) error {
	if err := this.incrementObjectCount(); err != nil {
		return err
	}

	switch this.getCurrentStateId() {
	case stateIdAwaitingMetadataObject:
		container, err := this.getFirstRealContainer()
		if err != nil {
			return err
		}
		if err = assertStateAllowsType(container, childType); err != nil {
			return err
		}
		this.unstackState()
		return this.onChildEnded(childType)
	case stateIdAwaitingMarkerObject:
		container, err := this.getFirstRealContainer()
		if err != nil {
			return err
		}
		if err = assertStateAllowsType(container, childType); err != nil {
			return err
		}
		markerID := this.unstackId()
		if _, exists := this.assignedIDs[markerID]; exists {
			return reportError("%v: Marker ID already defined", markerID)
		}
		this.assignedIDs[markerID] = childType
		this.unstackState()
		return this.onChildEnded(childType)
	case stateIdAwaitingReferenceID:
		container, err := this.getFirstRealContainer()
		if err != nil {
			return err
		}
		markerID := this.unstackId()

		_, ok := markerID.(url.URL)
		if ok {
			// We have no way to verify what the URL points to, so call it "anything".
			this.unstackState()
			return this.onChildEnded(typeAny)
		}

		referencedType, ok := this.assignedIDs[markerID]
		if !ok {
			return reportError("Referenced ID [%v] not found", markerID)
		}
		if err = assertStateAllowsType(container, referencedType); err != nil {
			return err
		}
		this.unstackState()
		return this.onChildEnded(referencedType)
	default:
		this.changeState(childEndStateChanges[this.getCurrentStateId()])
		return nil
	}
}

func (this *Rules) addScalar(scalarType event) error {
	if err := this.assertCurrentStateAllowsType(scalarType); err != nil {
		return err
	}

	return this.onChildEnded(scalarType)
}

func (this *Rules) beginContainer(containerType event, newState state) error {
	if err := this.assertCurrentStateAllowsType(containerType); err != nil {
		return err
	}

	return this.stackState(newState)
}

// ----------
// Public API
// ----------

func NewRules(version uint64, limits *Limits) *Rules {
	this := new(Rules)
	this.Init(version, limits)
	return this
}

func (this *Rules) Init(version uint64, limits *Limits) {
	this.version = version
	if limits == nil {
		this.limits = DefaultLimits()
	} else {
		limitsCopy := *limits
		this.limits = &limitsCopy
	}

	validateLimits(this.limits)
	this.limits.MaxContainerDepth += maxDepthAdjust
	this.stateStack = make([]state, 0, this.limits.MaxContainerDepth)

	this.Reset()
}

func (this *Rules) Reset() {
	this.stateStack = this.stateStack[:0]
	if err := this.stackState(stateAwaitingEndDocument); err != nil {
		// Should not happen
		panic(err)
	}
	if err := this.stackState(stateAwaitingVersion); err != nil {
		// should not happen
		panic(err)
	}
	this.unassignedIDs = this.unassignedIDs[:0]
	this.assignedIDs = make(map[interface{}]event)

	this.arrayType = typeNothing
	this.arrayData = this.arrayData[:0]
	this.chunkByteCount = 0
	this.chunkBytesWritten = 0
	this.arrayBytesWritten = 0
	this.isFinalChunk = false
	this.objectCount = 0
}

func (this *Rules) IsInArray() bool {
	return this.arrayType != typeNothing
}

func (this *Rules) IsDocumentComplete() bool {
	return this.getCurrentState() == stateAwaitingEndDocument
}

func (this *Rules) GetRemainingChunkByteCount() uint64 {
	return this.chunkByteCount - this.chunkBytesWritten
}

func (this *Rules) AddVersion(version uint64) (err error) {
	if err = this.assertCurrentStateAllowsType(typeVersion); err != nil {
		return
	}
	if version != this.version {
		return fmt.Errorf("Expected version %v but got version %v", this.version, version)
	}
	this.changeState(stateAwaitingTLO)
	return
}

func (this *Rules) AddPadding() (err error) {
	return this.assertCurrentStateAllowsType(typePadding)
}

func (this *Rules) AddNil() error {
	return this.addScalar(typeNil)
}

func (this *Rules) AddBool() error {
	return this.addScalar(typeBool)
}

// value is required because it could be used as a marker ID or reference ID in
// the document
func (this *Rules) AddPositiveInt(value uint64) error {
	if this.isAwaitingID() {
		this.stackId(value)
	}
	return this.addScalar(typePositiveInt)
}

func (this *Rules) AddNegativeInt() error {
	return this.addScalar(typeNegativeInt)
}

func (this *Rules) AddFloat(value float64) error {
	if math.IsNaN(value) {
		return this.AddNan()
	}
	return this.addScalar(typeFloat)
}

func (this *Rules) AddNan() error {
	return this.addScalar(typeNan)
}

func (this *Rules) AddTime() error {
	return this.addScalar(typeTime)
}

func (this *Rules) BeginList() error {
	return this.beginContainer(typeList, stateAwaitingListItem)
}

func (this *Rules) BeginMap() error {
	return this.beginContainer(typeMap, stateAwaitingMapKey)
}

func (this *Rules) BeginMarkup() error {
	return this.beginContainer(typeMarkup, stateAwaitingMarkupName)
}

func (this *Rules) BeginMetadata() error {
	return this.beginContainer(typeMetadata, stateAwaitingMetadataKey)
}

func (this *Rules) BeginComment() error {
	return this.beginContainer(typeComment, stateAwaitingCommentItem)
}

func (this *Rules) EndList() error {
	if err := this.assertCurrentStateAllowsType(typeEndList); err != nil {
		return err
	}
	this.unstackState()
	return this.onChildEnded(typeList)
}

func (this *Rules) EndMap() error {
	if err := this.assertCurrentStateAllowsType(typeEndMap); err != nil {
		return err
	}
	this.unstackState()
	return this.onChildEnded(typeMap)
}

func (this *Rules) EndMarkupAttributes() error {
	if err := this.assertCurrentStateAllowsType(typeEndMarkupAttributes); err != nil {
		return err
	}
	this.changeState(stateAwaitingMarkupContents)
	return nil
}

func (this *Rules) EndMarkup() error {
	if err := this.assertCurrentStateAllowsType(typeEndMarkup); err != nil {
		return err
	}
	this.unstackState()
	return this.onChildEnded(typeMarkup)
}

func (this *Rules) EndMetadata() error {
	if err := this.assertCurrentStateAllowsType(typeEndMetadata); err != nil {
		return err
	}
	this.changeState(stateAwaitingMetadataObject)
	return this.incrementObjectCount()
}

func (this *Rules) EndComment() error {
	if err := this.assertCurrentStateAllowsType(typeEndComment); err != nil {
		return err
	}
	this.unstackState()
	return this.incrementObjectCount()
}

func (this *Rules) EndContainer() error {
	if err := this.assertCurrentStateAllowsType(typeEndContainer); err != nil {
		return err
	}

	switch this.getCurrentStateId() {
	case stateIdAwaitingListItem:
		return this.EndList()
	case stateIdAwaitingMapKey:
		return this.EndMap()
	case stateIdAwaitingMarkupAttributeKey:
		return this.EndMarkupAttributes()
	case stateIdAwaitingMarkupContents:
		return this.EndMarkup()
	case stateIdAwaitingMetadataKey:
		return this.EndMetadata()
	case stateIdAwaitingCommentItem:
		return this.EndComment()
	default:
		return reportError("BUG: EndContainer() in state %x (%v) failed to trigger", getStateName(this.getCurrentState()), this.getCurrentState())
	}
}

func (this *Rules) BeginBytes() error {
	return this.beginArray(typeBytes)
}

func (this *Rules) BeginString() error {
	return this.beginArray(typeString)
}

func (this *Rules) BeginURI() error {
	return this.beginArray(typeURI)
}

func (this *Rules) BeginArrayChunk(length uint64, isFinalChunk bool) error {
	if err := this.assertCurrentStateAllowsType(typeArrayChunk); err != nil {
		return err
	}

	this.chunkByteCount = length
	this.chunkBytesWritten = 0
	this.isFinalChunk = isFinalChunk
	this.changeState(stateAwaitingArrayData)

	if length == 0 {
		return this.onArrayChunkEnded()
	}
	return nil
}

func (this *Rules) AddArrayData(data []byte) error {
	if err := this.assertCurrentStateAllowsType(typeArrayData); err != nil {
		return err
	}

	dataLength := uint64(len(data))
	if this.chunkBytesWritten+dataLength > this.chunkByteCount {
		return reportError("Chunk length %v exceeded by %v bytes",
			this.chunkByteCount, this.chunkBytesWritten+dataLength-this.chunkByteCount)
	}

	switch this.arrayType {
	case typeBytes:
		if this.arrayBytesWritten+dataLength > this.limits.MaxBytesLength {
			return reportError("Max byte array length (%v) exceeded", this.limits.MaxBytesLength)
		}
	case typeString:
		if this.arrayBytesWritten+dataLength > this.limits.MaxStringLength {
			return reportError("Max string length (%v) exceeded", this.limits.MaxStringLength)
		}
		if this.isStringInsideComment() {
			if err := this.validateCommentContents(data); err != nil {
				return err
			}
		} else {
			if err := this.validateStringContents(data); err != nil {
				return err
			}
		}
		if this.isAwaitingID() {
			this.arrayData = append(this.arrayData, data...)
		}
	case typeURI:
		if this.arrayBytesWritten+dataLength > this.limits.MaxURILength {
			return reportError("Max URI length (%v) exceeded", this.limits.MaxURILength)
		}
		// TODO: URI validation
	}

	this.arrayBytesWritten += dataLength
	this.chunkBytesWritten += dataLength
	if this.chunkBytesWritten == this.chunkByteCount {
		return this.onArrayChunkEnded()
	}
	return nil
}

func (this *Rules) BeginMarker() error {
	if uint64(len(this.assignedIDs)) >= this.limits.MaxReferenceCount {
		return reportError("Max number of marker IDs (%v) exceeded", this.limits.MaxReferenceCount)
	}
	return this.beginContainer(typeMarker, stateAwaitingMarkerID)
}

func (this *Rules) BeginReference() error {
	return this.beginContainer(typeReference, stateAwaitingReferenceID)
}

func (this *Rules) EndDocument() error {
	return this.assertCurrentStateAllowsType(typeEndDocument)
}

func (this *Rules) GetDecoderState() state {
	return this.getCurrentState() & state(idMask)
}
