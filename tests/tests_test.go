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
	"testing"

	"github.com/kstenerud/go-concise-encoding/test/test_runner"
)

func TestPrimary(t *testing.T) {
	test_runner.RunTests(t, "version.cte")
	test_runner.RunTests(t, "cbe-basic.cte")
	test_runner.RunTests(t, "cbe-arrays.cte")
	test_runner.RunTests(t, "cbe-containers.cte")
	test_runner.RunTests(t, "cte-basic.cte")
	test_runner.RunTests(t, "cte-arrays.cte")
	test_runner.RunTests(t, "cte-containers.cte")
	test_runner.RunTests(t, "cte-comment.cte")
	test_runner.RunTests(t, "cte-complex.cte")
}

func TestTemplates(t *testing.T) {

	// Make sure the general test template is valid
	test_runner.RunTests(t, "template.cte")

	// Make sure the bug report templates are valid
	test_runner.RunTests(t, "../bugreport/templates/incorrectly_allowed.cte")
	test_runner.RunTests(t, "../bugreport/templates/incorrectly_rejected.cte")

	// Make sure the default bug report test is valid
	test_runner.RunTests(t, "../bugreport/bugreport.cte")
}

func TestExamples(t *testing.T) {
	test_runner.RunTests(t, "website-examples.cte")
	test_runner.RunTests(t, "ce-specification-examples.cte")
}

func TestGithubIssues(t *testing.T) {
	test_runner.RunTests(t, "github-issues.cte")
}
