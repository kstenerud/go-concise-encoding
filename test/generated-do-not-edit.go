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

// Generated by github.com/kstenerud/go-concise-encoding/codegen
  // DO NOT EDIT THIS FILE. Contents will be overwritten.

package test

import (
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/events"
)

type EventArrayBit struct{ EventWithValue }

func AB(elements []bool) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []bool
	if v != nil {
		safeArg = v.([]bool)
	}

	return &EventArrayBit{
		EventWithValue: ConstructEventWithValue("ab", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeBit, uint64(len(safeArg)), arrayBitsToBytes(safeArg))
		}),
	}
}

type EventArrayDataBit struct{ EventWithValue }

func ADB(elements []bool) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []bool
	if v != nil {
		safeArg = v.([]bool)
	}

	return &EventArrayDataBit{
		EventWithValue: ConstructEventWithValue("adb", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayBitsToBytes(safeArg))
		}),
	}
}

type EventVersion struct{ EventWithValue }

func V(version uint64) Event {
	v := version
	safeArg := v

	return &EventVersion{
		EventWithValue: ConstructEventWithValue("v", v, func(receiver events.DataEventReceiver) {
			receiver.OnVersion(safeArg)
		}),
	}
}

type EventBoolean struct{ EventWithValue }

func B(value bool) Event {
	v := value
	safeArg := v

	return &EventBoolean{
		EventWithValue: ConstructEventWithValue("b", v, func(receiver events.DataEventReceiver) {
			receiver.OnBoolean(safeArg)
		}),
	}
}

type EventTime struct{ EventWithValue }

func T(value compact_time.Time) Event {
	v := value
	safeArg := v

	return &EventTime{
		EventWithValue: ConstructEventWithValue("t", v, func(receiver events.DataEventReceiver) {
			receiver.OnTime(safeArg)
		}),
	}
}

type EventUID struct{ EventWithValue }

func UID(value []byte) Event {
	v := copyOf(value)
	if len(value) == 0 {
		v = NoValue
	}
	var safeArg []byte
	if v != nil {
		safeArg = v.([]byte)
	}

	return &EventUID{
		EventWithValue: ConstructEventWithValue("uid", v, func(receiver events.DataEventReceiver) {
			receiver.OnUID(safeArg)
		}),
	}
}

type EventArrayChunkMore struct{ EventWithValue }

func ACM(length uint64) Event {
	v := length
	safeArg := v

	return &EventArrayChunkMore{
		EventWithValue: ConstructEventWithValue("acm", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayChunk(safeArg, true)
		}),
	}
}

type EventArrayChunkLast struct{ EventWithValue }

func ACL(length uint64) Event {
	v := length
	safeArg := v

	return &EventArrayChunkLast{
		EventWithValue: ConstructEventWithValue("acl", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayChunk(safeArg, false)
		}),
	}
}

type EventCommentMultiline struct{ EventWithValue }

func CM(comment string) Event {
	v := comment
	safeArg := v

	return &EventCommentMultiline{
		EventWithValue: ConstructEventWithValue("cm", v, func(receiver events.DataEventReceiver) {
			receiver.OnComment(true, []byte(safeArg))
		}),
	}
}

type EventCommentSingleLine struct{ EventWithValue }

func CS(comment string) Event {
	v := comment
	safeArg := v

	return &EventCommentSingleLine{
		EventWithValue: ConstructEventWithValue("cs", v, func(receiver events.DataEventReceiver) {
			receiver.OnComment(false, []byte(safeArg))
		}),
	}
}

type EventCustomBinary struct{ EventWithValue }

func CB(elements []byte) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []byte
	if v != nil {
		safeArg = v.([]byte)
	}

	return &EventCustomBinary{
		EventWithValue: ConstructEventWithValue("cb", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeCustomBinary, uint64(len(safeArg)), safeArg)
		}),
	}
}

type EventMarker struct{ EventWithValue }

func MARK(id string) Event {
	v := id
	safeArg := v

	return &EventMarker{
		EventWithValue: ConstructEventWithValue("mark", v, func(receiver events.DataEventReceiver) {
			receiver.OnMarker([]byte(safeArg))
		}),
	}
}

type EventReferenceLocal struct{ EventWithValue }

func REFL(id string) Event {
	v := id
	safeArg := v

	return &EventReferenceLocal{
		EventWithValue: ConstructEventWithValue("refl", v, func(receiver events.DataEventReceiver) {
			receiver.OnReferenceLocal([]byte(safeArg))
		}),
	}
}

type EventStructInstance struct{ EventWithValue }

func SI(id string) Event {
	v := id
	safeArg := v

	return &EventStructInstance{
		EventWithValue: ConstructEventWithValue("si", v, func(receiver events.DataEventReceiver) {
			receiver.OnStructInstance([]byte(safeArg))
		}),
	}
}

type EventStructTemplate struct{ EventWithValue }

func ST(id string) Event {
	v := id
	safeArg := v

	return &EventStructTemplate{
		EventWithValue: ConstructEventWithValue("st", v, func(receiver events.DataEventReceiver) {
			receiver.OnStructTemplate([]byte(safeArg))
		}),
	}
}

type EventString struct{ EventWithValue }

func S(str string) Event {
	v := str
	safeArg := v

	return &EventString{
		EventWithValue: ConstructEventWithValue("s", v, func(receiver events.DataEventReceiver) {
			receiver.OnStringlikeArray(events.ArrayTypeString, safeArg)
		}),
	}
}

type EventCustomText struct{ EventWithValue }

func CT(str string) Event {
	v := str
	safeArg := v

	return &EventCustomText{
		EventWithValue: ConstructEventWithValue("ct", v, func(receiver events.DataEventReceiver) {
			receiver.OnStringlikeArray(events.ArrayTypeCustomText, safeArg)
		}),
	}
}

type EventReferenceRemote struct{ EventWithValue }

func REFR(str string) Event {
	v := str
	safeArg := v

	return &EventReferenceRemote{
		EventWithValue: ConstructEventWithValue("refr", v, func(receiver events.DataEventReceiver) {
			receiver.OnStringlikeArray(events.ArrayTypeReferenceRemote, safeArg)
		}),
	}
}

type EventResourceID struct{ EventWithValue }

func RID(str string) Event {
	v := str
	safeArg := v

	return &EventResourceID{
		EventWithValue: ConstructEventWithValue("rid", v, func(receiver events.DataEventReceiver) {
			receiver.OnStringlikeArray(events.ArrayTypeResourceID, safeArg)
		}),
	}
}

type EventArrayDataInt8 struct{ EventWithValue }

func ADI8(elements []int8) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int8
	if v != nil {
		safeArg = v.([]int8)
	}

	return &EventArrayDataInt8{
		EventWithValue: ConstructEventWithValue("adi8", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayInt8ToBytes(safeArg))
		}),
	}
}

type EventArrayDataInt16 struct{ EventWithValue }

func ADI16(elements []int16) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int16
	if v != nil {
		safeArg = v.([]int16)
	}

	return &EventArrayDataInt16{
		EventWithValue: ConstructEventWithValue("adi16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayInt16ToBytes(safeArg))
		}),
	}
}

type EventArrayDataInt32 struct{ EventWithValue }

func ADI32(elements []int32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int32
	if v != nil {
		safeArg = v.([]int32)
	}

	return &EventArrayDataInt32{
		EventWithValue: ConstructEventWithValue("adi32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayInt32ToBytes(safeArg))
		}),
	}
}

type EventArrayDataInt64 struct{ EventWithValue }

func ADI64(elements []int64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int64
	if v != nil {
		safeArg = v.([]int64)
	}

	return &EventArrayDataInt64{
		EventWithValue: ConstructEventWithValue("adi64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayInt64ToBytes(safeArg))
		}),
	}
}

type EventArrayDataFloat16 struct{ EventWithValue }

func ADF16(elements []float32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float32
	if v != nil {
		safeArg = v.([]float32)
	}

	return &EventArrayDataFloat16{
		EventWithValue: ConstructEventWithValue("adf16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayFloat16ToBytes(safeArg))
		}),
	}
}

type EventArrayDataFloat32 struct{ EventWithValue }

func ADF32(elements []float32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float32
	if v != nil {
		safeArg = v.([]float32)
	}

	return &EventArrayDataFloat32{
		EventWithValue: ConstructEventWithValue("adf32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayFloat32ToBytes(safeArg))
		}),
	}
}

type EventArrayDataFloat64 struct{ EventWithValue }

func ADF64(elements []float64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float64
	if v != nil {
		safeArg = v.([]float64)
	}

	return &EventArrayDataFloat64{
		EventWithValue: ConstructEventWithValue("adf64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayFloat64ToBytes(safeArg))
		}),
	}
}

type EventArrayDataUint8 struct{ EventWithValue }

func ADU8(elements []uint8) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint8
	if v != nil {
		safeArg = v.([]uint8)
	}

	return &EventArrayDataUint8{
		EventWithValue: ConstructEventWithValue("adu8", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayUint8ToBytes(safeArg))
		}),
	}
}

type EventArrayDataUint16 struct{ EventWithValue }

func ADU16(elements []uint16) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint16
	if v != nil {
		safeArg = v.([]uint16)
	}

	return &EventArrayDataUint16{
		EventWithValue: ConstructEventWithValue("adu16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayUint16ToBytes(safeArg))
		}),
	}
}

type EventArrayDataUint32 struct{ EventWithValue }

func ADU32(elements []uint32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint32
	if v != nil {
		safeArg = v.([]uint32)
	}

	return &EventArrayDataUint32{
		EventWithValue: ConstructEventWithValue("adu32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayUint32ToBytes(safeArg))
		}),
	}
}

type EventArrayDataUint64 struct{ EventWithValue }

func ADU64(elements []uint64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint64
	if v != nil {
		safeArg = v.([]uint64)
	}

	return &EventArrayDataUint64{
		EventWithValue: ConstructEventWithValue("adu64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayUint64ToBytes(safeArg))
		}),
	}
}

type EventArrayDataUID struct{ EventWithValue }

func ADU(elements [][]byte) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg [][]byte
	if v != nil {
		safeArg = v.([][]byte)
	}

	return &EventArrayDataUID{
		EventWithValue: ConstructEventWithValue("adu", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayUIDToBytes(safeArg))
		}),
	}
}

type EventArrayDataText struct{ EventWithValue }

func ADT(elements string) Event {
	v := elements
	safeArg := v

	return &EventArrayDataText{
		EventWithValue: ConstructEventWithValue("adt", v, func(receiver events.DataEventReceiver) {
			receiver.OnArrayData(arrayTextToBytes(safeArg))
		}),
	}
}

type EventArrayInt8 struct{ EventWithValue }

func AI8(elements []int8) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int8
	if v != nil {
		safeArg = v.([]int8)
	}

	return &EventArrayInt8{
		EventWithValue: ConstructEventWithValue("ai8", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeInt8, uint64(len(safeArg)), arrayInt8ToBytes(safeArg))
		}),
	}
}

type EventArrayInt16 struct{ EventWithValue }

func AI16(elements []int16) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int16
	if v != nil {
		safeArg = v.([]int16)
	}

	return &EventArrayInt16{
		EventWithValue: ConstructEventWithValue("ai16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeInt16, uint64(len(safeArg)), arrayInt16ToBytes(safeArg))
		}),
	}
}

type EventArrayInt32 struct{ EventWithValue }

func AI32(elements []int32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int32
	if v != nil {
		safeArg = v.([]int32)
	}

	return &EventArrayInt32{
		EventWithValue: ConstructEventWithValue("ai32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeInt32, uint64(len(safeArg)), arrayInt32ToBytes(safeArg))
		}),
	}
}

type EventArrayInt64 struct{ EventWithValue }

func AI64(elements []int64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []int64
	if v != nil {
		safeArg = v.([]int64)
	}

	return &EventArrayInt64{
		EventWithValue: ConstructEventWithValue("ai64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeInt64, uint64(len(safeArg)), arrayInt64ToBytes(safeArg))
		}),
	}
}

type EventArrayFloat16 struct{ EventWithValue }

func AF16(elements []float32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float32
	if v != nil {
		safeArg = v.([]float32)
	}

	return &EventArrayFloat16{
		EventWithValue: ConstructEventWithValue("af16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeFloat16, uint64(len(safeArg)), arrayFloat16ToBytes(safeArg))
		}),
	}
}

type EventArrayFloat32 struct{ EventWithValue }

func AF32(elements []float32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float32
	if v != nil {
		safeArg = v.([]float32)
	}

	return &EventArrayFloat32{
		EventWithValue: ConstructEventWithValue("af32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeFloat32, uint64(len(safeArg)), arrayFloat32ToBytes(safeArg))
		}),
	}
}

type EventArrayFloat64 struct{ EventWithValue }

func AF64(elements []float64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []float64
	if v != nil {
		safeArg = v.([]float64)
	}

	return &EventArrayFloat64{
		EventWithValue: ConstructEventWithValue("af64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeFloat64, uint64(len(safeArg)), arrayFloat64ToBytes(safeArg))
		}),
	}
}

type EventArrayUint8 struct{ EventWithValue }

func AU8(elements []uint8) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint8
	if v != nil {
		safeArg = v.([]uint8)
	}

	return &EventArrayUint8{
		EventWithValue: ConstructEventWithValue("au8", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeUint8, uint64(len(safeArg)), arrayUint8ToBytes(safeArg))
		}),
	}
}

type EventArrayUint16 struct{ EventWithValue }

func AU16(elements []uint16) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint16
	if v != nil {
		safeArg = v.([]uint16)
	}

	return &EventArrayUint16{
		EventWithValue: ConstructEventWithValue("au16", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeUint16, uint64(len(safeArg)), arrayUint16ToBytes(safeArg))
		}),
	}
}

type EventArrayUint32 struct{ EventWithValue }

func AU32(elements []uint32) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint32
	if v != nil {
		safeArg = v.([]uint32)
	}

	return &EventArrayUint32{
		EventWithValue: ConstructEventWithValue("au32", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeUint32, uint64(len(safeArg)), arrayUint32ToBytes(safeArg))
		}),
	}
}

type EventArrayUint64 struct{ EventWithValue }

func AU64(elements []uint64) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg []uint64
	if v != nil {
		safeArg = v.([]uint64)
	}

	return &EventArrayUint64{
		EventWithValue: ConstructEventWithValue("au64", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeUint64, uint64(len(safeArg)), arrayUint64ToBytes(safeArg))
		}),
	}
}

type EventArrayUID struct{ EventWithValue }

func AU(elements [][]byte) Event {
	v := copyOf(elements)
	if len(elements) == 0 {
		v = NoValue
	}
	var safeArg [][]byte
	if v != nil {
		safeArg = v.([][]byte)
	}

	return &EventArrayUID{
		EventWithValue: ConstructEventWithValue("au", v, func(receiver events.DataEventReceiver) {
			receiver.OnArray(events.ArrayTypeUID, uint64(len(safeArg)), arrayUIDToBytes(safeArg))
		}),
	}
}

type EventBeginArrayBit struct{ EventWithValue }

func BAB() Event {
	return &EventBeginArrayBit{
		EventWithValue: ConstructEventWithValue("bab", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeBit)
		}),
	}
}

type EventBeginArrayFloat16 struct{ EventWithValue }

func BAF16() Event {
	return &EventBeginArrayFloat16{
		EventWithValue: ConstructEventWithValue("baf16", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeFloat16)
		}),
	}
}

type EventBeginArrayFloat32 struct{ EventWithValue }

func BAF32() Event {
	return &EventBeginArrayFloat32{
		EventWithValue: ConstructEventWithValue("baf32", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeFloat32)
		}),
	}
}

type EventBeginArrayFloat64 struct{ EventWithValue }

func BAF64() Event {
	return &EventBeginArrayFloat64{
		EventWithValue: ConstructEventWithValue("baf64", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeFloat64)
		}),
	}
}

type EventBeginArrayInt8 struct{ EventWithValue }

func BAI8() Event {
	return &EventBeginArrayInt8{
		EventWithValue: ConstructEventWithValue("bai8", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeInt8)
		}),
	}
}

type EventBeginArrayInt16 struct{ EventWithValue }

func BAI16() Event {
	return &EventBeginArrayInt16{
		EventWithValue: ConstructEventWithValue("bai16", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeInt16)
		}),
	}
}

type EventBeginArrayInt32 struct{ EventWithValue }

func BAI32() Event {
	return &EventBeginArrayInt32{
		EventWithValue: ConstructEventWithValue("bai32", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeInt32)
		}),
	}
}

type EventBeginArrayInt64 struct{ EventWithValue }

func BAI64() Event {
	return &EventBeginArrayInt64{
		EventWithValue: ConstructEventWithValue("bai64", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeInt64)
		}),
	}
}

type EventBeginArrayUID struct{ EventWithValue }

func BAU() Event {
	return &EventBeginArrayUID{
		EventWithValue: ConstructEventWithValue("bau", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeUID)
		}),
	}
}

type EventBeginArrayUint8 struct{ EventWithValue }

func BAU8() Event {
	return &EventBeginArrayUint8{
		EventWithValue: ConstructEventWithValue("bau8", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeUint8)
		}),
	}
}

type EventBeginArrayUint16 struct{ EventWithValue }

func BAU16() Event {
	return &EventBeginArrayUint16{
		EventWithValue: ConstructEventWithValue("bau16", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeUint16)
		}),
	}
}

type EventBeginArrayUint32 struct{ EventWithValue }

func BAU32() Event {
	return &EventBeginArrayUint32{
		EventWithValue: ConstructEventWithValue("bau32", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeUint32)
		}),
	}
}

type EventBeginArrayUint64 struct{ EventWithValue }

func BAU64() Event {
	return &EventBeginArrayUint64{
		EventWithValue: ConstructEventWithValue("bau64", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeUint64)
		}),
	}
}

type EventBeginCustomBinary struct{ EventWithValue }

func BCB() Event {
	return &EventBeginCustomBinary{
		EventWithValue: ConstructEventWithValue("bcb", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeCustomBinary)
		}),
	}
}

type EventBeginCustomText struct{ EventWithValue }

func BCT() Event {
	return &EventBeginCustomText{
		EventWithValue: ConstructEventWithValue("bct", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeCustomText)
		}),
	}
}

type EventBeginMedia struct{ EventWithValue }

func BMEDIA() Event {
	return &EventBeginMedia{
		EventWithValue: ConstructEventWithValue("bmedia", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeMedia)
		}),
	}
}

type EventBeginReferenceRemote struct{ EventWithValue }

func BREFR() Event {
	return &EventBeginReferenceRemote{
		EventWithValue: ConstructEventWithValue("brefr", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeReferenceRemote)
		}),
	}
}

type EventBeginResourceID struct{ EventWithValue }

func BRID() Event {
	return &EventBeginResourceID{
		EventWithValue: ConstructEventWithValue("brid", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeResourceID)
		}),
	}
}

type EventBeginString struct{ EventWithValue }

func BS() Event {
	return &EventBeginString{
		EventWithValue: ConstructEventWithValue("bs", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnArrayBegin(events.ArrayTypeString)
		}),
	}
}

type EventEdge struct{ EventWithValue }

func EDGE() Event {
	return &EventEdge{
		EventWithValue: ConstructEventWithValue("edge", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnEdge()
		}),
	}
}

type EventEnd struct{ EventWithValue }

func E() Event {
	return &EventEnd{
		EventWithValue: ConstructEventWithValue("e", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnEnd()
		}),
	}
}

type EventList struct{ EventWithValue }

func L() Event {
	return &EventList{
		EventWithValue: ConstructEventWithValue("l", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnList()
		}),
	}
}

type EventMap struct{ EventWithValue }

func M() Event {
	return &EventMap{
		EventWithValue: ConstructEventWithValue("m", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnMap()
		}),
	}
}

type EventNode struct{ EventWithValue }

func NODE() Event {
	return &EventNode{
		EventWithValue: ConstructEventWithValue("node", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnNode()
		}),
	}
}

type EventNull struct{ EventWithValue }

func NULL() Event {
	return &EventNull{
		EventWithValue: ConstructEventWithValue("null", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnNull()
		}),
	}
}

type EventPadding struct{ EventWithValue }

func PAD() Event {
	return &EventPadding{
		EventWithValue: ConstructEventWithValue("pad", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnPadding()
		}),
	}
}

type EventBeginDocument struct{ EventWithValue }

func BD() Event {
	return &EventBeginDocument{
		EventWithValue: ConstructEventWithValue("bd", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnBeginDocument()
		}),
	}
}

type EventEndDocument struct{ EventWithValue }

func ED() Event {
	return &EventEndDocument{
		EventWithValue: ConstructEventWithValue("ed", NoValue, func(receiver events.DataEventReceiver) {
			receiver.OnEndDocument()
		}),
	}
}

