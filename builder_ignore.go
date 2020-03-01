package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type ignoreBuilder struct {
	// Clone inserted data
	root   *RootBuilder
	parent ObjectBuilder
}

var globalIgnoreBuilder = &ignoreBuilder{}

func newIgnoreBuilder() ObjectBuilder {
	return globalIgnoreBuilder
}

func (this *ignoreBuilder) PostCacheInitBuilder() {
}

func (this *ignoreBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return &ignoreBuilder{
		parent: parent,
		root:   root,
	}
}

func (this *ignoreBuilder) Nil(dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Bool(value bool, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Int(value int64, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Uint(value uint64, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Float(value float64, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) String(value string, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Bytes(value []byte, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) URI(value *url.URL, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) Time(value time.Time, dst reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}

func (this *ignoreBuilder) List() {
	// TODO: Pre-generate these?
	container := reflect.ValueOf(make([]interface{}, 0))
	builder := getBuilderForType(container.Type())
	builder.CloneFromTemplate(this.root, this)
	builder.PrepareForListContents()
}

func (this *ignoreBuilder) Map() {
	builder := getBuilderForType(reflect.TypeOf(map[interface{}]interface{}{}))
	builder.CloneFromTemplate(this.root, this)
	builder.PrepareForMapContents()
}

func (this *ignoreBuilder) End() {
	builderPanicBadEvent(this, reflect.TypeOf([]interface{}{}).Elem(), "End")
}

func (this *ignoreBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *ignoreBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *ignoreBuilder) PrepareForListContents() {
	this.root.setCurrentBuilder(this)
}

func (this *ignoreBuilder) PrepareForMapContents() {
	this.root.setCurrentBuilder(this)
}

func (this *ignoreBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.root.setCurrentBuilder(this.parent)
}
