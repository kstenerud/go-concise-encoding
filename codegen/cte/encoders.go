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

package cte

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

var (
	Begin  = "Begin(ctx *EncoderContext)"
	End    = "End(ctx *EncoderContext)"
	Child  = "ChildContainerFinished(ctx *EncoderContext)"
	Bool   = "EncodeBool(ctx *EncoderContext, value bool)"
	True   = "EncodeTrue(ctx *EncoderContext)"
	False  = "EncodeFalse(ctx *EncoderContext)"
	PInt   = "EncodePositiveInt(ctx *EncoderContext, value uint64)"
	NInt   = "EncodeNegativeInt(ctx *EncoderContext, value uint64)"
	Int    = "EncodeInt(ctx *EncoderContext, value int64)"
	BInt   = "EncodeBigInt(ctx *EncoderContext, value *big.Int)"
	Float  = "EncodeFloat(ctx *EncoderContext, value float64)"
	BFloat = "EncodeBigFloat(ctx *EncoderContext, value *big.Float)"
	DFloat = "EncodeDecimalFloat(ctx *EncoderContext, value compact_float.DFloat)"
	BDF    = "EncodeBigDecimalFloat(ctx *EncoderContext, value *apd.Decimal)"
	Nan    = "EncodeNan(ctx *EncoderContext, signaling bool)"
	Time   = "EncodeTime(ctx *EncoderContext, value time.Time)"
	CTime  = "EncodeCompactTime(ctx *EncoderContext, value compact_time.Time)"
	UUID   = "EncodeUUID(ctx *EncoderContext, value []byte)"
	List   = "BeginList(ctx *EncoderContext)"
	Map    = "BeginMap(ctx *EncoderContext)"
	Markup = "BeginMarkup(ctx *EncoderContext)"
	Meta   = "BeginMetadata(ctx *EncoderContext)"
	Cmt    = "BeginComment(ctx *EncoderContext)"
	Marker = "BeginMarker(ctx *EncoderContext)"
	Ref    = "BeginReference(ctx *EncoderContext)"
	Cat    = "BeginConcatenate(ctx *EncoderContext)"
	Const  = "BeginConstant(ctx *EncoderContext, name []byte, explicitValue bool)"
	NA     = "BeginNA(ctx *EncoderContext)"
	Arr    = "EncodeArray(ctx *EncoderContext, arrayType events.ArrayType, elementCount uint64, data []uint8)"
	Str    = "EncodeStringlikeArray(ctx *EncoderContext, arrayType events.ArrayType, data string)"
	BArr   = "BeginArray(ctx *EncoderContext, arrayType events.ArrayType)"
	Chunk  = "BeginArrayChunk(ctx *EncoderContext, length uint64, moreChunksFollow bool)"
	Data   = "EncodeArrayData(ctx *EncoderContext, data []byte)"

	allMethods = []string{Begin, End, Child, Bool, True, False, PInt, NInt, Int,
		BInt, Float, BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map, Markup,
		Meta, Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr, Chunk, Data}
)

type Encoder struct {
	Name    string
	Methods []string
}

var encoders = []Encoder{
	{
		Name: "topLevel",
		Methods: []string{Child, Bool, True, False, PInt, NInt, Int, BInt, Float,
			BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map, Markup, Meta,
			Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "na",
		Methods: []string{Begin, Bool, True, False, PInt, NInt, Int, BInt, Float,
			BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map, Markup,
			Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "constant",
		Methods: []string{Begin, Bool, True, False, PInt, NInt, Int, BInt, Float,
			BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map, Markup,
			Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "postInvisible",
		Methods: []string{Bool, True, False, PInt, NInt, Int, BInt, Float,
			BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map, Markup, Meta,
			Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "list",
		Methods: []string{Child, Begin, End, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map,
			Markup, Meta, Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "mapKey",
		Methods: []string{Child, Begin, End, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Time, CTime, UUID, Meta, Cmt,
			Marker, Ref, Cat, Const, Arr, Str, BArr},
	},
	{
		Name: "mapValue",
		Methods: []string{Child, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map,
			Markup, Meta, Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "metadataKey",
		Methods: []string{Child, Begin, End, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Time, CTime, UUID, Meta, Cmt,
			Marker, Ref, Cat, Const, Arr, Str, BArr},
	},
	{
		Name: "metadataValue",
		Methods: []string{Child, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map,
			Markup, Meta, Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name: "markupName",
		Methods: []string{Child, Begin, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Time, CTime, UUID,
			Marker, Ref, Cat, Const, Arr, Str, BArr},
	},
	{
		Name: "markupKey",
		Methods: []string{Child, End, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Time, CTime, UUID, Meta, Cmt,
			Marker, Ref, Cat, Const, Arr, Str, BArr},
	},
	{
		Name: "markupValue",
		Methods: []string{Child, Bool, True, False, PInt, NInt, Int,
			BInt, Float, BFloat, DFloat, BDF, Nan, Time, CTime, UUID, List, Map,
			Markup, Meta, Cmt, Marker, Ref, Cat, Const, NA, Arr, Str, BArr},
	},
	{
		Name:    "markupContents",
		Methods: []string{Child, Begin, End, Markup, Cmt, Arr, Str, BArr},
	},
	{
		Name:    "comment",
		Methods: []string{Child, Begin, End, Cmt, Arr, Str, BArr},
	},
	{
		Name:    "array",
		Methods: []string{Chunk, Data},
	},
}

func GenerateCode(projectDir string) {
	generatedFilePath := filepath.Join(projectDir, "cte/generated-do-not-edit.go")
	writer, err := os.Create(generatedFilePath)
	if err != nil {
		panic(fmt.Errorf("could not open %s: %v", generatedFilePath, err))
	}
	defer writer.Close()

	_, err = writer.WriteString(codeHeader)
	if err != nil {
		panic(fmt.Errorf("could not write to %s: %v", generatedFilePath, err))
	}

	generateBadEventMethods(writer)
}

var codeHeader = standard.Header + `package cte

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
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

func openMethod(builder Encoder, methodSignature string, writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("func (_this *%sEncoder) %s {\n", builder.Name, methodSignature))); err != nil {
		panic(err)
	}
}

func closeMethod(writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintf("}\n"))); err != nil {
		panic(err)
	}
}

func generateBadEventMethod(builder Encoder, methodSignature string, writer io.Writer) {
	methodName := methodSignature[:strings.Index(methodSignature, "(")]
	openMethod(builder, methodSignature, writer)

	format := "\tpanic(fmt.Errorf(\"BUG: %%v cannot respond to %s\", reflect.TypeOf(_this)))\n"
	if _, err := writer.Write([]byte(fmt.Sprintf(format, methodName))); err != nil {
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
	for _, builder := range encoders {
		for _, methodSignature := range allMethods {
			if !contains(methodSignature, builder.Methods) {
				generateBadEventMethod(builder, methodSignature, writer)
			}
		}
	}
}
