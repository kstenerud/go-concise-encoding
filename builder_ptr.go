package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type ptrBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	elemBuilder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &ptrBuilder{
		dstType: dstType,
	}
}

func (this *ptrBuilder) PostCacheInitBuilder() {
	this.elemBuilder = getBuilderForType(this.dstType.Elem())
}

func (this *ptrBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &ptrBuilder{
		dstType: this.dstType,
		parent:  parent,
		root:    root,
	}
	that.elemBuilder = this.elemBuilder.CloneFromTemplate(root, that)
	return that
}

func (this *ptrBuilder) newElem() reflect.Value {
	return reflect.New(this.dstType.Elem())
}

func (this *ptrBuilder) Nil(dst reflect.Value) {
	dst.Set(reflect.Zero(this.dstType))
}

func (this *ptrBuilder) Bool(value bool, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Bool(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) Int(value int64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Int(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) Uint(value uint64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Uint(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) Float(value float64, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Float(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) String(value string, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.String(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) Bytes(value []byte, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Bytes(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) URI(value *url.URL, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.URI(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) Time(value time.Time, dst reflect.Value) {
	ptr := this.newElem()
	this.elemBuilder.Time(value, ptr.Elem())
	dst.Set(ptr)
}

func (this *ptrBuilder) List() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *ptrBuilder) Map() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *ptrBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *ptrBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *ptrBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *ptrBuilder) PrepareForListContents() {
	this.elemBuilder.PrepareForListContents()
}

func (this *ptrBuilder) PrepareForMapContents() {
	this.elemBuilder.PrepareForMapContents()
}

func (this *ptrBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.parent.NotifyChildContainerFinished(value.Addr())
}
