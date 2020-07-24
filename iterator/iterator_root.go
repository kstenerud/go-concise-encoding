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
	options         options.IteratorOptions
	session         *Session
}

// Create a new root object iterator that will send data events to eventReceiver.
// If options is nil, default options will be used.
func NewRootObjectIterator(session *Session, eventReceiver events.DataEventReceiver, options *options.IteratorOptions) *RootObjectIterator {
	_this := &RootObjectIterator{}
	_this.Init(session, eventReceiver, options)
	return _this
}

// Initialize this iterator to send data events to eventReceiver.
// If options is nil, default options will be used.
func (_this *RootObjectIterator) Init(session *Session, eventReceiver events.DataEventReceiver, options *options.IteratorOptions) {
	_this.options = *options.WithDefaultsApplied()
	_this.session = session
	_this.eventReceiver = eventReceiver
}

// Iterates over an object, sending events to the root iterator's
// DataEventReceiver as it visits all elements of the value.
func (_this *RootObjectIterator) Iterate(object interface{}) {
	if object == nil {
		_this.eventReceiver.OnVersion(_this.options.ConciseEncodingVersion)
		_this.eventReceiver.OnNil()
		_this.eventReceiver.OnEndDocument()
		return
	}
	_this.findReferences(object)
	rv := reflect.ValueOf(object)
	iterator := _this.session.GetIteratorForType(rv.Type())
	_this.eventReceiver.OnVersion(_this.options.ConciseEncodingVersion)
	iterator.IterateObject(rv, _this.eventReceiver, _this)
	_this.eventReceiver.OnEndDocument()
}

func (_this *RootObjectIterator) AddReference(v reflect.Value) (didGenerateReferenceEvent bool) {
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

// ============================================================================

// Internal

func (_this *RootObjectIterator) findReferences(value interface{}) {
	if _this.options.UseReferences {
		_this.foundReferences = duplicates.FindDuplicatePointers(value)
		_this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}
}
