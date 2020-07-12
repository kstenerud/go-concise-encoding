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
	"sync"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/version"
)

var defaultIteratorOptions = options.IteratorOptions{
	ConciseEncodingVersion: version.ConciseEncodingVersion,
}

func DefaultIteratorOptions() *options.IteratorOptions {
	opts := defaultIteratorOptions
	return &opts
}

// Iterate over an object (recursively), calling the eventReceiver as data is
// encountered. If options is nil, a zero value will be provided.
func IterateObject(value interface{}, eventReceiver events.DataEventReceiver, options *options.IteratorOptions) {
	iter := NewRootObjectIterator(eventReceiver, options)
	iter.Iterate(value, nil)
}

// ObjectIterator iterates through a value, calling callback methods as it goes.
type ObjectIterator interface {
	// Iterates over a value, potentially calling other iterators as it goes.
	Iterate(v reflect.Value, root *RootObjectIterator)

	// PostCacheInitIterator is called after the iterator template is saved to
	// cache but before use, so that lookups succeed on cyclic type references.
	PostCacheInitIterator()
}

var iterators sync.Map

func init() {
	// Pre-cache the most common iterators
	for _, t := range common.KeyableTypes {
		getIteratorForType(t)
		getIteratorForType(reflect.PtrTo(t))
		getIteratorForType(reflect.SliceOf(t))
		for _, u := range common.KeyableTypes {
			getIteratorForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
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
		case common.TypeTime:
			return newTimeIterator()
		case common.TypeCompactTime:
			return newCompactTimeIterator()
		case common.TypeDFloat:
			return newDFloatIterator()
		case common.TypeURL:
			return newURLIterator()
		case common.TypeBigInt:
			return newBigIntIterator()
		case common.TypeBigFloat:
			return newBigFloatIterator()
		case common.TypeBigDecimalFloat:
			return newBigDecimalFloatIterator()
		default:
			return newStructIterator(t)
		}
	case reflect.Ptr:
		switch t {
		case common.TypePCompactTime:
			return newPCompactTimeIterator()
		case common.TypePURL:
			return newPURLIterator()
		case common.TypePBigInt:
			return newPBigIntIterator()
		case common.TypePBigFloat:
			return newPBigFloatIterator()
		case common.TypePBigDecimalFloat:
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

func applyDefaultIteratorOptions(original *options.IteratorOptions) *options.IteratorOptions {
	var options options.IteratorOptions
	if original == nil {
		options = defaultIteratorOptions
	} else {
		options = *original
		if options.ConciseEncodingVersion < 1 {
			options.ConciseEncodingVersion = defaultIteratorOptions.ConciseEncodingVersion
		}
	}
	return &options
}
