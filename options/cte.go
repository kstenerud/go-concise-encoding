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

package options

import (
	"fmt"

	"github.com/kstenerud/go-concise-encoding/version"
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
// CTE Decoder

type CTEDecoderOptions struct {
	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	ConciseEncodingVersion uint64
}

func DefaultCTEDecoderOptions() *CTEDecoderOptions {
	return &CTEDecoderOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *CTEDecoderOptions) WithDefaultsApplied() *CTEDecoderOptions {
	if _this == nil {
		return DefaultCTEDecoderOptions()
	}

	if _this.ConciseEncodingVersion == 0 {
		_this.ConciseEncodingVersion = version.ConciseEncodingVersion
	}

	return _this
}

func (_this *CTEDecoderOptions) Validate() error {
	return nil
}

// ============================================================================
// CTE Encoder

type CTEEncoderOptions struct {
	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	ConciseEncodingVersion uint64

	// Indentation to use when pretty printing
	Indent string

	// TODO: Max column before forcing a newline (if possible)
	MaxColumn uint

	// TODO: Convert line endings to escapes
	EscapeLineEndings bool

	DefaultFormats struct {
		Int   CTEEncodingFormat
		Uint  CTEEncodingFormat
		Float CTEEncodingFormat
		Array struct {
			Int8    CTEEncodingFormat
			Int16   CTEEncodingFormat
			Int32   CTEEncodingFormat
			Int64   CTEEncodingFormat
			Uint8   CTEEncodingFormat
			Uint16  CTEEncodingFormat
			Uint32  CTEEncodingFormat
			Uint64  CTEEncodingFormat
			Float16 CTEEncodingFormat
			Float32 CTEEncodingFormat
			Float64 CTEEncodingFormat
		}
	}
}

type CTEEncodingFormat uint8

const (
	CTEEncodingFormatUnset          CTEEncodingFormat = 0
	CTEEncodingFormatFlagZeroFilled CTEEncodingFormat = 1
	CTEEncodingFormatBinary         CTEEncodingFormat = 2 + iota
	CTEEncodingFormatBinaryZeroFilled
	CTEEncodingFormatOctal
	CTEEncodingFormatOctalZeroFilled
	CTEEncodingFormatHexadecimal
	CTEEncodingFormatHexadecimalZeroFilled
	cteEncodingFormatCount
)

var cteEncodingFormatStrings = []string{
	CTEEncodingFormatUnset:                 "CTEEncodingFormatUnset",
	CTEEncodingFormatFlagZeroFilled:        "CTEEncodingFormatFlagZeroFilled",
	CTEEncodingFormatBinary:                "CTEEncodingFormatBinary",
	CTEEncodingFormatBinaryZeroFilled:      "CTEEncodingFormatBinaryZeroFilled",
	CTEEncodingFormatOctal:                 "CTEEncodingFormatOctal",
	CTEEncodingFormatOctalZeroFilled:       "CTEEncodingFormatOctalZeroFilled",
	CTEEncodingFormatHexadecimal:           "CTEEncodingFormatHexadecimal",
	CTEEncodingFormatHexadecimalZeroFilled: "CTEEncodingFormatHexadecimalZeroFilled",
}

func (_this CTEEncodingFormat) String() string {
	if _this < cteEncodingFormatCount {
		return cteEncodingFormatStrings[_this]
	}
	return fmt.Sprintf("CTEEncodingFormat(%d)", _this)
}

func DefaultCTEEncoderOptions() *CTEEncoderOptions {
	opts := &CTEEncoderOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
		Indent:                 "    ",
		MaxColumn:              0,
		EscapeLineEndings:      true,
	}

	opts.DefaultFormats.Float = CTEEncodingFormatUnset
	opts.DefaultFormats.Int = CTEEncodingFormatUnset
	opts.DefaultFormats.Uint = CTEEncodingFormatUnset
	opts.DefaultFormats.Array.Int8 = CTEEncodingFormatUnset
	opts.DefaultFormats.Array.Int16 = CTEEncodingFormatUnset
	opts.DefaultFormats.Array.Int32 = CTEEncodingFormatUnset
	opts.DefaultFormats.Array.Int64 = CTEEncodingFormatUnset
	opts.DefaultFormats.Array.Uint8 = CTEEncodingFormatHexadecimalZeroFilled
	opts.DefaultFormats.Array.Uint16 = CTEEncodingFormatHexadecimalZeroFilled
	opts.DefaultFormats.Array.Uint32 = CTEEncodingFormatHexadecimalZeroFilled
	opts.DefaultFormats.Array.Uint64 = CTEEncodingFormatHexadecimalZeroFilled
	opts.DefaultFormats.Array.Float16 = CTEEncodingFormatHexadecimal
	opts.DefaultFormats.Array.Float32 = CTEEncodingFormatHexadecimal
	opts.DefaultFormats.Array.Float64 = CTEEncodingFormatHexadecimal

	return opts
}

func (_this *CTEEncoderOptions) WithDefaultsApplied() *CTEEncoderOptions {
	if _this == nil {
		return DefaultCTEEncoderOptions()
	}

	if _this.ConciseEncodingVersion == 0 {
		_this.ConciseEncodingVersion = version.ConciseEncodingVersion
	}

	return _this
}

var cteValidFloatFormats = map[CTEEncodingFormat]bool{
	CTEEncodingFormatUnset:       true,
	CTEEncodingFormatHexadecimal: true,
}
var cteValidIntFormats = map[CTEEncodingFormat]bool{
	CTEEncodingFormatBinary:      true,
	CTEEncodingFormatUnset:       true,
	CTEEncodingFormatOctal:       true,
	CTEEncodingFormatHexadecimal: true,
}
var cteValidArrayIntFormats = map[CTEEncodingFormat]bool{
	CTEEncodingFormatBinary:                true,
	CTEEncodingFormatBinaryZeroFilled:      true,
	CTEEncodingFormatUnset:                 true,
	CTEEncodingFormatOctal:                 true,
	CTEEncodingFormatOctalZeroFilled:       true,
	CTEEncodingFormatHexadecimal:           true,
	CTEEncodingFormatHexadecimalZeroFilled: true,
}

func (_this *CTEEncoderOptions) Validate() (err error) {
	validate := func(name string, format CTEEncodingFormat, validFormats map[CTEEncodingFormat]bool) error {
		if validFormats[format] {
			return nil
		}
		return fmt.Errorf("%v is not a valid encoding format for %s", format, name)
	}

	if err = validate("Int", _this.DefaultFormats.Int, cteValidIntFormats); err != nil {
		return
	}
	if err = validate("Uint", _this.DefaultFormats.Uint, cteValidIntFormats); err != nil {
		return
	}
	if err = validate("Float", _this.DefaultFormats.Float, cteValidFloatFormats); err != nil {
		return
	}

	if err = validate("Int8", _this.DefaultFormats.Array.Int8, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int16", _this.DefaultFormats.Array.Int16, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int32", _this.DefaultFormats.Array.Int32, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Int64", _this.DefaultFormats.Array.Int64, cteValidArrayIntFormats); err != nil {
		return
	}

	if err = validate("Uint8", _this.DefaultFormats.Array.Uint8, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint16", _this.DefaultFormats.Array.Uint16, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint32", _this.DefaultFormats.Array.Uint32, cteValidArrayIntFormats); err != nil {
		return
	}
	if err = validate("Uint64", _this.DefaultFormats.Array.Uint64, cteValidArrayIntFormats); err != nil {
		return
	}

	if err = validate("Float16", _this.DefaultFormats.Array.Float16, cteValidFloatFormats); err != nil {
		return
	}
	if err = validate("Float32", _this.DefaultFormats.Array.Float32, cteValidFloatFormats); err != nil {
		return
	}
	if err = validate("Float64", _this.DefaultFormats.Array.Float64, cteValidFloatFormats); err != nil {
		return
	}

	return
}

// ============================================================================
// CTE Marshaler

type CTEMarshalerOptions struct {
	Encoder  CTEEncoderOptions
	Iterator IteratorOptions
	Session  IteratorSessionOptions
}

func DefaultCTEMarshalerOptions() *CTEMarshalerOptions {
	return &CTEMarshalerOptions{
		Encoder:  *DefaultCTEEncoderOptions(),
		Iterator: *DefaultIteratorOptions(),
		Session:  *DefaultIteratorSessionOptions(),
	}
}

func (_this *CTEMarshalerOptions) WithDefaultsApplied() *CTEMarshalerOptions {
	if _this == nil {
		return DefaultCTEMarshalerOptions()
	}

	_this.Encoder.WithDefaultsApplied()
	_this.Iterator.WithDefaultsApplied()
	_this.Session.WithDefaultsApplied()

	return _this
}

func (_this *CTEMarshalerOptions) Validate() error {
	if err := _this.Encoder.Validate(); err != nil {
		return err
	}
	if err := _this.Iterator.Validate(); err != nil {
		return err
	}
	return _this.Session.Validate()
}

// ============================================================================
// CTE Unmarshaler

type CTEUnmarshalerOptions struct {
	Decoder CTEDecoderOptions
	Builder BuilderOptions
	Session BuilderSessionOptions
	Rules   RuleOptions

	// If false, do not wrap a Rules object around the builder, disabling all rule checks.
	EnforceRules bool
}

func DefaultCTEUnmarshalerOptions() *CTEUnmarshalerOptions {
	return &CTEUnmarshalerOptions{
		Decoder:      *DefaultCTEDecoderOptions(),
		Builder:      *DefaultBuilderOptions(),
		Session:      *DefaultBuilderSessionOptions(),
		Rules:        *DefaultRuleOptions(),
		EnforceRules: true,
	}
}

func (_this *CTEUnmarshalerOptions) WithDefaultsApplied() *CTEUnmarshalerOptions {
	if _this == nil {
		return DefaultCTEUnmarshalerOptions()
	}

	_this.Decoder.WithDefaultsApplied()
	_this.Builder.WithDefaultsApplied()
	_this.Session.WithDefaultsApplied()
	_this.Rules.WithDefaultsApplied()

	return _this
}

func (_this *CTEUnmarshalerOptions) Validate() error {
	if err := _this.Builder.Validate(); err != nil {
		return err
	}
	if err := _this.Decoder.Validate(); err != nil {
		return err
	}
	if err := _this.Rules.Validate(); err != nil {
		return err
	}
	return _this.Session.Validate()
}
