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

type CBEDecoderOptions struct {
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
}

// Decode a CBE document, sending all data events to the specified event receiver.
func CBEDecode(document []byte, eventReceiver DataEventReceiver, options *CBEDecoderOptions) (err error) {
	defer func() {
		if !DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}
	}()

	decoder := NewCBEDecoder([]byte(document), eventReceiver, options)
	decoder.Decode()
	return
}

// Decodes CBE documents
type CBEDecoder struct {
	buffer       cbeDecodeBuffer
	nextReceiver DataEventReceiver
	options      CBEDecoderOptions
}

func NewCBEDecoder(document []byte, nextReceiver DataEventReceiver, options *CBEDecoderOptions) *CBEDecoder {
	_this := &CBEDecoder{}
	_this.Init(document, nextReceiver, options)
	return _this
}

func (_this *CBEDecoder) Init(document []byte, nextReceiver DataEventReceiver, options *CBEDecoderOptions) {
	_this.buffer.Init(document)
	if options != nil {
		_this.options = *options
	}
	_this.nextReceiver = nextReceiver
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *CBEDecoder) Decode() {
	_this.nextReceiver.OnVersion(_this.buffer.DecodeVersion())

	for _this.buffer.HasUnreadData() {
		cbeType := _this.buffer.DecodeType()
		switch cbeType {
		case cbeTypeDecimal:
			value, bigValue := _this.buffer.DecodeDecimalFloat()
			if bigValue != nil {
				_this.nextReceiver.OnBigDecimalFloat(bigValue)
			} else {
				_this.nextReceiver.OnDecimalFloat(value)
			}
		case cbeTypePosInt:
			asUint, asBig := _this.buffer.DecodeUint()
			if asBig != nil {
				_this.nextReceiver.OnBigInt(asBig)
			} else {
				_this.nextReceiver.OnPositiveInt(asUint)
			}
		case cbeTypeNegInt:
			asUint, asBig := _this.buffer.DecodeUint()
			if asBig != nil {
				_this.nextReceiver.OnBigInt(asBig.Neg(asBig))
			} else {
				_this.nextReceiver.OnNegativeInt(asUint)
			}
		case cbeTypePosInt8:
			_this.nextReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint8()))
		case cbeTypeNegInt8:
			_this.nextReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint8()))
		case cbeTypePosInt16:
			_this.nextReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint16()))
		case cbeTypeNegInt16:
			_this.nextReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint16()))
		case cbeTypePosInt32:
			_this.nextReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint32()))
		case cbeTypeNegInt32:
			_this.nextReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint32()))
		case cbeTypePosInt64:
			_this.nextReceiver.OnPositiveInt(_this.buffer.DecodeUint64())
		case cbeTypeNegInt64:
			_this.nextReceiver.OnNegativeInt(_this.buffer.DecodeUint64())
		case cbeTypeFloat32:
			_this.nextReceiver.OnFloat(float64(_this.buffer.DecodeFloat32()))
		case cbeTypeFloat64:
			_this.nextReceiver.OnFloat(_this.buffer.DecodeFloat64())
		case cbeTypeUUID:
			_this.nextReceiver.OnUUID(_this.buffer.DecodeBytes(16))
		case cbeTypeComment:
			_this.nextReceiver.OnComment()
		case cbeTypeMetadata:
			_this.nextReceiver.OnMetadata()
		case cbeTypeMarkup:
			_this.nextReceiver.OnMarkup()
		case cbeTypeMap:
			_this.nextReceiver.OnMap()
		case cbeTypeList:
			_this.nextReceiver.OnList()
		case cbeTypeEndContainer:
			_this.nextReceiver.OnEnd()
		case cbeTypeFalse:
			_this.nextReceiver.OnFalse()
		case cbeTypeTrue:
			_this.nextReceiver.OnTrue()
		case cbeTypeNil:
			_this.nextReceiver.OnNil()
		case cbeTypePadding:
			_this.nextReceiver.OnPadding(1)
		case cbeTypeString0:
			_this.nextReceiver.OnString("")
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			_this.nextReceiver.OnString(_this.decodeSmallString(int(cbeType - cbeTypeString0)))
		case cbeTypeString:
			_this.nextReceiver.OnString(string(_this.decodeArray()))
		case cbeTypeBytes:
			_this.nextReceiver.OnBytes(_this.decodeArray())
		case cbeTypeCustom:
			_this.nextReceiver.OnCustom(_this.decodeArray())
		case cbeTypeURI:
			_this.nextReceiver.OnURI(string(_this.decodeArray()))
		case cbeTypeMarker:
			_this.nextReceiver.OnMarker()
		case cbeTypeReference:
			_this.nextReceiver.OnReference()
		case cbeTypeDate:
			_this.nextReceiver.OnCompactTime(_this.buffer.DecodeDate())
		case cbeTypeTime:
			_this.nextReceiver.OnCompactTime(_this.buffer.DecodeTime())
		case cbeTypeTimestamp:
			_this.nextReceiver.OnCompactTime(_this.buffer.DecodeTimestamp())
		default:
			asSmallInt := int64(int8(cbeType))
			if asSmallInt < cbeSmallIntMin || asSmallInt > cbeSmallIntMax {
				panic(fmt.Errorf("Unknown type code 0x%02x", cbeType))
			}
			_this.nextReceiver.OnInt(asSmallInt)
		}
	}

	_this.nextReceiver.OnEndDocument()
	return
}

// ============================================================================

func (_this *CBEDecoder) possiblyZeroCopy(bytes []byte) []byte {
	if _this.options.ShouldZeroCopy {
		return bytes
	}
	bytesCopy := make([]byte, len(bytes), len(bytes))
	copy(bytesCopy, bytes)
	return bytesCopy
}

func (_this *CBEDecoder) decodeSmallString(length int) string {
	value := string(_this.possiblyZeroCopy(_this.buffer.DecodeBytes(length)))
	return value
}

func validateLength(length uint64) {
	const maxDefaultInt = uint64((^uint(0)) >> 1)
	if length > maxDefaultInt {
		panic(fmt.Errorf("%v > max int value (%v)", length, maxDefaultInt))
	}
}

func (_this *CBEDecoder) decodeUnichunkArray(length uint64) []byte {
	validateLength(length)
	// TODO:
	// _this.nextReceiver.OnArrayChunk(length, true)
	if length == 0 {
		return []byte{}
	}
	bytes := _this.possiblyZeroCopy(_this.buffer.DecodeBytes(int(length)))
	// _this.nextReceiver.OnArrayData(bytes)
	return bytes
}

func (_this *CBEDecoder) decodeMultichunkArray(initialLength uint64) []byte {
	length := initialLength
	isFinalChunk := false
	bytes := []byte{}
	for {
		validateLength(length)
		// TODO:
		// _this.nextReceiver.OnArrayChunk(length, isFinalChunk)
		nextBytes := _this.buffer.DecodeBytes(int(length))
		// _this.nextReceiver.OnArrayData(nextBytes)
		bytes = append(bytes, nextBytes...)
		if isFinalChunk {
			return bytes
		}
		length, isFinalChunk = _this.buffer.DecodeChunkHeader()
	}
}

func (_this *CBEDecoder) decodeArray() []byte {
	length, isFinalChunk := _this.buffer.DecodeChunkHeader()
	if isFinalChunk {
		return _this.decodeUnichunkArray(length)
	}

	return _this.decodeMultichunkArray(length)
}
