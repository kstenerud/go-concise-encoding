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

func BeginMarkerBuilder(parentBuilder ObjectBuilder, childBuilder ObjectBuilder, root *RootBuilder, options *BuilderOptions) {
	marker := &markerNameBuilder{
		root:    root,
		parent:  parentBuilder,
		options: options,
		child:   childBuilder,
	}

	root.SetCurrentBuilder(marker)
}

type markerNameBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions

	// Variable data (must be reset)
	child ObjectBuilder
	name  interface{}
}

func (_this *markerNameBuilder) IsContainerOnly() bool {
	return false
}

func (_this *markerNameBuilder) PostCacheInitBuilder() {
}

func (_this *markerNameBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return &markerNameBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
}

func (_this *markerNameBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *markerNameBuilder) prepareForMarkerObject(name interface{}) {
	mob := &markerObjectBuilder{
		root:           _this.root,
		parent:         _this.parent,
		options:        _this.options,
		name:           name,
		child:          _this.child,
		markerRegistry: _this.root.GetMarkerRegistry(),
	}
	_this.child.SetParent(mob)
	_this.root.SetCurrentBuilder(mob)
}

func (_this *markerNameBuilder) BuildFromNil(dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Nil")
}

func (_this *markerNameBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Bool")
}

func (_this *markerNameBuilder) BuildFromInt(value int64, dst reflect.Value) {
	if value < 0 {
		builderPanicBadEvent(_this, reflect.TypeOf(nil), "Int")
	}
	_this.prepareForMarkerObject(value)
}

func (_this *markerNameBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	_this.prepareForMarkerObject(value)
}

func (_this *markerNameBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	if isBigIntNegative(value) || !value.IsUint64() {
		builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigInt")
	}
	_this.prepareForMarkerObject(value.Uint64())
}

func (_this *markerNameBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Float")
}

func (_this *markerNameBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigFloat")
}

func (_this *markerNameBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "DecimalFloat")
}

func (_this *markerNameBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigDecimalFloat")
}

func (_this *markerNameBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "UUID")
}

func (_this *markerNameBuilder) BuildFromString(value string, dst reflect.Value) {
	_this.prepareForMarkerObject(value)
}

func (_this *markerNameBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Bytes")
}

func (_this *markerNameBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "URI")
}

func (_this *markerNameBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Time")
}

func (_this *markerNameBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "CompactTime")
}

func (_this *markerNameBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "List")
}

func (_this *markerNameBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Map")
}

func (_this *markerNameBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "End")
}

func (_this *markerNameBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: markerNameBuilder.Marker")
}

func (_this *markerNameBuilder) BuildFromReference(id interface{}) {
	panic("TODO: markerNameBuilder.Reference")
}

func (_this *markerNameBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "PrepareForListContents")
}

func (_this *markerNameBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "PrepareForMapContents")
}

func (_this *markerNameBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "NotifyChildContainerFinished")
}

// ============================================================================

type markerObjectBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions

	// Variable data (must be reset)
	name           interface{}
	child          ObjectBuilder
	markerRegistry *MarkerRegistry
}

func (_this *markerObjectBuilder) IsContainerOnly() bool {
	return false
}

func (_this *markerObjectBuilder) PostCacheInitBuilder() {
}

func (_this *markerObjectBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return &markerObjectBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
}

func (_this *markerObjectBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *markerObjectBuilder) BuildFromNil(dst reflect.Value) {
	_this.child.BuildFromNil(dst)
	_this.markerRegistry.NotifyMarker(_this.name, dst.Interface())
}

func (_this *markerObjectBuilder) BuildFromBool(value bool, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Bool")
}

func (_this *markerObjectBuilder) BuildFromInt(value int64, dst reflect.Value) {
	if value < 0 {
		builderPanicBadEvent(_this, reflect.TypeOf(nil), "Int")
	}
	panic("TODO: markerObjectBuilder.BuildFromInt")
}

func (_this *markerObjectBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	panic("TODO: markerObjectBuilder.BuildFromUint")
}

func (_this *markerObjectBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	if isBigIntNegative(value) {
		builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigInt")
	}
	panic("TODO: markerObjectBuilder.BuildFromBigInt")
}

func (_this *markerObjectBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Float")
}

func (_this *markerObjectBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigFloat")
}

func (_this *markerObjectBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "DecimalFloat")
}

func (_this *markerObjectBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "BigDecimalFloat")
}

func (_this *markerObjectBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "UUID")
}

func (_this *markerObjectBuilder) BuildFromString(value string, dst reflect.Value) {
	panic("TODO: markerObjectBuilder.BuildFromString")
}

func (_this *markerObjectBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Bytes")
}

func (_this *markerObjectBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "URI")
}

func (_this *markerObjectBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Time")
}

func (_this *markerObjectBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "CompactTime")
}

func (_this *markerObjectBuilder) BuildBeginList() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "List")
}

func (_this *markerObjectBuilder) BuildBeginMap() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "Map")
}

func (_this *markerObjectBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "End")
}

func (_this *markerObjectBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: markerObjectBuilder.Marker")
}

func (_this *markerObjectBuilder) BuildFromReference(id interface{}) {
	panic("TODO: markerObjectBuilder.Reference")
}

func (_this *markerObjectBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "PrepareForListContents")
}

func (_this *markerObjectBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "PrepareForMapContents")
}

func (_this *markerObjectBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, reflect.TypeOf(nil), "NotifyChildContainerFinished")
}
