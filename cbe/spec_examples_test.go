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

package cbe

import (
	"testing"

	"github.com/kstenerud/go-concise-encoding/test"
)

func assertCodecEvents(t *testing.T, documentContents []byte, documentEvents ...*test.TEvent) {
	document := append([]byte{header, ceVer}, documentContents...)
	events := append(append([]*test.TEvent{BD(), EvV}, documentEvents...), ED())
	assertDecode(t, nil, document, events...)
	assertEncode(t, nil, document, events...)
}

func TestSpecExamplesInteger(t *testing.T) {
	// [60] = 96
	assertCodecEvents(t, []byte{96}, I(96))

	// [00] = 0
	assertCodecEvents(t, []byte{0}, I(0))

	// [ca] = -54
	assertCodecEvents(t, []byte{0xca}, I(-54))

	// [68 7f] = 127
	assertCodecEvents(t, []byte{0x68, 0x7f}, PI(127))

	// [68 ff] = 255
	assertCodecEvents(t, []byte{0x68, 0xff}, PI(255))

	// [69 ff] = -255
	assertCodecEvents(t, []byte{0x69, 0xff}, NI(255))

	// [6c 80 96 98 00] = 10000000
	assertCodecEvents(t, []byte{0x6c, 0x80, 0x96, 0x98, 0x00}, PI(10000000))

	// [67 0f ff ee dd cc bb aa 99 88 77 66 55 44 33 22 11] = -0x112233445566778899aabbccddeeff
	assertCodecEvents(t, []byte{0x67, 0x0f, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11}, BI(NewBigInt("-0x112233445566778899aabbccddeeff")))
}

func TestSpecExamplesDecimalFloatingPoint(t *testing.T) {
	// [65 07 4b] = -7.5
	assertCodecEvents(t, []byte{0x65, 0x07, 0x4b}, DF(NewDFloat("-7.5")))

	// [65 ac 02 d0 9e 38] = 9.21424e+80
	assertCodecEvents(t, []byte{0x65, 0xac, 0x02, 0xd0, 0x9e, 0x38}, DF(NewDFloat("9.21424e+80")))
}

func TestSpecExamplesBinaryFloatingPoint(t *testing.T) {
	// [70 af 44] = 0x1.5ep+10
	assertCodecEvents(t, []byte{0x70, 0xaf, 0x44}, BF(0x1.5ep+10))

	// [71 00 e2 af 44] = 0x1.5fc4p+10
	assertCodecEvents(t, []byte{0x71, 0x00, 0xe2, 0xaf, 0x44}, BF(0x1.5fc4p+10))

	// [72 00 10 b4 3a 99 8f 32 46] = 0x1.28f993ab41p+100
	assertCodecEvents(t, []byte{0x72, 0x00, 0x10, 0xb4, 0x3a, 0x99, 0x8f, 0x32, 0x46}, BF(0x1.28f993ab41p+100))
}

func TestSpecExamplesUID(t *testing.T) {
	// [73 12 3e 45 67 e8 9b 12 d3 a4 56 42 66 55 44 00 00] = UID 123e4567-e89b-12d3-a456-426655440000
	assertCodecEvents(t, []byte{0x73, 0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44, 0x00, 0x00}, UID([]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x55, 0x44, 0x00, 0x00}))
}

func TestSpecExamplesTemporal(t *testing.T) {
	// [99 56 cd 00] = Oct 22, 2051
	assertCodecEvents(t, []byte{0x99, 0x56, 0xcd, 0x00}, CT(NewDate(2051, 10, 22)))

	// [9a f7 58 74 fc f6 a7 fd 10 45 2f 42 65 72 6c 69 6e] = 13:15:59.529435422/E/Berlin
	assertCodecEvents(t, []byte{0x9a, 0xf7, 0x58, 0x74, 0xfc, 0xf6, 0xa7, 0xfd, 0x10, 0x45, 0x2f, 0x42, 0x65, 0x72, 0x6c, 0x69, 0x6e}, CT(NewTime(13, 15, 59, 529435422, "E/Berlin")))

	// [9b 81 ac a0 b5 03 8f 1a ef d1] = Oct 26, 1985 1:22:16 at location 33.99, -117.93
	assertCodecEvents(t, []byte{0x9b, 0x81, 0xac, 0xa0, 0xb5, 0x03, 0x8f, 0x1a, 0xef, 0xd1}, CT(NewTSLL(1985, 10, 26, 1, 22, 16, 0, 3399, -11793)))
}

func TestSpecExamplesArray(t *testing.T) {
	// [83 61 62 63] = the string "abc" (short form array)
	assertCodecEvents(t, []byte{0x83, 0x61, 0x62, 0x63}, S("abc"))

	// [94 12 01 00 02 00] = unsigned 16-bit array with elements [1, 2] (long form array)
	assertCodecEvents(t, []byte{0x94, 0x12, 0x01, 0x00, 0x02, 0x00}, AU16([]uint16{1, 2}))

	// [90 21 m i s u n d e r s t a n d i n g 00 ...]
	assertCodecEvents(t, []byte{0x90, 0x21, 'm', 'i', 's', 'u', 'n', 'd', 'e', 'r', 's', 't', 'a', 'n', 'd', 'i', 'n', 'g', 0x00}, SB(), AC(16, true), AD([]byte("misunderstanding")), AC(0, false))

	// [8b 4d 61 69 6e 20 53 74 72 65 65 74] = Main Street
	assertCodecEvents(t, []byte{0x8b, 0x4d, 0x61, 0x69, 0x6e, 0x20, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74}, S("Main Street"))

	// [8d 52 c3 b6 64 65 6c 73 74 72 61 c3 9f 65] = Rödelstraße
	assertCodecEvents(t, []byte{0x8d, 0x52, 0xc3, 0xb6, 0x64, 0x65, 0x6c, 0x73, 0x74, 0x72, 0x61, 0xc3, 0x9f, 0x65}, S("Rödelstraße"))

	// [90 2a e8 a6 9a e7 8e 8b e5 b1 b1 e3 80 80 e6 97 a5 e6 b3 b0 e5 af ba] = 覚王山　日泰寺
	assertCodecEvents(t, []byte{0x90, 0x2a, 0xe8, 0xa6, 0x9a, 0xe7, 0x8e, 0x8b, 0xe5, 0xb1, 0xb1, 0xe3, 0x80, 0x80, 0xe6, 0x97, 0xa5, 0xe6, 0xb3, 0xb0, 0xe5, 0xaf, 0xba}, S("覚王山　日泰寺"))
}

func TestSpecExamplesIdentifier(t *testing.T) {
	// [97 07 73 6f 6d 65 5f 69 64 01] = value 1 marked with some_id
	assertCodecEvents(t, []byte{0x97, 0x07, 0x73, 0x6f, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x01}, MARK("some_id"), I(1))

	// [97 0f e7 99 bb e9 8c b2 e6 b8 88 e3 81 bf ef bc 95 01] = value 1 marked with 登録済み５
	assertCodecEvents(t, []byte{0x97, 0x0f, 0xe7, 0x99, 0xbb, 0xe9, 0x8c, 0xb2, 0xe6, 0xb8, 0x88, 0xe3, 0x81, 0xbf, 0xef, 0xbc, 0x95, 0x01}, MARK("登録済み５"), I(1))
}

func TestSpecExamplesResourceID(t *testing.T) {
	// [91 aa 01 68 74 74 70 73 3a 2f 2f 6a 6f 68 6e 2e 64 6f 65 40 77 77 77
	//  2e 65 78 61 6d 70 6c 65 2e 63 6f 6d 3a 31 32 33 2f 66 6f 72 75 6d 2f
	//  71 75 65 73 74 69 6f 6e 73 2f 3f 74 61 67 3d 6e 65 74 77 6f 72 6b 69
	//  6e 67 26 6f 72 64 65 72 3d 6e 65 77 65 73 74 23 74 6f 70]
	// = https://john.doe@www.example.com:123/forum/questions/?tag=networking&order=newest#top
	assertCodecEvents(t, []byte{0x91, 0xaa, 0x01, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a,
		0x2f, 0x2f, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x64, 0x6f, 0x65, 0x40, 0x77,
		0x77, 0x77, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63,
		0x6f, 0x6d, 0x3a, 0x31, 0x32, 0x33, 0x2f, 0x66, 0x6f, 0x72, 0x75, 0x6d,
		0x2f, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x3f,
		0x74, 0x61, 0x67, 0x3d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69,
		0x6e, 0x67, 0x26, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x3d, 0x6e, 0x65, 0x77,
		0x65, 0x73, 0x74, 0x23, 0x74, 0x6f, 0x70},
		RID("https://john.doe@www.example.com:123/forum/questions/?tag=networking&order=newest#top"))
}

func TestSpecExamplesCustomTypes(t *testing.T) {
	// [92 12 04 f6 28 3c 40 00 00 40 40]
	// = binary data representing an imaginary custom "cplx" struct
	//   {
	//       type:uint8 = 4
	//       real:float32 = 2.94
	//       imag:float32 = 3.0
	//   }
	assertCodecEvents(t, []byte{0x92, 0x12, 0x04, 0xf6, 0x28, 0x3c, 0x40, 0x00, 0x00, 0x40, 0x40},
		CUB([]byte{0x04, 0xf6, 0x28, 0x3c, 0x40, 0x00, 0x00, 0x40, 0x40}))

	// [93 1a 63 70 6c 78 28 32 2e 39 34 2b 33 69 29] = custom data encoded as the string "cplx(2.94+3i)"
	assertCodecEvents(t, []byte{0x93, 0x1a, 0x63, 0x70, 0x6c, 0x78, 0x28, 0x32, 0x2e,
		0x39, 0x34, 0x2b, 0x33, 0x69, 0x29}, CUT("cplx(2.94+3i)"))
}

func TestSpecExamplesTypedArray(t *testing.T) {
	// [95 04 01 02] = unsigned 8-bit array with elements 1, 2
	assertCodecEvents(t, []byte{0x95, 0x04, 0x01, 0x02}, AU8([]byte{1, 2}))

	// [94 12 01 00 02 00] = unsigned 16-bit array with elements 1, 2
	assertCodecEvents(t, []byte{0x94, 0x12, 0x01, 0x00, 0x02, 0x00}, AU16([]uint16{1, 2}))

	// [96 16 76 06] = bit array {0,1,1,0,1,1,1,0,0,1,1}
	assertCodecEvents(t, []byte{0x96, 0x16, 0x76, 0x06}, AB(11, []byte{0x76, 0x06}))
}

func TestSpecExamplesMedia(t *testing.T) {
	// [94 e1 20 61 70 70 6c 69 63 61 74 69 6f 6e 2f 78 2d 73 68 38 23 21 2f 62
	//  69 6e 2f 73 68 0a 0a 65 63 68 6f 20 68 65 6c 6c 6f 20 77 6f 72 6c 64 0a]
	assertCodecEvents(t,
		[]byte{0x94, 0xe1, 0x20, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
			0x69, 0x6f, 0x6e, 0x2f, 0x78, 0x2d, 0x73, 0x68, 0x38, 0x23, 0x21,
			0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73, 0x68, 0x0a, 0x0a, 0x65, 0x63,
			0x68, 0x6f, 0x20, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f,
			0x72, 0x6c, 0x64, 0x0a},
		MB(), AC(16, false), AD([]byte("application/x-sh")), AC(28, false),
		AD([]byte{0x23, 0x21, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73, 0x68, 0x0a,
			0x0a, 0x65, 0x63, 0x68, 0x6f, 0x20, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
			0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0a}))
}

func TestSpecExamplesContainers(t *testing.T) {
	// [7a 01 6a 88 13 7b] = A list containing integers (1, 5000)
	assertCodecEvents(t, []byte{0x7a, 0x01, 0x6a, 0x88, 0x13, 0x7b}, L(), I(1), PI(5000), E())

	// [79 81 61 01 81 62 02 7b] = A map containg the key-value pairs ("a", 1) ("b", 2)
	assertCodecEvents(t, []byte{0x79, 0x81, 0x61, 0x01, 0x81, 0x62, 0x02, 0x7b}, M(), S("a"), I(1), S("b"), I(2), E())

	// [76 91 24 68 74 74 70 3a 2f 2f 73 2e 67 6f 76 2f 68 6f 6d 65 72
	//  91 22 68 74 74 70 3a 2f 2f 65 2e 6f 72 67 2f 77 69 66 65
	//  91 24 68 74 74 70 3a 2f 2f 73 2e 67 6f 76 2f 6d 61 72 67 65]
	// = the relationship graph: @(@"http://s.gov/homer" @"http://e.org/wife" @"http://s.gov/marge")
	assertCodecEvents(t, []byte{0x76, 0x91, 0x24, 0x68, 0x74, 0x74, 0x70, 0x3a,
		0x2f, 0x2f, 0x73, 0x2e, 0x67, 0x6f, 0x76, 0x2f, 0x68, 0x6f, 0x6d, 0x65,
		0x72, 0x91, 0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x2e,
		0x6f, 0x72, 0x67, 0x2f, 0x77, 0x69, 0x66, 0x65, 0x91, 0x24, 0x68, 0x74,
		0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x73, 0x2e, 0x67, 0x6f, 0x76, 0x2f, 0x6d,
		0x61, 0x72, 0x67, 0x65},
		EDGE(), RID("http://s.gov/homer"), RID("http://e.org/wife"), RID("http://s.gov/marge"))

	// [77 01 77 03 77 05 7b 77 04 7b 7b 77 02 7b 7b]
	// = the binary tree:
	//   1
	//  / \
	// 2   3
	//    / \
	//   4   5
	assertCodecEvents(t, []byte{0x77, 0x01, 0x77, 0x03, 0x77, 0x05, 0x7b, 0x77, 0x04, 0x7b, 0x7b, 0x77, 0x02, 0x7b, 0x7b},
		NODE(), I(1), NODE(), I(3), NODE(), I(5), E(), NODE(), I(4), E(), E(), NODE(), I(2), E(), E())

	// [78 04 54 65 78 74 81 61 81 62 7b 89 53 6f 6d 65 20 74 65 78 74 7b] = <Text a=b,Some text>
	assertCodecEvents(t, []byte{0x78, 0x04, 0x54, 0x65, 0x78, 0x74, 0x81, 0x61, 0x81,
		0x62, 0x7b, 0x89, 0x53, 0x6f, 0x6d, 0x65, 0x20, 0x74, 0x65, 0x78, 0x74, 0x7b},
		MUP("Text"), S("a"), S("b"), E(), S("Some text"), E())
}

func TestSpecExamplesMarkerReference(t *testing.T) {
	// [97 01 61 79 8a 73 6f 6d 65 5f 76 61 6c 75 65 90 22 72
	//  65 70 65 61 74 20 74 68 69 73 20 76 61 6c 75 65 7b]
	// = the map {"some_value" = "repeat this value"}, tagged with the ID "a".
	assertCodecEvents(t, []byte{0x97, 0x01, 0x61, 0x79, 0x8a, 0x73, 0x6f, 0x6d, 0x65,
		0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x90, 0x22, 0x72, 0x65, 0x70, 0x65,
		0x61, 0x74, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x76, 0x61, 0x6c, 0x75,
		0x65, 0x7b},
		MARK("a"), M(), S("some_value"), S("repeat this value"), E())

	// [94 e2 24 63 6f 6d 6d 6f 6e 2e 63 65 23 6c 65 67 61 6c 65 73 65]
	// = reference to relative file "common.ce", ID "legalese" (common.ce#legalese)
	assertCodecEvents(t, []byte{0x94, 0xe0, 0x24, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
		0x2e, 0x63, 0x65, 0x23, 0x6c, 0x65, 0x67, 0x61, 0x6c, 0x65, 0x73, 0x65},
		RREF("common.ce#legalese"))

	// [94 e2 62 68 74 74 70 73 3a 2f 2f 73 6f 6d 65 77
	//  68 65 72 65 2e 63 6f 6d 2f 6d 79 5f 64 6f 63 75
	//  6d 65 6e 74 2e 63 62 65 3f 66 6f 72 6d 61 74 3d
	//  6c 6f 6e 67]
	// = reference to entire document at https://somewhere.com/my_document.cbe?format=long
	assertCodecEvents(t, []byte{0x94, 0xe0, 0x62, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a,
		0x2f, 0x2f, 0x73, 0x6f, 0x6d, 0x65, 0x77, 0x68, 0x65, 0x72, 0x65, 0x2e,
		0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x79, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d,
		0x65, 0x6e, 0x74, 0x2e, 0x63, 0x62, 0x65, 0x3f, 0x66, 0x6f, 0x72, 0x6d,
		0x61, 0x74, 0x3d, 0x6c, 0x6f, 0x6e, 0x67},
		RREF("https://somewhere.com/my_document.cbe?format=long"))
}

func TestSpecExamplesOtherPseudo(t *testing.T) {
	// Null is encoded as `[7e]`.
	assertCodecEvents(t, []byte{0x7e}, NULL())

	// [7f 7f 7f 6c 00 00 00 8f] = 0x8f000000, padded such that the 32-bit integer begins on a 4-byte boundary.
	assertCodecEvents(t, []byte{0x7f, 0x7f, 0x7f, 0x6c, 0x00, 0x00, 0x00, 0x8f},
		PAD(1), PAD(1), PAD(1), PI(0x8f000000))
}
