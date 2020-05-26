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
	"net/url"
	"testing"
	"time"
	"unsafe"

	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func assertCBEMarshalUnmarshal(t *testing.T, expected interface{}) {
	options := &CBEMarshalerOptions{
		IteratorOptions: IteratorOptions{
			UseReferences: true,
		},
	}
	document, err := MarshalCBE(expected, options)
	if err != nil {
		t.Error(err)
		return
	}

	var actual interface{}
	actual, err = UnmarshalCBE(document, expected, false)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertCTEMarshalUnmarshal(t *testing.T, expected interface{}) {
	options := &CTEMarshalerOptions{
		IteratorOptions: IteratorOptions{
			UseReferences: true,
		},
	}
	document, err := MarshalCTE(expected, options)
	if err != nil {
		t.Error(err)
		return
	}

	var actual interface{}
	actual, err = UnmarshalCTE(document, expected)
	if err != nil {
		fmt.Printf("While unmarshaling %v\n", string(document))
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}) {
	assertCBEMarshalUnmarshal(t, expected)
	assertCTEMarshalUnmarshal(t, expected)
}

type MarshalInnerStruct struct {
	Inner int
}

type MarshalTester struct {
	Bo   bool
	By   byte
	I    int
	I8   int8
	I16  int16
	I32  int32
	I64  int64
	U    uint
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	F32  float32
	F64  float64
	DF   compact_float.DFloat
	Ar   [4]byte
	St   string
	Ba   []byte
	Sl   []interface{}
	M    map[interface{}]interface{}
	Pi   *int
	IS   MarshalInnerStruct
	ISP  *MarshalInnerStruct
	Time time.Time
	// TODO: URL must have at least 2 chars... what to do here?
	// URL   url.URL
	PTime *time.Time
	PURL  *url.URL
}

func newMarshalTestStruct(baseValue int) *MarshalTester {
	this := new(MarshalTester)
	this.Init(baseValue)
	return this
}

func (this *MarshalTester) Init(baseValue int) {
	this.Bo = baseValue&1 == 1
	this.By = byte(baseValue + int(unsafe.Offsetof(this.By)))
	this.I = baseValue + int(unsafe.Offsetof(this.I))
	this.I8 = int8(baseValue + int(unsafe.Offsetof(this.I8)))
	this.I16 = int16(baseValue + int(unsafe.Offsetof(this.I16)))
	this.I32 = int32(baseValue + int(unsafe.Offsetof(this.I32)))
	this.I64 = int64(baseValue + int(unsafe.Offsetof(this.I64)))
	this.U = uint(baseValue + int(unsafe.Offsetof(this.U)))
	this.U8 = uint8(baseValue + int(unsafe.Offsetof(this.U8)))
	this.U16 = uint16(baseValue + int(unsafe.Offsetof(this.U16)))
	this.U32 = uint32(baseValue + int(unsafe.Offsetof(this.U32)))
	this.U64 = uint64(baseValue + int(unsafe.Offsetof(this.U64)))
	this.F32 = float32(baseValue+int(unsafe.Offsetof(this.F32))) + 0.5
	this.F64 = float64(baseValue+int(unsafe.Offsetof(this.F64))) + 0.5
	this.DF = compact_float.DFloat{
		Exponent:    -int32(baseValue),
		Coefficient: int64(baseValue + int(unsafe.Offsetof(this.DF))),
	}
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
	v := baseValue
	this.Pi = &v
	this.IS.Inner = baseValue + 15
	this.ISP = new(MarshalInnerStruct)
	this.ISP.Inner = baseValue + 16

	testTime := time.Date(2000+baseValue, time.Month(1), 1, 1, 1, 1, 0, time.UTC)
	this.PTime = &testTime
	this.PURL, _ = url.Parse(fmt.Sprintf("http://example.com/%v", baseValue))
}

func TestMarshalUnmarshal(t *testing.T) {
	assertMarshalUnmarshal(t, 101)
	assertMarshalUnmarshal(t, *newMarshalTestStruct(1))
	assertMarshalUnmarshal(t, *new(MarshalTester))
}
