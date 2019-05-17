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

func (callbacks *ignoredOptionalCallbacksStruct) OnCommentBegin(byteCount uint64) error {
	return nil
}

func (callbacks *ignoredOptionalCallbacksStruct) OnCommentData(bytes []byte) error {
	return nil
}

var ignoredOptionalCallbacks ignoredOptionalCallbacksStruct

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

type CbeDecoder struct {
	streamOffset     int64
	buffer           decodeBuffer
	partialBuffer    decodeBuffer
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

func (decoder *CbeDecoder) isExpectingMapKey() bool {
	return decoder.container.currentType[decoder.container.depth] == containerTypeMap &&
		!decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) isExpectingMapValue() bool {
	return decoder.container.currentType[decoder.container.depth] == containerTypeMap &&
		decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) flipMapKeyValueState() {
	decoder.container.hasProcessedMapKey[decoder.container.depth] = !decoder.container.hasProcessedMapKey[decoder.container.depth]
}

func (decoder *CbeDecoder) assertNotExpectingMapKey(keyType string) {
	if decoder.isExpectingMapKey() {
		panic(decoderError{fmt.Errorf("Cannot use type %v as a map key", keyType)})
	}
}

func (decoder *CbeDecoder) containerBegin(newContainerType containerType) {
	if decoder.container.depth+1 >= len(decoder.container.currentType) {
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
	decoder.charValidator.Reset()
	decoder.array.byteCountRemaining = length
	if decoder.array.onBegin != nil {
		panicOnCallbackError(decoder.array.onBegin(uint64(decoder.array.byteCountRemaining)))
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
		decoder.flipMapKeyValueState()
		return
	}

	decodeByteCount := len(buffer.data) - buffer.pos
	if int64(decodeByteCount) > decoder.array.byteCountRemaining {
		decodeByteCount = int(decoder.array.byteCountRemaining)
	}
	bytes := buffer.readPrimitiveBytes(decodeByteCount)
	if decoder.array.currentType == arrayTypeString || decoder.array.currentType == arrayTypeComment {
		// if len(bytes) < 20 {
		// 	fmt.Printf("### Validate [%v]\n", string(bytes))
		// }
		for _, ch := range bytes {
			if err := decoder.charValidator.AddByte(int(ch)); err != nil {
				panic(decoderError{err})
			}
			if decoder.charValidator.IsCompleteCharacter() && decoder.array.currentType == arrayTypeComment {
				if err := ValidateCommentCharacter(decoder.charValidator.Character()); err != nil {
					panic(decoderError{err})
				}
			}
		}
	}

	decoder.array.byteCountRemaining -= int64(decodeByteCount)
	if decoder.array.byteCountRemaining == 0 {
		decoder.array.currentType = arrayTypeNone
		decoder.flipMapKeyValueState()
	}
	if decoder.array.onData != nil {
		panicOnCallbackError(decoder.array.onData(bytes))
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
			panicOnCallbackError(decoder.callbacks.OnUint(uint64(dataType)))
		} else {
			panicOnCallbackError(decoder.callbacks.OnInt(int64(int8(dataType))))
		}
		decoder.flipMapKeyValueState()
		return
	}

	switch dataType {
	case typeTrue:
		panicOnCallbackError(decoder.callbacks.OnBool(true))
		decoder.flipMapKeyValueState()
	case typeFalse:
		panicOnCallbackError(decoder.callbacks.OnBool(false))
		decoder.flipMapKeyValueState()
	case typeFloat32:
		panicOnCallbackError(decoder.callbacks.OnFloat(float64(buffer.readFloat32())))
		decoder.flipMapKeyValueState()
	case typeFloat64:
		panicOnCallbackError(decoder.callbacks.OnFloat(buffer.readFloat64()))
		decoder.flipMapKeyValueState()
	case typePosInt8:
		panicOnCallbackError(decoder.callbacks.OnUint(uint64(buffer.readPrimitive8())))
		decoder.flipMapKeyValueState()
	case typePosInt16:
		panicOnCallbackError(decoder.callbacks.OnUint(uint64(buffer.readPrimitive16())))
		decoder.flipMapKeyValueState()
	case typePosInt32:
		panicOnCallbackError(decoder.callbacks.OnUint(uint64(buffer.readPrimitive32())))
		decoder.flipMapKeyValueState()
	case typePosInt64:
		panicOnCallbackError(decoder.callbacks.OnUint(buffer.readPrimitive64()))
		decoder.flipMapKeyValueState()
	case typeNegInt8:
		panicOnCallbackError(decoder.callbacks.OnInt(-int64(buffer.readPrimitive8())))
		decoder.flipMapKeyValueState()
	case typeNegInt16:
		panicOnCallbackError(decoder.callbacks.OnInt(-int64(buffer.readPrimitive16())))
		decoder.flipMapKeyValueState()
	case typeNegInt32:
		panicOnCallbackError(decoder.callbacks.OnInt(-int64(buffer.readPrimitive32())))
		decoder.flipMapKeyValueState()
	case typeNegInt64:
		panicOnCallbackError(decoder.callbacks.OnInt(buffer.readNegInt64()))
		decoder.flipMapKeyValueState()
	case typeSmalltime:
		// TODO: Specify time zone?
		panicOnCallbackError(decoder.callbacks.OnTime(buffer.readSmalltime().AsTime()))
		decoder.flipMapKeyValueState()
	case typeNanotime:
		// TODO: Specify time zone?
		panicOnCallbackError(decoder.callbacks.OnTime(buffer.readNanotime().AsTime()))
		decoder.flipMapKeyValueState()
	case typeNil:
		decoder.assertNotExpectingMapKey("nil")
		panicOnCallbackError(decoder.callbacks.OnNil())
		decoder.flipMapKeyValueState()
	case typePadding:
		// Ignore
	case typeList:
		decoder.assertNotExpectingMapKey("list")
		decoder.containerBegin(containerTypeList)
		panicOnCallbackError(decoder.callbacks.OnListBegin())
	case typeMap:
		decoder.assertNotExpectingMapKey("map")
		decoder.containerBegin(containerTypeMap)
		panicOnCallbackError(decoder.callbacks.OnMapBegin())
	case typeEndContainer:
		oldContainerType := decoder.containerEnd()
		switch oldContainerType {
		case containerTypeList:
			panicOnCallbackError(decoder.callbacks.OnListEnd())
		case containerTypeMap:
			panicOnCallbackError(decoder.callbacks.OnMapEnd())
		}
		decoder.flipMapKeyValueState()
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
		decoder.flipMapKeyValueState()
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
		objectType := buffer.readType()
		if decoder.container.depth == 0 && decoder.firstItemDecoded {
			panic(decoderError{fmt.Errorf("Extra top level object detected")})
		}
		decoder.decodeObject(buffer, objectType)
		decoder.firstItemDecoded = true
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
	if commentCallbacks, ok := callbacks.(CbeDecoderOptionalCallbacks); ok {
		decoder.commentCallbacks = commentCallbacks
	} else {
		decoder.commentCallbacks = &ignoredOptionalCallbacks
	}
	return decoder
}

// Feed bytes into the decoder to be decoded.
func (decoder *CbeDecoder) Feed(bytesToDecode []byte) (err error) {
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
				// Unexpected panics are passed as-is
				panic(r)
			}
		}
		decoder.streamOffset += int64(decoder.buffer.pos)
		decoder.buffer.data = nil
	}()

	// TODO: This needs to only check partial buffer for existing data, and decode 1 payload only.
	// Also, need to fetch remainder of data!
	// decoder.feedFromBuffer(&decoder.partialBuffer, panicHandler)

	decoder.buffer.data = bytesToDecode
	decoder.buffer.pos = 0

	err = decoder.feedFromBuffer(&decoder.buffer)

	return err
}

// End the decoding process, doing some final structural tests to make sure it's valid.
func (decoder *CbeDecoder) End() error {
	if decoder.container.depth > 0 {
		return fmt.Errorf("Document still has %v open container(s)", decoder.container.depth)
	}
	if decoder.array.byteCountRemaining > 0 {
		return fmt.Errorf("Array is still open, expecting %d more bytes", decoder.array.byteCountRemaining)
	}
	return nil
}

// Convenience function to decode an entire document in a single call.
func (decoder *CbeDecoder) Decode(document []byte) error {
	if err := decoder.Feed(document); err != nil {
		return err
	}
	return decoder.End()
}
