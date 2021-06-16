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
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/options"
)

type DecoderOp interface {
	Run(*DecoderContext)
}

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
	ctx.StackDecoder(global_decodeDocumentBegin)

	for !ctx.IsDocumentComplete {
		ctx.DecodeNext()
	}
	return
}

func (_this *Decoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error) {
	return _this.Decode(bytes.NewBuffer(document), eventReceiver)
}

var decoderOpsByFirstChar [0x101]DecoderOp

func init() {
	for i := 0; i < 0x100; i++ {
		decoderOpsByFirstChar[i] = global_decodeInvalidChar
	}

	decoderOpsByFirstChar['\r'] = global_decodeWhitespace
	decoderOpsByFirstChar['\n'] = global_decodeWhitespace
	decoderOpsByFirstChar['\t'] = global_decodeWhitespace
	decoderOpsByFirstChar[' '] = global_decodeWhitespace
	decoderOpsByFirstChar['"'] = global_advanceAndDecodeQuotedString
	decoderOpsByFirstChar['0'] = global_decode0Based
	for i := '1'; i <= '9'; i++ {
		decoderOpsByFirstChar[i] = global_decodeNumericPositive
	}
	for i := 'a'; i <= 'f'; i++ {
		decoderOpsByFirstChar[i] = global_decodeUID
	}
	for i := 'A'; i <= 'F'; i++ {
		decoderOpsByFirstChar[i] = global_decodeUID
	}
	decoderOpsByFirstChar['f'] = global_decodeFalseOrUID
	decoderOpsByFirstChar['F'] = global_decodeFalseOrUID
	decoderOpsByFirstChar['i'] = global_decodeNamedValueI
	decoderOpsByFirstChar['I'] = global_decodeNamedValueI
	decoderOpsByFirstChar['n'] = global_decodeNamedValueN
	decoderOpsByFirstChar['N'] = global_decodeNamedValueN
	decoderOpsByFirstChar['s'] = global_decodeNamedValueS
	decoderOpsByFirstChar['S'] = global_decodeNamedValueS
	decoderOpsByFirstChar['t'] = global_decodeNamedValueT
	decoderOpsByFirstChar['T'] = global_decodeNamedValueT
	decoderOpsByFirstChar['-'] = global_advanceAndDecodeNumericNegative
	decoderOpsByFirstChar['@'] = global_advanceAndDecodeResourceID
	decoderOpsByFirstChar['#'] = global_advanceAndDecodeConstant
	decoderOpsByFirstChar['$'] = global_advanceAndDecodeReference
	decoderOpsByFirstChar['&'] = global_advanceAndDecodeMarker
	decoderOpsByFirstChar['/'] = global_advanceAndDecodeComment
	decoderOpsByFirstChar[':'] = global_advanceAndDecodeSuffix
	decoderOpsByFirstChar['{'] = global_advanceAndDecodeMapBegin
	decoderOpsByFirstChar['}'] = global_advanceAndDecodeMapEnd
	decoderOpsByFirstChar['['] = global_advanceAndDecodeListBegin
	decoderOpsByFirstChar[']'] = global_advanceAndDecodeListEnd
	decoderOpsByFirstChar['<'] = global_advanceAndDecodeMarkupBegin
	decoderOpsByFirstChar[','] = global_advanceAndDecodeMarkupContentBegin
	decoderOpsByFirstChar['>'] = global_advanceAndDecodeMarkupEnd
	decoderOpsByFirstChar['*'] = global_advanceAndDecodeCommentEnd
	decoderOpsByFirstChar['|'] = global_advanceAndDecodeTypedArrayBegin
	decoderOpsByFirstChar['('] = global_advanceAndDecodeRelationshipBegin

	decoderOpsByFirstChar[chars.EOFMarker] = global_decodeInvalidChar
}
