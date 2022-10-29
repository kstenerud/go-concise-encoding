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

// ============================================================================
// Rules

type RuleConfiguration struct {
	// Limits before the ruleset artificially terminates with an error.
	MaxArrayByteLength      uint64
	MaxStringByteLength     uint64
	MaxResourceIDByteLength uint64
	MaxIdentifierLength     uint64
	MaxContainerDepth       uint64
	MaxObjectCount          uint64
	MaxLocalReferenceCount  uint64

	// TODO: Max bytes total for all array types
	MaxTotalArrayBytes uint64
}

func DefaultRuleConfiguration() RuleConfiguration {
	return defaultRuleConfiguration
}

var defaultRuleConfiguration = RuleConfiguration{
	MaxArrayByteLength:      1000000000,
	MaxStringByteLength:     100000000,
	MaxResourceIDByteLength: 10000,
	MaxIdentifierLength:     1000,
	MaxContainerDepth:       1000,
	MaxObjectCount:          10000000,
	MaxLocalReferenceCount:  100000,
	// TODO: References need to check for amplification attacks. Keep count of referenced things and their object counts
}

func (_this *RuleConfiguration) ApplyDefaults() {
	if _this.MaxArrayByteLength < 1 {
		_this.MaxArrayByteLength = defaultRuleConfiguration.MaxArrayByteLength
	}
	if _this.MaxStringByteLength < 1 {
		_this.MaxStringByteLength = defaultRuleConfiguration.MaxStringByteLength
	}
	if _this.MaxStringByteLength > _this.MaxArrayByteLength {
		_this.MaxStringByteLength = _this.MaxArrayByteLength
	}
	if _this.MaxResourceIDByteLength < 1 {
		_this.MaxResourceIDByteLength = defaultRuleConfiguration.MaxResourceIDByteLength
	}
	if _this.MaxResourceIDByteLength > _this.MaxArrayByteLength {
		_this.MaxResourceIDByteLength = _this.MaxArrayByteLength
	}
	if _this.MaxContainerDepth < 1 {
		_this.MaxContainerDepth = defaultRuleConfiguration.MaxContainerDepth
	}
	if _this.MaxObjectCount < 1 {
		_this.MaxObjectCount = defaultRuleConfiguration.MaxObjectCount
	}
	if _this.MaxLocalReferenceCount < 1 {
		_this.MaxLocalReferenceCount = defaultRuleConfiguration.MaxLocalReferenceCount
	}
}

func (_this *RuleConfiguration) Validate() error {
	return nil
}
