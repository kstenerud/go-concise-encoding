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

type arrayBuilder struct {
	// Static data
	dstType     reflect.Type
	elemBuilder ObjectBuilder

	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *options.BuilderOptions

	// Variable data (must be reset)
	container reflect.Value
	index     int
}

func newArrayBuilder(dstType reflect.Type) ObjectBuilder {
	return &arrayBuilder{
		dstType: dstType,
	}
}

func (_this *arrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.elemBuilder)
}

func (_this *arrayBuilder) InitTemplate(session *Session) {
	_this.elemBuilder = session.GetBuilderForType(_this.dstType.Elem())
}

func (_this *arrayBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	that := &arrayBuilder{
		dstType:     _this.dstType,
		elemBuilder: _this.elemBuilder,
		root:        root,
		parent:      parent,
		options:     options,
	}
	that.reset()
	return that
}

func (_this *arrayBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *arrayBuilder) reset() {
	_this.container = reflect.New(_this.dstType).Elem()
	_this.index = 0
}

func (_this *arrayBuilder) currentElem() reflect.Value {
	return _this.container.Index(_this.index)
}

func (_this *arrayBuilder) BuildFromNil(_ reflect.Value) {
	_this.elemBuilder.BuildFromNil(_this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromBool(value bool, _ reflect.Value) {
	_this.elemBuilder.BuildFromBool(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromInt(value int64, _ reflect.Value) {
	_this.elemBuilder.BuildFromInt(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	_this.elemBuilder.BuildFromUint(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	_this.elemBuilder.BuildFromBigInt(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromFloat(value float64, _ reflect.Value) {
	_this.elemBuilder.BuildFromFloat(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromBigFloat(value *big.Float, _ reflect.Value) {
	_this.elemBuilder.BuildFromBigFloat(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromDecimalFloat(value compact_float.DFloat, _ reflect.Value) {
	_this.elemBuilder.BuildFromDecimalFloat(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, _ reflect.Value) {
	_this.elemBuilder.BuildFromBigDecimalFloat(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromUUID(value []byte, _ reflect.Value) {
	_this.elemBuilder.BuildFromUUID(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromString(value []byte, _ reflect.Value) {
	_this.elemBuilder.BuildFromString(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromVerbatimString(value []byte, _ reflect.Value) {
	_this.elemBuilder.BuildFromVerbatimString(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	_this.elemBuilder.BuildFromBytes(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	_this.elemBuilder.BuildFromCustomBinary(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	_this.elemBuilder.BuildFromCustomText(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromURI(value *url.URL, _ reflect.Value) {
	_this.elemBuilder.BuildFromURI(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromTime(value time.Time, _ reflect.Value) {
	_this.elemBuilder.BuildFromTime(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildFromCompactTime(value *compact_time.Time, _ reflect.Value) {
	_this.elemBuilder.BuildFromCompactTime(value, _this.currentElem())
	_this.index++
}

func (_this *arrayBuilder) BuildBeginList() {
	_this.elemBuilder.PrepareForListContents()
}

func (_this *arrayBuilder) BuildBeginMap() {
	_this.elemBuilder.PrepareForMapContents()
}

func (_this *arrayBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *arrayBuilder) BuildBeginMarker(id interface{}) {
	origBuilder := _this.elemBuilder
	_this.elemBuilder = newMarkerObjectBuilder(_this, origBuilder, func(object reflect.Value) {
		_this.elemBuilder = origBuilder
		_this.root.GetMarkerRegistry().NotifyMarker(id, object)
	})
}

func (_this *arrayBuilder) BuildFromReference(id interface{}) {
	container := _this.container
	index := _this.index
	_this.index++
	_this.root.GetMarkerRegistry().NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, container.Index(index))
	})
}

func (_this *arrayBuilder) PrepareForListContents() {
	_this.elemBuilder = _this.elemBuilder.NewInstance(_this.root, _this, _this.options)
	_this.root.SetCurrentBuilder(_this)
}

func (_this *arrayBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForMapContents")
}

func (_this *arrayBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.currentElem().Set(value)
	_this.index++
}

// ============================================================================

type bytesArrayBuilder struct {
	// Static data
	session *Session
}

func newBytesArrayBuilder() ObjectBuilder {
	return &bytesArrayBuilder{}
}

func (_this *bytesArrayBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *bytesArrayBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *bytesArrayBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bytesArrayBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bytesArrayBuilder) BuildFromNil(_ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromNil")
}

func (_this *bytesArrayBuilder) BuildFromBool(value bool, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromBool")
}

func (_this *bytesArrayBuilder) BuildFromInt(value int64, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromInt")
}

func (_this *bytesArrayBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromUint")
}

func (_this *bytesArrayBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromBigInt")
}

func (_this *bytesArrayBuilder) BuildFromFloat(value float64, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromFloat")
}

func (_this *bytesArrayBuilder) BuildFromBigFloat(value *big.Float, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromBigFloat")
}

func (_this *bytesArrayBuilder) BuildFromDecimalFloat(value compact_float.DFloat, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromDecimalFloat")
}

func (_this *bytesArrayBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromBigDecimalFloat")
}

func (_this *bytesArrayBuilder) BuildFromUUID(value []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromUUID")
}

func (_this *bytesArrayBuilder) BuildFromString(value []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromString")
}

func (_this *bytesArrayBuilder) BuildFromVerbatimString(value []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromVerbatimString")
}

func (_this *bytesArrayBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	// TODO: Is there a more efficient way?
	for i := 0; i < len(value); i++ {
		elem := dst.Index(i)
		elem.SetUint(uint64(value[i]))
	}
}

func (_this *bytesArrayBuilder) BuildFromCustomBinary(value []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromCustomBinary")
}

func (_this *bytesArrayBuilder) BuildFromCustomText(value []byte, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromCustomText")
}

func (_this *bytesArrayBuilder) BuildFromURI(value *url.URL, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromURI")
}

func (_this *bytesArrayBuilder) BuildFromTime(value time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromTime")
}

func (_this *bytesArrayBuilder) BuildFromCompactTime(value *compact_time.Time, _ reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromCompactTime")
}

func (_this *bytesArrayBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildBeginList")
}

func (_this *bytesArrayBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildBeginMap")
}

func (_this *bytesArrayBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildEndContainer")
}

func (_this *bytesArrayBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildBeginMarker")
}

func (_this *bytesArrayBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "BuildFromReference")
}

func (_this *bytesArrayBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "PrepareForListContents")
}

func (_this *bytesArrayBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "PrepareForMapContents")
}

func (_this *bytesArrayBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBytes, "NotifyChildContainerFinished")
}
