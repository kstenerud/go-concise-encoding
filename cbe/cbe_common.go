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

const cbeDocumentHeader = byte(0x03)

type cbeTypeField uint8

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
	cbeTypeUUID         cbeTypeField = 0x73
	cbeTypeReserved74   cbeTypeField = 0x74
	cbeTypeReserved75   cbeTypeField = 0x75
	cbeTypeComment      cbeTypeField = 0x76
	cbeTypeMetadata     cbeTypeField = 0x77
	cbeTypeMarkup       cbeTypeField = 0x78
	cbeTypeMap          cbeTypeField = 0x79
	cbeTypeList         cbeTypeField = 0x7a
	cbeTypeEndContainer cbeTypeField = 0x7b
	cbeTypeFalse        cbeTypeField = 0x7c
	cbeTypeTrue         cbeTypeField = 0x7d
	cbeTypeNA           cbeTypeField = 0x7e
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
	cbeTypeCustomBinary cbeTypeField = 0x92
	cbeTypeCustomText   cbeTypeField = 0x93
	cbeTypePlane2       cbeTypeField = 0x94
	cbeTypeReserved95   cbeTypeField = 0x95
	cbeTypeReserved96   cbeTypeField = 0x96
	cbeTypeMarker       cbeTypeField = 0x97
	cbeTypeReference    cbeTypeField = 0x98
	cbeTypeDate         cbeTypeField = 0x99
	cbeTypeTime         cbeTypeField = 0x9a
	cbeTypeTimestamp    cbeTypeField = 0x9b

	// Plane 2 types
	cbeTypeArrayBoolean = cbeTypeTrue
	cbeTypeArrayUint8   = cbeTypePosInt8
	cbeTypeArrayUint16  = cbeTypePosInt16
	cbeTypeArrayUint32  = cbeTypePosInt32
	cbeTypeArrayUint64  = cbeTypePosInt64
	cbeTypeArrayInt8    = cbeTypeNegInt8
	cbeTypeArrayInt16   = cbeTypeNegInt16
	cbeTypeArrayInt32   = cbeTypeNegInt32
	cbeTypeArrayInt64   = cbeTypeNegInt64
	cbeTypeArrayFloat16 = cbeTypeFloat16
	cbeTypeArrayFloat32 = cbeTypeFloat32
	cbeTypeArrayFloat64 = cbeTypeFloat64
	cbeTypeArrayUUID    = cbeTypeUUID
	cbeTypeRIDCat       = cbeTypeRID
	cbeTypeNACat        = cbeTypeNA
)

const (
	cbeSmallIntMin int64 = -100
	cbeSmallIntMax int64 = 100
)

var isPlane2Array = []bool{
	events.ArrayTypeBoolean:          true,
	events.ArrayTypeUint8:            true,
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
	events.ArrayTypeUUID:             true,
	events.ArrayTypeString:           false,
	events.ArrayTypeResourceID:       false,
	events.ArrayTypeResourceIDConcat: true,
	events.ArrayTypeCustomBinary:     false,
	events.ArrayTypeCustomText:       false,
}

var arrayTypeToCBEType = []cbeTypeField{
	events.ArrayTypeBoolean:          cbeTypeTrue,
	events.ArrayTypeUint8:            cbeTypePosInt8,
	events.ArrayTypeUint16:           cbeTypePosInt16,
	events.ArrayTypeUint32:           cbeTypePosInt32,
	events.ArrayTypeUint64:           cbeTypePosInt64,
	events.ArrayTypeInt8:             cbeTypeNegInt8,
	events.ArrayTypeInt16:            cbeTypeNegInt16,
	events.ArrayTypeInt32:            cbeTypeNegInt32,
	events.ArrayTypeInt64:            cbeTypeNegInt64,
	events.ArrayTypeFloat16:          cbeTypeFloat16,
	events.ArrayTypeFloat32:          cbeTypeFloat32,
	events.ArrayTypeFloat64:          cbeTypeFloat64,
	events.ArrayTypeUUID:             cbeTypeUUID,
	events.ArrayTypeString:           cbeTypeString,
	events.ArrayTypeResourceID:       cbeTypeRID,
	events.ArrayTypeResourceIDConcat: cbeTypeRIDCat,
	events.ArrayTypeCustomBinary:     cbeTypeCustomBinary,
	events.ArrayTypeCustomText:       cbeTypeCustomText,
}

var cbePlane2TypeToArrayType = [256]events.ArrayType{
	cbeTypeArrayBoolean: events.ArrayTypeBoolean,
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
	cbeTypeArrayUUID:    events.ArrayTypeUUID,
	cbeTypeRIDCat:       events.ArrayTypeResourceIDConcat,
}
