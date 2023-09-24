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
	"testing"
	"time"

	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/nullevent"
	"github.com/kstenerud/go-concise-encoding/test"
)

// ===========
// Array Types
// ===========

func TestRulesBytes(t *testing.T) {
	testRulesBytes := func(t *testing.T, length int, byteCount ...int) {
		rules := newRulesWithMaxDepth(1)
		assertEventsSucceed(t, rules, BAU8())

		lastIndex := len(byteCount) - 1
		if lastIndex < 0 {
			assertEventsSucceed(t, rules, ACL(0))
		} else {
			for i, count := range byteCount {
				if i < lastIndex {
					assertEventsSucceed(t, rules, ACM(uint64(count)))
				} else {
					assertEventsSucceed(t, rules, ACL(uint64(count)))
				}
				if count > 0 {
					assertEventsSucceed(t, rules, ADU8(make([]byte, count)))
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
	assertEventsMaxDepth(t, 1, BAU8(), ACL(10), ADU8(NewBytes(5, 0)), ADU8(NewBytes(3, 0)), ADU8(NewBytes(2, 0)), ED())
}

func testRulesString(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, BS())

	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, ACL(0))
	} else {
		for i, count := range byteCount {
			if i < lastIndex {
				assertEventsSucceed(t, rules, ACM(uint64(count)))
			} else {
				assertEventsSucceed(t, rules, ACL(uint64(count)))
			}
			if count > 0 {
				assertEventsSucceed(t, rules, ADT(NewString(count, 0)))
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
	assertEventsMaxDepth(t, 1, BS(), ACL(20), ADT(NewString(5, 0)), ADT(NewString(3, 0)), ADT(NewString(12, 0)), ED())
}

func TestRulesRID(t *testing.T) {
	assertEventsMaxDepth(t, 1, BRID(), ACL(18), ADT("http://example.com"), ED())
}

func TestRulesResourceIDMultiChunk(t *testing.T) {
	assertEventsMaxDepth(t, 1, BRID(), ACL(13), ADT("http:"), ADT("test"), ADT(".net"), ED())
}

func testRulesCustomBinary(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, BCB(1))
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, ACL(0))
	} else {
		for i, count := range byteCount {
			if i < lastIndex {
				assertEventsSucceed(t, rules, ACM(uint64(count)))
			} else {
				assertEventsSucceed(t, rules, ACL(uint64(count)))
			}
			if count > 0 {
				assertEventsSucceed(t, rules, ADT(NewString(count, 0)))
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
	assertEventsMaxDepth(t, 1, BCB(1), ACL(10), ADU8(NewBytes(5, 0)), ADU8(NewBytes(3, 0)), ADU8(NewBytes(2, 0)), ED())
}

func testRulesCustomText(t *testing.T, length int, byteCount ...int) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, BCT(1))
	lastIndex := len(byteCount) - 1
	if lastIndex < 0 {
		assertEventsSucceed(t, rules, ACL(0))
	} else {
		for i, count := range byteCount {
			if i < lastIndex {
				assertEventsSucceed(t, rules, ACM(uint64(count)))
			} else {
				assertEventsSucceed(t, rules, ACL(uint64(count)))
			}
			if count > 0 {
				assertEventsSucceed(t, rules, ADU8(NewBytes(count, 0)))
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
	assertEventsMaxDepth(t, 1, BCT(1), ACL(10), ADT(NewString(5, 0)), ADT(NewString(3, 0)), ADT(NewString(2, 0)), ED())
}

func TestRulesInArrayBasic(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, BAU8())
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACL(4))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ADT(NewString(4, 0)))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, ACM(0))
}

func TestRulesInArray(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, BS())
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACM(10))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ADT(NewString(4, 0)))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ADT(NewString(4, 0)))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ADT(NewString(2, 0)))
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACL(5))
	assertEventsSucceed(t, rules, ADT(NewString(5, 0)))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, BRID())
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACL(5))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ADT("a:123"))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, E())
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, ACM(0))
}

func TestRulesInArrayEmpty(t *testing.T) {
	rules := newRulesWithMaxDepth(9)
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, BS())
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACL(0))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, BAU8())
	assertEventsSucceed(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ACL(0))
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, E())
	assertEventsFail(t, rules, ACM(0))
	assertEventsSucceed(t, rules, ED())
	assertEventsFail(t, rules, ACM(0))
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

func TestRulesNodeSingleItem(t *testing.T) {
	assertEventsMaxDepth(t, 1, NODE(), B(true), E(), ED())
	assertEventsMaxDepth(t, 1, NODE(), NULL(), B(true), E(), ED())
}

// ==================
// Containers: Filled
// ==================

func TestRulesListFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, L(), NULL(), NAN(), B(true), N(0.1), N(1), N(-1),
		T(compact_time.AsCompactTime(time.Now())), AU8(NewBytes(1, 0)), E(), ED())
}

func TestRulesMapFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), B(true), NULL(), N(1234567890), NAN(), N(1), N(-1),
		T(compact_time.AsCompactTime(time.Now())), AU8(NewBytes(1, 0)), E(), ED())
}

func TestRulesMapList(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), N(1), L(), E(), E(), ED())
}

func TestRulesMapMap(t *testing.T) {
	assertEventsMaxDepth(t, 2, M(), N(1), M(), E(), E(), ED())
}

func TestRulesDeepContainer(t *testing.T) {
	assertEventsMaxDepth(t, 6, L(), L(), M(), N(-1), M(), N(1), L(),
		S("0123456789"), E(), E(), E(), E(), E(), ED())
}

func TestRulesMarkerLocalReference(t *testing.T) {
	assertEventsMaxDepth(t, 9, M(),
		S("keys"),
		L(),
		MARK("1"), S("foo"),
		MARK("2"), S("bar"),
		E(),
		REFL("1"), N(1),
		REFL("2"), N(2),
		E())
}

func TestRulesNodeFilled(t *testing.T) {
	assertEventsMaxDepth(t, 2, NODE(), NULL(), NAN(), B(true), N(0.1), N(1), N(-1),
		T(compact_time.AsCompactTime(time.Now())), AU8(NewBytes(1, 0)), NODE(), NULL(), E(), E(), ED())
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
	assertEventsSucceed(t, rules, NODE())
	assertEventsFail(t, rules, ED())

	rules = newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, NODE(), NULL())
	assertEventsFail(t, rules, ED())
}

func TestRulesErrorArrayTooManyBytes(t *testing.T) {
	for _, arrayType := range test.ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, ACL(1))
		assertEventsFail(t, rules, ADU8([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}))
	}
}

func TestRulesErrorArrayTooFewBytes(t *testing.T) {
	for _, arrayType := range test.ArrayBeginTypes {
		rules := newRulesWithMaxDepth(10)
		assertEventsSucceed(t, rules, arrayType, ACL(2000),
			ADU8([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 1, 2, 3, 4, 5, 6}))
		assertEventsFail(t, rules, ED())
	}
}

func TestRulesErrorMarkerIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(2)
	assertEventsFail(t, rules, MARK(""))
}

func TestRulesErrorReferenceIDLength0(t *testing.T) {
	rules := newRulesWithMaxDepth(3)
	assertEventsFail(t, rules, L(), REFL(""))
}

func TestRulesErrorRefMissingMarker(t *testing.T) {
	rules := newRulesWithMaxDepth(5)
	assertEventsSucceed(t, rules, L(), REFL("test"), E())
	assertEventsFail(t, rules, ED())
}

func TestRulesResourceIDLength0_1(t *testing.T) {
	assertEventsMaxDepth(t, 2, RID(""))
	assertEventsMaxDepth(t, 2, RID("a"))
}

func TestRulesErrorDuplicateMarkerID(t *testing.T) {
	rules := newRulesWithMaxDepth(10)
	assertEventsSucceed(t, rules, L(), MARK("test"), B(true))
	// TODO: Can we check for duplicate markers before the marked object is read?
	assertEventsFail(t, rules, MARK("test"), B(false))
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
	config := configuration.New()
	config.Rules.MaxArraySizeBytes = 10
	rules := newRulesAfterVersion(config)
	assertEventsFail(t, rules, AU8(NewBytes(11, 0)))

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, BAU8(), ACM(8), ADU8(NewBytes(8, 0)))
	assertEventsFail(t, rules, ACL(4))
}

func TestRulesMaxStringLength(t *testing.T) {
	config := configuration.New()
	config.Rules.MaxArraySizeBytes = 10
	rules := newRulesAfterVersion(config)
	assertEventsFail(t, rules, S("12345678901"))

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, BS(), ACM(8), ADU8(NewBytes(8, 40)))
	assertEventsFail(t, rules, ACL(4))
}

func TestRulesMaxResourceIDLength(t *testing.T) {
	config := configuration.New()
	config.Rules.MaxArraySizeBytes = 10
	rules := newRulesAfterVersion(config)
	assertEventsFail(t, rules, RID("12345678901"))

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, BRID(), ACM(8), ADU8(NewBytes(8, 64)))
	assertEventsFail(t, rules, ACL(4))
}

func TestRulesMaxIDLength(t *testing.T) {
	maxIDLength := 200
	config := configuration.New()
	config.Rules.MaxIdentifierLength = uint64(maxIDLength)

	rules := newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK(string(NewString(maxIDLength+1, 0))))

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules,
		L(),
		MARK(string(NewString(maxIDLength, 0))),
		N(1),
		REFL(string(NewString(maxIDLength, 0))),
		E(),
	)
}

func TestRulesMaxContainerDepth(t *testing.T) {
	rules := newRulesWithMaxDepth(1)
	assertEventsSucceed(t, rules, L())
	assertEventsFail(t, rules, L())
}

func TestRulesMaxObjectCount(t *testing.T) {
	config := configuration.New()
	config.Rules.MaxObjectCount = 3
	rules := newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, L(), S("test"), B(true))
	assertEventsFail(t, rules, B(false))
}

func TestRulesMaxReferenceCount(t *testing.T) {
	config := configuration.New()
	config.Rules.MaxLocalReferenceCount = 2
	rules := newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, L(), MARK("test"), B(true), MARK("10"), B(true))
	assertEventsFail(t, rules, MARK("xx"), B(true))
}

func TestRulesReset(t *testing.T) {
	config := configuration.New()
	config.Rules.MaxContainerDepth = 2
	config.Rules.MaxObjectCount = 5
	config.Rules.MaxLocalReferenceCount = 2
	rules := newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, L())
	rules.Reset()
	assertEventsFail(t, rules, E())
	assertEventsSucceed(t, rules, BD(), EvV, L(), L(), N(1), N(1), N(1), E())
	assertEventsFail(t, rules, N(1))
	rules.Reset()
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK("1"), S("test"), MARK("2"), S("more tests"))
	assertEventsFail(t, rules, MARK("xx"), B(true))
	rules.Reset()
	assertEventsSucceed(t, rules, BD(), EvV, L(), MARK("1"), S("test"))
}

func TestTopLevelStringLikeReferenceID(t *testing.T) {
	config := configuration.New()
	rules := NewRules(nullevent.NewNullEventReceiver(), config)
	assertEventsSucceed(t, rules, BD(), EvV, REFR("http://x.y"), ED())
}

func TestRulesForwardLocalReference(t *testing.T) {
	rules := newRulesAfterVersion(configuration.New())
	assertEventsSucceed(t, rules, M(),
		S("marked"), MARK("x"), M(),
		S("recursive"), REFL("x"),
		E(),
		E())
}

func TestRulesIdentifier(t *testing.T) {
	config := configuration.New()

	rules := newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK("12\u0001abc"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK(""), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK("a+b"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK("a+b"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK(":a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK("a:"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsFail(t, rules, MARK("a::a"), N(1))

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("1"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a-a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("_a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a_"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a-"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("-a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a."), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK(".a"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("0_"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("0-"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a-_a12-_3_gy"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("人気"), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a...."), N(1))
	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, MARK("a\u0300"), N(1))
}

func TestRulesMultichunk(t *testing.T) {
	rules := newRulesAfterVersion(configuration.New())
	assertEventsSucceed(t, rules, BS(), ACM(1), ADT("a"), ACL(0))
}

func TestRulesEdge(t *testing.T) {
	config := configuration.New()

	rules := newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, EDGE(), RID("x"), RID("y"), N(1), E())

	rules = newRulesAfterVersion(config)
	assertEventsSucceed(t, rules, EDGE(), RID("a"), RID("b"), N(1), E(), ED())
}

// =============
// Allowed Types
// =============

func TestRulesAllowedTypesTLO(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{},
			test.Events{},
			test.Events{ED()},
			test.ValidTLOValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{},
			test.Events{},
			test.Events{ED()},
			test.InvalidTLOValues))
}

func TestRulesAllowedTypesList(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{L()},
			test.Events{E()},
			test.Events{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{L()},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidListValues))
}

func TestRulesAllowedTypesMapKey(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{M()},
			test.Events{N(1), E()},
			test.Events{ED()},
			test.ValidMapKeys))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{M()},
			test.Events{N(1), E()},
			test.Events{ED()},
			test.InvalidMapKeys))
}

func TestRulesAllowedTypesMapValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{M(), B(true)},
			test.Events{E()},
			test.Events{ED()},
			test.ValidMapValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{M(), B(true)},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidMapValues))
}

func TestRulesAllowedTypesNonStringArrayBegin(t *testing.T) {
	assertSuccess := func(events ...test.Event) {
		for _, arrayType := range test.NonStringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...test.Event) {
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
	assertSuccess := func(events ...test.Event) {
		for _, arrayType := range test.StringArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType)
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...test.Event) {
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
	assertSuccess := func(events ...test.Event) {
		for _, arrayType := range test.ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, ACL(1))
				assertEventsSucceed(t, rules, event)
			}
		}
	}
	assertFail := func(events ...test.Event) {
		for _, arrayType := range test.ArrayBeginTypes {
			for _, event := range events {
				rules := newRulesWithMaxDepth(10)
				assertEventsSucceed(t, rules, arrayType, ACL(1))
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
			test.Events{BD(), V(ceVer)},
			test.Events{MARK("1")},
			test.Events{},
			test.Events{ED()},
			test.ValidMarkerValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{MARK("1")},
			test.Events{},
			test.Events{ED()},
			test.InvalidMarkerValues))
}

func TestRulesAllowedTypesNodeValue(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE()},
			test.Events{E()},
			test.Events{ED()},
			test.ValidNodeValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE()},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidNodeValues))
}

func TestRulesAllowedTypesNodeChild(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE(), N(1)},
			test.Events{E()},
			test.Events{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE(), N(1)},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidListValues))

	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE(), N(1), N(1)},
			test.Events{E()},
			test.Events{ED()},
			test.ValidListValues))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{NODE(), N(1), N(1)},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidListValues))
}

func TestRulesAllowedTypesEdgeSource(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE()},
			test.Events{RID("a"), N(1), E()},
			test.Events{ED()},
			test.ValidEdgeSources))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE()},
			test.Events{RID("a"), N(1), E()},
			test.Events{ED()},
			test.InvalidEdgeSources))
}

func TestRulesAllowedTypesEdgeDescription(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE(), RID("a")},
			test.Events{N(1), E()},
			test.Events{ED()},
			test.ValidEdgeDescriptions))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE(), RID("a")},
			test.Events{N(1), E()},
			test.Events{ED()},
			test.InvalidEdgeDescriptions))
}

func TestRulesAllowedTypesEdgeDestination(t *testing.T) {
	assertEventStreamsSucceed(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE(), RID("a"), RID("b")},
			test.Events{E()},
			test.Events{ED()},
			test.ValidEdgeDestinations))

	assertEventStreamsFail(t,
		test.GenerateAllVariants(
			test.Events{BD(), V(ceVer)},
			test.Events{EDGE(), RID("a"), RID("b")},
			test.Events{E()},
			test.Events{ED()},
			test.InvalidEdgeDescriptions))
}

func TestMedia(t *testing.T) {
	rules := NewRules(nullevent.NewNullEventReceiver(), configuration.New())
	assertEventsSucceed(t, rules, BD(), V(0), BMEDIA("a/b"), ACL(0), ED())
}

func TestComment(t *testing.T) {
	assertEvents(t, BD(), V(ceVer), CM("a"), S("b"), ED())
	assertEvents(t, BD(), V(ceVer), CS("a"), S("b"), ED())

	assertEvents(t, BD(), V(ceVer), L(), CM("a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), L(), CM("a"), S("b"), E(), ED())
	assertEvents(t, BD(), V(ceVer), L(), S("b"), CM("a"), E(), ED())

	assertEvents(t, BD(), V(ceVer), M(), CM("a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), S("b"), S("b"), CM("a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), S("b"), CM("a"), S("b"), E(), ED())
	assertEvents(t, BD(), V(ceVer), M(), CM("a"), S("b"), S("b"), E(), ED())

	assertEvents(t, BD(), V(ceVer), NODE(), CM("a"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), NODE(), CM("a"), S("x"), CM("a"), E(), ED())
	assertEvents(t, BD(), V(ceVer), NODE(), CM("a"), S("x"), CM("a"), S("x"), S("x"), CM("a"), E(), ED())

	assertEvents(t, BD(), V(ceVer), EDGE(), CM("a"), S("x"), S("x"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), CM("a"), CM("a"), S("x"), S("x"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), CM("a"), S("x"), CM("a"), S("x"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), CM("a"), S("x"), S("x"), CM("a"), S("x"), E(), ED())
	assertEvents(t, BD(), V(ceVer), EDGE(), CM("a"), CM("a"), CM("a"), S("x"), CM("a"), CM("a"), S("x"), CM("a"), CM("a"), S("x"), E(), ED())

	// assertEvents(t, BD(), V(ceVer), AI8B(), ACM(1), AD([]byte{1}), CM("a"), ACL(1), AD([]byte{1}), ED())
}

func TestEdgeMaxDepth(t *testing.T) {
	assertEvents(t, BD(), V(ceVer),
		L(), MARK("a"), N(1), MARK("b"), N(2), MARK("c"), N(3), MARK("d"), N(4),
		EDGE(), REFL("a"), N(1), REFL("b"), E(),
		EDGE(), REFL("a"), N(1), REFL("c"), E(),
		EDGE(), REFL("b"), N(1), REFL("c"), E(),
		EDGE(), REFL("b"), N(1), REFL("d"), E(),
		EDGE(), REFL("c"), M(), E(), REFL("d"), E(),
		E(),
		ED())
}
