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
func (_this *BeginDocumentRule) OnBeginDocument(ctx *Context) { ctx.SwitchVersion() }

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
	ctx.SwitchTopLevel()
}

// =============================================================================

type TopLevelRule struct{}

func (_this *TopLevelRule) String() string                  { return "Top Level Rule" }
func (_this *TopLevelRule) OnKeyableObject(ctx *Context)    { ctx.SwitchEndDocument() }
func (_this *TopLevelRule) OnNonKeyableObject(ctx *Context) { ctx.SwitchEndDocument() }
func (_this *TopLevelRule) OnNA(ctx *Context) {
	_this.OnNonKeyableObject(ctx)
	ctx.BeginNA()
}
func (_this *TopLevelRule) OnChildContainerEnded(ctx *Context, _ DataType) { ctx.SwitchEndDocument() }
func (_this *TopLevelRule) OnPadding(ctx *Context)                         { /* Nothing to do */ }
func (_this *TopLevelRule) OnInt(ctx *Context, value int64)                { ctx.SwitchEndDocument() }
func (_this *TopLevelRule) OnPositiveInt(ctx *Context, value uint64) {
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnBigInt(ctx *Context, value *big.Int) {
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnFloat(ctx *Context, value float64) { ctx.SwitchEndDocument() }
func (_this *TopLevelRule) OnBigFloat(ctx *Context, value *big.Float) {
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) {
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal) {
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnList(ctx *Context)      { ctx.BeginList() }
func (_this *TopLevelRule) OnMap(ctx *Context)       { ctx.BeginMap() }
func (_this *TopLevelRule) OnMarkup(ctx *Context)    { ctx.BeginMarkup() }
func (_this *TopLevelRule) OnMetadata(ctx *Context)  { ctx.BeginMetadata() }
func (_this *TopLevelRule) OnComment(ctx *Context)   { ctx.BeginComment() }
func (_this *TopLevelRule) OnMarker(ctx *Context)    { ctx.BeginMarkerAnyType() }
func (_this *TopLevelRule) OnReference(ctx *Context) { ctx.BeginTopLevelReference() }
func (_this *TopLevelRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	ctx.SwitchEndDocument()
}
func (_this *TopLevelRule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}

// =============================================================================

type NARule struct{}

func (_this *NARule) String() string                                          { return "NA Rule" }
func (_this *NARule) OnKeyableObject(ctx *Context)                            { ctx.UnstackRule() }
func (_this *NARule) OnNonKeyableObject(ctx *Context)                         { ctx.UnstackRule() }
func (_this *NARule) OnNA(ctx *Context)                                       { ctx.UnstackRule() }
func (_this *NARule) OnChildContainerEnded(ctx *Context, _ DataType)          { ctx.UnstackRule() }
func (_this *NARule) OnPadding(ctx *Context)                                  { /* Nothing to do */ }
func (_this *NARule) OnInt(ctx *Context, value int64)                         { ctx.UnstackRule() }
func (_this *NARule) OnPositiveInt(ctx *Context, value uint64)                { ctx.UnstackRule() }
func (_this *NARule) OnBigInt(ctx *Context, value *big.Int)                   { ctx.UnstackRule() }
func (_this *NARule) OnFloat(ctx *Context, value float64)                     { ctx.UnstackRule() }
func (_this *NARule) OnBigFloat(ctx *Context, value *big.Float)               { ctx.UnstackRule() }
func (_this *NARule) OnDecimalFloat(ctx *Context, value compact_float.DFloat) { ctx.UnstackRule() }
func (_this *NARule) OnBigDecimalFloat(ctx *Context, value *apd.Decimal)      { ctx.UnstackRule() }
func (_this *NARule) OnList(ctx *Context)                                     { ctx.BeginList() }
func (_this *NARule) OnMap(ctx *Context)                                      { ctx.BeginMap() }
func (_this *NARule) OnMarkup(ctx *Context)                                   { ctx.BeginMarkup() }
func (_this *NARule) OnMetadata(ctx *Context)                                 { ctx.BeginMetadata() }
func (_this *NARule) OnComment(ctx *Context)                                  { ctx.BeginComment() }
func (_this *NARule) OnMarker(ctx *Context)                                   { ctx.BeginMarkerAnyType() }
func (_this *NARule) OnReference(ctx *Context)                                { ctx.BeginTopLevelReference() }
func (_this *NARule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	ctx.UnstackRule()
}
func (_this *NARule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	ctx.UnstackRule()
}
func (_this *NARule) OnArrayBegin(ctx *Context, arrayType events.ArrayType) {
	ctx.BeginArrayAnyType(arrayType)
}
