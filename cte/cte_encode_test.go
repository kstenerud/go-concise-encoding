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

package cte

import (
	"testing"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-equivalence"
)

func assertCTEDecodeEncode(t *testing.T, expected string) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	events, err := cteDecodeToEvents([]byte(expected))
	if err != nil {
		t.Error(err)
		return
	}
	actual := string(cteEncodeEvents(events...))
	if !equivalence.IsEquivalent(actual, expected) {
		t.Errorf("Expected [%v] but got [%v]", expected, actual)
	}
}

func assertCTEEncode(t *testing.T, v interface{}, expected string) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	actual := string(cteEncodeValue(v))
	if !equivalence.IsEquivalent(actual, expected) {
		t.Errorf("Expected [%v] but got [%v]", expected, actual)
	}
}

// ============================================================================

// Tests

func TestMapFloatKey(t *testing.T) {
	assertCTEDecodeEncode(t, "c1 {nil=@nil 1.5=1000}")
}

func TestMarkerReference(t *testing.T) {
	assertCTEDecodeEncode(t, "c1 {first=&1:1000 second=#1}")
}

func TestDuplicateEmptySliceInSlice(t *testing.T) {
	sl := []interface{}{}
	v := []interface{}{sl, sl, sl}
	assertCTEEncode(t, v, "c1 [[] [] []]")
}

// TODO: Fails intermittently due to different k-v pair ordering
// func TestDuplicateEmptySliceInMap(t *testing.T) {
// 	sl := []interface{}{}
// 	v := map[interface{}]interface{}{
// 		"a": sl,
// 		"b": sl,
// 		"c": sl,
// 	}
// 	assertCTEEncode(t, v, "c1 {a=[] b=[] c=[]}")
// }
