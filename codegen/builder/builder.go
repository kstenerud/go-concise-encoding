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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/kstenerud/go-concise-encoding/codegen/standard"
)

var typesAndInstances = map[string]string{
	"uint8":  "uint8(0)",
	"uint16": "uint16(0)",
	// "uint32":  "uint32(0)",
	// "uint64":  "uint64(0)",
	// "int8":    "int8(0)",
	// "int16":   "int16(0)",
	// "int32":   "int32(0)",
	// "int64":   "int64(0)",
	// "float16": "float16(0)",
	// "float32": "float32(0)",
	// "float64": "float64(0)",
	// "bool":    "true",
	// "UUID":    "???",
}

func GenerateCode(projectDir string) {
	tmpl, err := template.New("").Parse(templateContents)
	if err != nil {
		panic(err)
	}

	generatedFilePath := filepath.Join(projectDir, "builder/builder_generated.go")
	writer, err := os.Create(generatedFilePath)
	if err != nil {
		panic(fmt.Errorf("could not open %s: %v", generatedFilePath, err))
	}

	_, err = writer.WriteString(codeHeader)
	if err != nil {
		panic(fmt.Errorf("could not write to %s: %v", generatedFilePath, err))
	}

	keys := make([]string, 0, len(typesAndInstances))
	for k := range typesAndInstances {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, typeName := range keys {
		typeInstance := typesAndInstances[typeName]
		data := &generatorData{
			TypeName:     typeName,
			TypeTitle:    strings.Title(typeName),
			TypeInstance: typeInstance,
		}
		if err = tmpl.Execute(writer, data); err != nil {
			panic(fmt.Errorf("could not write template to %s: %v", generatedFilePath, err))
		}
	}
}

type generatorData struct {
	TypeName     string
	TypeTitle    string
	TypeInstance string
}

var codeHeader = standard.Header + `package builder

import (
	"math/big"
	"net/url"
	"reflect"
	"time"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)
`

var templateContents = `
// {{.TypeName}}

type {{.TypeName}}ArrayBuilder struct{}

func new{{.TypeTitle}}ArrayBuilder() ObjectBuilder       { return &{{.TypeName}}ArrayBuilder{} }
func (_this *{{.TypeName}}ArrayBuilder) String() string { return nameOf(_this) }
func (_this *{{.TypeName}}ArrayBuilder) badEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf({{.TypeInstance}}), name, args...)
}
func (_this *{{.TypeName}}ArrayBuilder) InitTemplate(_ *Session) {}
func (_this *{{.TypeName}}ArrayBuilder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *{{.TypeName}}ArrayBuilder) SetParent(_ ObjectBuilder) {}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromNil(_ reflect.Value) {
	_this.badEvent("BuildFromNil")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.badEvent("BuildFromBool")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.badEvent("BuildFromInt")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.badEvent("BuildFromUint")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.badEvent("BuildFromBigInt")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.badEvent("BuildFromFloat")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.badEvent("BuildFromBigFloat")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.badEvent("BuildFromDecimalFloat")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.badEvent("BuildFromBigDecimalFloat")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromUUID")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromString")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromVerbatimString(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromVerbatimString")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromURI(_ *url.URL, _ reflect.Value) {
	_this.badEvent("BuildFromURI")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromCustomBinary(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomBinary")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromCustomText(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromCustomText")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromTime")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromCompactTime")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildBeginList() {
	_this.badEvent("BuildBeginList")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildBeginMap() {
	_this.badEvent("BuildBeginMap")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildEndContainer() {
	_this.badEvent("BuildEndContainer")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildBeginMarker(_ interface{}) {
	_this.badEvent("BuildBeginMarker")
}
func (_this *{{.TypeName}}ArrayBuilder) BuildFromReference(_ interface{}) {
	_this.badEvent("BuildFromReference")
}
func (_this *{{.TypeName}}ArrayBuilder) PrepareForListContents() {
	_this.badEvent("PrepareForListContents")
}
func (_this *{{.TypeName}}ArrayBuilder) PrepareForMapContents() {
	_this.badEvent("PrepareForMapContents")
}
func (_this *{{.TypeName}}ArrayBuilder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.badEvent("NotifyChildContainerFinished")
}
`
