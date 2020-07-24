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

	"github.com/kstenerud/go-concise-encoding/version"
)

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
