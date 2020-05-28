// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package concise_encoding

import (
	"fmt"
)

func CBEDecode(document []byte, eventReceiver DataEventReceiver, shouldZeroCopy bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	decoder := NewCBEDecoder([]byte(document), eventReceiver, shouldZeroCopy)
	decoder.Decode()
	return
}

type CBEDecoder struct {
	buffer         cbeDecodeBuffer
	shouldZeroCopy bool
	nextReceiver   DataEventReceiver
}

func NewCBEDecoder(document []byte, nextReceiver DataEventReceiver, shouldZeroCopy bool) *CBEDecoder {
	this := &CBEDecoder{}
	this.Init(document, nextReceiver, shouldZeroCopy)
	return this
}

func (this *CBEDecoder) Init(document []byte, nextReceiver DataEventReceiver, shouldZeroCopy bool) {
	this.buffer.Init(document)
	this.shouldZeroCopy = shouldZeroCopy
	this.nextReceiver = nextReceiver
}

func (this *CBEDecoder) Decode() {
	this.nextReceiver.OnVersion(this.buffer.DecodeVersion())

	for this.buffer.HasUnreadData() {
		cbeType := this.buffer.DecodeType()
		switch cbeType {
		case cbeTypeDecimal:
			value, bigValue := this.buffer.DecodeDecimalFloat()
			if bigValue != nil {
				this.nextReceiver.OnBigDecimalFloat(bigValue)
			} else {
				this.nextReceiver.OnDecimalFloat(value)
			}
		case cbeTypePosInt:
			asUint, asBig := this.buffer.DecodeUint()
			if asBig != nil {
				this.nextReceiver.OnBigInt(asBig)
			} else {
				this.nextReceiver.OnPositiveInt(asUint)
			}
		case cbeTypeNegInt:
			asUint, asBig := this.buffer.DecodeUint()
			if asBig != nil {
				this.nextReceiver.OnBigInt(asBig.Neg(asBig))
			} else {
				this.nextReceiver.OnNegativeInt(asUint)
			}
		case cbeTypePosInt8:
			this.nextReceiver.OnPositiveInt(uint64(this.buffer.DecodeUint8()))
		case cbeTypeNegInt8:
			this.nextReceiver.OnNegativeInt(uint64(this.buffer.DecodeUint8()))
		case cbeTypePosInt16:
			this.nextReceiver.OnPositiveInt(uint64(this.buffer.DecodeUint16()))
		case cbeTypeNegInt16:
			this.nextReceiver.OnNegativeInt(uint64(this.buffer.DecodeUint16()))
		case cbeTypePosInt32:
			this.nextReceiver.OnPositiveInt(uint64(this.buffer.DecodeUint32()))
		case cbeTypeNegInt32:
			this.nextReceiver.OnNegativeInt(uint64(this.buffer.DecodeUint32()))
		case cbeTypePosInt64:
			this.nextReceiver.OnPositiveInt(this.buffer.DecodeUint64())
		case cbeTypeNegInt64:
			this.nextReceiver.OnNegativeInt(this.buffer.DecodeUint64())
		case cbeTypeFloat32:
			this.nextReceiver.OnFloat(float64(this.buffer.DecodeFloat32()))
		case cbeTypeFloat64:
			this.nextReceiver.OnFloat(this.buffer.DecodeFloat64())
		case cbeTypeUUID:
			this.nextReceiver.OnUUID(this.buffer.DecodeBytes(16))
		case cbeTypeComment:
			this.nextReceiver.OnComment()
		case cbeTypeMetadata:
			this.nextReceiver.OnMetadata()
		case cbeTypeMarkup:
			this.nextReceiver.OnMarkup()
		case cbeTypeMap:
			this.nextReceiver.OnMap()
		case cbeTypeList:
			this.nextReceiver.OnList()
		case cbeTypeEndContainer:
			this.nextReceiver.OnEnd()
		case cbeTypeFalse:
			this.nextReceiver.OnFalse()
		case cbeTypeTrue:
			this.nextReceiver.OnTrue()
		case cbeTypeNil:
			this.nextReceiver.OnNil()
		case cbeTypePadding:
			this.nextReceiver.OnPadding(1)
		case cbeTypeString0:
			this.nextReceiver.OnString("")
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			this.nextReceiver.OnString(this.decodeSmallString(int(cbeType - cbeTypeString0)))
		case cbeTypeString:
			this.nextReceiver.OnString(string(this.decodeArray()))
		case cbeTypeBytes:
			this.nextReceiver.OnBytes(this.decodeArray())
		case cbeTypeCustom:
			this.nextReceiver.OnCustom(this.decodeArray())
		case cbeTypeURI:
			this.nextReceiver.OnURI(string(this.decodeArray()))
		case cbeTypeMarker:
			this.nextReceiver.OnMarker()
		case cbeTypeReference:
			this.nextReceiver.OnReference()
		case cbeTypeDate:
			this.nextReceiver.OnCompactTime(this.buffer.DecodeDate())
		case cbeTypeTime:
			this.nextReceiver.OnCompactTime(this.buffer.DecodeTime())
		case cbeTypeTimestamp:
			this.nextReceiver.OnCompactTime(this.buffer.DecodeTimestamp())
		default:
			asSmallInt := int64(int8(cbeType))
			if asSmallInt < cbeSmallIntMin || asSmallInt > cbeSmallIntMax {
				panic(fmt.Errorf("Unknown type code 0x%02x", cbeType))
			}
			this.nextReceiver.OnInt(asSmallInt)
		}
	}

	this.nextReceiver.OnEndDocument()
	return
}

// ============================================================================

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
	// this.nextReceiver.OnArrayChunk(length, true)
	if length == 0 {
		return []byte{}
	}
	bytes := this.possiblyZeroCopy(this.buffer.DecodeBytes(int(length)))
	// this.nextReceiver.OnArrayData(bytes)
	return bytes
}

func (this *CBEDecoder) decodeMultichunkArray(initialLength uint64) []byte {
	length := initialLength
	isFinalChunk := false
	bytes := []byte{}
	for {
		validateLength(length)
		// TODO:
		// this.nextReceiver.OnArrayChunk(length, isFinalChunk)
		nextBytes := this.buffer.DecodeBytes(int(length))
		// this.nextReceiver.OnArrayData(nextBytes)
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
