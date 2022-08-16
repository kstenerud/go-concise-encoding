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

	"github.com/kstenerud/go-concise-encoding/events"
	"github.com/kstenerud/go-concise-encoding/options"
	"github.com/kstenerud/go-concise-encoding/version"
)

const noObjectCount = -1

type contextStackEntry struct {
	Rule                EventRule
	DataType            DataType
	CurrentObjectCount  int
	ExpectedObjectCount int // -1 means ignored
}

type Context struct {
	opts            *options.RuleOptions
	ExpectedVersion uint64

	objectCount uint64

	structTemplates    map[string]int
	structTemplateName string

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
	markerID               string
	markedObjects          map[interface{}]DataType
	forwardLocalReferences map[interface{}]DataType
	LocalReferenceCount    uint64
}

func (_this *Context) Init(opts *options.RuleOptions) {
	_this.opts = opts
	_this.ExpectedVersion = version.ConciseEncodingVersion
	_this.stack = make([]contextStackEntry, 0, 16)
	_this.Reset()
}

func (_this *Context) Reset() {
	_this.objectCount = 0
	_this.containerDepth = 0
	_this.LocalReferenceCount = 0
	_this.stack = _this.stack[:0]
	_this.structTemplates = make(map[string]int)
	if _this.markedObjects == nil || len(_this.markedObjects) > 0 {
		_this.markedObjects = make(map[interface{}]DataType)
	}
	if _this.forwardLocalReferences == nil || len(_this.forwardLocalReferences) > 0 {
		_this.forwardLocalReferences = make(map[interface{}]DataType)
	}
	_this.CurrentEntry = contextStackEntry{
		Rule:                &beginDocumentRule,
		DataType:            DataTypeInvalid,
		ExpectedObjectCount: noObjectCount,
		CurrentObjectCount:  0,
	}
}

func (_this *Context) ChangeRule(rule EventRule) {
	_this.CurrentEntry.Rule = rule
}

/**
 * Use expectedObjectCount == -1 to ignore it.
 */
func (_this *Context) stackRule(rule EventRule, dataType DataType, expectedObjectCount int) {
	_this.stack = append(_this.stack, _this.CurrentEntry)
	_this.CurrentEntry = contextStackEntry{
		Rule:                rule,
		DataType:            dataType,
		ExpectedObjectCount: expectedObjectCount,
		CurrentObjectCount:  0,
	}
}

func (_this *Context) UnstackRule() EventRule {
	unstackedRule := _this.CurrentEntry.Rule
	_this.CurrentEntry = _this.stack[len(_this.stack)-1]
	_this.stack = _this.stack[:len(_this.stack)-1]
	return unstackedRule
}

func (_this *Context) ParentRule() EventRule {
	return _this.stack[len(_this.stack)-1].Rule
}

func (_this *Context) NotifyNewObject(isRealObject bool) {
	if isRealObject {
		_this.CurrentEntry.CurrentObjectCount++
		if _this.CurrentEntry.ExpectedObjectCount >= 0 && _this.CurrentEntry.CurrentObjectCount > _this.CurrentEntry.ExpectedObjectCount {
			panic(fmt.Errorf("container exceeds expected object count of %d", _this.CurrentEntry.ExpectedObjectCount))
		}
	}
	_this.objectCount++
	if _this.objectCount > _this.opts.MaxObjectCount {
		panic(fmt.Errorf("exceeded max object count of %d", _this.opts.MaxObjectCount))
	}
}

/**
 * Use expectedObjectCount == -1 to ignore it.
 */
func (_this *Context) beginContainer(rule EventRule, dataType DataType, expectedObjectCount int) {
	_this.containerDepth++
	if _this.containerDepth > _this.opts.MaxContainerDepth {
		panic(fmt.Errorf("exceeded max container depth of %d", _this.opts.MaxContainerDepth))
	}
	_this.stackRule(rule, dataType, expectedObjectCount)
}

func (_this *Context) endContainerLike(notifyParent bool) {
	cType := _this.CurrentEntry.DataType
	_this.UnstackRule()
	if notifyParent {
		_this.CurrentEntry.Rule.OnChildContainerEnded(_this, cType)
	}
}

func (_this *Context) EndContainer(notifyParent bool) {
	if _this.containerDepth == 0 {
		panic("BUG: Too many end container calls")
	}
	if _this.CurrentEntry.ExpectedObjectCount >= 0 && _this.CurrentEntry.CurrentObjectCount != _this.CurrentEntry.ExpectedObjectCount {
		panic(fmt.Errorf("container has %v objects but expected object count of %d", _this.CurrentEntry.CurrentObjectCount, _this.CurrentEntry.ExpectedObjectCount))
	}
	if _this.CurrentEntry.DataType == DataTypeStructTemplate {
		_this.addTemplate(_this.structTemplateName, _this.CurrentEntry.CurrentObjectCount)
	}
	_this.containerDepth--
	_this.endContainerLike(notifyParent)
}

func (_this *Context) addTemplate(id string, objectCount int) {
	if _, exists := _this.structTemplates[id]; exists {
		panic(fmt.Errorf("struct template ID [%v] already exists", id))
	}

	_this.structTemplates[id] = objectCount
}

func (_this *Context) BeginList() {
	_this.beginContainer(&listRule, DataTypeList, noObjectCount)
}

func (_this *Context) BeginMap() {
	_this.beginContainer(&mapKeyRule, DataTypeMap, noObjectCount)
}

func (_this *Context) BeginStructTemplate(id []byte) {
	_this.beginContainer(&structTemplateRule, DataTypeStructTemplate, noObjectCount)
	_this.structTemplateName = string(id)
}

func (_this *Context) BeginStructInstance(id []byte) {
	expectedObjectCount, ok := _this.structTemplates[string(id)]
	if !ok {
		panic(fmt.Errorf("%v: no such struct template has been defined", string(id)))
	}
	_this.beginContainer(&structInstanceRule, DataTypeStructInstance, expectedObjectCount)
}

func (_this *Context) BeginEdge() {
	_this.beginContainer(&edgeSourceRule, DataTypeEdge, 3)
}

func (_this *Context) BeginNode() {
	_this.beginContainer(&nodeRule, DataTypeList, noObjectCount)
}

func (_this *Context) BeginMarkerKeyable(id []byte, dataType DataType) {
	_this.markerID = string(id)
	_this.stackRule(&markedObjectKeyableRule, dataType, noObjectCount)
}

func (_this *Context) BeginMarkerAnyType(id []byte, dataType DataType) {
	_this.markerID = string(id)
	_this.stackRule(&markedObjectAnyTypeRule, dataType, noObjectCount)
}

func (_this *Context) LocalReferenceKeyable(identifier []byte) {
	_this.LocalReferenceObject(identifier, AllowKeyable)
}

func (_this *Context) LocalReferenceAnyType(identifier []byte) {
	_this.LocalReferenceObject(identifier, AllowAny)
}

func (_this *Context) EndDocument() {
	if len(_this.forwardLocalReferences) > 0 {
		var sb strings.Builder
		sb.WriteString("Forward local references have not been resolved: [")
		for id := range _this.forwardLocalReferences {
			sb.WriteString(fmt.Sprintf("%v, ", id))
		}

		str := sb.String()
		str = str[:len(str)-2]
		panic(fmt.Errorf("%v]", str))
	}
	_this.ChangeRule(&terminalRule)
}

func (_this *Context) MarkObject(dataType DataType) {
	newLocalReferenceCount := _this.LocalReferenceCount + 1
	if newLocalReferenceCount > _this.opts.MaxLocalReferenceCount {
		panic(fmt.Errorf("too many marked objects (%d). Max is %d", newLocalReferenceCount, _this.opts.MaxLocalReferenceCount))
	}

	id := _this.markerID
	if _, exists := _this.markedObjects[id]; exists {
		panic(fmt.Errorf("marker ID [%v] already exists", id))
	}
	_this.LocalReferenceCount++
	_this.markedObjects[id] = dataType
	if allowedDataTypes, exists := _this.forwardLocalReferences[id]; exists {
		delete(_this.forwardLocalReferences, id)
		if allowedDataTypes&dataType == 0 {
			panic(fmt.Errorf("forward reference to marker ID [%v] cannot accept type %v", id, dataType))
		}
	}
}

func (_this *Context) LocalReferenceObject(id []byte, allowedDataTypes DataType) {
	idAsString := string(id)
	if dataType, exists := _this.markedObjects[idAsString]; exists {
		if dataType&allowedDataTypes == 0 {
			panic(fmt.Errorf("marked object id [%v] of type %v is not a valid type to be referenced here", idAsString, dataType))
		}
		return
	}

	current := _this.forwardLocalReferences[idAsString]
	if current == 0 {
		current = allowedDataTypes
	} else {
		current &= allowedDataTypes
	}
	_this.forwardLocalReferences[idAsString] = current
}
