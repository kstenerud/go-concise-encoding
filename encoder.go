/*
 */
package cbe

import (
	"fmt"
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

type CbeEncoder struct {
	containerDepth       int
	currentArrayType     arrayType
	remainingArrayLength int64
	currentContainerType []containerType
	encoded              []byte
}

func NewCbeEncoder(maxContainerDepth int) *CbeEncoder {
	encoder := new(CbeEncoder)
	encoder.currentContainerType = make([]containerType, maxContainerDepth)
	encoder.encoded = make([]byte, 0)
	return encoder
}

func (encoder *CbeEncoder) encodeBytes(bytes []byte) {
	encoder.encoded = append(encoder.encoded, bytes...)
}

func (encoder *CbeEncoder) encodePrimitive8(value byte) {
	encoder.encoded = append(encoder.encoded, value)
}

func (encoder *CbeEncoder) encodePrimitive16(value uint16) {
	encoder.encodeBytes([]byte{byte(value), byte(value >> 8)})
}

func (encoder *CbeEncoder) encodePrimitive32(value uint32) {
	encoder.encodeBytes([]byte{
		byte(value), byte(value >> 8),
		byte(value >> 16), byte(value >> 24),
	})
}

func (encoder *CbeEncoder) encodePrimitive64(value uint64) {
	encoder.encodeBytes([]byte{
		byte(value), byte(value >> 8), byte(value >> 16),
		byte(value >> 24), byte(value >> 32), byte(value >> 40),
		byte(value >> 48), byte(value >> 56),
	})
}

func (encoder *CbeEncoder) encodeTypeField(typeValue typeField) {
	encoder.encodePrimitive8(byte(typeValue))
}

func (encoder *CbeEncoder) encodeArrayLengthField(length int64) {
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

func (encoder *CbeEncoder) enterContainer(newContainerType containerType) {
	// TODO: Error if container depth >= max
	encoder.containerDepth++
	encoder.currentContainerType[encoder.containerDepth] = newContainerType
}

func (encoder *CbeEncoder) leaveContainer() {
	// TODO: Error if container depth == 0
	encoder.containerDepth--
}

func (encoder *CbeEncoder) Padding(byteCount int) error {
	for i := 0; i < byteCount; i++ {
		encoder.encodeTypeField(typePadding)
	}
	return nil
}

func (encoder *CbeEncoder) Nil() error {
	encoder.encodeTypeField(typeNil)
	return nil
}

func (encoder *CbeEncoder) Bool(value bool) error {
	if value {
		encoder.encodeTypeField(typeTrue)
	} else {
		encoder.encodeTypeField(typeFalse)
	}
	return nil
}

func (encoder *CbeEncoder) Uint(value uint64) error {
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
	return nil
}

func (encoder *CbeEncoder) Int(value int64) error {
	if value >= 0 {
		encoder.Uint(uint64(value))
		return nil
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
	return nil
}

func (encoder *CbeEncoder) Float(value float64) error {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		encoder.encodeTypeField(typeFloat32)
		encoder.encodePrimitive32(math.Float32bits(asfloat32))
	} else {
		encoder.encodeTypeField(typeFloat64)
		encoder.encodePrimitive64(math.Float64bits(value))
	}
	return nil
}

func (encoder *CbeEncoder) Time(value time.Time) error {
	if value.Nanosecond()%1000 == 0 {
		encoder.encodeTypeField(typeSmalltime)
		encoder.encodePrimitive64(uint64(smalltime.SmalltimeFromTime(value)))
	} else {
		encoder.encodeTypeField(typeNanotime)
		encoder.encodePrimitive64(uint64(smalltime.NanotimeFromTime(value)))
	}
	return nil
}

func (encoder *CbeEncoder) ListBegin() error {
	encoder.enterContainer(containerTypeList)
	encoder.encodeTypeField(typeList)
	return nil
}

func (encoder *CbeEncoder) containerEnd() error {
	encoder.leaveContainer()
	encoder.encodeTypeField(typeEndContainer)
	return nil
}

func (encoder *CbeEncoder) ListEnd() error {
	return encoder.containerEnd()
}

func (encoder *CbeEncoder) MapBegin() error {
	encoder.enterContainer(containerTypeMap)
	encoder.encodeTypeField(typeMap)
	return nil
}

func (encoder *CbeEncoder) MapEnd() error {
	return encoder.containerEnd()
}

func (encoder *CbeEncoder) arrayBegin(newArrayType arrayType, length uint64) error {
	if encoder.currentArrayType != arrayTypeNone && encoder.currentArrayType != newArrayType {
		return fmt.Errorf("Cannot start new array when already in an array")
	}
	encoder.currentArrayType = newArrayType
	encoder.remainingArrayLength = int64(length)
	return nil
}

func (encoder *CbeEncoder) arrayAddData(value []byte) error {
	length := int64(len(value))
	if length > encoder.remainingArrayLength {
		return fmt.Errorf("Data length exceeds array length by %v bytes", length-encoder.remainingArrayLength)
	}
	encoder.encodeBytes(value)
	encoder.remainingArrayLength -= length
	if encoder.remainingArrayLength == 0 {
		encoder.currentArrayType = arrayTypeNone
	}
	return nil
}

func (encoder *CbeEncoder) BinaryBegin(length uint64) error {
	if err := encoder.arrayBegin(arrayTypeBinary, length); err != nil {
		return err
	}
	encoder.encodeTypeField(typeBinary)
	encoder.encodeArrayLengthField(int64(length))
	return nil
}

func (encoder *CbeEncoder) BinaryData(value []byte) error {
	return encoder.arrayAddData(value)
}

func (encoder *CbeEncoder) Bytes(value []byte) error {
	if err := encoder.BinaryBegin(uint64(len(value))); err != nil {
		return err
	}
	return encoder.BinaryData(value)
}

func (encoder *CbeEncoder) StringBegin(length uint64) error {
	if err := encoder.arrayBegin(arrayTypeString, length); err != nil {
		return err
	}
	if length <= 15 {
		encoder.encodeTypeField(typeString0 + typeField(length))
	} else {
		encoder.encodeTypeField(typeString)
		encoder.encodeArrayLengthField(int64(length))
	}
	return nil
}

func (encoder *CbeEncoder) StringData(value []byte) error {
	return encoder.arrayAddData(value)
}

func (encoder *CbeEncoder) String(value string) error {
	if err := encoder.StringBegin(uint64(len(value))); err != nil {
		return err
	}

	return encoder.StringData([]byte(value))
}

func (encoder *CbeEncoder) CommentBegin(length uint64) error {
	if err := encoder.arrayBegin(arrayTypeComment, length); err != nil {
		return err
	}
	encoder.encodeTypeField(typeComment)
	encoder.encodeArrayLengthField(int64(length))
	return nil
}

func (encoder *CbeEncoder) CommentData(value []byte) error {
	return encoder.arrayAddData(value)
}

func (encoder *CbeEncoder) Comment(value string) error {
	if err := encoder.CommentBegin(uint64(len(value))); err != nil {
		return err
	}
	return encoder.CommentData([]byte(value))
}

func (encoder *CbeEncoder) End() error {
	if encoder.remainingArrayLength > 0 {
		return fmt.Errorf("Incomplete encode: Current array is unfinished")
	}
	return nil
}

func (encoder *CbeEncoder) Encoded() []byte {
	return encoder.encoded
}
