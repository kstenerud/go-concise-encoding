// +build go1.12

package concise_encoding

import (
	"reflect"
)

func mapRange(v reflect.Value) *reflect.MapIter {
	return v.MapRange()
}
