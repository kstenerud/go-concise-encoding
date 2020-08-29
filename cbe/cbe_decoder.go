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

package cbe

import (
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decodes CBE documents.
type Decoder struct {
	buffer        ReadBuffer
	eventReceiver events.DataEventReceiver
	opts          options.CBEDecoderOptions
}

// Create a new CBE decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func NewDecoder(opts *options.CBEDecoderOptions) *Decoder {
	_this := &Decoder{}
	_this.Init(opts)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func (_this *Decoder) Init(opts *options.CBEDecoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
}

func (_this *Decoder) reset() {
	_this.buffer.Reset()
	_this.eventReceiver = nil
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		_this.reset()
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				err = r.(error)
			}
		}
	}()

	_this.buffer.Init(reader, _this.opts.BufferSize, chooseLowWater(_this.opts.BufferSize))
	_this.eventReceiver = eventReceiver

	_this.eventReceiver.OnBeginDocument()

	_this.buffer.RefillIfNecessary()

	switch _this.opts.ImpliedStructure {
	case options.ImpliedStructureVersion:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
	case options.ImpliedStructureList:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
		_this.eventReceiver.OnList()
	case options.ImpliedStructureMap:
		_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
		_this.eventReceiver.OnMap()
	default:
		_this.eventReceiver.OnVersion(_this.buffer.DecodeVersion())
	}

	for _this.buffer.HasUnreadData() {
		_this.buffer.RefillIfNecessary()
		cbeType := _this.buffer.DecodeType()
		switch cbeType {
		case cbeTypeDecimal:
			value, bigValue := _this.buffer.DecodeDecimalFloat()
			if bigValue != nil {
				_this.eventReceiver.OnBigDecimalFloat(bigValue)
			} else {
				_this.eventReceiver.OnDecimalFloat(value)
			}
		case cbeTypePosInt:
			asUint, asBig := _this.buffer.DecodeUint()
			if asBig != nil {
				_this.eventReceiver.OnBigInt(asBig)
			} else {
				_this.eventReceiver.OnPositiveInt(asUint)
			}
		case cbeTypeNegInt:
			asUint, asBig := _this.buffer.DecodeUint()
			if asBig != nil {
				_this.eventReceiver.OnBigInt(asBig.Neg(asBig))
			} else {
				_this.eventReceiver.OnNegativeInt(asUint)
			}
		case cbeTypePosInt8:
			_this.eventReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint8()))
		case cbeTypeNegInt8:
			_this.eventReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint8()))
		case cbeTypePosInt16:
			_this.eventReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint16()))
		case cbeTypeNegInt16:
			_this.eventReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint16()))
		case cbeTypePosInt32:
			_this.eventReceiver.OnPositiveInt(uint64(_this.buffer.DecodeUint32()))
		case cbeTypeNegInt32:
			_this.eventReceiver.OnNegativeInt(uint64(_this.buffer.DecodeUint32()))
		case cbeTypePosInt64:
			_this.eventReceiver.OnPositiveInt(_this.buffer.DecodeUint64())
		case cbeTypeNegInt64:
			_this.eventReceiver.OnNegativeInt(_this.buffer.DecodeUint64())
		case cbeTypeFloat16:
			_this.eventReceiver.OnFloat(float64(_this.buffer.DecodeFloat16()))
		case cbeTypeFloat32:
			_this.eventReceiver.OnFloat(float64(_this.buffer.DecodeFloat32()))
		case cbeTypeFloat64:
			_this.eventReceiver.OnFloat(_this.buffer.DecodeFloat64())
		case cbeTypeUUID:
			_this.eventReceiver.OnUUID(_this.buffer.DecodeBytes(16))
		case cbeTypeComment:
			_this.eventReceiver.OnComment()
		case cbeTypeMetadata:
			_this.eventReceiver.OnMetadata()
		case cbeTypeMarkup:
			_this.eventReceiver.OnMarkup()
		case cbeTypeMap:
			_this.eventReceiver.OnMap()
		case cbeTypeList:
			_this.eventReceiver.OnList()
		case cbeTypeEndContainer:
			_this.eventReceiver.OnEnd()
		case cbeTypeFalse:
			_this.eventReceiver.OnFalse()
		case cbeTypeTrue:
			_this.eventReceiver.OnTrue()
		case cbeTypeNil:
			_this.eventReceiver.OnNil()
		case cbeTypePadding:
			_this.eventReceiver.OnPadding(1)
		case cbeTypeString0:
			_this.eventReceiver.OnArray(events.ArrayTypeString, 0, []byte{})
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			length := int(cbeType - cbeTypeString0)
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(length), _this.decodeSmallString(length))
		case cbeTypeString:
			bytes := _this.decodeArray(8)
			_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
		case cbeTypeVerbatimString:
			bytes := _this.decodeArray(8)
			_this.eventReceiver.OnArray(events.ArrayTypeVerbatimString, uint64(len(bytes)), bytes)
		case cbeTypeURI:
			bytes := _this.decodeArray(8)
			_this.eventReceiver.OnArray(events.ArrayTypeURI, uint64(len(bytes)), bytes)
		case cbeTypeCustomBinary:
			bytes := _this.decodeArray(8)
			_this.eventReceiver.OnArray(events.ArrayTypeCustomBinary, uint64(len(bytes)), bytes)
		case cbeTypeCustomText:
			bytes := _this.decodeArray(8)
			_this.eventReceiver.OnArray(events.ArrayTypeCustomText, uint64(len(bytes)), bytes)
		case cbeTypeArray:
			_this.decodeTypedArray()
		case cbeTypeMarker:
			_this.eventReceiver.OnMarker()
		case cbeTypeReference:
			_this.eventReceiver.OnReference()
		case cbeTypeDate:
			_this.eventReceiver.OnCompactTime(_this.buffer.DecodeDate())
		case cbeTypeTime:
			_this.eventReceiver.OnCompactTime(_this.buffer.DecodeTime())
		case cbeTypeTimestamp:
			_this.eventReceiver.OnCompactTime(_this.buffer.DecodeTimestamp())
		default:
			asSmallInt := int64(int8(cbeType))
			if asSmallInt < cbeSmallIntMin || asSmallInt > cbeSmallIntMax {
				panic(fmt.Errorf("unknown type code 0x%02x", cbeType))
			}
			_this.eventReceiver.OnInt(asSmallInt)
		}
	}

	switch _this.opts.ImpliedStructure {
	case options.ImpliedStructureList, options.ImpliedStructureMap:
		_this.eventReceiver.OnEnd()
	}

	_this.eventReceiver.OnEndDocument()
	return
}

// ============================================================================

// Internal

var typeToArrayType [256]uint8

const arrayTypeInvalid = events.ArrayType(0xff)

func init() {
	for i := 0; i < len(typeToArrayType); i++ {
		typeToArrayType[i] = uint8(arrayTypeInvalid)
	}
	typeToArrayType[uint8(cbeTypeTrue)] = uint8(events.ArrayTypeBoolean)
	typeToArrayType[uint8(cbeTypePosInt8)] = uint8(events.ArrayTypeUint8)
	typeToArrayType[uint8(cbeTypePosInt16)] = uint8(events.ArrayTypeUint16)
	typeToArrayType[uint8(cbeTypePosInt32)] = uint8(events.ArrayTypeUint32)
	typeToArrayType[uint8(cbeTypePosInt64)] = uint8(events.ArrayTypeUint64)
	typeToArrayType[uint8(cbeTypeNegInt8)] = uint8(events.ArrayTypeInt8)
	typeToArrayType[uint8(cbeTypeNegInt16)] = uint8(events.ArrayTypeInt16)
	typeToArrayType[uint8(cbeTypeNegInt32)] = uint8(events.ArrayTypeInt32)
	typeToArrayType[uint8(cbeTypeNegInt64)] = uint8(events.ArrayTypeInt64)
	typeToArrayType[uint8(cbeTypeFloat16)] = uint8(events.ArrayTypeFloat16)
	typeToArrayType[uint8(cbeTypeFloat32)] = uint8(events.ArrayTypeFloat32)
	typeToArrayType[uint8(cbeTypeFloat64)] = uint8(events.ArrayTypeFloat64)
	typeToArrayType[uint8(cbeTypeUUID)] = uint8(events.ArrayTypeUUID)
}

func (_this *Decoder) decodeTypedArray() {
	cbeType := _this.buffer.DecodeType()
	arrayType := events.ArrayType(typeToArrayType[uint8(cbeType)])
	if arrayType == arrayTypeInvalid {
		panic(fmt.Errorf("0x%x: Unsupported typed array type", cbeType))
	}

	elementBitWidth := arrayType.ElementSize()
	elementCount, moreChunksFollow := _this.buffer.DecodeChunkHeader()
	validateLength(elementCount)
	if !moreChunksFollow {
		bytes := _this.decodeUnichunkArray(elementBitWidth, elementCount)
		_this.eventReceiver.OnArray(arrayType, elementCount, bytes)
		return
	}

	for {
		_this.eventReceiver.OnArrayChunk(elementCount, moreChunksFollow)
		byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
		nextBytes := _this.buffer.DecodeBytes(int(byteCount))
		_this.eventReceiver.OnArrayData(nextBytes)
		if !moreChunksFollow {
			return
		}
		elementCount, moreChunksFollow = _this.buffer.DecodeChunkHeader()
		validateLength(elementCount)
	}
}

func chooseLowWater(bufferSize int) int {
	lowWater := bufferSize / 50
	if lowWater < 30 {
		lowWater = 30
	}
	return lowWater
}

func (_this *Decoder) decodeSmallString(length int) []byte {
	return _this.buffer.DecodeBytes(length)
}

func validateLength(length uint64) {
	const maxDefaultInt = uint64((^uint(0)) >> 1)
	if length > maxDefaultInt {
		panic(fmt.Errorf("%v > max int value (%v)", length, maxDefaultInt))
	}
}

func (_this *Decoder) decodeUnichunkArray(elementBitWidth int, elementCount uint64) []byte {
	validateLength(elementCount)
	if elementCount == 0 {
		return []byte{}
	}
	byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
	return _this.buffer.DecodeBytes(int(byteCount))
}

func (_this *Decoder) decodeMultichunkArray(elementBitWidth int, firstChunkElementCount uint64) []byte {
	elementCount := firstChunkElementCount
	moreChunksFollow := true
	var bytes []byte
	for {
		validateLength(elementCount)
		byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
		nextBytes := _this.buffer.DecodeBytes(int(byteCount))
		bytes = append(bytes, nextBytes...)
		if !moreChunksFollow {
			return bytes
		}
		elementCount, moreChunksFollow = _this.buffer.DecodeChunkHeader()
	}
}

func (_this *Decoder) decodeArray(elementBitWidth int) []byte {
	// TODO: Array chunking for string array types
	elementCount, moreChunksFollow := _this.buffer.DecodeChunkHeader()
	if moreChunksFollow {
		return _this.decodeMultichunkArray(elementBitWidth, elementCount)
	}

	return _this.decodeUnichunkArray(elementBitWidth, elementCount)
}
