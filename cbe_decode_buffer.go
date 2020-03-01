package concise_encoding

import (
	"errors"
	"math"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

var BufferExhaustedError = errors.New("Buffer exhausted")

type cbeDecodeBuffer struct {
	data     []byte
	position int
}

func (this *cbeDecodeBuffer) Init(data []byte) {
	this.data = data
	this.position = 0
}

func (this *cbeDecodeBuffer) HasUnreadData() bool {
	return this.position < len(this.data)
}

func (this *cbeDecodeBuffer) DecodeUint8() uint8 {
	value := uint8(this.data[this.position])
	this.position++
	return value
}

func (this *cbeDecodeBuffer) DecodeUint16() uint16 {
	value := uint16(this.data[this.position]) |
		uint16(this.data[this.position+1])<<8
	this.position += 2
	return value
}

func (this *cbeDecodeBuffer) DecodeUint32() uint32 {
	value := uint32(this.data[this.position]) |
		uint32(this.data[this.position+1])<<8 |
		uint32(this.data[this.position+2])<<16 |
		uint32(this.data[this.position+3])<<24
	this.position += 4
	return value
}

func (this *cbeDecodeBuffer) DecodeUint64() uint64 {
	value := uint64(this.data[this.position]) |
		uint64(this.data[this.position+1])<<8 |
		uint64(this.data[this.position+2])<<16 |
		uint64(this.data[this.position+3])<<24 |
		uint64(this.data[this.position+4])<<32 |
		uint64(this.data[this.position+5])<<40 |
		uint64(this.data[this.position+6])<<48 |
		uint64(this.data[this.position+7])<<56
	this.position += 8
	return value
}

func (this *cbeDecodeBuffer) DecodeUint() uint64 {
	value, bytesDecoded, isComplete := vlq.DecodeRvlqFrom(this.data[this.position:])
	if !isComplete {
		panic(BufferExhaustedError)
	}
	this.position += bytesDecoded
	return uint64(value)
}

func (this *cbeDecodeBuffer) DecodeVersion() uint64 {
	return this.DecodeUint()
}

func (this *cbeDecodeBuffer) DecodeType() cbeTypeField {
	return cbeTypeField(this.DecodeUint8())
}

func (this *cbeDecodeBuffer) DecodeFloat32() float32 {
	return math.Float32frombits(this.DecodeUint32())
}

func (this *cbeDecodeBuffer) DecodeFloat64() float64 {
	return math.Float64frombits(this.DecodeUint64())
}

func (this *cbeDecodeBuffer) DecodeFloat() float64 {
	value, _, bytesDecoded, isComplete := compact_float.Decode(this.data[this.position:])
	if !isComplete {
		panic(BufferExhaustedError)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeDate() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeDate(this.data[this.position:])
	if err := value.Validate(); err != nil {
		panic(err)
	}
	if !isComplete {
		panic(BufferExhaustedError)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeTime() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTime(this.data[this.position:])
	if err := value.Validate(); err != nil {
		panic(err)
	}
	if !isComplete {
		panic(BufferExhaustedError)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeTimestamp() *compact_time.Time {
	value, bytesDecoded, isComplete := compact_time.DecodeTimestamp(this.data[this.position:])
	if err := value.Validate(); err != nil {
		panic(err)
	}
	if !isComplete {
		panic(BufferExhaustedError)
	}
	this.position += bytesDecoded
	return value
}

func (this *cbeDecodeBuffer) DecodeBytes(byteCount int) []byte {
	value := this.data[this.position : this.position+byteCount]
	this.position += byteCount
	return value
}

func (this *cbeDecodeBuffer) DecodeChunkHeader() (length uint64, isFinalChunk bool) {
	value := this.DecodeUint()
	return value >> 1, value&1 == 1
}
