// decode_buffer
package cbe

import (
	"fmt"
	"math"

	"github.com/kstenerud/go-smalltime"
)

type failedByteCountReservation int

type decodeBuffer struct {
	data           []byte
	pos            int
	bytesToConsume int
}

func (buffer *decodeBuffer) readType() (result typeField, err error) {
	if buffer.pos >= len(buffer.data) {
		err = fmt.Errorf("End of data")
	} else {
		result = typeField(buffer.data[buffer.pos])
		buffer.pos++
	}
	return result, err
}

func (buffer *decodeBuffer) reserveBytes(byteCount int) {
	if buffer.pos+byteCount > len(buffer.data) {
		panic(failedByteCountReservation(byteCount))
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

func (buffer *decodeBuffer) readSmalltime() smalltime.Smalltime {
	return smalltime.Smalltime(buffer.readPrimitive64())
}

func (buffer *decodeBuffer) readNanotime() smalltime.Nanotime {
	return smalltime.Nanotime(buffer.readPrimitive64())
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
