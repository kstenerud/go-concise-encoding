// decode_buffer
package cbe

import (
	"fmt"
	"math"

	"github.com/kstenerud/go-smalltime"
)

type failedByteCountReservation int
type endOfData int

type decodeBuffer struct {
	data           []byte
	pos            int
	bytesToConsume int
}

func (buffer *decodeBuffer) readType() typeField {
	if buffer.pos >= len(buffer.data) {
		panic(endOfData(1))
	}

	result := typeField(buffer.data[buffer.pos])
	buffer.pos++
	return result
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

func (buffer *decodeBuffer) peekPrimitive8() uint {
	if buffer.pos == len(buffer.data) {
		panic(failedByteCountReservation(1))
	}
	return uint(buffer.data[buffer.pos])
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
	switch int64(buffer.peekPrimitive8() & 3) {
	case length6Bit:
		return int64(buffer.readPrimitive8() >> 2)
	case length14Bit:
		return int64(buffer.readPrimitive16() >> 2)
	case length30Bit:
		return int64(buffer.readPrimitive32() >> 2)
	default:
		return int64(buffer.readPrimitive64() >> 2)
	}
}

func (buffer *decodeBuffer) readNegInt64() int64 {
	value := buffer.readPrimitive64()
	// TODO: This won't be an error once 128 bit support is added
	if value&0x8000000000000000 != 0 && value != 0x8000000000000000 {
		panic(decoderError{fmt.Errorf("Value %016x is too big to be represented as negative", value)})
		return 0
	} else {
		return -int64(value)
	}
}
