package cbe

import (
	"fmt"
	"time"
)

// Callback functions that must be present in the receiver object.
type CbeDecoderCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnPositiveInt(value uint64) error
	OnNegativeInt(value uint64) error
	OnFloat(value float64) error
	OnDate(value time.Time) error
	OnTime(value time.Time) error
	OnTimestamp(value time.Time) error
	OnListBegin() error
	OnOrderedMapBegin() error
	OnUnorderedMapBegin() error
	OnMetadataMapBegin() error
	OnContainerEnd() error
	OnBytesBegin(byteCount uint64) error
	OnStringBegin(byteCount uint64) error
	OnURIBegin(byteCount uint64) error
	OnCommentBegin(byteCount uint64) error
	OnArrayData(bytes []byte) error
}

// Biggest item is timestamp (10 bytes), longest tz is "America/Argentina/ComodRivadavia"
const maxPartialReadLength = 50

type decoderError struct {
	err error
}

type callbackError struct {
	err error
}

type containerData struct {
	inlineContainerType        ContainerType
	inlineContainerInitialized bool
	depth                      int
	currentType                []ContainerType
	hasProcessedMapKey         []bool
}

type arrayData struct {
	currentType        arrayType
	remainingByteCount int64
}

type CbeDecoder struct {
	streamOffset     int64
	buffer           *decodeBuffer
	mainBuffer       *decodeBuffer
	underflowBuffer  *decodeBuffer
	container        containerData
	array            arrayData
	callbacks        CbeDecoderCallbacks
	firstItemDecoded bool
	charValidator    Utf8Validator
}

func panicOnCallbackError(err error) {
	if err != nil {
		panic(callbackError{err})
	}
}

func (this *CbeDecoder) isExpectingMapKey() bool {
	return this.container.currentType[this.container.depth] == ContainerTypeUnorderedMap &&
		!this.container.hasProcessedMapKey[this.container.depth]
}

func (this *CbeDecoder) isExpectingMapValue() bool {
	return this.container.currentType[this.container.depth] == ContainerTypeUnorderedMap &&
		this.container.hasProcessedMapKey[this.container.depth]
}

func (this *CbeDecoder) flipMapKeyValueState() {
	this.container.hasProcessedMapKey[this.container.depth] = !this.container.hasProcessedMapKey[this.container.depth]
}

func (this *CbeDecoder) assertNotExpectingMapKey(keyType string) {
	if this.isExpectingMapKey() {
		panic(decoderError{fmt.Errorf("Cannot use type %v as a map key", keyType)})
	}
}

func (this *CbeDecoder) containerBegin(newContainerType ContainerType) {
	if this.container.depth+1 >= len(this.container.currentType) {
		panic(decoderError{fmt.Errorf("Exceeded max container depth of %v", len(this.container.currentType))})
	}
	this.container.depth++
	this.container.currentType[this.container.depth] = newContainerType
	this.container.hasProcessedMapKey[this.container.depth] = false
}

func (this *CbeDecoder) containerEnd() ContainerType {
	if this.container.depth <= 0 {
		panic(decoderError{fmt.Errorf("Got container end but not in a container")})
	}
	if this.container.inlineContainerType != ContainerTypeNone && this.container.depth <= 1 {
		panic(decoderError{fmt.Errorf("Got container end but not in a container")})
	}
	if this.isExpectingMapValue() {
		panic(decoderError{fmt.Errorf("Expecting map value for already processed key")})
	}

	this.container.depth--
	return this.container.currentType[this.container.depth+1]
}

func (this *CbeDecoder) arrayBegin(newArrayType arrayType, length int64) {
	this.array.currentType = newArrayType
	this.charValidator.Reset()
	this.array.remainingByteCount = length

	switch newArrayType {
	case arrayTypeBytes:
		panicOnCallbackError(this.callbacks.OnBytesBegin(uint64(length)))
	case arrayTypeComment:
		panicOnCallbackError(this.callbacks.OnCommentBegin(uint64(length)))
	case arrayTypeString:
		panicOnCallbackError(this.callbacks.OnStringBegin(uint64(length)))
	default:
		panic(fmt.Errorf("BUG: Unhandled array type: %v", newArrayType))
	}
}

func (this *CbeDecoder) decodeArrayLength(buffer *decodeBuffer) {
	this.array.remainingByteCount = buffer.DecodeArrayLength()
}

func (this *CbeDecoder) decodeArrayData() {
	if this.array.currentType == arrayTypeNone {
		return
	}

	if this.array.remainingByteCount > 0 {
		decodeByteCount := this.buffer.RemainingByteCount()
		if int64(decodeByteCount) > this.array.remainingByteCount {
			decodeByteCount = int(this.array.remainingByteCount)
		}
		bytes := this.buffer.DecodeBytes(decodeByteCount)
		if this.array.currentType == arrayTypeString || this.array.currentType == arrayTypeComment {
			for _, ch := range bytes {
				if err := this.charValidator.AddByte(int(ch)); err != nil {
					panic(decoderError{err})
				}
				if this.charValidator.IsCompleteCharacter() && this.array.currentType == arrayTypeComment {
					if err := ValidateCommentCharacter(this.charValidator.Character()); err != nil {
						panic(decoderError{err})
					}
				}
			}
		}
		this.array.remainingByteCount -= int64(decodeByteCount)

		panicOnCallbackError(this.callbacks.OnArrayData(bytes))
		this.buffer.Commit()

		if this.array.remainingByteCount > 0 {
			panic(notEnoughBytesToDecodeArrayData(this.array.remainingByteCount))
		}
	}

	this.array.currentType = arrayTypeNone
	this.flipMapKeyValueState()
}

func (this *CbeDecoder) decodeStringOfLength(length int64) {
	this.arrayBegin(arrayTypeString, length)
	this.decodeArrayData()
}

func (this *CbeDecoder) decodeObject(dataType typeField) {
	asSmallInt := int8(dataType)
	if int64(asSmallInt) >= smallIntMin && int64(asSmallInt) <= smallIntMax {
		if asSmallInt >= 0 {
			panicOnCallbackError(this.callbacks.OnPositiveInt(uint64(asSmallInt)))
		} else {
			panicOnCallbackError(this.callbacks.OnNegativeInt(uint64(-asSmallInt)))
		}
		this.buffer.Commit()
		this.flipMapKeyValueState()
		return
	}

	switch dataType {
	case typeTrue:
		panicOnCallbackError(this.callbacks.OnBool(true))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeFalse:
		panicOnCallbackError(this.callbacks.OnBool(false))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeFloat32:
		panicOnCallbackError(this.callbacks.OnFloat(float64(this.buffer.DecodeFloat32())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeFloat64:
		panicOnCallbackError(this.callbacks.OnFloat(this.buffer.DecodeFloat64()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeDecimal:
		panicOnCallbackError(this.callbacks.OnFloat(this.buffer.DecodeFloat()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePosInt8:
		panicOnCallbackError(this.callbacks.OnPositiveInt(uint64(this.buffer.DecodeUint8())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePosInt16:
		panicOnCallbackError(this.callbacks.OnPositiveInt(uint64(this.buffer.DecodeUint16())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePosInt32:
		panicOnCallbackError(this.callbacks.OnPositiveInt(uint64(this.buffer.DecodeUint32())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePosInt64:
		panicOnCallbackError(this.callbacks.OnPositiveInt(this.buffer.DecodeUint64()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePosInt:
		panicOnCallbackError(this.callbacks.OnPositiveInt(this.buffer.DecodeUint()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNegInt8:
		panicOnCallbackError(this.callbacks.OnNegativeInt(uint64(this.buffer.DecodeUint8())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNegInt16:
		panicOnCallbackError(this.callbacks.OnNegativeInt(uint64(this.buffer.DecodeUint16())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNegInt32:
		panicOnCallbackError(this.callbacks.OnNegativeInt(uint64(this.buffer.DecodeUint32())))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNegInt64:
		panicOnCallbackError(this.callbacks.OnNegativeInt(this.buffer.DecodeUint64()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNegInt:
		panicOnCallbackError(this.callbacks.OnNegativeInt(this.buffer.DecodeUint()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeDate:
		panicOnCallbackError(this.callbacks.OnDate(this.buffer.DecodeDate()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeTime:
		panicOnCallbackError(this.callbacks.OnTime(this.buffer.DecodeTime()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeTimestamp:
		panicOnCallbackError(this.callbacks.OnTimestamp(this.buffer.DecodeTimestamp()))
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeNil:
		this.assertNotExpectingMapKey("nil")
		panicOnCallbackError(this.callbacks.OnNil())
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typePadding:
		// Ignore
	case typeList:
		this.assertNotExpectingMapKey("list")
		this.containerBegin(ContainerTypeList)
		panicOnCallbackError(this.callbacks.OnListBegin())
		this.buffer.Commit()
	case typeMapOrdered:
		this.assertNotExpectingMapKey("map")
		this.containerBegin(ContainerTypeOrderedMap)
		panicOnCallbackError(this.callbacks.OnOrderedMapBegin())
		this.buffer.Commit()
	case typeMapUnordered:
		this.assertNotExpectingMapKey("map")
		this.containerBegin(ContainerTypeUnorderedMap)
		panicOnCallbackError(this.callbacks.OnUnorderedMapBegin())
		this.buffer.Commit()
	case typeMapMetadata:
		this.assertNotExpectingMapKey("map")
		this.containerBegin(ContainerTypeMetadataMap)
		panicOnCallbackError(this.callbacks.OnMetadataMapBegin())
		this.buffer.Commit()
	case typeEndContainer:
		this.containerEnd()
		panicOnCallbackError(this.callbacks.OnContainerEnd())
		this.buffer.Commit()
		this.flipMapKeyValueState()
	case typeBytes:
		this.arrayBegin(arrayTypeBytes, this.buffer.DecodeArrayLength())
		this.decodeArrayData()
	case typeComment:
		this.arrayBegin(arrayTypeComment, this.buffer.DecodeArrayLength())
		this.decodeArrayData()
	case typeURI:
		this.arrayBegin(arrayTypeURI, this.buffer.DecodeArrayLength())
		this.decodeArrayData()
	case typeString:
		this.arrayBegin(arrayTypeString, this.buffer.DecodeArrayLength())
		this.decodeArrayData()
	case typeString0:
		this.decodeStringOfLength(0)
	case typeString1:
		this.decodeStringOfLength(1)
	case typeString2:
		this.decodeStringOfLength(2)
	case typeString3:
		this.decodeStringOfLength(3)
	case typeString4:
		this.decodeStringOfLength(4)
	case typeString5:
		this.decodeStringOfLength(5)
	case typeString6:
		this.decodeStringOfLength(6)
	case typeString7:
		this.decodeStringOfLength(7)
	case typeString8:
		this.decodeStringOfLength(8)
	case typeString9:
		this.decodeStringOfLength(9)
	case typeString10:
		this.decodeStringOfLength(10)
	case typeString11:
		this.decodeStringOfLength(11)
	case typeString12:
		this.decodeStringOfLength(12)
	case typeString13:
		this.decodeStringOfLength(13)
	case typeString14:
		this.decodeStringOfLength(14)
	case typeString15:
		this.decodeStringOfLength(15)
	}
	// TODO: 128 bit and decimal
}

func (this *CbeDecoder) beginInlineContainer() {
	if this.container.inlineContainerType != ContainerTypeNone && !this.container.inlineContainerInitialized {
		this.containerBegin(this.container.inlineContainerType)
		switch this.container.inlineContainerType {
		case ContainerTypeList:
			panicOnCallbackError(this.callbacks.OnListBegin())
		case ContainerTypeOrderedMap:
			panicOnCallbackError(this.callbacks.OnOrderedMapBegin())
		case ContainerTypeUnorderedMap:
			panicOnCallbackError(this.callbacks.OnUnorderedMapBegin())
		}
		this.container.inlineContainerInitialized = true
	}
}

func (this *CbeDecoder) endInlineContainer() {
	if this.container.inlineContainerInitialized {
		this.callbacks.OnContainerEnd()
		this.container.depth--
	}
}

func (this *CbeDecoder) assertOnlyOneTopLevelObject() {
	if this.container.depth == 0 && this.firstItemDecoded {
		panic(decoderError{fmt.Errorf("Extra top level object detected")})
	}
}

// ----------
// Public API
// ----------

func NewCbeDecoder(inlineContainerType ContainerType, maxContainerDepth int, callbacks CbeDecoderCallbacks) *CbeDecoder {
	this := new(CbeDecoder)
	this.container.inlineContainerType = inlineContainerType
	if inlineContainerType != ContainerTypeNone {
		maxContainerDepth++
	}
	this.underflowBuffer = NewDecodeBuffer(make([]byte, maxPartialReadLength))
	this.underflowBuffer.Clear()
	this.mainBuffer = NewDecodeBuffer(make([]byte, 0))
	this.buffer = this.mainBuffer
	this.callbacks = callbacks
	this.container.currentType = make([]ContainerType, maxContainerDepth)
	this.container.hasProcessedMapKey = make([]bool, maxContainerDepth)
	return this
}

// Feed bytes into the decoder to be decoded.
func (this *CbeDecoder) Feed(bytesToDecode []byte) (err error) {
	defer func() {
		this.streamOffset += int64(this.buffer.lastCommitPosition)
		if r := recover(); r != nil {
			switch r.(type) {
			case notEnoughBytesToDecodeType:
				// Return as if nothing's wrong
			case notEnoughBytesToDecodeArrayData:
				// Return as if nothing's wrong
			case notEnoughBytesToDecodeObject:
				this.underflowBuffer.AddContents(this.mainBuffer.GetUncommittedBytes())
				this.buffer = this.underflowBuffer
				this.mainBuffer.Clear()
			case callbackError:
				err = fmt.Errorf("cbe: offset %v: Error from callback: %v", this.streamOffset, r.(callbackError).err)
			case decoderError:
				err = fmt.Errorf("cbe: offset %v: Decode error: %v", this.streamOffset, r.(decoderError).err)
			default:
				// Unexpected panics are passed as-is
				panic(r)
			}
		}
	}()

	this.buffer.Rollback()
	this.mainBuffer.ReplaceBuffer(bytesToDecode)

	this.beginInlineContainer()

	if this.buffer == this.underflowBuffer && this.buffer.RemainingByteCount() > 0 {
		underflowByteCount := len(this.underflowBuffer.data)
		bytesFilled := this.buffer.FillFromBuffer(this.mainBuffer, maxPartialReadLength)
		this.mainBuffer.lastCommitPosition += bytesFilled
		this.mainBuffer.position = this.mainBuffer.lastCommitPosition
		objectType := this.buffer.DecodeType()
		this.assertOnlyOneTopLevelObject()
		this.decodeObject(objectType)
		this.firstItemDecoded = true
		mainBytesUsed := this.underflowBuffer.lastCommitPosition - underflowByteCount
		this.mainBuffer.lastCommitPosition = mainBytesUsed
		this.mainBuffer.position = this.mainBuffer.lastCommitPosition

		this.underflowBuffer.Clear()
		this.buffer = this.mainBuffer
	}

	// TODO: Does this handle end of buffer?
	this.decodeArrayData()

	for {
		objectType := this.buffer.DecodeType()
		this.assertOnlyOneTopLevelObject()
		this.decodeObject(objectType)
		this.firstItemDecoded = true
	}

	return err
}

// End the decoding process, doing some final structural tests to make sure it's valid.
func (this *CbeDecoder) End() error {
	this.endInlineContainer()

	if this.container.depth > 0 {
		return fmt.Errorf("Document still has %v open container(s)", this.container.depth)
	}
	if this.array.remainingByteCount > 0 {
		return fmt.Errorf("Array is still open, expecting %d more bytes", this.array.remainingByteCount)
	}
	if this.buffer == this.underflowBuffer {
		return fmt.Errorf("Document has not been completely decoded. %v bytes of underflow data remain", len(this.underflowBuffer.data))
	}
	return nil
}

// Convenience function to decode an entire document in a single call.
func (this *CbeDecoder) Decode(document []byte) error {
	if err := this.Feed(document); err != nil {
		return err
	}
	return this.End()
}
