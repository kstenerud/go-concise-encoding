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

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

// An iterator session holds a cache of known mappings of types to iterators.
// It is designed to be cloned so that any user-supplied custom iterators exist
// only in their own session, and don't pollute the base mapping and cause
// unintended behavior in codec activity elsewhere in the program.
type Session struct {
	iteratorFuncs sync.Map
	config        *configuration.IteratorConfiguration
	context       Context
}

// Start a new iterator session. It will inherit the iterators of its parent.
// If parent is nil, it will inherit from the root session, which has iterators
// for all basic go types.
// If config is nil, default configuration will be used.
func NewSession(parent *Session, config *configuration.IteratorConfiguration) *Session {
	_this := &Session{}
	_this.Init(parent, config)
	return _this
}

// Initialize an iterator session. It will inherit the iterators of its parent.
// If parent is nil, it will inherit from the root session, which has iterators
// for all basic go types.
// If config is nil, default configuration will be used.
func (_this *Session) Init(parent *Session, config *configuration.IteratorConfiguration) {
	if config == nil {
		defaultConfig := configuration.DefaultIteratorConfiguration()
		config = &defaultConfig
	} else {
		config.ApplyDefaults()
	}
	_this.config = config

	if parent == nil {
		parent = &rootSession
	}

	parent.iteratorFuncs.Range(func(k interface{}, v interface{}) bool {
		_this.iteratorFuncs.Store(k, v)
		return true
	})

	for t, converter := range _this.config.CustomBinaryConverters {
		_this.RegisterIteratorForType(t, newCustomBinaryIterator(converter))
	}
	for t, converter := range _this.config.CustomTextConverters {
		_this.RegisterIteratorForType(t, newCustomTextIterator(converter))
	}

	_this.context = sessionContext(_this.GetIteratorForType, _this.config)

	for i, entry := range _this.context.RecordTypeOrder {
		typeIter, recordIter := newRecordIterators(&_this.context, entry.Type, entry.Name)
		_this.context.RecordTypeOrder[i].Iterator = typeIter
		_this.RegisterIteratorForType(entry.Type, recordIter)
	}
}

// Creates a new iterator that sends data events to eventReceiver.
// If config is nil, default configuration will be used.
func (_this *Session) NewIterator(eventReceiver events.DataEventReceiver) *RootObjectIterator {

	return NewRootObjectIterator(&_this.context, eventReceiver, _this.config)
}

// Register a specific iterator for a type.
// If an iterator has already been registered for this type, it will be replaced.
func (_this *Session) RegisterIteratorForType(t reflect.Type, iterator IteratorFunction) {
	_this.iteratorFuncs.Store(t, iterator)
}

// Get an iterator template for the specified type. If a registered template
// doesn't yet exist, a new default template will be generated and registered.
func (_this *Session) GetIteratorForType(t reflect.Type) IteratorFunction {
	storedIterator, ok := _this.iteratorFuncs.Load(t)
	if ok {
		return storedIterator.(IteratorFunction)
	}

	var wg sync.WaitGroup
	var iterator IteratorFunction

	wg.Add(1)
	storedIterator, loaded := _this.iteratorFuncs.LoadOrStore(t, IteratorFunction(func(context *Context, value reflect.Value) {
		wg.Wait()
		iterator(context, value)
	}))
	if loaded {
		return storedIterator.(IteratorFunction)
	}

	iterator = _this.getDefaultIteratorForType(t)
	wg.Done()
	_this.iteratorFuncs.Store(t, iterator)
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
		rootSession.GetIteratorForType(t)
		rootSession.GetIteratorForType(reflect.PtrTo(t))
		rootSession.GetIteratorForType(reflect.SliceOf(t))
		rootSession.GetIteratorForType(reflect.SliceOf(reflect.PtrTo(t)))
		for _, u := range common.KeyableTypes {
			rootSession.GetIteratorForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
		rootSession.GetIteratorForType(t)
		rootSession.GetIteratorForType(reflect.PtrTo(t))
		rootSession.GetIteratorForType(reflect.SliceOf(t))
	}
}

func (_this *Session) getDefaultIteratorForType(t reflect.Type) IteratorFunction {
	switch t.Kind() {
	case reflect.Bool:
		return iterateBool
	case reflect.String:
		return iterateString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return iterateInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return iterateUint
	case reflect.Float32, reflect.Float64:
		return iterateFloat
	case reflect.Interface:
		return iterateInterface
	case reflect.Array:
		switch t {
		case common.TypeUID:
			return iterateUID
		}

		switch t.Elem().Kind() {
		case reflect.Uint8:
			return iterateArrayUint8
		case reflect.Uint16:
			return iterateSliceOrArrayUint16
		case reflect.Uint32:
			return iterateSliceOrArrayUint32
		case reflect.Uint64:
			return iterateSliceOrArrayUint64
		case reflect.Uint:
			return iterateSliceOrArrayUint
		case reflect.Int8:
			return iterateSliceOrArrayInt8
		case reflect.Int16:
			return iterateSliceOrArrayInt16
		case reflect.Int32:
			return iterateSliceOrArrayInt32
		case reflect.Int64:
			return iterateSliceOrArrayInt64
		case reflect.Int:
			return iterateSliceOrArrayInt
		case reflect.Float32:
			return iterateSliceOrArrayFloat32
		case reflect.Float64:
			return iterateSliceOrArrayFloat64
		case reflect.Bool:
			return iterateSliceOrArrayBool
		default:
			return newSliceOrArrayAsListIterator(&_this.context, t)
		}
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return iterateSliceUint8
		case reflect.Uint16:
			return iterateSliceOrArrayUint16
		case reflect.Uint32:
			return iterateSliceOrArrayUint32
		case reflect.Uint64:
			return iterateSliceOrArrayUint64
		case reflect.Uint:
			return iterateSliceOrArrayUint
		case reflect.Int8:
			return iterateSliceOrArrayInt8
		case reflect.Int16:
			return iterateSliceOrArrayInt16
		case reflect.Int32:
			return iterateSliceOrArrayInt32
		case reflect.Int64:
			return iterateSliceOrArrayInt64
		case reflect.Int:
			return iterateSliceOrArrayInt
		case reflect.Float32:
			return iterateSliceOrArrayFloat32
		case reflect.Float64:
			return iterateSliceOrArrayFloat64
		case reflect.Bool:
			return iterateSliceOrArrayBool
		default:
			return newSliceOrArrayAsListIterator(&_this.context, t)
		}
	case reflect.Map:
		return newMapIterator(&_this.context, t)
	case reflect.Struct:
		switch t {
		case common.TypeTime:
			return iterateTime
		case common.TypeCompactTime:
			return iterateCompactTime
		case common.TypeDFloat:
			return iterateDecimalFloat
		case common.TypeURL:
			return iterateURL
		case common.TypeBigInt:
			return iterateBigInt
		case common.TypeBigFloat:
			return iterateBigFloat
		case common.TypeBigDecimalFloat:
			return iterateBigDecimal
		case common.TypeMedia:
			return iterateMedia
		case common.TypeNode:
			return iterateNode
		case common.TypeEdge:
			return iterateEdge
		default:
			return newStructIterator(&_this.context, t)
		}
	case reflect.Ptr:
		switch t {
		case common.TypePURL:
			return iteratePURL
		case common.TypePBigInt:
			return iteratePBigInt
		case common.TypePBigFloat:
			return iteratePBigFloat
		case common.TypePBigDecimalFloat:
			return iteratePBigDecimal
		case common.TypePTime:
			return iteratePTime
		case common.TypePCompactTime:
			return iteratePCompactTime
		default:
			return newPointerIterator(&_this.context, t)
		}
	default:
		panic(fmt.Errorf("BUG: Unhandled type %v", t))
	}
}
