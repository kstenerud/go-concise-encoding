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

// +build purego

package arrays

func Int8SliceAsBytes(data []int8) []byte { return int8SliceToBytes(data) }

func Uint16SliceAsBytes(data []uint16) []byte { return uint16SliceToBytes(data) }

func Int16SliceAsBytes(data []int16) []byte { return int16SliceToBytes(data) }

func Uint32SliceAsBytes(data []uint32) []byte { return uint32SliceToBytes(data) }

func Int32SliceAsBytes(data []int32) []byte { return int32SliceToBytes(data) }

func Float32SliceAsBytes(data []float32) []byte { return float32SliceToBytes(data) }

func Uint64SliceAsBytes(data []uint64) []byte { return uint64SliceToBytes(data) }

func Int64SliceAsBytes(data []int64) []byte { return int64SliceToBytes(data) }

func Float64SliceAsBytes(data []float64) []byte { return float64SliceToBytes(data) }

func BytesToInt8Slice(data []byte) []int8 { return bytesToInt8Slice(data) }

func BytesToUint16Slice(data []byte) []uint16 { return bytesToUint16Slice(data) }

func BytesToInt16Slice(data []byte) []int16 { return bytesToInt16Slice(data) }

func BytesToUint32Slice(data []byte) []uint32 { return bytesToUint32Slice(data) }

func BytesToInt32Slice(data []byte) []int32 { return bytesToInt32Slice(data) }

func BytesToFloat32Slice(data []byte) []float32 { return bytesToFloat32Slice(data) }

func BytesToUint64Slice(data []byte) []uint64 { return bytesToUint64Slice(data) }

func BytesToInt64Slice(data []byte) []int64 { return bytesToInt64Slice(data) }

func BytesToFloat64Slice(data []byte) []float64 { return bytesToFloat64Slice(data) }

func UUIDSliceAsBytes(data [][]byte) []byte { return uuidSliceToBytes(data) }

func BytesToUUIDSlice(data []byte) [][]byte { return bytesToUUIDSlice(data) }
