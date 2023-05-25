// Copyright 2022 Karl Stenerud
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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/nullevent"
)

type EventCollection struct {
	Events Events
}

func (_this *EventCollection) Clear() {
	_this.Events = _this.Events[:0]
}

func (_this *EventCollection) IsEquivalentTo(events Events) bool {
	return AreEventsEquivalent(_this.Events, events)
}

func AreEventsEquivalent(a Events, b Events) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].IsEquivalentTo(b[i]) {
			return false
		}
	}
	return true
}

func NewEventCollector(next events.DataEventReceiver) (receiver events.DataEventReceiver, collection *EventCollection) {
	collection = &EventCollection{}
	receiver = NewEventReceiver(next, func(event Event) {
		collection.Events = append(collection.Events, event)
	})
	return
}

func NewEventPrinter(next events.DataEventReceiver) events.DataEventReceiver {
	return NewEventReceiver(next, func(event Event) {
		fmt.Println(event)
	})
}

func NewEventReceiver(next events.DataEventReceiver, iterate func(Event)) events.DataEventReceiver {
	if iterate == nil {
		iterate = func(Event) {}
	}
	if next == nil {
		next = nullevent.NewNullEventReceiver()
	}
	return &EventReceiver{
		next:    next,
		iterate: iterate,
	}
}

type EventReceiver struct {
	next                   events.DataEventReceiver
	iterate                func(Event)
	arrayType              events.ArrayType
	arrayElementsRemaining uint64
}

func (_this *EventReceiver) OnVersion(version uint64) {
	_this.iterate(V(version))
	_this.next.OnVersion(version)
}
func (_this *EventReceiver) OnPadding() {
	_this.iterate(PAD())
	_this.next.OnPadding()
}
func (_this *EventReceiver) OnComment(isMultiline bool, contents []byte) {
	if isMultiline {
		_this.iterate(CM(string(contents)))
	} else {
		_this.iterate(CS(string(contents)))
	}
	_this.next.OnComment(isMultiline, contents)
}
func (_this *EventReceiver) OnNull() {
	_this.iterate(NULL())
	_this.next.OnNull()
}
func (_this *EventReceiver) OnBoolean(value bool) {
	_this.iterate(B(value))
	_this.next.OnBoolean(value)
}
func (_this *EventReceiver) OnTrue() {
	_this.iterate(B(true))
	_this.next.OnTrue()
}
func (_this *EventReceiver) OnFalse() {
	_this.iterate(B(false))
	_this.next.OnFalse()
}
func (_this *EventReceiver) OnPositiveInt(value uint64) {
	_this.iterate(N(value))
	_this.next.OnPositiveInt(value)
}
func (_this *EventReceiver) OnNegativeInt(value uint64) {
	if value < 0x8000000000000000 {
		if value == 0 {
			_this.iterate(N(compact_float.NegativeZero()))
		} else {
			_this.iterate(N(-int64(value)))
		}
	} else {
		bi := big.NewInt(0)
		bi.SetUint64(value)
		bi = bi.Neg(bi)
		_this.iterate(N(bi))
	}
	_this.next.OnNegativeInt(value)
}
func (_this *EventReceiver) OnInt(value int64) {
	_this.iterate(N(value))
	_this.next.OnInt(value)
}
func (_this *EventReceiver) OnBigInt(value *big.Int) {
	_this.iterate(N(value))
	_this.next.OnBigInt(value)
}
func (_this *EventReceiver) OnFloat(value float64) {
	_this.iterate(N(value))
	_this.next.OnFloat(value)
}
func (_this *EventReceiver) OnBigFloat(value *big.Float) {
	_this.iterate(N(value))
	_this.next.OnBigFloat(value)
}
func (_this *EventReceiver) OnDecimalFloat(value compact_float.DFloat) {
	_this.iterate(N(value))
	_this.next.OnDecimalFloat(value)
}
func (_this *EventReceiver) OnBigDecimalFloat(value *apd.Decimal) {
	_this.iterate(N(value))
	_this.next.OnBigDecimalFloat(value)
}
func (_this *EventReceiver) OnUID(value []byte) {
	_this.iterate(UID(value))
	_this.next.OnUID(value)
}
func (_this *EventReceiver) OnTime(value compact_time.Time) {
	_this.iterate(T(value))
	_this.next.OnTime(value)
}
func (_this *EventReceiver) OnArray(arrayType events.ArrayType, elementCount uint64, value []byte) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.iterate(S(string(value)))
	case events.ArrayTypeResourceID:
		_this.iterate(RID(string(value)))
	case events.ArrayTypeReferenceRemote:
		_this.iterate(REFR(string(value)))
	case events.ArrayTypeBit:
		_this.iterate(AB(bytesToArrayBits(elementCount, value)))
	case events.ArrayTypeUint8:
		_this.iterate(AU8(value))
	case events.ArrayTypeUint16:
		_this.iterate(AU16(bytesToArrayUint16(value)))
	case events.ArrayTypeUint32:
		_this.iterate(AU32(bytesToArrayUint32(value)))
	case events.ArrayTypeUint64:
		_this.iterate(AU64(bytesToArrayUint64(value)))
	case events.ArrayTypeInt8:
		_this.iterate(AI8(bytesToArrayInt8(value)))
	case events.ArrayTypeInt16:
		_this.iterate(AI16(bytesToArrayInt16(value)))
	case events.ArrayTypeInt32:
		_this.iterate(AI32(bytesToArrayInt32(value)))
	case events.ArrayTypeInt64:
		_this.iterate(AI64(bytesToArrayInt64(value)))
	case events.ArrayTypeFloat16:
		_this.iterate(AF16(bytesToArrayFloat16(value)))
	case events.ArrayTypeFloat32:
		_this.iterate(AF32(bytesToArrayFloat32(value)))
	case events.ArrayTypeFloat64:
		_this.iterate(AF64(bytesToArrayFloat64(value)))
	case events.ArrayTypeUID:
		_this.iterate(AU(bytesToArrayUID(value)))
	default:
		panic(fmt.Errorf("unknown array type %v", arrayType))
	}
	_this.next.OnArray(arrayType, elementCount, value)
}
func (_this *EventReceiver) OnStringlikeArray(arrayType events.ArrayType, value string) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.iterate(S(value))
	case events.ArrayTypeResourceID:
		_this.iterate(RID(value))
	case events.ArrayTypeReferenceRemote:
		_this.iterate(REFR(value))
	default:
		panic(fmt.Errorf("unknown array type %v", arrayType))
	}
	_this.next.OnStringlikeArray(arrayType, value)
}
func (_this *EventReceiver) OnMedia(mediaType string, value []byte) {
	_this.iterate(MEDIA(mediaType, value))
	_this.next.OnMedia(mediaType, value)
}
func (_this *EventReceiver) OnCustomText(customType uint64, value string) {
	_this.iterate(CT(customType, value))
	_this.next.OnCustomText(customType, value)
}
func (_this *EventReceiver) OnCustomBinary(customType uint64, value []byte) {
	_this.iterate(CB(customType, value))
	_this.next.OnCustomBinary(customType, value)
}
func (_this *EventReceiver) OnArrayBegin(arrayType events.ArrayType) {
	switch arrayType {
	case events.ArrayTypeString:
		_this.iterate(BS())
	case events.ArrayTypeResourceID:
		_this.iterate(BRID())
	case events.ArrayTypeReferenceRemote:
		_this.iterate(BREFR())
	case events.ArrayTypeBit:
		_this.iterate(BAB())
	case events.ArrayTypeUint8:
		_this.iterate(BAU8())
	case events.ArrayTypeUint16:
		_this.iterate(BAU16())
	case events.ArrayTypeUint32:
		_this.iterate(BAU32())
	case events.ArrayTypeUint64:
		_this.iterate(BAU64())
	case events.ArrayTypeInt8:
		_this.iterate(BAI8())
	case events.ArrayTypeInt16:
		_this.iterate(BAI16())
	case events.ArrayTypeInt32:
		_this.iterate(BAI32())
	case events.ArrayTypeInt64:
		_this.iterate(BAI64())
	case events.ArrayTypeFloat16:
		_this.iterate(BAF16())
	case events.ArrayTypeFloat32:
		_this.iterate(BAF32())
	case events.ArrayTypeFloat64:
		_this.iterate(BAF64())
	case events.ArrayTypeUID:
		_this.iterate(BAU())
	default:
		panic(fmt.Errorf("unknown array type %v", arrayType))
	}
	_this.arrayType = arrayType
	_this.next.OnArrayBegin(arrayType)
}
func (_this *EventReceiver) OnMediaBegin(mediaType string) {
	_this.iterate(BMEDIA(mediaType))
	_this.arrayType = events.ArrayTypeMedia
	_this.next.OnMediaBegin(mediaType)
}
func (_this *EventReceiver) OnCustomBegin(arrayType events.ArrayType, customType uint64) {
	switch arrayType {
	case events.ArrayTypeCustomText:
		_this.iterate(BCT(customType))
	case events.ArrayTypeCustomBinary:
		_this.iterate(BCB(customType))
	default:
		panic(fmt.Errorf("unknown custom type %v", arrayType))
	}
	_this.arrayType = arrayType
	_this.next.OnCustomBegin(arrayType, customType)
}
func (_this *EventReceiver) OnArrayChunk(length uint64, moreChunks bool) {
	if moreChunks {
		_this.iterate(ACM(length))
	} else {
		_this.iterate(ACL(length))
	}
	_this.arrayElementsRemaining = length
	_this.next.OnArrayChunk(length, moreChunks)
}
func (_this *EventReceiver) OnArrayData(data []byte) {
	switch _this.arrayType {
	case events.ArrayTypeString,
		events.ArrayTypeResourceID,
		events.ArrayTypeReferenceRemote,
		events.ArrayTypeCustomText:
		_this.iterate(ADT(string(data)))
		_this.arrayElementsRemaining -= uint64(len(data))
	case events.ArrayTypeCustomBinary:
		_this.iterate(ADU8(data))
		_this.arrayElementsRemaining -= uint64(len(data))
	case events.ArrayTypeBit:
		elementCount := uint64(len(data)) * 8
		if elementCount > _this.arrayElementsRemaining {
			elementCount = _this.arrayElementsRemaining
		}
		_this.iterate(ADB(bytesToArrayBits(elementCount, data)))
		_this.arrayElementsRemaining -= elementCount
	case events.ArrayTypeUint8:
		_this.iterate(ADU8(data))
		_this.arrayElementsRemaining -= uint64(len(data))
	case events.ArrayTypeUint16:
		elements := bytesToArrayUint16(data)
		_this.iterate(ADU16(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeUint32:
		elements := bytesToArrayUint32(data)
		_this.iterate(ADU32(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeUint64:
		elements := bytesToArrayUint64(data)
		_this.iterate(ADU64(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeInt8:
		_this.iterate(ADI8(bytesToArrayInt8(data)))
		_this.arrayElementsRemaining -= uint64(len(data))
	case events.ArrayTypeInt16:
		elements := bytesToArrayInt16(data)
		_this.iterate(ADI16(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeInt32:
		elements := bytesToArrayInt32(data)
		_this.iterate(ADI32(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeInt64:
		elements := bytesToArrayInt64(data)
		_this.iterate(ADI64(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeFloat16:
		elements := bytesToArrayFloat16(data)
		_this.iterate(ADF16(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeFloat32:
		elements := bytesToArrayFloat32(data)
		_this.iterate(ADF32(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeFloat64:
		elements := bytesToArrayFloat64(data)
		_this.iterate(ADF64(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeUID:
		elements := bytesToArrayUID(data)
		_this.iterate(ADU(elements))
		_this.arrayElementsRemaining -= uint64(len(elements))
	case events.ArrayTypeMedia:
		_this.iterate(ADT(string(data)))
		_this.arrayElementsRemaining -= uint64(len(data))
		if _this.arrayElementsRemaining == 0 {
			_this.arrayType = events.ArrayTypeMediaData
		}
	case events.ArrayTypeMediaData:
		_this.iterate(ADU8(data))
		_this.arrayElementsRemaining -= uint64(len(data))
	default:
		panic(fmt.Errorf("unknown array type %v", _this.arrayType))
	}
	_this.next.OnArrayData(data)
}
func (_this *EventReceiver) OnList() {
	_this.iterate(L())
	_this.next.OnList()
}
func (_this *EventReceiver) OnMap() {
	_this.iterate(M())
	_this.next.OnMap()
}
func (_this *EventReceiver) OnRecordType(id []byte) {
	_this.iterate(RT(string(id)))
	_this.next.OnRecordType(id)
}
func (_this *EventReceiver) OnRecord(id []byte) {
	_this.iterate(REC(string(id)))
	_this.next.OnRecord(id)
}
func (_this *EventReceiver) OnEndContainer() {
	_this.iterate(E())
	_this.next.OnEndContainer()
}
func (_this *EventReceiver) OnNode() {
	_this.iterate(NODE())
	_this.next.OnNode()
}
func (_this *EventReceiver) OnEdge() {
	_this.iterate(EDGE())
	_this.next.OnEdge()
}
func (_this *EventReceiver) OnMarker(id []byte) {
	_this.iterate(MARK(string(id)))
	_this.next.OnMarker(id)
}
func (_this *EventReceiver) OnReferenceLocal(id []byte) {
	_this.iterate(REFL(string(id)))
	_this.next.OnReferenceLocal(id)
}
func (_this *EventReceiver) OnBeginDocument() {
	// Nothing to do
	_this.next.OnBeginDocument()
}
func (_this *EventReceiver) OnEndDocument() {
	// Nothing to do
	_this.next.OnEndDocument()
}
func (_this *EventReceiver) OnNan(signaling bool) {
	if signaling {
		_this.iterate(N(compact_float.SignalingNaN()))
	} else {
		_this.iterate(N(compact_float.QuietNaN()))
	}
	_this.next.OnNan(signaling)
}
