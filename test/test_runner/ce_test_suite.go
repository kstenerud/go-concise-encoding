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
	"fmt"
	"testing"
)

type CETestSuite struct {
	TestFile string
	Options  struct {
		FailFast bool
	}
	CETests  []*CETestRunner
	CTETests []*CTETestRunner
	CBETests []*CBETestRunner
}

func (_this *CETestSuite) String() string {
	return _this.TestFile
}

func (_this *CETestSuite) postDecodeInit() {
	postDecodeCETest := func(test *CETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "malformed CE test at index %v (%v)", index, test)
		}()
		test.postDecodeInit()
	}
	for index, test := range _this.CETests {
		postDecodeCETest(test, index)
	}

	postDecodeCTETest := func(test *CTETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "malformed CTE test at index %v (%v)", index, test)
		}()
		test.postDecodeInit()
	}
	for index, test := range _this.CTETests {
		postDecodeCTETest(test, index)
	}

	postDecodeCBETest := func(test *CBETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "malformed CBE test at index %v (%v)", index, test)
		}()
		test.postDecodeInit()
	}
	for index, test := range _this.CBETests {
		postDecodeCBETest(test, index)
	}
}

func (_this *CETestSuite) validate() {
	validateCETest := func(test *CETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "CE test at index %v (%v) failed validation", index, test)
		}()
		test.validate()
	}
	for index, test := range _this.CETests {
		validateCETest(test, index)
	}

	validateCTETest := func(test *CTETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "CTE test at index %v (%v) failed validation", index, test)
		}()
		test.validate()
	}
	for index, test := range _this.CTETests {
		validateCTETest(test, index)
	}

	validateCBETest := func(test *CBETestRunner, index int) {
		defer func() {
			wrapPanic(recover(), "CBE test at index %v (%v) failed validation", index, test)
		}()
		test.validate()
	}
	for index, test := range _this.CBETests {
		validateCBETest(test, index)
	}
}

func (_this *CETestSuite) run(t *testing.T) {

	fmt.Printf("Running test suite: %v\n", _this.TestFile)

	errorCount := 0

	runCETest := func(test *CETestRunner, index int) (success bool) {
		defer func() { wrapPanic(recover(), "CE test at index %v (%v)", index, test) }()
		defer func() {
			success = !reportAnyError(recover(), test.Trace, "❌ %v, CE test %v: %v", _this.TestFile, index, test)
		}()
		test.run()
		return
	}
	for index, test := range _this.CETests {
		if runCETest(test, index) {
			fmt.Printf("✅ %v, CE test %v: %v\n", _this.TestFile, index, test)
		} else {
			errorCount++
			if _this.Options.FailFast {
				fmt.Println("FailFast enabled - stopping tests.")
				t.Error("Test failed")
				return
			}
		}
	}

	runCTETest := func(test *CTETestRunner, index int) (success bool) {
		defer func() { wrapPanic(recover(), "CE test at index %v (%v)", index, test) }()
		success = true
		if !test.Trace {
			defer func() {
				success = !reportAnyError(recover(), test.Trace, "❌ %v, CTE test %v: %v", _this.TestFile, index, test)
			}()
		}
		test.run()
		return
	}
	for index, test := range _this.CTETests {
		if runCTETest(test, index) {
			fmt.Printf("✅ %v, CTE test %v: %v\n", _this.TestFile, index, test)
		} else {
			errorCount++
			if _this.Options.FailFast {
				fmt.Println("FailFast enabled - stopping tests.")
				t.Error("Test failed")
				return
			}
		}
	}

	runCBETest := func(test *CBETestRunner, index int) (success bool) {
		defer func() { wrapPanic(recover(), "CE test at index %v (%v)", index, test) }()
		defer func() {
			success = !reportAnyError(recover(), test.Trace, "❌ %v, CBE test %v: %v", _this.TestFile, index, test)
		}()
		test.run()
		return
	}
	for index, test := range _this.CBETests {
		if runCBETest(test, index) {
			fmt.Printf("✅ %v, CBE test %v: %v\n", _this.TestFile, index, test)
		} else {
			errorCount++
			if _this.Options.FailFast {
				fmt.Println("FailFast enabled - stopping tests.")
				t.Error("Test failed")
				return
			}
		}
	}

	if errorCount > 0 {
		t.Error("Test failed")
	}
}
