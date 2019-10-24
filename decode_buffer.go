// decode_buffer
package cbe

import (
	"math"
	"time"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

type notEnoughBytesToDecodeArrayData int64
type notEnoughBytesToDecodeObject int
type notEnoughBytesToDecodeType int

type decodeBuffer struct {
	data               []byte
	position           int
	lastCommitPosition int
	bytesToConsume     int
}

func NewDecodeBuffer(data []byte) *decodeBuffer {
	this := new(decodeBuffer)
	this.Init(data)
	return this
}

func (this *decodeBuffer) Init(data []byte) {
	this.data = data
	this.position = 0
	this.lastCommitPosition = 0
	this.bytesToConsume = 0
}

func (this *decodeBuffer) Clear() {
	this.data = this.data[0:0]
	this.position = 0
	this.lastCommitPosition = 0
	this.bytesToConsume = 0
}

func (this *decodeBuffer) AddContents(data []byte) {
	this.data = append(this.data, data...)
}

func (this *decodeBuffer) FillFromBuffer(buffer *decodeBuffer, fillToByteCount int) int {
	bytesToAppend := fillToByteCount - len(this.data)
	if bytesToAppend > len(buffer.data) {
		bytesToAppend = len(buffer.data)
	}
	if bytesToAppend > 0 {
		this.data = append(this.data, buffer.data[:bytesToAppend]...)
	}
	return bytesToAppend
}

func (this *decodeBuffer) ReplaceBuffer(newBuffer []byte) {
	this.data = newBuffer
	this.position = 0
	this.lastCommitPosition = 0
	this.bytesToConsume = 0
}

func (this *decodeBuffer) reserve(byteCount int) {
	if this.position+byteCount > len(this.data) {
		panic(notEnoughBytesToDecodeObject(byteCount))
	}
	this.bytesToConsume = byteCount
}

func (this *decodeBuffer) consumeReservedBytes() {
	this.position += this.bytesToConsume
}

func (this *decodeBuffer) consumeBytes(byteCount int) {
	this.position += byteCount
}

func (this *decodeBuffer) RemainingByteCount() int {
	return len(this.data) - this.position
}

func (this *decodeBuffer) Commit() {
	this.lastCommitPosition = this.position
}

func (this *decodeBuffer) Rollback() {
	this.position = this.lastCommitPosition
}

func (this *decodeBuffer) GetUncommittedBytes() []byte {
	return this.data[this.lastCommitPosition:]
}

func (this *decodeBuffer) DecodeType() typeField {
	if this.position >= len(this.data) {
		panic(notEnoughBytesToDecodeType(1))
	}

	result := typeField(this.data[this.position])
	this.position++
	return result
}

func (this *decodeBuffer) DecodeUint8() uint {
	this.reserve(1)
	value := uint(this.data[this.position])
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) DecodeUint16() uint {
	this.reserve(2)
	value := uint(this.data[this.position]) |
		uint(this.data[this.position+1])<<8
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) DecodeUint32() uint {
	this.reserve(4)
	value := uint(this.data[this.position]) |
		uint(this.data[this.position+1])<<8 |
		uint(this.data[this.position+2])<<16 |
		uint(this.data[this.position+3])<<24
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) DecodeUint64() uint64 {
	this.reserve(8)
	value := uint64(this.data[this.position]) |
		uint64(this.data[this.position+1])<<8 |
		uint64(this.data[this.position+2])<<16 |
		uint64(this.data[this.position+3])<<24 |
		uint64(this.data[this.position+4])<<32 |
		uint64(this.data[this.position+5])<<40 |
		uint64(this.data[this.position+6])<<48 |
		uint64(this.data[this.position+7])<<56
	this.consumeReservedBytes()
	return value
}

func (this *decodeBuffer) DecodeFloat32() float32 {
	bits := uint32(this.DecodeUint32())
	return math.Float32frombits(bits)
}

func (this *decodeBuffer) DecodeFloat64() float64 {
	bits := this.DecodeUint64()
	return math.Float64frombits(bits)
}

func (this *decodeBuffer) DecodeFloat() float64 {
	result, bytesDecoded, isComplete := compact_float.Decode(this.data[this.position:])
	if !isComplete {
		panic(notEnoughBytesToDecodeObject(bytesDecoded + 1))
	}
	this.consumeBytes(bytesDecoded)
	return result
}

func (this *decodeBuffer) DecodeDate() time.Time {
	result, bytesDecoded, isComplete := compact_time.DecodeDate(this.data[this.position:])
	if !isComplete {
		panic(notEnoughBytesToDecodeObject(bytesDecoded + 1))
	}
	this.consumeBytes(bytesDecoded)
	return result
}

func (this *decodeBuffer) DecodeTime() time.Time {
	result, bytesDecoded, isComplete, err := compact_time.DecodeTime(this.data[this.position:])
	if err != nil {
		panic(err)
	}
	if !isComplete {
		panic(notEnoughBytesToDecodeObject(bytesDecoded + 1))
	}
	this.consumeBytes(bytesDecoded)
	return result
}

func (this *decodeBuffer) DecodeTimestamp() time.Time {
	result, bytesDecoded, isComplete, err := compact_time.DecodeTimestamp(this.data[this.position:])
	if err != nil {
		panic(err)
	}
	if !isComplete {
		panic(notEnoughBytesToDecodeObject(bytesDecoded + 1))
	}
	this.consumeBytes(bytesDecoded)
	return result
}

func (this *decodeBuffer) DecodeUint() uint64 {
	result, bytesDecoded, isComplete := vlq.DecodeRvlqFrom(this.data[this.position:])
	if !isComplete {
		panic(notEnoughBytesToDecodeObject(bytesDecoded + 1))
	}
	this.consumeBytes(bytesDecoded)
	return uint64(result)
}

func (this *decodeBuffer) DecodeBytes(byteCount int) []byte {
	this.reserve(byteCount)
	bytes := this.data[this.position : this.position+byteCount]
	this.consumeReservedBytes()
	return bytes
}

func (this *decodeBuffer) DecodeArrayLength() int64 {
	return int64(this.DecodeUint())
}
