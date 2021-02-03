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
	"strings"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/kstenerud/go-concise-encoding/events"
)

type contextStackEntry struct {
	Rule     EventRule
	DataType DataType
}

type Context struct {
	opts            options.RuleOptions
	ExpectedVersion uint64

	objectCount uint64

	// Stack
	CurrentEntry   contextStackEntry
	stack          []contextStackEntry
	containerDepth uint64

	// Arrays
	arrayType              events.ArrayType
	moreChunksFollow       bool
	builtArrayBuffer       []byte
	arrayMaxByteCount      uint64
	arrayTotalByteCount    uint64
	chunkExpectedByteCount uint64
	chunkActualByteCount   uint64
	utf8RemainderBacking   [4]byte
	utf8RemainderBuffer    []byte
	ValidateArrayDataFunc  func(data []byte)

	// Marker/Reference
	currentMarkerID   interface{}
	markerObjectRule  EventRule
	markedObjects     map[interface{}]DataType
	forwardReferences map[interface{}]DataType
	referenceCount    uint64
}

func (_this *Context) Init(version uint64, opts *options.RuleOptions) {
	opts = opts.WithDefaultsApplied()
	_this.opts = *opts
	_this.ExpectedVersion = version
	_this.stack = make([]contextStackEntry, 0, 16)
	_this.Reset()
}

func (_this *Context) Reset() {
	_this.objectCount = 0
	_this.containerDepth = 0
	_this.referenceCount = 0
	_this.stack = _this.stack[:0]
	if _this.markedObjects == nil || len(_this.markedObjects) > 0 {
		_this.markedObjects = make(map[interface{}]DataType)
	}
	if _this.forwardReferences == nil || len(_this.forwardReferences) > 0 {
		_this.forwardReferences = make(map[interface{}]DataType)
	}
	_this.stackRule(&beginDocumentRule, DataTypeInvalid)
}

func (_this *Context) changeRule(rule EventRule) {
	_this.CurrentEntry.Rule = rule
	_this.stack[len(_this.stack)-1] = _this.CurrentEntry
}

func (_this *Context) stackRule(rule EventRule, dataType DataType) {
	_this.CurrentEntry = contextStackEntry{
		Rule:     rule,
		DataType: dataType,
	}
	_this.stack = append(_this.stack, _this.CurrentEntry)
}

func (_this *Context) UnstackRule() EventRule {
	unstackedRule := _this.CurrentEntry.Rule
	_this.stack = _this.stack[:len(_this.stack)-1]
	_this.CurrentEntry = _this.stack[len(_this.stack)-1]
	return unstackedRule
}

func (_this *Context) ParentRule() EventRule {
	return _this.stack[len(_this.stack)-2].Rule
}

func (_this *Context) NotifyNewObject() {
	_this.objectCount++
	if _this.objectCount > _this.opts.MaxObjectCount {
		panic(fmt.Errorf("Exceeded max object count of %d", _this.opts.MaxObjectCount))
	}
}

func (_this *Context) beginContainer(rule EventRule, dataType DataType) {
	_this.containerDepth++
	if _this.containerDepth > _this.opts.MaxContainerDepth {
		panic(fmt.Errorf("Exceeded max container depth of %d", _this.opts.MaxContainerDepth))
	}
	_this.stackRule(rule, dataType)
}

func (_this *Context) endContainerLike() {
	cType := _this.CurrentEntry.DataType
	_this.UnstackRule()
	_this.CurrentEntry.Rule.OnChildContainerEnded(_this, cType)
}

func (_this *Context) EndContainer() {
	_this.containerDepth--
	_this.endContainerLike()
}

func (_this *Context) BeginNA() {
	_this.stackRule(&naRule, DataTypeAnyType)
}

func (_this *Context) BeginList() {
	_this.beginContainer(&listRule, DataTypeAnyType)
}

func (_this *Context) BeginMap() {
	_this.beginContainer(&mapKeyRule, DataTypeAnyType)
}

func (_this *Context) BeginMarkup() {
	_this.beginContainer(&markupNameRule, DataTypeAnyType)
}

func (_this *Context) BeginMetadata() {
	_this.beginContainer(&metadataKeyRule, DataTypeInvalid)
}

func (_this *Context) BeginComment() {
	_this.beginContainer(&commentRule, DataTypeInvalid)
}

func (_this *Context) assertReferenceCount() {
	newReferenceCount := _this.referenceCount + 1
	if newReferenceCount > _this.opts.MaxReferenceCount {
		panic(fmt.Errorf("Too many marked objects (%d). Max is %d", newReferenceCount, _this.opts.MaxReferenceCount))
	}
}

func (_this *Context) BeginMarkerKeyable() {
	_this.assertReferenceCount()
	_this.stackRule(&markerIDKeyableRule, DataTypeKeyable)
}

func (_this *Context) BeginMarkerAnyType() {
	_this.assertReferenceCount()
	_this.stackRule(&markerIDAnyTypeRule, DataTypeAnyType)
}

func (_this *Context) BeginMarkedObjectKeyable(id interface{}) {
	_this.currentMarkerID = id
	_this.changeRule(&markedObjectKeyableRule)
}

func (_this *Context) BeginMarkedObjectAnyType(id interface{}) {
	_this.currentMarkerID = id
	_this.changeRule(&markedObjectAnyTypeRule)
}

func (_this *Context) BeginReferenceKeyable() {
	_this.stackRule(&referenceKeyableRule, DataTypeKeyable)
}

func (_this *Context) BeginReferenceAnyType() {
	_this.stackRule(&referenceAnyTypeRule, DataTypeAnyType)
}

func (_this *Context) BeginTopLevelReference() {
	_this.stackRule(&tlReferenceRIDRule, DataTypeKeyable)
}

func (_this *Context) SwitchVersion() {
	_this.changeRule(&versionRule)
}

func (_this *Context) SwitchTopLevel() {
	_this.changeRule(&topLevelRule)
}

func (_this *Context) SwitchEndDocument() {
	_this.changeRule(&endDocumentRule)
}

func (_this *Context) EndDocument() {
	if len(_this.forwardReferences) > 0 {
		var sb strings.Builder
		sb.WriteString("Forward references have not been resolved: [")
		for id := range _this.forwardReferences {
			sb.WriteString(fmt.Sprintf("%v, ", id))
		}

		str := sb.String()
		str = str[:len(str)-2]
		panic(fmt.Errorf("%v]", str))
	}
	_this.changeRule(&terminalRule)
}

func (_this *Context) SwitchMapKey() {
	_this.changeRule(&mapKeyRule)
}

func (_this *Context) SwitchMapValue() {
	_this.changeRule(&mapValueRule)
}

func (_this *Context) SwitchMetadataKey() {
	_this.changeRule(&metadataKeyRule)
}

func (_this *Context) SwitchMetadataValue() {
	_this.changeRule(&metadataValueRule)
}

func (_this *Context) SwitchMetadataCompletion() {
	_this.changeRule(&metadataCompleteRule)
}

func (_this *Context) SwitchMarkupKey() {
	_this.changeRule(&markupKeyRule)
}

func (_this *Context) SwitchMarkupValue() {
	_this.changeRule(&markupValueRule)
}

func (_this *Context) SwitchMarkupContents() {
	_this.changeRule(&markupContentsRule)
}

func (_this *Context) MarkObject(dataType DataType) {
	if _, exists := _this.markedObjects[_this.currentMarkerID]; exists {
		panic(fmt.Errorf("Marker ID [%v] already exists", _this.currentMarkerID))
	}
	_this.referenceCount++
	_this.markedObjects[_this.currentMarkerID] = dataType
	if allowedDataTypes, exists := _this.forwardReferences[_this.currentMarkerID]; exists {
		delete(_this.forwardReferences, _this.currentMarkerID)
		if allowedDataTypes&dataType == 0 {
			panic(fmt.Errorf("Forward reference to marker ID [%v] cannot accept type %v", _this.currentMarkerID, dataType))
		}
	}
}

func (_this *Context) ReferenceObject(id interface{}, allowedDataTypes DataType) {
	if dataType, exists := _this.markedObjects[id]; exists {
		if dataType&allowedDataTypes == 0 {
			panic(fmt.Errorf("Marked object id [%v] of type %v is not a valid type to be referenced here", id, dataType))
		}
		return
	}

	current := _this.forwardReferences[id]
	if current == 0 {
		current = allowedDataTypes
	} else {
		current &= allowedDataTypes
	}
	_this.forwardReferences[id] = current
}
