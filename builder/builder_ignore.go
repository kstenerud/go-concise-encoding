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

type ignoreBuilder struct {
	// Instance Data
	root   *RootBuilder
	parent ObjectBuilder
	opts   *options.BuilderOptions
}

var globalIgnoreBuilder = &ignoreBuilder{}

func newIgnoreBuilder() ObjectBuilder {
	return globalIgnoreBuilder
}

func (_this *ignoreBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *ignoreBuilder) InitTemplate(_ *Session) {
}

func (_this *ignoreBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &ignoreBuilder{
		parent: parent,
		root:   root,
		opts:   opts,
	}
}

func (_this *ignoreBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *ignoreBuilder) BuildFromNil(_ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromTypedArray(_ reflect.Type, _ []byte, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreBuilder) BuildBeginList() {
	builder := newIgnoreContainerBuilder().NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForListContents()
}

func (_this *ignoreBuilder) BuildBeginMap() {
	builder := newIgnoreContainerBuilder().NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *ignoreBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "End")
}

func (_this *ignoreBuilder) BuildBeginMarker(_ interface{}) {
	panic("TODO: ignoreBuilder.Marker")
}

func (_this *ignoreBuilder) BuildFromReference(_ interface{}) {
	// Ignore this directive
}

func (_this *ignoreBuilder) PrepareForListContents() {
	PanicBadEvent(_this, "PrepareForListContents")
}

func (_this *ignoreBuilder) PrepareForMapContents() {
	PanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *ignoreBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEvent(_this, "NotifyChildContainerFinished")
}

// ============================================================================

type ignoreContainerBuilder struct {
	// Instance Data
	root   *RootBuilder
	parent ObjectBuilder
	opts   *options.BuilderOptions
}

var globalIgnoreContainerBuilder = &ignoreContainerBuilder{}

func newIgnoreContainerBuilder() ObjectBuilder {
	return globalIgnoreContainerBuilder
}

func (_this *ignoreContainerBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *ignoreContainerBuilder) InitTemplate(_ *Session) {
}

func (_this *ignoreContainerBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &ignoreContainerBuilder{
		parent: parent,
		root:   root,
		opts:   opts,
	}
}

func (_this *ignoreContainerBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *ignoreContainerBuilder) BuildFromNil(_ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromTypedArray(_ reflect.Type, _ []byte, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) BuildBeginList() {
	builder := newIgnoreContainerBuilder().NewInstance(_this.root, _this, _this.opts)
	builder.PrepareForListContents()
}

func (_this *ignoreContainerBuilder) BuildBeginMap() {
	builder := newIgnoreContainerBuilder().NewInstance(_this.root, _this, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *ignoreContainerBuilder) BuildEndContainer() {
	_this.root.SetCurrentBuilder(_this.parent)
}

func (_this *ignoreContainerBuilder) BuildBeginMarker(_ interface{}) {
	panic("TODO: ignoreContainerBuilder.Marker")
}

func (_this *ignoreContainerBuilder) BuildFromReference(_ interface{}) {
	// Ignore this directive
}

func (_this *ignoreContainerBuilder) PrepareForListContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *ignoreContainerBuilder) PrepareForMapContents() {
	_this.root.SetCurrentBuilder(_this)
}

func (_this *ignoreContainerBuilder) NotifyChildContainerFinished(_ reflect.Value) {
}
