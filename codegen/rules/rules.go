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
	gen.AddCustom("AllowAny", AllowAny)
	gen.AddCustom("AllowKeyable", AllowKeyable)
	gen.AddCustom("AllowResource", AllowResource)
	gen.AddCustom("AllowSubject", AllowSubject)
	gen.AddCustom("AllowPredicate", AllowPredicate)
	gen.AddCustom("AllowObject", AllowObject)
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
			if !containsMethod(method, rule.Methods) {
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

func containsMethod(lookingFor *Method, inSlice []*Method) bool {
	for _, v := range inSlice {
		if v == lookingFor {
			return true
		}
	}
	return false
}

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
	DataTypeNA
	DataTypeNan
	DataTypeBool
	DataTypeInt
	DataTypeFloat
	DataTypeUID
	DataTypeTime
	DataTypeList
	DataTypeMap
	DataTypeMarkup
	DataTypeComment
	DataTypeRelationship
	DataTypeString
	DataTypeResourceID
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
	DataTypeMedia
	DataTypeRIDReference
	DataTypeCustomText
	DataTypeCustomBinary
	DataTypeResourceList
	EndDataTypes

	DataTypeInvalid = DataType(0)
	AllowAny        = ^DataType(0)
	AllowKeyable    = DataTypeBool | DataTypeInt | DataTypeFloat | DataTypeUID | DataTypeTime | DataTypeString | DataTypeResourceID
	AllowResource   = DataTypeMap | DataTypeRelationship | DataTypeResourceList | DataTypeResourceID
	AllowSubject    = AllowResource
	AllowPredicate  = DataTypeResourceID
	AllowObject     = AllowAny
)

func (_this DataType) String() string {
	return datatypes.FlagToString(dataTypeNames, _this)
}

var dataTypeNames = map[interface{}]string{
	DataTypeInvalid:      "DataTypeInvalid",
	DataTypeNil:          "DataTypeNil",
	DataTypeNA:           "DataTypeNA",
	DataTypeNan:          "DataTypeNan",
	DataTypeBool:         "DataTypeBool",
	DataTypeInt:          "DataTypeInt",
	DataTypeFloat:        "DataTypeFloat",
	DataTypeUID:          "DataTypeUID",
	DataTypeTime:         "DataTypeTime",
	DataTypeList:         "DataTypeList",
	DataTypeMap:          "DataTypeMap",
	DataTypeMarkup:       "DataTypeMarkup",
	DataTypeComment:      "DataTypeComment",
	DataTypeRelationship: "DataTypeRelationship",
	DataTypeString:       "DataTypeString",
	DataTypeResourceID:   "DataTypeResourceID",
	DataTypeArrayBit:     "DataTypeArrayBit",
	DataTypeArrayUint8:   "DataTypeArrayUint8",
	DataTypeArrayUint16:  "DataTypeArrayUint16",
	DataTypeArrayUint32:  "DataTypeArrayUint32",
	DataTypeArrayUint64:  "DataTypeArrayUint64",
	DataTypeArrayInt8:    "DataTypeArrayInt8",
	DataTypeArrayInt16:   "DataTypeArrayInt16",
	DataTypeArrayInt32:   "DataTypeArrayInt32",
	DataTypeArrayInt64:   "DataTypeArrayInt64",
	DataTypeArrayFloat16: "DataTypeArrayFloat16",
	DataTypeArrayFloat32: "DataTypeArrayFloat32",
	DataTypeArrayFloat64: "DataTypeArrayFloat64",
	DataTypeArrayUID:     "DataTypeArrayUID",
	DataTypeMedia:        "DataTypeMedia",
	DataTypeRIDReference: "DataTypeRIDReference",
	DataTypeCustomText:   "DataTypeCustomText",
	DataTypeCustomBinary: "DataTypeCustomBinary",
	DataTypeResourceList: "DataTypeResourceList",
}

type ContainerType int

const (
	ContainerTypeList = iota
	ContainerTypeMap
)

type MethodType int

const (
	MethodTypeOther = iota
	MethodTypeScalar
	MethodTypeArray
)

type Method struct {
	Name       string
	Signature  string
	MethodType MethodType
}

var (
	BDoc = &Method{
		Name:       "begin document",
		MethodType: MethodTypeOther,
		Signature:  "OnBeginDocument(ctx *Context)",
	}
	EDoc = &Method{
		Name:       "end document",
		MethodType: MethodTypeOther,
		Signature:  "OnEndDocument(ctx *Context)",
	}
	ECtr = &Method{
		Name:       "child end",
		MethodType: MethodTypeOther,
		Signature:  "OnChildContainerEnded(ctx *Context, containerType DataType)",
	}
	Ver = &Method{
		Name:       "version",
		MethodType: MethodTypeOther,
		Signature:  "OnVersion(ctx *Context, version uint64)",
	}
	NA = &Method{
		Name:       "NA",
		MethodType: MethodTypeOther,
		Signature:  "OnNA(ctx *Context)",
	}
	Pad = &Method{
		Name:       "padding",
		MethodType: MethodTypeOther,
		Signature:  "OnPadding(ctx *Context)",
	}
	Nil = &Method{
		Name:       "Nil",
		MethodType: MethodTypeOther,
		Signature:  "OnNil(ctx *Context)",
	}
	Key = &Method{
		Name:       "KEYABLE",
		MethodType: MethodTypeScalar,
		Signature:  "OnKeyableObject(ctx *Context, objType DataType)",
	}
	NonKey = &Method{
		Name:       "NONKEYABLE",
		MethodType: MethodTypeScalar,
		Signature:  "OnNonKeyableObject(ctx *Context, objType DataType)",
	}
	List = &Method{
		Name:       "list",
		MethodType: MethodTypeOther,
		Signature:  "OnList(ctx *Context)",
	}
	Map = &Method{
		Name:       "map",
		MethodType: MethodTypeOther,
		Signature:  "OnMap(ctx *Context)",
	}
	Markup = &Method{
		Name:       "markup",
		MethodType: MethodTypeOther,
		Signature:  "OnMarkup(ctx *Context, identifier []byte)",
	}
	Comment = &Method{
		Name:       "comment",
		MethodType: MethodTypeOther,
		Signature:  "OnComment(ctx *Context)",
	}
	End = &Method{
		Name:       "end container",
		MethodType: MethodTypeOther,
		Signature:  "OnEnd(ctx *Context)",
	}
	Rel = &Method{
		Name:       "relationship",
		MethodType: MethodTypeOther,
		Signature:  "OnRelationship(ctx *Context)",
	}
	Marker = &Method{
		Name:       "marker",
		MethodType: MethodTypeOther,
		Signature:  "OnMarker(ctx *Context, identifier []byte)",
	}
	Ref = &Method{
		Name:       "reference",
		MethodType: MethodTypeOther,
		Signature:  "OnReference(ctx *Context, identifier []byte)",
	}
	RIDRef = &Method{
		Name:       "RID reference",
		MethodType: MethodTypeOther,
		Signature:  "OnRIDReference(ctx *Context)",
	}
	Const = &Method{
		Name:       "constant",
		MethodType: MethodTypeOther,
		Signature:  "OnConstant(ctx *Context, identifier []byte)",
	}
	Array = &Method{
		Name:       "array",
		MethodType: MethodTypeArray,
		Signature:  "OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)",
	}
	SArray = &Method{
		Name:       "array",
		MethodType: MethodTypeArray,
		Signature:  "OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string)",
	}
	ABegin = &Method{
		Name:       "array begin",
		MethodType: MethodTypeArray,
		Signature:  "OnArrayBegin(ctx *Context, arrayType events.ArrayType)",
	}
	AChunk = &Method{
		Name:       "array chunk",
		MethodType: MethodTypeOther,
		Signature:  "OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)",
	}
	AData = &Method{
		Name:       "array data",
		MethodType: MethodTypeOther,
		Signature:  "OnArrayData(ctx *Context, data []byte)",
	}

	allMethods = []*Method{BDoc, EDoc, ECtr, Ver, NA, Pad, Nil, Key, NonKey, List,
		Map, Markup, Comment, End, Rel, Marker, Ref, RIDRef, Const, Array, SArray,
		ABegin, AChunk, AData}
)

type Rule struct {
	Name         string
	FriendlyName string
	Methods      []*Method
}

var allRules = []Rule{
	{
		Name:         "BeginDocumentRule",
		FriendlyName: "begin document",
		Methods:      []*Method{BDoc},
	},
	{
		Name:         "EndDocumentRule",
		FriendlyName: "end document",
		Methods:      []*Method{EDoc},
	},
	{
		Name:         "TerminalRule",
		FriendlyName: "terminal",
		Methods:      []*Method{},
	},
	{
		Name:         "VersionRule",
		FriendlyName: "version",
		Methods:      []*Method{Ver},
	},
	{
		Name:         "TopLevelRule",
		FriendlyName: "top level",
		Methods:      []*Method{ECtr, NA, Pad, Nil, Key, NonKey, List, Map, Markup, Comment, Rel, Marker, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:         "NARule",
		FriendlyName: "NA",
		Methods:      []*Method{ECtr, Pad, Nil, Key, NonKey, List, Map, Markup, Rel, Array, SArray, ABegin},
	},
	{
		Name:         "ListRule",
		FriendlyName: "list",
		Methods:      []*Method{ECtr, NA, Pad, Nil, Key, NonKey, List, Map, Markup, Comment, End, Rel, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:         "ResourceListRule",
		FriendlyName: "resource list",
		Methods:      []*Method{ECtr, Pad, Comment, Map, End, Rel, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:         "MapKeyRule",
		FriendlyName: "map key",
		Methods:      []*Method{ECtr, Pad, Key, Comment, End, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:         "MapValueRule",
		FriendlyName: "map value",
		Methods:      []*Method{ECtr, NA, Pad, Nil, Key, NonKey, List, Map, Markup, Comment, Rel, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:         "MarkupKeyRule",
		FriendlyName: "markup key",
		Methods:      []*Method{ECtr, Pad, Key, Comment, End, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:         "MarkupValueRule",
		FriendlyName: "markup value",
		Methods:      []*Method{ECtr, NA, Pad, Nil, Key, NonKey, List, Map, Markup, Comment, Rel, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:         "MarkupContentsRule",
		FriendlyName: "markup contents",
		Methods:      []*Method{ECtr, Pad, Markup, Comment, End, Array, SArray, ABegin},
	},
	{
		Name:         "CommentRule",
		FriendlyName: "comment",
		Methods:      []*Method{ECtr, Pad, Comment, End, Array, SArray, ABegin},
	},
	{
		Name:         "ArrayRule",
		FriendlyName: "array",
		Methods:      []*Method{AChunk},
	},
	{
		Name:         "ArrayChunkRule",
		FriendlyName: "array chunk",
		Methods:      []*Method{AData},
	},
	{
		Name:         "StringRule",
		FriendlyName: "string",
		Methods:      []*Method{AChunk},
	},
	{
		Name:         "StringChunkRule",
		FriendlyName: "string chunk",
		Methods:      []*Method{AData},
	},
	{
		Name:         "StringBuilderRule",
		FriendlyName: "string",
		Methods:      []*Method{AChunk},
	},
	{
		Name:         "StringBuilderChunkRule",
		FriendlyName: "string chunk",
		Methods:      []*Method{AData},
	},
	{
		Name:         "MediaTypeRule",
		FriendlyName: "media type",
		Methods:      []*Method{AChunk},
	},
	{
		Name:         "MediaTypeChunkRule",
		FriendlyName: "media type chunk",
		Methods:      []*Method{AData},
	},
	{
		Name:         "MarkedObjectKeyableRule",
		FriendlyName: "marked object",
		Methods:      []*Method{ECtr, Pad, Key, Array, SArray, ABegin},
	},
	{
		Name:         "MarkedObjectAnyTypeRule",
		FriendlyName: "marked object",
		Methods:      []*Method{ECtr, Pad, Nil, Key, NonKey, List, Map, Markup, Rel, Array, SArray, ABegin},
	},
	{
		Name:         "ConstantKeyableRule",
		FriendlyName: "constant",
		Methods:      []*Method{ECtr, Pad, Key, Array, SArray, ABegin},
	},
	{
		Name:         "ConstantAnyTypeRule",
		FriendlyName: "constant",
		Methods:      []*Method{ECtr, Pad, Nil, Key, NonKey, NA, List, Map, Markup, Rel, Array, SArray, ABegin},
	},
	{
		Name:         "RIDReferenceRule",
		FriendlyName: "RID reference",
		Methods:      []*Method{Pad, Array, SArray, ABegin, ECtr},
	},
	{
		Name:         "SubjectRule",
		FriendlyName: "relationship subject",
		Methods:      []*Method{ECtr, Pad, List, Map, Comment, Rel, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:         "PredicateRule",
		FriendlyName: "relationship predicate",
		Methods:      []*Method{ECtr, Pad, Comment, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:         "ObjectRule",
		FriendlyName: "relationship object",
		Methods:      []*Method{ECtr, NA, Pad, Nil, Key, NonKey, List, Map, Markup, Comment, Rel, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
}
