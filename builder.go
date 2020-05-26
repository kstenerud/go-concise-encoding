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
	"math/big"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

// TODO: BuilderOptions
// type BuilderOptions struct {
// 	// TODO: Something for lossy conversions?
// 	// TODO: Something for decimal floats?
// }

// Register a specific builder for a type.
// If a builder has already been registered for this type, it will be replaced.
func RegisterBuilderForType(dstType reflect.Type, builder ObjectBuilder) {
	builders.Store(dstType, builder)
}

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
	BuildFromNil(dst reflect.Value)
	BuildFromBool(value bool, dst reflect.Value)
	BuildFromInt(value int64, dst reflect.Value)
	BuildFromUint(value uint64, dst reflect.Value)
	BuildFromBigInt(value *big.Int, dst reflect.Value)
	BuildFromFloat(value float64, dst reflect.Value)
	BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value)
	BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value)
	BuildFromUUID(value []byte, dst reflect.Value)
	BuildFromString(value string, dst reflect.Value)
	BuildFromBytes(value []byte, dst reflect.Value)
	BuildFromURI(value *url.URL, dst reflect.Value)
	BuildFromTime(value time.Time, dst reflect.Value)
	BuildBeginList()
	BuildBeginMap()
	BuildEndContainer()
	BuildFromMarker(id interface{})
	BuildFromReference(id interface{})

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

// ============================================================================

var builders sync.Map

func init() {
	// Pre-cache the most common builders
	for _, t := range keyableTypes {
		getBuilderForType(t)
		getBuilderForType(reflect.PtrTo(t))
		getBuilderForType(reflect.SliceOf(t))
		for _, u := range keyableTypes {
			getBuilderForType(reflect.MapOf(t, u))
		}
		for _, u := range nonKeyableTypes {
			getBuilderForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range nonKeyableTypes {
		getBuilderForType(t)
		getBuilderForType(reflect.PtrTo(t))
		getBuilderForType(reflect.SliceOf(t))
	}
}

func builderPanicBadEvent(builder ObjectBuilder, dstType reflect.Type, containerMsg string) {
	panic(fmt.Errorf(`BUG: %v with type %v cannot respond to event "%v"`, reflect.TypeOf(builder), dstType, containerMsg))
}

func builderPanicCannotConvert(value interface{}, dstType reflect.Type) {
	panic(fmt.Errorf("Cannot convert %v (type %v) to type %v", value, reflect.TypeOf(value), dstType))
}

func builderPanicErrorConverting(value interface{}, dstType reflect.Type, err error) {
	panic(fmt.Errorf("Error converting %v (type %v) to type %v: %v", value, reflect.TypeOf(value), dstType, err))
}

func generateBuilderForType(dstType reflect.Type) ObjectBuilder {
	switch dstType.Kind() {
	case reflect.Bool:
		return newDirectBuilder(dstType)
	case reflect.String:
		return newStringBuilder()
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
			return newDirectPtrBuilder(dstType)
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
		switch dstType {
		case typeTime:
			return newDirectBuilder(dstType)
		case typeURL:
			return newDirectBuilder(dstType)
		case typeDFloat:
			return newDFloatBuilder()
		case typeBigInt:
			return newBigIntBuilder()
		case typeBigFloat:
			return newBigFloatBuilder()
		default:
			return newStructBuilder(dstType)
		}
	case reflect.Ptr:
		switch dstType {
		case typePURL:
			return newDirectPtrBuilder(dstType)
		case typePBigInt:
			return newpBigIntBuilder()
		case typePBigFloat:
			return newPBigFloatBuilder()
		default:
			return newPtrBuilder(dstType)
		}
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
		case typeTime, typeURL, typeDFloat:
			return getBuilderForType(dstType)
		default:
			return newTLContainerBuilder(dstType)
		}
	default:
		return getBuilderForType(dstType)
	}
}
