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

package builder

import (
	"math/big"
	"reflect"

	"github.com/cockroachdb/apd/v2"
)

// The matching generated code is in generated_code.go

func (_this *intBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}

func (_this *uintBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}

func (_this *floatBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}

func (_this *bigFloatBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}

func (_this *pBigFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*big.Float)(nil)))
}

func (_this *bigDecimalFloatBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}

func (_this *pBigDecimalFloatBuilder) BuildFromNil(dst reflect.Value) {
	dst.Set(reflect.ValueOf((*apd.Decimal)(nil)))
}

func (_this *decimalFloatBuilder) BuildFromNil(_ reflect.Value) {
	_this.panicBadEvent("BuildFromNil")
}
