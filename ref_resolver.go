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

package concise_encoding

import (
	"reflect"
)

// How does this work for adding to an array?
// Create default value, keep reflect.Value for it?
// What about pointers?
// Slice, array: need reflect.Value of actual element so it can be modified
// Map key: can't modify key directly. Also can't store zero value key...
// Map value: keep reflect.Value as ref
// Struct key: How would this even work?
// Struct value: keep reflect.Value as ref

// value may not yet be complete when it gets referenced (recursive)
// map and struct keys will always be complete when referenced.
// Only maps, structs, slices, arrays can be incomplete when referenced
// Of these, only slices are a problem
// Even with reflect, r.v after append can be different.
// If a subobject of a slice element, only the containing slice needs to be tracked because once the object is complete ...
// map, struct: containing object, key
// array, slice: containing object, index
// Or do I only need ref to builder that holds the obj being built?

// Do I need type-specific resolvers? slice,array,map,struct
// slice -> o, o, o, ref same slice ...
// If it's not finished building, the builder will still be alive
// - so ref builder needs reference to another builder. get built object.
// - for scalars, just keep a copy of the value. Generate and run builder.

// TODO: Refine this into something real
type ReferenceResolver struct {
	markers map[interface{}][]reflect.Value
}

// TODO
func (_this *ReferenceResolver) DefineMarker(id interface{}) {
	// TODO: Does this need to check for existence? Or is that better off in rules?
	_this.markers[id] = make([]reflect.Value, 0, 2)
}

// TODO
func (_this *ReferenceResolver) AddReference(id interface{}, object reflect.Value) {
	_this.markers[id] = append(_this.markers[id], object)
}

// TODO
func (_this *ReferenceResolver) ResolveReferences() {
	panic("TODO: ResolveReferences")
}
