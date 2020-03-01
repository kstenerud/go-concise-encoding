package concise_encoding

import (
	"math"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-compact-time"
)

func (this *RootBuilder) GetBuiltObject() interface{} {
	// TODO: Verify this behavior
	if this.object.IsValid() {
		v := this.object
		switch v.Kind() {
		case reflect.Struct, reflect.Array:
			return v.Addr().Interface()
		default:
			if this.object.CanInterface() {
				return this.object.Interface()
			}
		}
	}
	return nil
}

// RootBuilder adapts ObjectIteratorCallbacks to ObjectBuilder, coordinates the
// build, and provides GetBuiltObject() for fetching the final result.
type RootBuilder struct {
	dstType        reflect.Type
	currentBuilder ObjectBuilder
	object         reflect.Value
}

// -----------
// RootBuilder
// -----------

func newRootBuilder(dstType reflect.Type) *RootBuilder {
	this := &RootBuilder{
		dstType: dstType,
		object:  reflect.New(dstType).Elem(),
	}

	builder := getTopLevelBuilderForType(dstType)
	this.currentBuilder = builder.CloneFromTemplate(this, this)

	return this
}

func (this *RootBuilder) setCurrentBuilder(builder ObjectBuilder) {
	this.currentBuilder = builder
}

// -------------
// ObjectBuilder
// -------------

func (this *RootBuilder) PostCacheInitBuilder() {
}

func (this *RootBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder {
	return this
}
func (this *RootBuilder) Nil(dst reflect.Value) {
	this.currentBuilder.Nil(dst)
}
func (this *RootBuilder) Bool(value bool, dst reflect.Value) {
	this.currentBuilder.Bool(value, dst)
}
func (this *RootBuilder) Int(value int64, dst reflect.Value) {
	this.currentBuilder.Int(value, dst)
}
func (this *RootBuilder) Uint(value uint64, dst reflect.Value) {
	this.currentBuilder.Uint(value, dst)
}
func (this *RootBuilder) Float(value float64, dst reflect.Value) {
	this.currentBuilder.Float(value, dst)
}
func (this *RootBuilder) String(value string, dst reflect.Value) {
	this.currentBuilder.String(value, dst)
}
func (this *RootBuilder) Bytes(value []byte, dst reflect.Value) {
	this.currentBuilder.Bytes(value, dst)
}
func (this *RootBuilder) URI(value *url.URL, dst reflect.Value) {
	this.currentBuilder.URI(value, dst)
}
func (this *RootBuilder) Time(value time.Time, dst reflect.Value) {
	this.currentBuilder.Time(value, dst)
}
func (this *RootBuilder) List() {
	this.currentBuilder.List()
}
func (this *RootBuilder) Map() {
	this.currentBuilder.Map()
}
func (this *RootBuilder) End() {
	this.currentBuilder.End()
}
func (this *RootBuilder) Marker(id interface{}) {
	panic("TODO")
}
func (this *RootBuilder) Reference(id interface{}) {
	panic("TODO")
}
func (this *RootBuilder) PrepareForListContents() {
	panic("BUG")
}
func (this *RootBuilder) PrepareForMapContents() {
	panic("BUG")
}
func (this *RootBuilder) NotifyChildContainerFinished(value reflect.Value) {
	this.object = value
}

// -----------------------
// ObjectIteratorCallbacks
// -----------------------

func (this *RootBuilder) OnVersion(version uint64) {}
func (this *RootBuilder) OnPadding(count int)      {}
func (this *RootBuilder) OnNil() {
	this.Nil(this.object)
}
func (this *RootBuilder) OnBool(value bool) {
	this.Bool(value, this.object)
}
func (this *RootBuilder) OnTrue() {
	this.Bool(true, this.object)
}
func (this *RootBuilder) OnFalse() {
	this.Bool(false, this.object)
}
func (this *RootBuilder) OnPositiveInt(value uint64) {
	this.Uint(value, this.object)
}
func (this *RootBuilder) OnNegativeInt(value uint64) {
	this.Int(-int64(value), this.object)
}
func (this *RootBuilder) OnInt(value int64) {
	this.Int(value, this.object)
}
func (this *RootBuilder) OnFloat(value float64) {
	this.Float(value, this.object)
}
func (this *RootBuilder) OnNan() {
	this.Float(math.NaN(), this.object)
}
func (this *RootBuilder) OnComplex(value complex128) {
	panic("TODO")
}
func (this *RootBuilder) OnUUID(value []byte) {
	panic("TODO")
}
func (this *RootBuilder) OnTime(value time.Time) {
	this.Time(value, this.object)
}
func (this *RootBuilder) OnCompactTime(value *compact_time.Time) {
	t, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	this.Time(t, this.object)
}
func (this *RootBuilder) OnBytes(value []byte) {
	this.Bytes(value, this.object)
}
func (this *RootBuilder) OnString(value string) {
	this.String(value, this.object)
}
func (this *RootBuilder) OnURI(value string) {
	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	this.URI(u, this.object)
}
func (this *RootBuilder) OnCustom(value []byte) {
	panic("TODO")
}
func (this *RootBuilder) OnBytesBegin() {
	panic("TODO")
}
func (this *RootBuilder) OnStringBegin() {
	panic("TODO")
}
func (this *RootBuilder) OnURIBegin() {
	panic("TODO")
}
func (this *RootBuilder) OnCustomBegin() {
	panic("TODO")
}
func (this *RootBuilder) OnArrayChunk(length uint64, isFinalChunk bool) {
	panic("TODO")
}
func (this *RootBuilder) OnArrayData(data []byte) {
	panic("TODO")
}
func (this *RootBuilder) OnList() {
	this.List()
}
func (this *RootBuilder) OnMap() {
	this.Map()
}
func (this *RootBuilder) OnMarkup() {
	panic("TODO")
}
func (this *RootBuilder) OnMetadata() {
	panic("TODO")
}
func (this *RootBuilder) OnComment() {
	panic("TODO")
}
func (this *RootBuilder) OnEnd() {
	this.End()
}
func (this *RootBuilder) OnMarker() {
	panic("TODO")
}
func (this *RootBuilder) OnReference() {
	panic("TODO")
}
func (this *RootBuilder) OnEndDocument() {
	// Nothing to do
}
