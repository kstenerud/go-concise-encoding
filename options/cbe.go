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

// ============================================================================
// CBE Encoder

type CBEEncoderOptions struct {
}

func DefaultCBEEncoderOptions() CBEEncoderOptions {
	return defaultCBEEncoderOptions
}

var defaultCBEEncoderOptions = CBEEncoderOptions{}

func (_this *CBEEncoderOptions) ApplyDefaults() {
}

func (_this *CBEEncoderOptions) Validate() error {
	return nil
}

// ============================================================================
// CBE Marshaler

type CBEMarshalerOptions struct {
	Encoder  CBEEncoderOptions
	Iterator IteratorOptions
	Session  IteratorSessionOptions
}

func DefaultCBEMarshalerOptions() CBEMarshalerOptions {
	return defaultCBEMarshalerOptions
}

var defaultCBEMarshalerOptions = CBEMarshalerOptions{
	Encoder:  DefaultCBEEncoderOptions(),
	Iterator: DefaultIteratorOptions(),
	Session:  DefaultIteratorSessionOptions(),
}

func (_this *CBEMarshalerOptions) ApplyDefaults() {
	_this.Encoder.ApplyDefaults()
	_this.Iterator.ApplyDefaults()
	_this.Session.ApplyDefaults()
}

func (_this *CBEMarshalerOptions) Validate() error {
	if err := _this.Encoder.Validate(); err != nil {
		return err
	}
	if err := _this.Iterator.Validate(); err != nil {
		return err
	}
	return _this.Session.Validate()
}
