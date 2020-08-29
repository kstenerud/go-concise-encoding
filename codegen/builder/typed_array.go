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
	"sort"
	"strings"
	"text/template"
)

var typedArrayDefinitions = map[string]string{
	"string":      `""`,
	"uint8Array":  "uint8(0)",
	"uint16Array": "uint16(0)",
	// TODO: Remaining typed arrays
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

func generateTypedArrayBuilders(writer io.Writer) {
	tmpl, err := template.New("").Parse(typedArrayTemplate)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0, len(typedArrayDefinitions))
	for k := range typedArrayDefinitions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, typeName := range keys {
		typeInstance := typedArrayDefinitions[typeName]
		data := &typedArrayData{
			TypeName:     typeName,
			TypeTitle:    strings.Title(typeName),
			TypeInstance: typeInstance,
		}
		if err = tmpl.Execute(writer, data); err != nil {
			panic(fmt.Errorf("could not write template: %v", err))
		}
	}
}

type typedArrayData struct {
	TypeName     string
	TypeTitle    string
	TypeInstance string
}

var typedArrayTemplate = `
// {{.TypeName}}

type {{.TypeName}}Builder struct{}

func new{{.TypeTitle}}Builder() ObjectBuilder       { return &{{.TypeName}}Builder{} }
func (_this *{{.TypeName}}Builder) String() string { return nameOf(_this) }
func (_this *{{.TypeName}}Builder) badEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, reflect.TypeOf({{.TypeInstance}}), name, args...)
}
func (_this *{{.TypeName}}Builder) InitTemplate(_ *Session) {}
func (_this *{{.TypeName}}Builder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *{{.TypeName}}Builder) SetParent(_ ObjectBuilder) {}
func (_this *{{.TypeName}}Builder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.badEvent("BuildFromBool")
}
func (_this *{{.TypeName}}Builder) BuildFromInt(_ int64, _ reflect.Value) {
	_this.badEvent("BuildFromInt")
}
func (_this *{{.TypeName}}Builder) BuildFromUint(_ uint64, _ reflect.Value) {
	_this.badEvent("BuildFromUint")
}
func (_this *{{.TypeName}}Builder) BuildFromBigInt(_ *big.Int, _ reflect.Value) {
	_this.badEvent("BuildFromBigInt")
}
func (_this *{{.TypeName}}Builder) BuildFromFloat(_ float64, _ reflect.Value) {
	_this.badEvent("BuildFromFloat")
}
func (_this *{{.TypeName}}Builder) BuildFromBigFloat(_ *big.Float, _ reflect.Value) {
	_this.badEvent("BuildFromBigFloat")
}
func (_this *{{.TypeName}}Builder) BuildFromDecimalFloat(_ compact_float.DFloat, _ reflect.Value) {
	_this.badEvent("BuildFromDecimalFloat")
}
func (_this *{{.TypeName}}Builder) BuildFromBigDecimalFloat(_ *apd.Decimal, _ reflect.Value) {
	_this.badEvent("BuildFromBigDecimalFloat")
}
func (_this *{{.TypeName}}Builder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromUUID")
}
func (_this *{{.TypeName}}Builder) BuildFromTime(_ time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromTime")
}
func (_this *{{.TypeName}}Builder) BuildFromCompactTime(_ *compact_time.Time, _ reflect.Value) {
	_this.badEvent("BuildFromCompactTime")
}
func (_this *{{.TypeName}}Builder) BuildBeginList() {
	_this.badEvent("BuildBeginList")
}
func (_this *{{.TypeName}}Builder) BuildBeginMap() {
	_this.badEvent("BuildBeginMap")
}
func (_this *{{.TypeName}}Builder) BuildEndContainer() {
	_this.badEvent("BuildEndContainer")
}
func (_this *{{.TypeName}}Builder) BuildBeginMarker(_ interface{}) {
	_this.badEvent("BuildBeginMarker")
}
func (_this *{{.TypeName}}Builder) BuildFromReference(_ interface{}) {
	_this.badEvent("BuildFromReference")
}
func (_this *{{.TypeName}}Builder) PrepareForListContents() {
	_this.badEvent("PrepareForListContents")
}
func (_this *{{.TypeName}}Builder) PrepareForMapContents() {
	_this.badEvent("PrepareForMapContents")
}
func (_this *{{.TypeName}}Builder) NotifyChildContainerFinished(_ reflect.Value) {
	_this.badEvent("NotifyChildContainerFinished")
}
`
