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

package ce

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
)

// Decoder decodes a byte stream, converting it to events.
// Its operation is similar to a lexer.
type Decoder interface {
	// Decode the stream of bytes from reader, sending all events to eventReceiver.
	Decode(reader io.Reader, eventReceiver events.DataEventReceiver) error

	// Decode from the specified document, sending all events to eventReceiver.
	DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) (err error)
}

// A universal decoder automatically chooses a decoding method based on the first byte of the document:
// - 0x63 = Decode as CTE
// - 0x81 = Decode as CBE
type UniversalDecoder struct {
	opts *options.CEDecoderOptions
}

func (_this *UniversalDecoder) Decode(reader io.Reader, eventReceiver events.DataEventReceiver) error {
	bufReader := bufio.NewReader(reader)
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		return err
	}

	if decoder, err := chooseDecoder(firstByte[0], _this.opts); err == nil {
		return decoder.Decode(bufReader, eventReceiver)
	} else {
		return err
	}
}

func (_this *UniversalDecoder) DecodeDocument(document []byte, eventReceiver events.DataEventReceiver) error {
	if decoder, err := chooseDecoder(document[0], _this.opts); err == nil {
		return decoder.DecodeDocument(document, eventReceiver)
	} else {
		return err
	}
}

func chooseDecoder(identifier byte, opts *options.CEDecoderOptions) (decoder Decoder, err error) {
	switch identifier {
	case 'c':
		decoder = cte.NewDecoder(opts)
	case cbe.CBESignatureByte:
		decoder = cbe.NewDecoder(opts)
	default:
		err = fmt.Errorf("%02d: Unknown CE identifier", identifier)
	}
	return
}
