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

func TestEncodeDecodeNil(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), NA(), ED())
}

func TestEncodeDecodeTrue(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), TT(), ED())
}

func TestEncodeDecodeFalse(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), FF(), ED())
}

func TestEncodeDecodePositiveInt(t *testing.T) {
	assertEncodeDecodeCTE(t, BD(), V(1), PI(0), ED())
	assertEncodeDecodeCTE(t, BD(), V(1), PI(1), ED())
	assertEncodeDecodeCBE(t, BD(), V(1), I(0), ED())
	assertEncodeDecodeCBE(t, BD(), V(1), I(1), ED())
	assertEncodeDecode(t, BD(), V(1), PI(104), ED())
	assertEncodeDecode(t, BD(), V(1), PI(10405), ED())
	assertEncodeDecode(t, BD(), V(1), PI(999999), ED())
	assertEncodeDecode(t, BD(), V(1), PI(7234859234423), ED())
}

func TestEncodeDecodeNegativeInt(t *testing.T) {
	assertEncodeDecodeCTE(t, BD(), V(1), NI(1), ED())
	assertEncodeDecodeCBE(t, BD(), V(1), I(-1), ED())
	assertEncodeDecode(t, BD(), V(1), NI(104), ED())
	assertEncodeDecode(t, BD(), V(1), NI(10405), ED())
	assertEncodeDecode(t, BD(), V(1), NI(999999), ED())
	assertEncodeDecode(t, BD(), V(1), NI(7234859234423), ED())
}

func TestEncodeDecodeFloat(t *testing.T) {
	// CTE will convert to decimal float
	assertEncodeDecodeCBE(t, BD(), V(1), F(1.5), ED())
	assertEncodeDecode(t, BD(), V(1), DF(test.NewDFloat("1.5")), ED())
	assertEncodeDecodeCBE(t, BD(), V(1), F(-51.455e-16), ED())
	assertEncodeDecode(t, BD(), V(1), DF(test.NewDFloat("-51.455e-16")), ED())
}

func TestEncodeDecodeNan(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), NAN(), ED())
	assertEncodeDecode(t, BD(), V(1), SNAN(), ED())
}

func TestEncodeDecodeUUID(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), UUID([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}), ED())
}

func TestEncodeDecodeTime(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), CT(test.NewDate(2000, 1, 1)), ED())

	assertEncodeDecode(t, BD(), V(1), CT(test.NewTime(1, 45, 0, 0, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTime(23, 59, 59, 101000000, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTime(10, 0, 1, 930000000, "America/Los_Angeles")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTimeLL(10, 0, 1, 930000000, 8992, 110)), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTimeLL(10, 0, 1, 930000000, 0, 0)), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTimeLL(10, 0, 1, 930000000, 100, 100)), ED())

	assertEncodeDecode(t, BD(), V(1), CT(test.NewTS(2000, 1, 1, 19, 31, 44, 901554000, "")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTS(-50000, 12, 29, 1, 1, 1, 305, "Etc/UTC")), ED())
	assertEncodeDecode(t, BD(), V(1), CT(test.NewTSLL(2954, 8, 31, 12, 31, 15, 335523, 3154, 16004)), ED())
}

func TestEncodeDecodeBytes(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), AU8([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeCustom(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), CUB([]byte{1, 2, 3, 4, 5, 6, 7}), ED())
}

func TestEncodeDecodeRID(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, BD(), V(1), RID("http://example.com"), ED())
}

func TestEncodeDecodeString(t *testing.T) {
	// TODO: More complex tests
	assertEncodeDecode(t, BD(), V(1), S("A string"), ED())
}

func TestEncodeDecodeList(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), L(), E(), ED())
	assertEncodeDecode(t, BD(), V(1), L(), PI(1000), E(), ED())
}

func TestEncodeDecodeMap(t *testing.T) {
	assertEncodeDecode(t, BD(), V(1), M(), E(), ED())
	assertEncodeDecode(t, BD(), V(1), M(), S("a"), NI(1000), E(), ED())
	assertEncodeDecode(t, BD(), V(1), M(), S("some NA"), NA(), DF(test.NewDFloat("1.1")), S("somefloat"), E(), ED())
}

func TestWebsiteExampleNumericTypes(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c1
{
    boolean       = @true
    binary-int    = -0b10001011
    octal-int     = 0o644
    decimal-int   = -10000000
    hex-int       = 0xfffe0001
    decimal-float = -14.125
    hex-float     = 0x5.1ec4p20
    not-a-number  = @nan
    infinity      = @inf
    n-infinity    = -@inf
}`, []byte{0x03, 0x01, 0x79, 0x87, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x7d,
		0x8a, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x2d, 0x69, 0x6e, 0x74, 0x69,
		0x8b, 0x89, 0x6f, 0x63, 0x74, 0x61, 0x6c, 0x2d, 0x69, 0x6e, 0x74, 0x6a,
		0xa4, 0x01, 0x8b, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2d, 0x69,
		0x6e, 0x74, 0x6d, 0x80, 0x96, 0x98, 0x00, 0x87, 0x68, 0x65, 0x78, 0x2d,
		0x69, 0x6e, 0x74, 0x6c, 0x01, 0x00, 0xfe, 0xff, 0x8d, 0x64, 0x65, 0x63,
		0x69, 0x6d, 0x61, 0x6c, 0x2d, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x65, 0x0f,
		0xad, 0x6e, 0x89, 0x68, 0x65, 0x78, 0x2d, 0x66, 0x6c, 0x6f, 0x61, 0x74,
		0x71, 0x80, 0xd8, 0xa3, 0x4a, 0x8c, 0x6e, 0x6f, 0x74, 0x2d, 0x61, 0x2d,
		0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x65, 0x80, 0x00, 0x88, 0x69, 0x6e,
		0x66, 0x69, 0x6e, 0x69, 0x74, 0x79, 0x65, 0x82, 0x00, 0x8a, 0x6e, 0x2d,
		0x69, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x79, 0x65, 0x83, 0x00, 0x7b},
		BD(), V(1), M(),
		S("boolean"), TT(),
		S("binary-int"), NI(0b10001011),
		S("octal-int"), PI(0o644),
		S("decimal-int"), NI(10000000),
		S("hex-int"), PI(0xfffe0001),
		S("decimal-float"), DF(test.NewDFloat("-14.125")),
		S("hex-float"), F(0x5.1ec4p20),
		S("not-a-number"), NAN(),
		S("infinity"), F(math.Inf(1)),
		S("n-infinity"), F(math.Inf(-1)),
		E(), ED())
}

func TestWebsiteExampleStrings(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeEncode(t, cteEOpts, nil, nil, nil, `c1
{
    unquoted-string = no-quotes-needed
    quoted-string = "A string delimited by quotes"
    url = |r https://example.com/|
    email = |r mailto:me@somewhere.com|
}`, []byte{0x03, 0x01, 0x79, 0x8f, 0x75, 0x6e, 0x71, 0x75, 0x6f, 0x74, 0x65,
		0x64, 0x2d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x90, 0x20, 0x6e, 0x6f,
		0x2d, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x73, 0x2d, 0x6e, 0x65, 0x65, 0x64,
		0x65, 0x64, 0x8d, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x64, 0x2d, 0x73, 0x74,
		0x72, 0x69, 0x6e, 0x67, 0x90, 0x38, 0x41, 0x20, 0x73, 0x74, 0x72, 0x69,
		0x6e, 0x67, 0x20, 0x64, 0x65, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x65, 0x64,
		0x20, 0x62, 0x79, 0x20, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x73, 0x83, 0x75,
		0x72, 0x6c, 0x91, 0x28, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
		0x85, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x91, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
		0x74, 0x6f, 0x3a, 0x6d, 0x65, 0x40, 0x73, 0x6f, 0x6d, 0x65, 0x77, 0x68,
		0x65, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x7b},
		BD(), V(1), M(),
		S("unquoted-string"), S("no-quotes-needed"),
		S("quoted-string"), S("A string delimited by quotes"),
		S("url"), RID("https://example.com/"),
		S("email"), RID("mailto:me@somewhere.com"),
		E(), ED())
}

func TestWebsiteExampleOtherTypes(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeEncode(t, cteEOpts, nil, nil, nil, `c1
{
    uuid = @f1ce4567-e89b-12d3-a456-426655440000
    date = 2019-07-01
    time = 18:04:00.940231541/Europe/Prague
    timestamp = 2010-07-15/13:28:15.415942344
    not-available = @na
}`, []byte{0x03, 0x01, 0x79, 0x84, 0x75, 0x75, 0x69, 0x64, 0x73, 0xf1, 0xce,
		0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44,
		0x00, 0x00, 0x84, 0x64, 0x61, 0x74, 0x65, 0x99, 0xe1, 0x4c, 0x00, 0x84,
		0x74, 0x69, 0x6d, 0x65, 0x9a, 0xaf, 0x5b, 0x56, 0xc0, 0x01, 0x42, 0xfe,
		0x10, 0x45, 0x2f, 0x50, 0x72, 0x61, 0x67, 0x75, 0x65, 0x89, 0x74, 0x69,
		0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x9b, 0x46, 0x36, 0x56, 0xc6,
		0x1e, 0xae, 0xbd, 0xa3, 0x00, 0x8d, 0x6e, 0x6f, 0x74, 0x2d, 0x61, 0x76,
		0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x7e, 0x7b},
		BD(), V(1), M(),
		S("uuid"), UUID([]byte{0xf1, 0xce, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44, 0x00, 0x00}),
		S("date"), CT(test.NewDate(2019, 7, 1)),
		S("time"), CT(test.NewTime(18, 4, 0, 940231541, "Europe/Prague")),
		S("timestamp"), CT(test.NewTS(2010, 7, 15, 13, 28, 15, 415942344, "Etc/UTC")),
		S("not-available"), NA(),
		E(), ED())
}

func TestWebsiteExampleContainersArrays(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "
	cteEOpts.DefaultFormats.Array.Uint8 = options.CTEEncodingFormatHexadecimalZeroFilled
	cteEOpts.DefaultFormats.Array.Uint16 = options.CTEEncodingFormatHexadecimalZeroFilled
	cteEOpts.DefaultFormats.Array.Float32 = options.CTEEncodingFormatUnset

	assertDecodeEncode(t, cteEOpts, nil, nil, nil, `c1
{
    list = [
        1
        2.5
        "a string"
    ]
    map = {
        one = 1
        2 = two
        today = 2020-09-10
    }
    bytes = |u8x 01 ff de ad be ef|
    uint16-hex-array = |u16x 91fe 443a 9c15|
    int16-array = |i16 7374 17466 -9957|
    float32-array = |f32 1.5e+10 -8.31e-12|
}`, []byte{0x03, 0x01, 0x79, 0x84, 0x6c, 0x69, 0x73, 0x74, 0x7a, 0x01, 0x65,
		0x06, 0x19, 0x88, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x7b,
		0x83, 0x6d, 0x61, 0x70, 0x79, 0x83, 0x6f, 0x6e, 0x65, 0x01, 0x02, 0x83,
		0x74, 0x77, 0x6f, 0x85, 0x74, 0x6f, 0x64, 0x61, 0x79, 0x99, 0x2a, 0x51,
		0x00, 0x7b, 0x85, 0x62, 0x79, 0x74, 0x65, 0x73, 0x94, 0x68, 0x0c, 0x01,
		0xff, 0xde, 0xad, 0xbe, 0xef, 0x90, 0x20, 0x75, 0x69, 0x6e, 0x74, 0x31,
		0x36, 0x2d, 0x68, 0x65, 0x78, 0x2d, 0x61, 0x72, 0x72, 0x61, 0x79, 0x94,
		0x6a, 0x06, 0xfe, 0x91, 0x3a, 0x44, 0x15, 0x9c, 0x8b, 0x69, 0x6e, 0x74,
		0x31, 0x36, 0x2d, 0x61, 0x72, 0x72, 0x61, 0x79, 0x94, 0x6b, 0x06, 0xce,
		0x1c, 0x3a, 0x44, 0x1b, 0xd9, 0x8d, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x33,
		0x32, 0x2d, 0x61, 0x72, 0x72, 0x61, 0x79, 0x94, 0x71, 0x04, 0x76, 0x84,
		0x5f, 0x50, 0xea, 0x30, 0x12, 0xad, 0x7b},
		BD(), V(1), M(),
		S("list"), L(), PI(1), DF(NewDFloat("2.5")), S("a string"), E(),
		S("map"), M(), S("one"), PI(1), PI(2), S("two"), S("today"), CT(test.NewDate(2020, 9, 10)), E(),
		S("bytes"), AU8([]byte{0x01, 0xff, 0xde, 0xad, 0xbe, 0xef}),
		S("uint16-hex-array"), AU16([]uint16{0x91fe, 0x443a, 0x9c15}),
		S("int16-array"), AI16([]int16{7374, 17466, -9957}),
		S("float32-array"), AF32([]float32{1.5e+10, -8.31e-12}),
		E(), ED())
}

func TestWebsiteExampleMarkup(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c1
{
    main-view = <View,
        <Image src=|r images/avatar-image.jpg|>
        <Text id=HelloText,
            Hello! Please choose a name!
        >
        <TextInput id=NameInput style={height=40 borderColor=gray} OnChange="\.@@
            HelloText.SetText("Hello, " + NameInput.Text + "!")
        @@",
            Name me!
        >
    >
}`, []byte{0x03, 0x01, 0x79, 0x89, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x76, 0x69,
		0x65, 0x77, 0x78, 0x84, 0x56, 0x69, 0x65, 0x77, 0x7b, 0x78, 0x85, 0x49,
		0x6d, 0x61, 0x67, 0x65, 0x83, 0x73, 0x72, 0x63, 0x91, 0x2e, 0x69, 0x6d,
		0x61, 0x67, 0x65, 0x73, 0x2f, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x2d,
		0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x6a, 0x70, 0x67, 0x7b, 0x7b, 0x78,
		0x84, 0x54, 0x65, 0x78, 0x74, 0x82, 0x69, 0x64, 0x89, 0x48, 0x65, 0x6c,
		0x6c, 0x6f, 0x54, 0x65, 0x78, 0x74, 0x7b, 0x90, 0x38, 0x48, 0x65, 0x6c,
		0x6c, 0x6f, 0x21, 0x20, 0x50, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x20, 0x63,
		0x68, 0x6f, 0x6f, 0x73, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x61, 0x6d, 0x65,
		0x21, 0x7b, 0x78, 0x89, 0x54, 0x65, 0x78, 0x74, 0x49, 0x6e, 0x70, 0x75,
		0x74, 0x82, 0x69, 0x64, 0x89, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x70,
		0x75, 0x74, 0x85, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x79, 0x86, 0x68, 0x65,
		0x69, 0x67, 0x68, 0x74, 0x28, 0x8b, 0x62, 0x6f, 0x72, 0x64, 0x65, 0x72,
		0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x84, 0x67, 0x72, 0x61, 0x79, 0x7b, 0x88,
		0x4f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x90, 0x90, 0x01, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x48,
		0x65, 0x6c, 0x6c, 0x6f, 0x54, 0x65, 0x78, 0x74, 0x2e, 0x53, 0x65, 0x74,
		0x54, 0x65, 0x78, 0x74, 0x28, 0x22, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c,
		0x20, 0x22, 0x20, 0x2b, 0x20, 0x4e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x70,
		0x75, 0x74, 0x2e, 0x54, 0x65, 0x78, 0x74, 0x20, 0x2b, 0x20, 0x22, 0x21,
		0x22, 0x29, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7b,
		0x88, 0x4e, 0x61, 0x6d, 0x65, 0x20, 0x6d, 0x65, 0x21, 0x7b, 0x7b, 0x7b},
		BD(), V(1), M(),
		S("main-view"), MUP(), S("View"), E(),
		MUP(), S("Image"), S("src"), RID("images/avatar-image.jpg"), E(), E(),
		MUP(), S("Text"), S("id"), S("HelloText"), E(),
		S("Hello! Please choose a name!"),
		E(),
		MUP(), S("TextInput"), S("id"), S("NameInput"), S("style"), M(),
		S("height"), PI(40), S("borderColor"), S("gray"), E(),
		S("OnChange"), S(`            HelloText.SetText("Hello, " + NameInput.Text + "!")
        `), E(),
		S("Name me!"), E(),
		E(), E(), ED())
}

func TestWebsiteExampleReferences(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c1
{
    marked_object    = &id1:{
        description = "This map will be referenced later as $id1"
        value = -@inf
        child_elements = @na
        recursive = $id1
    }
    ref1             = $id1
    ref2             = $id1
    outside_ref      = $|r https://xyz.com/document.cte#some_id|
}`, []byte{0x03, 0x01, 0x79, 0x8d, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x64, 0x5f,
		0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x97, 0x83, 0x69, 0x64, 0x31, 0x79,
		0x8b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
		0x90, 0x52, 0x54, 0x68, 0x69, 0x73, 0x20, 0x6d, 0x61, 0x70, 0x20, 0x77,
		0x69, 0x6c, 0x6c, 0x20, 0x62, 0x65, 0x20, 0x72, 0x65, 0x66, 0x65, 0x72,
		0x65, 0x6e, 0x63, 0x65, 0x64, 0x20, 0x6c, 0x61, 0x74, 0x65, 0x72, 0x20,
		0x61, 0x73, 0x20, 0x24, 0x69, 0x64, 0x31, 0x85, 0x76, 0x61, 0x6c, 0x75,
		0x65, 0x65, 0x83, 0x00, 0x8e, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x5f, 0x65,
		0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x7e, 0x89, 0x72, 0x65, 0x63,
		0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x98, 0x83, 0x69, 0x64, 0x31, 0x7b,
		0x84, 0x72, 0x65, 0x66, 0x31, 0x98, 0x83, 0x69, 0x64, 0x31, 0x84, 0x72,
		0x65, 0x66, 0x32, 0x98, 0x83, 0x69, 0x64, 0x31, 0x8b, 0x6f, 0x75, 0x74,
		0x73, 0x69, 0x64, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x98, 0x91, 0x48, 0x68,
		0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x78, 0x79, 0x7a, 0x2e, 0x63,
		0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
		0x63, 0x74, 0x65, 0x23, 0x73, 0x6f, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x7b},
		BD(), V(1), M(),
		S("marked_object"), MARK(), S("id1"), M(),
		S("description"), S("This map will be referenced later as $id1"),
		S("value"), F(math.Inf(-1)),
		S("child_elements"), NA(),
		S("recursive"), REF(), S("id1"),
		E(),
		S("ref1"), REF(), S("id1"),
		S("ref2"), REF(), S("id1"),
		S("outside_ref"), REF(), RID("https://xyz.com/document.cte#some_id"),
		E(), ED())
}

func TestWebsiteExampleMetadataComments(t *testing.T) {
	defer test.PassThroughPanics(true)()
	cteEOpts := options.DefaultCTEEncoderOptions()
	cteEOpts.Indent = "    "

	assertDecodeCBECTE(t, cteEOpts, nil, nil, nil, `c1
// Metadata about the entire documents
(
    // _ct is the creation time, _d is description, _v is version.
    // See common generic metadata spec.
    _ct = 2019-9-1/22:14:01
    _d = "Some description"
    _v = "1.1.0"
    whatever = "some arbitrary data"
)
{
    /* Comments look very C-like, except:
       /* Nested comments are allowed! */
    */

    // Double-slash comments are also possible.

    (info = "something interesting about a_list")
    a_list           = [1 2 3]
}`, []byte{0x03, 0x01, 0x76, 0x90, 0x46, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
		0x74, 0x61, 0x20, 0x61, 0x62, 0x6f, 0x75, 0x74, 0x20, 0x74, 0x68, 0x65,
		0x20, 0x65, 0x6e, 0x74, 0x69, 0x72, 0x65, 0x20, 0x64, 0x6f, 0x63, 0x75,
		0x6d, 0x65, 0x6e, 0x74, 0x73, 0x7b, 0x77, 0x76, 0x90, 0x76, 0x5f, 0x63,
		0x74, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x72, 0x65,
		0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x2c, 0x20,
		0x5f, 0x64, 0x20, 0x69, 0x73, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
		0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x5f, 0x76, 0x20, 0x69, 0x73,
		0x20, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x7b, 0x76, 0x90,
		0x42, 0x53, 0x65, 0x65, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x20,
		0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x20, 0x6d, 0x65, 0x74, 0x61,
		0x64, 0x61, 0x74, 0x61, 0x20, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x7b, 0x83,
		0x5f, 0x63, 0x74, 0x9b, 0x08, 0x1c, 0x1b, 0xd2, 0x04, 0x82, 0x5f, 0x64,
		0x90, 0x20, 0x53, 0x6f, 0x6d, 0x65, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72,
		0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x82, 0x5f, 0x76, 0x85, 0x31, 0x2e,
		0x31, 0x2e, 0x30, 0x88, 0x77, 0x68, 0x61, 0x74, 0x65, 0x76, 0x65, 0x72,
		0x90, 0x26, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x61, 0x72, 0x62, 0x69, 0x74,
		0x72, 0x61, 0x72, 0x79, 0x20, 0x64, 0x61, 0x74, 0x61, 0x7b, 0x79, 0x76,
		0x90, 0x44, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x20, 0x6c,
		0x6f, 0x6f, 0x6b, 0x20, 0x76, 0x65, 0x72, 0x79, 0x20, 0x43, 0x2d, 0x6c,
		0x69, 0x6b, 0x65, 0x2c, 0x20, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x3a,
		0x76, 0x90, 0x38, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x20, 0x63, 0x6f,
		0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x20, 0x61, 0x72, 0x65, 0x20, 0x61,
		0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x21, 0x7b, 0x7b, 0x76, 0x90, 0x50,
		0x44, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x2d, 0x73, 0x6c, 0x61, 0x73, 0x68,
		0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x20, 0x61, 0x72,
		0x65, 0x20, 0x61, 0x6c, 0x73, 0x6f, 0x20, 0x70, 0x6f, 0x73, 0x73, 0x69,
		0x62, 0x6c, 0x65, 0x2e, 0x7b, 0x77, 0x84, 0x69, 0x6e, 0x66, 0x6f, 0x90,
		0x44, 0x73, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x20, 0x69,
		0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x61,
		0x62, 0x6f, 0x75, 0x74, 0x20, 0x61, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x7b,
		0x86, 0x61, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x7a, 0x01, 0x02, 0x03, 0x7b,
		0x7b},
		BD(), V(1),
		CMT(), S("Metadata about the entire documents"), E(),
		META(),
		CMT(), S("_ct is the creation time, _d is description, _v is version."), E(),
		CMT(), S("See common generic metadata spec."), E(),
		S("_ct"), CT(test.NewTS(2019, 9, 1, 22, 14, 01, 0, "")),
		S("_d"), S("Some description"),
		S("_v"), S("1.1.0"),
		S("whatever"), S("some arbitrary data"),
		E(),
		M(),
		CMT(), S("Comments look very C-like, except:"), CMT(), S("Nested comments are allowed!"), E(), E(),
		CMT(), S("Double-slash comments are also possible."), E(),
		META(),
		S("info"), S("something interesting about a_list"),
		E(),
		S("a_list"), L(), PI(1), PI(2), PI(3), E(),
		E(),
		ED())
}
