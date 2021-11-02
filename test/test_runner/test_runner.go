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
	"github.com/kstenerud/go-concise-encoding/debug"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"

	"github.com/kstenerud/go-equivalence"
)

// API

// Run a CE test suite, described by the CTE document at testDescriptorFile
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

	debug.DebugOptions.PassThroughPanics = true
	testSuite.run(t)
	debug.DebugOptions.PassThroughPanics = false
}

// Test Suite

type CETestSuite struct {
	TestFile string
	Tests    []*CETestRunner
}

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

	var unitTest *CETestRunner
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

	var unitTest *CETestRunner
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

	var unitTest *CETestRunner
	for index, unitTest = range _this.Tests {
		unitTest.run(t)
	}
}

// Unit Test

type CETestRunner struct {
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
	Trace     bool
	Skip      bool

	events             []*test.TEvent
	srcEventTypes      map[string]bool
	dstEventTypes      map[string]bool
	requiredEventTypes map[string]bool
	context            string
}

func (_this *CETestRunner) Description() string {
	if len(_this.context) > 0 {
		return fmt.Sprintf("%v index %v (%v): %v", _this.TestFile, _this.TestIndex, _this.Name, _this.context)
	}
	return fmt.Sprintf("%v index %v (%v)", _this.TestFile, _this.TestIndex, _this.Name)
}

func (_this *CETestRunner) postDecodeInit() {
	if _this.Skip {
		return
	}

	_this.events = event_parser.ParseEvents(_this.Events)
	_this.Cte = strings.TrimSpace(_this.Cte)
	_this.srcEventTypes = toMembershipSet(_this.From)
	_this.dstEventTypes = toMembershipSet(_this.To)
	_this.requiredEventTypes = toMembershipSet(_this.From, _this.To)
	_this.validate()
}

func (_this *CETestRunner) validate() {
	if _this.Skip {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("%v: %w", _this.Description(), v))
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

func (_this *CETestRunner) run(t *testing.T) {
	if _this.Skip {
		if _this.Debug {
			fmt.Printf("Skipping CE Test %v\n", _this.Description())
		}
		return
	}

	if _this.Debug {
		fmt.Printf("Running CE Test %v:\n", _this.Description())
	}

	defer func() {
		if r := recover(); r != nil {
			if !_this.Fail {
				if _this.Trace {
					panic(fmt.Errorf("%v: %w", _this.Description(), r))
				} else {
					t.Errorf("%v: %v", _this.Description(), r)
				}
			}
		}
	}()

	if _this.srcEventTypes["e"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runEventToNothing()
		}
		if _this.dstEventTypes["e"] {
			_this.runEventToEvent()
		}
		if _this.dstEventTypes["b"] {
			_this.runEventToCBE()
		}
		if _this.dstEventTypes["t"] {
			_this.runEventToCTE()
		}
	}

	if _this.srcEventTypes["b"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runCBEToNothing()
		}
		if _this.dstEventTypes["e"] {
			_this.runCBEToEvent()
		}
		if _this.dstEventTypes["b"] {
			_this.runCBEToCBE()
		}
		if _this.dstEventTypes["t"] {
			_this.runCBEToCTE()
		}
	}

	if _this.srcEventTypes["t"] {
		if len(_this.dstEventTypes) == 0 {
			_this.runCTEToNothing()
		}
		if _this.dstEventTypes["e"] {
			_this.runCTEToEvent()
		}
		if _this.dstEventTypes["b"] {
			_this.runCTEToCBE()
		}
		if _this.dstEventTypes["t"] {
			_this.runCTEToCTE()
		}
	}
}

func (_this *CETestRunner) capturePanic(operation func()) (err error) {
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

func (_this *CETestRunner) testFailed(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (_this *CETestRunner) assertOperation(receiver events.DataEventReceiver,
	operation func(receiver events.DataEventReceiver),
	describeSrc func() string) {

	if _this.Debug {
		receiver = test.NewStdoutTEventPrinter(receiver)
	}

	if !_this.Fail && _this.Trace {
		operation(receiver)
		return
	}

	eventStore := test.NewTEventStore(receiver)
	receiver = eventStore

	err := _this.capturePanic(func() {
		operation(receiver)
	})

	if !_this.Fail && err != nil {
		_this.testFailed("%v unexpectedly failed after producing events %v: %w", describeSrc(), eventStore.Events, err)
	} else if _this.Fail && err == nil {
		_this.testFailed("%v unexpectedly succeeded and produced events %v", describeSrc(), eventStore.Events)
	}
}

func (_this *CETestRunner) driveEvents(receiver events.DataEventReceiver, events ...*test.TEvent) {
	for _, event := range events {
		event.Invoke(receiver)
	}
}

func (_this *CETestRunner) beginTestPhase(phase string) {
	_this.context = phase
	if _this.Debug {
		fmt.Printf("Running test phase %v\n", _this.Description())
	}
}

func (_this *CETestRunner) runEventToNothing() {
	_this.beginTestPhase("E2N")

	receiver := rules.NewRules(events.NewNullEventReceiver(), nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			_this.driveEvents(recv, _this.events...)
		}, func() string {
			return fmt.Sprintf("Events %v", _this.events)
		})
}

func (_this *CETestRunner) runCBEToNothing() {
	_this.beginTestPhase("B2N")

	receiver := rules.NewRules(events.NewNullEventReceiver(), nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(_this.Cbe))
		})
}

func (_this *CETestRunner) runCTEToNothing() {
	_this.beginTestPhase("T2N")

	receiver := rules.NewRules(events.NewNullEventReceiver(), nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(_this.Cte))
		})
}

func (_this *CETestRunner) runEventToEvent() {
	_this.beginTestPhase("E2E")

	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			_this.driveEvents(recv, _this.events...)
		}, func() string {
			return fmt.Sprintf("Events %v", _this.events)
		})

	expectedEvents := _this.events
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		_this.testFailed("Expected events %v to produce events %v but got %v",
			_this.events, expectedEvents, actualEvents)
	}
}

func (_this *CETestRunner) runEventToCBE() {
	_this.beginTestPhase("E2B")

	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			_this.driveEvents(recv, _this.events...)
		}, func() string {
			return fmt.Sprintf("Events %v", _this.events)
		})

	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected events %v to produce CBE %v but got %v",
			_this.events, desc(expectedDocument), desc(actualDocument))
	}
}

func (_this *CETestRunner) runEventToCTE() {
	_this.beginTestPhase("E2T")

	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			_this.driveEvents(recv, _this.events...)
		}, func() string {
			return fmt.Sprintf("Events %v", _this.events)
		})

	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected events %v to produce CTE %v but got %v after events %v",
			_this.events, desc(expectedDocument), desc(actualDocument), _this.events)
	}
}

func (_this *CETestRunner) runCBEToEvent() {
	_this.beginTestPhase("B2E")

	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(_this.Cbe))
		})

	expectedEvents := _this.events
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		_this.testFailed("Expected CBE %v to produce events %v but got %v",
			desc(_this.Cbe), expectedEvents, actualEvents)
	}
}

func (_this *CETestRunner) runCBEToCBE() {
	_this.beginTestPhase("B2B")

	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(_this.Cbe))
		})

	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected CBE %v to produce CBE %v but got %v",
			desc(_this.Cbe), desc(expectedDocument), desc(actualDocument))
	}
}

func (_this *CETestRunner) runCBEToCTE() {
	_this.beginTestPhase("B2T")

	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cbe.NewDecoder(nil).Decode(bytes.NewBuffer(_this.Cbe), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(_this.Cbe))
		})

	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected CBE %v to produce CTE %v but got %v",
			desc(_this.Cbe), desc(expectedDocument), desc(actualDocument))
	}
}

func (_this *CETestRunner) runCTEToEvent() {
	_this.beginTestPhase("T2E")

	eventStore := test.NewTEventStore(events.NewNullEventReceiver())
	receiver := rules.NewRules(eventStore, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(_this.Cte))
		})

	expectedEvents := _this.events
	actualEvents := eventStore.Events
	if !test.AreAllEventsEqual(expectedEvents, actualEvents) {
		_this.testFailed("Expected CTE %v to produce events %v but got %v",
			desc(_this.Cte), expectedEvents, actualEvents)
	}
}

func (_this *CETestRunner) runCTEToCBE() {
	_this.beginTestPhase("T2B")

	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(_this.Cte))
		})

	expectedDocument := _this.Cbe
	actualDocument := buffer.Bytes()
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected CTE %v to produce CBE %v but got %v",
			desc(_this.Cte), desc(expectedDocument), desc(actualDocument))
	}
}

func (_this *CETestRunner) runCTEToCTE() {
	_this.beginTestPhase("T2T")

	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	receiver := rules.NewRules(encoder, nil)
	_this.assertOperation(receiver,
		func(recv events.DataEventReceiver) {
			cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(_this.Cte)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(_this.Cte))
		})

	expectedDocument := _this.Cte
	actualDocument := string(buffer.Bytes())
	if !equivalence.IsEquivalent(expectedDocument, actualDocument) {
		_this.testFailed("Expected CTE %v to produce CTE %v but got %v",
			desc(_this.Cte), desc(expectedDocument), desc(actualDocument))
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

func desc(v interface{}) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("[%v]", v)
	case []byte:
		sb := strings.Builder{}
		sb.WriteByte('[')
		for i, b := range v.([]byte) {
			sb.WriteString(fmt.Sprintf("%02x", b))
			if i < len(v.([]byte))-1 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte(']')
		return sb.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
