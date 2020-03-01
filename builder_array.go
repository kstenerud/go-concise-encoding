package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type arrayBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
	index     int
}

func newArrayBuilder(dstType reflect.Type) ObjectBuilder {
	return &arrayBuilder{
		dstType: dstType,
	}
}

func (this *arrayBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *arrayBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &arrayBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that)
	that.reset()
	return that
}

func (this *arrayBuilder) reset() {
	this.container = reflect.New(this.dstType).Elem()
	this.index = 0
}

func (this *arrayBuilder) currentElem() reflect.Value {
	return this.container.Index(this.index)
}

func (this *arrayBuilder) Nil(ignored reflect.Value) {
	this.elemBuilder.Nil(this.currentElem())
	this.index++
}

func (this *arrayBuilder) Bool(value bool, ignored reflect.Value) {
	this.elemBuilder.Bool(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) Int(value int64, ignored reflect.Value) {
	this.elemBuilder.Int(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) Uint(value uint64, ignored reflect.Value) {
	this.elemBuilder.Uint(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) Float(value float64, ignored reflect.Value) {
	this.elemBuilder.Float(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) String(value string, ignored reflect.Value) {
	this.elemBuilder.String(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) Bytes(value []byte, dst reflect.Value) {
	// TODO: Is this the right way to do this?
	for i := 0; i < len(value); i++ {
		elem := dst.Index(i + this.index)
		elem.SetUint(uint64(value[i]))
	}
}

func (this *arrayBuilder) URI(value *url.URL, ignored reflect.Value) {
	this.elemBuilder.URI(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) Time(value time.Time, ignored reflect.Value) {
	this.elemBuilder.Time(value, this.currentElem())
	this.index++
}

func (this *arrayBuilder) List() {
	this.elemBuilder.PrepareForListContents()
}

func (this *arrayBuilder) Map() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *arrayBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *arrayBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *arrayBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *arrayBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *arrayBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *arrayBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.currentElem().Set(value)
	this.index++
}
