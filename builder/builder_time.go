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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-compact-time"
)

// Go Time

type timeBuilder struct {
	// Template Data
	session *Session
}

func newTimeBuilder() ObjectBuilder {
	return &timeBuilder{}
}

func (_this *timeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *timeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypeTime, name, args...)
}

func (_this *timeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *timeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *timeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *timeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
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

// ============================================================================

type compactTimeBuilder struct {
	// Template Data
	session *Session
}

func newCompactTimeBuilder() ObjectBuilder {
	return &compactTimeBuilder{}
}

func (_this *compactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *compactTimeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypeCompactTime, name, args...)
}

func (_this *compactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *compactTimeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *compactTimeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *compactTimeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *compactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(*t))
}

func (_this *compactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value))
}

// ============================================================================

type pCompactTimeBuilder struct {
	// Template Data
	session *Session
}

func newPCompactTimeBuilder() ObjectBuilder {
	return &pCompactTimeBuilder{}
}

func (_this *pCompactTimeBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *pCompactTimeBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, common.TypePCompactTime, name, args...)
}

func (_this *pCompactTimeBuilder) InitTemplate(session *Session) {
	_this.session = session
}

func (_this *pCompactTimeBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *pCompactTimeBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *pCompactTimeBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*compact_time.Time)(nil)))
}

func (_this *pCompactTimeBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}

func (_this *pCompactTimeBuilder) BuildFromTime(value time.Time, dst reflect.Value) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(t))
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(value *compact_time.Time, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value))
}
