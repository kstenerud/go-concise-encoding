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
	"bytes"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

const floatStringMaxByteCount = 24 // From strconv.FormatFloat()
const uintStringMaxByteCount = 21  // Max uint as string: "18446744073709551616"

type Writer struct {
	writer       io.Writer
	stringWriter io.StringWriter
	adapter      StringWriterAdapter
	Buffer       []byte
	Column       int
}

func NewWriter() *Writer {
	_this := &Writer{}
	_this.Init()
	return _this
}

func (_this *Writer) Init() {
	_this.Buffer = make([]byte, writerStartBufferSize)
	_this.adapter.Init(_this)
	_this.Column = 0
}

// Set the actual io.Writer (or io.StringWriter) that will write the data.
func (_this *Writer) SetWriter(writer io.Writer) {
	_this.writer = writer
	if sw, ok := _this.writer.(io.StringWriter); ok {
		_this.stringWriter = sw
	} else {
		_this.stringWriter = &_this.adapter
	}
}

// Make sure the internal buffer is at least this big.
func (_this *Writer) ExpandBuffer(size int) {
	if len(_this.Buffer) < size {
		_this.Buffer = make([]byte, size*2)
	}
}

// Flush count bytes from the start of the buffer
func (_this *Writer) FlushBufferNotLF(count int) {
	_this.writeBytes(_this.Buffer[:count])
}

// Flush the buffer from the start to the end position
func (_this *Writer) FlushBufferPortionNotLF(start, end int) {
	_this.writeBytes(_this.Buffer[start:end])
}

func (_this *Writer) WriteLF() {
	_this.Buffer[0] = '\n'
	_this.FlushBufferNotLF(1)
	_this.Column = 0
}

func (_this *Writer) WriteByteNotLF(b byte) {
	_this.Buffer[0] = b
	_this.FlushBufferNotLF(1)
	_this.Column++
}

func (_this *Writer) WriteRuneNotLF(r rune) {
	n := utf8.EncodeRune(_this.Buffer, r)
	_this.FlushBufferNotLF(n)
	_this.Column += n
}

func (_this *Writer) WriteRunePossibleLF(r rune) {
	n := utf8.EncodeRune(_this.Buffer, r)
	_this.FlushBufferNotLF(n)
	if r == '\n' {
		_this.Column = 0
	} else {
		_this.Column += n
	}
}

func (_this *Writer) WriteBytesNotLF(b []byte) {
	_this.writeBytes(b)
	_this.Column += len(b)
}

func (_this *Writer) WriteBytesPossibleLF(b []byte) {
	_this.writeBytes(b)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == '\n' {
			_this.Column = len(b) - (i + 1)
			break
		}
	}
}

func (_this *Writer) WriteStringNotLF(str string) {
	if _, err := _this.stringWriter.WriteString(str); err != nil {
		panic(err)
	}
	_this.Column += len(str)
}

func (_this *Writer) WriteStringPossibleLF(str string) {
	if _, err := _this.stringWriter.WriteString(str); err != nil {
		panic(err)
	}
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '\n' {
			_this.Column = len(str) - (i + 1)
			break
		}
	}
}

func (_this *Writer) writeBytes(b []byte) {
	if _, err := _this.writer.Write(b); err != nil {
		panic(err)
	}
}

func (_this *Writer) WriteDecimalInt(value int64) {
	uintValue := uint64(value)
	if value < 0 {
		_this.WriteByteNotLF('-')
		uintValue = uint64(-value)
	}
	_this.WriteDecimalUint(uintValue)
}

func (_this *Writer) WriteDecimalUint(value uint64) {
	const maxUint64Digits = 21 // 18446744073709551615
	_this.ExpandBuffer(maxUint64Digits)
	var start int
	for start = maxUint64Digits - 1; start >= 0; start-- {
		_this.Buffer[start] = byte('0' + value%10)
		value /= 10
		if value == 0 {
			break
		}
	}
	_this.FlushBufferPortionNotLF(start, maxUint64Digits)
}

func (_this *Writer) WriteDecimalUintLeftLoaded(value uint64, digitCount int) {
	_this.ExpandBuffer(digitCount)
	lastNonZero := 0
	for i := digitCount - 1; i >= 0; i-- {
		digit := value % 10
		if lastNonZero == 0 && digit != 0 {
			lastNonZero = i
		}
		_this.Buffer[i] = byte('0' + digit)
		value /= 10
	}
	_this.FlushBufferNotLF(lastNonZero + 1)
}

func (_this *Writer) WriteDecimalUintDigits(value uint64, digits int) {
	_this.ExpandBuffer(digits)
	for i := digits - 1; i >= 0; i-- {
		_this.Buffer[i] = byte('0' + value%10)
		value /= 10
	}
	_this.FlushBufferNotLF(digits)
}

func (_this *Writer) WriteHexByte(value byte) {
	_this.ExpandBuffer(2)
	_this.Buffer[0] = chars.HexChars[value>>4]
	_this.Buffer[1] = chars.HexChars[value&15]
	_this.FlushBufferNotLF(2)
}

func (_this *Writer) WriteFmtNotLF(format string, args ...interface{}) {
	_this.WriteStringNotLF(fmt.Sprintf(format, args...))
}

// Add a formatted string, but strip the specified number of characters from the
// beginning of the result before adding.
func (_this *Writer) WriteFmtStrippedNotLF(stripByteCount int, format string, args ...interface{}) {
	_this.WriteStringNotLF(fmt.Sprintf(format, args...)[stripByteCount:])
}

func (_this *Writer) WriteConcat() {
	_this.WriteByteNotLF(':')
}

func (_this *Writer) WriteNull() {
	_this.WriteStringNotLF("null")
}

func (_this *Writer) WriteTrue() {
	_this.WriteStringNotLF("true")
}

func (_this *Writer) WriteFalse() {
	_this.WriteStringNotLF("false")
}

func (_this *Writer) WritePosInfinity() {
	_this.WriteStringNotLF("inf")
}

func (_this *Writer) WriteNegInfinity() {
	_this.WriteStringNotLF("-inf")
}

func (_this *Writer) WriteQuietNan() {
	_this.WriteStringNotLF("nan")
}

func (_this *Writer) WriteSignalingNan() {
	_this.WriteStringNotLF("snan")
}

func (_this *Writer) WriteVersion(value uint64) {
	_this.WritePositiveInt(value)
}

func (_this *Writer) WriteNan(signaling bool) {
	if signaling {
		_this.WriteSignalingNan()
	} else {
		_this.WriteQuietNan()
	}
}

func (_this *Writer) WriteBool(value bool) {
	if value {
		_this.WriteTrue()
	} else {
		_this.WriteFalse()
	}
}

func (_this *Writer) WriteInt(value int64) {
	if value >= 0 {
		_this.WritePositiveInt(uint64(value))
	} else {
		_this.WriteNegativeInt(uint64(-value))
	}
}

func (_this *Writer) WritePositiveInt(value uint64) {
	_this.ExpandBuffer(uintStringMaxByteCount)
	used := strconv.AppendUint(_this.Buffer[:0], value, 10)
	_this.FlushBufferNotLF(len(used))
}

func (_this *Writer) WriteNegativeInt(value uint64) {
	_this.WriteByteNotLF('-')
	_this.WritePositiveInt(value)
}

func (_this *Writer) WriteBigInt(value *big.Int) {
	if value == nil {
		_this.WriteNull()
		return
	}

	var buff [64]byte
	used := value.Append(buff[:0], 10)
	_this.WriteBytesNotLF(used)
}

func (_this *Writer) WriteFloat(value float64) {
	if math.IsNaN(value) {
		if common.HasQuietNanBitSet64(value) {
			_this.WriteQuietNan()
		} else {
			_this.WriteSignalingNan()
		}
		return
	}
	if math.IsInf(value, 0) {
		if value < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}
	if value == 0 {
		_this.ExpandBuffer(2)
		used := strconv.AppendFloat(_this.Buffer[:0], value, 'g', -1, 64)
		_this.FlushBufferNotLF(len(used))
		return
	}

	_this.ExpandBuffer(floatStringMaxByteCount)
	used := strconv.AppendFloat(_this.Buffer[:0], value, 'x', -1, 64)
	_this.FlushBufferNotLF(len(used))
}

func (_this *Writer) WriteFloat32UsingFormat(value float32, format string) {
	f64 := float64(value)
	if math.IsNaN(f64) {
		if common.HasQuietNanBitSet64(f64) {
			_this.WriteQuietNan()
		} else {
			_this.WriteSignalingNan()
		}
		return
	}
	if math.IsInf(f64, 0) {
		if f64 < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}

	_this.WriteFmtNotLF(format, value)
}

func (_this *Writer) WriteFloatUsingFormat(value float64, format string) {
	if math.IsNaN(value) {
		if common.HasQuietNanBitSet64(value) {
			_this.WriteQuietNan()
		} else {
			_this.WriteSignalingNan()
		}
		return
	}
	if math.IsInf(value, 0) {
		if value < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}

	_this.WriteFmtNotLF(format, value)
}

func (_this *Writer) WriteFloatHexNoPrefix(value float64) {
	_this.ExpandBuffer(floatStringMaxByteCount)
	if math.IsNaN(value) {
		if common.HasQuietNanBitSet64(value) {
			_this.WriteQuietNan()
		} else {
			_this.WriteSignalingNan()
		}
		return
	}
	if math.IsInf(value, 0) {
		if value < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}

	used := strconv.AppendFloat(_this.Buffer[:0], value, 'x', -1, 64)
	end := len(used)
	if bytes.HasSuffix(used, []byte("p+00")) {
		end -= 4
	}
	start := 2
	if value < 0 {
		used[start] = '-'
	}
	_this.FlushBufferPortionNotLF(start, end)
}

func (_this *Writer) WriteBigFloat(value *big.Float) {
	if value == nil {
		_this.WriteNull()
		return
	}
	if value.IsInf() {
		if value.Sign() < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}
	asFloat, accuracy := value.Float64()
	if accuracy == big.Exact && asFloat == 0 {
		_this.WriteFloat(asFloat)
		return
	}

	var buff [64]byte
	used := value.Append(buff[:0], 'x', -1)
	if len(used) > 3 {
		end := len(used) - 4
		if used[end] == 'p' &&
			// +-
			used[end+2] == '0' &&
			used[end+3] == '0' {
			used = used[:end]
		}
	}

	_this.WriteBytesNotLF(used)
}

func (_this *Writer) WriteDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			_this.WriteSignalingNan()
		} else {
			_this.WriteQuietNan()
		}
		return
	}
	if value.IsInfinity() {
		if value.IsNegativeInfinity() {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
		return
	}

	_this.WriteStringNotLF(value.Text('g'))
}

func (_this *Writer) WriteBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.WriteNull()
		return
	}
	switch value.Form {
	case apd.NaN:
		_this.WriteQuietNan()
	case apd.NaNSignaling:
		_this.WriteSignalingNan()
	case apd.Infinite:
		if value.Sign() < 0 {
			_this.WriteNegInfinity()
		} else {
			_this.WritePosInfinity()
		}
	default:
		var buff [64]byte
		used := value.Append(buff[:0], 'g')
		_this.WriteBytesNotLF(used)
	}
}

func (_this *Writer) WriteUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("expected UID length 16 but got %v", len(v)))
	}

	_this.WriteHexByte(v[0])
	_this.WriteHexByte(v[1])
	_this.WriteHexByte(v[2])
	_this.WriteHexByte(v[3])
	_this.WriteByteNotLF('-')
	_this.WriteHexByte(v[4])
	_this.WriteHexByte(v[5])
	_this.WriteByteNotLF('-')
	_this.WriteHexByte(v[6])
	_this.WriteHexByte(v[7])
	_this.WriteByteNotLF('-')
	_this.WriteHexByte(v[8])
	_this.WriteHexByte(v[9])
	_this.WriteByteNotLF('-')
	_this.WriteHexByte(v[10])
	_this.WriteHexByte(v[11])
	_this.WriteHexByte(v[12])
	_this.WriteHexByte(v[13])
	_this.WriteHexByte(v[14])
	_this.WriteHexByte(v[15])
}

func (_this *Writer) WriteTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.WriteNull()
		return
	}

	if value.Type == compact_time.TimeTypeDate || value.Type == compact_time.TimeTypeTimestamp {
		_this.WriteDecimalInt(int64(value.Year))
		_this.WriteByteNotLF('-')
		_this.WriteDecimalUintDigits(uint64(value.Month), 2)
		_this.WriteByteNotLF('-')
		_this.WriteDecimalUintDigits(uint64(value.Day), 2)
	}

	switch value.Type {
	case compact_time.TimeTypeDate:
		return
	case compact_time.TimeTypeTimestamp:
		_this.WriteByteNotLF('/')
	}

	_this.WriteDecimalUintDigits(uint64(value.Hour), 2)
	_this.WriteByteNotLF(':')
	_this.WriteDecimalUintDigits(uint64(value.Minute), 2)
	_this.WriteByteNotLF(':')
	_this.WriteDecimalUintDigits(uint64(value.Second), 2)

	if value.Nanosecond != 0 {
		_this.WriteByteNotLF('.')
		_this.WriteDecimalUintLeftLoaded(uint64(value.Nanosecond), 9)
	}

	switch value.Timezone.Type {
	case compact_time.TimezoneTypeUTC:
	case compact_time.TimezoneTypeAreaLocation, compact_time.TimezoneTypeLocal:
		_this.WriteByteNotLF('/')
		_this.WriteStringNotLF(value.Timezone.LongAreaLocation)
	case compact_time.TimezoneTypeLatitudeLongitude:
		_this.WriteFmtNotLF("/%.2f/%.2f", float64(value.Timezone.LatitudeHundredths)/100, float64(value.Timezone.LongitudeHundredths)/100)
	case compact_time.TimezoneTypeUTCOffset:
		minutes := int(value.Timezone.MinutesOffsetFromUTC)
		sign := '+'
		if value.Timezone.MinutesOffsetFromUTC < 0 {
			sign = '-'
			minutes = -minutes
		}
		_this.WriteFmtNotLF("%c%02d%02d", sign, minutes/60, minutes%60)
	default:
		panic(fmt.Errorf("BUG: unknown compact time timezone type %v", value.Timezone.Type))
	}
}

func (_this *Writer) WriteQuotedString(mayContainLF bool, value string) {
	if len(value) == 0 {
		_this.WriteStringNotLF(`""`)
		return
	}

	escapeCount := getEscapeCount(value)

	if escapeCount == 0 {
		_this.WriteByteNotLF('"')
		if mayContainLF {
			_this.WriteStringPossibleLF(value)
		} else {
			_this.WriteStringNotLF(value)
		}
		_this.WriteByteNotLF('"')
		return
	}

	_this.WriteEscapedQuotedString(mayContainLF, value, escapeCount)
}

func (_this *Writer) WriteQuotedStringBytes(mayContainLF bool, value []byte) {
	if len(value) == 0 {
		_this.WriteStringNotLF(`""`)
		return
	}

	escapeCount := getEscapeCountBytes(value)

	if escapeCount == 0 {
		_this.WriteByteNotLF('"')
		if mayContainLF {
			_this.WriteBytesPossibleLF(value)
		} else {
			_this.WriteBytesNotLF(value)
		}
		_this.WriteByteNotLF('"')
		return
	}

	_this.WriteEscapedQuotedStringBytes(mayContainLF, value, escapeCount)
}

func (_this *Writer) WriteEscapedQuotedString(mayContainLF bool, value string, escapeCount int) {
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	// data := _this.RequireBytes(len(value) + escapeCount*6 + 2)

	_this.WriteByteNotLF('"')
	if mayContainLF {
		for _, ch := range value {
			if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
				// TODO: render escapes into buffer to avoid allocs
				_this.WriteBytesNotLF(escapeCharQuoted(ch))
			} else {
				_this.WriteRunePossibleLF(ch)
			}
		}
	} else {
		for _, ch := range value {
			if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
				// TODO: render escapes into buffer to avoid allocs
				_this.WriteBytesNotLF(escapeCharQuoted(ch))
			} else {
				_this.WriteRuneNotLF(ch)
			}
		}
	}
	_this.WriteByteNotLF('"')
}

func (_this *Writer) WriteEscapedQuotedStringBytes(mayContainLF bool, value []byte, escapeCount int) {
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	// data := _this.RequireBytes(len(value) + escapeCount*6 + 2)

	_this.WriteByteNotLF('"')
	if mayContainLF {
		for _, ch := range string(value) {
			if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
				// TODO: render escapes into buffer to avoid allocs
				_this.WriteBytesNotLF(escapeCharQuoted(ch))
			} else {
				_this.WriteRunePossibleLF(ch)
			}
		}
	} else {
		for _, ch := range string(value) {
			if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
				// TODO: render escapes into buffer to avoid allocs
				_this.WriteBytesNotLF(escapeCharQuoted(ch))
			} else {
				_this.WriteRuneNotLF(ch)
			}
		}
	}
	_this.WriteByteNotLF('"')
}

func (_this *Writer) WriteHexBytes(value []byte) {
	length := len(value) * 3
	_this.ExpandBuffer(length)
	dst := _this.Buffer
	for i := 0; i < len(value); i++ {
		b := value[i]
		dst[i*3] = ' '
		dst[i*3+1] = chars.HexChars[b>>4]
		dst[i*3+2] = chars.HexChars[b&15]
	}
	_this.FlushBufferNotLF(length)
}

func (_this *Writer) WriteMarkerBegin(id []byte) {
	_this.WriteByteNotLF('&')
	_this.WriteBytesNotLF(id)
	_this.WriteByteNotLF(':')
}

func (_this *Writer) WriteRemoteReferenceBegin() {
	_this.WriteByteNotLF('$')
}

func (_this *Writer) WriteLocalReference(id []byte) {
	_this.WriteByteNotLF('$')
	_this.WriteBytesNotLF(id)
}

func (_this *Writer) WriteSeparator() {
	_this.WriteByteNotLF(':')
}

func (_this *Writer) WriteListBegin() {
	_this.WriteByteNotLF('[')
}

func (_this *Writer) WriteListEnd() {
	_this.WriteByteNotLF(']')
}

func (_this *Writer) WriteMapBegin() {
	_this.WriteByteNotLF('{')
}

var mapValueSeparator = []byte{' ', '=', ' '}

func (_this *Writer) WriteMapValueSeparator() {
	_this.WriteBytesNotLF(mapValueSeparator)
}

func (_this *Writer) WriteMapEnd() {
	_this.WriteByteNotLF('}')
}

func (_this *Writer) WriteStructTemplateBegin(id []byte) {
	_this.WriteByteNotLF('@')
	_this.WriteBytesNotLF(id)
	_this.WriteByteNotLF('<')
}

func (_this *Writer) WriteStructTemplateEnd() {
	_this.WriteByteNotLF('>')
}

func (_this *Writer) WriteStructInstanceBegin(id []byte) {
	_this.WriteByteNotLF('@')
	_this.WriteBytesNotLF(id)
	_this.WriteByteNotLF('(')
}

func (_this *Writer) WriteStructInstanceEnd() {
	_this.WriteByteNotLF(')')
}

func (_this *Writer) WriteArrayBegin() {
	_this.WriteByteNotLF('|')
}

func (_this *Writer) WriteArrayEnd() {
	_this.WriteByteNotLF('|')
}

func (_this *Writer) WriteEdgeBegin() {
	_this.WriteStringNotLF("@(")
}

func (_this *Writer) WriteEdgeEnd() {
	_this.WriteByteNotLF(')')
}

func (_this *Writer) WriteNodeBegin() {
	_this.WriteStringNotLF("(")
}

func (_this *Writer) WriteNodeEnd() {
	_this.WriteByteNotLF(')')
}

var commentBeginMultiline = []byte{'/', '*'}
var commentEndMultiline = []byte{'*', '/'}
var commentBeginSingle = []byte{'/', '/'}

func (_this *Writer) WriteCommentBegin(isMultiline bool) {
	if isMultiline {
		_this.WriteBytesNotLF(commentBeginMultiline)
	} else {
		_this.WriteBytesNotLF(commentBeginSingle)
	}
}

func (_this *Writer) WriteCommentEnd(isMultiline bool) {
	if isMultiline {
		_this.WriteBytesNotLF(commentEndMultiline)
	}
	// Don't write the single-line comment terminator because the decorator will do it.
}

// ============================================================================

const writerStartBufferSize = 32

type StringWriterAdapter struct {
	writer *Writer
}

func (_this *StringWriterAdapter) Init(writer *Writer) {
	_this.writer = writer
}

func (_this *StringWriterAdapter) WriteString(str string) (n int, err error) {
	n = len(str)
	_this.writer.ExpandBuffer(n)
	copy(_this.writer.Buffer, str)
	_this.writer.FlushBufferNotLF(n)
	return
}
