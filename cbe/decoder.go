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
	"math"

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

// Decodes CBE documents.
type Decoder struct {
	reader Reader
	config *configuration.Configuration
}

// Create a new CBE decoder.
func NewDecoder(config *configuration.Configuration) *Decoder {
	_this := &Decoder{}
	_this.Init(config)
	return _this
}

// Initialize this decoder.
func (_this *Decoder) Init(config *configuration.Configuration) {
	_this.config = config
	_this.reader.Init(config)
}

// Decode an already streamed document, sending all decoded events to eventReceiver.
func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	return _this.Decode(bytes.NewBuffer(document), eventReceiver)
}

// Read and decode a document from reader, sending all decoded events to eventReceiver.
func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !_this.config.Debug.PassThroughPanics {
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

	_this.reader.SetReader(reader)

	eventReceiver.OnBeginDocument()

	docHeader := _this.reader.ReadUint8()
	if docHeader != CBESignatureByte {
		_this.reader.errorf("First byte of CBE document must be 0x%02x (found 0x%02x)", CBESignatureByte, docHeader)
	}
	ver := _this.reader.ReadVersion()
	// TODO: Remove this when releasing V1
	if ver == 1 {
		ver = 0
	}
	eventReceiver.OnVersion(ver)

	return _this.runMainDecodeLoop(eventReceiver)
}

// ============================================================================

// Internal

func (_this *Decoder) runMainDecodeLoop(eventReceiver events.DataEventReceiver) (err error) {
	// TODO: Need to query rules to see when to stop
EOF:
	for {
		cbeType := _this.reader.ReadTypeOrEOF()
		switch cbeType {
		case cbeTypeDecimal:
			value, bigValue := _this.reader.ReadDecimalFloat()
			if bigValue != nil {
				eventReceiver.OnBigDecimalFloat(bigValue)
			} else {
				eventReceiver.OnDecimalFloat(value)
			}
		case cbeTypePosInt:
			asUint, asBig := _this.reader.ReadUint()
			if asBig != nil {
				eventReceiver.OnBigInt(asBig)
			} else {
				eventReceiver.OnPositiveInt(asUint)
			}
		case cbeTypeNegInt:
			asUint, asBig := _this.reader.ReadUint()
			if asBig != nil {
				eventReceiver.OnBigInt(asBig.Neg(asBig))
			} else {
				eventReceiver.OnNegativeInt(asUint)
			}
		case cbeTypePosInt8:
			eventReceiver.OnPositiveInt(uint64(_this.reader.ReadUint8()))
		case cbeTypeNegInt8:
			eventReceiver.OnNegativeInt(uint64(_this.reader.ReadUint8()))
		case cbeTypePosInt16:
			eventReceiver.OnPositiveInt(uint64(_this.reader.ReadUint16()))
		case cbeTypeNegInt16:
			eventReceiver.OnNegativeInt(uint64(_this.reader.ReadUint16()))
		case cbeTypePosInt32:
			eventReceiver.OnPositiveInt(uint64(_this.reader.ReadUint32()))
		case cbeTypeNegInt32:
			eventReceiver.OnNegativeInt(uint64(_this.reader.ReadUint32()))
		case cbeTypePosInt64:
			eventReceiver.OnPositiveInt(_this.reader.ReadUint64())
		case cbeTypeNegInt64:
			eventReceiver.OnNegativeInt(_this.reader.ReadUint64())
		case cbeTypeFloat16:
			f32Val := _this.reader.ReadFloat16()
			f64Val := float64(f32Val)
			if math.IsNaN(f64Val) && !common.HasQuietNanBitSet32(f32Val) {
				// Golang destroys the quiet bit status when converting to float64
				eventReceiver.OnFloat(common.Float64SignalingNan)
			} else {
				eventReceiver.OnFloat(f64Val)
			}
		case cbeTypeFloat32:
			f32Val := _this.reader.ReadFloat32()
			f64Val := float64(f32Val)
			if math.IsNaN(f64Val) && !common.HasQuietNanBitSet32(f32Val) {
				// Golang destroys the quiet bit status when converting to float64
				eventReceiver.OnFloat(common.Float64SignalingNan)
			} else {
				eventReceiver.OnFloat(f64Val)
			}
		case cbeTypeFloat64:
			eventReceiver.OnFloat(_this.reader.ReadFloat64())
		case cbeTypeUID:
			eventReceiver.OnUID(_this.reader.ReadBytes(16))
		case cbeTypeMap:
			eventReceiver.OnMap()
		case cbeTypeList:
			eventReceiver.OnList()
		case cbeTypeRecord:
			eventReceiver.OnRecord(_this.reader.ReadIdentifier())
		case cbeTypeEdge:
			eventReceiver.OnEdge()
		case cbeTypeNode:
			eventReceiver.OnNode()
		case cbeTypeEndContainer:
			eventReceiver.OnEndContainer()
		case cbeTypeFalse:
			eventReceiver.OnFalse()
		case cbeTypeTrue:
			eventReceiver.OnTrue()
		case cbeTypeNull:
			eventReceiver.OnNull()
		case cbeTypePadding:
			eventReceiver.OnPadding()
		case cbeTypeString0:
			eventReceiver.OnArray(events.ArrayTypeString, 0, []byte{})
		case cbeTypeString1, cbeTypeString2, cbeTypeString3, cbeTypeString4,
			cbeTypeString5, cbeTypeString6, cbeTypeString7, cbeTypeString8,
			cbeTypeString9, cbeTypeString10, cbeTypeString11, cbeTypeString12,
			cbeTypeString13, cbeTypeString14, cbeTypeString15:
			length := int(cbeType - cbeTypeString0)
			eventReceiver.OnArray(events.ArrayTypeString, uint64(length), _this.reader.ReadBytes(length))
		case cbeTypeString:
			_this.decodeArray(events.ArrayTypeString, eventReceiver)
		case cbeTypeRID:
			_this.decodeArray(events.ArrayTypeResourceID, eventReceiver)
		case cbeTypeCustomType:
			_this.decodeCustomType(eventReceiver)
		case cbeTypeEOF:
			break EOF
		case cbeTypePlane7f:
			_this.decodePlane7f(eventReceiver)
		case cbeTypeArrayBit:
			_this.decodeArray(events.ArrayTypeBit, eventReceiver)
		case cbeTypeArrayUint8:
			_this.decodeArray(events.ArrayTypeUint8, eventReceiver)
		case cbeTypeLocalReference:
			eventReceiver.OnReferenceLocal(_this.reader.ReadIdentifier())
		case cbeTypeDate:
			eventReceiver.OnTime(_this.reader.ReadDate())
		case cbeTypeTime:
			eventReceiver.OnTime(_this.reader.ReadTime())
		case cbeTypeTimestamp:
			eventReceiver.OnTime(_this.reader.ReadTimestamp())
		default:
			asSmallInt := int64(int8(cbeType))
			if asSmallInt < cbeSmallIntMin || asSmallInt > cbeSmallIntMax {
				panic(fmt.Errorf("0x%02x: Unsupported type", cbeType))
			}
			eventReceiver.OnInt(asSmallInt)
		}
	}

	eventReceiver.OnEndDocument()
	return
}

func (_this *Decoder) decodePlane7f(eventReceiver events.DataEventReceiver) {
	cbeType := _this.reader.ReadType()
	const lengthMask = 0x0f
	const shortTypeMask = 0xf0

	elementCount := int(cbeType) & lengthMask
	switch cbeType & shortTypeMask {
	case cbeTypeShortArrayInt8:
		bytes := _this.reader.ReadBytes(elementCount)
		eventReceiver.OnArray(events.ArrayTypeInt8, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayUint16:
		bytes := _this.reader.ReadBytes(elementCount * 2)
		eventReceiver.OnArray(events.ArrayTypeUint16, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayInt16:
		bytes := _this.reader.ReadBytes(elementCount * 2)
		eventReceiver.OnArray(events.ArrayTypeInt16, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayUint32:
		bytes := _this.reader.ReadBytes(elementCount * 4)
		eventReceiver.OnArray(events.ArrayTypeUint32, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayInt32:
		bytes := _this.reader.ReadBytes(elementCount * 4)
		eventReceiver.OnArray(events.ArrayTypeInt32, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayUint64:
		bytes := _this.reader.ReadBytes(elementCount * 8)
		eventReceiver.OnArray(events.ArrayTypeUint64, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayInt64:
		bytes := _this.reader.ReadBytes(elementCount * 8)
		eventReceiver.OnArray(events.ArrayTypeInt64, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayFloat16:
		bytes := _this.reader.ReadBytes(elementCount * 2)
		eventReceiver.OnArray(events.ArrayTypeFloat16, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayFloat32:
		bytes := _this.reader.ReadBytes(elementCount * 4)
		eventReceiver.OnArray(events.ArrayTypeFloat32, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayFloat64:
		bytes := _this.reader.ReadBytes(elementCount * 8)
		eventReceiver.OnArray(events.ArrayTypeFloat64, uint64(elementCount), bytes)
		return
	case cbeTypeShortArrayUID:
		bytes := _this.reader.ReadBytes(elementCount * 16)
		eventReceiver.OnArray(events.ArrayTypeUID, uint64(elementCount), bytes)
		return
	}

	switch cbeType {
	case cbeTypeMarker:
		eventReceiver.OnMarker(_this.reader.ReadIdentifier())
	case cbeTypeRecordType:
		eventReceiver.OnRecordType(_this.reader.ReadIdentifier())
	case cbeTypeRemoteReference:
		_this.decodeArray(events.ArrayTypeReferenceRemote, eventReceiver)
	case cbeTypeMedia:
		_this.decodeMedia(eventReceiver)
	default:
		arrayType := cbePlane7fTypeToArrayType[cbeType]
		if arrayType == events.ArrayTypeInvalid {
			panic(fmt.Errorf("0x%02x: Unsupported plane 2 type", cbeType))
		}
		_this.decodeArray(arrayType, eventReceiver)
	}
}

func (_this *Decoder) decodeArray(arrayType events.ArrayType, eventReceiver events.DataEventReceiver) {
	elementBitWidth := arrayType.ElementSize()
	eventReceiver.OnArrayBegin(arrayType)
	_this.decodeArrayChunks(eventReceiver, elementBitWidth)
}

func (_this *Decoder) decodeMedia(eventReceiver events.DataEventReceiver) {
	mediaTypeLength := _this.reader.readSmallULEB128("media type length", 0xffffffff)
	mediaType := string(_this.reader.ReadBytes(int(mediaTypeLength)))
	elementBitWidth := 8
	eventReceiver.OnMediaBegin(mediaType)
	_this.decodeArrayChunks(eventReceiver, elementBitWidth)
}

func (_this *Decoder) decodeCustomType(eventReceiver events.DataEventReceiver) {
	customType := _this.reader.readSmallULEB128("custom type code", 0xffffffff)
	elementBitWidth := 8
	eventReceiver.OnCustomBegin(events.ArrayTypeCustomBinary, customType)
	_this.decodeArrayChunks(eventReceiver, elementBitWidth)
}

func (_this *Decoder) decodeArrayChunks(eventReceiver events.DataEventReceiver, elementBitWidth int) {
	moreChunksFollow := true
	elementCount := uint64(0)

	for moreChunksFollow {
		elementCount, moreChunksFollow = _this.reader.ReadArrayChunkHeader()
		validateLength(elementCount)
		eventReceiver.OnArrayChunk(elementCount, moreChunksFollow)
		byteCount := common.ElementCountToByteCount(elementBitWidth, elementCount)
		if byteCount > 0 {
			nextBytes := _this.reader.ReadBytes(int(byteCount))
			eventReceiver.OnArrayData(nextBytes)
		}
	}
}

// Make sure we don't overflow max slice length
func validateLength(length uint64) {
	const maxDefaultInt = uint64((^uint(0)) >> 1)
	if length > maxDefaultInt {
		panic(fmt.Errorf("%v > max int value (%v)", length, maxDefaultInt))
	}
}
