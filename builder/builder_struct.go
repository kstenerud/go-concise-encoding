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

package builder

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/options"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

type structBuilderDesc struct {
	Name      string
	Builder   ObjectBuilder
	Index     int
	Omit      bool
	OmitEmpty bool
	OmitValue string
}

func (_this *structBuilderDesc) applyTags(tags string) {
	if tags == "" {
		return
	}

	requiresValue := func(kv []string, key string) {
		if len(kv) != 2 {
			panic(fmt.Errorf(`tag key "%s" requires a value`, key))
		}
	}

	for _, entry := range strings.Split(tags, ",") {
		kv := strings.Split(entry, "=")
		switch strings.TrimSpace(kv[0]) {
		// TODO: lossy/nolossy
		// TODO: lowercase/origcase
		case "-":
			_this.Omit = true
		case "omit":
			if len(kv) == 1 {
				_this.Omit = true
			} else {
				_this.OmitValue = strings.TrimSpace(kv[1])
			}
		case "omitempty":
			// TODO: Implement omitempty
			_this.OmitEmpty = true
		case "name":
			requiresValue(kv, "name")
			_this.Name = strings.TrimSpace(kv[1])
		default:
			panic(fmt.Errorf("%v: Unknown Concise Encoding struct tag field", entry))
		}
	}
}

type structBuilder struct {
	// Template Data
	dstType       reflect.Type
	builderDescs  map[string]*structBuilderDesc
	nameBuilder   ObjectBuilder
	ignoreBuilder ObjectBuilder

	// Instance Data
	root   *RootBuilder
	parent ObjectBuilder
	opts   *options.BuilderOptions

	// Variable data (must be reset)
	nextBuilder   ObjectBuilder
	container     reflect.Value
	nextValue     reflect.Value
	nextIsKey     bool
	nextIsIgnored bool
}

func newStructBuilder(dstType reflect.Type) ObjectBuilder {
	return &structBuilder{
		dstType: dstType,
	}
}

func (_this *structBuilder) String() string {
	return fmt.Sprintf("%v<%v>", reflect.TypeOf(_this), _this.dstType)
}

func (_this *structBuilder) InitTemplate(session *Session) {
	_this.nameBuilder = session.GetBuilderForType(reflect.TypeOf(""))
	_this.builderDescs = make(map[string]*structBuilderDesc)
	_this.ignoreBuilder = newIgnoreBuilder()
	for i := 0; i < _this.dstType.NumField(); i++ {
		field := _this.dstType.Field(i)
		if field.PkgPath == "" {
			builder := session.GetBuilderForType(field.Type)
			desc := &structBuilderDesc{
				Name:    field.Name,
				Builder: builder,
				Index:   i,
			}
			desc.applyTags(field.Tag.Get("ce"))
			_this.builderDescs[desc.Name] = desc
		}
	}
}

func (_this *structBuilder) NewInstance(root *RootBuilder, parent ObjectBuilder, opts *options.BuilderOptions) ObjectBuilder {
	builderDescs := _this.builderDescs
	if opts.CaseInsensitiveStructFieldNames {
		builderDescs = make(map[string]*structBuilderDesc)
		for name, desc := range _this.builderDescs {
			builderDescs[common.ASCIIToLower(name)] = desc
		}
	}

	that := &structBuilder{
		dstType:      _this.dstType,
		builderDescs: builderDescs,
		parent:       parent,
		root:         root,
		opts:         opts,
	}
	that.nameBuilder = _this.nameBuilder.NewInstance(root, that, opts)
	that.ignoreBuilder = _this.ignoreBuilder.NewInstance(root, that, opts)
	that.reset()
	return that
}

func (_this *structBuilder) SetParent(parent ObjectBuilder) {
	_this.parent = parent
}

func (_this *structBuilder) reset() {
	_this.nextBuilder = _this.nameBuilder
	_this.container = reflect.New(_this.dstType).Elem()
	_this.nextValue = reflect.Value{}
	_this.nextIsKey = true
	_this.nextIsIgnored = false
}

func (_this *structBuilder) swapKeyValue() {
	_this.nextIsKey = !_this.nextIsKey
}

func (_this *structBuilder) BuildFromNil(_ reflect.Value) {
	_this.nextBuilder.BuildFromNil(_this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBool(value bool, _ reflect.Value) {
	_this.nextBuilder.BuildFromBool(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromInt(value int64, _ reflect.Value) {
	_this.nextBuilder.BuildFromInt(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromUint(value uint64, _ reflect.Value) {
	_this.nextBuilder.BuildFromUint(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigInt(value *big.Int, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigInt(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromFloat(value float64, _ reflect.Value) {
	_this.nextBuilder.BuildFromFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigFloat(value *big.Float, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromDecimalFloat(value compact_float.DFloat, _ reflect.Value) {
	_this.nextBuilder.BuildFromDecimalFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromBigDecimalFloat(value *apd.Decimal, _ reflect.Value) {
	_this.nextBuilder.BuildFromBigDecimalFloat(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromUUID(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromUUID(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromString(value []byte, _ reflect.Value) {
	if _this.nextIsKey {
		if _this.opts.CaseInsensitiveStructFieldNames {
			common.ASCIIBytesToLower(value)
		}
		name := string(value)

		if builderDesc, ok := _this.builderDescs[name]; ok {
			_this.nextBuilder = builderDesc.Builder
			_this.nextValue = _this.container.Field(builderDesc.Index)
		} else {
			_this.root.SetCurrentBuilder(_this.ignoreBuilder)
			_this.nextBuilder = _this.ignoreBuilder
			_this.nextIsIgnored = true
			return
		}
	} else {
		_this.nextBuilder.BuildFromString(value, _this.nextValue)
	}

	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromVerbatimString(value []byte, _ reflect.Value) {
	if _this.nextIsKey {
		if builderDesc, ok := _this.builderDescs[string(value)]; ok {
			_this.nextBuilder = builderDesc.Builder
			_this.nextValue = _this.container.Field(builderDesc.Index)
		} else {
			_this.root.SetCurrentBuilder(_this.ignoreBuilder)
			_this.nextBuilder = _this.ignoreBuilder
			_this.nextIsIgnored = true
			return
		}
	} else {
		_this.nextBuilder.BuildFromVerbatimString(value, _this.nextValue)
	}

	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromURI(value *url.URL, _ reflect.Value) {
	_this.nextBuilder.BuildFromURI(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCustomBinary(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromCustomBinary(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCustomText(value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromCustomText(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromTypedArray(elemType reflect.Type, value []byte, _ reflect.Value) {
	_this.nextBuilder.BuildFromTypedArray(elemType, value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromTime(value time.Time, _ reflect.Value) {
	_this.nextBuilder.BuildFromTime(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildFromCompactTime(value *compact_time.Time, _ reflect.Value) {
	_this.nextBuilder.BuildFromCompactTime(value, _this.nextValue)
	_this.swapKeyValue()
}

func (_this *structBuilder) BuildBeginList() {
	_this.nextBuilder.PrepareForListContents()
}

func (_this *structBuilder) BuildBeginMap() {
	_this.nextBuilder.PrepareForMapContents()
}

func (_this *structBuilder) BuildEndContainer() {
	object := _this.container
	_this.reset()
	_this.parent.NotifyChildContainerFinished(object)
}

func (_this *structBuilder) BuildBeginMarker(id interface{}) {
	root := _this.root
	_this.nextBuilder = newMarkerObjectBuilder(_this, _this.nextBuilder, func(object reflect.Value) {
		root.NotifyMarker(id, object)
	})
}

func (_this *structBuilder) BuildFromReference(id interface{}) {
	nextValue := _this.nextValue
	_this.swapKeyValue()
	_this.root.NotifyReference(id, func(object reflect.Value) {
		setAnythingFromAnything(object, nextValue)
	})
}

func (_this *structBuilder) PrepareForListContents() {
	PanicBadEventWithType(_this, _this.dstType, "PrepareForListContents")
}

func (_this *structBuilder) PrepareForMapContents() {
	builderDescs := make(map[string]*structBuilderDesc)

	for k, builderElem := range _this.builderDescs {
		builderDescs[k] = &structBuilderDesc{
			Builder: builderElem.Builder.NewInstance(_this.root, _this, _this.opts),
			Index:   builderElem.Index,
		}
	}
	_this.builderDescs = builderDescs
	_this.root.SetCurrentBuilder(_this)
}

func (_this *structBuilder) NotifyChildContainerFinished(value reflect.Value) {
	_this.root.SetCurrentBuilder(_this)
	if _this.nextIsIgnored {
		_this.nextIsIgnored = false
		return
	}

	_this.nextValue.Set(value)
	_this.swapKeyValue()
}
