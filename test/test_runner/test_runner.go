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
	"os"
	"strings"
	"testing"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"

	"github.com/kstenerud/go-describe"
	"github.com/kstenerud/go-equivalence"
)

// API

func RunCEUnitTests(t *testing.T, testDescriptorFile string) {
	file, err := os.Open(testDescriptorFile)
	if err != nil {
		panic(fmt.Errorf("Unexpected error opening unit test file %v: %v", testDescriptorFile, err))
	}
	var testSuite *CETestSuite

	testSuiteIntf, err := ce.UnmarshalCTE(file, testSuite, nil)
	if err != nil {
		panic(fmt.Errorf("Malformed unit test: Unexpected CTE decode error in file %v: %w", testDescriptorFile, err))
	}
	testSuite = testSuiteIntf.(*CETestSuite)
	testSuite.TestFile = testDescriptorFile
	testSuite.postDecodeInit()
	testSuite.run(t)
}

// Test Suite

type CETestSuite struct {
	TestFile string
	// TODO: Decoding silently fails if nonexistent field is in CTE and is a list?
	// But int is ok, and string causes printing to show everything???
	References interface{}
	Tests      []*CETest
}

// func (_this *CETestSuite) String() string {
// 	return describe.D(_this)
// }

func (_this *CETestSuite) postDecodeInit() {
	index := 0
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Malformed unit test: Test index %v: %w", index, v))
			default:
				panic(v)
			}
		}
	}()

	var unitTest *CETest
	for index, unitTest = range _this.Tests {
		unitTest.TestFile = _this.TestFile
		unitTest.TestIndex = index
		unitTest.postDecodeInit()
	}
}

func (_this *CETestSuite) validate() {
	index := 0
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Malformed unit test: Test index %v: %w", index, v))
			default:
				panic(v)
			}
		}
	}()

	var unitTest *CETest
	for index, unitTest = range _this.Tests {
		unitTest.validate()
	}
}

func (_this *CETestSuite) run(t *testing.T) {
	index := 0
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Unit test failed: Test index %v: %w", index, v))
			default:
				panic(v)
			}
		}
	}()

	var unitTest *CETest
	for index, unitTest = range _this.Tests {
		unitTest.run(t)
	}
}

// Unit Test

type CETest struct {
	Name      string
	TestFile  string
	TestIndex int
	Events    []string
	Cte       string
	Cbe       []byte
	From      []string
	To        []string
	Fail      bool
	Debug     bool
	Panic     bool
	Skip      bool

	decodedEvents      []*test.TEvent
	srcEventTypes      map[string]bool
	dstEventTypes      map[string]bool
	requiredEventTypes map[string]bool
}

// func (_this *CETest) String() string {
// 	return describe.D(_this)
// }

func (_this *CETest) postDecodeInit() {
	_this.decodedEvents = event_parser.ParseEvents(_this.Events)
	_this.Cte = strings.TrimSpace(_this.Cte)
	_this.srcEventTypes = toMembershipSet(_this.From)
	_this.dstEventTypes = toMembershipSet(_this.To)
	_this.requiredEventTypes = toMembershipSet(_this.From, _this.To)
	_this.validate()
}

func (_this *CETest) validate() {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Test %v: %w", _this.Name, v))
			default:
				panic(v)
			}
		}
	}()

	if len(_this.Name) == 0 {
		panic(fmt.Errorf("Missing name"))
	}
	if len(_this.From) == 0 {
		panic(fmt.Errorf("Must have at least 1 \"from\" entry"))
	}

	for index, value := range _this.From {
		if value != "e" && value != "t" && value != "b" {
			panic(fmt.Errorf("%v: Unknown \"from\" type at index %v", value, index))
		}
	}

	for index, value := range _this.To {
		if value != "e" && value != "t" && value != "b" {
			panic(fmt.Errorf("%v: Unknown \"to\" type at index %v", value, index))
		}
	}

	if _this.requiredEventTypes["e"] && len(_this.Events) == 0 {
		panic(fmt.Errorf("Test calls for events as src or dst but does not provide any"))
	}
	if _this.requiredEventTypes["b"] && len(_this.Cbe) == 0 {
		panic(fmt.Errorf("Test calls for CBE as src or dst but does not provide any"))
	}
	if _this.requiredEventTypes["t"] && len(_this.Cte) == 0 {
		panic(fmt.Errorf("Test calls for CTE as src or dst but does not provide any"))
	}
}

func (_this *CETest) run(t *testing.T) {
	if _this.Skip {
		if _this.Debug {
			fmt.Printf("Skipping %v\n", _this.description())
		}
		return
	}

	if _this.Debug {
		fmt.Printf("CE Test %v:\n", _this.description())
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Test %v: %w", _this.Name, v))
			default:
				panic(v)
			}
		}
	}()

	if _this.srcEventTypes["e"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runEventToNothing(t)
		}
		if _this.dstEventTypes["e"] {
			_this.runEventToEvent(t)
		}
		if _this.dstEventTypes["b"] {
			_this.runEventToCBE(t)
		}
		if _this.dstEventTypes["t"] {
			_this.runEventToCTE(t)
		}
	}

	if _this.srcEventTypes["b"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runCBEToNothing(t)
		}
		if _this.dstEventTypes["e"] {
			_this.runCBEToEvent(t)
		}
		if _this.dstEventTypes["b"] {
			_this.runCBEToCBE(t)
		}
		if _this.dstEventTypes["t"] {
			_this.runCBEToCTE(t)
		}
	}

	if _this.srcEventTypes["t"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runCTEToNothing(t)
		}
		if _this.dstEventTypes["e"] {
			_this.runCTEToEvent(t)
		}
		if _this.dstEventTypes["b"] {
			_this.runCTEToCBE(t)
		}
		if _this.dstEventTypes["t"] {
			_this.runCTEToCTE(t)
		}
	}
}

func (_this *CETest) capturePanic(operation func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	operation()
	return
}

func (_this *CETest) description() string {
	return fmt.Sprintf("%v:%v (%v)", _this.TestFile, _this.TestIndex, _this.Name)
}

func (_this *CETest) fatalf(t *testing.T, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	t.Fatalf("%v: %v", _this.description(), message)
}

func (_this *CETest) assertPanics(t *testing.T, operation func()) {
	if err := _this.capturePanic(operation); err == nil {
		_this.fatalf(t, "Expected a panic")
	}
}

func (_this *CETest) assertNoPanic(t *testing.T, operation func()) {
	if _this.Panic {
		operation()
	}
	if err := _this.capturePanic(operation); err != nil {
		_this.fatalf(t, "Unexpected panic: %v", err)
	}
}

func (_this *CETest) assertEvents(t *testing.T, receiver events.DataEventReceiver, events ...*test.TEvent) {
	operation := func() {
		for _, event := range events {
			event.Invoke(receiver)
		}
	}
	if _this.Fail {
		_this.assertPanics(t, operation)
	} else {
		_this.assertNoPanic(t, operation)
	}
}

func (_this *CETest) runEventToNothing(t *testing.T) {
	var receiver events.DataEventReceiver = events.NewNullEventReceiver()
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runEventToNothing]\n")
	}
	receiver = rules.NewRules(receiver, nil)

	_this.assertEvents(t, receiver, _this.decodedEvents...)
}

func (_this *CETest) runCBEToNothing(t *testing.T) {
	eventStore := test.NewTEventStore()
	var receiver events.DataEventReceiver = eventStore
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCBEToNothing]\n")
	}
	receiver = rules.NewRules(receiver, nil)

	err := cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), receiver)
	if _this.Fail && err == nil {
		_this.fatalf(t, "B2N: Expected %v to fail, but produced %v", describe.D(_this.Cbe), eventStore.Events)
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "B2N: %v unexpectedly failed after producing %v: %v", describe.D(_this.Cbe), eventStore.Events, err)
	}
}

func (_this *CETest) runCTEToNothing(t *testing.T) {
	eventStore := test.NewTEventStore()
	var receiver events.DataEventReceiver = eventStore
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCTEToNothing]\n")
	}
	receiver = rules.NewRules(receiver, nil)

	err := cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), receiver)
	if _this.Fail && err == nil {
		_this.fatalf(t, "T2N: Expected [%v] to fail, but produced %v", _this.Cte, eventStore.Events)
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "T2N: [%v] unexpectedly failed after producing %v: %v", _this.Cte, eventStore.Events, err)
	}
}

func (_this *CETest) runEventToEvent(t *testing.T) {
	eventStore := test.NewTEventStore()
	var receiver events.DataEventReceiver = eventStore
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runEventToEvent]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	expectedEvents := _this.decodedEvents
	_this.assertEvents(t, receiver, expectedEvents...)
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(actualEvents, expectedEvents) {
		_this.fatalf(t, "E2E: Expected %v but got %v", expectedEvents, actualEvents)
	}
}

func (_this *CETest) runEventToCBE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runEventToCBE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	_this.assertEvents(t, receiver, _this.decodedEvents...)
	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		_this.fatalf(t, "E2B: Expected %v but got %v after events %v", describe.D(expectedDocument), describe.D(actualDocument), _this.decodedEvents)
	}
}

func (_this *CETest) runEventToCTE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runEventToCTE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	_this.assertEvents(t, receiver, _this.decodedEvents...)
	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if !equivalence.IsEquivalent(actualDocument, expectedDocument) {
		_this.fatalf(t, "E2T: Expected %v but got %v after events %v", expectedDocument, actualDocument, _this.decodedEvents)
	}
}

func (_this *CETest) runCBEToEvent(t *testing.T) {
	eventStore := test.NewTEventStore()
	var receiver events.DataEventReceiver = eventStore
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCBEToEvent]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), receiver)
	if !_this.Fail {
		if err != nil {
			_this.fatalf(t, "T2E: %v unexpectedly failed after events %v: %v", describe.D(_this.Cbe), eventStore.Events, err)
		}
		expectedEvents := _this.decodedEvents
		actualEvents := eventStore.Events
		if !test.AreAllEventsEqual(actualEvents, expectedEvents) {
			_this.fatalf(t, "T2E: Expected %v to produce %v but got %v", describe.D(_this.Cbe), expectedEvents, actualEvents)
		}
	} else if err == nil {
		_this.fatalf(t, "T2E: Expected %v to fail, but produced %v", describe.D(_this.Cbe), eventStore.Events)
	}
}

func (_this *CETest) runCBEToCBE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCBEToCBE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), receiver)
	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if _this.Fail && err == nil {
		_this.fatalf(t, "B2B: Expected %v to fail, but produced %v", describe.D(expectedDocument), describe.D(actualDocument))
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "B2B: %v unexpectedly failed: %v", describe.D(_this.Cbe), err)
	} else if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.fatalf(t, "B2B: Expected %v but got %v", describe.D(expectedDocument), describe.D(actualDocument))
	}
}

func (_this *CETest) runCBEToCTE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCBEToCTE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), receiver)
	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if _this.Fail && err == nil {
		_this.fatalf(t, "B2T: Expected %v to fail, but produced %v", describe.D(_this.Cbe), describe.D(actualDocument))
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "B2T: %v unexpectedly failed: %v", describe.D(_this.Cbe), err)
	} else if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.fatalf(t, "B2T: Expected %v but got %v", describe.D(expectedDocument), describe.D(actualDocument))
	}
}

func (_this *CETest) runCTEToEvent(t *testing.T) {
	eventStore := test.NewTEventStore()
	var receiver events.DataEventReceiver = eventStore
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCTEToEvent]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), receiver)
	if !_this.Fail {
		if err != nil {
			_this.fatalf(t, "T2E: %v unexpectedly failed after %v: %v", _this.Cte, eventStore.Events, err)
		}
		expectedEvents := _this.decodedEvents
		actualEvents := eventStore.Events
		if !test.AreAllEventsEqual(actualEvents, expectedEvents) {
			_this.fatalf(t, "T2E: Expected %v to produce %v but got %v", _this.Cte, expectedEvents, actualEvents)
		}
	} else if err == nil {
		_this.fatalf(t, "T2E: Expected %v to fail, but produced %v", _this.Cte, eventStore.Events)
	}
}

func (_this *CETest) runCTEToCBE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCTEToCBE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), receiver)
	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if _this.Fail && err == nil {
		_this.fatalf(t, "T2B: Expected %v to fail, but produced %v", _this.Cte, describe.D(actualDocument))
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "T2B: %v unexpectedly failed: %v", _this.Cte, err)
	} else if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.fatalf(t, "T2B: Expected %v but got %v", describe.D(expectedDocument), describe.D(actualDocument))
	}
}

func (_this *CETest) runCTEToCTE(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	var receiver events.DataEventReceiver = encoder
	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
		fmt.Printf("[runCTEToCTE]\n")
	}
	receiver = rules.NewRules(receiver, nil)
	err := cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), receiver)
	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if _this.Fail && err == nil {
		_this.fatalf(t, "T2T: Expected %v to fail, but produced %v", _this.Cte, actualDocument)
	} else if !_this.Fail && err != nil {
		_this.fatalf(t, "T2T: %v unexpectedly failed: %v", _this.Cte, err)
	} else if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.fatalf(t, "T2T: Expected %v but got %v", describe.D(expectedDocument), describe.D(actualDocument))
	}
}

func toMembershipSet(valueArrays ...[]string) map[string]bool {
	set := make(map[string]bool)

	for _, array := range valueArrays {
		for _, value := range array {
			set[value] = true
		}
	}
	return set
}
