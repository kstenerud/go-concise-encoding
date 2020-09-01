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

package ce

func int8SliceToBytes(data []int8) []byte {
	length := len(data)
	result := make([]byte, length, length)
	for i, v := range data {
		result[i] = byte(v)
	}
	return result
}

func uint16SliceToBytes(data []uint16) []byte {
	length := len(data) * 2
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*2] = byte(v)
		result[i*2+1] = byte(v >> 8)
	}
	return result
}

func int16SliceToBytes(data []int16) []byte {
	length := len(data) * 2
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*2] = byte(v)
		result[i*2+1] = byte(v >> 8)
	}
	return result
}

func uint32SliceToBytes(data []uint32) []byte {
	length := len(data) * 4
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*4] = byte(v)
		result[i*4+1] = byte(v >> 8)
		result[i*4+2] = byte(v >> 16)
		result[i*4+3] = byte(v >> 24)
	}
	return result
}

func int32SliceToBytes(data []int32) []byte {
	length := len(data) * 4
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*4] = byte(v)
		result[i*4+1] = byte(v >> 8)
		result[i*4+2] = byte(v >> 16)
		result[i*4+3] = byte(v >> 24)
	}
	return result
}

func uint64SliceToBytes(data []uint64) []byte {
	length := len(data) * 8
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*8] = byte(v)
		result[i*8+1] = byte(v >> 8)
		result[i*8+2] = byte(v >> 16)
		result[i*8+3] = byte(v >> 24)
		result[i*8+4] = byte(v >> 32)
		result[i*8+5] = byte(v >> 40)
		result[i*8+6] = byte(v >> 48)
		result[i*8+7] = byte(v >> 56)
	}
	return result
}

func int64SliceToBytes(data []int64) []byte {
	length := len(data) * 8
	result := make([]byte, length, length)
	for i, v := range data {
		result[i*8] = byte(v)
		result[i*8+1] = byte(v >> 8)
		result[i*8+2] = byte(v >> 16)
		result[i*8+3] = byte(v >> 24)
		result[i*8+4] = byte(v >> 32)
		result[i*8+5] = byte(v >> 40)
		result[i*8+6] = byte(v >> 48)
		result[i*8+7] = byte(v >> 56)
	}
	return result
}
