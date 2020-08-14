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

type markerIDBuilder struct {
	// Template Data
	onID func(interface{})
}

func newMarkerIDBuilder(onID func(interface{})) *markerIDBuilder {
	return &markerIDBuilder{
		onID: onID,
	}
}

func (_this *markerIDBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *markerIDBuilder) InitTemplate(_ *Session) {
}

func (_this *markerIDBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *markerIDBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *markerIDBuilder) BuildFromNil(_ reflect.Value) {
	PanicBadEvent(_this, "Nil")
}

func (_this *markerIDBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	PanicBadEvent(_this, "Bool")
}

func (_this *markerIDBuilder) BuildFromInt(value int64, _ reflect.Value) {
	if value < 0 {
		PanicBadEvent(_this, "Int")
	}
	_this.onID(value)
}

func (_this *markerIDBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	_this.onID(value)
}

func (_this *markerIDBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	if common.IsBigIntNegative(value) || !value.IsUint64() {
		PanicBadEvent(_this, "BigInt")
	}
	_this.onID(value.Uint64())
}

func (_this *markerIDBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	PanicBadEvent(_this, "Float")
}

func (_this *markerIDBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	PanicBadEvent(_this, "BigFloat")
}

func (_this *markerIDBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	PanicBadEvent(_this, "DecimalFloat")
}

func (_this *markerIDBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	PanicBadEvent(_this, "BigDecimalFloat")
}

func (_this *markerIDBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "UUID")
}

func (_this *markerIDBuilder) BuildFromString(value []byte, _ reflect.Value) {
	_this.onID(value)
}

func (_this *markerIDBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "VerbatimString")
}

func (_this *markerIDBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "Bytes")
}

func (_this *markerIDBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "CustomBinary")
}

func (_this *markerIDBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	PanicBadEvent(_this, "CustomText")
}

func (_this *markerIDBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	PanicBadEvent(_this, "URI")
}

func (_this *markerIDBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "Time")
}

func (_this *markerIDBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	PanicBadEvent(_this, "CompactTime")
}

func (_this *markerIDBuilder) BuildBeginList() {
	PanicBadEvent(_this, "List")
}

func (_this *markerIDBuilder) BuildBeginMap() {
	PanicBadEvent(_this, "Map")
}

func (_this *markerIDBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "End")
}

func (_this *markerIDBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEvent(_this, "Marker")
}

func (_this *markerIDBuilder) BuildFromReference(_ interface{}) {
	PanicBadEvent(_this, "Reference")
}

func (_this *markerIDBuilder) PrepareForListContents() {
	PanicBadEvent(_this, "PrepareForListContents")
}

func (_this *markerIDBuilder) PrepareForMapContents() {
	PanicBadEvent(_this, "PrepareForMapContents")
}

func (_this *markerIDBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	PanicBadEvent(_this, "NotifyChildContainerFinished")
}

// ============================================================================

type markerObjectBuilder struct {
	// Instance Data
	parent           ObjectBuilder
	child            ObjectBuilder
	onObjectComplete func(reflect.Value)
}

func newMarkerObjectBuilder(parent, child ObjectBuilder, onObjectComplete func(reflect.Value)) *markerObjectBuilder {
	return &markerObjectBuilder{
		parent:           parent,
		child:            child,
		onObjectComplete: onObjectComplete,
	}
}

func (_this *markerObjectBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.child)
}

func (_this *markerObjectBuilder) InitTemplate(_ *Session) {
}

func (_this *markerObjectBuilder) NewInstance(_ *RootBuilder, parent ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return &markerObjectBuilder{
		parent:           parent,
		child:            _this.child,
		onObjectComplete: _this.onObjectComplete,
	}
}

func (_this *markerObjectBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *markerObjectBuilder) BuildFromNil(dst reflect.Value) {
	_this.child.BuildFromNil(dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromBool(value bool, dst reflect.Value) {
	_this.child.BuildFromBool(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromInt(value int64, dst reflect.Value) {
	_this.child.BuildFromInt(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	_this.child.BuildFromUint(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	_this.child.BuildFromBigInt(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	_this.child.BuildFromFloat(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	_this.child.BuildFromBigFloat(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	_this.child.BuildFromDecimalFloat(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	_this.child.BuildFromBigDecimalFloat(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	_this.child.BuildFromUUID(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromString(value []byte, dst reflect.Value) {
	_this.child.BuildFromString(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	_this.child.BuildFromVerbatimString(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	_this.child.BuildFromBytes(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	_this.child.BuildFromCustomBinary(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	_this.child.BuildFromCustomText(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	_this.child.BuildFromURI(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	_this.child.BuildFromTime(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	_this.child.BuildFromCompactTime(value, dst)
	_this.onObjectComplete(dst)
}

func (_this *markerObjectBuilder) BuildBeginList() {
	PanicBadEvent(_this, "List")
}

func (_this *markerObjectBuilder) BuildBeginMap() {
	PanicBadEvent(_this, "Map")
}

func (_this *markerObjectBuilder) BuildEndContainer() {
	PanicBadEvent(_this, "End")
}

func (_this *markerObjectBuilder) BuildBeginMarker(_ interface{}) {
	PanicBadEvent(_this, "Marker")
}

func (_this *markerObjectBuilder) BuildFromReference(_ interface{}) {
	PanicBadEvent(_this, "Reference")
}

func (_this *markerObjectBuilder) PrepareForListContents() {
	_this.child.SetParent(_this)
	_this.child.PrepareForListContents()
}

func (_this *markerObjectBuilder) PrepareForMapContents() {
	_this.child.SetParent(_this)
	_this.child.PrepareForMapContents()
}

func (_this *markerObjectBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.onObjectComplete(value)
	_this.child.SetParent(_this.parent)
	_this.parent.NotifyChildContainerFinished(value)
}
