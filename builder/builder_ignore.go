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

type ignoreBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions
}

var globalIgnoreBuilder = &ignoreBuilder{}

func newIgnoreBuilder() ObjectBuilder {
	return globalIgnoreBuilder
}

func (_this *ignoreBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *ignoreBuilder) IsContainerOnly() bool {
	return false
}

func (_this *ignoreBuilder) PostCacheInitBuilder() {
}

func (_this *ignoreBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return &ignoreBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
}

func (_this *ignoreBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *ignoreBuilder) BuildFromNil(dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBool(value bool, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromInt(value int64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromString(value string, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildBeginList() {
	builder := newIgnoreContainerBuilder().CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForListContents()
}

func (_this *ignoreBuilder) BuildBeginMap() {
	builder := newIgnoreContainerBuilder().CloneFromTemplate(_this.root, _this.parent, _this.options)
	builder.PrepareForMapContents()
}

func (_this *ignoreBuilder) BuildEndContainer() {
	builderPanicBadEvent(_this, "End")
}

func (_this *ignoreBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: ignoreBuilder.Marker")
}

func (_this *ignoreBuilder) BuildFromReference(id interface{}) {
	// Ignore this directive
}

func (_this *ignoreBuilder) PrepareForListContents() {
	builderPanicBadEvent(_this, "PrepareForListContents")
}

func (_this *ignoreBuilder) PrepareForMapContents() {
	builderPanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *ignoreBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(_this, "NotifyChildContainerFinished")
}

// ============================================================================

type ignoreContainerBuilder struct {
	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *BuilderOptions
}

var globalIgnoreContainerBuilder = &ignoreContainerBuilder{}

func newIgnoreContainerBuilder() ObjectBuilder {
	return globalIgnoreContainerBuilder
}

func (_this *ignoreContainerBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *ignoreContainerBuilder) IsContainerOnly() bool {
	return true
}

func (_this *ignoreContainerBuilder) PostCacheInitBuilder() {
}

func (_this *ignoreContainerBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	return &ignoreContainerBuilder{
		parent:  parent,
		root:    root,
		options: options,
	}
}

func (_this *ignoreContainerBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *ignoreContainerBuilder) BuildFromNil(dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBool(value bool, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromInt(value int64, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromString(value string, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildBeginList() {
	builder := newIgnoreContainerBuilder().CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForListContents()
}

func (_this *ignoreContainerBuilder) BuildBeginMap() {
	builder := newIgnoreContainerBuilder().CloneFromTemplate(_this.root, _this, _this.options)
	builder.PrepareForMapContents()
}

func (_this *ignoreContainerBuilder) BuildEndContainer() {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreContainerBuilder) BuildBeginMarker(id interface{}) {
	panic("TODO: ignoreContainerBuilder.Marker")
}

func (_this *ignoreContainerBuilder) BuildFromReference(id interface{}) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) PrepareForListContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *ignoreContainerBuilder) PrepareForMapContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *ignoreContainerBuilder) NotifyChildContainerFinished(value reflect.Value) {
}
