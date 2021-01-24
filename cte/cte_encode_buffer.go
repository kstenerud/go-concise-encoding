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

type CTEEncodeBuffer struct {
	buffer.StreamingWriteBuffer
}

func (_this *CTEEncodeBuffer) WriteNA() {
	_this.AddString("@na")
}

func (_this *CTEEncodeBuffer) WriteTrue() {
	_this.AddString("@true")
}

func (_this *CTEEncodeBuffer) WriteFalse() {
	_this.AddString("@false")
}

func (_this *CTEEncodeBuffer) WritePosInfinity() {
	_this.AddString("@inf")
}

func (_this *CTEEncodeBuffer) WriteNegInfinity() {
	_this.AddString("-@inf")
}

func (_this *CTEEncodeBuffer) WriteQuietNan() {
	_this.AddString("@nan")
}

func (_this *CTEEncodeBuffer) WriteSignalingNan() {
	_this.AddString("@snan")
}

func (_this *CTEEncodeBuffer) WriteVersion(value uint64) {
	_this.AddByte('c')
	_this.WritePositiveInt(value)
}

func (_this *CTEEncodeBuffer) WriteBool(value bool) {
	if value {
		_this.WriteTrue()
	} else {
		_this.WriteFalse()
	}
}

func (_this *CTEEncodeBuffer) WriteInt(value int64) {
	if value >= 0 {
		_this.WritePositiveInt(uint64(value))
	} else {
		_this.WriteNegativeInt(uint64(-value))
	}
}

func (_this *CTEEncodeBuffer) WritePositiveInt(value uint64) {
	buff := _this.Allocate(uintStringMaxByteCount)[:0]
	used := strconv.AppendUint(buff, value, 10)
	_this.CorrectAllocation(len(used))
}

func (_this *CTEEncodeBuffer) WriteNegativeInt(value uint64) {
	_this.AddByte('-')
	_this.WritePositiveInt(value)
}

func (_this *CTEEncodeBuffer) WriteBigInt(value *big.Int) {
	if value == nil {
		_this.WriteNA()
		return
	}

	var buff [64]byte
	used := value.Append(buff[:0], 10)
	_this.AddBytes(used)
}

func (_this *CTEEncodeBuffer) WriteFloat(value float64) {
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

	buff := _this.Allocate(floatStringMaxByteCount)[:0]
	used := strconv.AppendFloat(buff, value, 'g', -1, 64)
	_this.CorrectAllocation(len(used))
}

func (_this *CTEEncodeBuffer) WriteFloatHexNoPrefix(value float64) {
	buff := _this.Allocate(floatStringMaxByteCount)[:0]
	used := strconv.AppendFloat(buff, value, 'x', -1, 64)
	copy(used, used[2:])
	if value < 0 {
		used[0] = '-'
	}
	_this.CorrectAllocation(len(used) - 2)
}

func (_this *CTEEncodeBuffer) WriteBigFloat(value *big.Float) {
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
}

func (_this *CTEEncodeBuffer) WriteDecimalFloat(value compact_float.DFloat) {
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
}

func (_this *CTEEncodeBuffer) WriteBigDecimalFloat(value *apd.Decimal) {
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
	}
}

func (_this *CTEEncodeBuffer) WriteUUID(v []byte) {
	if len(v) != 16 {
		panic(fmt.Errorf("expected UUID length 16 but got %v", len(v)))
	}
	_this.AddFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
}

func (_this *CTEEncodeBuffer) WriteTime(value time.Time) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		_this.unexpectedError(err, value)
	}
	_this.WriteCompactTime(t)
}

func (_this *CTEEncodeBuffer) WriteCompactTime(value *compact_time.Time) {
	if value == nil {
		_this.WriteNA()
		return
	}
	tz := func(v *compact_time.Time) string {
		switch v.TimezoneType {
		case compact_time.TypeZero:
			return ""
		case compact_time.TypeAreaLocation, compact_time.TypeLocal:
			return fmt.Sprintf("/%s", v.AreaLocation)
		case compact_time.TypeLatitudeLongitude:
			return fmt.Sprintf("/%.2f/%.2f", float64(v.LatitudeHundredths)/100, float64(v.LongitudeHundredths)/100)
		default:
			panic(fmt.Errorf("unknown compact time timezone type %v", value.TimezoneType))
			return ""
		}
	}
	subsec := func(v *compact_time.Time) string {
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
}

func (_this *CTEEncodeBuffer) WriteMarkerBegin(id interface{}) {
	_this.AddFmt("&%v:", id)
}

func (_this *CTEEncodeBuffer) WriteReference(id interface{}) {
	_this.AddFmt("$%v", id)
}

func (_this *CTEEncodeBuffer) WriteListBegin() {
	_this.AddByte('[')
}

func (_this *CTEEncodeBuffer) WriteListEnd() {
	_this.AddByte(']')
}

func (_this *CTEEncodeBuffer) WriteMapBegin() {
	_this.AddByte('{')
}

func (_this *CTEEncodeBuffer) WriteMapEnd() {
	_this.AddByte('}')
}

func (_this *CTEEncodeBuffer) WriteMetaBegin() {
	_this.AddByte('(')
}

func (_this *CTEEncodeBuffer) WriteMetaEnd() {
	_this.AddByte(')')
}

func (_this *CTEEncodeBuffer) WriteMarkupBegin() {
	_this.AddByte('<')
}

func (_this *CTEEncodeBuffer) WriteMarkupContentsBegin() {
	_this.AddByte(',')
}

func (_this *CTEEncodeBuffer) WriteMarkupEnd() {
	_this.AddByte('>')
}

func (_this *CTEEncodeBuffer) WriteArrayBegin() {
	_this.AddByte('|')
}

func (_this *CTEEncodeBuffer) WriteArrayEnd() {
	_this.AddByte('|')
}

func (_this *CTEEncodeBuffer) WriteCommentBegin() {
	_this.AddBytes([]byte{'/', '*'})
}

func (_this *CTEEncodeBuffer) WriteCommentEnd() {
	_this.AddBytes([]byte{'*', '/'})
}

func (_this *CTEEncodeBuffer) unexpectedError(err error, encoding interface{}) {
	panic(fmt.Errorf("unexpected error [%v] while encoding %v", err, encoding))
}
