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

type CBEDecoderOptions struct {
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

var defaultCBEDecoderOptions = CBEDecoderOptions{
	BufferSize: 2048,
}

func DefaultCBEDecoderOptions() *CBEDecoderOptions {
	options := defaultCBEDecoderOptions
	return &options
}

func (_this *CBEDecoderOptions) ApplyDefaults() *CBEDecoderOptions {
	defaults := &defaultCBEDecoderOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	return _this
}

// ============================================================================

type CBEEncoderOptions struct {
	BufferSize int
}

var defaultCBEEncoderOptions = CBEEncoderOptions{
	BufferSize: 1024,
}

func (_this *CBEEncoderOptions) ApplyDefaults() *CBEEncoderOptions {
	defaults := &defaultCBEEncoderOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	return _this
}

func DefaultCBEEncoderOptions() *CBEEncoderOptions {
	options := defaultCBEEncoderOptions
	return &options
}

// ============================================================================

type CBEMarshalerOptions struct {
	Encoder  CBEEncoderOptions
	Iterator IteratorOptions
}

var defaultCBEMarshalerOptions = CBEMarshalerOptions{
	Encoder:  defaultCBEEncoderOptions,
	Iterator: defaultIteratorOptions,
}

func DefaultCBEMarshalerOptions() *CBEMarshalerOptions {
	options := defaultCBEMarshalerOptions
	return &options
}

func (_this *CBEMarshalerOptions) ApplyDefaults() *CBEMarshalerOptions {
	defaults := &defaultCBEMarshalerOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	_this.Encoder.ApplyDefaults()
	_this.Iterator.ApplyDefaults()

	return _this
}

// ============================================================================

type CBEUnmarshalerOptions struct {
	Decoder CBEDecoderOptions
	Builder BuilderOptions
	Rules   RuleOptions
}

var defaultCBEUnmarshalerOptions = CBEUnmarshalerOptions{
	Decoder: defaultCBEDecoderOptions,
	Builder: defaultBuilderOptions,
	Rules:   defaultRuleOptions,
}

func DefaultCBEUnmarshalerOptions() *CBEUnmarshalerOptions {
	options := defaultCBEUnmarshalerOptions
	return &options
}

func (_this *CBEUnmarshalerOptions) ApplyDefaults() *CBEUnmarshalerOptions {
	defaults := &defaultCBEUnmarshalerOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	_this.Decoder.ApplyDefaults()
	_this.Builder.ApplyDefaults()
	_this.Rules.ApplyDefaults()

	return _this
}

// ============================================================================

type CTEDecoderOptions struct {
	// TODO: ShouldZeroCopy option
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

var defaultCTEDecoderOptions = CTEDecoderOptions{
	BufferSize: 4096,
}

func DefaultCTEDecoderOptions() *CTEDecoderOptions {
	options := defaultCTEDecoderOptions
	return &options
}

func (_this *CTEDecoderOptions) ApplyDefaults() *CTEDecoderOptions {
	defaults := &defaultCTEDecoderOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	if _this.BufferSize < 64 {
		_this.BufferSize = 64
	}

	return _this
}

// ============================================================================

type CTEEncoderOptions struct {
	BufferSize int
	Indent     string
	// TODO: BinaryFloatEncoding option
	BinaryFloatEncoding BinaryFloatEncodeAs
}

var defaultCTEEncoderOptions = CTEEncoderOptions{
	BufferSize: 1024,
}

func DefaultCTEEncoderOptions() *CTEEncoderOptions {
	opts := defaultCTEEncoderOptions
	return &opts
}

func (_this *CTEEncoderOptions) ApplyDefaults() *CTEEncoderOptions {
	defaults := &defaultCTEEncoderOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	// TODO: Check for default individual options

	return _this
}

// ============================================================================

type CTEMarshalerOptions struct {
	Encoder  CTEEncoderOptions
	Iterator IteratorOptions
}

var defaultCTEMarshalerOptions = CTEMarshalerOptions{
	Encoder:  defaultCTEEncoderOptions,
	Iterator: defaultIteratorOptions,
}

func DefaultCTEMarshalerOptions() *CTEMarshalerOptions {
	options := defaultCTEMarshalerOptions
	return &options
}

func (_this *CTEMarshalerOptions) ApplyDefaults() *CTEMarshalerOptions {
	defaults := &defaultCTEMarshalerOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	_this.Encoder.ApplyDefaults()
	_this.Iterator.ApplyDefaults()

	return _this
}

// ============================================================================

type CTEUnmarshalerOptions struct {
	Decoder CTEDecoderOptions
	Builder BuilderOptions
	Rules   RuleOptions
}

var defaultCTEUnmarshalerOptions = CTEUnmarshalerOptions{
	Decoder: defaultCTEDecoderOptions,
	Builder: defaultBuilderOptions,
	Rules:   defaultRuleOptions,
}

func DefaultCTEUnmarshalerOptions() *CTEUnmarshalerOptions {
	options := defaultCTEUnmarshalerOptions
	return &options
}

func (_this *CTEUnmarshalerOptions) ApplyDefaults() *CTEUnmarshalerOptions {
	defaults := &defaultCTEUnmarshalerOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	_this.Decoder.ApplyDefaults()
	_this.Builder.ApplyDefaults()
	_this.Rules.ApplyDefaults()

	return _this
}

// ============================================================================

type BuilderOptions struct {
	FloatToBigIntMaxBase10Exponent int
	FloatToBigIntMaxBase2Exponent  int
	// TODO: ErrorOnLossyFloatConversion option
	ErrorOnLossyFloatConversion bool
	// TODO: Something for decimal floats?
	// TODO: Error on unknown field
}

var defaultBuilderOptions = BuilderOptions{
	FloatToBigIntMaxBase10Exponent: 300,
	FloatToBigIntMaxBase2Exponent:  300 * 10 / 3,
}

func DefaultBuilderOptions() *BuilderOptions {
	options := defaultBuilderOptions
	return &options
}

func (_this *BuilderOptions) ApplyDefaults() *BuilderOptions {
	defaults := &defaultBuilderOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	// TODO: Check for default individual options

	return _this
}

// ============================================================================

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

var defaultIteratorOptions = IteratorOptions{
	ConciseEncodingVersion: version.ConciseEncodingVersion,
}

func DefaultIteratorOptions() *IteratorOptions {
	opts := defaultIteratorOptions
	return &opts
}

func (_this *IteratorOptions) ApplyDefaults() *IteratorOptions {
	defaults := &defaultIteratorOptions
	if _this == nil {
		options := *defaults
		return &options
	}

	if _this.ConciseEncodingVersion < 1 {
		_this.ConciseEncodingVersion = defaults.ConciseEncodingVersion
	}

	return _this
}

// ============================================================================

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

var defaultRuleOptions = RuleOptions{
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

func DefaultRuleOptions() *RuleOptions {
	options := defaultRuleOptions
	return &options
}

func (_this *RuleOptions) ApplyDefaults() *RuleOptions {
	defaults := &defaultRuleOptions
	if _this == nil {
		options := *defaults
		return &options
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
