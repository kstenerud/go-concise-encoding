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
	"github.com/kstenerud/go-concise-encoding/version"
)

// ============================================================================
// CE Decoder

type CEDecoderOptions struct {
	// Concise encoding spec version to adhere to. Uses latest if set to 0.
	ConciseEncodingVersion uint64
}

func DefaultCEDecoderOptions() *CEDecoderOptions {
	return &CEDecoderOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
	}
}

func (_this *CEDecoderOptions) WithDefaultsApplied() *CEDecoderOptions {
	if _this == nil {
		return DefaultCEDecoderOptions()
	}

	if _this.ConciseEncodingVersion == 0 {
		_this.ConciseEncodingVersion = version.ConciseEncodingVersion
	}

	return _this
}

func (_this *CEDecoderOptions) Validate() error {
	return nil
}

// ============================================================================
// CE Unmarshaler

type CEUnmarshalerOptions struct {
	Decoder CEDecoderOptions
	Builder BuilderOptions
	Session BuilderSessionOptions
	Rules   RuleOptions

	// If false, do not wrap a Rules object around the builder, disabling all rule checks.
	EnforceRules bool
}

func DefaultCEUnmarshalerOptions() *CEUnmarshalerOptions {
	return &CEUnmarshalerOptions{
		Decoder:      *DefaultCEDecoderOptions(),
		Builder:      *DefaultBuilderOptions(),
		Session:      *DefaultBuilderSessionOptions(),
		Rules:        *DefaultRuleOptions(),
		EnforceRules: true,
	}
}

func (_this *CEUnmarshalerOptions) WithDefaultsApplied() *CEUnmarshalerOptions {
	if _this == nil {
		return DefaultCEUnmarshalerOptions()
	}

	_this.Decoder.WithDefaultsApplied()
	_this.Builder.WithDefaultsApplied()
	_this.Session.WithDefaultsApplied()
	_this.Rules.WithDefaultsApplied()

	return _this
}

func (_this *CEUnmarshalerOptions) Validate() error {
	if err := _this.Builder.Validate(); err != nil {
		return err
	}
	if err := _this.Decoder.Validate(); err != nil {
		return err
	}
	if err := _this.Rules.Validate(); err != nil {
		return err
	}
	return _this.Session.Validate()
}
