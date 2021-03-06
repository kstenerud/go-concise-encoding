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
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	ceio "github.com/kstenerud/go-concise-encoding/internal/io"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

const floatStringMaxByteCount = 24 // From strconv.FormatFloat()
const uintStringMaxByteCount = 21  // Max uint as string: "18446744073709551616"

type Writer struct {
	ceio.Writer
}

func (_this *Writer) WriteConcat() {
	_this.WriteByte(':')
}

func (_this *Writer) WriteNA() {
	_this.WriteString("na:")
}

func (_this *Writer) WriteNil() {
	_this.WriteString("nil")
}

func (_this *Writer) WriteTrue() {
	_this.WriteString("true")
}

func (_this *Writer) WriteFalse() {
	_this.WriteString("false")
}

func (_this *Writer) WritePosInfinity() {
	_this.WriteString("inf")
}

func (_this *Writer) WriteNegInfinity() {
	_this.WriteString("-inf")
}

func (_this *Writer) WriteQuietNan() {
	_this.WriteString("nan")
}

func (_this *Writer) WriteSignalingNan() {
	_this.WriteString("snan")
}

func (_this *Writer) WriteVersion(value uint64) {
	_this.WriteByte('c')
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
	_this.FlushBuffer(len(used))
}

func (_this *Writer) WriteNegativeInt(value uint64) {
	_this.WriteByte('-')
	_this.WritePositiveInt(value)
}

func (_this *Writer) WriteBigInt(value *big.Int) {
	if value == nil {
		_this.WriteNil()
		return
	}

	var buff [64]byte
	used := value.Append(buff[:0], 10)
	_this.WriteBytes(used)
}

func (_this *Writer) WriteFloat(value float64) {
	if math.IsNaN(value) {
		if common.IsSignalingNan(value) {
			_this.WriteSignalingNan()
		} else {
			_this.WriteQuietNan()
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

	_this.ExpandBuffer(floatStringMaxByteCount)
	used := strconv.AppendFloat(_this.Buffer[:0], value, 'g', -1, 64)
	_this.FlushBuffer(len(used))
}

func (_this *Writer) WriteFloat32UsingFormat(value float32, format string) {
	f64 := float64(value)
	if math.IsNaN(f64) {
		if common.IsSignalingNan(f64) {
			_this.WriteSignalingNan()
		} else {
			_this.WriteQuietNan()
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

	_this.WriteFmt(format, value)
}

func (_this *Writer) WriteFloatUsingFormat(value float64, format string) {
	if math.IsNaN(value) {
		if common.IsSignalingNan(value) {
			_this.WriteSignalingNan()
		} else {
			_this.WriteQuietNan()
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

	_this.WriteFmt(format, value)
}

func (_this *Writer) WriteFloatHexNoPrefix(value float64) {
	_this.ExpandBuffer(floatStringMaxByteCount)
	if math.IsNaN(value) {
		if common.IsSignalingNan(value) {
			_this.WriteSignalingNan()
		} else {
			_this.WriteQuietNan()
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
	_this.FlushBufferPortion(start, end)
}

func (_this *Writer) WriteBigFloat(value *big.Float) {
	if value == nil {
		_this.WriteNil()
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

	var buff [64]byte
	used := value.Append(buff[:0], 'g', conversions.BitsToDecimalDigits(int(value.Prec())))
	_this.WriteBytes(used)
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

	_this.WriteString(value.Text('g'))
}

func (_this *Writer) WriteBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.WriteNil()
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
		_this.WriteBytes(used)
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
	_this.WriteByte('-')
	_this.WriteHexByte(v[4])
	_this.WriteHexByte(v[5])
	_this.WriteByte('-')
	_this.WriteHexByte(v[6])
	_this.WriteHexByte(v[7])
	_this.WriteByte('-')
	_this.WriteHexByte(v[8])
	_this.WriteHexByte(v[9])
	_this.WriteByte('-')
	_this.WriteHexByte(v[10])
	_this.WriteHexByte(v[11])
	_this.WriteHexByte(v[12])
	_this.WriteHexByte(v[13])
	_this.WriteHexByte(v[14])
	_this.WriteHexByte(v[15])
}

func (_this *Writer) WriteTime(value time.Time) {
	_this.WriteCompactTime(compact_time.AsCompactTime(value))
}

func (_this *Writer) WriteCompactTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.WriteNil()
		return
	}

	if value.Type == compact_time.TimeTypeDate || value.Type == compact_time.TimeTypeTimestamp {
		_this.WriteDecimalInt(int64(value.Year))
		_this.WriteByte('-')
		_this.WriteDecimalUintDigits(uint64(value.Month), 2)
		_this.WriteByte('-')
		_this.WriteDecimalUintDigits(uint64(value.Day), 2)
	}

	if value.Type == compact_time.TimeTypeDate {
		return
	}

	if value.Type == compact_time.TimeTypeTimestamp {
		_this.WriteByte('/')
	}

	_this.WriteDecimalUintDigits(uint64(value.Hour), 2)
	_this.WriteByte(':')
	_this.WriteDecimalUintDigits(uint64(value.Minute), 2)
	_this.WriteByte(':')
	_this.WriteDecimalUintDigits(uint64(value.Second), 2)

	if value.Nanosecond != 0 {
		_this.WriteByte('.')
		_this.WriteDecimalUintLeftLoaded(uint64(value.Nanosecond), 9)
	}

	switch value.Timezone.Type {
	case compact_time.TimezoneTypeUTC:
	case compact_time.TimezoneTypeAreaLocation, compact_time.TimezoneTypeLocal:
		_this.WriteByte('/')
		_this.WriteString(value.Timezone.LongAreaLocation)
	case compact_time.TimezoneTypeLatitudeLongitude:
		_this.WriteFmt("/%.2f/%.2f", float64(value.Timezone.LatitudeHundredths)/100, float64(value.Timezone.LongitudeHundredths)/100)
	case compact_time.TimezoneTypeUTCOffset:
		minutes := int(value.Timezone.MinutesOffsetFromUTC)
		sign := '+'
		if value.Timezone.MinutesOffsetFromUTC < 0 {
			sign = '-'
			minutes = -minutes
		}
		_this.WriteFmt("%c%02d%02d", sign, minutes/60, minutes%60)
	default:
		panic(fmt.Errorf("unknown compact time timezone type %v", value.Timezone.Type))
	}
}

func (_this *Writer) WriteQuotedString(value string) {
	if len(value) == 0 {
		_this.WriteString(`""`)
		return
	}

	escapeCount := getEscapeCount(value)

	if escapeCount == 0 {
		_this.WriteByte('"')
		_this.WriteString(value)
		_this.WriteByte('"')
		return
	}

	_this.WriteEscapedQuotedString(value, escapeCount)
}

func (_this *Writer) WriteQuotedStringBytes(value []byte) {
	if len(value) == 0 {
		_this.WriteString(`""`)
		return
	}

	escapeCount := getEscapeCountBytes(value)

	if escapeCount == 0 {
		_this.WriteByte('"')
		_this.WriteBytes(value)
		_this.WriteByte('"')
		return
	}

	_this.WriteEscapedQuotedStringBytes(value, escapeCount)
}

func (_this *Writer) WriteEscapedQuotedString(value string, escapeCount int) {
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	// data := _this.RequireBytes(len(value) + escapeCount*6 + 2)

	_this.WriteByte('"')
	for _, ch := range value {
		if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
			// TODO: render escapes into buffer to avoid allocs
			_this.WriteBytes(escapeCharQuoted(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
	_this.WriteByte('"')
}

func (_this *Writer) WriteEscapedQuotedStringBytes(value []byte, escapeCount int) {
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	// data := _this.RequireBytes(len(value) + escapeCount*6 + 2)

	_this.WriteByte('"')
	for _, ch := range string(value) {
		if !chars.IsRuneSafeFor(ch, chars.SafetyString) {
			// TODO: render escapes into buffer to avoid allocs
			_this.WriteBytes(escapeCharQuoted(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
	_this.WriteByte('"')
}

func (_this *Writer) WritePotentiallyEscapedStringArrayContentsBytes(value []byte) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesStringlikeArrayBytes(value) {
		_this.WriteBytes(value)
		return
	}

	for _, ch := range string(value) {
		if !chars.IsRuneSafeFor(ch, chars.SafetyArray) {
			_this.WriteBytes(escapeCharQuoted(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
}

func (_this *Writer) WritePotentiallyEscapedStringArrayContents(value string) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesStringlikeArray(value) {
		_this.WriteString(value)
		return
	}

	for _, ch := range value {
		if !chars.IsRuneSafeFor(ch, chars.SafetyArray) {
			_this.WriteBytes(escapeCharQuoted(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
}

func (_this *Writer) WritePotentiallyEscapedMarkupContents(value string) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesMarkup(value) {
		_this.WriteString(value)
		return
	}

	for _, ch := range value {
		if !chars.IsRuneSafeFor(ch, chars.SafetyMarkup) {
			_this.WriteBytes(escapeCharMarkup(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
}

func (_this *Writer) WritePotentiallyEscapedMarkupContentsBytes(value []byte) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesMarkupBytes(value) {
		_this.WriteBytes(value)
		return
	}

	for _, ch := range string(value) {
		if !chars.IsRuneSafeFor(ch, chars.SafetyMarkup) {
			_this.WriteBytes(escapeCharMarkup(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
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
	_this.FlushBuffer(length)
}

func (_this *Writer) WriteMarkerBegin(id []byte) {
	_this.WriteByte('&')
	_this.WriteBytes(id)
	_this.WriteByte(':')
}

func (_this *Writer) WriteReferenceBegin() {
	_this.WriteByte('$')
}

func (_this *Writer) WriteReference(id []byte) {
	_this.WriteByte('$')
	_this.WriteBytes(id)
}

func (_this *Writer) WriteConstant(id []byte) {
	_this.WriteByte('#')
	_this.WriteBytes(id)
}

func (_this *Writer) WriteSeparator() {
	_this.WriteByte(':')
}

func (_this *Writer) WriteListBegin() {
	_this.WriteByte('[')
}

func (_this *Writer) WriteListEnd() {
	_this.WriteByte(']')
}

func (_this *Writer) WriteMapBegin() {
	_this.WriteByte('{')
}

var mapValueSeparator = []byte{' ', '=', ' '}

func (_this *Writer) WriteMapValueSeparator() {
	_this.WriteBytes(mapValueSeparator)
}

func (_this *Writer) WriteMapEnd() {
	_this.WriteByte('}')
}

func (_this *Writer) WriteMarkupBegin(id []byte) {
	_this.WriteByte('<')
	_this.WriteBytes(id)
}

func (_this *Writer) WriteMarkupKeySeparator() {
	_this.WriteByte(' ')
}

func (_this *Writer) WriteMarkupValueSeparator() {
	_this.WriteByte('=')
}

func (_this *Writer) WriteMarkupContentsBegin() {
	_this.WriteByte(',')
}

func (_this *Writer) WriteMarkupEnd() {
	_this.WriteByte('>')
}

func (_this *Writer) WriteArrayBegin() {
	_this.WriteByte('|')
}

func (_this *Writer) WriteArrayEnd() {
	_this.WriteByte('|')
}

func (_this *Writer) WriteRelationshipBegin() {
	_this.WriteByte('(')
}

func (_this *Writer) WriteRelationshipEnd() {
	_this.WriteByte(')')
}

var commentBegin = []byte{'/', '*'}
var commentEnd = []byte{'*', '/'}

func (_this *Writer) WriteCommentBegin() {
	_this.WriteBytes(commentBegin)
}

func (_this *Writer) WriteCommentEnd() {
	_this.WriteBytes(commentEnd)
}

func (_this *Writer) unexpectedError(err error, encoding interface{}) {
	panic(fmt.Errorf("unexpected error [%v] while encoding %v", err, encoding))
}
