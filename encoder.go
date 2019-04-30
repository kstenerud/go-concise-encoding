/*
 */
package cbe

import (
	"math"
	"time"

	// "github.com/ericlagergren/decimal"
	"github.com/kstenerud/go-smalltime"
	// "github.com/mewmew/float"
	// "github.com/shabbyrobe/go-num"
)

const (
	maxValue6Bit  uint64 = 0x3f
	maxValue14Bit uint64 = 0x3fff
	maxValue30Bit uint64 = 0x3fffffff
)

type arrayType int

const (
	arrayTypeNone arrayType = iota
	arrayTypeBinary
	arrayTypeString
	arrayTypeComment
)

func is6BitLength(value uint64) bool {
	return value <= maxValue6Bit
}

func is14BitLength(value uint64) bool {
	return value <= maxValue14Bit
}

func is30BitLength(value uint64) bool {
	return value <= maxValue30Bit
}

func intFitsInSmallint(value int64) bool {
	return value >= smallIntMin && value <= smallIntMax
}

func uintFitsInSmallint(value uint64) bool {
	return value <= uint64(smallIntMax)
}

func fitsInUint8(value uint64) bool {
	return value <= math.MaxUint8
}

func fitsInUint16(value uint64) bool {
	return value <= math.MaxUint16
}

func fitsInUint32(value uint64) bool {
	return value <= math.MaxUint32
}

type Encoder struct {
	maxContainerDepth int
	arrayType         arrayType
	encoded           []byte
}

func New(maxContainerDepth int) *Encoder {
	encoder := new(Encoder)
	encoder.maxContainerDepth = maxContainerDepth
	encoder.encoded = make([]byte, 0)
	return encoder
}

func (encoder *Encoder) addBytes(bytes []byte) {
	encoder.encoded = append(encoder.encoded, bytes...)
}

func (encoder *Encoder) addPrimitive8(value byte) {
	encoder.encoded = append(encoder.encoded, value)
}

func (encoder *Encoder) addPrimitive16(value uint16) {
	encoder.addBytes([]byte{byte(value), byte(value >> 8)})
}

func (encoder *Encoder) addPrimitive32(value uint32) {
	encoder.addBytes([]byte{
		byte(value), byte(value >> 8),
		byte(value >> 16), byte(value >> 24),
	})
}

func (encoder *Encoder) addPrimitive64(value uint64) {
	encoder.addBytes([]byte{
		byte(value), byte(value >> 8), byte(value >> 16),
		byte(value >> 24), byte(value >> 32), byte(value >> 40),
		byte(value >> 48), byte(value >> 56),
	})
}

func (encoder *Encoder) addType(typeValue typeField) {
	encoder.addPrimitive8(byte(typeValue))
}

func (encoder *Encoder) addArrayLength(length uint64) {
	switch {
	case is6BitLength(length):
		encoder.addPrimitive8(byte(length<<2 | length6Bit))
	case is14BitLength(length):
		encoder.addPrimitive16(uint16(length<<2 | length14Bit))
	case is30BitLength(length):
		encoder.addPrimitive32(uint32(length<<2 | length30Bit))
	default:
		encoder.addPrimitive64(uint64(length<<2 | length62Bit))
	}
}

func (encoder *Encoder) enterArray(newArrayType arrayType) {
	// TODO: Rework API a bit so that String() etc don't have dual meanings
	if encoder.arrayType != arrayTypeNone && encoder.arrayType != newArrayType {
		panic("Cannot start new array when already in an array")
	}
	encoder.arrayType = newArrayType
}

func (encoder *Encoder) leaveArray() {
	encoder.arrayType = arrayTypeNone
}

func (encoder *Encoder) Padding(byteCount int) *Encoder {
	for i := 0; i < byteCount; i++ {
		encoder.addType(typePadding)
	}
	return encoder
}

func (encoder *Encoder) Nil() *Encoder {
	encoder.addType(typeNil)
	return encoder
}

func (encoder *Encoder) Uint(value uint64) *Encoder {
	switch {
	case uintFitsInSmallint(value):
		encoder.addPrimitive8(byte(value))
	case fitsInUint8(value):
		encoder.addType(typePosInt8)
		encoder.addPrimitive8(uint8(value))
	case fitsInUint16(value):
		encoder.addType(typePosInt16)
		encoder.addPrimitive16(uint16(value))
	case fitsInUint32(value):
		encoder.addType(typePosInt32)
		encoder.addPrimitive32(uint32(value))
	default:
		encoder.addType(typePosInt64)
		encoder.addPrimitive64(value)
	}
	return encoder
}

func (encoder *Encoder) Int(value int64) *Encoder {
	if value >= 0 {
		encoder.Uint(uint64(value))
		return encoder
	}

	uvalue := uint64(-value)

	switch {
	case intFitsInSmallint(value):
		encoder.addPrimitive8(byte(value))
	case fitsInUint8(uvalue):
		encoder.addType(typeNegInt8)
		encoder.addPrimitive8(uint8(uvalue))
	case fitsInUint16(uvalue):
		encoder.addType(typeNegInt16)
		encoder.addPrimitive16(uint16(uvalue))
	case fitsInUint32(uvalue):
		encoder.addType(typeNegInt32)
		encoder.addPrimitive32(uint32(uvalue))
	default:
		encoder.addType(typeNegInt64)
		encoder.addPrimitive64(uvalue)
	}
	return encoder
}

func (encoder *Encoder) Float(value float64) *Encoder {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		encoder.addType(typeFloat32)
		encoder.addPrimitive32(math.Float32bits(asfloat32))
	} else {
		encoder.addType(typeFloat64)
		encoder.addPrimitive64(math.Float64bits(value))
	}
	return encoder
}

func (encoder *Encoder) Time(value time.Time) *Encoder {
	encoder.addType(typeTime)
	encoder.addPrimitive64(uint64(smalltime.FromTime(value)))
	return encoder
}

func (encoder *Encoder) ListBegin() *Encoder {
	encoder.addType(typeList)
	return encoder
}

func (encoder *Encoder) endContainer() *Encoder {
	encoder.addType(typeEndContainer)
	return encoder
}

func (encoder *Encoder) ListEnd() *Encoder {
	encoder.endContainer()
	return encoder
}

func (encoder *Encoder) MapBegin() *Encoder {
	encoder.addType(typeMap)
	return encoder
}

func (encoder *Encoder) MapEnd() *Encoder {
	encoder.endContainer()
	return encoder
}

func (encoder *Encoder) BinaryBegin(length uint64) *Encoder {
	encoder.enterArray(arrayTypeBinary)
	encoder.addType(typeBinary)
	encoder.addArrayLength(length)
	return encoder
}

func (encoder *Encoder) BinaryData(value []byte) *Encoder {
	// TODO: sanity checks
	encoder.addBytes([]byte(value))
	// TODO: If all bytes written
	if true {
		encoder.leaveArray()
	}
	return encoder
}

func (encoder *Encoder) Binary(value []byte) *Encoder {
	return encoder.BinaryBegin(uint64(len(value))).BinaryData(value)
}

func (encoder *Encoder) StringBegin(length uint64) *Encoder {
	encoder.enterArray(arrayTypeString)
	if length <= 15 {
		encoder.addType(typeString0 + typeField(length))
	} else {
		encoder.addType(typeString)
		encoder.addArrayLength(length)
	}
	return encoder
}

func (encoder *Encoder) StringData(value []byte) *Encoder {
	// TODO: sanity checks
	encoder.addBytes([]byte(value))
	// TODO: If all bytes written
	if true {
		encoder.leaveArray()
	}
	return encoder
}

func (encoder *Encoder) String(value string) *Encoder {
	return encoder.StringBegin(uint64(len(value))).StringData([]byte(value))
}

func (encoder *Encoder) Comment(value string) *Encoder {
	encoder.enterArray(arrayTypeComment)
	encoder.addType(typeComment)
	encoder.addArrayLength(uint64(len(value)))
	encoder.addBytes([]byte(value))
	encoder.leaveArray()
	return encoder
}

func (encoder *Encoder) Encoded() []byte {
	return encoder.encoded
}
