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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/codegen/standard"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
)

const codePath = "test"
const parserBasePath = "test/event_parser"

var imports = []*standard.Import{
	&standard.Import{LocalName: "compact_time", Import: "github.com/kstenerud/go-compact-time"},
	&standard.Import{LocalName: "", Import: "github.com/kstenerud/go-concise-encoding/events"},
}

func GenerateCode(projectDir string) {
	generateAntlrCode(projectDir)
	generateTestFiles(projectDir)

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

	generateEvents(writer)
}

func generateCbe(events ...test.Event) []byte {
	buffer := bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(&buffer)
	encoder.OnBeginDocument()
	encoder.OnVersion(version.ConciseEncodingVersion)
	for _, event := range events {
		event.Invoke(encoder)
	}
	encoder.OnEndDocument()
	result := buffer.Bytes()
	return result[2:]
}

func generateCte(events ...test.Event) string {
	buffer := bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(&buffer)
	encoder.OnBeginDocument()
	encoder.OnVersion(version.ConciseEncodingVersion)
	for _, event := range events {
		event.Invoke(encoder)
	}
	encoder.OnEndDocument()
	result := buffer.String()
	return result[3:]
}

const (
	testTypeCbe    = "cbe"
	testTypeCte    = "cte"
	testTypeEvents = "events"
)

func stringifyEvents(events ...test.Event) string {
	sb := strings.Builder{}
	for i, event := range events {
		if i > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(event.String())
	}
	return sb.String()
}

func generateMustFailTest(testType string, events ...test.Event) map[string]interface{} {
	switch testType {
	case testTypeCbe:
		return map[string]interface{}{testType: generateCbe(events...)}
	case testTypeCte:
		return map[string]interface{}{testType: generateCte(events...)}
	case testTypeEvents:
		return map[string]interface{}{testType: stringifyEvents(events...)}
	default:
		panic(fmt.Errorf("%v: unknown mustFail test type", testType))
	}
}

func generateCustomMustFailTest(cteContents string) map[string]interface{} {
	return map[string]interface{}{
		"rawdocument": true,
		testTypeCte:   cteContents,
	}
}

func generateTest(name string, mustSucceed []interface{}, mustFail []interface{}) interface{} {
	m := map[string]interface{}{
		"name": name,
	}

	if mustSucceed != nil {
		m["mustSucceed"] = mustSucceed
	}
	if mustFail != nil {
		m["mustFail"] = mustFail
	}
	return m
}

func writeTestFile(path string, tests ...interface{}) {
	m := map[string]interface{}{
		"type": map[string]interface{}{
			"identifier": "ce-test",
			"version":    1,
		},
		"ceversion": version.ConciseEncodingVersion,
		"tests":     tests,
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	opts := options.DefaultCTEMarshalerOptions()
	opts.Encoder.DefaultNumericFormats.Array.Uint8 = options.CTEEncodingFormatHexadecimalZeroFilled
	if err := ce.MarshalCTE(m, f, &opts); err != nil {
		panic(err)
	}
}

func generateCteHeaderTests(path string) {
	wrongSentinelFailureTests := []interface{}{}
	for i := 0; i < 0x100; i++ {
		if i == 'c' || i == 'C' {
			continue
		}
		wrongSentinelFailureTests = append(wrongSentinelFailureTests, generateCustomMustFailTest(fmt.Sprintf("%c%v 0", rune(i), version.ConciseEncodingVersion)))
	}
	wrongSentinelTest := generateTest("Wrong sentinel", nil, wrongSentinelFailureTests)

	wrongVersionCharFailureTests := []interface{}{}
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		wrongVersionCharFailureTests = append(wrongVersionCharFailureTests, generateCustomMustFailTest(fmt.Sprintf("c%c 0", rune(i))))
	}
	wrongVersionCharTest := generateTest("Wrong version character", nil, wrongVersionCharFailureTests)

	wrongVersionFailureTests := []interface{}{}
	for i := 0; i < 0x100; i++ {
		// TODO: Remove i == 1 upon release
		if i == version.ConciseEncodingVersion || i == 1 {
			continue
		}
		wrongVersionFailureTests = append(wrongVersionFailureTests, generateCustomMustFailTest(fmt.Sprintf("c%v 0", i)))
	}
	wrongVersionTest := generateTest("Wrong version", nil, wrongVersionFailureTests)

	writeTestFile(path, wrongSentinelTest, wrongVersionCharTest, wrongVersionTest)
}

func generateRulesTests(path string) {
	noTests := generateTest("No tests", nil, []interface{}{})

	writeTestFile(path, noTests)
}

func generateTestFiles(projectDir string) {
	testsDir := filepath.Join(projectDir, "tests")

	generateCteHeaderTests(filepath.Join(testsDir, "cte-generated-do-not-edit.cte"))
	generateRulesTests(filepath.Join(testsDir, "rules-generated-do-not-edit.cte"))
}

func generateAntlrCode(projectDir string) {
	javaPath, err := exec.LookPath("java")
	if err != nil {
		panic(err)
	}
	dstPath := filepath.Join(projectDir, parserBasePath, "parser")
	if err := os.RemoveAll(dstPath); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		panic(err)
	}

	antlrPath := filepath.Join(projectDir, "codegen", "antlr-4.10.1-complete.jar")
	lexerPath := filepath.Join(projectDir, "codegen", "tests", "CEEventLexer.g4")
	parserPath := filepath.Join(projectDir, "codegen", "tests", "CEEventParser.g4")
	cmd := exec.Command(
		javaPath,
		"-cp", antlrPath,
		"org.antlr.v4.Tool",
		"-o", dstPath,
		"-Dlanguage=Go",
		lexerPath, parserPath,
	)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func generateEvents(writer io.Writer) {
	generateEventWithValue(writer, "ArrayBit", "ab", "elements", "[]bool", "receiver.OnArray(events.ArrayTypeBit, uint64(len(safeArg)), arrayBitsToBytes(safeArg))")
	generateEventWithValue(writer, "ArrayDataBit", "adb", "elements", "[]bool", "receiver.OnArrayData(arrayBitsToBytes(safeArg))")

	generateEventWithValue(writer, "Version", "v", "version", "uint64", "receiver.OnVersion(safeArg)")
	generateEventWithValue(writer, "Boolean", "b", "value", "bool", "receiver.OnBoolean(safeArg)")
	generateEventWithValue(writer, "Time", "t", "value", "compact_time.Time", "receiver.OnTime(safeArg)")
	generateEventWithValue(writer, "UID", "uid", "value", "[]byte", "receiver.OnUID(safeArg)")

	generateEventWithValue(writer, "ArrayChunkMore", "acm", "length", "uint64", "receiver.OnArrayChunk(safeArg, true)")
	generateEventWithValue(writer, "ArrayChunkLast", "acl", "length", "uint64", "receiver.OnArrayChunk(safeArg, false)")
	generateEventWithValue(writer, "CommentMultiline", "cm", "comment", "string", "receiver.OnComment(true, []byte(safeArg))")
	generateEventWithValue(writer, "CommentSingleLine", "cs", "comment", "string", "receiver.OnComment(false, []byte(safeArg))")
	generateEventWithValue(writer, "CustomBinary", "cb", "elements", "[]byte", "receiver.OnArray(events.ArrayTypeCustomBinary, uint64(len(safeArg)), safeArg)")

	generateIDEvent(writer, "Marker", "mark")
	generateIDEvent(writer, "ReferenceLocal", "refl")
	generateIDEvent(writer, "StructInstance", "si")
	generateIDEvent(writer, "StructTemplate", "st")

	generateStringArrayEvent(writer, "String", "s")
	generateStringArrayEvent(writer, "CustomText", "ct")
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
	generateArrayBeginEvent(writer, "CustomBinary", "CustomBinary", "bcb")
	generateArrayBeginEvent(writer, "CustomText", "CustomText", "bct")
	generateArrayBeginEvent(writer, "Media", "Media", "bmedia")
	generateArrayBeginEvent(writer, "ReferenceRemote", "ReferenceRemote", "brefr")
	generateArrayBeginEvent(writer, "ResourceID", "ResourceID", "brid")
	generateArrayBeginEvent(writer, "String", "String", "bs")

	generateBasicEvent(writer, "Edge", "edge")
	generateBasicEvent(writer, "End", "e")
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
	generateEventWithValue(writer, "ArrayData"+eventName, functionName, argName, elementType,
		fmt.Sprintf("receiver.OnArrayData(array%vToBytes(safeArg))", eventName))
}

func generateArrayEvent(writer io.Writer, eventName string, elementType string, functionName string) {
	argName := "elements"
	generateEventWithValue(writer, "Array"+eventName, functionName, argName, "[]"+elementType,
		fmt.Sprintf("receiver.OnArray(events.ArrayType%v, uint64(len(safeArg)), array%vToBytes(safeArg))", eventName, eventName))
}

func generateStringArrayEvent(writer io.Writer, arrayType string, functionName string) {
	eventName := arrayType
	argName := "str"
	generateEventWithValue(writer, eventName, functionName, argName, "string",
		fmt.Sprintf("receiver.OnStringlikeArray(events.ArrayType%v, safeArg)", arrayType))
}

func generateArrayBeginEvent(writer io.Writer, eventName string, arrayType string, functionName string) {
	generateZeroArgEvent(writer, "Begin"+eventName, functionName, fmt.Sprintf("receiver.OnArrayBegin(events.ArrayType%v)", arrayType))
}

func generateIDEvent(writer io.Writer, eventName string, functionName string) {
	argName := "id"
	generateEventWithValue(writer, eventName, functionName, argName, "string",
		fmt.Sprintf("receiver.On%v([]byte(safeArg))", eventName))
}

func generateBasicEvent(writer io.Writer, eventName string, functionName string) {
	generateZeroArgEvent(writer, eventName, functionName, fmt.Sprintf("receiver.On%v()", eventName))
}

func generateZeroArgEvent(writer io.Writer, eventName string, functionName string, invocation string) {
	functionUpper := strings.ToUpper(functionName)
	functionLower := strings.ToLower(functionName)

	writer.Write([]byte(fmt.Sprintf("type Event%v struct{ EventWithValue }\n\n", eventName)))
	writer.Write([]byte(fmt.Sprintf(`func %v() Event {
	return &Event%v{
		EventWithValue: ConstructEventWithValue("%v", NoValue, func(receiver events.DataEventReceiver) {
			%v
		}),
	}
}

`, functionUpper, eventName, functionLower, invocation)))
}

func generateEventWithValue(writer io.Writer, eventName string, functionName string, argName string, argType string, invocation string) {
	functionUpper := strings.ToUpper(functionName)
	functionLower := strings.ToLower(functionName)

	writer.Write([]byte(fmt.Sprintf("type Event%v struct{ EventWithValue }\n\n", eventName)))
	writer.Write([]byte(fmt.Sprintf("func %v(%v %v) Event {", functionUpper, argName, argType)))
	if isArrayType(argType) {
		writer.Write([]byte(fmt.Sprintf(`
	v := copyOf(%v)
	if len(%v) == 0 {
		v = NoValue
	}
	var safeArg %v
	if v != nil {
		safeArg = v.(%v)
	}
`, argName, argName, argType, argType)))
	} else {
		writer.Write([]byte(fmt.Sprintf(`
	v := %v
	safeArg := v
`, argName)))
	}

	writer.Write([]byte(fmt.Sprintf(`
	return &Event%v{
		EventWithValue: ConstructEventWithValue("%v", v, func(receiver events.DataEventReceiver) {
			%v
		}),
	}
}

`, eventName, functionLower, invocation)))
}

func isArrayType(argType string) bool {
	return strings.HasPrefix(argType, "[]")
}
