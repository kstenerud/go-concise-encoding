package cbe

import (
	"fmt"
	"reflect"
	"time"
)

type PrimitiveEncoder interface {
	Nil() error
	Bool(bool) error
	Uint(uint64) error
	Int(int64) error
	Float(float64) error
	Time(time.Time) error
	String(string) error
	Bytes([]byte) error
	ListBegin() error
	ListEnd() error
	MapBegin() error
	MapEnd() error
	// Comment() error
}

type UnsupportedTypeError error

func NewUnsupportedTypeError(unsupportedType reflect.Type) UnsupportedTypeError {
	return UnsupportedTypeError(fmt.Errorf("Unsupported type: %v", unsupportedType))
}

type Marshaler struct {
	encoder PrimitiveEncoder
}

func NewMarshaler(encoder PrimitiveEncoder) *Marshaler {
	marshaler := new(Marshaler)
	marshaler.encoder = encoder
	return marshaler
}

func (marshaler *Marshaler) Marshal(object interface{}) error {
	if object == nil {
		return marshaler.encoder.Nil()
	}

	rv := reflect.ValueOf(object)
	return marshaler.MarshalReflectValue(&rv)
}

func (marshaler *Marshaler) MarshalReflectValue(rv *reflect.Value) error {
	// TODO: IsNil is only for chan, func, interface, map, pointer, or slice
	// if rv.IsNil() {
	// 	return marshaler.encoder.Nil()
	// }

	switch rv.Kind() {
	// case reflect.Complex64, reflect.Complex128:
	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.UnsafePointer:
	case reflect.Bool:
		return marshaler.encoder.Bool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return marshaler.encoder.Int(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return marshaler.encoder.Uint(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return marshaler.encoder.Float(rv.Float())
	case reflect.String:
		return marshaler.encoder.String(rv.String())
	case reflect.Interface:
		v := rv.Elem()
		return marshaler.MarshalReflectValue(&v)
	case reflect.Struct:
		// TODO: anonymous structs?
		if err := marshaler.encoder.MapBegin(); err != nil {
			return err
		}
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			// TODO: tags: marshalKey, marshalShortKey? encodedKey?
			k := field.Name
			v := rv.Field(i)
			if err := marshaler.Marshal(k); err != nil {
				return err
			}
			if err := marshaler.Marshal(v); err != nil {
				return err
			}
		}
		return marshaler.encoder.MapEnd()
	case reflect.Map:
		if err := marshaler.encoder.MapBegin(); err != nil {
			return err
		}
		for iter := rv.MapRange(); iter.Next(); {
			k := iter.Key()
			v := iter.Value()
			if err := marshaler.MarshalReflectValue(&k); err != nil {
				return err
			}
			if err := marshaler.MarshalReflectValue(&v); err != nil {
				return err
			}
		}
		return marshaler.encoder.MapEnd()
	case reflect.Slice, reflect.Array:
		if rv.Elem().Kind() == reflect.Uint8 {
			return marshaler.encoder.Bytes(rv.Bytes())
		}
		if err := marshaler.encoder.ListBegin(); err != nil {
			return err
		}
		for i := 0; i < rv.Len(); i++ {
			v := rv.Index(i)
			if err := marshaler.MarshalReflectValue(&v); err != nil {
				return err
			}
		}
		return marshaler.encoder.ListEnd()
	case reflect.Ptr:
		v := rv.Elem()
		return marshaler.MarshalReflectValue(&v)
	default:
		return NewUnsupportedTypeError(rv.Type())
	}
	return nil
}
