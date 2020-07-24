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

// All options that can be used to fine-tune the behavior of various aspects of
// this library.
package options

import (
	"fmt"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/version"
)

// TODO: Opt: Convert line endings to escapes
// TODO: Opt: Don't convert escapes
// TODO: Don't use marker/ref on empty containers
// TODO: Builder that converts to string
// TODO: Iterator that converts from string to smaller type (numeric)
// TODO: Some method to notify that a string field should be encoded as a different type
// TODO: Optional spaces around `=` in maps

// The type of top-level container to assume is already opened (for implied
// structure documents). For normal operation, use TLContainerTypeNone.
type TLContainerType int

const (
	// Assume that no top-level container has already been opened.
	TLContainerTypeNone = iota
	// Assume that a list has already been opened.
	TLContainerTypeList
	// Assume that a map has already been opened.
	TLContainerTypeMap
)

// How to encode binary floats
type BinaryFloatEncodeAs int

const (
	// Use decimal encoding (1.2e4)
	BinaryFloatEncodeAsDecimal = iota
	// Use binary encoding (0x1.2p4)
	BinaryFloatEncodeAsBinary
)

// ============================================================================
// CBE Decoder

type CBEDecoderOptions struct {
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

func DefaultCBEDecoderOptions() *CBEDecoderOptions {
	return &CBEDecoderOptions{
		BufferSize: 2048,
	}
}

func (_this *CBEDecoderOptions) WithDefaultsApplied() *CBEDecoderOptions {
	if _this == nil {
		return DefaultCBEDecoderOptions()
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	return _this
}

// ============================================================================
// CBE Encoder

type CBEEncoderOptions struct {
	BufferSize int
}

func DefaultCBEEncoderOptions() *CBEEncoderOptions {
	return &CBEEncoderOptions{
		BufferSize: 1024,
	}
}

func (_this *CBEEncoderOptions) WithDefaultsApplied() *CBEEncoderOptions {
	if _this == nil {
		return DefaultCBEEncoderOptions()
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
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

// ============================================================================
// CTE Decoder

type CTEDecoderOptions struct {
	// TODO: ShouldZeroCopy option
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

func DefaultCTEDecoderOptions() *CTEDecoderOptions {
	return &CTEDecoderOptions{
		BufferSize: 4096,
	}
}

func (_this *CTEDecoderOptions) WithDefaultsApplied() *CTEDecoderOptions {
	if _this == nil {
		return DefaultCTEDecoderOptions()
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	return _this
}

// ============================================================================
// CTE Encoder

type CTEEncoderOptions struct {
	BufferSize int
	Indent     string
	// TODO: BinaryFloatEncoding option
	BinaryFloatEncoding BinaryFloatEncodeAs
}

func DefaultCTEEncoderOptions() *CTEEncoderOptions {
	return &CTEEncoderOptions{
		BufferSize: 1024,
	}
}

func (_this *CTEEncoderOptions) WithDefaultsApplied() *CTEEncoderOptions {
	if _this == nil {
		return DefaultCTEEncoderOptions()
	}

	// TODO: Check for default individual options

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

// ============================================================================
// Builder Session

// Fills out a value from custom binary data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-binary
type CustomBinaryBuildFunction func(src []byte, dst reflect.Value) error

// Fills out a value from custom text data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-text
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-text
type CustomTextBuildFunction func(src string, dst reflect.Value) error

type BuilderSessionOptions struct {
	// Specifies which types will be built using custom text/binary build
	// functions. You must also set one or both of CustomBinaryBuildFunction
	// and CustomTextBuildFunction in order to use this feature.
	// Both CBE and CTE will attempt to use either the binary or text version
	// depending on the data type (custom binary, custom text) encoded in the
	// source document.
	CustomBuiltTypes []reflect.Type

	// Build function to use when building from a custom binary source.
	CustomBinaryBuildFunction CustomBinaryBuildFunction

	// Build function to use when building from a custom text source.
	CustomTextBuildFunction CustomTextBuildFunction
}

func DefaultBuilderSessionOptions() *BuilderSessionOptions {
	return &BuilderSessionOptions{
		CustomBinaryBuildFunction: (func(src []byte, dst reflect.Value) error {
			return fmt.Errorf("No builder has been registered to handle custom binary data")
		}),
		CustomTextBuildFunction: (func(src string, dst reflect.Value) error {
			return fmt.Errorf("No builder has been registered to handle custom text data")
		}),
	}
}

func (_this *BuilderSessionOptions) WithDefaultsApplied() *BuilderSessionOptions {
	defaults := DefaultBuilderSessionOptions()
	if _this == nil {
		return defaults
	}

	if _this.CustomBinaryBuildFunction == nil {
		_this.CustomBinaryBuildFunction = defaults.CustomBinaryBuildFunction
	}
	if _this.CustomTextBuildFunction == nil {
		_this.CustomTextBuildFunction = defaults.CustomTextBuildFunction
	}
	if _this.CustomBuiltTypes == nil {
		_this.CustomBuiltTypes = []reflect.Type{}
	}

	return _this
}

// ============================================================================
// Builder

type BuilderOptions struct {
	FloatToBigIntMaxBase10Exponent int
	FloatToBigIntMaxBase2Exponent  int
	// TODO: ErrorOnLossyFloatConversion option
	ErrorOnLossyFloatConversion bool
	// TODO: Something for decimal floats?
	// TODO: Error on unknown field
}

func DefaultBuilderOptions() *BuilderOptions {
	return &BuilderOptions{
		FloatToBigIntMaxBase10Exponent: 300,
		FloatToBigIntMaxBase2Exponent:  300 * 10 / 3,
	}
}

func (_this *BuilderOptions) WithDefaultsApplied() *BuilderOptions {
	if _this == nil {
		return DefaultBuilderOptions()
	}

	// TODO: Check for default individual options

	return _this
}

// ============================================================================
// Iterator Session

// Converts a value to custom binary data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-binary
type ConvertToCustomBinaryFunction func(v reflect.Value) (asBytes []byte, err error)

// Converts a value to custom text data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-text
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-text
type ConvertToCustomTextFunction func(v reflect.Value) (asText string, err error)

type IteratorSessionOptions struct {
	// Specifies which types to convert to custom binary data, and how to do it.
	// Note: You should only fill out one of these maps, depending on your
	// indended encoding (binary or text). The iterator session will consult
	// the binary map first and the text map second, choosing the first match.
	CustomBinaryConverters map[reflect.Type]ConvertToCustomBinaryFunction

	// Specifies which types to convert to custom text data, and how to do it
	// Note: You should only fill out one of these maps, depending on your
	// indended encoding (binary or text). The iterator session will consult
	// the binary map first and the text map second, choosing the first match.
	CustomTextConverters map[reflect.Type]ConvertToCustomTextFunction
}

func DefaultIteratorSessionOptions() *IteratorSessionOptions {
	return &IteratorSessionOptions{
		CustomBinaryConverters: make(map[reflect.Type]ConvertToCustomBinaryFunction),
		CustomTextConverters:   make(map[reflect.Type]ConvertToCustomTextFunction),
	}
}

func (_this *IteratorSessionOptions) WithDefaultsApplied() *IteratorSessionOptions {
	if _this == nil {
		return DefaultIteratorSessionOptions()
	}

	if _this.CustomBinaryConverters == nil {
		_this.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomBinaryFunction)
	}
	if _this.CustomTextConverters == nil {
		_this.CustomTextConverters = make(map[reflect.Type]ConvertToCustomTextFunction)
	}

	return _this
}

// ============================================================================
// Iterator

type IteratorOptions struct {
	ConciseEncodingVersion uint64
	// If useReferences is true, the iterator will also look for duplicate
	// pointers to data, generating marker and reference events rather than
	// walking the object again. This is useful for cyclic or recursive data
	// structures.
	UseReferences bool
	// TODO
	OmitNilPointers bool
}

func DefaultIteratorOptions() *IteratorOptions {
	return &IteratorOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *IteratorOptions) WithDefaultsApplied() *IteratorOptions {
	if _this == nil {
		return DefaultIteratorOptions()
	}

	if _this.ConciseEncodingVersion < 1 {
		_this.ConciseEncodingVersion = DefaultIteratorOptions().ConciseEncodingVersion
	}

	return _this
}

// ============================================================================
// Rules

type RuleOptions struct {
	depthHasBeenAdjusted bool
	// Concise encoding spec version to adhere to
	ConciseEncodingVersion uint64

	// Limits before the ruleset artificially terminates with an error.
	MaxBytesLength      uint64
	MaxStringLength     uint64
	MaxURILength        uint64
	MaxIDLength         uint64
	MaxMarkupNameLength uint64
	MaxContainerDepth   uint64
	MaxObjectCount      uint64
	MaxReferenceCount   uint64
	// Max bytes total for all array types
}

func DefaultRuleOptions() *RuleOptions {
	return &RuleOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
		MaxBytesLength:         1000000000,
		MaxStringLength:        100000000,
		MaxURILength:           10000,
		MaxIDLength:            100,
		MaxMarkupNameLength:    100,
		MaxContainerDepth:      1000,
		MaxObjectCount:         10000000,
		MaxReferenceCount:      100000,
		// TODO: References need to check for amplification attacks. Keep count of referenced things and their object counts
	}
}

func (_this *RuleOptions) WithDefaultsApplied() *RuleOptions {
	defaults := DefaultRuleOptions()
	if _this == nil {
		return defaults
	}

	if _this.ConciseEncodingVersion < 1 {
		_this.ConciseEncodingVersion = defaults.ConciseEncodingVersion
	}
	if _this.MaxBytesLength < 1 {
		_this.MaxBytesLength = defaults.MaxBytesLength
	}
	if _this.MaxStringLength < 1 {
		_this.MaxStringLength = defaults.MaxStringLength
	}
	if _this.MaxURILength < 1 {
		_this.MaxURILength = defaults.MaxURILength
	}
	if _this.MaxIDLength < 1 {
		_this.MaxIDLength = defaults.MaxIDLength
	}
	if _this.MaxMarkupNameLength < 1 {
		_this.MaxMarkupNameLength = defaults.MaxMarkupNameLength
	}
	if _this.MaxContainerDepth < 1 {
		_this.MaxContainerDepth = defaults.MaxContainerDepth
	}
	if _this.MaxObjectCount < 1 {
		_this.MaxObjectCount = defaults.MaxObjectCount
	}
	if _this.MaxReferenceCount < 1 {
		_this.MaxReferenceCount = defaults.MaxReferenceCount
	}

	return _this
}
