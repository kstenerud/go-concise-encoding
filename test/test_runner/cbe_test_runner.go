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
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"
)

type CBETestRunner struct {
	Name              string
	DecodeMustSucceed []*CBEDecodeSuccessTest
	DecodeMustFail    []CBEDecodeFailTest
	Skip              bool
	Debug             bool
	Trace             bool
}

func (_this *CBETestRunner) String() string {
	return _this.Name
}

func (_this *CBETestRunner) postDecodeInit() {
	for _, t := range _this.DecodeMustSucceed {
		t.postDecodeInit()
		t.debug = _this.Debug
	}
}

func (_this *CBETestRunner) validate() {
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

	for _, test := range _this.DecodeMustSucceed {
		test.run()
	}

	for _, test := range _this.DecodeMustFail {
		test.run()
	}
}

type CBEDecodeFailTest []byte

func (_this CBEDecodeFailTest) run() {
	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	err := capturePanic(func() {
		cbe.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this)), receiver)
	})
	if err == nil {
		panic(fmt.Errorf("expected CBE %v to fail, but generated events %v", desc(_this), eventStore.Events))
	}
}

type CBEDecodeSuccessTest struct {
	Source []byte
	Events []string
	events []*test.TEvent
	debug  bool
}

func (_this *CBEDecodeSuccessTest) postDecodeInit() {
	_this.events = event_parser.ParseEvents(_this.Events)
}

func (_this *CBEDecodeSuccessTest) validate() {
}

func (_this *CBEDecodeSuccessTest) testFailed(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (_this *CBEDecodeSuccessTest) assertOperation(receiver events.DataEventReceiver,
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
		_this.testFailed("%v unexpecbedly failed after producing events %v: %w", describeSrc(), eventStore.Events, err)
	}
}

func (_this *CBEDecodeSuccessTest) run() {
	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cbe.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Source)), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(_this.Source))
		})

	expecbedEvents := _this.events
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(expecbedEvents, actualEvents) {
		_this.testFailed("Expecbed CBE %v to produce events %v but got %v",
			desc(_this.Source), expecbedEvents, actualEvents)
	}
}