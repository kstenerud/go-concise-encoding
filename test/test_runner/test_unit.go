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

type UnitTest struct {
	Name        string
	MustSucceed []*MustSucceedTest
	MustFail    []*MustFailTest
}

func (_this *UnitTest) PostDecodeInit(ceVersion int, context string, testIndex int) (errors []error) {
	context = fmt.Sprintf("%v, unit test #%v (%v)", context, testIndex+1, _this.Name)

	for index, test := range _this.MustSucceed {
		if err := test.PostDecodeInit(ceVersion, context, index); err != nil {
			errors = append(errors, err)
		}
	}

	for index, test := range _this.MustFail {
		if err := test.PostDecodeInit(ceVersion, context, index); err != nil {
			errors = append(errors, err)
		}
	}

	return
}

func (_this *UnitTest) Run() (errors []error) {
	for _, test := range _this.MustSucceed {
		if err := test.Run(); err != nil {
			errors = append(errors, err)
		}
	}

	for _, test := range _this.MustFail {
		if err := test.Run(); err != nil {
			errors = append(errors, err)
		}
	}

	return
}
