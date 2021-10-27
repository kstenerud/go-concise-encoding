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
	"fmt"
	"math/big"

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

type BeginDocumentRule struct{}

func (_this *BeginDocumentRule) String() string               { return "Begin Document Rule" }
func (_this *BeginDocumentRule) OnBeginDocument(ctx *Context) { ctx.ChangeRule(&versionRule) }

// =============================================================================

type EndDocumentRule struct{}

func (_this *EndDocumentRule) String() string             { return "End Document Rule" }
func (_this *EndDocumentRule) OnEndDocument(ctx *Context) { ctx.EndDocument() }

// =============================================================================

type TerminalRule struct{}

func (_this *TerminalRule) String() string { return "Terminal Rule" }

// =============================================================================

type VersionRule struct{}

func (_this *VersionRule) String() string { return "Version Rule" }
func (_this *VersionRule) OnVersion(ctx *Context, version uint64) {
	if version != ctx.ExpectedVersion {
		panic(fmt.Errorf("expected version %v but got version %v", ctx.ExpectedVersion, version))
	}
	ctx.ChangeRule(&topLevelRule)
}

// =============================================================================

type TopLevelRule struct{}

func (_this *TopLevelRule) String() string                              { return "Top Level Rule" }
func (_this *TopLevelRule) switchEndDocument(ctx *Context)              { ctx.ChangeRule(&endDocumentRule) }
func (_this *TopLevelRule) OnKeyableObject(ctx *Context, _ DataType)    { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnNonKeyableObject(ctx *Context, _ DataType) { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnPadding(ctx *Context)                    { /* Nothing to do */ }
func (_this *TopLevelRule) OnComment(ctx *Context)                    { /* Nothing to do */ }
func (_this *TopLevelRule) OnNil(ctx *Context)                        { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnInt(ctx *Context, value int64)           { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnPositiveInt(ctx *Context, value uint64)  { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnBigInt(ctx *Context, value *big.Int)     { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnFloat(ctx *Context, value float64)       { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnBigFloat(ctx *Context, value *big.Float) { _this.switchEndDocument(ctx) }
func (_this *TopLevelRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnList(ctx *Context)                      { ctx.BeginList() }
func (_this *TopLevelRule) OnMap(ctx *Context)                       { ctx.BeginMap() }
func (_this *TopLevelRule) OnMarkup(ctx *Context, identifier []byte) { ctx.BeginMarkup(identifier) }
func (_this *TopLevelRule) OnNode(ctx *Context)                      { ctx.BeginNode() }
func (_this *TopLevelRule) OnEdge(ctx *Context)                      { ctx.BeginEdge() }
func (_this *TopLevelRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
func (_this *TopLevelRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *TopLevelRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type NARule struct{}

func (_this *NARule) String() string                                          { return "NA Rule" }
func (_this *NARule) OnKeyableObject(ctx *Context, _ DataType)                { ctx.EndContainer() }
func (_this *NARule) OnNonKeyableObject(ctx *Context, _ DataType)             { ctx.EndContainer() }
func (_this *NARule) OnChildContainerEnded(ctx *Context, _ DataType)          { ctx.EndContainer() }
func (_this *NARule) OnPadding(ctx *Context)                                  { /* Nothing to do */ }
func (_this *NARule) OnNil(ctx *Context)                                      { ctx.EndContainer() }
func (_this *NARule) OnInt(ctx *Context, value int64)                         { ctx.EndContainer() }
func (_this *NARule) OnPositiveInt(ctx *Context, value uint64)                { ctx.EndContainer() }
func (_this *NARule) OnBigInt(ctx *Context, value *big.Int)                   { ctx.EndContainer() }
func (_this *NARule) OnFloat(ctx *Context, value float64)                     { ctx.EndContainer() }
func (_this *NARule) OnBigFloat(ctx *Context, value *big.Float)               { ctx.EndContainer() }
func (_this *NARule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) { ctx.EndContainer() }
func (_this *NARule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal)      { ctx.EndContainer() }
func (_this *NARule) OnList(ctx *Context)                                     { ctx.BeginList() }
func (_this *NARule) OnMap(ctx *Context)                                      { ctx.BeginMap() }
func (_this *NARule) OnMarkup(ctx *Context, identifier []byte)                { ctx.BeginMarkup(identifier) }
func (_this *NARule) OnNode(ctx *Context)                                     { ctx.BeginNode() }
func (_this *NARule) OnEdge(ctx *Context)                                     { ctx.BeginEdge() }
func (_this *NARule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.EndContainer()
}
func (_this *NARule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	ctx.EndContainer()
}
func (_this *NARule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type ConstantKeyableRule struct{}

func (_this *ConstantKeyableRule) String() string         { return "Keyable Constant Rule" }
func (_this *ConstantKeyableRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *ConstantKeyableRule) OnKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
	ctx.MarkObject(objType)
}
func (_this *ConstantKeyableRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArray(ctx, arrayType, elementCount, data)
	ctx.MarkObject(dataType)
}
func (_this *ConstantKeyableRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnStringlikeArray(ctx, arrayType, data)
	ctx.MarkObject(dataType)
}
func (_this *ConstantKeyableRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable("constant (keyable)", arrayType)
}
func (_this *ConstantKeyableRule) OnChildContainerEnded(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, objType)
}

// =============================================================================

type ConstantAnyTypeRule struct{}

func (_this *ConstantAnyTypeRule) String() string         { return "Constant Rule" }
func (_this *ConstantAnyTypeRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *ConstantAnyTypeRule) OnNonKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
}
func (_this *ConstantAnyTypeRule) OnNil(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnNil(ctx)
}
func (_this *ConstantAnyTypeRule) OnKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
}
func (_this *ConstantAnyTypeRule) OnList(ctx *Context) {
	ctx.ParentRule().OnList(ctx)
}
func (_this *ConstantAnyTypeRule) OnMap(ctx *Context) {
	ctx.ParentRule().OnMap(ctx)
}
func (_this *ConstantAnyTypeRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.ParentRule().OnMarkup(ctx, identifier)
}
func (_this *ConstantAnyTypeRule) OnNode(ctx *Context) {
	ctx.ParentRule().OnNode(ctx)
}
func (_this *ConstantAnyTypeRule) OnEdge(ctx *Context) {
	ctx.ParentRule().OnEdge(ctx)
}
func (_this *ConstantAnyTypeRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArray(ctx, arrayType, elementCount, data)
}
func (_this *ConstantAnyTypeRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnStringlikeArray(ctx, arrayType, data)
}
func (_this *ConstantAnyTypeRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.ParentRule().OnArrayBegin(ctx, arrayType)
}
func (_this *ConstantAnyTypeRule) OnChildContainerEnded(ctx *Context, cType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, cType)
}
