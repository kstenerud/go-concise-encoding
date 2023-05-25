// Copyright 2022 Karl Stenerud
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

package test

import (
	"encoding/binary"
	"math"
)

func arrayBitsToBytes(src []bool) []byte {
	dst := make([]byte, len(src)/8+1)
	if len(src)&3 == 0 {
		dst = dst[:len(dst)-1]
	}

	iDst := 0
	for iSrc := 0; iSrc < len(src); {
		bitCount := 8
		if len(src)-iSrc < 8 {
			bitCount = len(src) - iSrc
		}
		accum := byte(0)
		for iBit := 0; iBit < bitCount; iBit++ {
			if src[iSrc] {
				accum |= 1 << iBit
			}
			iSrc++
		}
		dst[iDst] = accum
		iDst++
	}

	return dst
}

func arrayInt8ToBytes(src []int8) []byte {
	bytes := make([]byte, len(src))
	for i, v := range src {
		bytes[i] = byte(v)
	}
	return bytes
}

func arrayInt16ToBytes(src []int16) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], uint16(v))
	}
	return bytes
}

func arrayInt32ToBytes(src []int32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], uint32(v))
	}
	return bytes
}

func arrayInt64ToBytes(src []int64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], uint64(v))
	}
	return bytes
}

func arrayUint8ToBytes(src []uint8) []byte {
	return src
}

func arrayUint16ToBytes(src []uint16) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], v)
	}
	return bytes
}

func arrayUint32ToBytes(src []uint32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], v)
	}
	return bytes
}

func arrayUint64ToBytes(src []uint64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], v)
	}
	return bytes
}

func arrayFloat16ToBytes(src []float32) []byte {
	const step = 2
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint16(bytes[i*step:], uint16(math.Float32bits(v)>>16))
	}
	return bytes
}

func arrayFloat32ToBytes(src []float32) []byte {
	const step = 4
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint32(bytes[i*step:], math.Float32bits(v))
	}
	return bytes
}

func arrayFloat64ToBytes(src []float64) []byte {
	const step = 8
	bytes := make([]byte, len(src)*step)
	for i, v := range src {
		binary.LittleEndian.PutUint64(bytes[i*step:], math.Float64bits(v))
	}
	return bytes
}

func arrayUIDToBytes(src [][]byte) []byte {
	const step = 16
	bytes := make([]byte, 0, len(src)*step)
	for _, v := range src {
		bytes = append(bytes, v...)
	}
	return bytes
}

func arrayTextToBytes(src string) []byte {
	return []byte(src)
}

// ===========================================================================

func bytesToArrayBits(length uint64, src []byte) []bool {
	dst := make([]bool, length)
	if length == 0 {
		return dst
	}

	iDst := 0
	iSrc := 0
	for ; iSrc < len(src)-1; iSrc++ {
		b := src[iSrc]
		for mask := 1; mask < 256; mask <<= 1 {
			dst[iDst] = b&byte(mask) != 0
			iDst++
		}
	}

	b := src[iSrc]
	length -= uint64(iDst)
	for iBit := uint64(0); iBit < length; iBit++ {
		dst[iDst] = b&byte(1<<iBit) != 0
		iDst++
	}

	return dst
}

func bytesToArrayInt8(src []byte) []int8 {
	array := make([]int8, len(src))
	for i, v := range src {
		array[i] = int8(v)
	}
	return array
}

func bytesToArrayInt16(src []byte) []int16 {
	const step = 2
	elemCount := len(src) / step
	array := make([]int16, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int16(binary.LittleEndian.Uint16(src[i*step:]))
	}
	return array
}

func bytesToArrayInt32(src []byte) []int32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]int32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int32(binary.LittleEndian.Uint32(src[i*step:]))
	}
	return array
}

func bytesToArrayInt64(src []byte) []int64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]int64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = int64(binary.LittleEndian.Uint64(src[i*step:]))
	}
	return array
}

func bytesToArrayUint16(src []byte) []uint16 {
	const step = 2
	elemCount := len(src) / step
	array := make([]uint16, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint16(src[i*step:])
	}
	return array
}

func bytesToArrayUint32(src []byte) []uint32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]uint32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint32(src[i*step:])
	}
	return array
}

func bytesToArrayUint64(src []byte) []uint64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]uint64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = binary.LittleEndian.Uint64(src[i*step:])
	}
	return array
}

func bytesToArrayFloat16(src []byte) []float32 {
	const step = 2
	elemCount := len(src) / step
	array := make([]float32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float32frombits(uint32(binary.LittleEndian.Uint16(src[i*step:])) << 16)
	}
	return array
}

func bytesToArrayFloat32(src []byte) []float32 {
	const step = 4
	elemCount := len(src) / step
	array := make([]float32, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float32frombits(binary.LittleEndian.Uint32(src[i*step:]))
	}
	return array
}

func bytesToArrayFloat64(src []byte) []float64 {
	const step = 8
	elemCount := len(src) / step
	array := make([]float64, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = math.Float64frombits(binary.LittleEndian.Uint64(src[i*step:]))
	}
	return array
}

func bytesToArrayUID(src []byte) [][]byte {
	const step = 16
	elemCount := len(src) / step
	array := make([][]byte, elemCount)
	for i := 0; i < elemCount; i++ {
		array[i] = src[i*step : (i+1)*step]
	}
	return array
}
