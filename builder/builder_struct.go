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

type structBuilderDesc struct {
	builder ObjectBuilder
	index   int
}

type structBuilder struct {
	// Static data
	dstType reflect.Type

	// Cloned data (must be populated)
	builderDescs  map[string]*structBuilderDesc
	nameBuilder   ObjectBuilder
	ignoreBuilder ObjectBuilder

	// Clone inserted data
	root    *RootBuilder
	parent  ObjectBuilder
	options *options.BuilderOptions

	// Variable data (must be reset)
	nextBuilder   ObjectBuilder
	container     reflect.Value
	nextValue     reflect.Value
	nextIsKey     bool
	nextIsIgnored bool
}

func newStructBuilder(dstType reflect.Type) ObjectBuilder {
	return &structBuilder{
		dstType: dstType,
	}
}

func (_this *structBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *structBuilder) PostCacheInitBuilder(session *Session) {
	_this.nameBuilder = session.GetBuilderForType(reflect.TypeOf(""))
	_this.builderDescs = make(map[string]*structBuilderDesc)
	_this.ignoreBuilder = newIgnoreBuilder()
	for i := 0; i < _this.dstType.NumField(); i++ {
		field := _this.dstType.Field(i)
		if field.PkgPath == "" {
			builder := session.GetBuilderForType(field.Type)
			_this.builderDescs[field.Name] = &structBuilderDesc{
				builder: builder,
				index:   i,
			}
		}
	}
}

func (_this *structBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	that := &structBuilder{
		dstType:      _this.dstType,
		builderDescs: _this.builderDescs,
		parent:       parent,
		root:         root,
		options:      options,
	}
	that.nameBuilder = _this.nameBuilder.CloneFromTemplate(root, that, options)
	that.ignoreBuilder = _this.ignoreBuilder.CloneFromTemplate(root, that, options)
	that.reset()
	return that
}

func (_this *structBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *structBuilder) reset() {
	_this.nextBuilder = _this.nameBuilder
	_this.container = reflect.New(_this.dstType).Elem()
	_this.nextValue = reflect.Value{}
	_this.nextIsKey = true
	_this.nextIsIgnored = false
}

func (_this *structBuilder) swapKeyValue() {
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *structBuilder) BuildFromNil(_ reflect.Value) {
	_this.nextBuilder.BuildFromNil(_this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBool(value bool, _ reflect.Value) {
	_this.nextBuilder.BuildFromBool(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromInt(value int64, _ reflect.Value) {
	_this.nextBuilder.BuildFromInt(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	_this.nextBuilder.BuildFromUint(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigInt(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromFloat(value float64, _ reflect.Value) {
	_this.nextBuilder.BuildFromFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigFloat(value *big.Float, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromDecimalFloat(value compact_float.DFloat, _ reflect.Value) {
	_this.nextBuilder.BuildFromDecimalFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigDecimalFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromUUID(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromUUID(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromString(value []byte, _ reflect.Value) {
	if _this.nextIsKey {
		if builderDesc, ok := _this.builderDescs[string(value)]; ok {
			_this.nextBuilder = builderDesc.builder
			_this.nextValue = _this.container.Field(builderDesc.index)
		} else {
			_this.root.SetCurrentBuilder(_this.ignoreBuilder)
			_this.nextBuilder = _this.ignoreBuilder
			_this.nextIsIgnored = true
			return
		}
	} else {
		_this.nextBuilder.BuildFromString(value, _this.nextValue)
	}

	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromVerbatimString(value []byte, _ reflect.Value) {
	if _this.nextIsKey {
		if builderDesc, ok := _this.builderDescs[string(value)]; ok {
			_this.nextBuilder = builderDesc.builder
			_this.nextValue = _this.container.Field(builderDesc.index)
		} else {
			_this.root.SetCurrentBuilder(_this.ignoreBuilder)
			_this.nextBuilder = _this.ignoreBuilder
			_this.nextIsIgnored = true
			return
		}
	} else {
		_this.nextBuilder.BuildFromVerbatimString(value, _this.nextValue)
	}

	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBytes(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromBytes(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCustomBinary(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromCustomBinary(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCustomText(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromCustomText(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromURI(value *url.URL, _ reflect.Value) {
	_this.nextBuilder.BuildFromURI(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromTime(value time.Time, _ reflect.Value) {
	_this.nextBuilder.BuildFromTime(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCompactTime(value *compact_time.Time, _ reflect.Value) {
	_this.nextBuilder.BuildFromCompactTime(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildBeginList() {
	_this.nextBuilder.PrepareForListContents()
}

func (_this *structBuilder) BuildBeginMap() {
	_this.nextBuilder.PrepareForMapContents()
}

func (_this *structBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *structBuilder) BuildBeginMarker(id interface{}) {
	root := _this.root
	_this.nextBuilder = newMarkerObjectBuilder(_this, _this.nextBuilder, func(object reflect.Value) {
		root.GetMarkerRegistry().NotifyMarker(id, object)
	})
}

func (_this *structBuilder) BuildFromReference(id interface{}) {
	nextValue := _this.nextValue
	_this.swapKeyValue()
	_this.root.GetMarkerRegistry().NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, nextValue)
	})
}

func (_this *structBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, _this.dstType, "PrepareForListContents")
}

func (_this *structBuilder) PrepareForMapContents() {
	builderDescs := make(map[string]*structBuilderDesc)

	for k, builderElem := range _this.builderDescs {
		builderDescs[k] = &structBuilderDesc{
			builder: builderElem.builder.CloneFromTemplate(_this.root, _this, _this.options),
			index:   builderElem.index,
		}
	}
	_this.builderDescs = builderDescs
	_this.root.SetCurrentBuilder(_this)
}

func (_this *structBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	if _this.nextIsIgnored {
		_this.nextIsIgnored = false
		return
	}

	_this.nextValue.Set(value)
	_this.swapKeyValue()
}
