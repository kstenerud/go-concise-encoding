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
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-equivalence"
)

func assertIterate(t *testing.T, obj interface{}, events ...*tevent) {
	expected := append([]*tevent{v(1)}, events...)
	expected = append(expected, ed())
	options := IteratorOptions{}
	receiver := NewTER()
	IterateObject(obj, receiver, &options)

	if !equivalence.IsEquivalent(expected, receiver.Events) {
		t.Errorf("Expected %v but got %v", expected, receiver.Events)
	}
}

// ============================================================================

func TestIterateBasic(t *testing.T) {
	pBigIntP := newBigInt("12345678901234567890123456789")
	pBigIntN := newBigInt("-999999999999999999999999999999")
	pBigFloat := newBigFloat("5.377345e-10000", 7)
	dfloat := newDFloat("1.23456e1000")
	pBigDFloat := newBDF("4.509e10000")
	gTimeNow := time.Now()
	pCTimeNow := compact_time.AsCompactTime(gTimeNow)
	pURL := newURL("http://x.com")

	assertIterate(t, nil, n())
	assertIterate(t, true, b(true))
	assertIterate(t, false, b(false))
	assertIterate(t, int(10), i(10))
	assertIterate(t, int8(10), i(10))
	assertIterate(t, int16(10), i(10))
	assertIterate(t, int32(10), i(10))
	assertIterate(t, int64(10), i(10))
	assertIterate(t, uint(10), i(10))
	assertIterate(t, uint8(10), i(10))
	assertIterate(t, uint16(10), i(10))
	assertIterate(t, uint32(10), i(10))
	assertIterate(t, uint64(10), i(10))
	assertIterate(t, 1, i(1))
	assertIterate(t, -1, i(-1))
	assertIterate(t, pBigIntP, bi(pBigIntP))
	assertIterate(t, *pBigIntP, bi(pBigIntP))
	assertIterate(t, pBigIntN, bi(pBigIntN))
	assertIterate(t, *pBigIntN, bi(pBigIntN))
	assertIterate(t, (*big.Int)(nil), n())
	assertIterate(t, float32(-1.25), f(-1.25))
	assertIterate(t, float64(-9.5e50), f(-9.5e50))
	assertIterate(t, pBigFloat, bf(pBigFloat))
	assertIterate(t, *pBigFloat, bf(pBigFloat))
	assertIterate(t, (*big.Float)(nil), n())
	assertIterate(t, dfloat, df(dfloat))
	assertIterate(t, pBigDFloat, bdf(pBigDFloat))
	assertIterate(t, *pBigDFloat, bdf(pBigDFloat))
	assertIterate(t, (*apd.Decimal)(nil), n())
	assertIterate(t, complex64(-1+5i), cplx(-1+5i))
	assertIterate(t, complex128(-1+5i), cplx(-1+5i))
	assertIterate(t, signalingNan, snan())
	assertIterate(t, quietNan, nan())
	assertIterate(t, gTimeNow, gt(gTimeNow))
	assertIterate(t, pCTimeNow, ct(pCTimeNow))
	assertIterate(t, *pCTimeNow, ct(pCTimeNow))
	assertIterate(t, []byte{1, 2, 3, 4}, bin([]byte{1, 2, 3, 4}))
	assertIterate(t, "test", s("test"))
	assertIterate(t, pURL, uri("http://x.com"))
	assertIterate(t, *pURL, uri("http://x.com"))
	assertIterate(t, (*url.URL)(nil), n())
}

func TestIterateSlice(t *testing.T) {
	assertIterate(t, []int{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []int8{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []int16{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []int32{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []int64{1, 2}, l(), i(1), i(2), e())

	assertIterate(t, []uint{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []uint16{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []uint32{1, 2}, l(), i(1), i(2), e())
	assertIterate(t, []uint64{1, 2}, l(), i(1), i(2), e())

	assertIterate(t, []float32{1, 2}, l(), f(1), f(2), e())
	assertIterate(t, []float64{1, 2}, l(), f(1), f(2), e())
	assertIterate(t, []float64(nil), n())
}

func TestIterateArray(t *testing.T) {
	a1 := [2]int{1, 2}
	assertIterate(t, a1, l(), i(1), i(2), e())
	assertIterate(t, &a1, l(), i(1), i(2), e())
	a2 := [2]byte{1, 2}
	assertIterate(t, a2, bin([]byte{1, 2}))
	assertIterate(t, &a2, bin([]byte{1, 2}))
}

func TestIterateInterface(t *testing.T) {
	assertIterate(t, []interface{}{1, nil, 5.5}, l(), i(1), n(), f(5.5), e())
}

func TestIteratePointer(t *testing.T) {
	v := 1
	assertIterate(t, &v, i(1))
	pv := (*int)(nil)
	assertIterate(t, pv, n())
}

func TestIterateMap(t *testing.T) {
	assertIterate(t, map[string]int{"a": 1}, m(), s("a"), i(1), e())
	assertIterate(t, (map[string]int)(nil), n())
}

type StructTestIterate struct {
	A int
}

func TestIterateStruct(t *testing.T) {
	assertIterate(t, new(StructTestIterate), m(), s("A"), i(0), e())
	assertIterate(t, (*StructTestIterate)(nil), n())
}

func TestIterateNilOpts(t *testing.T) {
	expected := []*tevent{v(1), i(1), ed()}
	receiver := NewTER()
	IterateObject(1, receiver, nil)

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

	expected := []*tevent{v(1), mark(), i(0), m(), s("I"), i(50), s("R"), ref(), i(0), e(), ed()}
	options := IteratorOptions{
		UseReferences: true,
	}
	receiver := NewTER()
	IterateObject(obj, receiver, &options)

	if !equivalence.IsEquivalent(expected, receiver.Events) {
		t.Errorf("Expected %v but got %v", expected, receiver.Events)
	}
}
