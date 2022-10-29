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

package cte

import (
	"fmt"
	"math"

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type arrayEncoderEngine struct {
	stream                 *Writer
	addElementsFunc        func(b []byte)
	onComplete             func()
	arrayElementBitWidth   int
	arrayElementByteWidth  int
	remainingChunkElements uint64
	hasWrittenElements     bool
	moreChunksFollow       bool
	arrayChunkBacking      [16]byte
	arrayChunkLeftover     []byte
	stringBuffer           []byte
	config                 *configuration.CTEEncoderConfiguration
}

func (_this *arrayEncoderEngine) Init(stream *Writer, config *configuration.CTEEncoderConfiguration) {
	_this.stream = stream
	_this.arrayChunkLeftover = _this.arrayChunkBacking[:]
	_this.config = config
}

func (_this *arrayEncoderEngine) setElementBitWidth(width int) {
	_this.arrayElementBitWidth = width
	_this.arrayElementByteWidth = width / 8
}

func (_this *arrayEncoderEngine) setElementByteWidth(width int) {
	_this.arrayElementBitWidth = width * 8
	_this.arrayElementByteWidth = width
}

func (_this *arrayEncoderEngine) EncodeStringlikeArray(arrayType events.ArrayType, data string) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.stream.WriteQuotedString(true, data)
	case events.ArrayTypeResourceID:
		_this.stream.WriteByteNotLF('@')
		_this.stream.WriteQuotedString(false, data)
	case events.ArrayTypeReferenceRemote:
		_this.stream.WriteByteNotLF('$')
		_this.stream.WriteQuotedString(false, data)
	default:
		panic(fmt.Errorf("BUG: EncodeStringlikeArray passed unhandled array type %v", arrayType))
	}
}

func (_this *arrayEncoderEngine) EncodeArray(arrayType events.ArrayType, elementCount uint64, data []uint8) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.stream.WriteQuotedStringBytes(true, data)
	case events.ArrayTypeResourceID:
		_this.stream.WriteByteNotLF('@')
		_this.stream.WriteQuotedStringBytes(false, data)
	case events.ArrayTypeReferenceRemote:
		_this.stream.WriteByteNotLF('$')
		_this.stream.WriteQuotedStringBytes(false, data)
	default:
		_this.BeginArray(arrayType, func() {})
		_this.BeginChunk(elementCount, false)
		if elementCount > 0 {
			_this.AddArrayData(data)
		}
	}
}

func (_this *arrayEncoderEngine) EncodeMedia(mediaType string, data []byte) {
	_this.stream.WriteFmtNotLF("|%v", mediaType)
	_this.stream.WriteHexBytes(data)
	_this.stream.WriteArrayEnd()
}

func (_this *arrayEncoderEngine) EncodeCustomBinary(customType uint64, data []byte) {
	_this.stream.WriteFmtNotLF("|c%v", customType)
	_this.stream.WriteHexBytes(data)
	_this.stream.WriteArrayEnd()
}

func (_this *arrayEncoderEngine) EncodeCustomText(customType uint64, data string) {
	_this.stream.WriteFmtNotLF("|c%v ", customType)
	_this.stream.WriteQuotedString(true, data)
	_this.stream.WriteArrayEnd()
}

func (_this *arrayEncoderEngine) BeginArray(arrayType events.ArrayType, onComplete func()) {
	_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
	_this.stringBuffer = _this.stringBuffer[:0]
	_this.remainingChunkElements = 0
	_this.hasWrittenElements = false

	// Default completion operation
	_this.onComplete = func() {
		_this.stream.WriteArrayEnd()
		onComplete()
	}

	beginOp := arrayEncodeBeginOps[arrayType]
	beginOp(_this, onComplete)
}

func (_this *arrayEncoderEngine) BeginMedia(mediaType string, onComplete func()) {
	_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
	_this.stringBuffer = _this.stringBuffer[:0]
	_this.remainingChunkElements = 0
	_this.hasWrittenElements = false

	_this.setElementByteWidth(1)
	_this.stream.WriteFmtNotLF("|%v", mediaType)
	_this.addElementsFunc = func(data []byte) { _this.stream.WriteHexBytes(data) }
	_this.onComplete = func() {
		_this.stream.WriteArrayEnd()
		onComplete()
	}
}

func (_this *arrayEncoderEngine) BeginCustomText(customType uint64, onComplete func()) {
	_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
	_this.stringBuffer = _this.stringBuffer[:0]
	_this.remainingChunkElements = 0
	_this.hasWrittenElements = false

	_this.setElementByteWidth(1)
	_this.stream.WriteFmtNotLF("|c%v ", customType)
	_this.addElementsFunc = func(data []byte) {
		_this.handleFirstElement(data)
		_this.appendStringbuffer(data)
	}
	_this.onComplete = func() {
		_this.stream.WriteQuotedStringBytes(true, _this.stringBuffer)
		_this.stream.WriteArrayEnd()
		onComplete()
	}
}

func (_this *arrayEncoderEngine) BeginCustomBinary(customType uint64, onComplete func()) {
	_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
	_this.stringBuffer = _this.stringBuffer[:0]
	_this.remainingChunkElements = 0
	_this.hasWrittenElements = false

	_this.setElementByteWidth(1)
	_this.stream.WriteFmtNotLF("|c%v", customType)
	_this.addElementsFunc = func(data []byte) { _this.stream.WriteHexBytes(data) }
	_this.onComplete = func() {
		_this.stream.WriteArrayEnd()
		onComplete()
	}
}

func (_this *arrayEncoderEngine) endArray() {
	_this.onComplete()
}

func (_this *arrayEncoderEngine) handleFirstElement(data []byte) {
	if !_this.hasWrittenElements && len(data) > 0 {
		_this.stream.WriteByteNotLF(' ')
		_this.hasWrittenElements = true
	}
}

func (_this *arrayEncoderEngine) BeginChunk(elementCount uint64, moreChunksFollow bool) {
	_this.remainingChunkElements = elementCount
	_this.moreChunksFollow = moreChunksFollow

	if elementCount == 0 && !moreChunksFollow {
		_this.endArray()
	}
}

func (_this *arrayEncoderEngine) addBooleanArrayData(data []byte) {
	_this.handleFirstElement(data)
	for _this.remainingChunkElements >= 8 && len(data) > 0 {
		b := data[0]
		for i := 0; i < 8; i++ {
			if (b & (1 << i)) != 0 {
				_this.stream.WriteByteNotLF('1')
			} else {
				_this.stream.WriteByteNotLF('0')
			}
		}
		data = data[1:]
		_this.remainingChunkElements -= 8
	}
	if _this.remainingChunkElements > 0 && len(data) > 0 {
		count := _this.remainingChunkElements
		b := data[0]
		for i := 0; i < int(count); i++ {
			if (b & (1 << i)) != 0 {
				_this.stream.WriteByteNotLF('1')
			} else {
				_this.stream.WriteByteNotLF('0')
			}
		}
		_this.remainingChunkElements -= count
	}
	if _this.remainingChunkElements == 0 && !_this.moreChunksFollow {
		_this.endArray()
	}
}

func (_this *arrayEncoderEngine) AddArrayData(data []byte) {
	if _this.arrayElementBitWidth == 1 {
		_this.addBooleanArrayData(data)
		return
	}

	if _this.arrayElementByteWidth > 1 {
		leftoverLength := len(_this.arrayChunkLeftover)
		if leftoverLength > 0 {
			fillCount := _this.arrayElementByteWidth - leftoverLength

			if len(data) < fillCount {
				_this.arrayChunkLeftover = append(_this.arrayChunkLeftover, data...)
				return
			}

			_this.arrayChunkLeftover = append(_this.arrayChunkLeftover, data[:fillCount]...)
			data = data[fillCount:]
			_this.addElementsFunc(_this.arrayChunkLeftover)
			_this.remainingChunkElements--
			_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
		}

		widthMask := _this.arrayElementByteWidth - 1
		remainderCount := len(data) & widthMask
		if remainderCount != 0 {
			_this.arrayChunkLeftover = append(_this.arrayChunkLeftover, data[len(data)-remainderCount:]...)
			data = data[:len(data)-remainderCount]
		}
	}
	_this.addElementsFunc(data)
	_this.remainingChunkElements -= uint64(len(data) / _this.arrayElementByteWidth)
	if _this.remainingChunkElements == 0 && !_this.moreChunksFollow {
		_this.endArray()
	}
}

// ============================================================================

// Utils

func (_this *arrayEncoderEngine) beginArrayBoolean(onComplete func()) {
	_this.setElementBitWidth(1)
	_this.stream.WriteStringNotLF("|b")
}

func (_this *arrayEncoderEngine) beginArrayString(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteQuotedStringBytes(true, _this.stringBuffer)
		onComplete()
	}
}

func (_this *arrayEncoderEngine) beginArrayResourceID(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteByteNotLF('@')
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteQuotedStringBytes(false, _this.stringBuffer)
		onComplete()
	}
}

func (_this *arrayEncoderEngine) beginArrayUint8(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteStringNotLF(arrayHeadersUint8[_this.config.DefaultNumericFormats.Array.Uint8])
	format := arrayFormats8[_this.config.DefaultNumericFormats.Array.Uint8]
	_this.addElementsFunc = func(data []byte) {
		for _, b := range data {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, b)
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersUint16[_this.config.DefaultNumericFormats.Array.Uint16])
	format := arrayFormats16[_this.config.DefaultNumericFormats.Array.Uint16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, uint(data[0])|(uint(data[1])<<8))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersUint32[_this.config.DefaultNumericFormats.Array.Uint32])
	format := arrayFormats32[_this.config.DefaultNumericFormats.Array.Uint32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, uint(data[0])|(uint(data[1])<<8)|(uint(data[2])<<16)|(uint(data[3])<<24))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersUint64[_this.config.DefaultNumericFormats.Array.Uint64])
	format := arrayFormats64[_this.config.DefaultNumericFormats.Array.Uint64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, uint64(data[0])|(uint64(data[1])<<8)|(uint64(data[2])<<16)|(uint64(data[3])<<24)|
				(uint64(data[4])<<32)|(uint64(data[5])<<40)|(uint64(data[6])<<48)|(uint64(data[7])<<56))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt8(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteStringNotLF(arrayHeadersInt8[_this.config.DefaultNumericFormats.Array.Int8])
	format := arrayFormats8[_this.config.DefaultNumericFormats.Array.Int8]
	_this.addElementsFunc = func(data []byte) {
		for _, b := range data {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, int8(b))
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersInt16[_this.config.DefaultNumericFormats.Array.Int16])
	format := arrayFormats16[_this.config.DefaultNumericFormats.Array.Int16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, int16(data[0])|(int16(data[1])<<8))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersInt32[_this.config.DefaultNumericFormats.Array.Int32])
	format := arrayFormats32[_this.config.DefaultNumericFormats.Array.Int32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, int32(data[0])|(int32(data[1])<<8)|(int32(data[2])<<16)|(int32(data[3])<<24))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersInt64[_this.config.DefaultNumericFormats.Array.Int64])
	format := arrayFormats64[_this.config.DefaultNumericFormats.Array.Int64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteFmtNotLF(format, int64(data[0])|(int64(data[1])<<8)|(int64(data[2])<<16)|(int64(data[3])<<24)|
				(int64(data[4])<<32)|(int64(data[5])<<40)|(int64(data[6])<<48)|(int64(data[7])<<56))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersFloat16[_this.config.DefaultNumericFormats.Array.Float16])
	if _this.config.DefaultNumericFormats.Array.Float16 == configuration.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByteNotLF(' ')
				bits := uint16(data[0]) | (uint16(data[1]) << 8)
				_this.stream.WriteFloatHexNoPrefix(common.Float64FromFloat16Bits(bits))
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.config.DefaultNumericFormats.Array.Float16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			bits := uint16(data[0]) | (uint16(data[1]) << 8)
			_this.stream.WriteFloatUsingFormat(common.Float64FromFloat16Bits(bits), format)
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersFloat32[_this.config.DefaultNumericFormats.Array.Float32])
	if _this.config.DefaultNumericFormats.Array.Float32 == configuration.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByteNotLF(' ')
				bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
				_this.stream.WriteFloatHexNoPrefix(common.Float64FromFloat32Bits(bits))
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.config.DefaultNumericFormats.Array.Float32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
			_this.stream.WriteFloatUsingFormat(common.Float64FromFloat32Bits(bits), format)
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF(arrayHeadersFloat64[_this.config.DefaultNumericFormats.Array.Float64])
	if _this.config.DefaultNumericFormats.Array.Float64 == configuration.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByteNotLF(' ')
				bits := uint64(data[0]) | (uint64(data[1]) << 8) | (uint64(data[2]) << 16) | (uint64(data[3]) << 24) |
					(uint64(data[4]) << 32) | (uint64(data[5]) << 40) | (uint64(data[6]) << 48) | (uint64(data[7]) << 56)
				_this.stream.WriteFloatHexNoPrefix(math.Float64frombits(bits))
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.config.DefaultNumericFormats.Array.Float64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			bits := uint64(data[0]) | (uint64(data[1]) << 8) | (uint64(data[2]) << 16) | (uint64(data[3]) << 24) |
				(uint64(data[4]) << 32) | (uint64(data[5]) << 40) | (uint64(data[6]) << 48) | (uint64(data[7]) << 56)
			_this.stream.WriteFloatUsingFormat(math.Float64frombits(bits), format)
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUID(onComplete func()) {
	const elemWidth = 16
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteStringNotLF("|u")
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByteNotLF(' ')
			_this.stream.WriteUID(data[:elemWidth])
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayMedia(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteStringNotLF("|m ")
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteBytesNotLF(_this.stringBuffer)
		_this.stringBuffer = _this.stringBuffer[:0]
		_this.addElementsFunc = func(data []byte) { _this.stream.WriteHexBytes(data) }
		_this.onComplete = func() {
			_this.stream.WriteByteNotLF('|')
			onComplete()
		}
	}
}

func (_this *arrayEncoderEngine) appendStringbuffer(data []byte) {
	_this.stringBuffer = append(_this.stringBuffer, data...)
}

// ============================================================================

// Data

var arrayEncodeBeginOps = []func(*arrayEncoderEngine, func()){
	events.ArrayTypeBit:        (*arrayEncoderEngine).beginArrayBoolean,
	events.ArrayTypeString:     (*arrayEncoderEngine).beginArrayString,
	events.ArrayTypeResourceID: (*arrayEncoderEngine).beginArrayResourceID,
	events.ArrayTypeUint8:      (*arrayEncoderEngine).beginArrayUint8,
	events.ArrayTypeUint16:     (*arrayEncoderEngine).beginArrayUint16,
	events.ArrayTypeUint32:     (*arrayEncoderEngine).beginArrayUint32,
	events.ArrayTypeUint64:     (*arrayEncoderEngine).beginArrayUint64,
	events.ArrayTypeInt8:       (*arrayEncoderEngine).beginArrayInt8,
	events.ArrayTypeInt16:      (*arrayEncoderEngine).beginArrayInt16,
	events.ArrayTypeInt32:      (*arrayEncoderEngine).beginArrayInt32,
	events.ArrayTypeInt64:      (*arrayEncoderEngine).beginArrayInt64,
	events.ArrayTypeFloat16:    (*arrayEncoderEngine).beginArrayFloat16,
	events.ArrayTypeFloat32:    (*arrayEncoderEngine).beginArrayFloat32,
	events.ArrayTypeFloat64:    (*arrayEncoderEngine).beginArrayFloat64,
	events.ArrayTypeUID:        (*arrayEncoderEngine).beginArrayUID,
	events.ArrayTypeMedia:      (*arrayEncoderEngine).beginArrayMedia,
}

var arrayFormatsGeneral = []string{
	configuration.CTEEncodingFormatDecimal:               "%v",
	configuration.CTEEncodingFormatBinary:                "%b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "%b",
	configuration.CTEEncodingFormatOctal:                 "%o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "%o",
	configuration.CTEEncodingFormatHexadecimal:           "%x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "%x",
}

var arrayFormats8 = []string{
	configuration.CTEEncodingFormatDecimal:               "%v",
	configuration.CTEEncodingFormatBinary:                "%b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "%08b",
	configuration.CTEEncodingFormatOctal:                 "%o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "%03o",
	configuration.CTEEncodingFormatHexadecimal:           "%x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "%02x",
}

var arrayFormats16 = []string{
	configuration.CTEEncodingFormatDecimal:               "%v",
	configuration.CTEEncodingFormatBinary:                "%b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "%016b",
	configuration.CTEEncodingFormatOctal:                 "%o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "%06o",
	configuration.CTEEncodingFormatHexadecimal:           "%x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "%04x",
}

var arrayFormats32 = []string{
	configuration.CTEEncodingFormatDecimal:               "%v",
	configuration.CTEEncodingFormatBinary:                "%b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "%032b",
	configuration.CTEEncodingFormatOctal:                 "%o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "%011o",
	configuration.CTEEncodingFormatHexadecimal:           "%x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "%08x",
}

var arrayFormats64 = []string{
	configuration.CTEEncodingFormatDecimal:               "%v",
	configuration.CTEEncodingFormatBinary:                "%b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "%064b",
	configuration.CTEEncodingFormatOctal:                 "%o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "%022o",
	configuration.CTEEncodingFormatHexadecimal:           "%x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "%016x",
}

var arrayHeadersUint8 = []string{
	configuration.CTEEncodingFormatDecimal:               "|u8",
	configuration.CTEEncodingFormatBinary:                "|u8b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|u8b",
	configuration.CTEEncodingFormatOctal:                 "|u8o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|u8o",
	configuration.CTEEncodingFormatHexadecimal:           "|u8x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|u8x",
}
var arrayHeadersUint16 = []string{
	configuration.CTEEncodingFormatDecimal:               "|u16",
	configuration.CTEEncodingFormatBinary:                "|u16b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|u16b",
	configuration.CTEEncodingFormatOctal:                 "|u16o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|u16o",
	configuration.CTEEncodingFormatHexadecimal:           "|u16x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|u16x",
}
var arrayHeadersUint32 = []string{
	configuration.CTEEncodingFormatDecimal:               "|u32",
	configuration.CTEEncodingFormatBinary:                "|u32b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|u32b",
	configuration.CTEEncodingFormatOctal:                 "|u32o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|u32o",
	configuration.CTEEncodingFormatHexadecimal:           "|u32x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|u32x",
}
var arrayHeadersUint64 = []string{
	configuration.CTEEncodingFormatDecimal:               "|u64",
	configuration.CTEEncodingFormatBinary:                "|u64b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|u64b",
	configuration.CTEEncodingFormatOctal:                 "|u64o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|u64o",
	configuration.CTEEncodingFormatHexadecimal:           "|u64x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|u64x",
}
var arrayHeadersInt8 = []string{
	configuration.CTEEncodingFormatDecimal:               "|i8",
	configuration.CTEEncodingFormatBinary:                "|i8b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|i8b",
	configuration.CTEEncodingFormatOctal:                 "|i8o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|i8o",
	configuration.CTEEncodingFormatHexadecimal:           "|i8x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|i8x",
}
var arrayHeadersInt16 = []string{
	configuration.CTEEncodingFormatDecimal:               "|i16",
	configuration.CTEEncodingFormatBinary:                "|i16b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|i16b",
	configuration.CTEEncodingFormatOctal:                 "|i16o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|i16o",
	configuration.CTEEncodingFormatHexadecimal:           "|i16x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|i16x",
}
var arrayHeadersInt32 = []string{
	configuration.CTEEncodingFormatDecimal:               "|i32",
	configuration.CTEEncodingFormatBinary:                "|i32b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|i32b",
	configuration.CTEEncodingFormatOctal:                 "|i32o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|i32o",
	configuration.CTEEncodingFormatHexadecimal:           "|i32x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|i32x",
}
var arrayHeadersInt64 = []string{
	configuration.CTEEncodingFormatDecimal:               "|i64",
	configuration.CTEEncodingFormatBinary:                "|i64b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|i64b",
	configuration.CTEEncodingFormatOctal:                 "|i64o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|i64o",
	configuration.CTEEncodingFormatHexadecimal:           "|i64x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|i64x",
}
var arrayHeadersFloat16 = []string{
	configuration.CTEEncodingFormatDecimal:               "|f16",
	configuration.CTEEncodingFormatBinary:                "|f16b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|f16b",
	configuration.CTEEncodingFormatOctal:                 "|f16o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|f16o",
	configuration.CTEEncodingFormatHexadecimal:           "|f16x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|f16x",
}
var arrayHeadersFloat32 = []string{
	configuration.CTEEncodingFormatDecimal:               "|f32",
	configuration.CTEEncodingFormatBinary:                "|f32b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|f32b",
	configuration.CTEEncodingFormatOctal:                 "|f32o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|f32o",
	configuration.CTEEncodingFormatHexadecimal:           "|f32x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|f32x",
}
var arrayHeadersFloat64 = []string{
	configuration.CTEEncodingFormatDecimal:               "|f64",
	configuration.CTEEncodingFormatBinary:                "|f64b",
	configuration.CTEEncodingFormatBinaryZeroFilled:      "|f64b",
	configuration.CTEEncodingFormatOctal:                 "|f64o",
	configuration.CTEEncodingFormatOctalZeroFilled:       "|f64o",
	configuration.CTEEncodingFormatHexadecimal:           "|f64x",
	configuration.CTEEncodingFormatHexadecimalZeroFilled: "|f64x",
}
