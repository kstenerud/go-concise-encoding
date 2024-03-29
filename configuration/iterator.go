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

	// Every struct type added here will have a record type with the associated
	// name generated at the top of the file, and all instances will be generated
	// as records instead of maps.
	// Only go struct types are allowed.
	// Don't map multiple types to the same name unless you're sure you know what you're doing.
	RecordTypes map[reflect.Type]string

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

func (_this *IteratorConfiguration) init() {
	_this.CustomBinaryConverters = make(map[reflect.Type]ConvertToCustomFunction)
	_this.CustomTextConverters = make(map[reflect.Type]ConvertToCustomFunction)
	_this.RecordTypes = make(map[reflect.Type]string)
}

var defaultIteratorConfiguration = IteratorConfiguration{
	FieldNameStyle:           FieldNameSnakeCase,
	RecursionSupport:         false,
	DefaultFieldOmitBehavior: OmitFieldEmpty,
}

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
