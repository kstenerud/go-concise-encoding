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
	"fmt"
	"reflect"
	"sync"
)

type IteratorOptions struct {
	// If useReferences is true, the iterator will also look for duplicate
	// pointers to data, generating marker and reference events rather than
	// walking the object again. This is useful for cyclic or recursive data
	// structures.
	UseReferences bool
}

// Iterate over an object (recursively), calling the eventReceiver as data is
// encountered. If options is nil, a zero value will be provided.
func IterateObject(value interface{}, eventReceiver DataEventReceiver, options *IteratorOptions) {
	iter := NewRootObjectIterator(eventReceiver, options)
	iter.Iterate(value)
}

// ObjectIterator iterates through a value, calling callback methods as it goes.
type ObjectIterator interface {
	// Iterate iterates over a value, potentially calling other iterators as
	// it goes.
	Iterate(v reflect.Value)

	// PostCacheInitIterator is called after the iterator template is saved to
	// cache but before use, so that lookups succeed on cyclic type references.
	PostCacheInitIterator()

	// CloneFromTemplate clones from this iterator as a template, adding contextual data.
	CloneFromTemplate(root *RootObjectIterator) ObjectIterator
}

var iterators sync.Map

func init() {
	// Pre-cache the most common iterators
	for _, t := range keyableTypes {
		getIteratorForType(t)
		getIteratorForType(reflect.PtrTo(t))
		getIteratorForType(reflect.SliceOf(t))
		for _, u := range keyableTypes {
			getIteratorForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range nonKeyableTypes {
		getIteratorForType(t)
		getIteratorForType(reflect.PtrTo(t))
		getIteratorForType(reflect.SliceOf(t))
	}
}

// Register a specific iterator for a type.
// If an iterator has already been registered for this type, it will be replaced.
func RegisterIteratorForType(t reflect.Type, iterator ObjectIterator) {
	iterators.Store(t, iterator)
}

func generateIteratorForType(t reflect.Type) ObjectIterator {
	switch t.Kind() {
	case reflect.Bool:
		return newBoolIterator()
	case reflect.String:
		return newStringIterator()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newIntIterator()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return newUintIterator()
	case reflect.Float32, reflect.Float64:
		return newFloatIterator()
	case reflect.Complex64, reflect.Complex128:
		return newComplexIterator()
	case reflect.Interface:
		return newInterfaceIterator(t)
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return newUInt8ArrayIterator()
		}
		return newArrayIterator(t)
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return newUInt8SliceIterator()
		}
		return newSliceIterator(t)
	case reflect.Map:
		return newMapIterator(t)
	case reflect.Struct:
		switch t {
		case typeTime:
			return newTimeIterator()
		case typeDFloat:
			return newDFloatIterator()
		case typeURL:
			return newURLIterator()
		case typeBigInt:
			return newBigIntIterator()
		case typeBigDecimalFloat:
			return newBigDecimalFloatIterator()
		default:
			return newStructIterator(t)
		}
	case reflect.Ptr:
		switch t {
		case typePURL:
			return newPURLIterator()
		case typePBigInt:
			return newPBigIntIterator()
		case typePBigDecimalFloat:
			return newPBigDecimalFloatIterator()
		default:
			return newPointerIterator(t)
		}
	default:
		panic(fmt.Errorf("BUG: Unhandled type %v", t))
	}
}

func getIteratorForType(t reflect.Type) ObjectIterator {
	if iterator, ok := iterators.Load(t); ok {
		return iterator.(ObjectIterator)
	}

	iterator, _ := iterators.LoadOrStore(t, generateIteratorForType(t))
	iterator.(ObjectIterator).PostCacheInitIterator()
	return iterator.(ObjectIterator)
}
