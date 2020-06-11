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
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func runBuild(expected interface{}, events ...*tevent) interface{} {
	builder := NewBuilderFor(expected)
	invokeEvents(builder, events...)
	return builder.GetBuiltObject()
}

func assertBuild(t *testing.T, expected interface{}, events ...*tevent) {
	actual := runBuild(expected, events...)
	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func assertBuildPanics(t *testing.T, template interface{}, events ...*tevent) {
	assertPanics(t, func() {
		runBuild(template, events...)
	})
}

// ============================================================================

func TestBuilderBasicTypes(t *testing.T) {
	pBigIntP := newBigInt("12345678901234567890123456789")
	pBigIntN := newBigInt("-999999999999999999999999999999")
	pBigFloat := newBigFloat("1.2345678901234567890123456789e10000", 29)
	dfloat := newDFloat("1.23456e1000")
	pBigDFloat := newBDF("4.509e10000")
	gTimeNow := time.Now()
	pCTimeNow := compact_time.AsCompactTime(gTimeNow)
	pCTime := compact_time.NewTimeLatLong(10, 5, 59, 100, 506, 107)
	pURL := newURL("http://x.com")

	assertBuildPanics(t, nil, n())
	assertBuild(t, true, b(true))
	assertBuild(t, false, b(false))
	assertBuild(t, int(10), i(10))
	assertBuild(t, int8(10), i(10))
	assertBuild(t, int16(-10), i(-10))
	assertBuild(t, int32(10), i(10))
	assertBuild(t, int64(-10), i(-10))
	assertBuild(t, uint(10), i(10))
	assertBuild(t, uint8(10), i(10))
	assertBuild(t, uint16(10), i(10))
	assertBuild(t, uint32(10), i(10))
	assertBuild(t, uint64(10), i(10))
	assertBuild(t, 1, i(1))
	assertBuild(t, -1, i(-1))
	assertBuild(t, pBigIntP, bi(pBigIntP))
	assertBuild(t, *pBigIntP, bi(pBigIntP))
	assertBuild(t, pBigIntN, bi(pBigIntN))
	assertBuild(t, *pBigIntN, bi(pBigIntN))
	assertBuild(t, (*big.Int)(nil), n())
	assertBuild(t, float32(-1.25), f(-1.25))
	assertBuild(t, float64(-9.5e50), f(-9.5e50))
	assertBuild(t, pBigFloat, bf(pBigFloat))
	assertBuild(t, *pBigFloat, bf(pBigFloat))
	assertBuild(t, (*big.Float)(nil), n())
	assertBuild(t, dfloat, df(dfloat))
	assertBuild(t, pBigDFloat, bdf(pBigDFloat))
	assertBuild(t, *pBigDFloat, bdf(pBigDFloat))
	assertBuild(t, (*apd.Decimal)(nil), n())
	// TODO
	// assertBuild(t, complex64(-1+5i), cplx(-1+5i))
	// assertBuild(t, complex128(-1+5i), cplx(-1+5i))
	assertBuild(t, signalingNan, snan())
	assertBuild(t, quietNan, nan())
	assertBuild(t, gTimeNow, gt(gTimeNow))
	assertBuild(t, pCTimeNow, ct(pCTimeNow))
	assertBuild(t, *pCTimeNow, ct(pCTimeNow))
	assertBuild(t, pCTime, ct(pCTime))
	assertBuild(t, *pCTime, ct(pCTime))
	assertBuild(t, []byte{1, 2, 3, 4}, bin([]byte{1, 2, 3, 4}))
	assertBuild(t, "test", s("test"))
	assertBuild(t, pURL, uri("http://x.com"))
	assertBuild(t, *pURL, uri("http://x.com"))
	assertBuild(t, (*url.URL)(nil), n())
	assertBuild(t, interface{}(1234), i(1234))
}

func TestBuilderConvertToBDF(t *testing.T) {
	pv := newBDF("1")
	nv := newBDF("-1")
	assertBuild(t, pv, pi(1))
	assertBuild(t, nv, ni(1))
	assertBuild(t, pv, bi(newBigInt("1")))
	assertBuild(t, pv, f(1))
	assertBuild(t, pv, bf(newBigFloat("1", 1)))
	assertBuild(t, pv, df(newDFloat("1")))
	assertBuild(t, pv, bdf(newBDF("1")))

	assertBuild(t, *pv, pi(1))
	assertBuild(t, *nv, ni(1))
	assertBuild(t, *pv, bi(newBigInt("1")))
	assertBuild(t, *pv, f(1))
	assertBuild(t, *pv, bf(newBigFloat("1", 1)))
	assertBuild(t, *pv, df(newDFloat("1")))
	assertBuild(t, *pv, bdf(newBDF("1")))
}

func TestBuilderConvertToBDFFail(t *testing.T) {
	v := newBDF("1")
	assertBuildPanics(t, v, b(true))
	assertBuildPanics(t, v, s("1"))
	assertBuildPanics(t, v, bin([]byte{1}))
	assertBuildPanics(t, v, uri("x://x"))
	assertBuildPanics(t, v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, gt(time.Now()))
	assertBuildPanics(t, v, l())
	assertBuildPanics(t, v, m())
	assertBuildPanics(t, v, e())

	assertBuildPanics(t, *v, n())
	assertBuildPanics(t, *v, b(true))
	assertBuildPanics(t, *v, s("1"))
	assertBuildPanics(t, *v, bin([]byte{1}))
	assertBuildPanics(t, *v, uri("x://x"))
	assertBuildPanics(t, *v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, gt(time.Now()))
	assertBuildPanics(t, *v, l())
	assertBuildPanics(t, *v, m())
	assertBuildPanics(t, *v, e())
}

func TestBuilderConvertToBF(t *testing.T) {
	pv := newBigFloat("1", 1)
	nv := newBigFloat("-1", 1)
	assertBuild(t, pv, pi(1))
	assertBuild(t, nv, ni(1))
	assertBuild(t, pv, bi(newBigInt("1")))
	assertBuild(t, pv, f(1))
	assertBuild(t, pv, bf(newBigFloat("1", 1)))
	assertBuild(t, pv, df(newDFloat("1")))
	assertBuild(t, pv, bdf(newBDF("1")))

	assertBuild(t, *pv, pi(1))
	assertBuild(t, *nv, ni(1))
	assertBuild(t, *pv, bi(newBigInt("1")))
	assertBuild(t, *pv, f(1))
	assertBuild(t, *pv, bf(newBigFloat("1", 1)))
	assertBuild(t, *pv, df(newDFloat("1")))
	assertBuild(t, *pv, bdf(newBDF("1")))
}

func TestBuilderConvertToBFFail(t *testing.T) {
	v := newBigFloat("1", 1)
	assertBuildPanics(t, v, b(true))
	assertBuildPanics(t, v, s("1"))
	assertBuildPanics(t, v, bin([]byte{1}))
	assertBuildPanics(t, v, uri("x://x"))
	assertBuildPanics(t, v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, gt(time.Now()))
	assertBuildPanics(t, v, l())
	assertBuildPanics(t, v, m())
	assertBuildPanics(t, v, e())

	assertBuildPanics(t, *v, n())
	assertBuildPanics(t, *v, b(true))
	assertBuildPanics(t, *v, s("1"))
	assertBuildPanics(t, *v, bin([]byte{1}))
	assertBuildPanics(t, *v, uri("x://x"))
	assertBuildPanics(t, *v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, gt(time.Now()))
	assertBuildPanics(t, *v, l())
	assertBuildPanics(t, *v, m())
	assertBuildPanics(t, *v, e())
}

func TestBuilderConvertToBI(t *testing.T) {
	pv := newBigInt("1")
	nv := newBigInt("-1")
	assertBuild(t, pv, pi(1))
	assertBuild(t, newBigInt("9223372036854775808"), pi(9223372036854775808))
	assertBuild(t, nv, ni(1))
	assertBuild(t, pv, bi(newBigInt("1")))
	assertBuild(t, pv, f(1))
	assertBuild(t, pv, bf(newBigFloat("1", 1)))
	assertBuild(t, pv, df(newDFloat("1")))
	assertBuild(t, pv, bdf(newBDF("1")))

	assertBuild(t, *pv, pi(1))
	assertBuild(t, *nv, ni(1))
	assertBuild(t, *pv, bi(newBigInt("1")))
	assertBuild(t, *pv, f(1))
	assertBuild(t, *pv, bf(newBigFloat("1", 1)))
	assertBuild(t, *pv, df(newDFloat("1")))
	assertBuild(t, *pv, bdf(newBDF("1")))
}

func TestBuilderConvertToBIFail(t *testing.T) {
	v := newBigInt("1")
	assertBuildPanics(t, v, b(true))
	assertBuildPanics(t, v, f(1.1))
	assertBuildPanics(t, v, bf(newBigFloat("1.1", 1)))
	assertBuildPanics(t, v, bf(newBigFloat("1.0e100000", 1)))
	assertBuildPanics(t, v, df(newDFloat("1.1")))
	assertBuildPanics(t, v, df(newDFloat("1.0e100000")))
	assertBuildPanics(t, v, bdf(newBDF("1.1")))
	assertBuildPanics(t, v, bdf(newBDF("1.0e100000")))
	assertBuildPanics(t, v, bdf(newBDF("nan")))
	assertBuildPanics(t, v, bdf(newBDF("snan")))
	assertBuildPanics(t, v, bdf(newBDF("inf")))
	assertBuildPanics(t, v, bdf(newBDF("-inf")))
	assertBuildPanics(t, v, s("1"))
	assertBuildPanics(t, v, bin([]byte{1}))
	assertBuildPanics(t, v, uri("x://x"))
	assertBuildPanics(t, v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, gt(time.Now()))
	assertBuildPanics(t, v, l())
	assertBuildPanics(t, v, m())
	assertBuildPanics(t, v, e())

	assertBuildPanics(t, *v, n())
	assertBuildPanics(t, *v, f(1.1))
	assertBuildPanics(t, *v, b(true))
	assertBuildPanics(t, *v, bf(newBigFloat("1.1", 1)))
	assertBuildPanics(t, *v, bf(newBigFloat("1.0e100000", 1)))
	assertBuildPanics(t, *v, df(newDFloat("1.1")))
	assertBuildPanics(t, *v, df(newDFloat("1.0e100000")))
	assertBuildPanics(t, *v, bdf(newBDF("1.1")))
	assertBuildPanics(t, *v, bdf(newBDF("1.0e100000")))
	assertBuildPanics(t, *v, bdf(newBDF("nan")))
	assertBuildPanics(t, *v, bdf(newBDF("snan")))
	assertBuildPanics(t, *v, bdf(newBDF("inf")))
	assertBuildPanics(t, *v, bdf(newBDF("-inf")))
	assertBuildPanics(t, *v, s("1"))
	assertBuildPanics(t, *v, bin([]byte{1}))
	assertBuildPanics(t, *v, uri("x://x"))
	assertBuildPanics(t, *v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, *v, gt(time.Now()))
	assertBuildPanics(t, *v, l())
	assertBuildPanics(t, *v, m())
	assertBuildPanics(t, *v, e())
}

func TestBuilderConvertToDecimalFloat(t *testing.T) {
	pv := newDFloat("1")
	nv := newDFloat("-1")
	assertBuild(t, pv, pi(1))
	assertBuild(t, nv, ni(1))
	assertBuild(t, pv, bi(newBigInt("1")))
	assertBuild(t, pv, f(1))
	assertBuild(t, pv, bf(newBigFloat("1", 1)))
	assertBuild(t, pv, df(newDFloat("1")))
	assertBuild(t, pv, bdf(newBDF("1")))
}

func TestBuilderDecimalFloatFail(t *testing.T) {
	v := newDFloat("1")
	assertBuildPanics(t, v, n())
	assertBuildPanics(t, v, b(true))
	assertBuildPanics(t, v, s("1"))
	assertBuildPanics(t, v, bin([]byte{1}))
	assertBuildPanics(t, v, uri("x://x"))
	assertBuildPanics(t, v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, gt(time.Now()))
	assertBuildPanics(t, v, l())
	assertBuildPanics(t, v, m())
	assertBuildPanics(t, v, e())
}

func TestBuilderConvertToFloat(t *testing.T) {
	pv := 1.0
	nv := -1.0
	assertBuild(t, pv, pi(1))
	assertBuild(t, nv, ni(1))
	assertBuild(t, pv, bi(newBigInt("1")))
	assertBuild(t, pv, f(1))
	assertBuild(t, pv, bf(newBigFloat("1", 1)))
	assertBuild(t, pv, df(newDFloat("1")))
	assertBuild(t, pv, bdf(newBDF("1")))
}

func TestBuilderConvertToFloatFail(t *testing.T) {
	// TODO: How to define required conversion accuracy?
	v := 1.0
	assertBuildPanics(t, v, b(true))
	assertBuildPanics(t, v, pi(0xffffffffffffffff))
	assertBuildPanics(t, v, i(-0x7fffffffffffffff))
	assertBuildPanics(t, v, bf(newBigFloat("1.0e309", 1)))
	assertBuildPanics(t, v, bf(newBigFloat("1.0e-311", 1)))
	// TODO: apd.Decimal and compact_float.DFloat don't handle float overflow
	assertBuildPanics(t, v, bi(newBigInt("1234567890123456789012345")))
	assertBuildPanics(t, v, s("1"))
	assertBuildPanics(t, v, bin([]byte{1}))
	assertBuildPanics(t, v, uri("x://x"))
	assertBuildPanics(t, v, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, v, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, v, gt(time.Now()))
	assertBuildPanics(t, v, l())
	assertBuildPanics(t, v, m())
	assertBuildPanics(t, v, e())
}

func TestBuilderConvertToInt(t *testing.T) {
	assertBuild(t, 1, pi(1))
	assertBuild(t, -1, ni(1))
	assertBuild(t, 1, bi(newBigInt("1")))
	assertBuild(t, 1, f(1))
	assertBuild(t, 1, bf(newBigFloat("1", 1)))
	assertBuild(t, 1, df(newDFloat("1")))
	assertBuild(t, 1, bdf(newBDF("1")))
}

func TestBuilderConvertToIntFail(t *testing.T) {
	assertBuildPanics(t, int(1), n())
	assertBuildPanics(t, int(1), b(true))
	assertBuildPanics(t, int(1), pi(0x8000000000000000))
	assertBuildPanics(t, int(1), f(1.1))
	assertBuildPanics(t, int(1), bf(newBigFloat("1.1", 1)))
	assertBuildPanics(t, int(1), df(newDFloat("1.1")))
	assertBuildPanics(t, int(1), bdf(newBDF("1.1")))
	assertBuildPanics(t, int(1), bf(newBigFloat("1e20", 1)))
	assertBuildPanics(t, int(1), df(newDFloat("1e20")))
	assertBuildPanics(t, int(1), bdf(newBDF("1e20")))
	assertBuildPanics(t, int(1), bf(newBigFloat("-1e20", 1)))
	assertBuildPanics(t, int(1), df(newDFloat("-1e20")))
	assertBuildPanics(t, int(1), bdf(newBDF("-1e20")))
	assertBuildPanics(t, int(1), bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, int(1), bi(newBigInt("-100000000000000000000")))
	assertBuildPanics(t, int(1), s("1"))
	assertBuildPanics(t, int(1), bin([]byte{1}))
	assertBuildPanics(t, int(1), uri("x://x"))
	assertBuildPanics(t, int(1), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, int(1), ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, int(1), gt(time.Now()))
	assertBuildPanics(t, int(1), l())
	assertBuildPanics(t, int(1), m())
	assertBuildPanics(t, int(1), e())

	assertBuildPanics(t, int8(1), pi(1000))
	assertBuildPanics(t, int8(1), ni(1000))
	assertBuildPanics(t, int8(1), i(1000))
	assertBuildPanics(t, int8(1), bi(newBigInt("1000")))
	assertBuildPanics(t, int8(1), f(1000))
	assertBuildPanics(t, int8(1), bf(newBigFloat("1000", 4)))
	assertBuildPanics(t, int8(1), df(newDFloat("1000")))
	assertBuildPanics(t, int8(1), bdf(newBDF("1000")))

	assertBuildPanics(t, int16(1), pi(100000))
	assertBuildPanics(t, int16(1), ni(100000))
	assertBuildPanics(t, int16(1), i(100000))
	assertBuildPanics(t, int16(1), bi(newBigInt("100000")))
	assertBuildPanics(t, int16(1), f(100000))
	assertBuildPanics(t, int16(1), bf(newBigFloat("100000", 6)))
	assertBuildPanics(t, int16(1), df(newDFloat("100000")))
	assertBuildPanics(t, int16(1), bdf(newBDF("100000")))

	assertBuildPanics(t, int32(1), pi(10000000000))
	assertBuildPanics(t, int32(1), ni(10000000000))
	assertBuildPanics(t, int32(1), i(10000000000))
	assertBuildPanics(t, int32(1), bi(newBigInt("10000000000")))
	assertBuildPanics(t, int32(1), f(10000000000))
	assertBuildPanics(t, int32(1), bf(newBigFloat("10000000000", 11)))
	assertBuildPanics(t, int32(1), df(newDFloat("10000000000")))
	assertBuildPanics(t, int32(1), bdf(newBDF("10000000000")))
}

func TestBuilderConvertToUint(t *testing.T) {
	assertBuild(t, uint(1), pi(1))
	assertBuild(t, uint(1), i(1))
	assertBuild(t, uint(1), bi(newBigInt("1")))
	assertBuild(t, uint(1), f(1))
	assertBuild(t, uint(1), bf(newBigFloat("1", 1)))
	assertBuild(t, uint(1), df(newDFloat("1")))
	assertBuild(t, uint(1), bdf(newBDF("1")))
}

func TestBuilderConvertToUintFail(t *testing.T) {
	assertBuildPanics(t, uint(1), n())
	assertBuildPanics(t, uint(1), b(true))
	assertBuildPanics(t, uint(1), ni(1))
	assertBuildPanics(t, uint(1), f(1.1))
	assertBuildPanics(t, uint(1), bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, uint(1), df(newDFloat("1.1")))
	assertBuildPanics(t, uint(1), bdf(newBDF("1.1")))
	assertBuildPanics(t, uint(1), bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, uint(1), df(newDFloat("1e20")))
	assertBuildPanics(t, uint(1), bdf(newBDF("1e20")))
	assertBuildPanics(t, uint8(1), bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, uint(1), s("1"))
	assertBuildPanics(t, uint(1), bin([]byte{1}))
	assertBuildPanics(t, uint(1), uri("x://x"))
	assertBuildPanics(t, uint(1), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, uint(1), ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, uint(1), gt(time.Now()))
	assertBuildPanics(t, uint(1), l())
	assertBuildPanics(t, uint(1), m())
	assertBuildPanics(t, uint(1), e())

	assertBuildPanics(t, uint8(1), pi(1000))
	assertBuildPanics(t, uint8(1), i(1000))
	assertBuildPanics(t, uint8(1), bi(newBigInt("1000")))
	assertBuildPanics(t, uint8(1), f(1000))
	assertBuildPanics(t, uint8(1), bf(newBigFloat("1000", 4)))
	assertBuildPanics(t, uint8(1), df(newDFloat("1000")))
	assertBuildPanics(t, uint8(1), bdf(newBDF("1000")))

	assertBuildPanics(t, uint16(1), pi(100000))
	assertBuildPanics(t, uint16(1), i(100000))
	assertBuildPanics(t, uint16(1), bi(newBigInt("100000")))
	assertBuildPanics(t, uint16(1), f(100000))
	assertBuildPanics(t, uint16(1), bf(newBigFloat("100000", 6)))
	assertBuildPanics(t, uint16(1), df(newDFloat("100000")))
	assertBuildPanics(t, uint16(1), bdf(newBDF("100000")))

	assertBuildPanics(t, uint32(1), pi(10000000000))
	assertBuildPanics(t, uint32(1), i(10000000000))
	assertBuildPanics(t, uint32(1), bi(newBigInt("10000000000")))
	assertBuildPanics(t, uint32(1), f(10000000000))
	assertBuildPanics(t, uint32(1), bf(newBigFloat("10000000000", 11)))
	assertBuildPanics(t, uint32(1), df(newDFloat("10000000000")))
	assertBuildPanics(t, uint32(1), bdf(newBDF("10000000000")))
}

//

func TestBuilderTime(t *testing.T) {
	testTime := time.Now()
	assertBuild(t, testTime, gt(testTime))
}

func TestBuilderURI(t *testing.T) {
	testURI := "https://x.com"
	assertBuild(t, newURI(testURI), uri(testURI))
}

func TestBuilderBasicTypeFail(t *testing.T) {
	assertBuildPanics(t, true, n())
	assertBuildPanics(t, true, i(1))
	assertBuildPanics(t, true, pi(1))
	assertBuildPanics(t, true, f(1))
	assertBuildPanics(t, true, s("1"))
	assertBuildPanics(t, true, bin([]byte{1}))
	assertBuildPanics(t, true, uri("x://x"))
	assertBuildPanics(t, true, gt(time.Now()))
	assertBuildPanics(t, true, l())
	assertBuildPanics(t, true, m())
	assertBuildPanics(t, true, e())

	assertBuildPanics(t, uint(1), n())
	assertBuildPanics(t, uint(1), b(true))
	assertBuildPanics(t, uint(1), s("1"))
	assertBuildPanics(t, uint(1), bin([]byte{1}))
	assertBuildPanics(t, uint(1), uri("x://x"))
	assertBuildPanics(t, uint(1), gt(time.Now()))
	assertBuildPanics(t, uint(1), l())
	assertBuildPanics(t, uint(1), m())
	assertBuildPanics(t, uint(1), e())

	assertBuildPanics(t, float64(1), n())
	assertBuildPanics(t, float64(1), b(true))
	assertBuildPanics(t, float64(1), s("1"))
	assertBuildPanics(t, float64(1), bin([]byte{1}))
	assertBuildPanics(t, float64(1), uri("x://x"))
	assertBuildPanics(t, float64(1), gt(time.Now()))
	assertBuildPanics(t, float64(1), l())
	assertBuildPanics(t, float64(1), m())
	assertBuildPanics(t, float64(1), e())

	assertBuildPanics(t, "", b(true))
	assertBuildPanics(t, "", i(1))
	assertBuildPanics(t, "", pi(1))
	assertBuildPanics(t, "", f(1))
	assertBuildPanics(t, "", bin([]byte{1}))
	assertBuildPanics(t, "", uri("x://x"))
	assertBuildPanics(t, "", gt(time.Now()))
	assertBuildPanics(t, "", l())
	assertBuildPanics(t, "", m())
	assertBuildPanics(t, "", e())

	assertBuildPanics(t, []byte{}, b(true))
	assertBuildPanics(t, []byte{}, i(1))
	assertBuildPanics(t, []byte{}, pi(1))
	assertBuildPanics(t, []byte{}, f(1))
	assertBuildPanics(t, []byte{}, s("1"))
	assertBuildPanics(t, []byte{}, uri("x://x"))
	assertBuildPanics(t, []byte{}, gt(time.Now()))
	assertBuildPanics(t, []byte{}, l())
	assertBuildPanics(t, []byte{}, m())
	assertBuildPanics(t, []byte{}, e())

	assertBuildPanics(t, newURI("x://x"), b(true))
	assertBuildPanics(t, newURI("x://x"), i(1))
	assertBuildPanics(t, newURI("x://x"), pi(1))
	assertBuildPanics(t, newURI("x://x"), f(1))
	assertBuildPanics(t, newURI("x://x"), s("1"))
	assertBuildPanics(t, newURI("x://x"), bin([]byte{1}))
	assertBuildPanics(t, newURI("x://x"), gt(time.Now()))
	assertBuildPanics(t, newURI("x://x"), l())
	assertBuildPanics(t, newURI("x://x"), m())
	assertBuildPanics(t, newURI("x://x"), e())

	assertBuildPanics(t, time.Now(), n())
	assertBuildPanics(t, time.Now(), b(true))
	assertBuildPanics(t, time.Now(), i(1))
	assertBuildPanics(t, time.Now(), pi(1))
	assertBuildPanics(t, time.Now(), f(1))
	assertBuildPanics(t, time.Now(), s("1"))
	assertBuildPanics(t, time.Now(), bin([]byte{1}))
	assertBuildPanics(t, time.Now(), uri("x://x"))
	assertBuildPanics(t, time.Now(), l())
	assertBuildPanics(t, time.Now(), m())
	assertBuildPanics(t, time.Now(), e())

	assertBuildPanics(t, []int{}, n())
	assertBuildPanics(t, []int{}, b(true))
	assertBuildPanics(t, []int{}, i(1))
	assertBuildPanics(t, []int{}, pi(1))
	assertBuildPanics(t, []int{}, f(1))
	assertBuildPanics(t, []int{}, s("1"))
	assertBuildPanics(t, []int{}, bin([]byte{1}))
	assertBuildPanics(t, []int{}, uri("x://x"))
	assertBuildPanics(t, []int{}, gt(time.Now()))
	assertBuildPanics(t, []int{}, m())
	assertBuildPanics(t, []int{}, e())

	assertBuildPanics(t, map[int]int{}, i(1))
	assertBuildPanics(t, map[int]int{}, pi(1))
	assertBuildPanics(t, map[int]int{}, f(1))
	assertBuildPanics(t, map[int]int{}, s("1"))
	assertBuildPanics(t, map[int]int{}, bin([]byte{1}))
	assertBuildPanics(t, map[int]int{}, uri("x://x"))
	assertBuildPanics(t, map[int]int{}, gt(time.Now()))
	assertBuildPanics(t, map[int]int{}, l())
	assertBuildPanics(t, map[int]int{}, e())
}

func TestBuilderNumericConversion(t *testing.T) {
	assertBuild(t, int8(50), pi(50))
	assertBuild(t, int16(50), f(50))
	assertBuild(t, uint32(50), i(50))
	assertBuild(t, uint64(50), f(50))
	assertBuild(t, float32(50), i(50))
	assertBuild(t, float64(50), pi(50))
	assertBuild(t, compact_float.DFloatValue(0, 50), pi(50))
}

func TestBuilderNumericConversionFail(t *testing.T) {
	assertBuildPanics(t, int8(0), i(300))
	assertBuildPanics(t, int(0), f(3.5))
	assertBuildPanics(t, uint(0), i(-1))
	assertBuildPanics(t, uint(0), f(3.5))
	assertBuildPanics(t, float32(0), i(0x7fffffffffffffff))
	assertBuildPanics(t, float64(0), pi(0xffffffffffffffff))
}

func TestBuilderSlice(t *testing.T) {
	assertBuild(t, []bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, []int8{1, 2, 3}, l(), i(1), pi(2), f(3), e())
	assertBuild(t, []interface{}{false, 1, "test"}, l(), b(false), i(1), s("test"), e())
}

func TestBuilderArray(t *testing.T) {
	assertBuild(t, [2]bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, [3]int8{1, 2, 3}, l(), i(1), pi(2), f(3), e())
	assertBuild(t, [3]interface{}{false, 1, "test"}, l(), b(false), i(1), s("test"), e())
}

func TestBuilderMap(t *testing.T) {
	assertBuild(t, map[string]bool{"true": true, "false": false}, m(), s("true"), b(true), s("false"), b(false), e())
	assertBuild(t, map[interface{}]interface{}{"false": false, 1: "one"}, m(), s("false"), b(false), i(1), s("one"), e())
}

func TestBuilderSliceSlice(t *testing.T) {
	assertBuild(t, [][]bool{{false, true}}, l(), l(), b(false), b(true), e(), e())
}

func TestBuilderMapMap(t *testing.T) {
	assertBuild(t, map[string]map[int]bool{"first": {1: true}}, m(), s("first"), m(), i(1), b(true), e(), e())
}

func TestBuilderSliceMap(t *testing.T) {
	assertBuild(t, []map[int]bool{{1: true}}, l(), m(), i(1), b(true), e(), e())
}

func TestBuilderMapSlice(t *testing.T) {
	assertBuild(t, map[string][]int{"first": {1}}, m(), s("first"), l(), i(1), e(), e())
}

type BuilderTestStruct struct {
	internal string
	ABool    bool
	AString  string
	AnInt    int
	AMap     map[int]int8
	ASlice   []string
}

func TestBuilderStruct(t *testing.T) {
	assertBuild(t,
		BuilderTestStruct{
			AString: "test",
			AnInt:   1,
			ABool:   true,
			AMap:    map[int]int8{1: 50},
			ASlice:  []string{"the slice"},
		},
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e())
}

func TestBuilderStructIgnored(t *testing.T) {
	assertBuild(t, BuilderTestStruct{
		AString: "test",
		AnInt:   1,
		ABool:   true,
	}, m(), s("AString"), s("test"), s("Something"), i(5), s("AnInt"), i(1), s("ABool"), b(true), e())
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
		l(),
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e(),
		e())
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
		m(),
		s("struct"),
		m(),
		s("AString"), s("test"),
		s("AMap"), m(), i(1), i(50), e(),
		s("AnInt"), i(1),
		s("ASlice"), l(), s("the slice"), e(),
		s("ABool"), b(true),
		e(),
		e())
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
			m(),
			s("AString"), s("test"),
			s("AMap"), m(), i(1), i(50), e(),
			s("AnInt"), i(1),
			s("ASlice"), l(), s("the slice"), e(),
			s("ABool"), b(true),
			e())
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
		m(),
		s("ABool"), b(true),
		s("AnInt"), i(1),
		s("AnInt8"), i(2),
		s("AnInt16"), i(3),
		s("AnInt32"), i(4),
		s("AnInt64"), i(5),
		s("AUint"), pi(6),
		s("AUint8"), pi(7),
		s("AUint16"), pi(8),
		s("AUint32"), pi(9),
		s("AUint64"), pi(10),
		s("AFloat32"), f(11.5),
		s("AFloat64"), f(12.5),
		s("AString"), s("test"),
		s("AnInterface"), i(100),
		e())
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
		l(),
		m(),
		s("AnInt"), l(), i(1), e(),
		e(),
		m(),
		s("AnInt"), l(), i(1), e(),
		e(),
		e())
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
		l(),
		m(),
		s("IValue"),
		i(5),
		e(),
		e())
}

type NilContainers struct {
	Bytes []byte
	Slice []interface{}
	Map   map[interface{}]interface{}
}

func TestBuilderNilContainers(t *testing.T) {
	v := NilContainers{}

	assertBuild(t, v,
		m(),
		s("Bytes"),
		n(),
		s("Slice"),
		n(),
		s("Map"),
		n(),
		e())
}

type PURLContainer struct {
	URL *url.URL
}

func TestBuilderPURLContainer(t *testing.T) {
	v := PURLContainer{newURI("http://x.com")}

	assertBuild(t, v,
		m(),
		s("URL"),
		uri("http://x.com"),
		e())
}

func TestBuilderNilPURLContainer(t *testing.T) {
	v := PURLContainer{}

	assertBuild(t, v,
		m(),
		s("URL"),
		n(),
		e())
}

func TestBuilderByteArrayBytes(t *testing.T) {
	assertBuild(t, [2]byte{1, 2},
		bin([]byte{1, 2}))
}

// TODO: Self-referential struct
// type SelfReferential struct {
// 	Self *SelfReferential
// }

// func TestBuilderSelfReferential(t *testing.T) {
// 	assertBuild(t, SelfReferential{},
// 		m(),
// 		s("Self"),
// 		n(),
// 		e())
// }
