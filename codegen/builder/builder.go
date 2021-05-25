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
	"os"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

const path = "builder"

var imports = []string{
	"fmt",
	"math/big",
	"reflect",
	"time",
	"github.com/kstenerud/go-concise-encoding/events",
	"github.com/cockroachdb/apd/v2",
	"github.com/kstenerud/go-compact-float",
	"github.com/kstenerud/go-compact-time",
}

var (
	Nil            = "BuildFromNil(ctx *Context, dst reflect.Value) reflect.Value"
	Bool           = "BuildFromBool(ctx *Context, value bool, dst reflect.Value) reflect.Value"
	Int            = "BuildFromInt(ctx *Context, value int64, dst reflect.Value) reflect.Value"
	Uint           = "BuildFromUint(ctx *Context, value uint64, dst reflect.Value) reflect.Value"
	BigInt         = "BuildFromBigInt(ctx *Context, value *big.Int, dst reflect.Value) reflect.Value"
	Float          = "BuildFromFloat(ctx *Context, value float64, dst reflect.Value) reflect.Value"
	BigFloat       = "BuildFromBigFloat(ctx *Context, value *big.Float, dst reflect.Value) reflect.Value"
	DFloat         = "BuildFromDecimalFloat(ctx *Context, value compact_float.DFloat, dst reflect.Value) reflect.Value"
	BigDFloat      = "BuildFromBigDecimalFloat(ctx *Context, value *apd.Decimal, dst reflect.Value) reflect.Value"
	UUID           = "BuildFromUUID(ctx *Context, value []byte, dst reflect.Value) reflect.Value"
	Array          = "BuildFromArray(ctx *Context, arrayType events.ArrayType, value []byte, dst reflect.Value) reflect.Value"
	SArray         = "BuildFromStringlikeArray(ctx *Context, arrayType events.ArrayType, value string, dst reflect.Value) reflect.Value"
	Time           = "BuildFromTime(ctx *Context, value time.Time, dst reflect.Value) reflect.Value"
	CTime          = "BuildFromCompactTime(ctx *Context, value compact_time.Time, dst reflect.Value) reflect.Value"
	Ref            = "BuildFromReference(ctx *Context, id []byte)"
	ListInit       = "BuildInitiateList(ctx *Context)"
	MapInit        = "BuildInitiateMap(ctx *Context)"
	End            = "BuildEndContainer(ctx *Context)"
	List           = "BuildBeginListContents(ctx *Context)"
	Map            = "BuildBeginMapContents(ctx *Context)"
	NotifyFinished = "NotifyChildContainerFinished(ctx *Context, container reflect.Value)"

	allMethods = []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat,
		BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, List,
		Map, End, Ref, NotifyFinished}
)

type Builder struct {
	Name    string
	Methods []string
}

var builders = []Builder{
	{
		Name:    "array",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, List, End, Ref, NotifyFinished},
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
		Methods: []string{Nil, Time, CTime},
	},
	{
		Name:    "custom",
		Methods: []string{Array},
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
		Methods: []string{Array, List},
	},
	{
		Name:    "float32Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "float64Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "float64Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "ignore",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, List, End, Map, Ref, NotifyFinished},
	},
	{
		Name:    "int",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "int8Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "int8Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "int16Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "int16Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "int32Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "int32Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "int64Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "int64Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "interface",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, Map, List, Ref, NotifyFinished},
	},
	{
		Name:    "map",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, Map, End, Ref, NotifyFinished},
	},
	{
		Name:    "markerObject",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, MapInit, ListInit, End, NotifyFinished},
	},
	{
		Name:    "pBigDecimalFloat",
		Methods: []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pBigFloat",
		Methods: []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pBigInt",
		Methods: []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "pCompactTime",
		Methods: []string{Nil, Time, CTime},
	},
	{
		Name:    "ptr",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, List, Map, NotifyFinished},
	},
	{
		Name:    "pRid",
		Methods: []string{Nil, Array, SArray},
	},
	{
		Name:    "slice",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, List, End, Ref, NotifyFinished},
	},
	{
		Name:    "string",
		Methods: []string{Nil, Array, SArray},
	},
	{
		Name:    "struct",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, Map, End, Ref, NotifyFinished},
	},
	{
		Name:    "time",
		Methods: []string{Time, CTime},
	},
	{
		Name:    "topLevel",
		Methods: []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, SArray, Time, CTime, ListInit, MapInit, NotifyFinished},
	},
	{
		Name:    "uint",
		Methods: []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	},
	{
		Name:    "uint8Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "uint8Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "uint16Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "uint16Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "uint32Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "uint32Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "uint64Array",
		Methods: []string{Array, List},
	},
	{
		Name:    "uint64Slice",
		Methods: []string{Nil, Array, List},
	},
	{
		Name:    "rid",
		Methods: []string{Array, SArray},
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

func hasDstValue(methodSignature string) bool {
	return strings.Contains(methodSignature, ", dst reflect.Value")
}

func openMethod(builder Builder, methodSignature string, writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("func (_this *%sBuilder) %s {\n", builder.Name, methodSignature))); err != nil {
		panic(err)
	}
}

func closeMethod(writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("}\n"))); err != nil {
		panic(err)
	}
}

func generateBadEventMethod(builder Builder, methodSignature string, writer io.Writer) {
	methodName := methodSignature[:strings.Index(methodSignature, "(")]
	openMethod(builder, methodSignature, writer)

	if hasDstValue(methodSignature) {
		format := "\tpanic(fmt.Errorf(\"BUG: %%v (building type %%v) cannot respond to %s\", reflect.TypeOf(_this), dst.Type()))\n"
		if _, err := writer.Write([]byte(fmt.Sprintf(format, methodName))); err != nil {
			panic(err)
		}
	} else {
		format := "\tpanic(fmt.Errorf(\"BUG: %%v cannot respond to %s\", reflect.TypeOf(_this)))\n"
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
