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

package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kstenerud/go-concise-encoding/test/test_runner"
)

func TestSuites(t *testing.T) {
	runTestsInPath(t, "suites")
}

func runTestsInPath(t *testing.T, testDir string) {
	filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".cte") {
			return nil
		}
		test_runner.RunTests(t, path)
		return nil
	})
}

func TestBugReportTemplates(t *testing.T) {
	// Make sure the bug report templates are valid
	test_runner.RunTests(t, "../bugreport/templates/incorrectly_allowed.cte")
	test_runner.RunTests(t, "../bugreport/templates/incorrectly_rejected.cte")

	// Make sure the default bug report test is valid
	test_runner.RunTests(t, "../bugreport/bugreport.cte")
}
