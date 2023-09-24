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
	"fmt"
	"reflect"
)

type BuilderConfiguration struct {
	// Max base-10 exponent allowed when converting from floating point to big integer.
	// As exponents get very large, it takes geometrically more CPU to convert.
	FloatToBigIntMaxBase10Exponent int

	// Max base-2 exponent allowed when converting from floating point to big integer.
	FloatToBigIntMaxBase2Exponent int

	// Match struct field names in a case insensitive manner
	CaseInsensitiveStructFieldNames bool

	// TODO: If true, don't raise an error on a lossy floating point conversion.
	AllowLossyFloatConversion bool

	// TODO: If true, don't raise an error on unknown fields
	IgnoreUnknownFields bool

	// Specifies which types will be built using custom text/binary build
	// functions. You must also set one or both of CustomBinaryBuildFunction
	// and CustomTextBuildFunction in order to use this feature.
	// Both CBE and CTE will attempt to use either the binary or text version
	// depending on the data type (custom binary, custom text) encoded in the
	// source document.
	CustomBuiltTypes []reflect.Type

	// Build function to use when building from a custom binary source.
	CustomBinaryBuildFunction CustomBinaryBuildFunction

	// Build function to use when building from a custom text source.
	CustomTextBuildFunction CustomTextBuildFunction
}

func (_this *BuilderConfiguration) init() {}

var defaultBuilderConfiguration = BuilderConfiguration{
	FloatToBigIntMaxBase10Exponent:  maxBase10Exp,
	FloatToBigIntMaxBase2Exponent:   maxBase10Exp * 10 / 3,
	AllowLossyFloatConversion:       true,
	IgnoreUnknownFields:             true,
	CaseInsensitiveStructFieldNames: true,
	CustomBinaryBuildFunction: func(customType uint64, src []byte, dst reflect.Value) error {
		return fmt.Errorf("no builder has been registered to handle custom binary data")
	},
	CustomTextBuildFunction: func(customType uint64, src string, dst reflect.Value) error {
		return fmt.Errorf("no builder has been registered to handle custom text data")
	},
	CustomBuiltTypes: []reflect.Type{},
}

// Fills out a value from custom data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-text
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-text
type CustomBinaryBuildFunction func(customType uint64, src []byte, dst reflect.Value) error
type CustomTextBuildFunction func(customType uint64, src string, dst reflect.Value) error

const maxBase10Exp = 50
