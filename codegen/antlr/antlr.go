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

package antlr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunAntlr(antlrPath string, grammarPath string, dstPath string) {

	javaPath, err := exec.LookPath("java")
	if err != nil {
		panic(err)
	}
	if err := os.RemoveAll(dstPath); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		panic(err)
	}

	antlrPath = filepath.Join(antlrPath, "antlr-4.12.0-complete.jar")
	lexerPath, err := findLexerFile(grammarPath)
	if err != nil {
		panic(err)
	}
	parserPath, err := findParserFile(grammarPath)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(
		javaPath,
		"-cp", antlrPath,
		"org.antlr.v4.Tool",
		"-o", dstPath,
		"-Dlanguage=Go",
		lexerPath, parserPath,
	)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("failed to run %v: %w\nStdout = [%v]\nStderr = [%v]", cmd.Args, err, stdout.String(), stderr.String()))
	}
}

func findLexerFile(basePath string) (path string, err error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	for _, finfo := range files {
		if strings.HasSuffix(finfo.Name(), "Lexer.g4") {
			return filepath.Join(basePath, finfo.Name()), nil
		}
	}
	return "", fmt.Errorf("could not find *Lexer.g4 in %v", basePath)
}

func findParserFile(basePath string) (path string, err error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	for _, finfo := range files {
		if strings.HasSuffix(finfo.Name(), "Parser.g4") {
			return filepath.Join(basePath, finfo.Name()), nil
		}
	}
	return "", fmt.Errorf("could not find *Parser.g4 in %v", basePath)
}
