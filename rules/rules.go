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

// Imposes the structural rules that enforce a well-formed concise encoding
// document.
package rules

import (
	"fmt"

	"github.com/kstenerud/go-concise-encoding/ce/events"
)

const maxMarkerIDRuneCount = 50
const maxMarkerIDByteCount = 4 * maxMarkerIDRuneCount // max 4 bytes per rune

type negint uint64
type rid string

type EventRule interface {
	OnBeginDocument(ctx *Context)
	OnEndDocument(ctx *Context)
	OnChildContainerEnded(ctx *Context, containerType DataType)
	OnVersion(ctx *Context, version uint64)
	OnPadding(ctx *Context)
	OnComment(ctx *Context)
	OnKeyableObject(ctx *Context, objType DataType, key interface{})
	OnNonKeyableObject(ctx *Context, objType DataType)
	OnNull(ctx *Context)
	OnList(ctx *Context)
	OnMap(ctx *Context)
	OnRecordType(ctx *Context, identifier []byte)
	OnRecord(ctx *Context, identifier []byte)
	OnEdge(ctx *Context)
	OnNode(ctx *Context)
	OnEnd(ctx *Context)
	OnMarker(ctx *Context, identifier []byte)
	OnReferenceLocal(ctx *Context, identifier []byte)
	OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)
	OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string)
	OnArrayBegin(ctx *Context, arrayType events.ArrayType)
	OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)
	OnArrayData(ctx *Context, data []byte)
}

var (
	beginDocumentRule       BeginDocumentRule
	endDocumentRule         EndDocumentRule
	terminalRule            TerminalRule
	versionRule             VersionRule
	topLevelRule            TopLevelRule
	listRule                ListRule
	mapKeyRule              MapKeyRule
	mapValueRule            MapValueRule
	recordTypeRule          RecordTypeRule
	recordRule              RecordRule
	arrayRule               ArrayRule
	arrayChunkRule          ArrayChunkRule
	stringRule              StringRule
	stringChunkRule         StringChunkRule
	markedObjectKeyableRule MarkedObjectKeyableRule
	markedObjectAnyTypeRule MarkedObjectAnyTypeRule
	stringBuilderRule       StringBuilderRule
	stringBuilderChunkRule  StringBuilderChunkRule
	edgeSourceRule          EdgeSourceRule
	edgeDescriptionRule     EdgeDescriptionRule
	edgeDestinationRule     EdgeDestinationRule
	nodeRule                NodeRule
	awaitEndRule            AwaitEndRule
)

var arrayTypeToDataType = []DataType{
	events.ArrayTypeInvalid:         DataTypeInvalid,
	events.ArrayTypeString:          DataTypeString,
	events.ArrayTypeResourceID:      DataTypeResourceID,
	events.ArrayTypeReferenceRemote: DataTypeRemoteReference,
	events.ArrayTypeCustomText:      DataTypeCustomText,
	events.ArrayTypeCustomBinary:    DataTypeCustomBinary,
	events.ArrayTypeBit:             DataTypeArrayBit,
	events.ArrayTypeUint8:           DataTypeArrayUint8,
	events.ArrayTypeUint16:          DataTypeArrayUint16,
	events.ArrayTypeUint32:          DataTypeArrayUint32,
	events.ArrayTypeUint64:          DataTypeArrayUint64,
	events.ArrayTypeInt8:            DataTypeArrayInt8,
	events.ArrayTypeInt16:           DataTypeArrayInt16,
	events.ArrayTypeInt32:           DataTypeArrayInt32,
	events.ArrayTypeInt64:           DataTypeArrayInt64,
	events.ArrayTypeFloat16:         DataTypeArrayFloat16,
	events.ArrayTypeFloat32:         DataTypeArrayFloat32,
	events.ArrayTypeFloat64:         DataTypeArrayFloat64,
	events.ArrayTypeUID:             DataTypeArrayUID,
	events.ArrayTypeMedia:           DataTypeMedia,
}

func wrongType(context interface{}, dataType interface{}) {
	panic(fmt.Errorf("%v is not allowed while processing %v", dataType, context))
}
