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

type DocumentLimits struct {
	MaxDocumentSizeBytes          uint64
	MaxArraySizeBytes             uint64
	MaxObjectCount                uint64
	MaxContainerDepth             uint64
	MaxIntegerDigitCount          uint64
	MaxFloatCoefficientDigitCount uint64
	MaxFloatExponentDigitCount    uint64
	MaxYearDigitCount             uint64
	MaxMarkerCount                uint64
	MaxLocalReferenceCount        uint64
}

var defaultDocumentLimits = DocumentLimits{
	MaxDocumentSizeBytes:          5368709120,
	MaxArraySizeBytes:             1073741824,
	MaxObjectCount:                1000000,
	MaxContainerDepth:             1000,
	MaxIntegerDigitCount:          100,
	MaxFloatCoefficientDigitCount: 100,
	MaxFloatExponentDigitCount:    5,
	MaxYearDigitCount:             11,
	MaxMarkerCount:                10000,
	MaxLocalReferenceCount:        10000,
}

func (_this *DocumentLimits) ApplyDefaults() {
	if _this.MaxDocumentSizeBytes == 0 {
		_this.MaxDocumentSizeBytes = defaultCEDecoderConfiguration.DocumentLimits.MaxDocumentSizeBytes
	}
	if _this.MaxArraySizeBytes == 0 {
		_this.MaxArraySizeBytes = defaultCEDecoderConfiguration.DocumentLimits.MaxArraySizeBytes
	}
	if _this.MaxObjectCount == 0 {
		_this.MaxObjectCount = defaultCEDecoderConfiguration.DocumentLimits.MaxObjectCount
	}
	if _this.MaxContainerDepth == 0 {
		_this.MaxContainerDepth = defaultCEDecoderConfiguration.DocumentLimits.MaxContainerDepth
	}
	if _this.MaxIntegerDigitCount == 0 {
		_this.MaxIntegerDigitCount = defaultCEDecoderConfiguration.DocumentLimits.MaxIntegerDigitCount
	}
	if _this.MaxFloatCoefficientDigitCount == 0 {
		_this.MaxFloatCoefficientDigitCount = defaultCEDecoderConfiguration.DocumentLimits.MaxFloatCoefficientDigitCount
	}
	if _this.MaxFloatExponentDigitCount == 0 {
		_this.MaxFloatExponentDigitCount = defaultCEDecoderConfiguration.DocumentLimits.MaxFloatExponentDigitCount
	}
	if _this.MaxYearDigitCount == 0 {
		_this.MaxYearDigitCount = defaultCEDecoderConfiguration.DocumentLimits.MaxYearDigitCount
	}
	if _this.MaxMarkerCount == 0 {
		_this.MaxMarkerCount = defaultCEDecoderConfiguration.DocumentLimits.MaxMarkerCount
	}
	if _this.MaxLocalReferenceCount == 0 {
		_this.MaxLocalReferenceCount = defaultCEDecoderConfiguration.DocumentLimits.MaxLocalReferenceCount
	}
}

// ============================================================================
// CE Decoder

type CEDecoderConfiguration struct {
	AllowRecursiveLocalReferences bool
	FollowRemoteReferences        bool
	CompleteTruncatedDocument     bool
	AllowNulCharacter             bool
	DocumentLimits                DocumentLimits
	DebugPanics                   bool
}

var defaultCEDecoderConfiguration = CEDecoderConfiguration{
	AllowRecursiveLocalReferences: false,
	FollowRemoteReferences:        false,
	CompleteTruncatedDocument:     false,
	AllowNulCharacter:             false,
	DocumentLimits:                defaultDocumentLimits,
	DebugPanics:                   false,
}

func DefaultCEDecoderConfiguration() CEDecoderConfiguration {
	return defaultCEDecoderConfiguration
}

func (_this *CEDecoderConfiguration) ApplyDefaults() {
	_this.DocumentLimits.ApplyDefaults()
}

func (_this *CEDecoderConfiguration) Validate() error {
	return nil
}

// ============================================================================
// CE Unmarshaler

type CEUnmarshalerConfiguration struct {
	Decoder CEDecoderConfiguration
	Builder BuilderConfiguration
	Session BuilderSessionConfiguration
	Rules   RuleConfiguration

	// If false, do not wrap a Rules object around the builder, disabling all rule checks.
	EnforceRules bool
	DebugPanics  bool
}

func DefaultCEUnmarshalerConfiguration() CEUnmarshalerConfiguration {
	return defaultCEUnmarshalerConfiguration
}

var defaultCEUnmarshalerConfiguration = CEUnmarshalerConfiguration{
	Decoder:      DefaultCEDecoderConfiguration(),
	Builder:      DefaultBuilderConfiguration(),
	Session:      DefaultBuilderSessionConfiguration(),
	Rules:        DefaultRuleConfiguration(),
	EnforceRules: true,
}

func (_this *CEUnmarshalerConfiguration) ApplyDefaults() {
	_this.Decoder.ApplyDefaults()
	_this.Builder.ApplyDefaults()
	_this.Session.ApplyDefaults()
	_this.Rules.ApplyDefaults()
	_this.Decoder.DebugPanics = _this.DebugPanics
}

func (_this *CEUnmarshalerConfiguration) Validate() error {
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
