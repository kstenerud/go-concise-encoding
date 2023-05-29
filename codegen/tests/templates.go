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
	"math"
	"math/big"

	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/test_runner"
)

func generateArrayInt32Tests() []*test_runner.UnitTest {
	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []int32
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(DirectionsAll, &config, AI32(nil)),
		newMustSucceedTest(DirectionsAll, &config, BAI32(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, int32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAI32(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, int32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAI32(), ACL(uint64(i)), ADI32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := math.MaxInt32 / 31
	for i := math.MinInt32; i < math.MaxInt32-31; i += multiple {
		contents = append(contents, int32(i))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAI32(), ACL(uint64(len(contents))), ADI32(contents)))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32 8 -17 0b1011 0B1010 -0b1100 -0B1101 0o31 0O33 -0o10 -0O5 0x10 0X15 -0x22 -0X33|",
		AI32([]int32{8, -17, 0b1011, 0b1010, -0b1100, -0b1101, 0o31, 0o33, -0o10, -0o5, 0x10, 0x15, -0x22, -0x33})))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 2, 8, 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = math.MaxUint32 / 7
	for i := math.MinInt32; i < math.MaxInt32-7; i += multiple {
		contents = append(contents, int32(i))
	}
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatBinary
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents)))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents[:0])))
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatOctal
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents)))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents[:0])))
	config.DefaultNumericFormats.Array.Int32 = configuration.CTEEncodingFormatHexadecimal
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents)))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AI32(contents[:0])))
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
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAI32(), ACM(4), ADI32(contents[:4]), ACL(3), ADI32(contents[:3])))
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
		newMustSucceedTest(DirectionsAll, &config, BAI32(), ACL(uint64(len(contents))), ADI32(contents))))

	// Whitespace at end of array
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32 |", AI32([]int32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32 1 |", AI32([]int32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32b |", AI32([]int32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32b 1 |", AI32([]int32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32o |", AI32([]int32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32o 1 |", AI32([]int32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32x |", AI32([]int32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|i32x 1 |", AI32([]int32{1})))
	unitTests = append(unitTests, newMustSucceedUnitTest("Whitespace at end of array", mustSucceed...))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Space before array type
	mustFail = nil
	mustFail = append(mustFail, newCTEMustFailTest("| i32|"))
	mustFail = append(mustFail, newCTEMustFailTest("| i32b|"))
	mustFail = append(mustFail, newCTEMustFailTest("| i32o|"))
	mustFail = append(mustFail, newCTEMustFailTest("| i32x|"))
	unitTests = append(unitTests, newMustFailUnitTest("Space before array type", mustFail...))

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
	for iB, base := range intBases {
		for iV, v := range baseOutOfRange {
			if iV <= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range intBases {
		for _, special := range nonIntSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	// Float value in int array
	mustFail = nil
	for _, base := range intBases {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|i32%v 1.2|", base)}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Float value in int array", mustFail...))

	return unitTests
}

func generateArrayUint32Tests() []*test_runner.UnitTest {
	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []uint32
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(DirectionsAll, &config, AU32(nil)),
		newMustSucceedTest(DirectionsAll, &config, BAU32(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, uint32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AU32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAU32(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, uint32(i-8))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAU32(), ACL(uint64(i)), ADU32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := uint64(math.MaxUint32 / 31)
	for i := uint64(0); i < math.MaxUint32-31; i += multiple {
		contents = append(contents, uint32(i))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAU32(), ACL(uint64(len(contents))), ADU32(contents)))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32 8 0b1011 0B1010 0o31 0O33 0x10 0X15|",
		AU32([]uint32{8, 0b1011, 0b1010, 0o31, 0o33, 0x10, 0x15})))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 2, 8, 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = uint64(math.MaxUint32 / 7)
	for i := uint64(0); i < math.MaxUint32-7; i += multiple {
		contents = append(contents, uint32(i))
	}
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatBinary
	t := newMustSucceedTest(DirectionsAll, &config, AU32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatOctal
	t = newMustSucceedTest(DirectionsAll, &config, AU32(contents))
	mustSucceed = append(mustSucceed, t)
	config.DefaultNumericFormats.Array.Uint32 = configuration.CTEEncodingFormatHexadecimal
	t = newMustSucceedTest(DirectionsAll, &config, AU32(contents))
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
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAU32(), ACM(4), ADU32(contents[:4]), ACL(3), ADU32(contents[:3])))
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
		newMustSucceedTest(DirectionsAll, &config, BAU32(), ACL(uint64(len(contents))), ADU32(contents))))

	// Whitespace at end of array
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32 |", AU32([]uint32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32 1 |", AU32([]uint32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32b |", AU32([]uint32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32b 1 |", AU32([]uint32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32o |", AU32([]uint32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32o 1 |", AU32([]uint32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32x |", AU32([]uint32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|u32x 1 |", AU32([]uint32{1})))
	unitTests = append(unitTests, newMustSucceedUnitTest("Whitespace at end of array", mustSucceed...))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Space before array type
	mustFail = nil
	mustFail = append(mustFail, newCTEMustFailTest("| u32|"))
	mustFail = append(mustFail, newCTEMustFailTest("| u32b|"))
	mustFail = append(mustFail, newCTEMustFailTest("| u32o|"))
	mustFail = append(mustFail, newCTEMustFailTest("| u32x|"))
	unitTests = append(unitTests, newMustFailUnitTest("Space before array type", mustFail...))

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
	for iB, base := range intBases {
		for iV, v := range baseOutOfRange {
			if iV <= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range intBases {
		for _, special := range nonIntSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	// Float value in int array
	mustFail = nil
	for _, base := range intBases {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|u32%v 1.2|", base)}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Float value in int array", mustFail...))

	return unitTests
}

func generateArrayFloat32Tests() []*test_runner.UnitTest {
	// Trick golang into producing negative zero
	float32ZeroValues[1] = -float32ZeroValues[0]

	var unitTests []*test_runner.UnitTest
	var mustSucceed []*test_runner.MustSucceedTest
	var contents []float32
	config := configuration.DefaultCTEEncoderConfiguration()

	// Empty array
	unitTests = append(unitTests, newMustSucceedUnitTest("Empty Array",
		newMustSucceedTest(DirectionsAll, &config, AF32(nil)),
		newMustSucceedTest(DirectionsAll, &config, BAF32(), ACL(0)),
	))

	// Short array
	contents = contents[:0]
	mustSucceed = nil
	for i := 1; i <= 15; i++ {
		contents = append(contents, floatCap32(float32(i-8)))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, AF32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Short Array", mustSucceed...))

	// Chunked array
	contents = contents[:0]
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAF32(), ACL(0)))
	for i := 1; i <= 20; i++ {
		contents = append(contents, floatCap32(float32(1)/float32(i-10)))
		mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAF32(), ACL(uint64(i)), ADF32(contents)))
	}
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunked Array", mustSucceed...))

	// Various element values
	contents = contents[:0]
	mustSucceed = nil
	multiple := math.MaxInt32 / 31
	for i := math.MinInt32; i < math.MaxInt32-31; i += multiple {
		contents = append(contents, floatCap32(float32(i)))
	}
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAF32(), ACL(uint64(len(contents))), ADF32(contents)))
	unitTests = append(unitTests, newMustSucceedUnitTest("Various Array Elements", mustSucceed...))

	// Base 16
	mustSucceed = nil
	contents = contents[:0]
	multiple = math.MaxUint32 / 7
	for i := math.MinInt32; i < math.MaxInt32-7; i += multiple {
		contents = append(contents, floatCap32(float32(i)))
	}
	config.DefaultNumericFormats.Array.Float32 = configuration.CTEEncodingFormatHexadecimal
	t := newMustSucceedTest(DirectionsAll, &config, AF32(contents))
	mustSucceed = append(mustSucceed, t)
	config = configuration.DefaultCTEEncoderConfiguration()
	unitTests = append(unitTests, newMustSucceedUnitTest("Base 16", mustSucceed...))

	// Chunking
	mustSucceed = nil
	events := []test.Event{BAF32()}
	for i := 0; i < 7; i++ {
		events = append(events, ACM(uint64(i)))
		if i > 0 {
			events = append(events, ADF32(contents[:i]))
		}
	}
	events = append(events, ACL(0))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, events...))
	mustSucceed = append(mustSucceed, newMustSucceedTest(DirectionsAll, &config, BAF32(), ACM(4), ADF32(contents[:4]), ACL(3), ADF32(contents[:3])))
	unitTests = append(unitTests, newMustSucceedUnitTest("Chunking Variations", mustSucceed...))

	// Edge-case element values
	contents = contents[:0]
	contents = append(contents, float32ZeroValues...)
	unitTests = append(unitTests, newMustSucceedUnitTest("Edge Case Element Values",
		newMustSucceedTest(DirectionsAll, &config, BAF32(), ACL(uint64(len(contents))), ADF32(contents))))

	// Edge-case element values
	contents = contents[:0]
	contents = append(contents, float32NanValues...)
	unitTests = append(unitTests, newMustSucceedUnitTest("NaN Element Values",
		newMustSucceedTest(DirectionsAll.except(DirectionToCBE), &config, BAF32(), ACL(uint64(len(contents))), ADF32(contents))))

	// Whitespace at end of array
	mustSucceed = nil
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|f32 |", AF32([]float32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|f32 1 |", AF32([]float32{1})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|f32x |", AF32([]float32{})))
	mustSucceed = append(mustSucceed, newCTEMustSucceedTest("|f32x 1 |", AF32([]float32{1})))
	unitTests = append(unitTests, newMustSucceedUnitTest("Whitespace at end of array", mustSucceed...))

	// Fail mode tests
	var mustFail []*test_runner.MustFailTest

	// Space before array type
	mustFail = nil
	mustFail = append(mustFail, newCTEMustFailTest("| f32|"))
	mustFail = append(mustFail, newCTEMustFailTest("| f32x|"))
	unitTests = append(unitTests, newMustFailUnitTest("Space before array type", mustFail...))

	// Truncated Array
	mustFail = nil
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAF32()))
	mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAF32(), ACL(1)))
	for i := 2; i <= 10; i++ {
		contents = contents[:i/2]
		mustFail = append(mustFail, newMustFailTest(testTypeCbe, BAF32(), ACL(uint64(i)), ADF32(contents)))
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: truncate(generateCBE(AF32(contents)), 1)}})
	}
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAF32())}})
	mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, BAF32(), ACL(uint64(2)), ADF32(contents[:1]))}})
	unitTests = append(unitTests, newMustFailUnitTest("Truncated Array", mustFail...))

	// Element value out of range
	mustFail = nil
	for _, v := range f32OutOfRange {
		mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: v}})
	}
	unitTests = append(unitTests, newMustFailUnitTest("Element value out of range", mustFail...))

	// Numeric digit out of range
	mustFail = nil
	for iB, base := range floatBases {
		for iV, v := range baseOutOfRange {
			if iV <= iB {
				mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|f32%v %v|", base, v)}})
			}
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Numeric digit out of range", mustFail...))

	// Invalid special values
	mustFail = nil
	for _, base := range intBases {
		for _, special := range nonFloatSpecials {
			mustFail = append(mustFail, &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: fmt.Sprintf("|f32%v %v|", base, special)}})
		}
	}
	unitTests = append(unitTests, newMustFailUnitTest("Invalid special values", mustFail...))

	return unitTests
}

func floatCap16(v float32) float32 {
	bits := math.Float32bits(v)
	bits &= 0xffff0000
	return math.Float32frombits(bits)
}

func floatCap32(v float32) float32 {
	return v
}

func floatCap64(v float64) float64 {
	return v
}

var f16OutOfRange = []string{
	"|f16x 1.23456p128|",
	"|f16 1.234567e40|",
	"|f16 0x1.23456p128|",
	"|f16x 1.23456p-151|",
	"|f16 1.234567e-50|",
	"|f16 -0x1.23456p-151|",
}

var f32OutOfRange = []string{
	"|f32x 1.23456p128|",
	"|f32 1.234567e40|",
	"|f32 0x1.23456p128|",
	"|f32x 1.23456p-151|",
	"|f32 1.234567e-50|",
	"|f32 -0x1.23456p-151|",
}

var f64OutOfRange = []string{
	"|f64x 1.23456p128000|",
	"|f64 1.234567e40000|",
	"|f64 0x1.23456p128000|",
	"|f64x 1.23456p-151000|",
	"|f64 1.234567e-50000|",
	"|f64 -0x1.23456p-151000|",
}

var uint8OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint8), big.NewInt(1))
var uint16OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint16), big.NewInt(1))
var uint32OutOfRange = big.NewInt(0).Add(big.NewInt(math.MaxUint32), big.NewInt(1))
var uint64OutOfRange = big.NewInt(0).Add(big.NewInt(0).Lsh(big.NewInt(math.MaxInt64), 1), big.NewInt(2))

var intBases = []string{
	"x",
	"",
	"o",
	"b",
}

var floatBases = []string{
	"x",
	"",
}

var baseOutOfRange = []string{
	"1g",
	"1a",
	"18",
	"12",
}

var nonIntSpecials = []string{
	"null",
	"-null",
	"true",
	"-true",
	"false",
	"-false",
	"nan",
	"-nan",
	"snan",
	"-snan",
	"inf",
	"-inf",
}

var nonFloatSpecials = []string{
	"null",
	"-null",
	"true",
	"-true",
	"false",
	"-false",
	"-nan",
	"-snan",
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

const (
	float32QuietNanBits     = uint32(0x7fe00000)
	float32SignalingNanBits = uint32(0x7fa00000)
	float64QuietNanBits     = uint64(0x7ffc000000000000)
	float64SignalingNanBits = uint64(0x7ff4000000000000)
)

var (
	float32SignalingNan = math.Float32frombits(float32SignalingNanBits)
	float32QuietNan     = math.Float32frombits(float32QuietNanBits)
	float64SignalingNan = math.Float64frombits(float64SignalingNanBits)
	float64QuietNan     = math.Float64frombits(float64QuietNanBits)
)

var float32ZeroValues = []float32{
	0,
	0,
}

var float64ZeroValues = []float64{
	0,
	0,
}

var float32NanValues = []float32{
	float32SignalingNan,
	float32QuietNan,
}

var float64NanValues = []float64{
	float64SignalingNan,
	float64QuietNan,
}
