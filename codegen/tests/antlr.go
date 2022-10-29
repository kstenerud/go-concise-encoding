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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const parserBasePath = "test/event_parser"

func generateAntlrCode(projectDir string) {
	javaPath, err := exec.LookPath("java")
	if err != nil {
		panic(err)
	}
	dstPath := filepath.Join(projectDir, parserBasePath, "parser")
	if err := os.RemoveAll(dstPath); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		panic(err)
	}

	antlrPath := filepath.Join(projectDir, "codegen", "antlr-4.10.1-complete.jar")
	lexerPath := filepath.Join(projectDir, "codegen", "tests", "CEEventLexer.g4")
	parserPath := filepath.Join(projectDir, "codegen", "tests", "CEEventParser.g4")
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
