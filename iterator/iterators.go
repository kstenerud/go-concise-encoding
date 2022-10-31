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

// Iterators iterate through go objects, producing data events.
package iterator

import (
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"sort"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/types"
)

func iterateTime(context *Context, v reflect.Value) {
	context.EventReceiver.OnTime(compact_time.AsCompactTime(v.Interface().(time.Time)))
}

func iteratePTime(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		t := v.Interface().(*time.Time)
		context.EventReceiver.OnTime(compact_time.AsCompactTime(*t))
	}
}

func iterateCompactTime(context *Context, v reflect.Value) {
	ct := v.Interface().(compact_time.Time)
	context.EventReceiver.OnTime(ct)
}

func iteratePCompactTime(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		t := v.Interface().(*compact_time.Time)
		context.EventReceiver.OnTime(*t)
	}
}

func iterateURL(context *Context, v reflect.Value) {
	vCopy := v.Interface().(url.URL)
	context.EventReceiver.OnStringlikeArray(events.ArrayTypeResourceID, (&vCopy).String())
}

func iteratePURL(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		str := v.Interface().(*url.URL).String()
		context.EventReceiver.OnStringlikeArray(events.ArrayTypeResourceID, str)
	}
}

func iterateBigInt(context *Context, v reflect.Value) {
	vCopy := v.Interface().(big.Int)
	context.EventReceiver.OnBigInt(&vCopy)
}

func iteratePBigInt(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		context.EventReceiver.OnBigInt(v.Interface().(*big.Int))
	}
}

func iterateBigFloat(context *Context, v reflect.Value) {
	vCopy := v.Interface().(big.Float)
	context.EventReceiver.OnBigFloat(&vCopy)
}

func iteratePBigFloat(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		context.EventReceiver.OnBigFloat(v.Interface().(*big.Float))
	}
}

func iterateBigDecimal(context *Context, v reflect.Value) {
	vCopy := v.Interface().(apd.Decimal)
	context.EventReceiver.OnBigDecimalFloat(&vCopy)
}

func iteratePBigDecimal(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		context.EventReceiver.OnBigDecimalFloat(v.Interface().(*apd.Decimal))
	}
}

func iterateDecimalFloat(context *Context, v reflect.Value) {
	context.EventReceiver.OnDecimalFloat(v.Interface().(compact_float.DFloat))
}

func iterateBool(context *Context, v reflect.Value) {
	context.EventReceiver.OnBoolean(v.Bool())
}

func iterateInt(context *Context, v reflect.Value) {
	context.EventReceiver.OnInt(v.Int())
}

func iterateUint(context *Context, v reflect.Value) {
	context.EventReceiver.OnPositiveInt(v.Uint())
}

func iterateFloat(context *Context, v reflect.Value) {
	context.EventReceiver.OnFloat(v.Float())
}

func iterateString(context *Context, v reflect.Value) {
	context.EventReceiver.OnStringlikeArray(events.ArrayTypeString, v.String())
}

func iterateUID(context *Context, v reflect.Value) {
	vCopy := v.Interface().(types.UID)
	context.EventReceiver.OnUID(vCopy[:])
}

func iterateMedia(context *Context, v reflect.Value) {
	vCopy := v.Interface().(types.Media)
	if len(vCopy.MediaType) == 0 {
		panic(fmt.Errorf("media cannot have an empty media type"))
	}

	context.EventReceiver.OnMedia(vCopy.MediaType, vCopy.Data)
}

func iterateNode(context *Context, v reflect.Value) {
	context.EventReceiver.OnNode()
	iterateInterface(context, v.Field(types.NodeFieldIndexValue))
	children := v.Field(types.NodeFieldIndexChildren)
	for i := 0; i < children.Len(); i++ {
		iterateInterface(context, children.Index(i))
	}
	context.EventReceiver.OnEndContainer()
}

func iterateEdge(context *Context, v reflect.Value) {
	context.EventReceiver.OnEdge()
	iterateInterface(context, v.Field(types.EdgeFieldIndexSource))
	iterateInterface(context, v.Field(types.EdgeFieldIndexDescription))
	iterateInterface(context, v.Field(types.EdgeFieldIndexDestination))
}

func iterateInterface(context *Context, v reflect.Value) {
	if common.IsNil(v) {
		context.NotifyNil()
	} else {
		elem := v.Elem()
		iterate := context.GetIteratorForType(elem.Type())
		iterate(context, elem)
	}
}

func newPointerIterator(ctx *Context, pointerType reflect.Type) IteratorFunction {
	iterate := ctx.GetIteratorForType(pointerType.Elem())

	return func(context *Context, v reflect.Value) {
		if common.IsNil(v) {
			context.NotifyNil()
			return
		}
		if context.TryAddLocalReference(v) {
			return
		}
		iterate(context, v.Elem())
	}
}

func newSliceOrArrayAsListIterator(ctx *Context, sliceType reflect.Type) IteratorFunction {
	iterate := ctx.GetIteratorForType(sliceType.Elem())

	return func(context *Context, v reflect.Value) {
		if common.IsNil(v) {
			context.NotifyNil()
			return
		}
		if context.TryAddLocalReference(v) {
			return
		}

		context.EventReceiver.OnList()
		length := v.Len()
		for i := 0; i < length; i++ {
			iterate(context, v.Index(i))
		}
		context.EventReceiver.OnEndContainer()
	}
}

func newMapIterator(ctx *Context, mapType reflect.Type) IteratorFunction {
	iterateKey := ctx.GetIteratorForType(mapType.Key())
	iterateValue := ctx.GetIteratorForType(mapType.Elem())

	return func(context *Context, v reflect.Value) {
		if common.IsNil(v) {
			context.NotifyNil()
			return
		}
		if context.TryAddLocalReference(v) {
			return
		}

		context.EventReceiver.OnMap()
		iter := common.MapRange(v)
		for iter.Next() {
			iterateKey(context, iter.Key())
			iterateValue(context, iter.Value())
		}
		context.EventReceiver.OnEndContainer()
	}
}

func newCustomBinaryIterator(convert configuration.ConvertToCustomFunction) IteratorFunction {
	return func(context *Context, v reflect.Value) {
		customType, asBytes, err := convert(v)
		if err != nil {
			panic(fmt.Errorf("error converting type %v to custom bytes: %v", v.Type(), err))
		}
		context.EventReceiver.OnCustomBinary(customType, asBytes)
	}
}

func newCustomTextIterator(convert configuration.ConvertToCustomFunction) IteratorFunction {
	return func(context *Context, v reflect.Value) {
		customType, asBytes, err := convert(v)
		if err != nil {
			panic(fmt.Errorf("error converting type %v to custom text: %v", v.Type(), err))
		}
		context.EventReceiver.OnCustomText(customType, string(asBytes))
	}
}

type structField struct {
	Identifier   string
	Name         string
	Type         reflect.Type
	Index        int
	Iterate      IteratorFunction
	IncludeName  bool
	OmitBehavior configuration.FieldOmitBehavior
	Order        int64
}

func (_this *structField) getValueFromStruct(fromValue reflect.Value) reflect.Value {
	return fromValue.Field(_this.Index)
}

func newStructField(fromField reflect.StructField, index int, fieldNameStyle configuration.FieldNameStyle) structField {
	tags := common.DecodeGoTags(fromField)
	name := tags.Name
	if fieldNameStyle == configuration.FieldNameSnakeCase {
		name = common.CamelCaseToSnakeCase(name)
	}

	return structField{
		Identifier:   common.ToStructFieldIdentifier(tags.Name),
		Name:         name,
		Type:         fromField.Type,
		Index:        index,
		OmitBehavior: tags.OmitBehavior,
		Order:        tags.Order,
	}
}

func isValueEmpty(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	switch value.Kind() {
	case reflect.Interface, reflect.Pointer:
		return value.IsNil()
	case reflect.Map, reflect.Slice:
		return value.IsNil() || value.Len() == 0
	case reflect.Array, reflect.String:
		return value.Len() == 0
	default:
		return false
	}
}

func isValueZero(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	if value.IsZero() {
		return true
	}

	return isValueEmpty(value)
}

func extractFields(ctx *Context, structType reflect.Type, fields []structField) []structField {
	for i := 0; i < structType.NumField(); i++ {
		reflectField := structType.Field(i)
		if common.IsFieldExported(reflectField.Name) {
			field := newStructField(reflectField, i, ctx.Configuration.FieldNameStyle)

			if field.OmitBehavior != configuration.OmitFieldAlways {
				if reflectField.Anonymous {
					field.Iterate = getEmbeddedFieldIterator(ctx, field.Type)
				} else {
					field.IncludeName = true
					field.Iterate = ctx.GetIteratorForType(field.Type)
				}
				fields = append(fields, field)
			}
		}
	}

	sort.SliceStable(fields, func(i, j int) bool {
		return fields[i].Order < fields[j].Order
	})

	return fields
}

func shouldIncludeField(field structField, value reflect.Value, defaultOmitBehavior configuration.FieldOmitBehavior) bool {
	omitBehavior := field.OmitBehavior
	if omitBehavior == configuration.OmitFieldChooseDefault {
		omitBehavior = defaultOmitBehavior
	}
	switch omitBehavior {
	case configuration.OmitFieldAlways:
		return false
	case configuration.OmitFieldNever:
		return true
	case configuration.OmitFieldEmpty:
		return !isValueEmpty(value)
	case configuration.OmitFieldZero:
		return !isValueZero(value)
	}
	// Should never happen
	return true
}

func newStructIterator(ctx *Context, structType reflect.Type) IteratorFunction {
	fields := extractFields(ctx, structType, make([]structField, 0, structType.NumField()))
	return func(context *Context, value reflect.Value) {
		context.EventReceiver.OnMap()
		for _, field := range fields {
			fieldValue := field.getValueFromStruct(value)
			if shouldIncludeField(field, fieldValue, ctx.Configuration.DefaultFieldOmitBehavior) {
				if field.IncludeName {
					context.EventReceiver.OnStringlikeArray(events.ArrayTypeString, field.Name)
				}
				field.Iterate(context, fieldValue)
			}
		}
		context.EventReceiver.OnEndContainer()
	}
}

func getEmbeddedFieldIterator(ctx *Context, structType reflect.Type) IteratorFunction {
	fields := extractFields(ctx, structType, make([]structField, 0, structType.NumField()))
	return func(context *Context, value reflect.Value) {
		for _, field := range fields {
			fieldValue := field.getValueFromStruct(value)
			if shouldIncludeField(field, fieldValue, ctx.Configuration.DefaultFieldOmitBehavior) {
				if field.IncludeName {
					context.EventReceiver.OnStringlikeArray(events.ArrayTypeString, field.Name)
				}
				field.Iterate(context, fieldValue)
			}
		}
	}
}

func iterateSliceUint8(context *Context, v reflect.Value) {
	context.EventReceiver.OnArray(events.ArrayTypeUint8, uint64(v.Len()), v.Bytes())
}

func iterateArrayUint8(context *Context, v reflect.Value) {
	if v.CanAddr() {
		bytes := v.Slice(0, v.Len()).Bytes()
		context.EventReceiver.OnArray(events.ArrayTypeUint8, uint64(len(bytes)), bytes)
	} else {
		byteCount := v.Len()
		bytes := make([]byte, byteCount)
		for i := 0; i < byteCount; i++ {
			bytes[i] = v.Index(i).Interface().(uint8)
		}
		context.EventReceiver.OnArray(events.ArrayTypeUint8, uint64(byteCount), bytes)
	}
}

func iterateSliceOrArrayUint16(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*2)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*2] = uint8(elem)
		data[i*2+1] = uint8(elem >> 8)
	}
	context.EventReceiver.OnArray(events.ArrayTypeUint16, uint64(elementCount), data)
}

func iterateSliceOrArrayUint32(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	context.EventReceiver.OnArray(events.ArrayTypeUint32, uint64(elementCount), data)
}

func iterateSliceOrArrayUint64(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*8] = uint8(elem)
		data[i*8+1] = uint8(elem >> 8)
		data[i*8+2] = uint8(elem >> 16)
		data[i*8+3] = uint8(elem >> 24)
		data[i*8+4] = uint8(elem >> 32)
		data[i*8+5] = uint8(elem >> 40)
		data[i*8+6] = uint8(elem >> 48)
		data[i*8+7] = uint8(elem >> 56)
	}
	context.EventReceiver.OnArray(events.ArrayTypeUint64, uint64(elementCount), data)
}

func iterateSliceOrArrayUint(context *Context, v reflect.Value) {
	if common.Is64BitArch() {
		iterateSliceOrArrayUint64(context, v)
	} else {
		iterateSliceOrArrayUint32(context, v)
	}
}

func iterateSliceOrArrayInt8(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount)
	for i := 0; i < elementCount; i++ {
		data[i] = uint8(v.Index(i).Int())
	}
	context.EventReceiver.OnArray(events.ArrayTypeInt8, uint64(elementCount), data)
}

func iterateSliceOrArrayInt16(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*2)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Int()
		data[i*2] = uint8(elem)
		data[i*2+1] = uint8(elem >> 8)
	}
	context.EventReceiver.OnArray(events.ArrayTypeInt16, uint64(elementCount), data)
}

func iterateSliceOrArrayInt32(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Int()
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	context.EventReceiver.OnArray(events.ArrayTypeInt32, uint64(elementCount), data)
}

func iterateSliceOrArrayInt64(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Int()
		data[i*8] = uint8(elem)
		data[i*8+1] = uint8(elem >> 8)
		data[i*8+2] = uint8(elem >> 16)
		data[i*8+3] = uint8(elem >> 24)
		data[i*8+4] = uint8(elem >> 32)
		data[i*8+5] = uint8(elem >> 40)
		data[i*8+6] = uint8(elem >> 48)
		data[i*8+7] = uint8(elem >> 56)
	}
	context.EventReceiver.OnArray(events.ArrayTypeInt64, uint64(elementCount), data)
}

func iterateSliceOrArrayInt(context *Context, v reflect.Value) {
	if common.Is64BitArch() {
		iterateSliceOrArrayInt64(context, v)
	} else {
		iterateSliceOrArrayInt32(context, v)
	}
}

func iterateSliceOrArrayFloat32(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := math.Float32bits(float32(v.Index(i).Float()))
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	context.EventReceiver.OnArray(events.ArrayTypeFloat32, uint64(elementCount), data)
}

func iterateSliceOrArrayFloat64(context *Context, v reflect.Value) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8)
	for i := 0; i < elementCount; i++ {
		elem := math.Float64bits(v.Index(i).Float())
		data[i*8] = uint8(elem)
		data[i*8+1] = uint8(elem >> 8)
		data[i*8+2] = uint8(elem >> 16)
		data[i*8+3] = uint8(elem >> 24)
		data[i*8+4] = uint8(elem >> 32)
		data[i*8+5] = uint8(elem >> 40)
		data[i*8+6] = uint8(elem >> 48)
		data[i*8+7] = uint8(elem >> 56)
	}
	context.EventReceiver.OnArray(events.ArrayTypeFloat64, uint64(elementCount), data)
}

func iterateSliceOrArrayBool(context *Context, v reflect.Value) {
	elementCount := v.Len()
	byteCount := common.ElementCountToByteCount(1, uint64(elementCount))
	data := make([]uint8, byteCount)
	if elementCount == 0 {
		context.EventReceiver.OnArray(events.ArrayTypeBit, uint64(elementCount), data)
		return
	}

	iDst := 0
	for iSrc := 0; iSrc < elementCount; {
		bitCount := 8
		if elementCount-iSrc < 8 {
			bitCount = elementCount - iSrc
		}
		accum := byte(0)
		for iBit := 0; iBit < bitCount; iBit++ {
			if v.Index(iBit).Bool() {
				accum |= 1 << iBit
			}
			iSrc++
		}
		data[iDst] = accum
		iDst++
	}

	context.EventReceiver.OnArray(events.ArrayTypeBit, uint64(elementCount), data)
}
