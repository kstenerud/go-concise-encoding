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

package builder

import (
	"fmt"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

type customBuilder struct {
	// Template Data
	session *Session
}

func newCustomBuilder(session *Session) ObjectBuilder {
	return &customBuilder{
		session: session,
	}
}

func (_this *customBuilder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}

func (_this *customBuilder) panicBadEvent(name string, args ...interface{}) {
	PanicBadEvent(_this, name, args...)
}

func (_this *customBuilder) InitTemplate(_ *Session) {
	_this.panicBadEvent("InitTemplate")
}

func (_this *customBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}

func (_this *customBuilder) SetParent(_ ObjectBuilder) {
}

func (_this *customBuilder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	if !_this.session.TryBuildFromCustom(_this, arrayType, value, dst) {
		_this.panicBadEvent("TypedArray(%v)", arrayType)
	}
}
