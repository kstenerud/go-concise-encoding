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
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/rules"
	"github.com/kstenerud/go-concise-encoding/test"
)

type MustSucceedTest struct {
	BaseTest
	FromCBE    []byte   `ce:"order=11,omitempty"`
	ToCBE      []byte   `ce:"order=12,omitempty"`
	FromCTE    string   `ce:"order=21,omitempty"`
	ToCTE      string   `ce:"order=22,omitempty"`
	FromEvents []string `ce:"order=31,omitempty"`
	ToEvents   []string `ce:"order=32,omitempty"`
	fromEvents test.Events
	toEvents   test.Events
}

func (_this *MustSucceedTest) PostDecodeInit(ceVersion int, context string, index int) error {
	if _this.Skip {
		return nil
	}
	context = fmt.Sprintf(`%v, "must succeed" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	if len(_this.FromCBE) == 0 {
		_this.FromCBE = _this.CBE
	} else {
		_this.FromCBE = _this.PostDecodeCBE(_this.FromCBE)
	}
	if len(_this.ToCBE) == 0 {
		_this.ToCBE = _this.CBE
	} else {
		_this.ToCBE = _this.PostDecodeCBE(_this.ToCBE)
	}

	if len(_this.FromCTE) == 0 {
		_this.FromCTE = _this.CTE
	} else {
		_this.FromCTE = _this.PostDecodeCTE(_this.FromCTE)
	}
	if len(_this.ToCTE) == 0 {
		_this.ToCTE = _this.CTE
	} else {
		_this.ToCTE = _this.PostDecodeCTE(_this.ToCTE)
	}

	if len(_this.FromEvents) == 0 {
		_this.FromEvents = _this.Events
		_this.fromEvents = _this.events
	} else {
		_this.fromEvents = _this.PostDecodeEvents(_this.FromEvents)
	}
	if len(_this.ToEvents) == 0 {
		_this.ToEvents = _this.Events
		_this.toEvents = _this.events
	} else {
		_this.toEvents = _this.PostDecodeEvents(_this.ToEvents)
	}

	if len(_this.FromCBE) == 0 && len(_this.FromCTE) == 0 && len(_this.FromEvents) == 0 {
		return _this.errorf("must specify source(s) to read from")
	}

	if len(_this.ToCBE) == 0 && len(_this.ToCTE) == 0 && len(_this.ToEvents) == 0 {
		return _this.errorf("must specify destination(s) to compare against")
	}

	return nil
}

func (_this *MustSucceedTest) Run() error {
	if _this.Skip {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("%v: %w", _this.context, v))
			default:
				panic(fmt.Errorf("%v: %v", _this.context, v))
			}
		}
	}()

	if len(_this.FromEvents) != 0 {
		if _this.Debug {
			fmt.Printf("%v: Convert events [%v] to ", _this.context, _this.fromEvents)
		}

		if len(_this.ToEvents) != 0 {
			if _this.Debug {
				fmt.Println("events")
			}
			result, err := _this.eventsToEvents(_this.fromEvents)
			if err != nil {
				return _this.wrapError(err, "converting events [%v] to events", _this.fromEvents)
			}
			if !_this.toEvents.AreEquivalentTo(result) {
				return _this.errorf("expected events [%v] to produce events [%v] but got [%v]",
					_this.fromEvents, _this.toEvents, result)
			}
		}

		if len(_this.ToCBE) != 0 {
			if _this.Debug {
				fmt.Println("CBE")
			}
			result, err := _this.eventsToCbe(_this.fromEvents)
			if err != nil {
				return _this.wrapError(err, "converting events [%v] to CBE", _this.fromEvents)
			}
			if !bytes.Equal(_this.ToCBE, result) {
				return _this.errorf("expected events [%v] to produce CBE [%v] but got [%v]",
					_this.fromEvents, asHex(_this.ToCBE), asHex(result))
			}
		}

		if len(_this.ToCTE) != 0 {
			if _this.Debug {
				fmt.Println("CTE")
			}
			result, err := _this.eventsToCte(_this.fromEvents)
			if err != nil {
				return _this.wrapError(err, "converting events [%v] to CTE", _this.fromEvents)
			}
			if _this.ToCTE != result {
				return _this.errorf("expected events [%v] to produce CTE [%v] but got [%v]",
					_this.fromEvents, _this.ToCTE, result)
			}
		}
	}

	if len(_this.FromCBE) != 0 {
		if _this.Debug {
			fmt.Printf("%v: Convert CBE [%v] to ", _this.context, asHex(_this.FromCBE))
		}

		if len(_this.ToEvents) != 0 {
			if _this.Debug {
				fmt.Println("events")
			}
			result, err := _this.cbeToEvents(_this.FromCBE)
			if err != nil {
				return _this.wrapError(err, "converting CBE [%v] to events", asHex(_this.FromCBE))
			}
			if !_this.toEvents.AreEquivalentTo(result) {
				return _this.errorf("expected CBE [%v] to produce events [%v] but got [%v]",
					asHex(_this.FromCBE), _this.toEvents, result)
			}
		}

		if len(_this.ToCBE) != 0 {
			if _this.Debug {
				fmt.Println("CBE")
			}
			result, err := _this.cbeToCbe(_this.FromCBE)
			if err != nil {
				return _this.wrapError(err, "converting CBE [%v] to CBE", asHex(_this.FromCBE))
			}
			if !bytes.Equal(_this.ToCBE, result) {
				return _this.errorf("expected CBE [%v] to produce CBE [%v] but got [%v]",
					asHex(_this.FromCBE), asHex(_this.ToCBE), asHex(result))
			}
		}

		if len(_this.ToCTE) != 0 {
			if _this.Debug {
				fmt.Println("CTE")
			}
			result, err := _this.cbeToCte(_this.FromCBE)
			if err != nil {
				return _this.wrapError(err, "converting CBE [%v] to CTE", asHex(_this.FromCBE))
			}
			if _this.ToCTE != result {
				return _this.errorf("expected CBE [%v] to produce CTE [%v] but got [%v]",
					asHex(_this.FromCBE), _this.ToCTE, result)
			}
		}
	}

	if len(_this.FromCTE) != 0 {
		if _this.Debug {
			fmt.Printf("%v: Convert CTE [%v] to ", _this.context, _this.FromCTE)
		}

		if len(_this.ToEvents) != 0 {
			if _this.Debug {
				fmt.Println("events")
			}
			result, err := _this.cteToEvents(_this.FromCTE)
			if err != nil {
				return _this.wrapError(err, "converting CTE [%v] to events", _this.FromCTE)
			}
			if !_this.toEvents.AreEquivalentTo(result) {
				return _this.errorf("expected CTE [%v] to produce events [%v] but got [%v]",
					_this.FromCTE, _this.toEvents, result)
			}
		}

		if len(_this.ToCBE) != 0 {
			if _this.Debug {
				fmt.Println("CBE")
			}
			result, err := _this.cteToCbe(_this.FromCTE)
			if err != nil {
				return _this.wrapError(err, "converting CTE [%v] to CBE", _this.FromCTE)
			}
			if !bytes.Equal(_this.ToCBE, result) {
				return _this.errorf("expected CTE [%v] to produce CBE [%v] but got [%v]",
					_this.FromCTE, asHex(_this.ToCBE), asHex(result))
			}
		}

		if len(_this.ToCTE) != 0 {
			if _this.Debug {
				fmt.Println("CTE")
			}
			result, err := _this.cteToCte(_this.FromCTE)
			if err != nil {
				return _this.wrapError(err, "converting CTE [%v] to CTE", _this.FromCTE)
			}
			if _this.ToCTE != result {
				return _this.errorf("expected CTE [%v] to produce CTE [%v] but got [%v]",
					_this.FromCTE, _this.ToCTE, result)
			}
		}
	}

	return nil
}

func (_this *MustSucceedTest) cteToCte(document string) (result string, err error) {
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

	decoderOpts := options.DefaultCEDecoderOptions()
	decoderOpts.DebugPanics = _this.Debug
	decoder := cte.NewDecoder(&decoderOpts)

	encoderOpts := options.DefaultCTEEncoderOptions()
	encoder := cte.NewEncoder(&encoderOpts)
	receiver := rules.NewRules(encoder, nil)

	inBuffer := bytes.NewBuffer([]byte(document))
	outBuffer := &strings.Builder{}
	encoder.PrepareToEncode(outBuffer)
	if err = decoder.Decode(inBuffer, receiver); err != nil {
		return
	}
	result = outBuffer.String()
	return
}

func (_this *MustSucceedTest) cbeToCte(document []byte) (result string, err error) {
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

	decoderOpts := options.DefaultCEDecoderOptions()
	decoderOpts.DebugPanics = _this.Debug
	decoder := cbe.NewDecoder(&decoderOpts)

	encoderOpts := options.DefaultCTEEncoderOptions()
	encoder := cte.NewEncoder(&encoderOpts)
	receiver := rules.NewRules(encoder, nil)

	inBuffer := bytes.NewBuffer([]byte(document))
	outBuffer := &strings.Builder{}
	encoder.PrepareToEncode(outBuffer)
	if err = decoder.Decode(inBuffer, receiver); err != nil {
		return
	}
	result = outBuffer.String()
	return
}

func (_this *MustSucceedTest) cteToCbe(document string) (result []byte, err error) {
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

	decoderOpts := options.DefaultCEDecoderOptions()
	decoderOpts.DebugPanics = _this.Debug
	decoder := cte.NewDecoder(&decoderOpts)

	encoderOpts := options.DefaultCBEEncoderOptions()
	encoder := cbe.NewEncoder(&encoderOpts)
	receiver := rules.NewRules(encoder, nil)

	inBuffer := bytes.NewBuffer([]byte(document))
	outBuffer := &bytes.Buffer{}
	encoder.PrepareToEncode(outBuffer)
	if err = decoder.Decode(inBuffer, receiver); err != nil {
		return
	}
	result = outBuffer.Bytes()
	return
}

func (_this *MustSucceedTest) cbeToCbe(document []byte) (result []byte, err error) {
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

	decoderOpts := options.DefaultCEDecoderOptions()
	decoderOpts.DebugPanics = _this.Debug
	decoder := cbe.NewDecoder(&decoderOpts)

	encoderOpts := options.DefaultCBEEncoderOptions()
	encoder := cbe.NewEncoder(&encoderOpts)
	receiver := rules.NewRules(encoder, nil)

	inBuffer := bytes.NewBuffer([]byte(document))
	outBuffer := &bytes.Buffer{}
	encoder.PrepareToEncode(outBuffer)
	if err = decoder.Decode(inBuffer, receiver); err != nil {
		return
	}
	result = outBuffer.Bytes()
	return
}

func (_this *MustSucceedTest) eventsToCte(events test.Events) (result string, err error) {
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

	encoder := cte.NewEncoder(nil)
	outBuffer := &strings.Builder{}
	encoder.PrepareToEncode(outBuffer)
	receiver := rules.NewRules(encoder, nil)
	receiver.OnBeginDocument()
	for _, event := range events {
		event.Invoke(receiver)
	}
	receiver.OnEndDocument()
	result = outBuffer.String()
	return
}

func (_this *MustSucceedTest) eventsToCbe(events test.Events) (result []byte, err error) {
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

	encoder := cbe.NewEncoder(nil)
	outBuffer := &bytes.Buffer{}
	encoder.PrepareToEncode(outBuffer)
	receiver := rules.NewRules(encoder, nil)
	receiver.OnBeginDocument()
	for _, event := range events {
		event.Invoke(receiver)
	}
	receiver.OnEndDocument()
	result = outBuffer.Bytes()
	return
}

func (_this *MustSucceedTest) eventsToEvents(events test.Events) (result test.Events, err error) {
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

	receiver, collection := test.NewEventCollector(nil)
	receiver = rules.NewRules(receiver, nil)
	receiver.OnBeginDocument()
	for _, event := range events {
		event.Invoke(receiver)
	}
	receiver.OnEndDocument()
	result = collection.Events
	return
}
