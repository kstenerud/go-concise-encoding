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
	Name() string
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
		case *apd.Decimal:
			comparable = fmt.Sprintf("%v=%v", shortName, v.Text('g'))
		default:
			comparable = fmt.Sprintf("%v=%v", shortName, value)
		}
		stringified = fmt.Sprintf("%v=%v", shortName, stringifyParam(value))
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
func (_this *EventWithValue) Name() string                             { return _this.shortName }
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

func (_this *EventCustomBinary) String() string {
	if _this.value == NoValue {
		return _this.shortName
	}
	sb := strings.Builder{}
	for i, b := range _this.value.([]byte) {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%02x", b))
	}

	return fmt.Sprintf("%v=%v", _this.shortName, sb.String())
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
		value = big.NewInt(0).Set(v)
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

func hasQuietBitSet64(value float64) bool {
	const quietBit = uint64(0x0008000000000000)
	return math.Float64bits(value)&quietBit != 0
}

func hasQuietBitSet32(value float32) bool {
	const quietBit = uint32(0x00400000)
	return math.Float32bits(value)&quietBit != 0
}
