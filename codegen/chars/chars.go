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

package chars

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

const path = "internal/chars"

var imports = []string{
	"unicode/utf8",
}

func GenerateCode(projectDir string, xmlPath string) {
	chars, err := loadUnicodeDB(xmlPath)
	standard.PanicIfError(err, "Error reading [%v]", xmlPath)
	classifyRunes(chars)

	generatedFilePath := standard.GetGeneratedCodePath(projectDir, path)
	writer, err := os.Create(generatedFilePath)
	standard.PanicIfError(err, "Error opening [%v] for writing", generatedFilePath)
	defer writer.Close()
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Errorf("Error while generating %v: %v", generatedFilePath, e))
		}
	}()

	standard.WriteHeader(writer, path, imports)

	generatePropertiesType(writer)
	generateSpacer(writer)

	generateSafetyFlagsType(writer)
	generateSpacer(writer)

	generateRuneByteCounts(writer)
	generateSpacer(writer)

	generateStringlikeUnsafeTable(writer)
	generateSpacer(writer)

	generatePropertiesTable(writer)
	generateSpacer(writer)

	generateIdentifierSafeTable(writer)
	generateSpacer(writer)

	generateStringlikeSafeTable(writer)
}

// -----
// Runes
// -----

func classifyRunes(chars CharSet) {
	const (
		safeForID           = true
		unsafeForID         = false
		safeForStringlike   = SafetyAll
		unsafeForStringlike = SafetyNone
	)

	// Character classes (https://unicodebook.readthedocs.io/unicode.html):

	// Private chars
	setSafety(unsafeForID, unsafeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		return char.Category == "Co"
	}))

	// Control chars
	setSafety(unsafeForID, unsafeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		return char.Category == "Cc"
	}))

	// Whitespace
	setSafety(unsafeForID, safeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		return char.MajorCategory == 'Z'
	}))

	// Letters, numbers, mark
	setSafety(safeForID, safeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		// TODO: Modifiers OK in identifiers?
		return char.MajorCategory == 'L' || char.MajorCategory == 'N' || char.MajorCategory == 'M'
	}))

	// Symbols, Punctuation
	setSafety(unsafeForID, safeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		return char.MajorCategory == 'P' || char.MajorCategory == 'S'
	}))

	// Format chars
	setSafety(unsafeForID, safeForStringlike, chars.RunesWithCriteria(func(char *Char) bool {
		// https://262.ecma-international.org/11.0/#sec-unicode-format-control-characters
		return char.Category == "Cf"
	}))

	// Structural whitespace
	setSafety(unsafeForID, safeForStringlike, charSet('\r', '\n', '\t', ' '))

	// Symbols allowed in identifiers
	setSafety(safeForID, safeForStringlike, charSet('-', '_'))

	// Stringlike safety
	markUnsafeFor(SafetyString, charsAndLookalikes(charSet('\\', '"')...))
	markUnsafeFor(SafetyArray, charsAndLookalikes(charSet('\\', '|', '\t', '\r', '\n')...))
	markUnsafeFor(SafetyMarkup, charsAndLookalikes(charSet('\\', '<', '>')...))

	// Structural char properties
	addProperties(StructWS, charSet('\r', '\n', '\t', ' '))
	addProperties(ObjectEnd, charSet('\r', '\n', '\t', ' ', ']', '}', ')', '>', ',', '=', ':', '|', '/'))
	addProperties(DigitBase2, charRange('0', '1'))
	addProperties(DigitBase8, charRange('0', '7'))
	addProperties(DigitBase10, charRange('0', '9'))
	addProperties(LowerAF, charRange('a', 'f'))
	addProperties(UpperAF, charRange('A', 'F'))
	addProperties(AreaLocation, charRange('a', 'z'), charRange('A', 'Z'), charSet('_', '-', '+', '/'))
	addProperties(UUID, charRange('0', '9'), charRange('a', 'f'), charRange('A', 'F'), charSet('-'))

	// Invalid chars:

	// Surrogates, Reserved
	markInvalid(chars.RunesWithCriteria(func(char *Char) bool {
		return char.Category == "Cs" || char.Category == "Cn"
	}))

	// Mark chars that can be printed in the generated comments
	markGoSafe(chars.RunesWithCriteria(func(char *Char) bool {
		return char.MajorCategory == 'L' || char.MajorCategory == 'N' || char.MajorCategory == 'P' || char.MajorCategory == 'S'
	}))
}

// Code Generators

func generateSpacer(writer io.Writer) {
	if _, err := fmt.Fprintf(writer, "\n"); err != nil {
		panic(err)
	}
}

func generatePropertiesType(writer io.Writer) {
	var propType string
	switch {
	case EndProperties <= 0x100:
		propType = "uint8"
	case EndProperties <= 0x10000:
		propType = "uint16"
	case EndProperties <= 0x100000000:
		propType = "uint32"
	default:
		propType = "uint64"
	}

	if _, err := fmt.Fprintf(writer, "type Properties %v\n\nconst (", propType); err != nil {
		panic(err)
	}

	for i := Properties(1); i < EndProperties; i <<= 1 {
		if _, err := fmt.Fprintf(writer, "\n\t%v", i); err != nil {
			panic(err)
		}
		if i == 1 {
			if _, err := fmt.Fprintf(writer, " Properties = 1 << iota"); err != nil {
				panic(err)
			}
		}
	}

	if _, err := fmt.Fprintf(writer, "\n\t%v = 0\n)\n", NoProperties); err != nil {
		panic(err)
	}
}

func generateSafetyFlagsType(writer io.Writer) {
	if _, err := fmt.Fprintf(writer, "type SafetyFlags uint8\n\nconst ("); err != nil {
		panic(err)
	}

	for i := SafetyFlags(1); i < EndSafetyFlags; i <<= 1 {
		if _, err := fmt.Fprintf(writer, "\n\t%v", i); err != nil {
			panic(err)
		}
		if i == 1 {
			if _, err := fmt.Fprintf(writer, " SafetyFlags = 1 << iota"); err != nil {
				panic(err)
			}
		}
	}

	if _, err := fmt.Fprintf(writer, "\n\tSafetyAll = %v", SafetyAll); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(writer, "\n\t%v = 0\n)\n", SafetyNone); err != nil {
		panic(err)
	}
}

func generatePropertiesTable(writer io.Writer) error {
	if _, err := fmt.Fprintf(writer, "var properties = [0x101]Properties{\n"); err != nil {
		return err
	}
	for r, props := range properties {
		rs := runeString(rune(r))
		if r >= 0x80 {
			rs = ""
		}
		if _, err := fmt.Fprintf(writer, "\t/* %-3s */ 0x%02x: %v,\n", rs, r, props); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(writer, "\t/* EOF */ 0x100: %v,\n}\n", ObjectEnd)
	return err
}

func generateStringlikeUnsafeTable(writer io.Writer) error {
	runes := make([]int, 0, len(stringlikeUnsafe))
	for r, flags := range stringlikeUnsafe {
		if flags != 0 {
			runes = append(runes, int(r))
		}
	}
	sort.Ints(runes)

	if _, err := fmt.Fprintf(writer, "var stringlikeUnsafe = map[rune]SafetyFlags{\n"); err != nil {
		return err
	}
	for _, k := range runes {
		r := rune(k)
		flags := stringlikeUnsafe[r]
		if flags != 0 && flags != SafetyAll {
			rs := runeString(rune(r))
			if _, err := fmt.Fprintf(writer, "\t/* %-3s */ 0x%02x: %v,\n", rs, r, flags); err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintf(writer, "}\n")
	return err
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

	if _, err := writer.Write([]byte("}\n")); err != nil {
		panic(err)
	}
}

func generateIdentifierSafeTable(writer io.Writer) {
	generateByteTable(writer, "identifierSafe", identifierSafe[:])
}

func generateStringlikeSafeTable(writer io.Writer) {
	generateByteTable(writer, "stringlikeSafe", stringlikeSafe[:])
}

func generateByteTable(writer io.Writer, name string, table []byte) {
	if _, err := writer.Write([]byte(fmt.Sprintf("var %s = [(utf8.MaxRune + 1) / 8]byte{", name))); err != nil {
		panic(err)
	}

	newlineAfter := 12
	var str string

	for i, v := range table {
		if i%newlineAfter == 0 {
			str = "\n\t"
		} else {
			str = " "
		}

		str := fmt.Sprintf("%s0x%02x,", str, v)

		if _, err := writer.Write([]byte(str)); err != nil {
			panic(err)
		}
	}

	if _, err := writer.Write([]byte("\n}\n")); err != nil {
		panic(err)
	}
}

// -------
// Utility
// -------

func runeString(r rune) string {
	switch r {
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	default:
		if isRuneGoSafe(r) {
			return fmt.Sprintf("[%c]", r)
		}
		return ""
	}
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

func setSafety(
	isIdentifierSafe bool,
	safeFor SafetyFlags,
	runes ...[]rune) {

	for _, r := range runes {
		for _, rr := range r {
			setRuneSafety(rr, isIdentifierSafe, safeFor)
		}
	}
}

func markUnsafeFor(unsafeFor SafetyFlags, runes ...[]rune) {

	for _, r := range runes {
		for _, rr := range r {
			markRuneUnsafeFor(rr, unsafeFor)
		}
	}
}

func addProperties(props Properties, runes ...[]rune) {

	for _, r := range runes {
		for _, rr := range r {
			addRuneProperties(rr, props)
		}
	}
}

func markInvalid(runes ...[]rune) {

	for _, r := range runes {
		for _, rr := range r {
			markRuneInvalid(rr)
		}
	}
}

func markGoSafe(runes ...[]rune) {

	for _, r := range runes {
		for _, rr := range r {
			markRuneGoSafe(rr)
		}
	}
}

func setBitArrayValue(array []byte, index int, value bool) {
	bit := byte(1 << (index & 7))
	bits := array[index>>3]
	if value {
		bits |= bit
	} else {
		bits &= ^bit
	}
	array[index>>3] = bits
}

func getBitArrayValue(array []byte, index int) bool {
	bits := array[index>>3]
	return bits&(1<<(index&7)) != 0
}

func setRuneSafety(r rune, isIdentifierSafe bool, safeFor SafetyFlags) {
	unsafeFor := ^safeFor & SafetyAll
	isStringlikeSafe := unsafeFor == 0

	setBitArrayValue(identifierSafe[:], int(r), isIdentifierSafe)
	setBitArrayValue(stringlikeSafe[:], int(r), isStringlikeSafe)
	stringlikeUnsafe[r] = unsafeFor
}

func markRuneUnsafeFor(r rune, unsafeFor SafetyFlags) {
	unsafeFor |= stringlikeUnsafe[r]
	isStringlikeSafe := unsafeFor == 0

	setBitArrayValue(stringlikeSafe[:], int(r), isStringlikeSafe)
	stringlikeUnsafe[r] = unsafeFor
}

func markRuneInvalid(r rune) {
	// TODO: really mark invalid in a findable way?
	setRuneSafety(r, false, SafetyAll)
}

func markRuneGoSafe(r rune) {
	setBitArrayValue(goSafe[:], int(r), true)
}

func isRuneGoSafe(r rune) bool {
	return getBitArrayValue(goSafe[:], int(r))
}

func addRuneProperties(r rune, props Properties) {

	if r < 0x100 {
		properties[r] |= props
	}
}

func charRange(low, high rune) (result []rune) {
	for i := low; i <= high; i++ {
		result = append(result, i)
	}
	return
}

func charSet(r ...rune) []rune {
	return r
}

func charsAndLookalikes(r ...rune) []rune {
	result := make([]rune, 0, len(r))
	for _, rr := range r {
		result = append(result, rr)
		for _, lookalike := range lookalikes[rr] {
			result = append(result, lookalike)
		}
	}
	return result
}

// ----------
// Unicode DB
// ----------

func loadUnicodeDB(path string) (chars CharSet, err error) {
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
		chars = append(chars, char.All()...)
	}
	for _, char := range dbWrapper.DB.Reserveds {
		chars = append(chars, char.All()...)
	}

	return
}

type CharSet []*Char

func (_this CharSet) RunesWithCriteria(criteria func(*Char) bool) (runes []rune) {
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
	XMLName   xml.Name `xml:"repertoire"`
	Chars     []*Char  `xml:"char"`
	Reserveds []*Char  `xml:"reserved"`
}

func (_this *UnicodeDB) PerformAction(criteria func(*Char) bool, action func(*Char)) {
	for _, char := range _this.Chars {
		if criteria(char) {
			action(char)
		}
	}
}

type Char struct {
	CodepointStr  string `xml:"cp,attr"`
	FirstCPStr    string `xml:"first-cp,attr"`
	LastCPStr     string `xml:"last-cp,attr"`
	Category      string `xml:"gc,attr"`
	MajorCategory byte
	Codepoint     rune
}

func (_this *Char) All() (result []*Char) {
	_this.MajorCategory = _this.Category[0]

	if _this.CodepointStr != "" {
		codepoint, err := strconv.ParseInt(_this.CodepointStr, 16, 32)
		if err != nil {
			return
		}
		_this.Codepoint = rune(codepoint)
		return []*Char{_this}
	}

	firstCP, err := strconv.ParseInt(_this.FirstCPStr, 16, 32)
	if err != nil {
		return
	}
	lastCP, err := strconv.ParseInt(_this.LastCPStr, 16, 32)
	if err != nil {
		return
	}

	for i := rune(firstCP); i <= rune(lastCP); i++ {
		result = append(result, &Char{
			Category:      _this.Category,
			MajorCategory: _this.MajorCategory,
			Codepoint:     i,
		})
	}
	return
}

// ----
// Data
// ----

type SafetyFlags byte

const (
	SafetyString SafetyFlags = 1 << iota
	SafetyArray
	SafetyMarkup
	SafetyComment

	EndSafetyFlags
	SafetyAll              = EndSafetyFlags - 1
	SafetyNone SafetyFlags = 0
)

func (_this SafetyFlags) String() string {
	if _this == 0 {
		return safetyNames[_this]
	}

	isFirst := true
	builder := strings.Builder{}
	for i := SafetyFlags(1); i < EndSafetyFlags; i <<= 1 {
		if _this&i != 0 {
			if isFirst {
				isFirst = false
			} else {
				builder.WriteString(" | ")
			}
			builder.WriteString(safetyNames[i])
		}
	}
	return builder.String()
}

var safetyNames = map[SafetyFlags]string{
	SafetyNone:     "SafetyNone",
	SafetyString:   "SafetyString",
	SafetyArray:    "SafetyArray",
	SafetyMarkup:   "SafetyMarkup",
	SafetyComment:  "SafetyComment",
	EndSafetyFlags: "EndSafetyFlags",
	SafetyAll:      "SafetyAll",
}

var goSafe [(utf8.MaxRune + 1) / 8]byte
var identifierSafe [(utf8.MaxRune + 1) / 8]byte
var stringlikeSafe [(utf8.MaxRune + 1) / 8]byte
var stringlikeUnsafe = make(map[rune]SafetyFlags)
var properties [0x100]Properties

type Properties uint64

const (
	StructWS Properties = 1 << iota
	DigitBase2
	DigitBase8
	DigitBase10
	LowerAF
	UpperAF
	AreaLocation
	ObjectEnd
	UUID

	EndProperties
	NoProperties Properties = 0
)

func (_this Properties) String() string {
	if _this == 0 {
		return propertyNames[_this]
	}

	isFirst := true
	builder := strings.Builder{}
	for i := Properties(1); i < EndProperties; i <<= 1 {
		if _this&i != 0 {
			if isFirst {
				isFirst = false
			} else {
				builder.WriteString(" | ")
			}
			builder.WriteString(propertyNames[i])
		}
	}
	return builder.String()
}

var propertyNames = map[Properties]string{
	NoProperties:  "NoProperties",
	StructWS:      "StructWS",
	DigitBase2:    "DigitBase2",
	DigitBase8:    "DigitBase8",
	DigitBase10:   "DigitBase10",
	LowerAF:       "LowerAF",
	UpperAF:       "UpperAF",
	AreaLocation:  "AreaLocation",
	ObjectEnd:     "ObjectEnd",
	UUID:          "UUID",
	EndProperties: "EndProperties",
}

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
