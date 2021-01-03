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

// Builders consume events to produce objects.
//
// Builders respond to builder events in order to build arbitrary objects.
// Generally, they take template types and generate objects of those types.
package builder

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// A builder session holds a cache of known mappings of types to builders.
// It is designed to be cloned so that any user-supplied custom builders exist
// only in their own session, and don't pollute the base mapping and cause
// unintended behavior in codec activity elsewhere in the program.
type Session struct {
	builderGenerators sync.Map
	opts              options.BuilderSessionOptions
}

// Start a new builder session. It will inherit the builders of its parent.
// If parent is nil, it will inherit from the root session, which has builders
// for all basic go types.
// If opts is nil, default options will be used.
func NewSession(parent *Session, opts *options.BuilderSessionOptions) *Session {
	_this := &Session{}
	_this.Init(parent, opts)
	return _this
}

// Initialize a builder session. It will inherit the builders of its parent.
// If parent is nil, it will inherit from the root session, which has builders
// for all basic go types.
// If opts is nil, default options will be used.
func (_this *Session) Init(parent *Session, opts *options.BuilderSessionOptions) {
	opts = opts.WithDefaultsApplied()
	if parent == nil {
		parent = &rootSession
	}
	if parent != _this {
		parent.builderGenerators.Range(func(key, value interface{}) bool {
			_this.builderGenerators.Store(key, value)
			return true
		})
	}

	_this.opts = *opts
	for _, t := range _this.opts.CustomBuiltTypes {
		_this.RegisterBuilderGeneratorForType(t, generateCustomBuilder)
	}
}

// NewBuilderFor creates a new builder that builds objects of the same type as
// the template object.
// If template is nil, a generic interface type will be used.
// If opts is nil, default options will be used.
func (_this *Session) NewBuilderFor(template interface{}, opts *options.BuilderOptions) *RootBuilder {
	rv := reflect.ValueOf(template)
	var t reflect.Type
	if rv.IsValid() {
		t = rv.Type()
	} else {
		t = common.TypeInterface
	}

	return NewRootBuilder(_this, t, opts)
}

// Register a specific builder for a type.
// If a builder has already been registered for this type, it will be replaced.
// This method is thread-safe.
func (_this *Session) RegisterBuilderGeneratorForType(dstType reflect.Type, builderGenerator BuilderGenerator) {
	_this.builderGenerators.Store(dstType, builderGenerator)
}

// Get a builder generator for the specified type. If a registered generator
// doesn't yet exist, a new default generator will be generated and registered.
// This method is thread-safe.
func (_this *Session) GetBuilderGeneratorForType(dstType reflect.Type) BuilderGenerator {
	storedIterator, ok := _this.builderGenerators.Load(dstType)
	if ok {
		return storedIterator.(BuilderGenerator)
	}

	var wg sync.WaitGroup
	var builderGenerator BuilderGenerator

	wg.Add(1)
	storedBuilderGenerator, loaded := _this.builderGenerators.LoadOrStore(dstType, BuilderGenerator(func() ObjectBuilder {
		wg.Wait()
		return builderGenerator()
	}))
	if loaded {
		return storedBuilderGenerator.(BuilderGenerator)
	}

	builderGenerator = _this.defaultBuilderGeneratorForType(dstType)
	wg.Done()
	_this.builderGenerators.Store(dstType, builderGenerator)
	return builderGenerator
}

// ============================================================================
// Internal

func (_this *Session) defaultBuilderGeneratorForType(dstType reflect.Type) BuilderGenerator {
	switch dstType.Kind() {
	case reflect.Bool:
		return generateDirectBuilder
	case reflect.String:
		return generateStringBuilder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return generateIntBuilder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return generateUintBuilder
	case reflect.Float32, reflect.Float64:
		return generateFloatBuilder
	case reflect.Interface:
		return generateInterfaceBuilder
	case reflect.Array:
		switch dstType.Elem().Kind() {
		case reflect.Uint8:
			return newUint8ArrayBuilderGenerator(dstType)
		case reflect.Uint16:
			return newUint16ArrayBuilderGenerator(dstType)
		case reflect.Uint32:
			return newUint32ArrayBuilderGenerator(dstType)
		case reflect.Uint64:
			return newUint64ArrayBuilderGenerator(dstType)
		case reflect.Int8:
			return newInt8ArrayBuilderGenerator(dstType)
		case reflect.Int16:
			return newInt16ArrayBuilderGenerator(dstType)
		case reflect.Int32:
			return newInt32ArrayBuilderGenerator(dstType)
		case reflect.Int64:
			return newInt64ArrayBuilderGenerator(dstType)
		case reflect.Float32:
			return newFloat32ArrayBuilderGenerator(dstType)
		case reflect.Float64:
			return newFloat64ArrayBuilderGenerator(dstType)
		default:
			return newArrayBuilderGenerator(_this.GetBuilderGeneratorForType, dstType)
		}
	case reflect.Slice:
		switch dstType.Elem().Kind() {
		case reflect.Uint8:
			return generateUint8SliceBuilder
		case reflect.Uint16:
			return generateUint16SliceBuilder
		case reflect.Uint32:
			return generateUint32SliceBuilder
		case reflect.Uint64:
			return generateUint64SliceBuilder
		case reflect.Int8:
			return generateInt8SliceBuilder
		case reflect.Int16:
			return generateInt16SliceBuilder
		case reflect.Int32:
			return generateInt32SliceBuilder
		case reflect.Int64:
			return generateInt64SliceBuilder
		case reflect.Float32:
			return generateFloat32SliceBuilder
		case reflect.Float64:
			return generateFloat64SliceBuilder
		default:
			return newSliceBuilderGenerator(_this.GetBuilderGeneratorForType, dstType)
		}
	case reflect.Map:
		return newMapBuilderGenerator(_this.GetBuilderGeneratorForType, dstType)
	case reflect.Struct:
		switch dstType {
		case common.TypeTime:
			return generateTimeBuilder
		case common.TypeCompactTime:
			return generateCompactTimeBuilder
		case common.TypeURL:
			return newUrlBuilderGenerator()
		case common.TypeDFloat:
			return generateDecimalFloatBuilder
		case common.TypeBigInt:
			return generateBigIntBuilder
		case common.TypeBigFloat:
			return generateBigFloatBuilder
		case common.TypeBigDecimalFloat:
			return generateBigDecimalFloatBuilder
		default:
			return newStructBuilderGenerator(_this.GetBuilderGeneratorForType, dstType)
		}
	case reflect.Ptr:
		switch dstType {
		case common.TypePURL:
			return newPUrlBuilderGenerator()
		case common.TypePBigInt:
			return generatePBigIntBuilder
		case common.TypePBigFloat:
			return generatePBigFloatBuilder
		case common.TypePBigDecimalFloat:
			return generatePBigDecimalFloatBuilder
		case common.TypePCompactTime:
			return generatePCompactTimeBuilder
		default:
			return newPtrBuilderGenerator(_this.GetBuilderGeneratorForType, dstType)
		}
	default:
		panic(fmt.Errorf("BUG: Unhandled type %v", dstType))
	}
}

// The root session caches the most common builders. All sessions inherit
// these cached values.
var rootSession Session
var interfaceSliceBuilderGenerator BuilderGenerator
var interfaceMapBuilderGenerator BuilderGenerator
var listToUint8SliceGenerator BuilderGenerator
var listToUint16SliceGenerator BuilderGenerator
var listToUint32SliceGenerator BuilderGenerator
var listToUint64SliceGenerator BuilderGenerator
var listToInt8SliceGenerator BuilderGenerator
var listToInt16SliceGenerator BuilderGenerator
var listToInt32SliceGenerator BuilderGenerator
var listToInt64SliceGenerator BuilderGenerator
var listToFloat32SliceGenerator BuilderGenerator
var listToFloat64SliceGenerator BuilderGenerator

func init() {
	rootSession.Init(nil, nil)

	for _, t := range common.KeyableTypes {
		rootSession.GetBuilderGeneratorForType(t)
		rootSession.GetBuilderGeneratorForType(reflect.PtrTo(t))
		rootSession.GetBuilderGeneratorForType(reflect.SliceOf(t))
		rootSession.GetBuilderGeneratorForType(reflect.SliceOf(reflect.PtrTo(t)))
		for _, u := range common.KeyableTypes {
			rootSession.GetBuilderGeneratorForType(reflect.MapOf(t, u))
		}
		for _, u := range common.NonKeyableTypes {
			rootSession.GetBuilderGeneratorForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
		rootSession.GetBuilderGeneratorForType(t)
		rootSession.GetBuilderGeneratorForType(reflect.PtrTo(t))
		rootSession.GetBuilderGeneratorForType(reflect.SliceOf(t))
	}

	interfaceMapBuilderGenerator = rootSession.GetBuilderGeneratorForType(common.TypeInterfaceMap)
	interfaceSliceBuilderGenerator = rootSession.GetBuilderGeneratorForType(common.TypeInterfaceSlice)
	listToUint8SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]uint8{}))
	listToUint16SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]uint16{}))
	listToUint32SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]uint32{}))
	listToUint64SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]uint64{}))
	listToInt8SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]int8{}))
	listToInt16SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]int16{}))
	listToInt32SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]int32{}))
	listToInt64SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]int64{}))
	listToFloat32SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]float32{}))
	listToFloat64SliceGenerator = newSliceBuilderGenerator(rootSession.GetBuilderGeneratorForType, reflect.TypeOf([]float64{}))
}
