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

	"github.com/kstenerud/go-compact-time"
)

func assertEncodeDecodeCBE(t *testing.T, expected ...*tevent) {
	actual, err := cbeEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("CBE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecodeCTE(t *testing.T, expected ...*tevent) {
	actual, err := cteEncodeDecode(expected...)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("CTE: Expected %v but got %v", expected, actual)
	}
}

func assertEncodeDecode(t *testing.T, expected ...*tevent) {
	assertEncodeDecodeCBE(t, expected...)
	assertEncodeDecodeCTE(t, expected...)
}

func TestEncodeDecodeVersion(t *testing.T) {
	assertEncodeDecode(t, v(1), ed())
}

func TestEncodeDecodeNil(t *testing.T) {
	assertEncodeDecode(t, v(1), n(), ed())
}

func TestEncodeDecodeTrue(t *testing.T) {
	assertEncodeDecode(t, v(1), tt(), ed())
}

func TestEncodeDecodeFalse(t *testing.T) {
	assertEncodeDecode(t, v(1), ff(), ed())
}

func TestEncodeDecodePositiveInt(t *testing.T) {
	assertEncodeDecode(t, v(1), pi(0), ed())
	assertEncodeDecode(t, v(1), pi(1), ed())
	assertEncodeDecode(t, v(1), pi(104), ed())
	assertEncodeDecode(t, v(1), pi(10405), ed())
	assertEncodeDecode(t, v(1), pi(999999), ed())
	assertEncodeDecode(t, v(1), pi(7234859234423), ed())
}

func TestEncodeDecodeNegativeInt(t *testing.T) {
	assertEncodeDecode(t, v(1), ni(1), ed())
	assertEncodeDecode(t, v(1), ni(104), ed())
	assertEncodeDecode(t, v(1), ni(10405), ed())
	assertEncodeDecode(t, v(1), ni(999999), ed())
	assertEncodeDecode(t, v(1), ni(7234859234423), ed())
}

func TestEncodeDecodeFloat(t *testing.T) {
	// CTE will convert to decimal float
	assertEncodeDecodeCBE(t, v(1), f(1.5), ed())
	assertEncodeDecode(t, v(1), df(newDFloat("1.5")), ed())
	assertEncodeDecodeCBE(t, v(1), f(-51.455e-16), ed())
	assertEncodeDecode(t, v(1), df(newDFloat("-51.455e-16")), ed())
}

func TestEncodeDecodeNan(t *testing.T) {
	assertEncodeDecode(t, v(1), nan(), ed())
	assertEncodeDecode(t, v(1), snan(), ed())
}

func TestEncodeDecodeUUID(t *testing.T) {
	assertEncodeDecode(t, v(1), uuid([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}), ed())
}

func TestEncodeDecodeTime(t *testing.T) {
	assertEncodeDecode(t, v(1), ct(compact_time.NewDate(2000, 1, 1)), ed())

	assertEncodeDecode(t, v(1), ct(compact_time.NewTime(1, 45, 0, 0, "")), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTime(23, 59, 59, 101000000, "")), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 8992, 110)), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 0, 0)), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTimeLatLong(10, 0, 1, 930000000, 100, 100)), ed())

	assertEncodeDecode(t, v(1), ct(compact_time.NewTimestamp(2000, 1, 1, 19, 31, 44, 901554000, "")), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTimestamp(-50000, 12, 29, 1, 1, 1, 305, "Etc/UTC")), ed())
	assertEncodeDecode(t, v(1), ct(compact_time.NewTimestampLatLong(2954, 8, 31, 12, 31, 15, 335523, 3154, 16004)), ed())
}

func TestEncodeDecodeBytes(t *testing.T) {
	assertEncodeDecode(t, v(1), bin([]byte{1, 2, 3, 4, 5, 6, 7}), ed())
}

func TestEncodeDecodeCustom(t *testing.T) {
	assertEncodeDecode(t, v(1), cust([]byte{1, 2, 3, 4, 5, 6, 7}), ed())
}

func TestEncodeDecodeURI(t *testing.T) {
	// TODO: More complex
	assertEncodeDecode(t, v(1), uri("http://example.com"), ed())
}

func TestEncodeDecodeString(t *testing.T) {
	// TODO: More complex
	assertEncodeDecode(t, v(1), s("A string"), ed())
}

func TestEncodeDecodeList(t *testing.T) {
	assertEncodeDecode(t, v(1), l(), e(), ed())
	assertEncodeDecode(t, v(1), l(), pi(1), e(), ed())
}

func TestEncodeDecodeMap(t *testing.T) {
	assertEncodeDecode(t, v(1), m(), e(), ed())
	assertEncodeDecode(t, v(1), m(), s("a"), ni(1), e(), ed())
}
