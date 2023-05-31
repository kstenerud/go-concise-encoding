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
	"bytes"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/ce/events"
)

func (_this *EventArrayBit) Expand() Events {
	begin := BAB()
	if len(_this.values) == 0 {
		return Events{begin, ACL(0)}
	}
	elements := _this.values[0].([]bool)
	if len(elements) == 0 {
		return Events{begin, ACL(0)}
	}
	return Events{begin, ACL(uint64(len(elements))), ADB(elements)}
}

func (_this *EventCustomBinary) Expand() Events {
	customType := _this.values[0].(uint64)
	begin := BCB(customType)
	if len(_this.values) < 2 {
		return Events{begin, ACL(0)}
	}
	elements := _this.values[1].([]byte)
	if len(elements) == 0 {
		return Events{begin, ACL(0)}
	}
	return Events{begin, ACL(uint64(len(elements))), ADU8(elements)}
}

func (_this *EventCustomText) Expand() Events {
	customType := _this.values[0].(uint64)
	begin := BCT(customType)
	if len(_this.values) < 2 {
		return Events{begin, ACL(0)}
	}
	elements := _this.values[1].(string)
	if len(elements) == 0 {
		return Events{begin, ACL(0)}
	}
	return Events{begin, ACL(uint64(len(elements))), ADT(elements)}
}

func (_this *EventMedia) Expand() Events {
	mediaType := _this.values[0].(string)
	begin := BMEDIA(mediaType)
	if len(_this.values) < 2 {
		return Events{begin, ACL(0)}
	}
	elements := _this.values[1].([]byte)
	if len(elements) == 0 {
		return Events{begin, ACL(0)}
	}
	return Events{begin, ACL(uint64(len(elements))), ADU8(elements)}
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

func (_this *EventNumeric) Expand() Events { return Events{_this} }

func NAN() Event  { return N(compact_float.QuietNaN()) }
func SNAN() Event { return N(compact_float.SignalingNaN()) }
