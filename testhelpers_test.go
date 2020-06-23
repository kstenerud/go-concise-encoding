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

package concise_encoding

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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

func newBigInt(str string) *big.Int {
	bi := new(big.Int)
	_, success := bi.SetString(str, 10)
	if !success {
		panic(fmt.Errorf("Cannot convert %v to big.Int", str))
	}
	return bi
}

func newBigFloat(str string, significantDigits int) *big.Float {
	f, _, err := big.ParseFloat(str, 10, uint(decimalDigitsToBits(significantDigits)), big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return f
}

func newDFloat(str string) compact_float.DFloat {
	df, err := compact_float.DFloatFromString(str)
	if err != nil {
		panic(err)
	}
	return df
}

func newBDF(str string) *apd.Decimal {
	v, _, err := apd.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return v
}

func newURL(str string) *url.URL {
	v, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return v
}

func newURI(uriString string) *url.URL {
	uri, err := url.Parse(uriString)
	if err != nil {
		fmt.Printf("ERROR ERROR ERROR BUG: Bad URL (%v): %v", uriString, err)
		panic(err)
	}
	return uri
}

func reportPanic(function func()) (err error) {
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

func assertNoPanic(t *testing.T, function func()) {
	if err := reportPanic(function); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func assertPanics(t *testing.T, function func()) {
	if err := reportPanic(function); err == nil {
		t.Errorf("Expected an error")
	}
}

func generateString(charCount int, startIndex int) string {
	charRange := int('z' - 'a')
	var object strings.Builder
	for i := 0; i < charCount; i++ {
		ch := 'a' + (i+charCount+startIndex)%charRange
		object.WriteByte(byte(ch))
	}
	return object.String()
}

func generateBytes(length int, startIndex int) []byte {
	return []byte(generateString(length, startIndex))
}

func invokeEvents(receiver DataEventReceiver, events ...*tevent) {
	for _, event := range events {
		event.Invoke(receiver)
	}
}

func cbeDecode(document []byte) (events []*tevent, err error) {
	receiver := NewTER()
	err = CBEDecode(document, receiver, nil)
	events = receiver.Events
	return
}

func cbeEncodeDecode(expected ...*tevent) (events []*tevent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	encoder := NewCBEEncoder(nil)
	invokeEvents(encoder, expected...)
	document := encoder.GetBuiltDocument()

	return cbeDecode(document)
}

func cteDecode(document []byte) (events []*tevent, err error) {
	receiver := NewTER()
	err = CTEDecode(document, receiver, nil)
	events = receiver.Events
	return
}

func cteEncode(events ...*tevent) []byte {
	encoder := NewCTEEncoder(nil)
	invokeEvents(encoder, events...)
	return encoder.GetBuiltDocument()
}

func cteEncodeDecode(events ...*tevent) (decodedEvents []*tevent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	return cteDecode(cteEncode(events...))
}

type teventType int

const (
	teventVersion teventType = iota
	teventPadding
	teventNil
	teventBool
	teventTrue
	teventFalse
	teventPInt
	teventNInt
	teventInt
	teventBigInt
	teventFloat
	teventBigFloat
	teventDecimalFloat
	teventBigDecimalFloat
	teventComplex
	teventNan
	teventSNan
	teventUUID
	teventTime
	teventCompactTime
	teventBytes
	teventString
	teventURI
	teventCustom
	teventBytesBegin
	teventStringBegin
	teventURIBegin
	teventCustomBegin
	teventArrayChunk
	teventArrayData
	teventList
	teventMap
	teventMarkup
	teventMetadata
	teventComment
	teventEnd
	teventMarker
	teventReference
	teventEndDocument
)

var teventNames = []string{
	"v",
	"pad",
	"n",
	"b",
	"tt",
	"ff",
	"pi",
	"ni",
	"i",
	"bi",
	"f",
	"bf",
	"df",
	"bdf",
	"cplx",
	"nan",
	"snan",
	"uuid",
	"gt",
	"ct",
	"bin",
	"s",
	"uri",
	"cust",
	"bb",
	"sb",
	"ub",
	"cb",
	"ac",
	"ad",
	"l",
	"m",
	"mup",
	"meta",
	"cmt",
	"e",
	"mark",
	"ref",
	"ed",
}

func (_this teventType) String() string {
	return teventNames[_this]
}

type tevent struct {
	Type teventType
	V1   interface{}
	V2   interface{}
}

func newTEvent(eventType teventType, v1 interface{}, v2 interface{}) *tevent {
	return &tevent{
		Type: eventType,
		V1:   v1,
		V2:   v2,
	}
}

func (_this *tevent) String() string {
	str := _this.Type.String()
	if _this.V1 != nil {
		if _this.V2 != nil {
			return fmt.Sprintf("%v(%v,%v)", str, _this.V1, _this.V2)
		}
		return fmt.Sprintf("%v(%v)", str, _this.V1)
	}
	return str
}

func (_this *tevent) Invoke(receiver DataEventReceiver) {
	switch _this.Type {
	case teventVersion:
		receiver.OnVersion(_this.V1.(uint64))
	case teventPadding:
		receiver.OnPadding(_this.V1.(int))
	case teventNil:
		receiver.OnNil()
	case teventBool:
		receiver.OnBool(_this.V1.(bool))
	case teventTrue:
		receiver.OnBool(true)
	case teventFalse:
		receiver.OnBool(false)
	case teventPInt:
		receiver.OnPositiveInt(_this.V1.(uint64))
	case teventNInt:
		receiver.OnNegativeInt(_this.V1.(uint64))
	case teventInt:
		v := _this.V1.(int64)
		if v >= 0 {
			receiver.OnPositiveInt(uint64(v))
		} else {
			receiver.OnNegativeInt(uint64(-v))
		}
	case teventBigInt:
		receiver.OnBigInt(_this.V1.(*big.Int))
	case teventFloat:
		receiver.OnFloat(_this.V1.(float64))
	case teventBigFloat:
		receiver.OnBigFloat(_this.V1.(*big.Float))
	case teventDecimalFloat:
		receiver.OnDecimalFloat(_this.V1.(compact_float.DFloat))
	case teventBigDecimalFloat:
		receiver.OnBigDecimalFloat(_this.V1.(*apd.Decimal))
	case teventComplex:
		receiver.OnComplex(_this.V1.(complex128))
	case teventNan:
		receiver.OnNan(false)
	case teventSNan:
		receiver.OnNan(true)
	case teventUUID:
		receiver.OnUUID(_this.V1.([]byte))
	case teventTime:
		receiver.OnTime(_this.V1.(time.Time))
	case teventCompactTime:
		receiver.OnCompactTime(_this.V1.(*compact_time.Time))
	case teventBytes:
		receiver.OnBytes(_this.V1.([]byte))
	case teventString:
		receiver.OnString(_this.V1.(string))
	case teventURI:
		receiver.OnURI(_this.V1.(string))
	case teventCustom:
		receiver.OnCustom(_this.V1.([]byte))
	case teventBytesBegin:
		receiver.OnBytesBegin()
	case teventStringBegin:
		receiver.OnStringBegin()
	case teventURIBegin:
		receiver.OnURIBegin()
	case teventCustomBegin:
		receiver.OnCustomBegin()
	case teventArrayChunk:
		receiver.OnArrayChunk(_this.V1.(uint64), _this.V2.(bool))
	case teventArrayData:
		receiver.OnArrayData(_this.V1.([]byte))
	case teventList:
		receiver.OnList()
	case teventMap:
		receiver.OnMap()
	case teventMarkup:
		receiver.OnMarkup()
	case teventMetadata:
		receiver.OnMetadata()
	case teventComment:
		receiver.OnComment()
	case teventEnd:
		receiver.OnEnd()
	case teventMarker:
		receiver.OnMarker()
	case teventReference:
		receiver.OnReference()
	case teventEndDocument:
		receiver.OnEndDocument()
	default:
		panic(fmt.Errorf("%v: Unhandled event type", _this.Type))
	}
}

func eventOrNil(eventType teventType, value interface{}) *tevent {
	if value == nil {
		eventType = teventNil
	}
	return newTEvent(eventType, value, nil)
}

func tt() *tevent { return newTEvent(teventBool, true, nil) }
func ff() *tevent { return newTEvent(teventBool, false, nil) }
func i(v int64) *tevent {
	if v >= 0 {
		return pi(uint64(v))
	}
	return ni(uint64(-v))
}
func f(v float64) *tevent {
	// TODO: Do I need to check for this? Doesn't the library handle it?
	if math.IsNaN(v) {
		if isSignalingNan(v) {
			return snan()
		}
		return nan()
	}
	return newTEvent(teventFloat, v, nil)
}

func bf(v *big.Float) *tevent           { return eventOrNil(teventBigFloat, v) }
func df(v compact_float.DFloat) *tevent { return newTEvent(teventDecimalFloat, v, nil) }
func bdf(v *apd.Decimal) *tevent        { return eventOrNil(teventBigDecimalFloat, v) }
func v(v uint64) *tevent                { return newTEvent(teventVersion, v, nil) }
func n() *tevent                        { return newTEvent(teventNil, nil, nil) }
func pad(v int) *tevent                 { return newTEvent(teventPadding, v, nil) }
func b(v bool) *tevent                  { return newTEvent(teventBool, v, nil) }
func pi(v uint64) *tevent               { return newTEvent(teventPInt, v, nil) }
func ni(v uint64) *tevent               { return newTEvent(teventNInt, v, nil) }
func bi(v *big.Int) *tevent             { return eventOrNil(teventBigInt, v) }
func cplx(v complex128) *tevent         { return newTEvent(teventComplex, v, nil) }
func nan() *tevent                      { return newTEvent(teventNan, nil, nil) }
func snan() *tevent                     { return newTEvent(teventSNan, nil, nil) }
func uuid(v []byte) *tevent             { return newTEvent(teventUUID, v, nil) }
func gt(v time.Time) *tevent            { return newTEvent(teventTime, v, nil) }
func ct(v *compact_time.Time) *tevent   { return eventOrNil(teventCompactTime, v) }
func bin(v []byte) *tevent              { return newTEvent(teventBytes, v, nil) }
func s(v string) *tevent                { return newTEvent(teventString, v, nil) }
func uri(v string) *tevent              { return newTEvent(teventURI, v, nil) }
func cust(v []byte) *tevent             { return newTEvent(teventCustom, v, nil) }
func bb() *tevent                       { return newTEvent(teventBytesBegin, nil, nil) }
func sb() *tevent                       { return newTEvent(teventStringBegin, nil, nil) }
func ub() *tevent                       { return newTEvent(teventURIBegin, nil, nil) }
func cb() *tevent                       { return newTEvent(teventCustomBegin, nil, nil) }
func ac(l uint64, term bool) *tevent    { return newTEvent(teventArrayChunk, l, term) }
func ad(v []byte) *tevent               { return newTEvent(teventArrayData, v, nil) }
func l() *tevent                        { return newTEvent(teventList, nil, nil) }
func m() *tevent                        { return newTEvent(teventMap, nil, nil) }
func mup() *tevent                      { return newTEvent(teventMarkup, nil, nil) }
func meta() *tevent                     { return newTEvent(teventMetadata, nil, nil) }
func cmt() *tevent                      { return newTEvent(teventComment, nil, nil) }
func e() *tevent                        { return newTEvent(teventEnd, nil, nil) }
func mark() *tevent                     { return newTEvent(teventMarker, nil, nil) }
func ref() *tevent                      { return newTEvent(teventReference, nil, nil) }
func ed() *tevent                       { return newTEvent(teventEndDocument, nil, nil) }

func eventForValue(value interface{}) *tevent {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return n()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return b(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return i(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return pi(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return f(rv.Float())
	case reflect.Complex64, reflect.Complex128:
		panic(fmt.Errorf("TODO: %v", rv.Type()))
	case reflect.String:
		return s(rv.String())
	case reflect.Slice:
		switch rv.Type() {
		case typeBytes:
			return bin(rv.Bytes())
		}
	case reflect.Ptr:
		if rv.IsNil() {
			return n()
		}
		switch rv.Type() {
		case typePBigDecimalFloat:
			return bdf(rv.Interface().(*apd.Decimal))
		case typePBigFloat:
			return bf(rv.Interface().(*big.Float))
		case typePBigInt:
			return bi(rv.Interface().(*big.Int))
		case typePCompactTime:
			return ct(rv.Interface().(*compact_time.Time))
		case typePURL:
			return uri(rv.Interface().(*url.URL).String())
		}
		return eventForValue(rv.Elem().Interface())
	case reflect.Struct:
		switch rv.Type() {
		case typeBigDecimalFloat:
			v := rv.Interface().(apd.Decimal)
			return bdf(&v)
		case typeBigFloat:
			v := rv.Interface().(big.Float)
			return bf(&v)
		case typeBigInt:
			v := rv.Interface().(big.Int)
			return bi(&v)
		case typeCompactTime:
			v := rv.Interface().(compact_time.Time)
			return ct(&v)
		case typeDFloat:
			v := rv.Interface().(compact_float.DFloat)
			return df(v)
		case typeTime:
			v := rv.Interface().(time.Time)
			return gt(v)
		case typeURL:
			v := rv.Interface().(url.URL)
			return uri(v.String())
		}
	}
	panic(fmt.Errorf("Testing BUG: Unhandled kind: %v", rv.Kind()))
}

type TER struct {
	Events []*tevent
}

func NewTER() *TER {
	return &TER{}
}
func (h *TER) add(event *tevent) {
	h.Events = append(h.Events, event)
}
func (h *TER) OnVersion(version uint64)   { h.add(v(version)) }
func (h *TER) OnPadding(count int)        { h.add(pad(count)) }
func (h *TER) OnNil()                     { h.add(n()) }
func (h *TER) OnBool(value bool)          { h.add(b(value)) }
func (h *TER) OnTrue()                    { h.add(tt()) }
func (h *TER) OnFalse()                   { h.add(ff()) }
func (h *TER) OnPositiveInt(value uint64) { h.add(pi(value)) }
func (h *TER) OnNegativeInt(value uint64) { h.add(ni(value)) }
func (h *TER) OnInt(value int64)          { h.add(i(value)) }
func (h *TER) OnBigInt(value *big.Int)    { h.add(bi(value)) }
func (h *TER) OnFloat(value float64) {
	if math.IsNaN(value) {
		if isSignalingNan(value) {
			h.add(snan())
		} else {
			h.add(nan())
		}
	} else {
		h.add(f(value))
	}
}
func (h *TER) OnBigFloat(value *big.Float) { h.add(newTEvent(teventBigFloat, value, nil)) }
func (h *TER) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			h.add(snan())
		} else {
			h.add(nan())
		}
	} else {
		h.add(df(value))
	}
}
func (h *TER) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		h.add(nan())
	case apd.NaNSignaling:
		h.add(snan())
	default:
		h.add(bdf(value))
	}
}
func (h *TER) OnComplex(value complex128)             { h.add(cplx(value)) }
func (h *TER) OnUUID(value []byte)                    { h.add(uuid(value)) }
func (h *TER) OnTime(value time.Time)                 { h.add(gt(value)) }
func (h *TER) OnCompactTime(value *compact_time.Time) { h.add(ct(value)) }
func (h *TER) OnBytes(value []byte)                   { h.add(bin(value)) }
func (h *TER) OnString(value string)                  { h.add(s(value)) }
func (h *TER) OnURI(value string)                     { h.add(uri((value))) }
func (h *TER) OnCustom(value []byte)                  { h.add(cust(value)) }
func (h *TER) OnBytesBegin()                          { h.add(bb()) }
func (h *TER) OnStringBegin()                         { h.add(sb()) }
func (h *TER) OnURIBegin()                            { h.add(ub()) }
func (h *TER) OnCustomBegin()                         { h.add(cb()) }
func (h *TER) OnArrayChunk(l uint64, final bool)      { h.add(ac(l, final)) }
func (h *TER) OnArrayData(data []byte)                { h.add(ad(data)) }
func (h *TER) OnList()                                { h.add(l()) }
func (h *TER) OnMap()                                 { h.add(m()) }
func (h *TER) OnMarkup()                              { h.add(mup()) }
func (h *TER) OnMetadata()                            { h.add(meta()) }
func (h *TER) OnComment()                             { h.add(cmt()) }
func (h *TER) OnEnd()                                 { h.add(e()) }
func (h *TER) OnMarker()                              { h.add(mark()) }
func (h *TER) OnReference()                           { h.add(ref()) }
func (h *TER) OnEndDocument()                         { h.add(ed()) }
func (h *TER) OnNan(signaling bool) {
	if signaling {
		h.add(snan())
	} else {
		h.add(nan())
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

func (_this *TestingOuterStruct) getRepresentativeEvents(includeFakes bool) (events []*tevent) {
	ade := func(e ...*tevent) {
		events = append(events, e...)
	}
	adv := func(value interface{}) {
		ade(eventForValue(value))
	}
	anv := func(name string, value interface{}) {
		adv(name)
		adv(value)
	}
	ane := func(name string, e ...*tevent) {
		adv(name)
		ade(e...)
	}

	ade(m())

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
	ane("Ar", bin(_this.Ar[:]))
	anv("St", _this.St)
	anv("Ba", _this.Ba)

	ane("Sl", l())
	for _, v := range _this.Sl {
		adv(v)
	}
	ade(e())

	ane("M", m())
	for k, v := range _this.M {
		adv(k)
		adv(v)
	}
	ade(e())

	ane("IS", m())
	anv("Inner", _this.IS.Inner)
	ade(e())

	if _this.PIS != nil {
		ane("PIS", m())
		anv("Inner", _this.PIS.Inner)
		ade(e())
	}

	anv("Time", _this.Time)
	anv("PTime", _this.PTime)
	anv("CTime", _this.CTime)
	anv("PCTime", _this.PCTime)
	anv("PURL", _this.PURL)

	if includeFakes {
		ane("F1", b(true))
		ane("F2", b(false))
		ane("F3", i(1))
		ane("F4", i(-1))
		ane("F5", f(1.1))
		ane("F6", bf(newBigFloat("1.1", 2)))
		ane("F7", df(newDFloat("1.1")))
		ane("F8", bdf(newBDF("1.1")))
		ane("F9", n())
		ane("F10", bi(newBigInt("1000")))
		// ane("F11", cplx(1+1i))
		ane("F12", nan())
		ane("F13", snan())
		ane("F14", uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
		ane("F15", gt(_this.Time))
		ane("F16", ct(_this.PCTime))
		ane("F17", bin([]byte{1}))
		ane("F18", s("xyz"))
		ane("F19", uri("http://example.com"))
		// ane("F20", cust([]byte{1}))
		ane("FakeList", l(), e())
		ane("FakeMap", m(), e())
		ane("FakeDeep", l(), m(), s("A"), l(),
			b(true),
			b(false),
			i(1),
			i(-1),
			f(1.1),
			bf(newBigFloat("1.1", 2)),
			df(newDFloat("1.1")),
			bdf(newBDF("1.1")),
			n(),
			bi(newBigInt("1000")),
			// cplx(1+1i),
			nan(),
			snan(),
			uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
			gt(_this.Time),
			ct(_this.PCTime),
			bin([]byte{1}),
			s("xyz"),
			uri("http://example.com"),
			// cust([]byte{1}),
			e(), e(), e())
	}

	ade(e())
	return
}

func newTestingOuterStruct(baseValue int) *TestingOuterStruct {
	_this := new(TestingOuterStruct)
	_this.Init(baseValue)
	return _this
}

func newBlankTestingOuterStruct() *TestingOuterStruct {
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
	_this.PBI = newBigInt(fmt.Sprintf("-10000000000000000000000000000000000000%v", unsafe.Offsetof(_this.PBI)))
	_this.BI = *_this.PBI
	_this.F32 = float32(1000000+baseValue+int(unsafe.Offsetof(_this.F32))) + 0.5
	_this.PF32 = &_this.F32
	_this.F64 = float64(1000000000000+baseValue+int(unsafe.Offsetof(_this.F64))) + 0.5
	_this.PF64 = &_this.F64
	_this.PBF = newBigFloat("12345678901234567890123.1234567", 30)
	_this.BF = *_this.PBF
	_this.DF = newDFloat(fmt.Sprintf("-100000000000000%ve-1000000", unsafe.Offsetof(_this.DF)))
	_this.PBDF = newBDF("-1.234567890123456789777777777777777777771234e-10000")
	_this.BDF = *_this.PBDF
	_this.Ar[0] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)))
	_this.Ar[1] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+1))
	_this.Ar[2] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+2))
	_this.Ar[3] = byte(baseValue + int(unsafe.Offsetof(_this.Ar)+3))
	_this.St = generateString(baseValue+5, baseValue)
	_this.Ba = generateBytes(baseValue+1, baseValue)
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
