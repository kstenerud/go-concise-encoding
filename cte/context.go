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

	"github.com/kstenerud/go-concise-encoding/options"
)

type EncoderContext struct {
	opts        options.CTEEncoderOptions
	DepthPrefix []byte

	// Stack
	CurrentPrefixer Prefixer
	prefixerStack   []Prefixer

	Stream EncodeBuffer

	// Arrays
	// arrayType              events.ArrayType
	// moreChunksFollow       bool
	// builtArrayBuffer       []byte
	// arrayMaxByteCount      uint64
	// arrayTotalByteCount    uint64
	// chunkExpectedByteCount uint64
	// chunkActualByteCount   uint64
	// utf8RemainderBacking   [4]byte
	// utf8RemainderBuffer    []byte
	// ValidateArrayDataFunc  func(data []byte)

	// // Marker/Reference
	// currentMarkerID   interface{}
	// markerObjectRule  EventRule
	// markedObjects     map[interface{}]DataType
	// forwardReferences map[interface{}]DataType
	// referenceCount    uint64
}

func (_this *EncoderContext) Init(opts *options.CTEEncoderOptions) {
	_this.opts = *opts
	_this.DepthPrefix = make([]byte, 0, 65)
	_this.DepthPrefix = append(_this.DepthPrefix, '\n')
	_this.stackPrefixer(&topLevelPrefixer)
}

func (_this *EncoderContext) Reset() {
}

func increaseDepthPrefix(ctx *EncoderContext) {
	ctx.DepthPrefix = append(ctx.DepthPrefix, ' ', ' ', ' ', ' ')
}

func decreaseDepthPrefix(ctx *EncoderContext) {
	ctx.DepthPrefix = ctx.DepthPrefix[:len(ctx.DepthPrefix)-4]
}

func (_this *EncoderContext) stackPrefixer(prefixer Prefixer) {
	_this.prefixerStack = append(_this.prefixerStack, prefixer)
	_this.CurrentPrefixer = prefixer
}

func (_this *EncoderContext) unstackPrefixer() {
	_this.prefixerStack = _this.prefixerStack[:len(_this.prefixerStack)-1]
	_this.CurrentPrefixer = _this.prefixerStack[len(_this.prefixerStack)-1]
}

func (_this *EncoderContext) changePrefixer(prefixer Prefixer) {
	_this.prefixerStack[len(_this.prefixerStack)-1] = prefixer
	_this.CurrentPrefixer = prefixer
}

func (_this *EncoderContext) ApplyPrefix() {
	_this.CurrentPrefixer.ApplyPrefix(_this)
}

func (_this *EncoderContext) NotifyObject() {
	_this.CurrentPrefixer.NotifyObject(_this)
}

func (_this *EncoderContext) EndContainer() {
	_this.CurrentPrefixer.End(_this)
	_this.unstackPrefixer()
}

func (_this *EncoderContext) BeginList() {
	_this.stackPrefixer(&listPrefixer)
	_this.CurrentPrefixer.Begin(_this)
}

func (_this *EncoderContext) BeginMap() {
	_this.stackPrefixer(&mapKeyPrefixer)
	_this.CurrentPrefixer.Begin(_this)
}

func (_this *EncoderContext) BeginMarkup() {
}

func (_this *EncoderContext) BeginMetadata() {
}

func (_this *EncoderContext) BeginComment() {
}

func (_this *EncoderContext) BeginMarker() {
}

func (_this *EncoderContext) BeginReference() {
}

func (_this *EncoderContext) BeginConcatenate() {
}

func (_this *EncoderContext) BeginNA() {
}

// pre-write (indent)
// post-write (lf?)
// list-type
// map-type (key section, value section)
// array-type
// metadata follow
// comment follow

// string in comment
// string in markup contents
// string as ID
// int as ID
// constant name
// custom string
// NA stuff

// Types that can be printed differently:
// - string (quoted, unquoted) (escape/noescape) (trim/notrim)
// - int (pos or neg) (bin, oct, dec, hex) (with/without prefix)
// - binary float (dec, hex) (with/without prefix)

type Prefixer interface {
	Begin(ctx *EncoderContext)
	End(ctx *EncoderContext)
	NotifyObject(ctx *EncoderContext)
	ApplyPrefix(ctx *EncoderContext)
}

var (
	topLevelPrefixer TopLevelPrefixer
	listPrefixer     ListPrefixer
	mapKeyPrefixer   MapKeyPrefixer
	mapValuePrefixer MapValuePrefixer
)

type TopLevelPrefixer struct{}

func (_this *TopLevelPrefixer) Begin(ctx *EncoderContext) {
	panic(fmt.Errorf("BUG: TopLevelPrefixer cannot respond to Begin"))
}

func (_this *TopLevelPrefixer) End(ctx *EncoderContext) {
	panic(fmt.Errorf("BUG: TopLevelPrefixer cannot respond to End"))
}

func (_this *TopLevelPrefixer) NotifyObject(ctx *EncoderContext) {
	// Nothing to do
}

func (_this *TopLevelPrefixer) ApplyPrefix(ctx *EncoderContext) {
	ctx.Stream.AddBytes(ctx.DepthPrefix)
}

type ListPrefixer struct{}

func (_this *ListPrefixer) Begin(ctx *EncoderContext) {
	ctx.Stream.AddByte('[')
	increaseDepthPrefix(ctx)
}

func (_this *ListPrefixer) End(ctx *EncoderContext) {
	decreaseDepthPrefix(ctx)
	ctx.Stream.AddByte(']')
}

func (_this *ListPrefixer) NotifyObject(ctx *EncoderContext) {
	// Nothing to do
}

func (_this *ListPrefixer) ApplyPrefix(ctx *EncoderContext) {
	ctx.Stream.AddBytes(ctx.DepthPrefix)
}

type MapKeyPrefixer struct{}

func (_this *MapKeyPrefixer) Begin(ctx *EncoderContext) {
	ctx.Stream.AddByte('{')
	increaseDepthPrefix(ctx)
}

func (_this *MapKeyPrefixer) End(ctx *EncoderContext) {
	decreaseDepthPrefix(ctx)
	ctx.Stream.AddByte('}')
}

func (_this *MapKeyPrefixer) NotifyObject(ctx *EncoderContext) {
	ctx.changePrefixer(&mapValuePrefixer)
}

func (_this *MapKeyPrefixer) ApplyPrefix(ctx *EncoderContext) {
	ctx.Stream.AddBytes(ctx.DepthPrefix)
}

type MapValuePrefixer struct{}

func (_this *MapValuePrefixer) Begin(ctx *EncoderContext) {
	panic(fmt.Errorf("BUG: MapValuePrefixer cannot respond to Begin"))
}

func (_this *MapValuePrefixer) End(ctx *EncoderContext) {
	panic(fmt.Errorf("BUG: MapValuePrefixer cannot respond to End"))
}

func (_this *MapValuePrefixer) NotifyObject(ctx *EncoderContext) {
	ctx.changePrefixer(&mapKeyPrefixer)
}

func (_this *MapValuePrefixer) ApplyPrefix(ctx *EncoderContext) {
	ctx.Stream.AddBytes([]byte{' ', '=', ' '})
}

// =============================================================================

type ArrayRenderer interface {
	RenderArrayPortion(ctx *EncoderContext, data []byte)
	RenderArrayComplete(ctx *EncoderContext, data []byte)
}

var (
	arrayRendererNone ArrayRendererNone
)

type ArrayRendererNone struct{}

func (_this *ArrayRendererNone) RenderArrayPortion(ctx *EncoderContext, data []byte) {
	panic(fmt.Errorf("BUG: ArrayRendererNone cannot respond to RenderArrayPortion"))
}

func (_this *ArrayRendererNone) RenderArrayComplete(ctx *EncoderContext, data []byte) {
	panic(fmt.Errorf("BUG: ArrayRendererNone cannot respond to RenderArrayComplete"))
}
