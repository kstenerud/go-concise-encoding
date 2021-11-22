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
	assertBuild(t, [2]uint16{0x1234, 0x2345}, L(), I(0x1234), I(0x2345), E())
	assertBuild(t, []uint16{0x12, 0x23}, L(), I(0x12), I(0x23), E())
}

func TestBuilderTypedArrayUint16(t *testing.T) {
	assertBuild(t, [2]uint16{0x1234, 0x2345}, AU16([]uint16{0x1234, 0x2345}))
	assertBuild(t, []uint16{0x1234, 0x2345}, AU16([]uint16{0x1234, 0x2345}))
	assertBuild(t, [2]uint16{0x1234, 0x2345}, L(), I(0x1234), I(0x2345), E())
	assertBuild(t, []uint16{0x1234, 0x2345}, L(), I(0x1234), I(0x2345), E())
}

func TestBuilderTypedArrayUint32(t *testing.T) {
	assertBuild(t, [2]uint32{0x12345678, 0x23456789}, AU32([]uint32{0x12345678, 0x23456789}))
	assertBuild(t, []uint32{0x12345678, 0x23456789}, AU32([]uint32{0x12345678, 0x23456789}))
	assertBuild(t, [2]uint32{0x12345678, 0x23456789}, L(), I(0x12345678), I(0x23456789), E())
	assertBuild(t, []uint32{0x12345678, 0x23456789}, L(), I(0x12345678), I(0x23456789), E())
}

func TestBuilderTypedArrayUint64(t *testing.T) {
	assertBuild(t, [2]uint64{0x123456789abcdef0, 0x23456789abcdef01}, AU64([]uint64{0x123456789abcdef0, 0x23456789abcdef01}))
	assertBuild(t, []uint64{0x123456789abcdef0, 0x23456789abcdef01}, AU64([]uint64{0x123456789abcdef0, 0x23456789abcdef01}))
	assertBuild(t, [2]uint64{0x123456789abcdef0, 0x23456789abcdef01}, L(), PI(0x123456789abcdef0), PI(0x23456789abcdef01), E())
	assertBuild(t, []uint64{0x123456789abcdef0, 0x23456789abcdef01}, L(), PI(0x123456789abcdef0), PI(0x23456789abcdef01), E())
}

func TestBuilderTypedArrayInt8(t *testing.T) {
	assertBuild(t, [2]int8{0x12, -0x23}, AI8([]int8{0x12, -0x23}))
	assertBuild(t, []int8{0x12, -0x23}, AI8([]int8{0x12, -0x23}))
	assertBuild(t, [2]int8{0x12, -0x23}, L(), I(0x12), I(-0x23), E())
	assertBuild(t, []int8{0x12, -0x23}, L(), I(0x12), I(-0x23), E())
}

func TestBuilderTypedArrayInt16(t *testing.T) {
	assertBuild(t, [2]int16{0x1234, -0x2345}, AI16([]int16{0x1234, -0x2345}))
	assertBuild(t, []int16{0x1234, -0x2345}, AI16([]int16{0x1234, -0x2345}))
	assertBuild(t, [2]int16{0x1234, -0x2345}, L(), I(0x1234), I(-0x2345), E())
	assertBuild(t, []int16{0x1234, -0x2345}, L(), I(0x1234), I(-0x2345), E())
}

func TestBuilderTypedArrayInt32(t *testing.T) {
	assertBuild(t, [2]int32{0x12345678, -0x23456789}, AI32([]int32{0x12345678, -0x23456789}))
	assertBuild(t, []int32{0x12345678, -0x23456789}, AI32([]int32{0x12345678, -0x23456789}))
	assertBuild(t, [2]int32{0x12345678, -0x23456789}, L(), I(0x12345678), I(-0x23456789), E())
	assertBuild(t, []int32{0x12345678, -0x23456789}, L(), I(0x12345678), I(-0x23456789), E())
}

func TestBuilderTypedArrayInt64(t *testing.T) {
	assertBuild(t, [2]int64{0x123456789abcdef0, -0x23456789abcdef01}, AI64([]int64{0x123456789abcdef0, -0x23456789abcdef01}))
	assertBuild(t, []int64{0x123456789abcdef0, -0x23456789abcdef01}, AI64([]int64{0x123456789abcdef0, -0x23456789abcdef01}))
	assertBuild(t, [2]int64{0x123456789abcdef0, -0x23456789abcdef01}, L(), PI(0x123456789abcdef0), I(-0x23456789abcdef01), E())
	assertBuild(t, []int64{0x123456789abcdef0, -0x23456789abcdef01}, L(), PI(0x123456789abcdef0), I(-0x23456789abcdef01), E())
}

func TestBuilderTypedArrayFloat32(t *testing.T) {
	assertBuild(t, [2]float32{-1.25, 9.5e10}, AF32([]float32{-1.25, 9.5e10}))
	assertBuild(t, []float32{-1.25, 9.5e10}, AF32([]float32{-1.25, 9.5e10}))
	assertBuild(t, [2]float32{-1.25, 9.5e10}, L(), F(-1.25), F(9.5e10), E())
	assertBuild(t, []float32{-1.25, 9.5e10}, L(), F(-1.25), F(9.5e10), E())
}

func TestBuilderTypedArrayFloat64(t *testing.T) {
	assertBuild(t, [2]float64{-1.25, 9.5e10}, AF64([]float64{-1.25, 9.5e10}))
	assertBuild(t, []float64{-1.25, 9.5e10}, AF64([]float64{-1.25, 9.5e10}))
	assertBuild(t, [2]float64{-1.25, 9.5e10}, L(), F(-1.25), F(9.5e10), E())
	assertBuild(t, []float64{-1.25, 9.5e10}, L(), F(-1.25), F(9.5e10), E())
}

// =================================================

func TestBuildUnknown(t *testing.T) {
	expected := []interface{}{1}
	actual := runBuild(NewSession(nil, nil), nil, L(), I(1), E())

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

var zero = float64(0)
var negZero = -zero

func TestBuilderBasicTypes(t *testing.T) {
	pBigIntP := NewBigInt("12345678901234567890123456789")
	pBigIntN := NewBigInt("-999999999999999999999999999999")
	pBigFloat := NewBigFloat("1.2345678901234567890123456789e10000")
	dfloat := NewDFloat("1.23456e1000")
	pBigDFloat := NewBDF("4.509e10000")
	gTimeNow := time.Now()
	cTimeNow := test.AsCompactTime(gTimeNow)
	cTime := test.NewTimeLL(10, 5, 59, 100, 506, 107)
	pURL := NewRID("http://x.com")
	pNode := NewNode("test", []interface{}{"a"})
	pEdge := NewEdge("a", "b", "c")

	assertBuild(t, true, B(true))
	assertBuild(t, false, B(false))
	assertBuild(t, int(10), I(10))
	assertBuild(t, int8(10), I(10))
	assertBuild(t, int16(-10), I(-10))
	assertBuild(t, int32(10), I(10))
	assertBuild(t, int64(-10), I(-10))
	assertBuild(t, uint(10), I(10))
	assertBuild(t, uint8(10), I(10))
	assertBuild(t, uint16(10), I(10))
	assertBuild(t, uint32(10), I(10))
	assertBuild(t, uint64(10), I(10))
	assertBuild(t, 1, I(1))
	assertBuild(t, -1, I(-1))
	assertBuild(t, pBigIntP, BI(pBigIntP))
	assertBuild(t, *pBigIntP, BI(pBigIntP))
	assertBuild(t, pBigIntN, BI(pBigIntN))
	assertBuild(t, *pBigIntN, BI(pBigIntN))
	assertBuild(t, (*big.Int)(nil), N())
	assertBuild(t, float32(negZero), F(negZero))
	assertBuild(t, float32(-1.25), F(-1.25))
	assertBuild(t, float64(-9.5e50), F(-9.5e50))
	assertBuild(t, pBigFloat, BF(pBigFloat))
	assertBuild(t, *pBigFloat, BF(pBigFloat))
	assertBuild(t, (*big.Float)(nil), N())
	assertBuild(t, dfloat, DF(dfloat))
	assertBuild(t, pBigDFloat, BDF(pBigDFloat))
	assertBuild(t, *pBigDFloat, BDF(pBigDFloat))
	assertBuild(t, (*apd.Decimal)(nil), N())
	assertBuild(t, common.SignalingNan, SNAN())
	assertBuild(t, common.QuietNan, NAN())
	assertBuild(t, gTimeNow, GT(gTimeNow))
	assertBuild(t, cTimeNow, CT(cTimeNow))
	assertBuild(t, cTime, CT(cTime))
	assertBuild(t, []byte{1, 2, 3, 4}, AU8([]byte{1, 2, 3, 4}))
	assertBuild(t, "test", S("test"))
	assertBuild(t, pURL, RID("http://x.com"))
	assertBuild(t, *pURL, RID("http://x.com"))
	assertBuild(t, (*url.URL)(nil), N())
	assertBuild(t, interface{}(1234), I(1234))
	assertBuild(t, pNode, NODE(), S("test"), S("a"), E())
	assertBuild(t, *pNode, NODE(), S("test"), S("a"), E())
	assertBuild(t, pEdge, EDGE(), S("a"), S("b"), S("c"))
	assertBuild(t, *pEdge, EDGE(), S("a"), S("b"), S("c"))
}

func TestBuilderConvertToBDF(t *testing.T) {
	pv := NewBDF("1")
	nv := NewBDF("-1")
	nz := NewBDF("-0")

	assertBuild(t, pv, PI(1))
	assertBuild(t, nv, NI(1))
	assertBuild(t, pv, BI(NewBigInt("1")))
	assertBuild(t, pv, F(1))
	assertBuild(t, pv, BF(NewBigFloat("1")))
	assertBuild(t, pv, DF(NewDFloat("1")))
	assertBuild(t, pv, BDF(NewBDF("1")))

	assertBuild(t, *pv, PI(1))
	assertBuild(t, *nv, NI(1))
	assertBuild(t, *pv, BI(NewBigInt("1")))
	assertBuild(t, *pv, F(1))
	assertBuild(t, *pv, BF(NewBigFloat("1")))
	assertBuild(t, *pv, DF(NewDFloat("1")))
	assertBuild(t, *pv, BDF(NewBDF("1")))

	assertBuild(t, nz, F(negZero))
	assertBuild(t, nz, BF(NewBigFloat("-0")))
	assertBuild(t, nz, DF(NewDFloat("-0")))
	assertBuild(t, nz, BDF(NewBDF("-0")))
}

func TestBuilderConvertToBDFFail(t *testing.T) {
	v := NewBDF("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, GT(time.Now()))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, N())
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, GT(time.Now()))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToBF(t *testing.T) {
	pv := NewBigFloat("1")
	nv := NewBigFloat("-1")
	nz := NewBigFloat("-0")

	assertBuild(t, pv, PI(1))
	assertBuild(t, nv, NI(1))
	assertBuild(t, pv, BI(NewBigInt("1")))
	assertBuild(t, pv, F(1))
	assertBuild(t, pv, BF(NewBigFloat("1")))
	assertBuild(t, pv, DF(NewDFloat("1")))
	assertBuild(t, pv, BDF(NewBDF("1")))

	assertBuild(t, *pv, PI(1))
	assertBuild(t, *nv, NI(1))
	assertBuild(t, *pv, BI(NewBigInt("1")))
	assertBuild(t, *pv, F(1))
	assertBuild(t, *pv, BF(NewBigFloat("1")))
	assertBuild(t, *pv, DF(NewDFloat("1")))
	assertBuild(t, *pv, BDF(NewBDF("1")))

	assertBuild(t, *nz, F(negZero))
	assertBuild(t, *nz, BF(NewBigFloat("-0")))
	assertBuild(t, *nz, DF(NewDFloat("-0")))
	assertBuild(t, *nz, BDF(NewBDF("-0")))
}

func TestBuilderConvertToBFFail(t *testing.T) {
	v := NewBigFloat("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, GT(time.Now()))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, N())
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, GT(time.Now()))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToBI(t *testing.T) {
	pv := NewBigInt("1")
	nv := NewBigInt("-1")
	assertBuild(t, pv, PI(1))
	assertBuild(t, NewBigInt("9223372036854775808"), PI(9223372036854775808))
	assertBuild(t, nv, NI(1))
	assertBuild(t, pv, BI(NewBigInt("1")))
	assertBuild(t, pv, F(1))
	assertBuild(t, pv, BF(NewBigFloat("1")))
	assertBuild(t, pv, DF(NewDFloat("1")))
	assertBuild(t, pv, BDF(NewBDF("1")))

	assertBuild(t, *pv, PI(1))
	assertBuild(t, *nv, NI(1))
	assertBuild(t, *pv, BI(NewBigInt("1")))
	assertBuild(t, *pv, F(1))
	assertBuild(t, *pv, BF(NewBigFloat("1")))
	assertBuild(t, *pv, DF(NewDFloat("1")))
	assertBuild(t, *pv, BDF(NewBDF("1")))
}

func TestBuilderConvertToBIFail(t *testing.T) {
	v := NewBigInt("1")
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, F(1.1))
	assertBuildPanics(t, v, BF(NewBigFloat("1.1")))
	assertBuildPanics(t, v, BF(NewBigFloat("1.0e100000")))
	assertBuildPanics(t, v, DF(NewDFloat("1.1")))
	assertBuildPanics(t, v, DF(NewDFloat("1.0e100000")))
	assertBuildPanics(t, v, BDF(NewBDF("1.1")))
	assertBuildPanics(t, v, BDF(NewBDF("1.0e100000")))
	assertBuildPanics(t, v, BDF(NewBDF("nan")))
	assertBuildPanics(t, v, BDF(NewBDF("snan")))
	assertBuildPanics(t, v, BDF(NewBDF("inf")))
	assertBuildPanics(t, v, BDF(NewBDF("-inf")))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, GT(time.Now()))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())

	assertBuildPanics(t, *v, N())
	assertBuildPanics(t, *v, F(1.1))
	assertBuildPanics(t, *v, B(true))
	assertBuildPanics(t, *v, BF(NewBigFloat("1.1")))
	assertBuildPanics(t, *v, BF(NewBigFloat("1.0e100000")))
	assertBuildPanics(t, *v, DF(NewDFloat("1.1")))
	assertBuildPanics(t, *v, DF(NewDFloat("1.0e100000")))
	assertBuildPanics(t, *v, BDF(NewBDF("1.1")))
	assertBuildPanics(t, *v, BDF(NewBDF("1.0e100000")))
	assertBuildPanics(t, *v, BDF(NewBDF("nan")))
	assertBuildPanics(t, *v, BDF(NewBDF("snan")))
	assertBuildPanics(t, *v, BDF(NewBDF("inf")))
	assertBuildPanics(t, *v, BDF(NewBDF("-inf")))
	assertBuildPanics(t, *v, S("1"))
	assertBuildPanics(t, *v, AU8([]byte{1}))
	assertBuildPanics(t, *v, RID("x://x"))
	assertBuildPanics(t, *v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, GT(time.Now()))
	assertBuildPanics(t, *v, L())
	assertBuildPanics(t, *v, M())
	assertBuildPanics(t, *v, E())
}

func TestBuilderConvertToDecimalFloat(t *testing.T) {
	pv := NewDFloat("1")
	nv := NewDFloat("-1")
	nz := NewDFloat("-0")

	assertBuild(t, pv, PI(1))
	assertBuild(t, nv, NI(1))
	assertBuild(t, pv, BI(NewBigInt("1")))
	assertBuild(t, pv, F(1))
	assertBuild(t, pv, BF(NewBigFloat("1")))
	assertBuild(t, pv, DF(NewDFloat("1")))
	assertBuild(t, pv, BDF(NewBDF("1")))
	assertBuild(t, nz, F(negZero))
	assertBuild(t, nz, BF(NewBigFloat("-0")))
	assertBuild(t, nz, DF(NewDFloat("-0")))
	assertBuild(t, nz, BDF(NewBDF("-0")))
}

func TestBuilderDecimalFloatFail(t *testing.T) {
	v := NewDFloat("1")
	assertBuildPanics(t, v, N())
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, GT(time.Now()))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())
}

func TestBuilderConvertToFloat(t *testing.T) {
	pv := 1.0
	nv := -1.0
	nz := negZero

	assertBuild(t, pv, PI(1))
	assertBuild(t, nv, NI(1))
	assertBuild(t, pv, BI(NewBigInt("1")))
	assertBuild(t, pv, F(1))
	assertBuild(t, pv, BF(NewBigFloat("1")))
	assertBuild(t, pv, DF(NewDFloat("1")))
	assertBuild(t, pv, BDF(NewBDF("1")))

	assertBuild(t, nz, F(negZero))
	assertBuild(t, nz, BF(NewBigFloat("-0")))
	assertBuild(t, nz, DF(NewDFloat("-0")))
	assertBuild(t, nz, BDF(NewBDF("-0")))
}

func TestBuilderConvertToFloatFail(t *testing.T) {
	// TODO: How to define required conversion accuracy?
	v := 1.0
	assertBuildPanics(t, v, N())
	assertBuildPanics(t, v, B(true))
	assertBuildPanics(t, v, PI(0xffffffffffffffff))
	assertBuildPanics(t, v, I(-0x7fffffffffffffff))
	assertBuildPanics(t, v, BF(NewBigFloat("1.0e309")))
	assertBuildPanics(t, v, BF(NewBigFloat("1.0e-311")))
	// TODO: apd.Decimal and compact_float.DFloat don't handle float overflow
	assertBuildPanics(t, v, BI(NewBigInt("1234567890123456789012345")))
	assertBuildPanics(t, v, S("1"))
	assertBuildPanics(t, v, AU8([]byte{1}))
	assertBuildPanics(t, v, RID("x://x"))
	assertBuildPanics(t, v, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, GT(time.Now()))
	assertBuildPanics(t, v, L())
	assertBuildPanics(t, v, M())
	assertBuildPanics(t, v, E())
}

func TestBuilderConvertToInt(t *testing.T) {
	assertBuild(t, 1, PI(1))
	assertBuild(t, -1, NI(1))
	assertBuild(t, 1, BI(NewBigInt("1")))
	assertBuild(t, 1, F(1))
	assertBuild(t, 1, BF(NewBigFloat("1")))
	assertBuild(t, 1, DF(NewDFloat("1")))
	assertBuild(t, 1, BDF(NewBDF("1")))
}

func TestBuilderConvertToIntFail(t *testing.T) {
	assertBuildPanics(t, int(1), N())
	assertBuildPanics(t, int(1), B(true))
	assertBuildPanics(t, int(1), PI(0x8000000000000000))
	assertBuildPanics(t, int(1), F(1.1))
	assertBuildPanics(t, int(1), BF(NewBigFloat("1.1")))
	assertBuildPanics(t, int(1), DF(NewDFloat("1.1")))
	assertBuildPanics(t, int(1), BDF(NewBDF("1.1")))
	assertBuildPanics(t, int(1), BF(NewBigFloat("1e20")))
	assertBuildPanics(t, int(1), DF(NewDFloat("1e20")))
	assertBuildPanics(t, int(1), BDF(NewBDF("1e20")))
	assertBuildPanics(t, int(1), BF(NewBigFloat("-1e20")))
	assertBuildPanics(t, int(1), DF(NewDFloat("-1e20")))
	assertBuildPanics(t, int(1), BDF(NewBDF("-1e20")))
	assertBuildPanics(t, int(1), BI(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, int(1), BI(NewBigInt("-100000000000000000000")))
	assertBuildPanics(t, int(1), S("1"))
	assertBuildPanics(t, int(1), AU8([]byte{1}))
	assertBuildPanics(t, int(1), RID("x://x"))
	assertBuildPanics(t, int(1), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, int(1), CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, int(1), GT(time.Now()))
	assertBuildPanics(t, int(1), L())
	assertBuildPanics(t, int(1), M())
	assertBuildPanics(t, int(1), E())

	assertBuildPanics(t, int8(1), PI(1000))
	assertBuildPanics(t, int8(1), NI(1000))
	assertBuildPanics(t, int8(1), I(1000))
	assertBuildPanics(t, int8(1), BI(NewBigInt("1000")))
	assertBuildPanics(t, int8(1), F(1000))
	assertBuildPanics(t, int8(1), BF(NewBigFloat("1000")))
	assertBuildPanics(t, int8(1), DF(NewDFloat("1000")))
	assertBuildPanics(t, int8(1), BDF(NewBDF("1000")))

	assertBuildPanics(t, int16(1), PI(100000))
	assertBuildPanics(t, int16(1), NI(100000))
	assertBuildPanics(t, int16(1), I(100000))
	assertBuildPanics(t, int16(1), BI(NewBigInt("100000")))
	assertBuildPanics(t, int16(1), F(100000))
	assertBuildPanics(t, int16(1), BF(NewBigFloat("100000")))
	assertBuildPanics(t, int16(1), DF(NewDFloat("100000")))
	assertBuildPanics(t, int16(1), BDF(NewBDF("100000")))

	assertBuildPanics(t, int32(1), PI(10000000000))
	assertBuildPanics(t, int32(1), NI(10000000000))
	assertBuildPanics(t, int32(1), I(10000000000))
	assertBuildPanics(t, int32(1), BI(NewBigInt("10000000000")))
	assertBuildPanics(t, int32(1), F(10000000000))
	assertBuildPanics(t, int32(1), BF(NewBigFloat("10000000000")))
	assertBuildPanics(t, int32(1), DF(NewDFloat("10000000000")))
	assertBuildPanics(t, int32(1), BDF(NewBDF("10000000000")))
}

func TestBuilderConvertToUint(t *testing.T) {
	assertBuild(t, uint(1), PI(1))
	assertBuild(t, uint(1), I(1))
	assertBuild(t, uint(1), BI(NewBigInt("1")))
	assertBuild(t, uint(1), F(1))
	assertBuild(t, uint(1), BF(NewBigFloat("1")))
	assertBuild(t, uint(1), DF(NewDFloat("1")))
	assertBuild(t, uint(1), BDF(NewBDF("1")))
}

func TestBuilderConvertToUintFail(t *testing.T) {
	assertBuildPanics(t, uint(1), N())
	assertBuildPanics(t, uint(1), B(true))
	assertBuildPanics(t, uint(1), NI(1))
	assertBuildPanics(t, uint(1), F(1.1))
	assertBuildPanics(t, uint(1), BF(NewBigFloat("1.1")))
	assertBuildPanics(t, uint(1), DF(NewDFloat("1.1")))
	assertBuildPanics(t, uint(1), BDF(NewBDF("1.1")))
	assertBuildPanics(t, uint(1), BF(NewBigFloat("1e20")))
	assertBuildPanics(t, uint(1), DF(NewDFloat("1e20")))
	assertBuildPanics(t, uint(1), BDF(NewBDF("1e20")))
	assertBuildPanics(t, uint8(1), BI(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, uint(1), S("1"))
	assertBuildPanics(t, uint(1), AU8([]byte{1}))
	assertBuildPanics(t, uint(1), RID("x://x"))
	assertBuildPanics(t, uint(1), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, uint(1), CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, uint(1), GT(time.Now()))
	assertBuildPanics(t, uint(1), L())
	assertBuildPanics(t, uint(1), M())
	assertBuildPanics(t, uint(1), E())

	assertBuildPanics(t, uint8(1), PI(1000))
	assertBuildPanics(t, uint8(1), I(1000))
	assertBuildPanics(t, uint8(1), BI(NewBigInt("1000")))
	assertBuildPanics(t, uint8(1), F(1000))
	assertBuildPanics(t, uint8(1), BF(NewBigFloat("1000")))
	assertBuildPanics(t, uint8(1), DF(NewDFloat("1000")))
	assertBuildPanics(t, uint8(1), BDF(NewBDF("1000")))

	assertBuildPanics(t, uint16(1), PI(100000))
	assertBuildPanics(t, uint16(1), I(100000))
	assertBuildPanics(t, uint16(1), BI(NewBigInt("100000")))
	assertBuildPanics(t, uint16(1), F(100000))
	assertBuildPanics(t, uint16(1), BF(NewBigFloat("100000")))
	assertBuildPanics(t, uint16(1), DF(NewDFloat("100000")))
	assertBuildPanics(t, uint16(1), BDF(NewBDF("100000")))

	assertBuildPanics(t, uint32(1), PI(10000000000))
	assertBuildPanics(t, uint32(1), I(10000000000))
	assertBuildPanics(t, uint32(1), BI(NewBigInt("10000000000")))
	assertBuildPanics(t, uint32(1), F(10000000000))
	assertBuildPanics(t, uint32(1), BF(NewBigFloat("10000000000")))
	assertBuildPanics(t, uint32(1), DF(NewDFloat("10000000000")))
	assertBuildPanics(t, uint32(1), BDF(NewBDF("10000000000")))
}

func TestBuilderString(t *testing.T) {
	assertBuild(t, "", N())
	assertBuild(t, "test", S("test"))
}

func TestBuilderStringFail(t *testing.T) {
	assertBuildPanics(t, "", B(false))
	assertBuildPanics(t, "", PI(1))
	assertBuildPanics(t, "", NI(1))
	assertBuildPanics(t, "", F(1.1))
	assertBuildPanics(t, "", BF(NewBigFloat("1.1")))
	assertBuildPanics(t, "", DF(NewDFloat("1.1")))
	assertBuildPanics(t, "", BDF(NewBDF("1.1")))
	assertBuildPanics(t, "", BF(NewBigFloat("1e20")))
	assertBuildPanics(t, "", DF(NewDFloat("1e20")))
	assertBuildPanics(t, "", BDF(NewBDF("1e20")))
	assertBuildPanics(t, "", BI(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, "", AU8([]byte{1}))
	assertBuildPanics(t, "", RID("x://x"))
	assertBuildPanics(t, "", UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, "", CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, "", GT(time.Now()))
	assertBuildPanics(t, "", L())
	assertBuildPanics(t, "", M())
	assertBuildPanics(t, "", E())
}

func TestBuilderChunkedBytes(t *testing.T) {
	assertBuild(t, []byte{1, 2, 3, 4}, AU8B(), AC(2, true), AD([]byte{1, 2}), AC(2, false), AD([]byte{3, 4}))
	assertBuild(t, []byte{1, 2, 3, 4}, AU8B(), AC(2, true), AD([]byte{1, 2}), AC(2, true), AD([]byte{3, 4}), AC(0, false))
}

func TestBuilderChunkedString(t *testing.T) {
	assertBuild(t, "test", SB(), AC(2, true), AD([]byte("te")), AC(2, false), AD([]byte("st")))
	assertBuild(t, "test", SB(), AC(2, true), AD([]byte("te")), AC(2, true), AD([]byte("st")), AC(0, false))
}

func TestBuilderChunkedRID(t *testing.T) {
	expected := NewRID("test")
	assertBuild(t, expected, RB(), AC(2, true), AD([]byte("te")), AC(2, false), AD([]byte("st")))
	assertBuild(t, expected, RB(), AC(2, true), AD([]byte("te")), AC(2, true), AD([]byte("st")), AC(0, false))
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
	session := NewSession(nil, opts)
	expected := CustomBinaryExampleType(0x01020304)
	assertBuildWithSession(t, session, expected, CBB(), AC(2, true), AD([]byte{1, 2}), AC(2, false), AD([]byte{3, 4}))
	assertBuildWithSession(t, session, expected, CBB(), AC(2, true), AD([]byte{1, 2}), AC(2, true), AD([]byte{3, 4}), AC(0, false))
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
	session := NewSession(nil, opts)
	expected := CustomTextExampleType(0x1234)
	assertBuildWithSession(t, session, expected, CTB(), AC(2, true), AD([]byte{'1', '2'}), AC(2, false), AD([]byte{'3', '4'}))
	assertBuildWithSession(t, session, expected, CTB(), AC(2, true), AD([]byte{'1', '2'}), AC(2, true), AD([]byte{'3', '4'}), AC(0, false))
}

func TestBuilderGoTime(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := test.AsCompactTime(gtime)

	assertBuild(t, gtime, GT(gtime))
	assertBuild(t, gtime, CT(ctime))
}

func TestBuilderGoTimeFail(t *testing.T) {
	gtime := time.Time{}
	ctime := test.NewTimeLL(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, gtime, N())
	assertBuildPanics(t, gtime, B(true))
	assertBuildPanics(t, gtime, PI(1))
	assertBuildPanics(t, gtime, NI(1))
	assertBuildPanics(t, gtime, F(1.1))
	assertBuildPanics(t, gtime, BF(NewBigFloat("1.1")))
	assertBuildPanics(t, gtime, DF(NewDFloat("1.1")))
	assertBuildPanics(t, gtime, BDF(NewBDF("1.1")))
	assertBuildPanics(t, gtime, BF(NewBigFloat("1e20")))
	assertBuildPanics(t, gtime, DF(NewDFloat("1e20")))
	assertBuildPanics(t, gtime, BDF(NewBDF("1e20")))
	assertBuildPanics(t, gtime, BI(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, gtime, S("1"))
	assertBuildPanics(t, gtime, AU8([]byte{1}))
	assertBuildPanics(t, gtime, RID("x://x"))
	assertBuildPanics(t, gtime, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, gtime, CT(ctime))
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

	assertBuild(t, ctime, GT(gtime))
	assertBuild(t, ctime, CT(ctime))
}

func TestBuilderCompactTimeFail(t *testing.T) {
	ctime := test.NewTimeLL(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, ctime, B(true))
	assertBuildPanics(t, ctime, PI(1))
	assertBuildPanics(t, ctime, NI(1))
	assertBuildPanics(t, ctime, F(1.1))
	assertBuildPanics(t, ctime, BF(NewBigFloat("1.1")))
	assertBuildPanics(t, ctime, DF(NewDFloat("1.1")))
	assertBuildPanics(t, ctime, BDF(NewBDF("1.1")))
	assertBuildPanics(t, ctime, BF(NewBigFloat("1e20")))
	assertBuildPanics(t, ctime, DF(NewDFloat("1e20")))
	assertBuildPanics(t, ctime, BDF(NewBDF("1e20")))
	assertBuildPanics(t, ctime, BI(NewBigInt("100000000000000000000")))
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
	assertBuild(t, []int8{-1, 2, 3, 4, 5, 6, 7}, L(), I(-1), PI(2), F(3),
		BI(NewBigInt("4")), BF(NewBigFloat("5")), DF(NewDFloat("6")),
		BDF(NewBDF("7")), E())
	assertBuild(t, []*int{nil}, L(), N(), E())
	assertBuild(t, []string{"test"}, L(), S("test"), E())
	assertBuild(t, [][]byte{[]byte{1}}, L(), AU8([]byte{1}), E())
	assertBuild(t, []*url.URL{NewRID("http://example.com")}, L(), RID("http://example.com"), E())
	assertBuild(t, []time.Time{gtime}, L(), GT(gtime), E())
	assertBuild(t, []compact_time.Time{ctime}, L(), CT(ctime), E())
	assertBuild(t, [][]int{[]int{1}}, L(), L(), I(1), E(), E())
	assertBuild(t, []map[int]int{map[int]int{1: 2}}, L(), M(), I(1), I(2), E(), E())
}

func TestBuilderSliceFail(t *testing.T) {
	assertBuildPanics(t, []int{}, N())
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
	assertBuild(t, [7]int8{-1, 2, 3, 4, 5, 6, 7}, L(), I(-1), PI(2), F(3),
		BI(NewBigInt("4")), BF(NewBigFloat("5")), DF(NewDFloat("6")),
		BDF(NewBDF("7")), E())
	assertBuild(t, [1]*int{nil}, L(), N(), E())
	assertBuild(t, [1]string{"test"}, L(), S("test"), E())
	assertBuild(t, [1][]byte{[]byte{1}}, L(), AU8([]byte{1}), E())
	assertBuild(t, [1]*url.URL{NewRID("http://example.com")}, L(), RID("http://example.com"), E())
	assertBuild(t, [1]time.Time{gtime}, L(), GT(gtime), E())
	assertBuild(t, [1]compact_time.Time{ctime}, L(), CT(ctime), E())
	assertBuild(t, [1][]int{[]int{1}}, L(), L(), I(1), E(), E())
	assertBuild(t, [1]map[int]int{map[int]int{1: 2}}, L(), M(), I(1), I(2), E(), E())
}

func TestBuilderArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]int{}, N())
	assertBuildPanics(t, [1]int{}, M())
	assertBuildPanics(t, [1][]int{}, L(), M())
}

func TestBuilderByteArray(t *testing.T) {
	assertBuild(t, [1]byte{1}, AU8([]byte{1}))
}

func TestBuilderByteArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]byte{}, N())
	assertBuildPanics(t, [1]byte{}, B(false))
	assertBuildPanics(t, [1]byte{}, PI(1))
	assertBuildPanics(t, [1]byte{}, NI(1))
	assertBuildPanics(t, [1]byte{}, F(1.1))
	assertBuildPanics(t, [1]byte{}, BF(NewBigFloat("1.1")))
	assertBuildPanics(t, [1]byte{}, DF(NewDFloat("1.1")))
	assertBuildPanics(t, [1]byte{}, BDF(NewBDF("1.1")))
	assertBuildPanics(t, [1]byte{}, BF(NewBigFloat("1e20")))
	assertBuildPanics(t, [1]byte{}, DF(NewDFloat("1e20")))
	assertBuildPanics(t, [1]byte{}, BDF(NewBDF("1e20")))
	assertBuildPanics(t, [1]byte{}, BI(NewBigInt("100000000000000000000")))
	assertBuildPanics(t, [1]byte{}, S(""))
	assertBuildPanics(t, [1]byte{}, RID("x://x"))
	assertBuildPanics(t, [1]byte{}, UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, [1]byte{}, CT(test.AsCompactTime(time.Now())))
	assertBuildPanics(t, [1]byte{}, GT(time.Now()))
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
		14: ctime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		M(),
		I(1), N(),
		I(2), B(true),
		I(3), PI(1),
		I(4), NI(1),
		I(5), F(1.1),
		I(6), BF(NewBigFloat("1.1")),
		I(7), DF(NewDFloat("1.1")),
		I(8), BDF(NewBDF("1.1")),
		I(9), BI(NewBigInt("100000000000000000000")),
		I(10), S("test"),
		I(11), RID("http://example.com"),
		I(12), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		I(13), GT(gtime),
		I(14), CT(ctime),
		I(15), L(), F(1), E(),
		I(16), M(), I(1), I(2), E(),
		I(17), AU8([]byte{1}),
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
		ctime,
		[]float64{1},
		map[int]int{1: 2},
		[]byte{1},
	}, L(),
		// n(),
		B(true),
		PI(1),
		NI(1),
		F(1.1),
		BF(NewBigFloat("1.1")),
		DF(NewDFloat("1.1")),
		BDF(NewBDF("1.1")),
		BI(NewBigInt("100000000000000000000")),
		S("test"),
		RID("http://example.com"),
		UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		GT(gtime),
		CT(ctime),
		L(), F(1), E(),
		M(), I(1), I(2), E(),
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
		14: ctime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		M(),
		I(1), N(),
		I(2), B(true),
		I(3), PI(1),
		I(4), NI(1),
		I(5), F(1.1),
		I(6), BF(NewBigFloat("1.1")),
		I(7), DF(NewDFloat("1.1")),
		I(8), BDF(NewBDF("1.1")),
		I(9), BI(NewBigInt("100000000000000000000")),
		I(10), S("test"),
		I(11), RID("http://example.com"),
		I(12), UID([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		I(13), GT(gtime),
		I(14), CT(ctime),
		I(15), L(), F(1), E(),
		I(16), M(), I(1), I(2), E(),
		I(17), AU8([]byte{1}),
		E())
}

// // Older tests

type BuilderTestStruct struct {
	internal string
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
		S("AMAP"), M(), I(1), I(50), E(),
		S("AnINT"), I(1),
		S("Aslice"), L(), S("the slice"), E(),
		S("abool"), B(true),
		E())
}

func TestBuilderStructIgnored(t *testing.T) {
	assertBuild(t, BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
	}, M(), S("AString"), S("test"), S("Something"), I(5), S("AnInt"), I(1), S("ABool"), B(true), E())
}

func TestBuilderListStruct(t *testing.T) {
	assertBuild(t,
		[]BuilderTestStruct{
			BuilderTestStruct{
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
		S("AMap"), M(), I(1), I(50), E(),
		S("AnInt"), I(1),
		S("ASlice"), L(), S("the slice"), E(),
		S("ABool"), B(true),
		E(),
		E())
}

func TestBuilderMapStruct(t *testing.T) {
	assertBuild(t,
		map[string]BuilderTestStruct{
			"struct": BuilderTestStruct{
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
		S("AMap"), M(), I(1), I(50), E(),
		S("AnInt"), I(1),
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
			S("AMap"), M(), I(1), I(50), E(),
			S("AnInt"), I(1),
			S("ASlice"), L(), S("the slice"), E(),
			S("ABool"), B(true),
			E())
	}
}

type BuilderPtrTestStruct struct {
	internal    string
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
	var anInterface interface{}
	anInterface = 100
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
		S("AnInt"), I(1),
		S("AnInt8"), I(2),
		S("AnInt16"), I(3),
		S("AnInt32"), I(4),
		S("AnInt64"), I(5),
		S("AUint"), PI(6),
		S("AUint8"), PI(7),
		S("AUint16"), PI(8),
		S("AUint32"), PI(9),
		S("AUint64"), PI(10),
		S("AFloat32"), F(11.5),
		S("AFloat64"), F(12.5),
		S("AString"), S("test"),
		S("AnInterface"), I(100),
		E())
}

type BuilderSliceTestStruct struct {
	internal    []string
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
		BuilderSliceTestStruct{
			AnInt: []int{1},
		},
		BuilderSliceTestStruct{
			AnInt: []int{1},
		},
	}

	assertBuild(t,
		v,
		L(),
		M(),
		S("AnInt"), L(), I(1), E(),
		E(),
		M(),
		S("AnInt"), L(), I(1), E(),
		E(),
		E())
}

type SimpleTestStruct struct {
	IValue int
}

func TestBuilderListOfStruct(t *testing.T) {
	v := []*SimpleTestStruct{
		&SimpleTestStruct{
			IValue: 5,
		},
	}

	assertBuild(t,
		v,
		L(),
		M(),
		S("IValue"),
		I(5),
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
		N(),
		S("Slice"),
		N(),
		S("Map"),
		N(),
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
		N(),
		E())
}

func TestBuilderByteArrayBytes(t *testing.T) {
	assertBuild(t, [2]byte{1, 2},
		AU8([]byte{1, 2}))
}

func TestBuilderMarkerSlice(t *testing.T) {
	// Reference in same container
	assertBuild(t, []int{100, 100}, L(), MARK("1"), PI(100), REF("1"), E())
	assertBuild(t, []int{100, 100}, L(), REF("1"), MARK("1"), PI(100), E())
	assertBuild(t, []string{"abcdef", "abcdef"}, L(), MARK("1"), S("abcdef"), REF("1"), E())
	assertBuild(t, [][]int{[]int{100}, []int{100}}, L(), MARK("1"), L(), PI(100), E(), REF("1"), E())
	assertBuild(t, [][]int{[]int{100}, []int{100}}, L(), REF("1"), MARK("1"), L(), PI(100), E(), E())
	assertBuild(t, [][]int{[]int{100, 100}, []int{100, 100}}, L(), REF("1"), MARK("1"), L(), PI(100), PI(100), E(), E())

	// Reference in different container
	assertBuild(t, [][]int{[]int{100, 100}, []int{100}}, L(), L(), REF("1"), REF("1"), E(), L(), MARK("1"), PI(100), E(), E())

	// Referenced containers
	assertBuild(t, [][]int{[]int{}, []int{}}, L(), MARK("1"), L(), E(), REF("1"), E())
	assertBuild(t, [][]int{[]int{}, []int{}}, L(), REF("1"), MARK("1"), L(), E(), E())
	assertBuild(t, []map[int]int{map[int]int{100: 100}, map[int]int{100: 100}},
		L(), REF("1"), MARK("1"), M(), PI(100), PI(100), E(), E())
	assertBuild(t, []map[int]int{map[int]int{100: 100}, map[int]int{100: 100}},
		L(), MARK("1"), M(), PI(100), PI(100), E(), REF("1"), E())

	// Interface
	assertBuild(t, []interface{}{100, 100}, L(), REF("1"), MARK("1"), PI(100), E())
	assertBuild(t, []interface{}{100, 100}, L(), MARK("1"), PI(100), REF("1"), E())

	// Recursive interface
	rintf := make([]interface{}, 1)
	rintf[0] = rintf
	assertBuild(t, rintf, MARK("1"), L(), REF("1"), E())
}

func TestBuilderMarkerArray(t *testing.T) {
	assertBuild(t, [2]int{100, 100}, L(), MARK("1"), PI(100), REF("1"), E())
	assertBuild(t, [2]int{100, 100}, L(), REF("1"), MARK("1"), PI(100), E())
}

func TestBuilderMarkerMap(t *testing.T) {
	assertBuild(t, map[int]int8{1: 100, 2: 100, 3: 100},
		M(), PI(1), REF("5"), PI(2), MARK("5"), PI(100), PI(3), REF("5"), E())

	rmap := make(map[int]interface{})
	rmap[0] = rmap
	assertBuild(t, rmap,
		MARK("1"), M(), PI(0), REF("1"), E())

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
		PI(100),
		S("Next"),
		REF("1"),
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
		S("I16"), REF("1"),
		S("F32"), MARK("1"), PI(1000),
		S("S1"), MARK("2"), S("test"),
		S("S2"), REF("2"),
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
	assertBuild(t, m, M(), I(1), UID([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}), E())
}

func TestBuilderMedia(t *testing.T) {
	media := types.Media{
		MediaType: "a",
		Data:      []byte{1},
	}

	assertBuild(t, media, MB(), AC(1, false), AD([]byte("a")), AC(1, false), AD([]byte{1}))
	list := []types.Media{media}
	assertBuild(t, list, L(), MB(), AC(1, false), AD([]byte("a")), AC(1, false), AD([]byte{1}), E())
	m := map[int]types.Media{1: media}
	assertBuild(t, m, M(), I(1), MB(), AC(1, false), AD([]byte("a")), AC(1, false), AD([]byte{1}), E())
}

func TestBuilderMarkup(t *testing.T) {
	m := types.Markup{
		Name: "a",
	}
	pm := &m
	assertBuild(t, m, MUP("a"), E(), E())
	assertBuild(t, pm, MUP("a"), E(), E())

	m.Attributes = map[interface{}]interface{}{
		"a": 1,
	}
	m.Content = []interface{}{}
	assertBuild(t, m, MUP("a"), S("a"), I(1), E(), E())
	assertBuild(t, pm, MUP("a"), S("a"), I(1), E(), E())

	m.Attributes = map[interface{}]interface{}{}
	m.Content = []interface{}{
		"a",
	}
	assertBuild(t, m, MUP("a"), E(), S("a"), E())
	assertBuild(t, pm, MUP("a"), E(), S("a"), E())

	m.Attributes = map[interface{}]interface{}{
		"a": 1,
	}
	m.Content = []interface{}{
		"a",
	}
	assertBuild(t, m, MUP("a"), S("a"), I(1), E(), S("a"), E())
	assertBuild(t, pm, MUP("a"), S("a"), I(1), E(), S("a"), E())

	m.Attributes = map[interface{}]interface{}{}
	m.Content = []interface{}{}
	m.AddMarkup(&types.Markup{
		Name:       "b",
		Attributes: map[interface{}]interface{}{},
	})
	assertBuild(t, m, MUP("a"), E(), MUP("b"), E(), E(), E())
	assertBuild(t, pm, MUP("a"), E(), MUP("b"), E(), E(), E())

	m.Attributes = map[interface{}]interface{}{
		"a": 1,
	}
	m.Content = []interface{}{
		"a",
	}
	m.AddMarkup(&types.Markup{
		Name: "b",
		Attributes: map[interface{}]interface{}{
			100: "x",
		},
		Content: []interface{}{
			"z",
		},
	})
	assertBuild(t, m, MUP("a"), S("a"), I(1), E(), S("a"), MUP("b"), I(100), S("x"), E(), S("z"), E(), E())
	assertBuild(t, pm, MUP("a"), S("a"), I(1), E(), S("a"), MUP("b"), I(100), S("x"), E(), S("z"), E(), E())
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
	assertBuild(t, r, EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), I(1))
	assertBuild(t, pr, EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), I(1))

	assertBuild(t, []interface{}{pr}, L(), EDGE(), L(), RID("a"), RID("b"), E(), RID("http://y.com"), I(1), E())
}

type MyStruct struct {
	Aaa int
	Bbb []int
}

func TestBuilderStructBadFieldName(t *testing.T) {
	s := &MyStruct{}
	assertBuild(t, s, M(), S("aaa"), PI(0), S("bbb"), L(), E(), E())
	assertBuild(t, s, M(), S("x"), PI(0), S("bbb"), L(), E(), E())
}
