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
	"reflect"
	"testing"

	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/kstenerud/go-compact-time"
)

func assertEncodeDecodeCBE(t *testing.T, expected ...*test.TEvent) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() {
		debug.DebugOptions.PassThroughPanics = false
	}()
	actual, err := cbeEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("CBE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecodeCTE(t *testing.T, expected ...*test.TEvent) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	actual, err := cteEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("CTE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecode(t *testing.T, expected ...*test.TEvent) {
	assertEncodeDecodeCBE(t, expected...)
	assertEncodeDecodeCTE(t, expected...)
}

// ============================================================================

// Tests

func TestEncodeDecodeVersion(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), ED())
}

func TestEncodeDecodeNil(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), N(), ED())
}

func TestEncodeDecodeTrue(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), TT(), ED())
}

func TestEncodeDecodeFalse(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), FF(), ED())
}

func TestEncodeDecodePositiveInt(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), PI(0), ED())
	assertEncodeDecode(t, BD(), V(1), PI(1), ED())
	assertEncodeDecode(t, BD(), V(1), PI(104), ED())
	assertEncodeDecode(t, BD(), V(1), PI(10405), ED())
	assertEncodeDecode(t, BD(), V(1), PI(999999), ED())
	assertEncodeDecode(t, BD(), V(1), PI(7234859234423), ED())
}

func TestEncodeDecodeNegativeInt(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), NI(1), ED())
	assertEncodeDecode(t, BD(), V(1), NI(104), ED())
	assertEncodeDecode(t, BD(), V(1), NI(10405), ED())
	assertEncodeDecode(t, BD(), V(1), NI(999999), ED())
	assertEncodeDecode(t, BD(), V(1), NI(7234859234423), ED())
}

func TestEncodeDecodeFloat(t *testing.T) {
	// CTE will convert to decimal float
	assertEncodeDecodeCBE(t, BD(), V(1), F(1.5), ED())
	assertEncodeDecode(t, BD(), V(1), DF(test.NewDFloat("1.5")), ED())
	assertEncodeDecodeCBE(t, BD(), V(1), F(-51.455e-16), ED())
	assertEncodeDecode(t, BD(), V(1), DF(test.NewDFloat("-51.455e-16")), ED())
}

func TestEncodeDecodeNan(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), NAN(), ED())
	assertEncodeDecode(t, BD(), V(1), SNAN(), ED())
}

func TestEncodeDecodeUUID(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), UUID([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}), ED())
}

func TestEncodeDecodeTime(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewDate(2000, 1, 1)), ED())

	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTime(1, 45, 0, 0, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTime(23, 59, 59, 101000000, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ED())

	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimestamp(-50000, 12, 29, 1, 1, 1, 305, "Etc/UTC")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(compact_time.NewTimestampLatLong(2954, 8, 31, 12, 31, 15, 335523, 3154, 16004)), ED())
}

func TestEncodeDecodeBytes(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), BIN([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeCustom(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), CUB([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeURI(t *testing.T) {
	// TODO: More complex
	assertEncodeDecode(t, BD(), V(1), URI("http://example.com"), ED())
}

func TestEncodeDecodeString(t *testing.T) {
	// TODO: More complex
	assertEncodeDecode(t, BD(), V(1), S("A string"), ED())
}

func TestEncodeDecodeList(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), L(), E(), ED())
	assertEncodeDecode(t, BD(), V(1), L(), PI(1), E(), ED())
}

func TestEncodeDecodeMap(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), M(), E(), ED())
	assertEncodeDecode(t, BD(), V(1), M(), S("a"), NI(1), E(), ED())
	assertEncodeDecode(t, BD(), V(1), M(), S("some nil"), N(), DF(test.NewDFloat("1.1")), S("somefloat"), E(), ED())
}
