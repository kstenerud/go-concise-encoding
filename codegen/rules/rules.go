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

package rules

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
	"github.com/kstenerud/go-concise-encoding/codegen/datatypes"
)

const path = "rules"

var imports = []*common.Import{
	{As: "", Import: "fmt"},
	{As: "", Import: "strings"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/ce/events"},
}

func GenerateCode(projectDir string) {
	common.GenerateGoFile(filepath.Join(projectDir, path), imports, func(writer io.Writer) {
		generateDataTypeType(writer)
		generateDefaultMethods(writer)
		generateBadEventMethods(writer)
	})
}

// ---------------
// Code Generators
// ---------------

func generateDataTypeType(writer io.Writer) {
	gen := datatypes.NewFlagDataTypeWriter(writer, "DataType", EndDataTypes)

	gen.BeginType()
	for i := DataType(1); i < EndDataTypes; i <<= 1 {
		gen.AddNamed(i)
	}
	gen.AddCustom(DataTypeInvalid, uint64(DataTypeInvalid))
	gen.AddCustom("AllowAny", DataTypesAll)
	gen.AddCustom("AllowNonNull", DataTypesNonNull)
	gen.AddCustom("AllowKeyable", DataTypesKeyable)
	gen.AddCustom("AllowMarkable", DataTypesMarkable)
	gen.AddCustom("AllowString", DataTypeString)
	gen.AddCustom("AllowResourceID", DataTypeResourceID)
	gen.EndType()

	gen.BeginStringer()
	for i := DataType(1); i < EndDataTypes; i <<= 1 {
		gen.AddStringer(i)
	}
	gen.EndStringer()
}

func generateBadEventMethods(writer io.Writer) {
	for _, rule := range allRules {
		for _, method := range allMethods {
			if !rule.hasMethod(method) {
				generateBadEventMethod(rule, method, writer)
			}
		}
	}
}

func generateBadEventMethod(rule Rule, method *Method, writer io.Writer) {
	openMethod(rule, method, writer)
	str := ""
	switch method.MethodType {
	case MethodTypeScalar:
		str = fmt.Sprintf("\twrongType(\"%v\", objType)\n", rule.FriendlyName)
	default:
		str = fmt.Sprintf("\twrongType(\"%v\", \"%v\")\n", rule.FriendlyName, method.Name)
	}

	if _, err := writer.Write([]byte(str)); err != nil {
		panic(err)
	}
	closeMethod(writer)
}

func generateDefaultMethods(writer io.Writer) {
	for _, rule := range allRules {
		for _, method := range allMethods {
			if rule.hasDefaultMethod(method) {
				generateDefaultMethod(rule, method, writer)
			}
		}
	}
}

func generateDefaultMethod(rule Rule, method *Method, writer io.Writer) {
	openMethod(rule, method, writer)
	// Just cut off the initial LF rather than complicate the openMethod function.
	if _, err := writer.Write([]byte(method.DefaultImplementation[1:])); err != nil {
		panic(err)
	}
	closeMethod(writer)
}

// -------
// Utility
// -------

func openMethod(rule Rule, method *Method, writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("func (_this *%s) %s {\n", rule.Name, method.Signature))); err != nil {
		panic(err)
	}
}

func closeMethod(writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("}\n"))); err != nil {
		panic(err)
	}
}

// ----
// Data
// ----

type DataType uint64

const (
	DataTypeNull DataType = 1 << iota
	DataTypeNan
	DataTypeBool
	DataTypeInt
	DataTypeFloat
	DataTypeUID
	DataTypeTime
	DataTypeList
	DataTypeMap
	DataTypeRecordType
	DataTypeRecord
	DataTypeEdge
	DataTypeNode
	DataTypeString
	DataTypeMedia
	DataTypeArrayBit
	DataTypeArrayUint8
	DataTypeArrayUint16
	DataTypeArrayUint32
	DataTypeArrayUint64
	DataTypeArrayInt8
	DataTypeArrayInt16
	DataTypeArrayInt32
	DataTypeArrayInt64
	DataTypeArrayFloat16
	DataTypeArrayFloat32
	DataTypeArrayFloat64
	DataTypeArrayUID
	DataTypeCustomText
	DataTypeCustomBinary
	DataTypeMarker
	DataTypeLocalReference
	DataTypeResourceID
	DataTypeRemoteReference
	DataTypeComment
	DataTypePadding
	EndDataTypes

	DataTypeInvalid  = DataType(0)
	DataTypesNone    = DataType(0)
	DataTypesAll     = ^DataTypesNone
	DataTypesNonNull = ^DataTypeNull
	DataTypesKeyable = DataTypeBool |
		DataTypeInt |
		DataTypeFloat |
		DataTypeUID |
		DataTypeTime |
		DataTypeString |
		DataTypeResourceID |
		DataTypeMarker |
		DataTypeLocalReference |
		DataTypePadding |
		DataTypeComment
	DataTypesNonKeyable = ^DataTypesKeyable
	DataTypesMarkable   = ^(DataTypeMarker | DataTypeLocalReference | DataTypeRemoteReference | DataTypeComment | DataTypeRecordType)
	DataTypesTopLevel   = ^(DataTypeLocalReference)
	DataTypesContainer  = DataTypeList | DataTypeMap | DataTypeEdge | DataTypeNode
	DataTypesStringlike = DataTypeString |
		DataTypeResourceID |
		DataTypeRemoteReference |
		DataTypeCustomText
	DataTypesBinary = DataTypeArrayBit |
		DataTypeArrayUint8 |
		DataTypeArrayUint16 |
		DataTypeArrayUint32 |
		DataTypeArrayUint64 |
		DataTypeArrayInt8 |
		DataTypeArrayInt16 |
		DataTypeArrayInt32 |
		DataTypeArrayInt64 |
		DataTypeArrayFloat16 |
		DataTypeArrayFloat32 |
		DataTypeArrayFloat64 |
		DataTypeArrayUID |
		DataTypeCustomBinary
	DataTypesAllArrays = DataTypesStringlike | DataTypesBinary
)

func (_this DataType) String() string {
	return datatypes.FlagToString(dataTypeNames, _this)
}

var dataTypeNames = map[interface{}]string{
	DataTypeInvalid:         "DataTypeInvalid",
	DataTypeNull:            "DataTypeNull",
	DataTypeNan:             "DataTypeNan",
	DataTypeBool:            "DataTypeBool",
	DataTypeInt:             "DataTypeInt",
	DataTypeFloat:           "DataTypeFloat",
	DataTypeUID:             "DataTypeUID",
	DataTypeTime:            "DataTypeTime",
	DataTypeList:            "DataTypeList",
	DataTypeMap:             "DataTypeMap",
	DataTypeRecordType:      "DataTypeRecordType",
	DataTypeRecord:          "DataTypeRecord",
	DataTypeEdge:            "DataTypeEdge",
	DataTypeNode:            "DataTypeNode",
	DataTypeString:          "DataTypeString",
	DataTypeMedia:           "DataTypeMedia",
	DataTypeArrayBit:        "DataTypeArrayBit",
	DataTypeArrayUint8:      "DataTypeArrayUint8",
	DataTypeArrayUint16:     "DataTypeArrayUint16",
	DataTypeArrayUint32:     "DataTypeArrayUint32",
	DataTypeArrayUint64:     "DataTypeArrayUint64",
	DataTypeArrayInt8:       "DataTypeArrayInt8",
	DataTypeArrayInt16:      "DataTypeArrayInt16",
	DataTypeArrayInt32:      "DataTypeArrayInt32",
	DataTypeArrayInt64:      "DataTypeArrayInt64",
	DataTypeArrayFloat16:    "DataTypeArrayFloat16",
	DataTypeArrayFloat32:    "DataTypeArrayFloat32",
	DataTypeArrayFloat64:    "DataTypeArrayFloat64",
	DataTypeArrayUID:        "DataTypeArrayUID",
	DataTypeCustomText:      "DataTypeCustomText",
	DataTypeCustomBinary:    "DataTypeCustomBinary",
	DataTypeMarker:          "DataTypeMarker",
	DataTypeLocalReference:  "DataTypeLocalReference",
	DataTypeResourceID:      "DataTypeResourceID",
	DataTypeRemoteReference: "DataTypeRemoteReference",
	DataTypeComment:         "DataTypeComment",
	DataTypePadding:         "DataTypePadding",
}

type MethodType int

const (
	MethodTypeOther = iota
	MethodTypeScalar
	MethodTypeArray
)

type Method struct {
	Name                  string
	Signature             string
	MethodType            MethodType
	AssociatedTypes       DataType
	DefaultImplementation string
}

var (
	nothingToDoImplementation = `
	/* Nothing to do */
`
	beginListImplementation = `
	ctx.BeginList()
`
	beginMapImplementation = `
	ctx.BeginMap()
`
	beginRecordTypeImplementation = `
	ctx.BeginRecordType(identifier)
`
	beginRecordImplementation = `
	ctx.BeginRecord(identifier)
`
	beginNodeImplementation = `
	ctx.BeginNode()
`
	beginEdgeImplementation = `
	ctx.BeginEdge()
`
	endContainerImplementation = `
	ctx.EndContainer(true)
`
	beginMarkerAnyImplementation = `
	ctx.BeginMarkerAnyType(identifier, AllowAny)
`
	LocalReferenceAnyImplementation = `
	ctx.LocalReferenceAnyType(identifier)
`
	arrayAnyImplementation = `
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
`
	stringlikeAnyImplementation = `
	ctx.ValidateFullArrayStringlike(arrayType, data)
`
	beginArrayAnyImplementation = `
	ctx.BeginArrayAnyType(arrayType)
`
)

var (
	BDoc = &Method{
		Name:            "begin document",
		MethodType:      MethodTypeOther,
		Signature:       "OnBeginDocument(ctx *Context)",
		AssociatedTypes: 0,
	}
	EDoc = &Method{
		Name:            "end document",
		MethodType:      MethodTypeOther,
		Signature:       "OnEndDocument(ctx *Context)",
		AssociatedTypes: 0,
	}
	Child = &Method{
		Name:                  "child end",
		MethodType:            MethodTypeOther,
		Signature:             "OnChildContainerEnded(ctx *Context, containerType DataType)",
		AssociatedTypes:       DataTypesContainer | DataTypesAllArrays,
		DefaultImplementation: nothingToDoImplementation,
	}
	Ver = &Method{
		Name:            "version",
		MethodType:      MethodTypeOther,
		Signature:       "OnVersion(ctx *Context, version uint64)",
		AssociatedTypes: 0,
	}
	Pad = &Method{
		Name:                  "padding",
		MethodType:            MethodTypeOther,
		Signature:             "OnPadding(ctx *Context)",
		AssociatedTypes:       DataTypePadding,
		DefaultImplementation: nothingToDoImplementation,
	}
	Comment = &Method{
		Name:                  "comment",
		MethodType:            MethodTypeOther,
		Signature:             "OnComment(ctx *Context)",
		AssociatedTypes:       DataTypeComment,
		DefaultImplementation: nothingToDoImplementation,
	}
	Null = &Method{
		Name:                  "Null",
		MethodType:            MethodTypeOther,
		Signature:             "OnNull(ctx *Context)",
		AssociatedTypes:       DataTypeNull,
		DefaultImplementation: nothingToDoImplementation,
	}
	Key = &Method{
		Name:                  "KEYABLE",
		MethodType:            MethodTypeScalar,
		Signature:             "OnKeyableObject(ctx *Context, objType DataType)",
		AssociatedTypes:       DataTypesKeyable,
		DefaultImplementation: nothingToDoImplementation,
	}
	NonKey = &Method{
		Name:                  "NONKEYABLE",
		MethodType:            MethodTypeScalar,
		Signature:             "OnNonKeyableObject(ctx *Context, objType DataType)",
		AssociatedTypes:       DataTypesNonKeyable,
		DefaultImplementation: nothingToDoImplementation,
		// TODO: What about stuff handled by other methods like list, map, etc?
	}
	List = &Method{
		Name:                  "list",
		MethodType:            MethodTypeOther,
		Signature:             "OnList(ctx *Context)",
		AssociatedTypes:       DataTypeList,
		DefaultImplementation: beginListImplementation,
	}
	Map = &Method{
		Name:                  "map",
		MethodType:            MethodTypeOther,
		Signature:             "OnMap(ctx *Context)",
		AssociatedTypes:       DataTypeMap,
		DefaultImplementation: beginMapImplementation,
	}
	RecordType = &Method{
		Name:                  "recordType",
		MethodType:            MethodTypeOther,
		Signature:             "OnRecordType(ctx *Context, identifier []byte)",
		AssociatedTypes:       DataTypeRecordType,
		DefaultImplementation: beginRecordTypeImplementation,
	}
	Record = &Method{
		Name:                  "record",
		MethodType:            MethodTypeOther,
		Signature:             "OnRecord(ctx *Context, identifier []byte)",
		AssociatedTypes:       DataTypeRecord,
		DefaultImplementation: beginRecordImplementation,
	}
	Edge = &Method{
		Name:                  "edge",
		MethodType:            MethodTypeOther,
		Signature:             "OnEdge(ctx *Context)",
		AssociatedTypes:       DataTypeEdge,
		DefaultImplementation: beginEdgeImplementation,
	}
	Node = &Method{
		Name:                  "node",
		MethodType:            MethodTypeOther,
		Signature:             "OnNode(ctx *Context)",
		AssociatedTypes:       DataTypeNode,
		DefaultImplementation: beginNodeImplementation,
	}
	End = &Method{
		Name:                  "end container",
		MethodType:            MethodTypeOther,
		Signature:             "OnEnd(ctx *Context)",
		AssociatedTypes:       0,
		DefaultImplementation: endContainerImplementation,
	}
	Marker = &Method{
		Name:                  "marker",
		MethodType:            MethodTypeOther,
		Signature:             "OnMarker(ctx *Context, identifier []byte)",
		AssociatedTypes:       DataTypeMarker,
		DefaultImplementation: beginMarkerAnyImplementation,
	}
	Ref = &Method{
		Name:                  "LocalReference",
		MethodType:            MethodTypeOther,
		Signature:             "OnReferenceLocal(ctx *Context, identifier []byte)",
		AssociatedTypes:       DataTypeLocalReference,
		DefaultImplementation: LocalReferenceAnyImplementation,
	}
	Array = &Method{
		Name:                  "array",
		MethodType:            MethodTypeArray,
		Signature:             "OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)",
		AssociatedTypes:       DataTypesAllArrays,
		DefaultImplementation: arrayAnyImplementation,
	}
	Stringlike = &Method{
		Name:                  "array",
		MethodType:            MethodTypeArray,
		Signature:             "OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string)",
		AssociatedTypes:       DataTypesStringlike,
		DefaultImplementation: stringlikeAnyImplementation,
	}
	ABegin = &Method{
		Name:                  "array begin",
		MethodType:            MethodTypeArray,
		Signature:             "OnArrayBegin(ctx *Context, arrayType events.ArrayType)",
		AssociatedTypes:       DataTypesAllArrays,
		DefaultImplementation: beginArrayAnyImplementation,
	}
	AChunk = &Method{
		Name:            "array chunk",
		MethodType:      MethodTypeOther,
		Signature:       "OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)",
		AssociatedTypes: 0,
	}
	AData = &Method{
		Name:            "array data",
		MethodType:      MethodTypeOther,
		Signature:       "OnArrayData(ctx *Context, data []byte)",
		AssociatedTypes: 0,
	}

	allMethods = []*Method{BDoc, EDoc, Child, Ver, Pad, Comment, Null, Key,
		NonKey, List, Map, RecordType, Record, Edge, Node, End,
		Marker, Ref, Array, Stringlike, ABegin, AChunk, AData}
)

type Rule struct {
	Name           string
	FriendlyName   string
	AllowedTypes   DataType
	IncludeMethods []*Method
	ExcludeMethods []*Method
	DefaultMethods []*Method
}

func (_this *Rule) hasMethod(method *Method) bool {
	for _, v := range _this.ExcludeMethods {
		if v == method {
			return false
		}
	}

	if _this.AllowedTypes&method.AssociatedTypes != 0 {
		return true
	}

	for _, v := range _this.IncludeMethods {
		if v == method {
			return true
		}
	}
	return false
}

func (_this *Rule) hasDefaultMethod(method *Method) bool {
	for _, v := range _this.DefaultMethods {
		if v == method {
			return true
		}
	}
	return false
}

var allRules = []Rule{
	{
		Name:           "BeginDocumentRule",
		FriendlyName:   "begin document",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{BDoc},
	},
	{
		Name:           "EndDocumentRule",
		FriendlyName:   "end document",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{EDoc},
	},
	{
		Name:         "TerminalRule",
		FriendlyName: "terminal",
		AllowedTypes: DataTypesNone,
	},
	{
		Name:           "VersionRule",
		FriendlyName:   "version",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{Ver},
	},
	{
		Name:         "TopLevelRule",
		FriendlyName: "top level",
		AllowedTypes: DataTypesTopLevel,
		DefaultMethods: []*Method{
			Pad, Comment, List, Map, RecordType,
			Record, Node, Edge, Marker, ABegin,
		},
	},
	{
		Name:           "ListRule",
		FriendlyName:   "list",
		AllowedTypes:   DataTypesAll,
		IncludeMethods: []*Method{End},
		DefaultMethods: []*Method{
			Child, Pad, Comment, Null, Key, NonKey, List, Map, RecordType,
			Record, Node, Edge, End, Marker, Ref, Array, Stringlike, ABegin,
		},
	},
	{
		Name:           "MapKeyRule",
		FriendlyName:   "map key",
		AllowedTypes:   DataTypesKeyable,
		IncludeMethods: []*Method{End},
		DefaultMethods: []*Method{Pad, Comment, End},
	},
	{
		Name:         "MapValueRule",
		FriendlyName: "map value",
		AllowedTypes: DataTypesAll,
		DefaultMethods: []*Method{
			Pad, Comment, List, Map, RecordType,
			Record, Node, Edge, Marker, ABegin,
		},
	},
	{
		Name:           "RecordTypeRule",
		FriendlyName:   "recordType",
		AllowedTypes:   DataTypesKeyable & (^DataTypeMarker),
		IncludeMethods: []*Method{End},
		ExcludeMethods: []*Method{Marker, Ref},
		DefaultMethods: []*Method{
			Child, Pad, Comment, Key,
		},
	},
	{
		Name:           "RecordRule",
		FriendlyName:   "record",
		AllowedTypes:   DataTypesAll,
		IncludeMethods: []*Method{End},
		DefaultMethods: []*Method{
			Child, Pad, Comment, Null, Key, NonKey, List, Map, RecordType,
			Record, Node, Edge, End, Marker, Ref, Array, Stringlike, ABegin,
		},
	},
	{
		Name:           "ArrayRule",
		FriendlyName:   "array",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AChunk, Comment},
	},
	{
		Name:           "ArrayChunkRule",
		FriendlyName:   "array chunk",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AData},
	},
	{
		Name:           "StringRule",
		FriendlyName:   "string",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AChunk},
	},
	{
		Name:           "StringChunkRule",
		FriendlyName:   "string chunk",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AData},
	},
	{
		Name:           "StringBuilderRule",
		FriendlyName:   "string",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AChunk},
	},
	{
		Name:           "StringBuilderChunkRule",
		FriendlyName:   "string chunk",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AData},
	},
	{
		Name:           "MarkedObjectKeyableRule",
		FriendlyName:   "marked object",
		AllowedTypes:   DataTypesMarkable & DataTypesKeyable,
		DefaultMethods: []*Method{Pad},
	},
	{
		Name:           "MarkedObjectAnyTypeRule",
		FriendlyName:   "marked object",
		AllowedTypes:   DataTypesMarkable,
		DefaultMethods: []*Method{Pad},
	},
	{
		Name:         "EdgeSourceRule",
		FriendlyName: "edge source",
		AllowedTypes: DataTypesNonNull,
		DefaultMethods: []*Method{
			Pad, Comment, List, Map, RecordType,
			Record, Node, Edge,
		},
	},
	{
		Name:         "EdgeDescriptionRule",
		FriendlyName: "edge description",
		AllowedTypes: DataTypesAll,
		DefaultMethods: []*Method{
			Pad, Comment, List, Map, RecordType,
			Record, Node, Edge,
		},
	},
	{
		Name:           "EdgeDestinationRule",
		FriendlyName:   "edge destination",
		AllowedTypes:   DataTypesNonNull,
		IncludeMethods: []*Method{End},
		DefaultMethods: []*Method{
			Child, Pad, Comment, Key, NonKey, List, Map, RecordType,
			Record, Node, Edge, End, Marker, Ref, Array, Stringlike, ABegin,
		},
	},
	{
		Name:         "NodeRule",
		FriendlyName: "node",
		AllowedTypes: DataTypesAll,
		DefaultMethods: []*Method{
			Pad, Comment, List, Map, RecordType,
			Record, Node, Edge,
		},
	},
	{
		Name:           "AwaitEndRule",
		FriendlyName:   "awaitEnd",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{Pad, Comment, End},
		DefaultMethods: []*Method{
			Pad, Comment, End,
		},
	},
}
