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
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-concise-encoding/cte"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/debug"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// Demonstration of the Concise Encoding "custom" data type.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom

// ============================================================================

// Assume the following custom serialized formats for complex64 and 128:
//
// | Offset | Size | Description                                             |
// | ------ | ---- | ------------------------------------------------------- |
// |      0 |    1 | Data type code (so we can have multiple custom types)   |
// |      1 |    4 | Real portion (float32, little endian)                   |
// |      5 |    4 | Imaginary portion (float32, little endian)              |
//
// | Offset | Size | Description                                             |
// | ------ | ---- | ------------------------------------------------------- |
// |      0 |    1 | Data type code (so we can have multiple custom types)   |
// |      1 |    8 | Real portion (float64, little endian)                   |
// |      9 |    8 | Imaginary portion (float64, little endian)              |

// Simple type mechanism: The first byte of the data is the type field
const (
	typeCodeComplex64  = 0
	typeCodeComplex128 = 1
)

// First piece: functions to convert from complex type to custom bytes.
// These functions each handle a single type only.
func convertComplex64ToCustom(rv reflect.Value) (asBytes []byte, err error) {
	cplx := complex64(rv.Complex())

	buff := bytes.Buffer{}

	buff.WriteByte(typeCodeComplex64)
	if err = binary.Write(&buff, binary.LittleEndian, real(cplx)); err != nil {
		return
	}
	if err = binary.Write(&buff, binary.LittleEndian, imag(cplx)); err != nil {
		return
	}

	asBytes = buff.Bytes()
	return
}

func convertComplex128ToCustom(rv reflect.Value) (asBytes []byte, err error) {
	cplx := rv.Complex()

	buff := bytes.Buffer{}

	buff.WriteByte(typeCodeComplex128)
	if err = binary.Write(&buff, binary.LittleEndian, real(cplx)); err != nil {
		return
	}
	if err = binary.Write(&buff, binary.LittleEndian, imag(cplx)); err != nil {
		return
	}

	asBytes = buff.Bytes()
	return
}

// Second piece: converter function to fill in an object from custom data.
// This same function will be used for ALL custom types.
func convertFromCustom(src []byte, dst reflect.Value) error {
	buff := bytes.NewBuffer(src)

	customType, _ := buff.ReadByte()
	switch customType {
	case typeCodeComplex64:
		var realPart float32
		var imagPart float32
		if err := binary.Read(buff, binary.LittleEndian, &realPart); err != nil {
			return err
		}
		if err := binary.Read(buff, binary.LittleEndian, &imagPart); err != nil {
			return err
		}
		dst.SetComplex(complex128(complex(realPart, imagPart)))
		return nil
	case typeCodeComplex128:
		var realPart float64
		var imagPart float64
		if err := binary.Read(buff, binary.LittleEndian, &realPart); err != nil {
			return err
		}
		if err := binary.Read(buff, binary.LittleEndian, &imagPart); err != nil {
			return err
		}
		dst.SetComplex(complex(realPart, imagPart))
		return nil
	default:
		return fmt.Errorf("Unknown custom type [0x%02x]", customType)
	}
}

// ============================================================================

func assertCBEMarshalUnmarshalComplex(t *testing.T, value interface{}) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	marshalOpts := options.DefaultCBEMarshalerOptions()
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float32(0), float32(0)))] = convertComplex64ToCustom
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float64(0), float64(0)))] = convertComplex128ToCustom
	unmarshalOpts := options.DefaultCBEUnmarshalerOptions()
	unmarshalOpts.Session.CustomBinaryBuildFunction = convertFromCustom
	unmarshalOpts.Session.CustomBuiltTypes = append(unmarshalOpts.Session.CustomBuiltTypes, reflect.TypeOf(value))

	marshaler := cbe.NewMarshaler(marshalOpts, unmarshalOpts)

	document, err := marshaler.MarshalToBytes(value)
	if err != nil {
		t.Error(err)
		return
	}

	template := value
	actual, err := marshaler.UnmarshalFromBytes(document, template)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actual, value) {
		t.Errorf("Expected %v but got %v", describe.D(value), describe.D(actual))
	}
}

func assertCTEMarshalUnmarshalComplex(t *testing.T, value interface{}) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	marshalOpts := options.DefaultCTEMarshalerOptions()
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float32(0), float32(0)))] = convertComplex64ToCustom
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float64(0), float64(0)))] = convertComplex128ToCustom
	unmarshalOpts := options.DefaultCTEUnmarshalerOptions()
	unmarshalOpts.Session.CustomBinaryBuildFunction = convertFromCustom
	unmarshalOpts.Session.CustomBuiltTypes = append(unmarshalOpts.Session.CustomBuiltTypes, reflect.TypeOf(value))

	marshaler := cte.NewMarshaler(marshalOpts, unmarshalOpts)

	document, err := marshaler.MarshalToBytes(value)
	if err != nil {
		t.Error(err)
		return
	}

	template := value
	actual, err := marshaler.UnmarshalFromBytes(document, template)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actual, value) {
		t.Errorf("Expected %v but got %v", describe.D(value), describe.D(actual))
	}
}

func assertMarshalUnmarshalComplex(t *testing.T, value interface{}) {
	assertCBEMarshalUnmarshalComplex(t, value)
	assertCTEMarshalUnmarshalComplex(t, value)
}

// ============================================================================

func TestCustomBuildIter(t *testing.T) {
	assertMarshalUnmarshalComplex(t, complex(1, 1))
	assertMarshalUnmarshalComplex(t, complex(float64(1.0000000000000000000000000001), float64(1)))
}

// TODO: custom text
