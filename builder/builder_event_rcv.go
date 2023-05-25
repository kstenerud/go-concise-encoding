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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/common"
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
// If config is nil, default configuration will be used.
func NewBuilderEventReceiver(session *Session, dstType reflect.Type, config *configuration.BuilderConfiguration) *BuilderEventReceiver {
	_this := &BuilderEventReceiver{}
	_this.Init(session, dstType, config)
	return _this
}

// Initialize to build objects of dstType.
// If config is nil, default configuration will be used.
func (_this *BuilderEventReceiver) Init(session *Session, dstType reflect.Type, config *configuration.BuilderConfiguration) {
	_this.context.Init(config,
		dstType,
		session.config.CustomBinaryBuildFunction,
		session.config.CustomTextBuildFunction,
		session.GetBuilderGeneratorForType)

	if dstType == nil {
		dstType = common.TypeInterface
	}
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
	switch v.Type() {
	case reflect.TypeOf(time.Time{}):
		return _this.object.Interface()
	}

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

func (_this *BuilderEventReceiver) OnBeginDocument()       {}
func (_this *BuilderEventReceiver) OnVersion(_ uint64)     {}
func (_this *BuilderEventReceiver) OnPadding()             {}
func (_this *BuilderEventReceiver) OnComment(bool, []byte) {}

func (_this *BuilderEventReceiver) OnNull() {
	_this.context.CurrentBuilder.BuildFromNull(&_this.context, _this.object)
}
func (_this *BuilderEventReceiver) OnBoolean(value bool) {
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
	if value == 0 {
		// Yes, this stupidity around negative zero literals in go is intentional. Blame them.
		const zero = float64(0)
		const negZero = -zero
		_this.OnFloat(negZero)
		return
	}
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
		_this.OnFloat(common.Float64SignalingNan)
	} else {
		_this.OnFloat(common.Float64QuietNan)
	}
}
func (_this *BuilderEventReceiver) OnUID(value []byte) {
	_this.context.CurrentBuilder.BuildFromUID(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnTime(value compact_time.Time) {
	_this.context.CurrentBuilder.BuildFromTime(&_this.context, value, _this.object)
}
func (_this *BuilderEventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.context.CurrentBuilder.BuildFromArray(&_this.context, arrayType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	_this.context.CurrentBuilder.BuildFromStringlikeArray(&_this.context, arrayType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnMedia(mediaType string, value []byte) {
	_this.context.CurrentBuilder.BuildFromMedia(&_this.context, mediaType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnCustomBinary(customType uint64, value []byte) {
	_this.context.CurrentBuilder.BuildFromCustomBinary(&_this.context, customType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnCustomText(customType uint64, value string) {
	_this.context.CurrentBuilder.BuildFromCustomText(&_this.context, customType, value, _this.object)
}
func (_this *BuilderEventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	_this.context.BeginArray(func(ctx *Context) {
		bytes := ctx.chunkedData
		elementCount := common.ByteCountToElementCount(arrayType.ElementSize(), uint64(len(bytes)))
		_this.OnArray(arrayType, elementCount, bytes)
	})
}
func (_this *BuilderEventReceiver) OnMediaBegin(mediaType string) {
	_this.context.BeginArray(func(ctx *Context) {
		bytes := ctx.chunkedData
		ctx.CurrentBuilder.BuildFromMedia(ctx, mediaType, bytes, _this.object)
	})
}
func (_this *BuilderEventReceiver) OnCustomBegin(arrayType events.ArrayType, customType uint64) {
	_this.context.BeginArray(func(ctx *Context) {
		bytes := ctx.chunkedData
		switch arrayType {
		case events.ArrayTypeCustomBinary:
			ctx.CurrentBuilder.BuildFromCustomBinary(ctx, customType, bytes, _this.object)
		case events.ArrayTypeCustomText:
			ctx.CurrentBuilder.BuildFromCustomText(ctx, customType, string(bytes), _this.object)
		default:
			panic(fmt.Errorf("BUG: Cannot handle type %v in OnCustomBegin", arrayType))
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
	_this.context.CurrentBuilder.BuildNewList(&_this.context)
}
func (_this *BuilderEventReceiver) OnMap() {
	_this.context.CurrentBuilder.BuildNewMap(&_this.context)
}
func (_this *BuilderEventReceiver) OnRecordType(id []byte) {
	_this.context.BeginRecordType(id)
}
func (_this *BuilderEventReceiver) OnRecord(id []byte) {
	_this.context.BeginRecord(id)
}
func (_this *BuilderEventReceiver) OnNode() {
	_this.context.CurrentBuilder.BuildNewNode(&_this.context)
}
func (_this *BuilderEventReceiver) OnEdge() {
	_this.context.CurrentBuilder.BuildNewEdge(&_this.context)
}
func (_this *BuilderEventReceiver) OnEndContainer() {
	_this.context.CurrentBuilder.BuildEndContainer(&_this.context)
}
func (_this *BuilderEventReceiver) OnMarker(id []byte) {
	_this.context.BeginMarkerObject(id)
}
func (_this *BuilderEventReceiver) OnReferenceLocal(id []byte) {
	_this.context.CurrentBuilder.BuildFromLocalReference(&_this.context, id)
}
func (_this *BuilderEventReceiver) OnEndDocument() {}
