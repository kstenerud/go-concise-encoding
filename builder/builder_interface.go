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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

var (
	builderIntfIntfMapType = reflect.TypeOf(map[interface{}]interface{}{})
	builderIntfSliceType   = reflect.TypeOf([]interface{}{})
	builderIntfType        = builderIntfSliceType.Elem()

	globalIntfBuilder        = &intfBuilder{}
	globalIntfSliceBuilder   = &intfSliceBuilder{}
	globalIntfIntfMapBuilder = &intfIntfMapBuilder{}
)

type intfBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions
}

func newInterfaceBuilder() ObjectBuilder {
	return globalIntfBuilder
}

func (_this *intfBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *intfBuilder) IsContainerOnly() bool {
	return false
}

func (_this *intfBuilder) PostCacheInitBuilder() {
}

func (_this *intfBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return &intfBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
}

func (_this *intfBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *intfBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(builderIntfType))
}

func (_this *intfBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromString(value string, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildBeginList() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForListContents()
}

func (_this *intfBuilder) BuildBeginMap() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForMapContents()
}

func (_this *intfBuilder) BuildEndContainer() {
	builderPanicBadEventType(_this, builderIntfType, "ContainerEnd")
}

func (_this *intfBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: intfBuilder.Marker")
}

func (_this *intfBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfBuilder.Reference")
}

func (_this *intfBuilder) PrepareForListContents() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForListContents()
}

func (_this *intfBuilder) PrepareForMapContents() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForMapContents()
}

func (_this *intfBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.parent.NotifyChildContainerFinished(value)
}

// ============================================================================

type intfSliceBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions

	// Variable data (must be reset)
	container reflect.Value
}

func newIntfSliceBuilder() ObjectBuilder {
	return globalIntfSliceBuilder
}

func (_this *intfSliceBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *intfSliceBuilder) IsContainerOnly() bool {
	return true
}

func (_this *intfSliceBuilder) PostCacheInitBuilder() {
}

func (_this *intfSliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &intfSliceBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
	that.reset()
	return that
}

func (_this *intfSliceBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *intfSliceBuilder) reset() {
	_this.container = reflect.MakeSlice(builderIntfSliceType, 0, defaultSliceCap)
}

func (_this *intfSliceBuilder) storeRValue(value reflect.Value) {
	_this.container = reflect.Append(_this.container, value)
}

func (_this *intfSliceBuilder) storeValue(value interface{}) {
	_this.storeRValue(reflect.ValueOf(value))
}

func (_this *intfSliceBuilder) BuildFromNil(ignored reflect.Value) {
	_this.storeRValue(reflect.New(builderIntfSliceType).Elem())
}

func (_this *intfSliceBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromString(value string, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	_this.storeValue(value)
}

func (_this *intfSliceBuilder) BuildBeginList() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForListContents()
}

func (_this *intfSliceBuilder) BuildBeginMap() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForMapContents()
}

func (_this *intfSliceBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *intfSliceBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: intfSliceBuilder.Marker")
}

func (_this *intfSliceBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfSliceBuilder.Reference")
}

func (_this *intfSliceBuilder) PrepareForListContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *intfSliceBuilder) PrepareForMapContents() {
	builderPanicBadEventType(_this, builderIntfType, "PrepareForMapContents")
}

func (_this *intfSliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.storeRValue(value)
}

// ============================================================================

type intfIntfMapBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions

	// Variable data (must be reset)
	container reflect.Value
	key       reflect.Value
	nextIsKey bool
}

func newIntfIntfMapBuilder() ObjectBuilder {
	return globalIntfIntfMapBuilder
}

func (_this *intfIntfMapBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *intfIntfMapBuilder) IsContainerOnly() bool {
	return true
}

func (_this *intfIntfMapBuilder) PostCacheInitBuilder() {
}

func (_this *intfIntfMapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	that := &intfIntfMapBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
	that.reset()
	return that
}

func (_this *intfIntfMapBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *intfIntfMapBuilder) reset() {
	_this.container = reflect.MakeMap(builderIntfIntfMapType)
	_this.key = reflect.Value{}
	_this.nextIsKey = true
}

func (_this *intfIntfMapBuilder) storeValue(value reflect.Value) {
	if _this.nextIsKey {
		_this.key = value
	} else {
		_this.container.SetMapIndex(_this.key, value)
	}
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *intfIntfMapBuilder) BuildFromNil(ignored reflect.Value) {
	_this.storeValue(reflect.Zero(builderIntfType))
}

func (_this *intfIntfMapBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromString(value string, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildFromCompactTime(value *compact_time.Time, ignored reflect.Value) {
	_this.storeValue(reflect.ValueOf(value))
}

func (_this *intfIntfMapBuilder) BuildBeginList() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForListContents()
}

func (_this *intfIntfMapBuilder) BuildBeginMap() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForMapContents()
}

func (_this *intfIntfMapBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *intfIntfMapBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: intfIntfMapBuilder.Marker")
}

func (_this *intfIntfMapBuilder) BuildFromReference(id interface{}) {
	panic("TODO: intfIntfMapBuilder.Reference")
}

func (_this *intfIntfMapBuilder) PrepareForListContents() {
	builderPanicBadEventType(_this, builderIntfType, "PrepareForListContents")
}

func (_this *intfIntfMapBuilder) PrepareForMapContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *intfIntfMapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.storeValue(value)
}
