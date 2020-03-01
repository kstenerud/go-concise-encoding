package concise_encoding

import (
	"reflect"

	"github.com/kstenerud/go-duplicates"
)

func NewRootObjectIterator(useReferences bool, eventHandler ConciseEncodingEventHandler) *RootObjectIterator {
	this := new(RootObjectIterator)
	this.Init(useReferences, eventHandler)
	return this
}

func (this *RootObjectIterator) Init(useReferences bool, eventHandler ConciseEncodingEventHandler) {
	this.useReferences = useReferences
	this.eventHandler = eventHandler
}

func (this *RootObjectIterator) Iterate(value interface{}) {
	this.findReferences(value)
	rv := reflect.ValueOf(value)
	iterator := getIteratorForType(rv.Type())
	iterator = iterator.CloneFromTemplate(this)
	// TODO: Move this somewhere else
	this.eventHandler.OnVersion(cbeCodecVersion)
	iterator.Iterate(rv)
	this.eventHandler.OnEndDocument()
}

// Iterates depth-first recursively through an object, notifying callbacks as it
// encounters data.
type RootObjectIterator struct {
	foundReferences map[duplicates.TypedPointer]bool
	namedReferences map[duplicates.TypedPointer]uint32
	nextMarkerName  uint32
	eventHandler    ConciseEncodingEventHandler
	useReferences   bool
}

func (this *RootObjectIterator) findReferences(value interface{}) {
	if this.useReferences {
		this.foundReferences = duplicates.FindDuplicatePointers(value)
		this.namedReferences = make(map[duplicates.TypedPointer]uint32)
	}
}

func (this *RootObjectIterator) addReference(v reflect.Value) (didAddReferenceObject bool) {
	if this.useReferences {
		ptr := duplicates.TypedPointerOfRV(v)
		if this.foundReferences[ptr] {
			var name uint32
			var exists bool
			if name, exists = this.namedReferences[ptr]; !exists {
				name = this.nextMarkerName
				this.nextMarkerName++
				this.namedReferences[ptr] = name
				this.eventHandler.OnMarker()
				this.eventHandler.OnPositiveInt(uint64(name))
				return false
			} else {
				this.eventHandler.OnReference()
				this.eventHandler.OnPositiveInt(uint64(name))
				return true
			}
		}
	}
	return false
}
