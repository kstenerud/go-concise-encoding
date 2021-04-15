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

// This package is used to populate `charProperties` and `asciiProperties` in
// github.com/kstenerud/go-concise-encoding/internal/chars/generated-do-not-edit.go
package chars

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

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

const generatedCodePath = "internal/chars/generated-do-not-edit.go"
const packageName = "chars"

type CharProperty uint64

const (
	CharIsWhitespace CharProperty = 1 << iota
	CharIsDigitBase2
	CharIsDigitBase8
	CharIsDigitBase10
	CharIsLowerAF
	CharIsUpperAF
	CharIsAZ
	CharIsAreaLocation
	CharIsObjectEnd
	CharNeedsEscapeQuoted
	CharNeedsEscapeArray
	CharNeedsEscapeMarkup
	CharIsCommentUnsafe
	CharPropertyEnd
	CharIsTextUnsafe
	NoProperties CharProperty = 0

// TODO: Valid timezone strings can contain alpha . - _ /
)

var charPropertyNames = map[CharProperty]string{
	NoProperties:          "NoProperties",
	CharIsWhitespace:      "CharIsWhitespace",
	CharIsDigitBase2:      "CharIsDigitBase2",
	CharIsDigitBase8:      "CharIsDigitBase8",
	CharIsDigitBase10:     "CharIsDigitBase10",
	CharIsLowerAF:         "CharIsLowerAF",
	CharIsUpperAF:         "CharIsUpperAF",
	CharIsAZ:              "CharIsAZ",
	CharIsAreaLocation:    "CharIsAreaLocation",
	CharIsObjectEnd:       "CharIsObjectEnd",
	CharNeedsEscapeQuoted: "CharNeedsEscapeQuoted",
	CharNeedsEscapeArray:  "CharNeedsEscapeArray",
	CharNeedsEscapeMarkup: "CharNeedsEscapeMarkup",
	CharIsCommentUnsafe:   "CharIsCommentUnsafe",
	CharIsTextUnsafe:      "CharIsTextUnsafe",
}

func GenerateCode(projectDir string, xmlPath string) {
	chars, reserveds, err := loadUnicodeDB(xmlPath)
	fatalIfError(err, "Error reading [%v]: %v\n", xmlPath, err)

	properties := extractCharProperties(chars, reserveds)

	outPath := filepath.Join(projectDir, generatedCodePath)
	writer, err := os.Create(outPath)
	fatalIfError(err, "Error opening [%v] for writing: %v", outPath, err)
	defer writer.Close()

	err = exportHeader(writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)

	err = exportCharProperties(properties, writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)
	_, err = fmt.Fprintf(writer, "\n")
	fatalIfError(err, "Error writing to %v: %v", outPath, err)
	err = exportAsciiProperties(properties, writer)
	fatalIfError(err, "Error writing to %v: %v", outPath, err)

	generateRuneByteCounts(writer)
}

func getUTF8ByteCount(firstByte byte) int {
	if firstByte&0x80 == 0 {
		return 1
	}
	if firstByte&0xe0 == 0xc0 {
		return 2
	}
	if firstByte&0xf0 == 0xe0 {
		return 3
	}
	if firstByte&0xf8 == 0xf0 {
		return 4
	}
	return 0
}

func generateRuneByteCounts(writer io.Writer) {
	if _, err := writer.Write([]byte("var runeByteCounts = [32]byte{\n")); err != nil {
		panic(err)
	}

	for i := 0; i < 32; i++ {
		byteCount := getUTF8ByteCount(byte(i << 3))
		if _, err := writer.Write([]byte(fmt.Sprintf("\t0x%02x: %d,\n", i, byteCount))); err != nil {
			panic(err)
		}
	}

	if _, err := writer.Write([]byte("}\n\n")); err != nil {
		panic(err)
	}
}

func charRange(low, high rune) (result []rune) {
	for i := low; i <= high; i++ {
		result = append(result, i)
	}
	return
}

func extractCharProperties(chars CharSet, reserveds ReservedSet) CharProperties {
	properties := CharProperties{}

	properties.Add(CharIsTextUnsafe|
		CharNeedsEscapeArray|
		CharNeedsEscapeMarkup|
		CharNeedsEscapeQuoted,
		chars.GetRunesWithCriteria(func(char *Char) bool {
			switch char.Codepoint {
			case '\r', '\n', '\t', ' ': // Allowed whitespace
				return false
			case 0x2028: // Line separator
				return true
			case 0x2029: // Paragraph separator
				return true
			case 0x1680: // Ogham space mark
				return true
			case 0xfeff: // BOM
				return true
			}
			if char.Codepoint < 0x20 {
				return true
			}

			// Whitespace
			if char.Category == "Zs" || char.Category == "Zl" || char.Category == "Zp" {
				return true
			}

			// Control
			return char.Category == "Cc" || char.Category == "Cf"
		})...)

	properties.Add(CharIsCommentUnsafe,
		chars.GetRunesWithCriteria(func(char *Char) bool {
			switch char.Codepoint {
			case '\r', '\n', '\t', ' ': // Allowed whitespace
				return false
			case 0x2028: // Line separator
				return true
			case 0x2029: // Paragraph separator
				return true
			case 0x1680: // Ogham space mark
				return true
			case 0xfeff: // BOM
				return true
			}
			if char.Codepoint < 0x20 {
				return true
			}

			// Control
			return char.Category == "Cc" || char.Category == "Cf"
		})...)

	properties.Add(CharIsWhitespace, '\t', '\r', '\n', ' ')

	properties.Add(CharIsDigitBase2, charRange('0', '1')...)
	properties.Add(CharIsDigitBase8, charRange('0', '7')...)
	properties.Add(CharIsDigitBase10, charRange('0', '9')...)
	properties.Add(CharIsLowerAF, charRange('a', 'f')...)
	properties.Add(CharIsUpperAF, charRange('A', 'F')...)
	properties.Add(CharIsAZ|CharIsAreaLocation, charRange('a', 'z')...)
	properties.Add(CharIsAZ|CharIsAreaLocation, charRange('A', 'Z')...)
	properties.Add(CharIsAreaLocation, '_', '-', '+', '/')

	properties.Add(CharIsObjectEnd, '\r', '\n', '\t', ' ', ']', '}', ')', '>', ',', '=', ':', '|', '/')

	properties.AddLL(CharNeedsEscapeQuoted, '\\', '"')
	properties.AddLL(CharNeedsEscapeArray, '\\', '|', '\t', '\r', '\n')
	properties.AddLL(CharNeedsEscapeMarkup, '\\', '<', '>')

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

	return properties
}

func exportHeader(writer io.Writer) (err error) {
	var propType string
	switch {
	case CharPropertyEnd&0x1ff != 0:
		propType = "uint8"
	case CharPropertyEnd&0x1fe00 != 0:
		propType = "uint16"
	case CharPropertyEnd&0x1fffe0000 != 0:
		propType = "uint32"
	default:
		propType = "uint64"
	}

	if _, err = fmt.Fprintf(writer, standard.Header+
		"package %v\n\ntype CharProperty %v\n\nconst (", packageName, propType); err != nil {
		return
	}

	for i := CharProperty(1); i < CharPropertyEnd; i <<= 1 {
		if _, err = fmt.Fprintf(writer, "\n\t%v", i); err != nil {
			return
		}
		if i == 1 {
			if _, err = fmt.Fprintf(writer, " CharProperty = 1 << iota"); err != nil {
				return
			}
		}
	}

	_, err = fmt.Fprintf(writer, "\n\t%v = 0\n)\n\n", NoProperties)
	return
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

	if _, err := fmt.Fprintf(writer, "var asciiProperties = [0x101]CharProperty{\n"); err != nil {
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

	_, err := fmt.Fprintf(writer, "\t/* EOF */ 0x100: CharIsObjectEnd|CharNeedsQuote,\n}\n")
	return err
}

func charValue(char rune, properties CharProperty) string {
	switch char {
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	}
	if char < 0x20 || char == 0x7f {
		return ""
	}
	if char >= 0x20 && char < 0x7f {
		return fmt.Sprintf("[%c]", char)
	}
	if properties&CharIsTextUnsafe == 0 {
		return fmt.Sprintf("[%c]", char)
	}
	return ""
}

func fatalIfError(err error, format string, args ...interface{}) {
	if err != nil {
		panic(fmt.Errorf(format, args))
	}
}

// ----------------------------------------------------------------------------

func (_this CharProperty) String() string {
	if _this == 0 {
		return charPropertyNames[_this]
	}

	isFirst := true
	builder := strings.Builder{}
	for i := CharProperty(1); i < CharPropertyEnd; i <<= 1 {
		if _this&i != 0 {
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

func (_this CharProperties) Add(properties CharProperty, chars ...rune) {
	for _, char := range chars {
		_this[char] |= properties
	}
}

func (_this CharProperties) AddLL(properties CharProperty, chars ...rune) {
	for _, char := range chars {
		_this[char] |= properties
		for _, ll := range lookalikes[char] {
			_this[ll] |= properties
		}
	}
}

func (_this CharProperties) Unmark(chars ...rune) {
	for _, char := range chars {
		delete(_this, char)
	}
}

func (_this CharProperties) UnmarkRange(start, end rune) {
	for i := start; i <= end; i++ {
		_this.Unmark(i)
	}
}

// ----------------------------------------------------------------------------

var lookalikes = map[rune][]rune{
	'!': []rune{0x01c3, 0x203c, 0x2048, 0x2049, 0x2d51, 0xfe15, 0xfe57, 0xff01},
	'"': []rune{0x02ba, 0x02ee, 0x201c, 0x201d, 0x201f, 0x2033, 0x2034, 0x2036, 0x2037, 0x2057, 0x3003, 0xff02},
	'#': []rune{0xfe5f, 0xff03},
	'$': []rune{0xfe69, 0xff04},
	'%': []rune{0x2052, 0xfe6a, 0xff05},
	'&': []rune{0xfe60, 0xff06},
	'\'': []rune{0x00b4, 0x02b9, 0x02bb, 0x02bc, 0x02bd, 0x02ca, 0x02c8, 0x0374, 0x2018, 0x2019,
		0x201a, 0x201b, 0x2032, 0x2035, 0xa78b, 0xa78c, 0xfe10, 0xfe50, 0xff07, 0x10107, 0x1d112},
	'(': []rune{0x2474, 0x2475, 0x2476, 0x2477, 0x2478, 0x2479, 0x247a, 0x247b, 0x247c, 0x247d,
		0x247e, 0x247f, 0x2480, 0x2481, 0x2482, 0x2483, 0x2484, 0x2485, 0x2486, 0x2487, 0x249c,
		0x249d, 0x249e, 0x249f, 0x24a0, 0x24a1, 0x24a2, 0x24a3, 0x24a4, 0x24a5, 0x24a6, 0x24a7,
		0x24a8, 0x24a9, 0x24aa, 0x24ab, 0x24ac, 0x24ad, 0x24ae, 0x24af, 0x24b0, 0x24b1, 0x24b2,
		0x24b3, 0x24b4, 0x24b5, 0xfe59, 0xff08},
	')': []rune{0x2474, 0x2475, 0x2476, 0x2477, 0x2478, 0x2479, 0x247a, 0x247b, 0x247c, 0x247d,
		0x247e, 0x247f, 0x2480, 0x2481, 0x2482, 0x2483, 0x2484, 0x2485, 0x2486, 0x2487, 0x249c,
		0x249d, 0x249e, 0x249f, 0x24a0, 0x24a1, 0x24a2, 0x24a3, 0x24a4, 0x24a5, 0x24a6, 0x24a7,
		0x24a8, 0x24a9, 0x24aa, 0x24ab, 0x24ac, 0x24ad, 0x24ae, 0x24af, 0x24b0, 0x24b1, 0x24b2,
		0x24b3, 0x24b4, 0x24b5, 0xfe5a, 0xff09},
	'*': []rune{0x204e, 0x2055, 0x2217, 0x22c6, 0x2b51, 0xfe61, 0xff0a},
	'+': []rune{0xfe62, 0xff0b},
	',': []rune{0x02cc, 0x02cf, 0x0375, 0xff0c, 0x10100},
	'-': []rune{0x02c9, 0x2010, 0x2011, 0x2012, 0x2013, 0x2014, 0x2015, 0x2212, 0x23af, 0x23bb,
		0x23bc, 0x23e4, 0x23fd, 0xfe58, 0xfe63, 0xff0d, 0xff70, 0x10110, 0x10191, 0x1d116},
	'.':  []rune{0xfe52, 0xff0e},
	'/':  []rune{0x2044, 0x2215, 0x27cb, 0x29f8, 0x3033, 0xff0f, 0x1d10d},
	':':  []rune{0x02f8, 0x205a, 0x2236, 0xa789, 0xfe13, 0xfe30, 0xfe55, 0xff1a, 0x1d108},
	';':  []rune{0x037e, 0xfe14, 0xfe54, 0xff1b},
	'<':  []rune{0x00ab, 0x02c2, 0x3111, 0x2039, 0x227a, 0x2329, 0x2d66, 0x3008, 0xfe64, 0xff1c, 0x1032d},
	'=':  []rune{0xa78a, 0xfe66, 0xff1d, 0x10190, 0x16fe3},
	'>':  []rune{0x00bb, 0x02c3, 0x203a, 0x227b, 0x232a, 0x3009, 0xfe65, 0xff1e},
	'?':  []rune{0x2047, 0x2048, 0x2049, 0xfe16, 0xfe56, 0xff1f},
	'@':  []rune{0xfe6b, 0xff20},
	'[':  []rune{0xfe5d, 0xff3b, 0x1d115},
	'\\': []rune{0x2216, 0x27cd, 0x29f5, 0x29f9, 0x3035, 0xfe68, 0xff3c},
	']':  []rune{0xfe5e, 0xff3d},
	'^':  []rune{0xff3e},
	'_':  []rune{0x02cd, 0x23bd, 0xff3f},
	'`':  []rune{0x02cb, 0xfe11, 0xfe45, 0xfe46, 0xfe51, 0xff40},
	'{':  []rune{0xfe5b, 0xff5b, 0x1d114},
	'|': []rune{0x00a6, 0x01c0, 0x2223, 0x2225, 0x239c, 0x239f, 0x23a2, 0x23a5, 0x23aa, 0x23ae,
		0x23b8, 0x23b9, 0x23d0, 0x2d4f, 0x3021, 0xfe31, 0xfe33, 0xff5c, 0xffdc, 0xffe4, 0xffe8,
		0x1028a, 0x10320, 0x10926, 0x10ce5, 0x10cfa, 0x1d100, 0x1d105, 0x1d1c1, 0x1d1c2},
	'}': []rune{0xfe5c, 0xff5d},
	'~': []rune{0x2053, 0x223c, 0x223f, 0x301c, 0xff5e},
	'0': []rune{0xff10, 0x1d7ce, 0x1d7d8, 0x1d7e2, 0x1d7ec, 0x1d7f6, 0x1f100, 0x1f101},
	'1': []rune{0x00b9, 0x2488, 0x2491, 0x2492, 0x2493, 0x2494, 0x2495, 0x2496, 0x2497, 0x2498, 0x2499, 0x249a, 0x249b, 0xff11, 0x1d7cf, 0x1d7d9, 0x1d7e3, 0x1d7ed, 0x1d7f7, 0x1f102},
	'2': []rune{0x00b2, 0x2489, 0xff12, 0x1d7d0, 0x1d7da, 0x1d7e4, 0x1d7ee, 0x1d7f8, 0x1f103},
	'3': []rune{0x00b3, 0x248a, 0xff13, 0x1d7d1, 0x1d7db, 0x1d7e5, 0x1d7ef, 0x1d7f9, 0x1f104},
	'4': []rune{0x248b, 0xff14, 0x1d7d2, 0x1d7dc, 0x1d7e6, 0x1d7f0, 0x1d7fa, 0x1f105},
	'5': []rune{0x248c, 0xff15, 0x1d7d3, 0x1d7dd, 0x1d7e7, 0x1d7f1, 0x1d7fb, 0x1f106},
	'6': []rune{0x248d, 0xff16, 0x1d7d4, 0x1d7de, 0x1d7e8, 0x1d7f2, 0x1d7fc, 0x1f107},
	'7': []rune{0x248f, 0xff17, 0x1d7d5, 0x1d7df, 0x1d7e9, 0x1d7f3, 0x1d7fd, 0x1f108},
	'8': []rune{0x248f, 0xff18, 0x10931, 0x1d7d6, 0x1d7e0, 0x1d7ea, 0x1d7f4, 0x1d7fe, 0x1f109},
	'9': []rune{0x2490, 0xff19, 0x1d7d7, 0x1d7e1, 0x1d7eb, 0x1d7f5, 0x1d7ff, 0x1f10a},
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
