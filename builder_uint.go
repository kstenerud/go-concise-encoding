package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type uintBuilder struct {
	// Const data
	dstType reflect.Type
}

func newUintBuilder(dstType reflect.Type) ObjectBuilder {
	return &uintBuilder{
		dstType: dstType,
	}
}

func (this *uintBuilder) PostCacheInitBuilder() {
}

func (this *uintBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *uintBuilder) Nil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *uintBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *uintBuilder) Int(value int64, dst reflect.Value) {
	if value < 0 {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(uint64(value))
	if int64(dst.Uint()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *uintBuilder) Uint(value uint64, dst reflect.Value) {
	dst.SetUint(value)
	if dst.Uint() != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *uintBuilder) Float(value float64, dst reflect.Value) {
	if value < 0 {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(uint64(value))
	if float64(dst.Uint()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *uintBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *uintBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *uintBuilder) URI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *uintBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *uintBuilder) List() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *uintBuilder) Map() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *uintBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *uintBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *uintBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *uintBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *uintBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *uintBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
