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
	"math"
	"math/big"
	"os"
	"path/filepath"

	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/codegen/standard"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
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
	writeTestFile(filepath.Join(testsDir, "array-int8-generated.cte"), generateArrayInt8Tests()...)
}

var testsImports = []*standard.Import{
	// {LocalName: "compact_time", Import: "github.com/kstenerud/go-compact-time"},
	// {LocalName: "", Import: "github.com/kstenerud/go-concise-encoding/ce/events"},
}

func generateTestGenerators(basePath string) {
	generatedFilePath := standard.GetGeneratedCodePath(basePath)
	writer, err := os.Create(generatedFilePath)
	standard.PanicIfError(err, "could not open %s", generatedFilePath)
	defer writer.Close()
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Errorf("error while generating %v: %v", generatedFilePath, e))
		}
	}()

	standard.WriteHeader(writer, "tests", testsImports)
	generateArrayTestGenerator(writer)
}

func generateArrayTestGenerator(writer io.Writer) {

}

func generateRulesTestFiles(testsDir string) {
	prefixes := test.Events{EvBAB, EvBAF16, EvBAF32, EvBAF64, EvBAI16, EvBAI32, EvBAI64, EvBAI8,
		EvBAU, EvBAU16, EvBAU32, EvBAU64, EvBAU8, EvBCB, EvBCT, EvBMEDIA, EvBRID, EvBS}
	for _, prefix := range prefixes {
		filename := fmt.Sprintf("rules-%v-generated.cte", prefix.Name())
		writeTestFile(filepath.Join(testsDir, filename), generateRulesInvalidArrayEventsTests(prefix)...)
	}
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

// Name:   Int8
// Type:   int8
// Fields: Int8
// Consts: MaxInt8, MinInt8
// Events: AI8, BAI8, ADI8
// Can probably just do "Int8", "I8" for int types
func generateArrayInt8Tests() []*test_runner.UnitTest {
	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []int8
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(&config, AI8(nil)),
		newMustSucceedTest(&config, BAI8(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, int8(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, AI8(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI8(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, int8(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI8(), ACL(uint64(i)), ADI8(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := math.MaxInt8 / 31
	for i := math.MinInt8; i < math.MaxInt8; i += multiple {
		contents = append(contents, int8(math.MinInt8+i))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI8(), ACL(uint64(len(contents))), ADI8(contents)))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 2, 8, 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = math.MaxUint8 / 7
	for i := math.MinInt8; i < math.MaxInt8; i += multiple {
		contents = append(contents, int8(math.MinInt8+i))
	}
	config.DefaultNumericFormats.Array.Int8 = configuration.CTEEncodingFormatBinary
	t := newMustSucceedTest(&config, AI8(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Int8 = configuration.CTEEncodingFormatOctal
	t = newMustSucceedTest(&config, AI8(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Int8 = configuration.CTEEncodingFormatHexadecimal
	t = newMustSucceedTest(&config, AI8(contents))
	mustSucceed = append(mustSucceed, t)
	config = configuration.DefaultCTEEncoderConfiguration()
	unitTests = append(unitTests, newMustSucceedUnitTest("Base 2, 8, 16", mustSucceed...))

	// Chunking
	mustSucceed = nil
	events := []test.Event{BAI8()}
	for i := 0; i < 7; i++ {
		events = append(events, ACM(uint64(i)))
		if i > 0 {
			events = append(events, ADI8(contents[:i]))
		}
	}
	events = append(events, ACL(0))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI8(), ACM(4), ADI8(contents[:4]), ACL(3), ADI8(contents[:3])))
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunking Variations", mustSucceed...))

	// Edge-case element values
	contents = contents[:0]
	for _, v := range intEdgeValues {
		contents = append(contents, int8(v))
	}
	for _, v := range uintEdgeValues {
		contents = append(contents, int8(v))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Edge Case Element Values",
		newMustSucceedTest(&config, BAI8(), ACL(uint64(len(contents))), ADI8(contents))))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Truncated Array
	mustFail = nil
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI8()))
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI8(), ACL(1)))
	for i := 2; i <= 10; i++ {
		contents = contents[:i/2]
		mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI8(), ACL(uint64(i)), ADI8(contents)))
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: truncate(generateCBE(AI8(contents)), 1)}})
	}
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAI8())}})
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAI8(), ACL(uint64(2)), ADI8(contents[:1]))}})
	unitTests = append(unitTests, newMustFailUnitTest("Truncated Array", mustFail...))

	// Element value out of range
	bigValue := big.NewInt(math.MaxInt8)
	bigValue.Add(bigValue, big.NewInt(1))
	smallValue := big.NewInt(math.MinInt8)
	smallValue.Sub(smallValue, big.NewInt(1))
	smallValue.Neg(smallValue)
	unitTests = append(unitTests, newMustFailUnitTest(
		"Element value out of range",
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8b %b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8o %o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8x %x|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 %d|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 0b%b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 0o%o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 0x%x|", bigValue)}},

		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8b -%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8o -%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8x -%x|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 -%d|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 -0b%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 -0o%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8 -0x%x|", smallValue)}},
	))

	// Numeric digit out of range
	mustFail = nil
	for iB, base := range bases {
		for iV, v := range outOfRange {
			if iV >= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range bases {
		for _, special := range nonIntSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	// Float value in int array
	mustFail = nil
	for _, base := range bases {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i8%v 1.2|", base)}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Float value in int array", mustFail...))

	return unitTests
}

func truncate(data []byte, count int) []byte {
	return data[:len(data)-count]
}

var bases = []string{
	"b",
	"o",
	"",
	"x",
}

var outOfRange = []string{
	"12",
	"18",
	"1a",
}

var nonIntSpecials = []string{
	"null",
	"true",
	"false",
	"nan",
	"snan",
	"inf",
	"-inf",
}

var nonFloatSpecials = []string{
	"null",
	"true",
	"false",
}

var intEdgeValues = []int64{
	0,
	1,
	0x7f,
	0x80,
	0x81,
	0xff,
	0x100,
	0x101,
	0x7fff,
	0x8000,
	0x8001,
	0xffff,
	0x10000,
	0x10001,
	0x7fffffff,
	0x80000000,
	0x80000001,
	0xffffffff,
	0x100000000,
	0x100000001,
	0x7fffffffffffffff,
	-1,
	-0x7f,
	-0x80,
	-0x81,
	-0xff,
	-0x100,
	-0x101,
	-0x7fff,
	-0x8000,
	-0x8001,
	-0xffff,
	-0x10000,
	-0x10001,
	-0x7fffffff,
	-0x80000000,
	-0x80000001,
	-0xffffffff,
	-0x100000000,
	-0x100000001,
	-0x7fffffffffffffff,
	-0x8000000000000000,
}
var uintEdgeValues = []uint64{
	0x8000000000000000,
	0x8000000000000001,
	0xffffffffffffffff,
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

	return generateUnitTest(name, mustSucceed, mustFail)
}

func generateCteHeaderTests() []*test_runner.UnitTest {
	wrongSentinelFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		if i == 'c' || i == 'C' {
			continue
		}
		wrongSentinelFailureTests = append(wrongSentinelFailureTests, generateCustomMustFailTest(fmt.Sprintf("%c%v 0", rune(i), version.ConciseEncodingVersion)))
	}
	wrongSentinelTest := generateUnitTest("Wrong sentinel", nil, wrongSentinelFailureTests)

	wrongVersionCharFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		if i >= '0' && i <= '9' {
			continue
		}
		wrongVersionCharFailureTests = append(wrongVersionCharFailureTests, generateCustomMustFailTest(fmt.Sprintf("c%c 0", rune(i))))
	}
	wrongVersionCharTest := generateUnitTest("Wrong version character", nil, wrongVersionCharFailureTests)

	wrongVersionFailureTests := []*test_runner.MustFailTest{}
	for i := 0; i < 0x100; i++ {
		// TODO: Remove i == 1 upon release
		if i == version.ConciseEncodingVersion || i == 1 {
			continue
		}
		wrongVersionFailureTests = append(wrongVersionFailureTests, generateCustomMustFailTest(fmt.Sprintf("c%v 0", i)))
	}
	wrongVersionTest := generateUnitTest("Wrong version", nil, wrongVersionFailureTests)

	return []*test_runner.UnitTest{wrongSentinelTest, wrongVersionCharTest, wrongVersionTest}
}

// ===========================================================================

func generateUnitTest(name string, mustSucceed []*test_runner.MustSucceedTest, mustFail []*test_runner.MustFailTest) *test_runner.UnitTest {
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

func generateCustomMustFailTest(cteContents string) *test_runner.MustFailTest {
	return &test_runner.MustFailTest{
		BaseTest: test_runner.BaseTest{
			CTE:         cteContents,
			RawDocument: true,
		},
	}
}

func generateCBE(events ...test.Event) []byte {
	if mustNotConvertToCBE(events...) {
		return []byte{}
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("in events [%v]: %w", test.Events(events), v))
			default:
				panic(fmt.Errorf("in events [%v]: %v", test.Events(events), v))
			}
		}
	}()

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

func generateCTE(config *configuration.CTEEncoderConfiguration, events ...test.Event) string {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("in events [%v]: %w", test.Events(events), v))
			default:
				panic(fmt.Errorf("in events [%v]: %v", test.Events(events), v))
			}
		}
	}()

	buffer := bytes.Buffer{}
	encoder := cte.NewEncoder(config)
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

type testType int

const (
	testTypeCbe testType = 1 << iota
	testTypeCte
	testTypeEvents
)

func stringifyEvents(events ...test.Event) (stringified []string) {
	for _, event := range events {
		stringified = append(stringified, event.String())
	}
	return
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

func AB(v []bool) test.Event              { return test.AB(v) }
func ACL(l uint64) test.Event             { return test.ACL(l) }
func ACM(l uint64) test.Event             { return test.ACM(l) }
func ADB(v []bool) test.Event             { return test.ADB(v) }
func ADF16(v []float32) test.Event        { return test.ADF16(v) }
func ADF32(v []float32) test.Event        { return test.ADF32(v) }
func ADF64(v []float64) test.Event        { return test.ADF64(v) }
func ADI16(v []int16) test.Event          { return test.ADI16(v) }
func ADI32(v []int32) test.Event          { return test.ADI32(v) }
func ADI64(v []int64) test.Event          { return test.ADI64(v) }
func ADI8(v []int8) test.Event            { return test.ADI8(v) }
func ADT(v string) test.Event             { return test.ADT(v) }
func ADU(v [][]byte) test.Event           { return test.ADU(v) }
func ADU16(v []uint16) test.Event         { return test.ADU16(v) }
func ADU32(v []uint32) test.Event         { return test.ADU32(v) }
func ADU64(v []uint64) test.Event         { return test.ADU64(v) }
func ADU8(v []uint8) test.Event           { return test.ADU8(v) }
func AF16(v []float32) test.Event         { return test.AF16(v) }
func AF32(v []float32) test.Event         { return test.AF32(v) }
func AF64(v []float64) test.Event         { return test.AF64(v) }
func AI16(v []int16) test.Event           { return test.AI16(v) }
func AI32(v []int32) test.Event           { return test.AI32(v) }
func AI64(v []int64) test.Event           { return test.AI64(v) }
func AI8(v []int8) test.Event             { return test.AI8(v) }
func AU(v [][]byte) test.Event            { return test.AU(v) }
func AU16(v []uint16) test.Event          { return test.AU16(v) }
func AU32(v []uint32) test.Event          { return test.AU32(v) }
func AU64(v []uint64) test.Event          { return test.AU64(v) }
func AU8(v []byte) test.Event             { return test.AU8(v) }
func B(v bool) test.Event                 { return test.B(v) }
func BAB() test.Event                     { return test.BAB() }
func BAF16() test.Event                   { return test.BAF16() }
func BAF32() test.Event                   { return test.BAF32() }
func BAF64() test.Event                   { return test.BAF64() }
func BAI16() test.Event                   { return test.BAI16() }
func BAI32() test.Event                   { return test.BAI32() }
func BAI64() test.Event                   { return test.BAI64() }
func BAI8() test.Event                    { return test.BAI8() }
func BAU() test.Event                     { return test.BAU() }
func BAU16() test.Event                   { return test.BAU16() }
func BAU32() test.Event                   { return test.BAU32() }
func BAU64() test.Event                   { return test.BAU64() }
func BAU8() test.Event                    { return test.BAU8() }
func BCB(v uint64) test.Event             { return test.BCB(v) }
func BCT(v uint64) test.Event             { return test.BCT(v) }
func BMEDIA(v string) test.Event          { return test.BMEDIA(v) }
func BREFR() test.Event                   { return test.BREFR() }
func BRID() test.Event                    { return test.BRID() }
func BS() test.Event                      { return test.BS() }
func CB(t uint64, v []byte) test.Event    { return test.CB(t, v) }
func CM(v string) test.Event              { return test.CM(v) }
func CS(v string) test.Event              { return test.CS(v) }
func CT(t uint64, v string) test.Event    { return test.CT(t, v) }
func E() test.Event                       { return test.E() }
func EDGE() test.Event                    { return test.EDGE() }
func L() test.Event                       { return test.L() }
func M() test.Event                       { return test.M() }
func MARK(id string) test.Event           { return test.MARK(id) }
func MEDIA(t string, v []byte) test.Event { return test.MEDIA(t, v) }
func N(v interface{}) test.Event          { return test.N(v) }
func NAN() test.Event                     { return test.NAN() }
func NODE() test.Event                    { return test.NODE() }
func NULL() test.Event                    { return test.NULL() }
func PAD() test.Event                     { return test.PAD() }
func REFL(id string) test.Event           { return test.REFL(id) }
func REFR(v string) test.Event            { return test.REFR(v) }
func RID(v string) test.Event             { return test.RID(v) }
func S(v string) test.Event               { return test.S(v) }
func SI(id string) test.Event             { return test.SI(id) }
func SNAN() test.Event                    { return test.SNAN() }
func ST(id string) test.Event             { return test.ST(id) }
func T(v compact_time.Time) test.Event    { return test.T(v) }
func UID(v []byte) test.Event             { return test.UID(v) }
func V(v uint64) test.Event               { return test.V(v) }
