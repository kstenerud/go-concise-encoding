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

// ============================================================================
// CBE Decoder

type CBEDecoderOptions struct {
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

func DefaultCBEDecoderOptions() *CBEDecoderOptions {
	return &CBEDecoderOptions{
		BufferSize:             2048,
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *CBEDecoderOptions) WithDefaultsApplied() *CBEDecoderOptions {
	if _this == nil {
		return DefaultCBEDecoderOptions()
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
// CBE Encoder

type CBEEncoderOptions struct {
	// The size of the underlying buffer to use when encoding a document.
	BufferSize int

	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	// This value is consulted if ImpliedStructure is anything other than
	// ImpliedStructureNone.
	ConciseEncodingVersion uint64

	// The implied structure that this encoder will assume.
	// Any implied structure will not actually be written to the document.
	ImpliedStructure ImpliedStructure
}

func DefaultCBEEncoderOptions() *CBEEncoderOptions {
	return &CBEEncoderOptions{
		BufferSize:             1024,
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *CBEEncoderOptions) WithDefaultsApplied() *CBEEncoderOptions {
	if _this == nil {
		return DefaultCBEEncoderOptions()
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
// CBE Marshaler

type CBEMarshalerOptions struct {
	Encoder  CBEEncoderOptions
	Iterator IteratorOptions
	Session  IteratorSessionOptions
}

func DefaultCBEMarshalerOptions() *CBEMarshalerOptions {
	return &CBEMarshalerOptions{
		Encoder:  *DefaultCBEEncoderOptions(),
		Iterator: *DefaultIteratorOptions(),
		Session:  *DefaultIteratorSessionOptions(),
	}
}

func (_this *CBEMarshalerOptions) WithDefaultsApplied() *CBEMarshalerOptions {
	if _this == nil {
		return DefaultCBEMarshalerOptions()
	}

	_this.Encoder.WithDefaultsApplied()
	_this.Iterator.WithDefaultsApplied()
	_this.Session.WithDefaultsApplied()

	return _this
}

// ============================================================================
// CBE Unmarshaler

type CBEUnmarshalerOptions struct {
	Decoder CBEDecoderOptions
	Builder BuilderOptions
	Session BuilderSessionOptions
	Rules   RuleOptions
}

func DefaultCBEUnmarshalerOptions() *CBEUnmarshalerOptions {
	return &CBEUnmarshalerOptions{
		Decoder: *DefaultCBEDecoderOptions(),
		Builder: *DefaultBuilderOptions(),
		Session: *DefaultBuilderSessionOptions(),
		Rules:   *DefaultRuleOptions(),
	}
}

func (_this *CBEUnmarshalerOptions) WithDefaultsApplied() *CBEUnmarshalerOptions {
	if _this == nil {
		return DefaultCBEUnmarshalerOptions()
	}

	_this.Decoder.WithDefaultsApplied()
	_this.Builder.WithDefaultsApplied()
	_this.Session.WithDefaultsApplied()
	_this.Rules.WithDefaultsApplied()

	return _this
}
