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
	"io"
	"math"
	"math/big"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
)

var testsImports = []*common.Import{
	{As: "", Import: "fmt"},
	{As: "", Import: "math"},
	{As: "", Import: "math/big"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/configuration"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/test"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/test/test_runner"},
}

func generateTestGenerators(basePath string) {
	common.GenerateGoFile(basePath, "tests", testsImports, func(writer io.Writer) {
		generateArrayTestGenerator(writer)
	})
}

func generateArrayTestGenerator(writer io.Writer) {
	writer.Write([]byte(strings.ReplaceAll(intArrayTestTemplate, "32", "8")))
	writer.Write([]byte(strings.ReplaceAll(intArrayTestTemplate, "32", "16")))
	writer.Write([]byte(intArrayTestTemplate))
	writer.Write([]byte(strings.ReplaceAll(intArrayTestTemplate, "32", "64")))

	writer.Write([]byte(strings.ReplaceAll(uintArrayTestTemplate, "32", "8")))
	writer.Write([]byte(strings.ReplaceAll(uintArrayTestTemplate, "32", "16")))
	writer.Write([]byte(uintArrayTestTemplate))
	writer.Write([]byte(strings.ReplaceAll(uintArrayTestTemplate, "32", "64")))
}

var intArrayTestTemplate = `func generateArrayInt32Tests() []*test_runner.UnitTest {
	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []int32
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(&config, AI32(nil)),
		newMustSucceedTest(&config, BAI32(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, int32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, AI32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI32(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, int32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI32(), ACL(uint64(i)), ADI32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := math.MaxInt32 / 31
	for i := math.MinInt32; i < math.MaxInt32-31; i += multiple {
		contents = append(contents, int32(math.MinInt32+i))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI32(), ACL(uint64(len(contents))), ADI32(contents)))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 2, 8, 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = math.MaxUint32 / 7
	for i := math.MinInt32; i < math.MaxInt32-7; i += multiple {
		contents = append(contents, int32(math.MinInt32+i))
	}
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatBinary
	t := newMustSucceedTest(&config, AI32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatOctal
	t = newMustSucceedTest(&config, AI32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatHexadecimal
	t = newMustSucceedTest(&config, AI32(contents))
	mustSucceed = append(mustSucceed, t)
	config = configuration.DefaultCTEEncoderConfiguration()
	unitTests = append(unitTests, newMustSucceedUnitTest("Base 2, 8, 16", mustSucceed...))

	// Chunking
	mustSucceed = nil
	events := []test.Event{BAI32()}
	for i := 0; i < 7; i++ {
		events = append(events, ACM(uint64(i)))
		if i > 0 {
			events = append(events, ADI32(contents[:i]))
		}
	}
	events = append(events, ACL(0))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAI32(), ACM(4), ADI32(contents[:4]), ACL(3), ADI32(contents[:3])))
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunking Variations", mustSucceed...))

	// Edge-case element values
	contents = contents[:0]
	for _, v := range intEdgeValues {
		contents = append(contents, int32(v))
	}
	for _, v := range uintEdgeValues {
		contents = append(contents, int32(v))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Edge Case Element Values",
		newMustSucceedTest(&config, BAI32(), ACL(uint64(len(contents))), ADI32(contents))))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Truncated Array
	mustFail = nil
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI32()))
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI32(), ACL(1)))
	for i := 2; i <= 10; i++ {
		contents = contents[:i/2]
		mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAI32(), ACL(uint64(i)), ADI32(contents)))
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: truncate(generateCBE(AI32(contents)), 1)}})
	}
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAI32())}})
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAI32(), ACL(uint64(2)), ADI32(contents[:1]))}})
	unitTests = append(unitTests, newMustFailUnitTest("Truncated Array", mustFail...))

	// Element value out of range
	bigValue := big.NewInt(math.MaxInt32)
	bigValue.Add(bigValue, big.NewInt(1))
	smallValue := big.NewInt(math.MinInt32)
	smallValue.Sub(smallValue, big.NewInt(1))
	smallValue.Neg(smallValue)
	unitTests = append(unitTests, newMustFailUnitTest(
		"Element value out of range",
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32b %b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32o %o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32x %x|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 %d|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 0b%b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 0o%o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 0x%x|", bigValue)}},

		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32b -%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32o -%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32x -%x|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 -%d|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 -0b%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 -0o%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32 -0x%x|", smallValue)}},
	))

	// Numeric digit out of range
	mustFail = nil
	for iB, base := range bases {
		for iV, v := range outOfRange {
			if iV >= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range bases {
		for _, special := range nonIntSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	// Float value in int array
	mustFail = nil
	for _, base := range bases {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v 1.2|", base)}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Float value in int array", mustFail...))

	return unitTests
}

`

var uintArrayTestTemplate = `func generateArrayUint32Tests() []*test_runner.UnitTest {
	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []uint32
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(&config, AU32(nil)),
		newMustSucceedTest(&config, BAU32(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, uint32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, AU32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAU32(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, uint32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAU32(), ACL(uint64(i)), ADU32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := uint64(math.MaxUint32 / 31)
	for i := uint64(0); i < math.MaxUint32-31; i += multiple {
		contents = append(contents, uint32(i))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAU32(), ACL(uint64(len(contents))), ADU32(contents)))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 2, 8, 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = uint64(math.MaxUint32 / 7)
	for i := uint64(0); i < math.MaxUint32-7; i += multiple {
		contents = append(contents, uint32(i))
	}
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatBinary
	t := newMustSucceedTest(&config, AU32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatOctal
	t = newMustSucceedTest(&config, AU32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatHexadecimal
	t = newMustSucceedTest(&config, AU32(contents))
	mustSucceed = append(mustSucceed, t)
	config = configuration.DefaultCTEEncoderConfiguration()
	unitTests = append(unitTests, newMustSucceedUnitTest("Base 2, 8, 16", mustSucceed...))

	// Chunking
	mustSucceed = nil
	events := []test.Event{BAU32()}
	for i := 0; i < 7; i++ {
		events = append(events, ACM(uint64(i)))
		if i > 0 {
			events = append(events, ADU32(contents[:i]))
		}
	}
	events = append(events, ACL(0))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(&config, BAU32(), ACM(4), ADU32(contents[:4]), ACL(3), ADU32(contents[:3])))
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunking Variations", mustSucceed...))

	// Edge-case element values
	contents = contents[:0]
	for _, v := range intEdgeValues {
		contents = append(contents, uint32(v))
	}
	for _, v := range uintEdgeValues {
		contents = append(contents, uint32(v))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Edge Case Element Values",
		newMustSucceedTest(&config, BAU32(), ACL(uint64(len(contents))), ADU32(contents))))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Truncated Array
	mustFail = nil
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAU32()))
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAU32(), ACL(1)))
	for i := 2; i <= 10; i++ {
		contents = contents[:i/2]
		mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAU32(), ACL(uint64(i)), ADU32(contents)))
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: truncate(generateCBE(AU32(contents)), 1)}})
	}
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAU32())}})
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAU32(), ACL(uint64(2)), ADU32(contents[:1]))}})
	unitTests = append(unitTests, newMustFailUnitTest("Truncated Array", mustFail...))

	// Element value out of range
	bigValue := uint32OutOfRange
	smallValue := big.NewInt(0)
	smallValue.Sub(smallValue, big.NewInt(1))
	smallValue.Neg(smallValue)
	unitTests = append(unitTests, newMustFailUnitTest(
		"Element value out of range",
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32b %b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32o %o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32x %x|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 %d|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 0b%b|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 0o%o|", bigValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 0x%x|", bigValue)}},

		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32b -%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32o -%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32x -%x|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 -%d|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 -0b%b|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 -0o%o|", smallValue)}},
		&test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32 -0x%x|", smallValue)}},
	))

	// Numeric digit out of range
	mustFail = nil
	for iB, base := range bases {
		for iV, v := range outOfRange {
			if iV >= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range bases {
		for _, special := range nonIntSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	// Float value in int array
	mustFail = nil
	for _, base := range bases {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v 1.2|", base)}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Float value in int array", mustFail...))

	return unitTests
}

`

var uint8OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint8), big.NewInt(1))
var uint16OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint16), big.NewInt(1))
var uint32OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint32), big.NewInt(1))
var uint64OutOfRange = big.NewInt(0).Add(big.NewInt(0).Lsh(big.NewInt(math.MaxInt64), 1), big.NewInt(2))

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
