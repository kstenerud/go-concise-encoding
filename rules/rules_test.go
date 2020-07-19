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

package rules

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

const rulesCodecVersion = 1

func assertRulesOnString(t *testing.T, rules *Rules, value string) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), true) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(value)) })
	}
}

func assertRulesAddBytes(t *testing.T, rules *Rules, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), true) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData(value) })
	}
}

func assertRulesAddURI(t *testing.T, rules *Rules, uri string) {
	length := len(uri)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), true) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(uri)) })
	}
}

func assertRulesAddCustom(t *testing.T, rules *Rules, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnCustomBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), true) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData(value) })
	}
}

func assertRulesInArray(t *testing.T, rules *Rules) {
	if rules.arrayType == eventTypeNothing {
		t.Errorf("Expected to be in array")
	}
}

func assertRulesNotInArray(t *testing.T, rules *Rules) {
	if rules.arrayType != eventTypeNothing {
		t.Errorf("Expected to not be in array")
	}
}

func newRulesAfterVersion(options *options.RuleOptions) *Rules {
	rules := NewRules(events.NewNullEventReceiver(), options)
	rules.OnVersion(options.ConciseEncodingVersion)
	return rules
}

func newRulesWithMaxDepth(maxDepth int) *Rules {
	options := options.DefaultRuleOptions()
	options.MaxContainerDepth = uint64(maxDepth)
	return newRulesAfterVersion(options)
}

// ============================================================================

// Tests

// ===========
// Basic Types
// ===========

func TestRulesVersion(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.ConciseEncodingVersion = 1
	rules := NewRules(events.NewNullEventReceiver(), options)
	test.AssertPanics(t, func() { rules.OnVersion(2) })
	test.AssertNoPanic(t, func() { rules.OnVersion(1) })
	test.AssertPanics(t, func() { rules.OnVersion(1) })
}

func TestRulesNil(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesNan(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNan(true) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesNan2(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnFloat(math.NaN()) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBool(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnTrue() })
	test.AssertNoPanic(t, func() { rules.OnFalse() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesInt(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnInt(-1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesPositiveInt(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesNegativeInt(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBigInt(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBigInt(big.NewInt(-1)) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	bi := big.NewInt(0x7fffffffffffffff)
	bi = bi.Mul(bi, big.NewInt(10000000))
	test.AssertNoPanic(t, func() { rules.OnBigInt(bi) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	bi = big.NewInt(-0x7fffffffffffffff)
	bi = bi.Mul(bi, big.NewInt(10000000))
	test.AssertNoPanic(t, func() { rules.OnBigInt(bi) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBigIntAsID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnBigInt(big.NewInt(10)) })
	test.AssertNoPanic(t, func() { rules.OnInt(0) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnInt(10) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	bi := big.NewInt(0x7fffffffffffffff)
	bi = bi.Mul(bi, big.NewInt(10000000))
	test.AssertPanics(t, func() { rules.OnBigInt(bi) })
}

func TestRulesFloat(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnFloat(math.NaN()) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBigFloat(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBigFloat(big.NewFloat(1.1)) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesDecimalFloat(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnDecimalFloat(compact_float.DFloatValue(2, 1000)) })
	test.AssertNoPanic(t, func() { rules.OnDecimalFloat(compact_float.QuietNaN()) })
	test.AssertNoPanic(t, func() { rules.OnDecimalFloat(compact_float.SignalingNaN()) })
	test.AssertNoPanic(t, func() { rules.OnDecimalFloat(compact_float.Infinity()) })
	test.AssertNoPanic(t, func() { rules.OnDecimalFloat(compact_float.NegativeInfinity()) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBigDecimalFloat(t *testing.T) {
	rules := newRulesWithMaxDepth(10)

	assertAPD := func(str string) {
		v, _, err := apd.NewFromString(str)
		if err != nil {
			panic(err)
		}
		test.AssertNoPanic(t, func() { rules.OnBigDecimalFloat(v) })
	}

	test.AssertNoPanic(t, func() { rules.OnList() })
	assertAPD("1.5")
	assertAPD("-10.544e10000")
	assertAPD("nan")
	assertAPD("snan")
	assertAPD("infinity")
	assertAPD("-infinity")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesUUID(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnUUID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesTime(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCompactTime(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnCompactTime(compact_time.AsCompactTime(time.Now())) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBytesOneshot(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytes([]byte{1, 2, 3, 4}) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesURIOneshot(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURI("http://example.com") })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCustomOneshot(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnCustom([]byte{1, 2, 3, 4}) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMarker(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(100) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
}

func TestRulesReference(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })

	test.AssertNoPanic(t, func() { rules.OnMarker() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })

	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(100) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })

	test.AssertNoPanic(t, func() { rules.OnReference() })
	assertRulesOnString(t, rules, "a")

	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnURI("http://example.com") })

	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(100) })

	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnPositiveInt(5) })

	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte("test")) })
}

// ===========
// Array Types
// ===========

func testRulesBytes(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayChunk(0, true) })
	} else {
		for i, count := range byteCount {
			test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(count), i == lastIndex) })
			if count > 0 {
				test.AssertNoPanic(t, func() { rules.OnArrayData(make([]byte, count)) })
			}
		}
	}
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesBytes(t *testing.T) {
	testRulesBytes(t, 0)
	testRulesBytes(t, 1, 1)
	testRulesBytes(t, 2, 2)
	testRulesBytes(t, 10, 10)
	testRulesBytes(t, 100, 14, 55, 20, 11)
}

func TestRulesBytesMultiChunk(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(10, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(5, 0)) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(3, 0)) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(2, 0)) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func testRulesString(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayChunk(0, true) })
	} else {
		for i, count := range byteCount {
			test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(count), i == lastIndex) })
			if count > 0 {
				test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(count, 0))) })
			}
		}
	}
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesString(t *testing.T) {
	testRulesString(t, 0)
	testRulesString(t, 1, 1)
	testRulesString(t, 2, 2)
	testRulesString(t, 10, 10)
	testRulesString(t, 111, 10, 50, 41, 10)
}

func testRulesSingleString(t *testing.T, value string) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnString(value) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesStringNonComment(t *testing.T) {
	testRulesSingleString(t, "\f")
}

func TestRulesStringMultiChunk(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(20, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(5, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(3, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(12, 0))) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesURI(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(18, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("http://example.com")) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesURIMultiChunk(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(13, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("http:")) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("test")) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(".net")) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func testRulesCustom(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnCustomBegin() })
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayChunk(0, true) })
	} else {
		for i, count := range byteCount {
			test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(count), i == lastIndex) })
			if count > 0 {
				test.AssertNoPanic(t, func() { rules.OnArrayData(make([]byte, count)) })
			}
		}
	}
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCustom(t *testing.T) {
	testRulesCustom(t, 0)
	testRulesCustom(t, 1, 1)
	testRulesCustom(t, 2, 2)
	testRulesCustom(t, 10, 10)
	testRulesCustom(t, 100, 14, 55, 20, 11)
}

func TestRulesCustomMultiChunk(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnCustomBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(10, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(5, 0)) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(3, 0)) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(2, 0)) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesInArrayBasic(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
	assertRulesNotInArray(t, rules)
}

func TestRulesInArray(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnList() })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(10, false) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(2, 0))) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(5, 0))) })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("a:123")) })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
	assertRulesNotInArray(t, rules)
}

func TestRulesInArrayEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnList() })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(0, true) })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	assertRulesInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(0, true) })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	assertRulesNotInArray(t, rules)
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
	assertRulesNotInArray(t, rules)
}

// =================
// Containers: Empty
// =================

func TestRulesListEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMapEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMarkupEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMetadataEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCommentEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

// =======================
// Containers: Single item
// =======================

func TestRulesListSingleItem(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMapPair(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMarkupSingleItem(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "abcdef")
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMetadataSingleItem(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCommentSingleItem(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

// ==================
// Containers: Filled
// ==================

func TestRulesListFilled(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	assertRulesAddBytes(t, rules, test.GenerateBytes(1, 0))
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMapFilled(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnNan(true) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	assertRulesAddBytes(t, rules, test.GenerateBytes(1, 0))
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMetadataFilled(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	assertRulesAddBytes(t, rules, test.GenerateBytes(1, 0))
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMapList(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMapMap(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesDeepContainer(t *testing.T) {
	rules := newRulesWithMaxDepth(6)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	assertRulesOnString(t, rules, "0123456789")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCommentInt(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	assertRulesOnString(t, rules, "blah\r\n\t\tblah")
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(1, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0x00}) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0x0b}) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0x7f}) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0x80}) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte{0x40}) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesCommentMap(t *testing.T) {
	rules := newRulesWithMaxDepth(3)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMetadataCommentMap(t *testing.T) {
	rules := newRulesWithMaxDepth(4)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMarkup(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

// =============
// Allowed Types
// =============

func TestRulesAllowedTypesTLO(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnArrayChunk(1, true) })
	rules = newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesList(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	assertRulesAddBytes(t, rules, []byte{})
	assertRulesOnString(t, rules, "")
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(2, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("a:")) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
}

func TestRulesAllowedTypesMapKey(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	assertRulesAddBytes(t, rules, []byte{})
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	assertRulesOnString(t, rules, "")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	assertRulesAddURI(t, rules, "a:")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMapValue(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesAddBytes(t, rules, []byte{})
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesOnString(t, rules, "")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesAddURI(t, rules, "a:")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMetadataKey(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	assertRulesAddBytes(t, rules, []byte{})
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	assertRulesOnString(t, rules, "")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	assertRulesAddURI(t, rules, "a:")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMetadataValue(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesAddBytes(t, rules, []byte{})
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesOnString(t, rules, "")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	assertRulesAddURI(t, rules, "a:")
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesComment(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	assertRulesOnString(t, rules, "")
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	test.AssertPanics(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
}

func TestRulesAllowedTypesMarkupName(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMarkupAttributeKey(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMarkupAttributeValue(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMarkupContents(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesArray(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(false) })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnComment() })
	test.AssertPanics(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	test.AssertPanics(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	test.AssertPanics(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnPadding(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	test.AssertNoPanic(t, func() { rules.OnArrayChunk(1, true) })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnComment() })
	test.AssertPanics(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	test.AssertPanics(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	test.AssertPanics(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnPadding(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(1, true) })

	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte{0}) })
}

func TestRulesAllowedTypesMarkerID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesMarkerValue(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNan(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesAllowedTypesReferenceID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnNil() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnNan(false) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnBool(true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnFloat(0.1) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnNegativeInt(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnTime(time.Now()) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnList() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnMap() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnMetadata() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnComment() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnEnd() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnBytesBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0}) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnMarker() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnReference() })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesEmptyDocument(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesOnlyComment(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesOnlyPadding(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnPadding(1) })
	test.AssertNoPanic(t, func() { rules.OnEndDocument() })
}

func TestRulesMarkerCommentObject(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
}

func TestRulesMarkerMetadataObject(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
}

func TestRulesMarkerReference(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(2) })
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
}

// ================
// Error conditions
// ================

func TestRulesErrorOnlyMetadata(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesErrorListOnlyMetadata(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesErrorNoContainerToEnd(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesErrorNoArrayToBeginChunk(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
}

func TestRulesErrorOnEndTooManyTimes(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesErrorUnendedContainer(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesErrorUnendedArray(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesErrorArrayTooManyBytes(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(6, 0)) })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(6, 0)) })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(6, 0)) })
}

func TestRulesErrorArrayTooFewBytes(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })

	rules = newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(5, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnEndDocument() })
}

func TestRulesErrorAddDataNotInArray(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
}

func TestRulesErrorInvalidMapKey(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(false) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(1, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(1, 0))) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
	test.AssertPanics(t, func() { rules.OnNil() })
	test.AssertPanics(t, func() { rules.OnNan(true) })
	test.AssertPanics(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnMap() })
	test.AssertPanics(t, func() { rules.OnMarkup() })
}

func TestRulesErrorMapMissingValue(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMap() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	assertRulesOnString(t, rules, "a")
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesErrorMarkupNameLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
}

func TestRulesErrorMarkerIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
}

func TestRulesErrorMarkerFollowedByOnEnd(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(1) })
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesErrorReferenceIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })
}

func TestRulesErrorURILength0_1(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertPanics(t, func() { rules.OnArrayChunk(0, true) })

	rules = newRulesWithMaxDepth(2)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(1, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte{0x40}) })
}

func TestRulesErrorDuplicateMarkerID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("test")) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("test")) })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })

	rules = newRulesWithMaxDepth(10)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(100) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(100) })
	test.AssertPanics(t, func() { rules.OnFloat(0.1) })
}

// ======
// Limits
// ======

func TestRulesMaxBytesLength(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxBytesLength = 10
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(11, 0)) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnBytesBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData(test.GenerateBytes(8, 0)) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData(test.GenerateBytes(4, 0)) })
}

func TestRulesMaxStringLength(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxStringLength = 10
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(8, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
}

func TestRulesMaxURILength(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxURILength = 10
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte("someuri:aaa")) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte("someuri:")) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
}

func TestRulesMaxIDLength(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxIDLength = 10
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(8, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnReference() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(8, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
}

func TestRulesMaxMarkupNameLength(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxMarkupNameLength = 10
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnMarkup() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(8, false) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(8, 0))) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(4, true) })
	test.AssertPanics(t, func() { rules.OnArrayData([]byte(test.GenerateString(4, 0))) })
}

func TestRulesMaxContainerDepth(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertPanics(t, func() { rules.OnList() })
}

func TestRulesMaxObjectCount(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxObjectCount = 3
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnEnd() })

	rules = newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnComment() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertNoPanic(t, func() { rules.OnMetadata() })
	test.AssertNoPanic(t, func() { rules.OnEnd() })
	test.AssertPanics(t, func() { rules.OnEnd() })
}

func TestRulesMaxReferenceCount(t *testing.T) {
	options := options.DefaultRuleOptions()
	options.MaxReferenceCount = 2
	rules := newRulesAfterVersion(options)
	test.AssertNoPanic(t, func() { rules.OnList() })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(11, true) })
	test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(test.GenerateString(11, 0))) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertNoPanic(t, func() { rules.OnMarker() })
	test.AssertNoPanic(t, func() { rules.OnPositiveInt(10) })
	test.AssertNoPanic(t, func() { rules.OnBool(true) })
	test.AssertPanics(t, func() { rules.OnMarker() })
}
