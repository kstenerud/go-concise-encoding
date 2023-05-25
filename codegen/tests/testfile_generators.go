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
	"path/filepath"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/test_runner"
	"github.com/kstenerud/go-concise-encoding/version"
)

func generateTestFiles(projectDir string) {
	testsDir := filepath.Join(projectDir, "tests/suites/generated")

	generateRulesTestFiles(testsDir)

	common.GenerateTestFile(filepath.Join(testsDir, "cte-header-generated.cte"), generateCteHeaderTests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "tlo-generated.cte"), generateTLOTests())
	common.GenerateTestFile(filepath.Join(testsDir, "list-generated.cte"), generateListTests())
	common.GenerateTestFile(filepath.Join(testsDir, "map-generated.cte"), generateMapKeyTests(), generateMapValueTests())
	common.GenerateTestFile(filepath.Join(testsDir, "edge-generated.cte"), generateEdgeSourceTests(), generateEdgeDescriptionTests(), generateEdgeDestinationTests())
	common.GenerateTestFile(filepath.Join(testsDir, "node-generated.cte"), generateNodeValueTests(), generateNodeChildTests())
	common.GenerateTestFile(filepath.Join(testsDir, "record-generated.cte"), generateRecordTests())
	generateArrayTestFiles(testsDir)
}

func generateRulesTestFiles(testsDir string) {
	prefixes := test.Events{EvBAB, EvBAF16, EvBAF32, EvBAF64, EvBAI16, EvBAI32, EvBAI64, EvBAI8,
		EvBAU, EvBAU16, EvBAU32, EvBAU64, EvBAU8, EvBCB, EvBCT, EvBMEDIA, EvBRID, EvBS}
	for _, prefix := range prefixes {
		filename := fmt.Sprintf("rules-%v-generated.cte", prefix.Name())
		common.GenerateTestFile(filepath.Join(testsDir, filename), generateRulesInvalidArrayEventsTests(prefix)...)
	}
}

func generateArrayTestFiles(testsDir string) {
	common.GenerateTestFile(filepath.Join(testsDir, "array-int8-generated.cte"), generateArrayInt8Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-int16-generated.cte"), generateArrayInt16Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-int32-generated.cte"), generateArrayInt32Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-int64-generated.cte"), generateArrayInt64Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-uint8-generated.cte"), generateArrayUint8Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-uint16-generated.cte"), generateArrayUint16Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-uint32-generated.cte"), generateArrayUint32Tests()...)
	common.GenerateTestFile(filepath.Join(testsDir, "array-uint64-generated.cte"), generateArrayUint64Tests()...)

	common.GenerateTestFile(filepath.Join(testsDir, "array-float32-generated.cte"), generateArrayFloat32Tests()...)
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
	invalidEvents := test.Events{EvV, EvE, EvACL, EvACM, EvREFL}
	validEvents := complementaryEvents(invalidEvents)

	return generateEncodeDecodeTest("Top-level objects", prefix, suffix, validEvents, invalidEvents)
}

func generateListTests() *test_runner.UnitTest {
	prefix := test.Events{EvL}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvRT}
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
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Map Value", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeSourceTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE}
	suffix := test.Events{EvN, EvN, EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvNULL, EvREFL, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Source", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeDescriptionTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE, EvN}
	suffix := test.Events{EvN, EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Description", prefix, suffix, validEvents, invalidEvents)
}

func generateEdgeDestinationTests() *test_runner.UnitTest {
	prefix := test.Events{EvEDGE, EvN, EvN}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvNULL, EvREFL, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Edge Destination", prefix, suffix, validEvents, invalidEvents)
}

func generateNodeValueTests() *test_runner.UnitTest {
	prefix := test.Events{EvNODE}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Node Value", prefix, suffix, validEvents, invalidEvents)
}

func generateNodeChildTests() *test_runner.UnitTest {
	prefix := test.Events{EvNODE, EvNULL}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Node Child", prefix, suffix, validEvents, invalidEvents)
}

func generateRecordTests() *test_runner.UnitTest {
	prefix := test.Events{EvREC}
	suffix := test.Events{EvE}
	invalidEvents := test.Events{EvV, EvACL, EvACM, EvREFL, EvREC, EvRT}
	validEvents := complementaryEvents(append(invalidEvents, EvE))

	return generateEncodeDecodeTest("Record", prefix, suffix, validEvents, invalidEvents)
}

func generateEncodeDecodeTest(name string, prefix test.Events, suffix test.Events, validEvents test.Events, invalidEvents test.Events) *test_runner.UnitTest {
	mustSucceed := []*test_runner.MustSucceedTest{}
	mustFail := []*test_runner.MustFailTest{}
	config := configuration.DefaultCTEEncoderConfiguration()

	for _, eventSet := range generateEventPrefixesAndFollowups(validEvents...) {
		events := append(prefix, eventSet...)
		events = append(events, suffix...)
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, events...))
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
