package concise_encoding

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"time"
)

// Iterate over an object (recursively), calling the eventHandler as data is
// encountered. If useReferences is true, it will also look for duplicate
// pointers to data, generating marker and reference events rather than walking
// the object again. This is useful for cyclic or recursive data structures.
func IterateObject(value interface{}, useReferences bool, eventHandler ConciseEncodingEventHandler) {
	iter := NewRootObjectIterator(useReferences, eventHandler)
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
	types := []reflect.Type{
		reflect.TypeOf((*bool)(nil)).Elem(),
		reflect.TypeOf((*int)(nil)).Elem(),
		reflect.TypeOf((*int8)(nil)).Elem(),
		reflect.TypeOf((*int16)(nil)).Elem(),
		reflect.TypeOf((*int32)(nil)).Elem(),
		reflect.TypeOf((*int64)(nil)).Elem(),
		reflect.TypeOf((*uint)(nil)).Elem(),
		reflect.TypeOf((*uint8)(nil)).Elem(),
		reflect.TypeOf((*uint16)(nil)).Elem(),
		reflect.TypeOf((*uint32)(nil)).Elem(),
		reflect.TypeOf((*uint64)(nil)).Elem(),
		reflect.TypeOf((*float32)(nil)).Elem(),
		reflect.TypeOf((*float64)(nil)).Elem(),
		reflect.TypeOf((*string)(nil)).Elem(),
		reflect.TypeOf((*url.URL)(nil)).Elem(),
		reflect.TypeOf((*time.Time)(nil)).Elem(),
		reflect.TypeOf((*interface{})(nil)).Elem(),
	}

	// Pre-cache the most common iterators
	for _, t := range types {
		getIteratorForType(t)
		getIteratorForType(reflect.PtrTo(t))
		getIteratorForType(reflect.SliceOf(t))
		for _, u := range types {
			getIteratorForType(reflect.MapOf(t, u))
		}
	}
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
		case timeType:
			return newTimeIterator()
		case urlType:
			return newURLIterator()
		default:
			return newStructIterator(t)
		}
	case reflect.Ptr:
		switch t {
		case pURLType:
			return newPURLIterator()
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
