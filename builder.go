package concise_encoding

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"time"
)

// NewBuilderFor creates a new builder that builds objects of the same type as
// the template object.
func NewBuilderFor(template interface{}) *RootBuilder {
	rv := reflect.ValueOf(template)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return newRootBuilder(rv.Type())
}

// ObjectBuilder responds to external events to progressively build an object.
type ObjectBuilder interface {
	// External data and structure events
	Nil(dst reflect.Value)
	Bool(value bool, dst reflect.Value)
	Int(value int64, dst reflect.Value)
	Uint(value uint64, dst reflect.Value)
	Float(value float64, dst reflect.Value)
	String(value string, dst reflect.Value)
	Bytes(value []byte, dst reflect.Value)
	URI(value *url.URL, dst reflect.Value)
	Time(value time.Time, dst reflect.Value)
	List()
	Map()
	End()
	Marker(id interface{})
	Reference(id interface{})

	// Prepare this builder for storing list contents, ultimately followed by End()
	PrepareForListContents()

	// Prepare this builder for storing map contents, ultimately followed by End()
	PrepareForMapContents()

	// Notify that a child builder has finished building a container
	NotifyChildContainerFinished(container reflect.Value)

	// Called after the builder template is saved to cache but before use, so
	// that lookups succeed on cyclic builder references
	PostCacheInitBuilder()

	// Clone from this builder as a template, adding contextual data
	CloneFromTemplate(root *RootBuilder, parent ObjectBuilder) ObjectBuilder
}

var builders sync.Map

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

	// Pre-cache the most common builders
	for _, t := range types {
		getBuilderForType(t)
		getBuilderForType(reflect.PtrTo(t))
		getBuilderForType(reflect.SliceOf(t))
		for _, u := range types {
			getBuilderForType(reflect.MapOf(t, u))
		}
	}
}

func builderPanicBadEvent(builder ObjectBuilder, dstType reflect.Type, containerMsg string) {
	panic(fmt.Sprintf(`%v with type %v cannot respond to event "%v"`, reflect.TypeOf(builder), dstType, containerMsg))
}

func builderPanicCannotConvert(value interface{}, dstType reflect.Type) {
	panic(fmt.Errorf("[%v] cannot be safely converted to %v", value, dstType))
}

func generateBuilderForType(dstType reflect.Type) ObjectBuilder {
	switch dstType.Kind() {
	case reflect.Bool, reflect.String:
		return newBasicBuilder(dstType)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newIntBuilder(dstType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return newUintBuilder(dstType)
	case reflect.Float32, reflect.Float64:
		return newFloatBuilder(dstType)
	case reflect.Interface:
		return newInterfaceBuilder()
	case reflect.Array:
		return newArrayBuilder(dstType)
	case reflect.Slice:
		switch dstType.Elem().Kind() {
		case reflect.Uint8:
			return newBytesBuilder()
		case reflect.Interface:
			return newIntfSliceBuilder()
		default:
			return newSliceBuilder(dstType)
		}
	case reflect.Map:
		if dstType.Elem().Kind() == reflect.Interface && dstType.Key().Kind() == reflect.Interface {
			return newIntfIntfMapBuilder()
		}
		return newMapBuilder(dstType)
	case reflect.Struct:
		if dstType == timeType {
			return newBasicBuilder(dstType)
		}
		if dstType == urlType {
			return newURLBuilder()
		}
		return newStructBuilder(dstType)
	case reflect.Ptr:
		if dstType == pURLType {
			return newPURLBuilder()
		}
		return newPtrBuilder(dstType)
	default:
		panic(fmt.Errorf("BUG: Unhandled type %v", dstType))
	}
}

func getBuilderForType(dstType reflect.Type) ObjectBuilder {
	if builder, ok := builders.Load(dstType); ok {
		return builder.(ObjectBuilder)
	}

	builder, _ := builders.LoadOrStore(dstType, generateBuilderForType(dstType))
	builder.(ObjectBuilder).PostCacheInitBuilder()
	return builder.(ObjectBuilder)
}

func getTopLevelBuilderForType(dstType reflect.Type) ObjectBuilder {
	switch dstType.Kind() {
	case reflect.Slice:
		if dstType.Elem().Kind() == reflect.Uint8 {
			return getBuilderForType(dstType)
		} else {
			return newTLContainerBuilder(dstType)
		}
	case reflect.Array, reflect.Map:
		return newTLContainerBuilder(dstType)
	case reflect.Struct:
		switch dstType {
		case timeType:
			return getBuilderForType(dstType)
		case urlType:
			return getBuilderForType(dstType)
		default:
			return newTLContainerBuilder(dstType)
		}
	default:
		return getBuilderForType(dstType)
	}
}
