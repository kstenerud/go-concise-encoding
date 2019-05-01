package cbe

import (
	"fmt"
	"math"
	"time"

	"github.com/kstenerud/go-smalltime"
)

type decoderError error
type callbackError error
type endOfBufferExit bool

type DecoderCallbacks struct {
	OnNil          func() error
	OnBool         func(value bool) error
	OnInt          func(value int) error
	OnInt64        func(value int64) error
	OnUint         func(value uint) error
	OnUint64       func(value uint64) error
	OnFloat32      func(value float32) error
	OnFloat64      func(value float64) error
	OnTime         func(value time.Time) error
	OnListBegin    func() error
	OnListEnd      func() error
	OnMapBegin     func() error
	OnMapEnd       func() error
	OnStringBegin  func(byteCount uint64) error
	OnStringData   func(bytes []byte) error
	OnCommentBegin func(byteCount uint64) error
	OnCommentData  func(bytes []byte) error
	OnBinaryBegin  func(byteCount uint64) error
	OnBinaryData   func(bytes []byte) error
}

type Decoder struct {
	buffer               []byte
	bufferPos            int
	streamOffset         uint64
	bytesToConsume       uint64
	containerDepth       int
	currentArrayType     arrayType
	currentContainerType []containerType
	arrayBytesRemaining  uint64
	arrayDecodeCallback  func([]byte) error
	arrayLengthCallback  func(uint64) error
	callbacks            *DecoderCallbacks
}

// TODO: Maybe allow these to be natural sized ints?
func (decoder *Decoder) reserveBytes(byteCount uint64) {
	if uint64(decoder.bufferPos)+byteCount > uint64(len(decoder.buffer)) {
		// TODO: panic
	}
	decoder.bytesToConsume = byteCount
}

func (decoder *Decoder) consumeReservedBytes() {
	decoder.bufferPos += int(decoder.bytesToConsume)
}

func (decoder *Decoder) readPrimitive8() uint {
	decoder.reserveBytes(1)
	value := uint(decoder.buffer[decoder.bufferPos])
	decoder.consumeReservedBytes()
	return value
}

func (decoder *Decoder) readPrimitive16() uint {
	decoder.reserveBytes(2)
	value := uint(decoder.buffer[decoder.bufferPos]) |
		uint(decoder.buffer[decoder.bufferPos+1])<<8
	decoder.consumeReservedBytes()
	return value
}

func (decoder *Decoder) readPrimitive32() uint {
	decoder.reserveBytes(4)
	value := uint(decoder.buffer[decoder.bufferPos]) |
		uint(decoder.buffer[decoder.bufferPos+1])<<8 |
		uint(decoder.buffer[decoder.bufferPos+2])<<16 |
		uint(decoder.buffer[decoder.bufferPos+3])<<24
	decoder.consumeReservedBytes()
	return value
}

func (decoder *Decoder) readPrimitive64() uint64 {
	decoder.reserveBytes(8)
	value := uint64(decoder.buffer[decoder.bufferPos]) |
		uint64(decoder.buffer[decoder.bufferPos+1])<<8 |
		uint64(decoder.buffer[decoder.bufferPos+2])<<16 |
		uint64(decoder.buffer[decoder.bufferPos+3])<<24 |
		uint64(decoder.buffer[decoder.bufferPos+4])<<32 |
		uint64(decoder.buffer[decoder.bufferPos+5])<<40 |
		uint64(decoder.buffer[decoder.bufferPos+6])<<48 |
		uint64(decoder.buffer[decoder.bufferPos+7])<<56
	decoder.consumeReservedBytes()
	return value
}

func (decoder *Decoder) readPrimitiveBytes(byteCount uint64) []byte {
	decoder.reserveBytes(byteCount)
	bytes := decoder.buffer[decoder.bufferPos : decoder.bufferPos+int(byteCount)]
	decoder.consumeReservedBytes()
	return bytes
}

func (decoder *Decoder) readInt8() int8 {
	return int8(decoder.readPrimitive8())
}

func (decoder *Decoder) readInt16() int16 {
	return int16(decoder.readPrimitive16())
}

func (decoder *Decoder) readInt32() int32 {
	return int32(decoder.readPrimitive32())
}

func (decoder *Decoder) readInt64() int64 {
	return int64(decoder.readPrimitive64())
}

func (decoder *Decoder) readFloat32() float32 {
	return math.Float32frombits(uint32(decoder.readPrimitive32()))
}

func (decoder *Decoder) readFloat64() float64 {
	return math.Float64frombits(decoder.readPrimitive64())
}

func (decoder *Decoder) readType() typeField {
	return typeField(decoder.readPrimitive8())
}

func (decoder *Decoder) readTime() smalltime.Smalltime {
	return smalltime.Smalltime(decoder.readPrimitive64())
}

func (decoder *Decoder) readArrayLength() uint64 {
	firstByte := decoder.readPrimitive8()
	switch uint64(firstByte & 3) {
	case length6Bit:
		return uint64(firstByte >> 2)
	case length14Bit:
		return uint64(firstByte>>2) |
			uint64(decoder.readPrimitive8())<<6
	case length30Bit:
		return uint64(firstByte>>2) |
			uint64(decoder.readPrimitive8())<<6 |
			uint64(decoder.readPrimitive8())<<14 |
			uint64(decoder.readPrimitive8())<<22
	case length62Bit:
		return uint64(firstByte>>2) |
			uint64(decoder.readPrimitive8())<<6 |
			uint64(decoder.readPrimitive8())<<14 |
			uint64(decoder.readPrimitive8())<<22 |
			uint64(decoder.readPrimitive8())<<30 |
			uint64(decoder.readPrimitive8())<<38 |
			uint64(decoder.readPrimitive8())<<46 |
			uint64(decoder.readPrimitive8())<<54
	default: // TODO: 128 bit
		return 0
	}
}

func (decoder *Decoder) decodeNegInt64() int64 {
	value := decoder.readPrimitive64()
	if value&0x8000000000000000 != 0 {
		// TODO: Error
		return 0
	} else {
		return -int64(value)
	}
}

func (decoder *Decoder) enterContainer(newContainerType containerType) {
	// TODO: Error if container depth >= max
	// TODO: Error if in array
	decoder.containerDepth++
	decoder.currentContainerType[decoder.containerDepth] = newContainerType
}

func (decoder *Decoder) leaveContainer() {
	// TODO: Error if not in container
	decoder.containerDepth--
}

func (decoder *Decoder) getCurrentContainerType() containerType {
	// TODO: Error if not in container
	return decoder.currentContainerType[decoder.containerDepth]
}

func (decoder *Decoder) beginArray(newArrayType arrayType) {
	// TODO: Error if already in array
	decoder.currentArrayType = newArrayType
	switch newArrayType {
	case arrayTypeBinary:
		decoder.arrayLengthCallback = decoder.callbacks.OnBinaryBegin
		decoder.arrayDecodeCallback = decoder.callbacks.OnBinaryData
	case arrayTypeComment:
		decoder.arrayLengthCallback = decoder.callbacks.OnCommentBegin
		decoder.arrayDecodeCallback = decoder.callbacks.OnCommentData
	case arrayTypeString:
		decoder.arrayLengthCallback = decoder.callbacks.OnStringBegin
		decoder.arrayDecodeCallback = decoder.callbacks.OnStringData
	}
}

func checkCallback(err error) {
	if err != nil {
		panic(callbackError(err))
	}
}

func (decoder *Decoder) setArrayLength(length uint64) {
	// TODO: Error if already in array
	decoder.arrayBytesRemaining = length
	checkCallback(decoder.arrayLengthCallback(decoder.arrayBytesRemaining))
}

func (decoder *Decoder) decodeArrayLength() {
	decoder.setArrayLength(decoder.readArrayLength())
}

func (decoder *Decoder) decodeArrayData() {
	// TODO: Make this more intelligent
	// Needs to decide how much data is available and only grab that,
	// leaving the rest for the next feed() call.
	bytes := decoder.readPrimitiveBytes(decoder.arrayBytesRemaining)
	checkCallback(decoder.arrayDecodeCallback(bytes))
}

func (decoder *Decoder) decodeStringOfLength(length uint64) {
	decoder.beginArray(arrayTypeString)
	decoder.setArrayLength(length)
	decoder.decodeArrayData()
}

func (decoder *Decoder) Feed(data []byte) (err error) {
	defer func() {
		decoder.streamOffset += uint64(decoder.bufferPos)
		decoder.buffer = nil
		if r := recover(); r != nil {
			switch r.(type) {
			case endOfBufferExit:
				err = nil
			case callbackError:
				offset := (decoder.streamOffset + uint64(decoder.bufferPos))
				err = fmt.Errorf("cbe: offset %v: error from callback: %v", offset, r)
			case decoderError:
				offset := (decoder.streamOffset + uint64(decoder.bufferPos))
				err = fmt.Errorf("cbe: offset %v: %v", offset, r)
			default:
				err = fmt.Errorf("cbe: internal error: %v", r)
			}
		}
	}()

	decoder.buffer = data
	decoder.bufferPos = 0

	if decoder.currentArrayType != arrayTypeNone {
		decoder.decodeArrayData()
	}

	for {
		dataType := decoder.readType()
		if int64(int8(dataType)) >= smallIntMin && int64(int8(dataType)) <= smallIntMax {
			checkCallback(decoder.callbacks.OnInt(int(int8(dataType))))
			continue
		}
		switch dataType {
		case typeTrue:
			checkCallback(decoder.callbacks.OnBool(true))
		case typeFalse:
			checkCallback(decoder.callbacks.OnBool(false))
		case typeFloat32:
			checkCallback(decoder.callbacks.OnFloat32(decoder.readFloat32()))
		case typeFloat64:
			checkCallback(decoder.callbacks.OnFloat64(decoder.readFloat64()))
		case typePosInt8:
			checkCallback(decoder.callbacks.OnUint(decoder.readPrimitive8()))
		case typePosInt16:
			checkCallback(decoder.callbacks.OnUint(decoder.readPrimitive16()))
		case typePosInt32:
			checkCallback(decoder.callbacks.OnUint(decoder.readPrimitive32()))
		case typePosInt64:
			checkCallback(decoder.callbacks.OnUint64(decoder.readPrimitive64()))
		case typeNegInt8:
			checkCallback(decoder.callbacks.OnInt(-int(decoder.readPrimitive8())))
		case typeNegInt16:
			checkCallback(decoder.callbacks.OnInt(-int(decoder.readPrimitive16())))
		case typeNegInt32:
			value := -int64(decoder.readPrimitive32())
			if value < math.MinInt32 {
				checkCallback(decoder.callbacks.OnInt64(value))
			} else {
				checkCallback(decoder.callbacks.OnInt(int(value)))
			}
		case typeNegInt64:
			checkCallback(decoder.callbacks.OnInt64(decoder.decodeNegInt64()))
		case typeTime:
			// TODO: Specify time zone?
			checkCallback(decoder.callbacks.OnTime(decoder.readTime().AsTime()))
		case typeNil:
			checkCallback(decoder.callbacks.OnNil())
		case typePadding:
			// Ignore
		case typeList:
			decoder.enterContainer(containerTypeList)
			checkCallback(decoder.callbacks.OnListBegin())
		case typeMap:
			decoder.enterContainer(containerTypeMap)
			checkCallback(decoder.callbacks.OnMapBegin())
		case typeEndContainer:
			oldContainerType := decoder.getCurrentContainerType()
			decoder.leaveContainer()
			switch oldContainerType {
			case containerTypeList:
				checkCallback(decoder.callbacks.OnListEnd())
			case containerTypeMap:
				checkCallback(decoder.callbacks.OnMapEnd())
			}
		case typeBinary:
			decoder.beginArray(arrayTypeBinary)
			decoder.decodeArrayLength()
			decoder.decodeArrayData()
		case typeComment:
			decoder.beginArray(arrayTypeComment)
			decoder.decodeArrayLength()
			decoder.decodeArrayData()
		case typeString:
			decoder.beginArray(arrayTypeString)
			decoder.decodeArrayLength()
			decoder.decodeArrayData()
		case typeString0:
			decoder.beginArray(arrayTypeString)
			decoder.setArrayLength(0)
		case typeString1:
			decoder.decodeStringOfLength(1)
		case typeString2:
			decoder.decodeStringOfLength(2)
		case typeString3:
			decoder.decodeStringOfLength(3)
		case typeString4:
			decoder.decodeStringOfLength(4)
		case typeString5:
			decoder.decodeStringOfLength(5)
		case typeString6:
			decoder.decodeStringOfLength(6)
		case typeString7:
			decoder.decodeStringOfLength(7)
		case typeString8:
			decoder.decodeStringOfLength(8)
		case typeString9:
			decoder.decodeStringOfLength(9)
		case typeString10:
			decoder.decodeStringOfLength(10)
		case typeString11:
			decoder.decodeStringOfLength(11)
		case typeString12:
			decoder.decodeStringOfLength(12)
		case typeString13:
			decoder.decodeStringOfLength(13)
		case typeString14:
			decoder.decodeStringOfLength(14)
		case typeString15:
			decoder.decodeStringOfLength(15)
		}
		// TODO: 128 bit and decimal
	}
	return nil
}

func (decoder *Decoder) End() error {
	// TODO
	return nil
}

func (decoder *Decoder) Decode(document []byte) error {
	if err := decoder.Feed(document); err != nil {
		return err
	}
	return decoder.End()
}
