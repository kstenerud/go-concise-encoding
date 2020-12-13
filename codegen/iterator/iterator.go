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

package iterator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

var types = [...]string{
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"float32",
	"float64",
	"bool",
	// No go types for these:
	// "float16",
	// "UUID",
}

func GenerateCode(projectDir string) {
	tmpl, err := template.New("").Parse(templateContents)
	if err != nil {
		panic(err)
	}

	generatedFilePath := filepath.Join(projectDir, "iterator/generated-do-not-edit.go")
	writer, err := os.Create(generatedFilePath)
	if err != nil {
		panic(fmt.Errorf("could not open %s: %v", generatedFilePath, err))
	}
	defer writer.Close()

	_, err = writer.WriteString(codeHeader)
	if err != nil {
		panic(fmt.Errorf("could not write to %s: %v", generatedFilePath, err))
	}

	sort.Strings(types[:])

	for _, typeName := range types {
		data := &generatorData{
			TypeName:  typeName,
			TypeTitle: strings.Title(typeName),
		}
		if err = tmpl.Execute(writer, data); err != nil {
			panic(fmt.Errorf("could not write template to %s: %v", generatedFilePath, err))
		}
	}
}

type generatorData struct {
	TypeName  string
	TypeTitle string
}

var codeHeader = standard.Header + `package iterator

import (
	"github.com/kstenerud/go-concise-encoding/options"
)
`

var templateContents = `
// {{.TypeName}}

type {{.TypeName}}ArrayIterator struct{}

func new{{.TypeTitle}}ArrayIterator() ObjectIterator                                                { return &{{.TypeName}}ArrayIterator{} }
func (_this *{{.TypeName}}ArrayIterator) InitTemplate(_ FetchIterator)                             {}
func (_this *{{.TypeName}}ArrayIterator) NewInstance() ObjectIterator                              { return _this }
func (_this *{{.TypeName}}ArrayIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {}

type {{.TypeName}}SliceIterator struct{}

func new{{.TypeTitle}}SliceIterator() ObjectIterator                                                { return &{{.TypeName}}SliceIterator{} }
func (_this *{{.TypeName}}SliceIterator) InitTemplate(_ FetchIterator)                             {}
func (_this *{{.TypeName}}SliceIterator) NewInstance() ObjectIterator                              { return _this }
func (_this *{{.TypeName}}SliceIterator) InitInstance(_ FetchIterator, _ *options.IteratorOptions) {}
`
