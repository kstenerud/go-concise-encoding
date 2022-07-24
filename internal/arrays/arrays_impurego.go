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

//go:build !purego
// +build !purego

package arrays

import (
	"unsafe"

	"github.com/kstenerud/go-concise-encoding/internal/common"
)

var isLittleEndian bool

func init() {
	v := uint16(1)
	bytes := (*[common.AddressSpace]byte)(unsafe.Pointer(&v))[:2]
	if bytes[1] == 1 {
		isLittleEndian = true
	}
}

func clonedBytesPtr(bytes []byte) unsafe.Pointer {
	dst := make([]byte, len(bytes))
	copy(dst, bytes)
	return unsafe.Pointer(&dst[0])
}

func asBytes(ptr unsafe.Pointer, length int) []byte {
	return (*[common.AddressSpace]byte)(ptr)[:length]
}

func Uint8SliceAsBytes(data []uint8) []byte {
	return []byte(data)
}

func Int8SliceAsBytes(data []int8) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data))
	}
	return int8SliceToBytes(data)
}

func Uint16SliceAsBytes(data []uint16) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*2)
	}
	return uint16SliceToBytes(data)
}

func Int16SliceAsBytes(data []int16) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*2)
	}
	return int16SliceToBytes(data)
}

func Float16SliceAsBytes(data []float32) []byte {
	return float16SliceToBytes(data)
}

func Uint32SliceAsBytes(data []uint32) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*4)
	}
	return uint32SliceToBytes(data)
}

func Int32SliceAsBytes(data []int32) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*4)
	}
	return int32SliceToBytes(data)
}

func Float32SliceAsBytes(data []float32) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*4)
	}
	return float32SliceToBytes(data)
}

func Uint64SliceAsBytes(data []uint64) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*8)
	}
	return uint64SliceToBytes(data)
}

func Int64SliceAsBytes(data []int64) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*8)
	}
	return int64SliceToBytes(data)
}

func Float64SliceAsBytes(data []float64) []byte {
	if isLittleEndian {
		return asBytes(unsafe.Pointer(&data[0]), len(data)*8)
	}
	return float64SliceToBytes(data)
}

func BytesToInt8Slice(data []byte) []int8 {
	if isLittleEndian {
		return (*[common.AddressSpace]int8)(clonedBytesPtr(data))[:len(data)]
	}
	return bytesToInt8Slice(data)
}

func BytesToUint16Slice(data []byte) []uint16 {
	if isLittleEndian {
		return (*[common.AddressSpace / 2]uint16)(clonedBytesPtr(data))[:len(data)*2]
	}
	return bytesToUint16Slice(data)
}

func BytesToInt16Slice(data []byte) []int16 {
	if isLittleEndian {
		return (*[common.AddressSpace / 2]int16)(clonedBytesPtr(data))[:len(data)*2]
	}
	return bytesToInt16Slice(data)
}

func BytesToFloat16Slice(data []byte) []float32 {
	return bytesToFloat16Slice(data)
}

func BytesToUint32Slice(data []byte) []uint32 {
	if isLittleEndian {
		return (*[common.AddressSpace / 4]uint32)(clonedBytesPtr(data))[:len(data)*4]
	}
	return bytesToUint32Slice(data)
}

func BytesToInt32Slice(data []byte) []int32 {
	if isLittleEndian {
		return (*[common.AddressSpace / 4]int32)(clonedBytesPtr(data))[:len(data)*4]
	}
	return bytesToInt32Slice(data)
}

func BytesToFloat32Slice(data []byte) []float32 {
	if isLittleEndian {
		return (*[common.AddressSpace / 4]float32)(clonedBytesPtr(data))[:len(data)*4]
	}
	return bytesToFloat32Slice(data)
}

func BytesToUint64Slice(data []byte) []uint64 {
	if isLittleEndian {
		return (*[common.AddressSpace / 8]uint64)(clonedBytesPtr(data))[:len(data)*8]
	}
	return bytesToUint64Slice(data)
}

func BytesToInt64Slice(data []byte) []int64 {
	if isLittleEndian {
		return (*[common.AddressSpace / 8]int64)(clonedBytesPtr(data))[:len(data)*8]
	}
	return bytesToInt64Slice(data)
}

func BytesToFloat64Slice(data []byte) []float64 {
	if isLittleEndian {
		return (*[common.AddressSpace / 8]float64)(clonedBytesPtr(data))[:len(data)*8]
	}
	return bytesToFloat64Slice(data)
}

func UUIDSliceAsBytes(data [][]byte) []byte { return uuidSliceToBytes(data) }

func BytesToUUIDSlice(data []byte) [][]byte { return bytesToUUIDSlice(data) }
