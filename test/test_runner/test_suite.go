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
	"os"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/version"
)

const TestSuiteVersion = 1

type TestSuiteType struct {
	Identifier string `ce:"order=1"`
	Version    int    `ce:"order=2"`
}

type TestSuite struct {
	Type      TestSuiteType `ce:"order=1"`
	CEVersion *int          `ce:"order=2"`
	Tests     []*UnitTest   `ce:"order=3"`
	context   string
}

func (_this *TestSuite) PostDecodeInit(sourceFile string) (errors []error) {
	_this.context = sourceFile

	if len(_this.Type.Identifier) == 0 {
		return []error{_this.errorf("missing type identifier field")}
	}

	if _this.Type.Identifier != "ce-test" {
		return []error{_this.errorf("%v: unrecognized type identifier", _this.Type.Identifier)}
	}

	if _this.Type.Version != TestSuiteVersion {
		return []error{_this.errorf("ce test format version %v runner cannot load from version %v file", TestSuiteVersion, _this.Type.Version)}
	}

	if _this.CEVersion == nil {
		return []error{_this.errorf("missing ceversion field")}
	}

	if *_this.CEVersion != version.ConciseEncodingVersion {
		return []error{_this.errorf("This codec is for Concise Encoding version %v and cannot run tests for version %v",
			version.ConciseEncodingVersion, _this.CEVersion)}
	}

	for index, test := range _this.Tests {
		if nextErrors := test.PostDecodeInit(*_this.CEVersion, _this.context, index); nextErrors != nil {
			errors = append(errors, nextErrors...)
		}
	}
	return
}

func (_this *TestSuite) Run() (errors []error) {
	for _, test := range _this.Tests {
		if nextErrors := test.Run(); nextErrors != nil {
			errors = append(errors, nextErrors...)
		}
	}
	return
}

func (_this *TestSuite) errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%v: %v", _this.context, message)
}

func loadTestSuite(testDescriptorFile string) (suite *TestSuite, errors []error) {
	defer func() { wrapPanic(recover(), "while loading test suite from %v", testDescriptorFile) }()

	file, err := os.Open(testDescriptorFile)
	if err != nil {
		panic(fmt.Errorf("unexpected error opening test suite file %v: %w", testDescriptorFile, err))
	}

	config := configuration.DefaultCEUnmarshalerConfiguration()
	config.DebugPanics = true
	loadedTest, err := ce.UnmarshalCTE(file, suite, &config)
	if err != nil {
		panic(fmt.Errorf("malformed unit test: Unexpected CTE decode error in test suite file %v: %w", testDescriptorFile, err))
	}

	suite = loadedTest.(*TestSuite)
	errors = suite.PostDecodeInit(testDescriptorFile)

	return
}
