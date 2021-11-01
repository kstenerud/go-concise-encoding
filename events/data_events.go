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
	ArrayTypeString
	ArrayTypeResourceID
	ArrayTypeResourceIDConcat
	ArrayTypeResourceIDConcat2
	ArrayTypeRemoteRef
	ArrayTypeCustomText
	ArrayTypeCustomBinary
	// Note: Boolean arrays are passed in the CBE boolean array representation,
	// meaning that the low bit is in the least significant position of the
	// first byte in the byte array, and the trailing upper bits of the last
	// byte are cleared.
	ArrayTypeBit
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
	ArrayTypeUID
	ArrayTypeMedia
	ArrayTypeMediaData
	NumArrayTypes
)

func (_this ArrayType) String() string {
	return arrayTypeNames[_this]
}

func (_this ArrayType) ElementSize() int {
	return arrayTypeElementSizes[_this]
}

var arrayTypeNames = [...]string{
	ArrayTypeInvalid:           "Invalid",
	ArrayTypeString:            "String",
	ArrayTypeResourceID:        "ResourceID",
	ArrayTypeResourceIDConcat:  "ResourceIDConcat",
	ArrayTypeResourceIDConcat2: "ResourceIDConcat 2",
	ArrayTypeRemoteRef:         "RemoteRef",
	ArrayTypeCustomText:        "Custom Text",
	ArrayTypeCustomBinary:      "Custom Binary",
	ArrayTypeBit:               "Boolean",
	ArrayTypeUint8:             "Uint8",
	ArrayTypeUint16:            "Uint16",
	ArrayTypeUint32:            "Uint32",
	ArrayTypeUint64:            "Uint64",
	ArrayTypeInt8:              "Int8",
	ArrayTypeInt16:             "Int16",
	ArrayTypeInt32:             "Int32",
	ArrayTypeInt64:             "Int64",
	ArrayTypeFloat16:           "Float16",
	ArrayTypeFloat32:           "Float32",
	ArrayTypeFloat64:           "Float64",
	ArrayTypeUID:               "UID",
	ArrayTypeMedia:             "Media",
	ArrayTypeMediaData:         "MediaData",
}

var arrayTypeElementSizes = [...]int{
	ArrayTypeInvalid:           0,
	ArrayTypeString:            8,
	ArrayTypeResourceID:        8,
	ArrayTypeResourceIDConcat:  8,
	ArrayTypeResourceIDConcat2: 8,
	ArrayTypeRemoteRef:         8,
	ArrayTypeCustomText:        8,
	ArrayTypeCustomBinary:      8,
	ArrayTypeBit:               1,
	ArrayTypeUint8:             8,
	ArrayTypeUint16:            16,
	ArrayTypeUint32:            32,
	ArrayTypeUint64:            64,
	ArrayTypeInt8:              8,
	ArrayTypeInt16:             16,
	ArrayTypeInt32:             32,
	ArrayTypeInt64:             64,
	ArrayTypeFloat16:           16,
	ArrayTypeFloat32:           32,
	ArrayTypeFloat64:           64,
	ArrayTypeUID:               128,
	ArrayTypeMedia:             8,
	ArrayTypeMediaData:         8,
}

// DataEventReceiver receives data events (int, string, etc) and performs
// actions based on those events. Generally, this is used to drive complex
// object builders, and also the encoders.
//
// WARNING: Do not directly store slice data! The underlying contents should be
// considered volatile and likely to change after the method returns (the
// decoders re-use memory).
//
// WARNING: Do not use OnArray or OnStringlikeArray to pass concatenated array
// types! Use OnArrayBegin and then two sets of OnArrayChunk sequences.
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
	OnComment(isMultiline bool, contents []byte)
	OnNil()
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
	OnUID(value []byte)
	OnNan(signaling bool)
	OnTime(value time.Time)
	OnCompactTime(value compact_time.Time)
	OnList()
	OnMap()
	OnEdge()
	OnNode()
	OnMarkup(identifier []byte)
	OnEnd()
	OnMarker(identifier []byte)
	OnReference(identifier []byte)
	OnConstant(identifier []byte)
	OnArray(arrayType ArrayType, elementCount uint64, data []uint8)
	OnArrayBegin(arrayType ArrayType)
	OnArrayChunk(length uint64, moreChunksFollow bool)
	OnArrayData(data []byte)

	// Stringlike array contents are safe to store directly.
	OnStringlikeArray(arrayType ArrayType, data string)
}

// NullEventReceiver receives events and does nothing with them.
type NullEventReceiver struct{}

func NewNullEventReceiver() *NullEventReceiver {
	return &NullEventReceiver{}
}
func (_this *NullEventReceiver) OnBeginDocument()                    {}
func (_this *NullEventReceiver) OnVersion(uint64)                    {}
func (_this *NullEventReceiver) OnComment(bool, []byte)              {}
func (_this *NullEventReceiver) OnPadding(int)                       {}
func (_this *NullEventReceiver) OnNil()                              {}
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
func (_this *NullEventReceiver) OnUID([]byte)                        {}
func (_this *NullEventReceiver) OnTime(time.Time)                    {}
func (_this *NullEventReceiver) OnCompactTime(compact_time.Time)     {}
func (_this *NullEventReceiver) OnArray(ArrayType, uint64, []byte)   {}
func (_this *NullEventReceiver) OnStringlikeArray(ArrayType, string) {}
func (_this *NullEventReceiver) OnArrayBegin(ArrayType)              {}
func (_this *NullEventReceiver) OnArrayChunk(uint64, bool)           {}
func (_this *NullEventReceiver) OnArrayData([]byte)                  {}
func (_this *NullEventReceiver) OnList()                             {}
func (_this *NullEventReceiver) OnMap()                              {}
func (_this *NullEventReceiver) OnEdge()                             {}
func (_this *NullEventReceiver) OnNode()                             {}
func (_this *NullEventReceiver) OnMarkup([]byte)                     {}
func (_this *NullEventReceiver) OnEnd()                              {}
func (_this *NullEventReceiver) OnMarker([]byte)                     {}
func (_this *NullEventReceiver) OnReference([]byte)                  {}
func (_this *NullEventReceiver) OnConstant(_ []byte)                 {}
func (_this *NullEventReceiver) OnEndDocument()                      {}
