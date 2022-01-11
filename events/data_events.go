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
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
)

type ArrayType uint8

const (
	ArrayTypeInvalid ArrayType = iota
	ArrayTypeString
	ArrayTypeResourceID
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
	ArrayTypeInvalid:      "Invalid",
	ArrayTypeString:       "String",
	ArrayTypeResourceID:   "ResourceID",
	ArrayTypeRemoteRef:    "RemoteRef",
	ArrayTypeCustomText:   "Custom Text",
	ArrayTypeCustomBinary: "Custom Binary",
	ArrayTypeBit:          "Boolean",
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
	ArrayTypeUID:          "UID",
	ArrayTypeMedia:        "Media",
	ArrayTypeMediaData:    "MediaData",
}

var arrayTypeElementSizes = [...]int{
	ArrayTypeInvalid:      0,
	ArrayTypeString:       8,
	ArrayTypeResourceID:   8,
	ArrayTypeRemoteRef:    8,
	ArrayTypeCustomText:   8,
	ArrayTypeCustomBinary: 8,
	ArrayTypeBit:          1,
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
	ArrayTypeUID:          128,
	ArrayTypeMedia:        8,
	ArrayTypeMediaData:    8,
}

// DataEventReceiver receives data events (int, string, etc) and performs
// actions based on those events. Generally, this is used to drive complex
// object builders, and also the encoders.
//
// WARNING: Do not directly store slice data! The underlying contents should be
// considered volatile and likely to change after the method returns (the
// decoders re-use the memory buffer).
//
// IMPORTANT: DataEventReceiver's methods signal errors via panics, NOT as returned errors.
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

	OnUID(value []byte)

	OnNan(signaling bool)

	OnTime(value time.Time)

	OnCompactTime(value compact_time.Time)

	OnList()

	OnMap()

	OnEdge()

	OnNode()

	OnMarkup(identifier []byte)

	// Ends the current container (list, map, markup, or node).
	// Note: Markup takes TWO ends (one for the attributes section,
	//       one for the contents section).
	// Note: Edge does NOT use end; it terminates automatically after three objects.
	OnEnd()

	OnMarker(identifier []byte)

	OnReference(identifier []byte)

	OnArray(arrayType ArrayType, elementCount uint64, data []uint8)

	// Stringlike array contents are safe to store directly.
	OnStringlikeArray(arrayType ArrayType, data string)

	OnArrayBegin(arrayType ArrayType)

	OnArrayChunk(length uint64, moreChunksFollow bool)

	// Chunked array data is always passed as bytes. After receiving, you must
	// cast the elements to the correct type yourself.
	//
	// Data will always end on an element boundary (the spec requires this).
	// So for example uint32 data sent via this call will always be a multiple of 4 bytes.
	OnArrayData(data []byte)
}
