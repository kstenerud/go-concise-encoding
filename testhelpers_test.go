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
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

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
	shouldZeroCopy := false
	err = CBEDecode(document, receiver, shouldZeroCopy)
	events = receiver.Events
	return
}

func cbeEncodeDecode(expected ...*tevent) (events []*tevent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	encoder := NewCBEEncoder()
	invokeEvents(encoder, expected...)
	document := encoder.Document()

	return cbeDecode(document)
}

func cteDecode(document []byte) (events []*tevent, err error) {
	receiver := NewTER()
	err = CTEDecode(document, receiver)
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
	"cmpl",
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
	if math.IsNaN(v) {
		if isSignalingNan(v) {
			return snan()
		}
		return nan()
	}
	return newTEvent(teventFloat, v, nil)
}
func bf(v string, significantDigits int) *tevent {
	bf, err := stringToBigFloat(v, significantDigits)
	if err != nil {
		panic(err)
	}
	return newTEvent(teventBigDecimalFloat, bf, nil)
}

func df(v string) *tevent {
	decimal, err := compact_float.DFloatFromString(v)
	if err != nil {
		panic(err)
	}
	return newTEvent(teventDecimalFloat, decimal, nil)
}

func bdf(v string) *tevent {
	decimal, _, err := apd.NewFromString(v)
	if err != nil {
		panic(err)
	}
	return newTEvent(teventBigDecimalFloat, decimal, nil)
}

func v(v uint64) *tevent              { return newTEvent(teventVersion, v, nil) }
func n() *tevent                      { return newTEvent(teventNil, nil, nil) }
func pad(v int) *tevent               { return newTEvent(teventPadding, v, nil) }
func b(v bool) *tevent                { return newTEvent(teventBool, v, nil) }
func pi(v uint64) *tevent             { return newTEvent(teventPInt, v, nil) }
func ni(v uint64) *tevent             { return newTEvent(teventNInt, v, nil) }
func bi(v *big.Int) *tevent           { return newTEvent(teventBigInt, v, nil) }
func cmpl(v complex128) *tevent       { return newTEvent(teventComplex, v, nil) }
func nan() *tevent                    { return newTEvent(teventNan, nil, nil) }
func snan() *tevent                   { return newTEvent(teventSNan, nil, nil) }
func uuid(v []byte) *tevent           { return newTEvent(teventUUID, v, nil) }
func gt(v time.Time) *tevent          { return newTEvent(teventTime, v, nil) }
func ct(v *compact_time.Time) *tevent { return newTEvent(teventCompactTime, v, nil) }
func bin(v []byte) *tevent            { return newTEvent(teventBytes, v, nil) }
func s(v string) *tevent              { return newTEvent(teventString, v, nil) }
func uri(v string) *tevent            { return newTEvent(teventURI, v, nil) }
func cust(v []byte) *tevent           { return newTEvent(teventCustom, v, nil) }
func bb() *tevent                     { return newTEvent(teventBytesBegin, nil, nil) }
func sb() *tevent                     { return newTEvent(teventStringBegin, nil, nil) }
func ub() *tevent                     { return newTEvent(teventURIBegin, nil, nil) }
func cb() *tevent                     { return newTEvent(teventCustomBegin, nil, nil) }
func ac(l uint64, term bool) *tevent  { return newTEvent(teventArrayChunk, l, term) }
func ad(v []byte) *tevent             { return newTEvent(teventArrayData, v, nil) }
func l() *tevent                      { return newTEvent(teventList, nil, nil) }
func m() *tevent                      { return newTEvent(teventMap, nil, nil) }
func mup() *tevent                    { return newTEvent(teventMarkup, nil, nil) }
func meta() *tevent                   { return newTEvent(teventMetadata, nil, nil) }
func cmt() *tevent                    { return newTEvent(teventComment, nil, nil) }
func e() *tevent                      { return newTEvent(teventEnd, nil, nil) }
func mark() *tevent                   { return newTEvent(teventMarker, nil, nil) }
func ref() *tevent                    { return newTEvent(teventReference, nil, nil) }
func ed() *tevent                     { return newTEvent(teventEndDocument, nil, nil) }

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
		h.add(df(value.String()))
	}
}
func (h *TER) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		h.add(nan())
	case apd.NaNSignaling:
		h.add(snan())
	default:
		h.add(bdf(value.String()))
	}
}
func (h *TER) OnComplex(value complex128)             { h.add(cmpl(value)) }
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
