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

type ridBuilder struct{}

var globalRidBuilder = &ridBuilder{}

func generateUrlBuilder(ctx *Context) Builder { return globalRidBuilder }
func (_this *ridBuilder) String() string      { return reflect.TypeOf(_this).String() }

func (_this *ridBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setRIDFromString(string(value), dst)
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *ridBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setRIDFromString(value, dst)
	default:
		PanicBadEvent(_this, "BuildFromStringlikeArray(%v)", arrayType)
	}
	return dst
}

// ============================================================================

type pRidBuilder struct{}

var globalPRidBuilder = &pRidBuilder{}

func generatePRidBuilder(ctx *Context) Builder { return globalPRidBuilder }
func (_this *pRidBuilder) String() string      { return reflect.TypeOf(_this).String() }

func (_this *pRidBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	return dst
}

func (_this *pRidBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setPRIDFromString(string(value), dst)
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *pRidBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeResourceID:
		setPRIDFromString(value, dst)
	default:
		PanicBadEvent(_this, "BuildFromStringlikeArray(%v)", arrayType)
	}
	return dst
}

// ============================================================================

type ridCatBuilder struct{}

var globalRidCatBuilder = &ridCatBuilder{}

func generateRidCatBuilder(ctx *Context) Builder { return globalRidCatBuilder }
func (_this *ridCatBuilder) String() string      { return reflect.TypeOf(_this).String() }

func (_this *ridCatBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		ctx.CompleteRIDCat(string(value))
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *ridCatBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		ctx.CompleteRIDCat(value)
	default:
		PanicBadEvent(_this, "BuildFromStringlikeArray(%v)", arrayType)
	}
	return dst
}
