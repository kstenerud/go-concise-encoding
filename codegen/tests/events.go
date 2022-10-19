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

package tests

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/codegen/standard"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
)

const codePath = "test"

var imports = []*standard.Import{
	{LocalName: "compact_time", Import: "github.com/kstenerud/go-compact-time"},
	{LocalName: "", Import: "github.com/kstenerud/go-concise-encoding/ce/events"},
}

func generateEvents(projectDir string) {
	generatedFilePath := standard.GetGeneratedCodePath(projectDir, codePath)
	writer, err := os.Create(generatedFilePath)
	standard.PanicIfError(err, "could not open %s", generatedFilePath)
	defer writer.Close()
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Errorf("error while generating %v: %v", generatedFilePath, e))
		}
	}()

	standard.WriteHeader(writer, codePath, imports)

	generateOneArgEvent(writer, "ArrayBit", "ab", "elements", "[]bool", "receiver.OnArray(events.ArrayTypeBit, uint64(len(safeArg)), arrayBitsToBytes(safeArg))")
	generateOneArgEvent(writer, "ArrayDataBit", "adb", "elements", "[]bool", "receiver.OnArrayData(arrayBitsToBytes(safeArg))")

	generateOneArgEvent(writer, "Version", "v", "version", "uint64", "receiver.OnVersion(safeArg)")
	generateOneArgEvent(writer, "Boolean", "b", "value", "bool", "receiver.OnBoolean(safeArg)")
	generateOneArgEvent(writer, "Time", "t", "value", "compact_time.Time", "receiver.OnTime(safeArg)")
	generateOneArgEvent(writer, "UID", "uid", "value", "[]byte", "receiver.OnUID(safeArg)")

	generateOneArgEvent(writer, "ArrayChunkMore", "acm", "length", "uint64", "receiver.OnArrayChunk(safeArg, true)")
	generateOneArgEvent(writer, "ArrayChunkLast", "acl", "length", "uint64", "receiver.OnArrayChunk(safeArg, false)")
	generateOneArgEvent(writer, "CommentMultiline", "cm", "comment", "string", "receiver.OnComment(true, []byte(safeArg))")
	generateOneArgEvent(writer, "CommentSingleLine", "cs", "comment", "string", "receiver.OnComment(false, []byte(safeArg))")

	generateTwoArgEvent(writer, "CustomBinary", "cb", "customType", "uint64", "data", "[]byte", "receiver.OnCustomBinary(customType, safeArg)")
	generateTwoArgEvent(writer, "CustomText", "ct", "customType", "uint64", "data", "string", "receiver.OnCustomText(customType, safeArg)")
	generateTwoArgEvent(writer, "Media", "media", "mediaType", "string", "data", "[]byte", "receiver.OnMedia(mediaType, safeArg)")

	generateIDEvent(writer, "Marker", "mark")
	generateIDEvent(writer, "ReferenceLocal", "refl")
	generateIDEvent(writer, "StructInstance", "si")
	generateIDEvent(writer, "StructTemplate", "st")

	generateStringArrayEvent(writer, "String", "s")
	generateStringArrayEvent(writer, "ReferenceRemote", "refr")
	generateStringArrayEvent(writer, "ResourceID", "rid")

	generateArrayDataEvent(writer, "Int8", "[]int8", "adi8")
	generateArrayDataEvent(writer, "Int16", "[]int16", "adi16")
	generateArrayDataEvent(writer, "Int32", "[]int32", "adi32")
	generateArrayDataEvent(writer, "Int64", "[]int64", "adi64")
	generateArrayDataEvent(writer, "Float16", "[]float32", "adf16")
	generateArrayDataEvent(writer, "Float32", "[]float32", "adf32")
	generateArrayDataEvent(writer, "Float64", "[]float64", "adf64")
	generateArrayDataEvent(writer, "Uint8", "[]uint8", "adu8")
	generateArrayDataEvent(writer, "Uint16", "[]uint16", "adu16")
	generateArrayDataEvent(writer, "Uint32", "[]uint32", "adu32")
	generateArrayDataEvent(writer, "Uint64", "[]uint64", "adu64")
	generateArrayDataEvent(writer, "UID", "[][]byte", "adu")
	generateArrayDataEvent(writer, "Text", "string", "adt")
	generateArrayEvent(writer, "Int8", "int8", "ai8")
	generateArrayEvent(writer, "Int16", "int16", "ai16")
	generateArrayEvent(writer, "Int32", "int32", "ai32")
	generateArrayEvent(writer, "Int64", "int64", "ai64")
	generateArrayEvent(writer, "Float16", "float32", "af16")
	generateArrayEvent(writer, "Float32", "float32", "af32")
	generateArrayEvent(writer, "Float64", "float64", "af64")
	generateArrayEvent(writer, "Uint8", "uint8", "au8")
	generateArrayEvent(writer, "Uint16", "uint16", "au16")
	generateArrayEvent(writer, "Uint32", "uint32", "au32")
	generateArrayEvent(writer, "Uint64", "uint64", "au64")
	generateArrayEvent(writer, "UID", "[]byte", "au")
	generateArrayBeginEvent(writer, "ArrayBit", "Bit", "bab")
	generateArrayBeginEvent(writer, "ArrayFloat16", "Float16", "baf16")
	generateArrayBeginEvent(writer, "ArrayFloat32", "Float32", "baf32")
	generateArrayBeginEvent(writer, "ArrayFloat64", "Float64", "baf64")
	generateArrayBeginEvent(writer, "ArrayInt8", "Int8", "bai8")
	generateArrayBeginEvent(writer, "ArrayInt16", "Int16", "bai16")
	generateArrayBeginEvent(writer, "ArrayInt32", "Int32", "bai32")
	generateArrayBeginEvent(writer, "ArrayInt64", "Int64", "bai64")
	generateArrayBeginEvent(writer, "ArrayUID", "UID", "bau")
	generateArrayBeginEvent(writer, "ArrayUint8", "Uint8", "bau8")
	generateArrayBeginEvent(writer, "ArrayUint16", "Uint16", "bau16")
	generateArrayBeginEvent(writer, "ArrayUint32", "Uint32", "bau32")
	generateArrayBeginEvent(writer, "ArrayUint64", "Uint64", "bau64")
	generateArrayBeginEvent(writer, "ReferenceRemote", "ReferenceRemote", "brefr")
	generateArrayBeginEvent(writer, "ResourceID", "ResourceID", "brid")
	generateArrayBeginEvent(writer, "String", "String", "bs")

	generateOneArgEvent(writer, "BeginCustomBinary", "bcb", "customType", "uint64", "receiver.OnCustomBegin(events.ArrayTypeCustomBinary, customType)")
	generateOneArgEvent(writer, "BeginCustomText", "bct", "customType", "uint64", "receiver.OnCustomBegin(events.ArrayTypeCustomText, customType)")
	generateOneArgEvent(writer, "BeginMedia", "bmedia", "mediaType", "string", "receiver.OnMediaBegin(mediaType)")

	generateBasicEvent(writer, "Edge", "edge")
	generateBasicEvent(writer, "EndContainer", "e")
	generateBasicEvent(writer, "List", "l")
	generateBasicEvent(writer, "Map", "m")
	generateBasicEvent(writer, "Node", "node")
	generateBasicEvent(writer, "Null", "null")
	generateBasicEvent(writer, "Padding", "pad")
	generateBasicEvent(writer, "BeginDocument", "bd")
	generateBasicEvent(writer, "EndDocument", "ed")
}

func generateArrayDataEvent(writer io.Writer, eventName string, elementType string, functionName string) {
	argName := "elements"
	generateOneArgEvent(writer, "ArrayData"+eventName, functionName, argName, elementType,
		fmt.Sprintf("receiver.OnArrayData(array%vToBytes(safeArg))", eventName))
}

func generateArrayEvent(writer io.Writer, eventName string, elementType string, functionName string) {
	argName := "elements"
	generateOneArgEvent(writer, "Array"+eventName, functionName, argName, "[]"+elementType,
		fmt.Sprintf("receiver.OnArray(events.ArrayType%v, uint64(len(safeArg)), array%vToBytes(safeArg))", eventName, eventName))
}

func generateStringArrayEvent(writer io.Writer, arrayType string, functionName string) {
	eventName := arrayType
	argName := "str"
	generateOneArgEvent(writer, eventName, functionName, argName, "string",
		fmt.Sprintf("receiver.OnStringlikeArray(events.ArrayType%v, safeArg)", arrayType))
}

func generateArrayBeginEvent(writer io.Writer, eventName string, arrayType string, functionName string) {
	generateZeroArgEvent(writer, "Begin"+eventName, functionName, fmt.Sprintf("receiver.OnArrayBegin(events.ArrayType%v)", arrayType))
}

func generateIDEvent(writer io.Writer, eventName string, functionName string) {
	argName := "id"
	generateOneArgEvent(writer, eventName, functionName, argName, "string",
		fmt.Sprintf("receiver.On%v([]byte(safeArg))", eventName))
}

func generateBasicEvent(writer io.Writer, eventName string, functionName string) {
	generateZeroArgEvent(writer, eventName, functionName, fmt.Sprintf("receiver.On%v()", eventName))
}

func generateZeroArgEvent(writer io.Writer, eventName string, functionName string, invocation string) {
	functionUpper := strings.ToUpper(functionName)
	functionLower := strings.ToLower(functionName)

	writer.Write([]byte(fmt.Sprintf("type Event%v struct{ BaseEvent }\n\n", eventName)))
	writer.Write([]byte(fmt.Sprintf(`func %v() Event {
	return &Event%v{
		BaseEvent: ConstructEvent("%v", func(receiver events.DataEventReceiver) {
			%v
		}),
	}
}

`, functionUpper, eventName, functionLower, invocation)))
}

func generateOneArgEvent(writer io.Writer, eventName string, functionName string, argName string, argType string, invocation string) {
	functionUpper := strings.ToUpper(functionName)
	functionLower := strings.ToLower(functionName)

	writer.Write([]byte(fmt.Sprintf("type Event%v struct{ BaseEvent }\n\n", eventName)))
	writer.Write([]byte(fmt.Sprintf("func %v(%v %v) Event {", functionUpper, argName, argType)))
	if isArrayType(argType) {
		writer.Write([]byte(fmt.Sprintf(`
	if len(%v) == 0 {
		return &Event%v{
			BaseEvent: ConstructEvent("%v", func(receiver events.DataEventReceiver) {
				var safeArg %v
				%v
			}),
		}
	}
	v := copyOf(%v)
	var safeArg %v
	if v != nil {
		safeArg = v.(%v)
	}
`, argName, eventName, functionLower, argType, invocation, argName, argType, argType)))
	} else {
		writer.Write([]byte(fmt.Sprintf(`
	v := %v
	safeArg := v
`, argName)))
	}

	writer.Write([]byte(fmt.Sprintf(`
	return &Event%v{
		BaseEvent: ConstructEvent("%v", func(receiver events.DataEventReceiver) {
			%v
		}, safeArg),
	}
}

`, eventName, functionLower, invocation)))
}

func generateTwoArgEvent(writer io.Writer, eventName string, functionName string, arg1Name string, arg1Type string, arg2Name string, arg2Type string, invocation string) {
	functionUpper := strings.ToUpper(functionName)
	functionLower := strings.ToLower(functionName)

	writer.Write([]byte(fmt.Sprintf("type Event%v struct{ BaseEvent }\n\n", eventName)))
	writer.Write([]byte(fmt.Sprintf("func %v(%v %v, %v %v) Event {", functionUpper, arg1Name, arg1Type, arg2Name, arg2Type)))
	if isArrayType(arg2Type) {
		writer.Write([]byte(fmt.Sprintf(`
	if len(%v) == 0 {
		return &Event%v{
			BaseEvent: ConstructEvent("%v", func(receiver events.DataEventReceiver) {
				var safeArg %v
				%v
			}, %v),
		}
	}
	v := copyOf(%v)
	var safeArg %v
	if v != nil {
		safeArg = v.(%v)
	}
`, arg2Name, eventName, functionLower, arg2Type, invocation, arg1Name, arg2Name, arg2Type, arg2Type)))
	} else {
		writer.Write([]byte(fmt.Sprintf(`
	v := %v
	safeArg := v
`, arg2Name)))
	}

	writer.Write([]byte(fmt.Sprintf(`
	return &Event%v{
		BaseEvent: ConstructEvent("%v", func(receiver events.DataEventReceiver) {
			%v
		}, %v, safeArg),
	}
}

`, eventName, functionLower, invocation, arg1Name)))
}

func isArrayType(argType string) bool {
	return strings.HasPrefix(argType, "[]")
}

// ===========================================================================

var allEvents = test.Events{
	EvAB, EvACL, EvACM, EvAF16, EvAF32, EvAF64, EvAI16, EvAI32, EvAI64, EvAI8,
	EvAU, EvAU16, EvAU32, EvAU64, EvAU8, EvB, EvBAB, EvBAF16, EvBAF32, EvBAF64,
	EvBAI16, EvBAI32, EvBAI64, EvBAI8, EvBAU, EvBAU16, EvBAU32, EvBAU64, EvBAU8,
	EvBCB /*EvBCT,*/, EvBMEDIA, EvBRID, EvBS, EvCB, EvCM, EvCS /*EvCT,*/, EvE,
	EvEDGE, EvINF, EvL, EvM, EvMARK, EvMEDIA, EvN, EvNAN, EvNINF, EvNODE, EvNULL,
	EvPAD, EvREFL, EvREFR, EvRID, EvS, EvSI, EvSNAN, EvST, EvT, EvUID, EvV,
}

var (
	prefixes = map[string]test.Events{
		EvSI.Name():   {EvST, EvS, EvE},
		EvREFL.Name(): {EvMARK, EvN},
	}
	followups = map[string]test.Events{
		EvL.Name():      {EvE},
		EvM.Name():      {EvE},
		EvSI.Name():     {EvS, EvE},
		EvST.Name():     {EvS, EvE, EvN},
		EvNODE.Name():   {EvN, EvE},
		EvEDGE.Name():   {EvRID, EvRID, EvN, EvE},
		EvBAB.Name():    {EvACL, EvADB},
		EvBAF16.Name():  {EvACL, EvADF16},
		EvBAF32.Name():  {EvACL, EvADF32},
		EvBAF64.Name():  {EvACL, EvADF64},
		EvBAI16.Name():  {EvACL, EvADI16},
		EvBAI32.Name():  {EvACL, EvADI32},
		EvBAI64.Name():  {EvACL, EvADI64},
		EvBAI8.Name():   {EvACL, EvADI8},
		EvBAU16.Name():  {EvACL, EvADU16},
		EvBAU32.Name():  {EvACL, EvADU32},
		EvBAU64.Name():  {EvACL, EvADU64},
		EvBAU8.Name():   {EvACL, EvADU8},
		EvBAU.Name():    {EvACL, EvADU},
		EvBCB.Name():    {EvACL, EvADU8},
		EvBCT.Name():    {EvACL, EvADT},
		EvBMEDIA.Name(): {EvACL, EvADU8},
		EvBRID.Name():   {EvACL, EvADT},
		EvBS.Name():     {EvACL, EvADT},
		EvCM.Name():     {EvN},
		EvCS.Name():     {EvN},
		EvMARK.Name():   {EvN},
		EvPAD.Name():    {EvN},
	}

	lossyCTE = map[string]bool{
		EvACL.Name():    true,
		EvACM.Name():    true,
		EvBAB.Name():    true,
		EvBAF16.Name():  true,
		EvBAF32.Name():  true,
		EvBAF64.Name():  true,
		EvBAI16.Name():  true,
		EvBAI32.Name():  true,
		EvBAI64.Name():  true,
		EvBAI8.Name():   true,
		EvBAU16.Name():  true,
		EvBAU32.Name():  true,
		EvBAU64.Name():  true,
		EvBAU8.Name():   true,
		EvBAU.Name():    true,
		EvBCB.Name():    true,
		EvBCT.Name():    true,
		EvBMEDIA.Name(): true,
		EvBRID.Name():   true,
		EvBS.Name():     true,

		// CTE doesn't have this
		EvPAD.Name(): true,
	}

	lossyCBE = map[string]bool{
		// Chunked arrays may have been optimized
		EvACL.Name():    true,
		EvACM.Name():    true,
		EvBAB.Name():    true,
		EvBAF16.Name():  true,
		EvBAF32.Name():  true,
		EvBAF64.Name():  true,
		EvBAI16.Name():  true,
		EvBAI32.Name():  true,
		EvBAI64.Name():  true,
		EvBAI8.Name():   true,
		EvBAU16.Name():  true,
		EvBAU32.Name():  true,
		EvBAU64.Name():  true,
		EvBAU8.Name():   true,
		EvBAU.Name():    true,
		EvBCB.Name():    true,
		EvBCT.Name():    true,
		EvBMEDIA.Name(): true,
		EvBRID.Name():   true,
		EvBS.Name():     true,

		// CBE doesn't have these
		EvCM.Name(): true,
		EvCS.Name(): true,
		EvCT.Name(): true,
	}
)

func hasLossyCTE(events ...test.Event) bool {
	for _, event := range events {
		if lossyCTE[event.Name()] {
			return true
		}
	}
	return false
}

func hasLossyCBE(events ...test.Event) bool {
	for _, event := range events {
		if lossyCBE[event.Name()] {
			return true
		}
	}
	return false
}

func generateEventPrefixesAndFollowups(events ...test.Event) (eventSets []test.Events) {
	for _, event := range events {
		eventSet := []test.Event{}
		if pre, ok := prefixes[event.Name()]; ok {
			eventSet = append(eventSet, pre...)
		}
		eventSet = append(eventSet, event)
		if post, ok := followups[event.Name()]; ok {
			eventSet = append(eventSet, post...)
		}
		eventSets = append(eventSets, eventSet)
	}
	return
}

func complementaryEvents(events test.Events) test.Events {
	complementary := make(test.Events, 0, len(allEvents)/2)
	for _, event := range allEvents {
		for _, compareEvent := range events {
			if event == compareEvent {
				goto Skip
			}
		}
		complementary = append(complementary, event)
	Skip:
	}
	return complementary
}

var (
	EvAB     = test.AB([]bool{true})
	EvACL    = test.ACL(1)
	EvACM    = test.ACM(1)
	EvADB    = test.ADB([]bool{true})
	EvADF16  = test.ADF16([]float32{1})
	EvADF32  = test.ADF32([]float32{1})
	EvADF64  = test.ADF64([]float64{1})
	EvADI16  = test.ADI16([]int16{1})
	EvADI32  = test.ADI32([]int32{1})
	EvADI64  = test.ADI64([]int64{1})
	EvADI8   = test.ADI8([]int8{1})
	EvADT    = test.ADT("a")
	EvADU    = test.ADU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvADU16  = test.ADU16([]uint16{1})
	EvADU32  = test.ADU32([]uint32{1})
	EvADU64  = test.ADU64([]uint64{1})
	EvADU8   = test.ADU8([]uint8{1})
	EvAF16   = test.AF16([]float32{1})
	EvAF32   = test.AF32([]float32{1})
	EvAF64   = test.AF64([]float64{1})
	EvAI16   = test.AI16([]int16{1})
	EvAI32   = test.AI32([]int32{1})
	EvAI64   = test.AI64([]int64{1})
	EvAI8    = test.AI8([]int8{1})
	EvAU     = test.AU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvAU16   = test.AU16([]uint16{1})
	EvAU32   = test.AU32([]uint32{1})
	EvAU64   = test.AU64([]uint64{1})
	EvAU8    = test.AU8([]uint8{1})
	EvB      = test.B(true)
	EvBAB    = test.BAB()
	EvBAF16  = test.BAF16()
	EvBAF32  = test.BAF32()
	EvBAF64  = test.BAF64()
	EvBAI16  = test.BAI16()
	EvBAI32  = test.BAI32()
	EvBAI64  = test.BAI64()
	EvBAI8   = test.BAI8()
	EvBAU    = test.BAU()
	EvBAU16  = test.BAU16()
	EvBAU32  = test.BAU32()
	EvBAU64  = test.BAU64()
	EvBAU8   = test.BAU8()
	EvBCB    = test.BCB(0)
	EvBCT    = test.BCT(0)
	EvBMEDIA = test.BMEDIA("a/b")
	EvBRID   = test.BRID()
	EvBS     = test.BS()
	EvCB     = test.CB(0, []byte{1})
	EvCM     = test.CM("a")
	EvCS     = test.CS("a")
	EvCT     = test.CT(0, "a")
	EvE      = test.E()
	EvEDGE   = test.EDGE()
	EvINF    = test.N(math.Inf(1))
	EvL      = test.L()
	EvM      = test.M()
	EvMARK   = test.MARK("a")
	EvMEDIA  = test.MEDIA("a/b", []byte{1})
	EvN      = test.N(-1)
	EvNAN    = test.N(compact_float.QuietNaN())
	EvNINF   = test.N(math.Inf(-1))
	EvNODE   = test.NODE()
	EvNULL   = test.NULL()
	EvPAD    = test.PAD()
	EvREFL   = test.REFL("a")
	EvREFR   = test.REFR("a")
	EvRID    = test.RID("http://z.com")
	EvS      = test.S("a")
	EvSI     = test.SI("a")
	EvSNAN   = test.N(compact_float.SignalingNaN())
	EvST     = test.ST("a")
	EvT      = test.T(compact_time.AsCompactTime(time.Date(2020, time.Month(1), 1, 1, 1, 1, 1, time.UTC)))
	EvUID    = test.UID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	EvV      = test.V(version.ConciseEncodingVersion)
)
