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
// Rules

type RuleOptions struct {
	depthHasBeenAdjusted bool
	// Concise encoding spec version to adhere to
	ConciseEncodingVersion uint64

	// Limits before the ruleset artificially terminates with an error.
	MaxBytesLength      uint64
	MaxStringLength     uint64
	MaxURILength        uint64
	MaxIDLength         uint64
	MaxMarkupNameLength uint64
	MaxContainerDepth   uint64
	MaxObjectCount      uint64
	MaxReferenceCount   uint64
	// Max bytes total for all array types
}

func DefaultRuleOptions() *RuleOptions {
	return &RuleOptions{
		ConciseEncodingVersion: version.ConciseEncodingVersion,
		MaxBytesLength:         1000000000,
		MaxStringLength:        100000000,
		MaxURILength:           10000,
		MaxIDLength:            100,
		MaxMarkupNameLength:    100,
		MaxContainerDepth:      1000,
		MaxObjectCount:         10000000,
		MaxReferenceCount:      100000,
		// TODO: References need to check for amplification attacks. Keep count of referenced things and their object counts
	}
}

func (_this *RuleOptions) WithDefaultsApplied() *RuleOptions {
	defaults := DefaultRuleOptions()
	if _this == nil {
		return defaults
	}

	if _this.ConciseEncodingVersion < 1 {
		_this.ConciseEncodingVersion = defaults.ConciseEncodingVersion
	}
	if _this.MaxBytesLength < 1 {
		_this.MaxBytesLength = defaults.MaxBytesLength
	}
	if _this.MaxStringLength < 1 {
		_this.MaxStringLength = defaults.MaxStringLength
	}
	if _this.MaxURILength < 1 {
		_this.MaxURILength = defaults.MaxURILength
	}
	if _this.MaxIDLength < 1 {
		_this.MaxIDLength = defaults.MaxIDLength
	}
	if _this.MaxMarkupNameLength < 1 {
		_this.MaxMarkupNameLength = defaults.MaxMarkupNameLength
	}
	if _this.MaxContainerDepth < 1 {
		_this.MaxContainerDepth = defaults.MaxContainerDepth
	}
	if _this.MaxObjectCount < 1 {
		_this.MaxObjectCount = defaults.MaxObjectCount
	}
	if _this.MaxReferenceCount < 1 {
		_this.MaxReferenceCount = defaults.MaxReferenceCount
	}

	return _this
}
