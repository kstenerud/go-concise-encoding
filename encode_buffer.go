package cbe

import (
	"fmt"
	"math"
	"time"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

type encodeBuffer struct {
	data               []byte
	lastCommitPosition int
	position           int
	isExternalBuffer   bool
}

func (this *encodeBuffer) Init(externalBuffer []byte) {
	if externalBuffer != nil {
		this.data = externalBuffer
		this.isExternalBuffer = true
	} else {
		// TODO: Default size should be what?
		this.data = make([]byte, 1024)
	}
}

func (this *encodeBuffer) extend(byteCount int) error {
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

func (this *encodeBuffer) reserve(byteCount int) (buffer []byte, err error) {
	err = this.extend(byteCount)
	if err != nil {
		return nil, err
	}
	lastPosition := this.position
	this.position = lastPosition + byteCount
	return this.data[lastPosition:this.position], nil
}

func (this *encodeBuffer) RemainingSpace() int {
	return len(this.data) - this.position
}

func (this *encodeBuffer) Commit() {
	this.lastCommitPosition = this.position
}

func (this *encodeBuffer) EncodedBytes() []byte {
	return this.data[:this.lastCommitPosition]
}

func (this *encodeBuffer) EncodeBytes(bytes []byte) error {
	dst, err := this.reserve(len(bytes))
	if err != nil {
		return err
	}
	copy(dst, bytes)
	return nil
}

// Encode the maximum number of bytes possible.
// Return value will be less than len(bytes) if there wasn't enough room.
func (this *encodeBuffer) EncodeMaxBytes(bytes []byte) (byteCount int) {
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

func (this *encodeBuffer) EncodeTypeField(value typeField) error {
	dst, err := this.reserve(1)
	if err != nil {
		return err
	}
	dst[0] = byte(value)
	return nil
}

func (this *encodeBuffer) EncodeUint8(typeValue typeField, value byte) error {
	dst, err := this.reserve(2)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = value
	return nil
}

func (this *encodeBuffer) EncodeUint16(typeValue typeField, value uint16) error {
	dst, err := this.reserve(3)
	if err != nil {
		return err
	}
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	return nil
}

func (this *encodeBuffer) EncodeUint32(typeValue typeField, value uint32) error {
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

func (this *encodeBuffer) EncodeUint64(typeValue typeField, value uint64) error {
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

func (this *encodeBuffer) EncodeFloat32(value float32) error {
	return this.EncodeUint32(typeFloat32, math.Float32bits(value))
}

func (this *encodeBuffer) EncodeFloat64(value float64) error {
	return this.EncodeUint64(typeFloat64, math.Float64bits(value))
}

func (this *encodeBuffer) EncodeUint(typeValue typeField, valueIn uint64) error {
	this.EncodeTypeField(typeValue)
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

func (this *encodeBuffer) EncodeFloat(value float64, significantDigits int) error {
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

func (this *encodeBuffer) EncodeDate(value time.Time) (err error) {
	this.EncodeTypeField(typeDate)
	var byteCount int
	var ok bool
	byteCount, ok = compact_time.EncodeDate(value, this.data[this.position:])
	if err != nil {
		return err
	}
	if !ok {
		if err := this.extend(byteCount); err != nil {
			return err
		}
		byteCount, ok = compact_time.EncodeDate(value, this.data[this.position:])
		if !ok {
			return fmt.Errorf("BUG: Failed encoding time value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}

func (this *encodeBuffer) EncodeTime(value time.Time) (err error) {
	this.EncodeTypeField(typeTime)
	var byteCount int
	var ok bool
	byteCount, ok, err = compact_time.EncodeTime(value, this.data[this.position:])
	if err != nil {
		return err
	}
	if !ok {
		if err := this.extend(byteCount); err != nil {
			return err
		}
		byteCount, ok, err = compact_time.EncodeTime(value, this.data[this.position:])
		if !ok {
			return fmt.Errorf("BUG: Failed encoding time value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}

func (this *encodeBuffer) EncodeTimestamp(value time.Time) (err error) {
	this.EncodeTypeField(typeTimestamp)
	var byteCount int
	var ok bool
	byteCount, ok, err = compact_time.EncodeTimestamp(value, this.data[this.position:])
	if err != nil {
		return err
	}
	if !ok {
		if err := this.extend(byteCount); err != nil {
			return err
		}
		byteCount, ok, err = compact_time.EncodeTimestamp(value, this.data[this.position:])
		if !ok {
			return fmt.Errorf("BUG: Failed encoding time value %v to buffer length %v from position %v",
				value, len(this.data), this.position)
		}
	}
	this.position += byteCount
	return nil
}
