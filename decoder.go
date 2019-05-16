package cbe

import (
	"fmt"
	"time"
)

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

type CbeDecoderCommentCallbacks interface {
	OnCommentBegin(byteCount uint64) error
	OnCommentData(bytes []byte) error
}

const maxPartialBufferSize = 16

type decoderError struct {
	err error
}

type callbackError struct {
	err error
}

type containerData struct {
	depth              int
	currentType        []containerType
	hasProcessedMapKey []bool
}

type arrayData struct {
	currentType        arrayType
	byteCountRemaining int64
	onBegin            func(uint64) error
	onData             func([]byte) error
}

type fakeDecoderCallbacksStruct struct {
}

var ignoredCommentCallbacks fakeDecoderCallbacksStruct

func (callbacks *fakeDecoderCallbacksStruct) OnCommentBegin(byteCount uint64) error {
	return nil
}

func (callbacks *fakeDecoderCallbacksStruct) OnCommentData(bytes []byte) error {
	return nil
}

type CbeDecoder struct {
	streamOffset     int64
	buffer           decodeBuffer
	partialBuffer    decodeBuffer
	container        containerData
	array            arrayData
	callbacks        CbeDecoderCallbacks
	commentCallbacks CbeDecoderCommentCallbacks
}

func checkCallback(err error) {
	if err != nil {
		panic(callbackError{err})
	}
}

func (decoder *CbeDecoder) isExpectingMapKey() bool {
	return decoder.container.currentType[decoder.container.depth] == containerTypeMap &&
		!decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) isExpectingMapValue() bool {
	return decoder.container.currentType[decoder.container.depth] == containerTypeMap &&
		decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) flipMapKeyStatus() {
	decoder.container.hasProcessedMapKey[decoder.container.depth] = !decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) assertNotExpectingMapKey(keyType string) {
	if decoder.isExpectingMapKey() {
		panic(fmt.Errorf("Cannot use type %v as a map key", keyType))
	}
}

func (decoder *CbeDecoder) containerBegin(newContainerType containerType) {
	if decoder.container.depth >= len(decoder.container.currentType) {
		panic(decoderError{fmt.Errorf("Exceeded max container depth of %v", len(decoder.container.currentType))})
	}
	decoder.container.depth++
	decoder.container.currentType[decoder.container.depth] = newContainerType
	decoder.container.hasProcessedMapKey[decoder.container.depth] = false
}

func (decoder *CbeDecoder) containerEnd() containerType {
	if decoder.container.depth <= 0 {
		panic(decoderError{fmt.Errorf("Got container end but not in a container")})
	}
	if decoder.isExpectingMapValue() {
		panic(decoderError{fmt.Errorf("Expecting map value for already processed key")})
	}

	decoder.container.depth--
	return decoder.container.currentType[decoder.container.depth+1]
}

func (decoder *CbeDecoder) arrayBegin(newArrayType arrayType) {
	decoder.array.currentType = newArrayType
	switch newArrayType {
	case arrayTypeBytes:
		decoder.array.onBegin = decoder.callbacks.OnBytesBegin
		decoder.array.onData = decoder.callbacks.OnBytesData
	case arrayTypeComment:
		decoder.array.onBegin = decoder.commentCallbacks.OnCommentBegin
		decoder.array.onData = decoder.commentCallbacks.OnCommentData
	case arrayTypeString:
		decoder.array.onBegin = decoder.callbacks.OnStringBegin
		decoder.array.onData = decoder.callbacks.OnStringData
	}
}

func (decoder *CbeDecoder) setArrayLength(length int64) {
	decoder.array.byteCountRemaining = length
	if decoder.array.onBegin != nil {
		checkCallback(decoder.array.onBegin(uint64(decoder.array.byteCountRemaining)))
	}
}

func (decoder *CbeDecoder) decodeArrayLength(buffer *decodeBuffer) {
	decoder.setArrayLength(buffer.readArrayLength())
}

func (decoder *CbeDecoder) getArrayDecodeByteCount(buffer *decodeBuffer) int {
	decodeByteCount := decoder.array.byteCountRemaining
	bytesRemaining := len(buffer.data) - buffer.pos
	if int64(bytesRemaining) < decodeByteCount {
		return bytesRemaining
	}
	return int(decodeByteCount)
}

func (decoder *CbeDecoder) decodeArrayData(buffer *decodeBuffer) {
	if decoder.array.currentType == arrayTypeNone {
		return
	}
	if decoder.array.byteCountRemaining == 0 {
		decoder.array.currentType = arrayTypeNone
		decoder.flipMapKeyStatus()
		return
	}

	decodeByteCount := len(buffer.data) - buffer.pos
	if int64(decodeByteCount) > decoder.array.byteCountRemaining {
		decodeByteCount = int(decoder.array.byteCountRemaining)
	}
	bytes := buffer.readPrimitiveBytes(decodeByteCount)
	decoder.array.byteCountRemaining -= int64(decodeByteCount)
	if decoder.array.byteCountRemaining == 0 {
		decoder.array.currentType = arrayTypeNone
		decoder.flipMapKeyStatus()
	}
	if decoder.array.onData != nil {
		checkCallback(decoder.array.onData(bytes))
	}
	if decoder.array.byteCountRemaining > 0 {
		// 0 because we don't want to reserve space in the partial buffer
		panic(failedByteCountReservation(0))
	}
}

func (decoder *CbeDecoder) decodeStringOfLength(buffer *decodeBuffer, length int64) {
	decoder.arrayBegin(arrayTypeString)
	decoder.setArrayLength(length)
	decoder.decodeArrayData(buffer)
}

func (decoder *CbeDecoder) decodeObject(buffer *decodeBuffer, dataType typeField) {
	if int64(int8(dataType)) >= smallIntMin && int64(int8(dataType)) <= smallIntMax {
		if int8(dataType) >= 0 {
			checkCallback(decoder.callbacks.OnUint(uint64(dataType)))
		} else {
			checkCallback(decoder.callbacks.OnInt(int64(int8(dataType))))
		}
		decoder.flipMapKeyStatus()
		return
	}

	switch dataType {
	case typeTrue:
		checkCallback(decoder.callbacks.OnBool(true))
		decoder.flipMapKeyStatus()
	case typeFalse:
		checkCallback(decoder.callbacks.OnBool(false))
		decoder.flipMapKeyStatus()
	case typeFloat32:
		checkCallback(decoder.callbacks.OnFloat(float64(buffer.readFloat32())))
		decoder.flipMapKeyStatus()
	case typeFloat64:
		checkCallback(decoder.callbacks.OnFloat(buffer.readFloat64()))
		decoder.flipMapKeyStatus()
	case typePosInt8:
		checkCallback(decoder.callbacks.OnUint(uint64(buffer.readPrimitive8())))
		decoder.flipMapKeyStatus()
	case typePosInt16:
		checkCallback(decoder.callbacks.OnUint(uint64(buffer.readPrimitive16())))
		decoder.flipMapKeyStatus()
	case typePosInt32:
		checkCallback(decoder.callbacks.OnUint(uint64(buffer.readPrimitive32())))
		decoder.flipMapKeyStatus()
	case typePosInt64:
		checkCallback(decoder.callbacks.OnUint(buffer.readPrimitive64()))
		decoder.flipMapKeyStatus()
	case typeNegInt8:
		checkCallback(decoder.callbacks.OnInt(-int64(buffer.readPrimitive8())))
		decoder.flipMapKeyStatus()
	case typeNegInt16:
		checkCallback(decoder.callbacks.OnInt(-int64(buffer.readPrimitive16())))
		decoder.flipMapKeyStatus()
	case typeNegInt32:
		checkCallback(decoder.callbacks.OnInt(-int64(buffer.readPrimitive32())))
		decoder.flipMapKeyStatus()
	case typeNegInt64:
		checkCallback(decoder.callbacks.OnInt(buffer.readNegInt64()))
		decoder.flipMapKeyStatus()
	case typeSmalltime:
		// TODO: Specify time zone?
		checkCallback(decoder.callbacks.OnTime(buffer.readSmalltime().AsTime()))
		decoder.flipMapKeyStatus()
	case typeNanotime:
		// TODO: Specify time zone?
		checkCallback(decoder.callbacks.OnTime(buffer.readNanotime().AsTime()))
		decoder.flipMapKeyStatus()
	case typeNil:
		decoder.assertNotExpectingMapKey("nil")
		checkCallback(decoder.callbacks.OnNil())
		decoder.flipMapKeyStatus()
	case typePadding:
		// Ignore
	case typeList:
		decoder.containerBegin(containerTypeList)
		checkCallback(decoder.callbacks.OnListBegin())
	case typeMap:
		decoder.containerBegin(containerTypeMap)
		checkCallback(decoder.callbacks.OnMapBegin())
	case typeEndContainer:
		oldContainerType := decoder.containerEnd()
		switch oldContainerType {
		case containerTypeList:
			checkCallback(decoder.callbacks.OnListEnd())
		case containerTypeMap:
			checkCallback(decoder.callbacks.OnMapEnd())
		}
		decoder.flipMapKeyStatus()
	case typeBytes:
		decoder.arrayBegin(arrayTypeBytes)
		decoder.decodeArrayLength(buffer)
		decoder.decodeArrayData(buffer)
	case typeComment:
		decoder.arrayBegin(arrayTypeComment)
		decoder.decodeArrayLength(buffer)
		decoder.decodeArrayData(buffer)
	case typeString:
		decoder.arrayBegin(arrayTypeString)
		decoder.decodeArrayLength(buffer)
		decoder.decodeArrayData(buffer)
	case typeString0:
		decoder.arrayBegin(arrayTypeString)
		decoder.setArrayLength(0)
		decoder.flipMapKeyStatus()
	case typeString1:
		decoder.decodeStringOfLength(buffer, 1)
	case typeString2:
		decoder.decodeStringOfLength(buffer, 2)
	case typeString3:
		decoder.decodeStringOfLength(buffer, 3)
	case typeString4:
		decoder.decodeStringOfLength(buffer, 4)
	case typeString5:
		decoder.decodeStringOfLength(buffer, 5)
	case typeString6:
		decoder.decodeStringOfLength(buffer, 6)
	case typeString7:
		decoder.decodeStringOfLength(buffer, 7)
	case typeString8:
		decoder.decodeStringOfLength(buffer, 8)
	case typeString9:
		decoder.decodeStringOfLength(buffer, 9)
	case typeString10:
		decoder.decodeStringOfLength(buffer, 10)
	case typeString11:
		decoder.decodeStringOfLength(buffer, 11)
	case typeString12:
		decoder.decodeStringOfLength(buffer, 12)
	case typeString13:
		decoder.decodeStringOfLength(buffer, 13)
	case typeString14:
		decoder.decodeStringOfLength(buffer, 14)
	case typeString15:
		decoder.decodeStringOfLength(buffer, 15)
		// TODO: 128 bit and decimal
	}
}

func (decoder *CbeDecoder) feedFromBuffer(buffer *decodeBuffer) error {
	decoder.decodeArrayData(buffer)
	for {
		decoder.decodeObject(buffer, buffer.readType())
	}
	return nil
}

func (decoder *CbeDecoder) handleFailedByteReservation(reservedByteCount int) {
	existingBytes := decoder.buffer.data[decoder.buffer.pos:len(decoder.buffer.data)]
	if decoder.partialBuffer.data == nil {
		decoder.partialBuffer.data = make([]byte, maxPartialBufferSize)
	}
	decoder.partialBuffer.data = decoder.partialBuffer.data[:len(existingBytes)]
	copy(existingBytes, decoder.partialBuffer.data)
	decoder.partialBuffer.bytesToConsume = reservedByteCount
	decoder.partialBuffer.pos = 0
}

// ----------
// Public API
// ----------

func NewCbeDecoder(maxContainerDepth int, callbacks CbeDecoderCallbacks) *CbeDecoder {
	decoder := new(CbeDecoder)
	decoder.callbacks = callbacks
	decoder.container.currentType = make([]containerType, maxContainerDepth)
	decoder.container.hasProcessedMapKey = make([]bool, maxContainerDepth)
	if commentCallbacks, ok := callbacks.(CbeDecoderCommentCallbacks); ok {
		decoder.commentCallbacks = commentCallbacks
	} else {
		decoder.commentCallbacks = &ignoredCommentCallbacks
	}
	return decoder
}

func (decoder *CbeDecoder) Feed(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case endOfData:
				// Nothing to do
			case failedByteCountReservation:
				decoder.handleFailedByteReservation(int(r.(failedByteCountReservation)))
				err = nil
			case callbackError:
				offset := (decoder.streamOffset + int64(decoder.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Error from callback: %v", offset, r.(callbackError).err)
			case decoderError:
				offset := (decoder.streamOffset + int64(decoder.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: Decode error: %v", offset, r.(decoderError).err)
			default:
				panic(r)
				// err = fmt.Errorf("cbe: internal error: %v", r)
			}
		}
		decoder.streamOffset += int64(decoder.buffer.pos)
		decoder.buffer.data = nil
	}()

	// TODO: This needs to only check partial buffer for existing data, and decode 1 payload only.
	// Also, need to fetch remainder of data!
	// decoder.feedFromBuffer(&decoder.partialBuffer, panicHandler)

	decoder.buffer.data = data
	decoder.buffer.pos = 0

	err = decoder.feedFromBuffer(&decoder.buffer)

	return err
}

func (decoder *CbeDecoder) End() error {
	if decoder.container.depth > 0 {
		return fmt.Errorf("Document still has %v open container(s)", decoder.container.depth)
	}
	if decoder.array.byteCountRemaining > 0 {
		return fmt.Errorf("Array is still open, expecting %d more bytes", decoder.array.byteCountRemaining)
	}
	return nil
}

func (decoder *CbeDecoder) Decode(document []byte) error {
	if err := decoder.Feed(document); err != nil {
		return err
	}
	return decoder.End()
}
