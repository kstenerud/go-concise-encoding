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

// Performs encoding and decoding of Concise Binary Encoding documents
// (https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md).
//
// The decoder decodes a document to produce data events, and the encoder
// consumes data events to produce a document.
package cbe

import (
	"github.com/kstenerud/go-concise-encoding/events"
)

// ============================================================================

// Internal

const CBESignatureByte = byte(0x81)

type cbeTypeField uint16

const (
	cbeTypeDecimal      cbeTypeField = 0x65
	cbeTypePosInt       cbeTypeField = 0x66
	cbeTypeNegInt       cbeTypeField = 0x67
	cbeTypePosInt8      cbeTypeField = 0x68
	cbeTypeNegInt8      cbeTypeField = 0x69
	cbeTypePosInt16     cbeTypeField = 0x6a
	cbeTypeNegInt16     cbeTypeField = 0x6b
	cbeTypePosInt32     cbeTypeField = 0x6c
	cbeTypeNegInt32     cbeTypeField = 0x6d
	cbeTypePosInt64     cbeTypeField = 0x6e
	cbeTypeNegInt64     cbeTypeField = 0x6f
	cbeTypeFloat16      cbeTypeField = 0x70
	cbeTypeFloat32      cbeTypeField = 0x71
	cbeTypeFloat64      cbeTypeField = 0x72
	cbeTypeUID          cbeTypeField = 0x73
	cbeTypeReserved74   cbeTypeField = 0x74 //nolint
	cbeTypeReserved75   cbeTypeField = 0x75 //nolint
	cbeTypeReserved76   cbeTypeField = 0x76
	cbeTypeEdge         cbeTypeField = 0x77
	cbeTypeNode         cbeTypeField = 0x78
	cbeTypeMap          cbeTypeField = 0x79
	cbeTypeList         cbeTypeField = 0x7a
	cbeTypeEndContainer cbeTypeField = 0x7b
	cbeTypeFalse        cbeTypeField = 0x7c
	cbeTypeTrue         cbeTypeField = 0x7d
	cbeTypeNull         cbeTypeField = 0x7e
	cbeTypePadding      cbeTypeField = 0x7f
	cbeTypeString0      cbeTypeField = 0x80
	cbeTypeString1      cbeTypeField = 0x81
	cbeTypeString2      cbeTypeField = 0x82
	cbeTypeString3      cbeTypeField = 0x83
	cbeTypeString4      cbeTypeField = 0x84
	cbeTypeString5      cbeTypeField = 0x85
	cbeTypeString6      cbeTypeField = 0x86
	cbeTypeString7      cbeTypeField = 0x87
	cbeTypeString8      cbeTypeField = 0x88
	cbeTypeString9      cbeTypeField = 0x89
	cbeTypeString10     cbeTypeField = 0x8a
	cbeTypeString11     cbeTypeField = 0x8b
	cbeTypeString12     cbeTypeField = 0x8c
	cbeTypeString13     cbeTypeField = 0x8d
	cbeTypeString14     cbeTypeField = 0x8e
	cbeTypeString15     cbeTypeField = 0x8f
	cbeTypeString       cbeTypeField = 0x90
	cbeTypeRID          cbeTypeField = 0x91
	cbeTypeCustomType   cbeTypeField = 0x92
	cbeTypeReserved93   cbeTypeField = 0x93 //nolint
	cbeTypePlane2       cbeTypeField = 0x94
	cbeTypeArrayUint8   cbeTypeField = 0x95
	cbeTypeArrayBit     cbeTypeField = 0x96
	cbeTypeMarker       cbeTypeField = 0x97
	cbeTypeReference    cbeTypeField = 0x98
	cbeTypeDate         cbeTypeField = 0x99
	cbeTypeTime         cbeTypeField = 0x9a
	cbeTypeTimestamp    cbeTypeField = 0x9b

	// Plane 2 types
	cbeTypeShortArrayInt8    cbeTypeField = 0x00
	cbeTypeShortArrayUint16  cbeTypeField = 0x10
	cbeTypeShortArrayInt16   cbeTypeField = 0x20
	cbeTypeShortArrayUint32  cbeTypeField = 0x30
	cbeTypeShortArrayInt32   cbeTypeField = 0x40
	cbeTypeShortArrayUint64  cbeTypeField = 0x50
	cbeTypeShortArrayInt64   cbeTypeField = 0x60
	cbeTypeShortArrayFloat16 cbeTypeField = 0x70
	cbeTypeShortArrayFloat32 cbeTypeField = 0x80
	cbeTypeShortArrayFloat64 cbeTypeField = 0x90
	cbeTypeShortArrayUID     cbeTypeField = 0xa0

	cbeTypeRemoteReference cbeTypeField = 0xe0
	cbeTypeMedia           cbeTypeField = 0xe1

	cbeTypeArrayUID     cbeTypeField = 0xf5
	cbeTypeArrayFloat64 cbeTypeField = 0xf6
	cbeTypeArrayFloat32 cbeTypeField = 0xf7
	cbeTypeArrayFloat16 cbeTypeField = 0xf8
	cbeTypeArrayInt64   cbeTypeField = 0xf9
	cbeTypeArrayUint64  cbeTypeField = 0xfa
	cbeTypeArrayInt32   cbeTypeField = 0xfb
	cbeTypeArrayUint32  cbeTypeField = 0xfc
	cbeTypeArrayInt16   cbeTypeField = 0xfd
	cbeTypeArrayUint16  cbeTypeField = 0xfe
	cbeTypeArrayInt8    cbeTypeField = 0xff

	// Special code to mark EOF
	cbeTypeEOF cbeTypeField = 0x100
)

const (
	cbeSmallIntMin int64 = -100
	cbeSmallIntMax int64 = 100
)

var isPlane2Array = []bool{
	events.ArrayTypeBit:          false,
	events.ArrayTypeUint8:        false,
	events.ArrayTypeUint16:       true,
	events.ArrayTypeUint32:       true,
	events.ArrayTypeUint64:       true,
	events.ArrayTypeInt8:         true,
	events.ArrayTypeInt16:        true,
	events.ArrayTypeInt32:        true,
	events.ArrayTypeInt64:        true,
	events.ArrayTypeFloat16:      true,
	events.ArrayTypeFloat32:      true,
	events.ArrayTypeFloat64:      true,
	events.ArrayTypeUID:          true,
	events.ArrayTypeString:       false,
	events.ArrayTypeResourceID:   false,
	events.ArrayTypeRemoteRef:    true,
	events.ArrayTypeMedia:        true,
	events.ArrayTypeCustomBinary: false,
	events.ArrayTypeCustomText:   false,
}

var arrayTypeToCBEType = []cbeTypeField{
	events.ArrayTypeBit:          cbeTypeArrayBit,
	events.ArrayTypeUint8:        cbeTypeArrayUint8,
	events.ArrayTypeUint16:       cbeTypeArrayUint16,
	events.ArrayTypeUint32:       cbeTypeArrayUint32,
	events.ArrayTypeUint64:       cbeTypeArrayUint64,
	events.ArrayTypeInt8:         cbeTypeArrayInt8,
	events.ArrayTypeInt16:        cbeTypeArrayInt16,
	events.ArrayTypeInt32:        cbeTypeArrayInt32,
	events.ArrayTypeInt64:        cbeTypeArrayInt64,
	events.ArrayTypeFloat16:      cbeTypeArrayFloat16,
	events.ArrayTypeFloat32:      cbeTypeArrayFloat32,
	events.ArrayTypeFloat64:      cbeTypeArrayFloat64,
	events.ArrayTypeUID:          cbeTypeArrayUID,
	events.ArrayTypeString:       cbeTypeString,
	events.ArrayTypeResourceID:   cbeTypeRID,
	events.ArrayTypeCustomBinary: cbeTypeCustomType,
	events.ArrayTypeCustomText:   cbeTypeCustomType,
	events.ArrayTypeRemoteRef:    cbeTypeRemoteReference,
	events.ArrayTypeMedia:        cbeTypeMedia,
}

var cbePlane2TypeToArrayType = [256]events.ArrayType{
	cbeTypeArrayBit:     events.ArrayTypeBit,
	cbeTypeArrayUint8:   events.ArrayTypeUint8,
	cbeTypeArrayUint16:  events.ArrayTypeUint16,
	cbeTypeArrayUint32:  events.ArrayTypeUint32,
	cbeTypeArrayUint64:  events.ArrayTypeUint64,
	cbeTypeArrayInt8:    events.ArrayTypeInt8,
	cbeTypeArrayInt16:   events.ArrayTypeInt16,
	cbeTypeArrayInt32:   events.ArrayTypeInt32,
	cbeTypeArrayInt64:   events.ArrayTypeInt64,
	cbeTypeArrayFloat16: events.ArrayTypeFloat16,
	cbeTypeArrayFloat32: events.ArrayTypeFloat32,
	cbeTypeArrayFloat64: events.ArrayTypeFloat64,
	cbeTypeArrayUID:     events.ArrayTypeUID,
	cbeTypeMedia:        events.ArrayTypeMedia,
}
