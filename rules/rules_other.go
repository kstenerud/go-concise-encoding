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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-concise-encoding/ce/events"
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
func (_this *TopLevelRule) OnNull(ctx *Context)                       { _this.switchEndDocument(ctx) }
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
func (_this *TopLevelRule) OnArray(ctx *Context, arrayType events.ArrayType, elementCount uint64, data []uint8) {
	ctx.ValidateFullArrayAnyType(arrayType, elementCount, data)
	_this.switchEndDocument(ctx)
}
func (_this *TopLevelRule) OnStringlikeArray(ctx *Context, arrayType events.ArrayType, data string) {
	ctx.ValidateFullArrayStringlike(arrayType, data)
	_this.switchEndDocument(ctx)
}
