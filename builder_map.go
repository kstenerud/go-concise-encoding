package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

const (
	kvBuilderKey   = 0
	kvBuilderValue = 1
)

type mapBuilder struct {
	// Const data
	dstType reflect.Type
	kvTypes [2]reflect.Type

	// Cloned data (must be populated)
	kvBuilders [2]ObjectBuilder

	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container    reflect.Value
	builderIndex int
	key          reflect.Value
}

func newMapBuilder(dstType reflect.Type) ObjectBuilder {
	return &mapBuilder{
		dstType: dstType,
		kvTypes: [2]reflect.Type{dstType.Key(), dstType.Elem()},
	}
}

func (this *mapBuilder) PostCacheInitBuilder() {
	this.kvBuilders[kvBuilderKey] = getBuilderForType(this.dstType.Key())
	this.kvBuilders[kvBuilderValue] = getBuilderForType(this.dstType.Elem())
}

func (this *mapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &mapBuilder{
		dstType: this.dstType,
		kvTypes: this.kvTypes,
		parent:  parent,
		root:    root,
	}
	that.kvBuilders[kvBuilderKey] = this.kvBuilders[kvBuilderKey].CloneFromTemplate(root, that)
	that.kvBuilders[kvBuilderValue] = this.kvBuilders[kvBuilderValue].CloneFromTemplate(root, that)
	that.reset()
	return that
}

func (this *mapBuilder) reset() {
	this.container = reflect.MakeMap(this.dstType)
	this.builderIndex = kvBuilderKey
	this.key = reflect.Value{}
}

func (this *mapBuilder) getBuilder() ObjectBuilder {
	return this.kvBuilders[this.builderIndex]
}

func (this *mapBuilder) storeValue(value reflect.Value) {
	if this.builderIndex == kvBuilderKey {
		this.key = value
	} else {
		this.container.SetMapIndex(this.key, value)
	}
	this.builderIndex = (this.builderIndex + 1) & 1
}

func (this *mapBuilder) newElem() reflect.Value {
	return reflect.New(this.kvTypes[this.builderIndex]).Elem()
}

func (this *mapBuilder) Nil(ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Nil(object)
	this.storeValue(object)
}

func (this *mapBuilder) Bool(value bool, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Bool(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) Int(value int64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Int(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) Uint(value uint64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Uint(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) Float(value float64, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Float(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) String(value string, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().String(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) Bytes(value []byte, ignored reflect.Value) {
	object := this.newElem()
	this.getBuilder().Bytes(value, object)
	this.storeValue(object)
}

func (this *mapBuilder) URI(value *url.URL, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *mapBuilder) Time(value time.Time, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *mapBuilder) List() {
	this.getBuilder().PrepareForListContents()
}

func (this *mapBuilder) Map() {
	this.getBuilder().PrepareForMapContents()
}

func (this *mapBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *mapBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *mapBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *mapBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForListContents")
}

func (this *mapBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *mapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
