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
)

// -------
// []uint8
// -------

type uint8SliceIterator struct {
	root *RootObjectIterator
}

func newUInt8SliceIterator() ObjectIterator {
	return &uint8SliceIterator{}
}

func (this *uint8SliceIterator) PostCacheInitIterator() {
}

func (this *uint8SliceIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &uint8SliceIterator{root: root}
}

func (this *uint8SliceIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnBytes(v.Bytes())
}

// ----
// Time
// ----

type timeIterator struct {
	root *RootObjectIterator
}

func newTimeIterator() ObjectIterator {
	return &timeIterator{}
}

func (this *timeIterator) PostCacheInitIterator() {
}

func (this *timeIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &timeIterator{root: root}
}

func (this *timeIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnTime(v.Interface().(time.Time))
}

// ----
// *URL
// ----

type pURLIterator struct {
	root *RootObjectIterator
}

func newPURLIterator() ObjectIterator {
	return &pURLIterator{}
}

func (this *pURLIterator) PostCacheInitIterator() {
}

func (this *pURLIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pURLIterator{root: root}
}

func (this *pURLIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
	} else {
		this.root.eventReceiver.OnURI(v.Interface().(*url.URL).String())
	}
}

// ---
// URL
// ---

type urlIterator struct {
	root *RootObjectIterator
}

func newURLIterator() ObjectIterator {
	return &urlIterator{}
}

func (this *urlIterator) PostCacheInitIterator() {
}

func (this *urlIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &urlIterator{root: root}
}

func (this *urlIterator) Iterate(v reflect.Value) {
	vCopy := v.Interface().(url.URL)
	this.root.eventReceiver.OnURI((&vCopy).String())
}

// --------
// *big.Int
// --------

type pBigIntIterator struct {
	root *RootObjectIterator
}

func newPBigIntIterator() ObjectIterator {
	return &pBigIntIterator{}
}

func (this *pBigIntIterator) PostCacheInitIterator() {
}

func (this *pBigIntIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pBigIntIterator{root: root}
}

func (this *pBigIntIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
	} else {
		this.root.eventReceiver.OnBigInt(v.Interface().(*big.Int))
	}
}

// -------
// big.Int
// -------

type bigIntIterator struct {
	root *RootObjectIterator
}

func newBigIntIterator() ObjectIterator {
	return &bigIntIterator{}
}

func (this *bigIntIterator) PostCacheInitIterator() {
}

func (this *bigIntIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &bigIntIterator{root: root}
}

func (this *bigIntIterator) Iterate(v reflect.Value) {
	vCopy := v.Interface().(big.Int)
	this.root.eventReceiver.OnBigInt((&vCopy))
}

// ----------
// *big.Float
// ----------

type pBigFloatIterator struct {
	root *RootObjectIterator
}

func newPBigFloatIterator() ObjectIterator {
	return &pBigFloatIterator{}
}

func (this *pBigFloatIterator) PostCacheInitIterator() {
}

func (this *pBigFloatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pBigFloatIterator{root: root}
}

func (this *pBigFloatIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
	} else {
		this.root.eventReceiver.OnBigFloat(v.Interface().(*big.Float))
	}
}

// ---------
// big.Float
// ---------

type bigFloatIterator struct {
	root *RootObjectIterator
}

func newBigFloatIterator() ObjectIterator {
	return &bigFloatIterator{}
}

func (this *bigFloatIterator) PostCacheInitIterator() {
}

func (this *bigFloatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &bigFloatIterator{root: root}
}

func (this *bigFloatIterator) Iterate(v reflect.Value) {
	vCopy := v.Interface().(big.Float)
	this.root.eventReceiver.OnBigFloat((&vCopy))
}

// ------------
// *apd.Decimal
// ------------

type pBigDecimalFloatIterator struct {
	root *RootObjectIterator
}

func newPBigDecimalFloatIterator() ObjectIterator {
	return &pBigDecimalFloatIterator{}
}

func (this *pBigDecimalFloatIterator) PostCacheInitIterator() {
}

func (this *pBigDecimalFloatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pBigDecimalFloatIterator{root: root}
}

func (this *pBigDecimalFloatIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
	} else {
		this.root.eventReceiver.OnBigDecimalFloat(v.Interface().(*apd.Decimal))
	}
}

// -----------
// apd.Decimal
// -----------

type bigDecimalFloatIterator struct {
	root *RootObjectIterator
}

func newBigDecimalFloatIterator() ObjectIterator {
	return &bigDecimalFloatIterator{}
}

func (this *bigDecimalFloatIterator) PostCacheInitIterator() {
}

func (this *bigDecimalFloatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &bigDecimalFloatIterator{root: root}
}

func (this *bigDecimalFloatIterator) Iterate(v reflect.Value) {
	vCopy := v.Interface().(apd.Decimal)
	this.root.eventReceiver.OnBigDecimalFloat((&vCopy))
}

// ------
// DFloat
// ------

type dfloatIterator struct {
	root *RootObjectIterator
}

func newDFloatIterator() ObjectIterator {
	return &dfloatIterator{}
}

func (this *dfloatIterator) PostCacheInitIterator() {
}

func (this *dfloatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &dfloatIterator{root: root}
}

func (this *dfloatIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnDecimalFloat(v.Interface().(compact_float.DFloat))
}

// ----
// Bool
// ----

type boolIterator struct {
	root *RootObjectIterator
}

func newBoolIterator() ObjectIterator {
	return &boolIterator{}
}

func (this *boolIterator) PostCacheInitIterator() {
}

func (this *boolIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &boolIterator{root: root}
}

func (this *boolIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnBool(v.Bool())
}

// ---
// Int
// ---

type intIterator struct {
	root *RootObjectIterator
}

func newIntIterator() ObjectIterator {
	return &intIterator{}
}

func (this *intIterator) PostCacheInitIterator() {
}

func (this *intIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &intIterator{root: root}
}

func (this *intIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnInt(v.Int())
}

// ----
// Uint
// ----

type uintIterator struct {
	root *RootObjectIterator
}

func newUintIterator() ObjectIterator {
	return &uintIterator{}
}

func (this *uintIterator) PostCacheInitIterator() {
}

func (this *uintIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &uintIterator{root: root}
}

func (this *uintIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnPositiveInt(v.Uint())
}

// -----
// Float
// -----

type floatIterator struct {
	root *RootObjectIterator
}

func newFloatIterator() ObjectIterator {
	return &floatIterator{}
}

func (this *floatIterator) PostCacheInitIterator() {
}

func (this *floatIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &floatIterator{root: root}
}

func (this *floatIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnFloat(v.Float())
}

// ------
// String
// ------

type stringIterator struct {
	root *RootObjectIterator
}

func newStringIterator() ObjectIterator {
	return &stringIterator{}
}

func (this *stringIterator) PostCacheInitIterator() {
}

func (this *stringIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &stringIterator{root: root}
}

func (this *stringIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnString(v.String())
}
