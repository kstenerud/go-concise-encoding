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
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kstenerud/go-concise-encoding/builder"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/iterator"
	"github.com/kstenerud/go-concise-encoding/nullevent"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
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

func benchmarkMarshal(b *testing.B, marshaler ce.Marshaler) {
	b.Helper()
	data := generate()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		o := data[i%len(data)]
		bytes, err := marshaler.MarshalToDocument(o)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(o))
		}
		serialSize += len(bytes)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func BenchmarkCTEMarshal(b *testing.B) {
	config := configuration.DefaultCTEMarshalerConfiguration()
	config.Iterator.RecursionSupport = false
	marshaler := ce.NewCTEMarshaler(&config)
	benchmarkMarshal(b, marshaler)
}

func BenchmarkCBEMarshal(b *testing.B) {
	config := configuration.DefaultCBEMarshalerConfiguration()
	config.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(&config)
	benchmarkMarshal(b, marshaler)
}

func BenchmarkJSONMarshal(b *testing.B) {
	b.Helper()
	data := generate()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	for i := 0; i < b.N; i++ {
		o := data[i%len(data)]
		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)
		err := enc.Encode(o)
		bytes := buff.Bytes()
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(o))
		}
		serialSize += len(bytes)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func benchmarkUnmarshal(b *testing.B, marshaler ce.Marshaler, unmarshaler ce.Unmarshaler) {
	b.Helper()
	expectedObjs := generate()
	actualObjs := make([]*A, len(expectedObjs))
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
		index := i % len(expectedObjs)
		document := documents[index]
		obj, err := unmarshaler.UnmarshalFromDocument(document, template)
		if err != nil {
			b.Fatalf("Unmarshal error: %s (while decoding [%v])", err, describe.D(document))
		}
		actualObjs[index] = obj.(*A)
	}
	b.StopTimer()
	for i, v := range actualObjs {
		if v != nil {
			if !equivalence.IsEquivalent(v, expectedObjs[i]) {
				b.Fatalf("Expected [%v] to produce %v but got %v", describe.D(documents[i]), describe.D(expectedObjs[i]), describe.D(v))
			}
		}
	}
}

func benchmarkDecode(b *testing.B, marshaler ce.Marshaler, decoder ce.Decoder) {
	b.Helper()
	expectedObjs := generate()
	documents := make([][]byte, 0, len(expectedObjs))
	for _, obj := range expectedObjs {
		bytes, err := marshaler.MarshalToDocument(obj)
		if err != nil {
			b.Fatalf("Marshal error: %s (while encoding %v)", err, describe.D(obj))
		}
		documents = append(documents, bytes)
	}
	nullReceiver := nullevent.NewNullEventReceiver()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % len(expectedObjs)
		document := documents[index]
		err := decoder.Decode(bytes.NewBuffer(document), nullReceiver)
		if err != nil {
			b.Fatalf("Unmarshal error: %s (while decoding [%v])", err, describe.D(document))
		}
	}
	b.StopTimer()
}

func BenchmarkCTEDecode(b *testing.B) {
	marshalConfig := configuration.DefaultCTEMarshalerConfiguration()
	marshalConfig.Iterator.RecursionSupport = false
	marshaler := ce.NewCTEMarshaler(&marshalConfig)
	decoder := ce.NewCTEDecoder(nil)
	benchmarkDecode(b, marshaler, decoder)
}

func BenchmarkCTEUnmarshalRules(b *testing.B) {
	marshalConfig := configuration.DefaultCTEMarshalerConfiguration()
	marshalConfig.Iterator.RecursionSupport = false
	marshaler := ce.NewCTEMarshaler(&marshalConfig)
	unmarshalConfig := configuration.DefaultCEUnmarshalerConfiguration()
	unmarshaler := ce.NewCTEUnmarshaler(&unmarshalConfig)
	benchmarkUnmarshal(b, marshaler, unmarshaler)
}

func BenchmarkCTEUnmarshalNoRules(b *testing.B) {
	marshalConfig := configuration.DefaultCTEMarshalerConfiguration()
	marshalConfig.Iterator.RecursionSupport = false
	marshaler := ce.NewCTEMarshaler(&marshalConfig)
	unmarshalConfig := configuration.DefaultCEUnmarshalerConfiguration()
	unmarshalConfig.EnforceRules = false
	unmarshaler := ce.NewCTEUnmarshaler(&unmarshalConfig)
	benchmarkUnmarshal(b, marshaler, unmarshaler)
}

func BenchmarkCBEUnmarshalRules(b *testing.B) {
	marshalConfig := configuration.DefaultCBEMarshalerConfiguration()
	marshalConfig.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(&marshalConfig)
	unmarshalConfig := configuration.DefaultCEUnmarshalerConfiguration()
	unmarshaler := ce.NewCBEUnmarshaler(&unmarshalConfig)
	benchmarkUnmarshal(b, marshaler, unmarshaler)
}

func BenchmarkCBEUnmarshalNoRules(b *testing.B) {
	marshalConfig := configuration.DefaultCBEMarshalerConfiguration()
	marshalConfig.Iterator.RecursionSupport = false
	marshaler := ce.NewCBEMarshaler(&marshalConfig)
	unmarshalConfig := configuration.DefaultCEUnmarshalerConfiguration()
	unmarshalConfig.EnforceRules = false
	unmarshaler := ce.NewCBEUnmarshaler(&unmarshalConfig)
	benchmarkUnmarshal(b, marshaler, unmarshaler)
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	b.Helper()
	expectedObjs := generate()
	actualObjs := make([]*A, len(expectedObjs))
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
		index := i % len(expectedObjs)
		document := documents[index]
		obj := &A{}
		decoder := json.NewDecoder(bytes.NewBuffer(document))
		err := decoder.Decode(obj)
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

func BenchmarkRules(b *testing.B) {
	b.Helper()
	receiver, store := test.NewEventCollector(nil)
	iterConfiguration := configuration.DefaultIteratorConfiguration()
	iterConfiguration.RecursionSupport = false
	iterSession := iterator.NewSession(nil, &iterConfiguration)
	iter := iterSession.NewIterator(receiver)

	objs := generate()
	documents := make([]test.Events, 0, len(objs))
	for _, obj := range objs {
		iter.Iterate(obj)
		documents = append(documents, store.Events)
		store.Clear()
	}

	b.ReportAllocs()
	b.ResetTimer()
	r := rules.NewRules(nullevent.NewNullEventReceiver(), nil)
	for i := 0; i < b.N; i++ {
		index := i % len(objs)
		r.Reset()
		test.InvokeEventsAsCompleteDocument(r, documents[index]...)
	}
	b.StopTimer()
}

func BenchmarkBuilder(b *testing.B) {
	b.Helper()
	receiver, store := test.NewEventCollector(nil)
	iterConfiguration := configuration.DefaultIteratorConfiguration()
	iterConfiguration.RecursionSupport = false
	iterSession := iterator.NewSession(nil, &iterConfiguration)
	iter := iterSession.NewIterator(receiver)

	objs := generate()
	documents := make([]test.Events, 0, len(objs))
	for _, obj := range objs {
		iter.Iterate(obj)
		documents = append(documents, store.Events)
	}

	b.ReportAllocs()
	b.ResetTimer()
	template := &A{}
	builderSession := builder.NewSession(nil, nil)
	for i := 0; i < b.N; i++ {
		index := i % len(objs)
		builder := builderSession.NewBuilderFor(template, nil)
		test.InvokeEventsAsCompleteDocument(builder, documents[index]...)
	}
	b.StopTimer()
}

func BenchmarkIterator(b *testing.B) {
	b.Helper()
	iterConfiguration := configuration.DefaultIteratorConfiguration()
	iterSession := iterator.NewSession(nil, &iterConfiguration)
	objs := generate()

	b.ReportAllocs()
	b.ResetTimer()
	iterConfiguration.RecursionSupport = false
	iter := iterSession.NewIterator(nullevent.NewNullEventReceiver())
	for i := 0; i < b.N; i++ {
		index := i % len(objs)
		iter.Iterate(objs[index])
	}
	b.StopTimer()
}
