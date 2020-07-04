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

package builder

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type BuilderOptions struct {
	// TODO: Currently handled in bigIntMaxBase10Exponent in conversions.go
	FloatToBigIntMaxExponent int
	// TODO: ErrorOnLossyFloatConversion option
	ErrorOnLossyFloatConversion bool
	// TODO: Something for decimal floats?
}

// Register a specific builder for a type.
// If a builder has already been registered for this type, it will be replaced.
func RegisterBuilderForType(dstType reflect.Type, builder ObjectBuilder) {
	builders.Store(dstType, builder)
}

// NewBuilderFor creates a new builder that builds objects of the same type as
// the template object.
func NewBuilderFor(template interface{}, options *BuilderOptions) *RootBuilder {
	rv := reflect.ValueOf(template)
	var t reflect.Type
	if rv.IsValid() {
		t = rv.Type()
	} else {
		t = common.TypeInterface
	}

	return NewRootBuilder(t, options)
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
	BuildFromBigFloat(value *big.Float, dst reflect.Value)
	BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value)
	BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value)
	BuildFromUUID(value []byte, dst reflect.Value)
	BuildFromString(value string, dst reflect.Value)
	BuildFromBytes(value []byte, dst reflect.Value)
	BuildFromURI(value *url.URL, dst reflect.Value)
	BuildFromTime(value time.Time, dst reflect.Value)
	BuildFromCompactTime(value *compact_time.Time, dst reflect.Value)
	BuildBeginList()
	BuildBeginMap()
	BuildEndContainer()
	BuildBeginMarker(id interface{})
	BuildFromReference(id interface{})

	IsContainerOnly() bool

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
	CloneFromTemplate(root *RootBuilder, parent ObjectBuilder, options *BuilderOptions) ObjectBuilder

	SetParent(newParent ObjectBuilder)
}

// ============================================================================

var builders sync.Map

func init() {
	// Pre-cache the most common builders
	for _, t := range common.KeyableTypes {
		getBuilderForType(t)
		getBuilderForType(reflect.PtrTo(t))
		getBuilderForType(reflect.SliceOf(t))
		for _, u := range common.KeyableTypes {
			getBuilderForType(reflect.MapOf(t, u))
		}
		for _, u := range common.NonKeyableTypes {
			getBuilderForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
		getBuilderForType(t)
		getBuilderForType(reflect.PtrTo(t))
		getBuilderForType(reflect.SliceOf(t))
	}
}

func builderPanicBadEvent(builder ObjectBuilder, event string) {
	panic(fmt.Errorf(`BUG: %v cannot respond to %v`, reflect.TypeOf(builder), event))
}

func builderPanicBadEventType(builder ObjectBuilder, dstType reflect.Type, event string) {
	panic(fmt.Errorf(`BUG: %v with type %v cannot respond to %v`, reflect.TypeOf(builder), dstType, event))
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
		switch dstType.Elem().Kind() {
		case reflect.Uint8:
			return newBytesArrayBuilder()
		default:
			return newArrayBuilder(dstType)
		}
	case reflect.Slice:
		switch dstType.Elem().Kind() {
		case reflect.Uint8:
			return newDirectPtrBuilder(dstType)
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
		case common.TypeTime:
			return newTimeBuilder()
		case common.TypeCompactTime:
			return newCompactTimeBuilder()
		case common.TypeURL:
			return newDirectBuilder(dstType)
		case common.TypeDFloat:
			return newDFloatBuilder()
		case common.TypeBigInt:
			return newBigIntBuilder()
		case common.TypeBigFloat:
			return newBigFloatBuilder()
		case common.TypeBigDecimalFloat:
			return newBigDecimalFloatBuilder()
		default:
			return newStructBuilder(dstType)
		}
	case reflect.Ptr:
		switch dstType {
		case common.TypePURL:
			return newDirectPtrBuilder(dstType)
		case common.TypePBigInt:
			return newPBigIntBuilder()
		case common.TypePBigFloat:
			return newPBigFloatBuilder()
		case common.TypePBigDecimalFloat:
			return newPBigDecimalFloatBuilder()
		case common.TypePCompactTime:
			return newPCompactTimeBuilder()
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

func applyDefaultBuilderOptions(original *BuilderOptions) *BuilderOptions {
	var options BuilderOptions
	if original != nil {
		options = *original
	}
	return &options
}
