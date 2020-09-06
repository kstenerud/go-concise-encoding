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

package iterator

import (
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-equivalence"
)

func TestIterateBasic(t *testing.T) {
	pBigIntP := NewBigInt("12345678901234567890123456789", 10)
	pBigIntN := NewBigInt("-999999999999999999999999999999", 10)
	pBigFloat := NewBigFloat("5.377345e-10000", 10, 7)
	dfloat := NewDFloat("1.23456e1000")
	pBigDFloat := NewBDF("4.509e10000")
	gTimeNow := time.Now()
	pCTimeNow := compact_time.AsCompactTime(gTimeNow)
	pURL := NewURI("http://x.com")

	assertIterate(t, nil, N())
	assertIterate(t, true, B(true))
	assertIterate(t, false, B(false))
	assertIterate(t, int(10), I(10))
	assertIterate(t, int8(10), I(10))
	assertIterate(t, int16(10), I(10))
	assertIterate(t, int32(10), I(10))
	assertIterate(t, int64(10), I(10))
	assertIterate(t, uint(10), I(10))
	assertIterate(t, uint8(10), I(10))
	assertIterate(t, uint16(10), I(10))
	assertIterate(t, uint32(10), I(10))
	assertIterate(t, uint64(10), I(10))
	assertIterate(t, 1, I(1))
	assertIterate(t, -1, I(-1))
	assertIterate(t, pBigIntP, BI(pBigIntP))
	assertIterate(t, *pBigIntP, BI(pBigIntP))
	assertIterate(t, pBigIntN, BI(pBigIntN))
	assertIterate(t, *pBigIntN, BI(pBigIntN))
	assertIterate(t, (*big.Int)(nil), N())
	assertIterate(t, float32(-1.25), F(-1.25))
	assertIterate(t, float64(-9.5e50), F(-9.5e50))
	assertIterate(t, pBigFloat, BF(pBigFloat))
	assertIterate(t, *pBigFloat, BF(pBigFloat))
	assertIterate(t, (*big.Float)(nil), N())
	assertIterate(t, dfloat, DF(dfloat))
	assertIterate(t, pBigDFloat, BDF(pBigDFloat))
	assertIterate(t, *pBigDFloat, BDF(pBigDFloat))
	assertIterate(t, (*apd.Decimal)(nil), N())
	assertIterate(t, common.SignalingNan, SNAN())
	assertIterate(t, common.QuietNan, NAN())
	assertIterate(t, gTimeNow, GT(gTimeNow))
	assertIterate(t, pCTimeNow, CT(pCTimeNow))
	assertIterate(t, *pCTimeNow, CT(pCTimeNow))
	assertIterate(t, []byte{1, 2, 3, 4}, AU8([]byte{1, 2, 3, 4}))
	assertIterate(t, "test", S("test"))
	assertIterate(t, pURL, URI("http://x.com"))
	assertIterate(t, *pURL, URI("http://x.com"))
	assertIterate(t, (*url.URL)(nil), N())
}

func TestIterateSlice(t *testing.T) {
	assertIterate(t, []int{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []int8{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []int16{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []int32{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []int64{1, 2}, L(), I(1), I(2), E())

	assertIterate(t, []uint{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []uint8{0x12, 0x34}, AU8([]byte{0x12, 0x34}))
	assertIterate(t, []uint16{0x1234, 0x5678}, AU16([]uint16{0x1234, 0x5678}))
	assertIterate(t, []uint32{1, 2}, L(), I(1), I(2), E())
	assertIterate(t, []uint64{1, 2}, L(), I(1), I(2), E())

	assertIterate(t, []float32{1, 2}, L(), F(1), F(2), E())
	assertIterate(t, []float64{1, 2}, L(), F(1), F(2), E())
	assertIterate(t, []float64(nil), N())
}

func TestIterateArray(t *testing.T) {
	a1 := [2]int{1, 2}
	assertIterate(t, a1, L(), I(1), I(2), E())
	assertIterate(t, &a1, L(), I(1), I(2), E())
	a2 := [2]byte{1, 2}
	assertIterate(t, a2, AU8([]byte{1, 2}))
	assertIterate(t, &a2, AU8([]byte{1, 2}))
}

func TestIterateInterface(t *testing.T) {
	assertIterate(t, []interface{}{1, nil, 5.5}, L(), I(1), N(), F(5.5), E())
}

func TestIteratePointer(t *testing.T) {
	v := 1
	assertIterate(t, &v, I(1))
	pv := (*int)(nil)
	assertIterate(t, pv, N())
}

func TestIterateMap(t *testing.T) {
	assertIterate(t, map[string]int{"a": 1}, M(), S("a"), I(1), E())
	assertIterate(t, (map[string]int)(nil), N())
}

type StructTestIterate struct {
	A int
}

func TestIterateStruct(t *testing.T) {
	sOpts := options.DefaultIteratorSessionOptions()
	iOpts := options.DefaultIteratorOptions()
	iOpts.LowercaseStructFieldNames = false

	assertIterate(t, new(StructTestIterate), M(), S("a"), I(0), E())
	assertIterateWithOptions(t, sOpts, iOpts, new(StructTestIterate), M(), S("A"), I(0), E())
	assertIterate(t, (*StructTestIterate)(nil), N())
}

func TestIterateNilOpts(t *testing.T) {
	expected := []*test.TEvent{BD(), V(1), I(1), ED()}
	receiver := test.NewTER()
	iterateObject(1, receiver, nil, nil)

	if !equivalence.IsEquivalent(expected, receiver.Events) {
		t.Errorf("Expected %v but got %v", expected, receiver.Events)
	}
}

type RecursiveStructTestIterate struct {
	I int
	R *RecursiveStructTestIterate
}

func TestIterateRecurse(t *testing.T) {
	obj := &RecursiveStructTestIterate{
		I: 50,
	}
	obj.R = obj

	expected := []*test.TEvent{BD(), V(1), MARK(), I(0), M(), S("i"), I(50), S("r"), REF(), I(0), E(), ED()}
	sessionOptions := options.DefaultIteratorSessionOptions()
	iteratorOptions := options.DefaultIteratorOptions()
	iteratorOptions.RecursionSupport = true
	receiver := test.NewTER()
	iterateObject(obj, receiver, sessionOptions, iteratorOptions)

	if !equivalence.IsEquivalent(expected, receiver.Events) {
		t.Errorf("Expected %v but got %v", expected, receiver.Events)
	}
}

type TagStruct struct {
	Omit1 string `ce:"-"`
	Omit2 string `ce:"omit"`
	Named string `ce:"name=test"`
}

type TagStruct2 struct {
	Omit1 string `ce:" - "`
	Omit2 string `ce:" omit "`
	Named string `ce:" name = test "`
}

func TestIterateTaggedStruct(t *testing.T) {
	obj := &TagStruct{
		Omit1: "Omit1 should be omitted",
		Omit2: "Omit2 should be omitted",
		Named: "Named should be present",
	}

	assertIterate(t, obj, M(), S("test"), S("Named should be present"), E())

	obj2 := &TagStruct2{
		Omit1: "Omit1 should be omitted",
		Omit2: "Omit2 should be omitted",
		Named: "Named should be present",
	}

	assertIterate(t, obj2, M(), S("test"), S("Named should be present"), E())
}
