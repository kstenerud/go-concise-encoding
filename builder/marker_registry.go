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
)

// TODO: Also keep a map of the costs of the references (bytes, depth, etc)

type MarkerRegistry struct {
	markedValues         map[interface{}]reflect.Value
	unresolvedReferences map[interface{}][]func(reflect.Value)
}

func NewMarkerRegistry() *MarkerRegistry {
	_this := new(MarkerRegistry)
	_this.Init()
	return _this
}

func (_this *MarkerRegistry) Init() {
	_this.markedValues = make(map[interface{}]reflect.Value)
	_this.unresolvedReferences = make(map[interface{}][]func(reflect.Value))
}

func (_this *MarkerRegistry) HasUnresolvedReferences() bool {
	return len(_this.unresolvedReferences) > 0
}

func (_this *MarkerRegistry) NotifyMarker(id interface{}, value reflect.Value) {
	_this.markedValues[id] = value

	if setters, ok := _this.unresolvedReferences[id]; ok {
		for _, setter := range setters {
			setter(value)
		}
		delete(_this.unresolvedReferences, id)
	}
}

func (_this *MarkerRegistry) NotifyReference(id interface{}, setter func(value reflect.Value)) {
	if value, ok := _this.markedValues[id]; ok {
		setter(value)
		return
	}
	_this.unresolvedReferences[id] = append(_this.unresolvedReferences[id], setter)
}
