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

// Iterate over an object (recursively), calling the eventReceiver as data is
// encountered. If options is nil, default options will be used.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this function.
func IterateObject(object interface{}, eventReceiver events.DataEventReceiver, options *options.IteratorOptions) {
	iter := NewRootObjectIterator(NewSession(), eventReceiver, options)
	iter.Iterate(object)
}

// ObjectIterator iterates through an object, calling DataEventReceiver callback
// methods as it encounters data.
type ObjectIterator interface {
	// Iterate over an object.
	IterateObject(object reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator)

	// PostCacheInitIterator is called after the iterator template is saved to
	// cache but before use, so that lookups succeed on cyclic type references.
	PostCacheInitIterator(session *Session)
}

// ConvertToCustomBytesFunction converts a value to a custom array of bytes.
// This allows fully user-configurable converting to custom Concise Encoding
// data.
//
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom
type ConvertToCustomBytesFunction func(v reflect.Value) (asBytes []byte, err error)
