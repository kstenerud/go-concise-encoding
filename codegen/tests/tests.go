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
	"os"
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
)

func generateTestFiles(projectDir string) {
	testsDir := filepath.Join(projectDir, "tests")

	generateCteHeaderTests(filepath.Join(testsDir, "cte-generated-do-not-edit.cte"))
	generateRulesTests(filepath.Join(testsDir, "rules-generated-do-not-edit.cte"))
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

// ===========================================================================

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
