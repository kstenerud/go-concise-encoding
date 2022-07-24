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

package test_runner

import (
	"bytes"
	"fmt"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"
)

type CBETestRunner struct {
	Name         string
	Encode       []*CBEEncodeTest
	Decode       []*CBEDecodeTest
	DecodeFail   []CBEDecodeFailTest
	EncodeDecode []*CBEEncodeDecodeTest
	DecodeEncode []*CBEDecodeEncodeTest
	Skip         bool
	Debug        bool
	Trace        bool
}

func (_this *CBETestRunner) String() string {
	return _this.Name
}

func (_this *CBETestRunner) postDecodeInit() {
	for _, t := range _this.Encode {
		t.postDecodeInit()
		t.debug = _this.Debug
		t.trace = _this.Trace
	}

	for _, t := range _this.Decode {
		t.postDecodeInit()
		t.debug = _this.Debug
		t.trace = _this.Trace
	}

	for _, t := range _this.EncodeDecode {
		t.postDecodeInit()
		t.debug = _this.Debug
		t.trace = _this.Trace
	}

	for _, t := range _this.DecodeEncode {
		t.postDecodeInit()
		t.debug = _this.Debug
		t.trace = _this.Trace
	}
}

func (_this *CBETestRunner) validate() {
	if _this.Skip {
		return
	}

	if len(_this.Name) == 0 {
		panic(fmt.Errorf("missing name"))
	}

	for _, t := range _this.Encode {
		t.validate()
	}

	for _, t := range _this.Decode {
		t.validate()
	}

	for _, t := range _this.EncodeDecode {
		t.validate()
	}

	for _, t := range _this.DecodeEncode {
		t.validate()
	}
}

func (_this *CBETestRunner) run() {
	if _this.Skip {
		if _this.Debug {
			fmt.Printf("Skipping CBE Test %v\n", _this)
		}
		return
	}

	if _this.Debug {
		fmt.Printf("Running CBE Test %v:\n", _this)
	}

	for _, test := range _this.Encode {
		test.run()
	}

	for _, test := range _this.Decode {
		test.run()
	}

	for _, test := range _this.DecodeFail {
		test.run()
	}

	for _, test := range _this.EncodeDecode {
		test.run()
	}

	for _, test := range _this.DecodeEncode {
		test.run()
	}
}

type CBEDecodeFailTest []byte

func (_this CBEDecodeFailTest) run() {
	receiver, eventStore := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	err := capturePanic(func() {
		// debug.DebugOptions.PassThroughPanics will be true, so we won't get an error
		_ = cbe.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this)), receiver)
	})
	if err == nil {
		panic(fmt.Errorf("expected CBE %v to fail, but generated events [%v]", desc(_this), eventStore.Events))
	}
}

type CBEDecodeTest struct {
	Document []byte
	Events   []string
	events   test.Events
	debug    bool
	trace    bool
}

func (_this *CBEDecodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events...)
}

func (_this *CBEDecodeTest) validate() {
}

func (_this *CBEDecodeTest) run() {
	expectedEvents := _this.events
	actualEvents := decodeCBE(_this.trace, _this.debug, _this.Document)
	if !test.AreEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events [%v] but got [%v]",
			desc(_this.Document), expectedEvents, actualEvents)
	}
}

type CBEEncodeTest struct {
	Document []byte
	Events   []string
	events   test.Events
	debug    bool
	trace    bool
}

func (_this *CBEEncodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events...)
}

func (_this *CBEEncodeTest) validate() {
}

func (_this *CBEEncodeTest) run() {
	expectedDocument := _this.Document
	actualDocument := encodeCBE(_this.trace, _this.debug, _this.events)
	if !bytes.Equal(expectedDocument, actualDocument) {
		testFailed("Expected events [%v] to encode to document %v but got %v",
			desc(_this.events), desc(expectedDocument), desc(actualDocument))
	}
}

type CBEEncodeDecodeTest struct {
	Document []byte
	Events   []string
	events   test.Events
	debug    bool
	trace    bool
}

func (_this *CBEEncodeDecodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events...)
}

func (_this *CBEEncodeDecodeTest) validate() {
}

func (_this *CBEEncodeDecodeTest) run() {
	expectedDocument := _this.Document
	actualDocument := encodeCBE(_this.trace, _this.debug, _this.events)
	if !bytes.Equal(expectedDocument, actualDocument) {
		testFailed("Expected events [%v] to encode to document %v but got %v",
			desc(_this.events), desc(expectedDocument), desc(actualDocument))
	}

	expectedEvents := _this.events
	actualEvents := decodeCBE(_this.trace, _this.debug, actualDocument)
	if !test.AreEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events [%v] but got [%v]",
			desc(actualDocument), expectedEvents, actualEvents)
	}
}

type CBEDecodeEncodeTest struct {
	Document []byte
	Events   []string
	events   test.Events
	debug    bool
	trace    bool
}

func (_this *CBEDecodeEncodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events...)
}

func (_this *CBEDecodeEncodeTest) validate() {
}

func (_this *CBEDecodeEncodeTest) run() {
	expectedEvents := _this.events
	actualEvents := decodeCBE(_this.trace, _this.debug, _this.Document)
	if !test.AreEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events [%v] but got [%v]",
			desc(_this.Document), expectedEvents, actualEvents)
	}

	encodedDocument := encodeCBE(_this.trace, _this.debug, actualEvents)
	secondRunEvents := decodeCBE(_this.trace, _this.debug, encodedDocument)
	if !test.AreEventsEquivalent(expectedEvents, secondRunEvents) {
		testFailed("Expected document %v re-encoded to %v to produce events [%v] but got [%v]",
			desc(_this.Document), desc(encodedDocument), expectedEvents, secondRunEvents)
	}
}
