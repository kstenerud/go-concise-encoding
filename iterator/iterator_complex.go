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

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/kstenerud/go-concise-encoding/internal/common"
)

// ---------
// Interface
// ---------

type interfaceIterator struct {
	session *Session
	opts    *options.IteratorOptions
}

func newInterfaceIterator(srcType reflect.Type) ObjectIterator {
	return &interfaceIterator{}
}

func (_this *interfaceIterator) PostCacheInitIterator(session *Session) {
	_this.session = session
}

func (_this *interfaceIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// return &interfaceIterator{
	// 	session: _this.session,
	// 	opts:    opts,
	// }
}

func (_this *interfaceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
	} else {
		elem := v.Elem()
		iter := _this.session.GetIteratorForType(elem.Type()).CloneFromTemplate(_this.opts)
		iter.IterateObject(elem, eventReceiver, references)
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

func (_this *pointerIterator) PostCacheInitIterator(session *Session) {
	_this.elemIter = session.GetIteratorForType(_this.srcType.Elem())
}

func (_this *pointerIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// return &pointerIterator{
	// 	srcType:  _this.srcType,
	// 	elemIter: _this.elemIter.CloneFromTemplate(opts),
	// }
}

func (_this *pointerIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if references.AddReference(v) {
		return
	}
	_this.elemIter.IterateObject(v.Elem(), eventReceiver, references)
}

// -----------
// uint8 array
// -----------

type uint8ArrayIterator struct {
}

func newUInt8ArrayIterator() ObjectIterator {
	return &uint8ArrayIterator{}
}

func (_this *uint8ArrayIterator) PostCacheInitIterator(session *Session) {
}

func (_this *uint8ArrayIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
}

func (_this *uint8ArrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.CanAddr() {
		eventReceiver.OnBytes(v.Slice(0, v.Len()).Bytes())
	} else {
		tempSlice := make([]byte, v.Len())
		tempLen := v.Len()
		for i := 0; i < tempLen; i++ {
			tempSlice[i] = v.Index(i).Interface().(uint8)
		}
		eventReceiver.OnBytes(tempSlice)
	}
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

func (_this *sliceIterator) PostCacheInitIterator(session *Session) {
	_this.elemIter = session.GetIteratorForType(_this.srcType.Elem())
}

func (_this *sliceIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// return &sliceIterator{
	// 	srcType:  _this.srcType,
	// 	elemIter: _this.elemIter.CloneFromTemplate(opts),
	// }
}

func (_this *sliceIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if references.AddReference(v) {
		return
	}

	eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.IterateObject(v.Index(i), eventReceiver, references)
	}
	eventReceiver.OnEnd()
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

func (_this *arrayIterator) PostCacheInitIterator(session *Session) {
	_this.elemIter = session.GetIteratorForType(_this.srcType.Elem())
}

func (_this *arrayIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// return &arrayIterator{
	// 	srcType:  _this.srcType,
	// 	elemIter: _this.elemIter.CloneFromTemplate(opts),
	// }
}

func (_this *arrayIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	eventReceiver.OnList()
	length := v.Len()
	for i := 0; i < length; i++ {
		_this.elemIter.IterateObject(v.Index(i), eventReceiver, references)
	}
	eventReceiver.OnEnd()
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

func (_this *mapIterator) PostCacheInitIterator(session *Session) {
	_this.keyIter = session.GetIteratorForType(_this.srcType.Key())
	_this.valueIter = session.GetIteratorForType(_this.srcType.Elem())
}

func (_this *mapIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// return &mapIterator{
	// 	srcType:   _this.srcType,
	// 	keyIter:   _this.keyIter.CloneFromTemplate(opts),
	// 	valueIter: _this.valueIter.CloneFromTemplate(opts),
	// }
}

func (_this *mapIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	if v.IsNil() {
		eventReceiver.OnNil()
		return
	}
	if references.AddReference(v) {
		return
	}

	eventReceiver.OnMap()

	iter := common.MapRange(v)
	for iter.Next() {
		_this.keyIter.IterateObject(iter.Key(), eventReceiver, references)
		_this.valueIter.IterateObject(iter.Value(), eventReceiver, references)
	}

	eventReceiver.OnEnd()
}

// ------
// Struct
// ------

type structField struct {
	Name      string
	Index     int
	Omit      bool
	OmitEmpty bool
	// TODO: OmitValue
	OmitValue string
}

func newStructField(name string, index int) *structField {
	return &structField{
		Name:  name,
		Index: index,
	}
}

func (_this *structField) applyTags(tags string) {
	if tags == "" {
		return
	}

	requiresValue := func(kv []string, key string) {
		if len(kv) != 2 {
			panic(fmt.Errorf(`Tag key "%s" requires a value`, key))
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
	srcType   reflect.Type
	fields    []*structField
	iterators []ObjectIterator
	opts      *options.IteratorOptions
}

func newStructIterator(srcType reflect.Type) ObjectIterator {
	return &structIterator{
		srcType: srcType,
	}
}

func (_this *structIterator) PostCacheInitIterator(session *Session) {
	for i := 0; i < _this.srcType.NumField(); i++ {
		reflectField := _this.srcType.Field(i)
		if common.IsFieldExported(reflectField.Name) {
			field := &structField{
				Name:  reflectField.Name,
				Index: i,
			}
			field.applyTags(reflectField.Tag.Get("ce"))

			if !field.Omit {
				_this.fields = append(_this.fields, field)
				_this.iterators = append(_this.iterators, session.GetIteratorForType(reflectField.Type))
			}
		}
	}
}

func (_this *structIterator) CloneFromTemplate(opts *options.IteratorOptions) ObjectIterator {
	return _this
	// iterators := make([]ObjectIterator, len(_this.iterators))
	// for i, v := range _this.iterators {
	// 	iterators[i] = v.CloneFromTemplate(opts)
	// }

	// return &structIterator{
	// 	srcType:   _this.srcType,
	// 	fields:    _this.fields,
	// 	iterators: iterators,
	// 	opts:      opts,
	// }
}

func (_this *structIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	eventReceiver.OnMap()

	for i, field := range _this.fields {
		name := field.Name
		// if _this.opts.LowercaseStructFieldNames {
		// 	name = strings.ToLower(name)
		// }
		eventReceiver.OnString([]byte(name))
		iterator := _this.iterators[i]
		iterator.IterateObject(v.Field(field.Index), eventReceiver, references)
	}

	eventReceiver.OnEnd()
}
