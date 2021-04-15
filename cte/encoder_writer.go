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
	_this.WriteString("@na:")
}

func (_this *Writer) WriteNil() {
	_this.WriteString("@nil")
}

func (_this *Writer) WriteTrue() {
	_this.WriteString("@true")
}

func (_this *Writer) WriteFalse() {
	_this.WriteString("@false")
}

func (_this *Writer) WritePosInfinity() {
	_this.WriteString("@inf")
}

func (_this *Writer) WriteNegInfinity() {
	_this.WriteString("-@inf")
}

func (_this *Writer) WriteQuietNan() {
	_this.WriteString("@nan")
}

func (_this *Writer) WriteSignalingNan() {
	_this.WriteString("@snan")
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

func (_this *Writer) WriteFloatHexNoPrefix(value float64) {
	_this.ExpandBuffer(floatStringMaxByteCount)
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

func (_this *Writer) WriteUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("expected UUID length 16 but got %v", len(v)))
	}
	_this.WriteFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
}

func (_this *Writer) WriteTime(value time.Time) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		_this.unexpectedError(err, value)
	}
	_this.WriteCompactTime(t)
}

func (_this *Writer) WriteCompactTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.WriteNil()
		return
	}

	tz := func(v compact_time.Time) string {
		switch v.TimezoneType {
		case compact_time.TypeZero:
			return ""
		case compact_time.TypeAreaLocation, compact_time.TypeLocal:
			return fmt.Sprintf("/%s", v.LongAreaLocation)
		case compact_time.TypeLatitudeLongitude:
			return fmt.Sprintf("/%.2f/%.2f", float64(v.LatitudeHundredths)/100, float64(v.LongitudeHundredths)/100)
		default:
			panic(fmt.Errorf("unknown compact time timezone type %v", value.TimezoneType))
			return ""
		}
	}
	subsec := func(v compact_time.Time) string {
		if v.Nanosecond == 0 {
			return ""
		}

		str := strconv.FormatFloat(float64(v.Nanosecond)/float64(1000000000), 'f', 9, 64)
		for str[len(str)-1] == '0' {
			str = str[:len(str)-1]
		}
		return str[1:]
	}
	switch value.TimeType {
	case compact_time.TypeDate:
		_this.WriteFmt("%d-%02d-%02d", value.Year, value.Month, value.Day)
	case compact_time.TypeTime:
		_this.WriteFmt("%02d:%02d:%02d%s%s", value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	case compact_time.TypeTimestamp:
		_this.WriteFmt("%d-%02d-%02d/%02d:%02d:%02d%s%s",
			value.Year, value.Month, value.Day, value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	default:
		panic(fmt.Errorf("unknown compact time type %v", value.TimeType))
	}
}

func (_this *Writer) WriteQuotedStringBytes(value []byte) {
	if len(value) == 0 {
		_this.WriteBytes([]byte{'"', '"'})
		return
	}

	escapeCount := getEscapeCount(value)

	if escapeCount == 0 {
		_this.WriteByte('"')
		_this.WriteBytes(value)
		_this.WriteByte('"')
		return
	}

	_this.WriteEscapedQuotedStringBytes(value, escapeCount)
}

func (_this *Writer) WriteEscapedQuotedStringBytes(value []byte, escapeCount int) {
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	// data := _this.RequireBytes(len(value) + escapeCount*6 + 2)

	_this.WriteByte('"')
	for _, ch := range string(value) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeQuoted) {
			// TODO: render escapes into buffer to avoid allocs
			_this.WriteBytes(escapeCharQuoted(ch))
		} else {
			_this.WriteRune(ch)
		}
	}
	_this.WriteByte('"')
}

func (_this *Writer) WritePotentiallyEscapedStringArrayContents(value []byte) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesStringlikeArray(value) {
		_this.WriteBytes(value)
		return
	}

	// TODO: Encode directly rather than using bytes.Buffer
	var bb bytes.Buffer
	bb.Grow(len(value))
	for _, ch := range string(value) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeArray) {
			// Note: StringBuilder's WriteXYZ() always return nil errors
			bb.Write(escapeCharStringArray(ch))
		} else {
			bb.WriteRune(ch)
		}
	}
	_this.WriteBytes(bb.Bytes())

}

func (_this *Writer) WritePotentiallyEscapedMarkupContents(value []byte) {
	if len(value) == 0 {
		return
	}

	if !needsEscapesMarkup(value) {
		_this.WriteBytes(value)
		return
	}

	// TODO: Encode directly rather than using bytes.Buffer
	var bb bytes.Buffer
	bb.Grow(len(value))
	// Note: StringBuilder's WriteXYZ() always return nil errors
	for _, ch := range string(value) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeMarkup) {
			bb.Write(escapeCharMarkup(ch))
		} else {
			bb.WriteRune(ch)
		}
	}
	_this.WriteBytes(bb.Bytes())
}

var hexToChar = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
}

func (_this *Writer) WriteHexBytes(value []byte) {
	length := len(value) * 3
	_this.ExpandBuffer(length)
	dst := _this.Buffer
	for i := 0; i < len(value); i++ {
		b := value[i]
		dst[i*3] = ' '
		dst[i*3+1] = hexToChar[b>>4]
		dst[i*3+2] = hexToChar[b&15]
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

func (_this *Writer) WriteMapValueSeparator() {
	_this.WriteBytes([]byte{' ', '=', ' '})
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

func (_this *Writer) WriteCommentBegin() {
	_this.WriteBytes([]byte{'/', '*'})
}

func (_this *Writer) WriteCommentEnd() {
	_this.WriteBytes([]byte{'*', '/'})
}

func (_this *Writer) unexpectedError(err error, encoding interface{}) {
	panic(fmt.Errorf("unexpected error [%v] while encoding %v", err, encoding))
}
