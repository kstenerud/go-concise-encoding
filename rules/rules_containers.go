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
	"github.com/kstenerud/go-concise-encoding/events"
)

type ListRule struct{}

func (_this *ListRule) String() string                                 { return "List Rule" }
func (_this *ListRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *ListRule) OnNA(ctx *Context)                              { ctx.BeginNA() }
func (_this *ListRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *ListRule) OnKeyableObject(ctx *Context, _ string)         { /* Nothing to do */ }
func (_this *ListRule) OnNonKeyableObject(ctx *Context, _ string)      { /* Nothing to do */ }
func (_this *ListRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *ListRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *ListRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *ListRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *ListRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *ListRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *ListRule) OnMarker(ctx *Context, identifier []byte)       { ctx.BeginMarkerAnyType(identifier) }
func (_this *ListRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
}
func (_this *ListRule) OnRIDReference(ctx *Context) {
	ctx.BeginRIDReference()
}
func (_this *ListRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *ListRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *ListRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *ListRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MapKeyRule struct{}

func (_this *MapKeyRule) String() string                                 { return "Map Key Rule" }
func (_this *MapKeyRule) switchMapValue(ctx *Context)                    { ctx.ChangeRule(&mapValueRule) }
func (_this *MapKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapKeyRule) OnKeyableObject(ctx *Context, _ string)         { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *MapKeyRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *MapKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier)
}
func (_this *MapKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMapValue(ctx)
}
func (_this *MapKeyRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantKeyable(name)
}
func (_this *MapKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	_this.switchMapValue(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MapKeyRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlikeKeyable(arrayType, data)
	_this.switchMapValue(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MapKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MapValueRule struct{}

func (_this *MapValueRule) String() string                                 { return "Map Value Rule" }
func (_this *MapValueRule) switchMapKey(ctx *Context)                      { ctx.ChangeRule(&mapKeyRule) }
func (_this *MapValueRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNA(ctx *Context) {
	_this.switchMapKey(ctx)
	ctx.BeginNA()
}
func (_this *MapValueRule) OnPadding(ctx *Context)                    { /* Nothing to do */ }
func (_this *MapValueRule) OnKeyableObject(ctx *Context, _ string)    { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNonKeyableObject(ctx *Context, _ string) { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnList(ctx *Context)                       { ctx.BeginList() }
func (_this *MapValueRule) OnMap(ctx *Context)                        { ctx.BeginMap() }
func (_this *MapValueRule) OnMarkup(ctx *Context, identifier []byte)  { ctx.BeginMarkup(identifier) }
func (_this *MapValueRule) OnComment(ctx *Context)                    { ctx.BeginComment() }
func (_this *MapValueRule) OnRelationship(ctx *Context)               { ctx.BeginRelationship() }
func (_this *MapValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier)
}
func (_this *MapValueRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.switchMapKey(ctx)
}
func (_this *MapValueRule) OnRIDReference(ctx *Context) {
	ctx.BeginRIDReference()
}
func (_this *MapValueRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *MapValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchMapKey(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MapValueRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchMapKey(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MapValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupKeyRule struct{}

func (_this *MarkupKeyRule) String() string                 { return "Markup Attribute Key Rule" }
func (_this *MarkupKeyRule) switchMarkupValue(ctx *Context) { ctx.ChangeRule(&markupValueRule) }
func (_this *MarkupKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnPadding(ctx *Context)                 { /* Nothing to do */ }
func (_this *MarkupKeyRule) OnKeyableObject(ctx *Context, _ string) { _this.switchMarkupValue(ctx) }
func (_this *MarkupKeyRule) OnComment(ctx *Context)                 { ctx.BeginComment() }
func (_this *MarkupKeyRule) OnEnd(ctx *Context)                     { ctx.ChangeRule(&markupContentsRule) }
func (_this *MarkupKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier)
}
func (_this *MarkupKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantKeyable(name)
}
func (_this *MarkupKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	_this.switchMarkupValue(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MarkupKeyRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlikeKeyable(arrayType, data)
	_this.switchMarkupValue(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MarkupKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MarkupValueRule struct{}

func (_this *MarkupValueRule) String() string               { return "Markup Attribute Value Rule" }
func (_this *MarkupValueRule) switchMarkupKey(ctx *Context) { ctx.ChangeRule(&markupKeyRule) }
func (_this *MarkupValueRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnNA(ctx *Context) {
	_this.switchMarkupKey(ctx)
	ctx.BeginNA()
}
func (_this *MarkupValueRule) OnPadding(ctx *Context)                    { /* Nothing to do */ }
func (_this *MarkupValueRule) OnKeyableObject(ctx *Context, _ string)    { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnNonKeyableObject(ctx *Context, _ string) { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnList(ctx *Context)                       { ctx.BeginList() }
func (_this *MarkupValueRule) OnMap(ctx *Context)                        { ctx.BeginMap() }
func (_this *MarkupValueRule) OnMarkup(ctx *Context, identifier []byte)  { ctx.BeginMarkup(identifier) }
func (_this *MarkupValueRule) OnComment(ctx *Context)                    { ctx.BeginComment() }
func (_this *MarkupValueRule) OnRelationship(ctx *Context)               { ctx.BeginRelationship() }
func (_this *MarkupValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier)
}
func (_this *MarkupValueRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnRIDReference(ctx *Context) {
	ctx.BeginRIDReference()
}
func (_this *MarkupValueRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *MarkupValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchMarkupKey(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MarkupValueRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchMarkupKey(ctx)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *MarkupValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupContentsRule struct{}

func (_this *MarkupContentsRule) String() string                                 { return "Markup Contents Rule" }
func (_this *MarkupContentsRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnMarkup(ctx *Context, identifier []byte) {
	ctx.BeginMarkup(identifier)
}
func (_this *MarkupContentsRule) OnComment(ctx *Context) { ctx.BeginComment() }
func (_this *MarkupContentsRule) OnEnd(ctx *Context)     { ctx.EndContainer() }
func (_this *MarkupContentsRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayMarkupContents(arrayType, elementCount, data)
}
func (_this *MarkupContentsRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayMarkupContentsString(arrayType, data)
}
func (_this *MarkupContentsRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayString(arrayType)
}

// =============================================================================

type CommentRule struct{}

func (_this *CommentRule) String() string                                 { return "Comment Rule" }
func (_this *CommentRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *CommentRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *CommentRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *CommentRule) OnEnd(ctx *Context)                             { ctx.UnstackRule() }
func (_this *CommentRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayComment(arrayType, elementCount, data)
}
func (_this *CommentRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayCommentString(arrayType, data)
}
func (_this *CommentRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayComment(arrayType)
}

// =============================================================================

type SubjectRule struct{}

func (_this *SubjectRule) String() string                                 { return "Subject Rule" }
func (_this *SubjectRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *SubjectRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *SubjectRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *SubjectRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *SubjectRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *SubjectRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *SubjectRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier)
}
func (_this *SubjectRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
}
func (_this *SubjectRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *SubjectRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *SubjectRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	ctx.BeginPotentialRIDCat(arrayType)
}
func (_this *SubjectRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginPotentialRIDCat(arrayType)
	ctx.BeginArrayAnyType(arrayType)
}
