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
	"math/big"
	"strconv"
	"time"

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

const floatStringMaxByteCount = 24 // From strconv.FormatFloat()
const uintStringMaxByteCount = 21  // 18446744073709551616

type EncodeBuffer struct {
	buffer.StreamingWriteBuffer
}

func (_this *EncodeBuffer) WriteNA() {
	_this.AddString("@na")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteTrue() {
	_this.AddString("@true")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteFalse() {
	_this.AddString("@false")
	_this.Flush()
}

func (_this *EncodeBuffer) WritePosInfinity() {
	_this.AddString("@inf")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteNegInfinity() {
	_this.AddString("-@inf")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteQuietNan() {
	_this.AddString("@nan")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteSignalingNan() {
	_this.AddString("@snan")
	_this.Flush()
}

func (_this *EncodeBuffer) WriteVersion(value uint64) {
	_this.AddByte('c')
	_this.WritePositiveInt(value)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteNan(signaling bool) {
	if signaling {
		_this.WriteSignalingNan()
	} else {
		_this.WriteQuietNan()
	}
}

func (_this *EncodeBuffer) WriteBool(value bool) {
	if value {
		_this.WriteTrue()
	} else {
		_this.WriteFalse()
	}
}

func (_this *EncodeBuffer) WriteInt(value int64) {
	if value >= 0 {
		_this.WritePositiveInt(uint64(value))
	} else {
		_this.WriteNegativeInt(uint64(-value))
	}
}

func (_this *EncodeBuffer) WritePositiveInt(value uint64) {
	buff := _this.RequireBytes(uintStringMaxByteCount)[:0]
	used := strconv.AppendUint(buff, value, 10)
	_this.UseBytes(len(used))
	_this.Flush()
}

func (_this *EncodeBuffer) WriteNegativeInt(value uint64) {
	_this.AddByte('-')
	_this.WritePositiveInt(value)
}

func (_this *EncodeBuffer) WriteBigInt(value *big.Int) {
	if value == nil {
		_this.WriteNA()
		return
	}

	var buff [64]byte
	used := value.Append(buff[:0], 10)
	_this.AddBytes(used)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteFloat(value float64) {
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

	buff := _this.RequireBytes(floatStringMaxByteCount)[:0]
	used := strconv.AppendFloat(buff, value, 'g', -1, 64)
	_this.UseBytes(len(used))
	_this.Flush()
}

func (_this *EncodeBuffer) WriteFloatHexNoPrefix(value float64) {
	buff := _this.RequireBytes(floatStringMaxByteCount)[:0]
	used := strconv.AppendFloat(buff, value, 'x', -1, 64)
	copy(used, used[2:])
	if value < 0 {
		used[0] = '-'
	}
	_this.UseBytes(len(used) - 2)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteBigFloat(value *big.Float) {
	if value == nil {
		_this.WriteNA()
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
	_this.AddBytes(used)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteDecimalFloat(value compact_float.DFloat) {
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

	_this.AddString(value.Text('g'))
	_this.Flush()
}

func (_this *EncodeBuffer) WriteBigDecimalFloat(value *apd.Decimal) {
	if value == nil {
		_this.WriteNA()
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
		_this.AddBytes(used)
		_this.Flush()
	}
}

func (_this *EncodeBuffer) WriteUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("expected UUID length 16 but got %v", len(v)))
	}
	_this.AddFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
	_this.Flush()
}

func (_this *EncodeBuffer) WriteTime(value time.Time) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		_this.unexpectedError(err, value)
	}
	_this.WriteCompactTime(t)
}

func (_this *EncodeBuffer) WriteCompactTime(value compact_time.Time) {
	if value.IsZeroValue() {
		_this.WriteNA()
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
		_this.AddFmt("%d-%02d-%02d", value.Year, value.Month, value.Day)
	case compact_time.TypeTime:
		_this.AddFmt("%02d:%02d:%02d%s%s", value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	case compact_time.TypeTimestamp:
		_this.AddFmt("%d-%02d-%02d/%02d:%02d:%02d%s%s",
			value.Year, value.Month, value.Day, value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	default:
		panic(fmt.Errorf("unknown compact time type %v", value.TimeType))
	}
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMarkerBegin(id interface{}) {
	_this.AddFmt("&%v:", id)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteReference(id interface{}) {
	_this.AddFmt("$%v", id)
	_this.Flush()
}

func (_this *EncodeBuffer) WriteListBegin() {
	_this.AddByte('[')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteListEnd() {
	_this.AddByte(']')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMapBegin() {
	_this.AddByte('{')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMapEnd() {
	_this.AddByte('}')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMetaBegin() {
	_this.AddByte('(')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMetaEnd() {
	_this.AddByte(')')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMarkupBegin() {
	_this.AddByte('<')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMarkupContentsBegin() {
	_this.AddByte(',')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteMarkupEnd() {
	_this.AddByte('>')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteArrayBegin() {
	_this.AddByte('|')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteArrayEnd() {
	_this.AddByte('|')
	_this.Flush()
}

func (_this *EncodeBuffer) WriteCommentBegin() {
	_this.AddBytes([]byte{'/', '*'})
	_this.Flush()
}

func (_this *EncodeBuffer) WriteCommentEnd() {
	_this.AddBytes([]byte{'*', '/'})
	_this.Flush()
}

func (_this *EncodeBuffer) unexpectedError(err error, encoding interface{}) {
	panic(fmt.Errorf("unexpected error [%v] while encoding %v", err, encoding))
}
