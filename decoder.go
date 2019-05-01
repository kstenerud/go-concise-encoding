package cbe

import (
	"math"
	"time"

	"github.com/kstenerud/go-smalltime"
)

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

func (decoder *Decoder) setArrayLength(length uint64) error {
	// TODO: Error if already in array
	decoder.arrayBytesRemaining = length
	return decoder.arrayLengthCallback(decoder.arrayBytesRemaining)
}

func (decoder *Decoder) decodeArrayLength() error {
	return decoder.setArrayLength(decoder.readArrayLength())
}

func (decoder *Decoder) decodeArrayData() error {
	// TODO: Make this more intelligent
	// Needs to decide how much data is available and only grab that,
	// leaving the rest for the next feed() call.
	bytes := decoder.readPrimitiveBytes(decoder.arrayBytesRemaining)
	if err := decoder.arrayDecodeCallback(bytes); err != nil {
		return err
	}
	return nil
}

func (decoder *Decoder) decodeStringOfLength(length uint64) error {
	decoder.beginArray(arrayTypeString)
	if err := decoder.setArrayLength(length); err != nil {
		return err
	}
	return decoder.decodeArrayData()
}

func (decoder *Decoder) feed(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Handle the situation
			// err = something
			// run off buffer: type runtime.errorString, value "runtime error: index out of range"
			// Could be brittle if callback code causes a panic...
		}
	}()

	// TODO: Check for unfinished array

	// TODO: replace this broken for loop
	for i := decoder.bufferPos; i < len(decoder.buffer); {
		dataType := decoder.readType()
		if int64(int8(dataType)) >= smallIntMin && int64(int8(dataType)) <= smallIntMax {
			if err := decoder.callbacks.OnInt(int(int8(dataType))); err != nil {
				return err
			}
			continue
		}
		switch dataType {
		case typeTrue:
			if err := decoder.callbacks.OnBool(true); err != nil {
				return err
			}
		case typeFalse:
			if err := decoder.callbacks.OnBool(false); err != nil {
				return err
			}
		case typeFloat32:
			if err := decoder.callbacks.OnFloat32(decoder.readFloat32()); err != nil {
				return err
			}
		case typeFloat64:
			if err := decoder.callbacks.OnFloat64(decoder.readFloat64()); err != nil {
				return err
			}
		case typePosInt8:
			if err := decoder.callbacks.OnUint(decoder.readPrimitive8()); err != nil {
				return err
			}
		case typePosInt16:
			if err := decoder.callbacks.OnUint(decoder.readPrimitive16()); err != nil {
				return err
			}
		case typePosInt32:
			if err := decoder.callbacks.OnUint(decoder.readPrimitive32()); err != nil {
				return err
			}
		case typePosInt64:
			if err := decoder.callbacks.OnUint64(decoder.readPrimitive64()); err != nil {
				return err
			}
		case typeNegInt8:
			if err := decoder.callbacks.OnInt(-int(decoder.readPrimitive8())); err != nil {
				return err
			}
		case typeNegInt16:
			if err := decoder.callbacks.OnInt(-int(decoder.readPrimitive16())); err != nil {
				return err
			}
		case typeNegInt32:
			value := -int64(decoder.readPrimitive32())
			if value < math.MinInt32 {
				if err := decoder.callbacks.OnInt64(value); err != nil {
					return err
				}
			} else {
				if err := decoder.callbacks.OnInt(int(value)); err != nil {
					return err
				}
			}
		case typeNegInt64:
			if err := decoder.callbacks.OnInt64(decoder.decodeNegInt64()); err != nil {
				return err
			}
		case typeTime:
			// TODO: Specify time zone?
			if err := decoder.callbacks.OnTime(decoder.readTime().AsTime()); err != nil {
				return err
			}
		case typeNil:
			if err := decoder.callbacks.OnNil(); err != nil {
				return err
			}
		case typePadding:
			// Ignore
		case typeList:
			decoder.enterContainer(containerTypeList)
			if err := decoder.callbacks.OnListBegin(); err != nil {
				return err
			}
		case typeMap:
			decoder.enterContainer(containerTypeMap)
			if err := decoder.callbacks.OnMapBegin(); err != nil {
				return err
			}
		case typeEndContainer:
			oldContainerType := decoder.getCurrentContainerType()
			decoder.leaveContainer()
			switch oldContainerType {
			case containerTypeList:
				if err := decoder.callbacks.OnListEnd(); err != nil {
					return err
				}
			case containerTypeMap:
				if err := decoder.callbacks.OnMapEnd(); err != nil {
					return err
				}
			}
		case typeBinary:
			decoder.beginArray(arrayTypeBinary)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeComment:
			decoder.beginArray(arrayTypeComment)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeString:
			decoder.beginArray(arrayTypeString)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeString0:
			decoder.beginArray(arrayTypeString)
			if err := decoder.setArrayLength(0); err != nil {
				return err
			}
		case typeString1:
			if err := decoder.decodeStringOfLength(1); err != nil {
				return err
			}
		case typeString2:
			if err := decoder.decodeStringOfLength(2); err != nil {
				return err
			}
		case typeString3:
			if err := decoder.decodeStringOfLength(3); err != nil {
				return err
			}
		case typeString4:
			if err := decoder.decodeStringOfLength(4); err != nil {
				return err
			}
		case typeString5:
			if err := decoder.decodeStringOfLength(5); err != nil {
				return err
			}
		case typeString6:
			if err := decoder.decodeStringOfLength(6); err != nil {
				return err
			}
		case typeString7:
			if err := decoder.decodeStringOfLength(7); err != nil {
				return err
			}
		case typeString8:
			if err := decoder.decodeStringOfLength(8); err != nil {
				return err
			}
		case typeString9:
			if err := decoder.decodeStringOfLength(9); err != nil {
				return err
			}
		case typeString10:
			if err := decoder.decodeStringOfLength(10); err != nil {
				return err
			}
		case typeString11:
			if err := decoder.decodeStringOfLength(11); err != nil {
				return err
			}
		case typeString12:
			if err := decoder.decodeStringOfLength(12); err != nil {
				return err
			}
		case typeString13:
			if err := decoder.decodeStringOfLength(13); err != nil {
				return err
			}
		case typeString14:
			if err := decoder.decodeStringOfLength(14); err != nil {
				return err
			}
		case typeString15:
			if err := decoder.decodeStringOfLength(15); err != nil {
				return err
			}
		}
		// TODO: 128 bit and decimal
	}
	return nil
}

func (decoder *Decoder) end() {
	// TODO
}
