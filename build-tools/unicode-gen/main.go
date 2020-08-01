// This program is used to populate `charProperties` and `asciiProperties` in
// github.com/kstenerud/go-concise-encoding/internal/common/unicode-generated.go
//
// Notes:
// - Private use chars are not included
// - Reserved chars are not included
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

func main() {
	if len(os.Args) != 2 {
		printUsage(os.Stderr)
		os.Exit(1)
	}

	xmlPath := os.Args[1]
	chars, reserveds, err := loadUnicodeDB(xmlPath)
	fatalIfError(err, "Error reading [%v]: %v\n", xmlPath, err)

	properties := extractCharProperties(chars, reserveds)

	outPath := getExeRelativePath("../../internal/unicode/unicode-generated.go")
	writer, err := os.Create(outPath)
	fatalIfError(err, "Error opening [%v] for writing: %v\n"+
		"This program expects to reside in the source repository at "+
		"github.com/kstenerud/go-concise-encoding/build-tools/unicode-gen\n", outPath, err)
	defer writer.Close()

	err = exportHeader(writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)

	err = exportCharProperties(properties, writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)
	_, err = fmt.Fprintf(writer, "\n")
	fatalIfError(err, "Error writing to %v: %v", outPath, err)
	err = exportAsciiProperties(properties, writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)
}

func extractCharProperties(chars CharSet, reserveds ReservedSet) CharProperties {
	properties := CharProperties{}

	properties.Add(CharIsQuotedOrCustomTextDelimiter, '"', '\\')
	properties.Add(CharIsMarkupDelimiter, '<', '>', '\\', '`')
	properties.AddRange(CharIsNumeralOrLookalike, '0', '9')

	properties.Add(CharIsTabReturnNewline, '\r', '\n', '\t')

	properties.Add(CharIsWhitespace, chars.GetRunesWithCriteria(func(char *Char) bool {
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

	properties.Add(CharIsControl, chars.GetRunesWithCriteria(func(char *Char) bool {
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

	properties.Add(CharIsLowSymbolOrLookalike,
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', ';', '<',
		'=', '>', '?', '@', '[', '\\', ']', '^', '`', '{', '|', '}', '~')

	properties.Clear('-', '.', ':', '_') // Unquoted safe

	// Latin punctuation https://unicode.org/charts/PDF/U0080.pdf
	properties.Add(CharIsLowSymbolOrLookalike,
		0x00a6, // |
		0x00b4, // '
	)

	// General punctuation https://unicode.org/charts/PDF/U2000.pdf
	properties.AddRange(CharIsLowSymbolOrLookalike, 0x2018, 0x201f) // "
	properties.AddRange(CharIsLowSymbolOrLookalike, 0x2032, 0x2037) // "
	properties.Add(CharIsLowSymbolOrLookalike,
		0x2039, // <
		0x203a, // >
		0x203c, // !
	)
	properties.AddRange(CharIsLowSymbolOrLookalike, 0x2047, 0x2049) // !?
	properties.Add(CharIsLowSymbolOrLookalike,
		0x204e, 0x2055, // *
		0x2052, // %
		0x2057, // "
	)

	// Mathematical operators https://unicode.org/charts/PDF/U2200.pdf
	properties.Add(CharIsLowSymbolOrLookalike,
		0x2217, // *
		0x2223, // |
		0x223c, // ~
	)

	// Miscellaneous technical https://unicode.org/charts/PDF/U2300.pdf
	properties.Add(CharIsLowSymbolOrLookalike,
		0x239c, 0x239f, 0x23a2, 0x23a5, 0x23aa, 0x23b8, 0x23b9, 0x23d0, 0x23fd, // |
	)

	// CJK vertical forms https://unicode.org/charts/PDF/UFE10.pdf
	properties.Add(CharIsLowSymbolOrLookalike, 0xfe10, 0xfe11)      // '
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xfe14, 0xfe16) // ;!?

	// CJK compatibility https://unicode.org/charts/PDF/UFE30.pdf
	properties.Add(CharIsLowSymbolOrLookalike, 0xfe31, 0xfe33) // |
	properties.Add(CharIsLowSymbolOrLookalike, 0xfe45, 0xfe46) // '

	// CJK small form variants https://unicode.org/charts/PDF/UFE50.pdf
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xfe50, 0xfe6b) // symbols
	properties.Clear(
		0xfe52, // .
		0xfe55, // :
		0xfe58, // -
		0xfe63, // -
	)

	// CJK halfwidth, fullwidth https://unicode.org/charts/PDF/UFF00.pdf
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xff00, 0xff0f) // symbols
	properties.AddRange(CharIsNumeralOrLookalike, 0xff10, 0xff19)   // 0-9
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xff1a, 0xff20) // symbols
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xff3b, 0xff40) // symbols
	properties.AddRange(CharIsLowSymbolOrLookalike, 0xff5b, 0xff5e) // symbols
	properties.Add(CharIsLowSymbolOrLookalike, 0xffe4, 0xffe8)      // |
	properties.Clear(
		0xff0d, // -
		0xff0e, // .
		0xff1a, // :
		0xff3f, // _
	)

	// Ancient symbols https://unicode.org/charts/PDF/U10190.pdf
	properties.Add(CharIsLowSymbolOrLookalike, 0x10190) // =

	// Ideographic punctuation https://unicode.org/charts/PDF/U16FE0.pdf
	properties.Add(CharIsControl, 0x16fe4) // invisible

	// Musical notation https://unicode.org/charts/PDF/U1D100.pdf
	properties.Add(CharIsLowSymbolOrLookalike, 0x1d1c1, 0x1d1c2) // |

	// Mathematical alphanumeric symbols https://unicode.org/charts/PDF/U1D400.pdf
	properties.AddRange(CharIsNumeralOrLookalike, 0x1d7ce, 0x1d7ff) // 0-9

	// Don't include reserved chars
	for _, r := range reserveds {
		if r.CPStr != "" {
			properties.Unmark(rune(r.CP))
		} else {
			properties.UnmarkRange(rune(r.FirstCP), rune(r.LastCP))
		}
	}

	// Don't include private chars
	properties.UnmarkRange(0xe000, 0xf8ff)
	properties.UnmarkRange(0xf0000, 0xffffd)
	properties.UnmarkRange(0x100000, 0x10fffd)

	// Allowed in marker and reference IDs
	properties.AddRange(CharIsMarkerIDSafe, '0', '9')
	properties.AddRange(CharIsMarkerIDSafe, 'a', 'z')
	properties.AddRange(CharIsMarkerIDSafe, 'A', 'Z')
	properties.Add(CharIsMarkerIDSafe, '_')

	return properties
}

func fatalIfError(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, args...)
		os.Exit(1)
	}
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

func exportHeader(writer io.Writer) error {
	_, err := fmt.Fprintln(writer, `package unicode

// Generated by github.com/kstenerud/go-concise-encoding/build-tools/unicode-gen/unicode-gen
  // DO NOT EDIT
  // IF THIS LINE SHOWS UP IN THE GIT DIFF, THIS FILE HAS BEEN EDITED

type CharProperty uint8

const (
	NumeralOrLookalike CharProperty = 1 << iota
	LowSymbolOrLookalike
	Whitespace
	Control
	TabReturnNewline
	QuotedOrCustomTextDelimiter
	MarkupDelimiter
	MarkerIDSafe
	NoProperties CharProperty = 0
)
`)
	return err
}

func charValue(char rune, properties CharProperty) string {
	const printableProperties = ^(CharIsControl)

	switch char {
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	}
	if char >= 0x20 && char < 0x7f {
		return fmt.Sprintf("[%c]", char)
	}
	if properties&printableProperties != 0 || properties == 0 {
		return fmt.Sprintf("[%c]", char)
	}
	return ""
}

func exportCharProperties(properties CharProperties, writer io.Writer) error {
	runes := make([]int, 0, len(properties))
	for k := range properties {
		runes = append(runes, int(k))
	}
	sort.Ints(runes)

	if _, err := fmt.Fprintf(writer, "var charProperties = map[rune]CharProperty{\n"); err != nil {
		return err
	}
	for _, k := range runes {
		ch := rune(k)
		props := properties[ch]
		if _, err := fmt.Fprintf(writer, "\t/* %-3s */ 0x%02x: %v,\n", charValue(ch, props), ch, props); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(writer, "}\n")
	return err
}

func exportAsciiProperties(properties CharProperties, writer io.Writer) error {
	runes := make([]int, 0, len(properties))
	for k := range properties {
		runes = append(runes, int(k))
	}
	sort.Ints(runes)

	if _, err := fmt.Fprintf(writer, "var asciiProperties = [256]CharProperty{\n"); err != nil {
		return err
	}
	for _, k := range runes {
		if k < 128 {
			ch := rune(k)
			props := properties[ch]
			if _, err := fmt.Fprintf(writer, "\t/* %-3s */ 0x%02x: %v,\n", charValue(ch, props), ch, props); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprintf(writer, "}\n")
	return err
}

// ----------------------------------------------------------------------------

type CharProperty uint16

const (
	CharIsNumeralOrLookalike CharProperty = 1 << iota
	CharIsLowSymbolOrLookalike
	CharIsWhitespace
	CharIsControl
	CharIsTabReturnNewline
	CharIsQuotedOrCustomTextDelimiter
	CharIsMarkupDelimiter
	CharIsMarkerIDSafe
	NoProperties CharProperty = 0
)

var charPropertyNames = []string{
	"NumeralOrLookalike",
	"LowSymbolOrLookalike",
	"Whitespace",
	"Control",
	"TabReturnNewline",
	"QuotedOrCustomTextDelimiter",
	"MarkupDelimiter",
	"MarkerIDSafe",
}

func (_this CharProperty) String() string {
	if _this == NoProperties {
		return "NoProperties"
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
		_this[char] = 0
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

func (_this CharProperties) Unmark(chars ...rune) {
	for _, char := range chars {
		// There are almost a million reserved chars, so just delete instead
		// _this[char] = CharIsReserved
		delete(_this, char)
	}
}

func (_this CharProperties) UnmarkRange(start, end rune) {
	for i := start; i <= end; i++ {
		_this.Unmark(i)
	}
}

// ----------------------------------------------------------------------------

func loadUnicodeDB(path string) (chars CharSet, reserved ReservedSet, err error) {
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

	reserved = make(ReservedSet, 0, len(dbWrapper.DB.Chars))
	for _, res := range dbWrapper.DB.Reserveds {
		if res.Validate() {
			reserved = append(reserved, res)
		}
	}
	return
}

type CharSet []*Char
type ReservedSet []*Reserved

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
	XMLName   xml.Name    `xml:"repertoire"`
	Chars     []*Char     `xml:"char"`
	Reserveds []*Reserved `xml:"reserved"`
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
	Codepoint    rune
}

func (_this *Char) Validate() bool {
	codepoint, err := strconv.ParseInt(_this.CodepointStr, 16, 32)
	if err != nil {
		return false
	}
	_this.Codepoint = rune(codepoint)
	return true
}

type Reserved struct {
	CPStr      string `xml:"cp,attr"`
	FirstCPStr string `xml:"first-cp,attr"`
	LastCPStr  string `xml:"last-cp,attr"`
	CP         rune
	FirstCP    rune
	LastCP     rune
}

func (_this *Reserved) Validate() bool {
	if _this.CPStr != "" {
		cp, err := strconv.ParseInt(_this.CPStr, 16, 32)
		if err != nil {
			return false
		}
		_this.CP = rune(cp)
		return true
	}

	firstCP, err := strconv.ParseInt(_this.FirstCPStr, 16, 32)
	if err != nil {
		return false
	}
	lastCP, err := strconv.ParseInt(_this.LastCPStr, 16, 32)
	if err != nil {
		return false
	}
	_this.FirstCP = rune(firstCP)
	_this.LastCP = rune(lastCP)
	return true
}
