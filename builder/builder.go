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

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-describe"
)

// ObjectBuilder responds to external events to progressively build an object.
type ObjectBuilder interface {

	// External data and structure events
	BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value
	BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value
	BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value
	BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value
	BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value
	BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value
	BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value
	BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value
	BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value
	BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value
	BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value
	BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value
	BuildFromCompactTime(ctx *Context, value *compact_time.Time, dst reflect.Value) reflect.Value
	BuildFromReference(ctx *Context, id interface{})
	BuildConcatenate(ctx *Context)

	// Signals that a new source container has begun.
	// This gets triggered from a data event.
	BuildInitiateList(ctx *Context)
	BuildInitiateMap(ctx *Context)

	// Signals that the source container is finished
	// This gets triggered from a data event.
	BuildEndContainer(ctx *Context)

	// Tells this builder to create a new container to receive the source container's objects.
	// This gets called by the parent builder.
	BuildBeginListContents(ctx *Context)
	BuildBeginMapContents(ctx *Context)

	// Notify that a child builder has finished building a container.
	// This gets triggered from the child builder when the container has ended and the builder unstacked.
	NotifyChildContainerFinished(ctx *Context, container reflect.Value)
}

type BuilderGenerator func(ctx *Context) ObjectBuilder
type BuilderGeneratorGetter func(reflect.Type) BuilderGenerator

// ============================================================================
// Error reporting

// Report that a builder was given an event that it can't handle.
// This indicates a bug in the implementation.
func PanicBadEvent(builder ObjectBuilder, eventFmt string, args ...interface{}) {
	panic(fmt.Errorf(`BUG: %v cannot respond to %v`, reflect.TypeOf(builder), fmt.Sprintf(eventFmt, args...)))
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
