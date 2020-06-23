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

	"github.com/kstenerud/go-duplicates"
)

// RootObjectIterator iterates recursively depth-first through an object,
// notifying the event receiver as it encounters data.
type RootObjectIterator struct {
	foundReferences map[duplicates.TypedPointer]bool
	namedReferences map[duplicates.TypedPointer]uint32
	nextMarkerName  uint32
	eventReceiver   DataEventReceiver
	options         IteratorOptions
}

func NewRootObjectIterator(eventReceiver DataEventReceiver, options *IteratorOptions) *RootObjectIterator {
	_this := new(RootObjectIterator)
	_this.Init(eventReceiver, options)
	return _this
}

func (_this *RootObjectIterator) Init(eventReceiver DataEventReceiver, options *IteratorOptions) {
	_this.options = *applyDefaultIteratorOptions(options)
	_this.eventReceiver = eventReceiver
}

// Iterates over an object, sending events to the root iterator's
// DataEventReceiver as it visits all elements of the value.
func (_this *RootObjectIterator) Iterate(value interface{}, _ *RootObjectIterator) {
	if value == nil {
		_this.eventReceiver.OnVersion(_this.options.ConciseEncodingVersion)
		_this.eventReceiver.OnNil()
		_this.eventReceiver.OnEndDocument()
		return
	}
	_this.findReferences(value)
	rv := reflect.ValueOf(value)
	iterator := getIteratorForType(rv.Type())
	_this.eventReceiver.OnVersion(_this.options.ConciseEncodingVersion)
	iterator.Iterate(rv, _this)
	_this.eventReceiver.OnEndDocument()
}

func (_this *RootObjectIterator) findReferences(value interface{}) {
	if _this.options.UseReferences {
		_this.foundReferences = duplicates.FindDuplicatePointers(value)
		_this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}
}

func (_this *RootObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if !_this.options.UseReferences {
		return false
	}

	ptr := duplicates.TypedPointerOfRV(v)
	if !_this.foundReferences[ptr] {
		return false
	}

	name, exists := _this.namedReferences[ptr]
	if !exists {
		name = _this.nextMarkerName
		_this.nextMarkerName++
		_this.namedReferences[ptr] = name
		_this.eventReceiver.OnMarker()
		_this.eventReceiver.OnPositiveInt(uint64(name))
		return false
	}

	_this.eventReceiver.OnReference()
	_this.eventReceiver.OnPositiveInt(uint64(name))
	return true
}
