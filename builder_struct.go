package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
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

func (this *structBuilder) Nil(ignored reflect.Value) {
	this.nextBuilder.Nil(this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) Bool(value bool, ignored reflect.Value) {
	this.nextBuilder.Bool(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) Int(value int64, ignored reflect.Value) {
	this.nextBuilder.Int(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) Uint(value uint64, ignored reflect.Value) {
	this.nextBuilder.Uint(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) Float(value float64, ignored reflect.Value) {
	this.nextBuilder.Float(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) String(value string, ignored reflect.Value) {
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
		this.nextBuilder.String(value, this.nextValue)
	}

	this.swapKeyValue()
}

func (this *structBuilder) Bytes(value []byte, ignored reflect.Value) {
	this.nextBuilder.Bytes(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) URI(value *url.URL, ignored reflect.Value) {
	this.nextBuilder.URI(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) Time(value time.Time, ignored reflect.Value) {
	this.nextBuilder.Time(value, this.nextValue)
	this.swapKeyValue()
}

func (this *structBuilder) List() {
	this.nextBuilder.PrepareForListContents()
}

func (this *structBuilder) Map() {
	this.nextBuilder.PrepareForMapContents()
}

func (this *structBuilder) End() {
	object := this.container
	this.reset()
	this.parent.NotifyChildContainerFinished(object)
}

func (this *structBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *structBuilder) Reference(id interface{}) {
	panic("TODO")
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
