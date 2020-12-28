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
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

type A struct {
	Name     string
	Birthday time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func randString(l int) string {
	buf := make([]byte, l)
	for i := 0; i < (l+1)/2; i++ {
		buf[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", buf)[:l]
}

func generate() []*A {
	a := make([]*A, 0, 1000)
	for i := 0; i < 1000; i++ {
		a = append(a, &A{
			Name:     randString(16),
			Birthday: time.Now().Truncate(-1),
			Phone:    randString(10),
			Siblings: rand.Intn(5),
			Spouse:   rand.Intn(2) == 1,
			Money:    rand.Float64(),
		})
	}
	return a
}

func BenchmarkCBEMarshal(b *testing.B) {
	b.Helper()
	opts := options.DefaultCBEMarshalerOptions()
	opts.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(opts)
	data := generate()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		o := data[rand.Intn(len(data))]
		bytes, err := marshaler.MarshalToDocument(o)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(o))
		}
		serialSize += len(bytes)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func BenchmarkJSONMarshal(b *testing.B) {
	b.Helper()
	data := generate()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		o := data[rand.Intn(len(data))]
		bytes, err := json.Marshal(o)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(o))
		}
		serialSize += len(bytes)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func BenchmarkCBEUnmarshal(b *testing.B) {
	b.Helper()
	marshalOpts := options.DefaultCBEMarshalerOptions()
	marshalOpts.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(marshalOpts)
	unmarshalOpts := options.DefaultCBEUnmarshalerOptions()
	unmarshalOpts.EnforceRules = false
	unmarshaler := ce.NewCBEUnmarshaler(unmarshalOpts)
	expectedObjs := generate()
	actualObjs := make([]*A, len(expectedObjs), len(expectedObjs))
	documents := make([][]byte, 0, len(expectedObjs))
	for _, obj := range expectedObjs {
		bytes, err := marshaler.MarshalToDocument(obj)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(obj))
		}
		documents = append(documents, bytes)
	}
	b.ReportAllocs()
	b.ResetTimer()
	template := &A{}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(expectedObjs))
		document := documents[index]
		obj, err := unmarshaler.UnmarshalFromDocument(document, template)
		if err != nil {
			b.Fatalf("Unmarshal error: %s (while decoding %v)", err, describe.D(document))
		}
		actualObjs[index] = obj.(*A)
	}
	b.StopTimer()
	for i, v := range actualObjs {
		if v != nil {
			if !equivalence.IsEquivalent(v, expectedObjs[i]) {
				b.Fatalf("Expected %v to produce %v but got %v", describe.D(documents[i]), describe.D(expectedObjs[i]), describe.D(v))
			}
		}
	}
}

func BenchmarkRules(b *testing.B) {
	b.Helper()
	marshalOpts := options.DefaultCBEMarshalerOptions()
	marshalOpts.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(marshalOpts)
	unmarshalOpts := options.DefaultCBEUnmarshalerOptions()
	unmarshaler := ce.NewCBEUnmarshaler(unmarshalOpts)
	expectedObjs := generate()
	actualObjs := make([]*A, len(expectedObjs), len(expectedObjs))
	documents := make([][]byte, 0, len(expectedObjs))
	for _, obj := range expectedObjs {
		bytes, err := marshaler.MarshalToDocument(obj)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(obj))
		}
		documents = append(documents, bytes)
	}
	b.ReportAllocs()
	b.ResetTimer()
	template := &A{}
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(expectedObjs))
		document := documents[index]
		obj, err := unmarshaler.UnmarshalFromDocument(document, template)
		if err != nil {
			b.Fatalf("Unmarshal error: %s (while decoding %v)", err, describe.D(document))
		}
		actualObjs[index] = obj.(*A)
	}
	b.StopTimer()
	for i, v := range actualObjs {
		if v != nil {
			if !equivalence.IsEquivalent(v, expectedObjs[i]) {
				b.Fatalf("Expected %v to produce %v but got %v", describe.D(documents[i]), describe.D(expectedObjs[i]), describe.D(v))
			}
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	b.Helper()
	expectedObjs := generate()
	actualObjs := make([]*A, len(expectedObjs), len(expectedObjs))
	documents := make([][]byte, 0, len(expectedObjs))
	for _, obj := range expectedObjs {
		bytes, err := json.Marshal(obj)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(obj))
		}
		documents = append(documents, bytes)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(expectedObjs))
		document := documents[index]
		obj := &A{}
		err := json.Unmarshal(document, obj)
		if err != nil {
			b.Fatalf("Unmarshal error: %s (while decoding %v)", err, describe.D(document))
		}
		actualObjs[index] = obj
	}
	b.StopTimer()
	for i, v := range actualObjs {
		if v != nil {
			if !equivalence.IsEquivalent(v, expectedObjs[i]) {
				b.Fatalf("Expected %v to produce %v but got %v", describe.D(documents[i]), describe.D(expectedObjs[i]), describe.D(v))
			}
		}
	}
}
