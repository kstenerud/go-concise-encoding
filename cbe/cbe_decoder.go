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
	"reflect"

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
			_this.eventReceiver.OnString([]byte{})
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			_this.eventReceiver.OnString(_this.decodeSmallString(int(cbeType - cbeTypeString0)))
		case cbeTypeString:
			_this.eventReceiver.OnString(_this.decodeArray())
		case cbeTypeVerbatimString:
			_this.eventReceiver.OnVerbatimString(_this.decodeArray())
		case cbeTypeURI:
			_this.eventReceiver.OnURI(_this.decodeArray())
		case cbeTypeCustomBinary:
			_this.eventReceiver.OnCustomBinary(_this.decodeArray())
		case cbeTypeCustomText:
			_this.eventReceiver.OnCustomText(_this.decodeArray())
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

func (_this *Decoder) decodeTypedArray() {
	cbeType := _this.buffer.DecodeType()
	switch cbeType {
	case cbeTypeTrue:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(true), _this.decodeArray())
	case cbeTypePosInt8:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(uint8(0)), _this.decodeArray())
	case cbeTypePosInt16:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(uint16(0)), _this.decodeArray())
	case cbeTypePosInt32:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(uint32(0)), _this.decodeArray())
	case cbeTypePosInt64:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(uint64(0)), _this.decodeArray())
	case cbeTypeNegInt8:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(int8(0)), _this.decodeArray())
	case cbeTypeNegInt16:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(int16(0)), _this.decodeArray())
	case cbeTypeNegInt32:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(int32(0)), _this.decodeArray())
	case cbeTypeNegInt64:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(int64(0)), _this.decodeArray())
	case cbeTypeFloat16:
		panic("TODO: Float16 array support")
	case cbeTypeFloat32:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(float32(0)), _this.decodeArray())
	case cbeTypeFloat64:
		_this.eventReceiver.OnTypedArray(reflect.TypeOf(float64(0)), _this.decodeArray())
	case cbeTypeUUID:
		panic("TODO: UUID array support")
	default:
		panic(fmt.Errorf("0x%x: Unsupported typed array type", cbeType))
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

func (_this *Decoder) decodeUnichunkArray(length uint64) []byte {
	validateLength(length)
	if length == 0 {
		return []byte{}
	}
	return _this.buffer.DecodeBytes(int(length))
}

func (_this *Decoder) decodeMultichunkArray(initialLength uint64) []byte {
	length := initialLength
	moreChunksFollow := true
	var bytes []byte
	for {
		validateLength(length)
		// TODO: array chunking instead of building a big slice
		// _this.nextReceiver.OnArrayChunk(length, moreChunksFollow)
		nextBytes := _this.buffer.DecodeBytes(int(length))
		// _this.nextReceiver.OnArrayData(nextBytes)
		bytes = append(bytes, nextBytes...)
		if !moreChunksFollow {
			return bytes
		}
		length, moreChunksFollow = _this.buffer.DecodeChunkHeader()
	}
}

func (_this *Decoder) decodeArray() []byte {
	length, moreChunksFollow := _this.buffer.DecodeChunkHeader()
	if moreChunksFollow {
		return _this.decodeMultichunkArray(length)
	}

	return _this.decodeUnichunkArray(length)
}
