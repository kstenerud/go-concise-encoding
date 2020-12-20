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
	"path/filepath"
	"strings"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

var (
	Nil            = "BuildFromNil(dst reflect.Value)"
	Bool           = "BuildFromBool(value bool, dst reflect.Value)"
	Int            = "BuildFromInt(value int64, dst reflect.Value)"
	Uint           = "BuildFromUint(value uint64, dst reflect.Value)"
	BigInt         = "BuildFromBigInt(value *big.Int, dst reflect.Value)"
	Float          = "BuildFromFloat(value float64, dst reflect.Value)"
	BigFloat       = "BuildFromBigFloat(value *big.Float, dst reflect.Value)"
	DFloat         = "BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value)"
	BigDFloat      = "BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value)"
	UUID           = "BuildFromUUID(value []byte, dst reflect.Value)"
	Array          = "BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value)"
	Time           = "BuildFromTime(value time.Time, dst reflect.Value)"
	CTime          = "BuildFromCompactTime(value *compact_time.Time, dst reflect.Value)"
	List           = "BuildBeginList()"
	Map            = "BuildBeginMap()"
	End            = "BuildEndContainer()"
	Marker         = "BuildBeginMarker(id interface{})"
	Ref            = "BuildFromReference(id interface{})"
	PrepList       = "PrepareForListContents()"
	PrepMap        = "PrepareForMapContents()"
	NotifyFinished = "NotifyChildContainerFinished(container reflect.Value)"
)

var allMethods = []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Marker, Ref, PrepList, PrepMap, NotifyFinished}

var validMethodsByClass = map[string][]string{
	"array":            []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Marker, Ref, PrepList, NotifyFinished},
	"bigDecimalFloat":  []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"bigFloat":         []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"bigInt":           []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"compactTime":      []string{Array, Time, CTime},
	"custom":           []string{Array},
	"decimalFloat":     []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"direct":           []string{Bool, UUID, Array, Time, CTime},
	"directPtr":        []string{Nil, UUID, Array, Time, CTime},
	"float":            []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"ignore":           []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, Ref},
	"ignoreContainer":  []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Ref, PrepList, PrepMap, NotifyFinished},
	"int":              []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"interface":        []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, PrepList, PrepMap, NotifyFinished},
	"map":              []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Marker, Ref, PrepMap, NotifyFinished},
	"markerID":         []string{Int, Uint, BigInt},
	"markerObject":     []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, PrepList, PrepMap, NotifyFinished},
	"pBigDecimalFloat": []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"pBigFloat":        []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"pBigInt":          []string{Nil, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"pCompactTime":     []string{Nil, Array, Time, CTime},
	"ptr":              []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, PrepList, PrepMap, NotifyFinished},
	"Root":             []string{NotifyFinished},
	"slice":            []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Marker, Ref, PrepList, NotifyFinished},
	"string":           []string{Nil, Array},
	"struct":           []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, End, Marker, Ref, PrepMap, NotifyFinished},
	"time":             []string{Array, Time, CTime},
	"topLevel":         []string{Nil, Bool, Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat, UUID, Array, Time, CTime, List, Map, Marker, NotifyFinished},
	"uint":             []string{Int, Uint, BigInt, Float, BigFloat, DFloat, BigDFloat},
	"uint16Array":      []string{Array},
	"uint8Array":       []string{Array},
}

func GenerateCode(projectDir string) {
	generatedFilePath := filepath.Join(projectDir, "builder/generated-do-not-edit.go")
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

var codeHeader = standard.Header + `package builder

import (
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

func generateBadEventMethod(className string, methodSignature string, writer io.Writer) {
	methodName := methodSignature[:strings.Index(methodSignature, "(")]
	_, err := writer.Write([]byte(fmt.Sprintf(`
func (_this *%sBuilder) %s {
	_this.panicBadEvent("%s")
}
`, className, methodSignature, methodName)))
	if err != nil {
		panic(err)
	}
}

func generateBadEventMethods(writer io.Writer) {
	for className, validMethods := range validMethodsByClass {
		for _, methodSignature := range allMethods {
			if !contains(methodSignature, validMethods) {
				generateBadEventMethod(className, methodSignature, writer)
			}
		}
	}
}
