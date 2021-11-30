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

type stringContext int

const (
	stringContextDefault stringContext = iota
	stringContextComment
	stringContextMarkup
)

// Encoder decorators take care of indentation and other pretty-printing details.
type EncoderDecorator interface {
	String() string
	GetStringContext() stringContext
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
func (_this TopLevelDecorator) GetStringContext() stringContext   { return stringContextDefault }
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

func (_this ListDecorator) String() string                  { return "ListDecorator" }
func (_this ListDecorator) GetStringContext() stringContext { return stringContextDefault }
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

func (_this MapKeyDecorator) String() string                  { return "MapKeyDecorator" }
func (_this MapKeyDecorator) GetStringContext() stringContext { return stringContextDefault }
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
func (_this MapValueDecorator) GetStringContext() stringContext { return stringContextDefault }
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

type MarkupKeyDecorator struct{}

var markupKeyDecorator MarkupKeyDecorator

func (_this MarkupKeyDecorator) String() string                  { return "MarkupKeyDecorator" }
func (_this MarkupKeyDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this MarkupKeyDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteIndentOrSpace()
}
func (_this MarkupKeyDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupValueSeparator()
	ctx.Switch(markupValueDecorator)
}
func (_this MarkupKeyDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteIndentOrSpace()
}
func (_this MarkupKeyDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteReturnToOrigin()
}
func (_this MarkupKeyDecorator) EndContainer(ctx *EncoderContext) {
	ctx.BeginContainer()
	ctx.Switch(markupContentsDecorator)
}

// ===========================================================================

type MarkupValueDecorator struct{}

var markupValueDecorator MarkupValueDecorator

func (_this MarkupValueDecorator) String() string                  { return "MarkupValueDecorator" }
func (_this MarkupValueDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this MarkupValueDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteIndentIfOrigin()
}
func (_this MarkupValueDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(markupKeyDecorator)
}
func (_this MarkupValueDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteIndentIfOrigin()
}
func (_this MarkupValueDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteReturnToOrigin()
}
func (_this MarkupValueDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type MarkupContentsDecorator struct{}

var markupContentsDecorator MarkupContentsDecorator

func (_this MarkupContentsDecorator) String() string                  { return "MarkupContentsDecorator" }
func (_this MarkupContentsDecorator) GetStringContext() stringContext { return stringContextMarkup }
func (_this MarkupContentsDecorator) BeforeValue(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.WriteMarkupContentsBegin()
	}
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this MarkupContentsDecorator) AfterValue(ctx *EncoderContext) {}
func (_this MarkupContentsDecorator) BeforeComment(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.WriteMarkupContentsBegin()
	}
	ctx.WriteNewlineAndOriginAndIndent()
}
func (_this MarkupContentsDecorator) AfterComment(ctx *EncoderContext) {}
func (_this MarkupContentsDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteNewlineAndOriginAndIndent()
	}
	ctx.Stream.WriteMarkupEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type ConcatDecorator struct{}

var concatDecorator ConcatDecorator

func (_this ConcatDecorator) String() string                  { return "ConcatDecorator" }
func (_this ConcatDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this ConcatDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this ConcatDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Unstack()
	ctx.AfterValue()
}
func (_this ConcatDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this ConcatDecorator) AfterComment(ctx *EncoderContext)  {}
func (_this ConcatDecorator) EndContainer(ctx *EncoderContext)  { errorBadEvent(_this, "End") }

// ===========================================================================

type EdgeSourceDecorator struct{}

var edgeSourceDecorator EdgeSourceDecorator

func (_this EdgeSourceDecorator) String() string                  { return "EdgeSourceDecorator" }
func (_this EdgeSourceDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this EdgeSourceDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this EdgeSourceDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(edgeDescriptionDecorator)
}
func (_this EdgeSourceDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this EdgeSourceDecorator) AfterComment(ctx *EncoderContext)  {}
func (_this EdgeSourceDecorator) EndContainer(ctx *EncoderContext)  { errorBadEvent(_this, "End") }

// ===========================================================================

type EdgeDescriptionDecorator struct{}

var edgeDescriptionDecorator EdgeDescriptionDecorator

func (_this EdgeDescriptionDecorator) String() string                  { return "EdgeDescriptionDecorator" }
func (_this EdgeDescriptionDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this EdgeDescriptionDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteElementSeparator()
}
func (_this EdgeDescriptionDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(edgeDestinationDecorator)
}
func (_this EdgeDescriptionDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteElementSeparator()
}
func (_this EdgeDescriptionDecorator) AfterComment(ctx *EncoderContext) {}
func (_this EdgeDescriptionDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type EdgeDestinationDecorator struct{}

var edgeDestinationDecorator EdgeDestinationDecorator

func (_this EdgeDestinationDecorator) String() string                  { return "EdgeDestinationDecorator" }
func (_this EdgeDestinationDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this EdgeDestinationDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteElementSeparator()
}
func (_this EdgeDestinationDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteEdgeEnd()
	ctx.Unstack()
	ctx.AfterValue()
}
func (_this EdgeDestinationDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteElementSeparator()
}
func (_this EdgeDestinationDecorator) AfterComment(ctx *EncoderContext) {}
func (_this EdgeDestinationDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type NodeValueDecorator struct{}

var nodeValueDecorator NodeValueDecorator

func (_this NodeValueDecorator) String() string                  { return "NodeValueDecorator" }
func (_this NodeValueDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this NodeValueDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this NodeValueDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(nodeChildrenDecorator)
	ctx.Indent()
}
func (_this NodeValueDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this NodeValueDecorator) AfterComment(ctx *EncoderContext)  {}
func (_this NodeValueDecorator) EndContainer(ctx *EncoderContext)  { errorBadEvent(_this, "End") }

// ===========================================================================

type NodeChildrenDecorator struct{}

var nodeChildrenDecorator NodeChildrenDecorator

func (_this NodeChildrenDecorator) String() string                  { return "NodeChildrenDecorator" }
func (_this NodeChildrenDecorator) GetStringContext() stringContext { return stringContextDefault }
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

func (_this NonStringArrayDecorator) String() string                  { return "NonStringArrayDecorator" }
func (_this NonStringArrayDecorator) GetStringContext() stringContext { return stringContextDefault }
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
