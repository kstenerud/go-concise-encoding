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

package cte

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/buffer"
	"github.com/kstenerud/go-concise-encoding/conversions"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Receives data events, constructing a CTE document from them.
//
// Note: This is a LOW LEVEL API. Error reporting is done via panics. Be sure
// to recover() at an appropriate location when calling Encoder's methods
// directly (with the exception of constructors and initializers, which are not
// designed to panic).
type Encoder struct {
	stream      buffer.StreamingWriteBuffer
	engine      encoderEngine
	arrayEngine arrayEncoderEngine
	opts        options.CTEEncoderOptions
}

// Create a new CTE encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func NewEncoder(opts *options.CTEEncoderOptions) *Encoder {
	_this := &Encoder{}
	_this.Init(opts)
	return _this
}

// Initialize this encoder, which will receive data events and write a document
// to writer. If opts is nil, default options will be used.
func (_this *Encoder) Init(opts *options.CTEEncoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.stream.Init(_this.opts.BufferSize)
	_this.engine.Init(&_this.stream, _this.opts.Indent)
	_this.arrayEngine.Init(&_this.engine, opts)
}

// Prepare the encoder for encoding. All events will be encoded to writer.
// PrepareToEncode MUST be called before using the encoder.
func (_this *Encoder) PrepareToEncode(writer io.Writer) {
	_this.stream.SetWriter(writer)
}

func (_this *Encoder) Reset() {
	_this.stream.Reset()
	_this.engine.Reset()
	_this.arrayEngine.Reset()
}

// ============================================================================

// DataEventReceiver

func (_this *Encoder) OnBeginDocument() {
	// Nothing to do
}

func (_this *Encoder) OnPadding(_ int) {
	// Nothing to do
}

func (_this *Encoder) OnVersion(version uint64) {
	_this.engine.AddVersion(version)
}

func (_this *Encoder) OnNull() {
	_this.engine.BeginObject()
	_this.stream.AddString("@null")
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnBool(value bool) {
	if value {
		_this.OnTrue()
	} else {
		_this.OnFalse()
	}
}

func (_this *Encoder) OnTrue() {
	_this.engine.BeginObject()
	_this.stream.AddString("@true")
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnFalse() {
	_this.engine.BeginObject()
	_this.stream.AddString("@false")
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnInt(value int64) {
	if value >= 0 {
		_this.OnPositiveInt(uint64(value))
	} else {
		_this.OnNegativeInt(uint64(-value))
	}
}

func (_this *Encoder) OnBigInt(value *big.Int) {
	_this.engine.BeginObject()
	_this.stream.AddString(value.String())
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnPositiveInt(value uint64) {
	switch _this.engine.Awaiting {
	case awaitingMarkerID:
		_this.engine.CompleteMarker(value)
	case awaitingReferenceID:
		_this.engine.CompleteReference(value)
	default:
		_this.engine.BeginObject()
		_this.stream.AddFmt("%d", value)
		_this.engine.CompleteObject()
	}
}

func (_this *Encoder) OnNegativeInt(value uint64) {
	_this.engine.BeginObject()
	_this.stream.AddFmt("-%d", value)
	_this.engine.CompleteObject()
}

func (_this *Encoder) onInfinity(isPositive bool) {
	_this.engine.BeginObject()
	if isPositive {
		_this.stream.AddString("@inf")
	} else {
		_this.stream.AddString("-@inf")
	}
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnFloat(value float64) {
	if math.IsNaN(value) {
		_this.OnNan(common.IsSignalingNan(value))
		return
	}
	if math.IsInf(value, 0) {
		_this.onInfinity(value >= 0)
		return
	}

	_this.engine.BeginObject()
	_this.stream.AddFmt("%g", value)
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnBigFloat(value *big.Float) {
	if value.IsInf() {
		_this.onInfinity(value.Sign() >= 0)
		return
	}

	_this.engine.BeginObject()
	_this.stream.AddString(conversions.BigFloatToString(value))
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		_this.OnNan(value.IsSignalingNan())
		return
	}
	if value.IsInfinity() {
		_this.onInfinity(!value.IsNegativeInfinity())
		return
	}

	_this.engine.BeginObject()
	_this.stream.AddString(value.Text('g'))
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		_this.OnNan(false)
	case apd.NaNSignaling:
		_this.OnNan(true)
	case apd.Infinite:
		_this.onInfinity(value.Sign() >= 0)
	default:
		_this.engine.BeginObject()
		_this.stream.AddString(value.Text('g'))
		_this.engine.CompleteObject()
	}
}

func (_this *Encoder) OnNan(signaling bool) {
	_this.engine.BeginObject()
	if signaling {
		_this.stream.AddString("@snan")
	} else {
		_this.stream.AddString("@nan")
	}
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnUUID(v []byte) {
	if len(v) != 16 {
		_this.errorf("expected UUID length 16 but got %v", len(v))
	}
	_this.engine.BeginObject()
	_this.stream.AddFmt("@%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7], v[8], v[9], v[10], v[11], v[12], v[13], v[14], v[15])
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnTime(value time.Time) {
	t, err := compact_time.AsCompactTime(value)
	if err != nil {
		_this.unexpectedError(err, value)
	}
	_this.OnCompactTime(t)
}

func (_this *Encoder) OnCompactTime(value *compact_time.Time) {
	tz := func(v *compact_time.Time) string {
		switch v.TimezoneType {
		case compact_time.TypeZero:
			return ""
		case compact_time.TypeAreaLocation, compact_time.TypeLocal:
			return fmt.Sprintf("/%s", v.AreaLocation)
		case compact_time.TypeLatitudeLongitude:
			return fmt.Sprintf("/%.2f/%.2f", float64(v.LatitudeHundredths)/100, float64(v.LongitudeHundredths)/100)
		default:
			_this.errorf("unknown compact time timezone type %v", value.TimezoneType)
			return ""
		}
	}
	subsec := func(v *compact_time.Time) string {
		if v.Nanosecond == 0 {
			return ""
		}

		str := fmt.Sprintf("%.9f", float64(v.Nanosecond)/float64(1000000000))
		for str[len(str)-1] == '0' {
			str = str[:len(str)-1]
		}
		return str[1:]
	}
	_this.engine.BeginObject()
	switch value.TimeType {
	case compact_time.TypeDate:
		_this.stream.AddFmt("%d-%02d-%02d", value.Year, value.Month, value.Day)
	case compact_time.TypeTime:
		_this.stream.AddFmt("%02d:%02d:%02d%s%s", value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	case compact_time.TypeTimestamp:
		_this.stream.AddFmt("%d-%02d-%02d/%02d:%02d:%02d%s%s",
			value.Year, value.Month, value.Day, value.Hour, value.Minute, value.Second, subsec(value), tz(value))
	default:
		_this.errorf("unknown compact time type %v", value.TimeType)
	}
	_this.engine.CompleteObject()
}

func (_this *Encoder) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	_this.OnArrayBegin(arrayType)
	_this.OnArrayChunk(elementCount, false)
	if elementCount > 0 {
		_this.OnArrayData(value)
	}
}

func (_this *Encoder) OnArrayBegin(arrayType events.ArrayType) {
	_this.arrayEngine.OnArrayBegin(arrayType)
}

func (_this *Encoder) OnArrayChunk(elementCount uint64, moreChunksFollow bool) {
	_this.arrayEngine.OnArrayChunk(elementCount, moreChunksFollow)
}

func (_this *Encoder) OnArrayData(data []byte) {
	_this.arrayEngine.OnArrayData(data)
}

func (_this *Encoder) OnList() {
	_this.engine.BeginContainer(awaitingListFirstItem, "[")
}

func (_this *Encoder) OnMap() {
	_this.engine.BeginContainer(awaitingMapFirstKey, "{")
}

func (_this *Encoder) OnMarkup() {
	_this.engine.BeginMarkup()
}

func (_this *Encoder) OnMetadata() {
	_this.engine.BeginPseudoContainer(awaitingMetaFirstKey, "(")
}

func (_this *Encoder) OnComment() {
	_this.engine.BeginComment()
}

func (_this *Encoder) OnEnd() {
	_this.engine.EndContainer()
}

func (_this *Encoder) OnMarker() {
	_this.engine.BeginMarker()
}

func (_this *Encoder) OnReference() {
	_this.engine.BeginReference()
}

func (_this *Encoder) OnConstant(name []byte, explicitValue bool) {
	_this.engine.BeginObject()
	_this.stream.AddByte('#')
	_this.stream.AddNonemptyBytes(name)
	if explicitValue {
		_this.stream.AddByte(':')
	} else {
		_this.engine.CompleteObject()
	}
}

func (_this *Encoder) OnEndDocument() {
	_this.stream.Flush()
}

func (_this *Encoder) unexpectedError(err error, encoding interface{}) {
	_this.errorf("unexpected error [%v] while encoding %v", err, encoding)
}

func (_this *Encoder) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
