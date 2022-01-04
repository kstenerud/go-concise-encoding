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
	"math"
)

func bytesToInt8Slice(data []byte) []int8 {
	length := len(data)
	result := make([]int8, length)
	for i, v := range data {
		result[i] = int8(v)
	}
	return result
}

func bytesToUint16Slice(data []byte) []uint16 {
	length := len(data) / 2
	result := make([]uint16, length)
	for i := 0; i < length; i++ {
		result[i] = uint16(data[i*2]) |
			uint16(data[i*2+1])<<8
	}
	return result
}

func bytesToInt16Slice(data []byte) []int16 {
	length := len(data) / 2
	result := make([]int16, length)
	for i := 0; i < length; i++ {
		result[i] = int16(data[i*2]) |
			int16(data[i*2+1])<<8
	}
	return result
}

func bytesToUint32Slice(data []byte) []uint32 {
	length := len(data) / 4
	result := make([]uint32, length)
	for i := 0; i < length; i++ {
		result[i] = uint32(data[i*4]) |
			uint32(data[i*4+1])<<8 |
			uint32(data[i*4+2])<<16 |
			uint32(data[i*4+3])<<24
	}
	return result
}

func bytesToInt32Slice(data []byte) []int32 {
	length := len(data) / 4
	result := make([]int32, length)
	for i := 0; i < length; i++ {
		result[i] = int32(data[i*4]) |
			int32(data[i*4+1])<<8 |
			int32(data[i*4+2])<<16 |
			int32(data[i*4+3])<<24
	}
	return result
}

func bytesToFloat16Slice(data []byte) []float32 {
	length := len(data) / 2
	result := make([]float32, length)
	for i := 0; i < length; i++ {
		result[i] = math.Float32frombits(uint32(uint32(data[i*2])<<16 |
			uint32(data[i*2+1])<<24))
	}
	return result
}

func bytesToFloat32Slice(data []byte) []float32 {
	length := len(data) / 4
	result := make([]float32, length)
	for i := 0; i < length; i++ {
		result[i] = math.Float32frombits(uint32(data[i*4]) |
			uint32(data[i*4+1])<<8 |
			uint32(data[i*4+2])<<16 |
			uint32(data[i*4+3])<<24)
	}
	return result
}

func bytesToUint64Slice(data []byte) []uint64 {
	length := len(data) / 8
	result := make([]uint64, length)
	for i := 0; i < length; i++ {
		result[i] = uint64(data[i*8]) |
			uint64(data[i*8+1])<<8 |
			uint64(data[i*8+2])<<16 |
			uint64(data[i*8+3])<<24 |
			uint64(data[i*8+4])<<32 |
			uint64(data[i*8+5])<<40 |
			uint64(data[i*8+6])<<48 |
			uint64(data[i*8+7])<<56
	}
	return result
}

func bytesToInt64Slice(data []byte) []int64 {
	length := len(data) / 8
	result := make([]int64, length)
	for i := 0; i < length; i++ {
		result[i] = int64(data[i*8]) |
			int64(data[i*8+1])<<8 |
			int64(data[i*8+2])<<16 |
			int64(data[i*8+3])<<24 |
			int64(data[i*8+4])<<32 |
			int64(data[i*8+5])<<40 |
			int64(data[i*8+6])<<48 |
			int64(data[i*8+7])<<56
	}
	return result
}

func bytesToFloat64Slice(data []byte) []float64 {
	length := len(data) / 8
	result := make([]float64, length)
	for i := 0; i < length; i++ {
		result[i] = math.Float64frombits(uint64(data[i*8]) |
			uint64(data[i*8+1])<<8 |
			uint64(data[i*8+2])<<16 |
			uint64(data[i*8+3])<<24 |
			uint64(data[i*8+4])<<32 |
			uint64(data[i*8+5])<<40 |
			uint64(data[i*8+6])<<48 |
			uint64(data[i*8+7])<<56)
	}
	return result
}

func int8SliceToBytes(data []int8) []byte {
	length := len(data)
	result := make([]byte, length)
	for i, v := range data {
		result[i] = byte(v)
	}
	return result
}

func uint16SliceToBytes(data []uint16) []byte {
	length := len(data) * 2
	result := make([]byte, length)
	for i, v := range data {
		result[i*2] = byte(v)
		result[i*2+1] = byte(v >> 8)
	}
	return result
}

func int16SliceToBytes(data []int16) []byte {
	length := len(data) * 2
	result := make([]byte, length)
	for i, v := range data {
		result[i*2] = byte(v)
		result[i*2+1] = byte(v >> 8)
	}
	return result
}

func float16SliceToBytes(data []float32) []byte {
	length := len(data) * 2
	result := make([]byte, length)
	for i, v := range data {
		f := math.Float32bits(v)
		result[i*2] = byte(f >> 16)
		result[i*2+1] = byte(f >> 24)
	}
	return result
}

func uint32SliceToBytes(data []uint32) []byte {
	length := len(data) * 4
	result := make([]byte, length)
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
	result := make([]byte, length)
	for i, v := range data {
		result[i*4] = byte(v)
		result[i*4+1] = byte(v >> 8)
		result[i*4+2] = byte(v >> 16)
		result[i*4+3] = byte(v >> 24)
	}
	return result
}

func float32SliceToBytes(data []float32) []byte {
	length := len(data) * 4
	result := make([]byte, length)
	for i, v := range data {
		f := math.Float32bits(v)
		result[i*4] = byte(f)
		result[i*4+1] = byte(f >> 8)
		result[i*4+2] = byte(f >> 16)
		result[i*4+3] = byte(f >> 24)
	}
	return result
}

func uint64SliceToBytes(data []uint64) []byte {
	length := len(data) * 8
	result := make([]byte, length)
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
	result := make([]byte, length)
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

func float64SliceToBytes(data []float64) []byte {
	length := len(data) * 8
	result := make([]byte, length)
	for i, v := range data {
		f := math.Float64bits(v)
		result[i*8] = byte(f)
		result[i*8+1] = byte(f >> 8)
		result[i*8+2] = byte(f >> 16)
		result[i*8+3] = byte(f >> 24)
		result[i*8+4] = byte(f >> 32)
		result[i*8+5] = byte(f >> 40)
		result[i*8+6] = byte(f >> 48)
		result[i*8+7] = byte(f >> 56)
	}
	return result
}

func uuidSliceToBytes(data [][]byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	result := data[0]
	for i := 1; i < len(data); i++ {
		result = append(result, data[i]...)
	}
	return result
}

func bytesToUUIDSlice(data []byte) [][]byte {
	result := [][]byte{}
	for i := 0; i < len(data); i += 16 {
		result = append(result, data[i:i+16])
	}
	return result
}
