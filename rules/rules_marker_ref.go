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
	"github.com/kstenerud/go-concise-encoding/ce/events"
)

type MarkedObjectKeyableRule struct{}

func (_this *MarkedObjectKeyableRule) String() string { return "Marked Keyable Object Rule" }
func (_this *MarkedObjectKeyableRule) OnKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
	ctx.MarkObject(objType)
}
func (_this *MarkedObjectKeyableRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArray(ctx, arrayType, elementCount, data)
	ctx.MarkObject(dataType)
}
func (_this *MarkedObjectKeyableRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnStringlikeArray(ctx, arrayType, data)
	ctx.MarkObject(dataType)
}
func (_this *MarkedObjectKeyableRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable("marked object (keyable)", arrayType)
}
func (_this *MarkedObjectKeyableRule) OnChildContainerEnded(ctx *Context, dataType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, dataType)
}

// =============================================================================

type MarkedObjectAnyTypeRule struct{}

func (_this *MarkedObjectAnyTypeRule) String() string { return "Marked Object Rule" }
func (_this *MarkedObjectAnyTypeRule) OnNull(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnNull(ctx)
	ctx.MarkObject(DataTypeNull)
}
func (_this *MarkedObjectAnyTypeRule) OnNonKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
	ctx.MarkObject(objType)
}
func (_this *MarkedObjectAnyTypeRule) OnKeyableObject(ctx *Context, objType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx, objType)
	ctx.MarkObject(objType)
}
func (_this *MarkedObjectAnyTypeRule) OnList(ctx *Context) {
	ctx.ParentRule().OnList(ctx)
}
func (_this *MarkedObjectAnyTypeRule) OnMap(ctx *Context) {
	ctx.ParentRule().OnMap(ctx)
}
func (_this *MarkedObjectAnyTypeRule) OnStructInstance(ctx *Context, identifier []byte) {
	ctx.ParentRule().OnStructInstance(ctx, identifier)
}
func (_this *MarkedObjectAnyTypeRule) OnNode(ctx *Context) { ctx.BeginNode() }
func (_this *MarkedObjectAnyTypeRule) OnEdge(ctx *Context) { ctx.BeginEdge() }
func (_this *MarkedObjectAnyTypeRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("marked object", arrayType, AllowMarkable)
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArray(ctx, arrayType, elementCount, data)
	switch arrayType {
	case events.ArrayTypeString, events.ArrayTypeResourceID:
		ctx.MarkObject(dataType)
	default:
		ctx.MarkObject(dataType)
	}
}
func (_this *MarkedObjectAnyTypeRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("marked object", arrayType, AllowMarkable)
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnStringlikeArray(ctx, arrayType, data)
	switch arrayType {
	case events.ArrayTypeString:
		ctx.MarkObject(dataType)
	default:
		ctx.MarkObject(dataType)
	}
}
func (_this *MarkedObjectAnyTypeRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("marked object", arrayType, AllowMarkable)
	ctx.ParentRule().OnArrayBegin(ctx, arrayType)
}
func (_this *MarkedObjectAnyTypeRule) OnChildContainerEnded(ctx *Context, cType DataType) {
	ctx.MarkObject(cType)
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, cType)
}
