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
	"fmt"
)

type MustFailTest struct {
	BaseTest
}

func (_this *MustFailTest) PostDecodeInit(ceVersion int, context string, index int) error {
	if _this.Skip {
		return nil
	}

	// A bit hacky, but it makes things easier to have Debug in the base class.
	_this.Debug = false

	context = fmt.Sprintf(`%v, "must fail" test #%v`, context, index+1)
	if err := _this.BaseTest.PostDecodeInit(ceVersion, context); err != nil {
		return err
	}

	total := 0
	if len(_this.CTE) > 0 {
		total++
	}
	if len(_this.CBE) > 0 {
		total++
	}
	if len(_this.Events) > 0 {
		total++
	}

	if total != 1 {
		return _this.errorf("must have one and only one of: cbe, cte, events")
	}
	return nil
}

func (_this *MustFailTest) Run() error {
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

	if len(_this.CTE) > 0 {
		if events, err := _this.cteToEvents(_this.CTE); err == nil {
			return _this.errorf("expected CTE document [%v] to fail but it produced [%v]",
				_this.CTE, events)
		}
	} else {
		if events, err := _this.cbeToEvents(_this.CBE); err == nil {
			return _this.errorf("expected CBE document [%v] to fail but it produced [%v]",
				asHex(_this.CBE), events)
		}
	}
	return nil
}
