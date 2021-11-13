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
	"reflect"

	"github.com/kstenerud/go-concise-encoding/events"
)

// Common function signatures
type GetIteratorForType func(reflect.Type) IteratorFunction
type TryAddReference func(reflect.Value) (didGenerateReferenceEvent bool)

type Context struct {
	// Per-session data
	GetIteratorForType        GetIteratorForType
	LowercaseStructFieldNames bool

	// Per-root-iterator data
	EventReceiver   events.DataEventReceiver
	TryAddReference TryAddReference
}

func (_this *Context) NotifyNil() {
	_this.EventReceiver.OnNull()
}

func sessionContext(getIteratorFunc GetIteratorForType, lowercaseStructFieldNames bool) Context {
	return Context{
		GetIteratorForType:        getIteratorFunc,
		LowercaseStructFieldNames: lowercaseStructFieldNames,
	}
}

func iteratorContext(sessionContext *Context,
	eventReceiver events.DataEventReceiver,
	tryAddReference TryAddReference) Context {

	return Context{
		GetIteratorForType:        sessionContext.GetIteratorForType,
		LowercaseStructFieldNames: sessionContext.LowercaseStructFieldNames,
		EventReceiver:             eventReceiver,
		TryAddReference:           tryAddReference,
	}
}
