// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package concise_encoding

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
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

// RootBuilder adapts DataEventReceiver to ObjectBuilder, coordinates the
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
func (this *RootBuilder) BuildFromNil(dst reflect.Value) {
	this.currentBuilder.BuildFromNil(dst)
}
func (this *RootBuilder) BuildFromBool(value bool, dst reflect.Value) {
	this.currentBuilder.BuildFromBool(value, dst)
}
func (this *RootBuilder) BuildFromInt(value int64, dst reflect.Value) {
	this.currentBuilder.BuildFromInt(value, dst)
}
func (this *RootBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	this.currentBuilder.BuildFromUint(value, dst)
}
func (this *RootBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	this.currentBuilder.BuildFromBigInt(value, dst)
}
func (this *RootBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	this.currentBuilder.BuildFromFloat(value, dst)
}
func (this *RootBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	this.currentBuilder.BuildFromDecimalFloat(value, dst)
}
func (this *RootBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	this.currentBuilder.BuildFromBigDecimalFloat(value, dst)
}
func (this *RootBuilder) BuildFromString(value string, dst reflect.Value) {
	this.currentBuilder.BuildFromString(value, dst)
}
func (this *RootBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	this.currentBuilder.BuildFromBytes(value, dst)
}
func (this *RootBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	this.currentBuilder.BuildFromURI(value, dst)
}
func (this *RootBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	this.currentBuilder.BuildFromTime(value, dst)
}
func (this *RootBuilder) BuildBeginList() {
	this.currentBuilder.BuildBeginList()
}
func (this *RootBuilder) BuildBeginMap() {
	this.currentBuilder.BuildBeginMap()
}
func (this *RootBuilder) BuildEndContainer() {
	this.currentBuilder.BuildEndContainer()
}
func (this *RootBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: RootBuilder.Marker")
}
func (this *RootBuilder) BuildFromReference(id interface{}) {
	panic("TODO: RootBuilder.Reference")
}
func (this *RootBuilder) PrepareForListContents() {
	panic(fmt.Errorf("BUG: %v cannot respond to event \"PrepareForListContents\"", reflect.TypeOf(this)))
}
func (this *RootBuilder) PrepareForMapContents() {
	panic(fmt.Errorf("BUG: %v cannot respond to event \"PrepareForMapContents\"", reflect.TypeOf(this)))
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
	this.BuildFromNil(this.object)
}
func (this *RootBuilder) OnBool(value bool) {
	this.BuildFromBool(value, this.object)
}
func (this *RootBuilder) OnTrue() {
	this.BuildFromBool(true, this.object)
}
func (this *RootBuilder) OnFalse() {
	this.BuildFromBool(false, this.object)
}
func (this *RootBuilder) OnPositiveInt(value uint64) {
	this.BuildFromUint(value, this.object)
}
func (this *RootBuilder) OnNegativeInt(value uint64) {
	// TODO: big int when too big
	this.BuildFromInt(-int64(value), this.object)
}
func (this *RootBuilder) OnInt(value int64) {
	this.BuildFromInt(value, this.object)
}
func (this *RootBuilder) OnBigInt(value *big.Int) {
	this.BuildFromBigInt(value, this.object)
}
func (this *RootBuilder) OnBinaryFloat(value float64) {
	this.BuildFromFloat(value, this.object)
}
func (this *RootBuilder) OnDecimalFloat(value compact_float.DFloat) {
	this.BuildFromDecimalFloat(value, this.object)
}
func (this *RootBuilder) OnBigDecimalFloat(value *apd.Decimal) {
	this.BuildFromBigDecimalFloat(value, this.object)
}
func (this *RootBuilder) OnNan(signaling bool) {
	nan := quietNan
	if signaling {
		nan = signalingNan
	}
	this.BuildFromFloat(nan, this.object)
}
func (this *RootBuilder) OnComplex(value complex128) {
	panic("TODO: RootBuilder.OnComplex")
}
func (this *RootBuilder) OnUUID(value []byte) {
	panic("TODO: RootBuilder.OnUUID")
}
func (this *RootBuilder) OnTime(value time.Time) {
	this.BuildFromTime(value, this.object)
}
func (this *RootBuilder) OnCompactTime(value *compact_time.Time) {
	t, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	this.BuildFromTime(t, this.object)
}
func (this *RootBuilder) OnBytes(value []byte) {
	this.BuildFromBytes(value, this.object)
}
func (this *RootBuilder) OnString(value string) {
	this.BuildFromString(value, this.object)
}
func (this *RootBuilder) OnURI(value string) {
	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	this.BuildFromURI(u, this.object)
}
func (this *RootBuilder) OnCustom(value []byte) {
	panic("TODO: RootBuilder.OnCustom")
}
func (this *RootBuilder) OnBytesBegin() {
	panic("TODO: RootBuilder.OnBytesBegin")
}
func (this *RootBuilder) OnStringBegin() {
	panic("TODO: RootBuilder.OnStringBegin")
}
func (this *RootBuilder) OnURIBegin() {
	panic("TODO: RootBuilder.OnURIBegin")
}
func (this *RootBuilder) OnCustomBegin() {
	panic("TODO: RootBuilder.OnCustomBegin")
}
func (this *RootBuilder) OnArrayChunk(length uint64, isFinalChunk bool) {
	panic("TODO: RootBuilder.OnArrayChunk")
}
func (this *RootBuilder) OnArrayData(data []byte) {
	panic("TODO: RootBuilder.OnArrayData")
}
func (this *RootBuilder) OnList() {
	this.BuildBeginList()
}
func (this *RootBuilder) OnMap() {
	this.BuildBeginMap()
}
func (this *RootBuilder) OnMarkup() {
	panic("TODO: RootBuilder.OnMarkup")
}
func (this *RootBuilder) OnMetadata() {
	panic("TODO: RootBuilder.OnMetadata")
}
func (this *RootBuilder) OnComment() {
	panic("TODO: RootBuilder.OnComment")
}
func (this *RootBuilder) OnEnd() {
	this.BuildEndContainer()
}
func (this *RootBuilder) OnMarker() {
	panic("TODO: RootBuilder.OnMarker")
}
func (this *RootBuilder) OnReference() {
	panic("TODO: RootBuilder.OnReference")
}
func (this *RootBuilder) OnEndDocument() {
	// Nothing to do
}
