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
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
)

var testsImports = []*common.Import{
	{As: "", Import: "fmt"},
	{As: "", Import: "math"},
	{As: "", Import: "math/big"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/configuration"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/test"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/test/test_runner"},
}

func generateTestGenerators(basePath string) {
	common.GenerateGoFile(basePath, testsImports, func(writer io.Writer) {
		generateArrayTestGenerator(basePath, writer)
	})
}

func readFunctionTemplate(basePath string, functionName string) string {
	var templateData string
	path := path.Join(basePath, "templates.go")
	if file, err := os.ReadFile(path); err != nil {
		panic(fmt.Errorf("could not read function template file %v: %w", path, err))
	} else {
		templateData = string(file)
	}

	iStart := strings.Index(templateData, fmt.Sprintf("\nfunc %v(", functionName))
	if iStart < 0 {
		panic(fmt.Errorf("could not find function %v in template file %v", functionName, path))
	}

	templateData = templateData[iStart+1:]
	iEnd := strings.Index(templateData, "\n}\n")
	if iEnd < 0 {
		panic(fmt.Errorf("could not find end of function %v in template file %v", functionName, path))
	}

	return templateData[:iEnd+3]
}

func generateArrayTestGenerator(basePath string, writer io.Writer) {
	template := readFunctionTemplate(basePath, "generateArrayInt32Tests")
	writer.Write([]byte(strings.ReplaceAll(template, "32", "8")))
	writer.Write([]byte(strings.ReplaceAll(template, "32", "16")))
	writer.Write([]byte(strings.ReplaceAll(template, "32", "64")))

	template = readFunctionTemplate(basePath, "generateArrayUint32Tests")
	writer.Write([]byte(strings.ReplaceAll(template, "32", "8")))
	writer.Write([]byte(strings.ReplaceAll(template, "32", "16")))
	writer.Write([]byte(strings.ReplaceAll(template, "32", "64")))
}
