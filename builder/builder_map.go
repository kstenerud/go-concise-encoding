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

const (
	kvBuilderKey   = 0
	kvBuilderValue = 1
)

type mapBuilder struct {
	// Template Data
	session    *Session
	dstType    reflect.Type
	kvTypes    [2]reflect.Type
	kvBuilders [2]ObjectBuilder

	// Instance Data
	root   *RootBuilder
	parent ObjectBuilder
	opts   *options.BuilderOptions

	// Variable data (must be reset)
	container       reflect.Value
	key             reflect.Value
	builderIndex    int
	nextBuilder     ObjectBuilder
	nextStoreMethod func(*mapBuilder, reflect.Value)
}

func newMapBuilder(dstType reflect.Type) ObjectBuilder {
	return &mapBuilder{
		dstType: dstType,
		kvTypes: [2]reflect.Type{dstType.Key(), dstType.Elem()},
	}
}

func (_this *mapBuilder) String() string {
	return fmt.Sprintf("%v<%v:%v>", reflect.TypeOf(_this), _this.kvBuilders[0], _this.kvBuilders[1])
}

func (_this *mapBuilder) InitTemplate(session *Session) {
	_this.session = session
	_this.kvBuilders[kvBuilderKey] = session.GetBuilderForType(_this.dstType.Key())
	_this.kvBuilders[kvBuilderValue] = session.GetBuilderForType(_this.dstType.Elem())
}

func (_this *mapBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	that := &mapBuilder{
		session: _this.session,
		dstType: _this.dstType,
		kvTypes: _this.kvTypes,
		parent:  parent,
		root:    root,
		opts:    opts,
	}
	that.kvBuilders[kvBuilderKey] = _this.kvBuilders[kvBuilderKey].NewInstance(root, that, opts)
	that.kvBuilders[kvBuilderValue] = _this.kvBuilders[kvBuilderValue]
	that.reset()
	return that
}

func (_this *mapBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *mapBuilder) reset() {
	_this.container = reflect.MakeMap(_this.dstType)
	_this.key = reflect.Value{}
	_this.builderIndex = kvBuilderKey
	_this.nextBuilder = _this.kvBuilders[_this.builderIndex]
	_this.nextStoreMethod = mapBuilderKVStoreMethods[_this.builderIndex]
}

func (_this *mapBuilder) storeKey(value reflect.Value) {
	_this.key = value
}

func (_this *mapBuilder) storeValue(value reflect.Value) {
	_this.container.SetMapIndex(_this.key, value)
}

var mapBuilderKVStoreMethods = []func(*mapBuilder, reflect.Value){
	(*mapBuilder).storeKey,
	(*mapBuilder).storeValue,
}

func (_this *mapBuilder) store(value reflect.Value) {
	_this.nextStoreMethod(_this, value)
	_this.swapKeyValue()
}

func (_this *mapBuilder) swapKeyValue() {
	_this.builderIndex = (_this.builderIndex + 1) & 1
	_this.nextBuilder = _this.kvBuilders[_this.builderIndex]
	_this.nextStoreMethod = mapBuilderKVStoreMethods[_this.builderIndex]
}

func (_this *mapBuilder) newElem() reflect.Value {
	return reflect.New(_this.kvTypes[_this.builderIndex]).Elem()
}

func (_this *mapBuilder) BuildFromNil(_ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromNil(object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromBool(value bool, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBool(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromInt(value int64, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromInt(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromUint(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigInt(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromFloat(value float64, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromFloat(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromBigFloat(value *big.Float, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigFloat(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromDecimalFloat(value compact_float.DFloat, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromDecimalFloat(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBigDecimalFloat(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromUUID(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromUUID(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromString(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromString(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromVerbatimString(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromVerbatimString(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromBytes(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromBytes(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromCustomBinary(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromCustomBinary(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromCustomText(value []byte, _ reflect.Value) {
	object := _this.newElem()
	_this.nextBuilder.BuildFromCustomText(value, object)
	_this.store(object)
}

func (_this *mapBuilder) BuildFromURI(value *url.URL, _ reflect.Value) {
	_this.store(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildFromTime(value time.Time, _ reflect.Value) {
	_this.store(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildFromCompactTime(value *compact_time.Time, _ reflect.Value) {
	_this.store(reflect.ValueOf(value))
}

func (_this *mapBuilder) BuildBeginList() {
	_this.nextBuilder.PrepareForListContents()
}

func (_this *mapBuilder) BuildBeginMap() {
	_this.nextBuilder.PrepareForMapContents()
}

func (_this *mapBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *mapBuilder) BuildBeginMarker(id interface{}) {
	root := _this.root
	_this.nextBuilder = newMarkerObjectBuilder(_this, _this.nextBuilder, func(object reflect.Value) {
		root.NotifyMarker(id, object)
	})
}

func (_this *mapBuilder) BuildFromReference(id interface{}) {
	container := _this.container
	key := _this.key
	tempValue := _this.newElem()
	_this.swapKeyValue()
	_this.root.NotifyReference(id, func(object reflect.Value) {
		if container.Type().Elem().Kind() == reflect.Interface || object.Type() == container.Type().Elem() {
			// In case of self-referencing pointers, we need to pass the original container, not a copy.
			container.SetMapIndex(key, object)
		} else {
			setAnythingFromAnything(object, tempValue)
			container.SetMapIndex(key, tempValue)
		}
	})
}

func (_this *mapBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, builderIntfType, "PrepareForListContents")
}

func (_this *mapBuilder) PrepareForMapContents() {
	_this.kvBuilders[kvBuilderValue] = _this.kvBuilders[kvBuilderValue].NewInstance(_this.root, _this, _this.opts)
	_this.root.SetCurrentBuilder(_this)
}

func (_this *mapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	_this.store(value)
}
