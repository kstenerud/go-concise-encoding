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
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/iterator"
	"github.com/kstenerud/go-concise-encoding/rules"
)

// ============================================================================
// Marshaler

// Marshaler is the top-level API for serializing objects. It maintains an
// iterator session so that cached iterator information is not lost between
// multiple calls to marshal.
type Marshaler struct {
	session iterator.Session
	encoder EncoderEventReceiver
	config  *configuration.Configuration
}

// Create a new marshaler with the specified configuration.
func NewMarshaler(config *configuration.Configuration) *Marshaler {
	_this := &Marshaler{}
	_this.Init(config)
	return _this
}

// Init a marshaler with the specified configuration.
func (_this *Marshaler) Init(config *configuration.Configuration) {
	_this.config = config
	_this.session.Init(nil, _this.config)
	_this.encoder.Init(_this.config)
}

// Marshal a go object into a CTE document, written to writer.
func (_this *Marshaler) Marshal(object interface{}, writer io.Writer) (err error) {
	if !_this.config.Debug.PassThroughPanics {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}()
	}

	_this.encoder.PrepareToEncode(writer)
	iterator := _this.session.NewIterator(&_this.encoder)
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
	session builder.Session
	decoder Decoder
	config  *configuration.Configuration
	rules   rules.RulesEventReceiver
}

// Create a new unmarshaler with the specified configuration.
func NewUnmarshaler(config *configuration.Configuration) *Unmarshaler {
	_this := &Unmarshaler{}
	_this.Init(config)
	return _this
}

// Init an unmarshaler with the specified configuration.
func (_this *Unmarshaler) Init(config *configuration.Configuration) {
	_this.config = config
	_this.session.Init(nil, _this.config)
	_this.decoder.Init(_this.config)
	_this.rules.Init(nil, _this.config)
}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
func (_this *Unmarshaler) Unmarshal(reader io.Reader, template interface{}) (decoded interface{}, err error) {
	if !_this.config.Debug.PassThroughPanics {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}()
	}

	builder := _this.session.NewBuilderFor(template)
	receiver := events.DataEventReceiver(builder)
	if _this.config.Marshal.EnforceRules {
		_this.rules.Reset()
		_this.rules.SetNextReceiver(receiver)
		receiver = &_this.rules
	}
	if err = _this.decoder.Decode(reader, receiver); err != nil {
		return
	}
	decoded = builder.GetBuiltObject()
	return
}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If template is nil, a best-guess type will be returned (likely a slice or map).
func (_this *Unmarshaler) UnmarshalFromDocument(document []byte, template interface{}) (decoded interface{}, err error) {
	return _this.Unmarshal(bytes.NewBuffer(document), template)
}
