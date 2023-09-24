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

// All configuration that can be used to fine-tune the behavior of various aspects of
// this library.
package configuration

func New() *Configuration {
	config := defaultConfiguration
	config.init()
	return &config
}

type Configuration struct {
	Rules    RuleConfiguration
	Iterator IteratorConfiguration
	Builder  BuilderConfiguration
	Decoder  DecoderConfiguration
	Encoder  EncoderConfiguration
	Marshal  MarshalConfiguration
	Debug    DebugConfiguration
}

func (_this *Configuration) init() {
	_this.Rules.init()
	_this.Iterator.init()
	_this.Builder.init()
	_this.Decoder.init()
	_this.Encoder.init()
	_this.Marshal.init()
	_this.Debug.init()
}

var defaultConfiguration = Configuration{
	Rules:    defaultRuleConfiguration,
	Iterator: defaultIteratorConfiguration,
	Builder:  defaultBuilderConfiguration,
	Decoder:  defaultDecoderConfiguration,
	Encoder:  defaultEncoderConfiguration,
	Marshal:  defaultMarshalConfiguration,
	Debug:    defaultDebugConfiguration,
}
