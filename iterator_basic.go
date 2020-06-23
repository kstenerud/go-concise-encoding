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
	"reflect"
	"time"

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

func (_this *uint8SliceIterator) PostCacheInitIterator() {
}

func (_this *uint8SliceIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnBytes(v.Bytes())
}

// ----
// Time
// ----

type timeIterator struct {
}

func newTimeIterator() ObjectIterator {
	return &timeIterator{}
}

func (_this *timeIterator) PostCacheInitIterator() {
}

func (_this *timeIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnTime(v.Interface().(time.Time))
}

// ------------
// Compact Time
// ------------

type compactTimeIterator struct {
}

func newCompactTimeIterator() ObjectIterator {
	return &compactTimeIterator{}
}

func (_this *compactTimeIterator) PostCacheInitIterator() {
}

func (_this *compactTimeIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	ct := v.Interface().(compact_time.Time)
	root.eventReceiver.OnCompactTime(&ct)
}

type pCompactTimeIterator struct {
}

func newPCompactTimeIterator() ObjectIterator {
	return &pCompactTimeIterator{}
}

func (_this *pCompactTimeIterator) PostCacheInitIterator() {
}

func (_this *pCompactTimeIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		root.eventReceiver.OnCompactTime(v.Interface().(*compact_time.Time))
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

func (_this *pURLIterator) PostCacheInitIterator() {
}

func (_this *pURLIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		root.eventReceiver.OnURI(v.Interface().(*url.URL).String())
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

func (_this *urlIterator) PostCacheInitIterator() {
}

func (_this *urlIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	vCopy := v.Interface().(url.URL)
	root.eventReceiver.OnURI((&vCopy).String())
}

// --------
// *big.Int
// --------

type pBigIntIterator struct {
}

func newPBigIntIterator() ObjectIterator {
	return &pBigIntIterator{}
}

func (_this *pBigIntIterator) PostCacheInitIterator() {
}

func (_this *pBigIntIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		root.eventReceiver.OnBigInt(v.Interface().(*big.Int))
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

func (_this *bigIntIterator) PostCacheInitIterator() {
}

func (_this *bigIntIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	vCopy := v.Interface().(big.Int)
	root.eventReceiver.OnBigInt((&vCopy))
}

// ----------
// *big.Float
// ----------

type pBigFloatIterator struct {
}

func newPBigFloatIterator() ObjectIterator {
	return &pBigFloatIterator{}
}

func (_this *pBigFloatIterator) PostCacheInitIterator() {
}

func (_this *pBigFloatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		root.eventReceiver.OnBigFloat(v.Interface().(*big.Float))
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

func (_this *bigFloatIterator) PostCacheInitIterator() {
}

func (_this *bigFloatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	vCopy := v.Interface().(big.Float)
	root.eventReceiver.OnBigFloat((&vCopy))
}

// ------------
// *apd.Decimal
// ------------

type pBigDecimalFloatIterator struct {
}

func newPBigDecimalFloatIterator() ObjectIterator {
	return &pBigDecimalFloatIterator{}
}

func (_this *pBigDecimalFloatIterator) PostCacheInitIterator() {
}

func (_this *pBigDecimalFloatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		root.eventReceiver.OnBigDecimalFloat(v.Interface().(*apd.Decimal))
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

func (_this *bigDecimalFloatIterator) PostCacheInitIterator() {
}

func (_this *bigDecimalFloatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	vCopy := v.Interface().(apd.Decimal)
	root.eventReceiver.OnBigDecimalFloat((&vCopy))
}

// ------
// DFloat
// ------

type dfloatIterator struct {
}

func newDFloatIterator() ObjectIterator {
	return &dfloatIterator{}
}

func (_this *dfloatIterator) PostCacheInitIterator() {
}

func (_this *dfloatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnDecimalFloat(v.Interface().(compact_float.DFloat))
}

// ----
// Bool
// ----

type boolIterator struct {
}

func newBoolIterator() ObjectIterator {
	return &boolIterator{}
}

func (_this *boolIterator) PostCacheInitIterator() {
}

func (_this *boolIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnBool(v.Bool())
}

// ---
// Int
// ---

type intIterator struct {
}

func newIntIterator() ObjectIterator {
	return &intIterator{}
}

func (_this *intIterator) PostCacheInitIterator() {
}

func (_this *intIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnInt(v.Int())
}

// ----
// Uint
// ----

type uintIterator struct {
}

func newUintIterator() ObjectIterator {
	return &uintIterator{}
}

func (_this *uintIterator) PostCacheInitIterator() {
}

func (_this *uintIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnPositiveInt(v.Uint())
}

// -----
// Float
// -----

type floatIterator struct {
}

func newFloatIterator() ObjectIterator {
	return &floatIterator{}
}

func (_this *floatIterator) PostCacheInitIterator() {
}

func (_this *floatIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnFloat(v.Float())
}

// ------
// String
// ------

type stringIterator struct {
}

func newStringIterator() ObjectIterator {
	return &stringIterator{}
}

func (_this *stringIterator) PostCacheInitIterator() {
}

func (_this *stringIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnString(v.String())
}
