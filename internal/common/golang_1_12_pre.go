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

// +build !go1.12

// TODO: Can't use < 1.13 anymore due to hex float notation
package common

import (
	"reflect"
)

type mapIter struct {
	mapInstance reflect.Value
	keys        []reflect.Value
	index       int
}

func (_this *mapIter) Key() reflect.Value {
	return _this.keys[_this.index]
}

func (_this *mapIter) Value() reflect.Value {
	return _this.mapInstance.MapIndex(_this.Key())
}

func (_this *mapIter) Next() bool {
	_this.index++
	return _this.index < len(_this.keys)
}

func MapRange(v reflect.Value) *mapIter {
	return &mapIter{
		mapInstance: v,
		keys:        v.MapKeys(),
		index:       -1,
	}
}
