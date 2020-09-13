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
	// The size of the underlying buffer to use when decoding a document.
	BufferSize int

	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	// This value is consulted if ImpliedStructure is anything other than
	// ImpliedStructureNone.
	ConciseEncodingVersion uint64

	// The implied structure that this decoder will assume.
	// Any implied structure will be automatically reported without being
	// present in the document.
	ImpliedStructure ImpliedStructure
}

func DefaultCTEDecoderOptions() *CTEDecoderOptions {
	return &CTEDecoderOptions{
		BufferSize:             4096,
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *CTEDecoderOptions) WithDefaultsApplied() *CTEDecoderOptions {
	if _this == nil {
		return DefaultCTEDecoderOptions()
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	if _this.ConciseEncodingVersion == 0 {
		_this.ConciseEncodingVersion = version.ConciseEncodingVersion
	}

	return _this
}

// ============================================================================
// CTE Encoder

type CTEEncoderOptions struct {
	BufferSize int

	// Indentation to use when pretty printing
	Indent string

	// TODO: Default encoding to use for floating point
	DefaultFloatEncoding FloatEncoding

	// TODO: Default encoding to use for integers
	DefaultIntEncoding IntEncoding

	// TODO: Max column before forcing a newline (if possible)
	MaxColumn uint

	// TODO: Convert line endings to escapes
	EscapeLineEndings bool

	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	// This value is consulted if ImpliedStructure is anything other than
	// ImpliedStructureNone.
	ConciseEncodingVersion uint64

	// The implied structure that this encoder will assume.
	// Any implied structure will not actually be written to the document.
	ImpliedStructure ImpliedStructure

	DefaultArrayEncodingBases ArrayEncodingBases
}

// The base to encode numbers in for numeric arrays.
// A value of 0 or 10 means decimal, 2 = binary, 8 = octal, 16 = hexadecimal.
type ArrayEncodingBases struct {
	Int8    int
	Int16   int
	Int32   int
	Int64   int
	Uint8   int
	Uint16  int
	Uint32  int
	Uint64  int
	Float16 int
	Float32 int
	Float64 int
}

func DefaultCTEEncoderOptions() *CTEEncoderOptions {
	return &CTEEncoderOptions{
		BufferSize:             1024,
		ConciseEncodingVersion: version.ConciseEncodingVersion,
		DefaultArrayEncodingBases: ArrayEncodingBases{
			Int8:    10,
			Int16:   10,
			Int32:   10,
			Int64:   10,
			Uint8:   16,
			Uint16:  16,
			Uint32:  16,
			Uint64:  16,
			Float16: 16,
			Float32: 16,
			Float64: 16,
		},
	}
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

// ============================================================================
// CTE Unmarshaler

type CTEUnmarshalerOptions struct {
	Decoder CTEDecoderOptions
	Builder BuilderOptions
	Session BuilderSessionOptions
	Rules   RuleOptions
}

func DefaultCTEUnmarshalerOptions() *CTEUnmarshalerOptions {
	return &CTEUnmarshalerOptions{
		Decoder: *DefaultCTEDecoderOptions(),
		Builder: *DefaultBuilderOptions(),
		Session: *DefaultBuilderSessionOptions(),
		Rules:   *DefaultRuleOptions(),
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
