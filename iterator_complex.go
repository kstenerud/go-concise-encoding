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

package concise_encoding

import (
	"reflect"
)

// ---------
// Interface
// ---------

type interfaceIterator struct {
}

func newInterfaceIterator(srcType reflect.Type) ObjectIterator {
	return &interfaceIterator{}
}

func (_this *interfaceIterator) PostCacheInitIterator() {
}

func (_this *interfaceIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
	} else {
		elem := v.Elem()
		iter := getIteratorForType(elem.Type())
		iter.Iterate(elem, root)
	}
}

// -------
// Pointer
// -------

type pointerIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
}

func newPointerIterator(srcType reflect.Type) ObjectIterator {
	return &pointerIterator{srcType: srcType}
}

func (_this *pointerIterator) PostCacheInitIterator() {
	_this.elemIter = getIteratorForType(_this.srcType.Elem())
}

func (_this *pointerIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
		return
	}
	if root.addReference(v) {
		return
	}
	_this.elemIter.Iterate(v.Elem(), root)
}

// -----------
// uint8 array
// -----------

type uint8ArrayIterator struct {
}

func newUInt8ArrayIterator() ObjectIterator {
	return &uint8ArrayIterator{}
}

func (_this *uint8ArrayIterator) PostCacheInitIterator() {
}

func (_this *uint8ArrayIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.CanAddr() {
		root.eventReceiver.OnBytes(v.Slice(0, v.Len()).Bytes())
	} else {
		tempSlice := make([]byte, v.Len())
		tempLen := v.Len()
		for i := 0; i < tempLen; i++ {
			tempSlice[i] = v.Index(i).Interface().(uint8)
		}
		root.eventReceiver.OnBytes(tempSlice)
	}
}

// -------
// Complex
// -------

type complexIterator struct {
}

func newComplexIterator() ObjectIterator {
	return &complexIterator{}
}

func (_this *complexIterator) PostCacheInitIterator() {
}

func (_this *complexIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnComplex(v.Complex())
}

// -----
// Slice
// -----

type sliceIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
}

func newSliceIterator(srcType reflect.Type) ObjectIterator {
	return &sliceIterator{
		srcType: srcType,
	}
}

func (_this *sliceIterator) PostCacheInitIterator() {
	_this.elemIter = getIteratorForType(_this.srcType.Elem())
}

func (_this *sliceIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
		return
	}
	if root.addReference(v) {
		return
	}

	root.eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.Iterate(v.Index(i), root)
	}
	root.eventReceiver.OnEnd()
}

// -----
// Array
// -----

type arrayIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
}

func newArrayIterator(srcType reflect.Type) ObjectIterator {
	return &arrayIterator{
		srcType: srcType,
	}
}

func (_this *arrayIterator) PostCacheInitIterator() {
	_this.elemIter = getIteratorForType(_this.srcType.Elem())
}

func (_this *arrayIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.Iterate(v.Index(i), root)
	}
	root.eventReceiver.OnEnd()
}

// ---
// Map
// ---

type mapIterator struct {
	srcType   reflect.Type
	keyIter   ObjectIterator
	valueIter ObjectIterator
}

func newMapIterator(srcType reflect.Type) ObjectIterator {
	return &mapIterator{
		srcType: srcType,
	}
}

func (_this *mapIterator) PostCacheInitIterator() {
	_this.keyIter = getIteratorForType(_this.srcType.Key())
	_this.valueIter = getIteratorForType(_this.srcType.Elem())
}

func (_this *mapIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	if v.IsNil() {
		root.eventReceiver.OnNil()
		return
	}
	if root.addReference(v) {
		return
	}

	root.eventReceiver.OnMap()

	iter := mapRange(v)
	for iter.Next() {
		_this.keyIter.Iterate(iter.Key(), root)
		_this.valueIter.Iterate(iter.Value(), root)
	}

	root.eventReceiver.OnEnd()
}

// ------
// Struct
// ------

type structIteratorField struct {
	Name     string
	Index    int
	Iterator ObjectIterator
}

func newStructIteratorField(name string, index int, iterator ObjectIterator) *structIteratorField {
	return &structIteratorField{
		Name:     name,
		Index:    index,
		Iterator: iterator,
	}
}

type structIterator struct {
	srcType        reflect.Type
	fieldIterators []*structIteratorField
}

func newStructIterator(srcType reflect.Type) ObjectIterator {
	return &structIterator{
		srcType: srcType,
	}
}

func (_this *structIterator) PostCacheInitIterator() {
	for i := 0; i < _this.srcType.NumField(); i++ {
		field := _this.srcType.Field(i)
		if isFieldExported(field.Name) {
			iterator := &structIteratorField{
				Name:     field.Name,
				Index:    i,
				Iterator: getIteratorForType(field.Type),
			}
			_this.fieldIterators = append(_this.fieldIterators, iterator)
		}
	}
}

func (_this *structIterator) Iterate(v reflect.Value, root *RootObjectIterator) {
	root.eventReceiver.OnMap()

	for _, iter := range _this.fieldIterators {
		root.eventReceiver.OnString(iter.Name)
		iter.Iterator.Iterate(v.Field(iter.Index), root)
	}

	root.eventReceiver.OnEnd()
}
