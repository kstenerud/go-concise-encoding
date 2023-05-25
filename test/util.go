// Copyright 2022 Karl Stenerud
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

package test

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
)

func toComparable(value interface{}) string {
	return toStringified(value)
}

func toStringified(value interface{}) string {
	if reflectV, ok := value.(reflect.Value); ok {
		return toStringified(reflectV.Interface())
	}
	switch v := value.(type) {
	case []bool:
		return stringifyArrayBoolean(v)
	case *big.Float:
		return v.Text('x', -1)
	case float64:
		return stringifyFloat64(v)
	case float32:
		return stringifyFloat32(v)
	case compact_float.DFloat:
		return stringifyDFloat(v)
	case *apd.Decimal:
		return stringifyAPD(v)
	case uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%v", value)
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		sb := strings.Builder{}
		v := reflect.ValueOf(value)
		if reflect.TypeOf(value).Elem().Kind() == reflect.Uint8 {
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("0x%02x", v.Index(i)))
			}
		} else {
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(toStringified(v.Index(i)))
			}
		}
		return sb.String()
	}

	return fmt.Sprintf("%v", value)
}

func stringifyArrayBoolean(v []bool) string {
	sb := strings.Builder{}
	for _, v := range v {
		if v {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func stringifyAPD(v *apd.Decimal) string {
	switch v.Form {
	case apd.NaNSignaling:
		return "snan"
	case apd.NaN:
		return "nan"
	default:
		return v.Text('g')
	}
}

func stringifyDFloat(v compact_float.DFloat) string {
	if v.IsNan() {
		if v.IsSignalingNan() {
			return "snan"
		} else {
			return "nan"
		}
	} else {
		return v.String()
	}
}

func stringifyFloat32(v float32) string {
	if math.IsNaN(float64(v)) {
		if hasQuietBitSet32(v) {
			return "nan"
		} else {
			return "snan"
		}
	}
	if math.IsInf(float64(v), 1) {
		return "inf"
	}
	if math.IsInf(float64(v), -1) {
		return "-inf"
	}
	if v == 0 {
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	}
	if asInt := int64(v); float32(asInt) == v {
		return strconv.FormatInt(asInt, 10)
	}
	return strconv.FormatFloat(float64(v), 'x', -1, 32)
}

func stringifyFloat64(v float64) string {
	if math.IsNaN(v) {
		if hasQuietBitSet64(v) {
			return "nan"
		} else {
			return "snan"
		}
	}
	if math.IsInf(v, 1) {
		return "inf"
	}
	if math.IsInf(v, -1) {
		return "-inf"
	}
	if v == 0 {
		return strconv.FormatFloat(float64(v), 'g', -1, 64)
	}
	if asInt := int64(v); float64(asInt) == v {
		return strconv.FormatInt(asInt, 10)
	}
	return strconv.FormatFloat(v, 'x', -1, 64)
}

func stringifyUIDArray(uids [][]byte) string {
	sb := strings.Builder{}
	for i, v := range uids {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(stringifyUID(v))
	}
	return sb.String()
}

func stringifyUID(uid []byte) string {
	if len(uid) != 16 {
		panic(fmt.Errorf("BUG: Not a UID: %v", uid))
	}
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		uid[0], uid[1], uid[2], uid[3], uid[4], uid[5], uid[6], uid[7], uid[8], uid[9], uid[10], uid[11], uid[12], uid[13], uid[14], uid[15])
}

func getArrayElementCount(values ...interface{}) int {
	for _, v := range values {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.String, reflect.Slice, reflect.Array:
			return rv.Len()
		}
	}
	return 0
}

func copyOf(value interface{}) interface{} {
	switch v := value.(type) {
	case []bool:
		return append([]bool{}, v...)
	case []int8:
		return append([]int8{}, v...)
	case []int16:
		return append([]int16{}, v...)
	case []int32:
		return append([]int32{}, v...)
	case []int64:
		return append([]int64{}, v...)
	case []float32:
		return append([]float32{}, v...)
	case []float64:
		return append([]float64{}, v...)
	case []uint8:
		return append([]uint8{}, v...)
	case []uint16:
		return append([]uint16{}, v...)
	case []uint32:
		return append([]uint32{}, v...)
	case []uint64:
		return append([]uint64{}, v...)
	case [][]uint8:
		var cp [][]uint8
		for _, elem := range v {
			cp = append(cp, append([]uint8{}, elem...))
		}
		return cp
	default:
		return value
	}
}

func hasQuietBitSet64(value float64) bool {
	const quietBit = uint64(0x0008000000000000)
	return math.Float64bits(value)&quietBit != 0
}

func hasQuietBitSet32(value float32) bool {
	const quietBit = uint32(0x00400000)
	return math.Float32bits(value)&quietBit != 0
}
