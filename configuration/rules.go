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
	MaxDocumentSizeBytes          uint64
	MaxArraySizeBytes             uint64
	MaxIdentifierLength           uint64
	MaxObjectCount                uint64
	MaxContainerDepth             uint64
	MaxIntegerDigitCount          uint64
	MaxFloatCoefficientDigitCount uint64
	MaxFloatExponentDigitCount    uint64
	MaxYearDigitCount             uint64
	MaxMarkerCount                uint64
	MaxLocalReferenceCount        uint64
	AllowRecursiveLocalReferences bool
	AutoCompleteTruncatedDocument bool
}

func (_this *RuleConfiguration) init() {}

var defaultRuleConfiguration = RuleConfiguration{
	MaxDocumentSizeBytes:          5 * gigabyte,
	MaxArraySizeBytes:             1 * gigabyte,
	MaxIdentifierLength:           1000,
	MaxObjectCount:                1000000,
	MaxContainerDepth:             1000,
	MaxIntegerDigitCount:          100,
	MaxFloatCoefficientDigitCount: 100,
	MaxFloatExponentDigitCount:    5,
	MaxYearDigitCount:             11,
	MaxMarkerCount:                10000,
	MaxLocalReferenceCount:        10000,
	AllowRecursiveLocalReferences: false,
	AutoCompleteTruncatedDocument: false,
}

const gigabyte = 1024 * 1024 * 1024
