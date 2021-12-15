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

func (_this *ListRule) String() string                                 { return "List Rule" }
func (_this *ListRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *ListRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *ListRule) OnComment(ctx *Context)                         { /* Nothing to do */ }
func (_this *ListRule) OnNull(ctx *Context)                            { /* Nothing to do */ }
func (_this *ListRule) OnKeyableObject(ctx *Context, _ DataType)       { /* Nothing to do */ }
func (_this *ListRule) OnNonKeyableObject(ctx *Context, _ DataType)    { /* Nothing to do */ }
func (_this *ListRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *ListRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *ListRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *ListRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *ListRule) OnNode(ctx *Context)                            { ctx.BeginNode() }
func (_this *ListRule) OnEdge(ctx *Context)                            { ctx.BeginEdge() }
func (_this *ListRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
func (_this *ListRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
}
func (_this *ListRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
}
func (_this *ListRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
}
func (_this *ListRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MapKeyRule struct{}

func (_this *MapKeyRule) String() string                                 { return "Map Key Rule" }
func (_this *MapKeyRule) switchMapValue(ctx *Context)                    { ctx.ChangeRule(&mapValueRule) }
func (_this *MapKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapKeyRule) OnComment(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapKeyRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
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
func (_this *MapValueRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapValueRule) OnComment(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapValueRule) OnNull(ctx *Context)                            { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNonKeyableObject(ctx *Context, _ DataType)    { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *MapValueRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *MapValueRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *MapValueRule) OnNode(ctx *Context)                            { ctx.BeginNode() }
func (_this *MapValueRule) OnEdge(ctx *Context)                            { ctx.BeginEdge() }
func (_this *MapValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
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
func (_this *MapValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupKeyRule struct{}

func (_this *MarkupKeyRule) String() string                 { return "Markup Attribute Key Rule" }
func (_this *MarkupKeyRule) switchMarkupValue(ctx *Context) { ctx.ChangeRule(&markupValueRule) }
func (_this *MarkupKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnPadding(ctx *Context)                   { /* Nothing to do */ }
func (_this *MarkupKeyRule) OnComment(ctx *Context)                   { /* Nothing to do */ }
func (_this *MarkupKeyRule) OnKeyableObject(ctx *Context, _ DataType) { _this.switchMarkupValue(ctx) }
func (_this *MarkupKeyRule) OnEnd(ctx *Context)                       { ctx.ChangeRule(&markupContentsRule) }
func (_this *MarkupKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier, AllowKeyable)
}
func (_this *MarkupKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable("markup attribute key", arrayType, elementCount, data)
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlikeKeyable("markup attribute key", arrayType, data)
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable("markup attribute key", arrayType)
}

// =============================================================================

type MarkupValueRule struct{}

func (_this *MarkupValueRule) String() string               { return "Markup Attribute Value Rule" }
func (_this *MarkupValueRule) switchMarkupKey(ctx *Context) { ctx.ChangeRule(&markupKeyRule) }
func (_this *MarkupValueRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnPadding(ctx *Context)                   { /* Nothing to do */ }
func (_this *MarkupValueRule) OnComment(ctx *Context)                   { /* Nothing to do */ }
func (_this *MarkupValueRule) OnNull(ctx *Context)                      { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnKeyableObject(ctx *Context, _ DataType) { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnNonKeyableObject(ctx *Context, _ DataType) {
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnList(ctx *Context)                      { ctx.BeginList() }
func (_this *MarkupValueRule) OnMap(ctx *Context)                       { ctx.BeginMap() }
func (_this *MarkupValueRule) OnMarkup(ctx *Context, identifier []byte) { ctx.BeginMarkup(identifier) }
func (_this *MarkupValueRule) OnNode(ctx *Context)                      { ctx.BeginNode() }
func (_this *MarkupValueRule) OnEdge(ctx *Context)                      { ctx.BeginEdge() }
func (_this *MarkupValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
func (_this *MarkupValueRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupContentsRule struct{}

func (_this *MarkupContentsRule) String() string                                 { return "Markup Contents Rule" }
func (_this *MarkupContentsRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnComment(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.BeginMarkup(identifier)
}
func (_this *MarkupContentsRule) OnEnd(ctx *Context) { ctx.EndContainer() }
func (_this *MarkupContentsRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayMarkupContents(arrayType, elementCount, data)
}
func (_this *MarkupContentsRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayMarkupContentsString(arrayType, data)
}
func (_this *MarkupContentsRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayString("markup contents", arrayType)
}

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
func (_this *EdgeSourceRule) OnPadding(ctx *Context)                    { /* Nothing to do */ }
func (_this *EdgeSourceRule) OnComment(ctx *Context)                    { /* Nothing to do */ }
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
func (_this *EdgeSourceRule) OnList(ctx *Context)                      { ctx.BeginList() }
func (_this *EdgeSourceRule) OnMap(ctx *Context)                       { ctx.BeginMap() }
func (_this *EdgeSourceRule) OnMarkup(ctx *Context, identifier []byte) { ctx.BeginMarkup(identifier) }
func (_this *EdgeSourceRule) OnNode(ctx *Context)                      { ctx.BeginNode() }
func (_this *EdgeSourceRule) OnEdge(ctx *Context)                      { ctx.BeginEdge() }
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
func (_this *EdgeDescriptionRule) OnPadding(ctx *Context)          { /* Nothing to do */ }
func (_this *EdgeDescriptionRule) OnComment(ctx *Context)          { /* Nothing to do */ }
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
func (_this *EdgeDescriptionRule) OnList(ctx *Context) { ctx.BeginList() }
func (_this *EdgeDescriptionRule) OnMap(ctx *Context)  { ctx.BeginMap() }
func (_this *EdgeDescriptionRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.BeginMarkup(identifier)
}
func (_this *EdgeDescriptionRule) OnNode(ctx *Context) { ctx.BeginNode() }
func (_this *EdgeDescriptionRule) OnEdge(ctx *Context) { ctx.BeginEdge() }
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

func (_this *EdgeDestinationRule) String() string                                 { return "Edge Destination Rule" }
func (_this *EdgeDestinationRule) end(ctx *Context)                               { ctx.EndContainer() }
func (_this *EdgeDestinationRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.end(ctx) }
func (_this *EdgeDestinationRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *EdgeDestinationRule) OnComment(ctx *Context)                         { /* Nothing to do */ }
func (_this *EdgeDestinationRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.end(ctx) }
func (_this *EdgeDestinationRule) OnNonKeyableObject(ctx *Context, _ DataType)    { _this.end(ctx) }
func (_this *EdgeDestinationRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *EdgeDestinationRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *EdgeDestinationRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.BeginMarkup(identifier)
}
func (_this *EdgeDestinationRule) OnNode(ctx *Context) { ctx.BeginNode() }
func (_this *EdgeDestinationRule) OnEdge(ctx *Context) { ctx.BeginEdge() }
func (_this *EdgeDestinationRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
func (_this *EdgeDestinationRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.end(ctx)
}
func (_this *EdgeDestinationRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("edge destination", arrayType, AllowNonNull)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.end(ctx)
}
func (_this *EdgeDestinationRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("edge destination", arrayType, AllowNonNull)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.end(ctx)
}
func (_this *EdgeDestinationRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("edge destination", arrayType, AllowNonNull)
	ctx.BeginArrayAnyType(arrayType)
}

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
func (_this *NodeRule) OnPadding(ctx *Context)          { /* Nothing to do */ }
func (_this *NodeRule) OnComment(ctx *Context)          { /* Nothing to do */ }
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
func (_this *NodeRule) OnList(ctx *Context) { ctx.BeginList() }
func (_this *NodeRule) OnMap(ctx *Context)  { ctx.BeginMap() }
func (_this *NodeRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.BeginMarkup(identifier)
}
func (_this *NodeRule) OnNode(ctx *Context) { ctx.BeginNode() }
func (_this *NodeRule) OnEdge(ctx *Context) { ctx.BeginEdge() }
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
