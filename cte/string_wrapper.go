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

	"github.com/kstenerud/go-concise-encoding/internal/common"
)

// Examines a string destined for a CTE document, wrapping it in quotes,
// escapes, or verbatim sequences as necessary.
func wrapString(str string) (finalString string) {
	if str == "" {
		return `""`
	}
	unquotedUnsafe := common.CharIsControl | common.CharIsWhitespace | common.CharIsLowSymbolOrLookalike
	unquotedUnsafeFirstChar := unquotedUnsafe | common.CharIsNumeralOrLookalike
	quotedUnsafe := common.CharIsControl | common.CharIsQuoteUnsafe

	requiresQuotes := false
	escapeCount := 0

	if common.GetCharProperty([]rune(str)[0])&unquotedUnsafeFirstChar != 0 {
		requiresQuotes = true
	}

	for _, ch := range str {
		props := common.GetCharProperty(ch)
		if props&unquotedUnsafe != 0 {
			requiresQuotes = true
		}
		if props&quotedUnsafe != 0 {
			escapeCount++
		}
	}

	// If more than 1/8 of the string requires escaping, use a verbatim string
	if len(str) > 8 && (escapeCount<<3) > len(str) {
		return convertToVerbatimString(str)
	}

	if requiresQuotes {
		if escapeCount == 0 {
			return `"` + str + `"`
		}
		return escapedQuoted(str, escapeCount)
	}

	return str
}

func escapedQuoted(str string, escapeCount int) string {
	quotedUnsafe := common.CharIsControl | common.CharIsQuoteUnsafe
	var sb strings.Builder
	// Worst case scenario: All characters that require escaping need a unicode
	// sequence `\uxxxx`. In this case, we'd need at least 6 bytes per escaped
	// character.
	sb.Grow(len([]byte(str)) + escapeCount*5 + 2)
	sb.WriteByte('"')
	for _, ch := range str {
		if common.GetCharProperty(ch)&quotedUnsafe == 0 {
			// Note: WriteRune always returns a nil error
			sb.WriteRune(ch)
		} else {
			// Note: WriteString always returns a nil error
			sb.WriteString(escapeChar(ch))
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

func escapeChar(ch rune) string {
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
	return fmt.Sprintf(`\u%04x`, ch)
}

func generateAlphabet() string {
	// Ordered from least common to most common, chosen to not be confused by
	// a human with other CTE document structural characters.
	return "#$%&*+/:;=^_|~23456789ZQXJVKBPYGCFMWULDHSNOIRATE10zqxjvkbpygcfmwuldhsnoirate"
}

var alphabet = []byte(generateAlphabet())

func generateVerbatimSentinel(str string) string {
	// Try all 1, 2, and 3-character sequences picked from a safe alphabet.

	usedChars := [256]bool{}
	for _, ch := range []byte(str) {
		usedChars[ch] = true
	}

	for _, ch := range alphabet {
		if !usedChars[ch] {
			return fmt.Sprintf("%c", ch)
		}
	}

	for _, ch0 := range alphabet {
		for _, ch1 := range alphabet {
			sentinel := fmt.Sprintf("%c%c", ch0, ch1)
			if !strings.Contains(str, sentinel) {
				return sentinel
			}
		}
	}

	for _, ch0 := range alphabet {
		for _, ch1 := range alphabet {
			for _, ch2 := range alphabet {
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
	panic(fmt.Errorf("Could not generate verbatim sentinel for malicious string [%v]", str))
}

var verbatimIdentifierChars = []rune{}

func convertToVerbatimString(str string) string {
	sentinel := generateVerbatimSentinel(str)
	return "`" + sentinel + str + sentinel
}
