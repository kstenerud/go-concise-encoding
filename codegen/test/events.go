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

package test

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
)

const codePath = "test"

var eventsImports = []*common.Import{
	{As: "compact_time", Import: "github.com/kstenerud/go-compact-time"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/ce/events"},
}

func generateEvents(projectDir string) {
	common.GenerateGoFile(filepath.Join(projectDir, codePath), eventsImports, func(writer io.Writer) {
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
	})
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
