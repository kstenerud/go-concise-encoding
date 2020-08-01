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
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// The direct builder has an unambiguous direct mapping from build event to
// a non-pointer destination type (for example, a bool is always a bool).
type directBuilder struct {
	// Static data
	session *Session
	dstType reflect.Type
}

func newDirectBuilder(dstType reflect.Type) ObjectBuilder {
	return &directBuilder{
		dstType: dstType,
	}
}

func (_this *directBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *directBuilder) IsContainerOnly() bool {
	return false
}

func (_this *directBuilder) PostCacheInitBuilder(session *Session) {
	_this.session = session
}

func (_this *directBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *directBuilder) BuildFromNil(dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Nil")
}

func (_this *directBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.SetBool(value)
}

func (_this *directBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Int")
}

func (_this *directBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Uint")
}

func (_this *directBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigInt")
}

func (_this *directBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Float")
}

func (_this *directBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigFloat")
}

func (_this *directBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "DecimalFloat")
}

func (_this *directBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigDecimalFloat")
}

func (_this *directBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directBuilder) BuildFromString(value []byte, dst reflect.Value) {
	dst.SetString(string(value))
}

func (_this *directBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	dst.SetString(string(value))
}

func (_this *directBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bytes")
}

func (_this *directBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *directBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *directBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value).Elem())
}

func (_this *directBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "List")
}

func (_this *directBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Map")
}

func (_this *directBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "ContainerEnd")
}

func (_this *directBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: directBuilder.BuildBeginMarker")
}

func (_this *directBuilder) BuildFromReference(id interface{}) {
	panic("TODO: directBuilder.BuildFromReference")
}

func (_this *directBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForListContents")
}

func (_this *directBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *directBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "NotifyChildContainerFinished")
}

// ============================================================================

// The direct builder has an unambiguous direct mapping from build event to
// a pointer destination type (for example, a *url is always a *url).
type directPtrBuilder struct {
	// Const data
	dstType reflect.Type
}

func newDirectPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &directPtrBuilder{
		dstType: dstType,
	}
}

func (_this *directPtrBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *directPtrBuilder) IsContainerOnly() bool {
	return false
}

func (_this *directPtrBuilder) PostCacheInitBuilder(session *Session) {
}

func (_this *directPtrBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directPtrBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *directPtrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(dst.Type()))
}

func (_this *directPtrBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Bool")
}

func (_this *directPtrBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Int")
}

func (_this *directPtrBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Uint")
}

func (_this *directPtrBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigInt")
}

func (_this *directPtrBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Float")
}

func (_this *directPtrBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigFloat")
}

func (_this *directPtrBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "DecimalFloat")
}

func (_this *directPtrBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "BigDecimalFloat")
}

func (_this *directPtrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (_this *directPtrBuilder) BuildFromString(value []byte, dst reflect.Value) {
	// String needs special handling since there's no such thing as a nil string
	// in go.
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "String")
}

func (_this *directPtrBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	// String needs special handling since there's no such thing as a nil string
	// in go.
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "VerbatimString")
}

func (_this *directPtrBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (_this *directPtrBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "CustomBinary")
}

func (_this *directPtrBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "CustomText")
}

func (_this *directPtrBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directPtrBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	// TODO: Should non-pointer stuff be here?
	dst.Set(reflect.ValueOf(value))
}

func (_this *directPtrBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directPtrBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "List")
}

func (_this *directPtrBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Map")
}

func (_this *directPtrBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "ContainerEnd")
}

func (_this *directPtrBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Marker")
}

func (_this *directPtrBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "Reference")
}

func (_this *directPtrBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForListContents")
}

func (_this *directPtrBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *directPtrBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "NotifyChildContainerFinished")
}
