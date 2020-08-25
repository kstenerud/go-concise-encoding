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
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
)

// The matching generated code is in iterator_generated.go

// -----
// uint8
// -----

func (_this *uint8ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	if v.CanAddr() {
		bytes := v.Slice(0, v.Len()).Bytes()
		eventReceiver.OnTypedArray(events.ArrayTypeUint8, uint64(len(bytes)), bytes)
	} else {
		byteCount := v.Len()
		bytes := make([]byte, byteCount)
		for i := 0; i < byteCount; i++ {
			bytes[i] = v.Index(i).Interface().(uint8)
		}
		eventReceiver.OnTypedArray(events.ArrayTypeUint8, uint64(byteCount), bytes)
	}
}

func (_this *uint8SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	eventReceiver.OnTypedArray(events.ArrayTypeUint8, uint64(v.Len()), v.Bytes())
}

// -----
// uint16
// -----

func iterateArrayUint16(v reflect.Value, eventReceiver events.DataEventReceiver) {
	elementCount := v.Len()
	data := make([]uint8, elementCount*2, elementCount*2)
	for i := 0; i < elementCount; i++ {
		elem := v.Index(i).Uint()
		data[i*2] = uint8(elem)
		data[i*2+1] = uint8(elem >> 8)
	}
	eventReceiver.OnTypedArray(events.ArrayTypeUint16, uint64(elementCount), data)
}

func (_this *uint16ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint16(v, eventReceiver)
}

func (_this *uint16SliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, _ AddReference) {
	iterateArrayUint16(v, eventReceiver)
}
