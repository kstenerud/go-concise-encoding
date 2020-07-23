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

package iterator

import (
	"fmt"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
)

type customIterator struct {
	convertFunction ConvertToCustomBytesFunction
}

func newCustomIterator(convertFunction ConvertToCustomBytesFunction) ObjectIterator {
	return &customIterator{
		convertFunction: convertFunction,
	}
}

func (_this *customIterator) PostCacheInitIterator(session *Session) {
}

func (_this *customIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	asBytes, err := _this.convertFunction(v)
	if err != nil {
		panic(fmt.Errorf("Error converting type %v to custom bytes: %v", v.Type(), err))
	}
	eventReceiver.OnCustomBinary(asBytes)
}
