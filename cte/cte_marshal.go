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
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/builder"
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/iterator"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
)

var defaultMarshalerOptions = options.CTEMarshalerOptions{
	Encoder: defaultEncoderOptions,
}

func DefaultMarshalerOptions() *options.CTEMarshalerOptions {
	opts := defaultMarshalerOptions
	return &opts
}

// Marshal a go object into a CTE document, written to writer.
// If options is nil, default options will be used.
func Marshal(object interface{}, writer io.Writer, opts *options.CTEMarshalerOptions) (err error) {
	if opts == nil {
		opts = &options.CTEMarshalerOptions{}
	}
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

	encoder := NewEncoder(writer, &opts.Encoder)
	iterator.IterateObject(object, encoder, &opts.Iterator)
	return
}

var defaultUnmarshalerOptions = options.CTEUnmarshalerOptions{
	Decoder: defaultDecoderOptions,
}

func DefaultUnmarshalerOptions() *options.CTEUnmarshalerOptions {
	opts := defaultUnmarshalerOptions
	return &opts
}

// Unmarshal a CTE document, creating an object of the same type as the template.
// If options is nil, default options will be used.
func Unmarshal(reader io.Reader, template interface{}, opts *options.CTEUnmarshalerOptions) (decoded interface{}, err error) {
	if opts == nil {
		opts = &options.CTEUnmarshalerOptions{}
	}
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

	builder := builder.NewBuilderFor(template, &opts.Builder)
	rules := rules.NewRules(&opts.Rules, builder)
	decoder := NewDecoder(reader, rules, &opts.Decoder)
	decoder.Decode()
	decoded = builder.GetBuiltObject()
	return
}
