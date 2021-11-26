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
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/arrays"
)

// Event receiver receives data events and stores them to an array which can be
// inspected, printed, or played back.
type TEventStore struct {
	Events   []*TEvent
	receiver events.DataEventReceiver
}

func NewTEventStore(receiver events.DataEventReceiver) *TEventStore {
	return &TEventStore{
		Events:   make([]*TEvent, 0, 1024),
		receiver: receiver,
	}
}
func (h *TEventStore) add(event *TEvent) {
	h.Events = append(h.Events, event)
}
func (h *TEventStore) OnVersion(version uint64) {
	h.add(V(version))
	h.receiver.OnVersion(version)
}
func (h *TEventStore) OnPadding(count int) {
	h.add(PAD(count))
	h.receiver.OnPadding(count)
}
func (h *TEventStore) OnComment(isMultiline bool, contents []byte) {
	h.add(COM(isMultiline, string(contents)))
	h.receiver.OnComment(isMultiline, contents)
}
func (h *TEventStore) OnNull() {
	h.add(NULL())
	h.receiver.OnNull()
}
func (h *TEventStore) OnBool(value bool) {
	h.add(B(value))
	h.receiver.OnBool(value)
}
func (h *TEventStore) OnTrue() {
	h.add(TT())
	h.receiver.OnTrue()
}
func (h *TEventStore) OnFalse() {
	h.add(FF())
	h.receiver.OnFalse()
}
func (h *TEventStore) OnPositiveInt(value uint64) {
	h.add(PI(value))
	h.receiver.OnPositiveInt(value)
}
func (h *TEventStore) OnNegativeInt(value uint64) {
	h.add(NI(value))
	h.receiver.OnNegativeInt(value)
}
func (h *TEventStore) OnInt(value int64) {
	h.add(I(value))
	h.receiver.OnInt(value)
}
func (h *TEventStore) OnBigInt(value *big.Int) {
	h.add(BI(value))
	h.receiver.OnBigInt(value)
}
func (h *TEventStore) OnFloat(value float64) {
	h.add(BF(value))
	h.receiver.OnFloat(value)
}
func (h *TEventStore) OnBigFloat(value *big.Float) {
	h.add(NewTEvent(TEventBigFloat, value, nil))
	h.receiver.OnBigFloat(value)
}
func (h *TEventStore) OnDecimalFloat(value compact_float.DFloat) {
	h.add(DF(value))
	h.receiver.OnDecimalFloat(value)
}
func (h *TEventStore) OnBigDecimalFloat(value *apd.Decimal) {
	h.add(BDF(value))
	h.receiver.OnBigDecimalFloat(value)
}
func (h *TEventStore) OnUID(value []byte) {
	h.add(UID(CloneBytes(value)))
	h.receiver.OnUID(value)
}
func (h *TEventStore) OnTime(value time.Time) {
	h.add(GT(value))
	h.receiver.OnTime(value)
}
func (h *TEventStore) OnCompactTime(value compact_time.Time) {
	h.add(T(value))
	h.receiver.OnCompactTime(value)
}
func (h *TEventStore) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	switch arrayType {
	case events.ArrayTypeString:
		h.add(S(string(value)))
	case events.ArrayTypeResourceID:
		h.add(RID(string(value)))
	case events.ArrayTypeRemoteRef:
		h.add(RREF(string(value)))
	case events.ArrayTypeCustomBinary:
		h.add(CB(CloneBytes(value)))
	case events.ArrayTypeCustomText:
		h.add(CT(string(value)))
	case events.ArrayTypeBit:
		h.add(AB(elementCount, CloneBytes(value)))
	case events.ArrayTypeInt8:
		h.add(AI8(arrays.BytesToInt8Slice(value)))
	case events.ArrayTypeInt16:
		h.add(AI16(arrays.BytesToInt16Slice(value)))
	case events.ArrayTypeInt32:
		h.add(AI32(arrays.BytesToInt32Slice(value)))
	case events.ArrayTypeInt64:
		h.add(AI64(arrays.BytesToInt64Slice(value)))
	case events.ArrayTypeUint8:
		h.add(AU8(CloneBytes(value)))
	case events.ArrayTypeUint16:
		h.add(AU16(arrays.BytesToUint16Slice(value)))
	case events.ArrayTypeUint32:
		h.add(AU32(arrays.BytesToUint32Slice(value)))
	case events.ArrayTypeUint64:
		h.add(AU64(arrays.BytesToUint64Slice(value)))
	case events.ArrayTypeFloat16:
		h.add(AF16(CloneBytes(value)))
	case events.ArrayTypeFloat32:
		h.add(AF32(arrays.BytesToFloat32Slice(value)))
	case events.ArrayTypeFloat64:
		h.add(AF64(arrays.BytesToFloat64Slice(value)))
	case events.ArrayTypeUID:
		h.add(AU(arrays.BytesToUUIDSlice(CloneBytes(value))))
	default:
		panic(fmt.Errorf("TODO: TEventStore.OnArray: Typed array support for %v", arrayType))
	}
	h.receiver.OnArray(arrayType, elementCount, value)
}
func (h *TEventStore) OnStringlikeArray(arrayType events.ArrayType, value string) {
	switch arrayType {
	case events.ArrayTypeString:
		h.add(S(value))
	case events.ArrayTypeResourceID:
		h.add(RID(value))
	case events.ArrayTypeRemoteRef:
		h.add(RREF(value))
	case events.ArrayTypeCustomText:
		h.add(CT(value))
	default:
		panic(fmt.Errorf("BUG: Array type %v is not stringlike", arrayType))
	}
	h.receiver.OnStringlikeArray(arrayType, value)
}
func (h *TEventStore) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString:
		h.add(SB())
	case events.ArrayTypeResourceID:
		h.add(RB())
	case events.ArrayTypeRemoteRef:
		h.add(RRB())
	case events.ArrayTypeCustomBinary:
		h.add(CBB())
	case events.ArrayTypeCustomText:
		h.add(CTB())
	case events.ArrayTypeBit:
		h.add(ABB())
	case events.ArrayTypeInt8:
		h.add(AI8B())
	case events.ArrayTypeInt16:
		h.add(AI16B())
	case events.ArrayTypeInt32:
		h.add(AI32B())
	case events.ArrayTypeInt64:
		h.add(AI64B())
	case events.ArrayTypeUint8:
		h.add(AU8B())
	case events.ArrayTypeUint16:
		h.add(AU16B())
	case events.ArrayTypeUint32:
		h.add(AU32B())
	case events.ArrayTypeUint64:
		h.add(AU64B())
	case events.ArrayTypeFloat16:
		h.add(AF16B())
	case events.ArrayTypeFloat32:
		h.add(AF32B())
	case events.ArrayTypeFloat64:
		h.add(AF64B())
	case events.ArrayTypeUID:
		h.add(AUB())
	case events.ArrayTypeMedia:
		h.add(MB())
	default:
		panic(fmt.Errorf("TODO: TEventStore.OnArrayBegin: Typed array support for %v", arrayType))
	}
	h.receiver.OnArrayBegin(arrayType)
}
func (h *TEventStore) OnArrayChunk(length uint64, moreChunks bool) {
	h.add(AC(length, moreChunks))
	h.receiver.OnArrayChunk(length, moreChunks)
}
func (h *TEventStore) OnArrayData(data []byte) {
	h.add(AD(CloneBytes(data)))
	h.receiver.OnArrayData(data)
}
func (h *TEventStore) OnList() {
	h.add(L())
	h.receiver.OnList()
}
func (h *TEventStore) OnMap() {
	h.add(M())
	h.receiver.OnMap()
}
func (h *TEventStore) OnMarkup(id []byte) {
	h.add(MU(string(id)))
	h.receiver.OnMarkup(id)
}
func (h *TEventStore) OnEnd() {
	h.add(E())
	h.receiver.OnEnd()
}
func (h *TEventStore) OnNode() {
	h.add(NODE())
	h.receiver.OnNode()
}
func (h *TEventStore) OnEdge() {
	h.add(EDGE())
	h.receiver.OnEdge()
}
func (h *TEventStore) OnMarker(id []byte) {
	h.add(MARK(string(id)))
	h.receiver.OnMarker(id)
}
func (h *TEventStore) OnReference(id []byte) {
	h.add(REF(string(id)))
	h.receiver.OnReference(id)
}
func (h *TEventStore) OnConstant(n []byte) {
	h.add(CONST(string(n)))
	h.receiver.OnConstant(n)
}
func (h *TEventStore) OnBeginDocument() {
	h.Events = h.Events[:0]
	h.add(BD())
	h.receiver.OnBeginDocument()
}
func (h *TEventStore) OnEndDocument() {
	h.add(ED())
	h.receiver.OnEndDocument()
}
func (h *TEventStore) OnNan(signaling bool) {
	if signaling {
		h.add(SNAN())
	} else {
		h.add(NAN())
	}
	h.receiver.OnNan(signaling)
}
