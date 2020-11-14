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

package cte

import (
	"math"
)

const (
	minInt8Value   = -0x80
	maxInt8Value   = 0x7f
	minInt16Value  = -0x8000
	maxInt16Value  = 0x7fff
	minInt32Value  = -0x80000000
	maxInt32Value  = 0x7fffffff
	maxUint8Value  = 0xff
	maxUint16Value = 0xffff
	maxUint32Value = 0xffffffff

	minFloat32Exponent = -126
	maxFloat32Exponent = 127

	minFloat64Exponent        = -1022
	maxFloat64Exponent        = 1023
	minFloat64DecimalExponent = -308 // Denormalized can technically go to -324
	maxFloat64DecimalExponent = 308
	maxFloat64Coefficient     = (uint64(1) << 54) - 1

	charNumericWhitespace = '_'
)

func extractFloat64Exponent(v float64) int {
	return int((math.Float64bits(v)>>52)&0x7ff) - 1023
}
