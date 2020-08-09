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

// Iterators iterate through go objects, producing data events.
package iterator

import (
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

// ObjectIterator iterates through an object, calling DataEventReceiver callback
// methods as it encounters data.
type ObjectIterator interface {

	// Iterate over an object (recursively), calling the eventReceiver as data is
	// encountered.
	//
	// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
	// to recover() at an appropriate location when calling this function.
	IterateObject(object reflect.Value, eventReceiver events.DataEventReceiver, addReference AddReference)

	// Initialize as a template. This is called after the template has been
	// cached in order to support recursive object graphs.
	InitTemplate(fetchTemplate FetchIterator)

	// Create an instance using this object as a template
	NewInstance() ObjectIterator

	// Initialize as an instance. This is called after the instance has been
	// cached in order to support recursive object graphs.
	// opts must not be nil.
	InitInstance(fetchInstance FetchIterator, opts *options.IteratorOptions)
}

// Fetches an interator for the specified type, building and caching one as needed.
type FetchIterator func(reflect.Type) ObjectIterator

// Attempt to add a reference for a possibly duplicate pointer (meaning
// that the same pointer exists elsewhere in the object graph currently
// being iterated). If the pointer is indeed a duplicate pointer, it will
// either generate a marker event (first encounter) or a reference event
// (subsequent encounters).
// Returns true if it generated a reference event, false if it generated a
// marker event or no event at all (i.e. not a duplicate pointer).
type AddReference func(pointer reflect.Value) (didGenerateReferenceEvent bool)
