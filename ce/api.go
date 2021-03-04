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

// Package ce contains the main API for concise encoding, with interfaces for
// the top two API layers:
//
// * Top layer: marshal/unmarshal
// * Next layer: encode/decode
//
// Layers lower than this are probably not what you're looking for, but they
// are described in ../docs.go
//
// Tags:
// Concise Encoding supports struct tags, using the identifier "ce". Tag
// contents are comma separated flags and key-value pairs:
//
// ce:"flag1,key1=value1,flag2,flag3,key2=value2"
//
// Note: Whitespace will be trimmed, so "flag1 , key1 = value1 " will work.
//
// The following tag values are recognized:
//
// omit: (flag) This struct field will not be written to a CE document, nor
//       will it be read from a CE document.
//
// -:    Shorthand for omit.
//
// name: (k=v) Specifies the name to use when encoding/decoding to a document.
//
package ce

import (
	"io"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/kstenerud/go-concise-encoding/rules"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/options"
)

// ============================================================================
// One-shot marshal/unmarshal API (binary format)

// Marshal a go object into a CBE document, written to writer.
// If opts is nil, default options will be used.
func MarshalCBE(object interface{}, writer io.Writer, opts *options.CBEMarshalerOptions) (err error) {
	return NewCBEMarshaler(opts).Marshal(object, writer)
}

// Marshal a go object into a CBE document, returned as a byte slice.
// If opts is nil, default options will be used.
func MarshalCBEToDocument(object interface{}, opts *options.CBEMarshalerOptions) (document []byte, err error) {
	return NewCBEMarshaler(opts).MarshalToDocument(object)
}

// Unmarshal a CBE document from a reader, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
// If opts is nil, default options will be used.
func UnmarshalCBE(reader io.Reader, template interface{}, opts *options.CBEUnmarshalerOptions) (decoded interface{}, err error) {
	return NewCBEUnmarshaler(opts).Unmarshal(reader, template)
}

// Unmarshal a CBE document from a byte slice, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
// If opts is nil, default options will be used.
func UnmarshalCBEFromDocument(document []byte, template interface{}, opts *options.CBEUnmarshalerOptions) (decoded interface{}, err error) {
	return NewCBEUnmarshaler(opts).UnmarshalFromDocument(document, template)
}

// ============================================================================
// One-shot marshal/unmarshal API (text format)

// Marshal a go object into a CTE document, written to writer.
// If opts is nil, default options will be used.
func MarshalCTE(object interface{}, writer io.Writer, opts *options.CTEMarshalerOptions) (err error) {
	return NewCTEMarshaler(opts).Marshal(object, writer)
}

// Marshal a go object into a CTE document, returned as a byte slice.
// If opts is nil, default options will be used.
func MarshalCTEToDocument(object interface{}, opts *options.CTEMarshalerOptions) (document []byte, err error) {
	return NewCTEMarshaler(opts).MarshalToDocument(object)

}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
// If opts is nil, default options will be used.
func UnmarshalCTE(reader io.Reader, template interface{}, opts *options.CTEUnmarshalerOptions) (decoded interface{}, err error) {
	return NewCTEUnmarshaler(opts).Unmarshal(reader, template)
}

// Unmarshal a CTE document from a byte slice, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
// If opts is nil, default options will be used.
func UnmarshalCTEFromDocument(document []byte, template interface{}, opts *options.CTEUnmarshalerOptions) (decoded interface{}, err error) {
	return NewCTEUnmarshaler(opts).UnmarshalFromDocument(document, template)
}

// ============================================================================
// Marshalers/Unmarshalers API

func NewCBEMarshaler(opts *options.CBEMarshalerOptions) Marshaler {
	return cbe.NewMarshaler(opts)
}

func NewCBEUnmarshaler(opts *options.CBEUnmarshalerOptions) Unmarshaler {
	return cbe.NewUnmarshaler(opts)
}

func NewCTEMarshaler(opts *options.CTEMarshalerOptions) Marshaler {
	return cte.NewMarshaler(opts)
}

func NewCTEUnmarshaler(opts *options.CTEUnmarshalerOptions) Unmarshaler {
	return cte.NewUnmarshaler(opts)
}

// ============================================================================
// Encoders/Decoders API

func NewCBEEncoder(opts *options.CBEEncoderOptions) Encoder {
	return cbe.NewEncoder(opts)
}

func NewCBEDecoder(opts *options.CBEDecoderOptions) Decoder {
	return cbe.NewDecoder(opts)
}

func NewCTEEncoder(opts *options.CTEEncoderOptions) Encoder {
	return cte.NewEncoder(opts)
}

func NewCTEDecoder(opts *options.CTEDecoderOptions) Decoder {
	return cte.NewDecoder(opts)
}

// Create a new rules data receiver, which will enforce proper concise encoding structure.
func NewRules(nextReceiver events.DataEventReceiver, opts *options.RuleOptions) *rules.RulesEventReceiver {
	return rules.NewRules(nextReceiver, opts)
}
