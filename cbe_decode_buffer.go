// decode_buffer
package cbe

import (
	"math"
	"runtime/debug"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

type cbeDecodeBuffer struct {
	data               []byte
	position           int
	lastCommitPosition int
	bytesToConsume     int
}

func NewDecodeBuffer(data []byte) *cbeDecodeBuffer {
	this := new(cbeDecodeBuffer)
	this.Init(data)
	return this
}

func (this *cbeDecodeBuffer) Init(data []byte) {
	this.data = data
	this.position = 0
	this.lastCommitPosition = 0
}

func (this *cbeDecodeBuffer) Clear() {
	this.data = this.data[0:0]
	this.position = 0
	this.lastCommitPosition = 0
}

func (this *cbeDecodeBuffer) AddContents(data []byte) {
	this.data = append(this.data, data...)
}

func (this *cbeDecodeBuffer) FillToByteCount(buffer *cbeDecodeBuffer, fillToByteCount int) int {
	bytesToAppend := fillToByteCount - len(this.data)
	if bytesToAppend > len(buffer.data) {
		bytesToAppend = len(buffer.data)
	}
	if bytesToAppend > 0 {
		this.data = append(this.data, buffer.data[:bytesToAppend]...)
	}
	return bytesToAppend
}

func (this *cbeDecodeBuffer) ReplaceBuffer(newBuffer []byte) {
	this.data = newBuffer
	this.position = 0
	this.lastCommitPosition = 0
}

func (this *cbeDecodeBuffer) GetUncommittedByteCount() int {
	return len(this.data) - this.lastCommitPosition
}

func (this *cbeDecodeBuffer) Commit() {
	this.lastCommitPosition = this.position
}

func (this *cbeDecodeBuffer) Rollback() {
	this.position = this.lastCommitPosition
}

func (this *cbeDecodeBuffer) GetUncommittedBytes() []byte {
	return this.data[this.lastCommitPosition:]
}

func (this *cbeDecodeBuffer) DecodeUint8() (value uint8, err error) {
	if this.position+1 > len(this.data) {
		err = BufferExhaustedError
		return
	}
	if this.position < 0 {
		debug.PrintStack()
	}
	value = uint8(this.data[this.position])
	this.position++
	return
}

func (this *cbeDecodeBuffer) DecodeUint16() (value uint16, err error) {
	if this.position+2 > len(this.data) {
		err = BufferExhaustedError
		return
	}

	value = uint16(this.data[this.position]) |
		uint16(this.data[this.position+1])<<8
	this.position += 2
	return
}

func (this *cbeDecodeBuffer) DecodeUint32() (value uint32, err error) {
	if this.position+4 > len(this.data) {
		err = BufferExhaustedError
		return
	}

	value = uint32(this.data[this.position]) |
		uint32(this.data[this.position+1])<<8 |
		uint32(this.data[this.position+2])<<16 |
		uint32(this.data[this.position+3])<<24
	this.position += 4
	return
}

func (this *cbeDecodeBuffer) DecodeUint64() (value uint64, err error) {
	if this.position+8 > len(this.data) {
		err = BufferExhaustedError
		return
	}

	value = uint64(this.data[this.position]) |
		uint64(this.data[this.position+1])<<8 |
		uint64(this.data[this.position+2])<<16 |
		uint64(this.data[this.position+3])<<24 |
		uint64(this.data[this.position+4])<<32 |
		uint64(this.data[this.position+5])<<40 |
		uint64(this.data[this.position+6])<<48 |
		uint64(this.data[this.position+7])<<56
	this.position += 8
	return
}

func (this *cbeDecodeBuffer) DecodeUint() (value uint64, err error) {
	readValue, bytesDecoded, isComplete := vlq.DecodeRvlqFrom(this.data[this.position:])
	if !isComplete {
		err = BufferExhaustedError
		return
	}
	value = uint64(readValue)
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeType() (value typeField, err error) {
	var decodedValue uint8
	decodedValue, err = this.DecodeUint8()
	value = typeField(decodedValue)
	return
}

func (this *cbeDecodeBuffer) DecodeFloat32() (value float32, err error) {
	var bits uint32
	bits, err = this.DecodeUint32()
	value = math.Float32frombits(bits)
	return
}

func (this *cbeDecodeBuffer) DecodeFloat64() (value float64, err error) {
	var bits uint64
	bits, err = this.DecodeUint64()
	value = math.Float64frombits(bits)
	return
}

func (this *cbeDecodeBuffer) DecodeFloat() (value float64, err error) {
	var bytesDecoded int
	var isComplete bool
	value, bytesDecoded, isComplete = compact_float.Decode(this.data[this.position:])
	if !isComplete {
		err = BufferExhaustedError
		return
	}
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeDate() (value *compact_time.Time, err error) {
	var bytesDecoded int
	var isComplete bool
	value, bytesDecoded, isComplete, err = compact_time.DecodeDate(this.data[this.position:])
	if err != nil {
		return
	}
	if !isComplete {
		err = BufferExhaustedError
		return
	}
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeTime() (value *compact_time.Time, err error) {
	var bytesDecoded int
	var isComplete bool
	value, bytesDecoded, isComplete, err = compact_time.DecodeTime(this.data[this.position:])
	if err != nil {
		return
	}
	if !isComplete {
		err = BufferExhaustedError
		return
	}
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeTimestamp() (value *compact_time.Time, err error) {
	var bytesDecoded int
	var isComplete bool
	value, bytesDecoded, isComplete, err = compact_time.DecodeTimestamp(this.data[this.position:])
	if err != nil {
		return
	}
	if !isComplete {
		err = BufferExhaustedError
		return
	}
	this.position += bytesDecoded
	return
}

func (this *cbeDecodeBuffer) DecodeBytes(byteCount int) (value []byte, err error) {
	if this.position+byteCount > len(this.data) {
		err = BufferExhaustedError
		return
	}

	value = this.data[this.position : this.position+byteCount]
	this.position += byteCount
	return
}

func (this *cbeDecodeBuffer) DecodeArrayLength() (value uint64, err error) {
	return this.DecodeUint()
}
