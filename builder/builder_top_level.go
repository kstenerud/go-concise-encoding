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
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// topLevelContainerBuilder proxies the first build instruction to make sure containers
// are properly built. See BuildBeginList and BuildBeginMap.
type topLevelBuilder struct {
	builder ObjectBuilder
	root    *RootBuilder
}

func newTopLevelBuilder(root *RootBuilder, builder ObjectBuilder) ObjectBuilder {
	return &topLevelBuilder{
		builder: builder,
		root:    root,
	}
}

func (_this *topLevelBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.builder)
}

func (_this *topLevelBuilder) InitTemplate(_ *Session) {
	PanicBadEvent(_this, "InitTemplate")
}

func (_this *topLevelBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	PanicBadEvent(_this, "NewInstance")
	return nil
}

func (_this *topLevelBuilder) SetParent(_ ObjectBuilder) {
	PanicBadEvent(_this, "SetParent")
}

func (_this *topLevelBuilder) BuildFromNil(dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromNil(dst)
}

func (_this *topLevelBuilder) BuildFromBool(value bool, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromBool(value, dst)
}

func (_this *topLevelBuilder) BuildFromInt(value int64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromInt(value, dst)
}

func (_this *topLevelBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromUint(value, dst)
}

func (_this *topLevelBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromBigInt(value, dst)
}

func (_this *topLevelBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromFloat(value, dst)
}

func (_this *topLevelBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromBigFloat(value, dst)
}

func (_this *topLevelBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromDecimalFloat(value, dst)
}

func (_this *topLevelBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromBigDecimalFloat(value, dst)
}

func (_this *topLevelBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromUUID(value, dst)
}

func (_this *topLevelBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromArray(arrayType, value, dst)
}

func (_this *topLevelBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromTime(value, dst)
}

func (_this *topLevelBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	_this.root.SetCurrentBuilder(_this.builder)
	_this.builder.BuildFromCompactTime(value, dst)
}

func (_this *topLevelBuilder) BuildBeginList() {
	_this.builder.PrepareForListContents()
}

func (_this *topLevelBuilder) BuildBeginMap() {
	_this.builder.PrepareForMapContents()
}

func (_this *topLevelBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "End")
}

func (_this *topLevelBuilder) PrepareForListContents() {
	PanicBadEvent(_this, "PrepareForListContents")
}

func (_this *topLevelBuilder) PrepareForMapContents() {
	PanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *topLevelBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.NotifyChildContainerFinished(value)
}

func (_this *topLevelBuilder) BuildBeginMarker(id interface{}) {
	origBuilder := _this.builder
	_this.builder = newMarkerObjectBuilder(_this, origBuilder, func(object reflect.Value) {
		_this.builder = origBuilder
		_this.root.NotifyMarker(id, object)
	})
}

func (_this *topLevelBuilder) BuildFromReference(_ interface{}) {
	PanicBadEvent(_this, "BuildFromReference")
}
