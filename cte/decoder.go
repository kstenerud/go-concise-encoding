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

package cte

import (
	"fmt"
	"io"
	"strings"

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
)

type Decoder struct {
	config *configuration.Configuration
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver.
func NewDecoder(config *configuration.Configuration) *Decoder {
	_this := &Decoder{}
	_this.Init(config)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver.
func (_this *Decoder) Init(config *configuration.Configuration) {
	_this.config = config
}

func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !_this.config.Debug.PassThroughPanics {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, reader); err != nil {
		return
	}

	_this.markBytesRead(buf.Len())

	return ParseDocument(buf.String(), eventReceiver)
}

func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !_this.config.Debug.PassThroughPanics {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	_this.markBytesRead(len(document))

	return ParseDocument(string(document), eventReceiver)
}

func (_this *Decoder) markBytesRead(byteCount int) {
	if uint64(byteCount) > _this.config.Rules.MaxDocumentSizeBytes {
		_this.errorf("exceeded maximum document size of %v", _this.config.Rules.MaxDocumentSizeBytes)
	}
}

func (_this *Decoder) errorf(format string, args ...interface{}) {
	// TODO: Diagnostics
	panic(fmt.Errorf(format, args...))
}
