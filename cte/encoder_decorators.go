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

package cte

import (
	"fmt"
)

// Encoder decorators take care of indentation and other pretty-printing details.
type EncoderDecorator interface {
	String() string
	BeforeValue(ctx *EncoderContext)
	AfterValue(ctx *EncoderContext)
	BeforeComment(ctx *EncoderContext)
	AfterComment(ctx *EncoderContext)
	EndContainer(ctx *EncoderContext)
}

func errorBadEvent(receiver interface{}, event string) {
	panic(fmt.Errorf("BUG: %v cannot respond to %v", receiver, event))
}

// ===========================================================================

type TopLevelDecorator struct{}

var topLevelDecorator TopLevelDecorator

func (_this TopLevelDecorator) String() string                    { return "TopLevelDecorator" }
func (_this TopLevelDecorator) BeforeValue(ctx *EncoderContext)   {}
func (_this TopLevelDecorator) AfterValue(ctx *EncoderContext)    {}
func (_this TopLevelDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this TopLevelDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this TopLevelDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type ListDecorator struct{}

var listDecorator ListDecorator

func (_this ListDecorator) String() string { return "ListDecorator" }
func (_this ListDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this ListDecorator) AfterValue(ctx *EncoderContext) {}
func (_this ListDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this ListDecorator) AfterComment(ctx *EncoderContext) {}
func (_this ListDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteListEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type MapKeyDecorator struct{}

var mapKeyDecorator MapKeyDecorator

func (_this MapKeyDecorator) String() string { return "MapKeyDecorator" }
func (_this MapKeyDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this MapKeyDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteMapValueSeparator()
	ctx.Switch(mapValueDecorator)
}
func (_this MapKeyDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this MapKeyDecorator) AfterComment(ctx *EncoderContext) {}
func (_this MapKeyDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteMapEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type MapValueDecorator struct{}

var mapValueDecorator MapValueDecorator

func (_this MapValueDecorator) String() string                  { return "MapValueDecorator" }
func (_this MapValueDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this MapValueDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(mapKeyDecorator)
}
func (_this MapValueDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this MapValueDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this MapValueDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type RecordTypeDecorator struct{}

var recordTypeDecorator RecordTypeDecorator

func (_this RecordTypeDecorator) String() string { return "RecordTypeDecorator" }
func (_this RecordTypeDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this RecordTypeDecorator) AfterValue(ctx *EncoderContext) {}
func (_this RecordTypeDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this RecordTypeDecorator) AfterComment(ctx *EncoderContext) {}
func (_this RecordTypeDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteRecordTypeEnd()
	ctx.Unstack()
	ctx.WriteNewlineAndOriginAndIndent()
}

// ===========================================================================

type RecordDecorator struct{}

var recordDecorator RecordDecorator

func (_this RecordDecorator) String() string { return "RecordDecorator" }
func (_this RecordDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this RecordDecorator) AfterValue(ctx *EncoderContext) {}
func (_this RecordDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this RecordDecorator) AfterComment(ctx *EncoderContext) {}
func (_this RecordDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteRecordEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

// ===========================================================================

type ConcatDecorator struct{}

var concatDecorator ConcatDecorator

func (_this ConcatDecorator) String() string                  { return "ConcatDecorator" }
func (_this ConcatDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this ConcatDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Unstack()
	ctx.AfterValue()
}
func (_this ConcatDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this ConcatDecorator) AfterComment(ctx *EncoderContext)  {}
func (_this ConcatDecorator) EndContainer(ctx *EncoderContext)  { errorBadEvent(_this, "End") }

// ===========================================================================

type EdgeDecorator struct{}

var edgeDecorator EdgeDecorator

func (_this EdgeDecorator) String() string { return "EdgeDecorator" }
func (_this EdgeDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this EdgeDecorator) AfterValue(ctx *EncoderContext) {}
func (_this EdgeDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this EdgeDecorator) AfterComment(ctx *EncoderContext) {}
func (_this EdgeDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteEdgeEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type NodeValueDecorator struct{}

var nodeValueDecorator NodeValueDecorator

func (_this NodeValueDecorator) String() string { return "NodeValueDecorator" }
func (_this NodeValueDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteIndentIfOrigin()
}
func (_this NodeValueDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(nodeChildrenDecorator)
}
func (_this NodeValueDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this NodeValueDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteReturnToOrigin()
}
func (_this NodeValueDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type NodeChildrenDecorator struct{}

var nodeChildrenDecorator NodeChildrenDecorator

func (_this NodeChildrenDecorator) String() string { return "NodeChildrenDecorator" }
func (_this NodeChildrenDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this NodeChildrenDecorator) AfterValue(ctx *EncoderContext) {}
func (_this NodeChildrenDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this NodeChildrenDecorator) AfterComment(ctx *EncoderContext) {}
func (_this NodeChildrenDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteNodeEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type NonStringArrayDecorator struct{}

var nonStringArrayDecorator NonStringArrayDecorator

func (_this NonStringArrayDecorator) String() string { return "NonStringArrayDecorator" }
func (_this NonStringArrayDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this NonStringArrayDecorator) AfterValue(ctx *EncoderContext) {}
func (_this NonStringArrayDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this NonStringArrayDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this NonStringArrayDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteListEnd()
	ctx.Unstack()
	ctx.AfterValue()
}
