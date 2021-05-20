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

	"github.com/kstenerud/go-concise-encoding/events"
)

type MarkedObjectKeyableRule struct{}

func (_this *MarkedObjectKeyableRule) String() string         { return "Marked Keyable Object Rule" }
func (_this *MarkedObjectKeyableRule) OnPadding(ctx *Context) { /* Nothing to do */ }
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
	ctx.BeginArrayKeyable(arrayType)
}
func (_this *MarkedObjectKeyableRule) OnChildContainerEnded(ctx *Context, dataType DataType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, dataType)
}

// =============================================================================

type MarkedObjectAnyTypeRule struct{}

func (_this *MarkedObjectAnyTypeRule) String() string         { return "Marked Object Rule" }
func (_this *MarkedObjectAnyTypeRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *MarkedObjectAnyTypeRule) OnNil(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnNil(ctx)
	ctx.MarkObject(DataTypeNil)
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
func (_this *MarkedObjectAnyTypeRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.ParentRule().OnMarkup(ctx, identifier)
}
func (_this *MarkedObjectAnyTypeRule) OnRelationship(ctx *Context) { ctx.BeginRelationship() }
func (_this *MarkedObjectAnyTypeRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
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
	dataType := arrayTypeToDataType[arrayType]
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnStringlikeArray(ctx, arrayType, data)
	switch arrayType {
	case events.ArrayTypeString, events.ArrayTypeResourceID:
		ctx.MarkObject(dataType)
	default:
		ctx.MarkObject(dataType)
	}
}
func (_this *MarkedObjectAnyTypeRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.ParentRule().OnArrayBegin(ctx, arrayType)
}
func (_this *MarkedObjectAnyTypeRule) OnChildContainerEnded(ctx *Context, cType DataType) {
	ctx.MarkObject(cType)
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, cType)
}

// =============================================================================

type RIDReferenceRule struct{}

func (_this *RIDReferenceRule) String() string         { return "Resource ID Reference Rule" }
func (_this *RIDReferenceRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *RIDReferenceRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	dataType := arrayTypeToDataType[arrayType]
	switch arrayType {
	case events.ArrayTypeResourceIDConcat:
		ctx.ValidateResourceID(data)
		ctx.ChangeRule(&ridCatRule)
	case events.ArrayTypeResourceID:
		ctx.ValidateResourceID(data)
		ctx.UnstackRule()
		ctx.CurrentEntry.Rule.OnChildContainerEnded(ctx, dataType)
	default:
		panic(fmt.Errorf("Reference Resource ID cannot be type %v", arrayType))
	}
}
func (_this *RIDReferenceRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayRIDReference(arrayType)
}
func (_this *RIDReferenceRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	// TODO: Make this properly
	_this.OnArray(ctx, arrayType, uint64(len(data)), []byte(data))
}
func (_this *RIDReferenceRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	// Toss out the result because it's a resource ID
}
