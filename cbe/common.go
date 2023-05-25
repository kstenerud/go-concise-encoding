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

import "github.com/kstenerud/go-concise-encoding/ce/events"

// ============================================================================

// Internal

const CBESignatureByte = byte(0x81)

type cbeTypeField uint16

const (
	cbeTypeUID            cbeTypeField = 0x65
	cbeTypePosInt         cbeTypeField = 0x66
	cbeTypeNegInt         cbeTypeField = 0x67
	cbeTypePosInt8        cbeTypeField = 0x68
	cbeTypeNegInt8        cbeTypeField = 0x69
	cbeTypePosInt16       cbeTypeField = 0x6a
	cbeTypeNegInt16       cbeTypeField = 0x6b
	cbeTypePosInt32       cbeTypeField = 0x6c
	cbeTypeNegInt32       cbeTypeField = 0x6d
	cbeTypePosInt64       cbeTypeField = 0x6e
	cbeTypeNegInt64       cbeTypeField = 0x6f
	cbeTypeFloat16        cbeTypeField = 0x70
	cbeTypeFloat32        cbeTypeField = 0x71
	cbeTypeFloat64        cbeTypeField = 0x72
	cbeTypeReserved73     cbeTypeField = 0x73 //nolint
	cbeTypeReserved74     cbeTypeField = 0x74 //nolint
	cbeTypeReserved75     cbeTypeField = 0x75 //nolint
	cbeTypeDecimal        cbeTypeField = 0x76
	cbeTypeLocalReference cbeTypeField = 0x77
	cbeTypeFalse          cbeTypeField = 0x78
	cbeTypeTrue           cbeTypeField = 0x79
	cbeTypeDate           cbeTypeField = 0x7a
	cbeTypeTime           cbeTypeField = 0x7b
	cbeTypeTimestamp      cbeTypeField = 0x7c
	cbeTypeNull           cbeTypeField = 0x7d
	cbeTypeReserved7e     cbeTypeField = 0x7e //nolint
	cbeTypePlane7f        cbeTypeField = 0x7f
	cbeTypeString0        cbeTypeField = 0x80
	cbeTypeString1        cbeTypeField = 0x81
	cbeTypeString2        cbeTypeField = 0x82
	cbeTypeString3        cbeTypeField = 0x83
	cbeTypeString4        cbeTypeField = 0x84
	cbeTypeString5        cbeTypeField = 0x85
	cbeTypeString6        cbeTypeField = 0x86
	cbeTypeString7        cbeTypeField = 0x87
	cbeTypeString8        cbeTypeField = 0x88
	cbeTypeString9        cbeTypeField = 0x89
	cbeTypeString10       cbeTypeField = 0x8a
	cbeTypeString11       cbeTypeField = 0x8b
	cbeTypeString12       cbeTypeField = 0x8c
	cbeTypeString13       cbeTypeField = 0x8d
	cbeTypeString14       cbeTypeField = 0x8e
	cbeTypeString15       cbeTypeField = 0x8f
	cbeTypeString         cbeTypeField = 0x90
	cbeTypeRID            cbeTypeField = 0x91
	cbeTypeCustomType     cbeTypeField = 0x92
	cbeTypeArrayUint8     cbeTypeField = 0x93
	cbeTypeArrayBit       cbeTypeField = 0x94
	cbeTypePadding        cbeTypeField = 0x95
	cbeTypeRecord         cbeTypeField = 0x96
	cbeTypeEdge           cbeTypeField = 0x97
	cbeTypeNode           cbeTypeField = 0x98
	cbeTypeMap            cbeTypeField = 0x99
	cbeTypeList           cbeTypeField = 0x9a
	cbeTypeEndContainer   cbeTypeField = 0x9b

	// Plane 2 types
	cbeTypeShortArrayUID     cbeTypeField = 0x00
	cbeTypeShortArrayInt8    cbeTypeField = 0x10
	cbeTypeShortArrayUint16  cbeTypeField = 0x20
	cbeTypeShortArrayInt16   cbeTypeField = 0x30
	cbeTypeShortArrayUint32  cbeTypeField = 0x40
	cbeTypeShortArrayInt32   cbeTypeField = 0x50
	cbeTypeShortArrayUint64  cbeTypeField = 0x60
	cbeTypeShortArrayInt64   cbeTypeField = 0x70
	cbeTypeShortArrayFloat16 cbeTypeField = 0x80
	cbeTypeShortArrayFloat32 cbeTypeField = 0x90
	cbeTypeShortArrayFloat64 cbeTypeField = 0xa0
	cbeTypeArrayUID          cbeTypeField = 0xe0
	cbeTypeArrayInt8         cbeTypeField = 0xe1
	cbeTypeArrayUint16       cbeTypeField = 0xe2
	cbeTypeArrayInt16        cbeTypeField = 0xe3
	cbeTypeArrayUint32       cbeTypeField = 0xe4
	cbeTypeArrayInt32        cbeTypeField = 0xe5
	cbeTypeArrayUint64       cbeTypeField = 0xe6
	cbeTypeArrayInt64        cbeTypeField = 0xe7
	cbeTypeArrayFloat16      cbeTypeField = 0xe8
	cbeTypeArrayFloat32      cbeTypeField = 0xe9
	cbeTypeArrayFloat64      cbeTypeField = 0xea
	cbeTypeMarker            cbeTypeField = 0xf0
	cbeTypeRecordType        cbeTypeField = 0xf1
	cbeTypeRemoteReference   cbeTypeField = 0xf2
	cbeTypeMedia             cbeTypeField = 0xf3

	// Special code to mark EOF
	cbeTypeEOF cbeTypeField = 0x100
)

const (
	cbeSmallIntMin int64 = -100
	cbeSmallIntMax int64 = 100
)

var isPlane7fArray = []bool{
	events.ArrayTypeBit:             false,
	events.ArrayTypeUint8:           false,
	events.ArrayTypeUint16:          true,
	events.ArrayTypeUint32:          true,
	events.ArrayTypeUint64:          true,
	events.ArrayTypeInt8:            true,
	events.ArrayTypeInt16:           true,
	events.ArrayTypeInt32:           true,
	events.ArrayTypeInt64:           true,
	events.ArrayTypeFloat16:         true,
	events.ArrayTypeFloat32:         true,
	events.ArrayTypeFloat64:         true,
	events.ArrayTypeUID:             true,
	events.ArrayTypeString:          false,
	events.ArrayTypeResourceID:      false,
	events.ArrayTypeReferenceRemote: true,
	events.ArrayTypeMedia:           true,
	events.ArrayTypeCustomBinary:    false,
	events.ArrayTypeCustomText:      false,
}

var arrayTypeToCBEType = []cbeTypeField{
	events.ArrayTypeBit:             cbeTypeArrayBit,
	events.ArrayTypeUint8:           cbeTypeArrayUint8,
	events.ArrayTypeUint16:          cbeTypeArrayUint16,
	events.ArrayTypeUint32:          cbeTypeArrayUint32,
	events.ArrayTypeUint64:          cbeTypeArrayUint64,
	events.ArrayTypeInt8:            cbeTypeArrayInt8,
	events.ArrayTypeInt16:           cbeTypeArrayInt16,
	events.ArrayTypeInt32:           cbeTypeArrayInt32,
	events.ArrayTypeInt64:           cbeTypeArrayInt64,
	events.ArrayTypeFloat16:         cbeTypeArrayFloat16,
	events.ArrayTypeFloat32:         cbeTypeArrayFloat32,
	events.ArrayTypeFloat64:         cbeTypeArrayFloat64,
	events.ArrayTypeUID:             cbeTypeArrayUID,
	events.ArrayTypeString:          cbeTypeString,
	events.ArrayTypeResourceID:      cbeTypeRID,
	events.ArrayTypeCustomBinary:    cbeTypeCustomType,
	events.ArrayTypeCustomText:      cbeTypeCustomType,
	events.ArrayTypeReferenceRemote: cbeTypeRemoteReference,
	events.ArrayTypeMedia:           cbeTypeMedia,
}

var cbePlane7fTypeToArrayType = [256]events.ArrayType{
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
