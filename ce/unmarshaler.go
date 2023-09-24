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
	"fmt"
	"io"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
)

// Unmarshaler decodes bytes from an input source, then creates, fills out, and
// returns an object that is compatible with the given template.
type Unmarshaler interface {
	// Unmarshal an object from the given reader, in a type compatible with template.
	// If template is nil, a best-guess type will be returned (likely a slice or map).
	Unmarshal(reader io.Reader, template interface{}) (decoded interface{}, err error)

	// Unmarshal an object from the given document, in a type compatible with template.
	// If template is nil, a best-guess type will be returned (likely a slice or map).
	UnmarshalFromDocument(document []byte, template interface{}) (decoded interface{}, err error)
}

func chooseUnmarshaler(identifier byte, config *configuration.Configuration) (unmarshaler Unmarshaler, err error) {
	switch identifier {
	case 'c':
		unmarshaler = cte.NewUnmarshaler(config)
	case cbe.CBESignatureByte:
		unmarshaler = cbe.NewUnmarshaler(config)
	default:
		err = fmt.Errorf("%02d: Unknown CE identifier", identifier)
	}
	return
}
