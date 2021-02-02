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
	"math"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/events"
)

type stringBuilder struct{}

var globalStringBuilder = &stringBuilder{}

func generateStringBuilder(ctx *Context) ObjectBuilder { return globalStringBuilder }
func (_this *stringBuilder) String() string            { return nameOf(_this) }

func (_this *stringBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	// Go doesn't have the concept of a nil string.
	dst.SetString("")
	ctx.NANext()
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

func (_this *stringBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		dst.SetString(value)
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

// ============================================================================

type uint8ArrayBuilder struct {
	dstType reflect.Type
}

func newUint8ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &uint8ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *uint8ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *uint8ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint8:
		for i := 0; i < len(value); i++ {
			elem := dst.Index(i)
			elem.SetUint(uint64(value[i]))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *uint8ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type uint8SliceBuilder struct{}

var globalUint8SliceBuilder = &uint8SliceBuilder{}

func generateUint8SliceBuilder(ctx *Context) ObjectBuilder { return globalUint8SliceBuilder }
func (_this *uint8SliceBuilder) String() string            { return nameOf(_this) }

func (_this *uint8SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *uint8SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint8:
		dst.SetBytes(common.CloneBytes(value))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *uint8SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToUint8SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type uint16ArrayBuilder struct {
	dstType reflect.Type
}

func newUint16ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &uint16ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *uint16ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

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

func (_this *uint16ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type uint16SliceBuilder struct{}

var globalUint16SliceBuilder = &uint16SliceBuilder{}

func generateUint16SliceBuilder(ctx *Context) ObjectBuilder { return globalUint16SliceBuilder }
func (_this *uint16SliceBuilder) String() string            { return nameOf(_this) }

func (_this *uint16SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

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

func (_this *uint16SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToUint16SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type uint32ArrayBuilder struct {
	dstType reflect.Type
}

func newUint32ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &uint32ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *uint32ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *uint32ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint32:
		elemCount := len(value) / 4
		for i := 0; i < elemCount; i++ {
			elemValue := uint32(value[i*4]) |
				(uint32(value[i*4+1]) << 8) |
				(uint32(value[i*4+2]) << 16) |
				(uint32(value[i*4+3]) << 24)
			elem := dst.Index(i)
			elem.SetUint(uint64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *uint32ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type uint32SliceBuilder struct{}

var globalUint32SliceBuilder = &uint32SliceBuilder{}

func generateUint32SliceBuilder(ctx *Context) ObjectBuilder { return globalUint32SliceBuilder }
func (_this *uint32SliceBuilder) String() string            { return nameOf(_this) }

func (_this *uint32SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *uint32SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint32:
		elemCount := len(value) / 4
		slice := make([]uint32, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = uint32(value[i*4]) |
				(uint32(value[i*4+1]) << 8) |
				(uint32(value[i*4+2]) << 16) |
				(uint32(value[i*4+3]) << 24)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *uint32SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToUint32SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type uint64ArrayBuilder struct {
	dstType reflect.Type
}

func newUint64ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &uint64ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *uint64ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *uint64ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint64:
		elemCount := len(value) / 8
		for i := 0; i < elemCount; i++ {
			elemValue := uint64(value[i*8]) |
				(uint64(value[i*8+1]) << 8) |
				(uint64(value[i*8+2]) << 16) |
				(uint64(value[i*8+3]) << 24) |
				(uint64(value[i*8+4]) << 32) |
				(uint64(value[i*8+5]) << 40) |
				(uint64(value[i*8+6]) << 48) |
				(uint64(value[i*8+7]) << 56)
			elem := dst.Index(i)
			elem.SetUint(uint64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *uint64ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type uint64SliceBuilder struct{}

var globalUint64SliceBuilder = &uint64SliceBuilder{}

func generateUint64SliceBuilder(ctx *Context) ObjectBuilder { return globalUint64SliceBuilder }
func (_this *uint64SliceBuilder) String() string            { return nameOf(_this) }

func (_this *uint64SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *uint64SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint64:
		elemCount := len(value) / 8
		slice := make([]uint64, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = uint64(value[i*8]) |
				(uint64(value[i*8+1]) << 8) |
				(uint64(value[i*8+2]) << 16) |
				(uint64(value[i*8+3]) << 24) |
				(uint64(value[i*8+4]) << 32) |
				(uint64(value[i*8+5]) << 40) |
				(uint64(value[i*8+6]) << 48) |
				(uint64(value[i*8+7]) << 56)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *uint64SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToUint64SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type int8ArrayBuilder struct {
	dstType reflect.Type
}

func newInt8ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &int8ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *int8ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *int8ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt8:
		elemCount := len(value)
		for i := 0; i < elemCount; i++ {
			elemValue := int8(value[i])
			elem := dst.Index(i)
			elem.SetInt(int64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *int8ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type int8SliceBuilder struct{}

var globalInt8SliceBuilder = &int8SliceBuilder{}

func generateInt8SliceBuilder(ctx *Context) ObjectBuilder { return globalInt8SliceBuilder }
func (_this *int8SliceBuilder) String() string            { return nameOf(_this) }

func (_this *int8SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *int8SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt8:
		elemCount := len(value)
		slice := make([]int8, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = int8(value[i])
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *int8SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToInt8SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type int16ArrayBuilder struct {
	dstType reflect.Type
}

func newInt16ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &int16ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *int16ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *int16ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt16:
		elemCount := len(value) / 2
		for i := 0; i < elemCount; i++ {
			elemValue := int16(value[i*2]) |
				(int16(value[i*2+1]) << 8)
			elem := dst.Index(i)
			elem.SetInt(int64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *int16ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type int16SliceBuilder struct{}

var globalInt16SliceBuilder = &int16SliceBuilder{}

func generateInt16SliceBuilder(ctx *Context) ObjectBuilder { return globalInt16SliceBuilder }
func (_this *int16SliceBuilder) String() string            { return nameOf(_this) }

func (_this *int16SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *int16SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt16:
		elemCount := len(value) / 2
		slice := make([]int16, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = int16(value[i*2]) |
				(int16(value[i*2+1]) << 8)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *int16SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToInt16SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type int32ArrayBuilder struct {
	dstType reflect.Type
}

func newInt32ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &int32ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *int32ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *int32ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt32:
		elemCount := len(value) / 4
		for i := 0; i < elemCount; i++ {
			elemValue := int32(value[i*4]) |
				(int32(value[i*4+1]) << 8) |
				(int32(value[i*4+2]) << 16) |
				(int32(value[i*4+3]) << 24)
			elem := dst.Index(i)
			elem.SetInt(int64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *int32ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type int32SliceBuilder struct{}

var globalInt32SliceBuilder = &int32SliceBuilder{}

func generateInt32SliceBuilder(ctx *Context) ObjectBuilder { return globalInt32SliceBuilder }
func (_this *int32SliceBuilder) String() string            { return nameOf(_this) }

func (_this *int32SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *int32SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt32:
		elemCount := len(value) / 4
		slice := make([]int32, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = int32(value[i*4]) |
				(int32(value[i*4+1]) << 8) |
				(int32(value[i*4+2]) << 16) |
				(int32(value[i*4+3]) << 24)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *int32SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToInt32SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type int64ArrayBuilder struct {
	dstType reflect.Type
}

func newInt64ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &int64ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *int64ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *int64ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt64:
		elemCount := len(value) / 8
		for i := 0; i < elemCount; i++ {
			elemValue := int64(value[i*8]) |
				(int64(value[i*8+1]) << 8) |
				(int64(value[i*8+2]) << 16) |
				(int64(value[i*8+3]) << 24) |
				(int64(value[i*8+4]) << 32) |
				(int64(value[i*8+5]) << 40) |
				(int64(value[i*8+6]) << 48) |
				(int64(value[i*8+7]) << 56)
			elem := dst.Index(i)
			elem.SetInt(int64(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *int64ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type int64SliceBuilder struct{}

var globalInt64SliceBuilder = &int64SliceBuilder{}

func generateInt64SliceBuilder(ctx *Context) ObjectBuilder { return globalInt64SliceBuilder }
func (_this *int64SliceBuilder) String() string            { return nameOf(_this) }

func (_this *int64SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *int64SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeInt64:
		elemCount := len(value) / 8
		slice := make([]int64, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			slice[i] = int64(value[i*8]) |
				(int64(value[i*8+1]) << 8) |
				(int64(value[i*8+2]) << 16) |
				(int64(value[i*8+3]) << 24) |
				(int64(value[i*8+4]) << 32) |
				(int64(value[i*8+5]) << 40) |
				(int64(value[i*8+6]) << 48) |
				(int64(value[i*8+7]) << 56)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *int64SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToInt64SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type float32ArrayBuilder struct {
	dstType reflect.Type
}

func newFloat32ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &float32ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *float32ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *float32ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeFloat32:
		elemCount := len(value) / 4
		for i := 0; i < elemCount; i++ {
			elemValue := uint32(value[i*4]) |
				(uint32(value[i*4+1]) << 8) |
				(uint32(value[i*4+2]) << 16) |
				(uint32(value[i*4+3]) << 24)
			elem := dst.Index(i)
			elem.SetFloat(float64(math.Float32frombits(elemValue)))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *float32ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type float32SliceBuilder struct{}

var globalFloat32SliceBuilder = &float32SliceBuilder{}

func generateFloat32SliceBuilder(ctx *Context) ObjectBuilder { return globalFloat32SliceBuilder }
func (_this *float32SliceBuilder) String() string            { return nameOf(_this) }

func (_this *float32SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *float32SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeFloat32:
		elemCount := len(value) / 4
		slice := make([]float32, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			elemValue := uint32(value[i*4]) |
				(uint32(value[i*4+1]) << 8) |
				(uint32(value[i*4+2]) << 16) |
				(uint32(value[i*4+3]) << 24)
			slice[i] = math.Float32frombits(elemValue)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *float32SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToFloat32SliceGenerator(ctx).BuildBeginListContents(ctx)
}

// ============================================================================

type float64ArrayBuilder struct {
	dstType reflect.Type
}

func newFloat64ArrayBuilderGenerator(dstType reflect.Type) BuilderGenerator {
	return func(ctx *Context) ObjectBuilder {
		return &float64ArrayBuilder{
			dstType: dstType,
		}
	}
}

func (_this *float64ArrayBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *float64ArrayBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeFloat64:
		elemCount := len(value) / 8
		for i := 0; i < elemCount; i++ {
			elemValue := uint64(value[i*8]) |
				(uint64(value[i*8+1]) << 8) |
				(uint64(value[i*8+2]) << 16) |
				(uint64(value[i*8+3]) << 24) |
				(uint64(value[i*8+4]) << 32) |
				(uint64(value[i*8+5]) << 40) |
				(uint64(value[i*8+6]) << 48) |
				(uint64(value[i*8+7]) << 56)
			elem := dst.Index(i)
			elem.SetFloat(math.Float64frombits(elemValue))
		}
	default:
		PanicBadEvent(_this, "BuildFromArray(%v)", arrayType)
	}
	return dst
}

func (_this *float64ArrayBuilder) BuildBeginListContents(ctx *Context) {
	generator := newArrayBuilderGenerator(ctx.GetBuilderGeneratorForType, _this.dstType)
	generator(ctx).BuildBeginListContents(ctx)
}

type float64SliceBuilder struct{}

var globalFloat64SliceBuilder = &float64SliceBuilder{}

func generateFloat64SliceBuilder(ctx *Context) ObjectBuilder { return globalFloat64SliceBuilder }
func (_this *float64SliceBuilder) String() string            { return nameOf(_this) }

func (_this *float64SliceBuilder) BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(dst.Type()))
	ctx.NANext()
	return dst
}

func (_this *float64SliceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeFloat64:
		elemCount := len(value) / 8
		slice := make([]float64, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			elemValue := uint64(value[i*8]) |
				(uint64(value[i*8+1]) << 8) |
				(uint64(value[i*8+2]) << 16) |
				(uint64(value[i*8+3]) << 24) |
				(uint64(value[i*8+4]) << 32) |
				(uint64(value[i*8+5]) << 40) |
				(uint64(value[i*8+6]) << 48) |
				(uint64(value[i*8+7]) << 56)
			slice[i] = math.Float64frombits(elemValue)
		}
		dst.Set(reflect.ValueOf(slice))
	default:
		PanicBadEvent(_this, "BuildFromSlice(%v)", arrayType)
	}
	return dst
}

func (_this *float64SliceBuilder) BuildBeginListContents(ctx *Context) {
	listToFloat64SliceGenerator(ctx).BuildBeginListContents(ctx)
}
