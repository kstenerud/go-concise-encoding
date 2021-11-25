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

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
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
	opts                   *options.CTEEncoderOptions
}

func (_this *arrayEncoderEngine) Init(stream *Writer, opts *options.CTEEncoderOptions) {
	_this.stream = stream
	_this.arrayChunkLeftover = _this.arrayChunkBacking[:]
	_this.opts = opts
}

func (_this *arrayEncoderEngine) setElementBitWidth(width int) {
	_this.arrayElementBitWidth = width
	_this.arrayElementByteWidth = width / 8
}

func (_this *arrayEncoderEngine) setElementByteWidth(width int) {
	_this.arrayElementBitWidth = width * 8
	_this.arrayElementByteWidth = width
}

func (_this *arrayEncoderEngine) EncodeStringlikeArray(stringContext stringContext, arrayType events.ArrayType, data string) {
	switch arrayType {
	case events.ArrayTypeString:
		switch stringContext {
		case stringContextDefault:
			_this.stream.WriteQuotedString(data)
		case stringContextComment:
			_this.stream.WriteString(data)
		case stringContextMarkup:
			_this.stream.WritePotentiallyEscapedMarkupContents(data)
		}
	case events.ArrayTypeResourceID:
		_this.stream.WriteByte('@')
		_this.stream.WriteQuotedString(data)
	case events.ArrayTypeRemoteRef:
		_this.stream.WriteByte('$')
		_this.stream.WriteQuotedString(data)
	case events.ArrayTypeCustomText:
		_this.stream.WriteString("|ct ")
		_this.stream.WritePotentiallyEscapedStringArrayContents(data)
		_this.stream.WriteArrayEnd()
	default:
		panic(fmt.Errorf("BUG: EncodeStringlikeArray passed unhandled array type %v", arrayType))
	}
}

func (_this *arrayEncoderEngine) EncodeArray(stringContext stringContext, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	switch arrayType {
	case events.ArrayTypeString:
		switch stringContext {
		case stringContextDefault:
			_this.stream.WriteQuotedStringBytes(data)
		case stringContextComment:
			_this.stream.WriteBytes(data)
		case stringContextMarkup:
			_this.stream.WritePotentiallyEscapedMarkupContentsBytes(data)
		}
	case events.ArrayTypeResourceID:
		_this.stream.WriteByte('@')
		_this.stream.WriteQuotedStringBytes(data)
	case events.ArrayTypeRemoteRef:
		_this.stream.WriteByte('$')
		_this.stream.WriteQuotedStringBytes(data)
	case events.ArrayTypeCustomText:
		_this.stream.WriteString("|ct ")
		_this.stream.WritePotentiallyEscapedStringArrayContentsBytes(data)
		_this.stream.WriteArrayEnd()
	default:
		_this.BeginArray(stringContext, arrayType, func() {})
		_this.BeginChunk(elementCount, false)
		if elementCount > 0 {
			_this.AddArrayData(data)
		}
	}
}

func (_this *arrayEncoderEngine) BeginArray(stringContext stringContext, arrayType events.ArrayType, onComplete func()) {
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

func (_this *arrayEncoderEngine) endArray() {
	_this.onComplete()
}

func (_this *arrayEncoderEngine) handleFirstElement(data []byte) {
	if !_this.hasWrittenElements && len(data) > 0 {
		_this.stream.WriteByte(' ')
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
				_this.stream.WriteByte('1')
			} else {
				_this.stream.WriteByte('0')
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
				_this.stream.WriteByte('1')
			} else {
				_this.stream.WriteByte('0')
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
	_this.stream.WriteString("|b")
}

func (_this *arrayEncoderEngine) beginArrayString(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteQuotedStringBytes(_this.stringBuffer)
		onComplete()
	}
}

func (_this *arrayEncoderEngine) beginArrayResourceID(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteByte('@')
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteQuotedStringBytes(_this.stringBuffer)
		onComplete()
	}
}

func (_this *arrayEncoderEngine) beginArrayCustomText(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteString("|ct")
	_this.addElementsFunc = func(data []byte) {
		_this.handleFirstElement(data)
		_this.appendStringbuffer(data)
	}
	_this.onComplete = func() {
		_this.stream.WritePotentiallyEscapedStringArrayContentsBytes(_this.stringBuffer)
		_this.stream.WriteArrayEnd()
		onComplete()
	}
}

func (_this *arrayEncoderEngine) beginArrayCustomBinary(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteString("|cb")
	_this.addElementsFunc = func(data []byte) { _this.stream.WriteHexBytes(data) }
}

func (_this *arrayEncoderEngine) beginArrayUint8(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteString(arrayHeadersUint8[_this.opts.DefaultFormats.Array.Uint8])
	format := arrayFormats8[_this.opts.DefaultFormats.Array.Uint8]
	_this.addElementsFunc = func(data []byte) {
		for _, b := range data {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, b)
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersUint16[_this.opts.DefaultFormats.Array.Uint16])
	format := arrayFormats16[_this.opts.DefaultFormats.Array.Uint16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, uint(data[0])|(uint(data[1])<<8))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersUint32[_this.opts.DefaultFormats.Array.Uint32])
	format := arrayFormats32[_this.opts.DefaultFormats.Array.Uint32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, uint(data[0])|(uint(data[1])<<8)|(uint(data[2])<<16)|(uint(data[3])<<24))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayUint64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersUint64[_this.opts.DefaultFormats.Array.Uint64])
	format := arrayFormats64[_this.opts.DefaultFormats.Array.Uint64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, uint64(data[0])|(uint64(data[1])<<8)|(uint64(data[2])<<16)|(uint64(data[3])<<24)|
				(uint64(data[4])<<32)|(uint64(data[5])<<40)|(uint64(data[6])<<48)|(uint64(data[7])<<56))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt8(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteString(arrayHeadersInt8[_this.opts.DefaultFormats.Array.Int8])
	format := arrayFormats8[_this.opts.DefaultFormats.Array.Int8]
	_this.addElementsFunc = func(data []byte) {
		for _, b := range data {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, int8(b))
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersInt16[_this.opts.DefaultFormats.Array.Int16])
	format := arrayFormats16[_this.opts.DefaultFormats.Array.Int16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, int16(data[0])|(int16(data[1])<<8))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersInt32[_this.opts.DefaultFormats.Array.Int32])
	format := arrayFormats32[_this.opts.DefaultFormats.Array.Int32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, int32(data[0])|(int32(data[1])<<8)|(int32(data[2])<<16)|(int32(data[3])<<24))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayInt64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersInt64[_this.opts.DefaultFormats.Array.Int64])
	format := arrayFormats64[_this.opts.DefaultFormats.Array.Int64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteFmt(format, int64(data[0])|(int64(data[1])<<8)|(int64(data[2])<<16)|(int64(data[3])<<24)|
				(int64(data[4])<<32)|(int64(data[5])<<40)|(int64(data[6])<<48)|(int64(data[7])<<56))
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat16(onComplete func()) {
	const elemWidth = 2
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersFloat16[_this.opts.DefaultFormats.Array.Float16])
	if _this.opts.DefaultFormats.Array.Float16 == options.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByte(' ')
				bits := (uint32(data[0]) << 16) | (uint32(data[1]) << 24)
				if bits == uint32(common.Bfloat16SignalingNanBits)<<16 {
					_this.stream.WriteSignalingNan()
				} else if bits == uint32(common.Bfloat16QuietNanBits)<<16 {
					_this.stream.WriteQuietNan()
				} else {
					v := math.Float32frombits(bits)
					if v < 0 {
						if math.IsInf(float64(v), -1) {
							_this.stream.WriteNegInfinity()
						} else {
							_this.stream.WriteByte('-')
							_this.stream.WriteFmtStripped(3, "%x", v)
						}
					} else {
						if math.IsInf(float64(v), 1) {
							_this.stream.WritePosInfinity()
						} else {
							_this.stream.WriteFmtStripped(2, "%x", v)
						}
					}
				}
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float16]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			bits := (uint32(data[0]) << 16) | (uint32(data[1]) << 24)
			if bits == uint32(common.Bfloat16SignalingNanBits)<<16 {
				_this.stream.WriteSignalingNan()
			} else if bits == uint32(common.Bfloat16QuietNanBits)<<16 {
				_this.stream.WriteQuietNan()
			} else {
				v := math.Float32frombits(bits)
				if math.IsInf(float64(v), 0) {
					if v < 0 {
						_this.stream.WriteNegInfinity()
					} else {
						_this.stream.WritePosInfinity()
					}
				} else {
					_this.stream.WriteFmt(format, v)
				}
			}
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat32(onComplete func()) {
	const elemWidth = 4
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersFloat32[_this.opts.DefaultFormats.Array.Float32])
	if _this.opts.DefaultFormats.Array.Float32 == options.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByte(' ')
				bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
				v := math.Float32frombits(bits)
				_this.stream.WriteFloatHexNoPrefix(float64(v))
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float32]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
			_this.stream.WriteFloat32UsingFormat(math.Float32frombits(bits), format)
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayFloat64(onComplete func()) {
	const elemWidth = 8
	_this.setElementByteWidth(elemWidth)
	_this.stream.WriteString(arrayHeadersFloat64[_this.opts.DefaultFormats.Array.Float64])
	if _this.opts.DefaultFormats.Array.Float64 == options.CTEEncodingFormatHexadecimal {
		_this.addElementsFunc = func(data []byte) {
			for len(data) > 0 {
				_this.stream.WriteByte(' ')
				bits := uint64(data[0]) | (uint64(data[1]) << 8) | (uint64(data[2]) << 16) | (uint64(data[3]) << 24) |
					(uint64(data[4]) << 32) | (uint64(data[5]) << 40) | (uint64(data[6]) << 48) | (uint64(data[7]) << 56)
				v := math.Float64frombits(bits)
				_this.stream.WriteFloatHexNoPrefix(v)
				data = data[elemWidth:]
			}
		}
		return
	}

	format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float64]
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
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
	_this.stream.WriteString("|u")
	_this.addElementsFunc = func(data []byte) {
		for len(data) > 0 {
			_this.stream.WriteByte(' ')
			_this.stream.WriteUID(data[:elemWidth])
			data = data[elemWidth:]
		}
	}
}

func (_this *arrayEncoderEngine) beginArrayMedia(onComplete func()) {
	_this.setElementByteWidth(1)
	_this.stream.WriteByte('|')
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.onComplete = func() {
		_this.stream.WriteBytes(_this.stringBuffer)
		_this.stringBuffer = _this.stringBuffer[:0]
		_this.addElementsFunc = func(data []byte) { _this.stream.WriteHexBytes(data) }
		_this.onComplete = func() {
			_this.stream.WriteByte('|')
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
	events.ArrayTypeBit:          (*arrayEncoderEngine).beginArrayBoolean,
	events.ArrayTypeString:       (*arrayEncoderEngine).beginArrayString,
	events.ArrayTypeResourceID:   (*arrayEncoderEngine).beginArrayResourceID,
	events.ArrayTypeCustomText:   (*arrayEncoderEngine).beginArrayCustomText,
	events.ArrayTypeCustomBinary: (*arrayEncoderEngine).beginArrayCustomBinary,
	events.ArrayTypeUint8:        (*arrayEncoderEngine).beginArrayUint8,
	events.ArrayTypeUint16:       (*arrayEncoderEngine).beginArrayUint16,
	events.ArrayTypeUint32:       (*arrayEncoderEngine).beginArrayUint32,
	events.ArrayTypeUint64:       (*arrayEncoderEngine).beginArrayUint64,
	events.ArrayTypeInt8:         (*arrayEncoderEngine).beginArrayInt8,
	events.ArrayTypeInt16:        (*arrayEncoderEngine).beginArrayInt16,
	events.ArrayTypeInt32:        (*arrayEncoderEngine).beginArrayInt32,
	events.ArrayTypeInt64:        (*arrayEncoderEngine).beginArrayInt64,
	events.ArrayTypeFloat16:      (*arrayEncoderEngine).beginArrayFloat16,
	events.ArrayTypeFloat32:      (*arrayEncoderEngine).beginArrayFloat32,
	events.ArrayTypeFloat64:      (*arrayEncoderEngine).beginArrayFloat64,
	events.ArrayTypeUID:          (*arrayEncoderEngine).beginArrayUID,
	events.ArrayTypeMedia:        (*arrayEncoderEngine).beginArrayMedia,
}

var arrayFormatsGeneral = []string{
	options.CTEEncodingFormatUnset:                 "%v",
	options.CTEEncodingFormatBinary:                "%b",
	options.CTEEncodingFormatBinaryZeroFilled:      "%b",
	options.CTEEncodingFormatOctal:                 "%o",
	options.CTEEncodingFormatOctalZeroFilled:       "%o",
	options.CTEEncodingFormatHexadecimal:           "%x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "%x",
}

var arrayFormats8 = []string{
	options.CTEEncodingFormatUnset:                 "%v",
	options.CTEEncodingFormatBinary:                "%b",
	options.CTEEncodingFormatBinaryZeroFilled:      "%08b",
	options.CTEEncodingFormatOctal:                 "%o",
	options.CTEEncodingFormatOctalZeroFilled:       "%03o",
	options.CTEEncodingFormatHexadecimal:           "%x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "%02x",
}

var arrayFormats16 = []string{
	options.CTEEncodingFormatUnset:                 "%v",
	options.CTEEncodingFormatBinary:                "%b",
	options.CTEEncodingFormatBinaryZeroFilled:      "%016b",
	options.CTEEncodingFormatOctal:                 "%o",
	options.CTEEncodingFormatOctalZeroFilled:       "%06o",
	options.CTEEncodingFormatHexadecimal:           "%x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "%04x",
}

var arrayFormats32 = []string{
	options.CTEEncodingFormatUnset:                 "%v",
	options.CTEEncodingFormatBinary:                "%b",
	options.CTEEncodingFormatBinaryZeroFilled:      "%032b",
	options.CTEEncodingFormatOctal:                 "%o",
	options.CTEEncodingFormatOctalZeroFilled:       "%011o",
	options.CTEEncodingFormatHexadecimal:           "%x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "%08x",
}

var arrayFormats64 = []string{
	options.CTEEncodingFormatUnset:                 "%v",
	options.CTEEncodingFormatBinary:                "%b",
	options.CTEEncodingFormatBinaryZeroFilled:      "%064b",
	options.CTEEncodingFormatOctal:                 "%o",
	options.CTEEncodingFormatOctalZeroFilled:       "%022o",
	options.CTEEncodingFormatHexadecimal:           "%x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "%016x",
}

var arrayHeadersUint8 = []string{
	options.CTEEncodingFormatUnset:                 "|u8",
	options.CTEEncodingFormatBinary:                "|u8b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|u8b",
	options.CTEEncodingFormatOctal:                 "|u8o",
	options.CTEEncodingFormatOctalZeroFilled:       "|u8o",
	options.CTEEncodingFormatHexadecimal:           "|u8x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|u8x",
}
var arrayHeadersUint16 = []string{
	options.CTEEncodingFormatUnset:                 "|u16",
	options.CTEEncodingFormatBinary:                "|u16b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|u16b",
	options.CTEEncodingFormatOctal:                 "|u16o",
	options.CTEEncodingFormatOctalZeroFilled:       "|u16o",
	options.CTEEncodingFormatHexadecimal:           "|u16x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|u16x",
}
var arrayHeadersUint32 = []string{
	options.CTEEncodingFormatUnset:                 "|u32",
	options.CTEEncodingFormatBinary:                "|u32b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|u32b",
	options.CTEEncodingFormatOctal:                 "|u32o",
	options.CTEEncodingFormatOctalZeroFilled:       "|u32o",
	options.CTEEncodingFormatHexadecimal:           "|u32x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|u32x",
}
var arrayHeadersUint64 = []string{
	options.CTEEncodingFormatUnset:                 "|u64",
	options.CTEEncodingFormatBinary:                "|u64b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|u64b",
	options.CTEEncodingFormatOctal:                 "|u64o",
	options.CTEEncodingFormatOctalZeroFilled:       "|u64o",
	options.CTEEncodingFormatHexadecimal:           "|u64x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|u64x",
}
var arrayHeadersInt8 = []string{
	options.CTEEncodingFormatUnset:                 "|i8",
	options.CTEEncodingFormatBinary:                "|i8b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|i8b",
	options.CTEEncodingFormatOctal:                 "|i8o",
	options.CTEEncodingFormatOctalZeroFilled:       "|i8o",
	options.CTEEncodingFormatHexadecimal:           "|i8x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|i8x",
}
var arrayHeadersInt16 = []string{
	options.CTEEncodingFormatUnset:                 "|i16",
	options.CTEEncodingFormatBinary:                "|i16b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|i16b",
	options.CTEEncodingFormatOctal:                 "|i16o",
	options.CTEEncodingFormatOctalZeroFilled:       "|i16o",
	options.CTEEncodingFormatHexadecimal:           "|i16x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|i16x",
}
var arrayHeadersInt32 = []string{
	options.CTEEncodingFormatUnset:                 "|i32",
	options.CTEEncodingFormatBinary:                "|i32b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|i32b",
	options.CTEEncodingFormatOctal:                 "|i32o",
	options.CTEEncodingFormatOctalZeroFilled:       "|i32o",
	options.CTEEncodingFormatHexadecimal:           "|i32x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|i32x",
}
var arrayHeadersInt64 = []string{
	options.CTEEncodingFormatUnset:                 "|i64",
	options.CTEEncodingFormatBinary:                "|i64b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|i64b",
	options.CTEEncodingFormatOctal:                 "|i64o",
	options.CTEEncodingFormatOctalZeroFilled:       "|i64o",
	options.CTEEncodingFormatHexadecimal:           "|i64x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|i64x",
}
var arrayHeadersFloat16 = []string{
	options.CTEEncodingFormatUnset:                 "|f16",
	options.CTEEncodingFormatBinary:                "|f16b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|f16b",
	options.CTEEncodingFormatOctal:                 "|f16o",
	options.CTEEncodingFormatOctalZeroFilled:       "|f16o",
	options.CTEEncodingFormatHexadecimal:           "|f16x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|f16x",
}
var arrayHeadersFloat32 = []string{
	options.CTEEncodingFormatUnset:                 "|f32",
	options.CTEEncodingFormatBinary:                "|f32b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|f32b",
	options.CTEEncodingFormatOctal:                 "|f32o",
	options.CTEEncodingFormatOctalZeroFilled:       "|f32o",
	options.CTEEncodingFormatHexadecimal:           "|f32x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|f32x",
}
var arrayHeadersFloat64 = []string{
	options.CTEEncodingFormatUnset:                 "|f64",
	options.CTEEncodingFormatBinary:                "|f64b",
	options.CTEEncodingFormatBinaryZeroFilled:      "|f64b",
	options.CTEEncodingFormatOctal:                 "|f64o",
	options.CTEEncodingFormatOctalZeroFilled:       "|f64o",
	options.CTEEncodingFormatHexadecimal:           "|f64x",
	options.CTEEncodingFormatHexadecimalZeroFilled: "|f64x",
}
