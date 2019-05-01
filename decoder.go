package cbe

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kstenerud/go-smalltime"
)

type decoderError error
type callbackError error
type endOfBufferExit error

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

type decodeBuffer struct {
	data           []byte
	pos            int
	bytesToConsume int
}

// TODO: Maybe allow these to be natural sized ints?
func (buffer *decodeBuffer) reserveBytes(byteCount int) {
	if buffer.pos+byteCount > len(buffer.data) {
		panic(endOfBufferExit(errors.New("")))
	}
	buffer.bytesToConsume = byteCount
}

func (buffer *decodeBuffer) consumeReservedBytes() {
	buffer.pos += buffer.bytesToConsume
}

func (buffer *decodeBuffer) readPrimitive8() uint {
	buffer.reserveBytes(1)
	value := uint(buffer.data[buffer.pos])
	buffer.consumeReservedBytes()
	return value
}

func (buffer *decodeBuffer) readPrimitive16() uint {
	buffer.reserveBytes(2)
	value := uint(buffer.data[buffer.pos]) |
		uint(buffer.data[buffer.pos+1])<<8
	buffer.consumeReservedBytes()
	return value
}

func (buffer *decodeBuffer) readPrimitive32() uint {
	buffer.reserveBytes(4)
	value := uint(buffer.data[buffer.pos]) |
		uint(buffer.data[buffer.pos+1])<<8 |
		uint(buffer.data[buffer.pos+2])<<16 |
		uint(buffer.data[buffer.pos+3])<<24
	buffer.consumeReservedBytes()
	return value
}

func (buffer *decodeBuffer) readPrimitive64() uint64 {
	buffer.reserveBytes(8)
	value := uint64(buffer.data[buffer.pos]) |
		uint64(buffer.data[buffer.pos+1])<<8 |
		uint64(buffer.data[buffer.pos+2])<<16 |
		uint64(buffer.data[buffer.pos+3])<<24 |
		uint64(buffer.data[buffer.pos+4])<<32 |
		uint64(buffer.data[buffer.pos+5])<<40 |
		uint64(buffer.data[buffer.pos+6])<<48 |
		uint64(buffer.data[buffer.pos+7])<<56
	buffer.consumeReservedBytes()
	return value
}

func (buffer *decodeBuffer) readPrimitiveBytes(byteCount int) []byte {
	buffer.reserveBytes(byteCount)
	bytes := buffer.data[buffer.pos : buffer.pos+byteCount]
	buffer.consumeReservedBytes()
	return bytes
}

func (buffer *decodeBuffer) readInt8() int8 {
	return int8(buffer.readPrimitive8())
}

func (buffer *decodeBuffer) readInt16() int16 {
	return int16(buffer.readPrimitive16())
}

func (buffer *decodeBuffer) readInt32() int32 {
	return int32(buffer.readPrimitive32())
}

func (buffer *decodeBuffer) readInt64() int64 {
	return int64(buffer.readPrimitive64())
}

func (buffer *decodeBuffer) readFloat32() float32 {
	return math.Float32frombits(uint32(buffer.readPrimitive32()))
}

func (buffer *decodeBuffer) readFloat64() float64 {
	return math.Float64frombits(buffer.readPrimitive64())
}

func (buffer *decodeBuffer) readType() typeField {
	return typeField(buffer.readPrimitive8())
}

func (buffer *decodeBuffer) readTime() smalltime.Smalltime {
	return smalltime.Smalltime(buffer.readPrimitive64())
}

func (buffer *decodeBuffer) readArrayLength() int64 {
	firstByte := buffer.readPrimitive8()
	switch int64(firstByte & 3) {
	case length6Bit:
		return int64(firstByte >> 2)
	case length14Bit:
		return int64(firstByte>>2) |
			int64(buffer.readPrimitive8())<<6
	case length30Bit:
		return int64(firstByte>>2) |
			int64(buffer.readPrimitive8())<<6 |
			int64(buffer.readPrimitive8())<<14 |
			int64(buffer.readPrimitive8())<<22
	case length62Bit:
		return int64(firstByte>>2) |
			int64(buffer.readPrimitive8())<<6 |
			int64(buffer.readPrimitive8())<<14 |
			int64(buffer.readPrimitive8())<<22 |
			int64(buffer.readPrimitive8())<<30 |
			int64(buffer.readPrimitive8())<<38 |
			int64(buffer.readPrimitive8())<<46 |
			int64(buffer.readPrimitive8())<<54
	default: // TODO: 128 bit
		return 0
	}
}

func (buffer *decodeBuffer) readNegInt64() int64 {
	value := buffer.readPrimitive64()
	// TODO: This won't be an error once 128 bit support is added
	if value&0x8000000000000000 != 0 {
		panic(decoderError(fmt.Errorf("Value %v is too big to be represented as negative", value)))
		return 0
	} else {
		return -int64(value)
	}
}

type Decoder struct {
	buffer               decodeBuffer
	streamOffset         int64
	containerDepth       int
	currentArrayType     arrayType
	currentContainerType []containerType
	arrayBytesRemaining  int64
	arrayDecodeCallback  func([]byte) error
	arrayLengthCallback  func(uint64) error
	callbacks            *DecoderCallbacks
}

func (decoder *Decoder) enterContainer(newContainerType containerType) {
	if decoder.containerDepth >= len(decoder.currentContainerType) {
		panic(decoderError(fmt.Errorf("Exceeded max container depth of %v", len(decoder.currentContainerType))))
	}
	decoder.containerDepth++
	decoder.currentContainerType[decoder.containerDepth] = newContainerType
}

func (decoder *Decoder) leaveContainer() containerType {
	if decoder.containerDepth <= 0 {
		panic(decoderError(fmt.Errorf("Got container end but not in a container")))
	}
	decoder.containerDepth--
	return decoder.currentContainerType[decoder.containerDepth+1]
}

func (decoder *Decoder) beginArray(newArrayType arrayType) {
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

func (decoder *Decoder) setArrayLength(length int64) {
	decoder.arrayBytesRemaining = length
	checkCallback(decoder.arrayLengthCallback(uint64(decoder.arrayBytesRemaining)))
}

func (decoder *Decoder) decodeArrayLength() {
	decoder.setArrayLength(decoder.buffer.readArrayLength())
}

func (decoder *Decoder) decodeArrayData() {
	bytesToDecode := decoder.arrayBytesRemaining
	bytesRemaining := len(decoder.buffer.data) - decoder.buffer.pos
	if int64(bytesRemaining) < bytesToDecode {
		bytesToDecode = int64(bytesRemaining)
	}
	bytes := decoder.buffer.readPrimitiveBytes(int(bytesToDecode))
	decoder.arrayBytesRemaining -= bytesToDecode
	if decoder.arrayBytesRemaining == 0 {
		decoder.currentArrayType = arrayTypeNone
	}
	checkCallback(decoder.arrayDecodeCallback(bytes))
}

func (decoder *Decoder) decodeStringOfLength(length int64) {
	decoder.beginArray(arrayTypeString)
	decoder.setArrayLength(length)
	decoder.decodeArrayData()
}

func (decoder *Decoder) Feed(data []byte) (err error) {
	defer func() {
		decoder.streamOffset += int64(decoder.buffer.pos)
		decoder.buffer.data = nil
		if r := recover(); r != nil {
			switch r.(type) {
			case endOfBufferExit:
				err = nil
			case callbackError:
				offset := (decoder.streamOffset + int64(decoder.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: error from callback: %v", offset, r)
			case decoderError:
				offset := (decoder.streamOffset + int64(decoder.buffer.pos))
				err = fmt.Errorf("cbe: offset %v: %v", offset, r)
			default:
				err = fmt.Errorf("cbe: internal error: %v", r)
			}
		}
	}()

	decoder.buffer.data = data
	decoder.buffer.pos = 0

	if decoder.currentArrayType != arrayTypeNone {
		decoder.decodeArrayData()
	}

	for {
		dataType := decoder.buffer.readType()
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
			checkCallback(decoder.callbacks.OnFloat32(decoder.buffer.readFloat32()))
		case typeFloat64:
			checkCallback(decoder.callbacks.OnFloat64(decoder.buffer.readFloat64()))
		case typePosInt8:
			checkCallback(decoder.callbacks.OnUint(decoder.buffer.readPrimitive8()))
		case typePosInt16:
			checkCallback(decoder.callbacks.OnUint(decoder.buffer.readPrimitive16()))
		case typePosInt32:
			checkCallback(decoder.callbacks.OnUint(decoder.buffer.readPrimitive32()))
		case typePosInt64:
			checkCallback(decoder.callbacks.OnUint64(decoder.buffer.readPrimitive64()))
		case typeNegInt8:
			checkCallback(decoder.callbacks.OnInt(-int(decoder.buffer.readPrimitive8())))
		case typeNegInt16:
			checkCallback(decoder.callbacks.OnInt(-int(decoder.buffer.readPrimitive16())))
		case typeNegInt32:
			value := -int64(decoder.buffer.readPrimitive32())
			if value < math.MinInt32 {
				checkCallback(decoder.callbacks.OnInt64(value))
			} else {
				checkCallback(decoder.callbacks.OnInt(int(value)))
			}
		case typeNegInt64:
			checkCallback(decoder.callbacks.OnInt64(decoder.buffer.readNegInt64()))
		case typeTime:
			// TODO: Specify time zone?
			checkCallback(decoder.callbacks.OnTime(decoder.buffer.readTime().AsTime()))
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
			oldContainerType := decoder.leaveContainer()
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
	if decoder.containerDepth > 0 {
		return fmt.Errorf("Document still has open containers")
	}
	return nil
}

func (decoder *Decoder) Decode(document []byte) error {
	if err := decoder.Feed(document); err != nil {
		return err
	}
	return decoder.End()
}
