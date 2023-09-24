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

// Snippet check looks through documents for CTE snippets (```cte ... ```)
// and attempts to parse it in order to verify that it is valid.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/kstenerud/go-concise-encoding/ce"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-describe"
)

func main() {
	quiet := flag.Bool("q", false, "quiet")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	args := flag.Args()
	verbosityLevel := getVerbosityLevel(*quiet, *verbose)

	if len(args) < 1 {
		printUsage()
		return
	}

	for _, path := range args {
		fi, err := os.Stat(path)
		if err != nil {
			fmt.Printf("Could not open %v: %v\n", path, err)
			return
		}
		if fi.Mode().IsRegular() {
			inspectFile(path, verbosityLevel)
		} else {
			fmt.Printf("Skipping %v because it is not a file.\n", path)
		}
	}
}

func printUsage() {
	fmt.Printf("Inspects files for inline CTE snippets (anything between [```cte] and [```]) and tries to parse them.\n")
	fmt.Printf("Usage: %v [opts] <files>\n", os.Args[0])
	flag.PrintDefaults()
}

type verbosity int

const (
	verbosityQuiet verbosity = iota
	verbosityNormal
	verbosityVerbose
)

func getVerbosityLevel(quiet bool, verbose bool) verbosity {
	if quiet {
		return verbosityQuiet
	}
	if verbose {
		return verbosityVerbose
	}
	return verbosityNormal
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
	if (data[0] == 'c' || data[0] == 'C') && (data[1] >= '0' && data[1] <= '9') {
		return data
	}

	return append([]byte{'c', '0', '\n'}, data...)
}

func reportError(snippet []byte, err error) {
	fmt.Printf("======================================================================\n")
	fmt.Printf("📜 Snippet:\n%v\n", string(snippet))
	fmt.Printf("----------------------------------------------------------------------\n")
	fmt.Printf("❌ Failed: %v\n", err)
	fmt.Printf("======================================================================\n")
}

func reportSuccess(snippet []byte, unmarshaled interface{}) {
	fmt.Printf("======================================================================\n")
	fmt.Printf("📜 Snippet:\n%v\n", string(snippet))
	fmt.Printf("----------------------------------------------------------------------\n")
	if unmarshaled != nil {
		fmt.Printf("✅ Unmarshaled to:\n%v\n", describe.Describe(unmarshaled, 4))
	} else {
		fmt.Printf("✅ Success\n")
	}
	fmt.Printf("======================================================================\n")
}

func inspectFile(path string, verbosityLevel verbosity) {
	contents, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read file [%v]: %v\n", path, err)
	}

	if verbosityLevel >= verbosityNormal {
		fmt.Printf("Inspecting %v\n", path)
	}

	config := configuration.New()
	for _, snippet := range getSnippets(contents) {
		snippet = addHeaderIfNeeded(snippet)
		unmarshaled, err := ce.UnmarshalCTE(bytes.NewBuffer(snippet), nil, config)
		if err != nil {
			decoder := ce.NewCTEDecoder(config)
			if err = decoder.DecodeDocument(snippet, ce.NewRules(nil, config)); err != nil {
				reportError(snippet, err)
				continue
			}

		}

		if verbosityLevel >= verbosityVerbose {
			reportSuccess(snippet, unmarshaled)
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
