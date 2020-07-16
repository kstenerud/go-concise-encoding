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

// TODO: Opt: Convert line endings to escapes
// TODO: Opt: Don't convert escapes
// TODO: Don't use marker/ref on empty containers
// TODO: Builder that converts to string
// TODO: Iterator that converts from string to smaller type (numeric)
// TODO: Some method to notify that a string field should be encoded as a different type
// TODO: Optional spaces around `=` in maps

type CTEMarshalerOptions struct {
	Encoder  CTEEncoderOptions
	Iterator IteratorOptions
}

type CTEUnmarshalerOptions struct {
	Decoder CTEDecoderOptions
	Builder BuilderOptions
	Rules   RuleOptions
	// TODO: Error on unknown field
}

type CBEMarshalerOptions struct {
	Encoder  CBEEncoderOptions
	Iterator IteratorOptions
}

type CBEUnmarshalerOptions struct {
	Decoder CBEDecoderOptions
	Builder BuilderOptions
	Rules   RuleOptions
	// TODO: Error on unknown field
}

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

type CBEDecoderOptions struct {
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

type CBEEncoderOptions struct {
	BufferSize int
}

// How to encode binary floats
type BinaryFloatEncodeAs int

const (
	// Use decimal encoding (1.2e4)
	BinaryFloatEncodeAsDecimal = iota
	// Use binary encoding (0x1.2p4)
	BinaryFloatEncodeAsBinary
)

// Where to place the opening brace on a struct, map, list
type BracePosition int

const (
	// Place the opening brace on the same line
	BracePositionAdjacent = iota
	// Place the opening brace on the next line
	BracePositionNextLine
)

type CTEEncoderOptions struct {
	BufferSize int
	Indent     string
	// TODO: BracePosition option
	BracePosition BracePosition
	// TODO: BinaryFloatEncoding option
	BinaryFloatEncoding BinaryFloatEncodeAs
}

type CTEDecoderOptions struct {
	// TODO: ShouldZeroCopy option
	ShouldZeroCopy bool
	// TODO: ImpliedVersion option
	ImpliedVersion uint
	// TODO: ImpliedTLContainer option
	ImpliedTLContainer TLContainerType
	BufferSize         int
}

type BuilderOptions struct {
	// TODO: Currently handled in bigIntMaxBase10Exponent in conversions.go
	FloatToBigIntMaxExponent int
	// TODO: ErrorOnLossyFloatConversion option
	ErrorOnLossyFloatConversion bool
	// TODO: Something for decimal floats?
}

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

type RuleOptions struct {
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
