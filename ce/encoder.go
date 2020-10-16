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

package ce

import (
	"io"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Encoder accepts events and encodes them to a byte stream.
type Encoder interface {
	// Prepare the encoder for encoding. All events will be encoded to writer.
	// PrepareToEncode MUST be called before using the encoder.
	PrepareToEncode(writer io.Writer)

	// Events

	OnBeginDocument()
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
	OnUUID(value []byte)
	OnTime(value time.Time)
	OnCompactTime(value *compact_time.Time)
	// Warning: Do not store a pointer to value! The underlying contents should
	// be considered volatile and likely to change after this method returns!
	OnArray(arrayType events.ArrayType, elementCount uint64, data []byte)
	OnArrayBegin(arrayType events.ArrayType)
	OnArrayChunk(length uint64, moreChunksFollow bool)
	OnArrayData(data []byte)
	OnList()
	OnMap()
	OnMarkup()
	OnMetadata()
	OnComment()
	OnEnd()
	OnMarker()
	OnReference()
}
