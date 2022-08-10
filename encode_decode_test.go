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

package concise_encoding

import (
	"testing"

	"github.com/kstenerud/go-concise-encoding/test"
)

func TestEncodeDecodeAllValidTLO(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{},
			test.Events{},
			test.Events{},
			test.RemoveEvents(test.ValidTLOValues, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidList(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{L()},
			test.Events{E()},
			test.Events{},
			test.RemoveEvents(test.ValidListValues, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidMapKey(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{M()},
			test.Events{B(true), E()},
			test.Events{},
			test.RemoveEvents(test.ValidMapKeys, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidMapValue(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{M(), B(true)},
			test.Events{E()},
			test.Events{},
			test.RemoveEvents(test.ValidMapValues, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidEdgeSources(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{EDGE()},
			test.Events{RID("x"), N(1), E()},
			test.Events{},
			test.RemoveEvents(test.ValidEdgeSources, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidEdgeDescriptions(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{EDGE(), RID("x")},
			test.Events{N(1), E()},
			test.Events{},
			test.RemoveEvents(test.ValidEdgeDescriptions, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}

func TestEncodeDecodeAllValidEdgeDestinations(t *testing.T) {
	assertEncodeDecodeEventStreams(t,
		test.GenerateAllVariants(
			test.Events{V(ceVer)},
			test.Events{EDGE(), RID("x"), RID("y")},
			test.Events{E()},
			test.Events{},
			test.RemoveEvents(test.ValidEdgeDestinations, append(test.ArrayBeginTypes, test.EvCTB, test.EvCUT)...)))
}
