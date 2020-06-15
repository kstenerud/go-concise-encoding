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

// RootBuilder adapts DataEventReceiver to ObjectBuilder, coordinating the build.
// Use GetBuiltObject() to fetch the final result.
type RootBuilder struct {
	dstType        reflect.Type
	currentBuilder ObjectBuilder
	object         reflect.Value
}

// -----------
// RootBuilder
// -----------

func NewRootBuilder(dstType reflect.Type, options *BuilderOptions) *RootBuilder {
	this := &RootBuilder{}
	this.Init(dstType, options)
	return this
}

func (this *RootBuilder) Init(dstType reflect.Type, options *BuilderOptions) {
	this.dstType = dstType
	this.object = reflect.New(dstType).Elem()

	builder := getBuilderForType(dstType).CloneFromTemplate(this, this, applyDefaultBuilderOptions(options))
	if builder.IsContainerOnly() {
		builder = newTopLevelContainerBuilder(this, builder)
	}
	this.currentBuilder = builder
}

func (this *RootBuilder) GetBuiltObject() interface{} {
	// TODO: Verify this behavior
	if !this.object.IsValid() {
		return nil
	}

	v := this.object
	switch v.Kind() {
	case reflect.Struct, reflect.Array:
		return v.Addr().Interface()
	default:
		if this.object.CanInterface() {
			return this.object.Interface()
		}
		return nil
	}
}

func (this *RootBuilder) setCurrentBuilder(builder ObjectBuilder) {
	this.currentBuilder = builder
}

// -------------
// ObjectBuilder
// -------------

func (this *RootBuilder) IsContainerOnly() bool {
	panic(fmt.Errorf("BUG: IsContainerOnly should never be called on RootBuilder"))
}

func (this *RootBuilder) PostCacheInitBuilder() {
	panic(fmt.Errorf("BUG: PostCacheInitBuilder should never be called on RootBuilder"))
}
func (this *RootBuilder) CloneFromTemplate(_ *RootBuilder, _ ObjectBuilder, _ *BuilderOptions) ObjectBuilder {
	panic(fmt.Errorf("BUG: CloneFromTemplate should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromNil(_ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromNil should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBool should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromInt should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUint should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigInt should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromFloat should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigFloat should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromDecimalFloat should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigDecimalFloat should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUUID should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromString(_ string, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromString should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromBytes(_ []byte, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBytes should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromURI should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromTime should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromCompactTime should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildBeginList() {
	panic(fmt.Errorf("BUG: BuildBeginList should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildBeginMap() {
	panic(fmt.Errorf("BUG: BuildBeginMap should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildEndContainer() {
	panic(fmt.Errorf("BUG: BuildEndContainer should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromMarker(_ interface{}) {
	panic(fmt.Errorf("BUG: BuildFromMarker should never be called on RootBuilder"))
}
func (this *RootBuilder) BuildFromReference(_ interface{}) {
	panic(fmt.Errorf("BUG: BuildFromReference should never be called on RootBuilder"))
}
func (this *RootBuilder) PrepareForListContents() {
	panic(fmt.Errorf("BUG: PrepareForListContents should never be called on RootBuilder"))
}
func (this *RootBuilder) PrepareForMapContents() {
	panic(fmt.Errorf("BUG: PrepareForMapContents should never be called on RootBuilder"))
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
	this.currentBuilder.BuildFromNil(this.object)
}
func (this *RootBuilder) OnBool(value bool) {
	this.currentBuilder.BuildFromBool(value, this.object)
}
func (this *RootBuilder) OnTrue() {
	this.currentBuilder.BuildFromBool(true, this.object)
}
func (this *RootBuilder) OnFalse() {
	this.currentBuilder.BuildFromBool(false, this.object)
}
func (this *RootBuilder) OnPositiveInt(value uint64) {
	this.currentBuilder.BuildFromUint(value, this.object)
}
func (this *RootBuilder) OnNegativeInt(value uint64) {
	if value <= 0x7fffffffffffffff {
		this.currentBuilder.BuildFromInt(-(int64(value)), this.object)
		return
	}
	bi := &big.Int{}
	bi.SetUint64(value)
	this.currentBuilder.BuildFromBigInt(bi.Neg(bi), this.object)
}
func (this *RootBuilder) OnInt(value int64) {
	this.currentBuilder.BuildFromInt(value, this.object)
}
func (this *RootBuilder) OnBigInt(value *big.Int) {
	this.currentBuilder.BuildFromBigInt(value, this.object)
}
func (this *RootBuilder) OnFloat(value float64) {
	this.currentBuilder.BuildFromFloat(value, this.object)
}
func (this *RootBuilder) OnBigFloat(value *big.Float) {
	this.currentBuilder.BuildFromBigFloat(value, this.object)
}
func (this *RootBuilder) OnDecimalFloat(value compact_float.DFloat) {
	this.currentBuilder.BuildFromDecimalFloat(value, this.object)
}
func (this *RootBuilder) OnBigDecimalFloat(value *apd.Decimal) {
	this.currentBuilder.BuildFromBigDecimalFloat(value, this.object)
}
func (this *RootBuilder) OnNan(signaling bool) {
	nan := quietNan
	if signaling {
		nan = signalingNan
	}
	this.currentBuilder.BuildFromFloat(nan, this.object)
}
func (this *RootBuilder) OnComplex(value complex128) {
	panic("TODO: RootBuilder.OnComplex")
}
func (this *RootBuilder) OnUUID(value []byte) {
	this.currentBuilder.BuildFromUUID(value, this.object)
}
func (this *RootBuilder) OnTime(value time.Time) {
	this.currentBuilder.BuildFromTime(value, this.object)
}
func (this *RootBuilder) OnCompactTime(value *compact_time.Time) {
	this.currentBuilder.BuildFromCompactTime(value, this.object)
}
func (this *RootBuilder) OnBytes(value []byte) {
	this.currentBuilder.BuildFromBytes(value, this.object)
}
func (this *RootBuilder) OnString(value string) {
	this.currentBuilder.BuildFromString(value, this.object)
}
func (this *RootBuilder) OnURI(value string) {
	u, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	this.currentBuilder.BuildFromURI(u, this.object)
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
	this.currentBuilder.BuildBeginList()
}
func (this *RootBuilder) OnMap() {
	this.currentBuilder.BuildBeginMap()
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
	this.currentBuilder.BuildEndContainer()
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
