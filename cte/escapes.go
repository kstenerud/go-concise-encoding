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
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
)

func getStringRequirements(str []byte) (escapeCount int, requiresQuotes bool) {
	if len(str) == 0 {
		return 0, true
	}

	firstRune, _ := utf8.DecodeRune(str)
	if chars.RuneHasProperty(firstRune, chars.CharNeedsQuoteFirst) {
		requiresQuotes = true
	}

	for _, ch := range string(str) {
		props := chars.GetRuneProperty(ch)
		if props.HasProperty(chars.CharNeedsQuote) {
			requiresQuotes = true
		}
		if props.HasProperty(chars.CharNeedsEscapeQuoted) {
			escapeCount++
		}
	}
	return
}

func needsEscapesStringLikeArray(str []byte) bool {
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeArray) {
			return true
		}
	}
	return false
}

func needsEscapesMarkup(str []byte) bool {
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeMarkup) {
			return true
		}
	}
	return false
}

func containsEscapes(str []byte) bool {
	for _, b := range str {
		if b == '\\' {
			return true
		}
	}
	return false
}

// Wraps a string destined for a CTE document, adding quotes or escapes as
// necessary.
func asPotentialQuotedString(str []byte) (finalString []byte) {
	if len(str) == 0 {
		return []byte{'"', '"'}
	}

	escapeCount, requiresQuotes := getStringRequirements(str)

	if !requiresQuotes {
		return str
	}

	if escapeCount == 0 {
		finalString = make([]byte, len(str)+2)
		copy(finalString[1:], str)
		finalString[0] = '"'
		finalString[len(finalString)-1] = '"'
		return finalString
	}

	return asEscapedQuotedString(str, escapeCount)
}

// Wraps a string-encoded array destined for a CTE document.
func asStringArray(elementType []byte, str []byte) []byte {
	if needsEscapesStringLikeArray(str) {
		str = asEscapedStringArrayContent(str)
	}

	bytes := make([]byte, len(str)+len(elementType)+3)
	bytes[0] = '|'
	bytes[len(bytes)-1] = '|'
	copy(bytes[1:], elementType)
	bytes[1+len(elementType)] = ' '
	copy(bytes[2+len(elementType):], str)
	return bytes
}

// Possibly escapes a string-encoded array destined for a CTE document.
func asStringArrayContents(str []byte) []byte {
	if needsEscapesStringLikeArray(str) {
		return asEscapedStringArrayContent(str)
	}

	return str
}

// Wraps markup content destined for a CTE document.
func asMarkupContents(str []byte) []byte {
	if needsEscapesMarkup(str) {
		return asEscapedMarkupContents(str)
	}

	return str
}

// ============================================================================

func asEscapedQuotedString(str []byte, escapeCount int) []byte {
	var bb bytes.Buffer
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	bb.Grow(len(str) + escapeCount*6 + 2)
	// Note: StringBuilder's WriteXYZ() always return nil errors
	bb.WriteByte('"')
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeQuoted) {
			bb.Write(escapeCharQuoted(ch))
		} else {
			bb.WriteRune(ch)
		}
	}
	bb.WriteByte('"')
	return bb.Bytes()
}

func escapeCharQuoted(ch rune) []byte {
	switch ch {
	case '\t':
		return []byte(`\t`)
	case '\r':
		return []byte(`\r`)
	case '\n':
		return []byte(`\n`)
	case '"':
		return []byte(`\"`)
	case '*':
		return []byte(`\*`)
	case '/':
		return []byte(`\/`)
	case '\\':
		return []byte(`\\`)
	}
	return unicodeEscape(ch)
}

func asEscapedStringArrayContent(str []byte) []byte {
	var bb bytes.Buffer
	bb.Grow(len(str))
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeArray) {
			// Note: StringBuilder's WriteXYZ() always return nil errors
			bb.Write(escapeCharStringArray(ch))
		} else {
			bb.WriteRune(ch)
		}
	}
	return bb.Bytes()
}

func unicodeEscape(ch rune) []byte {
	hex := fmt.Sprintf("%x", ch)
	return []byte(fmt.Sprintf("\\%d%s", len(hex), hex))
}

func escapeCharStringArray(ch rune) []byte {
	switch ch {
	case '|':
		return []byte(`\|`)
	case '\\':
		return []byte(`\\`)
	case '\t':
		return []byte(`\t`)
	case '\r':
		return []byte(`\r`)
	case '\n':
		return []byte(`\n`)
	}
	return unicodeEscape(ch)
}

func asEscapedMarkupContents(str []byte) []byte {
	var bb bytes.Buffer
	bb.Grow(len(str))
	// Note: StringBuilder's WriteXYZ() always return nil errors
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeMarkup) {
			bb.Write(escapeCharMarkup(ch))
		} else {
			bb.WriteRune(ch)
		}
	}
	return bb.Bytes()
}

func escapeCharMarkup(ch rune) []byte {
	switch ch {
	case '*':
		// TODO: Check ahead for /* */ instead of blindly escaping
		return []byte(`\*`)
	case '/':
		// TODO: Check ahead for /* */ instead of blindly escaping
		return []byte(`\/`)
	case '<':
		return []byte(`\<`)
	case '>':
		return []byte(`\>`)
	case 0xa0:
		return []byte(`\_`)
	case 0xad:
		return []byte(`\-`)
	case '\\':
		return []byte(`\\`)
	}
	return unicodeEscape(ch)
}

// Ordered from least common to most common, chosen to not be confused by
// a human with other CTE document structural characters.
var verbatimSentinelAlphabet = []byte("~%*+;=^_23456789ZQXJVKBPYGCFMWULDHSNOIRATE10zqxjvkbpygcfmwuldhsnoirate")

func generateVerbatimSentinel(str []byte) []byte {
	// Try all 1, 2, and 3-character sequences picked from a safe alphabet.

	usedChars := [256]bool{}
	for _, ch := range str {
		usedChars[ch] = true
	}

	var sentinelBuff [3]byte

	for _, ch := range verbatimSentinelAlphabet {
		if !usedChars[ch] {
			return []byte{ch}
		}
	}

	for _, ch0 := range verbatimSentinelAlphabet {
		for _, ch1 := range verbatimSentinelAlphabet {
			sentinelBuff[0] = ch0
			sentinelBuff[1] = ch1
			sentinel := sentinelBuff[:2]
			if !bytes.Contains(str, sentinel) {
				return sentinel
			}
		}
	}

	for _, ch0 := range verbatimSentinelAlphabet {
		for _, ch1 := range verbatimSentinelAlphabet {
			for _, ch2 := range verbatimSentinelAlphabet {
				sentinelBuff[0] = ch0
				sentinelBuff[1] = ch1
				sentinelBuff[2] = ch2
				sentinel := sentinelBuff[:3]
				if !bytes.Contains(str, sentinel) {
					return sentinel
				}
			}
		}
	}

	// If we're here, all 450,000 three-character sequences have occurred in
	// the string. At this point, we conclude that it's a specially crafted
	// attack string, and not naturally occurring.
	panic(fmt.Errorf("could not generate verbatim sentinel for malicious string [%v]", str))
}
