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
)

type MustSucceedTest struct {
	BaseTest
	LossyCTE    bool `ce:"order=101"`
	LossyCBE    bool `ce:"order=102"`
	LossyEvents bool `ce:"order=103"`
}

func (_this *MustSucceedTest) PostDecodeInit(ceVersion int, context string, index int) error {
	if _this.Skip {
		return nil
	}
	context = fmt.Sprintf(`%v, "must succeed" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	if len(_this.Cbe) == 0 && len(_this.Cte) == 0 {
		return _this.errorf("must have cbe and/or cte")
	}

	return nil
}

func (_this *MustSucceedTest) runCte() error {
	hasEvents := len(_this.events) > 0

	if _this.Debug {
		fmt.Printf("%v: Convert CTE to events: [%v]", _this.context, _this.Cte)
	}
	events, err := _this.cteToEvents(_this.Cte)
	if err != nil {
		return _this.wrapError(err, "decoding CTE [%v]", _this.Cte)
	}
	if hasEvents && !_this.LossyCTE {
		if !_this.events.AreEquivalentTo(events) {
			return _this.errorf("expected CTE [%v] to produce events [%v] but got [%v]",
				_this.Cte, _this.events, events)
		}
	}
	if _this.Debug {
		fmt.Printf("%v: Convert events to CTE: [%v]", _this.context, _this.events)
	}
	if hasEvents {
		events = _this.events
	}
	document, err := _this.eventsToCte(events)
	if err != nil {
		return _this.wrapError(err, "Encoding events [%v] to CTE", events)
	}
	if !_this.LossyCTE {
		if _this.Cte != document {
			return _this.errorf("re-encoding events [%v] from CTE [%v] produced unexpected CTE [%v]",
				events, _this.Cte, document)
		}
	}

	return nil
}

func (_this *MustSucceedTest) runCbe() error {
	hasEvents := len(_this.events) > 0

	if _this.Debug {
		fmt.Printf("%v: Convert CBE to events: [%v]", _this.context, asHex(_this.Cbe))
	}
	events, err := _this.cbeToEvents(_this.Cbe)
	if err != nil {
		return _this.wrapError(err, "decoding CBE [%v]", asHex(_this.Cbe))
	}
	if hasEvents && !_this.LossyCBE {
		if !_this.events.AreEquivalentTo(events) {
			return _this.errorf("expected CBE [%v] to produce events [%v] but got [%v]",
				asHex(_this.Cbe), _this.events, events)
		}
	}
	if _this.Debug {
		fmt.Printf("%v: Convert events to CBE: [%v]", _this.context, _this.events)
	}
	if hasEvents {
		events = _this.events
	}
	document, err := _this.eventsToCbe(events)
	if err != nil {
		return _this.wrapError(err, "Encoding events [%v] to CBE", events)
	}
	if !_this.LossyCBE {
		if !bytes.Equal(_this.Cbe, document) {
			return _this.errorf("re-encoding events [%v] from CBE [%v] produced unexpected CBE [%v]",
				events, asHex(_this.Cbe), asHex(document))
		}
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

	if len(_this.Cte) > 0 {
		if err := _this.runCte(); err != nil {
			return err
		}
	}

	if len(_this.Cbe) > 0 {
		if err := _this.runCbe(); err != nil {
			return err
		}
	}

	return nil
}
