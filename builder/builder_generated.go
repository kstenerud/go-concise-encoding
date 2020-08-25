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

// Generated by github.com/kstenerud/go-concise-encoding/codegen
  // DO NOT EDIT
  // IF THIS LINE SHOWS UP IN THE GIT DIFF, THIS FILE HAS BEEN EDITED

package builder

import (
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// uint16

type uint16ArrayBuilder struct{}

func newUint16ArrayBuilder() ObjectBuilder       { return &uint16ArrayBuilder{} }
func (_this *uint16ArrayBuilder) String() string { return nameOf(_this) }
func (_this *uint16ArrayBuilder) badEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf(uint16(0)), name, args...)
}
func (_this *uint16ArrayBuilder) InitTemplate(_ *Session) {}
func (_this *uint16ArrayBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *uint16ArrayBuilder) SetParent(_ ObjectBuilder) {}
func (_this *uint16ArrayBuilder) BuildFromNil(_ reflect.Value) {
	_this.badEvent("BuildFromNil")
}
func (_this *uint16ArrayBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.badEvent("BuildFromBool")
}
func (_this *uint16ArrayBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.badEvent("BuildFromInt")
}
func (_this *uint16ArrayBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.badEvent("BuildFromUint")
}
func (_this *uint16ArrayBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.badEvent("BuildFromBigInt")
}
func (_this *uint16ArrayBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.badEvent("BuildFromFloat")
}
func (_this *uint16ArrayBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.badEvent("BuildFromBigFloat")
}
func (_this *uint16ArrayBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.badEvent("BuildFromDecimalFloat")
}
func (_this *uint16ArrayBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.badEvent("BuildFromBigDecimalFloat")
}
func (_this *uint16ArrayBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromUUID")
}
func (_this *uint16ArrayBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromString")
}
func (_this *uint16ArrayBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromVerbatimString")
}
func (_this *uint16ArrayBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	_this.badEvent("BuildFromURI")
}
func (_this *uint16ArrayBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomBinary")
}
func (_this *uint16ArrayBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomText")
}
func (_this *uint16ArrayBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromTime")
}
func (_this *uint16ArrayBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromCompactTime")
}
func (_this *uint16ArrayBuilder) BuildBeginList() {
	_this.badEvent("BuildBeginList")
}
func (_this *uint16ArrayBuilder) BuildBeginMap() {
	_this.badEvent("BuildBeginMap")
}
func (_this *uint16ArrayBuilder) BuildEndContainer() {
	_this.badEvent("BuildEndContainer")
}
func (_this *uint16ArrayBuilder) BuildBeginMarker(_ interface{}) {
	_this.badEvent("BuildBeginMarker")
}
func (_this *uint16ArrayBuilder) BuildFromReference(_ interface{}) {
	_this.badEvent("BuildFromReference")
}
func (_this *uint16ArrayBuilder) PrepareForListContents() {
	_this.badEvent("PrepareForListContents")
}
func (_this *uint16ArrayBuilder) PrepareForMapContents() {
	_this.badEvent("PrepareForMapContents")
}
func (_this *uint16ArrayBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.badEvent("NotifyChildContainerFinished")
}

// uint8

type uint8ArrayBuilder struct{}

func newUint8ArrayBuilder() ObjectBuilder       { return &uint8ArrayBuilder{} }
func (_this *uint8ArrayBuilder) String() string { return nameOf(_this) }
func (_this *uint8ArrayBuilder) badEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf(uint8(0)), name, args...)
}
func (_this *uint8ArrayBuilder) InitTemplate(_ *Session) {}
func (_this *uint8ArrayBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *uint8ArrayBuilder) SetParent(_ ObjectBuilder) {}
func (_this *uint8ArrayBuilder) BuildFromNil(_ reflect.Value) {
	_this.badEvent("BuildFromNil")
}
func (_this *uint8ArrayBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.badEvent("BuildFromBool")
}
func (_this *uint8ArrayBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.badEvent("BuildFromInt")
}
func (_this *uint8ArrayBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.badEvent("BuildFromUint")
}
func (_this *uint8ArrayBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.badEvent("BuildFromBigInt")
}
func (_this *uint8ArrayBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.badEvent("BuildFromFloat")
}
func (_this *uint8ArrayBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.badEvent("BuildFromBigFloat")
}
func (_this *uint8ArrayBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.badEvent("BuildFromDecimalFloat")
}
func (_this *uint8ArrayBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.badEvent("BuildFromBigDecimalFloat")
}
func (_this *uint8ArrayBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromUUID")
}
func (_this *uint8ArrayBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromString")
}
func (_this *uint8ArrayBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromVerbatimString")
}
func (_this *uint8ArrayBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	_this.badEvent("BuildFromURI")
}
func (_this *uint8ArrayBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomBinary")
}
func (_this *uint8ArrayBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomText")
}
func (_this *uint8ArrayBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromTime")
}
func (_this *uint8ArrayBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromCompactTime")
}
func (_this *uint8ArrayBuilder) BuildBeginList() {
	_this.badEvent("BuildBeginList")
}
func (_this *uint8ArrayBuilder) BuildBeginMap() {
	_this.badEvent("BuildBeginMap")
}
func (_this *uint8ArrayBuilder) BuildEndContainer() {
	_this.badEvent("BuildEndContainer")
}
func (_this *uint8ArrayBuilder) BuildBeginMarker(_ interface{}) {
	_this.badEvent("BuildBeginMarker")
}
func (_this *uint8ArrayBuilder) BuildFromReference(_ interface{}) {
	_this.badEvent("BuildFromReference")
}
func (_this *uint8ArrayBuilder) PrepareForListContents() {
	_this.badEvent("PrepareForListContents")
}
func (_this *uint8ArrayBuilder) PrepareForMapContents() {
	_this.badEvent("PrepareForMapContents")
}
func (_this *uint8ArrayBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.badEvent("NotifyChildContainerFinished")
}
