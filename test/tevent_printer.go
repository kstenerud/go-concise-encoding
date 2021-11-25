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

// Test helper code.
package test

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/arrays"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type TEventPrinter struct {
	Next  events.DataEventReceiver
	Print func(event *TEvent)
}

// Create an event receiver that prints the event to stdout.
func NewStdoutTEventPrinter(next events.DataEventReceiver) *TEventPrinter {
	return &TEventPrinter{
		Next: next,
		Print: func(event *TEvent) {
			fmt.Printf("EVENT %v\n", event)
		},
	}
}

func (h *TEventPrinter) OnBeginDocument() {
	h.Print(BD())
	h.Next.OnBeginDocument()
}
func (h *TEventPrinter) OnVersion(version uint64) {
	h.Print(V(version))
	h.Next.OnVersion(version)
}
func (h *TEventPrinter) OnPadding(count int) {
	h.Print(PAD(count))
	h.Next.OnPadding(count)
}
func (h *TEventPrinter) OnComment(isMultiline bool, contents []byte) {
	h.Print(COM(isMultiline, string(contents)))
	h.Next.OnComment(isMultiline, contents)
}
func (h *TEventPrinter) OnNull() {
	h.Print(NULL())
	h.Next.OnNull()
}
func (h *TEventPrinter) OnBool(value bool) {
	h.Print(B(value))
	h.Next.OnBool(value)
}
func (h *TEventPrinter) OnTrue() {
	h.Print(TT())
	h.Next.OnTrue()
}
func (h *TEventPrinter) OnFalse() {
	h.Print(FF())
	h.Next.OnFalse()
}
func (h *TEventPrinter) OnPositiveInt(value uint64) {
	h.Print(PI(value))
	h.Next.OnPositiveInt(value)
}
func (h *TEventPrinter) OnNegativeInt(value uint64) {
	h.Print(NI(value))
	h.Next.OnNegativeInt(value)
}
func (h *TEventPrinter) OnInt(value int64) {
	h.Print(I(value))
	h.Next.OnInt(value)
}
func (h *TEventPrinter) OnBigInt(value *big.Int) {
	h.Print(BI(value))
	h.Next.OnBigInt(value)
}
func (h *TEventPrinter) OnFloat(value float64) {
	if math.IsNaN(value) {
		if common.IsSignalingNan(value) {
			h.Print(SNAN())
		} else {
			h.Print(NAN())
		}
	} else {
		h.Print(F(value))
	}
	h.Next.OnFloat(value)
}
func (h *TEventPrinter) OnBigFloat(value *big.Float) {
	h.Print(NewTEvent(TEventBigFloat, value, nil))
	h.Next.OnBigFloat(value)
}
func (h *TEventPrinter) OnDecimalFloat(value compact_float.DFloat) {
	if value.IsNan() {
		if value.IsSignalingNan() {
			h.Print(SNAN())
		} else {
			h.Print(NAN())
		}
	} else {
		h.Print(DF(value))
	}
	h.Next.OnDecimalFloat(value)
}
func (h *TEventPrinter) OnBigDecimalFloat(value *apd.Decimal) {
	switch value.Form {
	case apd.NaN:
		h.Print(NAN())
	case apd.NaNSignaling:
		h.Print(SNAN())
	default:
		h.Print(BDF(value))
	}
	h.Next.OnBigDecimalFloat(value)
}
func (h *TEventPrinter) OnUID(value []byte) {
	h.Print(UID(value))
	h.Next.OnUID(value)
}
func (h *TEventPrinter) OnTime(value time.Time) {
	h.Print(GT(value))
	h.Next.OnTime(value)
}
func (h *TEventPrinter) OnCompactTime(value compact_time.Time) {
	h.Print(CT(value))
	h.Next.OnCompactTime(value)
}
func (h *TEventPrinter) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	switch arrayType {
	case events.ArrayTypeString:
		h.Print(S(string(value)))
	case events.ArrayTypeResourceID:
		h.Print(RID(string(value)))
	case events.ArrayTypeRemoteRef:
		h.Print(RREF(string(value)))
	case events.ArrayTypeCustomBinary:
		h.Print(CUB(value))
	case events.ArrayTypeCustomText:
		h.Print(CUT(string(value)))
	case events.ArrayTypeBit:
		h.Print(AB(elementCount, value))
	case events.ArrayTypeInt8:
		h.Print(AI8(arrays.BytesToInt8Slice(value)))
	case events.ArrayTypeInt16:
		h.Print(AI16(arrays.BytesToInt16Slice(value)))
	case events.ArrayTypeInt32:
		h.Print(AI32(arrays.BytesToInt32Slice(value)))
	case events.ArrayTypeInt64:
		h.Print(AI64(arrays.BytesToInt64Slice(value)))
	case events.ArrayTypeUint8:
		h.Print(AU8(value))
	case events.ArrayTypeUint16:
		h.Print(AU16(arrays.BytesToUint16Slice(value)))
	case events.ArrayTypeUint32:
		h.Print(AU32(arrays.BytesToUint32Slice(value)))
	case events.ArrayTypeUint64:
		h.Print(AU64(arrays.BytesToUint64Slice(value)))
	case events.ArrayTypeFloat16:
		h.Print(AF16(value))
	case events.ArrayTypeFloat32:
		h.Print(AF32(arrays.BytesToFloat32Slice(value)))
	case events.ArrayTypeFloat64:
		h.Print(AF64(arrays.BytesToFloat64Slice(value)))
	case events.ArrayTypeUID:
		h.Print(AUU(arrays.BytesToUUIDSlice(value)))
	default:
		panic(fmt.Errorf("TODO: TEventPrinter.OnArray: Typed array support for %v", arrayType))
	}
	h.Next.OnArray(arrayType, elementCount, value)
}
func (h *TEventPrinter) OnStringlikeArray(arrayType events.ArrayType, value string) {
	switch arrayType {
	case events.ArrayTypeString:
		h.Print(S(value))
	case events.ArrayTypeResourceID:
		h.Print(RID(value))
	case events.ArrayTypeRemoteRef:
		h.Print(RREF(value))
	case events.ArrayTypeCustomText:
		h.Print(CUT(value))
	default:
		panic(fmt.Errorf("BUG: Array type %v is not stringlike", arrayType))
	}
	h.Next.OnStringlikeArray(arrayType, value)
}
func (h *TEventPrinter) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString:
		h.Print(SB())
	case events.ArrayTypeResourceID:
		h.Print(RB())
	case events.ArrayTypeRemoteRef:
		h.Print(RRB())
	case events.ArrayTypeCustomBinary:
		h.Print(CBB())
	case events.ArrayTypeCustomText:
		h.Print(CTB())
	case events.ArrayTypeBit:
		h.Print(ABB())
	case events.ArrayTypeInt8:
		h.Print(AI8B())
	case events.ArrayTypeInt16:
		h.Print(AI16B())
	case events.ArrayTypeInt32:
		h.Print(AI32B())
	case events.ArrayTypeInt64:
		h.Print(AI64B())
	case events.ArrayTypeUint8:
		h.Print(AU8B())
	case events.ArrayTypeUint16:
		h.Print(AU16B())
	case events.ArrayTypeUint32:
		h.Print(AU32B())
	case events.ArrayTypeUint64:
		h.Print(AU64B())
	case events.ArrayTypeFloat16:
		h.Print(AF16B())
	case events.ArrayTypeFloat32:
		h.Print(AF32B())
	case events.ArrayTypeFloat64:
		h.Print(AF64B())
	case events.ArrayTypeUID:
		h.Print(AUUB())
	case events.ArrayTypeMedia:
		h.Print(MB())
	default:
		panic(fmt.Errorf("TODO: TEventPrinter.OnArrayBegin: Typed array support for %v", arrayType))
	}
	h.Next.OnArrayBegin(arrayType)
}
func (h *TEventPrinter) OnArrayChunk(l uint64, moreChunksFollow bool) {
	h.Print(AC(l, moreChunksFollow))
	h.Next.OnArrayChunk(l, moreChunksFollow)
}
func (h *TEventPrinter) OnArrayData(data []byte) {
	h.Print(AD(data))
	h.Next.OnArrayData(data)
}
func (h *TEventPrinter) OnList() {
	h.Print(L())
	h.Next.OnList()
}
func (h *TEventPrinter) OnMap() {
	h.Print(M())
	h.Next.OnMap()
}
func (h *TEventPrinter) OnMarkup(id []byte) {
	h.Print(MUP(string(id)))
	h.Next.OnMarkup(id)
}
func (h *TEventPrinter) OnEnd() {
	h.Print(E())
	h.Next.OnEnd()
}
func (h *TEventPrinter) OnNode() {
	h.Print(NODE())
	h.Next.OnNode()
}
func (h *TEventPrinter) OnEdge() {
	h.Print(EDGE())
	h.Next.OnEdge()
}
func (h *TEventPrinter) OnMarker(id []byte) {
	h.Print(MARK(string(id)))
	h.Next.OnMarker(id)
}
func (h *TEventPrinter) OnReference(id []byte) {
	h.Print(REF(string(id)))
	h.Next.OnReference(id)
}
func (h *TEventPrinter) OnConstant(name []byte) {
	h.Print(CONST(string(name)))
	h.Next.OnConstant(name)
}
func (h *TEventPrinter) OnEndDocument() {
	h.Print(ED())
	h.Next.OnEndDocument()
}
func (h *TEventPrinter) OnNan(signaling bool) {
	if signaling {
		h.Print(SNAN())
	} else {
		h.Print(NAN())
	}
	h.Next.OnNan(signaling)
}
