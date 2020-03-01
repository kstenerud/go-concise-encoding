package concise_encoding

import (
	"fmt"
)

func CBEDecode(document []byte, eventHandler ConciseEncodingEventHandler, shouldZeroCopy bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	decoder := NewCBEDecoder([]byte(document), eventHandler, shouldZeroCopy)
	decoder.Decode()
	return
}

type CBEDecoder struct {
	buffer         cbeDecodeBuffer
	shouldZeroCopy bool
	nextHandler    ConciseEncodingEventHandler
}

func NewCBEDecoder(document []byte, nextHandler ConciseEncodingEventHandler, shouldZeroCopy bool) *CBEDecoder {
	this := &CBEDecoder{}
	this.Init(document, nextHandler, shouldZeroCopy)
	return this
}

func (this *CBEDecoder) Init(document []byte, nextHandler ConciseEncodingEventHandler, shouldZeroCopy bool) {
	this.buffer.Init(document)
	this.shouldZeroCopy = shouldZeroCopy
	this.nextHandler = nextHandler
}

func (this *CBEDecoder) Decode() {
	this.nextHandler.OnVersion(this.buffer.DecodeVersion())

	for this.buffer.HasUnreadData() {
		cbeType := this.buffer.DecodeType()
		switch cbeType {
		case cbeTypeDecimal:
			this.nextHandler.OnFloat(this.buffer.DecodeFloat())
		case cbeTypePosInt:
			this.nextHandler.OnPositiveInt(this.buffer.DecodeUint())
		case cbeTypeNegInt:
			this.nextHandler.OnNegativeInt(this.buffer.DecodeUint())
		case cbeTypePosInt8:
			this.nextHandler.OnPositiveInt(uint64(this.buffer.DecodeUint8()))
		case cbeTypeNegInt8:
			this.nextHandler.OnNegativeInt(uint64(this.buffer.DecodeUint8()))
		case cbeTypePosInt16:
			this.nextHandler.OnPositiveInt(uint64(this.buffer.DecodeUint16()))
		case cbeTypeNegInt16:
			this.nextHandler.OnNegativeInt(uint64(this.buffer.DecodeUint16()))
		case cbeTypePosInt32:
			this.nextHandler.OnPositiveInt(uint64(this.buffer.DecodeUint32()))
		case cbeTypeNegInt32:
			this.nextHandler.OnNegativeInt(uint64(this.buffer.DecodeUint32()))
		case cbeTypePosInt64:
			this.nextHandler.OnPositiveInt(this.buffer.DecodeUint64())
		case cbeTypeNegInt64:
			this.nextHandler.OnNegativeInt(this.buffer.DecodeUint64())
		case cbeTypeFloat32:
			this.nextHandler.OnFloat(float64(this.buffer.DecodeFloat32()))
		case cbeTypeFloat64:
			this.nextHandler.OnFloat(this.buffer.DecodeFloat64())
		case cbeTypeUUID:
			this.nextHandler.OnUUID(this.buffer.DecodeBytes(16))
		case cbeTypeComment:
			this.nextHandler.OnComment()
		case cbeTypeMetadata:
			this.nextHandler.OnMetadata()
		case cbeTypeMarkup:
			this.nextHandler.OnMarkup()
		case cbeTypeMap:
			this.nextHandler.OnMap()
		case cbeTypeList:
			this.nextHandler.OnList()
		case cbeTypeEndContainer:
			this.nextHandler.OnEnd()
		case cbeTypeFalse:
			this.nextHandler.OnFalse()
		case cbeTypeTrue:
			this.nextHandler.OnTrue()
		case cbeTypeNil:
			this.nextHandler.OnNil()
		case cbeTypePadding:
			this.nextHandler.OnPadding(1)
		case cbeTypeString0:
			this.nextHandler.OnString("")
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			this.nextHandler.OnString(this.decodeSmallString(int(cbeType - cbeTypeString0)))
		case cbeTypeString:
			this.nextHandler.OnString(string(this.decodeArray()))
		case cbeTypeBytes:
			this.nextHandler.OnBytes(this.decodeArray())
		case cbeTypeCustom:
			this.nextHandler.OnCustom(this.decodeArray())
		case cbeTypeURI:
			this.nextHandler.OnURI(string(this.decodeArray()))
		case cbeTypeMarker:
			this.nextHandler.OnMarker()
		case cbeTypeReference:
			this.nextHandler.OnReference()
		case cbeTypeDate:
			this.nextHandler.OnCompactTime(this.buffer.DecodeDate())
		case cbeTypeTime:
			this.nextHandler.OnCompactTime(this.buffer.DecodeTime())
		case cbeTypeTimestamp:
			this.nextHandler.OnCompactTime(this.buffer.DecodeTimestamp())
		default:
			asSmallInt := int64(int8(cbeType))
			if asSmallInt < cbeSmallIntMin || asSmallInt > cbeSmallIntMax {
				panic(fmt.Errorf("Unknown type code 0x%02x", cbeType))
			}
			this.nextHandler.OnInt(asSmallInt)
		}
	}

	this.nextHandler.OnEndDocument()
	return
}

func (this *CBEDecoder) possiblyZeroCopy(bytes []byte) []byte {
	if this.shouldZeroCopy {
		return bytes
	}
	bytesCopy := make([]byte, len(bytes), len(bytes))
	copy(bytesCopy, bytes)
	return bytesCopy
}

func (this *CBEDecoder) decodeSmallString(length int) string {
	value := string(this.possiblyZeroCopy(this.buffer.DecodeBytes(length)))
	return value
}

func validateLength(length uint64) {
	const maxDefaultInt = uint64((^uint(0)) >> 1)
	if length > maxDefaultInt {
		panic(fmt.Errorf("%v > max int value (%v)", length, maxDefaultInt))
	}
}

func (this *CBEDecoder) decodeUnichunkArray(length uint64) []byte {
	validateLength(length)
	// TODO:
	// this.nextHandler.OnArrayChunk(length, true)
	if length == 0 {
		return []byte{}
	}
	bytes := this.possiblyZeroCopy(this.buffer.DecodeBytes(int(length)))
	// this.nextHandler.OnArrayData(bytes)
	return bytes
}

func (this *CBEDecoder) decodeMultichunkArray(initialLength uint64) []byte {
	length := initialLength
	isFinalChunk := false
	bytes := []byte{}
	for {
		validateLength(length)
		// TODO:
		// this.nextHandler.OnArrayChunk(length, isFinalChunk)
		nextBytes := this.buffer.DecodeBytes(int(length))
		// this.nextHandler.OnArrayData(nextBytes)
		bytes = append(bytes, nextBytes...)
		if isFinalChunk {
			return bytes
		}
		length, isFinalChunk = this.buffer.DecodeChunkHeader()
	}
}

func (this *CBEDecoder) decodeArray() []byte {
	length, isFinalChunk := this.buffer.DecodeChunkHeader()
	if isFinalChunk {
		return this.decodeUnichunkArray(length)
	}

	return this.decodeMultichunkArray(length)
}
