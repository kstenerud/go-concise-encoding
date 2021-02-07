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
	"bytes"
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/debug"

	"github.com/kstenerud/go-concise-encoding/internal/chars"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-concise-encoding/events"
)

type DecoderFunc func(ctx *DecoderContext)

type Decoder struct {
	opts options.CTEDecoderOptions
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func NewDecoder(opts *options.CTEDecoderOptions) *Decoder {
	_this := &Decoder{}
	_this.Init(opts)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If opts is nil, default options will be used.
func (_this *Decoder) Init(opts *options.CTEDecoderOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
}

func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !debug.DebugOptions.PassThroughPanics {
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

	ctx := DecoderContext{}
	ctx.Init(&_this.opts, reader, eventReceiver)
	ctx.StackDecoder(decodeDocumentBegin)

	for !ctx.IsDocumentComplete {
		ctx.DecodeNext()
	}
	return
}

func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	return _this.Decode(bytes.NewBuffer(document), eventReceiver)
}

var decoderFuncsByFirstChar [0x101]DecoderFunc

func init() {
	for i := 0; i < 0x100; i++ {
		decoderFuncsByFirstChar[i] = decodeInvalidChar
	}

	decoderFuncsByFirstChar['\r'] = decodeWhitespace
	decoderFuncsByFirstChar['\n'] = decodeWhitespace
	decoderFuncsByFirstChar['\t'] = decodeWhitespace
	decoderFuncsByFirstChar[' '] = decodeWhitespace
	decoderFuncsByFirstChar['"'] = decodeQuotedString
	decoderFuncsByFirstChar['_'] = decodeUnquotedString
	for i := 'A'; i <= 'Z'; i++ {
		decoderFuncsByFirstChar[i] = decodeUnquotedString
	}
	for i := 'a'; i <= 'z'; i++ {
		decoderFuncsByFirstChar[i] = decodeUnquotedString
	}
	for i := 0xc0; i < 0xf8; i++ {
		decoderFuncsByFirstChar[i] = decodeUnquotedString
	}
	decoderFuncsByFirstChar['0'] = decodeOtherBasePositive
	for i := '1'; i <= '9'; i++ {
		decoderFuncsByFirstChar[i] = decodePositiveNumeric
	}
	decoderFuncsByFirstChar['-'] = decodeNegativeNumeric
	decoderFuncsByFirstChar['@'] = decodeNamedValueOrUUID
	decoderFuncsByFirstChar['#'] = decodeConstant
	decoderFuncsByFirstChar['$'] = decodeReference
	decoderFuncsByFirstChar['&'] = decodeMarker
	decoderFuncsByFirstChar['/'] = decodeComment
	// TODO: ':'
	// decoderFuncsByFirstChar[':'] =
	decoderFuncsByFirstChar['{'] = decodeMapBegin
	decoderFuncsByFirstChar['}'] = decodeMapEnd
	decoderFuncsByFirstChar['['] = decodeListBegin
	decoderFuncsByFirstChar[']'] = decodeListEnd
	decoderFuncsByFirstChar['<'] = decodeMarkupBegin
	decoderFuncsByFirstChar[','] = decodeMarkupContentBegin
	decoderFuncsByFirstChar['>'] = decodeMarkupEnd
	decoderFuncsByFirstChar['('] = decodeMetadataBegin
	decoderFuncsByFirstChar[')'] = decodeMetadataEnd
	decoderFuncsByFirstChar['|'] = decodeTypedArrayBegin

	decoderFuncsByFirstChar[chars.EndOfDocumentMarker] = decodeInvalidChar
}
