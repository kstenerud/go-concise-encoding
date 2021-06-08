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
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// BuilderEventReceiver adapts DataEventReceiver to ObjectBuilder, passing
// event contents to builder objects.
//
// Use GetBuiltObject() to fetch the final result.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling this struct's methods
// directly (with the exception of constructors, initializers, and
// GetBuildObject(), which are not designed to panic).
type BuilderEventReceiver struct {
	context Context
	object  reflect.Value
}

// Create a new builder event receiver that can build objects of dstType.
// If opts is nil, default options will be used.
func NewBuilder(session *Session, dstType reflect.Type, opts *options.BuilderOptions) *BuilderEventReceiver {
	_this := &BuilderEventReceiver{}
	_this.Init(session, dstType, opts)
	return _this
}

// Initialize to build objects of dstType.
// If opts is nil, default options will be used.
func (_this *BuilderEventReceiver) Init(session *Session, dstType reflect.Type, opts *options.BuilderOptions) {
	_this.context.Init(opts,
		dstType,
		session.opts.CustomBinaryBuildFunction,
		session.opts.CustomTextBuildFunction,
		session.GetBuilderGeneratorForType)

	_this.object = reflect.New(dstType).Elem()
	generator := session.GetBuilderGeneratorForType(dstType)
	_this.context.StackBuilder(newTopLevelBuilder(generator, func(value reflect.Value) {
		_this.object = value
	}))
}

func (_this *BuilderEventReceiver) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.context.CurrentBuilder)
}

// Get the object that was built after using this as a DataEventReceiver.
func (_this *BuilderEventReceiver) GetBuiltObject() interface{} {
	// TODO: Verify this behavior
	if !_this.object.IsValid() {
		return nil
	}

	v := _this.object
	switch v.Kind() {
	case reflect.Struct, reflect.Array:
		return v.Addr().Interface()
	default:
		if _this.object.CanInterface() {
			return _this.object.Interface()
		}
		return nil
	}
}

// ---------------------------
// DataEventReceiver Callbacks
// ---------------------------

func (_this *BuilderEventReceiver) OnBeginDocument()   {}
func (_this *BuilderEventReceiver) OnVersion(_ uint64) {}
func (_this *BuilderEventReceiver) OnPadding(_ int)    {}
func (_this *BuilderEventReceiver) OnNil() {
	_this.context.CurrentBuilder.BuildFromNil(&_this.context, _this.object)
}
func (_this *BuilderEventReceiver) OnNA() {
	_this.context.CurrentBuilder.BuildFromNil(&_this.context, _this.object)
	_this.context.IgnoreNext()
}
func (_this *BuilderEventReceiver) OnBool(value bool) {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnTrue() {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, true, _this.object)
}
func (_this *BuilderEventReceiver) OnFalse() {
	_this.context.CurrentBuilder.BuildFromBool(&_this.context, false, _this.object)
}
func (_this *BuilderEventReceiver) OnPositiveInt(value uint64) {
	_this.context.CurrentBuilder.BuildFromUint(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnNegativeInt(value uint64) {
	if value <= 0x7fffffffffffffff {
		_this.OnInt(-int64(value))
		return
	}
	bi := big.Int{}
	bi.SetUint64(value)
	_this.OnBigInt(&bi)
}
func (_this *BuilderEventReceiver) OnInt(value int64) {
	_this.context.CurrentBuilder.BuildFromInt(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnBigInt(value *big.Int) {
	_this.context.CurrentBuilder.BuildFromBigInt(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnFloat(value float64) {
	_this.context.CurrentBuilder.BuildFromFloat(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnBigFloat(value *big.Float) {
	_this.context.CurrentBuilder.BuildFromBigFloat(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	_this.context.CurrentBuilder.BuildFromDecimalFloat(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	_this.context.CurrentBuilder.BuildFromBigDecimalFloat(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnNan(signaling bool) {
	if signaling {
		_this.OnFloat(common.SignalingNan)
	} else {
		_this.OnFloat(common.QuietNan)
	}
}
func (_this *BuilderEventReceiver) OnUID(value []byte) {
	_this.context.CurrentBuilder.BuildFromUID(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnTime(value time.Time) {
	_this.context.CurrentBuilder.BuildFromTime(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnCompactTime(value compact_time.Time) {
	_this.context.CurrentBuilder.BuildFromCompactTime(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.CurrentBuilder.BuildFromArray(&_this.context, arrayType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.CurrentBuilder.BuildFromStringlikeArray(&_this.context, arrayType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.BeginArray(func(ctx *Context) {
		switch arrayType {
		case events.ArrayTypeResourceIDConcat:
			ctx.ContinueMultiComponentArray(func(ctx *Context) {
				bytes := ctx.chunkedData
				elementCount := common.ByteCountToElementCount(arrayType.ElementSize(), uint64(len(bytes)))
				_this.OnArray(arrayType, elementCount, bytes)
			})
		case events.ArrayTypeMedia:
			mediaType := string(ctx.chunkedData)
			ctx.chunkedData = ctx.chunkedData[:0]
			ctx.ContinueMultiComponentArray(func(ctx *Context) {
				data := ctx.chunkedData
				ctx.CurrentBuilder.BuildFromMedia(ctx, mediaType, data, _this.object)
			})
		default:
			bytes := ctx.chunkedData
			elementCount := common.ByteCountToElementCount(arrayType.ElementSize(), uint64(len(bytes)))
			_this.OnArray(arrayType, elementCount, bytes)
		}
	})
}
func (_this *BuilderEventReceiver) OnArrayChunk(length uint64, moreChunksFollow bool) {
	_this.context.BeginArrayChunk(length, moreChunksFollow)
}
func (_this *BuilderEventReceiver) OnArrayData(data []byte) {
	_this.context.AddArrayData(data)
}
func (_this *BuilderEventReceiver) OnList() {
	_this.context.CurrentBuilder.BuildInitiateList(&_this.context)
}
func (_this *BuilderEventReceiver) OnMap() {
	_this.context.CurrentBuilder.BuildInitiateMap(&_this.context)
}
func (_this *BuilderEventReceiver) OnMarkup(id []byte) {
	_this.context.CurrentBuilder.BuildInitiateMarkup(&_this.context, id)
}
func (_this *BuilderEventReceiver) OnComment() {
	_this.context.IgnoreNext()
}
func (_this *BuilderEventReceiver) OnEnd() {
	_this.context.CurrentBuilder.BuildEndContainer(&_this.context)
}
func (_this *BuilderEventReceiver) OnRelationship() {
	_this.context.BeginRelationship()
}
func (_this *BuilderEventReceiver) OnMarker(id []byte) {
	_this.context.BeginMarkerObject(id)
}
func (_this *BuilderEventReceiver) OnReference(id []byte) {
	_this.context.CurrentBuilder.BuildFromReference(&_this.context, id)
}
func (_this *BuilderEventReceiver) OnRIDReference() {
	panic("TODO: BuilderEventReceiver.OnRIDReference")
}
func (_this *BuilderEventReceiver) OnConstant(name []byte) {
	panic(fmt.Errorf("Cannot build from constant (%s)", string(name)))
}
func (_this *BuilderEventReceiver) OnEndDocument() {}
