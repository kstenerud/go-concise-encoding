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

type bigIntBuilder struct {
	// Static data
	session *Session

	// Clone inserted data
	options *options.BuilderOptions
}

func newBigIntBuilder() ObjectBuilder {
	return &bigIntBuilder{}
}

func (_this *bigIntBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *bigIntBuilder) PostCacheInitBuilder(session *Session) {
	_this.session = session
}

func (_this *bigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return &bigIntBuilder{
		session: _this.session,
		options: options,
	}
}

func (_this *bigIntBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *bigIntBuilder) BuildFromNil(dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Nil")
}

func (_this *bigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Bool")
}

func (_this *bigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setBigIntFromInt(value, dst)
}

func (_this *bigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setBigIntFromUint(value, dst)
}

func (_this *bigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *bigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setBigIntFromFloat(value, dst, _this.options.FloatToBigIntMaxBase2Exponent)
}

func (_this *bigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setBigIntFromBigFloat(value, dst, _this.options.FloatToBigIntMaxBase2Exponent)
}

func (_this *bigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setBigIntFromDecimalFloat(value, dst, _this.options.FloatToBigIntMaxBase10Exponent)
}

func (_this *bigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setBigIntFromBigDecimalFloat(value, dst, _this.options.FloatToBigIntMaxBase10Exponent)
}

func (_this *bigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "UUID")
}

func (_this *bigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "String")
}

func (_this *bigIntBuilder) BuildFromVerbatimString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "VerbatimString")
}

func (_this *bigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Bytes")
}

func (_this *bigIntBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *bigIntBuilder) BuildFromCustomText(value string, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *bigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "URI")
}

func (_this *bigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Time")
}

func (_this *bigIntBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "CompactTime")
}

func (_this *bigIntBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "List")
}

func (_this *bigIntBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Map")
}

func (_this *bigIntBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "ContainerEnd")
}

func (_this *bigIntBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Marker")
}

func (_this *bigIntBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "Reference")
}

func (_this *bigIntBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "PrepareForListContents")
}

func (_this *bigIntBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "PrepareForMapContents")
}

func (_this *bigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeBigInt, "NotifyChildContainerFinished")
}

// ============================================================================

type pBigIntBuilder struct {
	// Static data
	session *Session

	// Clone inserted data
	options *options.BuilderOptions
}

func newPBigIntBuilder() ObjectBuilder {
	return &pBigIntBuilder{}
}

func (_this *pBigIntBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pBigIntBuilder) PostCacheInitBuilder(session *Session) {
	_this.session = session
}

func (_this *pBigIntBuilder) CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return &pBigIntBuilder{
		session: _this.session,
		options: options,
	}
}

func (_this *pBigIntBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pBigIntBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Int)(nil)))
}

func (_this *pBigIntBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Bool")
}

func (_this *pBigIntBuilder) BuildFromInt(value int64, dst reflect.Value) {
	setPBigIntFromInt(value, dst)
}

func (_this *pBigIntBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	setPBigIntFromUint(value, dst)
}

func (_this *pBigIntBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pBigIntBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	setPBigIntFromFloat(value, dst, _this.options.FloatToBigIntMaxBase2Exponent)
}

func (_this *pBigIntBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	setPBigIntFromBigFloat(value, dst, _this.options.FloatToBigIntMaxBase2Exponent)
}

func (_this *pBigIntBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	setPBigIntFromDecimalFloat(value, dst, _this.options.FloatToBigIntMaxBase10Exponent)
}

func (_this *pBigIntBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	setPBigIntFromBigDecimalFloat(value, dst, _this.options.FloatToBigIntMaxBase10Exponent)
}

func (_this *pBigIntBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "UUID")
}

func (_this *pBigIntBuilder) BuildFromString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "String")
}

func (_this *pBigIntBuilder) BuildFromVerbatimString(value string, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "VerbatimString")
}

func (_this *pBigIntBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Bytes")
}

func (_this *pBigIntBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *pBigIntBuilder) BuildFromCustomText(value string, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *pBigIntBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "URI")
}

func (_this *pBigIntBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Time")
}

func (_this *pBigIntBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "CompactTime")
}

func (_this *pBigIntBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "List")
}

func (_this *pBigIntBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Map")
}

func (_this *pBigIntBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "ContainerEnd")
}

func (_this *pBigIntBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Marker")
}

func (_this *pBigIntBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "Reference")
}

func (_this *pBigIntBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "PrepareForListContents")
}

func (_this *pBigIntBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "PrepareForMapContents")
}

func (_this *pBigIntBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePBigInt, "NotifyChildContainerFinished")
}
