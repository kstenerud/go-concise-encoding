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
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
)

// ===========
// Basic Types
// ===========

func TestRulesBeginDocument(t *testing.T) {
	opts := options.DefaultRuleOptions()
	rules := NewRules(events.NewNullEventReceiver(), opts)
	assertEventsFail(t, rules, EvV)
	assertEventsSucceed(t, rules, BD(), EvV)
}

func TestRulesVersion(t *testing.T) {
	opts := options.DefaultRuleOptions()
	rules := NewRules(events.NewNullEventReceiver(), opts)
	assertEventsSucceed(t, rules, BD())
	assertEventsFail(t, rules, V(9))
	assertEventsSucceed(t, rules, EvV)
	assertEventsFail(t, rules, EvV)
}

func TestRulesNull(t *testing.T) {
	assertEventsMaxDepth(t, 1, NULL(), ED())
}

func TestRulesNan(t *testing.T) {
	assertEventsMaxDepth(t, 1, NAN(), ED())
	assertEventsMaxDepth(t, 1, SNAN(), ED())
}

func TestRulesNan2(t *testing.T) {
	assertEventsMaxDepth(t, 1, BF(math.NaN()), ED())
}

func TestRulesBool(t *testing.T) {
	assertEventsMaxDepth(t, 10, L(), B(true), TT(), FF(), E(), ED())
}

func TestRulesInt(t *testing.T) {
	assertEventsMaxDepth(t, 1, I(-1), ED())
}

func TestRulesPositiveInt(t *testing.T) {
	assertEventsMaxDepth(t, 1, PI(1), ED())
}

func TestRulesNegativeInt(t *testing.T) {
	assertEventsMaxDepth(t, 1, NI(1), ED())
}

func TestRulesBigInt(t *testing.T) {
	assertEventsMaxDepth(t, 1, BI(NewBigInt("1")), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("-1")), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("10000000000000000000000000000000000000000000")), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("-10000000000000000000000000000000000000000000")), ED())
}

func TestRulesFloat(t *testing.T) {
	assertEventsMaxDepth(t, 1, BF(0.1), ED())
	assertEventsMaxDepth(t, 1, BF(math.NaN()), ED())
}

func TestRulesBigFloat(t *testing.T) {
	assertEventsMaxDepth(t, 1, BBF(NewBigFloat("1.1")), ED())
}

func TestRulesDecimalFloat(t *testing.T) {
	assertEventsMaxDepth(t, 10, L(), DF(NewDFloat("2.1e2000")),
		DF(NewDFloat("nan")), DF(NewDFloat("snan")), DF(NewDFloat("inf")),
		DF(NewDFloat("-inf")), E(), ED())
}

func TestRulesBigDecimalFloat(t *testing.T) {
	assertEventsMaxDepth(t, 10, L(), BDF(NewBDF("1.5")), BDF(NewBDF("-10.544e10000")),
		BDF(NewBDF("nan")), BDF(NewBDF("snan")), BDF(NewBDF("infinity")),
		BDF(NewBDF("-infinity")), E(), ED())
}

func TestRulesUID(t *testing.T) {
	assertEventsMaxDepth(t, 1, UID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), ED())
}

func TestRulesTime(t *testing.T) {
	assertEventsMaxDepth(t, 1, GT(time.Now()), ED())
}

func TestRulesCompactTime(t *testing.T) {
	assertEventsMaxDepth(t, 1, T(NewDate(1, 1, 1)), ED())
}

func TestRulesArrayOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, AU8([]byte{1, 2, 3, 4}), ED())
	assertEventsMaxDepth(t, 1, AU16([]uint16{1, 2, 3, 4}), ED())
	// TODO: Other array types
}

func TestRulesResourceIDOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, RID("http://example.com"), ED())
}

func TestRulesCustomOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, CB([]byte{1, 2, 3, 4}), ED())
	assertEventsMaxDepth(t, 1, CT("test"), ED())
}

func TestRulesMarker(t *testing.T) {
	assertEventsMaxDepth(t, 10, L(), MARK("a"), BF(0.1), MARK("blah"), GT(time.Now()), E(), ED())
}

func TestRulesReference(t *testing.T) {
	rules := newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(), MARK("a"), BF(0.1), REF("a"), E())

	rules = newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(), RREF("http://example.com"), E())

	rules = newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(), MARK("a"), BF(0.1), MARK("blah"), GT(time.Now()),
		REF("a"),
		RREF("http://example.com"),
		REF("blah"),
		S("test"),
		E())

	rules = newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(),
		MARK("a"), BF(0.1),
		MARK("blah"), GT(time.Now()),
		REF("a"), RRB(), AC(18, false), AD([]byte("http://example.com")),
		REF("blah"),
		S("test"),
		E())
}

// ===========
// Array Types
// ===========

func TestRulesBytes(t *testing.T) {
	testRulesBytes := func(t *testing.T, length int, byteCount ...int) {
		rules := newRulesWithMaxDepth(1)
		assertEventsSucceed(t, rules, AU8B())

		lastIndex := len(byteCount) - 1
		if lastIndex < 0 {
			assertEventsSucceed(t, rules, AC(0, false))
		} else {
			for i, count := range byteCount {
				assertEventsSucceed(t, rules, AC(uint64(count), i != lastIndex))
				if count > 0 {
					assertEventsSucceed(t, rules, AD(make([]byte, count)))
				}
			}
		}
		assertEventsSucceed(t, rules, ED())
	}

	testRulesBytes(t, 0)
	testRulesBytes(t, 1, 1)
	testRulesBytes(t, 2, 2)
	testRulesBytes(t, 10, 10)
	testRulesBytes(t, 100, 14, 55, 20, 11)
}

func TestRulesBytesMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, AU8B(), AC(10, false), AD(NewBytes(5, 0)), AD(NewBytes(3, 0)), AD(NewBytes(2, 0)), ED())
}

func testRulesString(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, SB())

	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, AC(0, false))
	} else {
		for i, count := range byteCount {
			assertEventsSucceed(t, rules, AC(uint64(count), i != lastIndex))
			if count > 0 {
				assertEventsSucceed(t, rules, AD(NewString(count, 0)))
			}
		}
	}
	assertEventsSucceed(t, rules, ED())
}

func TestRulesString(t *testing.T) {
	testRulesString(t, 0)
	testRulesString(t, 1, 1)
	testRulesString(t, 2, 2)
	testRulesString(t, 10, 10)
	testRulesString(t, 111, 10, 50, 41, 10)

	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, S("新しい"))
}

func testRulesSingleString(t *testing.T, value string) {
	assertEventsMaxDepth(t, 1, S(value))
}

func TestRulesStringNonComment(t *testing.T) {
	testRulesSingleString(t, "\f")
}

func TestRulesStringMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, SB(), AC(20, false), AD(NewString(5, 0)), AD(NewString(3, 0)), AD(NewString(12, 0)), ED())
}

func TestRulesRID(t *testing.T) {
	assertEventsMaxDepth(t, 1, RB(), AC(18, false), AD([]byte("http://example.com")), ED())
}

func TestRulesResourceIDMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, RB(), AC(13, false), AD([]byte("http:")), AD([]byte("test")), AD([]byte(".net")), ED())
}

func testRulesCustomBinary(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, CBB())
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, AC(0, false))
	} else {
		for i, count := range byteCount {
			assertEventsSucceed(t, rules, AC(uint64(count), i != lastIndex))
			if count > 0 {
				assertEventsSucceed(t, rules, AD(NewString(count, 0)))
			}
		}
	}
	assertEventsSucceed(t, rules, ED())
}

func TestRulesCustomBinary(t *testing.T) {
	testRulesCustomBinary(t, 0)
	testRulesCustomBinary(t, 1, 1)
	testRulesCustomBinary(t, 2, 2)
	testRulesCustomBinary(t, 10, 10)
	testRulesCustomBinary(t, 100, 14, 55, 20, 11)
}

func TestRulesCustomBinaryMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, CBB(), AC(10, false), AD(NewBytes(5, 0)), AD(NewBytes(3, 0)), AD(NewBytes(2, 0)), ED())
}

func testRulesCustomText(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, CTB())
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, AC(0, false))
	} else {
		for i, count := range byteCount {
			assertEventsSucceed(t, rules, AC(uint64(count), i != lastIndex))
			if count > 0 {
				assertEventsSucceed(t, rules, AD(NewBytes(count, 0)))
			}
		}
	}
	assertEventsSucceed(t, rules, ED())
}

func TestRulesCustomText(t *testing.T) {
	testRulesCustomText(t, 0)
	testRulesCustomText(t, 1, 1)
	testRulesCustomText(t, 2, 2)
	testRulesCustomText(t, 10, 10)
	testRulesCustomText(t, 100, 14, 55, 20, 11)
}

func TestRulesCustomTextMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, CTB(), AC(10, false), AD(NewString(5, 0)), AD(NewString(3, 0)), AD(NewString(2, 0)), ED())
}

func TestRulesInArrayBasic(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AU8B())
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(4, false))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AD(NewString(4, 0)))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, AC(0, true))
}

func TestRulesInArray(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, SB())
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(10, true))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AD(NewString(4, 0)))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AD(NewString(4, 0)))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AD(NewString(2, 0)))
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(5, false))
	assertEventsSucceed(t, rules, AD(NewString(5, 0)))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, RB())
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(5, false))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AD([]byte("a:123")))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, E())
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, AC(0, true))
}

func TestRulesInArrayEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, SB())
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(0, false))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AU8B())
	assertEventsSucceed(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, AC(0, false))
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, E())
	assertEventsFail(t, rules, AC(0, true))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, AC(0, true))
}

// =================
// Containers: Empty
// =================

func TestRulesListEmpty(t *testing.T) {
	assertEventsMaxDepth(t, 1, L(), E(), ED())
}

func TestRulesMapEmpty(t *testing.T) {
	assertEventsMaxDepth(t, 1, M(), E(), ED())
}

func TestRulesMarkupEmpty(t *testing.T) {
	assertEventsMaxDepth(t, 2, MU("a"), E(), E(), ED())
}

func TestRulesNodeEmpty(t *testing.T) {
	assertEventsMaxDepth(t, 1, NODE(), NULL(), E(), ED())
}

// =======================
// Containers: Single item
// =======================

func TestRulesListSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 1, L(), B(true), E(), ED())
}

func TestRulesMapPair(t *testing.T) {
	assertEventsMaxDepth(t, 1, M(), B(true), NULL(), E(), ED())
}

func TestRulesMarkupSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 2, MU("abcdef"), I(-1), B(true), E(), S("a"), E(), ED())
}

func TestRulesNodeSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 1, NODE(), B(true), E(), ED())
	assertEventsMaxDepth(t, 1, NODE(), NULL(), B(true), E(), ED())
}

// ==================
// Containers: Filled
// ==================

func TestRulesListFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, L(), NULL(), NAN(), B(true), BF(0.1), I(1), I(-1),
		GT(time.Now()), AU8(NewBytes(1, 0)), E(), ED())
}

func TestRulesMapFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), B(true), NULL(), BF(0.1), NAN(), I(1), I(-1),
		GT(time.Now()), AU8(NewBytes(1, 0)), E(), ED())
}

func TestRulesMapList(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), I(1), L(), E(), E(), ED())
}

func TestRulesMapMap(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), I(1), M(), E(), E(), ED())
}

func TestRulesDeepContainer(t *testing.T) {
	assertEventsMaxDepth(t, 6, L(), L(), M(), I(-1), M(), I(1), L(),
		S("0123456789"), E(), E(), E(), E(), E(), ED())
}

func TestRulesMarkup(t *testing.T) {
	assertEventsMaxDepth(t, 2, MU("a"), I(1), I(-1), E(), S("a"), E(), ED())
}

func TestRulesMarkerReference(t *testing.T) {
	assertEventsMaxDepth(t, 9, M(),
		S("keys"),
		L(),
		MARK("1"), S("foo"),
		MARK("2"), S("bar"),
		E(),
		REF("1"), I(1),
		REF("2"), I(2),
		E())
}

func TestRulesNodeFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, NODE(), NULL(), NAN(), B(true), BF(0.1), I(1), I(-1),
		GT(time.Now()), AU8(NewBytes(1, 0)), NODE(), NULL(), E(), E(), ED())
}

// ================
// Error conditions
// ================

func TestRulesErrorOnEndTooManyTimes(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L(), E())
	assertEventsFail(t, rules, E())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, M(), E())
	assertEventsFail(t, rules, E())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MU("a"), E(), E())
	assertEventsFail(t, rules, E())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, NODE(), NULL(), E())
	assertEventsFail(t, rules, E())
}

func TestRulesErrorUnendedContainer(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, M())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MU("a"))
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MU("a"), E())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, NODE())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, NODE(), NULL())
	assertEventsFail(t, rules, ED())
}

func TestRulesErrorArrayTooManyBytes(t *testing.T) {
	for _, arrayType := range test.ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, AC(1, false))
		assertEventsFail(t, rules, AD([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}))
	}
}

func TestRulesErrorArrayTooFewBytes(t *testing.T) {
	for _, arrayType := range test.ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, AC(2000, false),
			AD([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6}))
		assertEventsFail(t, rules, ED())
	}
}

func TestRulesErrorMarkupNameLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	assertEventsFail(t, rules, MU(""))
}

func TestRulesErrorMarkerIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	assertEventsFail(t, rules, MARK(""))
}

func TestRulesErrorReferenceIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(3)
	assertEventsFail(t, rules, L(), REF(""))
}

func TestRulesErrorRefMissingMarker(t *testing.T) {
	rules := newRulesWithMaxDepth(5)
	assertEventsSucceed(t, rules, L(), REF("test"), E())
	assertEventsFail(t, rules, ED())
}

func TestRulesResourceIDLength0_1(t *testing.T) {
	assertEventsMaxDepth(t, 2, RID(""))
	assertEventsMaxDepth(t, 2, RID("a"))
}

func TestRulesErrorDuplicateMarkerID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L(), MARK("test"), TT())
	// TODO: Can we check for duplicate markers before the marked object is read?
	assertEventsFail(t, rules, MARK("test"), FF())
}

func TestRulesErrorInvalidMarkerID(t *testing.T) {
	prefix := func() events.DataEventReceiver {
		rules := newRulesWithMaxDepth(10)
		return rules
	}
	assertEachEventFails(t, prefix, MARK(""))

	for ch := 0; ch <= 0xff; ch++ {
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if ch == '_' {
			continue
		}
		// TODO: Fix marker ID validation test
		// assertEachEventFails(t, prefix, S(string([]byte{byte(ch)})))
	}
}

// ======
// Limits
// ======

// TODO: Other array byte lengths

func TestRulesMaxBytesLength(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxArrayByteLength = 10
	rules := newRulesAfterVersion(opts)
	assertEventsFail(t, rules, AU8(NewBytes(11, 0)))

	rules = newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, AU8B(), AC(8, true), AD(NewBytes(8, 0)))
	assertEventsFail(t, rules, AC(4, false))
}

func TestRulesMaxStringLength(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxStringByteLength = 10
	rules := newRulesAfterVersion(opts)
	assertEventsFail(t, rules, S("12345678901"))

	rules = newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, SB(), AC(8, true), AD(NewBytes(8, 40)))
	assertEventsFail(t, rules, AC(4, false))
}

func TestRulesMaxResourceIDLength(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxResourceIDByteLength = 10
	rules := newRulesAfterVersion(opts)
	assertEventsFail(t, rules, RID("12345678901"))

	rules = newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, RB(), AC(8, true), AD(NewBytes(8, 64)))
	assertEventsFail(t, rules, AC(4, false))
}

func TestRulesMaxIDLength(t *testing.T) {
	maxIDLength := 127
	opts := options.DefaultRuleOptions()
	rules := newRulesAfterVersion(opts)
	assertEventsFail(t, rules, MARK(string(NewString(maxIDLength+1, 0))))

	rules = newRulesAfterVersion(opts)
	assertEventsFail(t, rules, REF(string(NewString(maxIDLength+1, 0))))
}

func TestRulesMaxContainerDepth(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, L())
}

func TestRulesMaxObjectCount(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxObjectCount = 3
	rules := newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, L(), S("test"), TT())
	assertEventsFail(t, rules, FF())
}

func TestRulesMaxReferenceCount(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxReferenceCount = 2
	rules := newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, L(), MARK("test"), TT(), MARK("10"), TT())
	assertEventsFail(t, rules, MARK("xx"), TT())
}

func TestRulesReset(t *testing.T) {
	opts := options.DefaultRuleOptions()
	opts.MaxContainerDepth = 2
	opts.MaxObjectCount = 5
	opts.MaxReferenceCount = 2
	rules := newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, L())
	rules.Reset()
	assertEventsFail(t, rules, E())
	assertEventsSucceed(t, rules, BD(), EvV, L(), L(), I(1), I(1), I(1), E())
	assertEventsFail(t, rules, I(1))
	rules.Reset()
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK("1"), S("test"), MARK("2"), S("more tests"))
	assertEventsFail(t, rules, MARK("xx"), TT())
	rules.Reset()
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK("1"), S("test"))
}

func TestTopLevelStringLikeReferenceID(t *testing.T) {
	opts := options.DefaultRuleOptions()
	rules := NewRules(events.NewNullEventReceiver(), opts)
	assertEventsSucceed(t, rules, BD(), EvV, RREF("http://x.y"), ED())
}

func TestRulesForwardReference(t *testing.T) {
	rules := newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, M(),
		S("marked"), MARK("x"), M(),
		S("recursive"), REF("x"),
		E(),
		E())
}

func TestRulesIdentifier(t *testing.T) {
	rules := newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("12\u0001abc"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK(""), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("a+b"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("a+b"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("-a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK(":a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK(".a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("a:"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("a::a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsFail(t, rules, MARK("\u0300a"), I(1))

	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("1"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a-a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("_a"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a_"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a-"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("0_"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("0-"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a-_a12-_3_gy"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("人気"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a...."), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MARK("a\u0300"), I(1))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MU("a:b"))
	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, MU("a:b:c:d:e:f"))
}

func TestRulesMultichunk(t *testing.T) {
	rules := newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, SB(), AC(1, true), AD([]byte{'a'}), AC(0, false))
}

func TestRulesEdge(t *testing.T) {
	rules := newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, EDGE(), RID("x"), RID("y"), I(1))

	rules = newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, EDGE(), RID("a"), RID("b"), I(1), ED())
}

// =============
// Allowed Types
// =============

func TestRulesAllowedTypesTLO(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.ValidTLOValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.InvalidTLOValues))
}

func TestRulesAllowedTypesList(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{L()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{L()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidListValues))
}

func TestRulesAllowedTypesMapKey(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{M()},
			[]*test.TEvent{I(1), E()},
			[]*test.TEvent{ED()},
			test.ValidMapKeys))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{M()},
			[]*test.TEvent{I(1), E()},
			[]*test.TEvent{ED()},
			test.InvalidMapKeys))
}

func TestRulesAllowedTypesMapValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{M(), TT()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidMapValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{M(), TT()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidMapValues))
}

func TestRulesAllowedTypesMarkupAttributeKey(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a")},
			[]*test.TEvent{I(1), E(), E()},
			[]*test.TEvent{ED()},
			test.ValidMapKeys))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a")},
			[]*test.TEvent{I(1), E(), E()},
			[]*test.TEvent{ED()},
			test.InvalidMapKeys))
}

func TestRulesAllowedTypesMarkupAttributeValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a"), TT()},
			[]*test.TEvent{E(), E()},
			[]*test.TEvent{ED()},
			test.ValidMapValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a"), TT()},
			[]*test.TEvent{E(), E()},
			[]*test.TEvent{ED()},
			test.InvalidMapValues))
}

func TestRulesAllowedTypesMarkupContents(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a"), E()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidMarkupContents))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MU("a"), E()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidMarkupContents))
}

func TestRulesAllowedTypesNonStringArrayBegin(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, arrayType := range test.NonStringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, arrayType := range test.NonStringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsFail(t, rules, event)
			}
		}
	}

	assertSuccess(test.ValidAfterNonStringArrayBegin...)
	assertFail(test.InvalidAfterNonStringArrayBegin...)
}

func TestRulesAllowedTypesStringArrayBegin(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, arrayType := range test.StringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, arrayType := range test.StringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsFail(t, rules, event)
			}
		}
	}

	assertSuccess(test.ValidAfterStringArrayBegin...)
	assertFail(test.InvalidAfterStringArrayBegin...)
}

func TestRulesAllowedTypesArrayChunk(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, arrayType := range test.ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, AC(1, false))
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, arrayType := range test.ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, AC(1, false))
				assertEventsFail(t, rules, event)
			}
		}
	}

	// TODO: Make sure end of array is properly aligned to size width
	assertSuccess(test.ValidAfterArrayChunk...)
	assertFail(test.InvalidAfterArrayChunk...)
}

func TestRulesAllowedTypesMarkerValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MARK("1")},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.ValidMarkerValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{MARK("1")},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.InvalidMarkerValues))
}

func TestRulesAllowedTypesNodeValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidNodeValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE()},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidNodeValues))
}

func TestRulesAllowedTypesNodeChild(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE(), I(1)},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE(), I(1)},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidListValues))

	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE(), I(1), I(1)},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{NODE(), I(1), I(1)},
			[]*test.TEvent{E()},
			[]*test.TEvent{ED()},
			test.InvalidListValues))
}

func TestRulesAllowedTypesEdgeSource(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE()},
			[]*test.TEvent{RID("a"), I(1)},
			[]*test.TEvent{ED()},
			test.ValidEdgeSources))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE()},
			[]*test.TEvent{RID("a"), I(1)},
			[]*test.TEvent{ED()},
			test.InvalidEdgeSources))
}

func TestRulesAllowedTypesEdgeDescription(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE(), RID("a")},
			[]*test.TEvent{I(1)},
			[]*test.TEvent{ED()},
			test.ValidEdgeDescriptions))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE(), RID("a")},
			[]*test.TEvent{I(1)},
			[]*test.TEvent{ED()},
			test.InvalidEdgeDescriptions))
}

func TestRulesAllowedTypesEdgeDestination(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE(), RID("a"), RID("b")},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.ValidEdgeDestinations))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			[]*test.TEvent{BD(), V(ceVer)},
			[]*test.TEvent{EDGE(), RID("a"), RID("b")},
			[]*test.TEvent{},
			[]*test.TEvent{ED()},
			test.InvalidEdgeDescriptions))
}

func TestMedia(t *testing.T) {
	rules := NewRules(events.NewNullEventReceiver(), nil)
	assertEventsSucceed(t, rules, BD(), V(0), MB(), AC(1, false), AD([]byte("a")), AC(0, false), ED())
}

func TestComment(t *testing.T) {
	assertEvents(t, BD(), V(ceVer), COM(true, "a"), S("b"), ED())
	assertEvents(t, BD(), V(ceVer), COM(false, "a"), S("b"), ED())

	assertEvents(t, BD(), V(ceVer), L(), COM(true, "a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), L(), COM(true, "a"), S("b"), E(), ED())
	assertEvents(t, BD(), V(ceVer), L(), S("b"), COM(true, "a"), E(), ED())

	assertEvents(t, BD(), V(ceVer), M(), COM(true, "a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), S("b"), S("b"), COM(true, "a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), S("b"), COM(true, "a"), S("b"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), COM(true, "a"), S("b"), S("b"), E(), ED())

	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), E(), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), S("b"), S("b"), COM(true, "a"), E(), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), S("b"), COM(true, "a"), S("b"), E(), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), S("b"), S("b"), E(), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), E(), COM(false, "x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), E(), S("a"), COM(false, "x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), E(), S("a"), COM(false, "x"), S("a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), MU("a"), COM(true, "a"), E(), S("a"), COM(false, "x"), COM(false, "x"), S("a"), E(), ED())

	assertEvents(t, BD(), V(ceVer), NODE(), COM(true, "a"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), NODE(), COM(true, "a"), S("x"), COM(true, "a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), NODE(), COM(true, "a"), S("x"), COM(true, "a"), S("x"), S("x"), COM(true, "a"), E(), ED())

	assertEvents(t, BD(), V(ceVer), EDGE(), COM(true, "a"), S("x"), S("x"), S("x"), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), COM(true, "a"), COM(true, "a"), S("x"), S("x"), S("x"), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), COM(true, "a"), S("x"), COM(true, "a"), S("x"), S("x"), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), COM(true, "a"), S("x"), S("x"), COM(true, "a"), S("x"), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), COM(true, "a"), COM(true, "a"), COM(true, "a"), S("x"), COM(true, "a"), COM(true, "a"), S("x"), COM(true, "a"), COM(true, "a"), S("x"), ED())

	// assertEvents(t, BD(), V(ceVer), AI8B(), AC(1, true), AD([]byte{1}), COM(true, "a"), AC(1, false), AD([]byte{1}), ED())
}

func TestEdgeMaxDepth(t *testing.T) {
	assertEvents(t, BD(), V(ceVer),
		L(), MARK("a"), PI(1), MARK("b"), PI(2), MARK("c"), PI(3), MARK("d"), PI(4),
		EDGE(), REF("a"), PI(1), REF("b"),
		EDGE(), REF("a"), PI(1), REF("c"),
		EDGE(), REF("b"), PI(1), REF("c"),
		EDGE(), REF("b"), PI(1), REF("d"),
		EDGE(), REF("c"), M(), E(), REF("d"),
		E(),
		ED())
}
