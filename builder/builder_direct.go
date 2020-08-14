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

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// The direct builder has an unambiguous direct mapping from build event to
// a non-pointer destination type (for example, a bool is always a bool).
type directBuilder struct {
	// Template Data
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

func (_this *directBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *directBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *directBuilder) BuildFromNil(_ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Nil")
}

func (_this *directBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.SetBool(value)
}

func (_this *directBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Int")
}

func (_this *directBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Uint")
}

func (_this *directBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigInt")
}

func (_this *directBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Float")
}

func (_this *directBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigFloat")
}

func (_this *directBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "DecimalFloat")
}

func (_this *directBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigDecimalFloat")
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

func (_this *directBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Bytes")
}

func (_this *directBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *directBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
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
	PanicBadEventWithType(_this, _this.dstType, "List")
}

func (_this *directBuilder) BuildBeginMap() {
	PanicBadEventWithType(_this, _this.dstType, "Map")
}

func (_this *directBuilder) BuildEndContainer() {
	PanicBadEventWithType(_this, _this.dstType, "ContainerEnd")
}

func (_this *directBuilder) BuildBeginMarker(_ interface{}) {
	panic("TODO: directBuilder.BuildBeginMarker")
}

func (_this *directBuilder) BuildFromReference(_ interface{}) {
	panic("TODO: directBuilder.BuildFromReference")
}

func (_this *directBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, _this.dstType, "PrepareForListContents")
}

func (_this *directBuilder) PrepareForMapContents() {
	PanicBadEventWithType(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *directBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "NotifyChildContainerFinished")
}

// ============================================================================

// The direct builder has an unambiguous direct mapping from build event to
// a pointer destination type (for example, a *url is always a *url).
type directPtrBuilder struct {
	// Template Data
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

func (_this *directPtrBuilder) InitTemplate(_ *Session) {
}

func (_this *directPtrBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directPtrBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *directPtrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(dst.Type()))
}

func (_this *directPtrBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Bool")
}

func (_this *directPtrBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Int")
}

func (_this *directPtrBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Uint")
}

func (_this *directPtrBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigInt")
}

func (_this *directPtrBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "Float")
}

func (_this *directPtrBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigFloat")
}

func (_this *directPtrBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "DecimalFloat")
}

func (_this *directPtrBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "BigDecimalFloat")
}

func (_this *directPtrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (_this *directPtrBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	// String needs special handling since there's no such thing as a nil string
	// in go.
	PanicBadEventWithType(_this, _this.dstType, "String")
}

func (_this *directPtrBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	// String needs special handling since there's no such thing as a nil string
	// in go.
	PanicBadEventWithType(_this, _this.dstType, "VerbatimString")
}

func (_this *directPtrBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	dst.SetBytes(common.CloneBytes(value))
}

func (_this *directPtrBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "CustomBinary")
}

func (_this *directPtrBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "CustomText")
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
	PanicBadEventWithType(_this, _this.dstType, "List")
}

func (_this *directPtrBuilder) BuildBeginMap() {
	PanicBadEventWithType(_this, _this.dstType, "Map")
}

func (_this *directPtrBuilder) BuildEndContainer() {
	PanicBadEventWithType(_this, _this.dstType, "ContainerEnd")
}

func (_this *directPtrBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEventWithType(_this, _this.dstType, "Marker")
}

func (_this *directPtrBuilder) BuildFromReference(_ interface{}) {
	PanicBadEventWithType(_this, _this.dstType, "Reference")
}

func (_this *directPtrBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, _this.dstType, "PrepareForListContents")
}

func (_this *directPtrBuilder) PrepareForMapContents() {
	PanicBadEventWithType(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *directPtrBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEventWithType(_this, _this.dstType, "NotifyChildContainerFinished")
}
