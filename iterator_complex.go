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
	root *RootObjectIterator
}

func newInterfaceIterator(srcType reflect.Type) ObjectIterator {
	return &interfaceIterator{}
}

func (this *interfaceIterator) PostCacheInitIterator() {
}

func (this *interfaceIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &interfaceIterator{
		root: root,
	}
}

func (this *interfaceIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
	} else {
		elem := v.Elem()
		iter := getIteratorForType(elem.Type()).CloneFromTemplate(this.root)
		iter.Iterate(elem)
	}
}

// -------
// Pointer
// -------

type pointerIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
	root     *RootObjectIterator
}

func newPointerIterator(srcType reflect.Type) ObjectIterator {
	return &pointerIterator{srcType: srcType}
}

func (this *pointerIterator) PostCacheInitIterator() {
	this.elemIter = getIteratorForType(this.srcType.Elem())
}

func (this *pointerIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &pointerIterator{
		srcType:  this.srcType,
		root:     root,
		elemIter: this.elemIter.CloneFromTemplate(root),
	}
}

func (this *pointerIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
		return
	}
	if this.root.addReference(v) {
		return
	}
	this.elemIter.Iterate(v.Elem())
}

// -----------
// uint8 array
// -----------

type uint8ArrayIterator struct {
	root *RootObjectIterator
}

func newUInt8ArrayIterator() ObjectIterator {
	return &uint8ArrayIterator{}
}

func (this *uint8ArrayIterator) PostCacheInitIterator() {
}

func (this *uint8ArrayIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &uint8ArrayIterator{root: root}
}

func (this *uint8ArrayIterator) Iterate(v reflect.Value) {
	if v.CanAddr() {
		this.root.eventReceiver.OnBytes(v.Slice(0, v.Len()).Bytes())
	} else {
		tempSlice := make([]byte, v.Len())
		tempLen := v.Len()
		for i := 0; i < tempLen; i++ {
			tempSlice[i] = v.Index(i).Interface().(uint8)
		}
		this.root.eventReceiver.OnBytes(tempSlice)
	}
}

// -------
// Complex
// -------

type complexIterator struct {
	root *RootObjectIterator
}

func newComplexIterator() ObjectIterator {
	return &complexIterator{}
}

func (this *complexIterator) PostCacheInitIterator() {
}

func (this *complexIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &complexIterator{root: root}
}

func (this *complexIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnComplex(v.Complex())
}

// -----
// Slice
// -----

type sliceIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
	root     *RootObjectIterator
}

func newSliceIterator(srcType reflect.Type) ObjectIterator {
	return &sliceIterator{
		srcType: srcType,
	}
}

func (this *sliceIterator) PostCacheInitIterator() {
	this.elemIter = getIteratorForType(this.srcType.Elem())
}

func (this *sliceIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &sliceIterator{
		srcType:  this.srcType,
		root:     root,
		elemIter: this.elemIter.CloneFromTemplate(root),
	}
}

func (this *sliceIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
		return
	}
	if this.root.addReference(v) {
		return
	}

	this.root.eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		this.elemIter.Iterate(v.Index(i))
	}
	this.root.eventReceiver.OnEnd()
}

// -----
// Array
// -----

type arrayIterator struct {
	srcType  reflect.Type
	elemIter ObjectIterator
	root     *RootObjectIterator
}

func newArrayIterator(srcType reflect.Type) ObjectIterator {
	return &arrayIterator{
		srcType: srcType,
	}
}

func (this *arrayIterator) PostCacheInitIterator() {
	this.elemIter = getIteratorForType(this.srcType.Elem())
}

func (this *arrayIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &arrayIterator{
		srcType:  this.srcType,
		root:     root,
		elemIter: this.elemIter.CloneFromTemplate(root),
	}
}

func (this *arrayIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		this.elemIter.Iterate(v.Index(i))
	}
	this.root.eventReceiver.OnEnd()
}

// ---
// Map
// ---

type mapIterator struct {
	srcType   reflect.Type
	keyIter   ObjectIterator
	valueIter ObjectIterator
	root      *RootObjectIterator
}

func newMapIterator(srcType reflect.Type) ObjectIterator {
	return &mapIterator{
		srcType: srcType,
	}
}

func (this *mapIterator) PostCacheInitIterator() {
	this.keyIter = getIteratorForType(this.srcType.Key())
	this.valueIter = getIteratorForType(this.srcType.Elem())
}

func (this *mapIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	return &mapIterator{
		srcType:   this.srcType,
		keyIter:   this.keyIter.CloneFromTemplate(root),
		valueIter: this.valueIter.CloneFromTemplate(root),
		root:      root,
	}
}

func (this *mapIterator) Iterate(v reflect.Value) {
	if v.IsNil() {
		this.root.eventReceiver.OnNil()
		return
	}
	if this.root.addReference(v) {
		return
	}

	this.root.eventReceiver.OnMap()

	iter := mapRange(v)
	for iter.Next() {
		this.keyIter.Iterate(iter.Key())
		this.valueIter.Iterate(iter.Value())
	}

	this.root.eventReceiver.OnEnd()
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
	root           *RootObjectIterator
}

func newStructIterator(srcType reflect.Type) ObjectIterator {
	return &structIterator{
		srcType: srcType,
	}
}

func (this *structIterator) PostCacheInitIterator() {
	for i := 0; i < this.srcType.NumField(); i++ {
		field := this.srcType.Field(i)
		if isFieldExported(field.Name) {
			iterator := &structIteratorField{
				Name:     field.Name,
				Index:    i,
				Iterator: getIteratorForType(field.Type),
			}
			this.fieldIterators = append(this.fieldIterators, iterator)
		}
	}
}

func (this *structIterator) CloneFromTemplate(root *RootObjectIterator) ObjectIterator {
	that := &structIterator{
		srcType: this.srcType,
		root:    root,
	}
	that.fieldIterators = make([]*structIteratorField, 0, len(this.fieldIterators))
	for _, iter := range this.fieldIterators {
		that.fieldIterators = append(that.fieldIterators, newStructIteratorField(iter.Name, iter.Index, iter.Iterator.CloneFromTemplate(root)))
	}
	return that
}

func (this *structIterator) Iterate(v reflect.Value) {
	this.root.eventReceiver.OnMap()

	for _, iter := range this.fieldIterators {
		this.root.eventReceiver.OnString(iter.Name)
		iter.Iterator.Iterate(v.Field(iter.Index))
	}

	this.root.eventReceiver.OnEnd()
}
