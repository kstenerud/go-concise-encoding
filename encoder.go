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
	maxValue6Bit  int64 = 0x3f
	maxValue14Bit int64 = 0x3fff
	maxValue30Bit int64 = 0x3fffffff
)

func is6BitLength(value int64) bool {
	return value <= maxValue6Bit
}

func is14BitLength(value int64) bool {
	return value <= maxValue14Bit
}

func is30BitLength(value int64) bool {
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
	containerDepth       int
	currentArrayType     arrayType
	currentContainerType []containerType
	encoded              []byte
}

func New(maxContainerDepth int) *Encoder {
	encoder := new(Encoder)
	encoder.currentContainerType = make([]containerType, maxContainerDepth)
	encoder.encoded = make([]byte, 0)
	return encoder
}

func (encoder *Encoder) encodeBytes(bytes []byte) {
	encoder.encoded = append(encoder.encoded, bytes...)
}

func (encoder *Encoder) encodePrimitive8(value byte) {
	encoder.encoded = append(encoder.encoded, value)
}

func (encoder *Encoder) encodePrimitive16(value uint16) {
	encoder.encodeBytes([]byte{byte(value), byte(value >> 8)})
}

func (encoder *Encoder) encodePrimitive32(value uint32) {
	encoder.encodeBytes([]byte{
		byte(value), byte(value >> 8),
		byte(value >> 16), byte(value >> 24),
	})
}

func (encoder *Encoder) encodePrimitive64(value uint64) {
	encoder.encodeBytes([]byte{
		byte(value), byte(value >> 8), byte(value >> 16),
		byte(value >> 24), byte(value >> 32), byte(value >> 40),
		byte(value >> 48), byte(value >> 56),
	})
}

func (encoder *Encoder) encodeTypeField(typeValue typeField) {
	encoder.encodePrimitive8(byte(typeValue))
}

func (encoder *Encoder) encodeArrayLengthField(length int64) {
	switch {
	case is6BitLength(length):
		encoder.encodePrimitive8(byte(length<<2 | length6Bit))
	case is14BitLength(length):
		encoder.encodePrimitive16(uint16(length<<2 | length14Bit))
	case is30BitLength(length):
		encoder.encodePrimitive32(uint32(length<<2 | length30Bit))
	default:
		encoder.encodePrimitive64(uint64(length<<2 | length62Bit))
	}
}

func (encoder *Encoder) enterArray(newArrayType arrayType) {
	// TODO: Rework API a bit so that String() etc don't have dual meanings
	if encoder.currentArrayType != arrayTypeNone && encoder.currentArrayType != newArrayType {
		panic("Cannot start new array when already in an array")
	}
	encoder.currentArrayType = newArrayType
}

func (encoder *Encoder) leaveArray() {
	encoder.currentArrayType = arrayTypeNone
}

func (encoder *Encoder) enterContainer(newContainerType containerType) {
	// TODO: Error if container depth >= max
	encoder.containerDepth++
	encoder.currentContainerType[encoder.containerDepth] = newContainerType
}

func (encoder *Encoder) leaveContainer() {
	// TODO: Error if container depth == 0
	encoder.containerDepth--
}

func (encoder *Encoder) Padding(byteCount int) *Encoder {
	for i := 0; i < byteCount; i++ {
		encoder.encodeTypeField(typePadding)
	}
	return encoder
}

func (encoder *Encoder) Nil() *Encoder {
	encoder.encodeTypeField(typeNil)
	return encoder
}

func (encoder *Encoder) Uint(value uint64) *Encoder {
	switch {
	case uintFitsInSmallint(value):
		encoder.encodePrimitive8(byte(value))
	case fitsInUint8(value):
		encoder.encodeTypeField(typePosInt8)
		encoder.encodePrimitive8(uint8(value))
	case fitsInUint16(value):
		encoder.encodeTypeField(typePosInt16)
		encoder.encodePrimitive16(uint16(value))
	case fitsInUint32(value):
		encoder.encodeTypeField(typePosInt32)
		encoder.encodePrimitive32(uint32(value))
	default:
		encoder.encodeTypeField(typePosInt64)
		encoder.encodePrimitive64(value)
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
		encoder.encodePrimitive8(byte(value))
	case fitsInUint8(uvalue):
		encoder.encodeTypeField(typeNegInt8)
		encoder.encodePrimitive8(uint8(uvalue))
	case fitsInUint16(uvalue):
		encoder.encodeTypeField(typeNegInt16)
		encoder.encodePrimitive16(uint16(uvalue))
	case fitsInUint32(uvalue):
		encoder.encodeTypeField(typeNegInt32)
		encoder.encodePrimitive32(uint32(uvalue))
	default:
		encoder.encodeTypeField(typeNegInt64)
		encoder.encodePrimitive64(uvalue)
	}
	return encoder
}

func (encoder *Encoder) Float(value float64) *Encoder {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		encoder.encodeTypeField(typeFloat32)
		encoder.encodePrimitive32(math.Float32bits(asfloat32))
	} else {
		encoder.encodeTypeField(typeFloat64)
		encoder.encodePrimitive64(math.Float64bits(value))
	}
	return encoder
}

func (encoder *Encoder) Time(value time.Time) *Encoder {
	encoder.encodeTypeField(typeTime)
	encoder.encodePrimitive64(uint64(smalltime.FromTime(value)))
	return encoder
}

func (encoder *Encoder) ListBegin() *Encoder {
	encoder.enterContainer(containerTypeList)
	encoder.encodeTypeField(typeList)
	return encoder
}

func (encoder *Encoder) containerEnd() *Encoder {
	encoder.leaveContainer()
	encoder.encodeTypeField(typeEndContainer)
	return encoder
}

func (encoder *Encoder) ListEnd() *Encoder {
	encoder.containerEnd()
	return encoder
}

func (encoder *Encoder) MapBegin() *Encoder {
	encoder.enterContainer(containerTypeMap)
	encoder.encodeTypeField(typeMap)
	return encoder
}

func (encoder *Encoder) MapEnd() *Encoder {
	encoder.containerEnd()
	return encoder
}

func (encoder *Encoder) BinaryBegin(length uint64) *Encoder {
	encoder.enterArray(arrayTypeBinary)
	encoder.encodeTypeField(typeBinary)
	encoder.encodeArrayLengthField(int64(length))
	return encoder
}

func (encoder *Encoder) BinaryData(value []byte) *Encoder {
	// TODO: sanity checks
	encoder.encodeBytes([]byte(value))
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
		encoder.encodeTypeField(typeString0 + typeField(length))
	} else {
		encoder.encodeTypeField(typeString)
		encoder.encodeArrayLengthField(int64(length))
	}
	return encoder
}

func (encoder *Encoder) StringData(value []byte) *Encoder {
	// TODO: sanity checks
	encoder.encodeBytes([]byte(value))
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
	encoder.encodeTypeField(typeComment)
	encoder.encodeArrayLengthField(int64(len(value)))
	encoder.encodeBytes([]byte(value))
	encoder.leaveArray()
	return encoder
}

func (encoder *Encoder) Encoded() []byte {
	return encoder.encoded
}
