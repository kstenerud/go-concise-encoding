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
	"reflect"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/internal/arrays"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"
)

type interfaceBuilder struct{}

var globalInterfaceBuilder = &interfaceBuilder{}

func generateInterfaceBuilder(ctx *Context) Builder { return globalInterfaceBuilder }
func (_this *interfaceBuilder) String() string      { return reflect.TypeOf(_this).String() }

func (_this *interfaceBuilder) BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value {
	dst.Set(reflect.Zero(common.TypeInterface))
	return dst
}

func (_this *interfaceBuilder) BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value {
	dst.Set(reflect.ValueOf(value))
	return dst
}

func (_this *interfaceBuilder) BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value {
	setFromUID(value, dst)
	return dst
}

func (_this *interfaceBuilder) BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeUint8:
		dst.Set(reflect.ValueOf(common.CloneBytes(value)))
	case events.ArrayTypeUint16:
		dst.Set(reflect.ValueOf(arrays.BytesToUint16Slice(value)))
	case events.ArrayTypeUint32:
		dst.Set(reflect.ValueOf(arrays.BytesToUint32Slice(value)))
	case events.ArrayTypeUint64:
		dst.Set(reflect.ValueOf(arrays.BytesToUint64Slice(value)))
	case events.ArrayTypeInt8:
		dst.Set(reflect.ValueOf(arrays.BytesToInt8Slice(value)))
	case events.ArrayTypeInt16:
		dst.Set(reflect.ValueOf(arrays.BytesToInt16Slice(value)))
	case events.ArrayTypeInt32:
		dst.Set(reflect.ValueOf(arrays.BytesToInt32Slice(value)))
	case events.ArrayTypeInt64:
		dst.Set(reflect.ValueOf(arrays.BytesToInt64Slice(value)))
	case events.ArrayTypeFloat16:
		dst.Set(reflect.ValueOf(arrays.BytesToFloat16Slice(value)))
	case events.ArrayTypeFloat32:
		dst.Set(reflect.ValueOf(arrays.BytesToFloat32Slice(value)))
	case events.ArrayTypeFloat64:
		dst.Set(reflect.ValueOf(arrays.BytesToFloat64Slice(value)))
	case events.ArrayTypeString:
		dst.Set(reflect.ValueOf(string(value)))
	case events.ArrayTypeResourceID:
		setPRIDFromString(string(value), dst)
	default:
		panic(fmt.Errorf("TODO: Typed array support for %v", arrayType))
	}
	return dst
}

func (_this *interfaceBuilder) BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value {
	switch arrayType {
	case events.ArrayTypeString:
		dst.Set(reflect.ValueOf(value))
	case events.ArrayTypeResourceID:
		setPRIDFromString(value, dst)
	case events.ArrayTypeReferenceRemote:
		// TODO: What to do about remote ref??
		setPRIDFromString(value, dst)
	default:
		panic(fmt.Errorf("BUG: Array type %v is not stringlike", arrayType))
	}
	return dst
}

func (_this *interfaceBuilder) BuildFromCustomBinary(ctx *Context, customType uint64, data []byte, dst reflect.Value) reflect.Value {
	ctx.TryBuildFromCustomBinary(_this, customType, data, dst)

	return dst
}

func (_this *interfaceBuilder) BuildFromCustomText(ctx *Context, customType uint64, data string, dst reflect.Value) reflect.Value {
	ctx.TryBuildFromCustomText(_this, customType, data, dst)
	return dst
}

func (_this *interfaceBuilder) BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value {
	v := types.Media{
		MediaType: mediaType,
		Data:      common.CloneBytes(data),
	}
	dst.Set(reflect.ValueOf(v))
	return dst
}

func (_this *interfaceBuilder) BuildFromTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value {
	if gTime, err := value.AsGoTime(); err == nil {
		dst.Set(reflect.ValueOf(gTime))
	} else {
		dst.Set(reflect.ValueOf(value))
	}
	return dst
}

func (_this *interfaceBuilder) BuildNewList(ctx *Context) {
	interfaceSliceBuilderGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *interfaceBuilder) BuildNewMap(ctx *Context) {
	interfaceMapBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *interfaceBuilder) BuildNewEdge(ctx *Context) {
	interfaceEdgeBuilderGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *interfaceBuilder) BuildNewNode(ctx *Context) {
	interfaceNodeBuilderGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *interfaceBuilder) BuildBeginListContents(ctx *Context) {
	interfaceSliceBuilderGenerator(ctx).BuildBeginListContents(ctx)
}

func (_this *interfaceBuilder) BuildBeginMapContents(ctx *Context) {
	interfaceMapBuilderGenerator(ctx).BuildBeginMapContents(ctx)
}

func (_this *interfaceBuilder) BuildBeginEdgeContents(ctx *Context) {
	interfaceEdgeBuilderGenerator(ctx).BuildBeginEdgeContents(ctx)
}

func (_this *interfaceBuilder) BuildBeginNodeContents(ctx *Context) {
	interfaceNodeBuilderGenerator(ctx).BuildBeginNodeContents(ctx)
}

func (_this *interfaceBuilder) BuildBeginMarker(ctx *Context, id []byte) {
	panic("TODO: interfaceBuilder.BuildBeginMarker")
}

func (_this *interfaceBuilder) BuildFromLocalReference(ctx *Context, id []byte) {
	panic("TODO: interfaceBuilder.BuildFromLocalReference")
}

func (_this *interfaceBuilder) NotifyChildContainerFinished(ctx *Context, value reflect.Value) {
	ctx.UnstackBuilderAndNotifyChildFinished(value)
}

func (_this *interfaceBuilder) BuildArtificiallyEndContainer(ctx *Context) {
}
