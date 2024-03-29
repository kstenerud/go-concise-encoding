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

import (
	"fmt"
	"testing"
)

func RunTests(t *testing.T, sourceFile string) {
	suite, errors := loadTestSuite(sourceFile)
	if len(errors) > 0 {
		fmt.Printf("❌ %v\n", sourceFile)
		reportFailedTestLoad(errors)
		t.Fail()
		return
	}

	errors = suite.Run()
	if len(errors) > 0 {
		fmt.Printf("❌ %v\n", sourceFile)
		reportTestFailures(errors)
		t.Fail()
		return
	}
	fmt.Printf("✅ %v\n", sourceFile)
}

func reportFailedTestLoad(errors []error) {
	for _, err := range errors {
		fmt.Printf("  ❗ %v\n", err)
	}
}

func reportTestFailures(errors []error) {
	for _, err := range errors {
		fmt.Printf("  💣 %v\n", err)
	}
}
