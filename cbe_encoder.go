package concise_encoding

import (
	"math"
	"time"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-vlq"
)

type CBEEncoder struct {
	buff buffer
}

func NewCBEEncoder() *CBEEncoder {
	return &CBEEncoder{}
}

func (this *CBEEncoder) Document() []byte {
	return this.buff.bytes
}

func (this *CBEEncoder) OnPadding(count int) {
	dst := this.buff.Allocate(count)
	for i := 0; i < count; i++ {
		dst[i] = byte(cbeTypePadding)
	}
}

func (this *CBEEncoder) OnVersion(version uint64) {
	this.encodeRVLQ(version)
}

func (this *CBEEncoder) OnNil() {
	this.encodeTypeOnly(cbeTypeNil)
}

func (this *CBEEncoder) OnBool(value bool) {
	if value {
		this.OnTrue()
	} else {
		this.OnFalse()
	}
}

func (this *CBEEncoder) OnTrue() {
	this.encodeTypeOnly(cbeTypeTrue)
}

func (this *CBEEncoder) OnFalse() {
	this.encodeTypeOnly(cbeTypeFalse)
}

func (this *CBEEncoder) OnInt(value int64) {
	if value >= 0 {
		this.OnPositiveInt(uint64(value))
	} else {
		this.OnNegativeInt(uint64(-value))
	}
}

func (this *CBEEncoder) OnPositiveInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		this.encodeTypeOnly(cbeTypeField(value))
	case fitsInUint8(value):
		this.encodeTyped8Bits(cbeTypePosInt8, uint8(value))
	case fitsInUint16(value):
		this.encodeTyped16Bits(cbeTypePosInt16, uint16(value))
	case fitsInUint21(value):
		this.encodeUint(cbeTypePosInt, value)
	case fitsInUint32(value):
		this.encodeTyped32Bits(cbeTypePosInt32, uint32(value))
	case fitsInUint49(value):
		this.encodeUint(cbeTypePosInt, value)
	default:
		this.encodeTyped64Bits(cbeTypePosInt64, value)
	}
}

func (this *CBEEncoder) OnNegativeInt(value uint64) {
	switch {
	case fitsInSmallint(value):
		// Note: Must encode smallint using signed value
		this.encodeTypeOnly(cbeTypeField(-int64(value)))
	case fitsInUint8(value):
		this.encodeTyped8Bits(cbeTypeNegInt8, uint8(value))
	case fitsInUint16(value):
		this.encodeTyped16Bits(cbeTypeNegInt16, uint16(value))
	case fitsInUint21(value):
		this.encodeUint(cbeTypeNegInt, value)
	case fitsInUint32(value):
		this.encodeTyped32Bits(cbeTypeNegInt32, uint32(value))
	case fitsInUint49(value):
		this.encodeUint(cbeTypeNegInt, value)
	default:
		this.encodeTyped64Bits(cbeTypeNegInt64, value)
	}
}

func (this *CBEEncoder) OnFloat(value float64) {
	asfloat32 := float32(value)
	if float64(asfloat32) == value {
		this.encodeTyped32Bits(cbeTypeFloat32, math.Float32bits(asfloat32))
		// this.encodeFloat32(asfloat32)
	} else {
		this.encodeTyped64Bits(cbeTypeFloat64, math.Float64bits(value))
		// this.encodeFloat64(value)
	}
}

func (this *CBEEncoder) OnNan() {
	this.encodeDecimalFloat(math.NaN(), 0)
}

func (this *CBEEncoder) OnUUID(value []byte) {
	dst := this.buff.Allocate(17)
	dst[0] = byte(cbeTypeUUID)
	dst = dst[1:]
	copy(dst, value[:])
}

func (this *CBEEncoder) OnComplex(value complex128) {
	panic("TODO: complex. Requires custom type")
}

func (this *CBEEncoder) OnTime(value time.Time) {
	this.OnCompactTime(compact_time.AsCompactTime(value))
}

func (this *CBEEncoder) OnCompactTime(value *compact_time.Time) {
	var timeType cbeTypeField
	switch value.TimeIs {
	case compact_time.TypeDate:
		timeType = cbeTypeDate
	case compact_time.TypeTime:
		timeType = cbeTypeTime
	case compact_time.TypeTimestamp:
		timeType = cbeTypeTimestamp
	}
	dst := this.buff.Allocate(compact_time.MaxEncodeLength + 1)
	dst[0] = byte(timeType)
	dst = dst[1:]
	byteCount, _ := compact_time.Encode(value, dst)
	this.buff.TrimUnused(byteCount + 1)
}

func (this *CBEEncoder) OnBytes(value []byte) {
	this.encodeTypedBytes(cbeTypeBytes, value)
}

func (this *CBEEncoder) OnURI(value string) {
	this.encodeTypedBytes(cbeTypeURI, []byte(value))
}

func (this *CBEEncoder) OnString(value string) {
	bytes := []byte(value)
	stringLength := len(bytes)

	if stringLength > maxSmallStringLength {
		this.encodeTypedBytes(cbeTypeString, bytes)
		return
	}

	dst := this.buff.Allocate(stringLength + 1)
	dst[0] = byte(cbeTypeString0 + cbeTypeField(stringLength))
	dst = dst[1:]
	copy(dst, bytes)
}

func (this *CBEEncoder) OnCustom(value []byte) {
	this.encodeTypedBytes(cbeTypeCustom, value)
}

func (this *CBEEncoder) OnBytesBegin() {
	this.encodeTypeOnly(cbeTypeBytes)
}

func (this *CBEEncoder) OnStringBegin() {
	this.encodeTypeOnly(cbeTypeString)
}

func (this *CBEEncoder) OnURIBegin() {
	this.encodeTypeOnly(cbeTypeURI)
}

func (this *CBEEncoder) OnCustomBegin() {
	this.encodeTypeOnly(cbeTypeCustom)
}

func (this *CBEEncoder) OnArrayChunk(length uint64, isFinalChunk bool) {
	continuationBit := uint64(0)
	if isFinalChunk {
		continuationBit = 1
	}
	this.encodeRVLQ((uint64(length) << 1) | continuationBit)
}

func (this *CBEEncoder) OnArrayData(data []byte) {
	dst := this.buff.Allocate(len(data))
	copy(dst, data)
}

func (this *CBEEncoder) OnList() {
	this.encodeTypeOnly(cbeTypeList)
}

func (this *CBEEncoder) OnMap() {
	this.encodeTypeOnly(cbeTypeMap)
}

func (this *CBEEncoder) OnMarkup() {
	this.encodeTypeOnly(cbeTypeMarkup)
}

func (this *CBEEncoder) OnMetadata() {
	this.encodeTypeOnly(cbeTypeMetadata)
}

func (this *CBEEncoder) OnComment() {
	this.encodeTypeOnly(cbeTypeComment)
}

func (this *CBEEncoder) OnEnd() {
	this.encodeTypeOnly(cbeTypeEndContainer)
}

func (this *CBEEncoder) OnMarker() {
	this.encodeTypeOnly(cbeTypeMarker)
}

func (this *CBEEncoder) OnReference() {
	this.encodeTypeOnly(cbeTypeReference)
}

func (this *CBEEncoder) OnEndDocument() {
}

const (
	minBufferCap         = 64
	maxSmallStringLength = 15

	bitMask21 = uint64((1 << 21) - 1)
	bitMask49 = uint64((1 << 49) - 1)
)

func fitsInSmallint(value uint64) bool {
	return value <= uint64(cbeSmallIntMax)
}

func fitsInUint8(value uint64) bool {
	return value <= math.MaxUint8
}

func fitsInUint16(value uint64) bool {
	return value <= math.MaxUint16
}

func fitsInUint21(value uint64) bool {
	return value == (value & bitMask21)
}

func fitsInUint32(value uint64) bool {
	return value <= math.MaxUint32
}

func fitsInUint49(value uint64) bool {
	return value == (value & bitMask49)
}

func (this *CBEEncoder) encodeVersion(version uint64) {
	this.encodeRVLQ(version)
}

func (this *CBEEncoder) encodeTypeOnly(value cbeTypeField) {
	this.buff.Allocate(1)[0] = byte(value)
}

func (this *CBEEncoder) encodeRVLQ(valueIn uint64) {
	dst := this.buff.Allocate(vlq.MaxEncodeLength)
	value := vlq.Rvlq(valueIn)
	byteCount, _ := value.EncodeTo(dst)
	this.buff.TrimUnused(byteCount)
}

func (this *CBEEncoder) encodeTypedRVLQ(cbeType cbeTypeField, valueIn uint64) {
	dst := this.buff.Allocate(vlq.MaxEncodeLength + 1)
	dst[0] = byte(cbeType)
	dst = dst[1:]
	value := vlq.Rvlq(valueIn)
	byteCount, _ := value.EncodeTo(dst)
	this.buff.TrimUnused(byteCount + 1)
}

func (this *CBEEncoder) encodeTypedBytes(cbeType cbeTypeField, bytes []byte) {
	bytesLength := len(bytes)
	dst := this.buff.Allocate(vlq.MaxEncodeLength + bytesLength + 1)
	dst[0] = byte(cbeType)
	dst = dst[1:]
	lengthField := vlq.Rvlq((bytesLength << 1) | 1)
	byteCount, _ := lengthField.EncodeTo(dst)
	dst = dst[byteCount:]
	copy(dst, bytes)
	this.buff.TrimUnused(byteCount + bytesLength + 1)
}

func (this *CBEEncoder) encodeTyped8Bits(typeValue cbeTypeField, value byte) {
	dst := this.buff.Allocate(2)
	dst[0] = byte(typeValue)
	dst[1] = value
}

func (this *CBEEncoder) encodeTyped16Bits(typeValue cbeTypeField, value uint16) {
	dst := this.buff.Allocate(3)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
}

func (this *CBEEncoder) encodeTyped32Bits(typeValue cbeTypeField, value uint32) {
	dst := this.buff.Allocate(5)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
}

func (this *CBEEncoder) encodeTyped64Bits(typeValue cbeTypeField, value uint64) {
	dst := this.buff.Allocate(9)
	dst[0] = byte(typeValue)
	dst[1] = byte(value)
	dst[2] = byte(value >> 8)
	dst[3] = byte(value >> 16)
	dst[4] = byte(value >> 24)
	dst[5] = byte(value >> 32)
	dst[6] = byte(value >> 40)
	dst[7] = byte(value >> 48)
	dst[8] = byte(value >> 56)
}

func (this *CBEEncoder) encodeUint(typeValue cbeTypeField, value uint64) {
	this.encodeTypedRVLQ(typeValue, value)
}

// TODO
// func (this *CBEEncoder) encodeFloat32(value float32) {
// 	this.encodeTyped32Bits(cbeTypeFloat32, math.Float32bits(value))
// }

// func (this *CBEEncoder) encodeFloat64(value float64) {
// 	this.encodeTyped64Bits(cbeTypeFloat64, math.Float64bits(value))
// }

func (this *CBEEncoder) encodeDecimalFloat(value float64, significantDigits int) {
	dst := this.buff.Allocate(compact_float.MaxEncodeLength + 1)
	dst[0] = byte(cbeTypeDecimal)
	dst = dst[1:]
	byteCount, _ := compact_float.Encode(value, significantDigits, dst)
	this.buff.TrimUnused(byteCount + 1)
}
