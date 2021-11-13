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
	Name              string
	DecodeMustSucceed []*CTEDecodeSuccessTest
	DecodeMustFail    []CTEDecodeFailTest
	Skip              bool
	Debug             bool
	Trace             bool
}

func (_this *CTETestRunner) String() string {
	return _this.Name
}

func (_this *CTETestRunner) postDecodeInit() {
	for _, t := range _this.DecodeMustSucceed {
		t.postDecodeInit()
		t.debug = _this.Debug
	}
}

func (_this *CTETestRunner) validate() {
	if _this.Skip {
		return
	}

	if len(_this.Name) == 0 {
		panic(fmt.Errorf("missing name"))
	}

	for _, t := range _this.DecodeMustSucceed {
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

	for _, test := range _this.DecodeMustSucceed {
		test.run()
	}

	for _, test := range _this.DecodeMustFail {
		test.run()
	}
}

type CTEDecodeFailTest string

func (_this CTEDecodeFailTest) run() {
	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	err := capturePanic(func() {
		cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this)), receiver)
	})
	if err == nil {
		panic(fmt.Errorf("expected CTE [%v] to fail, but generated events %v", _this, eventStore.Events))
	}
}

type CTEDecodeSuccessTest struct {
	Source string
	Events []string
	events []*test.TEvent
	debug  bool
}

func (_this *CTEDecodeSuccessTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CTEDecodeSuccessTest) validate() {
}

func (_this *CTEDecodeSuccessTest) testFailed(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (_this *CTEDecodeSuccessTest) assertOperation(receiver events.DataEventReceiver,
	operation func(receiver events.DataEventReceiver),
	describeSrc func() string) {

	if _this.debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
	}

	eventStore := test.NewTEventStore(receiver)
	receiver = eventStore

	err := capturePanic(func() {
		operation(receiver)
	})

	if err != nil {
		_this.testFailed("%v unexpectedly failed after producing events %v: %w", describeSrc(), eventStore.Events, err)
	}
}

func (_this *CTEDecodeSuccessTest) run() {
	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Source)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(_this.Source))
		})

	expectedEvents := _this.events
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		_this.testFailed("Expected CTE %v to produce events %v but got %v",
			desc(_this.Source), expectedEvents, actualEvents)
	}
}
