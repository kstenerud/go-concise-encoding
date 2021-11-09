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
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/kstenerud/go-concise-encoding/ce"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Inspects files for inline CTE snippets (anything between [```cte] and [```]) and tries to parse them.\n")
		fmt.Printf("Usage: %v <files>\n", args[0])
		return
	}

	for i := 1; i < len(args); i++ {
		path := args[i]
		fi, err := os.Stat(path)
		if err != nil {
			fmt.Errorf("Could not open %v: %v\n", path, err)
			return
		}
		if fi.Mode().IsRegular() {
			inspectFile(path)
		} else {
			fmt.Printf("Skipping [%v] because it is not a file.\n", path)
		}
	}
}

func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\r', '\n', '\t':
		return true
	default:
		return false
	}
}

func addHeaderIfNeeded(data []byte) []byte {
	for len(data) > 0 && isWhitespace(data[0]) {
		data = data[1:]
	}
	if len(data) < 2 {
		return data
	}
	if data[0] == 'c' && (data[1] >= '0' && data[1] <= '9') {
		return data
	}

	return append([]byte{'c', '0', '\n'}, data...)
}

func inspectFile(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read file [%v]: %v\n", path, err)
	}

	fmt.Printf("Inspecting [%v]...\n", path)
	for _, snippet := range getSnippets(contents) {
		snippet = addHeaderIfNeeded(snippet)
		decoder := ce.NewCTEDecoder(nil)
		if err = decoder.DecodeDocument(snippet, ce.NewRules(nil, nil)); err != nil {

			// _, err := ce.UnmarshalCTE(bytes.NewBuffer(snippet), template, nil)
			// if err != nil {
			fmt.Printf("Snippet failed in %v:\n%v\n", path, err)
			fmt.Printf("----------------------------------------------------------------------\n")
			fmt.Printf("%v\n", string(snippet))
			fmt.Printf("======================================================================\n\n")
		} else {
			// fmt.Printf("%v\n----------------------------------------\n%v", string(snippet), describe.D(unmarshaled))
		}
	}
}

// var snippetMatcher = regexp.MustCompile("```cte\\s*(.*)```")
var snippetMatcher = regexp.MustCompile("(?s)```cte\\s*(.*?)```")

func getSnippets(data []byte) (snippets [][]byte) {
	for _, match := range snippetMatcher.FindAllSubmatch(data, -1) {
		snippets = append(snippets, match[1])
	}
	return
}
