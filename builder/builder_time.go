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
	"time"

	"github.com/kstenerud/go-compact-time"
)

// Go Time

var globalTimeBuilder = &timeBuilder{}

type timeBuilder struct{}

func generateTimeBuilder(ctx *Context) Builder { return globalTimeBuilder }
func (_this *timeBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *timeBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *timeBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	v, err := value.AsGoTime()
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(v))
	return dst
}

// ============================================================================

var globalCompactTimeBuilder = &compactTimeBuilder{}

type compactTimeBuilder struct{}

func generateCompactTimeBuilder(ctx *Context) Builder { return globalCompactTimeBuilder }
func (_this *compactTimeBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *compactTimeBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(compact_time.Time{}))
	ctx.NANext()
	return dst
}

func (_this *compactTimeBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(t))
	return dst
}

func (_this *compactTimeBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

// ============================================================================

var globalPCompactTimeBuilder = &pCompactTimeBuilder{}

type pCompactTimeBuilder struct{}

func generatePCompactTimeBuilder(ctx *Context) Builder { return &pCompactTimeBuilder{} }
func (_this *pCompactTimeBuilder) String() string            { return reflect.TypeOf(_this).String() }

func (_this *pCompactTimeBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(compact_time.Time{}))
	ctx.NANext()
	return dst
}

func (_this *pCompactTimeBuilder) BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		panic(err)
	}
	dst.Set(reflect.ValueOf(t))
	return dst
}

func (_this *pCompactTimeBuilder) BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}
