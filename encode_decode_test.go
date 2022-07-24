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

func TestEncodeDecodeTrue(t *testing.T) {
	assertEncodeDecode(t, EvV, B(true))
}

func TestEncodeDecodeFalse(t *testing.T) {
	assertEncodeDecode(t, EvV, B(false))
}

func TestEncodeDecodePositiveInt(t *testing.T) {
	assertEncodeDecodeCTE(t, EvV, N(0))
	assertEncodeDecodeCTE(t, EvV, N(1))
	assertEncodeDecode(t, EvV, N(104))
	assertEncodeDecode(t, EvV, N(10405))
	assertEncodeDecode(t, EvV, N(999999))
	assertEncodeDecode(t, EvV, N(7234859234423))
}

func TestEncodeDecodeNegativeInt(t *testing.T) {
	assertEncodeDecodeCTE(t, EvV, N(-1))
	// assertEncodeDecode(t, EvV, N(-104))
	// assertEncodeDecode(t, EvV, N(-10405))
	// assertEncodeDecode(t, EvV, N(-999999))
	// assertEncodeDecode(t, EvV, N(-7234859234423))
}

func TestEncodeDecodeFloat(t *testing.T) {
	// CTE will convert to decimal float
	assertEncodeDecodeCBE(t, EvV, N(1.5))
	assertEncodeDecode(t, EvV, N(test.NewDFloat("1.5")))
	assertEncodeDecodeCBE(t, EvV, N(-51.455e-16))
	assertEncodeDecode(t, EvV, N(test.NewDFloat("-51.455e-16")))
}

func TestEncodeDecodeNan(t *testing.T) {
	assertEncodeDecode(t, EvV, NAN())
	assertEncodeDecode(t, EvV, SNAN())
}

func TestEncodeDecodeUID(t *testing.T) {
	assertEncodeDecode(t, EvV, UID([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}))
}

func TestEncodeDecodeTime(t *testing.T) {
	assertEncodeDecode(t, EvV, T(test.NewDate(2000, 1, 1)))

	assertEncodeDecode(t, EvV, T(test.NewTime(1, 45, 0, 0, "")))
	assertEncodeDecode(t, EvV, T(test.NewTime(23, 59, 59, 101000000, "")))
	assertEncodeDecode(t, EvV, T(test.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")))
	assertEncodeDecode(t, EvV, T(test.NewTimeLL(10, 0, 1, 930000000, 8992, 110)))
	assertEncodeDecode(t, EvV, T(test.NewTimeLL(10, 0, 1, 930000000, 0, 0)))
	assertEncodeDecode(t, EvV, T(test.NewTimeLL(10, 0, 1, 930000000, 100, 100)))
	assertEncodeDecode(t, EvV, T(test.NewTimeOff(10, 0, 1, 930000000, 0)))
	assertEncodeDecode(t, EvV, T(test.NewTimeOff(10, 0, 1, 930000000, 120)))
	assertEncodeDecode(t, EvV, T(test.NewTimeOff(10, 0, 1, 930000000, -500)))

	assertEncodeDecode(t, EvV, T(test.NewTS(2000, 1, 1, 19, 31, 44, 901554000, "")))
	assertEncodeDecode(t, EvV, T(test.NewTS(-50000, 12, 29, 1, 1, 1, 305, "Etc/UTC")))
	assertEncodeDecode(t, EvV, T(test.NewTSLL(2954, 8, 31, 12, 31, 15, 335523, 3154, 16004)))
	assertEncodeDecode(t, EvV, T(test.NewTSOff(2954, 8, 31, 12, 31, 15, 335523, 1000)))
	assertEncodeDecode(t, EvV, T(test.NewTSOff(2954, 8, 31, 12, 31, 15, 335523, -1000)))
}

func TestEncodeDecodeBytes(t *testing.T) {
	assertEncodeDecode(t, EvV, AU8([]byte{1, 2, 3, 4, 5, 6, 7}))
}

func TestEncodeDecodeCustom(t *testing.T) {
	assertEncodeDecode(t, EvV, CB([]byte{1, 2, 3, 4, 5, 6, 7}))
}

func TestEncodeDecodeRID(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, EvV, RID("http://example.com"))
}

func TestEncodeDecodeString(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, EvV, S("A string"))
}

func TestEncodeDecodeList(t *testing.T) {
	assertEncodeDecode(t, EvV, L(), E())
	assertEncodeDecode(t, EvV, L(), N(1000), E())
}

func TestEncodeDecodeMap(t *testing.T) {
	assertEncodeDecode(t, EvV, M(), E())
	assertEncodeDecode(t, EvV, M(), S("a"), N(-1000), E())
	assertEncodeDecode(t, EvV, M(), S("some NA"), NULL(), N(test.NewDFloat("1.1")), S("somefloat"), E())
}

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

func TestDecodeEncodeMapReferences(t *testing.T) {
	assertDecodeEncode(t, nil, nil, nil, nil, `c0
{
    "keys" = [
        &1:"foo"
        &2:"bar"
    ]
    $1 = 1
    $2 = 2
}`,
		[]byte{0x81, 0x00, 0x79, 0x84, 0x6b, 0x65, 0x79, 0x73, 0x7a, 0x97, 0x01,
			0x31, 0x83, 0x66, 0x6f, 0x6f, 0x97, 0x01, 0x32, 0x83, 0x62, 0x61,
			0x72, 0x7b, 0x98, 0x01, 0x31, 0x01, 0x98, 0x01, 0x32, 0x02, 0x7b},
		EvV,
		M(),
		S("keys"),
		L(),
		MARK("1"), S("foo"),
		MARK("2"), S("bar"),
		E(),
		REFL("1"), N(1),
		REFL("2"), N(2),
		E())
}
