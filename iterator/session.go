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
	"fmt"
	"reflect"
	"sync"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// An iterator session holds a cache of known mappings of types to iterators.
// It is designed to be cloned so that any user-supplied custom iterators exist
// only in their own session, and don't pollute the base mapping and cause
// unintended behavior in codec activity elsewhere in the program.
type Session struct {
	iterators sync.Map
	opts      options.IteratorSessionOptions
}

// Start a new iterator session. It will inherit the iterators of its parent.
// If parent is nil, it will inherit from the root session, which has iterators
// for all basic go types.
// If opts is nil, default options will be used.
func NewSession(parent *Session, opts *options.IteratorSessionOptions) *Session {
	_this := &Session{}
	_this.Init(parent, opts)
	return _this
}

// Initialize an iterator session. It will inherit the iterators of its parent.
// If parent is nil, it will inherit from the root session, which has iterators
// for all basic go types.
// If opts is nil, default options will be used.
func (_this *Session) Init(parent *Session, opts *options.IteratorSessionOptions) {
	opts = opts.WithDefaultsApplied()
	if parent == nil {
		parent = &rootSession
	}

	parent.iterators.Range(func(k interface{}, v interface{}) bool {
		_this.iterators.Store(k, v)
		return true
	})

	_this.opts = *opts

	for t, converter := range _this.opts.CustomBinaryConverters {
		_this.RegisterIteratorForType(t, newCustomBinaryIterator(converter))
	}
	for t, converter := range _this.opts.CustomTextConverters {
		_this.RegisterIteratorForType(t, newCustomTextIterator(converter))
	}
}

// Creates a new iterator that sends data events to eventReceiver.
// If opts is nil, default options will be used.
func (_this *Session) NewIterator(eventReceiver events.DataEventReceiver, opts *options.IteratorOptions) *RootObjectIterator {
	return NewRootObjectIterator(_this.GetIteratorTemplateForType, eventReceiver, opts)
}

// Register a specific iterator for a type.
// If an iterator has already been registered for this type, it will be replaced.
// This function is thread-safe.
func (_this *Session) RegisterIteratorForType(t reflect.Type, iterator ObjectIterator) {
	_this.iterators.Store(t, iterator)
}

// Get an iterator template for the specified type. If a registered template
// doesn't yet exist, a new default template will be generated and registered.
// This method is thread-safe.
func (_this *Session) GetIteratorTemplateForType(t reflect.Type) ObjectIterator {
	if intf, loaded := _this.iterators.Load(t); loaded {
		return intf.(ObjectIterator)
	}

	intf, loaded := _this.iterators.LoadOrStore(t, defaultIteratorForType(t))
	iterator := intf.(ObjectIterator)
	if !loaded {
		iterator.InitTemplate(_this.GetIteratorTemplateForType)
	}
	return iterator
}

// ============================================================================
// Internal

// The root session caches the most common iterators. All sessions inherit
// these cached values.
var rootSession Session

func init() {
	rootSession.Init(nil, nil)

	for _, t := range common.KeyableTypes {
		rootSession.GetIteratorTemplateForType(t)
		rootSession.GetIteratorTemplateForType(reflect.PtrTo(t))
		rootSession.GetIteratorTemplateForType(reflect.SliceOf(t))
		for _, u := range common.KeyableTypes {
			rootSession.GetIteratorTemplateForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
		rootSession.GetIteratorTemplateForType(t)
		rootSession.GetIteratorTemplateForType(reflect.PtrTo(t))
		rootSession.GetIteratorTemplateForType(reflect.SliceOf(t))
	}
}

func defaultIteratorForType(t reflect.Type) ObjectIterator {
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
	case reflect.Interface:
		return newInterfaceIterator(t)
	case reflect.Array:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return newUint8ArrayIterator()
		case reflect.Uint16:
			return newUint16ArrayIterator()
		default:
			return newArrayIterator(t)
		}
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return newUint8SliceIterator()
		case reflect.Uint16:
			return newUint16SliceIterator()
		default:
			return newSliceIterator(t)
		}
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
