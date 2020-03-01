package concise_encoding

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"testing"
	"time"

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
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
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

func invokeEvents(handler ConciseEncodingEventHandler, events ...*tevent) {
	for _, event := range events {
		event.Invoke(handler)
	}
}

func cbeDecode(document []byte) (events []*tevent, err error) {
	handler := NewTEH()
	shouldZeroCopy := false
	err = CBEDecode(document, handler, shouldZeroCopy)
	events = handler.Events
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
	handler := NewTEH()
	err = CTEDecode(document, handler)
	events = handler.Events
	return
}

func cteEncode(events ...*tevent) []byte {
	encoder := NewCTEEncoder()
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
	teventFloat
	teventComplex
	teventNan
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
	"f",
	"cmpl",
	"nan",
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

func (this *tevent) Invoke(handler ConciseEncodingEventHandler) {
	switch this.Type {
	case teventVersion:
		handler.OnVersion(this.V1.(uint64))
	case teventPadding:
		handler.OnPadding(this.V1.(int))
	case teventNil:
		handler.OnNil()
	case teventBool:
		handler.OnBool(this.V1.(bool))
	case teventTrue:
		handler.OnBool(true)
	case teventFalse:
		handler.OnBool(false)
	case teventPInt:
		handler.OnPositiveInt(this.V1.(uint64))
	case teventNInt:
		handler.OnNegativeInt(this.V1.(uint64))
	case teventInt:
		v := this.V1.(int64)
		if v >= 0 {
			handler.OnPositiveInt(uint64(v))
		} else {
			handler.OnNegativeInt(uint64(-v))
		}
	case teventFloat:
		handler.OnFloat(this.V1.(float64))
	case teventComplex:
		handler.OnComplex(this.V1.(complex128))
	case teventNan:
		handler.OnNan()
	case teventUUID:
		handler.OnUUID(this.V1.([]byte))
	case teventTime:
		handler.OnTime(this.V1.(time.Time))
	case teventCompactTime:
		handler.OnCompactTime(this.V1.(*compact_time.Time))
	case teventBytes:
		handler.OnBytes(this.V1.([]byte))
	case teventString:
		handler.OnString(this.V1.(string))
	case teventURI:
		handler.OnURI(this.V1.(string))
	case teventCustom:
		handler.OnCustom(this.V1.([]byte))
	case teventBytesBegin:
		handler.OnBytesBegin()
	case teventStringBegin:
		handler.OnStringBegin()
	case teventURIBegin:
		handler.OnURIBegin()
	case teventCustomBegin:
		handler.OnCustomBegin()
	case teventArrayChunk:
		handler.OnArrayChunk(this.V1.(uint64), this.V2.(bool))
	case teventArrayData:
		handler.OnArrayData(this.V1.([]byte))
	case teventList:
		handler.OnList()
	case teventMap:
		handler.OnMap()
	case teventMarkup:
		handler.OnMarkup()
	case teventMetadata:
		handler.OnMetadata()
	case teventComment:
		handler.OnComment()
	case teventEnd:
		handler.OnEnd()
	case teventMarker:
		handler.OnMarker()
	case teventReference:
		handler.OnReference()
	case teventEndDocument:
		handler.OnEndDocument()
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
		return nan()
	}
	return newTEvent(teventFloat, v, nil)
}

func v(v uint64) *tevent              { return newTEvent(teventVersion, v, nil) }
func n() *tevent                      { return newTEvent(teventNil, nil, nil) }
func pad(v int) *tevent               { return newTEvent(teventPadding, v, nil) }
func b(v bool) *tevent                { return newTEvent(teventBool, v, nil) }
func pi(v uint64) *tevent             { return newTEvent(teventPInt, v, nil) }
func ni(v uint64) *tevent             { return newTEvent(teventNInt, v, nil) }
func cmpl(v complex128) *tevent       { return newTEvent(teventComplex, v, nil) }
func nan() *tevent                    { return newTEvent(teventNan, nil, nil) }
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

type TEH struct {
	Events []*tevent
}

func NewTEH() *TEH {
	return &TEH{}
}
func (h *TEH) add(event *tevent) {
	h.Events = append(h.Events, event)
}
func (h *TEH) OnVersion(version uint64)               { h.add(v(version)) }
func (h *TEH) OnPadding(count int)                    { h.add(pad(count)) }
func (h *TEH) OnNil()                                 { h.add(n()) }
func (h *TEH) OnBool(value bool)                      { h.add(b(value)) }
func (h *TEH) OnTrue()                                { h.add(tt()) }
func (h *TEH) OnFalse()                               { h.add(ff()) }
func (h *TEH) OnPositiveInt(value uint64)             { h.add(pi(value)) }
func (h *TEH) OnNegativeInt(value uint64)             { h.add(ni(value)) }
func (h *TEH) OnInt(value int64)                      { h.add(i(value)) }
func (h *TEH) OnFloat(value float64)                  { h.add(f(value)) }
func (h *TEH) OnComplex(value complex128)             { h.add(cmpl(value)) }
func (h *TEH) OnNan()                                 { h.add(nan()) }
func (h *TEH) OnUUID(value []byte)                    { h.add(uuid(value)) }
func (h *TEH) OnTime(value time.Time)                 { h.add(gt(value)) }
func (h *TEH) OnCompactTime(value *compact_time.Time) { h.add(ct(value)) }
func (h *TEH) OnBytes(value []byte)                   { h.add(bin(value)) }
func (h *TEH) OnString(value string)                  { h.add(s(value)) }
func (h *TEH) OnURI(value string)                     { h.add(uri((value))) }
func (h *TEH) OnCustom(value []byte)                  { h.add(cust(value)) }
func (h *TEH) OnBytesBegin()                          { h.add(bb()) }
func (h *TEH) OnStringBegin()                         { h.add(sb()) }
func (h *TEH) OnURIBegin()                            { h.add(ub()) }
func (h *TEH) OnCustomBegin()                         { h.add(cb()) }
func (h *TEH) OnArrayChunk(l uint64, final bool)      { h.add(ac(l, final)) }
func (h *TEH) OnArrayData(data []byte)                { h.add(ad(data)) }
func (h *TEH) OnList()                                { h.add(l()) }
func (h *TEH) OnMap()                                 { h.add(m()) }
func (h *TEH) OnMarkup()                              { h.add(mup()) }
func (h *TEH) OnMetadata()                            { h.add(meta()) }
func (h *TEH) OnComment()                             { h.add(cmt()) }
func (h *TEH) OnEnd()                                 { h.add(e()) }
func (h *TEH) OnMarker()                              { h.add(mark()) }
func (h *TEH) OnReference()                           { h.add(ref()) }
func (h *TEH) OnEndDocument()                         { h.add(ed()) }
