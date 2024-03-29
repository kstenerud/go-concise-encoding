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
	"math/big"
	"reflect"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
)

type recordTypeBuilder struct{}

var globalRecordTypeBuilder = &recordTypeBuilder{}

func generateRecordTypeBuilder(ctx *Context) Builder { return globalRecordTypeBuilder }
func (_this *recordTypeBuilder) String() string      { return reflect.TypeOf(_this).String() }

func (_this *recordTypeBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromBool(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromInt(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromUint(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromBigInt(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromFloat(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromBigFloat(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromDecimalFloat(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromBigDecimalFloat(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromUID(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromArray(c, arrayType, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromStringlikeArray(c, arrayType, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildFromTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	ctx.AddRecordTypeKey(func(c *Context, builder Builder) {
		builder.BuildFromTime(c, value, unusedValue)
	})
	return dst
}
func (_this *recordTypeBuilder) BuildEndContainer(ctx *Context) {
	ctx.EndRecordType()
}

func (_this *recordTypeBuilder) BuildArtificiallyEndContainer(ctx *Context) {
	_this.BuildEndContainer(ctx)
}
