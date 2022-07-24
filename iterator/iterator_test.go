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

	"github.com/cockroachdb/apd/v2"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/types"
)

func TestIterateBasic(t *testing.T) {
	pBigIntP := NewBigInt("12345678901234567890123456789")
	pBigIntN := NewBigInt("-999999999999999999999999999999")
	pBigFloat := NewBigFloat("5.377345e-10000")
	dfloat := NewDFloat("1.23456e1000")
	pBigDFloat := NewBDF("4.509e10000")
	gTimeNow := time.Now()
	cTimeNow := test.AsCompactTime(gTimeNow)
	pURL := NewRID("http://x.com")
	uid := types.UID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	media := types.Media{
		MediaType: "a",
		Data:      []byte{0x00},
	}
	pNode := NewNode("test", []interface{}{"a"})
	pEdge := NewEdge("a", "b", "c")

	assertIterate(t, nil, NULL())
	assertIterate(t, true, B(true))
	assertIterate(t, false, B(false))
	assertIterate(t, int(10), N(10))
	assertIterate(t, int8(10), N(10))
	assertIterate(t, int16(10), N(10))
	assertIterate(t, int32(10), N(10))
	assertIterate(t, int64(10), N(10))
	assertIterate(t, uint(10), N(10))
	assertIterate(t, uint8(10), N(10))
	assertIterate(t, uint16(10), N(10))
	assertIterate(t, uint32(10), N(10))
	assertIterate(t, uint64(10), N(10))
	assertIterate(t, 1, N(1))
	assertIterate(t, -1, N(-1))
	assertIterate(t, pBigIntP, N(pBigIntP))
	assertIterate(t, *pBigIntP, N(pBigIntP))
	assertIterate(t, pBigIntN, N(pBigIntN))
	assertIterate(t, *pBigIntN, N(pBigIntN))
	assertIterate(t, (*big.Int)(nil), NULL())
	assertIterate(t, float32(-1.25), N(-1.25))
	assertIterate(t, float64(-9.5e50), N(-9.5e50))
	assertIterate(t, pBigFloat, N(pBigFloat))
	assertIterate(t, *pBigFloat, N(pBigFloat))
	assertIterate(t, (*big.Float)(nil), NULL())
	assertIterate(t, dfloat, N(dfloat))
	assertIterate(t, pBigDFloat, N(pBigDFloat))
	assertIterate(t, *pBigDFloat, N(pBigDFloat))
	assertIterate(t, (*apd.Decimal)(nil), NULL())
	assertIterate(t, common.Float64SignalingNan, N(common.Float64SignalingNan))
	assertIterate(t, common.Float64QuietNan, N(common.Float64QuietNan))
	assertIterate(t, gTimeNow, T(compact_time.AsCompactTime(gTimeNow)))
	assertIterate(t, cTimeNow, T(cTimeNow))
	assertIterate(t, []byte{1, 2, 3, 4}, AU8([]byte{1, 2, 3, 4}))
	assertIterate(t, "test", S("test"))
	assertIterate(t, pURL, RID("http://x.com"))
	assertIterate(t, *pURL, RID("http://x.com"))
	assertIterate(t, (*url.URL)(nil), NULL())
	assertIterate(t, uid, UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}))
	assertIterate(t, &uid, UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}))
	assertIterate(t, media, BMEDIA(), ACL(1), ADT("a"), ACL(1), ADU8([]byte{0}))
	assertIterate(t, pNode, NODE(), S("test"), S("a"), E())
	assertIterate(t, *pNode, NODE(), S("test"), S("a"), E())
	assertIterate(t, pEdge, EDGE(), S("a"), S("b"), S("c"))
	assertIterate(t, *pEdge, EDGE(), S("a"), S("b"), S("c"))
}

func TestIterateArrayUint8(t *testing.T) {
	a := [2]byte{1, 2}
	assertIterate(t, a, AU8([]byte{1, 2}))
	assertIterate(t, &a, AU8([]byte{1, 2}))
	s := []byte{1, 2}
	assertIterate(t, s, AU8([]byte{1, 2}))
	assertIterate(t, &s, AU8([]byte{1, 2}))
}

func TestIterateArrayUint16(t *testing.T) {
	a := [2]uint16{0x1234, 0x5678}
	assertIterate(t, a, AU16([]uint16{0x1234, 0x5678}))
	assertIterate(t, &a, AU16([]uint16{0x1234, 0x5678}))
	s := []uint16{0x1234, 0x5678}
	assertIterate(t, s, AU16([]uint16{0x1234, 0x5678}))
	assertIterate(t, &s, AU16([]uint16{0x1234, 0x5678}))
}

func TestIterateArrayUint32(t *testing.T) {
	a := [2]uint32{0x12345678, 0xabcdef12}
	assertIterate(t, a, AU32([]uint32{0x12345678, 0xabcdef12}))
	assertIterate(t, &a, AU32([]uint32{0x12345678, 0xabcdef12}))
	s := []uint32{0x12345678, 0xabcdef12}
	assertIterate(t, s, AU32([]uint32{0x12345678, 0xabcdef12}))
	assertIterate(t, &s, AU32([]uint32{0x12345678, 0xabcdef12}))
}

func TestIterateArrayUint64(t *testing.T) {
	a := [2]uint64{0x123456789abcdef0, 0xfedcba9876543210}
	assertIterate(t, a, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	assertIterate(t, &a, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	s := []uint64{0x123456789abcdef0, 0xfedcba9876543210}
	assertIterate(t, s, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	assertIterate(t, &s, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
}

func TestIterateArrayUint(t *testing.T) {
	// Assuming 64-bit arch
	a := [2]uint{0x123456789abcdef0, 0xfedcba9876543210}
	assertIterate(t, a, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	assertIterate(t, &a, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	s := []uint{0x123456789abcdef0, 0xfedcba9876543210}
	assertIterate(t, s, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
	assertIterate(t, &s, AU64([]uint64{0x123456789abcdef0, 0xfedcba9876543210}))
}

func TestIterateArrayInt8(t *testing.T) {
	a := [2]int8{1, -2}
	assertIterate(t, a, AI8([]int8{1, -2}))
	assertIterate(t, &a, AI8([]int8{1, -2}))
	s := []int8{1, -2}
	assertIterate(t, s, AI8([]int8{1, -2}))
	assertIterate(t, &s, AI8([]int8{1, -2}))
}

func TestIterateArrayInt16(t *testing.T) {
	a := [2]int16{1000, -2000}
	assertIterate(t, a, AI16([]int16{1000, -2000}))
	assertIterate(t, &a, AI16([]int16{1000, -2000}))
	s := []int16{1000, -2000}
	assertIterate(t, s, AI16([]int16{1000, -2000}))
	assertIterate(t, &s, AI16([]int16{1000, -2000}))
}

func TestIterateArrayInt32(t *testing.T) {
	a := [2]int32{1000000, -2000000}
	assertIterate(t, a, AI32([]int32{1000000, -2000000}))
	assertIterate(t, &a, AI32([]int32{1000000, -2000000}))
	s := []int32{1000000, -2000000}
	assertIterate(t, s, AI32([]int32{1000000, -2000000}))
	assertIterate(t, &s, AI32([]int32{1000000, -2000000}))
}

func TestIterateArrayInt64(t *testing.T) {
	a := [2]int64{1000000000000, -2000000000000}
	assertIterate(t, a, AI64([]int64{1000000000000, -2000000000000}))
	assertIterate(t, &a, AI64([]int64{1000000000000, -2000000000000}))
	s := []int64{1000000000000, -2000000000000}
	assertIterate(t, s, AI64([]int64{1000000000000, -2000000000000}))
	assertIterate(t, &s, AI64([]int64{1000000000000, -2000000000000}))
}

func TestIterateArrayInt(t *testing.T) {
	// Assuming 64-bit arch
	a := [2]int{1000000000000, -2000000000000}
	assertIterate(t, a, AI64([]int64{1000000000000, -2000000000000}))
	assertIterate(t, &a, AI64([]int64{1000000000000, -2000000000000}))
	s := []int{1000000000000, -2000000000000}
	assertIterate(t, s, AI64([]int64{1000000000000, -2000000000000}))
	assertIterate(t, &s, AI64([]int64{1000000000000, -2000000000000}))
}

func TestIterateArrayFloat32(t *testing.T) {
	a := [2]float32{1.5, -1.5}
	assertIterate(t, a, AF32([]float32{1.5, -1.5}))
	assertIterate(t, &a, AF32([]float32{1.5, -1.5}))
	s := []float32{1.5, -1.5}
	assertIterate(t, s, AF32([]float32{1.5, -1.5}))
	assertIterate(t, &s, AF32([]float32{1.5, -1.5}))
}

func TestIterateArrayFloat64(t *testing.T) {
	a := [2]float64{1.5, -1.5}
	assertIterate(t, a, AF64([]float64{1.5, -1.5}))
	assertIterate(t, &a, AF64([]float64{1.5, -1.5}))
	s := []float64{1.5, -1.5}
	assertIterate(t, s, AF64([]float64{1.5, -1.5}))
	assertIterate(t, &s, AF64([]float64{1.5, -1.5}))
}

func TestIterateArrayBool(t *testing.T) {
	a := [2]bool{true, false}
	assertIterate(t, a, AB([]bool{true, false}))
	assertIterate(t, &a, AB([]bool{true, false}))
	s := []bool{true, false}
	assertIterate(t, s, AB([]bool{true, false}))
	assertIterate(t, &s, AB([]bool{true, false}))
}

func TestIterateInterface(t *testing.T) {
	assertIterate(t, []interface{}{1, nil, 5.5}, L(), N(1), NULL(), N(5.5), E())
}

func TestIteratePointer(t *testing.T) {
	v := 1
	assertIterate(t, &v, N(1))
	pv := (*int)(nil)
	assertIterate(t, pv, NULL())
}

func TestIterateMap(t *testing.T) {
	assertIterate(t, map[string]int{"a": 1}, M(), S("a"), N(1), E())
	assertIterate(t, (map[string]int)(nil), NULL())
}

type StructTestIterate struct {
	A int
}

func TestIterateStruct(t *testing.T) {
	sOpts := options.DefaultIteratorSessionOptions()
	sOpts.LowercaseStructFieldNames = false
	iOpts := options.DefaultIteratorOptions()

	assertIterate(t, new(StructTestIterate), M(), S("a"), N(0), E())
	assertIterateWithOptions(t, &sOpts, &iOpts, new(StructTestIterate), M(), S("A"), N(0), E())
	assertIterate(t, (*StructTestIterate)(nil), NULL())
}

func TestIterateNilOpts(t *testing.T) {
	expected := test.Events{EvV, N(1)}
	receiver, events := test.NewEventCollector(nil)
	iterateObject(1, receiver, nil, nil)

	if !events.IsEquivalentTo(expected) {
		t.Errorf("Expected %v but got %v", expected, events.Events)
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

	expected := test.Events{EvV, MARK("0"), M(), S("i"), N(50), S("r"), REFL("0"), E()}
	sessionOptions := options.DefaultIteratorSessionOptions()
	iteratorOptions := options.DefaultIteratorOptions()
	iteratorOptions.RecursionSupport = true
	receiver, events := test.NewEventCollector(nil)
	iterateObject(obj, receiver, &sessionOptions, &iteratorOptions)

	if !events.IsEquivalentTo(expected) {
		t.Errorf("Expected %v but got %v", expected, events.Events)
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
