package cbe

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

type PrimitiveEncoder interface {
	Nil() error
	Bool(bool) error
	Uint(uint64) error
	Int(int64) error
	Float(float64) error
	Timestamp(time.Time) error
	String(string) error
	URI(*url.URL) error
	Bytes([]byte) error
	ListBegin() error
	UnorderedMapBegin() error
	ContainerEnd() error
}

type UnsupportedTypeError error

func NewUnsupportedTypeError(unsupportedType reflect.Type) UnsupportedTypeError {
	return UnsupportedTypeError(fmt.Errorf("Unsupported type: %v", unsupportedType))
}

func Marshal(encoder PrimitiveEncoder, inlineContainerType ContainerType, object interface{}) error {
	rv := reflect.ValueOf(object)
	return marshalReflectValue(encoder, inlineContainerType, &rv)
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
			return encoder.Timestamp(realValue)
		}
		if rt.Name() == "URL" && rt.PkgPath() == "net/url" {
			realValue := rv.Interface().(url.URL)
			return encoder.URI(&realValue)
		}
		if inlineContainerType != ContainerTypeUnorderedMap {
			// TODO: anonymous structs?
			if err := encoder.UnorderedMapBegin(); err != nil {
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
		if inlineContainerType != ContainerTypeUnorderedMap {
			return encoder.ContainerEnd()
		}
		return nil
	case reflect.Map:
		if inlineContainerType != ContainerTypeUnorderedMap {
			if err := encoder.UnorderedMapBegin(); err != nil {
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
		if inlineContainerType != ContainerTypeUnorderedMap {
			return encoder.ContainerEnd()
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
				return encoder.ContainerEnd()
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
			return encoder.ContainerEnd()
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
