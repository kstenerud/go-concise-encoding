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

package concise_encoding

import (
	"math"
	"testing"

	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/test"
)

func TestWebsiteExampleNumericTypes(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    "boolean"       = true
    "binary-int"    = -0b10001011
    "octal-int"     = 0o644
    "decimal-int"   = -10000000
    "hex-int"       = 0xfffe0001
    "very-long-int" = 100000000000000000000000000000000000009
    "decimal-float" = -14.125
    "hex-float"     = 0x5.1ec4p+20
    "very-long-flt" = 4.957234990634579394723460546348e+100000
    "not-a-number"  = nan
    "infinity"      = inf
    "neg-infinity"  = -inf
}`, []byte{0x83, ceVer, 0x79, 0x87, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e,
		0x7d, 0x8a, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x2d, 0x69, 0x6e, 0x74,
		0x69, 0x8b, 0x89, 0x6f, 0x63, 0x74, 0x61, 0x6c, 0x2d, 0x69, 0x6e, 0x74,
		0x6a, 0xa4, 0x01, 0x8b, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2d,
		0x69, 0x6e, 0x74, 0x6d, 0x80, 0x96, 0x98, 0x00, 0x87, 0x68, 0x65, 0x78,
		0x2d, 0x69, 0x6e, 0x74, 0x6c, 0x01, 0x00, 0xfe, 0xff, 0x8d, 0x76, 0x65,
		0x72, 0x79, 0x2d, 0x6c, 0x6f, 0x6e, 0x67, 0x2d, 0x69, 0x6e, 0x74, 0x66,
		0x10, 0x09, 0x00, 0x00, 0x00, 0x40, 0x22, 0x8a, 0x09, 0x7a, 0xc4, 0x86,
		0x5a, 0xa8, 0x4c, 0x3b, 0x4b, 0x8d, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61,
		0x6c, 0x2d, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x65, 0x0f, 0xad, 0x6e, 0x89,
		0x68, 0x65, 0x78, 0x2d, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x71, 0x80, 0xd8,
		0xa3, 0x4a, 0x8d, 0x76, 0x65, 0x72, 0x79, 0x2d, 0x6c, 0x6f, 0x6e, 0x67,
		0x2d, 0x66, 0x6c, 0x74, 0x65, 0x88, 0xb4, 0x18, 0xac, 0xfe, 0x87, 0x98,
		0xb5, 0xa3, 0xd5, 0xe3, 0xdb, 0xac, 0xb4, 0x85, 0x9b, 0xd2, 0x0f, 0x8c,
		0x6e, 0x6f, 0x74, 0x2d, 0x61, 0x2d, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
		0x65, 0x80, 0x00, 0x88, 0x69, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x79,
		0x65, 0x82, 0x00, 0x8c, 0x6e, 0x65, 0x67, 0x2d, 0x69, 0x6e, 0x66, 0x69,
		0x6e, 0x69, 0x74, 0x79, 0x65, 0x83, 0x00, 0x7b},
		BD(), EvV, M(),
		S("boolean"), TT(),
		S("binary-int"), NI(0b10001011),
		S("octal-int"), PI(0o644),
		S("decimal-int"), NI(10000000),
		S("hex-int"), PI(0xfffe0001),
		S("very-long-int"), BI(test.NewBigInt("100000000000000000000000000000000000009", 10)),
		S("decimal-float"), DF(test.NewDFloat("-14.125")),
		S("hex-float"), F(0x5.1ec4p20),
		S("very-long-flt"), BDF(test.NewBDF("4.957234990634579394723460546348E+100000")),
		S("not-a-number"), NAN(),
		S("infinity"), F(math.Inf(1)),
		S("neg-infinity"), F(math.Inf(-1)),
		E(), ED())
}

func TestWebsiteExampleStrings(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    "string" = "Strings support escape sequences: \n \t \27f"
    "url"    = @"https://example.com/"
    "email"  = @"mailto:me@somewhere.com"
}`, []byte{0x83, ceVer, 0x79, 0x86, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x90,
		0x4e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x20, 0x73, 0x75, 0x70,
		0x70, 0x6f, 0x72, 0x74, 0x20, 0x65, 0x73, 0x63, 0x61, 0x70, 0x65, 0x20,
		0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x3a, 0x20, 0x0a,
		0x20, 0x09, 0x20, 0x7f, 0x83, 0x75, 0x72, 0x6c, 0x91, 0x28, 0x68, 0x74,
		0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
		0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x85, 0x65, 0x6d, 0x61, 0x69, 0x6c,
		0x91, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x74, 0x6f, 0x3a, 0x6d, 0x65, 0x40,
		0x73, 0x6f, 0x6d, 0x65, 0x77, 0x68, 0x65, 0x72, 0x65, 0x2e, 0x63, 0x6f,
		0x6d, 0x7b},
		BD(), EvV, M(),
		S("string"), S("Strings support escape sequences: \n \t \x7f"),
		S("url"), RID("https://example.com/"),
		S("email"), RID("mailto:me@somewhere.com"),
		E(), ED())
}

func TestWebsiteExampleOtherTypes(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeEncode(t, cteEOpts, nil, nil, nil, `c0
{
    "uuid" = f1ce4567-e89b-12d3-a456-426655440000
    "date" = 2019-07-01
    "time" = 18:04:00.940231541/Europe/Prague
    "timestamp" = 2010-07-15/13:28:15.415942344
    "null" = null
    "media" = |application/x-sh 23 21 2f 62 69 6e 2f 73 68 0a 0a 65 63 68 6f 20 68 65 6c 6c 6f 20 77 6f 72 6c 64 0a|
}`, []byte{0x83, ceVer, 0x79, 0x84, 0x75, 0x75, 0x69, 0x64, 0x73, 0xf1, 0xce,
		0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44,
		0x00, 0x00, 0x84, 0x64, 0x61, 0x74, 0x65, 0x99, 0xe1, 0x4c, 0x00, 0x84,
		0x74, 0x69, 0x6d, 0x65, 0x9a, 0xaf, 0x5b, 0x56, 0xc0, 0x01, 0x42, 0xfe,
		0x10, 0x45, 0x2f, 0x50, 0x72, 0x61, 0x67, 0x75, 0x65, 0x89, 0x74, 0x69,
		0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x9b, 0x46, 0x36, 0x56, 0xc6,
		0x1e, 0xae, 0xbd, 0xa3, 0x00, 0x84, 0x6e, 0x75, 0x6c, 0x6c, 0x7e, 0x85,
		0x6d, 0x65, 0x64, 0x69, 0x61, 0x94, 0xe3, 0x20, 0x61, 0x70, 0x70, 0x6c,
		0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x78, 0x2d, 0x73, 0x68,
		0x38, 0x23, 0x21, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73, 0x68, 0x0a, 0x0a,
		0x65, 0x63, 0x68, 0x6f, 0x20, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77,
		0x6f, 0x72, 0x6c, 0x64, 0x0a, 0x7b},
		BD(), EvV, M(),
		S("uuid"), UID([]byte{0xf1, 0xce, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44, 0x00, 0x00}),
		S("date"), CT(test.NewDate(2019, 7, 1)),
		S("time"), CT(test.NewTime(18, 4, 0, 940231541, "Europe/Prague")),
		S("timestamp"), CT(test.NewTS(2010, 7, 15, 13, 28, 15, 415942344, "Etc/UTC")),
		S("null"), N(),
		S("media"), MB(), AC(16, false), AD([]byte("application/x-sh")),
		AC(28, false), AD([]byte{0x23, 0x21, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73,
			0x68, 0x0a, 0x0a, 0x65, 0x63, 0x68, 0x6f, 0x20, 0x68, 0x65, 0x6c,
			0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0a}),
		E(), ED())
}

func TestWebsiteExampleContainersArrays(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "
	cteEOpts.DefaultFormats.Array.Uint8 = options.CTEEncodingFormatHexadecimalZeroFilled
	cteEOpts.DefaultFormats.Array.Uint16 = options.CTEEncodingFormatHexadecimalZeroFilled
	cteEOpts.DefaultFormats.Array.Float32 = options.CTEEncodingFormatUnset

	assertDecodeEncode(t, cteEOpts, nil, nil, nil, `c0
{
    "list" = [
        1
        2.5
        "a string"
    ]
    "map" = {
        "one" = 1
        2 = "two"
        "today" = 2020-09-10
    }
    "bytes" = |u8x 01 ff de ad be ef|
    "int16 array" = |i16 7374 17466 -9957|
    "uint16 hex" = |u16x 91fe 443a 9c15|
    "float32 array" = |f32 1.5e+10 -8.31e-12|
}`, []byte{0x83, ceVer, 0x79, 0x84, 0x6c, 0x69, 0x73, 0x74, 0x7a, 0x01, 0x65,
		0x06, 0x19, 0x88, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x7b,
		0x83, 0x6d, 0x61, 0x70, 0x79, 0x83, 0x6f, 0x6e, 0x65, 0x01, 0x02, 0x83,
		0x74, 0x77, 0x6f, 0x85, 0x74, 0x6f, 0x64, 0x61, 0x79, 0x99, 0x2a, 0x51,
		0x00, 0x7b, 0x85, 0x62, 0x79, 0x74, 0x65, 0x73, 0x95, 0x0c, 0x01, 0xff,
		0xde, 0xad, 0xbe, 0xef, 0x8b, 0x69, 0x6e, 0x74, 0x31, 0x36, 0x20, 0x61,
		0x72, 0x72, 0x61, 0x79, 0x94, 0x23, 0xce, 0x1c, 0x3a, 0x44, 0x1b, 0xd9,
		0x8a, 0x75, 0x69, 0x6e, 0x74, 0x31, 0x36, 0x20, 0x68, 0x65, 0x78, 0x94,
		0x13, 0xfe, 0x91, 0x3a, 0x44, 0x15, 0x9c, 0x8d, 0x66, 0x6c, 0x6f, 0x61,
		0x74, 0x33, 0x32, 0x20, 0x61, 0x72, 0x72, 0x61, 0x79, 0x94, 0x82, 0x76,
		0x84, 0x5f, 0x50, 0xea, 0x30, 0x12, 0xad, 0x7b},
		BD(), EvV, M(),
		S("list"), L(), PI(1), DF(NewDFloat("2.5")), S("a string"), E(),
		S("map"), M(), S("one"), PI(1), PI(2), S("two"), S("today"), CT(test.NewDate(2020, 9, 10)), E(),
		S("bytes"), AU8([]byte{0x01, 0xff, 0xde, 0xad, 0xbe, 0xef}),
		S("int16 array"), AI16([]int16{7374, 17466, -9957}),
		S("uint16 hex"), AU16([]uint16{0x91fe, 0x443a, 0x9c15}),
		S("float32 array"), AF32([]float32{1.5e+10, -8.31e-12}),
		E(), ED())
}

func TestWebsiteExampleMarkup(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    "main-view" = <View,
        <Image "src"=@"img/avatar-image.jpg">
        <Text "id"="HelloText",
            Hello! Please choose a name!
        >
        // OnChange contains code which might have problematic characters.
        // Use verbatim sequences (\.IDENTIFIER ... IDENTIFIER) to handle this.
        <TextInput "id"="NameInput" "style"={"height"=40 "color"="gray"} "OnChange"="\.@@
            NameInput.Parent.InsertRawAfter(NameInput, '<Image "src"=@"img/check.svg">')
            HelloText.SetText("Hello, " + NameInput.Text + "!")
            @@",
            Name me!
        >
    >
}`, []byte{0x83, ceVer, 0x79, 0x89, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x76, 0x69,
		0x65, 0x77, 0x78, 0x04, 0x56, 0x69, 0x65, 0x77, 0x7b, 0x78, 0x05, 0x49,
		0x6d, 0x61, 0x67, 0x65, 0x83, 0x73, 0x72, 0x63, 0x91, 0x28, 0x69, 0x6d,
		0x67, 0x2f, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x2d, 0x69, 0x6d, 0x61,
		0x67, 0x65, 0x2e, 0x6a, 0x70, 0x67, 0x7b, 0x7b, 0x78, 0x04, 0x54, 0x65,
		0x78, 0x74, 0x82, 0x69, 0x64, 0x89, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x54,
		0x65, 0x78, 0x74, 0x7b, 0x90, 0x38, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x21,
		0x20, 0x50, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x20, 0x63, 0x68, 0x6f, 0x6f,
		0x73, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x21, 0x7b, 0x78,
		0x09, 0x54, 0x65, 0x78, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x82, 0x69,
		0x64, 0x89, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x85,
		0x73, 0x74, 0x79, 0x6c, 0x65, 0x79, 0x86, 0x68, 0x65, 0x69, 0x67, 0x68,
		0x74, 0x28, 0x85, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x84, 0x67, 0x72, 0x61,
		0x79, 0x7b, 0x88, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x90,
		0xca, 0x02, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x2e,
		0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72,
		0x74, 0x52, 0x61, 0x77, 0x41, 0x66, 0x74, 0x65, 0x72, 0x28, 0x4e, 0x61,
		0x6d, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x2c, 0x20, 0x27, 0x3c, 0x49,
		0x6d, 0x61, 0x67, 0x65, 0x20, 0x22, 0x73, 0x72, 0x63, 0x22, 0x3d, 0x40,
		0x22, 0x69, 0x6d, 0x67, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x73,
		0x76, 0x67, 0x22, 0x3e, 0x27, 0x29, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x54, 0x65, 0x78, 0x74, 0x2e, 0x53, 0x65, 0x74, 0x54, 0x65, 0x78, 0x74,
		0x28, 0x22, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x22, 0x20, 0x2b,
		0x20, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x54,
		0x65, 0x78, 0x74, 0x20, 0x2b, 0x20, 0x22, 0x21, 0x22, 0x29, 0x0a, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b,
		0x88, 0x4e, 0x61, 0x6d, 0x65, 0x20, 0x6d, 0x65, 0x21, 0x7b, 0x7b, 0x7b},
		BD(), EvV, M(),
		S("main-view"), MUP("View"), E(),
		MUP("Image"), S("src"), RID("img/avatar-image.jpg"), E(), E(),
		MUP("Text"), S("id"), S("HelloText"), E(),
		S("Hello! Please choose a name!"),
		E(),
		COM(false, ` OnChange contains code which might have problematic characters.`),
		COM(false, ` Use verbatim sequences (\.IDENTIFIER ... IDENTIFIER) to handle this.`),
		MUP("TextInput"), S("id"), S("NameInput"), S("style"), M(),
		S("height"), PI(40), S("color"), S("gray"), E(),
		S("OnChange"), S(`            NameInput.Parent.InsertRawAfter(NameInput, '<Image "src"=@"img/check.svg">')
            HelloText.SetText("Hello, " + NameInput.Text + "!")
            `), E(),
		S("Name me!"), E(),
		E(), E(), ED())
}

func TestWebsiteExampleReferences(t *testing.T) {
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    // Entire map will be referenced later as $id1
    "marked object" = &id1:{
        "recursive" = $id1
    }
    "ref1" = $id1
    "ref2" = $id1

    // Reference pointing to part of another document.
    "outside ref" = $"https://xyz.com/document.cte#some_id"
}`, []byte{0x83, ceVer, 0x79, 0x8d, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x20,
		0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x97, 0x03, 0x69, 0x64, 0x31, 0x79,
		0x89, 0x72, 0x65, 0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x98, 0x03,
		0x69, 0x64, 0x31, 0x7b, 0x84, 0x72, 0x65, 0x66, 0x31, 0x98, 0x03, 0x69,
		0x64, 0x31, 0x84, 0x72, 0x65, 0x66, 0x32, 0x98, 0x03, 0x69, 0x64, 0x31,
		0x8b, 0x6f, 0x75, 0x74, 0x73, 0x69, 0x64, 0x65, 0x20, 0x72, 0x65, 0x66,
		0x94, 0xe2, 0x48, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x78,
		0x79, 0x7a, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d,
		0x65, 0x6e, 0x74, 0x2e, 0x63, 0x74, 0x65, 0x23, 0x73, 0x6f, 0x6d, 0x65,
		0x5f, 0x69, 0x64, 0x7b},
		BD(), EvV, M(),
		COM(false, " Entire map will be referenced later as $id1"),
		S("marked object"), MARK("id1"), M(),
		S("recursive"), REF("id1"),
		E(),
		S("ref1"), REF("id1"),
		S("ref2"), REF("id1"),
		COM(false, " Reference pointing to part of another document."),
		S("outside ref"), RIDREF("https://xyz.com/document.cte#some_id"),
		E(), ED())
}

func TestWebsiteGraphs(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    "vertices" = [
        &a:{}
        &b:{}
        &c:{}
        &d:{}
    ]
    "edges" = [
        @($a {"weight"=4 "direction"="both"} $b)
        @($a {"weight"=3 "direction"="both"} $c)
        @($b {"weight"=1 "direction"="both"} $c)
        @($b {"weight"=1 "direction"="both"} $d)
        @($c {"weight"=4 "direction"="both"} $d)
    ]
}`, []byte{0x83, ceVer, 0x79, 0x88, 0x76, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65,
		0x73, 0x7a, 0x97, 0x01, 0x61, 0x79, 0x7b, 0x97, 0x01, 0x62, 0x79, 0x7b,
		0x97, 0x01, 0x63, 0x79, 0x7b, 0x97, 0x01, 0x64, 0x79, 0x7b, 0x7b, 0x85,
		0x65, 0x64, 0x67, 0x65, 0x73, 0x7a, 0x76, 0x98, 0x01, 0x61, 0x79, 0x86,
		0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x04, 0x89, 0x64, 0x69, 0x72, 0x65,
		0x63, 0x74, 0x69, 0x6f, 0x6e, 0x84, 0x62, 0x6f, 0x74, 0x68, 0x7b, 0x98,
		0x01, 0x62, 0x76, 0x98, 0x01, 0x61, 0x79, 0x86, 0x77, 0x65, 0x69, 0x67,
		0x68, 0x74, 0x03, 0x89, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f,
		0x6e, 0x84, 0x62, 0x6f, 0x74, 0x68, 0x7b, 0x98, 0x01, 0x63, 0x76, 0x98,
		0x01, 0x62, 0x79, 0x86, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x01, 0x89,
		0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x84, 0x62, 0x6f,
		0x74, 0x68, 0x7b, 0x98, 0x01, 0x63, 0x76, 0x98, 0x01, 0x62, 0x79, 0x86,
		0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x01, 0x89, 0x64, 0x69, 0x72, 0x65,
		0x63, 0x74, 0x69, 0x6f, 0x6e, 0x84, 0x62, 0x6f, 0x74, 0x68, 0x7b, 0x98,
		0x01, 0x64, 0x76, 0x98, 0x01, 0x63, 0x79, 0x86, 0x77, 0x65, 0x69, 0x67,
		0x68, 0x74, 0x04, 0x89, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f,
		0x6e, 0x84, 0x62, 0x6f, 0x74, 0x68, 0x7b, 0x98, 0x01, 0x64, 0x7b, 0x7b},
		BD(), EvV, M(),
		S("vertices"), L(),
		MARK("a"), M(), E(),
		MARK("b"), M(), E(),
		MARK("c"), M(), E(),
		MARK("d"), M(), E(),
		E(),
		S("edges"), L(),
		EDGE(), REF("a"), M(), S("weight"), PI(4), S("direction"), S("both"), E(), REF("b"),
		EDGE(), REF("a"), M(), S("weight"), PI(3), S("direction"), S("both"), E(), REF("c"),
		EDGE(), REF("b"), M(), S("weight"), PI(1), S("direction"), S("both"), E(), REF("c"),
		EDGE(), REF("b"), M(), S("weight"), PI(1), S("direction"), S("both"), E(), REF("d"),
		EDGE(), REF("c"), M(), S("weight"), PI(4), S("direction"), S("both"), E(), REF("d"),
		E(),
		E(), ED())
}

func TestWebsiteTrees(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
(2
    (7
        2
        1
        (6
            5
            8
        )
    )
    (5
        (9
            4
        )
    )
)`, []byte{0x83, ceVer, 0x77, 0x02, 0x77, 0x07, 0x02, 0x01, 0x77, 0x06, 0x05,
		0x08, 0x7b, 0x7b, 0x77, 0x05, 0x77, 0x09, 0x04, 0x7b, 0x7b, 0x7b},
		BD(), EvV,
		NODE(), PI(2),
		NODE(), PI(7), PI(2), PI(1),
		NODE(), PI(6), PI(5), PI(8),
		E(),
		E(),
		NODE(), PI(5),
		NODE(), PI(9), PI(4),
		E(),
		E(),
		E(),
		ED())
}

func TestWebsiteCustomTypes(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c0
{
    // Custom types are user-defined, with user-supplied codecs.
    "custom text"   = |ct cplx(2.94+3i)|
    "custom binary" = |cb 04 f6 28 3c 40 00 00 40 40|
}`, []byte{0x83, ceVer, 0x79, 0x8b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x20,
		0x74, 0x65, 0x78, 0x74, 0x93, 0x1a, 0x63, 0x70, 0x6c, 0x78, 0x28, 0x32,
		0x2e, 0x39, 0x34, 0x2b, 0x33, 0x69, 0x29, 0x8d, 0x63, 0x75, 0x73, 0x74,
		0x6f, 0x6d, 0x20, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x92, 0x12, 0x04,
		0xf6, 0x28, 0x3c, 0x40, 0x00, 0x00, 0x40, 0x40, 0x7b},
		BD(), EvV, M(),
		COM(false, " Custom types are user-defined, with user-supplied codecs."),
		S("custom text"), CUT("cplx(2.94+3i)"),
		S("custom binary"), CUB([]byte{0x04, 0xf6, 0x28, 0x3c, 0x40, 0x00, 0x00, 0x40, 0x40}),
		E(), ED())
}
