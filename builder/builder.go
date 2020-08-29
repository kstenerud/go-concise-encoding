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
	"math/big"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
)

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
	BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value)
	BuildFromTime(value time.Time, dst reflect.Value)
	BuildFromCompactTime(value *compact_time.Time, dst reflect.Value)
	BuildBeginList()
	BuildBeginMap()
	BuildEndContainer()
	BuildBeginMarker(id interface{})
	BuildFromReference(id interface{})

	// Prepare this builder for storing list contents, ultimately followed by End()
	PrepareForListContents()

	// Prepare this builder for storing map contents, ultimately followed by End()
	PrepareForMapContents()

	// Notify that a child builder has finished building a container
	NotifyChildContainerFinished(container reflect.Value)

	// Initialize as a template. This is called after the template has been
	// cached in order to support recursive object graphs.
	InitTemplate(session *Session)

	// Create an instance using this object as a template.
	// opts must not be nil.
	NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder

	SetParent(newParent ObjectBuilder)
}

// Fetches a builder for the specified type, building and caching one as needed.
type FetchBuilder func(reflect.Type) ObjectBuilder

// ============================================================================
// Error reporting

// Report that a builder was given an event that it can't handle.
// This indicates a bug in the implementation.
func PanicBadEvent(builder ObjectBuilder, eventFmt string, args ...interface{}) {
	event := fmt.Sprintf(eventFmt, args...)
	panic(fmt.Errorf(`BUG: %v cannot respond to %v`, reflect.TypeOf(builder), event))
}

// Report that a builder with the specified type was given an event that it can't handle.
// This indicates a bug in the implementation.
func PanicBadEventWithType(builder ObjectBuilder, dstType reflect.Type, eventFmt string, args ...interface{}) {
	event := fmt.Sprintf(eventFmt, args...)
	panic(fmt.Errorf(`BUG: %v with type %v cannot respond to %v`, reflect.TypeOf(builder), dstType, event))
}

// Report that a builder couldn't convert between types. This can happen if
// source values are out of range, or incompatible with the destination type.
func PanicCannotConvert(value interface{}, dstType reflect.Type) {
	panic(fmt.Errorf("cannot convert %v (type %v) to type %v", describe.D(value), reflect.TypeOf(value), dstType))
}

// Report that a builder couldn't convert between types. This can happen if
// source values are out of range, or incompatible with the destination type.
func PanicCannotConvertRV(value reflect.Value, dstType reflect.Type) {
	panic(fmt.Errorf("cannot convert %v (type %v) to type %v", describe.D(value), value.Type(), dstType))
}

// Report that an error occurred while converting between types.
// This normally indicates a bug.
func PanicErrorConverting(value interface{}, dstType reflect.Type, err error) {
	panic(fmt.Errorf("error converting %v (type %v) to type %v: %v", describe.D(value), reflect.TypeOf(value), dstType, err))
}

// Report that an error occurred while building from custom binary data.
// This normally indicates a bug in your custom builder.
func PanicBuildFromCustomBinary(builder ObjectBuilder, src []byte, dstType reflect.Type, err error) {
	panic(fmt.Errorf("error converting custom binary data %v to type %v (via %v): %v", describe.D(src), dstType, reflect.TypeOf(builder), err))
}

// Report that an error occurred while building from custom text data.
// This normally indicates a bug in your custom builder.
func PanicBuildFromCustomText(builder ObjectBuilder, src []byte, dstType reflect.Type, err error) {
	panic(fmt.Errorf("error converting custom text data [%v] to type %v (via %v): %v", string(src), dstType, reflect.TypeOf(builder), err))
}

func nameOf(x interface{}) string {
	return fmt.Sprintf("%v", reflect.TypeOf(x))
}
