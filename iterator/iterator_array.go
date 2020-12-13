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

package iterator

import (
	"math"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/events"
)

// The matching generated code is in iterator_generated.go

// -----
// uint8
// -----

func (_this *uint8ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.CanAddr() {
		bytes := v.Slice(0, v.Len()).Bytes()
		eventReceiver.OnArray(events.ArrayTypeUint8, uint64(len(bytes)), bytes)
	} else {
		byteCount := v.Len()
		bytes := make([]byte, byteCount)
		for i := 0; i < byteCount; i++ {
			bytes[i] = v.Index(i).Interface().(uint8)
		}
		eventReceiver.OnArray(events.ArrayTypeUint8, uint64(byteCount), bytes)
	}
}

func (_this *uint8SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnArray(events.ArrayTypeUint8, uint64(v.Len()), v.Bytes())
}

// ------
// uint16
// ------

func iterateArrayUint16(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*2, elementCount*2)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*2] = uint8(elem)
		data[i*2+1] = uint8(elem >> 8)
	}
	eventReceiver.OnArray(events.ArrayTypeUint16, uint64(elementCount), data)
}

func (_this *uint16ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint16(v, eventReceiver)
}

func (_this *uint16SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint16(v, eventReceiver)
}

// ------
// uint32
// ------

func iterateArrayUint32(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	eventReceiver.OnArray(events.ArrayTypeUint32, uint64(elementCount), data)
}

func (_this *uint32ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint32(v, eventReceiver)
}

func (_this *uint32SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint32(v, eventReceiver)
}

// ------
// uint64
// ------

func iterateArrayUint64(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8, elementCount*8)
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
	eventReceiver.OnArray(events.ArrayTypeUint64, uint64(elementCount), data)
}

func (_this *uint64ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint64(v, eventReceiver)
}

func (_this *uint64SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint64(v, eventReceiver)
}

// ----
// int8
// ----

func iterateArrayInt8(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount, elementCount)
	for i := 0; i < elementCount; i++ {
		data[i] = uint8(v.Index(i).Int())
	}
	eventReceiver.OnArray(events.ArrayTypeInt8, uint64(elementCount), data)
}

func (_this *int8ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt8(v, eventReceiver)
}

func (_this *int8SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt8(v, eventReceiver)
}

// -----
// int16
// -----

func iterateArrayInt16(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*2, elementCount*2)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Int()
		data[i*2] = uint8(elem)
		data[i*2+1] = uint8(elem >> 8)
	}
	eventReceiver.OnArray(events.ArrayTypeInt16, uint64(elementCount), data)
}

func (_this *int16ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt16(v, eventReceiver)
}

func (_this *int16SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt16(v, eventReceiver)
}

// -----
// int32
// -----

func iterateArrayInt32(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Int()
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	eventReceiver.OnArray(events.ArrayTypeInt32, uint64(elementCount), data)
}

func (_this *int32ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt32(v, eventReceiver)
}

func (_this *int32SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt32(v, eventReceiver)
}

// -----
// int64
// -----

func iterateArrayInt64(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8, elementCount*8)
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
	eventReceiver.OnArray(events.ArrayTypeInt64, uint64(elementCount), data)
}

func (_this *int64ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt64(v, eventReceiver)
}

func (_this *int64SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayInt64(v, eventReceiver)
}

// -------
// float32
// -------

func iterateArrayFloat32(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*4, elementCount*4)
	for i := 0; i < elementCount; i++ {
		elem := math.Float32bits(float32(v.Index(i).Float()))
		data[i*4] = uint8(elem)
		data[i*4+1] = uint8(elem >> 8)
		data[i*4+2] = uint8(elem >> 16)
		data[i*4+3] = uint8(elem >> 24)
	}
	eventReceiver.OnArray(events.ArrayTypeFloat32, uint64(elementCount), data)
}

func (_this *float32ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayFloat32(v, eventReceiver)
}

func (_this *float32SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayFloat32(v, eventReceiver)
}

// -------
// float64
// -------

func iterateArrayFloat64(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*8, elementCount*8)
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
	eventReceiver.OnArray(events.ArrayTypeFloat64, uint64(elementCount), data)
}

func (_this *float64ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayFloat64(v, eventReceiver)
}

func (_this *float64SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayFloat64(v, eventReceiver)
}

// ----
// bool
// ----

func iterateArrayBool(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	byteCount := common.ElementCountToByteCount(1, uint64(elementCount))
	data := make([]uint8, byteCount, byteCount)
	if elementCount == 0 {
		eventReceiver.OnArray(events.ArrayTypeBoolean, uint64(elementCount), data)
		return
	}

	nextData := data
	var nextByte uint8
	if v.Index(0).Bool() {
		nextByte = 1
	}

	for i := 1; i < elementCount; i++ {
		if i&7 == 0 {
			nextData[0] = nextByte
			nextData = nextData[1:]
		}
		nextByte <<= 1
		if v.Index(i).Bool() {
			nextByte |= 1
		}
	}
	if elementCount&7 != 0 {
		nextData[0] = nextByte
	}
	eventReceiver.OnArray(events.ArrayTypeBoolean, uint64(elementCount), data)
}

func (_this *boolArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayBool(v, eventReceiver)
}

func (_this *boolSliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayBool(v, eventReceiver)
}
