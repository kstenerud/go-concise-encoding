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
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/version"
)

func demonstrateCTEMarshal() error {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCTE(dict, &buffer, configuration.New())
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return err
	}
	fmt.Printf("Marshaled CTE: %v\n", buffer.String())
	// Prints: Marshaled CTE: c0 {"a key"=0x1.4p+01 900=2020-01-15/13:41:00.000599}

	return nil
}

func demonstrateCTEUnmarshal() error {
	data := bytes.NewBuffer([]byte(`c0 {"a key"=2.5 900=2020-01-15/13:41:00.000599}`))

	value, err := ce.UnmarshalCTE(data, nil, configuration.New())
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return err
	}
	fmt.Printf("Unmarshaled CTE: %v\n", value)
	// Prints: Unmarshaled CTE: map[a key:2.5 900:2020-01-15/13:41:00.000599]

	return nil
}

func demonstrateCBEMarshal() error {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCBE(dict, &buffer, configuration.New())
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return err
	}
	fmt.Printf("Marshaled CBE: %v\n", buffer.Bytes())
	// Prints: Marshaled CBE: [129 0 153 106 132 3 124 188 18 0 32 109 47 80 0 133 97 32 107 101 121 112 32 64 155]

	return nil
}

const ceVer = version.ConciseEncodingVersion

func demonstrateCBEUnmarshal() error {
	data := bytes.NewBuffer([]byte{0x81, ceVer, 0x99, 0x85, 0x61, 0x20, 0x6b, 0x65,
		0x79, 0x70, 0x20, 0x40, 0x6a, 0x84, 0x03, 0x7c, 0xbc, 0x12,
		0x00, 0x20, 0x6d, 0x2f, 0x50, 0x00, 0x9b})

	value, err := ce.UnmarshalCBE(data, nil, configuration.New())
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		return err
	}
	fmt.Printf("Unmarshaled CBE: %v\n", value)
	// Prints: Unmarshaled CBE: map[a key:2.5 900:2020-01-15/13:41:00.000599]

	return nil
}

type SomeStruct struct {
	A int
	B string
	C *SomeStruct
}

func demonstrateRecursiveStructInMap() error {
	config := configuration.New()
	config.Iterator.RecursionSupport = true

	document := `c0 {"my-value" = &1:{"a"=100 "b"="test" "c"=$1}}`
	template := map[string]*SomeStruct{}
	result, err := ce.UnmarshalFromCTEDocument([]byte(document), template, config)
	if err != nil {
		fmt.Printf("Error unmarshaling CTE document: %v\n", err)
		return err
	}
	v := result.(map[string]*SomeStruct)
	s := v["my-value"]
	// Can't naively print a recursive structure in go (Printf will stack overflow), so we print each piece manually.
	fmt.Printf("A: %v, B: %v, Ptr to C: %p, ptr to s: %p\n", s.A, s.B, s.C, s)
	// Prints: A: 100, B: test, Ptr to C: 0xc0001f4600, ptr to s: 0xc0001f4600

	encodedDocument, err := ce.MarshalToCTEDocument(v, config)
	if err != nil {
		fmt.Printf("Error marshaling CTE document: %v\n", err)
		return err
	}
	fmt.Printf("Re-encoded CTE: %v\n", string(encodedDocument))
	// Prints: Re-encoded CTE: c0 {my-value=&0:{A=100 B=test C=$0}}

	encodedDocument, err = ce.MarshalToCBEDocument(v, config)
	if err != nil {
		fmt.Printf("Error marshaling CBE document: %v\n", err)
		return err
	}
	fmt.Printf("Re-encoded CBE: %v\n", encodedDocument)
	// Prints: Re-encoded CBE: [129 0 153 136 109 121 45 118 97 108 117 101 127 240 1 48 153 129 97 100 129 98 132 116 101 115 116 129 99 124 1 48 155 155]

	return nil
}

func TestDemonstrate(t *testing.T) {
	if err := demonstrateCTEMarshal(); err != nil {
		t.Error(err)
	}
	if err := demonstrateCTEUnmarshal(); err != nil {
		t.Error(err)
	}
	if err := demonstrateCBEMarshal(); err != nil {
		t.Error(err)
	}
	if err := demonstrateCBEUnmarshal(); err != nil {
		t.Error(err)
	}
	if err := demonstrateRecursiveStructInMap(); err != nil {
		t.Error(err)
	}
}
