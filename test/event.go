// Copyright 2022 Karl Stenerud
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

// Test helper code.
package test

import (
	"bytes"
	"fmt"

	"github.com/kstenerud/go-concise-encoding/ce/events"
)

type EventInvocation func(receiver events.DataEventReceiver)

type Event interface {
	Name() string
	String() string
	ArrayElementCount() int
	Invoke(events.DataEventReceiver)
	IsEquivalentTo(Event) bool
	Comparable() string
	Value() interface{}
}

type Events []Event

func (_this Events) String() string {
	sb := bytes.Buffer{}
	for i, v := range _this {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteByte('"')
		sb.WriteString(v.String())
		sb.WriteByte('"')
	}
	return sb.String()
}

func (_this Events) AreEquivalentTo(that Events) bool {
	return AreEventsEquivalent(_this, that)
}

type BaseEvent struct {
	shortName         string
	invocation        EventInvocation
	values            []interface{}
	comparable        string
	stringified       string
	arrayElementCount int
}

func (_this *BaseEvent) Invoke(receiver events.DataEventReceiver) { _this.invocation(receiver) }
func (_this *BaseEvent) Name() string                             { return _this.shortName }
func (_this *BaseEvent) String() string                           { return _this.stringified }
func (_this *BaseEvent) Comparable() string                       { return _this.comparable }
func (_this *BaseEvent) ArrayElementCount() int                   { return _this.arrayElementCount }
func (_this *BaseEvent) IsEquivalentTo(that Event) bool {
	return _this.Comparable() == that.Comparable()
}
func (_this *BaseEvent) Value() interface{} {
	if len(_this.values) == 0 {
		return nil
	}
	return _this.values[len(_this.values)-1]
}

func ConstructEvent(shortName string, invocation EventInvocation, values ...interface{}) BaseEvent {
	return BaseEvent{
		shortName:         shortName,
		invocation:        invocation,
		values:            values,
		arrayElementCount: getArrayElementCount(values...),
		comparable:        constructComparable(shortName, values...),
		stringified:       constructStringified(shortName, values...),
	}
}

func constructComparable(shortName string, values ...interface{}) string {
	switch len(values) {
	case 0:
		return shortName
	case 1:
		return fmt.Sprintf("%v=%v", shortName, toComparable(values[0]))
	case 2:
		return fmt.Sprintf("%v=%v %v", shortName, toComparable(values[0]), toComparable(values[1]))
	default:
		panic(fmt.Errorf("expected 0, 1, or 2 values but got %v", len(values)))
	}
}

func constructStringified(shortName string, values ...interface{}) string {
	switch len(values) {
	case 0:
		return shortName
	case 1:
		return fmt.Sprintf("%v=%v", shortName, toStringified(values[0]))
	case 2:
		return fmt.Sprintf("%v=%v %v", shortName, toStringified(values[0]), toStringified(values[1]))
	default:
		panic(fmt.Errorf("expected 0, 1, or 2 values but got %v", len(values)))
	}
}
