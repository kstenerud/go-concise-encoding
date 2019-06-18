package cbe

import (
	"fmt"
	"time"
)

// Callback functions that must be present in the receiver object.
type CbeDecoderCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnInt(value int64) error
	OnUint(value uint64) error
	OnFloat(value float64) error
	OnTime(value time.Time) error
	OnListBegin() error
	OnListEnd() error
	OnMapBegin() error
	OnMapEnd() error
	OnStringBegin(byteCount uint64) error
	OnStringData(bytes []byte) error
	OnBytesBegin(byteCount uint64) error
	OnBytesData(bytes []byte) error
}

// Callback functions that will be used if present in the receiver object.
type CbeDecoderOptionalCallbacks interface {
	OnCommentBegin(byteCount uint64) error
	OnCommentData(bytes []byte) error
}

type ignoredOptionalCallbacksStruct struct {
}

func (this *ignoredOptionalCallbacksStruct) OnCommentBegin(byteCount uint64) error {
	return nil
}

func (this *ignoredOptionalCallbacksStruct) OnCommentData(bytes []byte) error {
	return nil
}

var ignoredOptionalCallbacks ignoredOptionalCallbacksStruct

const maxPartialBufferSize = 17

type decoderError struct {
	err error
}

type callbackError struct {
	err error
}

type containerData struct {
	inlineContainerType           ContainerType
	hasInitializedInlineContainer bool
	depth                         int
	currentType                   []ContainerType
	hasProcessedMapKey            []bool
}

type arrayData struct {
	currentType        arrayType
	byteCountRemaining int64
	onBegin            func(uint64) error
	onData             func([]byte) error
}

type CbeDecoder struct {
	streamOffset     int64
	buffer           decodeBuffer
	underflowBuffer  decodeBuffer
	container        containerData
	array            arrayData
	callbacks        CbeDecoderCallbacks
	commentCallbacks CbeDecoderOptionalCallbacks
	firstItemDecoded bool
	charValidator    Utf8Validator
}

func panicOnCallbackError(err error) {
	if err != nil {
		panic(callbackError{err})
	}
}

func (this *CbeDecoder) isExpectingMapKey() bool {
	return this.container.currentType[this.container.depth] == ContainerTypeMap &&
		!this.container.hasProcessedMapKey[this.container.depth]
}

func (this *CbeDecoder) isExpectingMapValue() bool {
	return this.container.currentType[this.container.depth] == ContainerTypeMap &&
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

func (this *CbeDecoder) arrayBegin(newArrayType arrayType) {
	this.array.currentType = newArrayType
	switch newArrayType {
	case arrayTypeBytes:
		this.array.onBegin = this.callbacks.OnBytesBegin
		this.array.onData = this.callbacks.OnBytesData
	case arrayTypeComment:
		this.array.onBegin = this.commentCallbacks.OnCommentBegin
		this.array.onData = this.commentCallbacks.OnCommentData
	case arrayTypeString:
		this.array.onBegin = this.callbacks.OnStringBegin
		this.array.onData = this.callbacks.OnStringData
	}
}

func (this *CbeDecoder) setArrayLength(length int64) {
	this.charValidator.Reset()
	this.array.byteCountRemaining = length
	if this.array.onBegin != nil {
		panicOnCallbackError(this.array.onBegin(uint64(this.array.byteCountRemaining)))
	}
}

func (this *CbeDecoder) decodeArrayLength(buffer *decodeBuffer) {
	this.setArrayLength(buffer.readArrayLength())
}

func (this *CbeDecoder) getArrayDecodeByteCount(buffer *decodeBuffer) int {
	decodeByteCount := this.array.byteCountRemaining
	bytesRemaining := len(buffer.data) - buffer.pos
	if int64(bytesRemaining) < decodeByteCount {
		return bytesRemaining
	}
	return int(decodeByteCount)
}

func (this *CbeDecoder) decodeArrayData(buffer *decodeBuffer) {
	if this.array.currentType == arrayTypeNone {
		return
	}
	if this.array.byteCountRemaining == 0 {
		this.array.currentType = arrayTypeNone
		this.flipMapKeyValueState()
		return
	}

	decodeByteCount := len(buffer.data) - buffer.pos
	if int64(decodeByteCount) > this.array.byteCountRemaining {
		decodeByteCount = int(this.array.byteCountRemaining)
	}
	bytes := buffer.readPrimitiveBytes(decodeByteCount)
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

	this.array.byteCountRemaining -= int64(decodeByteCount)
	if this.array.byteCountRemaining == 0 {
		this.array.currentType = arrayTypeNone
		this.flipMapKeyValueState()
	}
	if this.array.onData != nil {
		panicOnCallbackError(this.array.onData(bytes))
	}
	if this.array.byteCountRemaining > 0 {
		// 0 because we don't want to reserve space in the partial buffer
		panic(failedByteCountReservation(0))
	}
}

func (this *CbeDecoder) decodeStringOfLength(buffer *decodeBuffer, length int64) {
	this.arrayBegin(arrayTypeString)
	this.setArrayLength(length)
	this.decodeArrayData(buffer)
}

func (this *CbeDecoder) decodeObject(buffer *decodeBuffer, dataType typeField) {
	if int64(int8(dataType)) >= smallIntMin && int64(int8(dataType)) <= smallIntMax {
		if int8(dataType) >= 0 {
			panicOnCallbackError(this.callbacks.OnUint(uint64(dataType)))
		} else {
			panicOnCallbackError(this.callbacks.OnInt(int64(int8(dataType))))
		}
		this.flipMapKeyValueState()
		return
	}

	switch dataType {
	case typeTrue:
		panicOnCallbackError(this.callbacks.OnBool(true))
		this.flipMapKeyValueState()
	case typeFalse:
		panicOnCallbackError(this.callbacks.OnBool(false))
		this.flipMapKeyValueState()
	case typeFloat32:
		panicOnCallbackError(this.callbacks.OnFloat(float64(buffer.readFloat32())))
		this.flipMapKeyValueState()
	case typeFloat64:
		panicOnCallbackError(this.callbacks.OnFloat(buffer.readFloat64()))
		this.flipMapKeyValueState()
	case typePosInt8:
		panicOnCallbackError(this.callbacks.OnUint(uint64(buffer.readPrimitive8())))
		this.flipMapKeyValueState()
	case typePosInt16:
		panicOnCallbackError(this.callbacks.OnUint(uint64(buffer.readPrimitive16())))
		this.flipMapKeyValueState()
	case typePosInt32:
		panicOnCallbackError(this.callbacks.OnUint(uint64(buffer.readPrimitive32())))
		this.flipMapKeyValueState()
	case typePosInt64:
		panicOnCallbackError(this.callbacks.OnUint(buffer.readPrimitive64()))
		this.flipMapKeyValueState()
	case typeNegInt8:
		panicOnCallbackError(this.callbacks.OnInt(-int64(buffer.readPrimitive8())))
		this.flipMapKeyValueState()
	case typeNegInt16:
		panicOnCallbackError(this.callbacks.OnInt(-int64(buffer.readPrimitive16())))
		this.flipMapKeyValueState()
	case typeNegInt32:
		panicOnCallbackError(this.callbacks.OnInt(-int64(buffer.readPrimitive32())))
		this.flipMapKeyValueState()
	case typeNegInt64:
		panicOnCallbackError(this.callbacks.OnInt(buffer.readNegInt64()))
		this.flipMapKeyValueState()
	case typeSmalltime:
		// TODO: Specify time zone?
		panicOnCallbackError(this.callbacks.OnTime(buffer.readSmalltime().AsTime()))
		this.flipMapKeyValueState()
	case typeNanotime:
		// TODO: Specify time zone?
		panicOnCallbackError(this.callbacks.OnTime(buffer.readNanotime().AsTime()))
		this.flipMapKeyValueState()
	case typeNil:
		this.assertNotExpectingMapKey("nil")
		panicOnCallbackError(this.callbacks.OnNil())
		this.flipMapKeyValueState()
	case typePadding:
		// Ignore
	case typeList:
		this.assertNotExpectingMapKey("list")
		this.containerBegin(ContainerTypeList)
		panicOnCallbackError(this.callbacks.OnListBegin())
	case typeMap:
		this.assertNotExpectingMapKey("map")
		this.containerBegin(ContainerTypeMap)
		panicOnCallbackError(this.callbacks.OnMapBegin())
	case typeEndContainer:
		oldContainerType := this.containerEnd()
		switch oldContainerType {
		case ContainerTypeList:
			panicOnCallbackError(this.callbacks.OnListEnd())
		case ContainerTypeMap:
			panicOnCallbackError(this.callbacks.OnMapEnd())
		}
		this.flipMapKeyValueState()
	case typeBytes:
		this.arrayBegin(arrayTypeBytes)
		this.decodeArrayLength(buffer)
		this.decodeArrayData(buffer)
	case typeComment:
		this.arrayBegin(arrayTypeComment)
		this.decodeArrayLength(buffer)
		this.decodeArrayData(buffer)
	case typeString:
		this.arrayBegin(arrayTypeString)
		this.decodeArrayLength(buffer)
		this.decodeArrayData(buffer)
	case typeString0:
		this.arrayBegin(arrayTypeString)
		this.setArrayLength(0)
		this.flipMapKeyValueState()
	case typeString1:
		this.decodeStringOfLength(buffer, 1)
	case typeString2:
		this.decodeStringOfLength(buffer, 2)
	case typeString3:
		this.decodeStringOfLength(buffer, 3)
	case typeString4:
		this.decodeStringOfLength(buffer, 4)
	case typeString5:
		this.decodeStringOfLength(buffer, 5)
	case typeString6:
		this.decodeStringOfLength(buffer, 6)
	case typeString7:
		this.decodeStringOfLength(buffer, 7)
	case typeString8:
		this.decodeStringOfLength(buffer, 8)
	case typeString9:
		this.decodeStringOfLength(buffer, 9)
	case typeString10:
		this.decodeStringOfLength(buffer, 10)
	case typeString11:
		this.decodeStringOfLength(buffer, 11)
	case typeString12:
		this.decodeStringOfLength(buffer, 12)
	case typeString13:
		this.decodeStringOfLength(buffer, 13)
	case typeString14:
		this.decodeStringOfLength(buffer, 14)
	case typeString15:
		this.decodeStringOfLength(buffer, 15)
		// TODO: 128 bit and decimal
	}
}

func (this *CbeDecoder) handleFailedByteReservation(objectType typeField, reservedByteCount int) {
	existingBytes := this.buffer.data[this.buffer.pos:len(this.buffer.data)]
	if this.underflowBuffer.data == nil {
		this.underflowBuffer.data = make([]byte, maxPartialBufferSize)
	}
	this.underflowBuffer.data = this.underflowBuffer.data[:len(existingBytes)+1]
	this.underflowBuffer.data[0] = byte(objectType)
	dst := this.underflowBuffer.data[1:len(this.underflowBuffer.data)]
	copy(dst, existingBytes)
	this.underflowBuffer.bytesToConsume = reservedByteCount - len(existingBytes)
	this.underflowBuffer.pos = 0
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
	this.callbacks = callbacks
	this.container.currentType = make([]ContainerType, maxContainerDepth)
	this.container.hasProcessedMapKey = make([]bool, maxContainerDepth)
	if commentCallbacks, ok := callbacks.(CbeDecoderOptionalCallbacks); ok {
		this.commentCallbacks = commentCallbacks
	} else {
		this.commentCallbacks = &ignoredOptionalCallbacks
	}
	return this
}

func (this *CbeDecoder) feedFromUnderflowBuffer() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case endOfData:
				// Nothing to do
			case failedByteCountReservation:
				// This would be a bug
				panic(r)
			case callbackError:
				offset := (this.streamOffset + int64(this.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Error from callback: %v", offset, r.(callbackError).err)
			case decoderError:
				offset := (this.streamOffset + int64(this.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Decode error: %v", offset, r.(decoderError).err)
			default:
				// Unexpected panics are passed as-is
				panic(r)
			}
		}
		this.streamOffset += int64(this.buffer.pos)
	}()

	objectType := this.underflowBuffer.readType()
	if this.container.depth == 0 && this.firstItemDecoded {
		panic(decoderError{fmt.Errorf("Extra top level object detected")})
	}
	this.decodeObject(&this.underflowBuffer, objectType)
	this.firstItemDecoded = true

	this.underflowBuffer.bytesToConsume = 0
	this.underflowBuffer.pos = 0
	// TODO: Offset gets screwed up here maybe?

	return err
}

func (this *CbeDecoder) feedFromMainBuffer() (err error) {
	var objectType typeField
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case endOfData:
				// Nothing to do
			case failedByteCountReservation:
				this.handleFailedByteReservation(objectType, int(r.(failedByteCountReservation)))
				err = nil
			case callbackError:
				offset := (this.streamOffset + int64(this.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Error from callback: %v", offset, r.(callbackError).err)
			case decoderError:
				offset := (this.streamOffset + int64(this.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Decode error: %v", offset, r.(decoderError).err)
			default:
				// Unexpected panics are passed as-is
				panic(r)
			}
		}
		this.streamOffset += int64(this.buffer.pos)
		this.buffer.data = nil
	}()

	this.decodeArrayData(&this.buffer)

	for {
		objectType = this.buffer.readType()
		if this.container.depth == 0 && this.firstItemDecoded {
			panic(decoderError{fmt.Errorf("Extra top level object detected")})
		}
		this.decodeObject(&this.buffer, objectType)
		this.firstItemDecoded = true
	}

	return err
}

// Feed bytes into the decoder to be decoded.
func (this *CbeDecoder) Feed(bytesToDecode []byte) error {
	if this.container.inlineContainerType != ContainerTypeNone && !this.container.hasInitializedInlineContainer {
		this.containerBegin(this.container.inlineContainerType)
		switch this.container.inlineContainerType {
		case ContainerTypeList:
			panicOnCallbackError(this.callbacks.OnListBegin())
		case ContainerTypeMap:
			panicOnCallbackError(this.callbacks.OnMapBegin())
		}
		this.container.hasInitializedInlineContainer = true
	}
	if bytesToCopy := this.underflowBuffer.bytesToConsume; bytesToCopy > 0 {
		if bytesAvailable := len(bytesToDecode); bytesAvailable < bytesToCopy {
			bytesToCopy = bytesAvailable
		}

		this.underflowBuffer.data = append(this.underflowBuffer.data, bytesToDecode[0:bytesToCopy]...)
		this.underflowBuffer.bytesToConsume -= bytesToCopy

		if this.underflowBuffer.bytesToConsume > 0 {
			return nil
		}

		if err := this.feedFromUnderflowBuffer(); err != nil {
			return err
		}

		bytesToDecode = bytesToDecode[bytesToCopy:len(bytesToDecode)]
	}

	this.buffer.data = bytesToDecode
	this.buffer.pos = 0

	return this.feedFromMainBuffer()
}

// End the decoding process, doing some final structural tests to make sure it's valid.
func (this *CbeDecoder) End() error {
	if this.container.hasInitializedInlineContainer {
		switch this.container.inlineContainerType {
		case ContainerTypeList:
			this.callbacks.OnListEnd()
			this.container.depth--
		case ContainerTypeMap:
			this.callbacks.OnMapEnd()
			this.container.depth--
		}
	}
	if this.container.depth > 0 {
		return fmt.Errorf("Document still has %v open container(s)", this.container.depth)
	}
	if this.array.byteCountRemaining > 0 {
		return fmt.Errorf("Array is still open, expecting %d more bytes", this.array.byteCountRemaining)
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
