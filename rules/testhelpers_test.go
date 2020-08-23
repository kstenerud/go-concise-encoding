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
	"reflect"
	"testing"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
)

const rulesCodecVersion = 1

var uint8Type = reflect.TypeOf(uint8(0))

func assertRulesOnString(t *testing.T, rules *Rules, value string) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnStringBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(value)) })
	}
}

func assertRulesAddBytes(t *testing.T, rules *Rules, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnTypedArrayBegin(events.ArrayTypeUint8) })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData(value) })
	}
}

func assertRulesAddURI(t *testing.T, rules *Rules, uri string) {
	length := len(uri)
	test.AssertNoPanic(t, func() { rules.OnURIBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData([]byte(uri)) })
	}
}

func assertRulesAddCustomBinary(t *testing.T, rules *Rules, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnCustomBinaryBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), false) })
	if length > 0 {
		test.AssertNoPanic(t, func() { rules.OnArrayData(value) })
	}
}

func assertRulesAddCustomText(t *testing.T, rules *Rules, value []byte) {
	length := len(value)
	test.AssertNoPanic(t, func() { rules.OnCustomTextBegin() })
	test.AssertNoPanic(t, func() { rules.OnArrayChunk(uint64(length), false) })
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

func newRulesAfterVersion(opts *options.RuleOptions) *Rules {
	rules := NewRules(events.NewNullEventReceiver(), opts)
	rules.OnBeginDocument()
	rules.OnVersion(opts.ConciseEncodingVersion)
	return rules
}

func newRulesWithMaxDepth(maxDepth int) *Rules {
	opts := options.DefaultRuleOptions()
	opts.MaxContainerDepth = uint64(maxDepth)
	return newRulesAfterVersion(opts)
}
