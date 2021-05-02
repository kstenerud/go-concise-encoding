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

func TestEncodeDecodeNA(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, NA(), TT(), ED())
}

func TestEncodeDecodeTrue(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, TT(), ED())
}

func TestEncodeDecodeFalse(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, FF(), ED())
}

func TestEncodeDecodePositiveInt(t *testing.T) {
	assertEncodeDecodeCTE(t, BD(), EvV, PI(0), ED())
	assertEncodeDecodeCTE(t, BD(), EvV, PI(1), ED())
	assertEncodeDecodeCBE(t, BD(), EvV, I(0), ED())
	assertEncodeDecodeCBE(t, BD(), EvV, I(1), ED())
	assertEncodeDecode(t, BD(), EvV, PI(104), ED())
	assertEncodeDecode(t, BD(), EvV, PI(10405), ED())
	assertEncodeDecode(t, BD(), EvV, PI(999999), ED())
	assertEncodeDecode(t, BD(), EvV, PI(7234859234423), ED())
}

func TestEncodeDecodeNegativeInt(t *testing.T) {
	assertEncodeDecodeCTE(t, BD(), EvV, NI(1), ED())
	assertEncodeDecodeCBE(t, BD(), EvV, I(-1), ED())
	assertEncodeDecode(t, BD(), EvV, NI(104), ED())
	assertEncodeDecode(t, BD(), EvV, NI(10405), ED())
	assertEncodeDecode(t, BD(), EvV, NI(999999), ED())
	assertEncodeDecode(t, BD(), EvV, NI(7234859234423), ED())
}

func TestEncodeDecodeFloat(t *testing.T) {
	// CTE will convert to decimal float
	assertEncodeDecodeCBE(t, BD(), EvV, F(1.5), ED())
	assertEncodeDecode(t, BD(), EvV, DF(test.NewDFloat("1.5")), ED())
	assertEncodeDecodeCBE(t, BD(), EvV, F(-51.455e-16), ED())
	assertEncodeDecode(t, BD(), EvV, DF(test.NewDFloat("-51.455e-16")), ED())
}

func TestEncodeDecodeNan(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, NAN(), ED())
	assertEncodeDecode(t, BD(), EvV, SNAN(), ED())
}

func TestEncodeDecodeUUID(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, UUID([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}), ED())
}

func TestEncodeDecodeTime(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, CT(test.NewDate(2000, 1, 1)), ED())

	assertEncodeDecode(t, BD(), EvV, CT(test.NewTime(1, 45, 0, 0, "")), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTime(23, 59, 59, 101000000, "")), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeLL(10, 0, 1, 930000000, 8992, 110)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeLL(10, 0, 1, 930000000, 0, 0)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeLL(10, 0, 1, 930000000, 100, 100)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeOff(10, 0, 1, 930000000, 0)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeOff(10, 0, 1, 930000000, 120)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTimeOff(10, 0, 1, 930000000, -500)), ED())

	assertEncodeDecode(t, BD(), EvV, CT(test.NewTS(2000, 1, 1, 19, 31, 44, 901554000, "")), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTS(-50000, 12, 29, 1, 1, 1, 305, "Etc/UTC")), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTSLL(2954, 8, 31, 12, 31, 15, 335523, 3154, 16004)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTSOff(2954, 8, 31, 12, 31, 15, 335523, 1000)), ED())
	assertEncodeDecode(t, BD(), EvV, CT(test.NewTSOff(2954, 8, 31, 12, 31, 15, 335523, -1000)), ED())
}

func TestEncodeDecodeBytes(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, AU8([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeCustom(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, CUB([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeRID(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, BD(), EvV, RID("http://example.com"), ED())
}

func TestEncodeDecodeString(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, BD(), EvV, S("A string"), ED())
}

func TestEncodeDecodeList(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, L(), E(), ED())
	assertEncodeDecode(t, BD(), EvV, L(), PI(1000), E(), ED())
}

func TestEncodeDecodeMap(t *testing.T) {
	assertEncodeDecode(t, BD(), EvV, M(), E(), ED())
	assertEncodeDecode(t, BD(), EvV, M(), S("a"), NI(1000), E(), ED())
	assertEncodeDecode(t, BD(), EvV, M(), S("some NA"), N(), DF(test.NewDFloat("1.1")), S("somefloat"), E(), ED())
}

func TestEncodeDecodeAllValidTLO(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV}
	suffix := []*test.TEvent{test.EvED}
	assertEncodeDecodeSetTLO(t, prefix, suffix, test.ValidTLOValues)
}

func TestEncodeDecodeAllValidList(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvL}
	suffix := []*test.TEvent{test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForContainer(test.ValidListValues))
}

func TestEncodeDecodeAllValidMapKey(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvM}
	suffix := []*test.TEvent{test.EvPI, test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForKey(test.ValidMapKeys))
}

func TestEncodeDecodeAllValidMapValue(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvM, test.EvPI}
	suffix := []*test.TEvent{test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForContainer(test.ValidMapValues))
}

func TestEncodeDecodeAllValidMarkupKey(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvMUP}
	suffix := []*test.TEvent{test.EvPI, test.EvE, test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForKey(test.ValidMapKeys))
}

func TestEncodeDecodeAllValidMarkupValue(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvMUP, test.EvPI}
	suffix := []*test.TEvent{test.EvE, test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForContainer(test.ValidMapValues))
}

func TestEncodeDecodeAllValidMarkupContents(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvMUP, test.EvE}
	suffix := []*test.TEvent{test.EvE, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForKey(test.ValidMarkupContents))
}

func TestEncodeDecodeAllValidCommentContents(t *testing.T) {
	prefix := []*test.TEvent{test.EvBD, test.EvV, test.EvCMT}
	suffix := []*test.TEvent{test.EvE, test.EvPI, test.EvED}
	assertEncodeDecodeSetContainer(t, prefix, suffix, test.FilterEventsForContainer(test.ValidCommentValues))
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
		[]byte{0x03, 0x00, 0x79, 0x84, 0x6b, 0x65, 0x79, 0x73, 0x7a, 0x97, 0x01,
			0x31, 0x83, 0x66, 0x6f, 0x6f, 0x97, 0x01, 0x32, 0x83, 0x62, 0x61,
			0x72, 0x7b, 0x98, 0x01, 0x31, 0x01, 0x98, 0x01, 0x32, 0x02, 0x7b},
		BD(), EvV,
		M(),
		S("keys"),
		L(),
		MARK("1"), S("foo"),
		MARK("2"), S("bar"),
		E(),
		REF("1"), PI(1),
		REF("2"), PI(2),
		E(),
		ED())
}
