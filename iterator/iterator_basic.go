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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// ----
// Time
// ----

type timeIterator struct {
}

func newTimeIterator() ObjectIterator {
	return &timeIterator{}
}

func (_this *timeIterator) InitTemplate(_ FetchIterator) {
}

func (_this *timeIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *timeIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *timeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnTime(v.Interface().(time.Time))
}

// ------------
// Compact Time
// ------------

type compactTimeIterator struct {
}

func newCompactTimeIterator() ObjectIterator {
	return &compactTimeIterator{}
}

func (_this *compactTimeIterator) InitTemplate(_ FetchIterator) {
}

func (_this *compactTimeIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *compactTimeIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *compactTimeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	ct := v.Interface().(compact_time.Time)
	eventReceiver.OnCompactTime(&ct)
}

type pCompactTimeIterator struct {
}

func newPCompactTimeIterator() ObjectIterator {
	return &pCompactTimeIterator{}
}

func (_this *pCompactTimeIterator) InitTemplate(_ FetchIterator) {
}

func (_this *pCompactTimeIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *pCompactTimeIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *pCompactTimeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.IsNil() {
		eventReceiver.OnNull()
	} else {
		eventReceiver.OnCompactTime(v.Interface().(*compact_time.Time))
	}
}

// ----
// *URL
// ----

type pURLIterator struct {
}

func newPURLIterator() ObjectIterator {
	return &pURLIterator{}
}

func (_this *pURLIterator) InitTemplate(_ FetchIterator) {
}

func (_this *pURLIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *pURLIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *pURLIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.IsNil() {
		eventReceiver.OnNull()
	} else {
		bytes := []byte(v.Interface().(*url.URL).String())
		eventReceiver.OnArray(events.ArrayTypeURI, uint64(len(bytes)), bytes)
	}
}

// ---
// URL
// ---

type urlIterator struct {
}

func newURLIterator() ObjectIterator {
	return &urlIterator{}
}

func (_this *urlIterator) InitTemplate(_ FetchIterator) {
}

func (_this *urlIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *urlIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *urlIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	vCopy := v.Interface().(url.URL)
	bytes := []byte((&vCopy).String())
	eventReceiver.OnArray(events.ArrayTypeURI, uint64(len(bytes)), bytes)
}

// --------
// *big.Int
// --------

type pBigIntIterator struct {
}

func newPBigIntIterator() ObjectIterator {
	return &pBigIntIterator{}
}

func (_this *pBigIntIterator) InitTemplate(_ FetchIterator) {
}

func (_this *pBigIntIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *pBigIntIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *pBigIntIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.IsNil() {
		eventReceiver.OnNull()
	} else {
		eventReceiver.OnBigInt(v.Interface().(*big.Int))
	}
}

// -------
// big.Int
// -------

type bigIntIterator struct {
}

func newBigIntIterator() ObjectIterator {
	return &bigIntIterator{}
}

func (_this *bigIntIterator) InitTemplate(_ FetchIterator) {
}

func (_this *bigIntIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *bigIntIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *bigIntIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	vCopy := v.Interface().(big.Int)
	eventReceiver.OnBigInt(&vCopy)
}

// ----------
// *big.Float
// ----------

type pBigFloatIterator struct {
}

func newPBigFloatIterator() ObjectIterator {
	return &pBigFloatIterator{}
}

func (_this *pBigFloatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *pBigFloatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *pBigFloatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *pBigFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.IsNil() {
		eventReceiver.OnNull()
	} else {
		eventReceiver.OnBigFloat(v.Interface().(*big.Float))
	}
}

// ---------
// big.Float
// ---------

type bigFloatIterator struct {
}

func newBigFloatIterator() ObjectIterator {
	return &bigFloatIterator{}
}

func (_this *bigFloatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *bigFloatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *bigFloatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *bigFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	vCopy := v.Interface().(big.Float)
	eventReceiver.OnBigFloat(&vCopy)
}

// ------------
// *apd.Decimal
// ------------

type pBigDecimalFloatIterator struct {
}

func newPBigDecimalFloatIterator() ObjectIterator {
	return &pBigDecimalFloatIterator{}
}

func (_this *pBigDecimalFloatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *pBigDecimalFloatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *pBigDecimalFloatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *pBigDecimalFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.IsNil() {
		eventReceiver.OnNull()
	} else {
		eventReceiver.OnBigDecimalFloat(v.Interface().(*apd.Decimal))
	}
}

// -----------
// apd.Decimal
// -----------

type bigDecimalFloatIterator struct {
}

func newBigDecimalFloatIterator() ObjectIterator {
	return &bigDecimalFloatIterator{}
}

func (_this *bigDecimalFloatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *bigDecimalFloatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *bigDecimalFloatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *bigDecimalFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	vCopy := v.Interface().(apd.Decimal)
	eventReceiver.OnBigDecimalFloat(&vCopy)
}

// ------
// DFloat
// ------

type dfloatIterator struct {
}

func newDFloatIterator() ObjectIterator {
	return &dfloatIterator{}
}

func (_this *dfloatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *dfloatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *dfloatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *dfloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnDecimalFloat(v.Interface().(compact_float.DFloat))
}

// ----
// Bool
// ----

type boolIterator struct {
}

func newBoolIterator() ObjectIterator {
	return &boolIterator{}
}

func (_this *boolIterator) InitTemplate(_ FetchIterator) {
}

func (_this *boolIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *boolIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *boolIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnBool(v.Bool())
}

// ---
// Int
// ---

type intIterator struct {
}

func newIntIterator() ObjectIterator {
	return &intIterator{}
}

func (_this *intIterator) InitTemplate(_ FetchIterator) {
}

func (_this *intIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *intIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *intIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnInt(v.Int())
}

// ----
// Uint
// ----

type uintIterator struct {
}

func newUintIterator() ObjectIterator {
	return &uintIterator{}
}

func (_this *uintIterator) InitTemplate(_ FetchIterator) {
}

func (_this *uintIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *uintIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *uintIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnPositiveInt(v.Uint())
}

// -----
// Float
// -----

type floatIterator struct {
}

func newFloatIterator() ObjectIterator {
	return &floatIterator{}
}

func (_this *floatIterator) InitTemplate(_ FetchIterator) {
}

func (_this *floatIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *floatIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *floatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnFloat(v.Float())
}

// ------
// String
// ------

type stringIterator struct {
}

func newStringIterator() ObjectIterator {
	return &stringIterator{}
}

func (_this *stringIterator) InitTemplate(_ FetchIterator) {
}

func (_this *stringIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *stringIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {
}

func (_this *stringIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	bytes := []byte(v.String())
	eventReceiver.OnArray(events.ArrayTypeString, uint64(len(bytes)), bytes)
}
