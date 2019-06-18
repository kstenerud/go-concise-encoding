package cbe

import (
	"fmt"
	"reflect"
	"time"
)

type PrimitiveEncoder interface {
	Nil() error
	NilSize() int
	Bool(bool) error
	BoolSize(bool) int
	Uint(uint64) error
	UintSize(uint64) int
	Int(int64) error
	IntSize(int64) int
	Float(float64) error
	FloatSize(float64) int
	Time(time.Time) error
	TimeSize(time.Time) int
	String(string) error
	StringSize(string) int
	Bytes([]byte) error
	BytesSize([]byte) int
	ListBegin() error
	ListBeginSize() int
	ListEnd() error
	ListEndSize() int
	MapBegin() error
	MapBeginSize() int
	MapEnd() error
	MapEndSize() int
}

type UnsupportedTypeError error

func NewUnsupportedTypeError(unsupportedType reflect.Type) UnsupportedTypeError {
	return UnsupportedTypeError(fmt.Errorf("Unsupported type: %v", unsupportedType))
}

func Marshal(encoder PrimitiveEncoder, inlineContainerType ContainerType, object interface{}) error {
	rv := reflect.ValueOf(object)
	return marshalReflectValue(encoder, inlineContainerType, &rv)
}

func MarshalSize(encoder PrimitiveEncoder, inlineContainerType ContainerType, object interface{}) int {
	rv := reflect.ValueOf(object)
	return reflectValueSize(encoder, inlineContainerType, &rv)
}

func marshalReflectValue(encoder PrimitiveEncoder, inlineContainerType ContainerType, rv *reflect.Value) error {
	if !rv.IsValid() {
		return encoder.Nil()
	}

	switch rv.Kind() {
	// case reflect.Complex64, reflect.Complex128:
	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.UnsafePointer:
	case reflect.Bool:
		return encoder.Bool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encoder.Int(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encoder.Uint(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return encoder.Float(rv.Float())
	case reflect.String:
		return encoder.String(rv.String())
	case reflect.Interface:
		v := rv.Elem()
		return marshalReflectValue(encoder, inlineContainerType, &v)
	case reflect.Struct:
		rt := rv.Type()
		if rt.Name() == "Time" && rt.PkgPath() == "time" {
			realValue := rv.Interface().(time.Time)
			return encoder.Time(realValue)
		}
		if inlineContainerType != ContainerTypeMap {
			// TODO: anonymous structs?
			if err := encoder.MapBegin(); err != nil {
				return err
			}
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			// TODO: tags: marshalKey, marshalShortKey? encodedKey?
			k := field.Name
			v := rv.Field(i)
			if v.CanInterface() {
				if err := Marshal(encoder, ContainerTypeNone, k); err != nil {
					return err
				}
				if err := marshalReflectValue(encoder, ContainerTypeNone, &v); err != nil {
					return err
				}
			}
		}
		if inlineContainerType != ContainerTypeMap {
			return encoder.MapEnd()
		}
		return nil
	case reflect.Map:
		if inlineContainerType != ContainerTypeMap {
			if err := encoder.MapBegin(); err != nil {
				return err
			}
		}
		for iter := rv.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			if err := marshalReflectValue(encoder, ContainerTypeNone, &k); err != nil {
				return err
			}
			if err := marshalReflectValue(encoder, ContainerTypeNone, &v); err != nil {
				return err
			}
		}
		if inlineContainerType != ContainerTypeMap {
			return encoder.MapEnd()
		}
		return nil
	case reflect.Array:
		if rv.CanAddr() {
			v := rv.Slice(0, rv.Len())
			return marshalReflectValue(encoder, inlineContainerType, &v)
		} else if rv.Type().Elem().Kind() == reflect.Uint8 {
			// TODO: Is there a better way to do this?
			tempSlice := make([]byte, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				tempSlice[i] = rv.Index(i).Interface().(uint8)
			}
			return Marshal(encoder, inlineContainerType, tempSlice)
		} else {
			if inlineContainerType != ContainerTypeList {
				if err := encoder.ListBegin(); err != nil {
					return err
				}
			}
			for i := 0; i < rv.Len(); i++ {
				v := rv.Index(i)
				if err := marshalReflectValue(encoder, ContainerTypeNone, &v); err != nil {
					return err
				}
			}
			if inlineContainerType != ContainerTypeList {
				return encoder.ListEnd()
			}
			return nil
		}
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return encoder.Bytes(rv.Bytes())
		}
		if inlineContainerType != ContainerTypeList {
			if err := encoder.ListBegin(); err != nil {
				return err
			}
		}
		for i := 0; i < rv.Len(); i++ {
			v := rv.Index(i)
			if err := marshalReflectValue(encoder, inlineContainerType, &v); err != nil {
				return err
			}
		}
		if inlineContainerType != ContainerTypeList {
			return encoder.ListEnd()
		}
		return nil
	case reflect.Ptr:
		v := rv.Elem()
		return marshalReflectValue(encoder, inlineContainerType, &v)
	default:
		return NewUnsupportedTypeError(rv.Type())
	}
	return nil
}

func reflectValueSize(encoder PrimitiveEncoder, inlineContainerType ContainerType, rv *reflect.Value) int {
	if !rv.IsValid() {
		return encoder.NilSize()
	}

	switch rv.Kind() {
	// case reflect.Complex64, reflect.Complex128:
	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.UnsafePointer:
	case reflect.Bool:
		return encoder.BoolSize(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encoder.IntSize(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return encoder.UintSize(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return encoder.FloatSize(rv.Float())
	case reflect.String:
		return encoder.StringSize(rv.String())
	case reflect.Interface:
		v := rv.Elem()
		return reflectValueSize(encoder, inlineContainerType, &v)
	case reflect.Struct:
		rt := rv.Type()
		if rt.Name() == "Time" && rt.PkgPath() == "time" {
			realValue := rv.Interface().(time.Time)
			return encoder.TimeSize(realValue)
		}
		var size int
		if inlineContainerType != ContainerTypeMap {
			size += encoder.MapBeginSize()
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			// TODO: tags: marshalKey, marshalShortKey? encodedKey?
			k := field.Name
			v := rv.Field(i)
			if v.CanInterface() {
				size += MarshalSize(encoder, ContainerTypeNone, k)
				size += reflectValueSize(encoder, ContainerTypeNone, &v)
			}
		}
		if inlineContainerType != ContainerTypeMap {
			size += encoder.MapEndSize()
		}
		return size
	case reflect.Map:
		var size int
		if inlineContainerType != ContainerTypeMap {
			size += encoder.MapBeginSize()
		}
		for iter := rv.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			size += reflectValueSize(encoder, ContainerTypeNone, &k)
			size += reflectValueSize(encoder, ContainerTypeNone, &v)
		}
		if inlineContainerType != ContainerTypeMap {
			size += encoder.MapEndSize()
		}
		return size
	case reflect.Array:
		if rv.CanAddr() {
			v := rv.Slice(0, rv.Len())
			return reflectValueSize(encoder, inlineContainerType, &v)
		} else if rv.Type().Elem().Kind() == reflect.Uint8 {
			// TODO: Is there a better way to do this?
			tempSlice := make([]byte, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				tempSlice[i] = rv.Index(i).Interface().(uint8)
			}
			return MarshalSize(encoder, inlineContainerType, tempSlice)
		} else {
			var size int
			if inlineContainerType != ContainerTypeList {
				size += encoder.ListBeginSize()
			}
			for i := 0; i < rv.Len(); i++ {
				v := rv.Index(i)
				size += reflectValueSize(encoder, ContainerTypeNone, &v)
			}
			if inlineContainerType != ContainerTypeList {
				size += encoder.ListEndSize()
			}
			return size
		}
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return encoder.BytesSize(rv.Bytes())
		}
		var size int
		if inlineContainerType != ContainerTypeList {
			size += encoder.ListBeginSize()
		}
		for i := 0; i < rv.Len(); i++ {
			v := rv.Index(i)
			size += reflectValueSize(encoder, inlineContainerType, &v)
		}
		if inlineContainerType != ContainerTypeList {
			size += encoder.ListEndSize()
		}
		return size
	case reflect.Ptr:
		v := rv.Elem()
		return reflectValueSize(encoder, inlineContainerType, &v)
	default:
		return 0
	}
	return 0
}
