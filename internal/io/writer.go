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

	"github.com/kstenerud/go-concise-encoding/internal/chars"
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
	_this.Buffer = make([]byte, writerStartBufferSize)
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

func (_this *Writer) FlushBuffer(count int) {
	_this.writeBytes(_this.Buffer[:count])
}

func (_this *Writer) FlushBufferPortion(start, end int) {
	_this.writeBytes(_this.Buffer[start:end])
}

func (_this *Writer) WriteByte(b byte) {
	_this.Buffer[0] = b
	_this.FlushBuffer(1)
}

func (_this *Writer) WriteRune(r rune) {
	n := utf8.EncodeRune(_this.Buffer, r)
	_this.FlushBuffer(n)
}

func (_this *Writer) WriteBytes(b []byte) {
	_this.writeBytes(b)
}

func (_this *Writer) WriteIdentifier(b []byte) {
	_this.WriteByte(byte(len(b)))
	_this.WriteBytes(b)
}

func (_this *Writer) WriteDecimalInt(value int64) {
	uintValue := uint64(value)
	if value < 0 {
		_this.WriteByte('-')
		uintValue = uint64(-value)
	}
	_this.WriteDecimalUint(uintValue)
}

func (_this *Writer) WriteDecimalUint(value uint64) {
	const maxUint64Digits = 21 // 18446744073709551615
	_this.ExpandBuffer(maxUint64Digits)
	var start int
	for start = maxUint64Digits - 1; start >= 0; start-- {
		_this.Buffer[start] = byte('0' + value%10)
		value /= 10
		if value == 0 {
			break
		}
	}
	_this.FlushBufferPortion(start, maxUint64Digits)
}

func (_this *Writer) WriteDecimalUintLeftLoaded(value uint64, digitCount int) {
	_this.ExpandBuffer(digitCount)
	lastNonZero := 0
	for i := digitCount - 1; i >= 0; i-- {
		digit := value % 10
		if lastNonZero == 0 && digit != 0 {
			lastNonZero = i
		}
		_this.Buffer[i] = byte('0' + digit)
		value /= 10
	}
	_this.FlushBuffer(lastNonZero + 1)
}

func (_this *Writer) WriteDecimalUintDigits(value uint64, digits int) {
	_this.ExpandBuffer(digits)
	for i := digits - 1; i >= 0; i-- {
		_this.Buffer[i] = byte('0' + value%10)
		value /= 10
	}
	_this.FlushBuffer(digits)
}

func (_this *Writer) WriteHexByte(value byte) {
	_this.ExpandBuffer(2)
	_this.Buffer[0] = chars.HexChars[value>>4]
	_this.Buffer[1] = chars.HexChars[value&15]
	_this.FlushBuffer(2)
}

func (_this *Writer) WriteFmt(format string, args ...interface{}) {
	_this.WriteString(fmt.Sprintf(format, args...))
}

// Add a formatted string, but strip the specified number of characters from the
// beginning of the result before adding.
func (_this *Writer) WriteFmtStripped(stripByteCount int, format string, args ...interface{}) {
	_this.WriteString(fmt.Sprintf(format, args...)[stripByteCount:])
}

func (_this *Writer) WriteString(str string) {
	if _, err := _this.stringWriter.WriteString(str); err != nil {
		panic(err)
	}
}

func (_this *Writer) writeBytes(b []byte) {
	if _, err := _this.writer.Write(b); err != nil {
		panic(err)
	}
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
	_this.writer.FlushBuffer(n)
	return
}
