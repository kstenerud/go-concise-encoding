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
	ctx.WriteIndent()
}
func (_this TopLevelDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type ListDecorator struct{}

var listDecorator ListDecorator

func (_this ListDecorator) String() string                  { return "ListDecorator" }
func (_this ListDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this ListDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteIndent()
}
func (_this ListDecorator) AfterValue(ctx *EncoderContext) {}
func (_this ListDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteIndent()
}
func (_this ListDecorator) AfterComment(ctx *EncoderContext) {}
func (_this ListDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteIndent()
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
	ctx.WriteIndent()
}
func (_this MapKeyDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteMapValueSeparator()
	ctx.Switch(mapValueDecorator)
}
func (_this MapKeyDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteIndent()
}
func (_this MapKeyDecorator) AfterComment(ctx *EncoderContext) {}
func (_this MapKeyDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteIndent()
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
	ctx.WriteIndent()
}
func (_this MapValueDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type MarkupKeyDecorator struct{}

var markupKeyDecorator MarkupKeyDecorator

func (_this MarkupKeyDecorator) String() string                  { return "MarkupKeyDecorator" }
func (_this MarkupKeyDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this MarkupKeyDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupKeySeparator()
}
func (_this MarkupKeyDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupValueSeparator()
	ctx.Switch(markupValueDecorator)
}
func (_this MarkupKeyDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.Stream.WriteMarkupKeySeparator()
}
func (_this MarkupKeyDecorator) AfterComment(ctx *EncoderContext) {}
func (_this MarkupKeyDecorator) EndContainer(ctx *EncoderContext) {
	ctx.BeginContainer()
	ctx.Switch(markupContentsDecorator)
	ctx.LastMarkupContentsWasComment = false
}

// ===========================================================================

type MarkupValueDecorator struct{}

var markupValueDecorator MarkupValueDecorator

func (_this MarkupValueDecorator) String() string                  { return "MarkupValueDecorator" }
func (_this MarkupValueDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this MarkupValueDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this MarkupValueDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(markupKeyDecorator)
}
func (_this MarkupValueDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this MarkupValueDecorator) AfterComment(ctx *EncoderContext) {
	ctx.WriteIndent()
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
	ctx.WriteIndent()
}
func (_this MarkupContentsDecorator) AfterValue(ctx *EncoderContext) {
	ctx.LastMarkupContentsWasComment = false
}
func (_this MarkupContentsDecorator) BeforeComment(ctx *EncoderContext) {
	if !ctx.ContainerHasObjects {
		ctx.Stream.WriteMarkupContentsBegin()
	}
	ctx.WriteIndent()
}
func (_this MarkupContentsDecorator) AfterComment(ctx *EncoderContext) {
	ctx.LastMarkupContentsWasComment = true
}
func (_this MarkupContentsDecorator) EndContainer(ctx *EncoderContext) {
	ctx.Unindent()
	if ctx.ContainerHasObjects {
		ctx.WriteIndent()
	}
	ctx.Stream.WriteMarkupEnd()
	ctx.Unstack()
	ctx.AfterValue()
}

// ===========================================================================

type CommentDecorator struct{}

var commentDecorator CommentDecorator

func (_this CommentDecorator) String() string                  { return "CommentDecorator" }
func (_this CommentDecorator) GetStringContext() stringContext { return stringContextComment }
func (_this CommentDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.Stream.WriteByte(' ')
}
func (_this CommentDecorator) AfterValue(ctx *EncoderContext) {}
func (_this CommentDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.Stream.WriteByte(' ')
}
func (_this CommentDecorator) AfterComment(ctx *EncoderContext) {}
func (_this CommentDecorator) EndContainer(ctx *EncoderContext) {
	if ctx.ContainerHasObjects {
		ctx.Stream.WriteByte(' ')
	}
	ctx.Stream.WriteCommentEnd()
	ctx.Unstack()
	ctx.AfterComment()
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

type SubjectDecorator struct{}

var subjectDecorator SubjectDecorator

func (_this SubjectDecorator) String() string                  { return "SubjectDecorator" }
func (_this SubjectDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this SubjectDecorator) BeforeValue(ctx *EncoderContext) {}
func (_this SubjectDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(predicateDecorator)
}
func (_this SubjectDecorator) BeforeComment(ctx *EncoderContext) {}
func (_this SubjectDecorator) AfterComment(ctx *EncoderContext)  {}
func (_this SubjectDecorator) EndContainer(ctx *EncoderContext)  { errorBadEvent(_this, "End") }

// ===========================================================================

type PredicateDecorator struct{}

var predicateDecorator PredicateDecorator

func (_this PredicateDecorator) String() string                  { return "PredicateDecorator" }
func (_this PredicateDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this PredicateDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteSpace()
}
func (_this PredicateDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Switch(objectDecorator)
}
func (_this PredicateDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteSpace()
}
func (_this PredicateDecorator) AfterComment(ctx *EncoderContext) {}
func (_this PredicateDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }

// ===========================================================================

type ObjectDecorator struct{}

var objectDecorator ObjectDecorator

func (_this ObjectDecorator) String() string                  { return "ObjectDecorator" }
func (_this ObjectDecorator) GetStringContext() stringContext { return stringContextDefault }
func (_this ObjectDecorator) BeforeValue(ctx *EncoderContext) {
	ctx.WriteSpace()
}
func (_this ObjectDecorator) AfterValue(ctx *EncoderContext) {
	ctx.Stream.WriteRelationshipEnd()
	ctx.Unstack()
	ctx.AfterValue()
}
func (_this ObjectDecorator) BeforeComment(ctx *EncoderContext) {
	ctx.WriteSpace()
}
func (_this ObjectDecorator) AfterComment(ctx *EncoderContext) {}
func (_this ObjectDecorator) EndContainer(ctx *EncoderContext) { errorBadEvent(_this, "End") }
