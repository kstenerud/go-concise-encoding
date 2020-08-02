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

package cte

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/builder"
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/iterator"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
)

// ============================================================================
// Marshaler

// Marshaler is the top-level API for serializing objects. It maintains an
// iterator session so that cached iterator information is not lost between
// multiple calls to marshal.
type Marshaler struct {
	options options.CTEMarshalerOptions
	Session *iterator.Session
	encoder Encoder
}

// Create a new marshaler with the specified options.
// If opts is nil, default options will be used.
func NewMarshaler(opts *options.CTEMarshalerOptions) *Marshaler {
	_this := &Marshaler{}
	_this.Init(opts)
	return _this
}

// Init a marshaler with the specified options.
// If opts is nil, default options will be used.
func (_this *Marshaler) Init(opts *options.CTEMarshalerOptions) {
	opts = opts.WithDefaultsApplied()
	_this.options = *opts
	_this.Session = iterator.NewSession(nil, &_this.options.Session)
	_this.encoder.Init(&_this.options.Encoder)
}

// Marshal a go object into a CTE document, written to writer.
func (_this *Marshaler) Marshal(object interface{}, writer io.Writer) (err error) {
	defer func() {
		_this.encoder.Reset()
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	_this.encoder.PrepareToEncode(writer)
	iterator := _this.Session.NewIterator(&_this.encoder, &_this.options.Iterator)
	iterator.Iterate(object)
	return
}

// Marshal a go object into a CTE document, returning the document as a byte slice.
func (_this *Marshaler) MarshalToDocument(object interface{}) (document []byte, err error) {
	var buff bytes.Buffer
	err = _this.Marshal(object, &buff)
	document = buff.Bytes()
	return
}

// ============================================================================
// Unmarshaler

// Unmarshaler is the top-level API for deserializing objects. It maintains a
// builder session so that cached builder information is not lost between
// multiple calls to unmarshal.
type Unmarshaler struct {
	options options.CTEUnmarshalerOptions
	Session *builder.Session
	decoder Decoder
}

// Create a new unmarshaler with the specified options.
// If options is nil, default options will be used.
func NewUnmarshaler(opts *options.CTEUnmarshalerOptions) *Unmarshaler {
	_this := &Unmarshaler{}
	_this.Init(opts)
	return _this
}

// Init an unmarshaler with the specified options.
// If options is nil, default options will be used.
func (_this *Unmarshaler) Init(opts *options.CTEUnmarshalerOptions) {
	opts = opts.WithDefaultsApplied()
	_this.options = *opts
	_this.Session = builder.NewSession(nil, &_this.options.Session)
	_this.decoder.Init(&_this.options.Decoder)
}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
func (_this *Unmarshaler) Unmarshal(reader io.Reader, template interface{}) (decoded interface{}, err error) {
	defer func() {
		if !debug.DebugOptions.PassThroughPanics {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	builder := _this.Session.NewBuilderFor(template, &_this.options.Builder)
	rules := rules.NewRules(builder, &_this.options.Rules)
	if err = _this.decoder.Decode(reader, rules); err != nil {
		return
	}
	decoded = builder.GetBuiltObject()
	return
}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, an interface type will be returned.
func (_this *Unmarshaler) UnmarshalFromDocument(document []byte, template interface{}) (decoded interface{}, err error) {
	return _this.Unmarshal(bytes.NewBuffer(document), template)
}
