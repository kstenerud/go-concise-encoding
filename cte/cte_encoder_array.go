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
	"github.com/kstenerud/go-concise-encoding/options"
)

type arrayEncoderEngineOld struct {
	engine                 *encoderEngine
	stream                 *EncodeBuffer
	addElementsFunc        func(b []byte)
	arrayCloseFunc         func()
	arrayElementWidth      int
	remainingChunkElements uint64
	moreChunksFollow       bool
	arrayChunkLeftover     []byte
	stringBuffer           []byte
	opts                   *options.CTEEncoderOptions
}

func (_this *arrayEncoderEngineOld) Init(engine *encoderEngine, opts *options.CTEEncoderOptions) {
	_this.engine = engine
	_this.stream = engine.stream
	_this.arrayChunkLeftover = make([]byte, 0, 16)
	_this.opts = opts
}

func (_this *arrayEncoderEngineOld) Reset() {
	_this.arrayChunkLeftover = _this.arrayChunkLeftover[:0]
	_this.stringBuffer = _this.stringBuffer[:0]
	_this.remainingChunkElements = 0
}

func (_this *arrayEncoderEngineOld) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.beginStringLikeArray(awaitingQuotedString,
			func(stringData []byte) {
				switch _this.engine.Awaiting {
				case awaitingMarkerID:
					_this.engine.CompleteMarker(string(stringData))
				case awaitingReferenceID:
					_this.engine.CompleteReference(string(stringData))
				case awaitingMarkupItem, awaitingMarkupFirstItem:
					_this.engine.BeginObject()
					_this.stream.AddBytes(asMarkupContents(stringData))
					_this.engine.CompleteObject()
				case awaitingCommentItem:
					_this.engine.AddCommentString(string(stringData))
				default:
					_this.engine.BeginObject()
					_this.stream.AddBytes(asPotentialQuotedString(stringData))
					_this.engine.CompleteObject()
				}
			})
	case events.ArrayTypeResourceID:
		_this.beginStringLikeArray(awaitingRID,
			func(stringData []byte) {
				if _this.engine.Awaiting == awaitingReferenceID {
					_this.engine.CompleteReference(string(asStringArray([]byte{'r'}, stringData)))
				} else {
					_this.engine.BeginObject()
					_this.encodeStringArray("r", stringData)
					_this.engine.CompleteObject()
				}
			})
	case events.ArrayTypeCustomBinary:
		// TODO: Remove these
		var hexToChar = [16]byte{
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
		}
		encodeHex := func(value []byte) {
			length := len(value) * 3
			dst := _this.stream.RequireBytes(length)
			dst[0] = '|'
			dst[1] = 'c'
			dst[2] = 'b'
			dst[3] = ' '
			dst = dst[3:]
			i := 0
			for ; i < len(value); i++ {
				b := value[i]
				dst[i*3] = ' '
				dst[i*3+1] = hexToChar[b>>4]
				dst[i*3+2] = hexToChar[b&15]
			}
			dst[i*3] = '|'
			_this.stream.UseBytes(length + 4)
		}

		_this.beginStringLikeArray(awaitingCustomBinary,
			func(stringData []byte) {
				_this.engine.BeginObject()
				encodeHex(stringData)
				_this.engine.CompleteObject()
			})
	case events.ArrayTypeCustomText:
		_this.beginStringLikeArray(awaitingCustomText,
			func(stringData []byte) {
				_this.engine.BeginObject()
				_this.encodeStringArray("ct", stringData)
				_this.engine.CompleteObject()
			})
	case events.ArrayTypeUint8:
		const elemWidth = 1
		opener := arrayHeadersUint8[_this.opts.DefaultFormats.Array.Uint8]
		format := arrayFormats8[_this.opts.DefaultFormats.Array.Uint8]
		_this.beginArray(awaitingArrayU8,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, data[0])
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeUint16:
		const elemWidth = 2
		opener := arrayHeadersUint16[_this.opts.DefaultFormats.Array.Uint16]
		format := arrayFormats16[_this.opts.DefaultFormats.Array.Uint16]
		_this.beginArray(awaitingArrayU16,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, uint(data[0])|(uint(data[1])<<8))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeUint32:
		const elemWidth = 4
		opener := arrayHeadersUint32[_this.opts.DefaultFormats.Array.Uint32]
		format := arrayFormats32[_this.opts.DefaultFormats.Array.Uint32]
		_this.beginArray(awaitingArrayU32,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, uint(data[0])|(uint(data[1])<<8)|(uint(data[2])<<16)|(uint(data[3])<<24))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeUint64:
		const elemWidth = 8
		opener := arrayHeadersUint64[_this.opts.DefaultFormats.Array.Uint64]
		format := arrayFormats64[_this.opts.DefaultFormats.Array.Uint64]
		_this.beginArray(awaitingArrayU64,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, uint64(data[0])|(uint64(data[1])<<8)|(uint64(data[2])<<16)|(uint64(data[3])<<24)|
						(uint64(data[4])<<32)|(uint64(data[5])<<40)|(uint64(data[6])<<48)|(uint64(data[7])<<56))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeInt8:
		const elemWidth = 1
		opener := arrayHeadersInt8[_this.opts.DefaultFormats.Array.Int8]
		format := arrayFormats8[_this.opts.DefaultFormats.Array.Int8]
		_this.beginArray(awaitingArrayI8,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, int8(data[0]))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeInt16:
		const elemWidth = 2
		opener := arrayHeadersInt16[_this.opts.DefaultFormats.Array.Int16]
		format := arrayFormats16[_this.opts.DefaultFormats.Array.Int16]
		_this.beginArray(awaitingArrayI16,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, int16(data[0])|(int16(data[1])<<8))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeInt32:
		const elemWidth = 4
		opener := arrayHeadersInt32[_this.opts.DefaultFormats.Array.Int32]
		format := arrayFormats32[_this.opts.DefaultFormats.Array.Int32]
		_this.beginArray(awaitingArrayI32,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, int32(data[0])|(int32(data[1])<<8)|(int32(data[2])<<16)|(int32(data[3])<<24))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeInt64:
		const elemWidth = 8
		opener := arrayHeadersInt64[_this.opts.DefaultFormats.Array.Int64]
		format := arrayFormats64[_this.opts.DefaultFormats.Array.Int64]
		_this.beginArray(awaitingArrayI64,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.AddFmt(format, int64(data[0])|(int64(data[1])<<8)|(int64(data[2])<<16)|(int64(data[3])<<24)|
						(int64(data[4])<<32)|(int64(data[5])<<40)|(int64(data[6])<<48)|(int64(data[7])<<56))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeFloat16:
		const elemWidth = 2
		opener := arrayHeadersFloat16[_this.opts.DefaultFormats.Array.Float16]
		if _this.opts.DefaultFormats.Array.Float16 == options.CTEEncodingFormatHexadecimal {
			_this.beginArray(awaitingArrayF16,
				elemWidth,
				opener,
				func(data []byte) {
					for len(data) > 0 {
						bits := (uint32(data[0]) << 16) | (uint32(data[1]) << 24)
						v := math.Float32frombits(bits)
						if v < 0 {
							_this.stream.AddString(" -")
							_this.stream.AddFmtStripped(3, "%x", v)
						} else {
							_this.stream.AddString(" ")
							_this.stream.AddFmtStripped(2, "%x", v)
						}
						data = data[elemWidth:]
					}
				})
			break
		}

		format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float16]
		_this.beginArray(awaitingArrayF16,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					bits := (uint32(data[0]) << 16) | (uint32(data[1]) << 24)
					_this.stream.AddFmt(format, math.Float32frombits(bits))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeFloat32:
		const elemWidth = 4
		opener := arrayHeadersFloat32[_this.opts.DefaultFormats.Array.Float32]
		if _this.opts.DefaultFormats.Array.Float32 == options.CTEEncodingFormatHexadecimal {
			_this.beginArray(awaitingArrayF32,
				elemWidth,
				opener,
				func(data []byte) {
					for len(data) > 0 {
						_this.stream.AddByte(' ')
						bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
						v := math.Float32frombits(bits)
						_this.stream.WriteFloatHexNoPrefix(float64(v))
						data = data[elemWidth:]
					}
				})
			break
		}

		format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float32]
		_this.beginArray(awaitingArrayF32,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					bits := uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
					_this.stream.AddFmt(format, math.Float32frombits(bits))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeFloat64:
		const elemWidth = 8
		opener := arrayHeadersFloat64[_this.opts.DefaultFormats.Array.Float64]
		if _this.opts.DefaultFormats.Array.Float64 == options.CTEEncodingFormatHexadecimal {
			_this.beginArray(awaitingArrayF64,
				elemWidth,
				opener,
				func(data []byte) {
					for len(data) > 0 {
						_this.stream.AddByte(' ')
						bits := uint64(data[0]) | (uint64(data[1]) << 8) | (uint64(data[2]) << 16) | (uint64(data[3]) << 24) |
							(uint64(data[4]) << 32) | (uint64(data[5]) << 40) | (uint64(data[6]) << 48) | (uint64(data[7]) << 56)
						v := math.Float64frombits(bits)
						_this.stream.WriteFloatHexNoPrefix(v)
						data = data[elemWidth:]
					}
				})
			break
		}

		format := arrayFormatsGeneral[_this.opts.DefaultFormats.Array.Float64]
		_this.beginArray(awaitingArrayF64,
			elemWidth,
			opener,
			func(data []byte) {
				for len(data) > 0 {
					bits := uint64(data[0]) | (uint64(data[1]) << 8) | (uint64(data[2]) << 16) | (uint64(data[3]) << 24) |
						(uint64(data[4]) << 32) | (uint64(data[5]) << 40) | (uint64(data[6]) << 48) | (uint64(data[7]) << 56)
					_this.stream.AddFmt(format, math.Float64frombits(bits))
					data = data[elemWidth:]
				}
			})
	case events.ArrayTypeUUID:
		const elemWidth = 16
		_this.beginArray(awaitingArrayUUID,
			elemWidth,
			"|u",
			func(data []byte) {
				for len(data) > 0 {
					_this.stream.WriteUUID(data)
					data = data[elemWidth:]
				}
			})
	default:
		panic(fmt.Errorf("%v: Unknown array type", arrayType))
	}
}

func (_this *arrayEncoderEngineOld) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.remainingChunkElements = elementCount
	_this.moreChunksFollow = moreChunksFollow

	if elementCount == 0 && !moreChunksFollow {
		_this.endArray()
	}
}

func (_this *arrayEncoderEngineOld) OnArrayData(data []byte) {
	if _this.arrayElementWidth > 1 {
		leftoverLength := len(_this.arrayChunkLeftover)
		if leftoverLength > 0 {
			fillCount := _this.arrayElementWidth - leftoverLength

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

		widthMask := _this.arrayElementWidth - 1
		remainderCount := len(data) & widthMask
		if remainderCount != 0 {
			_this.arrayChunkLeftover = append(_this.arrayChunkLeftover, data[len(data)-remainderCount:]...)
			data = data[:len(data)-remainderCount]
		}
	}
	_this.addElementsFunc(data)
	_this.remainingChunkElements -= uint64(len(data) / _this.arrayElementWidth)
	if _this.remainingChunkElements == 0 && !_this.moreChunksFollow {
		_this.endArray()
	}
}

// ============================================================================

// Utils

func (_this *arrayEncoderEngineOld) beginArray(newState awaiting, elementWidth int, opener string, addElementsFunc func([]byte)) {
	_this.engine.BeginArray(newState)
	_this.stream.AddString(opener)
	_this.arrayElementWidth = elementWidth
	_this.addElementsFunc = addElementsFunc
	_this.arrayCloseFunc = func() {
		_this.engine.CompleteArray()
	}
}

func (_this *arrayEncoderEngineOld) beginStringLikeArray(newState awaiting, closeFunc func(stringData []byte)) {
	_this.engine.BeginStringLikeArray(newState)
	_this.arrayElementWidth = 1
	_this.addElementsFunc = func(data []byte) { _this.appendStringbuffer(data) }
	_this.arrayCloseFunc = func() {
		_this.engine.EndStringLikeArray()
		closeFunc(_this.drainStringBuffer())
	}
}

func (_this *arrayEncoderEngineOld) endArray() {
	_this.arrayCloseFunc()
}

func (_this *arrayEncoderEngineOld) appendStringbuffer(data []byte) {
	_this.stringBuffer = append(_this.stringBuffer, data...)
}

func (_this *arrayEncoderEngineOld) drainStringBuffer() []byte {
	data := _this.stringBuffer
	_this.stringBuffer = _this.stringBuffer[:0]
	return data
}

func (_this *arrayEncoderEngineOld) encodeStringArray(arrayType string, contents []byte) {
	contents = asStringArrayContents(contents)
	contentsLen := len(contents)
	typeLen := len(arrayType)
	dst := _this.stream.RequireBytes(contentsLen)
	dst[0] = '|'
	dst = dst[1:]
	copy(dst, arrayType)
	dst = dst[typeLen:]
	dst[0] = ' '
	dst = dst[1:]
	copy(dst, contents)
	dst = dst[contentsLen:]
	dst[0] = '|'
	_this.stream.UseBytes(typeLen + len(contents) + 3)
}
