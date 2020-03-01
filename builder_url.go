package concise_encoding

import (
	"net/url"
	"reflect"
	"time"
)

type urlBuilder struct {
}

func newURLBuilder() ObjectBuilder {
	return &urlBuilder{}
}

func (this *urlBuilder) PostCacheInitBuilder() {
}

func (this *urlBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *urlBuilder) Nil(dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Nil")
}

func (this *urlBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Bool")
}

func (this *urlBuilder) Int(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Int")
}

func (this *urlBuilder) Uint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Uint")
}

func (this *urlBuilder) Float(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Float")
}

func (this *urlBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "String")
}

func (this *urlBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Bytes")
}

func (this *urlBuilder) URI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value).Elem())
}

func (this *urlBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, urlType, "Time")
}

func (this *urlBuilder) List() {
	builderPanicBadEvent(this, urlType, "List")
}

func (this *urlBuilder) Map() {
	builderPanicBadEvent(this, urlType, "Map")
}

func (this *urlBuilder) End() {
	builderPanicBadEvent(this, urlType, "End")
}

func (this *urlBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *urlBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *urlBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, urlType, "PrepareForListContents")
}

func (this *urlBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, urlType, "PrepareForMapContents")
}

func (this *urlBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, urlType, "NotifyChildContainerFinished")
}

// Pointer

type pURLBuilder struct {
}

func newPURLBuilder() ObjectBuilder {
	return &pURLBuilder{}
}

func (this *pURLBuilder) PostCacheInitBuilder() {
}

func (this *pURLBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}

func (this *pURLBuilder) Nil(dst reflect.Value) {
	dst.Set(reflect.Zero(pURLType))
}

func (this *pURLBuilder) Bool(value bool, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Bool")
}

func (this *pURLBuilder) Int(value int64, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Int")
}

func (this *pURLBuilder) Uint(value uint64, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Uint")
}

func (this *pURLBuilder) Float(value float64, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Float")
}

func (this *pURLBuilder) String(value string, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "String")
}

func (this *pURLBuilder) Bytes(value []byte, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Bytes")
}

func (this *pURLBuilder) URI(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (this *pURLBuilder) Time(value time.Time, dst reflect.Value) {
	builderPanicBadEvent(this, pURLType, "Time")
}

func (this *pURLBuilder) List() {
	builderPanicBadEvent(this, pURLType, "List")
}

func (this *pURLBuilder) Map() {
	builderPanicBadEvent(this, pURLType, "Map")
}

func (this *pURLBuilder) End() {
	builderPanicBadEvent(this, pURLType, "End")
}

func (this *pURLBuilder) Marker(id interface{}) {
	panic("TODO")
}

func (this *pURLBuilder) Reference(id interface{}) {
	panic("TODO")
}

func (this *pURLBuilder) PrepareForListContents() {
	builderPanicBadEvent(this, pURLType, "PrepareForListContents")
}

func (this *pURLBuilder) PrepareForMapContents() {
	builderPanicBadEvent(this, pURLType, "PrepareForMapContents")
}

func (this *pURLBuilder) NotifyChildContainerFinished(value reflect.Value) {
	builderPanicBadEvent(this, pURLType, "NotifyChildContainerFinished")
}
