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
)

type structBuilderDesc struct {
	builder ObjectBuilder
	index   int
}

type structBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	builderDescs  map[string]*structBuilderDesc
	nameBuilder   ObjectBuilder
	ignoreBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

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

func (this *structBuilder) IsContainerOnly() bool {
	return true
}

func (this *structBuilder) PostCacheInitBuilder() {
	this.nameBuilder = getBuilderForType(reflect.TypeOf(""))
	this.builderDescs = make(map[string]*structBuilderDesc)
	this.ignoreBuilder = newIgnoreBuilder()
	for i := 0; i < this.dstType.NumField(); i++ {
		field := this.dstType.Field(i)
		if field.PkgPath == "" {
			builder := getBuilderForType(field.Type)
			this.builderDescs[field.Name] = &structBuilderDesc{
				builder: builder,
				index:   i,
			}
		}
	}
}

func (this *structBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &structBuilder{
		dstType:      this.dstType,
		builderDescs: make(map[string]*structBuilderDesc),
		parent:       parent,
		root:         root,
	}
	that.nameBuilder = this.nameBuilder.CloneFromTemplate(root, that)
	that.ignoreBuilder = this.ignoreBuilder.CloneFromTemplate(root, that)
	for k, builderElem := range this.builderDescs {
		that.builderDescs[k] = &structBuilderDesc{
			builder: builderElem.builder.CloneFromTemplate(root, that),
			index:   builderElem.index,
		}
	}
	that.reset()
	return that
}

func (this *structBuilder) reset() {
	this.nextBuilder = this.nameBuilder
	this.container = reflect.New(this.dstType).Elem()
	this.nextValue = reflect.Value{}
	this.nextIsKey = true
	this.nextIsIgnored = false
}

func (this *structBuilder) swapKeyValue() {
	this.nextIsKey = !this.nextIsKey
}

func (this *structBuilder) BuildFromNil(ignored reflect.Value) {
	this.nextBuilder.BuildFromNil(this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromBool(value bool, ignored reflect.Value) {
	this.nextBuilder.BuildFromBool(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromInt(value int64, ignored reflect.Value) {
	this.nextBuilder.BuildFromInt(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromUint(value uint64, ignored reflect.Value) {
	this.nextBuilder.BuildFromUint(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromBigInt(value *big.Int, ignored reflect.Value) {
	this.nextBuilder.BuildFromBigInt(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromFloat(value float64, ignored reflect.Value) {
	this.nextBuilder.BuildFromFloat(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromBigFloat(value *big.Float, ignored reflect.Value) {
	this.nextBuilder.BuildFromBigFloat(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromDecimalFloat(value compact_float.DFloat, ignored reflect.Value) {
	this.nextBuilder.BuildFromDecimalFloat(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, ignored reflect.Value) {
	this.nextBuilder.BuildFromBigDecimalFloat(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromUUID(value []byte, ignored reflect.Value) {
	this.nextBuilder.BuildFromUUID(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromString(value string, ignored reflect.Value) {
	if this.nextIsKey {
		if builderDesc, ok := this.builderDescs[value]; ok {
			this.nextBuilder = builderDesc.builder
			this.nextValue = this.container.Field(builderDesc.index)
		} else {
			this.root.setCurrentBuilder(this.ignoreBuilder)
			this.nextBuilder = this.ignoreBuilder
			this.nextIsIgnored = true
			return
		}
	} else {
		this.nextBuilder.BuildFromString(value, this.nextValue)
	}

	this.swapKeyValue()
}

func (this *structBuilder) BuildFromBytes(value []byte, ignored reflect.Value) {
	this.nextBuilder.BuildFromBytes(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromURI(value *url.URL, ignored reflect.Value) {
	this.nextBuilder.BuildFromURI(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildFromTime(value time.Time, ignored reflect.Value) {
	this.nextBuilder.BuildFromTime(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) BuildBeginList() {
	this.nextBuilder.PrepareForListContents()
}

func (this *structBuilder) BuildBeginMap() {
	this.nextBuilder.PrepareForMapContents()
}

func (this *structBuilder) BuildEndContainer() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *structBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: structBuilder.Marker")
}

func (this *structBuilder) BuildFromReference(id interface{}) {
	panic("TODO: structBuilder.Reference")
}

func (this *structBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *structBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *structBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	if this.nextIsIgnored {
		this.nextIsIgnored = false
		return
	}

	this.nextValue.Set(value)
	this.swapKeyValue()
}
