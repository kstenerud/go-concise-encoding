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
	"bytes"
	"strings"
	"testing"

	"github.com/kstenerud/go-concise-encoding/builder"
	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

func TestCEUnmarshalCTE(t *testing.T) {
	document := "c0 [1 2 3]"
	expected := []interface{}{1, 2, 3}
	actual, err := UnmarshalFromCEDocument([]byte(document), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}

	stream := strings.NewReader(document)
	actual, err = UnmarshalCE(stream, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func TestCEDecodeCTE(t *testing.T) {
	document := "c1 [1 2 3]"
	expected := []interface{}{1, 2, 3}

	decoder := NewCEDecoder(nil)

	builderSession := builder.NewSession(nil, nil)
	builder := builder.NewBuilderEventReceiver(builderSession, nil, nil)
	err := decoder.DecodeDocument([]byte(document), builder)
	if err != nil {
		t.Fatal(err)
	}
	actual := builder.GetBuiltObject()
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}

	stream := strings.NewReader(document)
	err = decoder.Decode(stream, builder)
	if err != nil {
		t.Fatal(err)
	}
	actual = builder.GetBuiltObject()
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func TestCEUnmarshalCBE(t *testing.T) {
	document := []byte{0x81, 0x00, 0x9a, 0x01, 0x02, 0x03, 0x9b}
	expected := []interface{}{1, 2, 3}
	actual, err := UnmarshalFromCEDocument([]byte(document), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !equivalence.IsEquivalent(expected, actual) {
		t.Errorf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}

	stream := bytes.NewBuffer(document)
	actual, err = UnmarshalCE(stream, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}

func TestCEDecodeCBE(t *testing.T) {
	document := []byte{0x81, 0x00, 0x9a, 0x01, 0x02, 0x03, 0x9b}
	expected := []interface{}{1, 2, 3}

	decoder := NewCEDecoder(nil)

	builderSession := builder.NewSession(nil, nil)
	builder := builder.NewBuilderEventReceiver(builderSession, nil, nil)
	err := decoder.DecodeDocument([]byte(document), builder)
	if err != nil {
		t.Fatal(err)
	}
	actual := builder.GetBuiltObject()
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}

	stream := bytes.NewBuffer(document)
	err = decoder.Decode(stream, builder)
	if err != nil {
		t.Fatal(err)
	}
	actual = builder.GetBuiltObject()
	if !equivalence.IsEquivalent(expected, actual) {
		t.Fatalf("Expected %v but got %v", describe.D(expected), describe.D(actual))
	}
}
