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

func TestRulesNA(t *testing.T) {
	assertEventsMaxDepth(t, 1, NA(), I(1), ED())
}

func TestRulesNil(t *testing.T) {
	assertEventsMaxDepth(t, 1, N(), ED())
}

func TestRulesNan(t *testing.T) {
	assertEventsMaxDepth(t, 1, NAN(), ED())
	assertEventsMaxDepth(t, 1, SNAN(), ED())
}

func TestRulesNan2(t *testing.T) {
	assertEventsMaxDepth(t, 1, F(math.NaN()), ED())
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
	assertEventsMaxDepth(t, 1, BI(NewBigInt("1", 10)), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("-1", 10)), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("10000000000000000000000000000000000000000000", 10)), ED())
	assertEventsMaxDepth(t, 1, BI(NewBigInt("-10000000000000000000000000000000000000000000", 10)), ED())
}

func TestRulesBigIntAsID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsFail(t, rules, L(), MARK(), BI(NewBigInt("10", 10)), I(0), REF(), I(10), E(), ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MARK())
	assertEventsFail(t, rules, BI(NewBigInt("10000000000000000000000000000000000000000000", 10)))
}

func TestRulesFloat(t *testing.T) {
	assertEventsMaxDepth(t, 1, F(0.1), ED())
	assertEventsMaxDepth(t, 1, F(math.NaN()), ED())
}

func TestRulesBigFloat(t *testing.T) {
	assertEventsMaxDepth(t, 1, BF(NewBigFloat("1.1", 10, 2)), ED())
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

func TestRulesUUID(t *testing.T) {
	assertEventsMaxDepth(t, 1, UUID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), ED())
}

func TestRulesTime(t *testing.T) {
	assertEventsMaxDepth(t, 1, GT(time.Now()), ED())
}

func TestRulesCompactTime(t *testing.T) {
	assertEventsMaxDepth(t, 1, CT(NewDate(1, 1, 1)), ED())
}

func TestRulesArrayOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, AU8([]byte{1, 2, 3, 4}), ED())
	assertEventsMaxDepth(t, 1, AU16([]uint16{1, 2, 3, 4}), ED())
	// TODO: Other array types
}

func TestRulesResourceIDOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, RID("http://example.com"), ED())
	assertEventsMaxDepth(t, 1, RIDCat("http://example.com"), I(1), ED())
}

func TestRulesCustomOneshot(t *testing.T) {
	assertEventsMaxDepth(t, 1, CUB([]byte{1, 2, 3, 4}), ED())
	assertEventsMaxDepth(t, 1, CUT("test"), ED())
}

func TestRulesMarker(t *testing.T) {
	assertEventsMaxDepth(t, 10, L(), MARK(), ID("a"), F(0.1), MARK(), ID("blah"), GT(time.Now()), E(), ED())
}

func TestRulesReference(t *testing.T) {
	rules := newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(), MARK(), ID("a"), F(0.1), MARK(), ID("blah"), GT(time.Now()),
		REF(), ID("a"),
		REF(), RID("http://example.com"),
		REF(), ID("blah"),
		S("test"),
		E())

	rules = newRulesWithMaxDepth(10)

	assertEventsSucceed(t, rules, L(),
		MARK(), ID("a"), F(0.1),
		MARK(), ID("blah"), GT(time.Now()),
		REF(), ID("a"), REF(), RIDCat("http://example.com"), S("1"),
		REF(), ID("blah"),
		S("test"),
		E())
}

func TestRulesConstant(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, M(), S("a"), CONST("something", true), PI(100), E())

	rules = newRulesWithMaxDepth(10)
	assertEventsFail(t, rules, M(), S("a"), CONST("something", true), E())

	rules = newRulesWithMaxDepth(10)
	assertEventsFail(t, rules, M(), S("a"), CONST("something", false), E())
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
	assertEventsMaxDepth(t, 2, MUP(), ID("a"), E(), E(), ED())
}

func TestRulesCommentEmpty(t *testing.T) {
	assertEventsMaxDepth(t, 1, CMT(), E(), I(1), ED())
}

// =======================
// Containers: Single item
// =======================

func TestRulesListSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 1, L(), B(true), E(), ED())
}

func TestRulesMapPair(t *testing.T) {
	assertEventsMaxDepth(t, 1, M(), B(true), N(), E(), ED())
}

func TestRulesMarkupSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 2, MUP(), ID("abcdef"), I(-1), B(true), E(), S("a"), E(), ED())
}

func TestRulesCommentSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 2, CMT(), S("a"), E(), I(1), ED())
}

// ==================
// Containers: Filled
// ==================

func TestRulesListFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, L(), N(), NAN(), B(true), F(0.1), I(1), I(-1),
		GT(time.Now()), AU8(NewBytes(1, 0)), E(), ED())
}

func TestRulesMapFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), B(true), N(), F(0.1), NAN(), I(1), I(-1),
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

func TestRulesCommentInt(t *testing.T) {
	// TODO: Recheck comment validation
	rules := newRulesWithMaxDepth(2)
	assertEventsSucceed(t, rules, CMT(), S("blah\r\n\t\tblah"), SB(), AC(1, false))
	// assertEventsFail(t, rules, AD([]byte{0x00}))
	// assertEventsFail(t, rules, AD([]byte{0x0b}))
	// assertEventsFail(t, rules, AD([]byte{0x7f}))
	// assertEventsFail(t, rules, AD([]byte{0x80}))
	// assertEvents(t, rules, AD([]byte{0x40}), E(), I(1), ED())
}

func TestRulesCommentMap(t *testing.T) {
	assertEventsMaxDepth(t, 3, M(), CMT(), S("a"), E(), I(1), I(-1), E(), ED())
}

func TestRulesMarkup(t *testing.T) {
	assertEventsMaxDepth(t, 2, MUP(), ID("a"), I(1), I(-1), E(), S("a"), E(), ED())
}

// =============
// Allowed Types
// =============

func TestRulesAllowedTypesTLO(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			assertEventsMaxDepth(t, 1, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(1)
			assertEventsFail(t, rules, event)
		}
	}
	assertSuccess(ValidTLOValues...)
	assertFail(InvalidTLOValues...)
}

func TestRulesAllowedTypesList(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, L())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, L())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidListValues...)
	assertFail(InvalidListValues...)
}

func TestRulesAllowedTypesMapKey(t *testing.T) {
	prefix := func() events.DataEventReceiver {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, M())
		return rules
	}
	assertEachEventSucceeds(t, prefix, ValidMapKeys...)
	assertEachEventFails(t, prefix, InvalidMapKeys...)

	prefix = func() events.DataEventReceiver {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, M(), REF())
		return rules
	}
	assertEachEventSucceeds(t, prefix, KeyableReferenceIDTypes...)
	assertEachEventFails(t, prefix, InvalidKeyableReferenceIDTypes...)
}

func TestRulesAllowedTypesMapValue(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, M(), TT())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, M(), TT())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMapValues...)
	assertFail(InvalidMapValues...)
}

func TestRulesAllowedTypesComment(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, CMT())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, CMT())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidCommentValues...)
	assertFail(InvalidCommentValues...)
}

func TestRulesAllowedTypesMarkupName(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMarkupNames...)
	assertFail(InvalidMarkupNames...)
}

func TestRulesAllowedTypesMarkupAttributeKey(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"))
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"))
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMapKeys...)
	assertFail(InvalidMapKeys...)
}

func TestRulesAllowedTypesMarkupAttributeValue(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"), TT())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"), TT())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMapValues...)
	assertFail(InvalidMapValues...)
}

func TestRulesAllowedTypesMarkupContents(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"), E())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MUP(), ID("a"), E())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMarkupContents...)
	assertFail(InvalidMarkupContents...)
}

func TestRulesAllowedTypesArrayBegin(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, arrayType := range ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, arrayType := range ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsFail(t, rules, event)
			}
		}
	}

	assertSuccess(ValidAfterArrayBegin...)
	assertFail(InvalidAfterArrayBegin...)
}

func TestRulesAllowedTypesArrayChunk(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, arrayType := range ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, AC(1, false))
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, arrayType := range ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, AC(1, false))
				assertEventsFail(t, rules, event)
			}
		}
	}

	// TODO: Make sure end of array is properly aligned to size width
	assertSuccess(ValidAfterArrayChunk...)
	assertFail(InvalidAfterArrayChunk...)
}

func TestRulesAllowedTypesMarkerID(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MARK())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MARK())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMarkerIDs...)
	assertFail(InvalidMarkerIDs...)
}

func TestRulesAllowedTypesMarkerValue(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MARK(), ID("1"))
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, MARK(), ID("1"))
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidMarkerValues...)
	assertFail(InvalidMarkerValues...)
}

func TestRulesAllowedTypesReferenceID(t *testing.T) {
	assertSuccess := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, L(), REF())
			assertEventsSucceed(t, rules, event)
		}
	}
	assertFail := func(events ...*test.TEvent) {
		for _, event := range events {
			rules := newRulesWithMaxDepth(10)
			assertEventsSucceed(t, rules, L(), REF())
			assertEventsFail(t, rules, event)
		}
	}

	assertSuccess(ValidReferenceIDs...)
	assertFail(InvalidReferenceIDs...)
}

func TestRulesMarkerReference(t *testing.T) {
	assertEventsMaxDepth(t, 9, L(), MARK(), ID("1"), TT(), MARK(), ID("2"), REF(), ID("1"), E())
}

func TestRulesMarkerReference2(t *testing.T) {
	assertEventsMaxDepth(t, 9, M(),
		S("keys"),
		L(),
		MARK(), ID("1"), S("foo"),
		MARK(), ID("2"), S("bar"),
		E(),
		REF(), ID("1"), I(1),
		REF(), ID("2"), I(2),
		E())
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
	assertEventsSucceed(t, rules, CMT(), E())
	assertEventsFail(t, rules, E())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MUP(), ID("a"), E(), E())
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
	assertEventsSucceed(t, rules, CMT())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MUP())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MUP(), ID("a"))
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MUP(), ID("a"), E())
	assertEventsFail(t, rules, ED())
}

func TestRulesErrorArrayTooManyBytes(t *testing.T) {
	for _, arrayType := range ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, AC(1, false))
		assertEventsFail(t, rules, AD([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}))
	}
}

func TestRulesErrorArrayTooFewBytes(t *testing.T) {
	for _, arrayType := range ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, AC(2000, false),
			AD([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6}))
		assertEventsFail(t, rules, ED())
	}
}

func TestRulesErrorMarkupNameLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	assertEventsSucceed(t, rules, MUP())
	assertEventsFail(t, rules, ID(""))
}

func TestRulesErrorMarkerIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	assertEventsSucceed(t, rules, MARK())
	assertEventsFail(t, rules, ID(""))
}

func TestRulesErrorReferenceIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(3)
	assertEventsSucceed(t, rules, L(), REF())
	assertEventsFail(t, rules, ID(""))
}

func TestRulesErrorRefMissingMarker(t *testing.T) {
	rules := newRulesWithMaxDepth(5)
	assertEventsSucceed(t, rules, L(), REF(), ID("test"), E())
	assertEventsFail(t, rules, ED())
}

func TestRulesResourceIDLength0_1(t *testing.T) {
	assertEventsMaxDepth(t, 2, RID(""))
	assertEventsMaxDepth(t, 2, RID("a"))
}

func TestRulesErrorDuplicateMarkerID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L(), MARK(), ID("test"), TT(), MARK())
	// TODO: Can we check for duplicate markers before the marked object is read?
	assertEventsFail(t, rules, ID("test"), FF())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L(), MARK(), ID("100"), TT(), MARK())
	assertEventsFail(t, rules, ID("100"), FF())
}

func TestRulesErrorInvalidMarkerID(t *testing.T) {
	prefix := func() events.DataEventReceiver {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, MARK())
		return rules
	}
	assertEachEventFails(t, prefix, S(""), S("1z"), S("a:z"))

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
		// assertEachEventFails(t, prefix, S(string([]byte{byte(ch)})))
	}

	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, MARK())
	assertEventsFail(t, rules, SB(), AC(0, false), AD([]byte{}))
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
	assertEventsSucceed(t, rules, MARK())
	assertEventsFail(t, rules, ID(string(NewString(maxIDLength+1, 0))))

	rules = newRulesAfterVersion(opts)
	assertEventsSucceed(t, rules, REF())
	assertEventsFail(t, rules, ID(string(NewString(maxIDLength+1, 0))))
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
	assertEventsSucceed(t, rules, L(), MARK(), ID("test"), TT(), MARK(), ID("10"), TT())
	assertEventsFail(t, rules, MARK())
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
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK(), ID("1"), S("test"), MARK(), ID("2"), S("more tests"))
	assertEventsFail(t, rules, MARK())
	rules.Reset()
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK(), ID("1"), S("test"))
}

func TestTopLevelStringLikeReferenceID(t *testing.T) {
	opts := options.DefaultRuleOptions()
	rules := NewRules(events.NewNullEventReceiver(), opts)
	assertEventsSucceed(t, rules, BD(), EvV, REF(), RID("http://x.y"), ED())
}

func TestRulesForwardReference(t *testing.T) {
	rules := newRulesAfterVersion(nil)
	assertEventsSucceed(t, rules, M(),
		S("marked"), MARK(), ID("x"), M(),
		S("recursive"), REF(), ID("x"),
		E(),
		E())
}
