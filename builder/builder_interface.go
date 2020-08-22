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

var (
	builderIntfIntfMapType = reflect.TypeOf(map[interface{}]interface{}{})
	builderIntfSliceType   = reflect.TypeOf([]interface{}{})
	builderIntfType        = builderIntfSliceType.Elem()

	globalIntfBuilder = &intfBuilder{}
)

type intfBuilder struct {
	// Template Data
	session *Session

	// Instance Data
	root   *RootBuilder
	parent ObjectBuilder
	opts   *options.BuilderOptions
}

func newInterfaceBuilder() ObjectBuilder {
	return globalIntfBuilder
}

func (_this *intfBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *intfBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *intfBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &intfBuilder{
		session: _this.session,
		parent:  parent,
		root:    root,
		opts:    opts,
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

func (_this *intfBuilder) BuildFromString(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(string(value)))
}

func (_this *intfBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(string(value)))
}

func (_this *intfBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *intfBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *intfBuilder) BuildFromTypedArray(elemType reflect.Type, value []byte, dst reflect.Value) {
	switch elemType.Kind() {
	case reflect.Uint8:
		dst.Set(reflect.ValueOf(common.CloneBytes(value)))
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", elemType))
	}
}

func (_this *intfBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *intfBuilder) BuildBeginList() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForListContents()
}

func (_this *intfBuilder) BuildBeginMap() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *intfBuilder) BuildEndContainer() {
	PanicBadEventWithType(_this, builderIntfType, "ContainerEnd")
}

func (_this *intfBuilder) BuildBeginMarker(_ interface{}) {
	panic("TODO: intfBuilder.Marker")
}

func (_this *intfBuilder) BuildFromReference(_ interface{}) {
	panic("TODO: intfBuilder.Reference")
}

func (_this *intfBuilder) PrepareForListContents() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForListContents()
}

func (_this *intfBuilder) PrepareForMapContents() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceMap)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *intfBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.parent.NotifyChildContainerFinished(value)
}
