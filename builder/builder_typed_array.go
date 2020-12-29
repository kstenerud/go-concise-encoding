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
)

type stringBuilder struct{}

var globalStringBuilder = &stringBuilder{}

func generateStringBuilder() ObjectBuilder  { return globalStringBuilder }
func (_this *stringBuilder) String() string { return nameOf(_this) }

func (_this *stringBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	// Go doesn't have the concept of a nil string.
	dst.SetString("")
	return dst
}
func (_this *stringBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		dst.SetString(string(value))
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

// ============================================================================

type uint8ArrayBuilder struct{}

var globalUint8ArrayBuilder = &uint8ArrayBuilder{}

func generateUint8ArrayBuilder() ObjectBuilder  { return globalUint8ArrayBuilder }
func (_this *uint8ArrayBuilder) String() string { return nameOf(_this) }

func (_this *uint8ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint8:
		// TODO: Is there a more efficient way?
		for i := 0; i < len(value); i++ {
			elem := dst.Index(i)
			elem.SetUint(uint64(value[i]))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

// ============================================================================

type uint16ArrayBuilder struct{}

var globalUint16ArrayBuilder = &uint16ArrayBuilder{}

func generateUint16ArrayBuilder() ObjectBuilder  { return globalUint16ArrayBuilder }
func (_this *uint16ArrayBuilder) String() string { return nameOf(_this) }

func (_this *uint16ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint16:
		elemCount := len(value) / 2
		for i := 0; i < elemCount; i++ {
			elemValue := uint16(value[i*2]) |
				(uint16(value[i*2+1]) << 8)
			elem := dst.Index(i)
			elem.SetUint(uint64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

type uint16SliceBuilder struct{}

var globalUint16SliceBuilder = &uint16SliceBuilder{}

func generateUint16SliceBuilder() ObjectBuilder  { return globalUint16SliceBuilder }
func (_this *uint16SliceBuilder) String() string { return nameOf(_this) }

func (_this *uint16SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint16:
		elemCount := len(value) / 2
		slice := make([]uint16, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = uint16(value[i*2]) |
				(uint16(value[i*2+1]) << 8)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}
