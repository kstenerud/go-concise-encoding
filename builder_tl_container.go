package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type tlContainerBuilder struct {
	// Const data
	dstType reflect.Type

	// Cloned data (must be populated)
	builder ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newTLContainerBuilder(dstType reflect.Type) ObjectBuilder {
	return &tlContainerBuilder{
		dstType: dstType,
		builder: getBuilderForType(dstType),
	}
}

func (this *tlContainerBuilder) PostCacheInitBuilder() {
	// TLContainer is not cached, so we must init on creation
}

func (this *tlContainerBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &tlContainerBuilder{
		dstType: this.dstType,
		builder: this.builder.CloneFromTemplate(root, parent),
		parent:  parent,
		root:    root,
	}
	return that
}

func (this *tlContainerBuilder) Nil(ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *tlContainerBuilder) Bool(value bool, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *tlContainerBuilder) Int(value int64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Int")
}

func (this *tlContainerBuilder) Uint(value uint64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Uint")
}

func (this *tlContainerBuilder) Float(value float64, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Float")
}

func (this *tlContainerBuilder) String(value string, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *tlContainerBuilder) Bytes(value []byte, dst reflect.Value) {
	// TODO: Is this the right way to do this?
	this.root.setCurrentBuilder(this.builder)
	this.builder.Bytes(value, dst)
}

func (this *tlContainerBuilder) URI(value *url.URL, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *tlContainerBuilder) Time(value time.Time, ignored reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *tlContainerBuilder) List() {
	this.root.setCurrentBuilder(this.builder)
	this.builder.PrepareForListContents()
}

func (this *tlContainerBuilder) Map() {
	this.root.setCurrentBuilder(this.builder)
	this.builder.PrepareForMapContents()
}

func (this *tlContainerBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "End")
}

func (this *tlContainerBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *tlContainerBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *tlContainerBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}

func (this *tlContainerBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *tlContainerBuilder) Reference(id interface{}) {
	panic("TODO")
}
