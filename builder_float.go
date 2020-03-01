package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type floatBuilder struct {
	// Const data
	dstType reflect.Type
}

func newFloatBuilder(dstType reflect.Type) ObjectBuilder {
	return &floatBuilder{
		dstType: dstType,
	}
}

func (this *floatBuilder) PostCacheInitBuilder() {
}

func (this *floatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *floatBuilder) Nil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *floatBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bool")
}

func (this *floatBuilder) Int(value int64, dst reflect.Value) {
	dst.SetFloat(float64(value))
	if int64(dst.Float()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *floatBuilder) Uint(value uint64, dst reflect.Value) {
	dst.SetFloat(float64(value))
	if uint64(dst.Float()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func (this *floatBuilder) Float(value float64, dst reflect.Value) {
	// Note: We just silently truncate 64-to-32 bit conversions.
	dst.SetFloat(value)
}

func (this *floatBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "String")
}

func (this *floatBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *floatBuilder) URI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *floatBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Time")
}

func (this *floatBuilder) List() {
	builderPanicBadEvent(this, this.dstType, "ListBegin")
}

func (this *floatBuilder) Map() {
	builderPanicBadEvent(this, this.dstType, "MapBegin")
}

func (this *floatBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *floatBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *floatBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *floatBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *floatBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *floatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
