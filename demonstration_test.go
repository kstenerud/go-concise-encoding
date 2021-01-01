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

package concise_encoding

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/ce"
)

func demonstrateCTEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCTE(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CTE: %v\n", string(buffer.Bytes()))
	// Prints: Marshaled CTE: c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}
}

func demonstrateCTEUnmarshal() {
	data := bytes.NewBuffer([]byte(`c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}`))

	value, err := ce.UnmarshalCTE(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CTE: %v\n", value)
	// Prints: Unmarshaled CTE: map[a key:2.5 900:2020-01-15/13:41:00.000599]
}

func demonstrateCBEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCBE(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CBE: %v\n", buffer.Bytes())
	// Prints: Marshaled CBE: [3 1 121 133 97 32 107 101 121 112 32 64 106 132 3 155 188 18 0 32 109 47 80 0 123]
}

func demonstrateCBEUnmarshal() {
	data := bytes.NewBuffer([]byte{0x03, 0x01, 0x79, 0x85, 0x61, 0x20, 0x6b, 0x65,
		0x79, 0x70, 0x20, 0x40, 0x6a, 0x84, 0x03, 0x9b, 0xbc, 0x12,
		0x00, 0x20, 0x6d, 0x2f, 0x50, 0x00, 0x7b})

	value, err := ce.UnmarshalCBE(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CBE: %v\n", value)
	// Prints: Unmarshaled CBE: map[a key:2.5 900:2020-01-15/13:41:00.000599]
}

func TestDemonstrate(t *testing.T) {
	demonstrateCTEMarshal()
	demonstrateCTEUnmarshal()
	demonstrateCBEMarshal()
	demonstrateCBEUnmarshal()
}
