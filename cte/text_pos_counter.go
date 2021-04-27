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

package cte

import (
	"fmt"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
)

type TextPositionCounter struct {
	Position    int
	LineCount   int
	ColumnCount int

	lastReadChar    chars.ByteWithEOF
	lastLineCount   int
	lastColumnCount int
}

// Advance with the specified character.
func (_this *TextPositionCounter) Advance(b chars.ByteWithEOF) {
	// TODO: Handle more line terminators?
	// LF:    Line Feed, U+000A
	// VT:    Vertical Tab, U+000B
	// FF:    Form Feed, U+000C
	// CR:    Carriage Return, U+000D
	// NEL:   Next Line, U+0085                      : C2 85
	// LS:    Line Separator, U+2028                 : E2 80 A8
	// PS:    Paragraph Separator, U+2029            : E2 80 A9

	_this.lastLineCount = _this.LineCount
	_this.lastColumnCount = _this.ColumnCount
	_this.lastReadChar = b

	switch b {
	case '\n':
		_this.LineCount++
		_this.ColumnCount = 0
	case chars.EOFMarker:
		// Do nothing
	default:
		_this.ColumnCount++
	}

	_this.Position++
}

// Retreat by a single character.
// This can handle retreating over a newline.
// DO NOT call this multiple times in a row!
func (_this *TextPositionCounter) RetreatOneChar() {
	_this.Position--
	_this.LineCount = _this.lastLineCount
	_this.ColumnCount = _this.lastColumnCount
}

// Retreat by the specified number of chars on the current line.
// DO NOT use this to backtrack over a newline!
func (_this *TextPositionCounter) Retreat(charCount int, currentByte chars.ByteWithEOF) {
	_this.Position -= charCount
	_this.ColumnCount -= charCount
	_this.lastReadChar = currentByte
}

func (_this *TextPositionCounter) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	panic(fmt.Errorf("offset %v (line %v, col %v): %v", _this.Position, _this.LineCount+1, _this.ColumnCount+1, msg))
}

func (_this *TextPositionCounter) UnexpectedChar(decoding string) {
	if _this.lastReadChar == chars.EOFMarker {
		_this.UnexpectedEOF(decoding)
	} else {
		_this.Errorf("unexpected [%v] while decoding %v", _this.DescribeCurrentChar(), decoding)
	}
}

func (_this *TextPositionCounter) UnexpectedEOF(decoding string) {
	_this.Errorf("unexpected end of document while decoding %v", decoding)
}

func (_this *TextPositionCounter) DescribeCurrentChar() string {
	b := _this.lastReadChar
	switch {
	case b == chars.EOFMarker:
		return "EOF"
	case b == ' ':
		return "SP"
	case b > ' ' && b <= '~':
		return fmt.Sprintf("%c", b)
	default:
		return fmt.Sprintf("0x%02x", b)
	}
}
