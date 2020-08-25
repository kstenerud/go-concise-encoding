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
	"fmt"
	"reflect"
	"strings"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// ---------
// Interface
// ---------

type interfaceIterator struct {
	// Instance Data
	fetchInstance FetchIterator
}

func newInterfaceIterator(_ reflect.Type) ObjectIterator {
	return &interfaceIterator{}
}

func (_this *interfaceIterator) InitTemplate(_ FetchIterator) {
}

func (_this *interfaceIterator) NewInstance() ObjectIterator {
	return _this
}

func (_this *interfaceIterator) InitInstance(fetchInstance FetchIterator, _ *options.IteratorOptions) {
	_this.fetchInstance = fetchInstance
}

func (_this *interfaceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	if v.IsNil() {
		eventReceiver.OnNil()
	} else {
		elem := v.Elem()
		iter := _this.fetchInstance(elem.Type())
		iter.IterateObject(elem, eventReceiver, addReference)
	}
}

// -------
// Pointer
// -------

type pointerIterator struct {
	// Template Data
	pointerType reflect.Type

	// Instance Data
	elemIter ObjectIterator
}

func newPointerIterator(pointerType reflect.Type) ObjectIterator {
	return &pointerIterator{pointerType: pointerType}
}

func (_this *pointerIterator) InitTemplate(fetchTemplate FetchIterator) {
	// Fetch to cache the template.
	fetchTemplate(_this.pointerType.Elem())
}

func (_this *pointerIterator) NewInstance() ObjectIterator {
	return &pointerIterator{
		pointerType: _this.pointerType,
	}
}

func (_this *pointerIterator) InitInstance(fetchInstance FetchIterator, _ *options.IteratorOptions) {
	_this.elemIter = fetchInstance(_this.pointerType.Elem())
}

func (_this *pointerIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if addReference(v) {
		return
	}
	_this.elemIter.IterateObject(v.Elem(), eventReceiver, addReference)
}

// -----
// Slice
// -----

type sliceIterator struct {
	// Template Data
	sliceType reflect.Type

	// Instance Data
	elemIter ObjectIterator
}

func newSliceIterator(sliceType reflect.Type) ObjectIterator {
	return &sliceIterator{
		sliceType: sliceType,
	}
}

func (_this *sliceIterator) InitTemplate(fetchTemplate FetchIterator) {
	// Fetch to cache the template.
	fetchTemplate(_this.sliceType.Elem())
}

func (_this *sliceIterator) NewInstance() ObjectIterator {
	return &sliceIterator{
		sliceType: _this.sliceType,
	}
}

func (_this *sliceIterator) InitInstance(fetchInstance FetchIterator, _ *options.IteratorOptions) {
	_this.elemIter = fetchInstance(_this.sliceType.Elem())
}

func (_this *sliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if addReference(v) {
		return
	}

	eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.IterateObject(v.Index(i), eventReceiver, addReference)
	}
	eventReceiver.OnEnd()
}

// -----
// Array
// -----

type arrayIterator struct {
	// Template Data
	arrayType reflect.Type

	// Instance Data
	elemIter ObjectIterator
}

func newArrayIterator(arrayType reflect.Type) ObjectIterator {
	return &arrayIterator{
		arrayType: arrayType,
	}
}

func (_this *arrayIterator) InitTemplate(fetchTemplate FetchIterator) {
	// Fetch to cache the template.
	fetchTemplate(_this.arrayType.Elem())
}

func (_this *arrayIterator) NewInstance() ObjectIterator {
	return &arrayIterator{
		arrayType: _this.arrayType,
	}
}

func (_this *arrayIterator) InitInstance(fetchInstance FetchIterator, _ *options.IteratorOptions) {
	_this.elemIter = fetchInstance(_this.arrayType.Elem())
}

func (_this *arrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.IterateObject(v.Index(i), eventReceiver, addReference)
	}
	eventReceiver.OnEnd()
}

// ---
// Map
// ---

type mapIterator struct {
	// Template Data
	mapType reflect.Type

	// Instance Data
	keyIter   ObjectIterator
	valueIter ObjectIterator
}

func newMapIterator(mapType reflect.Type) ObjectIterator {
	return &mapIterator{
		mapType: mapType,
	}
}

func (_this *mapIterator) InitTemplate(fetchTemplate FetchIterator) {
	// Fetch to cache the template.
	fetchTemplate(_this.mapType.Key())
	fetchTemplate(_this.mapType.Elem())
}

func (_this *mapIterator) NewInstance() ObjectIterator {
	return &mapIterator{
		mapType: _this.mapType,
	}
}

func (_this *mapIterator) InitInstance(fetchInstance FetchIterator, _ *options.IteratorOptions) {
	_this.keyIter = fetchInstance(_this.mapType.Key())
	_this.valueIter = fetchInstance(_this.mapType.Elem())
}

func (_this *mapIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if addReference(v) {
		return
	}

	eventReceiver.OnMap()

	iter := common.MapRange(v)
	for iter.Next() {
		_this.keyIter.IterateObject(iter.Key(), eventReceiver, addReference)
		_this.valueIter.IterateObject(iter.Value(), eventReceiver, addReference)
	}

	eventReceiver.OnEnd()
}

// ------
// Struct
// ------

type structField struct {
	Name      string
	Type      reflect.Type
	Index     int
	Omit      bool
	OmitEmpty bool
	// TODO: OmitValue
	OmitValue string
}

func newStructField(name string, t reflect.Type, index int) *structField {
	return &structField{
		Name:  name,
		Type:  t,
		Index: index,
	}
}

func (_this *structField) applyTags(tags string) {
	if tags == "" {
		return
	}

	requiresValue := func(kv []string, key string) {
		if len(kv) != 2 {
			panic(fmt.Errorf(`tag key "%s" requires a value`, key))
		}
	}

	for _, entry := range strings.Split(tags, ",") {
		kv := strings.Split(entry, "=")
		switch strings.TrimSpace(kv[0]) {
		// TODO: lowercase/origcase
		// TODO: recurse/norecurse?
		// TODO: nil?
		// TODO: type=f16, f10.x, i2, i8, i10, i16, string, vstring
		case "-":
			_this.Omit = true
		case "omit":
			if len(kv) == 1 {
				_this.Omit = true
			} else {
				_this.OmitValue = strings.TrimSpace(kv[1])
			}
		case "omitempty":
			// TODO: Implement omitempty
			_this.OmitEmpty = true
		case "name":
			requiresValue(kv, "name")
			_this.Name = strings.TrimSpace(kv[1])
		default:
			panic(fmt.Errorf("%v: Unknown Concise Encoding struct tag field", entry))
		}
	}
}

type structIterator struct {
	// Template Data
	structType reflect.Type
	fields     []*structField

	// Instance Data
	iterators []ObjectIterator
	opts      *options.IteratorOptions
}

func newStructIterator(structType reflect.Type) ObjectIterator {
	return &structIterator{
		structType: structType,
	}
}

func (_this *structIterator) InitTemplate(_ FetchIterator) {
	for i := 0; i < _this.structType.NumField(); i++ {
		reflectField := _this.structType.Field(i)
		if common.IsFieldExported(reflectField.Name) {
			field := newStructField(reflectField.Name, reflectField.Type, i)
			field.applyTags(reflectField.Tag.Get("ce"))

			if !field.Omit {
				_this.fields = append(_this.fields, field)
			}
		}
	}
}

func (_this *structIterator) NewInstance() ObjectIterator {
	return &structIterator{
		structType: _this.structType,
		fields:     _this.fields,
	}
}

func (_this *structIterator) InitInstance(fetchInstance FetchIterator, opts *options.IteratorOptions) {
	for _, field := range _this.fields {
		_this.iterators = append(_this.iterators, fetchInstance(field.Type))
	}
	_this.opts = opts
}

func (_this *structIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference) {
	eventReceiver.OnMap()

	for i, field := range _this.fields {
		name := field.Name
		if _this.opts.LowercaseStructFieldNames {
			name = common.ASCIIToLower(name)
		}
		eventReceiver.OnString([]byte(name))
		iterator := _this.iterators[i]
		iterator.IterateObject(v.Field(field.Index), eventReceiver, addReference)
	}

	eventReceiver.OnEnd()
}
