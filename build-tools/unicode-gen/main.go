// This program is used to populate `charProperties` and `asciiProperties` in
// github.com/kstenerud/go-concise-encoding/internal/common/unicode-generated.go
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// _ - . : are allowed in unquoted strings
// Whitespace, control chars disallowed in unquoted strings
// 0-9 disallowed as first char of unquoted strings
// " \ and control chars must be escaped in quoted strings
// " % must be escaped in URIs
// /* */ < > \ ` must be escaped in markup (exception: entity refs)

// Special considerations:
// - quoted string continuation \
// - markup comment begin/end /* */
// - markup entity ref

func main() {
	if len(os.Args) != 2 {
		printUsage(os.Stderr)
		os.Exit(1)
	}

	xmlPath := os.Args[1]
	db, err := loadDB(xmlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading [%v]: %v\n", xmlPath, err)
		os.Exit(1)
	}

	properties := CharProperties{}

	charIsNumeralOrLookalike := CharIsUnquotedFirstCharUnsafe | CharIsPrintable
	charIsWhitespace := CharIsUnquotedUnsafe | CharIsPrintable
	charIsLowSymbolOrLookalike := CharIsUnquotedFirstCharUnsafe | CharIsUnquotedUnsafe | CharIsPrintable
	charIsControl := CharIsUnquotedFirstCharUnsafe | CharIsUnquotedUnsafe |
		CharIsQuotedUnsafe | CharIsMarkupUnsafe | CharIsCustomTextUnsafe
	charIsReserved := CharIsUnquotedFirstCharUnsafe | CharIsUnquotedUnsafe |
		CharIsQuotedUnsafe | CharIsMarkupUnsafe | CharIsCustomTextUnsafe

	properties.Add(CharIsQuotedUnsafe, '"', '\\')
	properties.Add(CharIsCustomTextUnsafe, '"', '\\')
	properties.Add(CharIsMarkupUnsafe, '<', '>', '\\', '`')
	properties.AddRange(charIsNumeralOrLookalike, '0', '9')

	properties.Add(charIsWhitespace, db.GetRunesWithCriteria(func(char *Char) bool {
		switch char.Codepoint {
		case '\r', '\n', '\t':
			return true
		case 0x2028: // Line separator
			return false
		case 0x2029: // Paragraph separator
			return false
		case 0x1680: // Ogham space mark
			return false
		}
		if char.Codepoint < 0x20 {
			return false
		}
		return char.Category == "Zs" || char.Category == "Zl" || char.Category == "Zp"
	})...)

	properties.Add(charIsControl, db.GetRunesWithCriteria(func(char *Char) bool {
		switch char.Codepoint {
		case '\r', '\n', '\t':
			return false
		case 0x2028: // Line separator
			return true
		case 0x2029: // Paragraph separator
			return true
		}
		return char.Category == "Cc" || char.Category == "Cf"
	})...)

	properties.Add(charIsLowSymbolOrLookalike,
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', ';', '<',
		'=', '>', '?', '@', '[', '\\', ']', '^', '`', '{', '|', '}', '~')

	properties.Clear('-', '.', ':', '_') // Unquoted safe

	// Latin punctuation https://unicode.org/charts/PDF/U0080.pdf
	properties.Add(charIsLowSymbolOrLookalike,
		0x00a6, // |
		0x00b4, // '
	)

	// General punctuation https://unicode.org/charts/PDF/U2000.pdf
	properties.AddRange(charIsLowSymbolOrLookalike, 0x2018, 0x201f) // "
	properties.AddRange(charIsLowSymbolOrLookalike, 0x2032, 0x2037) // "
	properties.Add(charIsLowSymbolOrLookalike,
		0x2039, // <
		0x203a, // >
		0x203c, // !
	)
	properties.AddRange(charIsLowSymbolOrLookalike, 0x2047, 0x2049) // !?
	properties.Add(charIsLowSymbolOrLookalike,
		0x204e, 0x2055, // *
		0x2052, // %
		0x2057, // "
	)

	// Mathematical operators https://unicode.org/charts/PDF/U2200.pdf
	properties.Add(charIsLowSymbolOrLookalike,
		0x2217, // *
		0x2223, // |
		0x223c, // ~
	)

	// Miscellaneous technical https://unicode.org/charts/PDF/U2300.pdf
	properties.Add(charIsLowSymbolOrLookalike,
		0x239c, 0x239f, 0x23a2, 0x23a5, 0x23aa, 0x23b8, 0x23b9, 0x23d0, 0x23fd, // |
	)

	// CJK vertical forms https://unicode.org/charts/PDF/UFE10.pdf
	properties.Add(charIsLowSymbolOrLookalike, 0xfe10, 0xfe11)      // '
	properties.AddRange(charIsLowSymbolOrLookalike, 0xfe14, 0xfe16) // ;!?

	// CJK compatibility https://unicode.org/charts/PDF/UFE30.pdf
	properties.Add(charIsLowSymbolOrLookalike, 0xfe31, 0xfe33) // |
	properties.Add(charIsLowSymbolOrLookalike, 0xfe45, 0xfe46) // '

	// CJK small form variants https://unicode.org/charts/PDF/UFE50.pdf
	properties.AddRange(charIsLowSymbolOrLookalike, 0xfe50, 0xfe6b) // symbols
	properties.Set(charIsReserved, 0xfe53, 0xfe67)
	properties.Clear(
		0xfe52, // .
		0xfe55, // :
		0xfe58, // -
		0xfe63, // -
	)

	// CJK halfwidth, fullwidth https://unicode.org/charts/PDF/UFF00.pdf
	properties.AddRange(charIsLowSymbolOrLookalike, 0xff00, 0xff0f) // symbols
	properties.AddRange(charIsNumeralOrLookalike, 0xff10, 0xff19)   // 0-9
	properties.AddRange(charIsLowSymbolOrLookalike, 0xff1a, 0xff20) // symbols
	properties.AddRange(charIsLowSymbolOrLookalike, 0xff3b, 0xff40) // symbols
	properties.AddRange(charIsLowSymbolOrLookalike, 0xff5b, 0xff5e) // symbols
	properties.Add(charIsLowSymbolOrLookalike, 0xffe4, 0xffe8)      // |
	properties.Set(charIsReserved, 0xff00)
	properties.Clear(
		0xff0d, // -
		0xff0e, // .
		0xff1a, // :
		0xff3f, // _
	)

	// Ancient symbols https://unicode.org/charts/PDF/U10190.pdf
	properties.Add(charIsLowSymbolOrLookalike, 0x10190) // =

	// Ideopgraphic punctuation https://unicode.org/charts/PDF/U16FE0.pdf
	properties.Add(charIsControl, 0x16fe4) // invisible

	// Musical notation https://unicode.org/charts/PDF/U1D100.pdf
	properties.Add(charIsLowSymbolOrLookalike, 0x1d1c1, 0x1d1c2) // |

	// Mathematical alphanumeric symbols https://unicode.org/charts/PDF/U1D400.pdf
	properties.AddRange(charIsNumeralOrLookalike, 0x1d7ce, 0x1d7ff) // 0-9

	// Escape sequences
	properties.Add(CharIsCustomTextEscapeChar, '"', '\\')
	properties.Add(CharIsQuotedEscapeChar, '\n', '\r', 't', 'n', 'r', '"', '*', '/', '\\', 'u')
	properties.Add(CharIsMarkupEscapeChar, '*', '/', '<', '>', '`', '_', 'u')

	outPath := getExeRelativePath("../../internal/common/unicode-generated.go")
	writer, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening [%v] for writing: %v\n", outPath, err)
		fmt.Fprintf(os.Stderr, "This program expects to reside in the source "+
			"repository at github.com/kstenerud/go-concise-encoding/build-tools/unicode-gen\n")
		os.Exit(1)
	}

	exportProperties(properties, writer)
}

func getExeRelativePath(path string) string {
	return filepath.Join(getExePath(), path)
}

func getExePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func printUsage(writer io.Writer) {
	fmt.Fprintf(writer, "Usage: %v /path/to/ucd.all.flat.xml\n", os.Args[0])
	fmt.Fprintln(writer, "Requires ucd.all.flat.xml from https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip")
}

func exportProperties(properties CharProperties, writer io.Writer) {
	exportHeader(writer)
	exportCharProperties(properties, writer)
	fmt.Fprintf(writer, "\n")
	exportAsciiProperties(properties, writer)
}

func exportHeader(writer io.Writer) {
	fmt.Fprintln(writer, `package common

// Generated by github.com/kstenerud/go-concise-encoding/build-tools/unicode-gen/unicode-gen
  // DO NOT EDIT
  // IF THIS LINE SHOWS IN THE GIT DIFF, YOU HAVE EDITED THIS FILE

type CharProperty uint8

const (
	CharIsUnquotedUnsafe CharProperty = 1 << iota
	CharIsUnquotedFirstCharUnsafe
	CharIsQuotedUnsafe
	CharIsMarkupUnsafe
	CharIsCustomTextUnsafe
	CharIsQuotedEscapeChar
	CharIsMarkupEscapeChar
	CharIsCustomTextEscapeChar
	CharNoProperties CharProperty = 0
)
`)
}

func charValue(char rune, properties CharProperty) string {
	switch char {
	case '\r':
		return `'\r':`
	case '\n':
		return `'\n':`
	case '\t':
		return `'\t':`
	case '\'':
		return `'\'':`
	case '\\':
		return `'\\':`
	}
	if char >= 0x20 && char < 0x7f {
		return fmt.Sprintf(`'%c':`, char)
	}
	if properties&CharIsPrintable != 0 {
		return fmt.Sprintf(`'%c':`, char)
	}
	return fmt.Sprintf("0x%04x:", char)
}

func exportCharProperties(properties CharProperties, writer io.Writer) {
	runes := make([]int, 0, len(properties))
	for k := range properties {
		runes = append(runes, int(k))
	}
	sort.Ints(runes)

	fmt.Fprintf(writer, "var charProperties = map[rune]CharProperty{\n")
	for _, k := range runes {
		ch := rune(k)
		props := properties[ch]
		value := fmt.Sprintf("%02x", ch)
		fmt.Fprintf(writer, "\t/* %5s */ %-8s %v,\n", value, charValue(ch, props), props)
	}
	fmt.Fprintf(writer, "}\n")
}

func exportAsciiProperties(properties CharProperties, writer io.Writer) {
	runes := make([]int, 0, len(properties))
	for k := range properties {
		runes = append(runes, int(k))
	}
	sort.Ints(runes)

	fmt.Fprintf(writer, "var asciiProperties = [256]CharProperty{\n")
	for _, k := range runes {
		if k < 128 {
			ch := rune(k)
			props := properties[ch]
			fmt.Fprintf(writer, "\t/* %02x */ %-7s %v,\n", ch, charValue(ch, props), props)
		}
	}
	fmt.Fprintf(writer, "}\n")
}

// ----------------------------------------------------------------------------

type CharProperty uint16

const (
	CharIsUnquotedUnsafe CharProperty = 1 << iota
	CharIsUnquotedFirstCharUnsafe
	CharIsQuotedUnsafe
	CharIsMarkupUnsafe
	CharIsCustomTextUnsafe
	CharIsQuotedEscapeChar
	CharIsMarkupEscapeChar
	CharIsCustomTextEscapeChar
	CharIsPrintable
	CharNoProperties CharProperty = 0
)

const GeneratorOnlyBits = CharIsPrintable

var charPropertyNames = []string{
	"CharIsUnquotedUnsafe",
	"CharIsUnquotedFirstCharUnsafe",
	"CharIsQuotedUnsafe",
	"CharIsMarkupUnsafe",
	"CharIsCustomTextUnsafe",
	"CharIsQuotedEscapeChar",
	"CharIsMarkupEscapeChar",
	"CharIsCustomTextEscapeChar",
}

func (_this CharProperty) String() string {
	if _this&^GeneratorOnlyBits == CharNoProperties {
		return "CharNoProperties"
	}

	isFirst := true
	builder := strings.Builder{}
	for i := 0; i < 8; i++ {
		if _this&CharProperty(1<<i) != 0 {
			if isFirst {
				isFirst = false
			} else {
				builder.WriteString(" | ")
			}
			builder.WriteString(charPropertyNames[i])
		}
	}
	return builder.String()
}

type CharProperties map[rune]CharProperty

func (_this CharProperties) Clear(chars ...rune) {
	for _, char := range chars {
		_this[char] &= GeneratorOnlyBits
	}
}

func (_this CharProperties) Add(properties CharProperty, chars ...rune) {
	for _, char := range chars {
		_this[char] |= properties
	}
}

func (_this CharProperties) Set(properties CharProperty, chars ...rune) {
	for _, char := range chars {
		_this[char] = properties
	}
}

func (_this CharProperties) Remove(properties CharProperty, chars ...rune) {
	for _, char := range chars {
		_this[char] &= ^properties
	}
}

func (_this CharProperties) AddRange(properties CharProperty, start, end rune) {
	for i := start; i <= end; i++ {
		_this.Add(properties, i)
	}
}

// ----------------------------------------------------------------------------

func loadDB(path string) (chars CharSet, err error) {
	document, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	var dbWrapper DBWrapper
	if err = xml.Unmarshal(document, &dbWrapper); err != nil {
		return
	}

	chars = make(CharSet, 0, len(dbWrapper.DB.Chars))
	for _, char := range dbWrapper.DB.Chars {
		if char.Validate() {
			chars = append(chars, char)
		}
	}
	return
}

type CharSet []*Char

func (_this CharSet) GetRunesWithCriteria(criteria func(*Char) bool) (runes []rune) {
	for _, char := range _this {
		if criteria(char) {
			runes = append(runes, rune(char.Codepoint))
		}
	}
	return
}

type DBWrapper struct {
	XMLName xml.Name   `xml:"ucd"`
	DB      *UnicodeDB `xml:"repertoire"`
}

type UnicodeDB struct {
	XMLName xml.Name `xml:"repertoire"`
	Chars   []*Char  `xml:"char"`
}

func (_this *UnicodeDB) PerformAction(criteria func(*Char) bool, action func(*Char)) {
	for _, char := range _this.Chars {
		if criteria(char) {
			action(char)
		}
	}
}

type Char struct {
	XMLName      xml.Name `xml:"char"`
	CodepointStr string   `xml:"cp,attr"`
	Category     string   `xml:"gc,attr"`
	BidiCategory string   `xml:"bc,attr"`
	Codepoint    int
}

func (_this *Char) Validate() bool {
	codepoint, err := strconv.ParseInt(_this.CodepointStr, 16, 32)
	if err != nil {
		return false
	}
	_this.Codepoint = int(codepoint)
	return true
}
