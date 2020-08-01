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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// -------
// []uint8
// -------

type uint8SliceIterator struct {
}

func newUInt8SliceIterator() ObjectIterator {
	return &uint8SliceIterator{}
}

func (_this *uint8SliceIterator) PostCacheInitIterator(session *Session) {
}

func (_this *uint8SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	eventReceiver.OnBytes(v.Bytes())
}

// ----
// Time
// ----

type timeIterator struct {
}

func newTimeIterator() ObjectIterator {
	return &timeIterator{}
}

func (_this *timeIterator) PostCacheInitIterator(session *Session) {
}

func (_this *timeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *compactTimeIterator) PostCacheInitIterator(session *Session) {
}

func (_this *compactTimeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	ct := v.Interface().(compact_time.Time)
	eventReceiver.OnCompactTime(&ct)
}

type pCompactTimeIterator struct {
}

func newPCompactTimeIterator() ObjectIterator {
	return &pCompactTimeIterator{}
}

func (_this *pCompactTimeIterator) PostCacheInitIterator(session *Session) {
}

func (_this *pCompactTimeIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
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

func (_this *pURLIterator) PostCacheInitIterator(session *Session) {
}

func (_this *pURLIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
	} else {
		eventReceiver.OnURI([]byte(v.Interface().(*url.URL).String()))
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

func (_this *urlIterator) PostCacheInitIterator(session *Session) {
}

func (_this *urlIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	vCopy := v.Interface().(url.URL)
	eventReceiver.OnURI([]byte((&vCopy).String()))
}

// --------
// *big.Int
// --------

type pBigIntIterator struct {
}

func newPBigIntIterator() ObjectIterator {
	return &pBigIntIterator{}
}

func (_this *pBigIntIterator) PostCacheInitIterator(session *Session) {
}

func (_this *pBigIntIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
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

func (_this *bigIntIterator) PostCacheInitIterator(session *Session) {
}

func (_this *bigIntIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	vCopy := v.Interface().(big.Int)
	eventReceiver.OnBigInt((&vCopy))
}

// ----------
// *big.Float
// ----------

type pBigFloatIterator struct {
}

func newPBigFloatIterator() ObjectIterator {
	return &pBigFloatIterator{}
}

func (_this *pBigFloatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *pBigFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
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

func (_this *bigFloatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *bigFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	vCopy := v.Interface().(big.Float)
	eventReceiver.OnBigFloat((&vCopy))
}

// ------------
// *apd.Decimal
// ------------

type pBigDecimalFloatIterator struct {
}

func newPBigDecimalFloatIterator() ObjectIterator {
	return &pBigDecimalFloatIterator{}
}

func (_this *pBigDecimalFloatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *pBigDecimalFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
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

func (_this *bigDecimalFloatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *bigDecimalFloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	vCopy := v.Interface().(apd.Decimal)
	eventReceiver.OnBigDecimalFloat((&vCopy))
}

// ------
// DFloat
// ------

type dfloatIterator struct {
}

func newDFloatIterator() ObjectIterator {
	return &dfloatIterator{}
}

func (_this *dfloatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *dfloatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *boolIterator) PostCacheInitIterator(session *Session) {
}

func (_this *boolIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *intIterator) PostCacheInitIterator(session *Session) {
}

func (_this *intIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *uintIterator) PostCacheInitIterator(session *Session) {
}

func (_this *uintIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *floatIterator) PostCacheInitIterator(session *Session) {
}

func (_this *floatIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
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

func (_this *stringIterator) PostCacheInitIterator(session *Session) {
}

func (_this *stringIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	eventReceiver.OnString([]byte(v.String()))
}
