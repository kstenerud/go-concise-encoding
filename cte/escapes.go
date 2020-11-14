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
	"strings"

	"github.com/kstenerud/go-concise-encoding/internal/chars"
)

func containsEscapes(str string) bool {
	for _, b := range []byte(str) {
		if b == '\\' {
			return true
		}
	}
	return false
}

// Wraps a string destined for a CTE document, adding quotes or escapes as
// necessary.
func asPotentialQuotedString(str []byte) (finalString string) {
	asString := string(str)
	if asString == "" {
		return `""`
	}

	requiresQuotes := false
	escapeCount := 0

	if chars.RuneHasProperty([]rune(asString)[0], chars.CharNeedsQuoteFirst) {
		requiresQuotes = true
	}

	for _, ch := range asString {
		props := chars.GetRuneProperty(ch)
		if props.HasProperty(chars.CharNeedsQuote) {
			requiresQuotes = true
		}
		if props.HasProperty(chars.CharNeedsEscapeQuoted) {
			escapeCount++
		}
	}

	if !requiresQuotes {
		return asString
	}

	if escapeCount == 0 {
		return `"` + asString + `"`
	}

	return escapedQuoted(asString, escapeCount)
}

// Wraps a string-encoded array destined for a CTE document.
func asStringArray(elementType string, str []byte) string {
	for _, ch := range string(str) {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeArray) {
			return escapedStringArray(elementType, string(str))
		}
	}

	return "|" + elementType + " " + string(str) + "|"
}

// Wraps markup content destined for a CTE document.
func asMarkupContent(str []byte) string {
	asString := string(str)
	for _, ch := range asString {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeMarkup) {
			return escapedMarkupText(asString)
		}
	}

	return asString
}

// ============================================================================

func escapedQuoted(str string, escapeCount int) string {
	var sb strings.Builder
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence. In this case, we'd need at least 7 bytes per escaped character.
	sb.Grow(len([]byte(str)) + escapeCount*6 + 2)
	// Note: StringBuilder's WriteXYZ() always return nil errors
	sb.WriteByte('"')
	for _, ch := range str {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeQuoted) {
			sb.WriteString(escapeCharQuoted(ch))
		} else {
			sb.WriteRune(ch)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

func escapeCharQuoted(ch rune) string {
	switch ch {
	case '\t':
		return `\t`
	case '\r':
		return `\r`
	case '\n':
		return `\n`
	case '"':
		return `\"`
	case '*':
		return `\*`
	case '/':
		return `\/`
	case '\\':
		return `\\`
	}
	return unicodeEscape(ch)
}

func escapedStringArray(elementType string, str string) string {
	var sb strings.Builder
	sb.Grow(len([]byte(str)))
	// Note: StringBuilder's WriteXYZ() always return nil errors
	sb.WriteByte('|')
	sb.WriteString(elementType)
	sb.WriteByte(' ')
	for _, ch := range str {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeArray) {
			sb.WriteString(escapeCharStringArray(ch))
		} else {
			sb.WriteRune(ch)
		}
	}
	sb.WriteByte('|')
	return sb.String()
}

func unicodeEscape(ch rune) string {
	hex := fmt.Sprintf("%x", ch)
	return fmt.Sprintf("\\%d%s", len(hex), hex)
}

func escapeCharStringArray(ch rune) string {
	switch ch {
	case '|':
		return `\|`
	case '\\':
		return `\\`
	case '\t':
		return `\t`
	case '\r':
		return `\r`
	case '\n':
		return `\n`
	}
	return unicodeEscape(ch)
}

func escapedMarkupText(str string) string {
	var sb strings.Builder
	sb.Grow(len([]byte(str)))
	// Note: StringBuilder's WriteXYZ() always return nil errors
	for _, ch := range str {
		if chars.RuneHasProperty(ch, chars.CharNeedsEscapeMarkup) {
			sb.WriteString(escapeCharMarkup(ch))
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

func escapeCharMarkup(ch rune) string {
	switch ch {
	case '*':
		// TODO: Check ahead for /* */ instead of blindly escaping
		return `\*`
	case '/':
		// TODO: Check ahead for /* */ instead of blindly escaping
		return `\/`
	case '<':
		return `\<`
	case '>':
		return `\>`
	case 0xa0:
		return `\_`
	case 0xad:
		return `\-`
	case '\\':
		return `\\`
	}
	return unicodeEscape(ch)
}

// Ordered from least common to most common, chosen to not be confused by
// a human with other CTE document structural characters.
var verbatimSentinelAlphabet = []byte("~%*+;=^_23456789ZQXJVKBPYGCFMWULDHSNOIRATE10zqxjvkbpygcfmwuldhsnoirate")

func generateVerbatimSentinel(str string) string {
	// Try all 1, 2, and 3-character sequences picked from a safe alphabet.

	usedChars := [256]bool{}
	for _, ch := range []byte(str) {
		usedChars[ch] = true
	}

	for _, ch := range verbatimSentinelAlphabet {
		if !usedChars[ch] {
			return fmt.Sprintf("%c", ch)
		}
	}

	for _, ch0 := range verbatimSentinelAlphabet {
		for _, ch1 := range verbatimSentinelAlphabet {
			sentinel := fmt.Sprintf("%c%c", ch0, ch1)
			if !strings.Contains(str, sentinel) {
				return sentinel
			}
		}
	}

	for _, ch0 := range verbatimSentinelAlphabet {
		for _, ch1 := range verbatimSentinelAlphabet {
			for _, ch2 := range verbatimSentinelAlphabet {
				sentinel := fmt.Sprintf("%c%c%c", ch0, ch1, ch2)
				if !strings.Contains(str, sentinel) {
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

func unescape(str []byte) []byte {
	iDst := 0
	for iSrc := 0; iSrc < len(str); iSrc++ {
		b := str[iSrc]
		if b == '\\' {
			iSrc++
			if iSrc >= len(str) {
				// TODO: Error
			}
			b := str[iSrc]
			switch b {
			case 't':
				str[iDst] = '\t'
			case 'n':
				str[iDst] = '\n'
			case 'r':
				str[iDst] = '\r'
			case '"', '*', '/', '<', '>', '\\', '|':
				str[iDst] = b
			case '_':
				// Non-breaking space
				str[iDst] = 0xc2
				iDst++
				str[iDst] = 0xa0
			case '-':
				// Soft hyphen
				str[iDst] = 0xc2
				iDst++
				str[iDst] = 0xad
			case '\r', '\n':
				// TODO: Continuation
			case '0':
				str[iDst] = 0
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// TODO: Unicode escape
			case '.':
				// TODO: Verbatim
			default:
				// TODO: Error
			}
		} else {
			str[iDst] = b
		}
		iDst++
	}
	return str[:iDst]
}
