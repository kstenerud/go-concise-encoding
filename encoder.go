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

// ------------
// Array Length
// ------------

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

// -----------
// CBE Encoder
// -----------

type CbeEncoder struct {
	containerDepth       int
	currentArrayType     arrayType
	remainingArrayLength int64
	currentContainerType []containerType
	hasStoredMapKey      []bool
	encoded              []byte
	charValidator        Utf8Validator
}

// --------
// Internal
// --------

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
	if length == 0 {
		encoder.encodePrimitive8(0)
	}
	for length > 0 {
		b := byte(length & 0x7f)
		length >>= 7
		if length > 0 {
			b |= 0x80
		}
		encoder.encodePrimitive8(b)
	}
}

func (encoder *CbeEncoder) containerBegin(newContainerType containerType) error {
	if encoder.containerDepth+1 >= len(encoder.currentContainerType) {
		return fmt.Errorf("Max container depth exceeded")
	}
	encoder.containerDepth++
	encoder.currentContainerType[encoder.containerDepth] = newContainerType
	encoder.hasStoredMapKey[encoder.containerDepth] = false
	return nil
}

func (encoder *CbeEncoder) containerEnd() error {
	if encoder.containerDepth <= 0 {
		return fmt.Errorf("No containers are open")
	}
	encoder.containerDepth--
	encoder.encodeTypeField(typeEndContainer)
	encoder.flipMapKeyStatus()
	return nil
}

func (encoder *CbeEncoder) isExpectingMapKey() bool {
	return encoder.currentContainerType[encoder.containerDepth] == containerTypeMap &&
		!encoder.hasStoredMapKey[encoder.containerDepth]
}

func (encoder *CbeEncoder) isExpectingMapValue() bool {
	return encoder.currentContainerType[encoder.containerDepth] == containerTypeMap &&
		encoder.hasStoredMapKey[encoder.containerDepth]
}

func (encoder *CbeEncoder) flipMapKeyStatus() {
	encoder.hasStoredMapKey[encoder.containerDepth] = !encoder.hasStoredMapKey[encoder.containerDepth]
}

func (encoder *CbeEncoder) assertNotExpectingMapKey(keyType string) error {
	if encoder.isExpectingMapKey() {
		return fmt.Errorf("Cannot use type %v as a map key", keyType)
	}
	return nil
}

func (encoder *CbeEncoder) arrayBegin(newArrayType arrayType, length uint64) error {
	if encoder.currentArrayType != arrayTypeNone && encoder.currentArrayType != newArrayType {
		return fmt.Errorf("Cannot start new array when already in an array")
	}
	encoder.charValidator.Reset()
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
		if encoder.currentArrayType != arrayTypeComment {
			encoder.flipMapKeyStatus()
		}
		encoder.currentArrayType = arrayTypeNone
	}
	return nil
}

// ----------
// Public API
// ----------

func NewCbeEncoder(maxContainerDepth int) *CbeEncoder {
	encoder := new(CbeEncoder)
	encoder.currentContainerType = make([]containerType, maxContainerDepth+1)
	encoder.hasStoredMapKey = make([]bool, maxContainerDepth+1)
	encoder.encoded = make([]byte, 0)
	return encoder
}

func (encoder *CbeEncoder) Padding(byteCount int) error {
	for i := 0; i < byteCount; i++ {
		encoder.encodeTypeField(typePadding)
	}
	return nil
}

func (encoder *CbeEncoder) Nil() error {
	if err := encoder.assertNotExpectingMapKey("nil"); err != nil {
		return err
	}
	encoder.encodeTypeField(typeNil)
	encoder.flipMapKeyStatus()
	return nil
}

func (encoder *CbeEncoder) Bool(value bool) error {
	if value {
		encoder.encodeTypeField(typeTrue)
	} else {
		encoder.encodeTypeField(typeFalse)
	}
	encoder.flipMapKeyStatus()
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
	encoder.flipMapKeyStatus()
	return nil
}

func (encoder *CbeEncoder) Int(value int64) error {
	uvalue := uint64(-value)

	switch {
	case intFitsInSmallint(value):
		encoder.encodePrimitive8(byte(value))
	case value >= 0:
		return encoder.Uint(uint64(value))
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
	encoder.flipMapKeyStatus()
	return nil
}

func (encoder *CbeEncoder) Float(value float64) error {
	asfloat32 := float32(value)
	// TODO: Check if it fits in an int/uint
	if float64(asfloat32) == value {
		encoder.encodeTypeField(typeFloat32)
		encoder.encodePrimitive32(math.Float32bits(asfloat32))
	} else {
		encoder.encodeTypeField(typeFloat64)
		encoder.encodePrimitive64(math.Float64bits(value))
	}
	encoder.flipMapKeyStatus()
	return nil
}

// Add a time value. Times are converted to their UTC equivalents before storage.
func (encoder *CbeEncoder) Time(value time.Time) error {
	if value.Nanosecond()%1000 == 0 {
		encoder.encodeTypeField(typeSmalltime)
		encoder.encodePrimitive64(uint64(smalltime.SmalltimeFromTime(value)))
	} else {
		encoder.encodeTypeField(typeNanotime)
		encoder.encodePrimitive64(uint64(smalltime.NanotimeFromTime(value)))
	}
	encoder.flipMapKeyStatus()
	return nil
}

func (encoder *CbeEncoder) ListBegin() error {
	if err := encoder.assertNotExpectingMapKey("list"); err != nil {
		return err
	}
	if err := encoder.containerBegin(containerTypeList); err != nil {
		return err
	}
	encoder.encodeTypeField(typeList)
	return nil
}

func (encoder *CbeEncoder) ListEnd() error {
	return encoder.containerEnd()
}

// Begin a map. Any subsequent objects added are assumed to alternate between
// key and value entries in the map, until MapEnd() is called.
func (encoder *CbeEncoder) MapBegin() error {
	if err := encoder.assertNotExpectingMapKey("map"); err != nil {
		return err
	}
	if err := encoder.containerBegin(containerTypeMap); err != nil {
		return err
	}
	encoder.encodeTypeField(typeMap)
	return nil
}

func (encoder *CbeEncoder) MapEnd() error {
	if encoder.isExpectingMapValue() {
		return fmt.Errorf("Expecting map value for already stored key")
	}
	return encoder.containerEnd()
}

// Begin a byte array. Encoder expects subsequent calls to BytesData to provide
// a total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (encoder *CbeEncoder) BytesBegin(length uint64) error {
	if err := encoder.arrayBegin(arrayTypeBytes, length); err != nil {
		return err
	}
	encoder.encodeTypeField(typeBytes)
	encoder.encodeArrayLengthField(int64(length))
	return nil
}

func (encoder *CbeEncoder) BytesData(value []byte) error {
	return encoder.arrayAddData(value)
}

// Convenience function to completely fill a byte array in one call.
func (encoder *CbeEncoder) Bytes(value []byte) error {
	if err := encoder.BytesBegin(uint64(len(value))); err != nil {
		return err
	}
	return encoder.BytesData(value)
}

// Begin a string. Encoder expects subsequent calls to StringData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
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
	for _, ch := range value {
		if err := encoder.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
	}
	return encoder.arrayAddData(value)
}

// Convenience function to completely fill a string in one call.
func (encoder *CbeEncoder) String(value string) error {
	if err := encoder.StringBegin(uint64(len(value))); err != nil {
		return err
	}

	return encoder.StringData([]byte(value))
}

// Begin a comment. Encoder expects subsequent calls to CommentData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (encoder *CbeEncoder) CommentBegin(length uint64) error {
	if err := encoder.arrayBegin(arrayTypeComment, length); err != nil {
		return err
	}
	encoder.encodeTypeField(typeComment)
	encoder.encodeArrayLengthField(int64(length))
	return nil
}

func (encoder *CbeEncoder) CommentData(value []byte) error {
	for _, ch := range value {
		if err := encoder.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
		if encoder.charValidator.IsCompleteCharacter() {
			if err := ValidateCommentCharacter(encoder.charValidator.Character()); err != nil {
				return err
			}
		}
	}
	return encoder.arrayAddData(value)
}

// Convenience function to completely fill a comment in one call.
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
	if encoder.containerDepth > 0 {
		return fmt.Errorf("Not all containers have been closed")
	}
	return nil
}

func (encoder *CbeEncoder) EncodedBytes() []byte {
	return encoder.encoded
}
