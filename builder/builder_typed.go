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
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
)

// The matching generated code is in builder_generated.go

func (_this *uint8ArrayBuilder) BuildFromTypedArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	switch arrayType {
	case events.ArrayTypeUint8:
		// TODO: Is there a more efficient way?
		for i := 0; i < len(value); i++ {
			elem := dst.Index(i)
			elem.SetUint(uint64(value[i]))
		}
	default:
		_this.badEvent("BuildFromTypedArray(%v)", arrayType)
	}
}

func (_this *uint16ArrayBuilder) BuildFromTypedArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	panic("TODO")
}
