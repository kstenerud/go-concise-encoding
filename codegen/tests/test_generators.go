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
	"os"
	"path/filepath"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/test_runner"
	"github.com/kstenerud/go-concise-encoding/version"
)

func generateTestFiles(projectDir string) {
	generateTestGenerators(filepath.Join(projectDir, "codegen/tests"))

	testsDir := filepath.Join(projectDir, "tests/suites/generated")

	generateRulesTestFiles(testsDir)

	writeTestFile(filepath.Join(testsDir, "cte-header-generated.cte"), generateCteHeaderTests()...)
	writeTestFile(filepath.Join(testsDir, "tlo-generated.cte"), generateTLOTests())
	writeTestFile(filepath.Join(testsDir, "list-generated.cte"), generateListTests())
	writeTestFile(filepath.Join(testsDir, "map-generated.cte"), generateMapKeyTests(), generateMapValueTests())
	writeTestFile(filepath.Join(testsDir, "edge-generated.cte"), generateEdgeSourceTests(), generateEdgeDescriptionTests(), generateEdgeDestinationTests())
	writeTestFile(filepath.Join(testsDir, "node-generated.cte"), generateNodeValueTests(), generateNodeChildTests())
	writeTestFile(filepath.Join(testsDir, "struct-generated.cte"), generateStructTemplateTests(), generateStructInstanceTests())
	generateArrayTestFiles(testsDir)
}

func generateRulesTestFiles(testsDir string) {
	prefixes := test.Events{EvBAB, EvBAF16, EvBAF32, EvBAF64, EvBAI16, EvBAI32, EvBAI64, EvBAI8,
		EvBAU, EvBAU16, EvBAU32, EvBAU64, EvBAU8, EvBCB, EvBCT, EvBMEDIA, EvBRID, EvBS}
	for _, prefix := range prefixes {
		filename := fmt.Sprintf("rules-%v-generated.cte", prefix.Name())
		writeTestFile(filepath.Join(testsDir, filename), generateRulesInvalidArrayEventsTests(prefix)...)
	}
}

func generateArrayTestFiles(testsDir string) {
	writeTestFile(filepath.Join(testsDir, "array-int8-generated.cte"), generateArrayInt8Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-int16-generated.cte"), generateArrayInt16Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-int32-generated.cte"), generateArrayInt32Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-int64-generated.cte"), generateArrayInt64Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-uint8-generated.cte"), generateArrayUint8Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-uint16-generated.cte"), generateArrayUint16Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-uint32-generated.cte"), generateArrayUint32Tests()...)
	writeTestFile(filepath.Join(testsDir, "array-uint64-generated.cte"), generateArrayUint64Tests()...)
}

func generateRulesInvalidArrayEventsTests(prefix test.Event) []*test_runner.UnitTest {
	var mustFail []*test_runner.MustFailTest
	for _, event := range complementaryEvents(test.Events{EvACL, EvACM}) {
		events := test.Events{prefix, event}
		mustFail = append(mustFail, newMustFailTest(testTypeEvents, events...))
	}
	name := fmt.Sprintf("Invlid %v Event Sequences", prefix.Name())
	return []*test_runner.UnitTest{newMustFailUnitTest(name, mustFail...)}
}

func generateTLOTests() *test_runner.UnitTest {
	prefix := test.Events{}
	suffix := test.Events{}
	invalidEvents := test.Events{EvV, EvE, EvACL, EvACM, EvREFL, EvSI}
	validEvents := complementaryEvents(invalidEvents)

	return generateEncodeDecodeTest("Top-level objects", prefix, suffix, validEvents, invalidEvents)
}

func generateListTests() *test_runner.UnitTest {
	prefix := test.Events{EvL}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("List", prefix, suffix, validEvents, invalidEvents)
}

func generateMapKeyTests() *test_runner.UnitTest {
	prefix := test.Events{EvM}
	suffix := test.Events{EvN, EvE}
	validEvents := test.Events{EvB, EvBRID, EvBS, EvCM, EvCS, EvINF, EvN, EvNINF, EvPAD, EvRID, EvS, EvT, EvUID}
	invalidEvents := complementaryEvents(validEvents)

	return generateEncodeDecodeTest("Map Key", prefix, suffix, validEvents, invalidEvents)
}

func generateMapValueTests() *test_runner.UnitTest {
	prefix := test.Events{EvM, EvN}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Map Value", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeSourceTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE}
	suffix := test.Events{EvN, EvN, EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvNULL, EvREFL}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Source", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeDescriptionTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE, EvN}
	suffix := test.Events{EvN, EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Description", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeDestinationTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE, EvN, EvN}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvNULL, EvREFL}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Destination", prefix, suffix, validEvents, invalidEvents)
}

func generateNodeValueTests() *test_runner.UnitTest {
	prefix := test.Events{EvNODE}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Node Value", prefix, suffix, validEvents, invalidEvents)
}

func generateNodeChildTests() *test_runner.UnitTest {
	prefix := test.Events{EvNODE, EvNULL}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Node Child", prefix, suffix, validEvents, invalidEvents)
}

func generateStructTemplateTests() *test_runner.UnitTest {
	prefix := test.Events{EvST}
	suffix := test.Events{EvE, EvN}
	validEvents := test.Events{EvB, EvBRID, EvBS, EvCM, EvCS, EvINF, EvN, EvNINF, EvPAD, EvRID, EvS, EvT, EvUID}
	invalidEvents := complementaryEvents(validEvents)

	return generateEncodeDecodeTest("Struct Template", prefix, suffix, validEvents, invalidEvents)
}

func generateStructInstanceTests() *test_runner.UnitTest {
	prefix := test.Events{EvST, EvS, EvE, EvSI}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL, EvSI, EvST}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Struct Instance", prefix, suffix, validEvents, invalidEvents)
}

func generateEncodeDecodeTest(name string, prefix test.Events, suffix test.Events, validEvents test.Events, invalidEvents test.Events) *test_runner.UnitTest {
	mustSucceed := []*test_runner.MustSucceedTest{}
	mustFail := []*test_runner.MustFailTest{}
	config := configuration.DefaultCTEEncoderConfiguration()

	for _, eventSet := range generateEventPrefixesAndFollowups(validEvents...) {
		events := append(prefix, eventSet...)
		events = append(events, suffix...)
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, events...))
	}

	for _, event := range invalidEvents {
		events := append(prefix, event)
		events = append(events, suffix...)
		mustFail = append(mustFail, newMustFailTest(testTypeEvents, events...))
	}

	return newUnitTest(name, mustSucceed, mustFail)
}

func generateCteHeaderTests() []*test_runner.UnitTest {
	wrongSentinelFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		if i == 'c' || i == 'C' {
			continue
		}
		wrongSentinelFailureTests = append(wrongSentinelFailureTests, newCustomMustFailTest(fmt.Sprintf("%c%v 0", rune(i), version.ConciseEncodingVersion)))
	}
	wrongSentinelTest := newUnitTest("Wrong sentinel", nil, wrongSentinelFailureTests)

	wrongVersionCharFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		wrongVersionCharFailureTests = append(wrongVersionCharFailureTests, newCustomMustFailTest(fmt.Sprintf("c%c 0", rune(i))))
	}
	wrongVersionCharTest := newUnitTest("Wrong version character", nil, wrongVersionCharFailureTests)

	wrongVersionFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		// TODO: Remove i == 1 upon release
		if i == version.ConciseEncodingVersion || i == 1 {
			continue
		}
		wrongVersionFailureTests = append(wrongVersionFailureTests, newCustomMustFailTest(fmt.Sprintf("c%v 0", i)))
	}
	wrongVersionTest := newUnitTest("Wrong version", nil, wrongVersionFailureTests)

	return []*test_runner.UnitTest{wrongSentinelTest, wrongVersionCharTest, wrongVersionTest}
}

// ===========================================================================

func newUnitTest(name string, mustSucceed []*test_runner.MustSucceedTest, mustFail []*test_runner.MustFailTest) *test_runner.UnitTest {
	unitTest := &test_runner.UnitTest{
		Name: name,
	}

	if mustSucceed != nil {
		unitTest.MustSucceed = mustSucceed
	}
	if mustFail != nil {
		unitTest.MustFail = mustFail
	}

	return unitTest
}

func newMustSucceedUnitTest(name string, mustSucceed ...*test_runner.MustSucceedTest) *test_runner.UnitTest {
	return &test_runner.UnitTest{
		Name:        name,
		MustSucceed: mustSucceed,
	}
}

func newMustFailUnitTest(name string, mustFail ...*test_runner.MustFailTest) *test_runner.UnitTest {
	return &test_runner.UnitTest{
		Name:     name,
		MustFail: mustFail,
	}
}

func newMustSucceedTest(config *configuration.CTEEncoderConfiguration, events ...test.Event) *test_runner.MustSucceedTest {
	hasFromCBE := canConvertFromCBE(config, events...)
	hasToCBE := canConvertToCBE(config, events...)
	cbeDocument := generateCBE(events...)
	cbe := []byte{}
	fromCBE := []byte{}
	toCBE := []byte{}
	if hasFromCBE && hasToCBE {
		cbe = cbeDocument
	} else if hasFromCBE {
		fromCBE = cbeDocument
	} else if hasToCBE {
		toCBE = cbeDocument
	}

	hasFromCTE := canConvertFromCTE(config, events...)
	hasToCTE := canConvertToCTE(config, events...)
	cteDocument := generateCTE(config, events...)
	cte := ""
	fromCTE := ""
	toCTE := ""
	if hasFromCTE && hasToCTE {
		cte = cteDocument
	} else if hasFromCTE {
		fromCTE = cteDocument
	} else if hasToCTE {
		toCTE = cteDocument
	}

	return &test_runner.MustSucceedTest{
		BaseTest: test_runner.BaseTest{
			CBE:    cbe,
			CTE:    cte,
			Events: stringifyEvents(events...),
		},
		FromCBE: fromCBE,
		ToCBE:   toCBE,
		FromCTE: fromCTE,
		ToCTE:   toCTE,
	}
}

func newMustFailTest(testType testType, events ...test.Event) *test_runner.MustFailTest {
	switch testType {
	case testTypeCbe:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: generateCBE(events...)}}
	case testTypeCte:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, events...)}}
	case testTypeEvents:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{Events: stringifyEvents(events...)}}
	default:
		panic(fmt.Errorf("%v: unknown mustFail test type", testType))
	}
}

func newCustomMustFailTest(cteContents string) *test_runner.MustFailTest {
	return &test_runner.MustFailTest{
		BaseTest: test_runner.BaseTest{
			CTE:         cteContents,
			RawDocument: true,
		},
	}
}

func writeTestFile(path string, tests ...*test_runner.UnitTest) {
	ceVersion := version.ConciseEncodingVersion
	suite := &test_runner.TestSuite{
		Type: test_runner.TestSuiteType{
			Identifier: "ce-test",
			Version:    1,
		},
		CEVersion: &ceVersion,
		Tests:     tests,
	}

	config := configuration.DefaultCTEMarshalerConfiguration()
	config.Iterator.FieldNameStyle = configuration.FieldNameSnakeCase
	config.Encoder.DefaultNumericFormats.Array.Uint8 = configuration.CTEEncodingFormatHexadecimalZeroFilled
	config.DebugPanics = true
	document, err := ce.MarshalToCTEDocument(suite, &config)
	if err != nil {
		panic(err)
	}

	comment := "// GENERATED FILE, DO NOT EDIT!\n// Generated by https://github.com/kstenerud/go-concise-encoding/tree/master/codegen/tests\n"
	commentedDocument := make([]byte, 0, len(document)+len(comment))
	commentedDocument = append(commentedDocument, document[:3]...)
	commentedDocument = append(commentedDocument, comment...)
	commentedDocument = append(commentedDocument, document[3:]...)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(commentedDocument)
}
