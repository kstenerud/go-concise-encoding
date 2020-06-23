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

// topLevelContainerBuilder proxies the first build instruction to make sure containers
// are properly built. See BuildBeginList and BuildBeginMap.
type topLevelContainerBuilder struct {
	builder ObjectBuilder
	root    *RootBuilder
}

func newTopLevelContainerBuilder(root *RootBuilder, builder ObjectBuilder) ObjectBuilder {
	return &topLevelContainerBuilder{
		builder: builder,
		root:    root,
	}
}

func (_this *topLevelContainerBuilder) IsContainerOnly() bool {
	panic(fmt.Errorf("BUG: IsContainerOnly should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) PostCacheInitBuilder() {
	panic(fmt.Errorf("BUG: PostCacheInitBuilder should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder {
	panic(fmt.Errorf("BUG: CloneFromTemplate should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromNil(dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromNil should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromBool(value bool, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBool should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromInt(value int64, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromInt should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUint should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigInt should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromFloat should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigFloat should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromDecimalFloat should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBigDecimalFloat should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromUUID should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromString(value string, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromString should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromBytes should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromURI should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	panic(fmt.Errorf("BUG: BuildFromTime should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	panic(fmt.Errorf("BUG: topLevelContainerBuilder should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildBeginList() {
	_this.root.setCurrentBuilder(_this.builder)
	_this.builder.PrepareForListContents()
}

func (_this *topLevelContainerBuilder) BuildBeginMap() {
	_this.root.setCurrentBuilder(_this.builder)
	_this.builder.PrepareForMapContents()
}

func (_this *topLevelContainerBuilder) BuildEndContainer() {
	panic(fmt.Errorf("BUG: BuildEndContainer should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) PrepareForListContents() {
	panic(fmt.Errorf("BUG: PrepareForListContents should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) PrepareForMapContents() {
	panic(fmt.Errorf("BUG: PrepareForMapContents should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) NotifyChildContainerFinished(value reflect.Value) {
	panic(fmt.Errorf("BUG: NotifyChildContainerFinished should never be called on topLevelContainerBuilder"))
}

func (_this *topLevelContainerBuilder) BuildFromMarker(id interface{}) {
	panic("TODO: BuildFromMarker")
}

func (_this *topLevelContainerBuilder) BuildFromReference(id interface{}) {
	panic("TODO: BuildFromReference")
}
