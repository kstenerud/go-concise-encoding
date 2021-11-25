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

import (
	"github.com/kstenerud/go-concise-encoding/internal/arrays"
)

// Copy byte slice data, interpreting as little endian data of larger types.
func BytesToInt8Slice(data []byte) []int8 { return arrays.BytesToInt8Slice(data) }

func BytesToUint16Slice(data []byte) []uint16 { return arrays.BytesToUint16Slice(data) }

func BytesToInt16Slice(data []byte) []int16 { return arrays.BytesToInt16Slice(data) }

func BytesToUint32Slice(data []byte) []uint32 { return arrays.BytesToUint32Slice(data) }

func BytesToInt32Slice(data []byte) []int32 { return arrays.BytesToInt32Slice(data) }

func BytesToFloat32Slice(data []byte) []float32 { return arrays.BytesToFloat32Slice(data) }

func BytesToUint64Slice(data []byte) []uint64 { return arrays.BytesToUint64Slice(data) }

func BytesToInt64Slice(data []byte) []int64 { return arrays.BytesToInt64Slice(data) }

func BytesToFloat64Slice(data []byte) []float64 { return arrays.BytesToFloat64Slice(data) }

// Reinterpret slices as byte slices. On little endian machines, these are zero
// cost. On big endian machines, these will allocate space and copy the data.
// Warning: Do not write to or store the resulting byte slice! The contents
//          should be considered volatile and likely to change if the source
//          slice changes.
func Int8SliceAsBytes(data []int8) []byte { return arrays.Int8SliceAsBytes(data) }

func Uint16SliceAsBytes(data []uint16) []byte { return arrays.Uint16SliceAsBytes(data) }

func Int16SliceAsBytes(data []int16) []byte { return arrays.Int16SliceAsBytes(data) }

func Uint32SliceAsBytes(data []uint32) []byte { return arrays.Uint32SliceAsBytes(data) }

func Int32SliceAsBytes(data []int32) []byte { return arrays.Int32SliceAsBytes(data) }

func Float32SliceAsBytes(data []float32) []byte { return arrays.Float32SliceAsBytes(data) }

func Uint64SliceAsBytes(data []uint64) []byte { return arrays.Uint64SliceAsBytes(data) }

func Int64SliceAsBytes(data []int64) []byte { return arrays.Int64SliceAsBytes(data) }

func Float64SliceAsBytes(data []float64) []byte { return arrays.Float64SliceAsBytes(data) }
