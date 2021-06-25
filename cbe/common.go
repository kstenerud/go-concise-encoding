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

const cbeDocumentHeader = byte(0x83)

type cbeTypeField uint16

const (
	cbeTypeDecimal      cbeTypeField = 0x65
	cbeTypePosInt                    = 0x66
	cbeTypeNegInt                    = 0x67
	cbeTypePosInt8                   = 0x68
	cbeTypeNegInt8                   = 0x69
	cbeTypePosInt16                  = 0x6a
	cbeTypeNegInt16                  = 0x6b
	cbeTypePosInt32                  = 0x6c
	cbeTypeNegInt32                  = 0x6d
	cbeTypePosInt64                  = 0x6e
	cbeTypeNegInt64                  = 0x6f
	cbeTypeFloat16                   = 0x70
	cbeTypeFloat32                   = 0x71
	cbeTypeFloat64                   = 0x72
	cbeTypeUID                       = 0x73
	cbeTypeReserved74                = 0x74
	cbeTypeReserved75                = 0x75
	cbeTypeRelationship              = 0x76
	cbeTypeComment                   = 0x77
	cbeTypeMarkup                    = 0x78
	cbeTypeMap                       = 0x79
	cbeTypeList                      = 0x7a
	cbeTypeEndContainer              = 0x7b
	cbeTypeFalse                     = 0x7c
	cbeTypeTrue                      = 0x7d
	cbeTypeNil                       = 0x7e
	cbeTypePadding                   = 0x7f
	cbeTypeString0                   = 0x80
	cbeTypeString1                   = 0x81
	cbeTypeString2                   = 0x82
	cbeTypeString3                   = 0x83
	cbeTypeString4                   = 0x84
	cbeTypeString5                   = 0x85
	cbeTypeString6                   = 0x86
	cbeTypeString7                   = 0x87
	cbeTypeString8                   = 0x88
	cbeTypeString9                   = 0x89
	cbeTypeString10                  = 0x8a
	cbeTypeString11                  = 0x8b
	cbeTypeString12                  = 0x8c
	cbeTypeString13                  = 0x8d
	cbeTypeString14                  = 0x8e
	cbeTypeString15                  = 0x8f
	cbeTypeString                    = 0x90
	cbeTypeRID                       = 0x91
	cbeTypeCustomBinary              = 0x92
	cbeTypeCustomText                = 0x93
	cbeTypePlane2                    = 0x94
	cbeTypeArrayUint8                = 0x95
	cbeTypeArrayBit                  = 0x96
	cbeTypeMarker                    = 0x97
	cbeTypeReference                 = 0x98
	cbeTypeDate                      = 0x99
	cbeTypeTime                      = 0x9a
	cbeTypeTimestamp                 = 0x9b

	// Plane 2 types
	cbeTypeNA           = 0xe0
	cbeTypeRIDCat       = 0xe1
	cbeTypeRIDReference = 0xe2
	cbeTypeMedia        = 0xe3
	cbeTypeArrayInt8    = 0xff
	cbeTypeArrayUint16  = 0xfe
	cbeTypeArrayInt16   = 0xfd
	cbeTypeArrayUint32  = 0xfc
	cbeTypeArrayInt32   = 0xfb
	cbeTypeArrayUint64  = 0xfa
	cbeTypeArrayInt64   = 0xf9
	cbeTypeArrayFloat16 = 0xf8
	cbeTypeArrayFloat32 = 0xf7
	cbeTypeArrayFloat64 = 0xf6
	cbeTypeArrayUID     = 0xf5

	cbeTypeShortArrayInt8    = 0x00
	cbeTypeShortArrayUint16  = 0x10
	cbeTypeShortArrayInt16   = 0x20
	cbeTypeShortArrayUint32  = 0x30
	cbeTypeShortArrayInt32   = 0x40
	cbeTypeShortArrayUint64  = 0x50
	cbeTypeShortArrayInt64   = 0x60
	cbeTypeShortArrayFloat16 = 0x70
	cbeTypeShortArrayFloat32 = 0x80
	cbeTypeShortArrayFloat64 = 0x90
	cbeTypeShortArrayUID     = 0xa0

	// Special code to mark EOF
	cbeTypeEOF = 0x100
)

const (
	cbeSmallIntMin int64 = -100
	cbeSmallIntMax int64 = 100
)

var isPlane2Array = []bool{
	events.ArrayTypeBit:              false,
	events.ArrayTypeUint8:            false,
	events.ArrayTypeUint16:           true,
	events.ArrayTypeUint32:           true,
	events.ArrayTypeUint64:           true,
	events.ArrayTypeInt8:             true,
	events.ArrayTypeInt16:            true,
	events.ArrayTypeInt32:            true,
	events.ArrayTypeInt64:            true,
	events.ArrayTypeFloat16:          true,
	events.ArrayTypeFloat32:          true,
	events.ArrayTypeFloat64:          true,
	events.ArrayTypeUID:              true,
	events.ArrayTypeString:           false,
	events.ArrayTypeResourceID:       false,
	events.ArrayTypeResourceIDConcat: true,
	events.ArrayTypeMedia:            true,
	events.ArrayTypeCustomBinary:     false,
	events.ArrayTypeCustomText:       false,
}

var arrayTypeToCBEType = []cbeTypeField{
	events.ArrayTypeBit:              cbeTypeArrayBit,
	events.ArrayTypeUint8:            cbeTypeArrayUint8,
	events.ArrayTypeUint16:           cbeTypeArrayUint16,
	events.ArrayTypeUint32:           cbeTypeArrayUint32,
	events.ArrayTypeUint64:           cbeTypeArrayUint64,
	events.ArrayTypeInt8:             cbeTypeArrayInt8,
	events.ArrayTypeInt16:            cbeTypeArrayInt16,
	events.ArrayTypeInt32:            cbeTypeArrayInt32,
	events.ArrayTypeInt64:            cbeTypeArrayInt64,
	events.ArrayTypeFloat16:          cbeTypeArrayFloat16,
	events.ArrayTypeFloat32:          cbeTypeArrayFloat32,
	events.ArrayTypeFloat64:          cbeTypeArrayFloat64,
	events.ArrayTypeUID:              cbeTypeArrayUID,
	events.ArrayTypeString:           cbeTypeString,
	events.ArrayTypeResourceID:       cbeTypeRID,
	events.ArrayTypeCustomBinary:     cbeTypeCustomBinary,
	events.ArrayTypeCustomText:       cbeTypeCustomText,
	events.ArrayTypeResourceIDConcat: cbeTypeRIDCat,
	events.ArrayTypeMedia:            cbeTypeMedia,
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
	cbeTypeRIDCat:       events.ArrayTypeResourceIDConcat,
	cbeTypeMedia:        events.ArrayTypeMedia,
}
