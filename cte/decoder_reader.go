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
	"io"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

type Reader struct {
	reader    io.Reader
	lastByte  chars.ByteWithEOF
	hasUnread bool
	isEOF     bool
	byteBuff  [1]byte

	TextPos TextPositionCounter

	token            []byte
	verbatimSentinel []byte
}

func NewReader(reader io.Reader) *Reader {
	_this := &Reader{}
	_this.Init(reader)
	return _this
}

// Init the read buffer. You may call this again to re-initialize the buffer.
func (_this *Reader) Init(reader io.Reader) {
	_this.reader = reader
	if cap(_this.token) == 0 {
		_this.token = make([]byte, 0, 64)
	}
	_this.hasUnread = false
	_this.isEOF = false
}

// Bytes

func (_this *Reader) readNext() {
	// Can't create local [1]byte here because it mallocs? WTF???
	if _, err := _this.reader.Read(_this.byteBuff[:]); err == nil {
		_this.lastByte = chars.ByteWithEOF(_this.byteBuff[0])
	} else {
		if err == io.EOF {
			_this.isEOF = true
			_this.lastByte = chars.EOFMarker
		} else {
			panic(err)
		}
	}
}

func (_this *Reader) UnreadByte() {
	if _this.hasUnread {
		panic("Cannot unread twice")
	}
	_this.hasUnread = true
	_this.TextPos.RetreatOneChar()
}

func (_this *Reader) PeekByteAllowEOF() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
		_this.hasUnread = true
	}
	return _this.lastByte
}

func (_this *Reader) ReadByteAllowEOF() chars.ByteWithEOF {
	if !_this.hasUnread {
		_this.readNext()
	}

	_this.hasUnread = false
	b := _this.lastByte
	_this.TextPos.Advance(b)
	return b
}

func (_this *Reader) PeekByteNoEOF() byte {
	b := _this.PeekByteAllowEOF()
	if b == chars.EOFMarker {
		_this.unexpectedEOF()
	}
	return byte(b)
}

func (_this *Reader) ReadByteNoEOF() byte {
	b := _this.ReadByteAllowEOF()
	if b == chars.EOFMarker {
		_this.unexpectedEOF()
	}
	return byte(b)
}

func (_this *Reader) AdvanceByte() {
	_this.ReadByteNoEOF()
}

func (_this *Reader) SkipWhileProperty(property chars.Properties) {
	for _this.ReadByteAllowEOF().HasProperty(property) {
	}
	_this.UnreadByte()
}

// Tokens

func (_this *Reader) TokenBegin() {
	_this.token = _this.token[:0]
}

func (_this *Reader) TokenStripLastByte() {
	_this.token = _this.token[:len(_this.token)-1]
}

func (_this *Reader) TokenStripLastBytes(count int) {
	_this.token = _this.token[:len(_this.token)-count]
}

func (_this *Reader) TokenAppendByte(b byte) {
	_this.token = append(_this.token, b)
}

func (_this *Reader) TokenAppendBytes(b []byte) {
	_this.token = append(_this.token, b...)
}

var emptyRuneBytes = []byte{0, 0, 0, 0, 0}

func (_this *Reader) TokenAppendRune(r rune) {
	if r < utf8.RuneSelf {
		_this.TokenAppendByte(byte(r))
	} else {
		pos := len(_this.token)
		_this.TokenAppendBytes(emptyRuneBytes)
		length := utf8.EncodeRune(_this.token[pos:], r)
		_this.token = _this.token[:pos+length]
	}
}

func (_this *Reader) TokenGet() Token {
	return _this.token
}

func (_this *Reader) TokenReadByteNoEOF() byte {
	b := _this.ReadByteNoEOF()
	_this.TokenAppendByte(b)
	return b
}

func (_this *Reader) TokenReadByteAllowEOF() chars.ByteWithEOF {
	b := _this.ReadByteAllowEOF()
	if b != chars.EOFMarker {
		_this.TokenAppendByte(byte(b))
	}
	return b
}

func (_this *Reader) TokenReadUntilAndIncludingByte(untilByte byte) {
	for {
		b := _this.ReadByteNoEOF()
		if b == untilByte {
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *Reader) TokenReadUntilPropertyNoEOF(property chars.Properties) {
	for {
		b := _this.ReadByteNoEOF()
		if chars.ByteHasProperty(b, property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(b)
	}
}

func (_this *Reader) TokenReadUntilPropertyAllowEOF(property chars.Properties) {
	for {
		b := _this.ReadByteAllowEOF()
		if b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

func (_this *Reader) TokenReadWhilePropertyAllowEOF(property chars.Properties) {
	for {
		b := _this.ReadByteAllowEOF()
		if !b.HasProperty(property) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
}

// Errors

func (_this *Reader) AssertAtObjectEnd(decoding string) {
	if !_this.PeekByteAllowEOF().HasProperty(chars.ObjectEnd) {
		_this.unexpectedChar(decoding)
	}
}

func (_this *Reader) errorf(format string, args ...interface{}) {
	_this.TextPos.Errorf(format, args...)
}

func (_this *Reader) unexpectedEOF() {
	_this.errorf("unexpected end of document")
}

func (_this *Reader) unexpectedError(err error, decoding string) {
	_this.errorf("unexpected error [%v] while decoding %v", err, decoding)
}

func (_this *Reader) unexpectedChar(decoding string) {
	_this.TextPos.UnexpectedChar(decoding)
}

// Decoders

func (_this *Reader) SkipWhitespace() {
	_this.SkipWhileProperty(chars.StructWS)
}

func (_this *Reader) ReadToken() Token {
	_this.TokenBegin()
	_this.TokenReadUntilPropertyAllowEOF(chars.ObjectEnd)
	return _this.TokenGet()
}

func (_this *Reader) ReadNamedValue() []byte {
	_this.TokenBegin()
	_this.TokenReadWhilePropertyAllowEOF(chars.AZ)
	namedValue := _this.TokenGet()
	if len(namedValue) == 0 {
		_this.unexpectedChar("name")
	}
	common.ASCIIBytesToLower(namedValue)
	return namedValue
}

func (_this *Reader) TokenReadVerbatimSequence() {
	_this.verbatimSentinel = _this.verbatimSentinel[:0]
	for {
		b := _this.ReadByteNoEOF()
		if chars.ByteHasProperty(b, chars.StructWS) {
			if b == '\r' {
				if _this.ReadByteNoEOF() != '\n' {
					_this.unexpectedChar("verbatim sentinel")
				}
			}
			break
		}
		_this.verbatimSentinel = append(_this.verbatimSentinel, b)
	}

	sentinelLength := len(_this.verbatimSentinel)

Outer:
	for {
		_this.TokenReadByteNoEOF()
		for i := 1; i <= sentinelLength; i++ {
			if _this.token[len(_this.token)-i] != _this.verbatimSentinel[sentinelLength-i] {
				continue Outer
			}
		}

		_this.TokenStripLastBytes(sentinelLength)
		return
	}
}

var nbspBytes = []byte{0xc0, 0xa0}
var shyBytes = []byte{0xc0, 0xad}

func (_this *Reader) TokenReadEscape() {
	escape := _this.ReadByteNoEOF()
	switch escape {
	case 't':
		_this.TokenAppendByte('\t')
	case 'n':
		_this.TokenAppendByte('\n')
	case 'r':
		_this.TokenAppendByte('\r')
	case '"', '*', '/', '<', '>', '\\', '|':
		_this.TokenAppendByte(escape)
	case '_':
		// Non-breaking space
		_this.TokenAppendBytes(nbspBytes)
	case '-':
		// Soft hyphen
		_this.TokenAppendBytes(shyBytes)
	case '\r', '\n':
		// Continuation
		_this.SkipWhitespace()
	case '0':
		_this.TokenAppendByte(0)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		length := int(escape - '0')
		codepoint := rune(0)
		for i := 0; i < length; i++ {
			b := _this.ReadByteNoEOF()
			switch {
			case chars.ByteHasProperty(b, chars.DigitBase10):
				codepoint = (codepoint << 4) | (rune(b) - '0')
			case chars.ByteHasProperty(b, chars.LowerAF):
				codepoint = (codepoint << 4) | (rune(b) - 'a' + 10)
			case chars.ByteHasProperty(b, chars.UpperAF):
				codepoint = (codepoint << 4) | (rune(b) - 'A' + 10)
			default:
				_this.unexpectedChar("unicode escape")
			}
		}

		_this.TokenAppendRune(codepoint)
	case '.':
		_this.TokenReadVerbatimSequence()
	default:
		_this.unexpectedChar("escape sequence")
	}
}

func (_this *Reader) ReadQuotedString() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '"':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()
		}
	}
}

func (_this *Reader) ReadIdentifier() []byte {
	_this.TokenBegin()
	for {
		b := _this.ReadByteAllowEOF()
		// Only do a per-byte check here. The rules will do a per-rune check.
		if b == chars.EOFMarker || (b < 0x80 && !chars.IsRuneValidIdentifier(rune(b))) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
	return _this.TokenGet()
}

func (_this *Reader) ReadMarkerIdentifier() []byte {
	_this.TokenBegin()
	for {
		b := _this.ReadByteAllowEOF()
		// Only do a per-byte check here. The rules will do a per-rune check.
		if b == chars.EOFMarker || (b < 0x80 && !chars.IsRuneValidMarkerID(rune(b))) {
			_this.UnreadByte()
			break
		}
		_this.TokenAppendByte(byte(b))
	}
	return _this.TokenGet()
}

func (_this *Reader) ReadStringArray() []byte {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '|':
			_this.TokenStripLastByte()
			return _this.TokenGet()
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()
		}
	}
}

func trimWhitespace(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.StructWS) {
		str = str[1:]
	}
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.StructWS) {
		str = str[:len(str)-1]
	}
	return str
}

func trimWhitespaceMarkupContent(str []byte) []byte {
	for len(str) > 0 && chars.ByteHasProperty(str[0], chars.StructWS) {
		str = str[1:]
	}
	hasTrailingWS := false
	for len(str) > 0 && chars.ByteHasProperty(str[len(str)-1], chars.StructWS) {
		str = str[:len(str)-1]
		hasTrailingWS = true
	}
	if hasTrailingWS {
		str = append(str, ' ')
	}
	return str
}

func trimWhitespaceMarkupEnd(str []byte) []byte {
	return trimWhitespace(str)
}

func (_this *Reader) ReadSingleLineComment() []byte {
	_this.TokenBegin()
	_this.TokenReadUntilAndIncludingByte('\n')
	contents := _this.TokenGet()

	return trimWhitespace(contents)
}

func (_this *Reader) ReadMultilineComment() ([]byte, nextType) {
	_this.TokenBegin()
	lastByte := _this.TokenReadByteNoEOF()

	for {
		firstByte := lastByte
		lastByte = _this.TokenReadByteNoEOF()

		if firstByte == '*' && lastByte == '/' {
			_this.TokenStripLastBytes(2)
			contents := _this.TokenGet()
			return trimWhitespace(contents), nextIsCommentEnd
		}

		if firstByte == '/' && lastByte == '*' {
			_this.TokenStripLastBytes(2)
			contents := _this.TokenGet()
			return trimWhitespace(contents), nextIsCommentBegin
		}
	}
}

func (_this *Reader) ReadMarkupContent() ([]byte, nextType) {
	_this.TokenBegin()
	for {
		b := _this.TokenReadByteNoEOF()
		switch b {
		case '<':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsMarkupBegin
		case '>':
			_this.TokenStripLastByte()
			return trimWhitespaceMarkupEnd(_this.TokenGet()), nextIsMarkupEnd
		case '/':
			switch _this.TokenReadByteAllowEOF() {
			case '*':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsCommentBegin
			case '/':
				_this.TokenStripLastBytes(2)
				return trimWhitespaceMarkupContent(_this.TokenGet()), nextIsSingleLineComment
			}
		case '\\':
			_this.TokenStripLastByte()
			_this.TokenReadEscape()

		}
	}
}

// ============================================================================

// Internal

type nextType int

const (
	nextIsCommentBegin nextType = iota
	nextIsCommentEnd
	nextIsSingleLineComment
	nextIsMarkupBegin
	nextIsMarkupEnd
)

var subsecondMagnitudes = []int{
	1000000000,
	100000000,
	10000000,
	1000000,
	100000,
	10000,
	1000,
	100,
	10,
	1,
}
