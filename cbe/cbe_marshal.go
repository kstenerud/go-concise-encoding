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

package cbe

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

// Marshal a go object into a CBE document, written to writer.
// If options is nil, default options will be used.
func Marshal(object interface{}, writer io.Writer, opts *options.CBEMarshalerOptions) (err error) {
	var marshaler Marshaler
	marshaler.Init(opts, nil)
	return marshaler.Marshal(object, writer)
}

// Marshal a go object into a CBE document, returned as a byte slice.
// If options is nil, default options will be used.
func MarshalToBytes(object interface{}, opts *options.CBEMarshalerOptions) (document []byte, err error) {
	var marshaler Marshaler
	marshaler.Init(opts, nil)
	return marshaler.MarshalToBytes(object)
}

// Unmarshal a CBE document, creating an object of the same type as the template.
// If options is nil, default options will be used.
func Unmarshal(reader io.Reader, template interface{}, opts *options.CBEUnmarshalerOptions) (decoded interface{}, err error) {
	var marshaler Marshaler
	marshaler.Init(nil, opts)
	return marshaler.Unmarshal(reader, template)
}

// Unmarshal CBE from a byte slice, creating an object of the same type as the template.
// If options is nil, default options will be used.
func UnmarshalFromBytes(document []byte, template interface{}, opts *options.CBEUnmarshalerOptions) (decoded interface{}, err error) {
	var marshaler Marshaler
	marshaler.Init(nil, opts)
	return marshaler.UnmarshalFromBytes(document, template)
}

// A marshaler keeps builder and iterator sessions so that cached builder &
// iterator information is not lost between multiple calls to marshal/unmarshal.
//
// If all you want is a one-off marshal or unmarshal, use the standalone functions.
type Marshaler struct {
	BuilderSession  *builder.Session
	IteratorSession *iterator.Session
	MarshalOpts     options.CBEMarshalerOptions
	UnmarshalOpts   options.CBEUnmarshalerOptions
}

// Create a new marshaler with the specified options.
// If options is nil, default options will be used.
func NewMarshaler(marshalOpts *options.CBEMarshalerOptions, unmarshalOpts *options.CBEUnmarshalerOptions) *Marshaler {
	_this := &Marshaler{}
	_this.Init(marshalOpts, unmarshalOpts)
	return _this
}

// Init a marshaler with the specified options.
// If options is nil, default options will be used.
func (_this *Marshaler) Init(marshalOpts *options.CBEMarshalerOptions, unmarshalOpts *options.CBEUnmarshalerOptions) {
	_this.MarshalOpts = *marshalOpts.WithDefaultsApplied()
	_this.UnmarshalOpts = *unmarshalOpts.WithDefaultsApplied()
	_this.BuilderSession = builder.NewSession(nil, &_this.UnmarshalOpts.Session)
	_this.IteratorSession = iterator.NewSession(nil, &_this.MarshalOpts.Session)
}

// Marshal a go object into a CBE document, written to writer.
func (_this *Marshaler) Marshal(object interface{}, writer io.Writer) (err error) {
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

	encoder := NewEncoder(writer, &_this.MarshalOpts.Encoder)
	iterator := _this.IteratorSession.NewIterator(encoder, &_this.MarshalOpts.Iterator)
	iterator.Iterate(object)
	return
}

// Marshal a go object into a CBE document, returning the document as a byte slice.
func (_this *Marshaler) MarshalToBytes(object interface{}) (document []byte, err error) {
	var buff bytes.Buffer
	err = _this.Marshal(object, &buff)
	document = buff.Bytes()
	return
}

// Unmarshal a CBE document, creating an object of the same type as the template.
func (_this *Marshaler) Unmarshal(reader io.Reader, template interface{}) (decoded interface{}, err error) {
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

	builder := _this.BuilderSession.NewBuilderFor(template, &_this.UnmarshalOpts.Builder)
	rules := rules.NewRules(builder, &_this.UnmarshalOpts.Rules)
	decoder := NewDecoder(reader, rules, &_this.UnmarshalOpts.Decoder)
	if err = decoder.Decode(); err != nil {
		return
	}
	decoded = builder.GetBuiltObject()
	return
}

// Unmarshal a CBE document, creating an object of the same type as the template.
func (_this *Marshaler) UnmarshalFromBytes(document []byte, template interface{}) (decoded interface{}, err error) {
	return _this.Unmarshal(bytes.NewBuffer(document), template)
}
