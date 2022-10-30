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

// Test helper code.
package test

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/ce/events"
)

// ----------------------------------------------------------------------------
// Stored events
// ----------------------------------------------------------------------------

type EventInvocation func(receiver events.DataEventReceiver)

type Event interface {
	Name() string
	String() string
	ArrayElementCount() int
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

type BaseEvent struct {
	shortName         string
	invocation        EventInvocation
	values            []interface{}
	comparable        string
	stringified       string
	arrayElementCount int
}

func (_this *BaseEvent) Invoke(receiver events.DataEventReceiver) { _this.invocation(receiver) }
func (_this *BaseEvent) Name() string                             { return _this.shortName }
func (_this *BaseEvent) String() string                           { return _this.stringified }
func (_this *BaseEvent) Comparable() string                       { return _this.comparable }
func (_this *BaseEvent) ArrayElementCount() int                   { return _this.arrayElementCount }
func (_this *BaseEvent) IsEquivalentTo(that Event) bool {
	return _this.Comparable() == that.Comparable()
}

func ConstructEvent(shortName string, invocation EventInvocation, values ...interface{}) BaseEvent {
	return BaseEvent{
		shortName:         shortName,
		invocation:        invocation,
		values:            values,
		arrayElementCount: getArrayElementCount(values...),
		comparable:        constructComparable(shortName, values...),
		stringified:       constructStringified(shortName, values...),
	}
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

func constructComparable(shortName string, values ...interface{}) string {
	switch len(values) {
	case 0:
		return shortName
	case 1:
		return fmt.Sprintf("%v=%v", shortName, paramToComparable(values[0]))
	case 2:
		return fmt.Sprintf("%v=%v %v", shortName, paramToComparable(values[0]), paramToComparable(values[1]))
	default:
		panic(fmt.Errorf("expected 0, 1, or 2 values but got %v", len(values)))
	}
}

func constructStringified(shortName string, values ...interface{}) string {
	switch len(values) {
	case 0:
		return shortName
	case 1:
		return fmt.Sprintf("%v=%v", shortName, stringifyParam(values[0]))
	case 2:
		return fmt.Sprintf("%v=%v %v", shortName, stringifyParam(values[0]), stringifyParam(values[1]))
	default:
		panic(fmt.Errorf("expected 0, 1, or 2 values but got %v", len(values)))
	}
}

func paramToComparable(value interface{}) string {
	switch v := value.(type) {
	case float32:
		asInt := int64(v)
		if float32(asInt) == v {
			return strconv.FormatInt(asInt, 10)
		} else {
			return strconv.FormatFloat(float64(v), 'x', -1, 32)
		}
	case float64:
		asInt := int64(v)
		if float64(asInt) == v {
			return strconv.FormatInt(asInt, 10)
		} else {
			return strconv.FormatFloat(float64(v), 'x', -1, 64)
		}
	case *big.Float:
		if v.IsInt() {
			bi := &big.Int{}
			v.Int(bi)
			return bi.String()
		} else {
			return v.Text('x', -1)
		}
	case *apd.Decimal:
		return v.Text('g')
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (_this *EventUID) String() string {
	if len(_this.values) == 0 {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUID(_this.values[0].([]byte)))
}
func (_this *EventUID) IsEquivalentTo(that Event) bool {
	return _this.Comparable() == that.Comparable()
}

func (_this *EventArrayUID) String() string {
	if len(_this.values) == 0 {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUIDArray(_this.values[0].([][]byte)))
}

func (_this *EventArrayDataUID) String() string {
	if len(_this.values) == 0 {
		return _this.shortName
	}
	return fmt.Sprintf("%v=%v", _this.shortName, stringifyUIDArray(_this.values[0].([][]byte)))
}

func (_this *EventArrayDataText) IsEquivalentTo(that Event) bool {
	switch v := that.(type) {
	case *EventArrayDataText:
		return bytes.Equal([]byte(_this.values[0].(string)), []byte(v.values[0].(string)))
	case *EventArrayDataUint8:
		return bytes.Equal([]byte(_this.values[0].(string)), v.values[0].([]byte))
	default:
		return false
	}
}

func (_this *EventArrayDataUint8) IsEquivalentTo(that Event) bool {
	switch v := that.(type) {
	case *EventArrayDataText:
		return bytes.Equal(_this.values[0].([]byte), []byte(v.values[0].(string)))
	case *EventArrayDataUint8:
		return bytes.Equal(_this.values[0].([]byte), v.values[0].([]byte))
	default:
		return false
	}
}

func (_this *EventCustomBinary) String() string {
	switch len(_this.values) {
	case 1:
		return fmt.Sprintf("%v=%v", _this.shortName, _this.values[0])
	case 2:
		sb := strings.Builder{}
		for i, b := range _this.values[1].([]byte) {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%02x", b))
		}

		return fmt.Sprintf("%v=%v %v", _this.shortName, _this.values[0], sb.String())
	default:
		panic(fmt.Errorf("expected 1 or 2 values but got %v", len(_this.values)))
	}
}

func (_this *EventMedia) String() string {
	switch len(_this.values) {
	case 1:
		return fmt.Sprintf("%v=%v", _this.shortName, _this.values[0])
	case 2:
		sb := strings.Builder{}
		for i, b := range _this.values[1].([]byte) {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%02x", b))
		}

		return fmt.Sprintf("%v=%v %v", _this.shortName, _this.values[0], sb.String())
	default:
		panic(fmt.Errorf("expected 1 or 2 values but got %v", len(_this.values)))
	}
}

type NegFFFFFFFFFFFFFFFF struct{}

func (_this NegFFFFFFFFFFFFFFFF) String() string { return "-18446744073709551615" }

type EventNumeric struct{ BaseEvent }

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
			if hasQuietBitSet64(v) {
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
		} else {
			value = v.Copy(v)
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

		default:
			value = &apd.Decimal{
				Form:     v.Form,
				Negative: v.Negative,
				Exponent: v.Exponent,
				Coeff:    v.Coeff,
			}
		}
	case *big.Int:
		vCopy := big.NewInt(0).Set(v)
		if u := vCopy.Uint64(); u == 0xffffffffffffffff && vCopy.Sign() < 0 {
			vCopy.Neg(vCopy)
			if vCopy.IsUint64() {
				// Workaround for big.Int bug where the sign gets reset
				value = NegFFFFFFFFFFFFFFFF{}
				break
			}
			vCopy.Set(v)
		}
		value = vCopy
	}

	return &EventNumeric{
		BaseEvent: ConstructEvent("n", func(receiver events.DataEventReceiver) {
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
			case NegFFFFFFFFFFFFFFFF:
				receiver.OnNegativeInt(0xffffffffffffffff)
			default:
				panic(fmt.Errorf("unexpected numeric type %v for value %v", reflect.TypeOf(value), value))
			}
		}, value),
	}
}

func (_this *EventNumeric) String() string {
	if v, ok := _this.values[0].(compact_float.DFloat); ok {
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
	return _this.BaseEvent.String()
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

func hasQuietBitSet64(value float64) bool {
	const quietBit = uint64(0x0008000000000000)
	return math.Float64bits(value)&quietBit != 0
}

func hasQuietBitSet32(value float32) bool {
	const quietBit = uint32(0x00400000)
	return math.Float32bits(value)&quietBit != 0
}
