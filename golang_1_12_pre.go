// +build !go1.12

package concise_encoding

import (
	"reflect"
)

type mapIter struct {
	mapInstance reflect.Value
	keys        []reflect.Value
	index       int
}

func (this *mapIter) Key() reflect.Value {
	return this.keys[this.index]
}

func (this *mapIter) Value() reflect.Value {
	return this.mapInstance.MapIndex(this.Key())
}

func (this *mapIter) Next() bool {
	this.index++
	return this.index < len(this.keys)
}

func mapRange(v reflect.Value) *mapIter {
	return &mapIter{
		mapInstance: v,
		keys:        v.MapKeys(),
		index:       -1,
	}
}
