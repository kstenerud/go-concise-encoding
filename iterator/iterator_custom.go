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
	"github.com/kstenerud/go-concise-encoding/options"
)

type customBinaryIterator struct {
	convertFunction options.ConvertToCustomBinaryFunction
}

func newCustomBinaryIterator(convertFunction options.ConvertToCustomBinaryFunction) ObjectIterator {
	return &customBinaryIterator{
		convertFunction: convertFunction,
	}
}

func (_this *customBinaryIterator) PostCacheInitIterator(session *Session) {
}

func (_this *customBinaryIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	asBytes, err := _this.convertFunction(v)
	if err != nil {
		panic(fmt.Errorf("Error converting type %v to custom bytes: %v", v.Type(), err))
	}
	eventReceiver.OnCustomBinary(asBytes)
}

type customTextIterator struct {
	convertFunction options.ConvertToCustomTextFunction
}

func newCustomTextIterator(convertFunction options.ConvertToCustomTextFunction) ObjectIterator {
	return &customTextIterator{
		convertFunction: convertFunction,
	}
}

func (_this *customTextIterator) PostCacheInitIterator(session *Session) {
}

func (_this *customTextIterator) IterateObject(v reflect.Value, eventReceiver events.DataEventReceiver, references ReferenceEventGenerator) {
	asString, err := _this.convertFunction(v)
	if err != nil {
		panic(fmt.Errorf("Error converting type %v to custom text: %v", v.Type(), err))
	}
	eventReceiver.OnCustomText(asString)
}
