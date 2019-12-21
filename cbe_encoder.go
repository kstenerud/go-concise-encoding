package cbe

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/kstenerud/go-cbe/rules"
	"github.com/kstenerud/go-compact-time"
)

// TODO: Some ints would store better as float, and vice versa

// ---------
// Utilities
// ---------

const bitMask21 = uint64((1 << 21) - 1)
const bitMask49 = uint64((1 << 49) - 1)

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
	return value == (value & bitMask21)
}

func fitsInUint49(value uint64) bool {
	return value == (value & bitMask49)
}

// -----------
// CBE Encoder
// -----------

type CBEEncoder struct {
	buffer              cbeEncodeBuffer
	rules               rules.Rules
	inlineContainerType InlineContainerType
	context             fmt.Stringer // TODO
}

const defaultBufferSize = 1024

// ---------
// Utilities
// ---------

func (this *CBEEncoder) beginBytes() (err error) {
	return this.buffer.EncodeTypeField(typeBytes)
}

func (this *CBEEncoder) beginString() (err error) {
	return this.buffer.EncodeTypeField(typeString)
}

func (this *CBEEncoder) beginURI() (err error) {
	return this.buffer.EncodeTypeField(typeURI)
}

func (this *CBEEncoder) beginChunk(length uint64, isFinalChunk bool) (err error) {
	return this.buffer.EncodeArrayChunkHeader(length, isFinalChunk)
}

func (this *CBEEncoder) arrayData(value []byte) (bytesEncoded int, err error) {
	if len(value) > this.buffer.RemainingSpace() && this.buffer.isExternalBuffer {
		bytesEncoded = this.buffer.EncodeMaxBytes(value)
	} else {
		if err = this.buffer.EncodeBytes(value); err != nil {
			return
		}
		bytesEncoded = len(value)
	}
	return
}

func (this *CBEEncoder) arrayEntireData(value []byte) (err error) {
	bytesToEncode := len(value)
	if err = this.beginChunk(uint64(bytesToEncode), true); err != nil {
		return
	}
	bytesEncoded := 0
	bytesEncoded, err = this.arrayData(value)
	if err != nil {
		return
	}
	if bytesEncoded != bytesToEncode {
		return fmt.Errorf("Not enough room to encode %v bytes of binary data", len(value))
	}
	return
}

func (this *CBEEncoder) smallEntireString(value string) (err error) {
	byteCount := len(value)
	if err = this.buffer.EncodeTypeField(typeString0 + typeField(byteCount)); err != nil {
		return
	}
	bytesEncoded := 0
	bytesEncoded, err = this.arrayData([]byte(value))
	if err != nil {
		return
	}
	if bytesEncoded != byteCount {
		return fmt.Errorf("Not enough room to encode %v bytes of binary data", len(value))
	}
	return
}

func (this *CBEEncoder) largeEntireString(value string) (err error) {
	if err = this.beginString(); err != nil {
		return
	}
	return this.arrayEntireData([]byte(value))
}

// ----------
// Public API
// ----------

// Create a new encoder. if buffer is nil, the encoder allocates its own buffer.
func NewCBEEncoder(inlineContainerType InlineContainerType, buffer []byte, limits *rules.Limits) *CBEEncoder {
	this := new(CBEEncoder)
	this.Init(inlineContainerType, buffer, limits)
	return this
}

func (this *CBEEncoder) Init(inlineContainerType InlineContainerType, externalBuffer []byte, limits *rules.Limits) {
	this.rules.Init(cbeCodecVersion, limits)
	this.inlineContainerType = inlineContainerType
	switch inlineContainerType {
	case InlineContainerTypeList:
		if err := this.rules.BeginList(); err != nil {
			panic(fmt.Errorf("BUG: This should never happen: %v", err))
		}
	case InlineContainerTypeMap:
		if err := this.rules.BeginMap(); err != nil {
			panic(fmt.Errorf("BUG: This should never happen: %v", err))
		}
	}

	this.buffer.Init(externalBuffer)

	// TODO: Init context

	if err := this.buffer.EncodeVersion(cbeCodecVersion); err != nil {
		panic(fmt.Errorf("BUG: This should never happen: %v", err))
	}
	if err := this.rules.AddVersion(cbeCodecVersion); err != nil {
		panic(fmt.Errorf("BUG: This should never happen: %v", err))
	}
	this.buffer.Commit()
}

func (this *CBEEncoder) Padding(byteCount int) (err error) {
	for i := 0; i < byteCount; i++ {
		if err = this.buffer.EncodeTypeField(typePadding); err != nil {
			return
		}
	}
	if err = this.rules.AddPadding(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return nil
}

func (this *CBEEncoder) Nil() (err error) {
	if err = this.buffer.EncodeTypeField(typeNil); err != nil {
		return
	}
	if err = this.rules.AddNil(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return nil
}

func (this *CBEEncoder) Bool(value bool) (err error) {
	typeValue := typeTrue
	if !value {
		typeValue = typeFalse
	}
	if err = this.buffer.EncodeTypeField(typeValue); err != nil {
		return
	}
	if err = this.rules.AddBool(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) PositiveInt(value uint64) (err error) {
	switch {
	case uintFitsInSmallint(value):
		if err = this.buffer.EncodeTypeField(typeField(value)); err != nil {
			return
		}
	case fitsInUint8(value):
		if err = this.buffer.EncodeUint8(typePosInt8, uint8(value)); err != nil {
			return
		}
	case fitsInUint16(value):
		if err = this.buffer.EncodeUint16(typePosInt16, uint16(value)); err != nil {
			return
		}
	case fitsInUint21(value):
		if err = this.buffer.EncodeUint(typePosInt, value); err != nil {
			return
		}
	case fitsInUint32(value):
		if err = this.buffer.EncodeUint32(typePosInt32, uint32(value)); err != nil {
			return
		}
	case fitsInUint49(value):
		if err = this.buffer.EncodeUint(typePosInt, value); err != nil {
			return
		}
	default:
		if err = this.buffer.EncodeUint64(typePosInt64, value); err != nil {
			return
		}
	}
	if err = this.rules.AddPositiveInt(value); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) NegativeInt(value uint64) (err error) {
	switch {
	case intFitsInSmallint(-int64(value)):
		if err = this.buffer.EncodeTypeField(typeField(-int64(value))); err != nil {
			return
		}
	case fitsInUint8(value):
		if err = this.buffer.EncodeUint8(typeNegInt8, uint8(value)); err != nil {
			return
		}
	case fitsInUint16(value):
		if err = this.buffer.EncodeUint16(typeNegInt16, uint16(value)); err != nil {
			return
		}
	case fitsInUint21(value):
		if err = this.buffer.EncodeUint(typeNegInt, value); err != nil {
			return
		}
	case fitsInUint32(value):
		if err = this.buffer.EncodeUint32(typeNegInt32, uint32(value)); err != nil {
			return
		}
	case fitsInUint49(value):
		if err = this.buffer.EncodeUint(typeNegInt, value); err != nil {
			return
		}
	default:
		if err = this.buffer.EncodeUint64(typeNegInt64, value); err != nil {
			return
		}
	}
	if err = this.rules.AddNegativeInt(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) Int(value int64) error {
	if value >= 0 {
		return this.PositiveInt(uint64(value))
	}
	return this.NegativeInt(uint64(-value))
}

func (this *CBEEncoder) FloatRounded(value float64, significantDigits int) (err error) {
	if significantDigits < 1 || significantDigits > 15 {
		return this.Float(value)
	}

	if err = this.buffer.EncodeFloat(value, significantDigits); err != nil {
		return
	}
	if err = this.rules.AddFloat(value); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) Float(value float64) (err error) {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		if err = this.buffer.EncodeFloat32(asfloat32); err != nil {
			return
		}
	} else {
		if err = this.buffer.EncodeFloat64(value); err != nil {
			return
		}
	}
	if err = this.rules.AddFloat(value); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) CompactTime(value *compact_time.Time) (err error) {
	if err = this.buffer.EncodeCompactTime(value); err != nil {
		return
	}
	if err = this.rules.AddTime(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) Time(value time.Time) (err error) {
	if err = this.buffer.EncodeTime(value); err != nil {
		return
	}
	if err = this.rules.AddTime(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) BeginMarker() (err error) {
	if err = this.buffer.EncodeTypeField(typeMarker); err != nil {
		return
	}
	if err = this.rules.BeginMarker(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) BeginReference() (err error) {
	if err = this.buffer.EncodeTypeField(typeReference); err != nil {
		return
	}
	if err = this.rules.BeginReference(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) EndContainer() (err error) {
	if err = this.buffer.EncodeTypeField(typeEndContainer); err != nil {
		return
	}
	if err = this.rules.EndContainer(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) BeginList() (err error) {
	if err = this.buffer.EncodeTypeField(typeList); err != nil {
		return
	}
	if err = this.rules.BeginList(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Begin a map. Any subsequent objects added are assumed to alternate
// between key and value entries in the map, until EndContainer() is called.
func (this *CBEEncoder) BeginMap() (err error) {
	if err = this.buffer.EncodeTypeField(typeMap); err != nil {
		return
	}
	if err = this.rules.BeginMap(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) BeginMarkup() (err error) {
	if err = this.buffer.EncodeTypeField(typeMarkup); err != nil {
		return
	}
	if err = this.rules.BeginMarkup(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Begin a metadata map. Any subsequent objects added are assumed to alternate
// between key and value entries in the map, until EndContainer() is called.
func (this *CBEEncoder) BeginMetadata() (err error) {
	if err = this.buffer.EncodeTypeField(typeMetadata); err != nil {
		return
	}
	if err = this.rules.BeginMetadata(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

func (this *CBEEncoder) BeginComment() (err error) {
	if err = this.buffer.EncodeTypeField(typeComment); err != nil {
		return
	}
	if err = this.rules.BeginComment(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Begin an array chunk. Encoder expects subsequent calls to ArrayData to
// provide a total of exactly the length provided here.
func (this *CBEEncoder) BeginChunk(length uint64, isFinalChunk bool) (err error) {
	if err = this.beginChunk(length, isFinalChunk); err != nil {
		return
	}
	if err = this.rules.BeginArrayChunk(length, isFinalChunk); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Add binary data to fill in the currently open array chunk.
func (this *CBEEncoder) ArrayData(value []byte) (bytesEncoded int, err error) {
	if bytesEncoded, err = this.arrayData(value); err != nil {
		return
	}
	if err = this.rules.AddArrayData(value); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Begin a byte array. Encoder expects subsequent calls to BeginChunk.
func (this *CBEEncoder) BeginBytes() (err error) {
	if err = this.beginBytes(); err != nil {
		return
	}
	if err = this.rules.BeginBytes(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Convenience function to completely fill a byte array in one call.
// This will fail if there's not enough room in the buffer to completely encode
// the byte array.
func (this *CBEEncoder) Bytes(value []byte) (err error) {
	length := len(value)
	isFinalChunk := true

	if err = this.beginBytes(); err != nil {
		return
	}
	if err = this.arrayEntireData(value); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.BeginBytes(); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.BeginArrayChunk(uint64(length), isFinalChunk); err != nil {
		this.buffer.Rollback()
		return
	}
	if length > 0 {
		if err = this.rules.AddArrayData(value); err != nil {
			this.buffer.Rollback()
			return
		}
	}
	this.buffer.Commit()
	return
}

// Begin a string. Encoder expects subsequent calls to BeginChunk.
func (this *CBEEncoder) BeginString() (err error) {
	if err = this.beginString(); err != nil {
		return
	}
	if err = this.rules.BeginString(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Convenience function to completely fill a string in one call.
// This will fail if there's not enough room in the buffer to completely encode
// the string.
func (this *CBEEncoder) String(value string) (err error) {
	asBytes := []byte(value)
	length := len(asBytes)
	isFinalChunk := true

	if length <= 15 {
		if err = this.smallEntireString(value); err != nil {
			this.buffer.Rollback()
			return
		}
	} else {
		if err = this.largeEntireString(value); err != nil {
			this.buffer.Rollback()
			return
		}
	}
	if err = this.rules.BeginString(); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.BeginArrayChunk(uint64(length), isFinalChunk); err != nil {
		this.buffer.Rollback()
		return
	}
	if length > 0 {
		if err = this.rules.AddArrayData(asBytes); err != nil {
			this.buffer.Rollback()
			return
		}
	}
	this.buffer.Commit()
	return
}

// Begin a URI. Encoder expects subsequent calls to BeginChunk.
func (this *CBEEncoder) BeginURI() (err error) {
	if err = this.beginURI(); err != nil {
		return
	}
	if err = this.rules.BeginURI(); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// Convenience function to completely fill a string in one call.
// This will fail if there's not enough room in the buffer to completely encode
// the URI.
func (this *CBEEncoder) URI(value *url.URL) (err error) {
	asBytes := []byte(value.String())
	length := len(asBytes)
	isFinalChunk := true

	if err = this.beginURI(); err != nil {
		return
	}
	if err = this.arrayEntireData(asBytes); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.BeginURI(); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.BeginArrayChunk(uint64(length), isFinalChunk); err != nil {
		this.buffer.Rollback()
		return
	}
	if err = this.rules.AddArrayData(asBytes); err != nil {
		this.buffer.Rollback()
		return
	}
	this.buffer.Commit()
	return
}

// End the current document.
// Normally, the document will automatically end, except in two cases:
// - You're using inline containers
// - You're creating an empty document
// In these cases, you should manually call End() to make sure there are no errors.
// Note: This method is idempotent, so you can safely call it even if you don't
// need to.
func (this *CBEEncoder) End() (err error) {
	if this.inlineContainerType != InlineContainerTypeNone {
		if err = this.rules.EndContainer(); err != nil {
			return
		}
	}
	return this.rules.EndDocument()
}

func (this *CBEEncoder) EncodedBytes() []byte {
	return this.buffer.EncodedBytes()
}
