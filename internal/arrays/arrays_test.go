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

package arrays

import (
	"reflect"
	"testing"
)

func TestInt8SliceToBytes(t *testing.T) {
	original := []int8{1, -1, 10, -10}
	expected := []byte{0x01, 0xff, 0x0a, 0xf6}

	actual := Int8SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToInt8Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestUint16SliceToBytes(t *testing.T) {
	original := []uint16{0x1234, 0x5678}
	expected := []byte{0x34, 0x12, 0x78, 0x56}

	actual := Uint16SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToUint16Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestInt16SliceToBytes(t *testing.T) {
	original := []int16{0x1234, -0x5678}
	expected := []byte{0x34, 0x12, 0x88, 0xa9}

	actual := Int16SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToInt16Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestUint32SliceToBytes(t *testing.T) {
	original := []uint32{0x12345678, 0x9abcdef0}
	expected := []byte{0x78, 0x56, 0x34, 0x12, 0xf0, 0xde, 0xbc, 0x9a}

	actual := Uint32SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToUint32Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestInt32SliceToBytes(t *testing.T) {
	original := []int32{0x12345678, -0x1abcdef0}
	expected := []byte{0x78, 0x56, 0x34, 0x12, 0x10, 0x21, 0x43, 0xe5}

	actual := Int32SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToInt32Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestFloat32SliceToBytes(t *testing.T) {
	original := []float32{0x1.234p-50, -0x9.3p20}
	expected := []byte{0x00, 0xa0, 0x91, 0x26, 0x00, 0x00, 0x13, 0xcb}

	actual := Float32SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToFloat32Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestUint64SliceToBytes(t *testing.T) {
	original := []uint64{0x123456789abcdef0, 0x0fedcba987654321}
	expected := []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12, 0x21, 0x43, 0x65, 0x87, 0xa9, 0xcb, 0xed, 0x0f}

	actual := Uint64SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToUint64Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestInt64SliceToBytes(t *testing.T) {
	original := []int64{0x123456789abcdef0, -0x0fedcba987654321}
	expected := []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12, 0xdf, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12, 0xf0}

	actual := Int64SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToInt64Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}

func TestFloat64SliceToBytes(t *testing.T) {
	original := []float64{0x123456789abcdef0, -0x0fedcba987654321}
	expected := []byte{0xdf, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0xb2, 0x43, 0x86, 0xca, 0x0e, 0x53, 0x97, 0xdb, 0xaf, 0xc3}

	actual := Float64SliceAsBytes(original)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual2 := BytesToFloat64Slice(expected)
	if !reflect.DeepEqual(original, actual2) {
		t.Errorf("Expected %v but got %v", original, actual2)
	}
}
