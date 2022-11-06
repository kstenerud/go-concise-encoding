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

package builder

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/common"
)

const path = "builder"

var imports = []*common.Import{
	{As: "", Import: "math/big"},
	{As: "", Import: "reflect"},
	{As: "", Import: "github.com/kstenerud/go-concise-encoding/ce/events"},
	{As: "", Import: "github.com/cockroachdb/apd/v2"},
	{As: "", Import: "github.com/kstenerud/go-compact-float"},
	{As: "", Import: "github.com/kstenerud/go-compact-time"},
}

var (
	Null           = "BuildFromNull(ctx *Context, dst reflect.Value) reflect.Value"
	Bool           = "BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value"
	Int            = "BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value"
	Uint           = "BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value"
	BigInt         = "BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value"
	Float          = "BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value"
	BigFloat       = "BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value"
	DFloat         = "BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value"
	BigDFloat      = "BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value"
	UID            = "BuildFromUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value"
	Array          = "BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value"
	SArray         = "BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value"
	CustomBinary   = "BuildFromCustomBinary(ctx *Context, customType uint64, data []byte, dst reflect.Value) reflect.Value"
	CustomText     = "BuildFromCustomText(ctx *Context, customType uint64, data string, dst reflect.Value) reflect.Value"
	Media          = "BuildFromMedia(ctx *Context, mediaType string, data []byte, dst reflect.Value) reflect.Value"
	Time           = "BuildFromTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value"
	Ref            = "BuildFromLocalReference(ctx *Context, id []byte)"
	List           = "BuildNewList(ctx *Context)"
	Map            = "BuildNewMap(ctx *Context)"
	Node           = "BuildNewNode(ctx *Context)"
	Edge           = "BuildNewEdge(ctx *Context)"
	End            = "BuildEndContainer(ctx *Context)"
	ListContents   = "BuildBeginListContents(ctx *Context)"
	MapContents    = "BuildBeginMapContents(ctx *Context)"
	NodeContents   = "BuildBeginNodeContents(ctx *Context)"
	EdgeContents   = "BuildBeginEdgeContents(ctx *Context)"
	NotifyFinished = "NotifyChildContainerFinished(ctx *Context, container reflect.Value)"

	allMethods = []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
		BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
		Node, Edge, ListContents, MapContents, NodeContents,
		EdgeContents, End, Ref, NotifyFinished}
)

type Builder struct {
	Name    string
	Methods []string
}

var builders = []Builder{
	{
		Name: "array",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, Media, Time, List, Map,
			Node, Edge, ListContents, End, Ref, NotifyFinished},
	},
	{
		Name:    "bigDecimalFloat",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "bigFloat",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "bigInt",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "bool",
		Methods: []string{Bool},
	},
	{
		Name:    "compactTime",
		Methods: []string{Null, Time},
	},
	{
		Name:    "custom",
		Methods: []string{CustomBinary, CustomText},
	},
	{
		Name:    "decimalFloat",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "float",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "float32Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "float32Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "float64Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "float64Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name: "ignore",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, Ref, NotifyFinished},
	},
	{
		Name: "ignoreXTimes",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, Ref, NotifyFinished},
	},
	{
		Name: "ignoreContainer",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, End, ListContents, MapContents,
			NodeContents, EdgeContents, Ref, NotifyFinished},
	},
	{
		Name:    "int",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "int8Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "int8Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "int16Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "int16Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "int32Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "int32Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "int64Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "int64Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name: "interface",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, MapContents, ListContents,
			NodeContents, EdgeContents, Ref, NotifyFinished},
	},
	{
		Name: "map",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, MapContents, End, Ref, NotifyFinished},
	},
	{
		Name: "structTemplate",
		Methods: []string{Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, Time, End},
	},
	{
		Name: "structInstance",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, End, Ref, NotifyFinished},
	},
	{
		Name: "markerObject",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, Map,
			Node, Edge, List, End, NotifyFinished},
	},
	{
		Name:    "pBigDecimalFloat",
		Methods: []string{Null, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pBigFloat",
		Methods: []string{Null, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pBigInt",
		Methods: []string{Null, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pCompactTime",
		Methods: []string{Null, Time},
	},
	{
		Name: "ptr",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, ListContents,
			MapContents, NodeContents, EdgeContents, NotifyFinished},
	},
	{
		Name:    "pRid",
		Methods: []string{Null, Array, SArray},
	},
	{
		Name: "slice",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, ListContents, End, Ref, NotifyFinished},
	},
	{
		Name:    "string",
		Methods: []string{Null, Array, SArray},
	},
	{
		Name: "struct",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, MapContents, End, Ref, NotifyFinished},
	},
	{
		Name:    "time",
		Methods: []string{Time},
	},
	{
		Name: "topLevel",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, NotifyFinished},
	},
	{
		Name:    "uint",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "uint8Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "uint8Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "uint16Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "uint16Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "uint32Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "uint32Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "uint64Array",
		Methods: []string{Array, ListContents},
	},
	{
		Name:    "uint64Slice",
		Methods: []string{Null, Array, ListContents},
	},
	{
		Name:    "rid",
		Methods: []string{Array, SArray},
	},
	{
		Name:    "uid",
		Methods: []string{UID},
	},
	{
		Name:    "media",
		Methods: []string{Media},
	},
	{
		Name: "edge",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, EdgeContents, Ref, NotifyFinished},
	},
	{
		Name: "node",
		Methods: []string{Null, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
			BigDFloat, UID, Array, SArray, CustomBinary, CustomText, Media, Time, List, Map,
			Node, Edge, NodeContents, Ref, NotifyFinished},
	},
}

func GenerateCode(projectDir string) {
	common.GenerateGoFile(filepath.Join(projectDir, path), imports, func(writer io.Writer) {
		generateBadEventMethods(writer)
	})
}

func contains(lookingFor string, inSlice []string) bool {
	for _, v := range inSlice {
		if v == lookingFor {
			return true
		}
	}
	return false
}

func hasDstValue(methodSignature string) bool {
	return strings.Contains(methodSignature, ", dst reflect.Value")
}

func openMethod(builder Builder, methodSignature string, writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("func (_this *%sBuilder) %s {\n", builder.Name, methodSignature))); err != nil {
		panic(err)
	}
}

func closeMethod(writer io.Writer) {
	if _, err := writer.Write([]byte("}\n")); err != nil {
		panic(err)
	}
}

func generateBadEventMethod(builder Builder, methodSignature string, writer io.Writer) {
	methodName := methodSignature[:strings.Index(methodSignature, "(")]
	openMethod(builder, methodSignature, writer)

	if hasDstValue(methodSignature) {
		format := "\treturn PanicBadEventBuildingValue(_this, dst, \"%v\")\n"
		// format := "\tpanic(fmt.Errorf(\"BUG: %%v (building type %%v) cannot respond to %s\", reflect.TypeOf(_this), dst.Type()))\n"
		if _, err := writer.Write([]byte(fmt.Sprintf(format, methodName))); err != nil {
			panic(err)
		}
	} else {
		format := "\tPanicBadEvent(_this, \"%v\")\n"
		// format := "\tpanic(fmt.Errorf(\"BUG: %%v cannot respond to %s\", reflect.TypeOf(_this)))\n"
		if _, err := writer.Write([]byte(fmt.Sprintf(format, methodName))); err != nil {
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
	for _, builder := range builders {
		for _, methodSignature := range allMethods {
			if !contains(methodSignature, builder.Methods) {
				generateBadEventMethod(builder, methodSignature, writer)
			}
		}
	}
}
