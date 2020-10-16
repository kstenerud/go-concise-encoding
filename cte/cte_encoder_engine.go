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
	"strings"

	"github.com/kstenerud/go-concise-encoding/buffer"
)

const prefixInitialBufferSize = 32

type encoderEngine struct {
	stream                 *buffer.StreamingWriteBuffer
	initialState           awaiting
	indent                 string
	prefixSetters          []prefixSetter
	transitionFromSuffixes []string

	Awaiting       awaiting
	ItemCount      int
	awaitingStack  []awaiting
	itemCountStack []int
	prefix         []byte
	ContainerDepth int
}

func (_this *encoderEngine) Init(stream *buffer.StreamingWriteBuffer, indent string) {
	_this.awaitingStack = []awaiting{}
	_this.itemCountStack = []int{}
	_this.prefix = make([]byte, 0, prefixInitialBufferSize)
	_this.stream = stream
	_this.indent = indent
	if len(indent) == 0 {
		_this.prefixSetters = prefixSettersCompact[:]
		_this.transitionFromSuffixes = transitionFromSuffixesCompact[:]
	} else {
		_this.prefixSetters = prefixSettersPretty[:]
		_this.transitionFromSuffixes = transitionFromSuffixesPretty[:]
	}
	_this.Reset()
}

func (_this *encoderEngine) Reset() {
	_this.Awaiting = awaitingVersion
	_this.awaitingStack = _this.awaitingStack[:0]
	_this.ItemCount = 0
	_this.itemCountStack = _this.itemCountStack[:0]
	_this.clearPrefix()
}

func (_this *encoderEngine) AddVersion(version uint64) {
	_this.stream.AddFmt("c%d", version)
	if _this.isPretty() {
		_this.setPrefix("\n")
	} else {
		_this.setPrefix(" ")
	}
	_this.transition()
}

func (_this *encoderEngine) BeginObject() {
	_this.flushPrefix()
	_this.ItemCount++
}

func (_this *encoderEngine) CompleteObject() {
	_this.transition()
}

func (_this *encoderEngine) BeginPseudoObject() {
	_this.flushPrefix()
	_this.ItemCount++
}

func (_this *encoderEngine) CompletePseudoObject() {
}

func (_this *encoderEngine) BeginContainer(newState awaiting, opener string) {
	_this.appendPrefix(opener)
	_this.BeginObject()
	_this.stack(newState)
	_this.ContainerDepth++
	_this.autoSetPrefix()
}

func (_this *encoderEngine) EndContainer() {
	switch _this.Awaiting {
	case awaitingMarkupFirstKey, awaitingMarkupKey:
		_this.EndMarkupAttrs()
	case awaitingMetaFirstKey, awaitingMetaKey:
		_this.CompletePseudoContainer()
	case awaitingCommentItem:
		_this.CompleteComment()
	default:
		_this.CompleteContainer()
	}
}

func (_this *encoderEngine) CompleteContainer() {
	itemCount := _this.ItemCount
	closer := containerTerminators[_this.Awaiting]
	_this.clearPrefix()
	_this.unstack()
	_this.ContainerDepth--
	if itemCount > 0 {
		_this.stream.AddString(_this.generateIndent())
	}
	_this.stream.AddString(closer)
	_this.CompleteObject()
}

func (_this *encoderEngine) BeginPseudoContainer(newState awaiting, opener string) {
	_this.appendPrefix(opener)
	_this.BeginPseudoObject()
	_this.stack(newState)
	_this.ContainerDepth++
	_this.autoSetPrefix()
}

func (_this *encoderEngine) CompletePseudoContainer() {
	itemCount := _this.ItemCount
	closer := containerTerminators[_this.Awaiting]
	_this.clearPrefix()
	_this.unstack()
	_this.ContainerDepth--
	if itemCount > 0 {
		_this.stream.AddString(_this.generateIndent())
	}
	_this.stream.AddString(closer)
	_this.CompletePseudoObject()
}

func (_this *encoderEngine) BeginArray(newState awaiting) {
	_this.BeginObject()
	_this.stack(newState)
	_this.ContainerDepth++
}

func (_this *encoderEngine) CompleteArray() {
	_this.clearPrefix()
	_this.unstack()
	_this.ContainerDepth--
	_this.stream.AddString("|")
	_this.CompleteObject()
}

func (_this *encoderEngine) BeginStringLikeArray(newState awaiting) {
	_this.stack(newState)
}

func (_this *encoderEngine) EndStringLikeArray() {
	_this.unstack()
}

func (_this *encoderEngine) BeginMarkup() {
	_this.appendPrefix("<")
	_this.BeginObject()
	_this.stack(awaitingMarkupFirstItem)
	_this.ContainerDepth++
	_this.stack(awaitingMarkupName)
}

func (_this *encoderEngine) EndMarkupAttrs() {
	_this.clearPrefix()
	_this.setPrefix(";" + _this.generateIndent())
	_this.unstack()
}

func (_this *encoderEngine) BeginComment() {
	if _this.isPretty() && _this.isPrefixEmpty() {
		_this.setPrefix(" ")
	}
	_this.BeginContainer(awaitingCommentItem, "/*")
}

func (_this *encoderEngine) AddCommentString(value string) {
	_this.BeginPseudoObject()
	if _this.isPretty() {
		_this.stream.AddString(" ")
	}
	_this.stream.AddString(value)
	_this.CompletePseudoObject()
}

func (_this *encoderEngine) CompleteComment() {
	itemCount := _this.ItemCount
	closer := containerTerminators[_this.Awaiting]
	_this.clearPrefix()
	_this.unstack()
	_this.ContainerDepth--
	if itemCount > 0 {
		if _this.isPretty() {
			_this.stream.AddString(" ")
		}
	}
	_this.stream.AddString(closer)
	_this.CompletePseudoObject()
	if _this.isPretty() {
		_this.setPrefix(_this.generateIndent())
	}
}

func (_this *encoderEngine) BeginMarker() {
	_this.stack(awaitingMarkerID)
}

func (_this *encoderEngine) CompleteMarker(markerValue interface{}) {
	_this.unstack()
	_this.appendPrefix(fmt.Sprintf("&%v:", markerValue))
}

func (_this *encoderEngine) BeginReference() {
	_this.stack(awaitingReferenceID)
}

func (_this *encoderEngine) CompleteReference(markerValue interface{}) {
	_this.unstack()
	_this.BeginObject()
	_this.stream.AddString(fmt.Sprintf("$%v", markerValue))
	_this.CompleteObject()
}

// ============================================================================

// Util

func (_this *encoderEngine) stack(newState awaiting) {
	_this.awaitingStack = append(_this.awaitingStack, _this.Awaiting)
	_this.itemCountStack = append(_this.itemCountStack, _this.ItemCount)
	_this.Awaiting = newState
	_this.ItemCount = 0
}

func (_this *encoderEngine) unstack() {
	newEnd := len(_this.awaitingStack) - 1

	_this.Awaiting = _this.awaitingStack[newEnd]
	_this.ItemCount = _this.itemCountStack[newEnd]
	_this.awaitingStack = _this.awaitingStack[:newEnd]
	_this.itemCountStack = _this.itemCountStack[:newEnd]
}

func (_this *encoderEngine) transition() {
	_this.stream.AddString(_this.transitionFromSuffixes[_this.Awaiting])
	_this.Awaiting = stateTransitions[_this.Awaiting]
	_this.autoSetPrefix()
}

func (_this *encoderEngine) isPretty() bool {
	return len(_this.indent) > 0
}

func (_this *encoderEngine) clearPrefix() {
	_this.prefix = _this.prefix[:0]
}

func (_this *encoderEngine) isPrefixEmpty() bool {
	return len(_this.prefix) == 0
}

func (_this *encoderEngine) setPrefix(value string) {
	_this.prefix = append(_this.prefix[:0], value...)
}

func (_this *encoderEngine) appendPrefix(value string) {
	_this.prefix = append(_this.prefix, value...)
}

func (_this *encoderEngine) autoSetPrefix() {
	_this.prefixSetters[_this.Awaiting](_this)
}

func (_this *encoderEngine) flushPrefix() {
	dst := _this.stream.Allocate(len(_this.prefix))
	copy(dst, _this.prefix)
	_this.clearPrefix()
}

func (_this *encoderEngine) setNoPrefix() {
}

func (_this *encoderEngine) setSpacePrefix() {
	_this.setPrefix(" ")
}

func (_this *encoderEngine) setIndentPrefix() {
	_this.setPrefix(_this.generatePrettyIndent())
}

func (_this *encoderEngine) generatePrettyIndent() string {
	return "\n" + strings.Repeat(_this.indent, _this.ContainerDepth)
}

func (_this *encoderEngine) generateIndent() string {
	if _this.isPretty() {
		return _this.generatePrettyIndent()
	}
	return ""
}

// ============================================================================

// Data

type awaiting int64

const (
	awaitingVersion awaiting = iota
	awaitingTLO
	awaitingListFirstItem
	awaitingListItem
	awaitingMapFirstKey
	awaitingMapKey
	awaitingMapValue
	awaitingMetaFirstKey
	awaitingMetaKey
	awaitingMetaValue
	awaitingMarkupName
	awaitingMarkupFirstKey
	awaitingMarkupKey
	awaitingMarkupValue
	awaitingMarkupFirstItem
	awaitingMarkupItem
	awaitingCommentItem
	awaitingMarkerID
	awaitingReferenceID
	awaitingQuotedString
	awaitingVerbatimString
	awaitingURI
	awaitingCustomBinary
	awaitingCustomText
	awaitingArrayBool
	awaitingArrayU8
	awaitingArrayU16
	awaitingArrayU32
	awaitingArrayU64
	awaitingArrayI8
	awaitingArrayI16
	awaitingArrayI32
	awaitingArrayI64
	awaitingArrayF16
	awaitingArrayF32
	awaitingArrayF64
	awaitingArrayUUID
	awaitingCount
)

var awaitingNames = [awaitingCount]string{
	awaitingVersion:         "Version",
	awaitingTLO:             "TLO",
	awaitingListFirstItem:   "ListFirstItem",
	awaitingListItem:        "ListItem",
	awaitingMapFirstKey:     "MapFirstKey",
	awaitingMapKey:          "MapKey",
	awaitingMapValue:        "MapValue",
	awaitingMetaFirstKey:    "MetaFirstKey",
	awaitingMetaKey:         "MetaKey",
	awaitingMetaValue:       "MetaValue",
	awaitingMarkupName:      "MarkupName",
	awaitingMarkupFirstKey:  "MarkupFirstKey",
	awaitingMarkupKey:       "MarkupKey",
	awaitingMarkupValue:     "MarkupValue",
	awaitingMarkupFirstItem: "MarkupFirstItem",
	awaitingMarkupItem:      "MarkupItem",
	awaitingCommentItem:     "CommentItem",
	awaitingMarkerID:        "MarkerID",
	awaitingReferenceID:     "ReferenceID",
	awaitingQuotedString:    "QuotedString",
	awaitingVerbatimString:  "VerbatimString",
	awaitingURI:             "URI",
	awaitingCustomBinary:    "CustomBinary",
	awaitingCustomText:      "CustomText",
	awaitingArrayBool:       "ArrayBool",
	awaitingArrayU8:         "ArrayU8",
	awaitingArrayU16:        "ArrayU16",
	awaitingArrayU32:        "ArrayU32",
	awaitingArrayU64:        "ArrayU64",
	awaitingArrayI8:         "ArrayI8",
	awaitingArrayI16:        "ArrayI16",
	awaitingArrayI32:        "ArrayI32",
	awaitingArrayI64:        "ArrayI64",
	awaitingArrayF16:        "ArrayF16",
	awaitingArrayF32:        "ArrayF32",
	awaitingArrayF64:        "ArrayF64",
	awaitingArrayUUID:       "ArrayUUID",
}

func (_this awaiting) String() string {
	return awaitingNames[_this]
}

var stateTransitions = [awaitingCount]awaiting{
	awaitingVersion:         awaitingTLO,
	awaitingTLO:             awaitingTLO,
	awaitingListFirstItem:   awaitingListItem,
	awaitingListItem:        awaitingListItem,
	awaitingMapFirstKey:     awaitingMapValue,
	awaitingMapKey:          awaitingMapValue,
	awaitingMapValue:        awaitingMapKey,
	awaitingMetaFirstKey:    awaitingMetaValue,
	awaitingMetaKey:         awaitingMetaValue,
	awaitingMetaValue:       awaitingMetaKey,
	awaitingMarkupName:      awaitingMarkupFirstKey,
	awaitingMarkupFirstKey:  awaitingMarkupValue,
	awaitingMarkupKey:       awaitingMarkupValue,
	awaitingMarkupValue:     awaitingMarkupKey,
	awaitingMarkupFirstItem: awaitingMarkupItem,
	awaitingMarkupItem:      awaitingMarkupItem,
	awaitingCommentItem:     awaitingCommentItem,
}

type prefixSetter func(*encoderEngine)

var prefixSettersCompact = [awaitingCount]prefixSetter{
	awaitingVersion:         (*encoderEngine).setNoPrefix,
	awaitingTLO:             (*encoderEngine).setNoPrefix,
	awaitingListFirstItem:   (*encoderEngine).setNoPrefix,
	awaitingListItem:        (*encoderEngine).setSpacePrefix,
	awaitingMapFirstKey:     (*encoderEngine).setNoPrefix,
	awaitingMapKey:          (*encoderEngine).setSpacePrefix,
	awaitingMapValue:        (*encoderEngine).setNoPrefix,
	awaitingMetaFirstKey:    (*encoderEngine).setNoPrefix,
	awaitingMetaKey:         (*encoderEngine).setSpacePrefix,
	awaitingMetaValue:       (*encoderEngine).setNoPrefix,
	awaitingMarkupName:      (*encoderEngine).setNoPrefix,
	awaitingMarkupFirstKey:  (*encoderEngine).setSpacePrefix,
	awaitingMarkupKey:       (*encoderEngine).setSpacePrefix,
	awaitingMarkupValue:     (*encoderEngine).setNoPrefix,
	awaitingMarkupFirstItem: (*encoderEngine).setNoPrefix,
	awaitingMarkupItem:      (*encoderEngine).setNoPrefix,
	awaitingCommentItem:     (*encoderEngine).setNoPrefix,
	awaitingMarkerID:        (*encoderEngine).setNoPrefix,
	awaitingReferenceID:     (*encoderEngine).setNoPrefix,
	awaitingQuotedString:    (*encoderEngine).setNoPrefix,
	awaitingVerbatimString:  (*encoderEngine).setNoPrefix,
	awaitingURI:             (*encoderEngine).setNoPrefix,
	awaitingCustomBinary:    (*encoderEngine).setNoPrefix,
	awaitingCustomText:      (*encoderEngine).setNoPrefix,
	awaitingArrayBool:       (*encoderEngine).setNoPrefix,
	awaitingArrayU8:         (*encoderEngine).setNoPrefix,
	awaitingArrayU16:        (*encoderEngine).setNoPrefix,
	awaitingArrayU32:        (*encoderEngine).setNoPrefix,
	awaitingArrayU64:        (*encoderEngine).setNoPrefix,
	awaitingArrayI8:         (*encoderEngine).setNoPrefix,
	awaitingArrayI16:        (*encoderEngine).setNoPrefix,
	awaitingArrayI32:        (*encoderEngine).setNoPrefix,
	awaitingArrayI64:        (*encoderEngine).setNoPrefix,
	awaitingArrayF16:        (*encoderEngine).setNoPrefix,
	awaitingArrayF32:        (*encoderEngine).setNoPrefix,
	awaitingArrayF64:        (*encoderEngine).setNoPrefix,
	awaitingArrayUUID:       (*encoderEngine).setNoPrefix,
}

var prefixSettersPretty = [awaitingCount]prefixSetter{
	awaitingVersion:         (*encoderEngine).setNoPrefix,
	awaitingTLO:             (*encoderEngine).setNoPrefix,
	awaitingListFirstItem:   (*encoderEngine).setIndentPrefix,
	awaitingListItem:        (*encoderEngine).setIndentPrefix,
	awaitingMapFirstKey:     (*encoderEngine).setIndentPrefix,
	awaitingMapKey:          (*encoderEngine).setIndentPrefix,
	awaitingMapValue:        (*encoderEngine).setNoPrefix,
	awaitingMetaFirstKey:    (*encoderEngine).setIndentPrefix,
	awaitingMetaKey:         (*encoderEngine).setIndentPrefix,
	awaitingMetaValue:       (*encoderEngine).setNoPrefix,
	awaitingMarkupName:      (*encoderEngine).setNoPrefix,
	awaitingMarkupFirstKey:  (*encoderEngine).setSpacePrefix,
	awaitingMarkupKey:       (*encoderEngine).setSpacePrefix,
	awaitingMarkupValue:     (*encoderEngine).setNoPrefix,
	awaitingMarkupFirstItem: (*encoderEngine).setIndentPrefix,
	awaitingMarkupItem:      (*encoderEngine).setNoPrefix,
	awaitingCommentItem:     (*encoderEngine).setNoPrefix,
	awaitingMarkerID:        (*encoderEngine).setNoPrefix,
	awaitingReferenceID:     (*encoderEngine).setNoPrefix,
	awaitingQuotedString:    (*encoderEngine).setNoPrefix,
	awaitingVerbatimString:  (*encoderEngine).setNoPrefix,
	awaitingURI:             (*encoderEngine).setNoPrefix,
	awaitingCustomBinary:    (*encoderEngine).setNoPrefix,
	awaitingCustomText:      (*encoderEngine).setNoPrefix,
	awaitingArrayBool:       (*encoderEngine).setNoPrefix,
	awaitingArrayU8:         (*encoderEngine).setNoPrefix,
	awaitingArrayU16:        (*encoderEngine).setNoPrefix,
	awaitingArrayU32:        (*encoderEngine).setNoPrefix,
	awaitingArrayU64:        (*encoderEngine).setNoPrefix,
	awaitingArrayI8:         (*encoderEngine).setNoPrefix,
	awaitingArrayI16:        (*encoderEngine).setNoPrefix,
	awaitingArrayI32:        (*encoderEngine).setNoPrefix,
	awaitingArrayI64:        (*encoderEngine).setNoPrefix,
	awaitingArrayF16:        (*encoderEngine).setNoPrefix,
	awaitingArrayF32:        (*encoderEngine).setNoPrefix,
	awaitingArrayF64:        (*encoderEngine).setNoPrefix,
	awaitingArrayUUID:       (*encoderEngine).setNoPrefix,
}

var containerTerminators = [awaitingCount]string{
	awaitingListFirstItem:   "]",
	awaitingListItem:        "]",
	awaitingMapFirstKey:     "}",
	awaitingMapKey:          "}",
	awaitingMetaFirstKey:    ")",
	awaitingMetaKey:         ")",
	awaitingMarkupFirstItem: ">",
	awaitingMarkupItem:      ">",
	awaitingCommentItem:     "*/",
}

var transitionFromSuffixesCompact = [awaitingCount]string{
	awaitingMapFirstKey:    "=",
	awaitingMapKey:         "=",
	awaitingMetaFirstKey:   "=",
	awaitingMetaKey:        "=",
	awaitingMarkupFirstKey: "=",
	awaitingMarkupKey:      "=",
}

var transitionFromSuffixesPretty = [awaitingCount]string{
	awaitingMapFirstKey:    " = ",
	awaitingMapKey:         " = ",
	awaitingMetaFirstKey:   " = ",
	awaitingMetaKey:        " = ",
	awaitingMarkupFirstKey: "=", // No spaces for markup k-v pairs
	awaitingMarkupKey:      "=", // No spaces for markup k-v pairs
}
