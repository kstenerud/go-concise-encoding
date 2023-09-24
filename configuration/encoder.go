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
	"math"
)

type EncoderConfiguration struct {
	CTE CTEEncoderConfiguration
}

func (_this *EncoderConfiguration) init() {
	_this.CTE.init()
}

var defaultEncoderConfiguration = EncoderConfiguration{
	CTE: defaultCTEEncoderConfiguration,
}

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

func (_this *CTEEncoderConfiguration) init() {}

var defaultCTEEncoderConfiguration = CTEEncoderConfiguration{
	Indent:            "    ",
	MaxColumn:         math.MaxUint64,
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
