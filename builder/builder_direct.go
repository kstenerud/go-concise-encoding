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
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-compact-time"
)

// The direct builder has an unambiguous direct mapping from build event to
// a non-pointer destination type (for example, a bool is always a bool).
type directBuilder struct {
	// Template Data
	session *Session
	dstType reflect.Type
}

func newDirectBuilder(dstType reflect.Type) ObjectBuilder {
	return &directBuilder{
		dstType: dstType,
	}
}

func (_this *directBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *directBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}

func (_this *directBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *directBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *directBuilder) BuildFromBool(value bool, dst reflect.Value) {
	dst.SetBool(value)
}

func (_this *directBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	value = common.CloneBytes(value)
	dst.Set(reflect.ValueOf(value))
}

func (_this *directBuilder) BuildFromString(value []byte, dst reflect.Value) {
	dst.SetString(string(value))
}

func (_this *directBuilder) BuildFromRID(value *url.URL, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value).Elem())
}

func (_this *directBuilder) BuildFromCustomBinary(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
	}
}

func (_this *directBuilder) BuildFromCustomText(value []byte, dst reflect.Value) {
	if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
		PanicBuildFromCustomText(_this, value, dst.Type(), err)
	}
}

func (_this *directBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setRIDFromString(string(value), dst)
	default:
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *directBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

func (_this *directBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}

// ============================================================================

// The direct builder has an unambiguous direct mapping from build event to
// a pointer destination type (for example, a *url is always a *url).
type directPtrBuilder struct {
	// Template Data
	session *Session
	dstType reflect.Type
}

func newDirectPtrBuilder(dstType reflect.Type) ObjectBuilder {
	return &directPtrBuilder{
		dstType: dstType,
	}
}

func (_this *directPtrBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *directPtrBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}

func (_this *directPtrBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *directPtrBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *directPtrBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *directPtrBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.Zero(dst.Type()))
}

func (_this *directPtrBuilder) BuildFromUUID(value []byte, dst reflect.Value) {
	dst.SetBytes(common.CloneBytes(value))
}

func (_this *directPtrBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		if err := _this.session.GetCustomBinaryBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
		}
	case events.ArrayTypeCustomText:
		if err := _this.session.GetCustomTextBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomText(_this, value, dst.Type(), err)
		}
	case events.ArrayTypeUint8:
		dst.SetBytes(common.CloneBytes(value))
	case events.ArrayTypeString:
		dst.SetString(string(value))
	case events.ArrayTypeResourceID:
		setPRIDFromString(string(value), dst)
	default:
		panic(fmt.Errorf("TODO: Add typed array support for %v", arrayType))
	}
}

func (_this *directPtrBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	// TODO: Should non-pointer stuff be here?
	dst.Set(reflect.ValueOf(value))
}

func (_this *directPtrBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}
