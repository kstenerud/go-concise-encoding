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
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func runBuild(template interface{}, events ...*tevent) interface{} {
	builder := NewBuilderFor(template, nil)
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

func TestBuildUnknown(t *testing.T) {
	expected := []interface{}{1}
	actual := runBuild(nil, l(), i(1), e())

	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

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
	assertBuildPanics(t, v, n())
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

func TestBuilderString(t *testing.T) {
	assertBuild(t, "", n())
	assertBuild(t, "test", s("test"))
}

func TestBuilderStringFail(t *testing.T) {
	assertBuildPanics(t, "", b(false))
	assertBuildPanics(t, "", pi(1))
	assertBuildPanics(t, "", ni(1))
	assertBuildPanics(t, "", f(1.1))
	assertBuildPanics(t, "", bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, "", df(newDFloat("1.1")))
	assertBuildPanics(t, "", bdf(newBDF("1.1")))
	assertBuildPanics(t, "", bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, "", df(newDFloat("1e20")))
	assertBuildPanics(t, "", bdf(newBDF("1e20")))
	assertBuildPanics(t, "", bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, "", bin([]byte{1}))
	assertBuildPanics(t, "", uri("x://x"))
	assertBuildPanics(t, "", uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, "", ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, "", gt(time.Now()))
	assertBuildPanics(t, "", l())
	assertBuildPanics(t, "", m())
	assertBuildPanics(t, "", e())
}

func TestBuilderGoTime(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, gtime, gt(gtime))
	assertBuild(t, gtime, ct(ctime))
}

func TestBuilderGoTimeFail(t *testing.T) {
	gtime := time.Time{}
	ctime := compact_time.NewTimeLatLong(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, gtime, n())
	assertBuildPanics(t, gtime, b(true))
	assertBuildPanics(t, gtime, pi(1))
	assertBuildPanics(t, gtime, ni(1))
	assertBuildPanics(t, gtime, f(1.1))
	assertBuildPanics(t, gtime, bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, gtime, df(newDFloat("1.1")))
	assertBuildPanics(t, gtime, bdf(newBDF("1.1")))
	assertBuildPanics(t, gtime, bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, gtime, df(newDFloat("1e20")))
	assertBuildPanics(t, gtime, bdf(newBDF("1e20")))
	assertBuildPanics(t, gtime, bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, gtime, s("1"))
	assertBuildPanics(t, gtime, bin([]byte{1}))
	assertBuildPanics(t, gtime, uri("x://x"))
	assertBuildPanics(t, gtime, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, gtime, ct(ctime))
	assertBuildPanics(t, gtime, l())
	assertBuildPanics(t, gtime, m())
	assertBuildPanics(t, gtime, e())
}

func TestBuilderCompactTime(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, (*compact_time.Time)(nil), n())
	assertBuild(t, ctime, gt(gtime))
	assertBuild(t, ctime, ct(ctime))

	assertBuild(t, *ctime, gt(gtime))
	assertBuild(t, *ctime, ct(ctime))
}

func TestBuilderCompactTimeFail(t *testing.T) {
	ctime := compact_time.NewTimeLatLong(1, 1, 1, 1, 100, 0)
	assertBuildPanics(t, ctime, b(true))
	assertBuildPanics(t, ctime, pi(1))
	assertBuildPanics(t, ctime, ni(1))
	assertBuildPanics(t, ctime, f(1.1))
	assertBuildPanics(t, ctime, bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, ctime, df(newDFloat("1.1")))
	assertBuildPanics(t, ctime, bdf(newBDF("1.1")))
	assertBuildPanics(t, ctime, bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, ctime, df(newDFloat("1e20")))
	assertBuildPanics(t, ctime, bdf(newBDF("1e20")))
	assertBuildPanics(t, ctime, bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, ctime, s("1"))
	assertBuildPanics(t, ctime, bin([]byte{1}))
	assertBuildPanics(t, ctime, uri("x://x"))
	assertBuildPanics(t, ctime, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, ctime, l())
	assertBuildPanics(t, ctime, m())
	assertBuildPanics(t, ctime, e())

	assertBuildPanics(t, *ctime, n())
	assertBuildPanics(t, *ctime, b(true))
	assertBuildPanics(t, *ctime, pi(1))
	assertBuildPanics(t, *ctime, ni(1))
	assertBuildPanics(t, *ctime, f(1.1))
	assertBuildPanics(t, *ctime, bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, *ctime, df(newDFloat("1.1")))
	assertBuildPanics(t, *ctime, bdf(newBDF("1.1")))
	assertBuildPanics(t, *ctime, bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, *ctime, df(newDFloat("1e20")))
	assertBuildPanics(t, *ctime, bdf(newBDF("1e20")))
	assertBuildPanics(t, *ctime, bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, *ctime, s("1"))
	assertBuildPanics(t, *ctime, bin([]byte{1}))
	assertBuildPanics(t, *ctime, uri("x://x"))
	assertBuildPanics(t, *ctime, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, *ctime, l())
	assertBuildPanics(t, *ctime, m())
	assertBuildPanics(t, *ctime, e())
}

func TestBuilderSlice(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, []bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, []int8{-1, 2, 3, 4, 5, 6, 7}, l(), i(-1), pi(2), f(3),
		bi(newBigInt("4")), bf(newBigFloat("5", 1)), df(newDFloat("6")),
		bdf(newBDF("7")), e())
	assertBuild(t, []*int{nil}, l(), n(), e())
	assertBuild(t, [][]byte{[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, l(), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}), e())
	assertBuild(t, []string{"test"}, l(), s("test"), e())
	assertBuild(t, [][]byte{[]byte{1}}, l(), bin([]byte{1}), e())
	assertBuild(t, []*url.URL{newURI("http://example.com")}, l(), uri("http://example.com"), e())
	assertBuild(t, []time.Time{gtime}, l(), gt(gtime), e())
	assertBuild(t, []*compact_time.Time{ctime}, l(), ct(ctime), e())
	assertBuild(t, [][]int{[]int{1}}, l(), l(), i(1), e(), e())
	assertBuild(t, []map[int]int{map[int]int{1: 2}}, l(), m(), i(1), i(2), e(), e())
}

func TestBuilderSliceFail(t *testing.T) {
	assertBuildPanics(t, []int{}, n())
	assertBuildPanics(t, []int{}, m())
	assertBuildPanics(t, [][]int{}, l(), m())
}

func TestBuilderArray(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, [2]bool{false, true}, l(), b(false), b(true), e())
	assertBuild(t, [7]int8{-1, 2, 3, 4, 5, 6, 7}, l(), i(-1), pi(2), f(3),
		bi(newBigInt("4")), bf(newBigFloat("5", 1)), df(newDFloat("6")),
		bdf(newBDF("7")), e())
	assertBuild(t, [1]*int{nil}, l(), n(), e())
	assertBuild(t, [1][]byte{[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}, l(), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}), e())
	assertBuild(t, [1]string{"test"}, l(), s("test"), e())
	assertBuild(t, [1][]byte{[]byte{1}}, l(), bin([]byte{1}), e())
	assertBuild(t, [1]*url.URL{newURI("http://example.com")}, l(), uri("http://example.com"), e())
	assertBuild(t, [1]time.Time{gtime}, l(), gt(gtime), e())
	assertBuild(t, [1]*compact_time.Time{ctime}, l(), ct(ctime), e())
	assertBuild(t, [1][]int{[]int{1}}, l(), l(), i(1), e(), e())
	assertBuild(t, [1]map[int]int{map[int]int{1: 2}}, l(), m(), i(1), i(2), e(), e())
}

func TestBuilderArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]int{}, n())
	assertBuildPanics(t, [1]int{}, m())
	assertBuildPanics(t, [1][]int{}, l(), m())
}

func TestBuilderByteArray(t *testing.T) {
	assertBuild(t, [1]byte{1}, bin([]byte{1}))
}

func TestBuilderByteArrayFail(t *testing.T) {
	assertBuildPanics(t, [1]byte{}, n())
	assertBuildPanics(t, [1]byte{}, b(false))
	assertBuildPanics(t, [1]byte{}, pi(1))
	assertBuildPanics(t, [1]byte{}, ni(1))
	assertBuildPanics(t, [1]byte{}, f(1.1))
	assertBuildPanics(t, [1]byte{}, bf(newBigFloat("1.1", 2)))
	assertBuildPanics(t, [1]byte{}, df(newDFloat("1.1")))
	assertBuildPanics(t, [1]byte{}, bdf(newBDF("1.1")))
	assertBuildPanics(t, [1]byte{}, bf(newBigFloat("1e20", 2)))
	assertBuildPanics(t, [1]byte{}, df(newDFloat("1e20")))
	assertBuildPanics(t, [1]byte{}, bdf(newBDF("1e20")))
	assertBuildPanics(t, [1]byte{}, bi(newBigInt("100000000000000000000")))
	assertBuildPanics(t, [1]byte{}, s(""))
	assertBuildPanics(t, [1]byte{}, uri("x://x"))
	assertBuildPanics(t, [1]byte{}, uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}))
	assertBuildPanics(t, [1]byte{}, ct(compact_time.AsCompactTime(time.Now())))
	assertBuildPanics(t, [1]byte{}, gt(time.Now()))
	assertBuildPanics(t, [1]byte{}, l())
	assertBuildPanics(t, [1]byte{}, m())
	assertBuildPanics(t, [1]byte{}, e())
}

func TestBuilderMap(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, map[int]interface{}{
		1:  nil,
		2:  true,
		3:  1,
		4:  -1,
		5:  1.1,
		6:  newBigFloat("1.1", 2),
		7:  newDFloat("1.1"),
		8:  newBDF("1.1"),
		9:  newBigInt("100000000000000000000"),
		10: "test",
		11: newURI("http://example.com"),
		12: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		13: gtime,
		14: ctime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		m(),
		i(1), n(),
		i(2), b(true),
		i(3), pi(1),
		i(4), ni(1),
		i(5), f(1.1),
		i(6), bf(newBigFloat("1.1", 2)),
		i(7), df(newDFloat("1.1")),
		i(8), bdf(newBDF("1.1")),
		i(9), bi(newBigInt("100000000000000000000")),
		i(10), s("test"),
		i(11), uri("http://example.com"),
		i(12), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		i(13), gt(gtime),
		i(14), ct(ctime),
		i(15), l(), f(1), e(),
		i(16), m(), i(1), i(2), e(),
		i(17), bin([]byte{1}),
		e())
}

func TestBuilderStruct(t *testing.T) {
	s := newTestingOuterStruct(1)
	includeFakes := true
	assertBuild(t, s, s.getRepresentativeEvents(includeFakes)...)
}

func TestBuilderInterfaceSlice(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, []interface{}{
		// nil,
		true,
		1,
		-1,
		1.1,
		newBigFloat("1.1", 2),
		newDFloat("1.1"),
		newBDF("1.1"),
		newBigInt("100000000000000000000"),
		"test",
		newURI("http://example.com"),
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		gtime,
		ctime,
		[]float64{1},
		map[int]int{1: 2},
		[]byte{1},
	}, l(),
		// n(),
		b(true),
		pi(1),
		ni(1),
		f(1.1),
		bf(newBigFloat("1.1", 2)),
		df(newDFloat("1.1")),
		bdf(newBDF("1.1")),
		bi(newBigInt("100000000000000000000")),
		s("test"),
		uri("http://example.com"),
		uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		gt(gtime),
		ct(ctime),
		l(), f(1), e(),
		m(), i(1), i(2), e(),
		bin([]byte{1}),
		e())
}

func TestBuilderInterfaceMap(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}
	gtime := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, loc)
	ctime := compact_time.AsCompactTime(gtime)

	assertBuild(t, map[interface{}]interface{}{
		1:  nil,
		2:  true,
		3:  1,
		4:  -1,
		5:  1.1,
		6:  newBigFloat("1.1", 2),
		7:  newDFloat("1.1"),
		8:  newBDF("1.1"),
		9:  newBigInt("100000000000000000000"),
		10: "test",
		11: newURI("http://example.com"),
		12: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		13: gtime,
		14: ctime,
		15: []float64{1},
		16: map[int]int{1: 2},
		17: []byte{1},
	},
		m(),
		i(1), n(),
		i(2), b(true),
		i(3), pi(1),
		i(4), ni(1),
		i(5), f(1.1),
		i(6), bf(newBigFloat("1.1", 2)),
		i(7), df(newDFloat("1.1")),
		i(8), bdf(newBDF("1.1")),
		i(9), bi(newBigInt("100000000000000000000")),
		i(10), s("test"),
		i(11), uri("http://example.com"),
		i(12), uuid([]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}),
		i(13), gt(gtime),
		i(14), ct(ctime),
		i(15), l(), f(1), e(),
		i(16), m(), i(1), i(2), e(),
		i(17), bin([]byte{1}),
		e())
}

// Older tests

type BuilderTestStruct struct {
	internal string
	ABool    bool
	AString  string
	AnInt    int
	AMap     map[int]int8
	ASlice   []string
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
