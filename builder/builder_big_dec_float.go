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

package builder

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type bigDecimalFloatBuilder struct {
	// Static data
	session *Session
}

func newBigDecimalFloatBuilder() ObjectBuilder {
	return &bigDecimalFloatBuilder{}
}

func (_this *bigDecimalFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *bigDecimalFloatBuilder) PostCacheInitBuilder(session *Session) {
	_this.session = session
}

func (_this *bigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *bigDecimalFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Nil")
}

func (_this *bigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Bool")
}

func (_this *bigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigDecimalFloatFromInt(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigDecimalFloatFromUint(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setBigDecimalFloatFromBigInt(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigDecimalFloatFromFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigDecimalFloatFromBigFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigDecimalFloatFromDecimalFloat(value, dst)
}

func (_this *bigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *bigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "UUID")
}

func (_this *bigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "String")
}

func (_this *bigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Bytes")
}

func (_this *bigDecimalFloatBuilder) BuildFromCustom(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustom(_this, value, dst.Type(), err)
	}
}

func (_this *bigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "URI")
}

func (_this *bigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Time")
}

func (_this *bigDecimalFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "CompactTime")
}

func (_this *bigDecimalFloatBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "List")
}

func (_this *bigDecimalFloatBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Map")
}

func (_this *bigDecimalFloatBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "ContainerEnd")
}

func (_this *bigDecimalFloatBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Marker")
}

func (_this *bigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "Reference")
}

func (_this *bigDecimalFloatBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "PrepareForListContents")
}

func (_this *bigDecimalFloatBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "PrepareForMapContents")
}

func (_this *bigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "NotifyChildContainerFinished")
}

// ============================================================================

type pBigDecimalFloatBuilder struct {
	// Static data
	session *Session
}

func newPBigDecimalFloatBuilder() ObjectBuilder {
	return &pBigDecimalFloatBuilder{}
}

func (_this *pBigDecimalFloatBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pBigDecimalFloatBuilder) PostCacheInitBuilder(session *Session) {
	_this.session = session
}

func (_this *pBigDecimalFloatBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pBigDecimalFloatBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pBigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
}

func (_this *pBigDecimalFloatBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Bool")
}

func (_this *pBigDecimalFloatBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigDecimalFloatFromInt(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigDecimalFloatFromUint(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	setPBigDecimalFloatFromBigInt(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigDecimalFloatFromFloat(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigDecimalFloatFromBigFloat(value, dst)
}

func (_this *pBigDecimalFloatBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value.APD()))
}

func (_this *pBigDecimalFloatBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pBigDecimalFloatBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "UUID")
}

func (_this *pBigDecimalFloatBuilder) BuildFromString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "String")
}

func (_this *pBigDecimalFloatBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Bytes")
}

func (_this *pBigDecimalFloatBuilder) BuildFromCustom(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustom(_this, value, dst.Type(), err)
	}
}

func (_this *pBigDecimalFloatBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "URI")
}

func (_this *pBigDecimalFloatBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Time")
}

func (_this *pBigDecimalFloatBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigDecimalFloat, "CompactTime")
}

func (_this *pBigDecimalFloatBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "List")
}

func (_this *pBigDecimalFloatBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Map")
}

func (_this *pBigDecimalFloatBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "ContainerEnd")
}

func (_this *pBigDecimalFloatBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Marker")
}

func (_this *pBigDecimalFloatBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "Reference")
}

func (_this *pBigDecimalFloatBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "PrepareForListContents")
}

func (_this *pBigDecimalFloatBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "PrepareForMapContents")
}

func (_this *pBigDecimalFloatBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigDecimalFloat, "NotifyChildContainerFinished")
}
