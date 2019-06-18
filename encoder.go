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

func fitsInFloat32(value float64) bool {
	smaller := float32(value)
	return float64(smaller) == value
}

// -----------
// CBE Encoder
// -----------

type CbeEncoder struct {
	hasInlineContainer   bool
	containerDepth       int
	currentArrayType     arrayType
	remainingArrayLength int64
	currentContainerType []ContainerType
	hasStoredMapKey      []bool
	encodedBuffer        []byte
	charValidator        Utf8Validator
}

// --------
// Internal
// --------

func (this *CbeEncoder) encodeBytes(bytes []byte) {
	this.encodedBuffer = append(this.encodedBuffer, bytes...)
}

func (this *CbeEncoder) encodePrimitive8(value byte) {
	this.encodedBuffer = append(this.encodedBuffer, value)
}

func (this *CbeEncoder) encodePrimitive16(value uint16) {
	this.encodeBytes([]byte{byte(value), byte(value >> 8)})
}

func (this *CbeEncoder) encodePrimitive32(value uint32) {
	this.encodeBytes([]byte{
		byte(value), byte(value >> 8),
		byte(value >> 16), byte(value >> 24),
	})
}

func (this *CbeEncoder) encodePrimitive64(value uint64) {
	this.encodeBytes([]byte{
		byte(value), byte(value >> 8), byte(value >> 16),
		byte(value >> 24), byte(value >> 32), byte(value >> 40),
		byte(value >> 48), byte(value >> 56),
	})
}

func (this *CbeEncoder) encodeTypeField(typeValue typeField) {
	this.encodePrimitive8(byte(typeValue))
}

func (this *CbeEncoder) encodeArrayLengthField(length int64) {
	if length == 0 {
		this.encodePrimitive8(0)
	}
	for length > 0 {
		b := byte(length & 0x7f)
		length >>= 7
		if length > 0 {
			b |= 0x80
		}
		this.encodePrimitive8(b)
	}
}

func (this *CbeEncoder) containerBegin(newContainerType ContainerType) error {
	if this.containerDepth+1 >= len(this.currentContainerType) {
		return fmt.Errorf("Max container depth exceeded")
	}
	this.containerDepth++
	this.currentContainerType[this.containerDepth] = newContainerType
	this.hasStoredMapKey[this.containerDepth] = false
	return nil
}

func (this *CbeEncoder) containerEnd() error {
	if this.containerDepth <= 0 {
		return fmt.Errorf("No containers are open")
	}
	if this.hasInlineContainer && this.containerDepth <= 1 {
		return fmt.Errorf("No containers are open")
	}
	this.containerDepth--
	this.encodeTypeField(typeEndContainer)
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) isExpectingMapKey() bool {
	return this.currentContainerType[this.containerDepth] == ContainerTypeMap &&
		!this.hasStoredMapKey[this.containerDepth]
}

func (this *CbeEncoder) isExpectingMapValue() bool {
	return this.currentContainerType[this.containerDepth] == ContainerTypeMap &&
		this.hasStoredMapKey[this.containerDepth]
}

func (this *CbeEncoder) flipMapKeyStatus() {
	this.hasStoredMapKey[this.containerDepth] = !this.hasStoredMapKey[this.containerDepth]
}

func (this *CbeEncoder) assertNotExpectingMapKey(keyType string) error {
	if this.isExpectingMapKey() {
		return fmt.Errorf("Cannot use type %v as a map key", keyType)
	}
	return nil
}

func (this *CbeEncoder) arrayBegin(newArrayType arrayType, length uint64) error {
	if this.currentArrayType != arrayTypeNone && this.currentArrayType != newArrayType {
		return fmt.Errorf("Cannot start new array when already in an array")
	}
	this.charValidator.Reset()
	this.currentArrayType = newArrayType
	this.remainingArrayLength = int64(length)
	return nil
}

func (this *CbeEncoder) arrayAddData(value []byte) error {
	length := int64(len(value))
	if length > this.remainingArrayLength {
		return fmt.Errorf("Data length exceeds array length by %v bytes", length-this.remainingArrayLength)
	}
	this.encodeBytes(value)
	this.remainingArrayLength -= length
	if this.remainingArrayLength == 0 {
		if this.currentArrayType != arrayTypeComment {
			this.flipMapKeyStatus()
		}
		this.currentArrayType = arrayTypeNone
	}
	return nil
}

// ----------
// Public API
// ----------

func NewCbeEncoder(inlineContainerType ContainerType, maxContainerDepth int) *CbeEncoder {
	this := new(CbeEncoder)
	if inlineContainerType != ContainerTypeNone {
		maxContainerDepth++
		this.hasInlineContainer = true
	}
	this.currentContainerType = make([]ContainerType, maxContainerDepth+1)
	this.hasStoredMapKey = make([]bool, maxContainerDepth+1)
	this.encodedBuffer = make([]byte, 0)
	if inlineContainerType != ContainerTypeNone {
		this.containerBegin(inlineContainerType)
	}
	return this
}

func (this *CbeEncoder) Padding(byteCount int) error {
	for i := 0; i < byteCount; i++ {
		this.encodeTypeField(typePadding)
	}
	return nil
}

func (this *CbeEncoder) Nil() error {
	if err := this.assertNotExpectingMapKey("nil"); err != nil {
		return err
	}
	this.encodeTypeField(typeNil)
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) Bool(value bool) error {
	if value {
		this.encodeTypeField(typeTrue)
	} else {
		this.encodeTypeField(typeFalse)
	}
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) Uint(value uint64) error {
	switch {
	case uintFitsInSmallint(value):
		this.encodePrimitive8(byte(value))
	case fitsInUint8(value):
		this.encodeTypeField(typePosInt8)
		this.encodePrimitive8(uint8(value))
	case fitsInUint16(value):
		this.encodeTypeField(typePosInt16)
		this.encodePrimitive16(uint16(value))
	case fitsInUint32(value):
		this.encodeTypeField(typePosInt32)
		this.encodePrimitive32(uint32(value))
	default:
		this.encodeTypeField(typePosInt64)
		this.encodePrimitive64(value)
	}
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) Int(value int64) error {
	uvalue := uint64(-value)

	switch {
	case intFitsInSmallint(value):
		this.encodePrimitive8(byte(value))
	case value >= 0:
		return this.Uint(uint64(value))
	case fitsInUint8(uvalue):
		this.encodeTypeField(typeNegInt8)
		this.encodePrimitive8(uint8(uvalue))
	case fitsInUint16(uvalue):
		this.encodeTypeField(typeNegInt16)
		this.encodePrimitive16(uint16(uvalue))
	case fitsInUint32(uvalue):
		this.encodeTypeField(typeNegInt32)
		this.encodePrimitive32(uint32(uvalue))
	default:
		this.encodeTypeField(typeNegInt64)
		this.encodePrimitive64(uvalue)
	}
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) Float(value float64) error {
	asfloat32 := float32(value)
	// TODO: Check if it fits in an int/uint
	if float64(asfloat32) == value {
		this.encodeTypeField(typeFloat32)
		this.encodePrimitive32(math.Float32bits(asfloat32))
	} else {
		this.encodeTypeField(typeFloat64)
		this.encodePrimitive64(math.Float64bits(value))
	}
	this.flipMapKeyStatus()
	return nil
}

// Add a time value. Times are converted to their UTC equivalents before storage.
func (this *CbeEncoder) Time(value time.Time) error {
	if value.Nanosecond()%1000 == 0 {
		this.encodeTypeField(typeSmalltime)
		this.encodePrimitive64(uint64(smalltime.SmalltimeFromTime(value)))
	} else {
		this.encodeTypeField(typeNanotime)
		this.encodePrimitive64(uint64(smalltime.NanotimeFromTime(value)))
	}
	this.flipMapKeyStatus()
	return nil
}

func (this *CbeEncoder) ListBegin() error {
	if err := this.assertNotExpectingMapKey("list"); err != nil {
		return err
	}
	if err := this.containerBegin(ContainerTypeList); err != nil {
		return err
	}
	this.encodeTypeField(typeList)
	return nil
}

func (this *CbeEncoder) ListEnd() error {
	return this.containerEnd()
}

// Begin a map. Any subsequent objects added are assumed to alternate between
// key and value entries in the map, until MapEnd() is called.
func (this *CbeEncoder) MapBegin() error {
	if err := this.assertNotExpectingMapKey("map"); err != nil {
		return err
	}
	if err := this.containerBegin(ContainerTypeMap); err != nil {
		return err
	}
	this.encodeTypeField(typeMap)
	return nil
}

func (this *CbeEncoder) MapEnd() error {
	if this.isExpectingMapValue() {
		return fmt.Errorf("Expecting map value for already stored key")
	}
	return this.containerEnd()
}

// Begin a byte array. Encoder expects subsequent calls to BytesData to provide
// a total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *CbeEncoder) BytesBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeBytes, length); err != nil {
		return err
	}
	this.encodeTypeField(typeBytes)
	this.encodeArrayLengthField(int64(length))
	return nil
}

func (this *CbeEncoder) BytesData(value []byte) error {
	return this.arrayAddData(value)
}

// Convenience function to completely fill a byte array in one call.
func (this *CbeEncoder) Bytes(value []byte) error {
	if err := this.BytesBegin(uint64(len(value))); err != nil {
		return err
	}
	return this.BytesData(value)
}

// Begin a string. Encoder expects subsequent calls to StringData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *CbeEncoder) StringBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeString, length); err != nil {
		return err
	}
	if length <= 15 {
		this.encodeTypeField(typeString0 + typeField(length))
	} else {
		this.encodeTypeField(typeString)
		this.encodeArrayLengthField(int64(length))
	}
	return nil
}

func (this *CbeEncoder) StringData(value []byte) error {
	for _, ch := range value {
		if err := this.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
	}
	return this.arrayAddData(value)
}

// Convenience function to completely fill a string in one call.
func (this *CbeEncoder) String(value string) error {
	if err := this.StringBegin(uint64(len(value))); err != nil {
		return err
	}

	return this.StringData([]byte(value))
}

// Begin a comment. Encoder expects subsequent calls to CommentData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *CbeEncoder) CommentBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeComment, length); err != nil {
		return err
	}
	this.encodeTypeField(typeComment)
	this.encodeArrayLengthField(int64(length))
	return nil
}

func (this *CbeEncoder) CommentData(value []byte) error {
	for _, ch := range value {
		if err := this.charValidator.AddByte(int(ch)); err != nil {
			return err
		}
		if this.charValidator.IsCompleteCharacter() {
			if err := ValidateCommentCharacter(this.charValidator.Character()); err != nil {
				return err
			}
		}
	}
	return this.arrayAddData(value)
}

// Convenience function to completely fill a comment in one call.
func (this *CbeEncoder) Comment(value string) error {
	if err := this.CommentBegin(uint64(len(value))); err != nil {
		return err
	}
	return this.CommentData([]byte(value))
}

func (this *CbeEncoder) End() error {
	if this.remainingArrayLength > 0 {
		return fmt.Errorf("Incomplete encode: Current array is unfinished")
	}
	if this.containerDepth > 0 {
		if !(this.containerDepth == 1 && this.hasInlineContainer) {
			return fmt.Errorf("Not all containers have been closed")
		}
	}
	return nil
}

func (this *CbeEncoder) EncodedBytes() []byte {
	return this.encodedBuffer
}
