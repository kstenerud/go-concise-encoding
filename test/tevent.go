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
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/arrays"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-equivalence"
)

// ----------------------------------------------------------------------------
// Event types and pretty-print names
// ----------------------------------------------------------------------------

type TEventType int

const (
	TEventBeginDocument TEventType = iota
	TEventVersion
	TEventPadding
	TEventComment
	TEventNull
	TEventBool
	TEventTrue
	TEventFalse
	TEventPInt
	TEventNInt
	TEventInt
	TEventBigInt
	TEventFloat
	TEventBigFloat
	TEventDecimalFloat
	TEventBigDecimalFloat
	TEventQNan
	TEventSNan
	TEventUID
	TEventTime
	TEventCompactTime
	TEventString
	TEventResourceID
	TEventRemoteRef
	TEventCustomBinary
	TEventCustomText
	TEventArrayBoolean
	TEventArrayInt8
	TEventArrayInt16
	TEventArrayInt32
	TEventArrayInt64
	TEventArrayUint8
	TEventArrayUint16
	TEventArrayUint32
	TEventArrayUint64
	TEventArrayFloat16
	TEventArrayFloat32
	TEventArrayFloat64
	TEventArrayUID
	TEventStringBegin
	TEventResourceIDBegin
	TEventRemoteRefBegin
	TEventCustomBinaryBegin
	TEventCustomTextBegin
	TEventArrayBooleanBegin
	TEventArrayInt8Begin
	TEventArrayInt16Begin
	TEventArrayInt32Begin
	TEventArrayInt64Begin
	TEventArrayUint8Begin
	TEventArrayUint16Begin
	TEventArrayUint32Begin
	TEventArrayUint64Begin
	TEventArrayFloat16Begin
	TEventArrayFloat32Begin
	TEventArrayFloat64Begin
	TEventArrayUIDBegin
	TEventMediaBegin
	TEventArrayChunk
	TEventArrayData
	TEventList
	TEventMap
	TEventMarkup
	TEventEnd
	TEventNode
	TEventEdge
	TEventMarker
	TEventReference
	TEventConstant
	TEventEndDocument
)

var TEventNames = []string{
	TEventBeginDocument:     "BD",
	TEventVersion:           "V",
	TEventPadding:           "PAD",
	TEventComment:           "COM",
	TEventNull:              "N",
	TEventBool:              "B",
	TEventTrue:              "B(true)",
	TEventFalse:             "B(false)",
	TEventPInt:              "PI",
	TEventNInt:              "NI",
	TEventInt:               "I",
	TEventBigInt:            "BI",
	TEventFloat:             "BF",
	TEventBigFloat:          "BBF",
	TEventDecimalFloat:      "DF",
	TEventBigDecimalFloat:   "BDF",
	TEventQNan:              "QNAN",
	TEventSNan:              "SNAN",
	TEventUID:               "UID",
	TEventTime:              "GT",
	TEventCompactTime:       "T",
	TEventString:            "S",
	TEventResourceID:        "RID",
	TEventRemoteRef:         "RREF",
	TEventCustomBinary:      "CB",
	TEventCustomText:        "CT",
	TEventArrayBoolean:      "AB",
	TEventArrayInt8:         "AI8",
	TEventArrayInt16:        "AI16",
	TEventArrayInt32:        "AI32",
	TEventArrayInt64:        "AI64",
	TEventArrayUint8:        "AU8",
	TEventArrayUint16:       "AU16",
	TEventArrayUint32:       "AU32",
	TEventArrayUint64:       "AU64",
	TEventArrayFloat16:      "AF16",
	TEventArrayFloat32:      "AF32",
	TEventArrayFloat64:      "AF64",
	TEventArrayUID:          "AU",
	TEventStringBegin:       "SB",
	TEventResourceIDBegin:   "RB",
	TEventRemoteRefBegin:    "RRB",
	TEventCustomBinaryBegin: "CBB",
	TEventCustomTextBegin:   "CTB",
	TEventArrayBooleanBegin: "ABB",
	TEventArrayInt8Begin:    "AI8B",
	TEventArrayInt16Begin:   "AI16B",
	TEventArrayInt32Begin:   "AI32B",
	TEventArrayInt64Begin:   "AI64B",
	TEventArrayUint8Begin:   "AU8B",
	TEventArrayUint16Begin:  "AU16B",
	TEventArrayUint32Begin:  "AU32B",
	TEventArrayUint64Begin:  "AU64B",
	TEventArrayFloat16Begin: "AF16B",
	TEventArrayFloat32Begin: "AF32B",
	TEventArrayFloat64Begin: "AF64B",
	TEventArrayUIDBegin:     "AUB",
	TEventMediaBegin:        "MB",
	TEventArrayChunk:        "AC",
	TEventArrayData:         "AD",
	TEventList:              "L",
	TEventMap:               "M",
	TEventMarkup:            "MU",
	TEventNode:              "NODE",
	TEventEdge:              "EDGE",
	TEventEnd:               "E",
	TEventMarker:            "MARK",
	TEventReference:         "REF",
	TEventConstant:          "CONST",
	TEventEndDocument:       "ED",
}

func (_this TEventType) String() string {
	return TEventNames[_this]
}

func (_this TEventType) IsBoolean() bool {
	switch _this {
	case TEventTrue, TEventFalse, TEventBool:
		return true
	default:
		return false
	}
}

func (_this TEventType) IsNumeric() bool {
	switch _this {
	case TEventPInt, TEventNInt, TEventInt, TEventBigInt, TEventFloat,
		TEventBigFloat, TEventDecimalFloat, TEventBigDecimalFloat, TEventQNan,
		TEventSNan:
		return true
	default:
		return false
	}
}

// ----------------------------------------------------------------------------
// Stored events
// ----------------------------------------------------------------------------

func AreAllEventsEquivalent(a []*TEvent, b []*TEvent) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ev := range a {
		if !ev.IsEquivalentTo(b[i]) {
			return false
		}
	}
	return true
}

type TEvent struct {
	Type TEventType
	V1   interface{}
	V2   interface{}
}

func NewTEvent(eventType TEventType, v1 interface{}, v2 interface{}) *TEvent {
	return &TEvent{
		Type: eventType,
		V1:   v1,
		V2:   v2,
	}
}

func hexChar(v byte) byte {
	if v < 10 {
		return '0' + v
	}
	return 'a' + v - 10
}

func (_this *TEvent) stringify(value interface{}) string {
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
	case string:
		return fmt.Sprintf("\"%v\"", value)
	case *big.Float:
		return v.Text('x', -1)
	case float64:
		if math.IsNaN(v) {
			if common.HasQuietNanBitSet64(v) {
				return "qnan"
			} else {
				return "snan"
			}
		}
		return strconv.FormatFloat(v, 'x', -1, 64)
	case float32:
		if math.IsNaN(float64(v)) {
			if common.HasQuietNanBitSet32(v) {
				return "qnan"
			} else {
				return "snan"
			}
		}
		return strconv.FormatFloat(float64(v), 'x', -1, 64)
	case uint, uint64, uint32, uint16, uint8:
		if _this.Type == TEventNInt {
			return fmt.Sprintf("-%v", _this.V1)
		} else {
			return fmt.Sprintf("%v", _this.V1)
		}
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (_this *TEvent) String() string {
	if _this.V1 != nil {
		if _this.V2 != nil {
			return fmt.Sprintf("%v(%v,%v)", _this.Type.String(), _this.stringify(_this.V1), _this.stringify(_this.V2))
		}
		return fmt.Sprintf("%v(%v)", _this.Type.String(), _this.stringify(_this.V1))
	}
	return _this.Type.String()
}

func (_this *TEvent) IsTrue() bool {
	switch _this.Type {
	case TEventTrue:
		return true
	case TEventBool:
		return _this.V1.(bool)
	default:
		return false
	}
}

func (_this *TEvent) IsQuietNan() bool {
	switch _this.Type {
	case TEventQNan:
		return true
	case TEventFloat:
		f64 := _this.V1.(float64)
		return math.IsNaN(f64) && common.HasQuietNanBitSet64(f64)
	case TEventDecimalFloat:
		cf := _this.V1.(compact_float.DFloat)
		return cf.IsNan() && !cf.IsSignalingNan()
	case TEventBigDecimalFloat:
		return _this.V1.(*apd.Decimal).Form == apd.NaN
	default:
		return false
	}
}

func (_this *TEvent) IsSignalingNan() bool {
	switch _this.Type {
	case TEventSNan:
		return true
	case TEventFloat:
		f64 := _this.V1.(float64)
		return math.IsNaN(f64) && !common.HasQuietNanBitSet64(f64)
	case TEventDecimalFloat:
		return _this.V1.(compact_float.DFloat).IsSignalingNan()
	case TEventBigDecimalFloat:
		return _this.V1.(*apd.Decimal).Form == apd.NaNSignaling
	default:
		return false
	}
}

func (_this *TEvent) IsEffectivelyNull() bool {
	return _this.Type == TEventNull ||
		_this == EvBINull ||
		_this == EvBBFNull ||
		_this == EvBDFNull
}

func (_this *TEvent) IsEquivalentTo(that *TEvent) bool {
	if equivalence.IsEquivalent(_this, that) {
		return true
	}

	if _this.IsEffectivelyNull() && that.IsEffectivelyNull() {
		return true
	}

	if _this.Type.IsBoolean() && that.Type.IsBoolean() {
		return _this.IsTrue() == that.IsTrue()
	}

	if _this.Type.IsNumeric() && that.Type.IsNumeric() {
		if _this.IsQuietNan() && that.IsQuietNan() {
			return true
		}
		if _this.IsSignalingNan() && that.IsSignalingNan() {
			return true
		}
		switch _this.Type {
		case TEventPInt:
			switch that.Type {
			case TEventPInt:
				return _this.V1 == that.V1
			case TEventNInt:
				return false
			case TEventInt:
				return pintEqualsInt(_this.V1.(uint64), that.V1.(int64))
			case TEventBigInt:
				return pintEqualsBigInt(_this.V1.(uint64), that.V1.(*big.Int))
			case TEventFloat:
				return pintEqualsFloat(_this.V1.(uint64), that.V1.(float64))
			case TEventBigFloat:
				return pintEqualsBigFloat(_this.V1.(uint64), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return pintEqualsDFloat(_this.V1.(uint64), that.V1.(compact_float.DFloat))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventNInt:
			switch that.Type {
			case TEventPInt:
				return false
			case TEventNInt:
				return _this.V1 == that.V1
			case TEventInt:
				return nintEqualsInt(_this.V1.(uint64), that.V1.(int64))
			case TEventBigInt:
				return nintEqualsBigInt(_this.V1.(uint64), that.V1.(*big.Int))
			case TEventFloat:
				return nintEqualsFloat(_this.V1.(uint64), that.V1.(float64))
			case TEventBigFloat:
				return nintEqualsBigFloat(_this.V1.(uint64), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return nintEqualsDFloat(_this.V1.(uint64), that.V1.(compact_float.DFloat))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventInt:
			switch that.Type {
			case TEventPInt:
				return pintEqualsInt(that.V1.(uint64), _this.V1.(int64))
			case TEventNInt:
				return nintEqualsInt(that.V1.(uint64), _this.V1.(int64))
			case TEventInt:
				return _this.V1 == that.V1
			case TEventBigInt:
				return intEqualsBigInt(_this.V1.(int64), that.V1.(*big.Int))
			case TEventFloat:
				return intEqualsFloat(_this.V1.(int64), that.V1.(float64))
			case TEventBigFloat:
				return intEqualsBigFloat(_this.V1.(int64), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return intEqualsDFloat(_this.V1.(int64), that.V1.(compact_float.DFloat))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventBigInt:
			switch that.Type {
			case TEventPInt:
				return pintEqualsBigInt(that.V1.(uint64), _this.V1.(*big.Int))
			case TEventNInt:
				return nintEqualsBigInt(that.V1.(uint64), _this.V1.(*big.Int))
			case TEventInt:
				return intEqualsBigInt(that.V1.(int64), _this.V1.(*big.Int))
			case TEventBigInt:
				return _this.V1.(*big.Int).Cmp(that.V1.(*big.Int)) == 0
			case TEventFloat:
				return floatEqualsBigInt(_this.V1.(float64), that.V1.(*big.Int))
			case TEventBigFloat:
				return bigIntEqualsBigFloat(_this.V1.(*big.Int), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return bigIntEqualsDFloat(_this.V1.(*big.Int), that.V1.(compact_float.DFloat))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventFloat:
			switch that.Type {
			case TEventPInt:
				return pintEqualsFloat(that.V1.(uint64), _this.V1.(float64))
			case TEventNInt:
				return nintEqualsFloat(that.V1.(uint64), _this.V1.(float64))
			case TEventInt:
				return intEqualsFloat(that.V1.(int64), _this.V1.(float64))
			case TEventBigInt:
				return floatEqualsBigInt(_this.V1.(float64), that.V1.(*big.Int))
			case TEventFloat:
				return _this.V1.(float64) == that.V1.(float64)
			case TEventBigFloat:
				return floatEqualsBigFloat(_this.V1.(float64), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return floatEqualsDFloat(_this.V1.(float64), that.V1.(compact_float.DFloat))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventBigFloat:
			switch that.Type {
			case TEventPInt:
				return pintEqualsBigFloat(that.V1.(uint64), _this.V1.(*big.Float))
			case TEventNInt:
				return nintEqualsBigFloat(that.V1.(uint64), _this.V1.(*big.Float))
			case TEventInt:
				return intEqualsBigFloat(that.V1.(int64), _this.V1.(*big.Float))
			case TEventBigInt:
				return bigIntEqualsBigFloat(that.V1.(*big.Int), _this.V1.(*big.Float))
			case TEventFloat:
				return floatEqualsBigFloat(that.V1.(float64), _this.V1.(*big.Float))
			case TEventBigFloat:
				return _this.V1.(*big.Float).Cmp(that.V1.(*big.Float)) == 0
			case TEventDecimalFloat:
				return dfloatEqualsBigFloat(that.V1.(compact_float.DFloat), _this.V1.(*big.Float))
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventDecimalFloat:
			switch that.Type {
			case TEventPInt:
				return pintEqualsDFloat(that.V1.(uint64), _this.V1.(compact_float.DFloat))
			case TEventNInt:
				return nintEqualsDFloat(that.V1.(uint64), _this.V1.(compact_float.DFloat))
			case TEventInt:
				return intEqualsDFloat(that.V1.(int64), _this.V1.(compact_float.DFloat))
			case TEventBigInt:
				return bigIntEqualsDFloat(that.V1.(*big.Int), _this.V1.(compact_float.DFloat))
			case TEventFloat:
				return floatEqualsDFloat(that.V1.(float64), _this.V1.(compact_float.DFloat))
			case TEventBigFloat:
				return dfloatEqualsBigFloat(_this.V1.(compact_float.DFloat), that.V1.(*big.Float))
			case TEventDecimalFloat:
				return _this.V1.(compact_float.DFloat) == that.V1.(compact_float.DFloat)
			case TEventBigDecimalFloat:
				return _this.stringify(_this.V1) == that.stringify(that.V1)
			default:
				panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
			}
		case TEventBigDecimalFloat:
			return _this.stringify(_this.V1) == that.stringify(that.V1)
		default:
			panic(fmt.Errorf("TEST BUG: Cannot compare %v to %v", _this.Type, that.Type))
		}
	}

	if _this.Type == TEventTime || _this.Type == TEventCompactTime {
		var err error
		var a compact_time.Time
		var b compact_time.Time

		switch _this.Type {
		case TEventCompactTime:
			a = _this.V1.(compact_time.Time)
		default:
			a = compact_time.AsCompactTime(_this.V1.(time.Time))
			if err = a.Validate(); err != nil {
				panic(err)
			}
		}

		switch that.Type {
		case TEventCompactTime:
			b = that.V1.(compact_time.Time)
		case TEventTime:
			b = compact_time.AsCompactTime(that.V1.(time.Time))
			if err = b.Validate(); err != nil {
				panic(err)
			}
		default:
			return false
		}

		return a.IsEquivalentTo(b)
	}

	return false
}

// Invoking a stored event generates the appropriate data event message to
// the receiver.
func (_this *TEvent) Invoke(receiver events.DataEventReceiver) {
	switch _this.Type {
	case TEventBeginDocument:
		receiver.OnBeginDocument()
	case TEventVersion:
		receiver.OnVersion(_this.V1.(uint64))
	case TEventPadding:
		receiver.OnPadding(_this.V1.(int))
	case TEventComment:
		receiver.OnComment(_this.V1.(bool), []byte(_this.V2.(string)))
	case TEventNull:
		receiver.OnNull()
	case TEventBool:
		receiver.OnBool(_this.V1.(bool))
	case TEventTrue:
		receiver.OnTrue()
	case TEventFalse:
		receiver.OnFalse()
	case TEventPInt:
		receiver.OnPositiveInt(_this.V1.(uint64))
	case TEventNInt:
		receiver.OnNegativeInt(_this.V1.(uint64))
	case TEventInt:
		receiver.OnInt(_this.V1.(int64))
	case TEventBigInt:
		receiver.OnBigInt(_this.V1.(*big.Int))
	case TEventFloat:
		receiver.OnFloat(_this.V1.(float64))
	case TEventBigFloat:
		receiver.OnBigFloat(_this.V1.(*big.Float))
	case TEventDecimalFloat:
		receiver.OnDecimalFloat(_this.V1.(compact_float.DFloat))
	case TEventBigDecimalFloat:
		receiver.OnBigDecimalFloat(_this.V1.(*apd.Decimal))
	case TEventQNan:
		receiver.OnNan(false)
	case TEventSNan:
		receiver.OnNan(true)
	case TEventUID:
		receiver.OnUID(_this.V1.([]byte))
	case TEventTime:
		receiver.OnTime(_this.V1.(time.Time))
	case TEventCompactTime:
		receiver.OnCompactTime(_this.V1.(compact_time.Time))
	case TEventString:
		receiver.OnStringlikeArray(events.ArrayTypeString, _this.V1.(string))
	case TEventResourceID:
		receiver.OnStringlikeArray(events.ArrayTypeResourceID, _this.V1.(string))
	case TEventRemoteRef:
		receiver.OnStringlikeArray(events.ArrayTypeRemoteRef, _this.V1.(string))
	case TEventCustomBinary:
		bytes := _this.V1.([]byte)
		receiver.OnArray(events.ArrayTypeCustomBinary, uint64(len(bytes)), bytes)
	case TEventCustomText:
		receiver.OnStringlikeArray(events.ArrayTypeCustomText, _this.V1.(string))
	case TEventArrayBoolean:
		bitCount := _this.V1.(uint64)
		bytes := _this.V2.([]byte)
		receiver.OnArray(events.ArrayTypeBit, bitCount, bytes)
	case TEventArrayInt8:
		bytes := arrays.Int8SliceAsBytes(_this.V1.([]int8))
		receiver.OnArray(events.ArrayTypeInt8, uint64(len(bytes)), bytes)
	case TEventArrayInt16:
		bytes := arrays.Int16SliceAsBytes(_this.V1.([]int16))
		receiver.OnArray(events.ArrayTypeInt16, uint64(len(bytes)/2), bytes)
	case TEventArrayInt32:
		bytes := arrays.Int32SliceAsBytes(_this.V1.([]int32))
		receiver.OnArray(events.ArrayTypeInt32, uint64(len(bytes)/4), bytes)
	case TEventArrayInt64:
		bytes := arrays.Int64SliceAsBytes(_this.V1.([]int64))
		receiver.OnArray(events.ArrayTypeInt64, uint64(len(bytes)/8), bytes)
	case TEventArrayUint8:
		bytes := _this.V1.([]byte)
		receiver.OnArray(events.ArrayTypeUint8, uint64(len(bytes)), bytes)
	case TEventArrayUint16:
		bytes := arrays.Uint16SliceAsBytes(_this.V1.([]uint16))
		receiver.OnArray(events.ArrayTypeUint16, uint64(len(bytes)/2), bytes)
	case TEventArrayUint32:
		bytes := arrays.Uint32SliceAsBytes(_this.V1.([]uint32))
		receiver.OnArray(events.ArrayTypeUint32, uint64(len(bytes)/4), bytes)
	case TEventArrayUint64:
		bytes := arrays.Uint64SliceAsBytes(_this.V1.([]uint64))
		receiver.OnArray(events.ArrayTypeUint64, uint64(len(bytes)/8), bytes)
	case TEventArrayFloat16:
		// TODO: How to handle float16 in go code?
		bytes := _this.V1.([]byte)
		receiver.OnArray(events.ArrayTypeFloat16, uint64(len(bytes)/2), bytes)
	case TEventArrayFloat32:
		bytes := arrays.Float32SliceAsBytes(_this.V1.([]float32))
		receiver.OnArray(events.ArrayTypeFloat32, uint64(len(bytes)/4), bytes)
	case TEventArrayFloat64:
		bytes := arrays.Float64SliceAsBytes(_this.V1.([]float64))
		receiver.OnArray(events.ArrayTypeFloat64, uint64(len(bytes)/8), bytes)
	case TEventArrayUID:
		bytes := arrays.UUIDSliceAsBytes(_this.V1.([][]byte))
		receiver.OnArray(events.ArrayTypeUID, uint64(len(bytes)/16), bytes)
	case TEventStringBegin:
		receiver.OnArrayBegin(events.ArrayTypeString)
	case TEventResourceIDBegin:
		receiver.OnArrayBegin(events.ArrayTypeResourceID)
	case TEventRemoteRefBegin:
		receiver.OnArrayBegin(events.ArrayTypeRemoteRef)
	case TEventCustomBinaryBegin:
		receiver.OnArrayBegin(events.ArrayTypeCustomBinary)
	case TEventCustomTextBegin:
		receiver.OnArrayBegin(events.ArrayTypeCustomText)
	case TEventArrayBooleanBegin:
		receiver.OnArrayBegin(events.ArrayTypeBit)
	case TEventArrayInt8Begin:
		receiver.OnArrayBegin(events.ArrayTypeInt8)
	case TEventArrayInt16Begin:
		receiver.OnArrayBegin(events.ArrayTypeInt16)
	case TEventArrayInt32Begin:
		receiver.OnArrayBegin(events.ArrayTypeInt32)
	case TEventArrayInt64Begin:
		receiver.OnArrayBegin(events.ArrayTypeInt64)
	case TEventArrayUint8Begin:
		receiver.OnArrayBegin(events.ArrayTypeUint8)
	case TEventArrayUint16Begin:
		receiver.OnArrayBegin(events.ArrayTypeUint16)
	case TEventArrayUint32Begin:
		receiver.OnArrayBegin(events.ArrayTypeUint32)
	case TEventArrayUint64Begin:
		receiver.OnArrayBegin(events.ArrayTypeUint64)
	case TEventArrayFloat16Begin:
		receiver.OnArrayBegin(events.ArrayTypeFloat16)
	case TEventArrayFloat32Begin:
		receiver.OnArrayBegin(events.ArrayTypeFloat32)
	case TEventArrayFloat64Begin:
		receiver.OnArrayBegin(events.ArrayTypeFloat64)
	case TEventArrayUIDBegin:
		receiver.OnArrayBegin(events.ArrayTypeUID)
	case TEventMediaBegin:
		receiver.OnArrayBegin(events.ArrayTypeMedia)
	case TEventArrayChunk:
		receiver.OnArrayChunk(_this.V1.(uint64), _this.V2.(bool))
	case TEventArrayData:
		receiver.OnArrayData(_this.V1.([]byte))
	case TEventList:
		receiver.OnList()
	case TEventMap:
		receiver.OnMap()
	case TEventMarkup:
		receiver.OnMarkup([]byte(_this.V1.(string)))
	case TEventEnd:
		receiver.OnEnd()
	case TEventNode:
		receiver.OnNode()
	case TEventEdge:
		receiver.OnEdge()
	case TEventMarker:
		receiver.OnMarker([]byte(_this.V1.(string)))
	case TEventReference:
		receiver.OnReference([]byte(_this.V1.(string)))
	case TEventConstant:
		receiver.OnConstant([]byte(_this.V1.(string)))
	case TEventEndDocument:
		receiver.OnEndDocument()
	default:
		panic(fmt.Errorf("%v: Unhandled event type", _this.Type))
	}
}

// Comparators

func pintEqualsInt(a uint64, b int64) bool {
	if b < 0 {
		return false
	}
	return a == uint64(b)
}

func pintEqualsFloat(a uint64, b float64) bool {
	return float64(a) == b
}

func pintEqualsDFloat(a uint64, b compact_float.DFloat) bool {
	df, err := compact_float.DFloatFromUInt(a)
	if err != nil {
		return false
	}
	return df == b
}

func pintEqualsBigInt(a uint64, b *big.Int) bool {
	bi := &big.Int{}
	bi.SetUint64(a)
	return bi.Cmp(b) == 0
}

func pintEqualsBigFloat(a uint64, b *big.Float) bool {
	bf := &big.Float{}
	bf.SetUint64(a)
	return bf.Cmp(b) == 0
}

func nintEqualsInt(a uint64, b int64) bool {
	if b >= 0 {
		return false
	}
	if a > math.MaxInt64+1 {
		return false
	}
	return -int64(a) == b
}

func nintEqualsFloat(a uint64, b float64) bool {
	return -float64(a) == b
}

func nintEqualsDFloat(a uint64, b compact_float.DFloat) bool {
	var df compact_float.DFloat
	var err error

	if a == 0 && b == compact_float.NegativeZero() {
		return true
	}

	if a <= math.MaxInt64 {
		df = compact_float.DFloatValue(0, -int64(a))
	} else {
		df, err = compact_float.DFloatFromString(fmt.Sprintf("-%d", a))
		if err != nil {
			return false
		}
	}
	return df == b
}

func nintEqualsBigInt(a uint64, b *big.Int) bool {
	bi := &big.Int{}
	bi.SetUint64(a)
	bi.Neg(bi)
	return bi.Cmp(b) == 0
}

func nintEqualsBigFloat(a uint64, b *big.Float) bool {
	bf := &big.Float{}
	bf.SetUint64(a)
	bf.Neg(bf)
	return bf.Cmp(b) == 0
}

func intEqualsFloat(a int64, b float64) bool {
	return float64(a) == b
}

func intEqualsDFloat(a int64, b compact_float.DFloat) bool {
	return compact_float.DFloatValue(0, a) == b
}

func intEqualsBigInt(a int64, b *big.Int) bool {
	bi := &big.Int{}
	bi.SetInt64(a)
	return bi.Cmp(b) == 0
}

func intEqualsBigFloat(a int64, b *big.Float) bool {
	bf := &big.Float{}
	bf.SetInt64(a)
	return bf.Cmp(b) == 0
}

func floatEqualsBigInt(a float64, b *big.Int) bool {
	fa, cond, err := apd.NewFromString(strconv.FormatFloat(a, 'x', -1, 64))
	if err != nil {
		panic(err)
	}
	if cond != 0 {
		return false
	}
	fb := apd.NewWithBigInt(b, 0)
	return fa.Cmp(fb) == 0
}

func floatEqualsDFloat(a float64, b compact_float.DFloat) bool {
	df, err := compact_float.DFloatFromFloat64(a, 20)
	if err != nil {
		return false
	}
	return df == b
}

func floatEqualsBigFloat(a float64, b *big.Float) bool {
	return big.NewFloat(a).Cmp(b) == 0
}

func dfloatEqualsBigFloat(a compact_float.DFloat, b *big.Float) bool {
	return a.BigFloat().Cmp(b) == 0
}

func bigIntEqualsBigFloat(a *big.Int, b *big.Float) bool {
	bf := &big.Float{}
	bf.SetInt(a)
	return bf.Cmp(b) == 0
}

func bigIntEqualsDFloat(a *big.Int, b compact_float.DFloat) bool {
	df, err := compact_float.DFloatFromBigInt(a)
	if err != nil {
		return false
	}
	return df == b
}
