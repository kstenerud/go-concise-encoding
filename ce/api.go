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
//   - name: (k=v) Specifies the name to use when encoding/decoding to a document.
//   - omit: (flag) Never read or write this field to/from a CE document.
//   - omit_never: (flag) Never omit this field.
//   - omit_empty: (flag) Omit only if it's empty (empty slice, array, map, nil value).
//   - omit_empty: (flag) Omit only if it's a golang zero value.
package ce

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/rules"
)

// ============================================================================
// One-shot universal unmarshal API

// Unmarshal a CE document (CBE or CTE) from a reader, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
func UnmarshalCE(reader io.Reader, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	bufReader := bufio.NewReader(reader)
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		return
	}

	unmarshaler, err := chooseUnmarshaler(firstByte[0], config)
	if err != nil {
		return
	}
	return unmarshaler.Unmarshal(bufReader, template)
}

// Unmarshal a CE document (CBE or CTE) from a byte slice, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
func UnmarshalFromCEDocument(document []byte, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	if len(document) == 0 {
		err = fmt.Errorf("no data")
		return
	}
	unmarshaler, err := chooseUnmarshaler(document[0], config)
	if err != nil {
		return
	}
	return unmarshaler.UnmarshalFromDocument(document, template)
}

// ============================================================================
// One-shot marshal/unmarshal API (binary format)

// Marshal a go object into a CBE document, written to writer.
func MarshalCBE(object interface{}, writer io.Writer, config *configuration.Configuration) (err error) {
	return NewCBEMarshaler(config).Marshal(object, writer)
}

// Marshal a go object into a CBE document, returned as a byte slice.
func MarshalToCBEDocument(object interface{}, config *configuration.Configuration) (document []byte, err error) {
	return NewCBEMarshaler(config).MarshalToDocument(object)
}

// Unmarshal a CBE document from a reader, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
// If an error occurs while decoding, err will be non-null, and decoded will contain whatever was successfully decoded thus far.
func UnmarshalCBE(reader io.Reader, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	return NewCBEUnmarshaler(config).Unmarshal(reader, template)
}

// Unmarshal a CBE document from a byte slice, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
// If an error occurs while decoding, err will be non-null, and decoded will contain whatever was successfully decoded thus far.
func UnmarshalFromCBEDocument(document []byte, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	return NewCBEUnmarshaler(config).UnmarshalFromDocument(document, template)
}

// ============================================================================
// One-shot marshal/unmarshal API (text format)

// Marshal a go object into a CTE document, written to writer.
func MarshalCTE(object interface{}, writer io.Writer, config *configuration.Configuration) (err error) {
	return NewCTEMarshaler(config).Marshal(object, writer)
}

// Marshal a go object into a CTE document, returned as a byte slice.
func MarshalToCTEDocument(object interface{}, config *configuration.Configuration) (document []byte, err error) {
	return NewCTEMarshaler(config).MarshalToDocument(object)

}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
// If an error occurs while decoding, err will be non-null, and decoded will contain whatever was successfully decoded thus far.
func UnmarshalCTE(reader io.Reader, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	return NewCTEUnmarshaler(config).Unmarshal(reader, template)
}

// Unmarshal a CTE document from a byte slice, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
// If an error occurs while decoding, err will be non-null, and decoded will contain whatever was successfully decoded thus far.
func UnmarshalFromCTEDocument(document []byte, template interface{}, config *configuration.Configuration) (decoded interface{}, err error) {
	return NewCTEUnmarshaler(config).UnmarshalFromDocument(document, template)
}

// ============================================================================
// Marshalers/Unmarshalers API

func NewCBEMarshaler(config *configuration.Configuration) Marshaler {
	return cbe.NewMarshaler(config)
}

func NewCBEUnmarshaler(config *configuration.Configuration) Unmarshaler {
	return cbe.NewUnmarshaler(config)
}

func NewCTEMarshaler(config *configuration.Configuration) Marshaler {
	return cte.NewMarshaler(config)
}

func NewCTEUnmarshaler(config *configuration.Configuration) Unmarshaler {
	return cte.NewUnmarshaler(config)
}

// ============================================================================
// Encoders/Decoders API

// Create a new universal CE decoder
func NewCEDecoder(config *configuration.Configuration) Decoder {
	return &UniversalDecoder{config: config}
}

func NewCBEEncoder(config *configuration.Configuration) Encoder {
	return cbe.NewEncoder(config)
}

func NewCBEDecoder(config *configuration.Configuration) Decoder {
	return cbe.NewDecoder(config)
}

func NewCTEEncoder(config *configuration.Configuration) Encoder {
	return cte.NewEncoder(config)
}

func NewCTEDecoder(config *configuration.Configuration) Decoder {
	return cte.NewDecoder(config)
}

// Create a new rules data receiver, which will enforce proper concise encoding structure.
func NewRules(nextReceiver events.DataEventReceiver, config *configuration.Configuration) *rules.RulesEventReceiver {
	return rules.NewRules(nextReceiver, config)
}
