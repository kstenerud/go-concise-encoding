package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

var (
	builderIntfIntfMapType = reflect.TypeOf(map[interface{}]interface{}{})
	builderIntfSliceType   = reflect.TypeOf([]interface{}{})
	builderIntfType        = builderIntfSliceType.Elem()

	globalIntfBuilder        = &intfBuilder{}
	globalIntfSliceBuilder   = &intfSliceBuilder{}
	globalIntfIntfMapBuilder = &intfIntfMapBuilder{}
)

type intfBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

func newInterfaceBuilder() ObjectBuilder {
	return globalIntfBuilder
}

func (this *intfBuilder) PostCacheInitBuilder() {
}

func (this *intfBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return &intfBuilder{
		parent: parent,
		root:   root,
	}
}

func (this *intfBuilder) Nil(dst reflect.Value) {
	dst.Set(reflect.Zero(builderIntfType))
}

func (this *intfBuilder) Bool(value bool, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) Int(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) Uint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) Float(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) String(value string, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) Bytes(value []byte, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) URI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) Time(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *intfBuilder) List() {
	builderPanicBadEvent(this, builderIntfType, "List")
}

func (this *intfBuilder) Map() {
	builderPanicBadEvent(this, builderIntfType, "Map")
}

func (this *intfBuilder) End() {
	builderPanicBadEvent(this, builderIntfType, "ContainerEnd")
}

func (this *intfBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *intfBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *intfBuilder) PrepareForListContents() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this.parent)
	builder.PrepareForListContents()
}

func (this *intfBuilder) PrepareForMapContents() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this.parent)
	builder.PrepareForMapContents()
}

func (this *intfBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.parent.NotifyChildContainerFinished(value)
}

// -----
// Slice
// -----

type intfSliceBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
}

func newIntfSliceBuilder() ObjectBuilder {
	return globalIntfSliceBuilder
}

func (this *intfSliceBuilder) PostCacheInitBuilder() {
}

func (this *intfSliceBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &intfSliceBuilder{
		parent: parent,
		root:   root,
	}
	that.reset()
	return that
}

func (this *intfSliceBuilder) reset() {
	this.container = reflect.MakeSlice(builderIntfSliceType, 0, defaultSliceCap)
}

func (this *intfSliceBuilder) storeRValue(value reflect.Value) {
	this.container = reflect.Append(this.container, value)
}

func (this *intfSliceBuilder) storeValue(value interface{}) {
	this.storeRValue(reflect.ValueOf(value))
}

func (this *intfSliceBuilder) Nil(ignored reflect.Value) {
	this.storeRValue(reflect.New(builderIntfSliceType).Elem())
}

func (this *intfSliceBuilder) Bool(value bool, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) Int(value int64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) Uint(value uint64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) Float(value float64, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) String(value string, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) Bytes(value []byte, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) URI(value *url.URL, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) Time(value time.Time, ignored reflect.Value) {
	this.storeValue(value)
}

func (this *intfSliceBuilder) List() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForListContents()
}

func (this *intfSliceBuilder) Map() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForMapContents()
}

func (this *intfSliceBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *intfSliceBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *intfSliceBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *intfSliceBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *intfSliceBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForMapContents")
}

func (this *intfSliceBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeRValue(value)
}

// ---
// Map
// ---

type intfIntfMapBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder

	// Variable data (must be reset)
	container reflect.Value
	key       reflect.Value
	nextIsKey bool
}

func newIntfIntfMapBuilder() ObjectBuilder {
	return globalIntfIntfMapBuilder
}

func (this *intfIntfMapBuilder) PostCacheInitBuilder() {
}

func (this *intfIntfMapBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	that := &intfIntfMapBuilder{
		parent: parent,
		root:   root,
	}
	that.reset()
	return that
}

func (this *intfIntfMapBuilder) reset() {
	this.container = reflect.MakeMap(builderIntfIntfMapType)
	this.key = reflect.Value{}
	this.nextIsKey = true
}

func (this *intfIntfMapBuilder) storeValue(value reflect.Value) {
	if this.nextIsKey {
		this.key = value
	} else {
		this.container.SetMapIndex(this.key, value)
	}
	this.nextIsKey = !this.nextIsKey
}

func (this *intfIntfMapBuilder) Nil(ignored reflect.Value) {
	this.storeValue(reflect.Zero(builderIntfType))
}

func (this *intfIntfMapBuilder) Bool(value bool, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) Int(value int64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) Uint(value uint64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) Float(value float64, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) String(value string, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) Bytes(value []byte, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) URI(value *url.URL, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) Time(value time.Time, ignored reflect.Value) {
	this.storeValue(reflect.ValueOf(value))
}

func (this *intfIntfMapBuilder) List() {
	builder := globalIntfSliceBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForListContents()
}

func (this *intfIntfMapBuilder) Map() {
	builder := globalIntfIntfMapBuilder.CloneFromTemplate(this.root, this)
	builder.PrepareForMapContents()
}

func (this *intfIntfMapBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *intfIntfMapBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *intfIntfMapBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *intfIntfMapBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, builderIntfType, "PrepareForListContents")
}

func (this *intfIntfMapBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *intfIntfMapBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this)
	this.storeValue(value)
}
