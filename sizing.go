package cbe

import (
	"reflect"
)

func arrayLengthFieldSize(length int) int {
	size := 1
	length >>= 7
	for length > 0 {
		length >>= 7
		size++
	}
	return size
}

func paddingSize(byteCount int) int {
	return byteCount
}

func nilSize() int {
	return 1
}

func boolSize() int {
	return 1
}

func uintSize(value uint64) int {
	switch {
	case uintFitsInSmallint(value):
		return 1
	case fitsInUint8(value):
		return 2
	case fitsInUint16(value):
		return 3
	case fitsInUint32(value):
		return 5
	default:
		return 9
	}
}

func intSize(value int64) int {
	uvalue := uint64(-value)

	switch {
	case intFitsInSmallint(value):
		return 1
	case value >= 0:
		return uintSize(uint64(value))
	case fitsInUint8(uvalue):
		return 2
	case fitsInUint16(uvalue):
		return 3
	case fitsInUint32(uvalue):
		return 5
	default:
		return 9
	}
}

func floatSize(value float64) int {
	asFloat32 := float64(float32(value))
	if value == asFloat32 {
		return 5
	}
	return 9
}

func timeSize() int {
	return 9
}

func listBeginSize() int {
	return 1
}

func listEndSize() int {
	return 1
}

func mapBeginSize() int {
	return 1
}

func mapEndSize() int {
	return 1
}

func bytesSize(value []byte) int {
	return 1 + arrayLengthFieldSize(len(value)) + len(value)
}

func stringSize(value string) int {
	fieldSize := 1
	valueLength := len(value)
	if valueLength > 15 {
		fieldSize = 1 + arrayLengthFieldSize(valueLength)
	}
	return fieldSize + valueLength
}

func commentSize(value []byte) int {
	return 1 + arrayLengthFieldSize(len(value)) + len(value)
}

func reflectValueSize(inlineContainerType ContainerType, rv *reflect.Value) int {
	if !rv.IsValid() {
		return nilSize()
	}

	switch rv.Kind() {
	// case reflect.Complex64, reflect.Complex128:
	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.UnsafePointer:
	case reflect.Bool:
		return boolSize()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intSize(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintSize(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return floatSize(rv.Float())
	case reflect.String:
		return stringSize(rv.String())
	case reflect.Interface:
		v := rv.Elem()
		return reflectValueSize(inlineContainerType, &v)
	case reflect.Struct:
		rt := rv.Type()
		if rt.Name() == "Time" && rt.PkgPath() == "time" {
			return timeSize()
		}
		var size int
		if inlineContainerType != ContainerTypeUnorderedMap {
			size += mapBeginSize()
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			// TODO: tags: marshalKey, marshalShortKey? encodedKey?
			k := field.Name
			v := rv.Field(i)
			if v.CanInterface() {
				size += EncodedSize(ContainerTypeNone, k)
				size += reflectValueSize(ContainerTypeNone, &v)
			}
		}
		if inlineContainerType != ContainerTypeUnorderedMap {
			size += mapEndSize()
		}
		return size
	case reflect.Map:
		var size int
		if inlineContainerType != ContainerTypeUnorderedMap {
			size += mapBeginSize()
		}
		for iter := rv.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			size += reflectValueSize(ContainerTypeNone, &k)
			size += reflectValueSize(ContainerTypeNone, &v)
		}
		if inlineContainerType != ContainerTypeUnorderedMap {
			size += mapEndSize()
		}
		return size
	case reflect.Array:
		if rv.CanAddr() {
			v := rv.Slice(0, rv.Len())
			return reflectValueSize(inlineContainerType, &v)
		} else if rv.Type().Elem().Kind() == reflect.Uint8 {
			// TODO: Is there a better way to do this?
			tempSlice := make([]byte, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				tempSlice[i] = rv.Index(i).Interface().(uint8)
			}
			return EncodedSize(inlineContainerType, tempSlice)
		} else {
			var size int
			if inlineContainerType != ContainerTypeList {
				size += listBeginSize()
			}
			for i := 0; i < rv.Len(); i++ {
				v := rv.Index(i)
				size += reflectValueSize(ContainerTypeNone, &v)
			}
			if inlineContainerType != ContainerTypeList {
				size += listEndSize()
			}
			return size
		}
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return bytesSize(rv.Bytes())
		}
		var size int
		if inlineContainerType != ContainerTypeList {
			size += listBeginSize()
		}
		for i := 0; i < rv.Len(); i++ {
			v := rv.Index(i)
			size += reflectValueSize(inlineContainerType, &v)
		}
		if inlineContainerType != ContainerTypeList {
			size += listEndSize()
		}
		return size
	case reflect.Ptr:
		v := rv.Elem()
		return reflectValueSize(inlineContainerType, &v)
	default:
		return 0
	}
	return 0
}

func EncodedSize(inlineContainerType ContainerType, object interface{}) int {
	rv := reflect.ValueOf(object)
	return reflectValueSize(inlineContainerType, &rv)
}
