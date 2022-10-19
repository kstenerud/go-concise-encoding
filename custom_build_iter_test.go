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
	"strings"
	"testing"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// Demonstration of the Concise Encoding "custom" data type.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom

// ============================================================================
// Implementation Code
// ============================================================================

// Assume the following custom serialized formats for complex64 and 128:
//
// | Offset | Size | Description                                             |
// | ------ | ---- | ------------------------------------------------------- |
// |      0 |    4 | Real portion (float32, little endian)                   |
// |      4 |    4 | Imaginary portion (float32, little endian)              |
//
// | Offset | Size | Description                                             |
// | ------ | ---- | ------------------------------------------------------- |
// |      0 |    8 | Real portion (float64, little endian)                   |
// |      8 |    8 | Imaginary portion (float64, little endian)              |

// We'll assign the following custom type codes:
const (
	typeCodeComplex64  = 0
	typeCodeComplex128 = 1
)

// First piece: functions to convert from complex type to custom bytes.
// These functions each handle a single type only.
func convertComplex64ToCustomBinary(rv reflect.Value) (customType uint64, asBytes []byte, err error) {
	cplx := complex64(rv.Complex())

	buff := bytes.Buffer{}
	if err = binary.Write(&buff, binary.LittleEndian, real(cplx)); err != nil {
		return
	}
	if err = binary.Write(&buff, binary.LittleEndian, imag(cplx)); err != nil {
		return
	}
	customType = typeCodeComplex64
	asBytes = buff.Bytes()
	return
}

func convertComplex128ToCustomBinary(rv reflect.Value) (customType uint64, asBytes []byte, err error) {
	cplx := rv.Complex()

	buff := bytes.Buffer{}
	if err = binary.Write(&buff, binary.LittleEndian, real(cplx)); err != nil {
		return
	}
	if err = binary.Write(&buff, binary.LittleEndian, imag(cplx)); err != nil {
		return
	}
	customType = typeCodeComplex128
	asBytes = buff.Bytes()
	return
}

// Second piece: converter function to fill in an object from custom data.
// This same function will be used for ALL custom types.
func convertFromCustomBinary(customType uint64, src []byte, dst reflect.Value) error {
	buff := bytes.NewBuffer(src)

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
		return fmt.Errorf("unknown custom type [0x%02x]", customType)
	}
}

// ============================================================================
// Test Code
// ============================================================================

func assertCBEMarshalUnmarshalComplexFromBinary(t *testing.T, value interface{}) {
	marshalOpts := options.DefaultCBEMarshalerOptions()
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float32(0), float32(0)))] = convertComplex64ToCustomBinary
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float64(0), float64(0)))] = convertComplex128ToCustomBinary
	unmarshalOpts := options.DefaultCEUnmarshalerOptions()
	unmarshalOpts.Session.CustomBinaryBuildFunction = convertFromCustomBinary
	unmarshalOpts.Session.CustomBuiltTypes = append(unmarshalOpts.Session.CustomBuiltTypes, reflect.TypeOf(value))

	marshaler := ce.NewCBEMarshaler(&marshalOpts)
	document, err := marshaler.MarshalToDocument(value)
	if err != nil {
		t.Error(err)
		return
	}

	template := value
	unmarshaler := ce.NewCBEUnmarshaler(&unmarshalOpts)
	actual, err := unmarshaler.UnmarshalFromDocument(document, template)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actual, value) {
		t.Errorf("Expected %v but got %v", describe.D(value), describe.D(actual))
	}
}

func assertCTEMarshalUnmarshalComplexFromBinary(t *testing.T, value interface{}) {
	marshalOpts := options.DefaultCTEMarshalerOptions()
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float32(0), float32(0)))] = convertComplex64ToCustomBinary
	marshalOpts.Session.CustomBinaryConverters[reflect.TypeOf(complex(float64(0), float64(0)))] = convertComplex128ToCustomBinary
	unmarshalOpts := options.DefaultCEUnmarshalerOptions()
	unmarshalOpts.Session.CustomBinaryBuildFunction = convertFromCustomBinary
	unmarshalOpts.Session.CustomBuiltTypes = append(unmarshalOpts.Session.CustomBuiltTypes, reflect.TypeOf(value))

	marshaler := ce.NewCTEMarshaler(&marshalOpts)
	document, err := marshaler.MarshalToDocument(value)
	if err != nil {
		t.Error(err)
		return
	}

	template := value
	unmarshaler := ce.NewCTEUnmarshaler(&unmarshalOpts)
	actual, err := unmarshaler.UnmarshalFromDocument(document, template)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actual, value) {
		t.Errorf("Expected %v but got %v", describe.D(value), describe.D(actual))
	}
}

func assertMarshalUnmarshalComplexFromBinary(t *testing.T, value interface{}) {
	assertCBEMarshalUnmarshalComplexFromBinary(t, value)
	assertCTEMarshalUnmarshalComplexFromBinary(t, value)
}

// ============================================================================

func convertComplexToCustomText(rv reflect.Value) (customType uint64, asString []byte, err error) {
	cplx := rv.Complex()
	switch rv.Kind() {
	case reflect.Complex64:
		customType = typeCodeComplex64
	case reflect.Complex128:
		customType = typeCodeComplex128
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%g+%gi", real(cplx), imag(cplx)))
	asString = []byte(builder.String())
	return
}

func convertFromCustomText(customType uint64, src string, dst reflect.Value) error {
	var r, i float64
	if _, err := fmt.Sscanf(src, "%f+%fi", &r, &i); err != nil {
		return err
	}

	dst.SetComplex(complex(r, i))
	return nil
}

func assertCTEMarshalUnmarshalComplexFromText(t *testing.T, value interface{}) {
	marshalOpts := options.DefaultCTEMarshalerOptions()
	marshalOpts.Session.CustomTextConverters[reflect.TypeOf(complex(float32(0), float32(0)))] = convertComplexToCustomText
	marshalOpts.Session.CustomTextConverters[reflect.TypeOf(complex(float64(0), float64(0)))] = convertComplexToCustomText
	unmarshalOpts := options.DefaultCEUnmarshalerOptions()
	unmarshalOpts.Session.CustomTextBuildFunction = convertFromCustomText
	unmarshalOpts.Session.CustomBuiltTypes = append(unmarshalOpts.Session.CustomBuiltTypes, reflect.TypeOf(value))

	marshaler := ce.NewCTEMarshaler(&marshalOpts)
	document, err := marshaler.MarshalToDocument(value)
	if err != nil {
		t.Error(err)
		return
	}

	template := value
	unmarshaler := ce.NewCTEUnmarshaler(&unmarshalOpts)
	actual, err := unmarshaler.UnmarshalFromDocument(document, template)
	if err != nil {
		t.Error(err)
		return
	}

	if !equivalence.IsEquivalent(actual, value) {
		t.Errorf("Expected %v but got %v", describe.D(value), describe.D(actual))
	}
}

func assertMarshalUnmarshalComplexFromText(t *testing.T, value interface{}) {
	assertCTEMarshalUnmarshalComplexFromText(t, value)
}

func assertMarshalUnmarshalComplex(t *testing.T, value interface{}) {
	assertMarshalUnmarshalComplexFromBinary(t, value)
	assertMarshalUnmarshalComplexFromText(t, value)
}

// ============================================================================

// Tests

func TestCustomBuildIter(t *testing.T) {
	assertMarshalUnmarshalComplex(t, complex(1, 1))
	assertMarshalUnmarshalComplex(t, complex(float64(1.0000000000000000000000000001), float64(1)))
}
