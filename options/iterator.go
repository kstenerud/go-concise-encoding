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
	"reflect"
)

// ============================================================================
// Iterator Session

// Converts a value to custom binary data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-text
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-text
type ConvertToCustomFunction func(v reflect.Value) (customType uint64, asBytes []byte, err error)

type IteratorSessionOptions struct {

	// Use lowercase struct field names
	LowercaseStructFieldNames bool

	// Specifies which types to convert to custom binary data, and how to do it.
	// Note: You should only fill out one of these maps, depending on your
	// intended encoding (binary or text). The iterator session will consult
	// the binary map first and the text map second, choosing the first match.
	CustomBinaryConverters map[reflect.Type]ConvertToCustomFunction

	// Specifies which types to convert to custom text data, and how to do it
	// Note: You should only fill out one of these maps, depending on your
	// intended encoding (binary or text). The iterator session will consult
	// the binary map first and the text map second, choosing the first match.
	CustomTextConverters map[reflect.Type]ConvertToCustomFunction
}

func DefaultIteratorSessionOptions() IteratorSessionOptions {
	opts := defaultIteratorSessionOptions
	opts.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomFunction)
	opts.CustomTextConverters = make(map[reflect.Type]ConvertToCustomFunction)
	return opts
}

var defaultIteratorSessionOptions = IteratorSessionOptions{
	LowercaseStructFieldNames: true,
}

func (_this *IteratorSessionOptions) ApplyDefaults() {
	if _this.CustomBinaryConverters == nil {
		_this.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomFunction)
	}
	if _this.CustomTextConverters == nil {
		_this.CustomTextConverters = make(map[reflect.Type]ConvertToCustomFunction)
	}
}

func (_this *IteratorSessionOptions) Validate() error {
	return nil
}

// ============================================================================
// Iterator

type IteratorOptions struct {
	// If RecursionSupport is true, the iterator will also look for duplicate
	// pointers to data, generating marker and reference events rather than
	// walking the object again. This is useful for cyclic or recursive data
	// structures, but has a performance cost.
	RecursionSupport bool

	// TODO If true, don't write a nil object when a nil pointer is encountered.
	OmitNilPointers bool
}

func DefaultIteratorOptions() IteratorOptions {
	return defaultIteratorOptions
}

var defaultIteratorOptions = IteratorOptions{
	RecursionSupport: true,
	OmitNilPointers:  true,
}

func (_this *IteratorOptions) ApplyDefaults() {
	// Nothing to do
}

func (_this *IteratorOptions) Validate() error {
	return nil
}
