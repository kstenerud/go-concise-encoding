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
	"strconv"

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/version"
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
	context         Context
	config          *configuration.Configuration
	referenceIdBuff []byte
}

// Create a new root object iterator that will send data events to eventReceiver.
func NewRootObjectIterator(context *Context,
	eventReceiver events.DataEventReceiver,
	config *configuration.Configuration) *RootObjectIterator {

	_this := &RootObjectIterator{}
	_this.Init(context, eventReceiver, config)
	return _this
}

// Initialize this iterator to send data events to eventReceiver.
func (_this *RootObjectIterator) Init(context *Context,
	eventReceiver events.DataEventReceiver,
	config *configuration.Configuration) {

	_this.config = config
	_this.context = iteratorContext(context,
		eventReceiver,
		_this.addLocalReference)
	_this.referenceIdBuff = make([]byte, 0, 16)
}

// Iterates over an object, sending events to the root iterator's
// DataEventReceiver as it visits all elements of the value.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this function.
func (_this *RootObjectIterator) Iterate(object interface{}) {
	_this.context.EventReceiver.OnBeginDocument()
	_this.context.EventReceiver.OnVersion(version.ConciseEncodingVersion)

	if object == nil {
		_this.context.NotifyNil()
		_this.context.EventReceiver.OnEndDocument()
		return
	}

	if _this.config.Iterator.RecursionSupport {
		_this.foundReferences = duplicates.FindDuplicatePointers(object)
		_this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}

	// Generate all record types at the top of the document
	for _, entry := range _this.context.RecordTypeOrder {
		entry.Iterator(&_this.context, reflect.ValueOf(nil))
	}

	rv := reflect.ValueOf(object)
	iterate := _this.context.GetIteratorForType(rv.Type())
	iterate(&_this.context, rv)
	_this.context.EventReceiver.OnEndDocument()
}

// ============================================================================
// Internal

func (_this *RootObjectIterator) getNamedLocalReference(ptr duplicates.TypedPointer) (name []byte, exists bool) {
	num, exists := _this.namedReferences[ptr]
	if !exists {
		num = _this.nextMarkerName
		_this.namedReferences[ptr] = num
		_this.nextMarkerName++
	}

	_this.referenceIdBuff = strconv.AppendInt(_this.referenceIdBuff[:0], int64(num), 10)
	return _this.referenceIdBuff, exists
}

func (_this *RootObjectIterator) addLocalReference(v reflect.Value) (didGenerateReferenceEvent bool) {
	if !_this.config.Iterator.RecursionSupport {
		return false
	}

	ptr := duplicates.TypedPointerOfRV(v)
	if !_this.foundReferences[ptr] {
		return false
	}

	name, exists := _this.getNamedLocalReference(ptr)
	if !exists {
		_this.context.EventReceiver.OnMarker(name)
		return false
	}

	_this.context.EventReceiver.OnReferenceLocal(name)
	return true
}
