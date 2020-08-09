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

// Go Time

type timeBuilder struct {
	// Static data
	session *Session
}

func newTimeBuilder() ObjectBuilder {
	return &timeBuilder{}
}

func (_this *timeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *timeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *timeBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *timeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *timeBuilder) BuildFromNil(dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Nil")
}

func (_this *timeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Bool")
}

func (_this *timeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Int")
}

func (_this *timeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Uint")
}

func (_this *timeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "BigInt")
}

func (_this *timeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Float")
}

func (_this *timeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "BigFloat")
}

func (_this *timeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "DecimalFloat")
}

func (_this *timeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "BigDecimalFloat")
}

func (_this *timeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "UUID")
}

func (_this *timeBuilder) BuildFromString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "String")
}

func (_this *timeBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "VerbatimString")
}

func (_this *timeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Bytes")
}

func (_this *timeBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *timeBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *timeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "URI")
}

func (_this *timeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *timeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	v, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(v))
}

func (_this *timeBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "List")
}

func (_this *timeBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Map")
}

func (_this *timeBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "ContainerEnd")
}

func (_this *timeBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Marker")
}

func (_this *timeBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "Reference")
}

func (_this *timeBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "PrepareForListContents")
}

func (_this *timeBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "PrepareForMapContents")
}

func (_this *timeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeTime, "NotifyChildContainerFinished")
}

// ============================================================================

type compactTimeBuilder struct {
	// Static data
	session *Session
}

func newCompactTimeBuilder() ObjectBuilder {
	return &compactTimeBuilder{}
}

func (_this *compactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *compactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *compactTimeBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *compactTimeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *compactTimeBuilder) BuildFromNil(dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Nil")
}

func (_this *compactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Bool")
}

func (_this *compactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Int")
}

func (_this *compactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Uint")
}

func (_this *compactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "BigInt")
}

func (_this *compactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Float")
}

func (_this *compactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "BigFloat")
}

func (_this *compactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "DecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "BigDecimalFloat")
}

func (_this *compactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "UUID")
}

func (_this *compactTimeBuilder) BuildFromString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "String")
}

func (_this *compactTimeBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "VerbatimString")
}

func (_this *compactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Bytes")
}

func (_this *compactTimeBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *compactTimeBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *compactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "URI")
}

func (_this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*compact_time.AsCompactTime(value)))
}

func (_this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

func (_this *compactTimeBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "List")
}

func (_this *compactTimeBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Map")
}

func (_this *compactTimeBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "ContainerEnd")
}

func (_this *compactTimeBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Marker")
}

func (_this *compactTimeBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "Reference")
}

func (_this *compactTimeBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "PrepareForListContents")
}

func (_this *compactTimeBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "PrepareForMapContents")
}

func (_this *compactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypeCompactTime, "NotifyChildContainerFinished")
}

// ============================================================================

type pCompactTimeBuilder struct {
	// Static data
	session *Session
}

func newPCompactTimeBuilder() ObjectBuilder {
	return &pCompactTimeBuilder{}
}

func (_this *pCompactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pCompactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *pCompactTimeBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, options *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pCompactTimeBuilder) SetParent(parent ObjectBuilder) {
}

func (_this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (_this *pCompactTimeBuilder) BuildFromBool(value bool, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Bool")
}

func (_this *pCompactTimeBuilder) BuildFromInt(value int64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Int")
}

func (_this *pCompactTimeBuilder) BuildFromUint(value uint64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Uint")
}

func (_this *pCompactTimeBuilder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "BigInt")
}

func (_this *pCompactTimeBuilder) BuildFromFloat(value float64, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Float")
}

func (_this *pCompactTimeBuilder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "BigFloat")
}

func (_this *pCompactTimeBuilder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "DecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "BigDecimalFloat")
}

func (_this *pCompactTimeBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "UUID")
}

func (_this *pCompactTimeBuilder) BuildFromString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "String")
}

func (_this *pCompactTimeBuilder) BuildFromVerbatimString(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "VerbatimString")
}

func (_this *pCompactTimeBuilder) BuildFromBytes(value []byte, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Bytes")
}

func (_this *pCompactTimeBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *pCompactTimeBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		BuilderPanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *pCompactTimeBuilder) BuildFromURI(value *url.URL, dst reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "URI")
}

func (_this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(compact_time.AsCompactTime(value)))
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *pCompactTimeBuilder) BuildBeginList() {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "List")
}

func (_this *pCompactTimeBuilder) BuildBeginMap() {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Map")
}

func (_this *pCompactTimeBuilder) BuildEndContainer() {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "ContainerEnd")
}

func (_this *pCompactTimeBuilder) BuildBeginMarker(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Marker")
}

func (_this *pCompactTimeBuilder) BuildFromReference(id interface{}) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "Reference")
}

func (_this *pCompactTimeBuilder) PrepareForListContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "PrepareForListContents")
}

func (_this *pCompactTimeBuilder) PrepareForMapContents() {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "PrepareForMapContents")
}

func (_this *pCompactTimeBuilder) NotifyChildContainerFinished(value reflect.Value) {
	BuilderWithTypePanicBadEvent(_this, common.TypePCompactTime, "NotifyChildContainerFinished")
}
