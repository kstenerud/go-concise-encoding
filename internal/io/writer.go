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

package io

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Writer struct {
	writer       io.Writer
	stringWriter io.StringWriter
	adapter      StringWriterAdapter
	Buffer       []byte
}

func NewWriter() *Writer {
	_this := &Writer{}
	_this.Init()
	return _this
}

func (_this *Writer) Init() {
	_this.Buffer = make([]byte, writerStartBufferSize, writerStartBufferSize)
	_this.adapter.Init(_this)
}

func (_this *Writer) SetWriter(writer io.Writer) {
	_this.writer = writer
	if sw, ok := _this.writer.(io.StringWriter); ok {
		_this.stringWriter = sw
	} else {
		_this.stringWriter = &_this.adapter
	}
}

func (_this *Writer) ExpandBuffer(size int) {
	if len(_this.Buffer) < size {
		_this.Buffer = make([]byte, size*2)
	}
}

func (_this *Writer) FlushBuffer(start int, end int) {
	_this.WriteBytes(_this.Buffer[start:end])
}

func (_this *Writer) WriteByte(b byte) {
	_this.Buffer[0] = b
	_this.FlushBuffer(0, 1)
}

func (_this *Writer) WriteRune(r rune) {
	n := utf8.EncodeRune(_this.Buffer, r)
	_this.FlushBuffer(0, n)
}

func (_this *Writer) WriteString(str string) {
	_this.stringWriter.WriteString(str)
}

func (_this *Writer) WriteBytes(b []byte) {
	_this.writer.Write(b)
}

func (_this *Writer) WriteFmt(format string, args ...interface{}) {
	_this.WriteString(fmt.Sprintf(format, args...))
}

// Add a formatted string, but strip the specified number of characters from the
// beginning of the result before adding.
func (_this *Writer) WriteFmtStripped(stripByteCount int, format string, args ...interface{}) {
	_this.WriteString(fmt.Sprintf(format, args...)[stripByteCount:])
}

// ============================================================================

const writerStartBufferSize = 32

type StringWriterAdapter struct {
	writer *Writer
}

func (_this *StringWriterAdapter) Init(writer *Writer) {
	_this.writer = writer
}

func (_this *StringWriterAdapter) WriteString(str string) (n int, err error) {
	n = len(str)
	_this.writer.ExpandBuffer(n)
	copy(_this.writer.Buffer, str)
	_this.writer.FlushBuffer(0, n)
	return
}
