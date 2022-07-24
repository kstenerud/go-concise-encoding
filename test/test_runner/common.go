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
	"runtime/debug"
	"strings"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
)

func testFailed(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func capturePanic(operation func()) (err error) {
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

func wrapPanic(recovery interface{}, format string, args ...interface{}) {
	if recovery != nil {
		message := fmt.Sprintf(format, args...)

		switch e := recovery.(type) {
		case error:
			panic(fmt.Errorf("%v: %w", message, e))
		default:
			panic(fmt.Errorf("%v: %v", message, e))
		}
	}
}

func reportAnyError(r interface{}, trace bool, format string, args ...interface{}) (reportedError bool) {
	if r == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("%v: %v\n", message, r)
	if trace {
		fmt.Println(string(debug.Stack()))
	}
	return true
}

func desc(v interface{}) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprintf("[%v]", v)
	case []byte:
		sb := strings.Builder{}
		sb.WriteByte('[')
		for i, b := range v {
			sb.WriteString(fmt.Sprintf("%02x", b))
			if i < len(v)-1 {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte(']')
		return sb.String()
	default:
		return fmt.Sprintf("%v", v)
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

func assertEventProducingOperation(trace bool, debug bool, receiver events.DataEventReceiver,
	operation func(receiver events.DataEventReceiver),
	describeSrc func() string) {

	if debug {
		receiver = test.NewEventPrinter(receiver)
	}

	var eventStore *test.EventCollection
	receiver, eventStore = test.NewEventCollector(receiver)
	op := func() {
		operation(receiver)
	}

	if trace {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("%v unexpectedly failed after producing events [%v]", describeSrc(), eventStore.Events)

				switch e := r.(type) {
				case error:
					panic(fmt.Errorf("%v: %w", message, e))
				default:
					panic(fmt.Errorf("%v: %v", message, r))
				}
			}
		}()
		operation(receiver)
	} else {
		if err := capturePanic(op); err != nil {
			testFailed("%v unexpectedly failed after producing events [%v]: %w", describeSrc(), eventStore.Events, err)
		}
	}
}

func decodeCTE(trace bool, debug bool, document string) test.Events {
	receiver, eventStore := test.NewEventCollector(nil)
	rules := rules.NewRules(receiver, nil)
	assertEventProducingOperation(trace, debug, rules,
		func(recv events.DataEventReceiver) {
			// debug.DebugOptions.PassThroughPanics will be true, so we won't get an error
			_ = cte.NewDecoder(nil).Decode(bytes.NewBuffer([]byte(document)), recv)
		}, func() string {
			return fmt.Sprintf("CTE %v", desc(document))
		})
	return eventStore.Events
}

func encodeCTE(trace bool, debug bool, eventStore test.Events) string {
	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	rules := rules.NewRules(encoder, nil)

	assertEventProducingOperation(trace, debug, rules,
		func(recv events.DataEventReceiver) {
			test.InvokeEventsAsCompleteDocument(recv, eventStore...)
		}, func() string {
			return fmt.Sprintf("events [%v]", desc(eventStore))
		})

	return buffer.String()
}

func decodeCBE(trace bool, debug bool, document []byte) test.Events {
	receiver, eventStore := test.NewEventCollector(nil)
	rules := rules.NewRules(receiver, nil)
	assertEventProducingOperation(trace, debug, rules,
		func(recv events.DataEventReceiver) {
			// debug.DebugOptions.PassThroughPanics will be true, so we won't get an error
			_ = cbe.NewDecoder(nil).Decode(bytes.NewBuffer(document), recv)
		}, func() string {
			return fmt.Sprintf("CBE %v", desc(document))
		})
	return eventStore.Events
}

func encodeCBE(trace bool, debug bool, eventStore test.Events) []byte {
	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(buffer)
	rules := rules.NewRules(encoder, nil)

	assertEventProducingOperation(trace, debug, rules,
		func(recv events.DataEventReceiver) {
			test.InvokeEventsAsCompleteDocument(recv, eventStore...)
		}, func() string {
			return fmt.Sprintf("events [%v]", desc(eventStore))
		})

	return buffer.Bytes()
}
