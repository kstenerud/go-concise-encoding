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

	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/internal/chars"
)

type DecoderOp func(*DecoderContext)

type Decoder struct {
	config *configuration.CEDecoderConfiguration
}

// Create a new CTE decoder, which will read from reader and send data events
// to nextReceiver. If config is nil, default configuration will be used.
func NewDecoder(config *configuration.CEDecoderConfiguration) *Decoder {
	_this := &Decoder{}
	_this.Init(config)
	return _this
}

// Initialize this decoder, which will read from reader and send data events
// to nextReceiver. If config is nil, default configuration will be used.
func (_this *Decoder) Init(config *configuration.CEDecoderConfiguration) {
	if config == nil {
		defaultConfig := configuration.DefaultCEDecoderConfiguration()
		config = &defaultConfig
	} else {
		config.ApplyDefaults()
	}
	_this.config = config
}

func (_this *Decoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) (err error) {
	defer func() {
		if !_this.config.DebugPanics {
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
	ctx.Init(_this.config, reader, eventReceiver)
	ctx.StackDecoder(decodeDocumentBegin)

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
		decoderOpsByFirstChar[i] = decodeInvalidChar
	}

	decoderOpsByFirstChar['\r'] = decodeWhitespace
	decoderOpsByFirstChar['\n'] = decodeWhitespace
	decoderOpsByFirstChar['\t'] = decodeWhitespace
	decoderOpsByFirstChar[' '] = decodeWhitespace
	decoderOpsByFirstChar['"'] = advanceAndDecodeQuotedString
	decoderOpsByFirstChar['0'] = decode0Based
	for i := '1'; i <= '9'; i++ {
		decoderOpsByFirstChar[i] = decodeNumericPositive
	}
	for i := 'a'; i <= 'f'; i++ {
		decoderOpsByFirstChar[i] = decodeUID
	}
	for i := 'A'; i <= 'F'; i++ {
		decoderOpsByFirstChar[i] = decodeUID
	}
	decoderOpsByFirstChar['f'] = decodeFalseOrUID
	decoderOpsByFirstChar['F'] = decodeFalseOrUID
	decoderOpsByFirstChar['i'] = decodeNamedValueI
	decoderOpsByFirstChar['I'] = decodeNamedValueI
	decoderOpsByFirstChar['n'] = decodeNamedValueN
	decoderOpsByFirstChar['N'] = decodeNamedValueN
	decoderOpsByFirstChar['s'] = decodeNamedValueS
	decoderOpsByFirstChar['S'] = decodeNamedValueS
	decoderOpsByFirstChar['t'] = decodeNamedValueT
	decoderOpsByFirstChar['T'] = decodeNamedValueT
	decoderOpsByFirstChar['-'] = advanceAndDecodeNumericNegative
	decoderOpsByFirstChar['@'] = advanceAndDecodeAt
	decoderOpsByFirstChar['$'] = advanceAndDecodeLocalReference
	decoderOpsByFirstChar['&'] = advanceAndDecodeMarker
	decoderOpsByFirstChar['/'] = advanceAndDecodeComment
	decoderOpsByFirstChar['{'] = advanceAndDecodeMapBegin
	decoderOpsByFirstChar['}'] = advanceAndDecodeMapEnd
	decoderOpsByFirstChar['['] = advanceAndDecodeListBegin
	decoderOpsByFirstChar[']'] = advanceAndDecodeListEnd
	decoderOpsByFirstChar['|'] = advanceAndDecodeTypedArrayBegin
	decoderOpsByFirstChar['('] = advanceAndDecodeNodeBegin
	decoderOpsByFirstChar[')'] = advanceAndDecodeNodeOrEdgeOrStructTemplateEnd
	decoderOpsByFirstChar['>'] = advanceAndDecodeStructTemplateEnd

	decoderOpsByFirstChar[chars.EOFMarker] = decodeInvalidChar
}
