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
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

type stringBuilder struct{}

func newStringBuilder() ObjectBuilder       { return &stringBuilder{} }
func (_this *stringBuilder) String() string { return nameOf(_this) }
func (_this *stringBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf(""), name, args...)
}
func (_this *stringBuilder) InitTemplate(_ *Session) {}
func (_this *stringBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *stringBuilder) SetParent(_ ObjectBuilder) {}

type uint8ArrayBuilder struct{}

func newUint8ArrayBuilder() ObjectBuilder       { return &uint8ArrayBuilder{} }
func (_this *uint8ArrayBuilder) String() string { return nameOf(_this) }
func (_this *uint8ArrayBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf(uint8(0)), name, args...)
}
func (_this *uint8ArrayBuilder) InitTemplate(_ *Session) {}
func (_this *uint8ArrayBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *uint8ArrayBuilder) SetParent(_ ObjectBuilder) {}

func (_this *uint8ArrayBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeUint8:
		// TODO: Is there a more efficient way?
		for i := 0; i < len(value); i++ {
			elem := dst.Index(i)
			elem.SetUint(uint64(value[i]))
		}
	default:
		_this.panicBadEvent("BuildFromArray(%v)", arrayType)
	}
}

type uint16ArrayBuilder struct{}

func newUint16ArrayBuilder() ObjectBuilder       { return &uint16ArrayBuilder{} }
func (_this *uint16ArrayBuilder) String() string { return nameOf(_this) }
func (_this *uint16ArrayBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf(uint16(0)), name, args...)
}
func (_this *uint16ArrayBuilder) InitTemplate(_ *Session) {}
func (_this *uint16ArrayBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *uint16ArrayBuilder) SetParent(_ ObjectBuilder) {}

func (_this *uint16ArrayBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	panic("TODO")
}

func (_this *stringBuilder) BuildFromNil(dst reflect.Value) {
	// Go doesn't have the concept of a nil string.
	dst.SetString("")
}
func (_this *stringBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeString:
		dst.SetString(string(value))
	default:
		_this.panicBadEvent("BuildFromArray(%v)", arrayType)
	}
}
