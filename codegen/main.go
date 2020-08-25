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

// Package build generates code for other parts of the library. The lack of
// generics and inheritance makes a number of things tedious and error prone,
// which these generators attempt to deal with.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kstenerud/go-concise-encoding/codegen/builder"
	"github.com/kstenerud/go-concise-encoding/codegen/iterator"
	"github.com/kstenerud/go-concise-encoding/codegen/unicode"
)

func main() {
	projectPath := getProjectPath()

	unicodePath := flag.String("unicode", "", "/path/to/ucd.all.flat.xml. Get it from https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip")
	flag.Parse()

	if *unicodePath != "" {
		unicode.GenerateCode(projectPath, *unicodePath)
	}

	builder.GenerateCode(projectPath)
	iterator.GenerateCode(projectPath)
}

func printHelp() {
	ex, err := os.Executable()
	if err != nil {
		ex = "build"
	}
	fmt.Printf(`Usage: %s [/path/to/ucd.all.flat.xml]

/path/to/ucd.all.flat.xml, ff specified, can be retrieved from the latest https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip.
`, ex)
}

func getUnicodeXMLPath() string {
	if len(os.Args) == 2 {
		return os.Args[1]
	}
	return ""
}

func getExePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func getProjectPath() string {
	return filepath.Dir(getExePath())
}
