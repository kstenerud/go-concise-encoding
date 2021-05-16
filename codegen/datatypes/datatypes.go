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

package datatypes

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type FlagDataTypeWriter struct {
	writer       io.Writer
	typeName     string
	highestValue interface{}
	isFirst      bool
}

func NewFlagDataTypeWriter(writer io.Writer, typeName string, endMarkerValue interface{}) *FlagDataTypeWriter {
	endValue := reflect.ValueOf(endMarkerValue)
	highestValue := reflect.New(reflect.TypeOf(endMarkerValue)).Elem()
	highestValue.SetUint(endValue.Uint() >> 1)
	return &FlagDataTypeWriter{
		writer:       writer,
		typeName:     typeName,
		highestValue: highestValue.Interface(),
		isFirst:      true,
	}
}

func (_this *FlagDataTypeWriter) BeginType() {
	endMarkerValue := reflect.ValueOf(_this.highestValue).Uint() << 1
	var baseType string
	switch {
	case endMarkerValue <= 0x100:
		baseType = "uint8"
	case endMarkerValue <= 0x10000:
		baseType = "uint16"
	case endMarkerValue <= 0x100000000:
		baseType = "uint32"
	default:
		baseType = "uint64"
	}

	if _, err := fmt.Fprintf(_this.writer, "type %v %v\n\nconst (", _this.typeName, baseType); err != nil {
		panic(err)
	}
}

func (_this *FlagDataTypeWriter) AddNamed(valueName interface{}) {
	if _, err := fmt.Fprintf(_this.writer, "\n\t%v", valueName); err != nil {
		panic(err)
	}
	if _this.isFirst {
		_this.isFirst = false
		if _, err := fmt.Fprintf(_this.writer, " %v = 1 << iota", _this.typeName); err != nil {
			panic(err)
		}
	}
}

func (_this *FlagDataTypeWriter) AddCustom(valueName interface{}, value interface{}) {
	if reflect.TypeOf(value) == reflect.TypeOf(uint64(0)) {
		value = fmt.Sprintf("0x%x", value)
	}
	if _, err := fmt.Fprintf(_this.writer, "\n\t%v = %v", valueName, value); err != nil {
		panic(err)
	}
}

func (_this *FlagDataTypeWriter) EndType() {
	if _, err := fmt.Fprintf(_this.writer, "\n)\n"); err != nil {
		panic(err)
	}
}

func (_this *FlagDataTypeWriter) BeginStringer() {
	typeName := _this.typeName
	lowerName := strings.ToLower(typeName)
	fmt.Fprintf(_this.writer, `
func (_this %v) String() string {
	asString := ""
	if _this == 0 {
		asString = %vNames[_this]
	} else {
		isFirst := true
		builder := strings.Builder{}
		for i := %v(1); i <= %v; i <<= 1 {
			if _this&i != 0 {
				if isFirst {
					isFirst = false
				} else {
					builder.WriteString(" | ")
				}
				builder.WriteString(%vNames[i])
			}
		}
		asString = builder.String()
	}
	if asString == "" {
		asString = fmt.Sprintf("%%d", _this)
	}
	return asString
}

var %vNames = map[%v]string{
`, typeName, lowerName, typeName, _this.highestValue, lowerName, lowerName, typeName)
}

func (_this *FlagDataTypeWriter) AddStringer(typeName interface{}) {
	fmt.Fprintf(_this.writer, "\n\t%v: \"%v\",", typeName, typeName)
}

func (_this *FlagDataTypeWriter) EndStringer() {
	fmt.Fprintf(_this.writer, "\n}\n")
}

func FlagToString(valueMap map[interface{}]string, value interface{}) string {
	uintValue := reflect.ValueOf(value).Uint()
	if uintValue == 0 {
		return valueMap[value]
	}

	isFirst := true
	builder := strings.Builder{}
	lookupValue := reflect.New(reflect.TypeOf(value)).Elem()
	for i := 0; i < 64; i++ {
		testValue := uint64(1) << i
		if uintValue&testValue != 0 {
			lookupValue.SetUint(testValue)
			stringValue, ok := valueMap[lookupValue.Interface()]
			if !ok {
				continue
			}

			if isFirst {
				isFirst = false
			} else {
				builder.WriteString(" | ")
			}
			builder.WriteString(stringValue)
		}
	}
	return builder.String()
}
