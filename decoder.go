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
	buffer              []byte
	bufferPos           int
	currentArrayType    arrayType
	arrayBytesRemaining uint64
	arrayDecodeCallback func([]byte) error
	arrayLengthCallback func(uint64) error
	callbacks           *DecoderCallbacks
}

func (decoder *Decoder) readPrimitive8() uint {
	value := uint(decoder.buffer[decoder.bufferPos])
	decoder.bufferPos += 1
	return value
}

func (decoder *Decoder) readPrimitive16() uint {
	value := uint(decoder.buffer[decoder.bufferPos]) |
		uint(decoder.buffer[decoder.bufferPos+1])<<8
	decoder.bufferPos += 2
	return value
}

func (decoder *Decoder) readPrimitive32() uint {
	value := uint(decoder.buffer[decoder.bufferPos]) |
		uint(decoder.buffer[decoder.bufferPos+1])<<8 |
		uint(decoder.buffer[decoder.bufferPos+2])<<16 |
		uint(decoder.buffer[decoder.bufferPos+3])<<24
	decoder.bufferPos += 4
	return value
}

func (decoder *Decoder) readPrimitive64() uint64 {
	value := uint64(decoder.buffer[decoder.bufferPos]) |
		uint64(decoder.buffer[decoder.bufferPos+1])<<8 |
		uint64(decoder.buffer[decoder.bufferPos+2])<<16 |
		uint64(decoder.buffer[decoder.bufferPos+3])<<24 |
		uint64(decoder.buffer[decoder.bufferPos+4])<<32 |
		uint64(decoder.buffer[decoder.bufferPos+5])<<40 |
		uint64(decoder.buffer[decoder.bufferPos+6])<<48 |
		uint64(decoder.buffer[decoder.bufferPos+7])<<56
	decoder.bufferPos += 8
	return value
}

func (decoder *Decoder) readPrimitiveBytes(byteCount uint64) []byte {
	// TODO handle end of input from feed()
	bytes := decoder.buffer[decoder.bufferPos : decoder.bufferPos+int(byteCount)]
	decoder.bufferPos += int(byteCount)
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
	return float32(decoder.readPrimitive32())
}

func (decoder *Decoder) readFloat64() float64 {
	return float64(decoder.readPrimitive64())
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

func (decoder *Decoder) startArray(newArrayType arrayType) {
	if decoder.currentArrayType != arrayTypeNone {
		// TODO: error
	}
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
	decoder.arrayBytesRemaining = length
	return decoder.arrayLengthCallback(decoder.arrayBytesRemaining)
}

func (decoder *Decoder) decodeArrayLength() error {
	return decoder.setArrayLength(decoder.readArrayLength())
}

func (decoder *Decoder) decodeArrayData() error {
	// TODO: Make this more intelligent
	bytes := decoder.readPrimitiveBytes(decoder.arrayBytesRemaining)
	if err := decoder.arrayDecodeCallback(bytes); err != nil {
		return err
	}
	return nil
}

func (decoder *Decoder) feed(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Handle the situation
		}
	}()

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
			value := decoder.readPrimitive64()
			if value&0x8000000000000000 != 0 {
				// TODO: Error
			} else {
				if err := decoder.callbacks.OnInt64(-int64(value)); err != nil {
					return err
				}
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
			if err := decoder.callbacks.OnListBegin(); err != nil {
				return err
			}
		case typeMap:
			if err := decoder.callbacks.OnMapBegin(); err != nil {
				return err
			}
		case typeEndContainer:
			// TODO
		case typeBinary:
			decoder.startArray(arrayTypeBinary)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeComment:
			decoder.startArray(arrayTypeComment)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeString:
			decoder.startArray(arrayTypeString)
			if err := decoder.decodeArrayLength(); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeString0:
			if err := decoder.setArrayLength(0); err != nil {
				return err
			}
		case typeString1:
			if err := decoder.setArrayLength(1); err != nil {
				return err
			}
			if err := decoder.decodeArrayData(); err != nil {
				return err
			}
		case typeString2:
		case typeString3:
		case typeString4:
		case typeString5:
		case typeString6:
		case typeString7:
		case typeString8:
		case typeString9:
		case typeString10:
		case typeString11:
		case typeString12:
		case typeString13:
		case typeString14:
		case typeString15:
		}
		// TODO: 128 bit and decimal
	}
	return nil
}

func (decoder *Decoder) end() {

}
