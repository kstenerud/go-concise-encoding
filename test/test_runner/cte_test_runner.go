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

	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"
)

type CTETestRunner struct {
	Name         string
	Encode       []*CTEEncodeTest
	Decode       []*CTEDecodeTest
	DecodeFail   []CTEDecodeFailTest
	EncodeDecode []*CTEEncodeDecodeTest
	DecodeEncode []*CTEDecodeEncodeTest
	Skip         bool
	Debug        bool
	Trace        bool
}

func (_this *CTETestRunner) String() string {
	return _this.Name
}

func (_this *CTETestRunner) postDecodeInit() {
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

func (_this *CTETestRunner) validate() {
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

func (_this *CTETestRunner) run() {
	if _this.Skip {
		if _this.Debug {
			fmt.Printf("Skipping CTE Test %v\n", _this)
		}
		return
	}

	if _this.Debug {
		fmt.Printf("Running CTE Test %v:\n", _this)
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

type CTEDecodeFailTest string

func (_this CTEDecodeFailTest) run() {
	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	err := capturePanic(func() {
		// debug.DebugOptions.PassThroughPanics will be true, so we won't get an error
		_ = cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this)), receiver)
	})
	if err == nil {
		panic(fmt.Errorf("expected CTE [%v] to fail, but generated events %v", _this, eventStore.Events))
	}
}

type CTEDecodeTest struct {
	Document string
	Events   []string
	events   []*test.TEvent
	debug    bool
	trace    bool
}

func (_this *CTEDecodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CTEDecodeTest) validate() {
}

func (_this *CTEDecodeTest) run() {
	expectedEvents := _this.events
	actualEvents := decodeCTE(_this.trace, _this.debug, _this.Document)
	if !test.AreAllEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events %v but got %v",
			desc(_this.Document), expectedEvents, actualEvents)
	}
}

type CTEEncodeTest struct {
	Document string
	Events   []string
	events   []*test.TEvent
	debug    bool
	trace    bool
}

func (_this *CTEEncodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CTEEncodeTest) validate() {
}

func (_this *CTEEncodeTest) run() {
	expectedDocument := _this.Document
	actualDocument := encodeCTE(_this.trace, _this.debug, _this.events)
	if expectedDocument != actualDocument {
		testFailed("Expected events %v to encode to document %v but got %v",
			desc(_this.events), desc(expectedDocument), desc(actualDocument))
	}
}

type CTEEncodeDecodeTest struct {
	Document string
	Events   []string
	events   []*test.TEvent
	debug    bool
	trace    bool
}

func (_this *CTEEncodeDecodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CTEEncodeDecodeTest) validate() {
}

func (_this *CTEEncodeDecodeTest) run() {
	expectedDocument := _this.Document
	actualDocument := encodeCTE(_this.trace, _this.debug, _this.events)
	if expectedDocument != actualDocument {
		testFailed("Expected events %v to encode to document %v but got %v",
			desc(_this.events), desc(expectedDocument), desc(actualDocument))
	}

	expectedEvents := _this.events
	actualEvents := decodeCTE(_this.trace, _this.debug, actualDocument)
	if !test.AreAllEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events %v but got %v",
			desc(actualDocument), expectedEvents, actualEvents)
	}
}

type CTEDecodeEncodeTest struct {
	Document string
	Events   []string
	events   []*test.TEvent
	debug    bool
	trace    bool
}

func (_this *CTEDecodeEncodeTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CTEDecodeEncodeTest) validate() {
}

func (_this *CTEDecodeEncodeTest) run() {
	expectedEvents := _this.events
	actualEvents := decodeCTE(_this.trace, _this.debug, _this.Document)
	if !test.AreAllEventsEquivalent(expectedEvents, actualEvents) {
		testFailed("Expected document %v to produce events %v but got %v",
			desc(_this.Document), expectedEvents, actualEvents)
	}

	encodedDocument := encodeCTE(_this.trace, _this.debug, actualEvents)
	secondRunEvents := decodeCTE(_this.trace, _this.debug, encodedDocument)
	if !test.AreAllEventsEquivalent(expectedEvents, secondRunEvents) {
		testFailed("Expected document %v re-encoded to %v to produce events %v but got %v",
			desc(_this.Document), desc(encodedDocument), expectedEvents, secondRunEvents)
	}
}
