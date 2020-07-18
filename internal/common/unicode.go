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

package common

type CharProperty uint8

const (
	CharNoProperties = iota << 1
	CharIsControl
	CharIsWhitespace
	CharIsNumeralOrLookalike
	CharIsLowSymbolOrLookalike
	CharIsURIUnsafe
	CharIsQuoteUnsafe
	CharIsMarkupUnsafe
)

func GetProperty(char rune) CharProperty {
	return charProperties[char]
}

var charProperties = map[rune]CharProperty{}

func init() {
	// _ - . : / are allowed in unquoted strings
	// Whitespace, control chars disallowed in unquoted strings
	// 0-9 disallowed as first char of unquoted strings
	// " \ and control chars must be escaped in quoted strings
	// " % must be escaped in URIs
	// /* */ < > \ ` must be escaped in markup (exception: entity refs)

	// Special considerations:
	// - quoted string continuation \
	// - markup comment begin/end /* */
	// - markup entity ref

	assignRange := func(start, end rune, property CharProperty) {
		for i := start; i <= end; i++ {
			charProperties[i] = property
		}
	}

	// Symbols
	charProperties['!'] = CharIsLowSymbolOrLookalike
	charProperties['"'] = CharIsLowSymbolOrLookalike | CharIsQuoteUnsafe | CharIsURIUnsafe
	charProperties['#'] = CharIsLowSymbolOrLookalike
	charProperties['$'] = CharIsLowSymbolOrLookalike
	charProperties['%'] = CharIsLowSymbolOrLookalike | CharIsURIUnsafe
	charProperties['&'] = CharIsLowSymbolOrLookalike
	charProperties['\''] = CharIsLowSymbolOrLookalike
	charProperties['('] = CharIsLowSymbolOrLookalike
	charProperties[')'] = CharIsLowSymbolOrLookalike
	charProperties['*'] = CharIsLowSymbolOrLookalike
	charProperties['+'] = CharIsLowSymbolOrLookalike
	charProperties[','] = CharIsLowSymbolOrLookalike
	charProperties['-'] = CharNoProperties
	charProperties['.'] = CharNoProperties
	charProperties['/'] = CharNoProperties
	charProperties[':'] = CharNoProperties
	charProperties[';'] = CharIsLowSymbolOrLookalike
	charProperties['<'] = CharIsLowSymbolOrLookalike | CharIsMarkupUnsafe
	charProperties['='] = CharIsLowSymbolOrLookalike
	charProperties['>'] = CharIsLowSymbolOrLookalike | CharIsMarkupUnsafe
	charProperties['?'] = CharIsLowSymbolOrLookalike
	charProperties['@'] = CharIsLowSymbolOrLookalike
	charProperties['['] = CharIsLowSymbolOrLookalike
	charProperties['\\'] = CharIsLowSymbolOrLookalike | CharIsQuoteUnsafe | CharIsMarkupUnsafe
	charProperties[']'] = CharIsLowSymbolOrLookalike
	charProperties['^'] = CharIsLowSymbolOrLookalike
	charProperties['_'] = CharNoProperties
	charProperties['`'] = CharIsLowSymbolOrLookalike | CharIsMarkupUnsafe
	charProperties['{'] = CharIsLowSymbolOrLookalike
	charProperties['|'] = CharIsLowSymbolOrLookalike
	charProperties['}'] = CharIsLowSymbolOrLookalike
	charProperties['~'] = CharIsLowSymbolOrLookalike

	assignRange('0', '9', CharIsNumeralOrLookalike)

	// Latin punctuation https://unicode.org/charts/PDF/U0080.pdf
	charProperties[0x00a6] = CharIsLowSymbolOrLookalike // |
	charProperties[0x00b4] = CharIsLowSymbolOrLookalike // '

	// General punctuation https://unicode.org/charts/PDF/U2000.pdf
	assignRange(0x2018, 0x201f, CharIsLowSymbolOrLookalike) // "
	assignRange(0x2032, 0x2037, CharIsLowSymbolOrLookalike) // "
	charProperties[0x2039] = CharIsLowSymbolOrLookalike     // <
	charProperties[0x203a] = CharIsLowSymbolOrLookalike     // >
	charProperties[0x203c] = CharIsLowSymbolOrLookalike     // !
	assignRange(0x2047, 0x2049, CharIsLowSymbolOrLookalike) // !?
	charProperties[0x204e] = CharIsLowSymbolOrLookalike     // *
	charProperties[0x2052] = CharIsLowSymbolOrLookalike     // %
	charProperties[0x2055] = CharIsLowSymbolOrLookalike     // *
	charProperties[0x2057] = CharIsLowSymbolOrLookalike     // "

	// Mathematical operators https://unicode.org/charts/PDF/U2200.pdf
	charProperties[0x2217] = CharIsLowSymbolOrLookalike // *
	charProperties[0x2223] = CharIsLowSymbolOrLookalike // |
	charProperties[0x223c] = CharIsLowSymbolOrLookalike // ~

	// Miscellaneous technical https://unicode.org/charts/PDF/U2300.pdf
	charProperties[0x239c] = CharIsLowSymbolOrLookalike // |
	charProperties[0x239f] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23a2] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23a5] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23aa] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23b8] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23b9] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23d0] = CharIsLowSymbolOrLookalike // |
	charProperties[0x23fd] = CharIsLowSymbolOrLookalike // |

	// CJK vertical forms https://unicode.org/charts/PDF/UFE10.pdf
	charProperties[0xfe10] = CharIsLowSymbolOrLookalike     // '
	charProperties[0xfe11] = CharIsLowSymbolOrLookalike     // '
	assignRange(0xfe14, 0xfe16, CharIsLowSymbolOrLookalike) // ;!?

	// CJK compatibility https://unicode.org/charts/PDF/UFE30.pdf
	charProperties[0xfe31] = CharIsLowSymbolOrLookalike // |
	charProperties[0xfe33] = CharIsLowSymbolOrLookalike // |
	charProperties[0xfe45] = CharIsLowSymbolOrLookalike // '
	charProperties[0xfe46] = CharIsLowSymbolOrLookalike // '

	// CJK small form variants https://unicode.org/charts/PDF/UFE50.pdf
	assignRange(0xfe50, 0xfe6b, CharIsLowSymbolOrLookalike) // symbols
	charProperties[0xfe52] = CharNoProperties               // .
	charProperties[0xfe53] = CharNoProperties               // reserved
	charProperties[0xfe55] = CharNoProperties               // :
	charProperties[0xfe58] = CharNoProperties               // -
	charProperties[0xfe63] = CharNoProperties               // -
	charProperties[0xfe67] = CharNoProperties               // reserved

	// CJK halfwidth, fullwidth https://unicode.org/charts/PDF/UFF00.pdf
	assignRange(0xff00, 0xff0f, CharIsLowSymbolOrLookalike) // symbols
	assignRange(0xff10, 0xff19, CharIsNumeralOrLookalike)   // 0-9
	assignRange(0xff1a, 0xff20, CharIsLowSymbolOrLookalike) // symbols
	assignRange(0xff3b, 0xff40, CharIsLowSymbolOrLookalike) // symbols
	assignRange(0xff5b, 0xff5e, CharIsLowSymbolOrLookalike) // symbols
	charProperties[0xffe4] = CharIsLowSymbolOrLookalike     // |
	charProperties[0xffe8] = CharIsLowSymbolOrLookalike     // |
	charProperties[0xff0d] = CharNoProperties               // -
	charProperties[0xff0e] = CharNoProperties               // .
	charProperties[0xff0f] = CharNoProperties               // /
	charProperties[0xff1a] = CharNoProperties               // :
	charProperties[0xff3f] = CharNoProperties               // _

	// Ancient symbols https://unicode.org/charts/PDF/U10190.pdf
	charProperties[0x10190] = CharIsLowSymbolOrLookalike // =

	// Ideopgraphic punctuation https://unicode.org/charts/PDF/U16FE0.pdf
	charProperties[0x16fe4] = CharIsLowSymbolOrLookalike // invisible

	// Musical notation https://unicode.org/charts/PDF/U1D100.pdf
	charProperties[0x1d1c1] = CharIsLowSymbolOrLookalike // |
	charProperties[0x1d1c2] = CharIsLowSymbolOrLookalike // |

	// Mathematical alphanumeric symbols https://unicode.org/charts/PDF/U1D400.pdf
	assignRange(0x1d7ce, 0x1d7ff, CharIsNumeralOrLookalike) // 0-9
}
