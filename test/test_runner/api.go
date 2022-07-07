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
	"os"
	"testing"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/debug"
)

// Run a CE test suite, described by the CTE document at testDescriptorFile
func RunCEUnitTests(t *testing.T, testDescriptorFile string) {
	testSuite := loadTestSuite(testDescriptorFile)

	defer func() { wrapPanic("while running test suite %v", testSuite) }()

	debug.DebugOptions.PassThroughPanics = true
	testSuite.run(t)
	debug.DebugOptions.PassThroughPanics = false
}

func loadTestSuite(testDescriptorFile string) *CETestSuite {
	defer func() { wrapPanic("while loading test suite from %v", testDescriptorFile) }()

	file, err := os.Open(testDescriptorFile)
	if err != nil {
		panic(fmt.Errorf("unexpected error opening test suite file %v: %v", testDescriptorFile, err))
	}
	var testSuite *CETestSuite

	testSuiteIntf, err := ce.UnmarshalCTE(file, testSuite, nil)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Unexpected CTE decode error in test suite file %v: %w", testDescriptorFile, err))
	}
	testSuite = testSuiteIntf.(*CETestSuite)
	testSuite.TestFile = testDescriptorFile
	testSuite.postDecodeInit()
	testSuite.validate()

	return testSuite
}
