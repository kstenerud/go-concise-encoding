package cbe

import (
	"math"
	"reflect"
)

func asInt(value reflect.Value) (result int64, ok bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue := value.Uint()
		if uintValue <= math.MaxInt64 {
			return int64(uintValue), true
		}
	case reflect.Float32, reflect.Float64:
		floatValue := value.Float()
		intValue := int64(floatValue)
		if float64(intValue) == floatValue {
			return intValue, true
		}
	}
	return 0, false
}

func asUint(value reflect.Value) (result uint64, ok bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue := value.Int()
		if intValue >= 0 {
			return uint64(intValue), true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint(), true
	case reflect.Float32, reflect.Float64:
		floatValue := value.Float()
		uintValue := uint64(floatValue)
		if float64(uintValue) == floatValue {
			return uintValue, true
		}
	}
	return 0, false
}

func asFloat(value reflect.Value) (result float64, ok bool) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue := value.Int()
		floatValue := float64(intValue)
		if int64(floatValue) == intValue {
			return floatValue, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue := value.Uint()
		floatValue := float64(uintValue)
		if uint64(floatValue) == uintValue {
			return floatValue, true
		}
	case reflect.Float32, reflect.Float64:
		return value.Float(), true
	}
	return 0, false
}

func deepEquivalence(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Interface, reflect.Ptr:
		return deepEquivalence(a.Elem(), b.Elem())
	case reflect.Bool:
		return a.Bool() == b.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bValue, ok := asInt(b)
		return ok && a.Int() == bValue
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bValue, ok := asUint(b)
		return ok && a.Uint() == bValue
	case reflect.Float32, reflect.Float64:
		bValue, ok := asFloat(b)
		return ok && a.Float() == bValue
	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()
	case reflect.String:
		return a.String() == b.String()
	case reflect.Uintptr:
		return a.Pointer() == b.Pointer()
	case reflect.UnsafePointer:
		return a.UnsafeAddr() == b.UnsafeAddr()
	case reflect.Array, reflect.Slice:
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !deepEquivalence(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}
		iter := a.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if !deepEquivalence(v, b.MapIndex(k)) {
				// TODO: Search equivalent type keys
				return false
			}
		}
		return true
	case reflect.Struct:
		if a.NumField() != b.NumField() {
			return false
		}
		for i := 0; i < a.NumField(); i++ {
			return deepEquivalence(a.Field(i), b.Field(i))
		}
	// case reflect.Chan:
	// case reflect.Func:
	default:
		return false
	}
	return false
}

// Test if two objects are equivalent.
// Equivalence means that they are either equal, or one can be converted to the
// other's type without loss and be equal.
// For arrays, maps, and structs, it will compare elements.
func DeepEquivalence(a, b interface{}) (equal bool) {
	defer func() {
		if r := recover(); r != nil {
			equal = false
		}
	}()
	if a == nil && b == nil {
		return true
	}
	return deepEquivalence(reflect.ValueOf(a), reflect.ValueOf(b))
}
