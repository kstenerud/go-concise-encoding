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

// TODO: Also keep a map of the costs of the references (bytes, depth, etc) (in rules)

// ReferenceFiller tracks markers and references encountered in a document,
// filling out references when their corresponding markers are found.
type ReferenceFiller struct {
	markedValues         map[interface{}]reflect.Value
	unresolvedReferences map[interface{}][]func(reflect.Value)
}

// Create and initialize a new ReferenceFiller
func NewReferenceFiller() *ReferenceFiller {
	_this := new(ReferenceFiller)
	_this.Init()
	return _this
}

// Initialize an existing ReferenceFiller
func (_this *ReferenceFiller) Init() {
	_this.markedValues = make(map[interface{}]reflect.Value)
	_this.unresolvedReferences = make(map[interface{}][]func(reflect.Value))
}

// Notify that a new marker has been found.
func (_this *ReferenceFiller) NotifyMarker(id interface{}, value reflect.Value) {
	_this.markedValues[id] = value
	if setters, ok := _this.unresolvedReferences[id]; ok {
		for _, setter := range setters {
			setter(value)
		}
		delete(_this.unresolvedReferences, id)
	}
}

// Notify that a new reference has been found. valueSetter will be called when
// the marker with lookingForID is found (possibly immediately if it already has
// been found).
func (_this *ReferenceFiller) NotifyReference(lookingForID interface{}, valueSetter func(value reflect.Value)) {
	if value, ok := _this.markedValues[lookingForID]; ok {
		valueSetter(value)
		return
	}
	_this.unresolvedReferences[lookingForID] = append(_this.unresolvedReferences[lookingForID], valueSetter)
}
