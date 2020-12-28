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

	"github.com/kstenerud/go-compact-time"
)

// The direct builder has an unambiguous direct mapping from build event to
// a non-pointer destination type (for example, a bool is always a bool).
type directBuilder struct{}

var globalDirectBuilder = &directBuilder{}

func generateDirectBuilder() ObjectBuilder  { return globalDirectBuilder }
func (_this *directBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *directBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	dst.SetBool(value)
	return dst
}

func (_this *directBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	value = common.CloneBytes(value)
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *directBuilder) BuildFromString(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	dst.SetString(string(value))
	return dst
}

func (_this *directBuilder) BuildFromRID(ctx *Context, value *url.URL, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value).Elem())
	return dst
}

func (_this *directBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setRIDFromString(string(value), dst)
	default:
		PanicBadEvent(_this, "TypedArray(%v)", arrayType)
	}
	return dst
}

func (_this *directBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *directBuilder) BuildFromCompactTime(ctx *Context, value *compact_time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

// ============================================================================

// The direct builder has an unambiguous direct mapping from build event to
// a pointer destination type (for example, a *url is always a *url).
type directPtrBuilder struct{}

var globalDirectPtrBuilder = &directPtrBuilder{}

func generateDirectPtrBuilder() ObjectBuilder  { return globalDirectPtrBuilder }
func (_this *directPtrBuilder) String() string { return reflect.TypeOf(_this).String() }

func (_this *directPtrBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	return dst
}

func (_this *directPtrBuilder) BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	dst.SetBytes(common.CloneBytes(value))
	return dst
}

func (_this *directPtrBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		if err := ctx.CustomBinaryBuildFunction(value, dst); err != nil {
			PanicBuildFromCustomBinary(_this, value, dst.Type(), err)
		}
	case events.ArrayTypeCustomText:
		if err := ctx.CustomTextBuildFunction(value, dst); err != nil {
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
	return dst
}

func (_this *directPtrBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	// TODO: Should non-pointer stuff be here?
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *directPtrBuilder) BuildFromCompactTime(ctx *Context, value *compact_time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}
