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
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

var (
	BDoc    = "OnBeginDocument(ctx *Context)"
	EDoc    = "OnEndDocument(ctx *Context)"
	ECtr    = "OnChildContainerEnded(ctx *Context, cType DataType)"
	Ver     = "OnVersion(ctx *Context, version uint64)"
	Pad     = "OnPadding(ctx *Context)"
	Key     = "OnKeyableObject(ctx *Context)"
	NonKey  = "OnNonKeyableObject(ctx *Context)"
	NA      = "OnNA(ctx *Context)"
	Int     = "OnInt(ctx *Context, value int64)"
	PInt    = "OnPositiveInt(ctx *Context, value uint64)"
	BInt    = "OnBigInt(ctx *Context, value *big.Int)"
	Float   = "OnFloat(ctx *Context, value float64)"
	BFloat  = "OnBigFloat(ctx *Context, value *big.Float)"
	DFloat  = "OnDecimalFloat(ctx *Context, value compact_float.DFloat)"
	BDFloat = "OnBigDecimalFloat(ctx *Context, value *apd.Decimal)"
	List    = "OnList(ctx *Context)"
	Map     = "OnMap(ctx *Context)"
	Markup  = "OnMarkup(ctx *Context)"
	Meta    = "OnMetadata(ctx *Context)"
	Comment = "OnComment(ctx *Context)"
	End     = "OnEnd(ctx *Context)"
	Marker  = "OnMarker(ctx *Context)"
	Ref     = "OnReference(ctx *Context)"
	Array   = "OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8)"
	ABegin  = "OnArrayBegin(ctx *Context, arrayType events.ArrayType)"
	AChunk  = "OnArrayChunk(ctx *Context, length uint64, moreChunksFollow bool)"
	AData   = "OnArrayData(ctx *Context, data []byte)"

	allMethods = []string{BDoc, EDoc, ECtr, Ver, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, End, Marker, Ref, Array, ABegin, AChunk, AData}
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
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "NARule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "ListRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, End, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MapKeyRule",
		Methods: []string{ECtr, Pad, Key, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, Meta, Comment, End, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MapValueRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MarkupNameRule",
		Methods: []string{ECtr, Pad, Key, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MarkupKeyRule",
		Methods: []string{ECtr, Pad, Key, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, Meta, Comment, End, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MarkupValueRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MarkupContentsRule",
		Methods: []string{ECtr, Pad, Markup, Comment, End, Array, ABegin},
	},
	{
		Name:    "CommentRule",
		Methods: []string{Pad, Comment, End, Array, ABegin},
	},
	{
		Name:    "MetaKeyRule",
		Methods: []string{ECtr, Pad, Key, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, Meta, Comment, End, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MetaValueRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
	},
	{
		Name:    "MetaCompletionRule",
		Methods: []string{Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Meta, Comment, Marker, Ref, Array, ABegin},
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
		Name:    "MarkerIDKeyableRule",
		Methods: []string{Pad, Int, PInt, BInt, Array, ABegin, ECtr},
	},
	{
		Name:    "MarkerIDAnyTypeRule",
		Methods: []string{Pad, Int, PInt, BInt, Array, ABegin, ECtr},
	},
	{
		Name:    "MarkedObjectKeyableRule",
		Methods: []string{ECtr, Pad, Key, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, Ref, Array, ABegin},
	},
	{
		Name:    "MarkedObjectAnyTypeRule",
		Methods: []string{ECtr, Pad, Key, NonKey, NA, Int, PInt, BInt, Float, BFloat, DFloat, BDFloat, List, Map, Markup, Ref, Array, ABegin},
	},
	{
		Name:    "ReferenceKeyableRule",
		Methods: []string{Pad, Int, PInt, BInt, Array, ABegin, ECtr},
	},
	{
		Name:    "ReferenceAnyTypeRule",
		Methods: []string{Pad, Int, PInt, BInt, Array, ABegin, ECtr},
	},
	{
		Name:    "TLReferenceRIDRule",
		Methods: []string{Pad, Array, ABegin, ECtr},
	},
}

func GenerateCode(projectDir string) {
	generatedFilePath := filepath.Join(projectDir, "rules/generated-do-not-edit.go")
	writer, err := os.Create(generatedFilePath)
	if err != nil {
		panic(fmt.Errorf("could not open %s: %v", generatedFilePath, err))
	}
	defer writer.Close()

	if _, err := writer.WriteString(codeHeader); err != nil {
		panic(fmt.Errorf("could not write to %s: %v", generatedFilePath, err))
	}

	generateBadEventMethods(writer)
}

var codeHeader = standard.Header + `package rules

import (
	"fmt"
	"math/big"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

`

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

	format := "\tpanic(fmt.Errorf(\"%%v does not allow %s\", _this))\n"
	if _, err := writer.Write([]byte(fmt.Sprintf(format, id))); err != nil {
		panic(err)
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
