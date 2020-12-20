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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

var globalIntfBuilder = &interfaceBuilder{}

type interfaceBuilder struct {
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

func (_this *interfaceBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *interfaceBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypeInterface, name, args...)
}

func (_this *interfaceBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *interfaceBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	return &interfaceBuilder{
		session: _this.session,
		parent:  parent,
		root:    root,
		opts:    opts,
	}
}

func (_this *interfaceBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *interfaceBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(common.TypeInterface))
}

func (_this *interfaceBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	value = common.CloneBytes(value)
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
		}
	case events.ArrayTypeCustomText:
		if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomText(_this, value, dst.Type(), err)
		}
	case events.ArrayTypeUint8:
		dst.Set(reflect.ValueOf(common.CloneBytes(value)))
	case events.ArrayTypeString:
		dst.Set(reflect.ValueOf(string(value)))
	case events.ArrayTypeResourceID:
		setPRIDFromString(string(value), dst)
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", arrayType))
	}
}

func (_this *interfaceBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *interfaceBuilder) BuildBeginList() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForListContents()
}

func (_this *interfaceBuilder) BuildBeginMap() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *interfaceBuilder) PrepareForListContents() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceSlice)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForListContents()
}

func (_this *interfaceBuilder) PrepareForMapContents() {
	builder := _this.session.GetBuilderForType(common.TypeInterfaceMap)
	builder = builder.NewInstance(_this.root, _this.parent, _this.opts)
	builder.PrepareForMapContents()
}

func (_this *interfaceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.parent.NotifyChildContainerFinished(value)
}
