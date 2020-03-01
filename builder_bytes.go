package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type bytesBuilder struct {
}

var globalBytesBuilder = &bytesBuilder{}

func newBytesBuilder() ObjectBuilder {
	return globalBytesBuilder
}

func (this *bytesBuilder) PostCacheInitBuilder() {
}

func (this *bytesBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *bytesBuilder) Nil(dst reflect.Value) {
	dst.Set(reflect.Zero(dst.Type()))
}

func (this *bytesBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "Bool")
}

func (this *bytesBuilder) Int(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "Int")
}

func (this *bytesBuilder) Uint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "Uint")
}

func (this *bytesBuilder) Float(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "Float")
}

func (this *bytesBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "String")
}

func (this *bytesBuilder) Bytes(value []byte, dst reflect.Value) {
	dst.SetBytes(value)
}

func (this *bytesBuilder) URI(value *url.URL, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "URI")
}

func (this *bytesBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, bytesType, "Time")
}

func (this *bytesBuilder) List() {
	builderPanicBadEvent(this, bytesType, "List")
}

func (this *bytesBuilder) Map() {
	builderPanicBadEvent(this, bytesType, "Map")
}

func (this *bytesBuilder) End() {
	builderPanicBadEvent(this, bytesType, "ContainerEnd")
}

func (this *bytesBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *bytesBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *bytesBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, bytesType, "PrepareForListContents")
}

func (this *bytesBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, bytesType, "PrepareForMapContents")
}

func (this *bytesBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, bytesType, "NotifyChildContainerFinished")
}
