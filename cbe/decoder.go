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
	"bytes"
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decodes CBE documents.
type Decoder struct {
	buffer        Reader
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
	_this.buffer.Init()
}

// Run the complete decode process. The document and data receiver specified
// when initializing the decoder will be used.
func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	_this.buffer.SetReader(reader)
	_this.eventReceiver = eventReceiver

	_this.eventReceiver.OnBeginDocument()

	docHeader := _this.buffer.DecodeUint8()
	if docHeader != cbeDocumentHeader {
		_this.buffer.errorf("First byte of CBE document must be 0x%02x (found 0x%02x)", cbeDocumentHeader, docHeader)
	}
	ver := _this.buffer.DecodeVersion()
	// TODO: Remove this when releasing V1
	if ver == 1 {
		ver = 0
	}
	_this.eventReceiver.OnVersion(ver)

	// TODO: Need to query rules to see when to stop
	for {
		cbeType, isEOF := _this.buffer.DecodeTypeWithEOFCheck()
		if isEOF {
			break
		}
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
			_this.decodeArray(events.ArrayTypeString)
		case cbeTypeRID:
			_this.decodeArray(events.ArrayTypeResourceID)
		case cbeTypeCustomBinary:
			_this.decodeArray(events.ArrayTypeCustomBinary)
		case cbeTypeCustomText:
			_this.decodeArray(events.ArrayTypeCustomText)
		case cbeTypePlane2:
			cbeType := _this.buffer.DecodeType()
			switch cbeType {
			case cbeTypeNA:
				_this.eventReceiver.OnNA()
			default:
				arrayType := cbePlane2TypeToArrayType[cbeType]
				if arrayType == events.ArrayTypeInvalid {
					panic(fmt.Errorf("0x%02x: Unsupported plane 2 type", cbeType))
				}
				_this.decodeArray(arrayType)
			}
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
				panic(fmt.Errorf("0x%02x: Unsupported type", cbeType))
			}
			_this.eventReceiver.OnInt(asSmallInt)
		}
	}

	_this.eventReceiver.OnEndDocument()
	return
}

func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	return _this.Decode(bytes.NewBuffer(document), eventReceiver)
}

// ============================================================================

// Internal

func (_this *Decoder) decodeArray(arrayType events.ArrayType) {
	elementBitWidth := arrayType.ElementSize()
	elementCount, moreChunksFollow := _this.buffer.DecodeArrayChunkHeader()
	validateLength(elementCount)

	if !moreChunksFollow {
		if elementCount == 0 {
			_this.eventReceiver.OnArrayBegin(arrayType)
			_this.eventReceiver.OnArrayChunk(0, false)
			return
		}
		byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
		bytes := _this.buffer.DecodeBytes(int(byteCount))
		_this.eventReceiver.OnArray(arrayType, elementCount, bytes)
		return
	}

	_this.eventReceiver.OnArrayBegin(arrayType)

	for {
		_this.eventReceiver.OnArrayChunk(elementCount, moreChunksFollow)
		byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
		nextBytes := _this.buffer.DecodeBytes(int(byteCount))
		_this.eventReceiver.OnArrayData(nextBytes)
		if !moreChunksFollow {
			return
		}
		elementCount, moreChunksFollow = _this.buffer.DecodeArrayChunkHeader()
		validateLength(elementCount)
	}
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
