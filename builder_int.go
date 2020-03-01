package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type intBuilder struct {
	// Const data
	dstType reflect.Type
}

func newIntBuilder(dstType reflect.Type) ObjectBuilder {
	return &intBuilder{
		dstType: dstType,
	}
}

func (this *intBuilder) PostCacheInitBuilder() {
}

func (this *intBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *intBuilder) Nil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *intBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *intBuilder) Int(value int64, dst reflect.Value) {
	dst.SetInt(value)
	if dst.Int() != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *intBuilder) Uint(value uint64, dst reflect.Value) {
	if value > 0x7fffffffffffffff {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetInt(int64(value))
	if uint64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *intBuilder) Float(value float64, dst reflect.Value) {
	dst.SetInt(int64(value))
	if float64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *intBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *intBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *intBuilder) URI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *intBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *intBuilder) List() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *intBuilder) Map() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *intBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *intBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *intBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *intBuilder) PrepareForListContents() {
}

func (this *intBuilder) PrepareForMapContents() {
}

func (this *intBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
