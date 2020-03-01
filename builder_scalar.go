package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type scalarBuilder struct {
	// Const data
	dstType reflect.Type
}

func newBasicBuilder(dstType reflect.Type) ObjectBuilder {
	return &scalarBuilder{
		dstType: dstType,
	}
}

func (this *scalarBuilder) PostCacheInitBuilder() {
}

func (this *scalarBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *scalarBuilder) Nil(dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Nil")
}

func (this *scalarBuilder) Bool(value bool, dst reflect.Value) {
	dst.SetBool(value)
}

func (this *scalarBuilder) Int(value int64, dst reflect.Value) {
	dst.SetInt(value)
}

func (this *scalarBuilder) Uint(value uint64, dst reflect.Value) {
	dst.SetUint(value)
}

func (this *scalarBuilder) Float(value float64, dst reflect.Value) {
	dst.SetFloat(value)
}

func (this *scalarBuilder) String(value string, dst reflect.Value) {
	dst.SetString(value)
}

func (this *scalarBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "Bytes")
}

func (this *scalarBuilder) URI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "URI")
}

func (this *scalarBuilder) Time(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *scalarBuilder) List() {
	builderPanicBadEvent(this, this.dstType, "List")
}

func (this *scalarBuilder) Map() {
	builderPanicBadEvent(this, this.dstType, "Map")
}

func (this *scalarBuilder) End() {
	builderPanicBadEvent(this, this.dstType, "ContainerEnd")
}

func (this *scalarBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *scalarBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *scalarBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForListContents")
}

func (this *scalarBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, this.dstType, "PrepareForMapContents")
}

func (this *scalarBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, this.dstType, "NotifyChildContainerFinished")
}
