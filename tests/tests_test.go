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

func TestNew(t *testing.T) {
	// test_runner.RunTests(t, "version.cte")
	// test_runner.RunTests(t, "cbe-basic.cte")
	// test_runner.RunTests(t, "cbe-arrays.cte")
	// test_runner.RunTests(t, "cbe-containers.cte")
	// test_runner.RunTests(t, "cte-basic.cte")
	test_runner.RunTests(t, "cte-arrays.cte")
}

func TestTemplates(t *testing.T) {

	// Make sure the general test template is valid
	test_runner.RunCEUnitTests(t, "template.cte")

	// Make sure the bug report templates are valid
	test_runner.RunCEUnitTests(t, "../bugreport/templates/cbe_decoded_incorrectly.cte")
	test_runner.RunCEUnitTests(t, "../bugreport/templates/cbe_output_incorrect.cte")
	test_runner.RunCEUnitTests(t, "../bugreport/templates/cte_decoded_incorrectly.cte")
	test_runner.RunCEUnitTests(t, "../bugreport/templates/cte_output_incorrect.cte")
	test_runner.RunCEUnitTests(t, "../bugreport/templates/doc_wrongfully_allowed.cte")
	test_runner.RunCEUnitTests(t, "../bugreport/templates/doc_wrongfully_rejected.cte")

	// Make sure the default bug report test is valid
	test_runner.RunCEUnitTests(t, "../bugreport/bugreport.cte")
}

func TestCE(t *testing.T) {
	test_runner.RunCEUnitTests(t, "ce-basic.cte")
}

func TestCTEBasic(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cte-basic.cte")
}

func TestCTEArrays(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cte-arrays.cte")
}

func TestCTEComment(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cte-comment.cte")
}

func TestCTEComplex(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cte-complex.cte")
}

func TestCTEContainers(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cte-containers.cte")
}

func TestCBEBasic(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cbe-basic.cte")
}

func TestCBEArrays(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cbe-arrays.cte")
}

func TestCBEContainers(t *testing.T) {
	test_runner.RunCEUnitTests(t, "cbe-containers.cte")
}

func TestWebsiteExamples(t *testing.T) {
	test_runner.RunCEUnitTests(t, "website-examples.cte")
}

func TestCESpecificationExamples(t *testing.T) {
	test_runner.RunCEUnitTests(t, "ce-specification-examples.cte")
}

func TestGithubIssues(t *testing.T) {
	test_runner.RunCEUnitTests(t, "github-issues.cte")
}
