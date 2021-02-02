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

	"github.com/kstenerud/go-concise-encoding/events"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

type ListRule struct{}

func (_this *ListRule) String() string                  { return "List Rule" }
func (_this *ListRule) OnKeyableObject(ctx *Context)    { /* Nothing to do */ }
func (_this *ListRule) OnNonKeyableObject(ctx *Context) { /* Nothing to do */ }
func (_this *ListRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *ListRule) OnChildContainerEnded(ctx *Context, _ DataType)          { /* Nothing to do */ }
func (_this *ListRule) OnPadding(ctx *Context)                                  { /* Nothing to do */ }
func (_this *ListRule) OnInt(ctx *Context, value int64)                         { /* Nothing to do */ }
func (_this *ListRule) OnPositiveInt(ctx *Context, value uint64)                { /* Nothing to do */ }
func (_this *ListRule) OnBigInt(ctx *Context, value *big.Int)                   { /* Nothing to do */ }
func (_this *ListRule) OnFloat(ctx *Context, value float64)                     { /* Nothing to do */ }
func (_this *ListRule) OnBigFloat(ctx *Context, value *big.Float)               { /* Nothing to do */ }
func (_this *ListRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) { /* Nothing to do */ }
func (_this *ListRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal)      { /* Nothing to do */ }
func (_this *ListRule) OnList(ctx *Context)                                     { ctx.BeginList() }
func (_this *ListRule) OnMap(ctx *Context)                                      { ctx.BeginMap() }
func (_this *ListRule) OnMarkup(ctx *Context)                                   { ctx.BeginMarkup() }
func (_this *ListRule) OnMetadata(ctx *Context)                                 { ctx.BeginMetadata() }
func (_this *ListRule) OnComment(ctx *Context)                                  { ctx.BeginComment() }
func (_this *ListRule) OnEnd(ctx *Context)                                      { ctx.EndContainer() }
func (_this *ListRule) OnMarker(ctx *Context)                                   { ctx.BeginMarkerAnyType() }
func (_this *ListRule) OnReference(ctx *Context)                                { ctx.BeginReferenceAnyType() }
func (_this *ListRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
}
func (_this *ListRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MapKeyRule struct{}

func (_this *MapKeyRule) String() string                                 { return "Map Key Rule" }
func (_this *MapKeyRule) OnKeyableObject(ctx *Context)                   { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapKeyRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnPositiveInt(ctx *Context, value uint64)       { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnBigInt(ctx *Context, value *big.Int)          { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnFloat(ctx *Context, value float64)            { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnBigFloat(ctx *Context, value *big.Float)      { ctx.SwitchMapValue() }
func (_this *MapKeyRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMapValue()
}
func (_this *MapKeyRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMapValue()
}
func (_this *MapKeyRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MapKeyRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MapKeyRule) OnEnd(ctx *Context)       { ctx.EndContainer() }
func (_this *MapKeyRule) OnMarker(ctx *Context)    { ctx.BeginMarkerKeyable() }
func (_this *MapKeyRule) OnReference(ctx *Context) { ctx.BeginReferenceKeyable() }
func (_this *MapKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	ctx.SwitchMapValue()
}
func (_this *MapKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MapValueRule struct{}

func (_this *MapValueRule) String() string                  { return "Map Value Rule" }
func (_this *MapValueRule) OnKeyableObject(ctx *Context)    { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnNonKeyableObject(ctx *Context) { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *MapValueRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MapValueRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnPositiveInt(ctx *Context, value uint64) {
	ctx.SwitchMapKey()
}
func (_this *MapValueRule) OnBigInt(ctx *Context, value *big.Int) { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnFloat(ctx *Context, value float64)   { ctx.SwitchMapKey() }
func (_this *MapValueRule) OnBigFloat(ctx *Context, value *big.Float) {
	ctx.SwitchMapKey()
}
func (_this *MapValueRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMapKey()
}
func (_this *MapValueRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMapKey()
}
func (_this *MapValueRule) OnList(ctx *Context)      { ctx.BeginList() }
func (_this *MapValueRule) OnMap(ctx *Context)       { ctx.BeginMap() }
func (_this *MapValueRule) OnMarkup(ctx *Context)    { ctx.BeginMarkup() }
func (_this *MapValueRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MapValueRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MapValueRule) OnMarker(ctx *Context)    { ctx.BeginMarkerAnyType() }
func (_this *MapValueRule) OnReference(ctx *Context) { ctx.BeginReferenceAnyType() }
func (_this *MapValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.SwitchMapKey()
}
func (_this *MapValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupNameRule struct{}

func (_this *MarkupNameRule) String() string                                 { return "Markup Name Rule" }
func (_this *MarkupNameRule) OnKeyableObject(ctx *Context)                   { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupNameRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnPositiveInt(ctx *Context, value uint64)       { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnBigInt(ctx *Context, value *big.Int)          { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnFloat(ctx *Context, value float64)            { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnBigFloat(ctx *Context, value *big.Float)      { ctx.SwitchMarkupKey() }
func (_this *MarkupNameRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupNameRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupNameRule) OnMarker(ctx *Context)    { ctx.BeginMarkerKeyable() }
func (_this *MarkupNameRule) OnReference(ctx *Context) { ctx.BeginReferenceKeyable() }
func (_this *MarkupNameRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	ctx.SwitchMarkupKey()
}
func (_this *MarkupNameRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MarkupKeyRule struct{}

func (_this *MarkupKeyRule) String() string                                 { return "Markup Attribute Key Rule" }
func (_this *MarkupKeyRule) OnKeyableObject(ctx *Context)                   { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupKeyRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnPositiveInt(ctx *Context, value uint64)       { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnBigInt(ctx *Context, value *big.Int)          { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnFloat(ctx *Context, value float64)            { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnBigFloat(ctx *Context, value *big.Float)      { ctx.SwitchMarkupValue() }
func (_this *MarkupKeyRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMarkupValue()
}
func (_this *MarkupKeyRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMarkupValue()
}
func (_this *MarkupKeyRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MarkupKeyRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MarkupKeyRule) OnEnd(ctx *Context)       { ctx.SwitchMarkupContents() }
func (_this *MarkupKeyRule) OnMarker(ctx *Context)    { ctx.BeginMarkerKeyable() }
func (_this *MarkupKeyRule) OnReference(ctx *Context) { ctx.BeginReferenceKeyable() }
func (_this *MarkupKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	ctx.SwitchMarkupValue()
}
func (_this *MarkupKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MarkupValueRule struct{}

func (_this *MarkupValueRule) String() string                  { return "Markup Attribute Value Rule" }
func (_this *MarkupValueRule) OnKeyableObject(ctx *Context)    { ctx.SwitchMarkupKey() }
func (_this *MarkupValueRule) OnNonKeyableObject(ctx *Context) { ctx.SwitchMarkupKey() }
func (_this *MarkupValueRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *MarkupValueRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMarkupKey() }
func (_this *MarkupValueRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupValueRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMarkupKey() }
func (_this *MarkupValueRule) OnPositiveInt(ctx *Context, value uint64) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnBigInt(ctx *Context, value *big.Int) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnFloat(ctx *Context, value float64) { ctx.SwitchMarkupKey() }
func (_this *MarkupValueRule) OnBigFloat(ctx *Context, value *big.Float) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnList(ctx *Context)      { ctx.BeginList() }
func (_this *MarkupValueRule) OnMap(ctx *Context)       { ctx.BeginMap() }
func (_this *MarkupValueRule) OnMarkup(ctx *Context)    { ctx.BeginMarkup() }
func (_this *MarkupValueRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MarkupValueRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MarkupValueRule) OnMarker(ctx *Context)    { ctx.BeginMarkerAnyType() }
func (_this *MarkupValueRule) OnReference(ctx *Context) { ctx.BeginReferenceAnyType() }
func (_this *MarkupValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.SwitchMarkupKey()
}
func (_this *MarkupValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MarkupContentsRule struct{}

func (_this *MarkupContentsRule) String() string                                 { return "List Rule" }
func (_this *MarkupContentsRule) OnChildContainerEnded(ctx *Context, _ DataType) { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MarkupContentsRule) OnMarkup(ctx *Context)                          { ctx.BeginMarkup() }
func (_this *MarkupContentsRule) OnComment(ctx *Context)                         { ctx.BeginComment() }
func (_this *MarkupContentsRule) OnEnd(ctx *Context)                             { ctx.EndContainer() }
func (_this *MarkupContentsRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayString(arrayType, elementCount, data)
}
func (_this *MarkupContentsRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayString(arrayType)
}

// =============================================================================

type CommentRule struct{}

func (_this *CommentRule) String() string                   { return "Comment Rule" }
func (_this *CommentRule) OnChildCommentEnded(ctx *Context) { /* Nothing to do */ }
func (_this *CommentRule) OnPadding(ctx *Context)           { /* Nothing to do */ }
func (_this *CommentRule) OnComment(ctx *Context)           { ctx.BeginComment() }
func (_this *CommentRule) OnEnd(ctx *Context)               { ctx.UnstackRule() }
func (_this *CommentRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayComment(arrayType, elementCount, data)
}
func (_this *CommentRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayComment(arrayType)
}

// =============================================================================

type MetaKeyRule struct{}

func (_this *MetaKeyRule) String() string                                 { return "Metadata Key Rule" }
func (_this *MetaKeyRule) OnKeyableObject(ctx *Context)                   { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MetaKeyRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnPositiveInt(ctx *Context, value uint64)       { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnBigInt(ctx *Context, value *big.Int)          { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnFloat(ctx *Context, value float64)            { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnBigFloat(ctx *Context, value *big.Float)      { ctx.SwitchMetadataValue() }
func (_this *MetaKeyRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMetadataValue()
}
func (_this *MetaKeyRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMetadataValue()
}
func (_this *MetaKeyRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MetaKeyRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MetaKeyRule) OnEnd(ctx *Context)       { ctx.SwitchMetadataCompletion() }
func (_this *MetaKeyRule) OnMarker(ctx *Context)    { ctx.BeginMarkerKeyable() }
func (_this *MetaKeyRule) OnReference(ctx *Context) { ctx.BeginReferenceKeyable() }
func (_this *MetaKeyRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayKeyable(arrayType, elementCount, data)
	ctx.SwitchMetadataValue()
}
func (_this *MetaKeyRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayKeyable(arrayType)
}

// =============================================================================

type MetaValueRule struct{}

func (_this *MetaValueRule) String() string                  { return "Metadata Value Rule" }
func (_this *MetaValueRule) OnKeyableObject(ctx *Context)    { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnNonKeyableObject(ctx *Context) { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *MetaValueRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *MetaValueRule) OnInt(ctx *Context, value int64)                { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnPositiveInt(ctx *Context, value uint64)       { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnBigInt(ctx *Context, value *big.Int)          { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnFloat(ctx *Context, value float64)            { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnBigFloat(ctx *Context, value *big.Float)      { ctx.SwitchMetadataKey() }
func (_this *MetaValueRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchMetadataKey()
}
func (_this *MetaValueRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchMetadataKey()
}
func (_this *MetaValueRule) OnList(ctx *Context)      { ctx.BeginList() }
func (_this *MetaValueRule) OnMap(ctx *Context)       { ctx.BeginMap() }
func (_this *MetaValueRule) OnMarkup(ctx *Context)    { ctx.BeginMarkup() }
func (_this *MetaValueRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *MetaValueRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *MetaValueRule) OnMarker(ctx *Context)    { ctx.BeginMarkerAnyType() }
func (_this *MetaValueRule) OnReference(ctx *Context) { ctx.BeginReferenceAnyType() }
func (_this *MetaValueRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.SwitchMetadataKey()
}
func (_this *MetaValueRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type MetaCompletionRule struct{}

func (_this *MetaCompletionRule) String() string { return "Metadata Completion Rule" }
func (_this *MetaCompletionRule) OnKeyableObject(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnKeyableObject(ctx)
}
func (_this *MetaCompletionRule) OnNonKeyableObject(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnNonKeyableObject(ctx)
}
func (_this *MetaCompletionRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *MetaCompletionRule) OnPadding(ctx *Context) { /* Nothing to do */ }
func (_this *MetaCompletionRule) OnInt(ctx *Context, value int64) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnInt(ctx, value)
}
func (_this *MetaCompletionRule) OnPositiveInt(ctx *Context, value uint64) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnPositiveInt(ctx, value)
}
func (_this *MetaCompletionRule) OnBigInt(ctx *Context, value *big.Int) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnBigInt(ctx, value)
}
func (_this *MetaCompletionRule) OnFloat(ctx *Context, value float64) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnFloat(ctx, value)
}
func (_this *MetaCompletionRule) OnBigFloat(ctx *Context, value *big.Float) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnBigFloat(ctx, value)
}
func (_this *MetaCompletionRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnDecimalFloat(ctx, value)
}
func (_this *MetaCompletionRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnBigDecimalFloat(ctx, value)
}
func (_this *MetaCompletionRule) OnList(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnList(ctx)
}
func (_this *MetaCompletionRule) OnMap(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnMap(ctx)
}
func (_this *MetaCompletionRule) OnMarkup(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnMarkup(ctx)
}
func (_this *MetaCompletionRule) OnMetadata(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnMetadata(ctx)
}
func (_this *MetaCompletionRule) OnComment(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnComment(ctx)
}
func (_this *MetaCompletionRule) OnMarker(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnMarker(ctx)
}
func (_this *MetaCompletionRule) OnReference(ctx *Context) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnReference(ctx)
}
func (_this *MetaCompletionRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArray(ctx, arrayType, elementCount, data)
}
func (_this *MetaCompletionRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.UnstackRule()
	ctx.CurrentEntry.Rule.OnArrayBegin(ctx, arrayType)
}
