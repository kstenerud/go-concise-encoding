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
	document := encoder.Document()

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
	return encoder.Document()
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

func (this teventType) String() string {
	return teventNames[this]
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

func (this *tevent) String() string {
	str := this.Type.String()
	if this.V1 != nil {
		if this.V2 != nil {
			return fmt.Sprintf("%v(%v,%v)", str, this.V1, this.V2)
		}
		return fmt.Sprintf("%v(%v)", str, this.V1)
	}
	return str
}

func (this *tevent) Invoke(receiver DataEventReceiver) {
	switch this.Type {
	case teventVersion:
		receiver.OnVersion(this.V1.(uint64))
	case teventPadding:
		receiver.OnPadding(this.V1.(int))
	case teventNil:
		receiver.OnNil()
	case teventBool:
		receiver.OnBool(this.V1.(bool))
	case teventTrue:
		receiver.OnBool(true)
	case teventFalse:
		receiver.OnBool(false)
	case teventPInt:
		receiver.OnPositiveInt(this.V1.(uint64))
	case teventNInt:
		receiver.OnNegativeInt(this.V1.(uint64))
	case teventInt:
		v := this.V1.(int64)
		if v >= 0 {
			receiver.OnPositiveInt(uint64(v))
		} else {
			receiver.OnNegativeInt(uint64(-v))
		}
	case teventBigInt:
		receiver.OnBigInt(this.V1.(*big.Int))
	case teventFloat:
		receiver.OnFloat(this.V1.(float64))
	case teventBigFloat:
		receiver.OnBigFloat(this.V1.(*big.Float))
	case teventDecimalFloat:
		receiver.OnDecimalFloat(this.V1.(compact_float.DFloat))
	case teventBigDecimalFloat:
		receiver.OnBigDecimalFloat(this.V1.(*apd.Decimal))
	case teventComplex:
		receiver.OnComplex(this.V1.(complex128))
	case teventNan:
		receiver.OnNan(false)
	case teventSNan:
		receiver.OnNan(true)
	case teventUUID:
		receiver.OnUUID(this.V1.([]byte))
	case teventTime:
		receiver.OnTime(this.V1.(time.Time))
	case teventCompactTime:
		receiver.OnCompactTime(this.V1.(*compact_time.Time))
	case teventBytes:
		receiver.OnBytes(this.V1.([]byte))
	case teventString:
		receiver.OnString(this.V1.(string))
	case teventURI:
		receiver.OnURI(this.V1.(string))
	case teventCustom:
		receiver.OnCustom(this.V1.([]byte))
	case teventBytesBegin:
		receiver.OnBytesBegin()
	case teventStringBegin:
		receiver.OnStringBegin()
	case teventURIBegin:
		receiver.OnURIBegin()
	case teventCustomBegin:
		receiver.OnCustomBegin()
	case teventArrayChunk:
		receiver.OnArrayChunk(this.V1.(uint64), this.V2.(bool))
	case teventArrayData:
		receiver.OnArrayData(this.V1.([]byte))
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
		panic(fmt.Errorf("%v: Unhandled event type", this.Type))
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
	} else {
		return ni(uint64(-v))
	}
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
	// TODO: URL must have at least 2 chars... what to do here?
	// URL   url.URL
}

func (this *TestingOuterStruct) Events(includeFakes bool) (events []*tevent) {
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

	anv("Bo", this.Bo)
	anv("PBo", this.PBo)
	anv("By", this.By)
	anv("PBy", this.PBy)
	anv("I", this.I)
	anv("PI", this.PI)
	anv("I8", this.I8)
	anv("PI8", this.PI8)
	anv("I16", this.I16)
	anv("PI16", this.PI16)
	anv("I32", this.I32)
	anv("PI32", this.PI32)
	anv("I64", this.I64)
	anv("PI64", this.PI64)
	anv("U", this.U)
	anv("PU", this.PU)
	anv("U8", this.U8)
	anv("PU8", this.PU8)
	anv("U16", this.U16)
	anv("PU16", this.PU16)
	anv("U32", this.U32)
	anv("PU32", this.PU32)
	anv("U64", this.U64)
	anv("PU64", this.PU64)
	anv("BI", this.BI)
	anv("PBI", this.PBI)
	anv("F32", this.F32)
	anv("PF32", this.PF32)
	anv("F64", this.F64)
	anv("PF64", this.PF64)
	anv("BF", this.BF)
	anv("PBF", this.PBF)
	anv("DF", this.DF)
	anv("BDF", this.BDF)
	anv("PBDF", this.PBDF)
	ane("Ar", bin(this.Ar[:]))
	anv("St", this.St)
	anv("Ba", this.Ba)

	ane("Sl", l())
	for _, v := range this.Sl {
		adv(v)
	}
	ade(e())

	ane("M", m())
	for k, v := range this.M {
		adv(k)
		adv(v)
	}
	ade(e())

	ane("IS", m())
	anv("Inner", this.IS.Inner)
	ade(e())

	if this.PIS != nil {
		ane("PIS", m())
		anv("Inner", this.PIS.Inner)
		ade(e())
	}

	anv("Time", this.Time)
	anv("PTime", this.PTime)
	anv("CTime", this.CTime)
	anv("PCTime", this.PCTime)
	anv("PURL", this.PURL)

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
		ane("F15", gt(this.Time))
		ane("F16", ct(this.PCTime))
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
			gt(this.Time),
			ct(this.PCTime),
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
	this := new(TestingOuterStruct)
	this.Init(baseValue)
	return this
}

func newBlankTestingOuterStruct() *TestingOuterStruct {
	this := new(TestingOuterStruct)
	this.CTime.Year = 1
	this.CTime.Month = 1
	this.CTime.Day = 1
	return this
}

func (this *TestingOuterStruct) Init(baseValue int) {
	this.Bo = baseValue&1 == 1
	this.PBo = &this.Bo
	this.By = byte(baseValue + int(unsafe.Offsetof(this.By)))
	this.PBy = &this.By
	this.I = 100000 + baseValue + int(unsafe.Offsetof(this.I))
	this.PI = &this.I
	this.I8 = int8(baseValue + int(unsafe.Offsetof(this.I8)))
	this.PI8 = &this.I8
	this.I16 = int16(-1000 - baseValue - int(unsafe.Offsetof(this.I16)))
	this.PI16 = &this.I16
	this.I32 = int32(1000000000 + baseValue + int(unsafe.Offsetof(this.I32)))
	this.PI32 = &this.I32
	this.I64 = int64(1000000000000000 + baseValue + int(unsafe.Offsetof(this.I64)))
	this.PI64 = &this.I64
	this.U = uint(1000000 + baseValue + int(unsafe.Offsetof(this.U)))
	this.PU = &this.U
	this.U8 = uint8(baseValue + int(unsafe.Offsetof(this.U8)))
	this.PU8 = &this.U8
	this.U16 = uint16(10000 + baseValue + int(unsafe.Offsetof(this.U16)))
	this.PU16 = &this.U16
	this.U32 = uint32(100000000 + baseValue + int(unsafe.Offsetof(this.U32)))
	this.PU32 = &this.U32
	this.U64 = uint64(1000000000000 + baseValue + int(unsafe.Offsetof(this.U64)))
	this.PU64 = &this.U64
	this.PBI = newBigInt(fmt.Sprintf("-10000000000000000000000000000000000000%v", unsafe.Offsetof(this.PBI)))
	this.BI = *this.PBI
	this.F32 = float32(1000000+baseValue+int(unsafe.Offsetof(this.F32))) + 0.5
	this.PF32 = &this.F32
	this.F64 = float64(1000000000000+baseValue+int(unsafe.Offsetof(this.F64))) + 0.5
	this.PF64 = &this.F64
	this.PBF = newBigFloat("12345678901234567890123.1234567", 30)
	this.BF = *this.PBF
	this.DF = newDFloat(fmt.Sprintf("-100000000000000%ve-1000000", unsafe.Offsetof(this.DF)))
	this.PBDF = newBDF("-1.234567890123456789777777777777777777771234e-10000")
	this.BDF = *this.PBDF
	this.Ar[0] = byte(baseValue + int(unsafe.Offsetof(this.Ar)))
	this.Ar[1] = byte(baseValue + int(unsafe.Offsetof(this.Ar)+1))
	this.Ar[2] = byte(baseValue + int(unsafe.Offsetof(this.Ar)+2))
	this.Ar[3] = byte(baseValue + int(unsafe.Offsetof(this.Ar)+3))
	this.St = generateString(baseValue+5, baseValue)
	this.Ba = generateBytes(baseValue+1, baseValue)
	this.M = make(map[interface{}]interface{})
	for i := 0; i < baseValue+2; i++ {
		this.Sl = append(this.Sl, i)
		this.M[fmt.Sprintf("key%v", i)] = i
	}
	this.IS.Inner = baseValue + 15
	this.PIS = new(TestingInnerStruct)
	this.PIS.Inner = baseValue + 16
	testTime := time.Date(30000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	this.PTime = &testTime
	this.PCTime = compact_time.NewTimestamp(-1000, 1, 1, 1, 1, 1, 1, "Europe/Berlin")
	this.CTime = *this.PCTime
	this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}
