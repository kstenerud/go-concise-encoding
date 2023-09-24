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

package test_runner

// Contrived file name to get around the idiotic go "feature" that gives
// special meaning to filenames ending in _test. Every experienced engineer
// knows that you NEVER add extra constraints to existing published standards
// (no matter how "clever" you think you are) because it always bites you in
// the ass eventually, and requires ugly workarounds. Simplicity indeed...

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser"
)

type BaseTest struct {
	ceVersion   int
	CBE         []byte   `ce:"order=10,omit_empty"`
	CTE         string   `ce:"order=20,omit_empty"`
	Events      []string `ce:"order=30,omit_empty"`
	RawDocument bool     `ce:"order=40,omit_zero"`
	Skip        bool     `ce:"order=50,omit_zero"`
	Debug       bool     `ce:"order=60,omit_zero"` // When true, don't convert panics to errors.
	events      test.Events
	context     string
}

func (_this *BaseTest) PostDecodeInit(ceVersion int, context string) error {
	if _this.Skip {
		return nil
	}

	_this.ceVersion = ceVersion
	_this.context = context
	_this.CTE = _this.PostDecodeCTE(_this.CTE)
	_this.CBE = _this.PostDecodeCBE(_this.CBE)
	_this.events = _this.PostDecodeEvents(_this.Events)

	return nil
}

func (_this *BaseTest) PostDecodeCTE(cte string) string {
	cte = strings.TrimSpace(cte)
	if !_this.RawDocument && len(cte) > 0 {
		cte = fmt.Sprintf("c%v\n%v", _this.ceVersion, cte)
	}
	return cte
}

func (_this *BaseTest) PostDecodeCBE(cbe []byte) []byte {
	if !_this.RawDocument && len(cbe) > 0 {
		b := make([]byte, len(cbe)+2)
		b[0] = 0x81
		b[1] = byte(_this.ceVersion)
		copy(b[2:], cbe)
		cbe = b
	}
	return cbe
}

func (_this *BaseTest) PostDecodeEvents(eventStrings []string) test.Events {
	events := event_parser.ParseEvents(eventStrings...)
	if !_this.RawDocument && len(events) > 0 {
		newEvents := make(test.Events, 0, len(events)+1)
		newEvents = append(newEvents, test.V(uint64(_this.ceVersion)))
		newEvents = append(newEvents, events...)
		events = newEvents
	}
	return events
}

func (_this *BaseTest) errorf(format string, args ...interface{}) error {
	if message := fmt.Sprintf(format, args...); len(message) > 0 {
		return fmt.Errorf("%v: %v", _this.context, message)
	} else {
		return fmt.Errorf("%v", _this.context)
	}
}

func (_this *BaseTest) wrapError(err error, format string, args ...interface{}) error {
	if message := fmt.Sprintf(format, args...); len(message) > 0 {
		return fmt.Errorf("%v: %v: %w", _this.context, message, err)
	} else {
		return fmt.Errorf("%v: %w", _this.context, err)
	}
}

func (_this *BaseTest) cteToEvents(document string) (result test.Events, err error) {
	if !_this.Debug {
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
	}

	config := configuration.New()
	config.Debug.PassThroughPanics = _this.Debug
	decoder := cte.NewDecoder(config)
	buffer := bytes.NewBuffer([]byte(document))
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, config)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}

func (_this *BaseTest) cbeToEvents(document []byte) (result test.Events, err error) {
	if !_this.Debug {
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
	}

	config := configuration.New()
	config.Debug.PassThroughPanics = _this.Debug
	decoder := cbe.NewDecoder(config)
	buffer := bytes.NewBuffer(document)
	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, config)
	if err = decoder.Decode(buffer, receiver); err != nil {
		return
	}
	result = collection.Events
	return
}
