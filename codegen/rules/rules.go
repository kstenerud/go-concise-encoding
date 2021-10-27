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
	"os"

	"github.com/kstenerud/go-concise-encoding/codegen/datatypes"
	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

const path = "rules"

var imports = []string{
	"fmt",
	"strings",
	"github.com/kstenerud/go-concise-encoding/events",
}

func GenerateCode(projectDir string) {
	generatedFilePath := standard.GetGeneratedCodePath(projectDir, path)
	writer, err := os.Create(generatedFilePath)
	standard.PanicIfError(err, "could not open %s", generatedFilePath)
	defer writer.Close()
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Errorf("Error while generating %v: %v", generatedFilePath, e))
		}
	}()

	standard.WriteHeader(writer, path, imports)
	generateDataTypeType(writer)
	generateBadEventMethods(writer)
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
	gen.AddCustom("AllowNonNil", DataTypesNonNil)
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
	DataTypeNil DataType = 1 << iota
	DataTypeNan
	DataTypeBool
	DataTypeInt
	DataTypeFloat
	DataTypeUID
	DataTypeTime
	DataTypeList
	DataTypeMap
	DataTypeEdge
	DataTypeNode
	DataTypeMarkup
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
	DataTypeReference
	DataTypeResourceID
	DataTypeResourceIDCat
	DataTypeResourceIDRef
	DataTypeConstant
	DataTypeComment
	DataTypePadding
	EndDataTypes

	DataTypeInvalid  = DataType(0)
	DataTypesNone    = DataType(0)
	DataTypesAll     = ^DataTypesNone
	DataTypesNonNil  = ^DataTypeNil
	DataTypesKeyable = DataTypeBool |
		DataTypeInt |
		DataTypeFloat |
		DataTypeUID |
		DataTypeTime |
		DataTypeString |
		DataTypeResourceID |
		DataTypeResourceIDCat |
		DataTypeMarker |
		DataTypeReference |
		DataTypeConstant |
		DataTypePadding |
		DataTypeComment
	DataTypesNonKeyable     = ^DataTypesKeyable
	DataTypesMarkable       = ^(DataTypeMarker | DataTypeReference | DataTypeResourceIDRef | DataTypeConstant | DataTypeComment)
	DataTypesTopLevel       = ^(DataTypeReference)
	DataTypesContainer      = DataTypeList | DataTypeMap | DataTypeEdge | DataTypeNode | DataTypeMarkup
	DataTypesMarkupContents = DataTypeMarkup | DataTypeString | DataTypeComment | DataTypePadding
	DataTypesStringlike     = DataTypeString |
		DataTypeResourceID |
		DataTypeResourceIDCat |
		DataTypeResourceIDRef |
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
	DataTypeInvalid:       "DataTypeInvalid",
	DataTypeNil:           "DataTypeNil",
	DataTypeNan:           "DataTypeNan",
	DataTypeBool:          "DataTypeBool",
	DataTypeInt:           "DataTypeInt",
	DataTypeFloat:         "DataTypeFloat",
	DataTypeUID:           "DataTypeUID",
	DataTypeTime:          "DataTypeTime",
	DataTypeList:          "DataTypeList",
	DataTypeMap:           "DataTypeMap",
	DataTypeEdge:          "DataTypeEdge",
	DataTypeNode:          "DataTypeNode",
	DataTypeMarkup:        "DataTypeMarkup",
	DataTypeString:        "DataTypeString",
	DataTypeMedia:         "DataTypeMedia",
	DataTypeArrayBit:      "DataTypeArrayBit",
	DataTypeArrayUint8:    "DataTypeArrayUint8",
	DataTypeArrayUint16:   "DataTypeArrayUint16",
	DataTypeArrayUint32:   "DataTypeArrayUint32",
	DataTypeArrayUint64:   "DataTypeArrayUint64",
	DataTypeArrayInt8:     "DataTypeArrayInt8",
	DataTypeArrayInt16:    "DataTypeArrayInt16",
	DataTypeArrayInt32:    "DataTypeArrayInt32",
	DataTypeArrayInt64:    "DataTypeArrayInt64",
	DataTypeArrayFloat16:  "DataTypeArrayFloat16",
	DataTypeArrayFloat32:  "DataTypeArrayFloat32",
	DataTypeArrayFloat64:  "DataTypeArrayFloat64",
	DataTypeArrayUID:      "DataTypeArrayUID",
	DataTypeCustomText:    "DataTypeCustomText",
	DataTypeCustomBinary:  "DataTypeCustomBinary",
	DataTypeMarker:        "DataTypeMarker",
	DataTypeReference:     "DataTypeReference",
	DataTypeResourceID:    "DataTypeResourceID",
	DataTypeResourceIDCat: "DataTypeResourceIDCat",
	DataTypeResourceIDRef: "DataTypeResourceIDRef",
	DataTypeConstant:      "DataTypeConstant",
	DataTypeComment:       "DataTypeComment",
	DataTypePadding:       "DataTypePadding",
}

type MethodType int

const (
	MethodTypeOther = iota
	MethodTypeScalar
	MethodTypeArray
)

type Method struct {
	Name            string
	Signature       string
	MethodType      MethodType
	AssociatedTypes DataType
}

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
		Name:            "child end",
		MethodType:      MethodTypeOther,
		Signature:       "OnChildContainerEnded(ctx *Context, containerType DataType)",
		AssociatedTypes: DataTypesContainer | DataTypesAllArrays,
	}
	Ver = &Method{
		Name:            "version",
		MethodType:      MethodTypeOther,
		Signature:       "OnVersion(ctx *Context, version uint64)",
		AssociatedTypes: 0,
	}
	Pad = &Method{
		Name:            "padding",
		MethodType:      MethodTypeOther,
		Signature:       "OnPadding(ctx *Context)",
		AssociatedTypes: DataTypePadding,
	}
	Comment = &Method{
		Name:            "comment",
		MethodType:      MethodTypeOther,
		Signature:       "OnComment(ctx *Context)",
		AssociatedTypes: DataTypeComment,
	}
	Nil = &Method{
		Name:            "Nil",
		MethodType:      MethodTypeOther,
		Signature:       "OnNil(ctx *Context)",
		AssociatedTypes: DataTypeNil,
	}
	Key = &Method{
		Name:            "KEYABLE",
		MethodType:      MethodTypeScalar,
		Signature:       "OnKeyableObject(ctx *Context, objType DataType)",
		AssociatedTypes: DataTypesKeyable,
	}
	NonKey = &Method{
		Name:            "NONKEYABLE",
		MethodType:      MethodTypeScalar,
		Signature:       "OnNonKeyableObject(ctx *Context, objType DataType)",
		AssociatedTypes: DataTypesNonKeyable,
		// TODO: What about stuff handled by other methods like list, map, etc?
	}
	List = &Method{
		Name:            "list",
		MethodType:      MethodTypeOther,
		Signature:       "OnList(ctx *Context)",
		AssociatedTypes: DataTypeList,
	}
	Map = &Method{
		Name:            "map",
		MethodType:      MethodTypeOther,
		Signature:       "OnMap(ctx *Context)",
		AssociatedTypes: DataTypeMap,
	}
	Edge = &Method{
		Name:            "edge",
		MethodType:      MethodTypeOther,
		Signature:       "OnEdge(ctx *Context)",
		AssociatedTypes: DataTypeEdge,
	}
	Node = &Method{
		Name:            "node",
		MethodType:      MethodTypeOther,
		Signature:       "OnNode(ctx *Context)",
		AssociatedTypes: DataTypeNode,
	}
	Markup = &Method{
		Name:            "markup",
		MethodType:      MethodTypeOther,
		Signature:       "OnMarkup(ctx *Context, identifier []byte)",
		AssociatedTypes: DataTypeMarkup,
	}
	End = &Method{
		Name:            "end container",
		MethodType:      MethodTypeOther,
		Signature:       "OnEnd(ctx *Context)",
		AssociatedTypes: 0,
	}
	Marker = &Method{
		Name:            "marker",
		MethodType:      MethodTypeOther,
		Signature:       "OnMarker(ctx *Context, identifier []byte)",
		AssociatedTypes: DataTypeMarker,
	}
	Ref = &Method{
		Name:            "reference",
		MethodType:      MethodTypeOther,
		Signature:       "OnReference(ctx *Context, identifier []byte)",
		AssociatedTypes: DataTypeReference,
	}
	Const = &Method{
		Name:            "constant",
		MethodType:      MethodTypeOther,
		Signature:       "OnConstant(ctx *Context, identifier []byte)",
		AssociatedTypes: DataTypeConstant,
	}
	Array = &Method{
		Name:            "array",
		MethodType:      MethodTypeArray,
		Signature:       "OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)",
		AssociatedTypes: DataTypesAllArrays,
	}
	Stringlike = &Method{
		Name:            "array",
		MethodType:      MethodTypeArray,
		Signature:       "OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string)",
		AssociatedTypes: DataTypesStringlike,
	}
	ABegin = &Method{
		Name:            "array begin",
		MethodType:      MethodTypeArray,
		Signature:       "OnArrayBegin(ctx *Context, arrayType events.ArrayType)",
		AssociatedTypes: DataTypesAllArrays,
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

	allMethods = []*Method{BDoc, EDoc, Child, Ver, Pad, Comment, Nil, Key,
		NonKey, List, Map, Edge, Node, Markup, End, Marker, Ref, Const,
		Array, Stringlike, ABegin, AChunk, AData}
)

type Rule struct {
	Name           string
	FriendlyName   string
	AllowedTypes   DataType
	IncludeMethods []*Method
	ExcludeMethods []*Method
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
	},
	{
		Name:           "ListRule",
		FriendlyName:   "list",
		AllowedTypes:   DataTypesAll,
		IncludeMethods: []*Method{End},
	},
	{
		Name:           "MapKeyRule",
		FriendlyName:   "map key",
		AllowedTypes:   DataTypesKeyable,
		IncludeMethods: []*Method{End},
	},
	{
		Name:         "MapValueRule",
		FriendlyName: "map value",
		AllowedTypes: DataTypesAll,
	},
	{
		Name:           "MarkupKeyRule",
		FriendlyName:   "markup key",
		AllowedTypes:   DataTypesKeyable,
		IncludeMethods: []*Method{End},
	},
	{
		Name:         "MarkupValueRule",
		FriendlyName: "markup value",
		AllowedTypes: DataTypesAll,
	},
	{
		Name:           "MarkupContentsRule",
		FriendlyName:   "markup contents",
		AllowedTypes:   DataTypesMarkupContents,
		IncludeMethods: []*Method{End},
		ExcludeMethods: []*Method{Key, NonKey},
	},
	{
		Name:         "ArrayRule",
		FriendlyName: "array",
		// TODO: AllowedTypes:   DataTypeComment,
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AChunk},
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
		Name:           "MediaTypeRule",
		FriendlyName:   "media type",
		AllowedTypes:   DataTypesNone,
		IncludeMethods: []*Method{AChunk},
	},
	{
		Name:           "MediaTypeChunkRule",
		IncludeMethods: []*Method{AData},
		FriendlyName:   "media type chunk",
	},
	{
		Name:         "MarkedObjectKeyableRule",
		FriendlyName: "marked object",
		AllowedTypes: DataTypesMarkable & DataTypesKeyable,
	},
	{
		Name:         "MarkedObjectAnyTypeRule",
		FriendlyName: "marked object",
		AllowedTypes: DataTypesMarkable,
	},
	{
		Name:         "ConstantKeyableRule",
		FriendlyName: "constant",
		AllowedTypes: DataTypesKeyable,
	},
	{
		Name:         "ConstantAnyTypeRule",
		FriendlyName: "constant",
		AllowedTypes: DataTypesAll,
	},
	{
		Name:         "EdgeSourceRule",
		FriendlyName: "edge source",
		AllowedTypes: DataTypesNonNil,
	},
	{
		Name:         "EdgeDescriptionRule",
		FriendlyName: "edge description",
		AllowedTypes: DataTypesAll,
	},
	{
		Name:         "EdgeDestinationRule",
		FriendlyName: "edge destination",
		AllowedTypes: DataTypesNonNil,
	},
	{
		Name:         "NodeRule",
		FriendlyName: "node",
		AllowedTypes: DataTypesAll,
	},
}
