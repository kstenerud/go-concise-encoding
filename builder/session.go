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

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"
)

// A builder session holds a cache of known mappings of types to builders.
// It is designed to be cloned so that any user-supplied custom builders exist
// only in their own session, and don't pollute the base mapping and cause
// unintended behavior in codec activity elsewhere in the program.
type Session struct {
	builders map[reflect.Type]ObjectBuilder
	opts     options.BuilderSessionOptions
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
	_this.builders = make(map[reflect.Type]ObjectBuilder)
	opts = opts.WithDefaultsApplied()
	if parent == nil {
		parent = &rootSession
	}
	for k, v := range parent.builders {
		_this.builders[k] = v
	}

	_this.opts = *opts
	for _, t := range _this.opts.CustomBuiltTypes {
		_this.RegisterBuilderForType(t, newCustomBuilder(_this))
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
func (_this *Session) RegisterBuilderForType(dstType reflect.Type, builder ObjectBuilder) {
	_this.builders[dstType] = builder
}

// Get a builder for the specified type. If a registered builder doesn't yet
// exist, a new default builder will be generated and registered.
// This method is thread-safe.
func (_this *Session) GetBuilderForType(dstType reflect.Type) ObjectBuilder {
	if builder, ok := _this.builders[dstType]; ok {
		return builder
	}

	builder, ok := _this.builders[dstType]
	if !ok {
		builder = _this.defaultBuilderForType(dstType)
		_this.builders[dstType] = builder
	}
	builder.InitTemplate(_this)
	return builder.(ObjectBuilder)
}

func (_this *Session) GetCustomBinaryBuildFunction() options.CustomBuildFunction {
	return _this.opts.CustomBinaryBuildFunction
}

func (_this *Session) GetCustomTextBuildFunction() options.CustomBuildFunction {
	return _this.opts.CustomTextBuildFunction
}

func (_this *Session) TryBuildFromCustom(builder ObjectBuilder, arrayType events.ArrayType, value []byte, dst reflect.Value) bool {
	switch arrayType {
	case events.ArrayTypeCustomBinary:
		if err := _this.GetCustomBinaryBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomBinary(builder, value, dst.Type(), err)
		}
		return true
	case events.ArrayTypeCustomText:
		if err := _this.GetCustomTextBuildFunction()(value, dst); err != nil {
			PanicBuildFromCustomText(builder, value, dst.Type(), err)
		}
		return true
	default:
		return false
	}
}

// ============================================================================
// Internal

func (_this *Session) defaultBuilderForType(dstType reflect.Type) ObjectBuilder {
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
			return newUint8ArrayBuilder()
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
			return newDecimalFloatBuilder(common.TypeDFloat)
		case common.TypeBigInt:
			return newBigIntBuilder(common.TypeBigInt)
		case common.TypeBigFloat:
			return newBigFloatBuilder(common.TypeBigFloat)
		case common.TypeBigDecimalFloat:
			return newBigDecimalFloatBuilder(common.TypeBigDecimalFloat)
		default:
			return newStructBuilder(dstType)
		}
	case reflect.Ptr:
		switch dstType {
		case common.TypePURL:
			return newDirectPtrBuilder(dstType)
		case common.TypePBigInt:
			return newPBigIntBuilder(common.TypePBigInt)
		case common.TypePBigFloat:
			return newPBigFloatBuilder(common.TypePBigFloat)
		case common.TypePBigDecimalFloat:
			return newPBigDecimalFloatBuilder(common.TypePBigDecimalFloat)
		case common.TypePCompactTime:
			return newPCompactTimeBuilder()
		default:
			return newPtrBuilder(dstType)
		}
	default:
		panic(fmt.Errorf("BUG: Unhandled type %v", dstType))
	}
}

// The root session caches the most common builders. All sessions inherit
// these cached values.
var rootSession Session

func init() {
	rootSession.Init(nil, nil)

	for _, t := range common.KeyableTypes {
		rootSession.GetBuilderForType(t)
		rootSession.GetBuilderForType(reflect.PtrTo(t))
		rootSession.GetBuilderForType(reflect.SliceOf(t))
		for _, u := range common.KeyableTypes {
			rootSession.GetBuilderForType(reflect.MapOf(t, u))
		}
		for _, u := range common.NonKeyableTypes {
			rootSession.GetBuilderForType(reflect.MapOf(t, u))
		}
	}

	for _, t := range common.NonKeyableTypes {
		rootSession.GetBuilderForType(t)
		rootSession.GetBuilderForType(reflect.PtrTo(t))
		rootSession.GetBuilderForType(reflect.SliceOf(t))
	}
}
