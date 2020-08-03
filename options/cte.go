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
}

func DefaultCTEEncoderOptions() *CTEEncoderOptions {
	return &CTEEncoderOptions{
		BufferSize:             1024,
		ConciseEncodingVersion: version.ConciseEncodingVersion,
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
