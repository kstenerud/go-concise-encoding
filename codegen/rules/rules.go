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
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

const path = "rules"

var imports = []string{
	"fmt",
	"github.com/kstenerud/go-concise-encoding/events",
}

var (
	BDoc    = "OnBeginDocument(ctx *Context)"
	EDoc    = "OnEndDocument(ctx *Context)"
	ECtr    = "OnChildContainerEnded(ctx *Context, cType DataType)"
	Ver     = "OnVersion(ctx *Context, version uint64)"
	NA      = "OnNA(ctx *Context)"
	Pad     = "OnPadding(ctx *Context)"
	Key     = "OnKeyableObject(ctx *Context, objType string)"
	NonKey  = "OnNonKeyableObject(ctx *Context, objType string)"
	List    = "OnList(ctx *Context)"
	Map     = "OnMap(ctx *Context)"
	Markup  = "OnMarkup(ctx *Context, identifier []byte)"
	Comment = "OnComment(ctx *Context)"
	End     = "OnEnd(ctx *Context)"
	Marker  = "OnMarker(ctx *Context, identifier []byte)"
	Ref     = "OnReference(ctx *Context, identifier []byte)"
	RIDRef  = "OnRIDReference(ctx *Context)"
	Const   = "OnConstant(ctx *Context, identifier []byte)"
	Array   = "OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)"
	SArray  = "OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string)"
	ABegin  = "OnArrayBegin(ctx *Context, arrayType events.ArrayType)"
	AChunk  = "OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)"
	AData   = "OnArrayData(ctx *Context, data []byte)"

	allMethods = []string{BDoc, EDoc, ECtr, Ver, NA, Pad, Key, NonKey, List, Map,
		Markup, Comment, End, Marker, Ref, RIDRef, Const, Array, SArray, ABegin,
		AChunk, AData}
)

type RuleClass struct {
	Name    string
	Methods []string
}

var ruleClasses = []RuleClass{
	{
		Name:    "BeginDocumentRule",
		Methods: []string{BDoc},
	},
	{
		Name:    "EndDocumentRule",
		Methods: []string{EDoc},
	},
	{
		Name:    "TerminalRule",
		Methods: []string{},
	},
	{
		Name:    "VersionRule",
		Methods: []string{Ver},
	},
	{
		Name:    "TopLevelRule",
		Methods: []string{ECtr, NA, Pad, Key, NonKey, List, Map, Markup, Comment, Marker, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:    "NARule",
		Methods: []string{ECtr, Pad, Key, NonKey, List, Map, Markup, Array, SArray, ABegin},
	},
	{
		Name:    "ListRule",
		Methods: []string{ECtr, NA, Pad, Key, NonKey, List, Map, Markup, Comment, End, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:    "MapKeyRule",
		Methods: []string{ECtr, Pad, Key, Comment, End, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:    "MapValueRule",
		Methods: []string{ECtr, NA, Pad, Key, NonKey, List, Map, Markup, Comment, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:    "MarkupKeyRule",
		Methods: []string{ECtr, Pad, Key, Comment, End, Marker, Ref, Const, Array, SArray, ABegin},
	},
	{
		Name:    "MarkupValueRule",
		Methods: []string{ECtr, NA, Pad, Key, NonKey, List, Map, Markup, Comment, Marker, Ref, RIDRef, Const, Array, SArray, ABegin},
	},
	{
		Name:    "MarkupContentsRule",
		Methods: []string{ECtr, Pad, Markup, Comment, End, Array, SArray, ABegin},
	},
	{
		Name:    "CommentRule",
		Methods: []string{ECtr, Pad, Comment, End, Array, SArray, ABegin},
	},
	{
		Name:    "ArrayRule",
		Methods: []string{AChunk},
	},
	{
		Name:    "ArrayChunkRule",
		Methods: []string{AData},
	},
	{
		Name:    "StringRule",
		Methods: []string{AChunk},
	},
	{
		Name:    "StringChunkRule",
		Methods: []string{AData},
	},
	{
		Name:    "StringBuilderRule",
		Methods: []string{AChunk},
	},
	{
		Name:    "StringBuilderChunkRule",
		Methods: []string{AData},
	},
	{
		Name:    "MarkedObjectKeyableRule",
		Methods: []string{ECtr, Pad, Key, Array, SArray, ABegin},
	},
	{
		Name:    "MarkedObjectAnyTypeRule",
		Methods: []string{ECtr, Pad, Key, NonKey, List, Map, Markup, Array, SArray, ABegin},
	},
	{
		Name:    "ConstantKeyableRule",
		Methods: []string{ECtr, Pad, Key, Array, SArray, ABegin},
	},
	{
		Name:    "ConstantAnyTypeRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, List, Map, Markup, Array, SArray, ABegin},
	},
	{
		Name:    "RIDReferenceRule",
		Methods: []string{Pad, Array, SArray, ABegin, ECtr},
	},
	{
		Name:    "RIDCatRule",
		Methods: []string{Pad, Array, SArray, ABegin, ECtr},
	},
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
	generateBadEventMethods(writer)
}

func contains(lookingFor string, inSlice []string) bool {
	for _, v := range inSlice {
		if v == lookingFor {
			return true
		}
	}
	return false
}

func openMethod(rule RuleClass, methodSignature string, writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("func (_this *%s) %s {\n", rule.Name, methodSignature))); err != nil {
		panic(err)
	}
}

func closeMethod(writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("}\n"))); err != nil {
		panic(err)
	}
}

func generateBadEventMethod(rule RuleClass, methodSignature string, writer io.Writer) {
	methodName := methodSignature[:strings.Index(methodSignature, "(")]
	openMethod(rule, methodSignature, writer)
	id := methodName[2:]
	if id == "KeyableObject" || id == "NonKeyableObject" {
		format := "\tpanic(fmt.Errorf(\"%%v does not allow %%s\", _this, objType))\n"
		if _, err := writer.Write([]byte(fmt.Sprintf(format))); err != nil {
			panic(err)
		}
	} else {
		format := "\tpanic(fmt.Errorf(\"%%v does not allow %s\", _this))\n"
		if _, err := writer.Write([]byte(fmt.Sprintf(format, id))); err != nil {
			panic(err)
		}
	}

	closeMethod(writer)
}

type ContainerType int

const (
	ContainerTypeList = iota
	ContainerTypeMap
)

func generateBadEventMethods(writer io.Writer) {
	for _, rule := range ruleClasses {
		for _, methodSignature := range allMethods {
			if !contains(methodSignature, rule.Methods) {
				generateBadEventMethod(rule, methodSignature, writer)
			}
		}
	}
}
