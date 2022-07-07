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

// Imposes the structural rules that enforce a well-formed concise encoding
// document.
package rules

import (
	"math/big"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/events"
)

type ListRule struct{}

func (_this *ListRule) String() string { return "List Rule" }

// =============================================================================

type MapKeyRule struct{}

func (_this *MapKeyRule) String() string                                 { return "Map Key Rule" }
func (_this *MapKeyRule) switchMapValue(ctx *Context)                    { ctx.ChangeRule(&mapValueRule) }
func (_this *MapKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier, AllowKeyable)
}
func (_this *MapKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMapValue(ctx)
}
func (_this *MapKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable("map key", arrayType, elementCount, data)
	_this.switchMapValue(ctx)
}
func (_this *MapKeyRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlikeKeyable("map key", arrayType, data)
	_this.switchMapValue(ctx)
}
func (_this *MapKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable("map key", arrayType)
}

// =============================================================================

type MapValueRule struct{}

func (_this *MapValueRule) String() string                                 { return "Map Value Rule" }
func (_this *MapValueRule) switchMapKey(ctx *Context)                      { ctx.ChangeRule(&mapKeyRule) }
func (_this *MapValueRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNull(ctx *Context)                            { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNonKeyableObject(ctx *Context, _ DataType)    { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.switchMapKey(ctx)
}
func (_this *MapValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchMapKey(ctx)
}
func (_this *MapValueRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchMapKey(ctx)
}

// =============================================================================

type StructTemplateRule struct{}

func (_this *StructTemplateRule) String() string { return "Struct Template Rule" }
func (_this *StructTemplateRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable("struct template", arrayType, elementCount, data)
}
func (_this *StructTemplateRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlikeKeyable("struct template", arrayType, data)
}
func (_this *StructTemplateRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable("struct template", arrayType)
}
func (_this *StructTemplateRule) OnEnd(ctx *Context) {
	ctx.EndContainer(false)
}

// =============================================================================

type StructInstanceRule struct{}

func (_this *StructInstanceRule) String() string { return "Struct Instance Rule" }

// =============================================================================

type EdgeSourceRule struct{}

func (_this *EdgeSourceRule) String() string { return "Edge Source Rule" }
func (_this *EdgeSourceRule) moveToNextRule(ctx *Context) {
	ctx.ChangeRule(&edgeDescriptionRule)
}
func (_this *EdgeSourceRule) OnKeyableObject(ctx *Context, _ DataType)    { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnNonKeyableObject(ctx *Context, _ DataType) { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnInt(ctx *Context, value int64)           { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnPositiveInt(ctx *Context, value uint64)  { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnBigInt(ctx *Context, value *big.Int)     { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnFloat(ctx *Context, value float64)       { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnBigFloat(ctx *Context, value *big.Float) { _this.moveToNextRule(ctx) }
func (_this *EdgeSourceRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowNonNull)
}
func (_this *EdgeSourceRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("edge source", arrayType, AllowNonNull)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("edge source", arrayType, AllowNonNull)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeSourceRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("edge source", arrayType, AllowNonNull)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type EdgeDescriptionRule struct{}

func (_this *EdgeDescriptionRule) String() string { return "Edge Description Rule" }
func (_this *EdgeDescriptionRule) moveToNextRule(ctx *Context) {
	ctx.ChangeRule(&edgeDestinationRule)
}
func (_this *EdgeDescriptionRule) OnKeyableObject(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnNonKeyableObject(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnNull(ctx *Context)             { _this.moveToNextRule(ctx) }
func (_this *EdgeDescriptionRule) OnInt(ctx *Context, value int64) { _this.moveToNextRule(ctx) }
func (_this *EdgeDescriptionRule) OnPositiveInt(ctx *Context, value uint64) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnBigInt(ctx *Context, value *big.Int) { _this.moveToNextRule(ctx) }
func (_this *EdgeDescriptionRule) OnFloat(ctx *Context, value float64)   { _this.moveToNextRule(ctx) }
func (_this *EdgeDescriptionRule) OnBigFloat(ctx *Context, value *big.Float) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowNonNull)
}
func (_this *EdgeDescriptionRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("edge description", arrayType, AllowAny)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("edge description", arrayType, AllowAny)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.moveToNextRule(ctx)
}
func (_this *EdgeDescriptionRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("edge description", arrayType, AllowAny)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type EdgeDestinationRule struct{}

func (_this *EdgeDestinationRule) String() string { return "Edge Destination Rule" }

// =============================================================================

type NodeRule struct{}

func (_this *NodeRule) String() string { return "Node Rule" }
func (_this *NodeRule) moveToNextRule(ctx *Context) {
	ctx.ChangeRule(&listRule)
}
func (_this *NodeRule) OnKeyableObject(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnNonKeyableObject(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnNull(ctx *Context)             { _this.moveToNextRule(ctx) }
func (_this *NodeRule) OnInt(ctx *Context, value int64) { _this.moveToNextRule(ctx) }
func (_this *NodeRule) OnPositiveInt(ctx *Context, value uint64) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnBigInt(ctx *Context, value *big.Int) { _this.moveToNextRule(ctx) }
func (_this *NodeRule) OnFloat(ctx *Context, value float64)   { _this.moveToNextRule(ctx) }
func (_this *NodeRule) OnBigFloat(ctx *Context, value *big.Float) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowNonNull)
}
func (_this *NodeRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("node", arrayType, AllowAny)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("node", arrayType, AllowAny)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.moveToNextRule(ctx)
}
func (_this *NodeRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("node", arrayType, AllowAny)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type AwaitEndRule struct{}

func (_this *AwaitEndRule) String() string { return "Await Container End Rule" }
