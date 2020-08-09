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
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-duplicates"
)

// RootObjectIterator acts as a top-level iterator, coordinating iteration
// through an arbitrary object via sub-iterators.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type RootObjectIterator struct {
	foundReferences map[duplicates.TypedPointer]bool
	namedReferences map[duplicates.TypedPointer]uint32
	nextMarkerName  uint32
	eventReceiver   events.DataEventReceiver
	opts            options.IteratorOptions
	fetchTemplate   FetchIterator
	iterators       map[reflect.Type]ObjectIterator
}

// Create a new root object iterator that will send data events to eventReceiver.
// If opts is nil, default options will be used.
func NewRootObjectIterator(fetchTemplate FetchIterator, eventReceiver events.DataEventReceiver, opts *options.IteratorOptions) *RootObjectIterator {
	_this := &RootObjectIterator{}
	_this.Init(fetchTemplate, eventReceiver, opts)
	return _this
}

// Initialize this iterator to send data events to eventReceiver.
// If opts is nil, default options will be used.
func (_this *RootObjectIterator) Init(fetchTemplate FetchIterator, eventReceiver events.DataEventReceiver, opts *options.IteratorOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.fetchTemplate = fetchTemplate
	_this.eventReceiver = eventReceiver
	_this.iterators = make(map[reflect.Type]ObjectIterator)
}

// Iterates over an object, sending events to the root iterator's
// DataEventReceiver as it visits all elements of the value.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this function.
func (_this *RootObjectIterator) Iterate(object interface{}) {
	_this.eventReceiver.OnBeginDocument()
	_this.eventReceiver.OnVersion(_this.opts.ConciseEncodingVersion)
	if object == nil {
		_this.eventReceiver.OnNil()
		_this.eventReceiver.OnEndDocument()
		return
	}

	if _this.opts.RecursionSupport {
		_this.foundReferences = duplicates.FindDuplicatePointers(object)
		_this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}

	rv := reflect.ValueOf(object)
	iterator := _this.getIteratorInstanceForType(rv.Type())
	iterator.IterateObject(rv, _this.eventReceiver, _this.addReference)
	_this.eventReceiver.OnEndDocument()
}

// ============================================================================
// Internal

func (_this *RootObjectIterator) addReference(v reflect.Value) (didGenerateReferenceEvent bool) {
	if !_this.opts.RecursionSupport {
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

func (_this *RootObjectIterator) getIteratorInstanceForType(t reflect.Type) ObjectIterator {
	iterator := _this.iterators[t]
	if iterator == nil {
		template := _this.fetchTemplate(t)
		iterator = template.NewInstance()
		_this.iterators[t] = iterator
		iterator.InitInstance(_this.getIteratorInstanceForType, &_this.opts)
	}
	return iterator
}
