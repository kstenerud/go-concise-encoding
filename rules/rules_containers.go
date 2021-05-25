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
func (_this *ListRule) OnNil(ctx *Context)                             { /* Nothing to do */ }
func (_this *ListRule) OnKeyableObject(ctx *Context, _ DataType)       { /* Nothing to do */ }
func (_this *ListRule) OnNonKeyableObject(ctx *Context, _ DataType)    { /* Nothing to do */ }
func (_this *ListRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *ListRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *ListRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *ListRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *ListRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *ListRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *ListRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
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
}
func (_this *ListRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
}
func (_this *ListRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type ResourceListRule struct{}

func (_this *ResourceListRule) String() string                                 { return "Resource List Rule" }
func (_this *ResourceListRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *ResourceListRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *ResourceListRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *ResourceListRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *ResourceListRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *ResourceListRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *ResourceListRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowResource)
}
func (_this *ResourceListRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceObject(identifier, AllowResource)
}
func (_this *ResourceListRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *ResourceListRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("resource list", arrayType, AllowResource)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
}
func (_this *ResourceListRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("resource list", arrayType, AllowResource)
	ctx.ValidateFullArrayStringlike(arrayType, data)
}
func (_this *ResourceListRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("resource list", arrayType, AllowResource)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MapKeyRule struct{}

func (_this *MapKeyRule) String() string                                 { return "Map Key Rule" }
func (_this *MapKeyRule) switchMapValue(ctx *Context)                    { ctx.ChangeRule(&mapValueRule) }
func (_this *MapKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapKeyRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapValue(ctx) }
func (_this *MapKeyRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *MapKeyRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *MapKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier, AllowKeyable)
}
func (_this *MapKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMapValue(ctx)
}
func (_this *MapKeyRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantKeyable(name)
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
func (_this *MapValueRule) OnNA(ctx *Context)                              { ctx.BeginNA() }
func (_this *MapValueRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapValueRule) OnNil(ctx *Context)                             { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnNonKeyableObject(ctx *Context, _ DataType)    { _this.switchMapKey(ctx) }
func (_this *MapValueRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *MapValueRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *MapValueRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *MapValueRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *MapValueRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *MapValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
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
func (_this *MarkupKeyRule) OnKeyableObject(ctx *Context, _ DataType) { _this.switchMarkupValue(ctx) }
func (_this *MarkupKeyRule) OnComment(ctx *Context)                   { ctx.BeginComment() }
func (_this *MarkupKeyRule) OnEnd(ctx *Context)                       { ctx.ChangeRule(&markupContentsRule) }
func (_this *MarkupKeyRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerKeyable(identifier, AllowKeyable)
}
func (_this *MarkupKeyRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceKeyable(identifier)
	_this.switchMarkupValue(ctx)
}
func (_this *MarkupKeyRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantKeyable(name)
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
func (_this *MarkupValueRule) OnNA(ctx *Context)                        { ctx.BeginNA() }
func (_this *MarkupValueRule) OnPadding(ctx *Context)                   { /* Nothing to do */ }
func (_this *MarkupValueRule) OnNil(ctx *Context)                       { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnKeyableObject(ctx *Context, _ DataType) { _this.switchMarkupKey(ctx) }
func (_this *MarkupValueRule) OnNonKeyableObject(ctx *Context, _ DataType) {
	_this.switchMarkupKey(ctx)
}
func (_this *MarkupValueRule) OnList(ctx *Context)                      { ctx.BeginList() }
func (_this *MarkupValueRule) OnMap(ctx *Context)                       { ctx.BeginMap() }
func (_this *MarkupValueRule) OnMarkup(ctx *Context, identifier []byte) { ctx.BeginMarkup(identifier) }
func (_this *MarkupValueRule) OnComment(ctx *Context)                   { ctx.BeginComment() }
func (_this *MarkupValueRule) OnRelationship(ctx *Context)              { ctx.BeginRelationship() }
func (_this *MarkupValueRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
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
	ctx.BeginArrayString("markup contents", arrayType)
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

func (_this *SubjectRule) String() string { return "Subject Rule" }
func (_this *SubjectRule) moveToNextRule(ctx *Context) {
	ctx.ChangeRule(&predicateRule)
}
func (_this *SubjectRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.moveToNextRule(ctx) }
func (_this *SubjectRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *SubjectRule) OnList(ctx *Context)                            { ctx.BeginResourceList() }
func (_this *SubjectRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *SubjectRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *SubjectRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *SubjectRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowSubject)
}
func (_this *SubjectRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.moveToNextRule(ctx)
}
func (_this *SubjectRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *SubjectRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("relationship subject", arrayType, AllowSubject)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.moveToNextRule(ctx)
}
func (_this *SubjectRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("relationship subject", arrayType, AllowSubject)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.moveToNextRule(ctx)
}
func (_this *SubjectRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("relationship subject", arrayType, AllowSubject)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type PredicateRule struct{}

func (_this *PredicateRule) String() string { return "Predicate Rule" }
func (_this *PredicateRule) moveToNextRule(ctx *Context) {
	ctx.ChangeRule(&objectRule)
}
func (_this *PredicateRule) OnChildContainerEnded(ctx *Context, _ DataType) {
	_this.moveToNextRule(ctx)
}
func (_this *PredicateRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *PredicateRule) OnComment(ctx *Context) { ctx.BeginComment() }
func (_this *PredicateRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowPredicate)
}
func (_this *PredicateRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.moveToNextRule(ctx)
}
func (_this *PredicateRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *PredicateRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("relationship predicate", arrayType, AllowPredicate)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.moveToNextRule(ctx)
}
func (_this *PredicateRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("relationship predicate", arrayType, AllowPredicate)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.moveToNextRule(ctx)
}
func (_this *PredicateRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("relationship predicate", arrayType, AllowPredicate)
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type ObjectRule struct{}

func (_this *ObjectRule) String() string                                 { return "List Rule" }
func (_this *ObjectRule) end(ctx *Context)                               { ctx.EndContainer() }
func (_this *ObjectRule) OnChildContainerEnded(ctx *Context, _ DataType) { _this.end(ctx) }
func (_this *ObjectRule) OnNA(ctx *Context)                              { ctx.BeginNA() }
func (_this *ObjectRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *ObjectRule) OnNil(ctx *Context)                             { _this.end(ctx) }
func (_this *ObjectRule) OnKeyableObject(ctx *Context, _ DataType)       { _this.end(ctx) }
func (_this *ObjectRule) OnNonKeyableObject(ctx *Context, _ DataType)    { _this.end(ctx) }
func (_this *ObjectRule) OnList(ctx *Context)                            { ctx.BeginList() }
func (_this *ObjectRule) OnMap(ctx *Context)                             { ctx.BeginMap() }
func (_this *ObjectRule) OnMarkup(ctx *Context, identifier []byte)       { ctx.BeginMarkup(identifier) }
func (_this *ObjectRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *ObjectRule) OnRelationship(ctx *Context)                    { ctx.BeginRelationship() }
func (_this *ObjectRule) OnMarker(ctx *Context, identifier []byte) {
	ctx.BeginMarkerAnyType(identifier, AllowAny)
}
func (_this *ObjectRule) OnReference(ctx *Context, identifier []byte) {
	ctx.ReferenceAnyType(identifier)
	_this.end(ctx)
}
func (_this *ObjectRule) OnRIDReference(ctx *Context) {
	ctx.BeginRIDReference()
}
func (_this *ObjectRule) OnConstant(ctx *Context, name []byte) {
	ctx.BeginConstantAnyType(name)
}
func (_this *ObjectRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.AssertArrayType("relationship object", arrayType, AllowObject)
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.end(ctx)
}
func (_this *ObjectRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.AssertArrayType("relationship object", arrayType, AllowObject)
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.end(ctx)
}
func (_this *ObjectRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.AssertArrayType("relationship object", arrayType, AllowObject)
	ctx.BeginArrayAnyType(arrayType)
}
