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
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

func NewBigInt(str string) *big.Int {
	bi := new(big.Int)
	_, success := bi.SetString(str, 10)
	if !success {
		panic(fmt.Errorf("Cannot convert %v to big.Int", str))
	}
	return bi
}

func NewBigFloat(str string, significantDigits int) *big.Float {
	f, _, err := big.ParseFloat(str, 10, uint(conversions.DecimalDigitsToBits(significantDigits)), big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return f
}

func NewDFloat(str string) compact_float.DFloat {
	df, err := compact_float.DFloatFromString(str)
	if err != nil {
		panic(err)
	}
	return df
}

func NewBDF(str string) *apd.Decimal {
	v, _, err := apd.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return v
}

func NewURL(str string) *url.URL {
	v, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return v
}

func NewURI(uriString string) *url.URL {
	uri, err := url.Parse(uriString)
	if err != nil {
		fmt.Printf("TEST CODE BUG: Bad URL (%v): %v", uriString, err)
		panic(err)
	}
	return uri
}

func ReportPanic(function func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	function()
	return
}

func AssertNoPanic(t *testing.T, function func()) {
	if err := ReportPanic(function); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func AssertPanics(t *testing.T, function func()) {
	if err := ReportPanic(function); err == nil {
		t.Errorf("Expected an error")
	}
}

func GenerateString(charCount int, startIndex int) string {
	charRange := int('z' - 'a')
	var object strings.Builder
	for i := 0; i < charCount; i++ {
		ch := 'a' + (i+charCount+startIndex)%charRange
		object.WriteByte(byte(ch))
	}
	return object.String()
}

func GenerateBytes(length int, startIndex int) []byte {
	return []byte(GenerateString(length, startIndex))
}

func InvokeEvents(receiver events.DataEventReceiver, events ...*TEvent) {
	for _, event := range events {
		event.Invoke(receiver)
	}
}

type TEventType int

const (
	TEventVersion TEventType = iota
	TEventPadding
	TEventNil
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
	TEventNan
	TEventSNan
	TEventUUID
	TEventTime
	TEventCompactTime
	TEventBytes
	TEventString
	TEventURI
	TEventCustom
	TEventBytesBegin
	TEventStringBegin
	TEventURIBegin
	TEventCustomBegin
	TEventArrayChunk
	TEventArrayData
	TEventList
	TEventMap
	TEventMarkup
	TEventMetadata
	TEventComment
	TEventEnd
	TEventMarker
	TEventReference
	TEventEndDocument
)

var TEventNames = []string{
	"V",
	"PAD",
	"N",
	"B",
	"TT",
	"FF",
	"PI",
	"NI",
	"I",
	"BI",
	"F",
	"BF",
	"DF",
	"BDF",
	"CPLX",
	"NAN",
	"SNAN",
	"UUID",
	"GT",
	"CT",
	"BIN",
	"S",
	"URI",
	"CUST",
	"BB",
	"SB",
	"UB",
	"CB",
	"AC",
	"AD",
	"L",
	"M",
	"MUP",
	"META",
	"CMT",
	"E",
	"MARK",
	"REF",
	"ED",
}

func (_this TEventType) String() string {
	return TEventNames[_this]
}

type TEvent struct {
	Type TEventType
	V1   interface{}
	V2   interface{}
}

func newTEvent(eventType TEventType, v1 interface{}, v2 interface{}) *TEvent {
	return &TEvent{
		Type: eventType,
		V1:   v1,
		V2:   v2,
	}
}

func (_this *TEvent) String() string {
	str := _this.Type.String()
	if _this.V1 != nil {
		if _this.V2 != nil {
			return fmt.Sprintf("%v(%v,%v)", str, _this.V1, _this.V2)
		}
		return fmt.Sprintf("%v(%v)", str, _this.V1)
	}
	return str
}

func (_this *TEvent) Invoke(receiver events.DataEventReceiver) {
	switch _this.Type {
	case TEventVersion:
		receiver.OnVersion(_this.V1.(uint64))
	case TEventPadding:
		receiver.OnPadding(_this.V1.(int))
	case TEventNil:
		receiver.OnNil()
	case TEventBool:
		receiver.OnBool(_this.V1.(bool))
	case TEventTrue:
		receiver.OnBool(true)
	case TEventFalse:
		receiver.OnBool(false)
	case TEventPInt:
		receiver.OnPositiveInt(_this.V1.(uint64))
	case TEventNInt:
		receiver.OnNegativeInt(_this.V1.(uint64))
	case TEventInt:
		v := _this.V1.(int64)
		if v >= 0 {
			receiver.OnPositiveInt(uint64(v))
		} else {
			receiver.OnNegativeInt(uint64(-v))
		}
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
	case TEventNan:
		receiver.OnNan(false)
	case TEventSNan:
		receiver.OnNan(true)
	case TEventUUID:
		receiver.OnUUID(_this.V1.([]byte))
	case TEventTime:
		receiver.OnTime(_this.V1.(time.Time))
	case TEventCompactTime:
		receiver.OnCompactTime(_this.V1.(*compact_time.Time))
	case TEventBytes:
		receiver.OnBytes(_this.V1.([]byte))
	case TEventString:
		receiver.OnString(_this.V1.(string))
	case TEventURI:
		receiver.OnURI(_this.V1.(string))
	case TEventCustom:
		receiver.OnCustom(_this.V1.([]byte))
	case TEventBytesBegin:
		receiver.OnBytesBegin()
	case TEventStringBegin:
		receiver.OnStringBegin()
	case TEventURIBegin:
		receiver.OnURIBegin()
	case TEventCustomBegin:
		receiver.OnCustomBegin()
	case TEventArrayChunk:
		receiver.OnArrayChunk(_this.V1.(uint64), _this.V2.(bool))
	case TEventArrayData:
		receiver.OnArrayData(_this.V1.([]byte))
	case TEventList:
		receiver.OnList()
	case TEventMap:
		receiver.OnMap()
	case TEventMarkup:
		receiver.OnMarkup()
	case TEventMetadata:
		receiver.OnMetadata()
	case TEventComment:
		receiver.OnComment()
	case TEventEnd:
		receiver.OnEnd()
	case TEventMarker:
		receiver.OnMarker()
	case TEventReference:
		receiver.OnReference()
	case TEventEndDocument:
		receiver.OnEndDocument()
	default:
		panic(fmt.Errorf("%v: Unhandled event type", _this.Type))
	}
}

func EventOrNil(eventType TEventType, value interface{}) *TEvent {
	if value == nil {
		eventType = TEventNil
	}
	return newTEvent(eventType, value, nil)
}

func TT() *TEvent { return newTEvent(TEventBool, true, nil) }
func FF() *TEvent { return newTEvent(TEventBool, false, nil) }
func I(v int64) *TEvent {
	if v >= 0 {
		return PI(uint64(v))
	}
	return NI(uint64(-v))
}
func F(v float64) *TEvent {
	// TODO: Do I need to check for this? Doesn't the library handle it?
	if math.IsNaN(v) {
		if common.IsSignalingNan(v) {
			return SNAN()
		}
		return NAN()
	}
	return newTEvent(TEventFloat, v, nil)
}

func BF(v *big.Float) *TEvent           { return EventOrNil(TEventBigFloat, v) }
func DF(v compact_float.DFloat) *TEvent { return newTEvent(TEventDecimalFloat, v, nil) }
func BDF(v *apd.Decimal) *TEvent        { return EventOrNil(TEventBigDecimalFloat, v) }
func V(v uint64) *TEvent                { return newTEvent(TEventVersion, v, nil) }
func N() *TEvent                        { return newTEvent(TEventNil, nil, nil) }
func PAD(v int) *TEvent                 { return newTEvent(TEventPadding, v, nil) }
func B(v bool) *TEvent                  { return newTEvent(TEventBool, v, nil) }
func PI(v uint64) *TEvent               { return newTEvent(TEventPInt, v, nil) }
func NI(v uint64) *TEvent               { return newTEvent(TEventNInt, v, nil) }
func BI(v *big.Int) *TEvent             { return EventOrNil(TEventBigInt, v) }
func NAN() *TEvent                      { return newTEvent(TEventNan, nil, nil) }
func SNAN() *TEvent                     { return newTEvent(TEventSNan, nil, nil) }
func UUID(v []byte) *TEvent             { return newTEvent(TEventUUID, v, nil) }
func GT(v time.Time) *TEvent            { return newTEvent(TEventTime, v, nil) }
func CT(v *compact_time.Time) *TEvent   { return EventOrNil(TEventCompactTime, v) }
func BIN(v []byte) *TEvent              { return newTEvent(TEventBytes, v, nil) }
func S(v string) *TEvent                { return newTEvent(TEventString, v, nil) }
func URI(v string) *TEvent              { return newTEvent(TEventURI, v, nil) }
func CUST(v []byte) *TEvent             { return newTEvent(TEventCustom, v, nil) }
func BB() *TEvent                       { return newTEvent(TEventBytesBegin, nil, nil) }
func SB() *TEvent                       { return newTEvent(TEventStringBegin, nil, nil) }
func UB() *TEvent                       { return newTEvent(TEventURIBegin, nil, nil) }
func CB() *TEvent                       { return newTEvent(TEventCustomBegin, nil, nil) }
func AC(l uint64, term bool) *TEvent    { return newTEvent(TEventArrayChunk, l, term) }
func AD(v []byte) *TEvent               { return newTEvent(TEventArrayData, v, nil) }
func L() *TEvent                        { return newTEvent(TEventList, nil, nil) }
func M() *TEvent                        { return newTEvent(TEventMap, nil, nil) }
func MUP() *TEvent                      { return newTEvent(TEventMarkup, nil, nil) }
func META() *TEvent                     { return newTEvent(TEventMetadata, nil, nil) }
func CMT() *TEvent                      { return newTEvent(TEventComment, nil, nil) }
func E() *TEvent                        { return newTEvent(TEventEnd, nil, nil) }
func MARK() *TEvent                     { return newTEvent(TEventMarker, nil, nil) }
func REF() *TEvent                      { return newTEvent(TEventReference, nil, nil) }
func ED() *TEvent                       { return newTEvent(TEventEndDocument, nil, nil) }

func EventForValue(value interface{}) *TEvent {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return N()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return B(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return I(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return PI(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return F(rv.Float())
	case reflect.String:
		return S(rv.String())
	case reflect.Slice:
		switch rv.Type() {
		case common.TypeBytes:
			return BIN(rv.Bytes())
		}
	case reflect.Ptr:
		if rv.IsNil() {
			return N()
		}
		switch rv.Type() {
		case common.TypePBigDecimalFloat:
			return BDF(rv.Interface().(*apd.Decimal))
		case common.TypePBigFloat:
			return BF(rv.Interface().(*big.Float))
		case common.TypePBigInt:
			return BI(rv.Interface().(*big.Int))
		case common.TypePCompactTime:
			return CT(rv.Interface().(*compact_time.Time))
		case common.TypePURL:
			return URI(rv.Interface().(*url.URL).String())
		}
		return EventForValue(rv.Elem().Interface())
	case reflect.Struct:
		switch rv.Type() {
		case common.TypeBigDecimalFloat:
			v := rv.Interface().(apd.Decimal)
			return BDF(&v)
		case common.TypeBigFloat:
			v := rv.Interface().(big.Float)
			return BF(&v)
		case common.TypeBigInt:
			v := rv.Interface().(big.Int)
			return BI(&v)
		case common.TypeCompactTime:
			v := rv.Interface().(compact_time.Time)
			return CT(&v)
		case common.TypeDFloat:
			v := rv.Interface().(compact_float.DFloat)
			return DF(v)
		case common.TypeTime:
			v := rv.Interface().(time.Time)
			return GT(v)
		case common.TypeURL:
			v := rv.Interface().(url.URL)
			return URI(v.String())
		}
	}
	panic(fmt.Errorf("TEST CODE BUG: Unhandled kind: %v", rv.Kind()))
}

type TER struct {
	Events []*TEvent
}

func NewTER() *TER {
	return &TER{}
}
func (h *TER) add(event *TEvent) {
	h.Events = append(h.Events, event)
}
func (h *TER) OnVersion(version uint64)   { h.add(V(version)) }
func (h *TER) OnPadding(count int)        { h.add(PAD(count)) }
func (h *TER) OnNil()                     { h.add(N()) }
func (h *TER) OnBool(value bool)          { h.add(B(value)) }
func (h *TER) OnTrue()                    { h.add(TT()) }
func (h *TER) OnFalse()                   { h.add(FF()) }
func (h *TER) OnPositiveInt(value uint64) { h.add(PI(value)) }
func (h *TER) OnNegativeInt(value uint64) { h.add(NI(value)) }
func (h *TER) OnInt(value int64)          { h.add(I(value)) }
func (h *TER) OnBigInt(value *big.Int)    { h.add(BI(value)) }
func (h *TER) OnFloat(value float64) {
	if math.IsNaN(value) {
		if common.IsSignalingNan(value) {
			h.add(SNAN())
		} else {
			h.add(NAN())
		}
	} else {
		h.add(F(value))
	}
}
func (h *TER) OnBigFloat(value *big.Float) { h.add(newTEvent(TEventBigFloat, value, nil)) }
func (h *TER) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			h.add(SNAN())
		} else {
			h.add(NAN())
		}
	} else {
		h.add(DF(value))
	}
}
func (h *TER) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		h.add(NAN())
	case apd.NaNSignaling:
		h.add(SNAN())
	default:
		h.add(BDF(value))
	}
}
func (h *TER) OnUUID(value []byte)                    { h.add(UUID(value)) }
func (h *TER) OnTime(value time.Time)                 { h.add(GT(value)) }
func (h *TER) OnCompactTime(value *compact_time.Time) { h.add(CT(value)) }
func (h *TER) OnBytes(value []byte)                   { h.add(BIN(value)) }
func (h *TER) OnString(value string)                  { h.add(S(value)) }
func (h *TER) OnURI(value string)                     { h.add(URI((value))) }
func (h *TER) OnCustom(value []byte)                  { h.add(CUST(value)) }
func (h *TER) OnBytesBegin()                          { h.add(BB()) }
func (h *TER) OnStringBegin()                         { h.add(SB()) }
func (h *TER) OnURIBegin()                            { h.add(UB()) }
func (h *TER) OnCustomBegin()                         { h.add(CB()) }
func (h *TER) OnArrayChunk(l uint64, final bool)      { h.add(AC(l, final)) }
func (h *TER) OnArrayData(data []byte)                { h.add(AD(data)) }
func (h *TER) OnList()                                { h.add(L()) }
func (h *TER) OnMap()                                 { h.add(M()) }
func (h *TER) OnMarkup()                              { h.add(MUP()) }
func (h *TER) OnMetadata()                            { h.add(META()) }
func (h *TER) OnComment()                             { h.add(CMT()) }
func (h *TER) OnEnd()                                 { h.add(E()) }
func (h *TER) OnMarker()                              { h.add(MARK()) }
func (h *TER) OnReference()                           { h.add(REF()) }
func (h *TER) OnEndDocument()                         { h.add(ED()) }
func (h *TER) OnNan(signaling bool) {
	if signaling {
		h.add(SNAN())
	} else {
		h.add(NAN())
	}
}

type TestingInnerStruct struct {
	Inner int
}

type TestingOuterStruct struct {
	Bo     bool
	PBo    *bool
	By     byte
	PBy    *byte
	I      int
	PI     *int
	I8     int8
	PI8    *int8
	I16    int16
	PI16   *int16
	I32    int32
	PI32   *int32
	I64    int64
	PI64   *int64
	U      uint
	PU     *uint
	U8     uint8
	PU8    *uint8
	U16    uint16
	PU16   *uint16
	U32    uint32
	PU32   *uint32
	U64    uint64
	PU64   *uint64
	BI     big.Int
	PBI    *big.Int
	F32    float32
	PF32   *float32
	F64    float64
	PF64   *float64
	BF     big.Float
	PBF    *big.Float
	DF     compact_float.DFloat
	BDF    apd.Decimal
	PBDF   *apd.Decimal
	Ar     [4]byte
	St     string
	Ba     []byte
	Sl     []interface{}
	M      map[interface{}]interface{}
	IS     TestingInnerStruct
	PIS    *TestingInnerStruct
	Time   time.Time
	PTime  *time.Time
	CTime  compact_time.Time
	PCTime *compact_time.Time
	PURL   *url.URL
	// TODO: If URI is relative, doesn't need scheme. Can be 0 length...
	// URL   url.URL
}

func (_this *TestingOuterStruct) GetRepresentativeEvents(includeFakes bool) (events []*TEvent) {
	ade := func(e ...*TEvent) {
		events = append(events, e...)
	}
	adv := func(value interface{}) {
		ade(EventForValue(value))
	}
	anv := func(name string, value interface{}) {
		adv(name)
		adv(value)
	}
	ane := func(name string, e ...*TEvent) {
		adv(name)
		ade(e...)
	}

	ade(M())

	anv("Bo", _this.Bo)
	anv("PBo", _this.PBo)
	anv("By", _this.By)
	anv("PBy", _this.PBy)
	anv("I", _this.I)
	anv("PI", _this.PI)
	anv("I8", _this.I8)
	anv("PI8", _this.PI8)
	anv("I16", _this.I16)
	anv("PI16", _this.PI16)
	anv("I32", _this.I32)
	anv("PI32", _this.PI32)
	anv("I64", _this.I64)
	anv("PI64", _this.PI64)
	anv("U", _this.U)
	anv("PU", _this.PU)
	anv("U8", _this.U8)
	anv("PU8", _this.PU8)
	anv("U16", _this.U16)
	anv("PU16", _this.PU16)
	anv("U32", _this.U32)
	anv("PU32", _this.PU32)
	anv("U64", _this.U64)
	anv("PU64", _this.PU64)
	anv("BI", _this.BI)
	anv("PBI", _this.PBI)
	anv("F32", _this.F32)
	anv("PF32", _this.PF32)
	anv("F64", _this.F64)
	anv("PF64", _this.PF64)
	anv("BF", _this.BF)
	anv("PBF", _this.PBF)
	anv("DF", _this.DF)
	anv("BDF", _this.BDF)
	anv("PBDF", _this.PBDF)
	ane("Ar", BIN(_this.Ar[:]))
	anv("St", _this.St)
	anv("Ba", _this.Ba)

	ane("Sl", L())
	for _, v := range _this.Sl {
		adv(v)
	}
	ade(E())

	ane("M", M())
	for k, v := range _this.M {
		adv(k)
		adv(v)
	}
	ade(E())

	ane("IS", M())
	anv("Inner", _this.IS.Inner)
	ade(E())

	if _this.PIS != nil {
		ane("PIS", M())
		anv("Inner", _this.PIS.Inner)
		ade(E())
	}

	anv("Time", _this.Time)
	anv("PTime", _this.PTime)
	anv("CTime", _this.CTime)
	anv("PCTime", _this.PCTime)
	anv("PURL", _this.PURL)

	if includeFakes {
		ane("F1", B(true))
		ane("F2", B(false))
		ane("F3", I(1))
		ane("F4", I(-1))
		ane("F5", F(1.1))
		ane("F6", BF(NewBigFloat("1.1", 2)))
		ane("F7", DF(NewDFloat("1.1")))
		ane("F8", BDF(NewBDF("1.1")))
		ane("F9", N())
		ane("F10", BI(NewBigInt("1000")))
		// ane("F11", cplx(1+1i))
		ane("F12", NAN())
		ane("F13", SNAN())
		ane("F14", UUID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
		ane("F15", GT(_this.Time))
		ane("F16", CT(_this.PCTime))
		ane("F17", BIN([]byte{1}))
		ane("F18", S("xyz"))
		ane("F19", URI("http://example.com"))
		// ane("F20", cust([]byte{1}))
		ane("FakeList", L(), E())
		ane("FakeMap", M(), E())
		ane("FakeDeep", L(), M(), S("A"), L(),
			B(true),
			B(false),
			I(1),
			I(-1),
			F(1.1),
			BF(NewBigFloat("1.1", 2)),
			DF(NewDFloat("1.1")),
			BDF(NewBDF("1.1")),
			N(),
			BI(NewBigInt("1000")),
			// cplx(1+1i),
			NAN(),
			SNAN(),
			UUID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
			GT(_this.Time),
			CT(_this.PCTime),
			BIN([]byte{1}),
			S("xyz"),
			URI("http://example.com"),
			// cust([]byte{1}),
			E(), E(), E())
	}

	ade(E())
	return
}

func NewTestingOuterStruct(baseValue int) *TestingOuterStruct {
	_this := new(TestingOuterStruct)
	_this.Init(baseValue)
	return _this
}

func NewBlankTestingOuterStruct() *TestingOuterStruct {
	_this := new(TestingOuterStruct)
	_this.CTime.Year = 1
	_this.CTime.Month = 1
	_this.CTime.Day = 1
	return _this
}

func (_this *TestingOuterStruct) Init(baseValue int) {
	_this.Bo = baseValue&1 == 1
	_this.PBo = &_this.Bo
	_this.By = byte(baseValue + int(unsafe.Offsetof(_this.By)))
	_this.PBy = &_this.By
	_this.I = 100000 + baseValue + int(unsafe.Offsetof(_this.I))
	_this.PI = &_this.I
	_this.I8 = int8(baseValue + int(unsafe.Offsetof(_this.I8)))
	_this.PI8 = &_this.I8
	_this.I16 = int16(-1000 - baseValue - int(unsafe.Offsetof(_this.I16)))
	_this.PI16 = &_this.I16
	_this.I32 = int32(1000000000 + baseValue + int(unsafe.Offsetof(_this.I32)))
	_this.PI32 = &_this.I32
	_this.I64 = int64(1000000000000000 + baseValue + int(unsafe.Offsetof(_this.I64)))
	_this.PI64 = &_this.I64
	_this.U = uint(1000000 + baseValue + int(unsafe.Offsetof(_this.U)))
	_this.PU = &_this.U
	_this.U8 = uint8(baseValue + int(unsafe.Offsetof(_this.U8)))
	_this.PU8 = &_this.U8
	_this.U16 = uint16(10000 + baseValue + int(unsafe.Offsetof(_this.U16)))
	_this.PU16 = &_this.U16
	_this.U32 = uint32(100000000 + baseValue + int(unsafe.Offsetof(_this.U32)))
	_this.PU32 = &_this.U32
	_this.U64 = uint64(1000000000000 + baseValue + int(unsafe.Offsetof(_this.U64)))
	_this.PU64 = &_this.U64
	_this.PBI = NewBigInt(fmt.Sprintf("-10000000000000000000000000000000000000%v", unsafe.Offsetof(_this.PBI)))
	_this.BI = *_this.PBI
	_this.F32 = float32(1000000+baseValue+int(unsafe.Offsetof(_this.F32))) + 0.5
	_this.PF32 = &_this.F32
	_this.F64 = float64(1000000000000+baseValue+int(unsafe.Offsetof(_this.F64))) + 0.5
	_this.PF64 = &_this.F64
	_this.PBF = NewBigFloat("12345678901234567890123.1234567", 30)
	_this.BF = *_this.PBF
	_this.DF = NewDFloat(fmt.Sprintf("-100000000000000%ve-1000000", unsafe.Offsetof(_this.DF)))
	_this.PBDF = NewBDF("-1.234567890123456789777777777777777777771234e-10000")
	_this.BDF = *_this.PBDF
	_this.Ar[0] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)))
	_this.Ar[1] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+1))
	_this.Ar[2] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+2))
	_this.Ar[3] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+3))
	_this.St = GenerateString(baseValue+5, baseValue)
	_this.Ba = GenerateBytes(baseValue+1, baseValue)
	_this.M = make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		_this.Sl = append(_this.Sl, i)
		_this.M[fmt.Sprintf("key%v", i)] = i
	}
	_this.IS.Inner = baseValue + 15
	_this.PIS = new(TestingInnerStruct)
	_this.PIS.Inner = baseValue + 16
	testTime := time.Date(30000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	_this.PTime = &testTime
	_this.PCTime = compact_time.NewTimestamp(-1000, 1, 1, 1, 1, 1, 1, "Europe/Berlin")
	_this.CTime = *_this.PCTime
	_this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}
