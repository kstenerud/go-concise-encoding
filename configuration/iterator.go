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

package configuration

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

type FieldNameStyle int

const (
	FieldNameCamelCase FieldNameStyle = iota
	FieldNameSnakeCase
)

type IteratorConfiguration struct {
	FieldNameStyle FieldNameStyle

	// If RecursionSupport is true, the iterator will also look for duplicate
	// pointers to data, generating marker and reference events rather than
	// walking the object again. This is useful for cyclic or recursive data
	// structures, but has a performance cost.
	RecursionSupport bool

	// What to do by default when an empty field is encountered.
	// This can be overridden at the field level using the ce tags
	// "omit", "omit_never", "omit_empty", and "omit_zero".
	//
	// Defaults to OmitFieldEmpty
	DefaultFieldOmitBehavior FieldOmitBehavior

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

func DefaultIteratorConfiguration() IteratorConfiguration {
	config := defaultIteratorConfiguration
	config.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomFunction)
	config.CustomTextConverters = make(map[reflect.Type]ConvertToCustomFunction)
	return config
}

var defaultIteratorConfiguration = IteratorConfiguration{
	FieldNameStyle:           FieldNameSnakeCase,
	RecursionSupport:         true,
	DefaultFieldOmitBehavior: OmitFieldEmpty,
}

func (_this *IteratorConfiguration) ApplyDefaults() {
	if _this.CustomBinaryConverters == nil {
		_this.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomFunction)
	}
	if _this.CustomTextConverters == nil {
		_this.CustomTextConverters = make(map[reflect.Type]ConvertToCustomFunction)
	}
}

func (_this *IteratorConfiguration) Validate() error {
	return nil
}

type FieldOmitBehavior int

const (
	OmitFieldChooseDefault FieldOmitBehavior = iota

	OmitFieldNever

	OmitFieldAlways

	// An "empty" field is:
	//  * A container with no contents
	//  * An empty array
	//  * A nil pointer
	OmitFieldEmpty

	// A zero field ebcompasses OmitFieldEmpty and also omits
	// any golang zero value.
	OmitFieldZero
)
