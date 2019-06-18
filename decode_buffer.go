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

func (this *decodeBuffer) readType() typeField {
	if this.pos >= len(this.data) {
		panic(endOfData(1))
	}

	result := typeField(this.data[this.pos])
	this.pos++
	return result
}

func (this *decodeBuffer) reserveBytes(byteCount int) {
	if this.pos+byteCount > len(this.data) {
		panic(failedByteCountReservation(byteCount))
	}
	this.bytesToConsume = byteCount
}

func (this *decodeBuffer) consumeReservedBytes() {
	this.pos += this.bytesToConsume
}

func (this *decodeBuffer) peekPrimitive8(offset int) uint {
	if this.pos+offset >= len(this.data) {
		panic(failedByteCountReservation(offset + 1))
	}
	return uint(this.data[this.pos+offset])
}

func (this *decodeBuffer) readPrimitive8() uint {
	this.reserveBytes(1)
	value := uint(this.data[this.pos])
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) readPrimitive16() uint {
	this.reserveBytes(2)
	value := uint(this.data[this.pos]) |
		uint(this.data[this.pos+1])<<8
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) readPrimitive32() uint {
	this.reserveBytes(4)
	value := uint(this.data[this.pos]) |
		uint(this.data[this.pos+1])<<8 |
		uint(this.data[this.pos+2])<<16 |
		uint(this.data[this.pos+3])<<24
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) readPrimitive64() uint64 {
	this.reserveBytes(8)
	value := uint64(this.data[this.pos]) |
		uint64(this.data[this.pos+1])<<8 |
		uint64(this.data[this.pos+2])<<16 |
		uint64(this.data[this.pos+3])<<24 |
		uint64(this.data[this.pos+4])<<32 |
		uint64(this.data[this.pos+5])<<40 |
		uint64(this.data[this.pos+6])<<48 |
		uint64(this.data[this.pos+7])<<56
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) readPrimitiveBytes(byteCount int) []byte {
	this.reserveBytes(byteCount)
	bytes := this.data[this.pos : this.pos+byteCount]
	this.consumeReservedBytes()
	return bytes
}

func (this *decodeBuffer) readFloat32() float32 {
	return math.Float32frombits(uint32(this.readPrimitive32()))
}

func (this *decodeBuffer) readFloat64() float64 {
	return math.Float64frombits(this.readPrimitive64())
}

func (this *decodeBuffer) readSmalltime() smalltime.Smalltime {
	return smalltime.Smalltime(this.readPrimitive64())
}

func (this *decodeBuffer) readNanotime() smalltime.Nanotime {
	return smalltime.Nanotime(this.readPrimitive64())
}

func (this *decodeBuffer) readArrayLength() int64 {
	arrayLength := int64(0)
	lengthPos := 0
	for {
		next := this.peekPrimitive8(lengthPos)
		arrayLength |= int64(next&0x7f) << uint(7*lengthPos)
		lengthPos++
		if next&0x80 == 0 {
			break
		}
	}
	this.pos += lengthPos
	return arrayLength
}

func (this *decodeBuffer) readNegInt64() int64 {
	value := this.readPrimitive64()
	// TODO: This won't be an error once 128 bit support is added
	if value&0x8000000000000000 != 0 && value != 0x8000000000000000 {
		panic(decoderError{fmt.Errorf("Value %016x is too big to be represented as negative", value)})
		return 0
	} else {
		return -int64(value)
	}
}
