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

	"github.com/kstenerud/go-concise-encoding/cte"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/debug"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// Demonstration of the Concise Encoding "custom" data type.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom

// Demonstration custom type. Technically, this would be better off encoded as
// a map, but it's best to demonstrate with something simple!
type Measure struct {
	Units string
	Value float64
}

// ============================================================================

// Assume the following custom serialized format for Measure:
//
// | Offset | Size | Description                                             |
// | ------ | ---- | ------------------------------------------------------- |
// |      0 |    1 | Data type code (so we can have multiple custom types)   |
// |      1 |    8 | Value (float64, little endian)                          |
// |      9 |    1 | Units name length (max 255)                             |
// |     10 |    n | Units name (string with length from units name length)  |

// Simple type mechanism: The first byte of the data is the type field
const typeCodeMeasure = 1

// First piece: function to convert from an object to custom bytes.
// This function handles a single type only.
func convertToCustom(v reflect.Value) (asBytes []byte, err error) {
	value := v.FieldByName("Value").Float()
	units := v.FieldByName("Units").String()
	unitsLength := len(units)

	buff := bytes.Buffer{}
	buff.WriteByte(typeCodeMeasure)
	if err = binary.Write(&buff, binary.LittleEndian, value); err != nil {
		return
	}
	buff.WriteByte(byte(unitsLength))
	buff.Write([]byte(units))
	asBytes = buff.Bytes()
	return
}

// Second piece: converter function to fill in an object from custom data.
// This same function will be used for ALL custom types.
func convertFromCustom(src []byte, dst reflect.Value) error {
	buff := bytes.NewBuffer(src)

	customType, _ := buff.ReadByte()
	switch customType {
	case typeCodeMeasure:
		var value float64
		if err := binary.Read(buff, binary.LittleEndian, &value); err != nil {
			return err
		}
		unitsLength, _ := buff.ReadByte()
		units := make([]byte, unitsLength)
		byteCount, err := buff.Read(units)
		if err != nil {
			return err
		}
		if byteCount < len(units) {
			return fmt.Errorf("Incomplete data")
		}
		dst.FieldByName("Value").SetFloat(value)
		dst.FieldByName("Units").SetString(string(units))
		return nil
	default:
		return fmt.Errorf("%x: Unknown custom type", customType)
	}
}

// ============================================================================

func assertCBEMarshalUnmarshalMeasure(t *testing.T, value *Measure) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	marshaler := cbe.NewMarshaler(nil, nil)
	marshaler.BuilderSession.SetCustomBuildFunction(convertFromCustom)
	marshaler.BuilderSession.UseCustomBuildFunctionForType(reflect.TypeOf(*value))
	marshaler.IteratorSession.RegisterCustomConverterForType(reflect.TypeOf(*value), convertToCustom)

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

func assertCTEMarshalUnmarshalMeasure(t *testing.T, value *Measure) {
	debug.DebugOptions.PassThroughPanics = true
	defer func() { debug.DebugOptions.PassThroughPanics = false }()

	marshaler := cte.NewMarshaler(nil, nil)
	marshaler.BuilderSession.SetCustomBuildFunction(convertFromCustom)
	marshaler.BuilderSession.UseCustomBuildFunctionForType(reflect.TypeOf(*value))
	marshaler.IteratorSession.RegisterCustomConverterForType(reflect.TypeOf(*value), convertToCustom)

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

func assertMarshalUnmarshalMeasure(t *testing.T, value *Measure) {
	assertCBEMarshalUnmarshalMeasure(t, value)
	assertCTEMarshalUnmarshalMeasure(t, value)
}

// ============================================================================

func TestCustomBuildIter(t *testing.T) {
	assertMarshalUnmarshalMeasure(t, &Measure{
		Units: "cm",
		Value: 1.8,
	})
}
