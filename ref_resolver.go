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
	"fmt"
	"reflect"
)

// src accumulator
// dst value setter

// src is:
// not defined at all yet
// completely defined (save copy or ptr)
// being built (save ref to src builder)

// Container builders will need to know that there's a marked object being built
// On object complete, it registers the marked object
// on registering, marker registry checks for references and fills them
// create a reflect.Value style interface so that setting it sets the r.v and also fills a marker?
// This means an extra indirection for all settters. Better than if statement?
// Alternative: Fit in a wrapper builder that passes on build instr and also fills marker?
// marker, scalar: insert marker-builder-scalar in front of current builder, trap all scalar builds
// marker, container: insert marker-builder-container in front of inner builder, trap child finished.
// what about interface, where scalar or container not known yet?
// - outer marker builder can somehow notify that a container build is needed?
// - OnMarkerContainer?
/*
marker:
- insert marker name builder, which accepts string or int
- on string or int finished, remove name builder and insert marked object builder, passing in name
- on scalar, pass on the build message, store value in reference manager
- on list/map, wait for child container finished, store value in reference manager

reference:
- insert reference name builder, which accepts string or int
- on finished, either fetch value from ref manager, or ...

Reference has pointer to thing that needs to be set:
- ptr to object
- map w/ key
- ptr to array w/ index
- slice ... ? pointer to builder maybe?
-- use ptr to slice in builder?

--------------------

Or maybe have builder methods return a reflect.Value of the object just built?
- would require chain of returns. Maybe not


OnMarker:
- replace root's builder with marker ID builder

OnString:
- marker ID builder calls BuildFromMarker

BuildFromMarker:
- insert marker object builder before elem builder...
-- how to deal with struct elem builders?
-- maybe just use if statement to call marker builder chained to real builder...
-- or build parallel elem builder sets, one with marker wrapper...

*/

type MarkerRegistry struct {
	markedValues         map[interface{}]interface{}
	unresolvedReferences map[interface{}][]interface{}
}

func (_this *MarkerRegistry) HasUnresolvedReferences() bool {
	return len(_this.unresolvedReferences) > 0
}

func (_this *MarkerRegistry) NotifyMarker(id interface{}, value interface{}) {
	_this.markedValues[id] = value

	if references, ok := _this.unresolvedReferences[id]; ok {
		for _, reference := range references {
			fillReference(reference, value)
		}
		delete(_this.unresolvedReferences, id)
	}
}

func (_this *MarkerRegistry) NotifyReference(id interface{}, reference interface{}) {
	if value, ok := _this.markedValues[id]; ok {
		fillReference(reference, value)
		return
	}
	_this.unresolvedReferences[id] = append(_this.unresolvedReferences[id], reference)
}

// Maybe use ptr to reflect.value?
func fillReference(reference interface{}, value interface{}) {
	src := reflect.ValueOf(value)
	dst := reflect.ValueOf(reference)
	if dst.Type() == src.Type() {
		dst.Set(src)
	}
	panic(fmt.Errorf("TODO: Support converting %v to %v", src.Type(), dst.Type()))
}

//

// dst is:
// asking for completely built src (just set it)
// asking for unfinished src (save ref to dst builder? Special setter?)

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
