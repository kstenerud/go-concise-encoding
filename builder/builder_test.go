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

package builder

import (
	"math/big"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/types"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func TestBuilderTypedArrayUint8(t *testing.T) {
	assertBuild(t, [2]uint8{0x12, 0x23}, AU8([]uint8{0x12, 0x23}))
	assertBuild(t, []uint8{0x12, 0x23}, AU8([]uint8{0x12, 0x23}))
	assertBuild(t, [2]uint16{0x1234, 0x2345}, L(), N(0x1234), N(0x2345), E())
	assertBuild(t, []uint16{0x12, 0x23}, L(), N(0x12), N(0x23), E())
}

func TestBuilderTypedArrayUint16(t *testing.T) {
	assertBuild(t, [2]uint16{0x1234, 0x2345}, AU16([]uint16{0x1234, 0x2345}))
	assertBuild(t, []uint16{0x1234, 0x2345}, AU16([]uint16{0x1234, 0x2345}))
	assertBuild(t, [2]uint16{0x1234, 0x2345}, L(), N(0x1234), N(0x2345), E())
	assertBuild(t, []uint16{0x1234, 0x2345}, L(), N(0x1234), N(0x2345), E())
}

func TestBuilderTypedArrayUint32(t *testing.T) {
	assertBuild(t, [2]uint32{0x12345678, 0x23456789}, AU32([]uint32{0x12345678, 0x23456789}))
	assertBuild(t, []uint32{0x12345678, 0x23456789}, AU32([]uint32{0x12345678, 0x23456789}))
	assertBuild(t, [2]uint32{0x12345678, 0x23456789}, L(), N(0x12345678), N(0x23456789), E())
	assertBuild(t, []uint32{0x12345678, 0x23456789}, L(), N(0x12345678), N(0x23456789), E())
}

func TestBuilderTypedArrayUint64(t *testing.T) {
	assertBuild(t, [2]uint64{0x123456789abcdef0, 0x23456789abcdef01}, AU64([]uint64{0x123456789abcdef0, 0x23456789abcdef01}))
	assertBuild(t, []uint64{0x123456789abcdef0, 0x23456789abcdef01}, AU64([]uint64{0x123456789abcdef0, 0x23456789abcdef01}))
	assertBuild(t, [2]uint64{0x123456789abcdef0, 0x23456789abcdef01}, L(), N(0x123456789abcdef0), N(0x23456789abcdef01), E())
	assertBuild(t, []uint64{0x123456789abcdef0, 0x23456789abcdef01}, L(), N(0x123456789abcdef0), N(0x23456789abcdef01), E())
}

func TestBuilderTypedArrayInt8(t *testing.T) {
	assertBuild(t, [2]int8{0x12, -0x23}, AI8([]int8{0x12, -0x23}))
	assertBuild(t, []int8{0x12, -0x23}, AI8([]int8{0x12, -0x23}))
	assertBuild(t, [2]int8{0x12, -0x23}, L(), N(0x12), N(-0x23), E())
	assertBuild(t, []int8{0x12, -0x23}, L(), N(0x12), N(-0x23), E())
}

func TestBuilderTypedArrayInt16(t *testing.T) {
	assertBuild(t, [2]int16{0x1234, -0x2345}, AI16([]int16{0x1234, -0x2345}))
	assertBuild(t, []int16{0x1234, -0x2345}, AI16([]int16{0x1234, -0x2345}))
	assertBuild(t, [2]int16{0x1234, -0x2345}, L(), N(0x1234), N(-0x2345), E())
	assertBuild(t, []int16{0x1234, -0x2345}, L(), N(0x1234), N(-0x2345), E())
}

func TestBuilderTypedArrayInt32(t *testing.T) {
	assertBuild(t, [2]int32{0x12345678, -0x23456789}, AI32([]int32{0x12345678, -0x23456789}))
	assertBuild(t, []int32{0x12345678, -0x23456789}, AI32([]int32{0x12345678, -0x23456789}))
	assertBuild(t, [2]int32{0x12345678, -0x23456789}, L(), N(0x12345678), N(-0x23456789), E())
	assertBuild(t, []int32{0x12345678, -0x23456789}, L(), N(0x12345678), N(-0x23456789), E())
}

func TestBuilderTypedArrayInt64(t *testing.T) {
	assertBuild(t, [2]int64{0x123456789abcdef0, -0x23456789abcdef01}, AI64([]int64{0x123456789abcdef0, -0x23456789abcdef01}))
	assertBuild(t, []int64{0x123456789abcdef0, -0x23456789abcdef01}, AI64([]int64{0x123456789abcdef0, -0x23456789abcdef01}))
	assertBuild(t, [2]int64{0x123456789abcdef0, -0x23456789abcdef01}, L(), N(0x123456789abcdef0), N(-0x23456789abcdef01), E())
	assertBuild(t, []int64{0x123456789abcdef0, -0x23456789abcdef01}, L(), N(0x123456789abcdef0), N(-0x23456789abcdef01), E())
}

func TestBuilderTypedArrayFloat32(t *testing.T) {
	assertBuild(t, [2]float32{-1.25, 9.5e10}, AF32([]float32{-1.25, 9.5e10}))
	assertBuild(t, []float32{-1.25, 9.5e10}, AF32([]float32{-1.25, 9.5e10}))
	assertBuild(t, [2]float32{-1.25, 9.5e10}, L(), N(-1.25), N(9.5e10), E())
	assertBuild(t, []float32{-1.25, 9.5e10}, L(), N(-1.25), N(9.5e10), E())
}

func TestBuilderTypedArrayFloat64(t *testing.T) {
	assertBuild(t, [2]float64{-1.25, 9.5e10}, AF64([]float64{-1.25, 9.5e10}))
	assertBuild(t, []float64{-1.25, 9.5e10}, AF64([]float64{-1.25, 9.5e10}))
	assertBuild(t, [2]float64{-1.25, 9.5e10}, L(), N(-1.25), N(9.5e10), E())
	assertBuild(t, []float64{-1.25, 9.5e10}, L(), N(-1.25), N(9.5e10), E())
}

// =================================================

func TestBuildUnknown(t *testing.T) {
	expected := []interface{}{1}
	actual := runBuild(NewSession(nil, nil), nil, L(), N(1), E())

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

var zero = float64(0)
var negZero = -zero

func TestBuilderBasicTypes(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}

	pBigIntP := NewBigInt("12345678901234567890123456789")
	pBigIntN := NewBigInt("-999999999999999999999999999999")
	pBigFloat := NewBigFloat("1.2345678901234567890123456789e10000")
	dfloat := NewDFloat("1.23456e1000")
	pBigDFloat := NewBDF("4.509e10000")
	gTime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	cTime := test.AsCompactTime(gTime)
	cTimeLL := test.NewTimeLL(10, 5, 59, 100, 506, 107)
	pURL := NewRID("http://x.com")
	pNode := NewNode("test", []interface{}{"a"})
	pEdge := NewEdge("a", "b", "c")

	assertBuild(t, true, B(true))
	assertBuild(t, false, B(false))
	assertBuild(t, int(10), N(10))
	assertBuild(t, int8(10), N(10))
	assertBuild(t, int16(-10), N(-10))
	assertBuild(t, int32(10), N(10))
	assertBuild(t, int64(-10), N(-10))
	assertBuild(t, uint(10), N(10))
	assertBuild(t, uint8(10), N(10))
	assertBuild(t, uint16(10), N(10))
	assertBuild(t, uint32(10), N(10))
	assertBuild(t, uint64(10), N(10))
	assertBuild(t, 1, N(1))
	assertBuild(t, -1, N(-1))
	assertBuild(t, pBigIntP, N(pBigIntP))
	assertBuild(t, *pBigIntP, N(pBigIntP))
	assertBuild(t, pBigIntN, N(pBigIntN))
	assertBuild(t, *pBigIntN, N(pBigIntN))
	assertBuild(t, (*big.Int)(nil), NULL())
	assertBuild(t, float32(negZero), N(negZero))
	assertBuild(t, float32(-1.25), N(-1.25))
	assertBuild(t, float64(-9.5e50), N(-9.5e50))
	assertBuild(t, pBigFloat, N(pBigFloat))
	assertBuild(t, *pBigFloat, N(pBigFloat))
	assertBuild(t, (*big.Float)(nil), NULL())
	assertBuild(t, dfloat, N(dfloat))
	assertBuild(t, pBigDFloat, N(pBigDFloat))
	assertBuild(t, *pBigDFloat, N(pBigDFloat))
	assertBuild(t, (*apd.Decimal)(nil), NULL())
	assertBuild(t, common.Float64SignalingNan, SNAN())
	assertBuild(t, common.Float64QuietNan, NAN())
	assertBuild(t, gTime, T(compact_time.AsCompactTime(gTime)))
	assertBuild(t, cTime, T(cTime))
	assertBuild(t, cTimeLL, T(cTimeLL))
	assertBuild(t, []byte{1, 2, 3, 4}, AU8([]byte{1, 2, 3, 4}))
	assertBuild(t, "test", S("test"))
	assertBuild(t, pURL, RID("http://x.com"))
	assertBuild(t, *pURL, RID("http://x.com"))
	assertBuild(t, (*url.URL)(nil), NULL())
	assertBuild(t, interface{}(1234), N(1234))
	assertBuild(t, pNode, NODE(), S("test"), S("a"), E())
	assertBuild(t, *pNode, NODE(), S("test"), S("a"), E())
	assertBuild(t, pEdge, EDGE(), S("a"), S("b"), S("c"))
	assertBuild(t, *pEdge, EDGE(), S("a"), S("b"), S("c"))
}

func TestBuilderConvertToBDF(t *testing.T) {
	pv := NewBDF("1")
	nv := NewBDF("-1")
	nz := NewBDF("-0")

	assertBuild(t, pv, N(1))
	assertBuild(t, nv, N(-1))
	assertBuild(t, pv, N(NewBigInt("1")))
	assertBuild(t, pv, N(1))
	assertBuild(t, pv, N(NewBigFloat("1")))
	assertBuild(t, pv, N(NewDFloat("1")))
	assertBuild(t, pv, N(NewBDF("1")))

	assertBuild(t, *pv, N(1))
	assertBuild(t, *nv, N(-1))
	assertBuild(t, *pv, N(NewBigInt("1")))
	assertBuild(t, *pv, N(1))
	assertBuild(t, *pv, N(NewBigFloat("1")))
	assertBuild(t, *pv, N(NewDFloat("1")))
	assertBuild(t, *pv, N(NewBDF("1")))

	assertBuild(t, nz, N(negZero))
	assertBuild(t, nz, N(NewBigFloat("-0")))
	assertBuild(t, nz, N(NewDFloat("-0")))
	assertBuild(t, nz, N(NewBDF("-0")))
}

func TestBuilderConvertToBDFFail(t *testing.T) {
	v := NewBDF("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, NULL())
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToBF(t *testing.T) {
	pv := NewBigFloat("1")
	nv := NewBigFloat("-1")
	nz := NewBigFloat("-0")

	assertBuild(t, pv, N(1))
	assertBuild(t, nv, N(-1))
	assertBuild(t, pv, N(NewBigInt("1")))
	assertBuild(t, pv, N(1))
	assertBuild(t, pv, N(NewBigFloat("1")))
	assertBuild(t, pv, N(NewDFloat("1")))
	assertBuild(t, pv, N(NewBDF("1")))

	assertBuild(t, *pv, N(1))
	assertBuild(t, *nv, N(-1))
	assertBuild(t, *pv, N(NewBigInt("1")))
	assertBuild(t, *pv, N(1))
	assertBuild(t, *pv, N(NewBigFloat("1")))
	assertBuild(t, *pv, N(NewDFloat("1")))
	assertBuild(t, *pv, N(NewBDF("1")))

	assertBuild(t, *nz, N(negZero))
	assertBuild(t, *nz, N(NewBigFloat("-0")))
	assertBuild(t, *nz, N(NewDFloat("-0")))
	assertBuild(t, *nz, N(NewBDF("-0")))
}

func TestBuilderConvertToBFFail(t *testing.T) {
	v := NewBigFloat("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, NULL())
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToBI(t *testing.T) {
	pv := NewBigInt("1")
	nv := NewBigInt("-1")
	assertBuild(t, pv, N(1))
	assertBuild(t, NewBigInt("9223372036854775808"), N(uint64(9223372036854775808)))
	assertBuild(t, nv, N(-1))
	assertBuild(t, pv, N(NewBigInt("1")))
	assertBuild(t, pv, N(1))
	assertBuild(t, pv, N(NewBigFloat("1")))
	assertBuild(t, pv, N(NewDFloat("1")))
	assertBuild(t, pv, N(NewBDF("1")))

	assertBuild(t, *pv, N(1))
	assertBuild(t, *nv, N(-1))
	assertBuild(t, *pv, N(NewBigInt("1")))
	assertBuild(t, *pv, N(1))
	assertBuild(t, *pv, N(NewBigFloat("1")))
	assertBuild(t, *pv, N(NewDFloat("1")))
	assertBuild(t, *pv, N(NewBDF("1")))
}

func TestBuilderConvertToBIFail(t *testing.T) {
	v := NewBigInt("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, N(1.1))
	assertBuildPanics(t, v, N(NewBigFloat("1.1")))
	assertBuildPanics(t, v, N(NewBigFloat("1.0e100000")))
	assertBuildPanics(t, v, N(NewDFloat("1.1")))
	assertBuildPanics(t, v, N(NewDFloat("1.0e100000")))
	assertBuildPanics(t, v, N(NewBDF("1.1")))
	assertBuildPanics(t, v, N(NewBDF("1.0e100000")))
	assertBuildPanics(t, v, N(NewBDF("nan")))
	assertBuildPanics(t, v, N(NewBDF("snan")))
	assertBuildPanics(t, v, N(NewBDF("inf")))
	assertBuildPanics(t, v, N(NewBDF("-inf")))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, NULL())
	assertBuildPanics(t, *v, N(1.1))
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, N(NewBigFloat("1.1")))
	assertBuildPanics(t, *v, N(NewBigFloat("1.0e100000")))
	assertBuildPanics(t, *v, N(NewDFloat("1.1")))
	assertBuildPanics(t, *v, N(NewDFloat("1.0e100000")))
	assertBuildPanics(t, *v, N(NewBDF("1.1")))
	assertBuildPanics(t, *v, N(NewBDF("1.0e100000")))
	assertBuildPanics(t, *v, N(NewBDF("nan")))
	assertBuildPanics(t, *v, N(NewBDF("snan")))
	assertBuildPanics(t, *v, N(NewBDF("inf")))
	assertBuildPanics(t, *v, N(NewBDF("-inf")))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToDecimalFloat(t *testing.T) {
	pv := NewDFloat("1")
	nv := NewDFloat("-1")
	nz := NewDFloat("-0")

	assertBuild(t, pv, N(1))
	assertBuild(t, nv, N(-1))
	assertBuild(t, pv, N(NewBigInt("1")))
	assertBuild(t, pv, N(1))
	assertBuild(t, pv, N(NewBigFloat("1")))
	assertBuild(t, pv, N(NewDFloat("1")))
	assertBuild(t, pv, N(NewBDF("1")))
	assertBuild(t, nz, N(negZero))
	assertBuild(t, nz, N(NewBigFloat("-0")))
	assertBuild(t, nz, N(NewDFloat("-0")))
	assertBuild(t, nz, N(NewBDF("-0")))
}

func TestBuilderDecimalFloatFail(t *testing.T) {
	v := NewDFloat("1")
	assertBuildPanics(t, v, NULL())
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())
}

func TestBuilderConvertToFloat(t *testing.T) {
	pv := 1.0
	nv := -1.0
	nz := negZero

	assertBuild(t, pv, N(1))
	assertBuild(t, nv, N(-1))
	assertBuild(t, pv, N(NewBigInt("1")))
	assertBuild(t, pv, N(1))
	assertBuild(t, pv, N(NewBigFloat("1")))
	assertBuild(t, pv, N(NewDFloat("1")))
	assertBuild(t, pv, N(NewBDF("1")))

	assertBuild(t, nz, N(negZero))
	assertBuild(t, nz, N(NewBigFloat("-0")))
	assertBuild(t, nz, N(NewDFloat("-0")))
	assertBuild(t, nz, N(NewBDF("-0")))
}

func TestBuilderConvertToFloatFail(t *testing.T) {
	// TODO: How to define required conversion accuracy?
	v := 1.0
	assertBuildPanics(t, v, NULL())
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, N(uint64(0xffffffffffffffff)))
	assertBuildPanics(t, v, N(-0x7fffffffffffffff))
	assertBuildPanics(t, v, N(NewBigFloat("1.0e309")))
	assertBuildPanics(t, v, N(NewBigFloat("1.0e-311")))
	// TODO: apd.Decimal and compact_float.DFloat don't handle float overflow
	assertBuildPanics(t, v, N(NewBigInt("1234567890123456789012345")))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())
}

func TestBuilderConvertToInt(t *testing.T) {
	assertBuild(t, 1, N(1))
	assertBuild(t, -1, N(-1))
	assertBuild(t, 1, N(NewBigInt("1")))
	assertBuild(t, 1, N(1))
	assertBuild(t, 1, N(NewBigFloat("1")))
	assertBuild(t, 1, N(NewDFloat("1")))
	assertBuild(t, 1, N(NewBDF("1")))
}

func TestBuilderConvertToIntFail(t *testing.T) {
	assertBuildPanics(t, int(1), NULL())
	assertBuildPanics(t, int(1), B(true))
	assertBuildPanics(t, int(1), N(uint64(0x8000000000000000)))
	assertBuildPanics(t, int(1), N(1.1))
	assertBuildPanics(t, int(1), N(NewBigFloat("1.1")))
	assertBuildPanics(t, int(1), N(NewDFloat("1.1")))
	assertBuildPanics(t, int(1), N(NewBDF("1.1")))
	assertBuildPanics(t, int(1), N(NewBigFloat("1e20")))
	assertBuildPanics(t, int(1), N(NewDFloat("1e20")))
	assertBuildPanics(t, int(1), N(NewBDF("1e20")))
	assertBuildPanics(t, int(1), N(NewBigFloat("-1e20")))
	assertBuildPanics(t, int(1), N(NewDFloat("-1e20")))
	assertBuildPanics(t, int(1), N(NewBDF("-1e20")))
	assertBuildPanics(t, int(1), N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, int(1), N(NewBigInt("-100000000000000000000")))
	assertBuildPanics(t, int(1), S("1"))
	assertBuildPanics(t, int(1), AU8([]byte{1}))
	assertBuildPanics(t, int(1), RID("x://x"))
	assertBuildPanics(t, int(1), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, int(1), T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, int(1), L())
	assertBuildPanics(t, int(1), M())
	assertBuildPanics(t, int(1), E())

	assertBuildPanics(t, int8(1), N(1000))
	assertBuildPanics(t, int8(1), N(-1000))
	assertBuildPanics(t, int8(1), N(1000))
	assertBuildPanics(t, int8(1), N(NewBigInt("1000")))
	assertBuildPanics(t, int8(1), N(1000))
	assertBuildPanics(t, int8(1), N(NewBigFloat("1000")))
	assertBuildPanics(t, int8(1), N(NewDFloat("1000")))
	assertBuildPanics(t, int8(1), N(NewBDF("1000")))

	assertBuildPanics(t, int16(1), N(100000))
	assertBuildPanics(t, int16(1), N(-100000))
	assertBuildPanics(t, int16(1), N(100000))
	assertBuildPanics(t, int16(1), N(NewBigInt("100000")))
	assertBuildPanics(t, int16(1), N(100000))
	assertBuildPanics(t, int16(1), N(NewBigFloat("100000")))
	assertBuildPanics(t, int16(1), N(NewDFloat("100000")))
	assertBuildPanics(t, int16(1), N(NewBDF("100000")))

	assertBuildPanics(t, int32(1), N(10000000000))
	assertBuildPanics(t, int32(1), N(-10000000000))
	assertBuildPanics(t, int32(1), N(10000000000))
	assertBuildPanics(t, int32(1), N(NewBigInt("10000000000")))
	assertBuildPanics(t, int32(1), N(10000000000))
	assertBuildPanics(t, int32(1), N(NewBigFloat("10000000000")))
	assertBuildPanics(t, int32(1), N(NewDFloat("10000000000")))
	assertBuildPanics(t, int32(1), N(NewBDF("10000000000")))
}

func TestBuilderConvertToUint(t *testing.T) {
	assertBuild(t, uint(1), N(1))
	assertBuild(t, uint(1), N(1))
	assertBuild(t, uint(1), N(NewBigInt("1")))
	assertBuild(t, uint(1), N(1))
	assertBuild(t, uint(1), N(NewBigFloat("1")))
	assertBuild(t, uint(1), N(NewDFloat("1")))
	assertBuild(t, uint(1), N(NewBDF("1")))
}

func TestBuilderConvertToUintFail(t *testing.T) {
	assertBuildPanics(t, uint(1), NULL())
	assertBuildPanics(t, uint(1), B(true))
	assertBuildPanics(t, uint(1), N(-1))
	assertBuildPanics(t, uint(1), N(1.1))
	assertBuildPanics(t, uint(1), N(NewBigFloat("1.1")))
	assertBuildPanics(t, uint(1), N(NewDFloat("1.1")))
	assertBuildPanics(t, uint(1), N(NewBDF("1.1")))
	assertBuildPanics(t, uint(1), N(NewBigFloat("1e20")))
	assertBuildPanics(t, uint(1), N(NewDFloat("1e20")))
	assertBuildPanics(t, uint(1), N(NewBDF("1e20")))
	assertBuildPanics(t, uint8(1), N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, uint(1), S("1"))
	assertBuildPanics(t, uint(1), AU8([]byte{1}))
	assertBuildPanics(t, uint(1), RID("x://x"))
	assertBuildPanics(t, uint(1), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, uint(1), T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, uint(1), L())
	assertBuildPanics(t, uint(1), M())
	assertBuildPanics(t, uint(1), E())

	assertBuildPanics(t, uint8(1), N(1000))
	assertBuildPanics(t, uint8(1), N(1000))
	assertBuildPanics(t, uint8(1), N(NewBigInt("1000")))
	assertBuildPanics(t, uint8(1), N(1000))
	assertBuildPanics(t, uint8(1), N(NewBigFloat("1000")))
	assertBuildPanics(t, uint8(1), N(NewDFloat("1000")))
	assertBuildPanics(t, uint8(1), N(NewBDF("1000")))

	assertBuildPanics(t, uint16(1), N(100000))
	assertBuildPanics(t, uint16(1), N(100000))
	assertBuildPanics(t, uint16(1), N(NewBigInt("100000")))
	assertBuildPanics(t, uint16(1), N(100000))
	assertBuildPanics(t, uint16(1), N(NewBigFloat("100000")))
	assertBuildPanics(t, uint16(1), N(NewDFloat("100000")))
	assertBuildPanics(t, uint16(1), N(NewBDF("100000")))

	assertBuildPanics(t, uint32(1), N(10000000000))
	assertBuildPanics(t, uint32(1), N(10000000000))
	assertBuildPanics(t, uint32(1), N(NewBigInt("10000000000")))
	assertBuildPanics(t, uint32(1), N(10000000000))
	assertBuildPanics(t, uint32(1), N(NewBigFloat("10000000000")))
	assertBuildPanics(t, uint32(1), N(NewDFloat("10000000000")))
	assertBuildPanics(t, uint32(1), N(NewBDF("10000000000")))
}

func TestBuilderString(t *testing.T) {
	assertBuild(t, "", NULL())
	assertBuild(t, "test", S("test"))
}

func TestBuilderStringFail(t *testing.T) {
	assertBuildPanics(t, "", B(false))
	assertBuildPanics(t, "", N(1))
	assertBuildPanics(t, "", N(-1))
	assertBuildPanics(t, "", N(1.1))
	assertBuildPanics(t, "", N(NewBigFloat("1.1")))
	assertBuildPanics(t, "", N(NewDFloat("1.1")))
	assertBuildPanics(t, "", N(NewBDF("1.1")))
	assertBuildPanics(t, "", N(NewBigFloat("1e20")))
	assertBuildPanics(t, "", N(NewDFloat("1e20")))
	assertBuildPanics(t, "", N(NewBDF("1e20")))
	assertBuildPanics(t, "", N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, "", AU8([]byte{1}))
	assertBuildPanics(t, "", RID("x://x"))
	assertBuildPanics(t, "", UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, "", T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, "", L())
	assertBuildPanics(t, "", M())
	assertBuildPanics(t, "", E())
}

func TestBuilderChunkedBytes(t *testing.T) {
	assertBuild(t, []byte{1, 2, 3, 4}, BAU8(), ACM(2), ADU8([]byte{1, 2}), ACL(2), ADU8([]byte{3, 4}))
	assertBuild(t, []byte{1, 2, 3, 4}, BAU8(), ACM(2), ADU8([]byte{1, 2}), ACM(2), ADU8([]byte{3, 4}), ACL(0))
}

func TestBuilderChunkedString(t *testing.T) {
	assertBuild(t, "test", BS(), ACM(2), ADT("te"), ACL(2), ADT("st"))
	assertBuild(t, "test", BS(), ACM(2), ADT("te"), ACM(2), ADT("st"), ACL(0))
}

func TestBuilderChunkedRID(t *testing.T) {
	expected := NewRID("test")
	assertBuild(t, expected, BRID(), ACM(2), ADT("te"), ACL(2), ADT("st"))
	assertBuild(t, expected, BRID(), ACM(2), ADT("te"), ACM(2), ADT("st"), ACL(0))
}

type CustomBinaryExampleType uint32

func TestBuilderChunkedCustomBinary(t *testing.T) {
	opts := options.DefaultBuilderSessionOptions()
	opts.CustomBinaryBuildFunction = func(src []byte, dst reflect.Value) error {
		var accum CustomBinaryExampleType
		for _, b := range src {
			accum = (accum << 8) | CustomBinaryExampleType(b)
		}
		dst.SetUint(uint64(accum))
		return nil
	}
	opts.CustomBuiltTypes = append(opts.CustomBuiltTypes, reflect.TypeOf(CustomBinaryExampleType(0)))
	session := NewSession(nil, &opts)
	expected := CustomBinaryExampleType(0x01020304)
	assertBuildWithSession(t, session, expected, BCB(), ACM(2), ADU8([]byte{1, 2}), ACL(2), ADU8([]byte{3, 4}))
	assertBuildWithSession(t, session, expected, BCB(), ACM(2), ADU8([]byte{1, 2}), ACM(2), ADU8([]byte{3, 4}), ACL(0))
}

type CustomTextExampleType uint32

func TestBuilderChunkedCustomText(t *testing.T) {
	opts := options.DefaultBuilderSessionOptions()
	opts.CustomTextBuildFunction = func(src []byte, dst reflect.Value) error {
		v, err := strconv.ParseUint(string(src), 16, 64)
		if err != nil {
			return err
		}
		dst.SetUint(v)
		return nil
	}
	opts.CustomBuiltTypes = append(opts.CustomBuiltTypes, reflect.TypeOf(CustomTextExampleType(0)))
	session := NewSession(nil, &opts)
	expected := CustomTextExampleType(0x1234)
	assertBuildWithSession(t, session, expected, BCT(), ACM(2), ADU8([]byte{'1', '2'}), ACL(2), ADU8([]byte{'3', '4'}))
	assertBuildWithSession(t, session, expected, BCT(), ACM(2), ADU8([]byte{'1', '2'}), ACM(2), ADU8([]byte{'3', '4'}), ACL(0))
}

func TestBuilderGoTime(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, gtime, T(compact_time.AsCompactTime(gtime)))
	assertBuild(t, gtime, T(ctime))
}

func TestBuilderGoTimeFail(t *testing.T) {
	gtime := time.Time{}
	ctime := test.NewTimeLL(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, gtime, NULL())
	assertBuildPanics(t, gtime, B(true))
	assertBuildPanics(t, gtime, N(1))
	assertBuildPanics(t, gtime, N(-1))
	assertBuildPanics(t, gtime, N(1.1))
	assertBuildPanics(t, gtime, N(NewBigFloat("1.1")))
	assertBuildPanics(t, gtime, N(NewDFloat("1.1")))
	assertBuildPanics(t, gtime, N(NewBDF("1.1")))
	assertBuildPanics(t, gtime, N(NewBigFloat("1e20")))
	assertBuildPanics(t, gtime, N(NewDFloat("1e20")))
	assertBuildPanics(t, gtime, N(NewBDF("1e20")))
	assertBuildPanics(t, gtime, N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, gtime, S("1"))
	assertBuildPanics(t, gtime, AU8([]byte{1}))
	assertBuildPanics(t, gtime, RID("x://x"))
	assertBuildPanics(t, gtime, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, gtime, T(ctime))
	assertBuildPanics(t, gtime, L())
	assertBuildPanics(t, gtime, M())
	assertBuildPanics(t, gtime, E())
}

func TestBuilderCompactTime(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, ctime, T(compact_time.AsCompactTime(gtime)))
	assertBuild(t, ctime, T(ctime))
}

func TestBuilderCompactTimeFail(t *testing.T) {
	ctime := test.NewTimeLL(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, ctime, B(true))
	assertBuildPanics(t, ctime, N(1))
	assertBuildPanics(t, ctime, N(-1))
	assertBuildPanics(t, ctime, N(1.1))
	assertBuildPanics(t, ctime, N(NewBigFloat("1.1")))
	assertBuildPanics(t, ctime, N(NewDFloat("1.1")))
	assertBuildPanics(t, ctime, N(NewBDF("1.1")))
	assertBuildPanics(t, ctime, N(NewBigFloat("1e20")))
	assertBuildPanics(t, ctime, N(NewDFloat("1e20")))
	assertBuildPanics(t, ctime, N(NewBDF("1e20")))
	assertBuildPanics(t, ctime, N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, ctime, S("1"))
	assertBuildPanics(t, ctime, AU8([]byte{1}))
	assertBuildPanics(t, ctime, RID("x://x"))
	assertBuildPanics(t, ctime, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, ctime, L())
	assertBuildPanics(t, ctime, M())
	assertBuildPanics(t, ctime, E())
}

func TestBuilderSlice(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, []bool{false, true}, L(), B(false), B(true), E())
	assertBuild(t, []int8{-1, 2, 3, 4, 5, 6, 7}, L(), N(-1), N(2), N(3),
		N(NewBigInt("4")), N(NewBigFloat("5")), N(NewDFloat("6")),
		N(NewBDF("7")), E())
	assertBuild(t, []*int{nil}, L(), NULL(), E())
	assertBuild(t, []string{"test"}, L(), S("test"), E())
	assertBuild(t, [][]byte{{1}}, L(), AU8([]byte{1}), E())
	assertBuild(t, []*url.URL{NewRID("http://example.com")}, L(), RID("http://example.com"), E())
	assertBuild(t, []time.Time{gtime}, L(), T(compact_time.AsCompactTime(gtime)), E())
	assertBuild(t, []compact_time.Time{ctime}, L(), T(ctime), E())
	assertBuild(t, [][]int{{1}}, L(), L(), N(1), E(), E())
	assertBuild(t, []map[int]int{{1: 2}}, L(), M(), N(1), N(2), E(), E())
}

func TestBuilderSliceFail(t *testing.T) {
	assertBuildPanics(t, []int{}, NULL())
	assertBuildPanics(t, []int{}, M())
	assertBuildPanics(t, [][]int{}, L(), M())
}

func TestBuilderArray(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, [2]bool{false, true}, L(), B(false), B(true), E())
	assertBuild(t, [7]int8{-1, 2, 3, 4, 5, 6, 7}, L(), N(-1), N(2), N(3),
		N(NewBigInt("4")), N(NewBigFloat("5")), N(NewDFloat("6")),
		N(NewBDF("7")), E())
	assertBuild(t, [1]*int{nil}, L(), NULL(), E())
	assertBuild(t, [1]string{"test"}, L(), S("test"), E())
	assertBuild(t, [1][]byte{{1}}, L(), AU8([]byte{1}), E())
	assertBuild(t, [1]*url.URL{NewRID("http://example.com")}, L(), RID("http://example.com"), E())
	assertBuild(t, [1]time.Time{gtime}, L(), T(compact_time.AsCompactTime(gtime)), E())
	assertBuild(t, [1]compact_time.Time{ctime}, L(), T(ctime), E())
	assertBuild(t, [1][]int{{1}}, L(), L(), N(1), E(), E())
	assertBuild(t, [1]map[int]int{{1: 2}}, L(), M(), N(1), N(2), E(), E())
}

func TestBuilderArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]int{}, NULL())
	assertBuildPanics(t, [1]int{}, M())
	assertBuildPanics(t, [1][]int{}, L(), M())
}

func TestBuilderByteArray(t *testing.T) {
	assertBuild(t, [1]byte{1}, AU8([]byte{1}))
}

func TestBuilderByteArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]byte{}, NULL())
	assertBuildPanics(t, [1]byte{}, B(false))
	assertBuildPanics(t, [1]byte{}, N(1))
	assertBuildPanics(t, [1]byte{}, N(-1))
	assertBuildPanics(t, [1]byte{}, N(1.1))
	assertBuildPanics(t, [1]byte{}, N(NewBigFloat("1.1")))
	assertBuildPanics(t, [1]byte{}, N(NewDFloat("1.1")))
	assertBuildPanics(t, [1]byte{}, N(NewBDF("1.1")))
	assertBuildPanics(t, [1]byte{}, N(NewBigFloat("1e20")))
	assertBuildPanics(t, [1]byte{}, N(NewDFloat("1e20")))
	assertBuildPanics(t, [1]byte{}, N(NewBDF("1e20")))
	assertBuildPanics(t, [1]byte{}, N(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, [1]byte{}, S(""))
	assertBuildPanics(t, [1]byte{}, RID("x://x"))
	assertBuildPanics(t, [1]byte{}, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, [1]byte{}, T(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, [1]byte{}, T(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, [1]byte{}, M())
	assertBuildPanics(t, [1]byte{}, E())
}

func TestBuilderMap(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, map[int]interface{}{
		1:  nil,
		2:  true,
		3:  1,
		4:  -1,
		5:  1.1,
		6:  NewBigFloat("1.1"),
		7:  NewDFloat("1.1"),
		8:  NewBDF("1.1"),
		9:  NewBigInt("100000000000000000000"),
		10: "test",
		11: NewRID("http://example.com"),
		12: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		13: gtime,
		14: gtime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		M(),
		N(1), NULL(),
		N(2), B(true),
		N(3), N(1),
		N(4), N(-1),
		N(5), N(1.1),
		N(6), N(NewBigFloat("1.1")),
		N(7), N(NewDFloat("1.1")),
		N(8), N(NewBDF("1.1")),
		N(9), N(NewBigInt("100000000000000000000")),
		N(10), S("test"),
		N(11), RID("http://example.com"),
		N(12), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		N(13), T(compact_time.AsCompactTime(gtime)),
		N(14), T(ctime),
		N(15), L(), N(1), E(),
		N(16), M(), N(1), N(2), E(),
		N(17), AU8([]byte{1}),
		E())
}

func TestBuilderStruct(t *testing.T) {
	s := NewTestingOuterStruct(1)
	includeFakes := true
	assertBuild(t, s, s.GetRepresentativeEvents(includeFakes)...)
}

func TestBuilderInterfaceSlice(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, []interface{}{
		// nil,
		true,
		1,
		-1,
		1.1,
		NewBigFloat("1.1"),
		NewDFloat("1.1"),
		NewBDF("1.1"),
		NewBigInt("100000000000000000000"),
		"test",
		NewRID("http://example.com"),
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		gtime,
		gtime,
		[]float64{1},
		map[int]int{1: 2},
		[]byte{1},
	}, L(),
		// n(),
		B(true),
		N(1),
		N(-1),
		N(1.1),
		N(NewBigFloat("1.1")),
		N(NewDFloat("1.1")),
		N(NewBDF("1.1")),
		N(NewBigInt("100000000000000000000")),
		S("test"),
		RID("http://example.com"),
		UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		T(compact_time.AsCompactTime(gtime)),
		T(ctime),
		L(), N(1), E(),
		M(), N(1), N(2), E(),
		AU8([]byte{1}),
		E())
}

func TestBuilderInterfaceMap(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, map[interface{}]interface{}{
		1:  nil,
		2:  true,
		3:  1,
		4:  -1,
		5:  1.1,
		6:  NewBigFloat("1.1"),
		7:  NewDFloat("1.1"),
		8:  NewBDF("1.1"),
		9:  NewBigInt("100000000000000000000"),
		10: "test",
		11: NewRID("http://example.com"),
		12: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		13: gtime,
		14: gtime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		M(),
		N(1), NULL(),
		N(2), B(true),
		N(3), N(1),
		N(4), N(-1),
		N(5), N(1.1),
		N(6), N(NewBigFloat("1.1")),
		N(7), N(NewDFloat("1.1")),
		N(8), N(NewBDF("1.1")),
		N(9), N(NewBigInt("100000000000000000000")),
		N(10), S("test"),
		N(11), RID("http://example.com"),
		N(12), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		N(13), T(compact_time.AsCompactTime(gtime)),
		N(14), T(ctime),
		N(15), L(), N(1), E(),
		N(16), M(), N(1), N(2), E(),
		N(17), AU8([]byte{1}),
		E())
}

// // Older tests

type BuilderTestStruct struct {
	internal string //nolint
	ABool    bool
	AString  string
	AnInt    int
	AMap     map[int]int8
	ASlice   []string
}

func TestBuilderStructCaseInsensitive(t *testing.T) {
	assertBuild(t,
		&BuilderTestStruct{
			AString: "test",
			AnInt:   1,
			ABool:   true,
			AMap:    map[int]int8{1: 50},
			ASlice:  []string{"the slice"},
		},
		M(),
		S("astring"), S("test"),
		S("AMAP"), M(), N(1), N(50), E(),
		S("AnINT"), N(1),
		S("Aslice"), L(), S("the slice"), E(),
		S("abool"), B(true),
		E())
}

func TestBuilderStructIgnored(t *testing.T) {
	assertBuild(t, BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
	}, M(), S("AString"), S("test"), S("Something"), N(5), S("AnInt"), N(1), S("ABool"), B(true), E())
}

func TestBuilderListStruct(t *testing.T) {
	assertBuild(t,
		[]BuilderTestStruct{
			{
				AString: "test",
				AnInt:   1,
				ABool:   true,
				AMap:    map[int]int8{1: 50},
				ASlice:  []string{"the slice"},
			},
		},
		L(),
		M(),
		S("AString"), S("test"),
		S("AMap"), M(), N(1), N(50), E(),
		S("AnInt"), N(1),
		S("ASlice"), L(), S("the slice"), E(),
		S("ABool"), B(true),
		E(),
		E())
}

func TestBuilderMapStruct(t *testing.T) {
	assertBuild(t,
		map[string]BuilderTestStruct{
			"struct": {
				AString: "test",
				AnInt:   1,
				ABool:   true,
				AMap:    map[int]int8{1: 50},
				ASlice:  []string{"the slice"},
			},
		},
		M(),
		S("struct"),
		M(),
		S("AString"), S("test"),
		S("AMap"), M(), N(1), N(50), E(),
		S("AnInt"), N(1),
		S("ASlice"), L(), S("the slice"), E(),
		S("ABool"), B(true),
		E(),
		E())
}

func TestBuilderMultipleComplexBuilds(t *testing.T) {
	v := BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
		AMap:    map[int]int8{1: 50},
		ASlice:  []string{"the slice"},
	}

	for idx := 0; idx < 10; idx++ {
		assertBuild(t,
			v,
			M(),
			S("AString"), S("test"),
			S("AMap"), M(), N(1), N(50), E(),
			S("AnInt"), N(1),
			S("ASlice"), L(), S("the slice"), E(),
			S("ABool"), B(true),
			E())
	}
}

type BuilderPtrTestStruct struct {
	internal    string //nolint
	ABool       *bool
	AnInt       *int
	AnInt8      *int8
	AnInt16     *int16
	AnInt32     *int32
	AnInt64     *int64
	AUint       *uint
	AUint8      *uint8
	AUint16     *uint16
	AUint32     *uint32
	AUint64     *uint64
	AFloat32    *float32
	AFloat64    *float64
	AString     *string
	AnInterface *interface{}
}

func TestBuilderPtr(t *testing.T) {
	aBool := true
	anInt := int(1)
	anInt8 := int8(2)
	anInt16 := int16(3)
	anInt32 := int32(4)
	anInt64 := int64(5)
	aUint := uint(6)
	aUint8 := uint8(7)
	aUint16 := uint16(8)
	aUint32 := uint32(9)
	aUint64 := uint64(10)
	aFloat32 := float32(11.5)
	aFloat64 := float64(12.5)
	aString := "test"
	var anInterface interface{} = 100
	v := BuilderPtrTestStruct{
		ABool:       &aBool,
		AnInt:       &anInt,
		AnInt8:      &anInt8,
		AnInt16:     &anInt16,
		AnInt32:     &anInt32,
		AnInt64:     &anInt64,
		AUint:       &aUint,
		AUint8:      &aUint8,
		AUint16:     &aUint16,
		AUint32:     &aUint32,
		AUint64:     &aUint64,
		AFloat32:    &aFloat32,
		AFloat64:    &aFloat64,
		AString:     &aString,
		AnInterface: &anInterface,
	}
	assertBuild(t,
		v,
		M(),
		S("ABool"), B(true),
		S("AnInt"), N(1),
		S("AnInt8"), N(2),
		S("AnInt16"), N(3),
		S("AnInt32"), N(4),
		S("AnInt64"), N(5),
		S("AUint"), N(6),
		S("AUint8"), N(7),
		S("AUint16"), N(8),
		S("AUint32"), N(9),
		S("AUint64"), N(10),
		S("AFloat32"), N(11.5),
		S("AFloat64"), N(12.5),
		S("AString"), S("test"),
		S("AnInterface"), N(100),
		E())
}

type BuilderSliceTestStruct struct {
	internal    []string //nolint
	ABool       []bool
	AnInt       []int
	AnInt8      []int8
	AnInt16     []int16
	AnInt32     []int32
	AnInt64     []int64
	AUint       []uint
	AUint8      []uint8
	AUint16     []uint16
	AUint32     []uint32
	AUint64     []uint64
	AFloat32    []float32
	AFloat64    []float64
	AString     []string
	AnInterface []interface{}
}

func TestBuilderSliceOfStructs(t *testing.T) {
	v := []BuilderSliceTestStruct{
		{
			AnInt: []int{1},
		},
		{
			AnInt: []int{1},
		},
	}

	assertBuild(t,
		v,
		L(),
		M(),
		S("AnInt"), L(), N(1), E(),
		E(),
		M(),
		S("AnInt"), L(), N(1), E(),
		E(),
		E())
}

type SimpleTestStruct struct {
	IValue int
}

func TestBuilderListOfStruct(t *testing.T) {
	v := []*SimpleTestStruct{
		{
			IValue: 5,
		},
	}

	assertBuild(t,
		v,
		L(),
		M(),
		S("IValue"),
		N(5),
		E(),
		E())
}

type NullContainers struct {
	Bytes []byte
	Slice []interface{}
	Map   map[interface{}]interface{}
}

func TestBuilderNullContainers(t *testing.T) {
	v := NullContainers{}

	assertBuild(t, v,
		M(),
		S("Bytes"),
		NULL(),
		S("Slice"),
		NULL(),
		S("Map"),
		NULL(),
		E())
}

type PURLContainer struct {
	URL *url.URL
}

func TestBuilderPURLContainer(t *testing.T) {
	v := PURLContainer{NewRID("http://x.com")}

	assertBuild(t, v,
		M(),
		S("URL"),
		RID("http://x.com"),
		E())
}

func TestBuilderNullPURLContainer(t *testing.T) {
	v := PURLContainer{}

	assertBuild(t, v,
		M(),
		S("URL"),
		NULL(),
		E())
}

func TestBuilderByteArrayBytes(t *testing.T) {
	assertBuild(t, [2]byte{1, 2},
		AU8([]byte{1, 2}))
}

func TestBuilderMarkerSlice(t *testing.T) {
	// Reference in same container
	assertBuild(t, []int{100, 100}, L(), MARK("1"), N(100), REFL("1"), E())
	assertBuild(t, []int{100, 100}, L(), REFL("1"), MARK("1"), N(100), E())
	assertBuild(t, []string{"abcdef", "abcdef"}, L(), MARK("1"), S("abcdef"), REFL("1"), E())
	assertBuild(t, [][]int{{100}, {100}}, L(), MARK("1"), L(), N(100), E(), REFL("1"), E())
	assertBuild(t, [][]int{{100}, {100}}, L(), REFL("1"), MARK("1"), L(), N(100), E(), E())
	assertBuild(t, [][]int{{100, 100}, {100, 100}}, L(), REFL("1"), MARK("1"), L(), N(100), N(100), E(), E())

	// Reference in different container
	assertBuild(t, [][]int{{100, 100}, {100}}, L(), L(), REFL("1"), REFL("1"), E(), L(), MARK("1"), N(100), E(), E())

	// Referenced containers
	assertBuild(t, [][]int{{}, {}}, L(), MARK("1"), L(), E(), REFL("1"), E())
	assertBuild(t, [][]int{{}, {}}, L(), REFL("1"), MARK("1"), L(), E(), E())
	assertBuild(t, []map[int]int{{100: 100}, {100: 100}},
		L(), REFL("1"), MARK("1"), M(), N(100), N(100), E(), E())
	assertBuild(t, []map[int]int{{100: 100}, {100: 100}},
		L(), MARK("1"), M(), N(100), N(100), E(), REFL("1"), E())

	// Interface
	assertBuild(t, []interface{}{100, 100}, L(), REFL("1"), MARK("1"), N(100), E())
	assertBuild(t, []interface{}{100, 100}, L(), MARK("1"), N(100), REFL("1"), E())

	// Recursive interface
	rintf := make([]interface{}, 1)
	rintf[0] = rintf
	assertBuild(t, rintf, MARK("1"), L(), REFL("1"), E())
}

func TestBuilderMarkerArray(t *testing.T) {
	assertBuild(t, [2]int{100, 100}, L(), MARK("1"), N(100), REFL("1"), E())
	assertBuild(t, [2]int{100, 100}, L(), REFL("1"), MARK("1"), N(100), E())
}

func TestBuilderMarkerMap(t *testing.T) {
	assertBuild(t, map[int]int8{1: 100, 2: 100, 3: 100},
		M(), N(1), REFL("5"), N(2), MARK("5"), N(100), N(3), REFL("5"), E())

	rmap := make(map[int]interface{})
	rmap[0] = rmap
	assertBuild(t, rmap,
		MARK("1"), M(), N(0), REFL("1"), E())

}

type SelfReferential struct {
	Value int
	Next  *SelfReferential
}

func TestBuilderSelfReferential(t *testing.T) {
	v := &SelfReferential{
		Value: 100,
	}
	v.Next = v
	assertBuild(t, v,
		MARK("1"),
		M(),
		S("Value"),
		N(100),
		S("Next"),
		REFL("1"),
		E())
}

type RefStruct struct {
	I16 int16
	F32 float32
	S1  string
	S2  string
}

func TestBuilderRefStruct(t *testing.T) {
	assertBuild(t, &RefStruct{
		I16: 1000,
		F32: 1000,
		S1:  "test",
		S2:  "test",
	}, M(),
		S("I16"), REFL("1"),
		S("F32"), MARK("1"), N(1000),
		S("S1"), MARK("2"), S("test"),
		S("S2"), REFL("2"),
		E())
}

type TagStruct struct {
	Omit1 string `ce:"-"`
	Omit2 string `ce:"omit"`
	Named string `ce:"name=test"`
}

type TagStruct2 struct {
	Omit1 string `ce:"  - "`
	Omit2 string `ce:"  omit "`
	Named string `ce:" name = test "`
}

func TestBuilderStructTags(t *testing.T) {
	assertBuild(t, &TagStruct{
		Named: "Something",
	}, M(),
		S("test"), S("Something"),
		E())

	assertBuild(t, &TagStruct2{
		Named: "Something",
	}, M(),
		S("test"), S("Something"),
		E())
}

func TestBuilderUID(t *testing.T) {
	uid := types.UID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	assertBuild(t, uid, UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}))
	list := []types.UID{uid}
	assertBuild(t, list, L(), UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}), E())
	m := map[int]types.UID{1: uid}
	assertBuild(t, m, M(), N(1), UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}), E())
}

func TestBuilderMedia(t *testing.T) {
	media := types.Media{
		MediaType: "a",
		Data:      []byte{1},
	}

	assertBuild(t, media, BMEDIA(), ACL(1), ADT("a"), ACL(1), ADU8([]byte{1}))
	list := []types.Media{media}
	assertBuild(t, list, L(), BMEDIA(), ACL(1), ADT("a"), ACL(1), ADU8([]byte{1}), E())
	m := map[int]types.Media{1: media}
	assertBuild(t, m, M(), N(1), BMEDIA(), ACL(1), ADT("a"), ACL(1), ADU8([]byte{1}), E())
}

func TestBuilderEdge(t *testing.T) {
	r := types.Edge{
		Source:      NewRID("http://x.com"),
		Description: NewRID("http://y.com"),
		Destination: NewRID("http://z.com"),
	}
	pr := &r

	assertBuild(t, r, EDGE(), RID("http://x.com"), RID("http://y.com"), RID("http://z.com"))
	assertBuild(t, pr, EDGE(), RID("http://x.com"), RID("http://y.com"), RID("http://z.com"))

	r.Source = []interface{}{NewRID("a"), NewRID("b")}
	r.Destination = 1
	assertBuild(t, r, EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), N(1))
	assertBuild(t, pr, EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), N(1))

	assertBuild(t, []interface{}{pr}, L(), EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), N(1), E())
}

type MyStruct struct {
	Aaa int
	Bbb []int
}

func TestBuilderStructBadFieldName(t *testing.T) {
	s := &MyStruct{}
	assertBuild(t, s, M(), S("aaa"), N(0), S("bbb"), L(), E(), E())
	assertBuild(t, s, M(), S("x"), N(0), S("bbb"), L(), E(), E())
}
