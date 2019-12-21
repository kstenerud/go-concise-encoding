package cbe

import (
	"fmt"
	"math"
	"time"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

type cbeEncodeBuffer struct {
	data               []byte
	lastCommitPosition int
	position           int
	isExternalBuffer   bool
}

func (this *cbeEncodeBuffer) Init(externalBuffer []byte) {
	if externalBuffer != nil {
		this.data = externalBuffer
		this.isExternalBuffer = true
	} else {
		// TODO: Default size should be what?
		this.data = make([]byte, 1024)
	}
}

func (this *cbeEncodeBuffer) extend(byteCount int) error {
	targetLength := this.position + byteCount
	updatedLength := len(this.data)
	if updatedLength >= targetLength {
		return nil
	}

	if this.isExternalBuffer {
		return fmt.Errorf("External buffer capacity exhausted. Tried to extend by %v bytes but only %v bytes left",
			byteCount, len(this.data)-this.position)
	}

	for updatedLength < targetLength {
		if updatedLength < 1 {
			return fmt.Errorf("Integer overflow attempting to extend buffer")
		}
		updatedLength = updatedLength * 2
	}
	newBuffer := make([]byte, updatedLength)
	copy(newBuffer, this.data)
	this.data = newBuffer
	return nil
}

func (this *cbeEncodeBuffer) reserve(byteCount int) (buffer []byte, err error) {
	err = this.extend(byteCount)
	if err != nil {
		return nil, err
	}
	lastPosition := this.position
	this.position = lastPosition + byteCount
	return this.data[lastPosition:this.position], nil
}

func (this *cbeEncodeBuffer) RemainingSpace() int {
	return len(this.data) - this.position
}

func (this *cbeEncodeBuffer) Commit() {
	this.lastCommitPosition = this.position
}

func (this *cbeEncodeBuffer) Rollback() {
	this.position = this.lastCommitPosition
}

func (this *cbeEncodeBuffer) encodeRVLQ(valueIn uint64) error {
	value := vlq.Rvlq(valueIn)
	byteCount, ok := value.EncodeTo(this.data[this.position:])
	if !ok {
		this.extend(byteCount)
		byteCount, ok = value.EncodeTo(this.data[this.position:])
		if !ok {
			return fmt.Errorf("BUG: Failed encoding vlq value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}

func (this *cbeEncodeBuffer) EncodedBytes() []byte {
	return this.data[:this.lastCommitPosition]
}

func (this *cbeEncodeBuffer) EncodeBytes(bytes []byte) error {
	dst, err := this.reserve(len(bytes))
	if err != nil {
		return err
	}
	copy(dst, bytes)
	return nil
}

// Encode the maximum number of bytes possible.
// Return value will be less than len(bytes) if there wasn't enough room.
func (this *cbeEncodeBuffer) EncodeMaxBytes(bytes []byte) (byteCount int) {
	dst := this.data[this.position:]
	dstCount := len(dst)
	byteCount = len(bytes)
	if dstCount < byteCount {
		byteCount = dstCount
		bytes = bytes[:byteCount]
	}
	copy(dst, bytes)
	this.position += byteCount
	return byteCount
}

func (this *cbeEncodeBuffer) EncodeVersion(version uint64) error {
	return this.encodeRVLQ(version)
}

func (this *cbeEncodeBuffer) EncodeTypeField(value typeField) error {
	dst, err := this.reserve(1)
	if err != nil {
		return err
	}
	dst[0] = byte(value)
	return nil
}

func (this *cbeEncodeBuffer) EncodeUint8(typeValue typeField, value byte) error {
	dst, err := this.reserve(2)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = value
	return nil
}

func (this *cbeEncodeBuffer) EncodeUint16(typeValue typeField, value uint16) error {
	dst, err := this.reserve(3)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	return nil
}

func (this *cbeEncodeBuffer) EncodeUint32(typeValue typeField, value uint32) error {
	dst, err := this.reserve(5)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	return nil
}

func (this *cbeEncodeBuffer) EncodeUint64(typeValue typeField, value uint64) error {
	dst, err := this.reserve(9)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	dst[5] = byte(value >> 32)
	dst[6] = byte(value >> 40)
	dst[7] = byte(value >> 48)
	dst[8] = byte(value >> 56)
	return nil
}

func (this *cbeEncodeBuffer) EncodeFloat32(value float32) error {
	return this.EncodeUint32(typeFloat32, math.Float32bits(value))
}

func (this *cbeEncodeBuffer) EncodeFloat64(value float64) error {
	return this.EncodeUint64(typeFloat64, math.Float64bits(value))
}

func (this *cbeEncodeBuffer) EncodeUint(typeValue typeField, value uint64) error {
	if err := this.EncodeTypeField(typeValue); err != nil {
		return err
	}
	return this.encodeRVLQ(value)
}

func (this *cbeEncodeBuffer) EncodeFloat(value float64, significantDigits int) error {
	this.EncodeTypeField(typeDecimal)
	byteCount, ok := compact_float.Encode(value, significantDigits, this.data[this.position:])
	if !ok {
		this.extend(byteCount)
		byteCount, ok = compact_float.Encode(value, significantDigits, this.data[this.position:])
		if !ok {
			return fmt.Errorf("BUG: Failed encoding float value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}

func (this *cbeEncodeBuffer) EncodeTime(value time.Time) (err error) {
	return this.EncodeCompactTime(compact_time.AsCompactTime(value))
}

func (this *cbeEncodeBuffer) EncodeCompactTime(value *compact_time.Time) (err error) {
	var timeType typeField
	switch value.TimeIs {
	case compact_time.TypeDate:
		timeType = typeDate
	case compact_time.TypeTime:
		timeType = typeTime
	case compact_time.TypeTimestamp:
		timeType = typeTimestamp
	}
	this.EncodeTypeField(timeType)
	var byteCount int
	var ok bool
	byteCount, ok, err = compact_time.Encode(value, this.data[this.position:])
	if err != nil {
		return err
	}
	if !ok {
		if err := this.extend(byteCount); err != nil {
			return err
		}
		byteCount, ok, err = compact_time.Encode(value, this.data[this.position:])
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("BUG: Failed encoding time value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}

func (this *cbeEncodeBuffer) EncodeArrayChunkHeader(length uint64, isFinalChunk bool) error {
	length <<= 1
	if !isFinalChunk {
		length |= 1
	}
	return this.encodeRVLQ(length)
}
