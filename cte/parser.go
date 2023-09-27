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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	antlr "github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/ce/events"
	"github.com/kstenerud/go-concise-encoding/cte/parser"
	"github.com/kstenerud/go-concise-encoding/internal/common"
)

// See Antlr grammar at https://github.com/kstenerud/go-concise-encoding/tree/master/codegen/cte

func ParseDocument(document string, eventReceiver events.DataEventReceiver) error {
	errorListener := new(reportingErrorListener)

	is := antlr.NewInputStream(document)
	lexer := parser.NewContextualCTELexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewCTEParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)
	p.SetErrorHandler(new(bailErrorStrategy))

	listener := newCteListener(eventReceiver)

	antlr.ParseTreeWalkerDefault.Walk(listener, p.Cte())
	return errorListener.Error
}

type reportingErrorListener struct {
	*antlr.DefaultErrorListener
	Error error
}

func (_this *reportingErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	_this.Error = fmt.Errorf("line %v, col %v: %v", line, column, msg)
}

type bailErrorStrategy struct {
	antlr.DefaultErrorStrategy
}

func (b *bailErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {}
func (b *bailErrorStrategy) Sync(recognizer antlr.Parser)                                  {}

// --------------------------------------------------------------------

type cteListener struct {
	*parser.BaseCTEParserListener
	eventReceiver   events.DataEventReceiver
	elementSizeBits int
	customType      uint64
	mediaType       string
	arrayData       []uint8
	hexFloatRegex   *regexp.Regexp
}

func newCteListener(eventReceiver events.DataEventReceiver) *cteListener {
	return &cteListener{
		eventReceiver: eventReceiver,
		hexFloatRegex: regexp.MustCompile(`0[xX]([0-9a-fA-F]+)(\.([0-9a-fA-F]+))?`),
	}
}

func (_this *cteListener) wrapPanic(r interface{}, ctx *antlr.BaseParserRuleContext) {
	if r == nil {
		return
	}

	switch v := r.(type) {
	case error:
		panic(fmt.Errorf("line %v, col %v: %w", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()+1, v))
	default:
		panic(fmt.Errorf("line %v, col %v: %v", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()+1, v))
	}
}

func (_this *cteListener) errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

func (_this *cteListener) clearArrayData() {
	_this.arrayData = _this.arrayData[:0]
}

func (_this *cteListener) beginArray(elementSizeBits int) {
	_this.elementSizeBits = elementSizeBits
	_this.clearArrayData()
}

func (_this *cteListener) appendCodepoint(codepoint rune) {
	buff := bytes.NewBuffer(_this.arrayData)
	if _, err := buff.WriteRune(codepoint); err != nil {
		_this.errorf("error writing rune %x: %v", codepoint, err)
	}
	_this.arrayData = buff.Bytes()
}

func (_this *cteListener) EnterCte(ctx *parser.CteContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnBeginDocument()
}

func (_this *cteListener) ExitCte(ctx *parser.CteContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndDocument()
}

func (_this *cteListener) ExitVersion(ctx *parser.VersionContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	versionStr := ctx.GetText()
	if len(versionStr) < 2 {
		panic(fmt.Errorf("expected a version string"))
	}
	versionStr = versionStr[1:]
	_this.eventReceiver.OnVersion(parseSmallUint(versionStr))
}

func (_this *cteListener) ExitValueNull(ctx *parser.ValueNullContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnNull()
}

func (_this *cteListener) ExitValueUid(ctx *parser.ValueUidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnUID(appendUID(ctx.GetText(), make([]byte, 0, 16)))
}

func (_this *cteListener) ExitValueBool(ctx *parser.ValueBoolContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnBoolean(strings.ToLower(ctx.GetText()) == "true")
}

func (_this *cteListener) ExitValueInt(ctx *parser.ValueIntContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	str := ctx.GetText()
	str = strings.ReplaceAll(str, "_", "")
	isNegative := false
	if str[0] == '-' {
		isNegative = true
	}

	if v, err := strconv.ParseInt(str, 0, 64); err == nil {
		if v == 0 && isNegative {
			_this.eventReceiver.OnNegativeInt(0)
		} else {
			_this.eventReceiver.OnInt(v)
		}
		return
	}

	bigInt := &big.Int{}
	if _, success := bigInt.SetString(str, 0); success {
		_this.eventReceiver.OnBigInt(bigInt)
		return
	}

	panic(fmt.Errorf("BUG: Expected an integer but got \"%v\"", str))
}

func countFloatSignificantDigits(str string) (count uint) {
	for _, ch := range str {
		if ch != '0' {
			break
		}
		count++
	}
	str = str[count:]

	for _, ch := range str {
		if ch == '.' {
			continue
		}
		if ch == 'p' || ch == 'P' {
			break
		}
		count++
	}
	return count
}

func (_this *cteListener) ExitValueFloat(ctx *parser.ValueFloatContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	str := normalizeFloatString(ctx.GetText(), 0)
	strNoSign := str
	sign := 1.0
	if str[0] == '-' {
		// Split out negation because go doesn't handle -0 properly
		sign = -sign
		strNoSign = str[1:]
	}

	if strings.HasPrefix(strNoSign, "0x") || strings.HasPrefix(strNoSign, "0X") {
		significantBits := countFloatSignificantDigits(strNoSign[2:]) * 4

		if v, _, err := big.ParseFloat(strNoSign[2:], 16, significantBits, big.ToNearestEven); err == nil {
			if f64, accuracy := v.Float64(); accuracy == big.Exact {
				_this.eventReceiver.OnFloat(f64 * sign)
				return
			}

			if sign < 0 {
				v = v.Neg(v)
			}
			_this.eventReceiver.OnBigFloat(v)
			return
		} else {
			panic(err)
		}
	}

	if value, err := compact_float.DFloatFromString(str); err == nil {
		_this.eventReceiver.OnDecimalFloat(value)
		return
	}

	decimal, cond, err := apd.NewFromString(strNoSign)
	if err != nil {
		panic(err)
	}
	if cond != 0 {
		panic(fmt.Errorf("APD Condition %v", cond))
	}
	if sign < 0 {
		decimal = decimal.Neg(decimal)
	}
	_this.eventReceiver.OnBigDecimalFloat(decimal)
}

func (_this *cteListener) ExitValueInf(ctx *parser.ValueInfContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnDecimalFloat(compact_float.Infinity())
}

func (_this *cteListener) ExitValueNinf(ctx *parser.ValueNinfContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnDecimalFloat(compact_float.NegativeInfinity())
}

func (_this *cteListener) ExitValueNan(ctx *parser.ValueNanContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnDecimalFloat(compact_float.QuietNaN())
}

func (_this *cteListener) ExitValueSnan(ctx *parser.ValueSnanContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnDecimalFloat(compact_float.SignalingNaN())
}

func (_this *cteListener) ExitValueDate(ctx *parser.ValueDateContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnTime(parseDateOrDateTime(ctx.GetText()))
}

func (_this *cteListener) ExitValueTime(ctx *parser.ValueTimeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnTime(parseTime(ctx.GetText()))
}

func (_this *cteListener) ExitCodepointContents(ctx *parser.CodepointContentsContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()
	text = text[:len(text)-1]
	_this.appendCodepoint(parseHexCodepoint(text))
}

func (_this *cteListener) EnterValueRid(ctx *parser.ValueRidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) EnterValueRemoteRef(ctx *parser.ValueRemoteRefContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) EnterValueString(ctx *parser.ValueStringContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitValueString(ctx *parser.ValueStringContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeString, uint64(len(_this.arrayData)), _this.arrayData)
}

func (_this *cteListener) ExitStringContents(ctx *parser.StringContentsContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = append(_this.arrayData, []byte(ctx.GetText())...)
}

func (_this *cteListener) ExitVerbatimContents(ctx *parser.VerbatimContentsContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = append(_this.arrayData, []byte(ctx.GetText())...)
}

func (_this *cteListener) ExitEscapeChar(ctx *parser.EscapeCharContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	ch := ctx.GetText()[0]
	var appendCh rune
	switch ch {
	case 'r', 'R':
		appendCh = '\r'
	case 'n', 'N':
		appendCh = '\n'
	case 't', 'T':
		appendCh = '\t'
	case '"':
		appendCh = '"'
	case '*':
		appendCh = '*'
	case '/':
		appendCh = '/'
	case '\\':
		appendCh = '\\'
	case '-':
		appendCh = '\u00ad' // NBSP
	case '_':
		appendCh = '\u00a0' // SHY
	default:
		_this.errorf("BUG: Invalid escape char: [%c]", ch)
	}
	_this.appendCodepoint(appendCh)
}

func (_this *cteListener) ExitValueRid(ctx *parser.ValueRidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeResourceID, uint64(len(_this.arrayData)), _this.arrayData)
}

func (_this *cteListener) ExitCustomTextBegin(ctx *parser.CustomTextBeginContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()
	_this.customType = parseSmallUint(text[1 : len(text)-1])
}

func (_this *cteListener) EnterCustomText(ctx *parser.CustomTextContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitCustomText(ctx *parser.CustomTextContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnCustomText(_this.customType, string(_this.arrayData))
}

func (_this *cteListener) ExitCustomBinaryBegin(ctx *parser.CustomBinaryBeginContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()
	_this.customType = parseSmallUint(text[1 : len(text)-1])
}

func (_this *cteListener) EnterCustomBinary(ctx *parser.CustomBinaryContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitCustomBinary(ctx *parser.CustomBinaryContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnCustomBinary(_this.customType, _this.arrayData)
}

func (_this *cteListener) ExitMediaTextBegin(ctx *parser.MediaTextBeginContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()
	_this.mediaType = text[1 : len(text)-1]
}

func (_this *cteListener) EnterMediaText(ctx *parser.MediaTextContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitMediaText(ctx *parser.MediaTextContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnMedia(_this.mediaType, _this.arrayData)
}

func (_this *cteListener) ExitMediaBinaryBegin(ctx *parser.MediaBinaryBeginContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()
	_this.mediaType = text[1 : len(text)-1]
}

func (_this *cteListener) EnterMediaBinary(ctx *parser.MediaBinaryContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitMediaBinary(ctx *parser.MediaBinaryContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnMedia(_this.mediaType, _this.arrayData)
}

func (_this *cteListener) ExitMarkerID(ctx *parser.MarkerIDContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	contents := ctx.GetText()
	contents = contents[1 : len(contents)-1]
	_this.eventReceiver.OnMarker([]byte(contents))
}

func (_this *cteListener) ExitReference(ctx *parser.ReferenceContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	contents := ctx.GetText()
	contents = contents[1:]
	_this.eventReceiver.OnReferenceLocal([]byte(contents))
}

func (_this *cteListener) ExitValueRemoteRef(ctx *parser.ValueRemoteRefContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeReferenceRemote, uint64(len(_this.arrayData)), _this.arrayData)
}

func (_this *cteListener) EnterContainerMap(ctx *parser.ContainerMapContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnMap()
}

func (_this *cteListener) ExitContainerMap(ctx *parser.ContainerMapContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterContainerList(ctx *parser.ContainerListContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnList()
}

func (_this *cteListener) ExitContainerList(ctx *parser.ContainerListContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterContainerEdge(ctx *parser.ContainerEdgeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEdge()
}

func (_this *cteListener) ExitContainerEdge(ctx *parser.ContainerEdgeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterContainerNode(ctx *parser.ContainerNodeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnNode()
}

func (_this *cteListener) ExitContainerNode(ctx *parser.ContainerNodeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterContainerRecordType(ctx *parser.ContainerRecordTypeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	// TODO: ExitRecordTypeBegin isn't getting called???

	identifier := ctx.GetText()
	cutoff := strings.IndexByte(identifier, '<')
	_this.eventReceiver.OnRecordType([]byte(identifier[1:cutoff]))
}

func (_this *cteListener) ExitContainerRecordType(ctx *parser.ContainerRecordTypeContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterContainerRecord(ctx *parser.ContainerRecordContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	identifier := ctx.GetText()
	cutoff := strings.IndexByte(identifier, '{')
	_this.eventReceiver.OnRecord([]byte(identifier[1:cutoff]))
}

func (_this *cteListener) ExitContainerRecord(ctx *parser.ContainerRecordContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnEndContainer()
}

func (_this *cteListener) EnterArrayI8(ctx *parser.ArrayI8Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitArrayI8(ctx *parser.ArrayI8Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeInt8, uint64(len(_this.arrayData)), _this.arrayData)
}

func (_this *cteListener) EnterArrayI16(ctx *parser.ArrayI16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(16)
}

func (_this *cteListener) ExitArrayI16(ctx *parser.ArrayI16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeInt16, uint64(len(_this.arrayData)/2), _this.arrayData)
}

func (_this *cteListener) EnterArrayI32(ctx *parser.ArrayI32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(32)
}

func (_this *cteListener) ExitArrayI32(ctx *parser.ArrayI32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeInt32, uint64(len(_this.arrayData)/4), _this.arrayData)
}

func (_this *cteListener) EnterArrayI64(ctx *parser.ArrayI64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(64)
}

func (_this *cteListener) ExitArrayI64(ctx *parser.ArrayI64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeInt64, uint64(len(_this.arrayData)/8), _this.arrayData)
}

func (_this *cteListener) EnterArrayU8(ctx *parser.ArrayU8Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(8)
}

func (_this *cteListener) ExitArrayU8(ctx *parser.ArrayU8Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeUint8, uint64(len(_this.arrayData)), _this.arrayData)
}

func (_this *cteListener) EnterArrayU16(ctx *parser.ArrayU16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(16)
}

func (_this *cteListener) ExitArrayU16(ctx *parser.ArrayU16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeUint16, uint64(len(_this.arrayData)/2), _this.arrayData)
}

func (_this *cteListener) EnterArrayU32(ctx *parser.ArrayU32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(32)
}

func (_this *cteListener) ExitArrayU32(ctx *parser.ArrayU32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeUint32, uint64(len(_this.arrayData)/4), _this.arrayData)
}

func (_this *cteListener) EnterArrayU64(ctx *parser.ArrayU64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(64)
}

func (_this *cteListener) ExitArrayU64(ctx *parser.ArrayU64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeUint64, uint64(len(_this.arrayData)/8), _this.arrayData)
}

func (_this *cteListener) EnterArrayF16(ctx *parser.ArrayF16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(16)
}

func (_this *cteListener) ExitArrayF16(ctx *parser.ArrayF16Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeFloat16, uint64(len(_this.arrayData)/2), _this.arrayData)
}

func (_this *cteListener) EnterArrayF32(ctx *parser.ArrayF32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(32)
}

func (_this *cteListener) ExitArrayF32(ctx *parser.ArrayF32Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeFloat32, uint64(len(_this.arrayData)/4), _this.arrayData)
}

func (_this *cteListener) EnterArrayF64(ctx *parser.ArrayF64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(64)
}

func (_this *cteListener) ExitArrayF64(ctx *parser.ArrayF64Context) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeFloat64, uint64(len(_this.arrayData)/8), _this.arrayData)
}

func (_this *cteListener) EnterArrayUid(ctx *parser.ArrayUidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(128)
}

func (_this *cteListener) ExitArrayUid(ctx *parser.ArrayUidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.eventReceiver.OnArray(events.ArrayTypeUID, uint64(len(_this.arrayData)/16), _this.arrayData)
}

func (_this *cteListener) EnterArrayBit(ctx *parser.ArrayBitContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.beginArray(1)
}

func (_this *cteListener) ExitArrayBit(ctx *parser.ArrayBitContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	bitCount := len(_this.arrayData)
	if bitCount == 0 {
		_this.eventReceiver.OnArray(events.ArrayTypeBit, 0, []byte{})
		return
	}
	bits := make([]byte, 0, bitCount/8+1)
	current := _this.arrayData
	for {
		nextByte := uint8(0)
		for i := 0; i < 8; i++ {
			if len(current) == 0 {
				break
			}
			if current[0] == '1' {
				nextByte |= 1 << i
			}
			current = current[1:]
		}
		bits = append(bits, nextByte)
		if len(current) == 0 {
			break
		}
	}
	_this.eventReceiver.OnArray(events.ArrayTypeBit, uint64(bitCount), bits)
}

func (_this *cteListener) ExitArrayElemBits(ctx *parser.ArrayElemBitsContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = append(_this.arrayData, []uint8(ctx.GetText())...)
}

func (_this *cteListener) ExitArrayElemInt(ctx *parser.ArrayElemIntContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseIntElement(ctx.GetText(), 0, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemIntB(ctx *parser.ArrayElemIntBContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseIntElement(ctx.GetText(), 2, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemIntO(ctx *parser.ArrayElemIntOContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseIntElement(ctx.GetText(), 8, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemIntX(ctx *parser.ArrayElemIntXContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseIntElement(ctx.GetText(), 16, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemUint(ctx *parser.ArrayElemUintContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseUintElement(ctx.GetText(), 0, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemUintB(ctx *parser.ArrayElemUintBContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseUintElement(ctx.GetText(), 2, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemUintO(ctx *parser.ArrayElemUintOContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseUintElement(ctx.GetText(), 8, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemUintX(ctx *parser.ArrayElemUintXContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseUintElement(ctx.GetText(), 16, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemFloat(ctx *parser.ArrayElemFloatContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseFloatElement(ctx.GetText(), 0, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemFloatX(ctx *parser.ArrayElemFloatXContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseFloatElement(ctx.GetText(), 16, _this.elementSizeBits, _this.arrayData)
}

func (_this *cteListener) ExitArrayElemNan(ctx *parser.ArrayElemNanContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	switch _this.elementSizeBits {
	case 16:
		_this.addUint16Element(common.Bfloat16QuietNanBits)
	case 32:
		_this.addUint32Element(common.Float32QuietNanBits)
	case 64:
		_this.addUint64Element(common.Float64QuietNanBits)
	default:
		panic(fmt.Errorf("BUG: tried to generate NAN for element size %v", _this.elementSizeBits))
	}
}

func (_this *cteListener) ExitArrayElemSnan(ctx *parser.ArrayElemSnanContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	switch _this.elementSizeBits {
	case 16:
		_this.addUint16Element(common.Bfloat16SignalingNanBits)
	case 32:
		_this.addUint32Element(common.Float32SignalingNanBits)
	case 64:
		_this.addUint64Element(common.Float64SignalingNanBits)
	default:
		panic(fmt.Errorf("BUG: tried to generate signaling NAN for element size %v", _this.elementSizeBits))
	}
}

func (_this *cteListener) ExitArrayElemInf(ctx *parser.ArrayElemInfContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	switch _this.elementSizeBits {
	case 16:
		_this.addUint16Element(common.Bfloat16InfBits)
	case 32:
		_this.addUint32Element(common.Float32InfBits)
	case 64:
		_this.addUint64Element(common.Float64InfBits)
	default:
		panic(fmt.Errorf("BUG: tried to generate INF for element size %v", _this.elementSizeBits))
	}
}

func (_this *cteListener) ExitArrayElemNinf(ctx *parser.ArrayElemNinfContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	switch _this.elementSizeBits {
	case 16:
		_this.addUint16Element(common.Bfloat16NegInfBits)
	case 32:
		_this.addUint32Element(common.Float32NegInfBits)
	case 64:
		_this.addUint64Element(common.Float64NegInfBits)
	default:
		panic(fmt.Errorf("BUG: tried to generate negative INF for element size %v", _this.elementSizeBits))
	}
}

func (_this *cteListener) ExitArrayElemUid(ctx *parser.ArrayElemUidContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = appendUID(ctx.GetText(), _this.arrayData)
}

func (_this *cteListener) ExitArrayElemByteX(ctx *parser.ArrayElemByteXContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	_this.arrayData = parseUintElement(ctx.GetText(), 16, 8, _this.arrayData)
}

func (_this *cteListener) ExitCommentLine(ctx *parser.CommentLineContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()[2:]
	text = text[:len(text)-1]
	if len(text) > 0 && text[len(text)-1] == '\r' {
		text = text[:len(text)-1]
	}
	_this.eventReceiver.OnComment(false, []byte(text))
}

func (_this *cteListener) ExitCommentBlock(ctx *parser.CommentBlockContext) {
	defer func() {
		_this.wrapPanic(recover(), ctx.BaseParserRuleContext)
	}()

	text := ctx.GetText()[2:]
	text = text[:len(text)-2]
	_this.eventReceiver.OnComment(true, []byte(text))
}

// ---------------------------------------------------------------------------

func parseSmallUint(str string) uint64 {
	if v, err := strconv.ParseUint(str, 0, 64); err == nil {
		return v
	} else {
		panic(err)
	}
}

func parseHexCodepoint(str string) rune {
	if v, err := strconv.ParseUint(str, 16, 32); err == nil {
		return rune(v)
	} else {
		panic(err)
	}
}

func appendUID(str string, dst []byte) []byte {
	endpoints := []int{-1, 8, 13, 18, 23, 36}
	for iEndpoint := 0; iEndpoint < len(endpoints)-1; iEndpoint++ {
		p1 := endpoints[iEndpoint] + 1
		p2 := endpoints[iEndpoint+1]
		newBytes, err := hex.DecodeString(str[p1:p2])
		if err != nil {
			panic(fmt.Errorf("error parsing hex bytes: %v", err))
		}
		dst = append(dst, newBytes...)
	}
	return dst
}

func parseUintElement(str string, base int, bitSize int, result []byte) []byte {
	element, err := strconv.ParseUint(str, base, bitSize)
	if err != nil {
		panic(fmt.Errorf("error parsing uint element: %v", err))
	}
	switch bitSize {
	case 8:
		return append(result, byte(element))
	case 16:
		return binary.LittleEndian.AppendUint16(result, uint16(element))
	case 32:
		return binary.LittleEndian.AppendUint32(result, uint32(element))
	case 64:
		return binary.LittleEndian.AppendUint64(result, element)
	default:
		panic(fmt.Errorf("BUG: passed invalid bit size of %v to parseUintElement", bitSize))
	}
}

func parseIntElement(str string, base int, bitSize int, result []byte) []byte {
	element, err := strconv.ParseInt(str, base, bitSize)
	if err != nil {
		panic(fmt.Errorf("error parsing int element: %v", err))
	}
	switch bitSize {
	case 8:
		return append(result, byte(element))
	case 16:
		return binary.LittleEndian.AppendUint16(result, uint16(element))
	case 32:
		return binary.LittleEndian.AppendUint32(result, uint32(element))
	case 64:
		return binary.LittleEndian.AppendUint64(result, uint64(element))
	default:
		panic(fmt.Errorf("BUG: passed invalid bit size of %v to parseIntElement", bitSize))
	}
}

var floatZero = big.Float{}

func isFloatZero(str string) bool {
	f, _, err := big.ParseFloat(str, 0, uint(len(str)), big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	return f.Cmp(&floatZero) == 0
}

func normalizeFloatString(str string, base int) string {
	str = strings.ReplaceAll(str, "_", "")
	strNoSign := str
	isNegative := false
	if str[0] == '-' {
		isNegative = true
		strNoSign = str[1:]
	}
	if base == 16 {
		if isNegative {
			str = fmt.Sprintf("-0x%v", str[1:])
		} else {
			str = fmt.Sprintf("0x%v", str)
		}
	} else if base == 0 {
		if len(strNoSign) > 2 && strNoSign[0] == '0' && (strNoSign[1] == 'x' || strNoSign[1] == 'X') {
			base = 16
		}
	}

	if base == 16 && !strings.ContainsRune(str, 'p') && !strings.ContainsRune(str, 'P') {
		// Add a zero exponent to keep strconv.ParseFloat() happy
		str += "p0"
	}
	return str
}

func parseFloatElement(str string, base int, bitSize int, result []byte) []byte {
	str = normalizeFloatString(str, base)
	strNoSign := str
	isNegative := false
	if str[0] == '-' {
		// Split out negation because go doesn't handle -0 properly
		isNegative = true
		strNoSign = str[1:]
	}
	parseBitSize := bitSize
	if parseBitSize == 16 {
		parseBitSize = 32
	}
	element, err := strconv.ParseFloat(strNoSign, parseBitSize)
	if err != nil {
		panic(fmt.Errorf("error parsing float element: %v", err))
	}
	if isNegative {
		element = -element
	}
	switch bitSize {
	case 16:
		elem32 := float32(element)
		// TODO: Check for float16 overflow
		if math.IsInf(float64(elem32), 0) {
			panic(fmt.Errorf("float element %v is too big for float16", str))
		}
		if elem32 == 0 && !isFloatZero(str) {
			panic(fmt.Errorf("float element %v is too small for float16", str))
		}
		return binary.LittleEndian.AppendUint16(result, uint16(math.Float32bits(elem32)>>16))
	case 32:
		elem32 := float32(element)
		if math.IsInf(float64(elem32), 0) {
			panic(fmt.Errorf("float element %v is too big for float32", str))
		}
		if elem32 == 0 && !isFloatZero(str) {
			panic(fmt.Errorf("float element %v is too small for float32", str))
		}
		return binary.LittleEndian.AppendUint32(result, math.Float32bits(elem32))
	case 64:
		if math.IsInf(element, 0) {
			panic(fmt.Errorf("float element %v is too big for float64", str))
		}
		if element == 0 && !isFloatZero(str) {
			panic(fmt.Errorf("float element %v is too small for float64", str))
		}
		return binary.LittleEndian.AppendUint64(result, math.Float64bits(element))
	default:
		panic(fmt.Errorf("BUG: passed invalid bit size of %v to parseFloatElement", bitSize))
	}
}

func (_this *cteListener) addUint16Element(element uint16) {
	_this.arrayData = binary.LittleEndian.AppendUint16(_this.arrayData, element)
}

func (_this *cteListener) addUint32Element(element uint32) {
	_this.arrayData = binary.LittleEndian.AppendUint32(_this.arrayData, element)
}

func (_this *cteListener) addUint64Element(element uint64) {
	_this.arrayData = binary.LittleEndian.AppendUint64(_this.arrayData, element)
}

func parseDecimalInt(str string) int64 {
	if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		return v
	} else {
		panic(err)
	}
}

func parseCoord(str string) float64 {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	} else {
		panic(err)
	}
}

func parseDate(str string) compact_time.Time {
	r, err := regexp.Compile(`(-?\d+)-(\d\d?)-(\d\d?)$`)
	if err != nil {
		panic(err)
	}
	strs := r.FindAllStringSubmatch(str, -1)[0]
	year := int(parseDecimalInt(strs[1]))
	month := int(parseDecimalInt(strs[2]))
	day := int(parseDecimalInt(strs[3]))
	v := compact_time.NewDate(year, month, day)
	if err := v.Validate(); err != nil {
		panic(err)
	}
	return v
}

func parseTimeNanoseconds(str string) int {
	if len(str) == 0 {
		return 0
	}
	subsecStr := str[1:]
	digitCount := len(subsecStr)
	subseconds := int(parseDecimalInt(subsecStr))
	return subseconds * int(math.Pow10(9-digitCount))
}

func parseTimezone(str string) compact_time.Timezone {
	if len(str) == 0 {
		return compact_time.TZAtUTC()
	}
	switch str[0] {
	case '/':
		if (str[1] >= '0' && str[1] <= '9') || str[1] == '-' {
			r, err := regexp.Compile(`(-?\d+(\.\d+)?)/(-?\d+(\.\d+)?)$`)
			if err != nil {
				panic(err)
			}
			strs := r.FindAllStringSubmatch(str[1:], -1)[0]
			latitude := parseCoord(strs[1])
			longitude := parseCoord(strs[3])
			return compact_time.TZAtLatLong(int(latitude*100), int(longitude*100))
		} else {
			return compact_time.TZAtAreaLocation(str[1:])
		}
	case '+':
		hours := parseDecimalInt(str[1:3])
		minutes := parseDecimalInt(str[3:])
		return compact_time.TZWithMiutesOffsetFromUTC(int(hours*60 + minutes))
	case '-':
		hours := parseDecimalInt(str[1:3])
		minutes := parseDecimalInt(str[3:])
		return compact_time.TZWithMiutesOffsetFromUTC(-int(hours*60 + minutes))
	default:
		panic(fmt.Errorf("BUG: Unknown time separator '%c'", str[0]))
	}
}

func parseTime(str string) compact_time.Time {
	r, err := regexp.Compile(`(\d+):(\d\d):(\d\d)(\.\d+)?([+-/].+)?$`)
	if err != nil {
		panic(err)
	}
	strs := r.FindAllStringSubmatch(str, -1)[0]
	hour := int(parseDecimalInt(strs[1]))
	minute := int(parseDecimalInt(strs[2]))
	second := int(parseDecimalInt(strs[3]))
	var nanosecond int
	timezone := compact_time.TZAtUTC()
	if len(strs[4]) > 0 {
		nanosecond = parseTimeNanoseconds(strs[4])
	}
	if len(strs[5]) > 0 {
		timezone = parseTimezone(strs[5])
	}

	v := compact_time.NewTime(hour, minute, second, nanosecond, timezone)
	if err := v.Validate(); err != nil {
		panic(err)
	}
	return v
}

func parseDateTime(str string) compact_time.Time {
	r, err := regexp.Compile(`(-?\d+)-(\d\d)-(\d\d)/(\d+):(\d\d):(\d\d)(\.\d+)?([+-/].+)?$`)
	if err != nil {
		panic(err)
	}
	strs := r.FindAllStringSubmatch(str, -1)[0]
	year := int(parseDecimalInt(strs[1]))
	month := int(parseDecimalInt(strs[2]))
	day := int(parseDecimalInt(strs[3]))
	hour := int(parseDecimalInt(strs[4]))
	minute := int(parseDecimalInt(strs[5]))
	second := int(parseDecimalInt(strs[6]))
	nanosecond := parseTimeNanoseconds(strs[7])
	timezone := parseTimezone(strs[8])

	v := compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, timezone)
	if err := v.Validate(); err != nil {
		panic(err)
	}
	return v
}

func parseDateOrDateTime(str string) compact_time.Time {
	if strings.ContainsRune(str, '/') {
		return parseDateTime(str)
	} else {
		return parseDate(str)
	}
}
