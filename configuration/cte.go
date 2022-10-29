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

package configuration

import (
	"fmt"
)

// How to encode binary floats
type FloatEncoding int

const (
	// Use decimal encoding (1.2e4)
	FloatEncodingDecimal = iota
	// Use binary encoding (0x1.2p4)
	FloatEncodingBinary
)

type IntEncoding int

const (
	IntEncodingDecimal = iota
	IntEncodingBinary
	IntEncodingOctal
	IntEncodingHexadecimal
)

// ============================================================================
// CTE Encoder

type CTEEncoderConfiguration struct {
	// Indentation to use when pretty printing
	Indent string

	// TODO: Max column before forcing a newline (if possible)
	MaxColumn uint

	// TODO: Convert line endings to escapes
	EscapeLineEndings bool

	DefaultNumericFormats CTEEncoderDefaultNumericFormats
}

type CTEEncoderDefaultNumericFormats struct {
	Int         CTENumericFormat
	Uint        CTENumericFormat
	BinaryFloat CTENumericFormat
	Array       CTEEncoderDefaultArrayFormats
}

type CTEEncoderDefaultArrayFormats struct {
	Int8    CTENumericFormat
	Int16   CTENumericFormat
	Int32   CTENumericFormat
	Int64   CTENumericFormat
	Uint8   CTENumericFormat
	Uint16  CTENumericFormat
	Uint32  CTENumericFormat
	Uint64  CTENumericFormat
	Float16 CTENumericFormat
	Float32 CTENumericFormat
	Float64 CTENumericFormat
}

type CTENumericFormat uint8

const (
	CTEEncodingFormatDecimal        CTENumericFormat = 0
	CTEEncodingFormatFlagZeroFilled CTENumericFormat = 1
	CTEEncodingFormatBinary         CTENumericFormat = 2 + iota
	CTEEncodingFormatBinaryZeroFilled
	CTEEncodingFormatOctal
	CTEEncodingFormatOctalZeroFilled
	CTEEncodingFormatHexadecimal
	CTEEncodingFormatHexadecimalZeroFilled
	cteEncodingFormatCount
)

var cteEncodingFormatStrings = []string{
	CTEEncodingFormatDecimal:               "CTEEncodingFormatDecimal",
	CTEEncodingFormatFlagZeroFilled:        "CTEEncodingFormatFlagZeroFilled",
	CTEEncodingFormatBinary:                "CTEEncodingFormatBinary",
	CTEEncodingFormatBinaryZeroFilled:      "CTEEncodingFormatBinaryZeroFilled",
	CTEEncodingFormatOctal:                 "CTEEncodingFormatOctal",
	CTEEncodingFormatOctalZeroFilled:       "CTEEncodingFormatOctalZeroFilled",
	CTEEncodingFormatHexadecimal:           "CTEEncodingFormatHexadecimal",
	CTEEncodingFormatHexadecimalZeroFilled: "CTEEncodingFormatHexadecimalZeroFilled",
}

func (_this CTENumericFormat) String() string {
	if _this < cteEncodingFormatCount {
		return cteEncodingFormatStrings[_this]
	}
	return fmt.Sprintf("CTEEncodingFormat(%d)", _this)
}

func DefaultCTEEncoderConfiguration() CTEEncoderConfiguration {
	return defaultCTEEncoderConfiguration
}

var defaultCTEEncoderConfiguration = CTEEncoderConfiguration{
	Indent:            "    ",
	MaxColumn:         0,
	EscapeLineEndings: true,
	DefaultNumericFormats: CTEEncoderDefaultNumericFormats{
		BinaryFloat: CTEEncodingFormatHexadecimal,
		Int:         CTEEncodingFormatDecimal,
		Uint:        CTEEncodingFormatDecimal,
		Array: CTEEncoderDefaultArrayFormats{
			Int8:    CTEEncodingFormatDecimal,
			Int16:   CTEEncodingFormatDecimal,
			Int32:   CTEEncodingFormatDecimal,
			Int64:   CTEEncodingFormatDecimal,
			Uint8:   CTEEncodingFormatDecimal,
			Uint16:  CTEEncodingFormatDecimal,
			Uint32:  CTEEncodingFormatDecimal,
			Uint64:  CTEEncodingFormatDecimal,
			Float16: CTEEncodingFormatHexadecimal,
			Float32: CTEEncodingFormatHexadecimal,
			Float64: CTEEncodingFormatHexadecimal,
		},
	},
}

func (_this *CTEEncoderConfiguration) ApplyDefaults() {
	// Nothing to do
}

var cteValidFloatFormats = map[CTENumericFormat]bool{
	CTEEncodingFormatDecimal:     true,
	CTEEncodingFormatHexadecimal: true,
}
var cteValidIntFormats = map[CTENumericFormat]bool{
	CTEEncodingFormatBinary:      true,
	CTEEncodingFormatDecimal:     true,
	CTEEncodingFormatOctal:       true,
	CTEEncodingFormatHexadecimal: true,
}
var cteValidArrayIntFormats = map[CTENumericFormat]bool{
	CTEEncodingFormatBinary:                true,
	CTEEncodingFormatBinaryZeroFilled:      true,
	CTEEncodingFormatDecimal:               true,
	CTEEncodingFormatOctal:                 true,
	CTEEncodingFormatOctalZeroFilled:       true,
	CTEEncodingFormatHexadecimal:           true,
	CTEEncodingFormatHexadecimalZeroFilled: true,
}

func (_this *CTEEncoderConfiguration) Validate() (err error) {
	validate := func(name string, format CTENumericFormat, validFormats map[CTENumericFormat]bool) error {
		if validFormats[format] {
			return nil
		}
		return fmt.Errorf("%v is not a valid encoding format for %s", format, name)
	}

	if err = validate("Int", _this.DefaultNumericFormats.Int, cteValidIntFormats); err != nil {
		return
	}
	if err = validate("Uint", _this.DefaultNumericFormats.Uint, cteValidIntFormats); err != nil {
		return
	}
	if err = validate("Binary float", _this.DefaultNumericFormats.BinaryFloat, cteValidFloatFormats); err != nil {
		return
	}

	if err = validate("Int8", _this.DefaultNumericFormats.Array.Int8, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int16", _this.DefaultNumericFormats.Array.Int16, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int32", _this.DefaultNumericFormats.Array.Int32, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int64", _this.DefaultNumericFormats.Array.Int64, cteValidArrayIntFormats); err != nil {
		return
	}

	if err = validate("Uint8", _this.DefaultNumericFormats.Array.Uint8, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint16", _this.DefaultNumericFormats.Array.Uint16, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint32", _this.DefaultNumericFormats.Array.Uint32, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint64", _this.DefaultNumericFormats.Array.Uint64, cteValidArrayIntFormats); err != nil {
		return
	}

	if err = validate("Float16", _this.DefaultNumericFormats.Array.Float16, cteValidFloatFormats); err != nil {
		return
	}
	if err = validate("Float32", _this.DefaultNumericFormats.Array.Float32, cteValidFloatFormats); err != nil {
		return
	}
	if err = validate("Float64", _this.DefaultNumericFormats.Array.Float64, cteValidFloatFormats); err != nil {
		return
	}

	return
}

// ============================================================================
// CTE Marshaler

type CTEMarshalerConfiguration struct {
	Encoder     CTEEncoderConfiguration
	Iterator    IteratorConfiguration
	DebugPanics bool
}

func DefaultCTEMarshalerConfiguration() CTEMarshalerConfiguration {
	return CTEMarshalerConfiguration{
		Encoder:  DefaultCTEEncoderConfiguration(),
		Iterator: DefaultIteratorConfiguration(),
	}
}

func (_this *CTEMarshalerConfiguration) ApplyDefaults() {
	_this.Encoder.ApplyDefaults()
	_this.Iterator.ApplyDefaults()
}

func (_this *CTEMarshalerConfiguration) Validate() error {
	if err := _this.Encoder.Validate(); err != nil {
		return err
	}
	return _this.Iterator.Validate()
}
