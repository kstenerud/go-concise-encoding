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
	"fmt"
	"reflect"
)

// ============================================================================
// Builder Session

// Fills out a value from custom binary data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-binary
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-binary
type CustomBinaryBuildFunction func(src []byte, dst reflect.Value) error

// Fills out a value from custom text data.
// See https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md#custom-text
// See https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#custom-text
type CustomTextBuildFunction func(src string, dst reflect.Value) error

type BuilderSessionOptions struct {
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

func DefaultBuilderSessionOptions() *BuilderSessionOptions {
	return &BuilderSessionOptions{
		CustomBinaryBuildFunction: (func(src []byte, dst reflect.Value) error {
			return fmt.Errorf("No builder has been registered to handle custom binary data")
		}),
		CustomTextBuildFunction: (func(src string, dst reflect.Value) error {
			return fmt.Errorf("No builder has been registered to handle custom text data")
		}),
	}
}

func (_this *BuilderSessionOptions) WithDefaultsApplied() *BuilderSessionOptions {
	defaults := DefaultBuilderSessionOptions()
	if _this == nil {
		return defaults
	}

	if _this.CustomBinaryBuildFunction == nil {
		_this.CustomBinaryBuildFunction = defaults.CustomBinaryBuildFunction
	}
	if _this.CustomTextBuildFunction == nil {
		_this.CustomTextBuildFunction = defaults.CustomTextBuildFunction
	}
	if _this.CustomBuiltTypes == nil {
		_this.CustomBuiltTypes = []reflect.Type{}
	}

	return _this
}

// ============================================================================
// Builder

type BuilderOptions struct {
	FloatToBigIntMaxBase10Exponent int
	FloatToBigIntMaxBase2Exponent  int
	// TODO: ErrorOnLossyFloatConversion option
	ErrorOnLossyFloatConversion bool
	// TODO: Something for decimal floats?
	// TODO: Error on unknown field
}

func DefaultBuilderOptions() *BuilderOptions {
	return &BuilderOptions{
		FloatToBigIntMaxBase10Exponent: 300,
		FloatToBigIntMaxBase2Exponent:  300 * 10 / 3,
	}
}

func (_this *BuilderOptions) WithDefaultsApplied() *BuilderOptions {
	if _this == nil {
		return DefaultBuilderOptions()
	}

	// TODO: Check for default individual options

	return _this
}