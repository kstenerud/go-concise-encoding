package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

const defaultSliceCap = 4

type sliceBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
}

func newSliceBuilder(dstType reflect.Type) ObjectBuilder {
	return &sliceBuilder{
		dstType: dstType,
	}
}

func (this *sliceBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *sliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &sliceBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that)
	that.reset()
	return that
}

func (this *sliceBuilder) reset() {
	this.container = reflect.MakeSlice(this.dstType, 0, defaultSliceCap)
}

func (this *sliceBuilder) newElem() reflect.Value {
	return reflect.New(this.dstType.Elem()).Elem()
}

func (this *sliceBuilder) storeValue(value reflect.Value) {
	this.container = reflect.Append(this.container, value)
}

func (this *sliceBuilder) Nil(ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Nil(object)
	this.storeValue(object)
}

func (this *sliceBuilder) Bool(value bool, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Bool(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) Int(value int64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Int(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) Uint(value uint64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Uint(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) Float(value float64, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Float(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) String(value string, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.String(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) Bytes(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Bytes(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) URI(value *url.URL, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.URI(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) Time(value time.Time, ignored reflect.Value) {
	object := this.newElem()
	this.elemBuilder.Time(value, object)
	this.storeValue(object)
}

func (this *sliceBuilder) List() {
	this.elemBuilder.PrepareForListContents()
}

func (this *sliceBuilder) Map() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *sliceBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *sliceBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *sliceBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *sliceBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *sliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *sliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
