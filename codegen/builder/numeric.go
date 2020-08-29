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

var numericDefinitions = [...]string{
	"int",
	"uint",
	"float",
	// "bigInt",
	// "pBigInt",
	"bigFloat",
	"pBigFloat",
	"bigDecimalFloat",
	"pBigDecimalFloat",
	"decimalFloat",
}

func generateNumericBuilders(writer io.Writer) {
	tmpl, err := template.New("").Parse(numericTemplate)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0, len(numericDefinitions))
	for _, k := range numericDefinitions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, typeName := range keys {
		data := &numericData{
			TypeName:  typeName,
			TypeTitle: strings.Title(typeName),
		}
		if err = tmpl.Execute(writer, data); err != nil {
			panic(fmt.Errorf("could not write template: %v", err))
		}
	}
}

type numericData struct {
	TypeName  string
	TypeTitle string
}

var numericTemplate = `
// {{.TypeName}}

type {{.TypeName}}Builder struct {
	// Template Data
	dstType reflect.Type
}

func new{{.TypeTitle}}Builder(dstType reflect.Type) ObjectBuilder {
	return &{{.TypeName}}Builder{
		dstType: dstType,
	}
}
func (_this *{{.TypeName}}Builder) String() string {
	return fmt.Sprintf("%v", reflect.TypeOf(_this))
}
func (_this *{{.TypeName}}Builder) badEvent(name string, args ...interface{}) {
	PanicBadEventWithType(_this, _this.dstType, name, args...)
}
func (_this *{{.TypeName}}Builder) InitTemplate(_ *Session) {}
func (_this *{{.TypeName}}Builder) NewInstance(_ *RootBuilder, _ ObjectBuilder, _ *options.BuilderOptions) ObjectBuilder {
	return _this
}
func (_this *{{.TypeName}}Builder) SetParent(_ ObjectBuilder) {}
func (_this *{{.TypeName}}Builder) BuildFromBool(_ bool, _ reflect.Value) {
	_this.badEvent("BuildFromBool")
}
func (_this *{{.TypeName}}Builder) BuildFromInt(value int64, dst reflect.Value) {
	set{{.TypeTitle}}FromInt(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromUint(value uint64, dst reflect.Value) {
	set{{.TypeTitle}}FromUint(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromBigInt(value *big.Int, dst reflect.Value) {
	set{{.TypeTitle}}FromBigInt(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromFloat(value float64, dst reflect.Value) {
	set{{.TypeTitle}}FromFloat(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	set{{.TypeTitle}}FromDecimalFloat(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromBigFloat(value *big.Float, dst reflect.Value) {
	set{{.TypeTitle}}FromBigFloat(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	set{{.TypeTitle}}FromBigDecimalFloat(value, dst)
}
func (_this *{{.TypeName}}Builder) BuildFromUUID(_ []byte, _ reflect.Value) {
	_this.badEvent("BuildFromUUID")
}
func (_this *{{.TypeName}}Builder) BuildFromArray(arrayType events.ArrayType, value []byte, dst reflect.Value) {
	_this.badEvent("TypedArray(%v)", arrayType)
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
