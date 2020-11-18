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

// Data events, which are primarily produced by iterators and consumed by builders.
//
// Data events form the backbone of the library. Everything is built upon the
// concepts of producing and consuming data events in order to create and
// interpret Concise Encoding documents.
package events

import (
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type ArrayType uint8

const (
	ArrayTypeInvalid ArrayType = iota
	ArrayTypeBoolean
	ArrayTypeUint8
	ArrayTypeUint16
	ArrayTypeUint32
	ArrayTypeUint64
	ArrayTypeInt8
	ArrayTypeInt16
	ArrayTypeInt32
	ArrayTypeInt64
	ArrayTypeFloat16
	ArrayTypeFloat32
	ArrayTypeFloat64
	ArrayTypeUUID
	ArrayTypeString
	ArrayTypeURI
	ArrayTypeCustomBinary
	ArrayTypeCustomText
)

func (_this ArrayType) String() string {
	return arrayTypeNames[_this]
}

func (_this ArrayType) ElementSize() int {
	return arrayTypeElementSizes[_this]
}

var arrayTypeNames = [...]string{
	ArrayTypeInvalid:      "Invalid",
	ArrayTypeBoolean:      "Boolean",
	ArrayTypeUint8:        "Uint8",
	ArrayTypeUint16:       "Uint16",
	ArrayTypeUint32:       "Uint32",
	ArrayTypeUint64:       "Uint64",
	ArrayTypeInt8:         "Int8",
	ArrayTypeInt16:        "Int16",
	ArrayTypeInt32:        "Int32",
	ArrayTypeInt64:        "Int64",
	ArrayTypeFloat16:      "Float16",
	ArrayTypeFloat32:      "Float32",
	ArrayTypeFloat64:      "Float64",
	ArrayTypeUUID:         "UUID",
	ArrayTypeString:       "String",
	ArrayTypeURI:          "URI",
	ArrayTypeCustomBinary: "Custom Binary",
	ArrayTypeCustomText:   "Custom Text",
}

var arrayTypeElementSizes = [...]int{
	ArrayTypeInvalid:      0,
	ArrayTypeBoolean:      1,
	ArrayTypeUint8:        8,
	ArrayTypeUint16:       16,
	ArrayTypeUint32:       32,
	ArrayTypeUint64:       64,
	ArrayTypeInt8:         8,
	ArrayTypeInt16:        16,
	ArrayTypeInt32:        32,
	ArrayTypeInt64:        64,
	ArrayTypeFloat16:      16,
	ArrayTypeFloat32:      32,
	ArrayTypeFloat64:      64,
	ArrayTypeUUID:         128,
	ArrayTypeString:       8,
	ArrayTypeURI:          8,
	ArrayTypeCustomBinary: 8,
	ArrayTypeCustomText:   8,
}

// DataEventReceiver receives data events (int, string, etc) and performs
// actions based on those events. Generally, this is used to drive complex
// object builders, and also the encoders.
//
// IMPORTANT: DataEventReceiver's methods signal errors via panics, not
// returned errors.
// You must use recover() in code that calls DataEventReceiver methods. The
// recover statements should not be inside of loops, as this causes slow defers.
// (https://go.googlesource.com/proposal/+/refs/heads/master/design/34481-opencoded-defers.md)
// (https://github.com/golang/go/commit/be64a19d99918c843f8555aad580221207ea35bc)
type DataEventReceiver interface {
	// Must be called before any other event.
	OnBeginDocument()
	// Must be called last of all. No other events may be sent after this call.
	OnEndDocument()

	OnVersion(version uint64)
	OnPadding(count int)
	OnNull()
	OnBool(value bool)
	OnTrue()
	OnFalse()
	OnPositiveInt(value uint64)
	OnNegativeInt(value uint64)
	OnInt(value int64)
	OnBigInt(value *big.Int)
	OnFloat(value float64)
	OnBigFloat(value *big.Float)
	OnDecimalFloat(value compact_float.DFloat)
	OnBigDecimalFloat(value *apd.Decimal)
	OnNan(signaling bool)
	OnTime(value time.Time)
	OnCompactTime(value *compact_time.Time)
	OnList()
	OnMap()
	OnMarkup()
	OnMetadata()
	OnComment()
	OnEnd()
	OnMarker()
	OnReference()

	// WARNING: Do not directly store pointers to the data passed via array or
	// UUID handlers! The underlying contents should be considered volatile and
	// likely to change after the method returns (the decoders re-use memory).
	OnUUID(value []byte)
	OnArray(arrayType ArrayType, elementCount uint64, data []uint8)
	OnArrayBegin(arrayType ArrayType)
	OnArrayChunk(length uint64, moreChunksFollow bool)
	OnArrayData(data []byte)
}

// NullEventReceiver receives events and does nothing with them.
type NullEventReceiver struct{}

func NewNullEventReceiver() *NullEventReceiver {
	return &NullEventReceiver{}
}
func (_this *NullEventReceiver) OnBeginDocument()                    {}
func (_this *NullEventReceiver) OnVersion(uint64)                    {}
func (_this *NullEventReceiver) OnPadding(int)                       {}
func (_this *NullEventReceiver) OnNull()                             {}
func (_this *NullEventReceiver) OnBool(bool)                         {}
func (_this *NullEventReceiver) OnTrue()                             {}
func (_this *NullEventReceiver) OnFalse()                            {}
func (_this *NullEventReceiver) OnPositiveInt(uint64)                {}
func (_this *NullEventReceiver) OnNegativeInt(uint64)                {}
func (_this *NullEventReceiver) OnInt(int64)                         {}
func (_this *NullEventReceiver) OnBigInt(*big.Int)                   {}
func (_this *NullEventReceiver) OnFloat(float64)                     {}
func (_this *NullEventReceiver) OnBigFloat(*big.Float)               {}
func (_this *NullEventReceiver) OnDecimalFloat(compact_float.DFloat) {}
func (_this *NullEventReceiver) OnBigDecimalFloat(*apd.Decimal)      {}
func (_this *NullEventReceiver) OnNan(bool)                          {}
func (_this *NullEventReceiver) OnUUID([]byte)                       {}
func (_this *NullEventReceiver) OnTime(time.Time)                    {}
func (_this *NullEventReceiver) OnCompactTime(*compact_time.Time)    {}
func (_this *NullEventReceiver) OnArray(ArrayType, uint64, []byte)   {}
func (_this *NullEventReceiver) OnArrayBegin(ArrayType)              {}
func (_this *NullEventReceiver) OnArrayChunk(uint64, bool)           {}
func (_this *NullEventReceiver) OnArrayData([]byte)                  {}
func (_this *NullEventReceiver) OnList()                             {}
func (_this *NullEventReceiver) OnMap()                              {}
func (_this *NullEventReceiver) OnMarkup()                           {}
func (_this *NullEventReceiver) OnMetadata()                         {}
func (_this *NullEventReceiver) OnComment()                          {}
func (_this *NullEventReceiver) OnEnd()                              {}
func (_this *NullEventReceiver) OnMarker()                           {}
func (_this *NullEventReceiver) OnReference()                        {}
func (_this *NullEventReceiver) OnEndDocument()                      {}
