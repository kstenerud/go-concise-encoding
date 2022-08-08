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

// Test helper code.
package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/events"
)

// ----------------------------------------------------------------------------
// Stored events
// ----------------------------------------------------------------------------

type EventInvocation func(receiver events.DataEventReceiver)

type Event interface {
	String() string
	Invoke(events.DataEventReceiver)
	IsEquivalentTo(Event) bool
	Comparable() string
}

type Events []Event

func (_this Events) String() string {
	sb := bytes.Buffer{}
	for i, v := range _this {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('"')
		sb.WriteString(v.String())
		sb.WriteByte('"')
	}
	return sb.String()
}

func (_this Events) AreEquivalentTo(that Events) bool {
	return AreEventsEquivalent(_this, that)
}

var NoValue interface{}

type EventWithValue struct {
	shortName   string
	invocation  EventInvocation
	value       interface{}
	comparable  string
	stringified string
}

func ConstructEventWithValue(shortName string, value interface{}, invocation EventInvocation) EventWithValue {
	comparable := shortName
	stringified := shortName
	if value != NoValue {
		switch v := value.(type) {
		case float32:
			asInt := int(v)
			if float32(asInt) == v {
				comparable = fmt.Sprintf("%v=%v", shortName, asInt)
			} else {
				comparable = fmt.Sprintf("%v=%x", shortName, v)
			}
		case float64:
			asInt := int(v)
			if float64(asInt) == v {
				comparable = fmt.Sprintf("%v=%v", shortName, asInt)
			} else {
				comparable = fmt.Sprintf("%v=%x", shortName, v)
			}
		case *big.Float:
			if v.IsInt() {
				bi := &big.Int{}
				v.Int(bi)
				comparable = fmt.Sprintf("%v=%v", shortName, bi)
			} else {
				comparable = fmt.Sprintf("%v=%v", shortName, v.Text('x', -1))
			}
		default:
			comparable = fmt.Sprintf("%v=%v", shortName, value)
		}
		// stringified = fmt.Sprintf("%v=%v", shortName, stringifyParam(value))
		stringified = comparable
	}

	return EventWithValue{
		shortName:   shortName,
		invocation:  invocation,
		value:       value,
		comparable:  comparable,
		stringified: stringified,
	}
}
func (_this *EventWithValue) Invoke(receiver events.DataEventReceiver) { _this.invocation(receiver) }
func (_this *EventWithValue) String() string                           { return _this.stringified }
func (_this *EventWithValue) Comparable() string                       { return _this.comparable }
func (_this *EventWithValue) IsEquivalentTo(that Event) bool {
	return _this.Comparable() == that.Comparable()
}

func (_this *EventUID) String() string {
	if _this.value == NoValue {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUID(_this.value.([]byte)))
}
func (_this *EventUID) IsEquivalentTo(that Event) bool {
	return _this.Comparable() == that.Comparable()
}

func (_this *EventArrayUID) String() string {
	if _this.value == NoValue {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUIDArray(_this.value.([][]byte)))
}

func (_this *EventArrayDataUID) String() string {
	if _this.value == NoValue {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUIDArray(_this.value.([][]byte)))
}

func (_this *EventArrayDataText) IsEquivalentTo(that Event) bool {
	switch v := that.(type) {
	case *EventArrayDataText:
		return bytes.Equal([]byte(_this.value.(string)), []byte(v.value.(string)))
	case *EventArrayDataUint8:
		return bytes.Equal([]byte(_this.value.(string)), v.value.([]byte))
	default:
		return false
	}
}

func (_this *EventArrayDataUint8) IsEquivalentTo(that Event) bool {
	switch v := that.(type) {
	case *EventArrayDataText:
		return bytes.Equal(_this.value.([]byte), []byte(v.value.(string)))
	case *EventArrayDataUint8:
		return bytes.Equal(_this.value.([]byte), v.value.([]byte))
	default:
		return false
	}
}

type EventNumeric struct{ EventWithValue }

func N(value interface{}) Event {
	if value == nil {
		return NULL()
	}
	switch v := value.(type) {
	case float32:
		if math.IsNaN(float64(v)) {
			if hasQuietBitSet32(v) {
				value = compact_float.QuietNaN()
			} else {
				value = compact_float.SignalingNaN()
			}
		} else if math.IsInf(float64(v), 1) {
			value = compact_float.Infinity()
		} else if math.IsInf(float64(v), -1) {
			value = compact_float.NegativeInfinity()
		} else if v == 0 && math.Float32bits(v) == 0x80000000 {
			value = compact_float.NegativeZero()
		}
	case float64:
		if math.IsNaN(v) {
			if HasQuietBitSet64(v) {
				value = compact_float.QuietNaN()
			} else {
				value = compact_float.SignalingNaN()
			}
		} else if math.IsInf(v, 1) {
			value = compact_float.Infinity()
		} else if math.IsInf(v, -1) {
			value = compact_float.NegativeInfinity()
		} else if v == 0 && math.Float64bits(v) == 0x8000000000000000 {
			value = compact_float.NegativeZero()
		}
	case *big.Float:
		if v.IsInf() {
			if v.Sign() > 0 {
				value = compact_float.Infinity()
			} else {
				value = compact_float.NegativeInfinity()
			}
		}
	case *apd.Decimal:
		switch v.Form {
		case apd.NaN:
			value = compact_float.QuietNaN()
		case apd.NaNSignaling:
			value = compact_float.SignalingNaN()
		case apd.Infinite:
			if v.Negative {
				value = compact_float.NegativeInfinity()
			} else {
				value = compact_float.Infinity()
			}
		}
	}

	return &EventNumeric{
		EventWithValue: ConstructEventWithValue("n", value, func(receiver events.DataEventReceiver) {
			switch v := value.(type) {
			case int:
				receiver.OnInt(int64(v))
			case int8:
				receiver.OnInt(int64(v))
			case int16:
				receiver.OnInt(int64(v))
			case int32:
				receiver.OnInt(int64(v))
			case int64:
				receiver.OnInt(int64(v))
			case uint:
				receiver.OnPositiveInt(uint64(v))
			case uint8:
				receiver.OnPositiveInt(uint64(v))
			case uint16:
				receiver.OnPositiveInt(uint64(v))
			case uint32:
				receiver.OnPositiveInt(uint64(v))
			case uint64:
				receiver.OnPositiveInt(uint64(v))
			case float32:
				receiver.OnFloat(float64(v))
			case float64:
				receiver.OnFloat(v)
			case *big.Int:
				receiver.OnBigInt(v)
			case *big.Float:
				receiver.OnBigFloat(v)
			case compact_float.DFloat:
				receiver.OnDecimalFloat(v)
			case *apd.Decimal:
				receiver.OnBigDecimalFloat(v)
			default:
				panic(fmt.Errorf("unexpected numeric type %v for value %v", reflect.TypeOf(value), value))
			}
		}),
	}
}

func (_this *EventNumeric) String() string {
	if v, ok := _this.value.(compact_float.DFloat); ok {
		switch v {
		case compact_float.Infinity():
			return fmt.Sprintf("%v=inf", _this.shortName)
		case compact_float.NegativeInfinity():
			return fmt.Sprintf("%v=-inf", _this.shortName)
		case compact_float.QuietNaN():
			return fmt.Sprintf("%v=nan", _this.shortName)
		case compact_float.SignalingNaN():
			return fmt.Sprintf("%v=snan", _this.shortName)
		}
	}
	return _this.EventWithValue.String()
}

func NAN() Event  { return N(compact_float.QuietNaN()) }
func SNAN() Event { return N(compact_float.SignalingNaN()) }

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

// func areBytesEquivalentTo(a []byte, b interface{}) bool {
// 	switch v := b.(type) {
// 	case string:
// 		return string(a) == v
// 	case []byte:
// 		return bytes.Equal(a, v)
// 	case []int8:
// 		return bytes.Equal(a, ArrayInt8ToBytes(v))
// 	case []int16:
// 		return bytes.Equal(a, ArrayInt16ToBytes(v))
// 	case []int32:
// 		return bytes.Equal(a, ArrayInt32ToBytes(v))
// 	case []int64:
// 		return bytes.Equal(a, ArrayInt64ToBytes(v))
// 	case []uint16:
// 		return bytes.Equal(a, ArrayUint16ToBytes(v))
// 	case []uint32:
// 		return bytes.Equal(a, ArrayUint32ToBytes(v))
// 	case []uint64:
// 		return bytes.Equal(a, ArrayUint64ToBytes(v))
// 	case []float32:
// 		return bytes.Equal(a, ArrayFloat32ToBytes(v))
// 	case []float64:
// 		return bytes.Equal(a, ArrayFloat64ToBytes(v))
// 	case []bool:
// 		return bytes.Equal(a, ArrayBitsToBytes(v))
// 	case [][]byte:
// 		return bytes.Equal(a, ArrayUIDToBytes(v))
// 	default:
// 		return false
// 	}
// }

// ----------------------------------------------------------------------------
// Conversions
// ----------------------------------------------------------------------------

func ArrayBitsToBytes(src []bool) []byte {
	dst := make([]byte, len(src)/8+1)
	if len(src)&3 == 0 {
		dst = dst[:len(dst)-1]
	}

	iDst := 0
	for iSrc := 0; iSrc < len(src); {
		bitCount := 8
		if len(src)-iSrc < 8 {
			bitCount = len(src) - iSrc
		}
		accum := byte(0)
		for iBit := 0; iBit < bitCount; iBit++ {
			if src[iSrc] {
				accum |= 1 << iBit
			}
			iSrc++
		}
		dst[iDst] = accum
		iDst++
	}

	return dst
}

func ArrayInt8ToBytes(src []int8) []byte {
	bytes := make([]byte, len(src))
	for i, v := range src {
		bytes[i] = byte(v)
	}
	return bytes
}

func ArrayInt16ToBytes(src []int16) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], uint16(v))
	}
	return bytes
}

func ArrayInt32ToBytes(src []int32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], uint32(v))
	}
	return bytes
}

func ArrayInt64ToBytes(src []int64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], uint64(v))
	}
	return bytes
}

func ArrayUint8ToBytes(src []uint8) []byte {
	return src
}

func ArrayUint16ToBytes(src []uint16) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], v)
	}
	return bytes
}

func ArrayUint32ToBytes(src []uint32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], v)
	}
	return bytes
}

func ArrayUint64ToBytes(src []uint64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], v)
	}
	return bytes
}

func ArrayFloat16ToBytes(src []float32) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], uint16(math.Float32bits(v)>>16))
	}
	return bytes
}

func ArrayFloat32ToBytes(src []float32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], math.Float32bits(v))
	}
	return bytes
}

func ArrayFloat64ToBytes(src []float64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], math.Float64bits(v))
	}
	return bytes
}

func ArrayUIDToBytes(src [][]byte) []byte {
	const step = 16
	bytes := make([]byte, 0, len(src)*step)
	for _, v := range src {
		bytes = append(bytes, v...)
	}
	return bytes
}

func ArrayTextToBytes(src string) []byte {
	return []byte(src)
}

func BytesToArrayBits(length uint64, src []byte) []bool {
	dst := make([]bool, length)
	if length == 0 {
		return dst
	}

	iDst := 0
	iSrc := 0
	for ; iSrc < len(src)-1; iSrc++ {
		b := src[iSrc]
		for mask := 1; mask < 256; mask <<= 1 {
			dst[iDst] = b&byte(mask) != 0
			iDst++
		}
	}

	b := src[iSrc]
	length -= uint64(iDst)
	for iBit := uint64(0); iBit < length; iBit++ {
		dst[iDst] = b&byte(1<<iBit) != 0
		iDst++
	}

	return dst
}

func BytesToArrayInt8(src []byte) []int8 {
	array := make([]int8, len(src))
	for i, v := range src {
		array[i] = int8(v)
	}
	return array
}

func BytesToArrayInt16(src []byte) []int16 {
	const step = 2
	elemCount := len(src) / step
	array := make([]int16, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int16(binary.LittleEndian.Uint16(src[i*step:]))
	}
	return array
}

func BytesToArrayInt32(src []byte) []int32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]int32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int32(binary.LittleEndian.Uint32(src[i*step:]))
	}
	return array
}

func BytesToArrayInt64(src []byte) []int64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]int64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int64(binary.LittleEndian.Uint64(src[i*step:]))
	}
	return array
}

func BytesToArrayUint8(src []byte) []uint8 {
	return src
}

func BytesToArrayUint16(src []byte) []uint16 {
	const step = 2
	elemCount := len(src) / step
	array := make([]uint16, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint16(src[i*step:])
	}
	return array
}

func BytesToArrayUint32(src []byte) []uint32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]uint32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint32(src[i*step:])
	}
	return array
}

func BytesToArrayUint64(src []byte) []uint64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]uint64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint64(src[i*step:])
	}
	return array
}

func BytesToArrayFloat16(src []byte) []float32 {
	const step = 2
	elemCount := len(src) / step
	array := make([]float32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float32frombits(uint32(binary.LittleEndian.Uint16(src[i*step:])) << 16)
	}
	return array
}

func BytesToArrayFloat32(src []byte) []float32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]float32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float32frombits(binary.LittleEndian.Uint32(src[i*step:]))
	}
	return array
}

func BytesToArrayFloat64(src []byte) []float64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]float64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float64frombits(binary.LittleEndian.Uint64(src[i*step:]))
	}
	return array
}

func BytesToArrayUID(src []byte) [][]byte {
	const step = 16
	elemCount := len(src) / step
	array := make([][]byte, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = src[i*step : (i+1)*step]
	}
	return array
}

func BytesToArrayText(src []byte) string {
	return string(src)
}

// ----------------------------------------------------------------------------
// String conversion
// ----------------------------------------------------------------------------

func hexChar(v byte) byte {
	if v < 10 {
		return '0' + v
	}
	return 'a' + v - 10
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case []byte:
		var builder strings.Builder
		builder.WriteByte('[')
		for i, b := range v {
			builder.WriteByte(hexChar(b >> 4))
			builder.WriteByte(hexChar(b & 15))
			if i < len(v) {
				builder.WriteByte(' ')
			}
		}
		builder.WriteByte(']')
		return builder.String()
	case *big.Float:
		return v.Text('x', -1)
	case float64:
		if math.IsNaN(v) {
			if HasQuietBitSet64(v) {
				return "nan"
			} else {
				return "snan"
			}
		}
		return strconv.FormatFloat(v, 'x', -1, 64)
	case float32:
		if math.IsNaN(float64(v)) {
			if hasQuietBitSet32(v) {
				return "nan"
			} else {
				return "snan"
			}
		}
		return strconv.FormatFloat(float64(v), 'x', -1, 64)
	case compact_float.DFloat:
		if v.IsNan() {
			if v.IsSignalingNan() {
				return "snan"
			} else {
				return "nan"
			}
		} else {
			return v.String()
		}
	case *apd.Decimal:
		switch v.Form {
		case apd.NaNSignaling:
			return "snan"
		case apd.NaN:
			return "nan"
		default:
			return v.String()
		}
	case uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%v", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func stringifyBitArray(array []bool) string {
	sb := strings.Builder{}
	for _, v := range array {
		if v {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
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

func stringifyParam(param interface{}) string {
	switch v := param.(type) {
	case []bool:
		return stringifyBitArray(v)
	}

	switch reflect.TypeOf(param).Kind() {
	case reflect.Slice:
		sb := strings.Builder{}
		v := reflect.ValueOf(param)
		if reflect.TypeOf(param).Elem().Kind() == reflect.Uint8 {
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
				sb.WriteString(stringify(v.Index(i)))
			}
		}
		return sb.String()
	default:
		return stringify(param)
	}
}

func HasQuietBitSet64(value float64) bool {
	const quietBit = uint64(0x0008000000000000)
	return math.Float64bits(value)&quietBit != 0
}

func hasQuietBitSet32(value float32) bool {
	const quietBit = uint32(0x00400000)
	return math.Float32bits(value)&quietBit != 0
}
