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

	"github.com/kstenerud/go-duplicates"
)

func NewRootObjectIterator(eventReceiver DataEventReceiver, options *IteratorOptions) *RootObjectIterator {
	this := new(RootObjectIterator)
	this.Init(eventReceiver, options)
	return this
}

func (this *RootObjectIterator) Init(eventReceiver DataEventReceiver, options *IteratorOptions) {
	if options == nil {
		options = &IteratorOptions{}
	}
	this.options = options
	this.eventReceiver = eventReceiver
}

func (this *RootObjectIterator) Iterate(value interface{}) {
	if value == nil {
		this.eventReceiver.OnVersion(cbeCodecVersion)
		this.eventReceiver.OnNil()
		this.eventReceiver.OnEndDocument()
		return
	}
	this.findReferences(value)
	rv := reflect.ValueOf(value)
	iterator := getIteratorForType(rv.Type())
	iterator = iterator.CloneFromTemplate(this)
	// TODO: Move this somewhere else
	this.eventReceiver.OnVersion(cbeCodecVersion)
	iterator.Iterate(rv)
	this.eventReceiver.OnEndDocument()
}

// Iterates depth-first recursively through an object, notifying callbacks as it
// encounters data.
type RootObjectIterator struct {
	foundReferences map[duplicates.TypedPointer]bool
	namedReferences map[duplicates.TypedPointer]uint32
	nextMarkerName  uint32
	eventReceiver   DataEventReceiver
	options         *IteratorOptions
}

func (this *RootObjectIterator) findReferences(value interface{}) {
	if this.options.UseReferences {
		this.foundReferences = duplicates.FindDuplicatePointers(value)
		this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}
}

func (this *RootObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if this.options.UseReferences {
		ptr := duplicates.TypedPointerOfRV(v)
		if this.foundReferences[ptr] {
			var name uint32
			var exists bool
			if name, exists = this.namedReferences[ptr]; !exists {
				name = this.nextMarkerName
				this.nextMarkerName++
				this.namedReferences[ptr] = name
				this.eventReceiver.OnMarker()
				this.eventReceiver.OnPositiveInt(uint64(name))
				return false
			} else {
				this.eventReceiver.OnReference()
				this.eventReceiver.OnPositiveInt(uint64(name))
				return true
			}
		}
	}
	return false
}
