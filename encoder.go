package cbe

import (
	"fmt"
	"math"
	"net/url"
	"time"
)

// ------------
// Array Length
// ------------

func bitMask(bitCount uint) uint64 {
	return (1 << bitCount) - 1
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

func fitsInUint21(value uint64) bool {
	return value == (value & bitMask(21))
}

func fitsInUint49(value uint64) bool {
	return value == (value & bitMask(49))
}

// -----------
// CBE Encoder
// -----------

type Encoder struct {
	hasInlineContainer   bool
	containerDepth       int
	currentArrayType     arrayType
	remainingArrayLength int64
	currentContainerType []typeField
	hasStoredMapKey      []bool
	charValidator        Utf8Validator
	buffer               encodeBuffer
}

const defaultBufferSize = 1024

var containerTypes = map[ContainerType]typeField{
	ContainerTypeList:         typeList,
	ContainerTypeOrderedMap:   typeMapOrdered,
	ContainerTypeUnorderedMap: typeMapUnordered,
	ContainerTypeMetadataMap:  typeMapMetadata,
}

// --------
// Internal
// --------

func (this *Encoder) containerBegin(newContainerType typeField) error {
	if this.containerDepth+1 >= len(this.currentContainerType) {
		return fmt.Errorf("Max container depth exceeded")
	}
	this.containerDepth++
	this.currentContainerType[this.containerDepth] = newContainerType
	this.hasStoredMapKey[this.containerDepth] = false
	return nil
}

func (this *Encoder) isInsideMap() bool {
	switch this.currentContainerType[this.containerDepth] {
	case typeMapUnordered:
		return true
	case typeMapOrdered:
		return true
	case typeMapMetadata:
		return true
	}
	return false
}

func (this *Encoder) isExpectingMapKey() bool {
	return this.isInsideMap() && !this.hasStoredMapKey[this.containerDepth]
}

func (this *Encoder) isExpectingMapValue() bool {
	return this.isInsideMap() && this.hasStoredMapKey[this.containerDepth]
}

func (this *Encoder) flipMapKeyStatus() {
	this.hasStoredMapKey[this.containerDepth] = !this.hasStoredMapKey[this.containerDepth]
}

func (this *Encoder) assertNotExpectingMapKey(keyType string) error {
	if this.isExpectingMapKey() {
		return fmt.Errorf("Cannot use type %v as a map key", keyType)
	}
	return nil
}

func (this *Encoder) arrayBegin(newArrayType arrayType, length uint64) error {
	if this.currentArrayType != arrayTypeNone && this.currentArrayType != newArrayType {
		return fmt.Errorf("Cannot start new array when already in an array")
	}
	this.charValidator.Reset()
	this.currentArrayType = newArrayType
	this.remainingArrayLength = int64(length)
	return nil
}

func (this *Encoder) arrayAddData(value []byte) (bytesEncoded int, err error) {
	if int64(len(value)) > this.remainingArrayLength {
		return 0, fmt.Errorf("Data length exceeds array length by %v bytes", int64(len(value))-this.remainingArrayLength)
	}

	if len(value) > this.buffer.RemainingSpace() && this.buffer.isExternalBuffer {
		bytesEncoded = this.buffer.EncodeMaxBytes(value)
	} else {
		err = this.buffer.EncodeBytes(value)
		if err != nil {
			return bytesEncoded, err
		}
		bytesEncoded = len(value)
	}

	this.remainingArrayLength -= int64(bytesEncoded)
	if this.remainingArrayLength == 0 {
		if this.currentArrayType != arrayTypeComment {
			this.buffer.Commit()
			this.flipMapKeyStatus()
		}
		this.currentArrayType = arrayTypeNone
	}
	return bytesEncoded, err
}

// ----------
// Public API
// ----------

// Create a new encoder. if buffer is nil, the encoder allocates its own buffer.
func NewCbeEncoder(inlineContainerType ContainerType, buffer []byte, maxContainerDepth int) *Encoder {
	this := new(Encoder)
	this.Init(inlineContainerType, buffer, maxContainerDepth)
	return this
}

func (this *Encoder) Init(inlineContainerType ContainerType, externalBuffer []byte, maxContainerDepth int) {
	if inlineContainerType != ContainerTypeNone {
		maxContainerDepth++
		this.hasInlineContainer = true
	}
	this.buffer.Init(externalBuffer)
	this.currentContainerType = make([]typeField, maxContainerDepth+1)
	this.hasStoredMapKey = make([]bool, maxContainerDepth+1)

	if inlineContainerType != ContainerTypeNone {
		this.containerBegin(containerTypes[inlineContainerType])
	}
}

func (this *Encoder) Padding(byteCount int) error {
	for i := 0; i < byteCount; i++ {
		if err := this.buffer.EncodeTypeField(typePadding); err != nil {
			return err
		}
	}
	this.buffer.Commit()
	return nil
}

func (this *Encoder) Nil() error {
	if err := this.assertNotExpectingMapKey("nil"); err != nil {
		return err
	}
	if err := this.buffer.EncodeTypeField(typeNil); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Bool(value bool) error {
	typeValue := typeTrue
	if !value {
		typeValue = typeFalse
	}
	if err := this.buffer.EncodeTypeField(typeValue); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Uint(value uint64) error {
	switch {
	// TODO: vlq int
	// TODO: pos neg int
	case uintFitsInSmallint(value):
		if err := this.buffer.EncodeTypeField(typeField(value)); err != nil {
			return err
		}
	case fitsInUint8(value):
		if err := this.buffer.EncodeUint8(typePosInt8, uint8(value)); err != nil {
			return err
		}
	case fitsInUint16(value):
		if err := this.buffer.EncodeUint16(typePosInt16, uint16(value)); err != nil {
			return err
		}
	case fitsInUint21(value):
		if err := this.buffer.EncodeUint(typePosInt, value); err != nil {
			return err
		}
	case fitsInUint32(value):
		if err := this.buffer.EncodeUint32(typePosInt32, uint32(value)); err != nil {
			return err
		}
	case fitsInUint49(value):
		if err := this.buffer.EncodeUint(typePosInt, value); err != nil {
			return err
		}
	default:
		if err := this.buffer.EncodeUint64(typePosInt64, value); err != nil {
			return err
		}
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Int(value int64) error {
	uvalue := uint64(-value)

	switch {
	case intFitsInSmallint(value):
		if err := this.buffer.EncodeTypeField(typeField(value)); err != nil {
			return err
		}
	case value >= 0:
		return this.Uint(uint64(value))
	case fitsInUint8(uvalue):
		if err := this.buffer.EncodeUint8(typeNegInt8, uint8(uvalue)); err != nil {
			return err
		}
	case fitsInUint16(uvalue):
		if err := this.buffer.EncodeUint16(typeNegInt16, uint16(uvalue)); err != nil {
			return err
		}
	case fitsInUint21(uvalue):
		if err := this.buffer.EncodeUint(typeNegInt, uvalue); err != nil {
			return err
		}
	case fitsInUint32(uvalue):
		if err := this.buffer.EncodeUint32(typeNegInt32, uint32(uvalue)); err != nil {
			return err
		}
	case fitsInUint49(uvalue):
		if err := this.buffer.EncodeUint(typeNegInt, uvalue); err != nil {
			return err
		}
	default:
		if err := this.buffer.EncodeUint64(typeNegInt64, uvalue); err != nil {
			return err
		}
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) FloatRounded(value float64, significantDigits int) error {
	if significantDigits < 1 || significantDigits > 15 {
		return this.Float(value)
	}

	if err := this.buffer.EncodeFloat(value, significantDigits); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Float(value float64) error {
	asfloat32 := float32(value)
	// TODO: Check if it fits in an int/uint
	if float64(asfloat32) == value {
		if err := this.buffer.EncodeUint32(typeFloat32, math.Float32bits(asfloat32)); err != nil {
			return err
		}
	} else {
		if err := this.buffer.EncodeUint64(typeFloat64, math.Float64bits(value)); err != nil {
			return err
		}
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Date(value time.Time) error {
	if err := this.buffer.EncodeDate(value); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Time(value time.Time) error {
	if err := this.buffer.EncodeTime(value); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) Timestamp(value time.Time) error {
	if err := this.buffer.EncodeTimestamp(value); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) ContainerEnd() error {
	if this.containerDepth <= 0 {
		return fmt.Errorf("No containers are open")
	}
	if this.hasInlineContainer && this.containerDepth <= 1 {
		return fmt.Errorf("No containers are open")
	}
	switch this.currentContainerType[this.containerDepth] {
	case typeList:
		break
	case typeMapMetadata:
		if this.isExpectingMapValue() {
			return fmt.Errorf("Expecting map value for already stored key")
		}
	case typeMapOrdered:
		if this.isExpectingMapValue() {
			return fmt.Errorf("Expecting map value for already stored key")
		}
	case typeMapUnordered:
		if this.isExpectingMapValue() {
			return fmt.Errorf("Expecting map value for already stored key")
		}
	}
	this.containerDepth--
	if err := this.buffer.EncodeTypeField(typeEndContainer); err != nil {
		return err
	}
	this.buffer.Commit()
	this.flipMapKeyStatus()
	return nil
}

func (this *Encoder) ListBegin() error {
	if err := this.assertNotExpectingMapKey("list"); err != nil {
		return err
	}
	if err := this.containerBegin(typeList); err != nil {
		return err
	}
	err := this.buffer.EncodeTypeField(typeList)
	if err != nil {
		return err
	}
	this.buffer.Commit()
	return nil
}

// Begin an unordered map. Any subsequent objects added are assumed to alternate
// between key and value entries in the map, until MapEnd() is called.
func (this *Encoder) UnorderedMapBegin() error {
	if err := this.assertNotExpectingMapKey("map"); err != nil {
		return err
	}
	if err := this.containerBegin(typeMapUnordered); err != nil {
		return err
	}
	err := this.buffer.EncodeTypeField(typeMapUnordered)
	if err != nil {
		return err
	}
	this.buffer.Commit()
	return nil
}

// Begin an ordered map. Any subsequent objects added are assumed to alternate
// between key and value entries in the map, until MapEnd() is called.
func (this *Encoder) OrderedMapBegin() error {
	if err := this.assertNotExpectingMapKey("map"); err != nil {
		return err
	}
	if err := this.containerBegin(typeMapOrdered); err != nil {
		return err
	}
	err := this.buffer.EncodeTypeField(typeMapOrdered)
	if err != nil {
		return err
	}
	this.buffer.Commit()
	return nil
}

// Begin a metadata map. Any subsequent objects added are assumed to alternate
// between key and value entries in the map, until MapEnd() is called.
func (this *Encoder) MetadataMapBegin() error {
	if err := this.assertNotExpectingMapKey("map"); err != nil {
		return err
	}
	if err := this.containerBegin(typeMapMetadata); err != nil {
		return err
	}
	err := this.buffer.EncodeTypeField(typeMapMetadata)
	if err != nil {
		return err
	}
	this.buffer.Commit()
	return nil
}

// Begin a byte array. Encoder expects subsequent calls to BytesData to provide
// a total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *Encoder) BytesBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeBytes, length); err != nil {
		return err
	}
	if err := this.buffer.EncodeUint(typeBytes, uint64(length)); err != nil {
		return err
	}
	return nil
}

func (this *Encoder) validateArrayData(value []byte) error {
	switch this.currentArrayType {
	case arrayTypeBytes:
		return nil
	case arrayTypeComment:
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
	case arrayTypeString:
		for _, ch := range value {
			if err := this.charValidator.AddByte(int(ch)); err != nil {
				return err
			}
		}
	case arrayTypeURI:
		// TODO: URI validation
		return nil
	}
	return nil
}

func (this *Encoder) ArrayData(value []byte) (byteCount int, err error) {
	if err = this.validateArrayData(value); err != nil {
		return 0, err
	}
	byteCount, err = this.arrayAddData(value)
	if err == nil {
		this.buffer.Commit()
	}
	return byteCount, err
}

// Convenience function to completely fill a byte array in one call.
func (this *Encoder) Bytes(value []byte) error {
	bytesToEncode := len(value)
	if err := this.BytesBegin(uint64(bytesToEncode)); err != nil {
		return err
	}
	if err := this.validateArrayData(value); err != nil {
		return err
	}
	bytesEncoded, err := this.arrayAddData(value)
	if err != nil {
		return err
	}
	if bytesEncoded != bytesToEncode {
		return fmt.Errorf("Not enough room to encode %v bytes of binary data", len(value))
	}
	this.buffer.Commit()
	return nil
}

// Begin a string. Encoder expects subsequent calls to StringData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *Encoder) StringBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeString, length); err != nil {
		return err
	}
	if length <= 15 {
		if err := this.buffer.EncodeTypeField(typeString0 + typeField(length)); err != nil {
			return err
		}
	} else {
		if err := this.buffer.EncodeUint(typeString, uint64(length)); err != nil {
			return err
		}
	}
	return nil
}

// Convenience function to completely fill a string in one call.
func (this *Encoder) String(value string) error {
	bytesToEncode := len(value)
	if err := this.StringBegin(uint64(bytesToEncode)); err != nil {
		return err
	}
	if err := this.validateArrayData([]byte(value)); err != nil {
		return err
	}
	bytesEncoded, err := this.arrayAddData([]byte(value))
	if err != nil {
		return err
	}
	if bytesEncoded != bytesToEncode {
		return fmt.Errorf("Not enough room to encode %v bytes of string data", len(value))
	}
	this.buffer.Commit()
	return nil
}

// Begin a URI. Encoder expects subsequent calls to StringData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *Encoder) URIBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeURI, length); err != nil {
		return err
	}
	return this.buffer.EncodeUint(typeURI, uint64(length))
}

// Convenience function to completely fill a string in one call.
func (this *Encoder) URI(value *url.URL) error {
	asString := value.String()
	bytesToEncode := len(asString)
	if err := this.URIBegin(uint64(bytesToEncode)); err != nil {
		return err
	}
	if err := this.validateArrayData([]byte(asString)); err != nil {
		return err
	}
	bytesEncoded, err := this.arrayAddData([]byte(asString))
	if err != nil {
		return err
	}
	if bytesEncoded != bytesToEncode {
		return fmt.Errorf("Not enough room to encode %v bytes of string data", len(asString))
	}
	this.buffer.Commit()
	return nil
}

// Begin a comment. Encoder expects subsequent calls to CommentData to provide a
// total of exactly the length provided here.
// Only lengths up to 0x3fffffffffffffff are supported.
func (this *Encoder) CommentBegin(length uint64) error {
	if err := this.arrayBegin(arrayTypeComment, length); err != nil {
		return err
	}
	return this.buffer.EncodeUint(typeComment, uint64(length))
}

// Convenience function to completely fill a comment in one call.
func (this *Encoder) Comment(value string) error {
	bytesToEncode := len(value)
	if err := this.CommentBegin(uint64(bytesToEncode)); err != nil {
		return err
	}
	if err := this.validateArrayData([]byte(value)); err != nil {
		return err
	}
	bytesEncoded, err := this.arrayAddData([]byte(value))
	if err != nil {
		return err
	}
	if bytesEncoded != bytesToEncode {
		return fmt.Errorf("Not enough room to encode %v bytes of comment data", len(value))
	}
	this.buffer.Commit()
	return nil
}

func (this *Encoder) End() error {
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

func (this *Encoder) EncodedBytes() []byte {
	return this.buffer.EncodedBytes()
}
