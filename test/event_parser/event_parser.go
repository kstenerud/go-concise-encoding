package event_parser

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/internal/common"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/event_parser/parser"
)

func ParseEvents(eventStrings ...string) test.Events {
	var index = 0
	var eventStr string

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("event index %v \"%v\": %w", index, eventStr, v))
			default:
				panic(v)
			}
		}
	}()

	events := test.Events{}
	for index, eventStr = range eventStrings {
		evts, err := parseEventString(eventStr)
		if err != nil {
			panic(err)
		}
		events = append(events, evts...)
	}

	return events
}

func parseEventString(eventStr string) (events test.Events, err error) {
	errorListener := new(reportingErrorListener)

	is := antlr.NewInputStream(eventStr)
	lexer := parser.NewCEEventLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewCEEventParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)
	p.SetErrorHandler(new(bailErrorStrategy))

	listener := &eventListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, p.Start())

	return listener.Events, errorListener.Error
}

type reportingErrorListener struct {
	*antlr.DefaultErrorListener
	Error error
}

func (_this *reportingErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	_this.Error = fmt.Errorf("col %v: %v", column, msg)
}

type bailErrorStrategy struct {
	antlr.DefaultErrorStrategy
}

func (b *bailErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {}
func (b *bailErrorStrategy) Sync(recognizer antlr.Parser)                                  {}

func iterateChildren(children []antlr.Tree, iterate func(string)) {
	for j := 1; j < len(children); j++ {
		child := children[j]
		if c2, ok := child.(antlr.TerminalNode); ok {
			iterate(c2.GetText())
		}
	}
}

func getTokenInfo(node antlr.Tree) (tokenType int, token string) {
	if t, ok := node.(antlr.TerminalNode); ok {
		return t.GetSymbol().GetTokenType(), t.GetText()
	}
	panic(fmt.Errorf("BUG: Expected terminal node but got %v", reflect.TypeOf(node)))
}

func getTokenText(node antlr.Tree) string {
	if t, ok := node.(antlr.TerminalNode); ok {
		return t.GetText()
	}
	panic(fmt.Errorf("BUG: Expected terminal node but got %v", reflect.TypeOf(node)))
}

func parseInt(str string) (smallInt int64, bigInt *big.Int) {
	var err error
	if smallInt, err = strconv.ParseInt(str, 0, 64); err == nil {
		return
	}

	bigInt = &big.Int{}
	if _, success := bigInt.SetString(str, 0); success {
		return
	}

	panic(fmt.Errorf("BUG: Expected an integer but got \"%v\"", str))
}

func parseSmallInt(str string) int64 {
	si, bi := parseInt(str)
	if bi != nil {
		panic(fmt.Errorf("unexpectedly large integer \"%v\"", str))
	}
	return si
}

func parseDecimalInt(str string) int64 {
	if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		return v
	} else {
		panic(err)
	}
}

func parseSmallUint(str string) uint64 {
	if v, err := strconv.ParseUint(str, 0, 64); err == nil {
		return v
	} else {
		panic(err)
	}
}

func parseSmallUintX(str string) uint64 {
	if v, err := strconv.ParseUint(str, 16, 64); err == nil {
		return v
	} else {
		panic(err)
	}
}

func parseBinaryFloat(str string) (smallFloat float64, bigFloat *big.Float) {
	switch str {
	case "nan":
		return math.Float64frombits(common.Float64QuietNanBits), nil
	case "snan":
		return math.Float64frombits(common.Float64SignalingNanBits), nil
	}

	var accuracy big.Accuracy
	if strings.Contains(str, "0x") {
		digitCount := len(strings.Split(str, "x")[1])
		bigFloat = &big.Float{}
		bigFloat.SetPrec(uint(digitCount) * 4)
		if _, success := bigFloat.SetString(str); success {
			smallFloat, accuracy = bigFloat.Float64()
			if accuracy == big.Exact {
				bigFloat = nil
			}
			return
		} else {
			panic(fmt.Errorf("BUG: Could not convert [%v] to float", str))
		}
	}

	digitCount := len(str)
	bigFloat = &big.Float{}
	bigFloat.SetPrec(uint(common.DecimalDigitsToBits(digitCount)))
	if _, success := bigFloat.SetString(str); success {
		smallFloat, accuracy = bigFloat.Float64()
		if accuracy == big.Exact {
			// big.Float to float64 introduces false precision, so parse with a better tested function
			smallFloat, _ = strconv.ParseFloat(str, 64)
			bigFloat = nil
		}
	} else {
		panic(fmt.Errorf("BUG: Could not convert [%v] to float", str))
	}

	return
}

func parseFloat64(str string) float64 {
	sf, bf := parseBinaryFloat(str)
	if bf != nil {
		panic(fmt.Errorf("unexpectedly large binary float \"%v\"", str))
	}
	return float64(sf)
}

func parseFloat32(str string) float32 {
	switch str {
	case "nan":
		return math.Float32frombits(common.Float32QuietNanBits)
	case "snan":
		return math.Float32frombits(common.Float32SignalingNanBits)
	default:
		return float32(parseFloat64(str))
	}
}

func parseDecimalFloat(str string) (smallFloat compact_float.DFloat, bigFloat *apd.Decimal) {
	var err error
	if smallFloat, err = compact_float.DFloatFromString(str); err == nil {
		return
	}

	if bigFloat, _, err = apd.NewFromString(str); err == nil {
		return
	}

	panic(fmt.Errorf("BUG: Expected a decimal float but got \"%v\"", str))
}

func parseHex(str string) (result uint64) {
	for _, b := range str {
		switch b {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			result = (result << 4) | uint64(b-'0')
		case 'a', 'b', 'c', 'd', 'e', 'f':
			result = (result << 4) | uint64(b-'a'+10)
		case 'A', 'B', 'C', 'D', 'E', 'F':
			result = (result << 4) | uint64(b-'A'+10)
		default:
			panic(fmt.Errorf("BUG: error parsing hexadecimal: Invalid char [%c] in [%v]", b, str))
		}
	}
	return
}

func parseCoord(str string) float64 {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	} else {
		panic(err)
	}
}

func parseDate(str string) compact_time.Time {
	r, err := regexp.Compile(`(-?\d+)-(\d\d)-(\d\d)$`)
	if err != nil {
		panic(err)
	}
	strs := r.FindAllStringSubmatch(str, -1)[0]
	year := int(parseDecimalInt(strs[1]))
	month := int(parseDecimalInt(strs[2]))
	day := int(parseDecimalInt(strs[3]))
	return compact_time.NewDate(year, month, day)
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

	return compact_time.NewTime(hour, minute, second, nanosecond, timezone)
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

	return compact_time.NewTimestamp(year, month, day, hour, minute, second, nanosecond, timezone)
}

func parseUUID(str string) []byte {
	buff := bytes.Buffer{}
	endpoints := []int{8, 13, 18, 23, 36}
	iData := 0
	for iEndpoint := 0; iEndpoint < len(endpoints); iEndpoint++ {
		endpoint := endpoints[iEndpoint]
		for ; iData < endpoint; iData += 2 {
			buff.WriteByte(byte(parseHex(str[iData : iData+2])))
		}
		iData++
	}
	return buff.Bytes()
}

func getStringArg(tree []antlr.Tree) string {
	if len(tree) > 1 {
		return getTokenText(tree[1])
	}
	return ""
}

func parseArrayElementsBit(children []antlr.Tree) []bool {
	elements := make([]bool, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		switch text {
		case "1":
			elements = append(elements, true)
		case "0":
			elements = append(elements, false)
		default:
			panic(fmt.Errorf("BUG: %v: unexpected boolean element value", text))
		}
	})
	return elements
}

func parseArrayElementsUUID(children []antlr.Tree) [][]byte {
	elements := make([][]byte, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, parseUUID(text))
	})
	return elements
}

func parseArrayElementsInt16(children []antlr.Tree) []int16 {
	elements := make([]int16, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, int16(parseSmallInt(text)))
	})
	return elements
}

func parseArrayElementsInt32(children []antlr.Tree) []int32 {
	elements := make([]int32, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, int32(parseSmallInt(text)))
	})
	return elements
}

func parseArrayElementsInt64(children []antlr.Tree) []int64 {
	elements := make([]int64, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, int64(parseSmallInt(text)))
	})
	return elements
}

func parseArrayElementsInt8(children []antlr.Tree) []int8 {
	elements := make([]int8, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, int8(parseSmallInt(text)))
	})
	return elements
}

func parseArrayElementsUint8(children []antlr.Tree) []uint8 {
	elements := make([]uint8, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint8(parseSmallUint(text)))
	})
	return elements
}

func parseArrayElementsUint8X(children []antlr.Tree) []uint8 {
	elements := make([]uint8, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint8(parseSmallUintX(text)))
	})
	return elements
}

func parseArrayElementsUint16(children []antlr.Tree) []uint16 {
	elements := make([]uint16, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint16(parseSmallUint(text)))
	})
	return elements
}

func parseArrayElementsUint16X(children []antlr.Tree) []uint16 {
	elements := make([]uint16, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint16(parseSmallUintX(text)))
	})
	return elements
}

func parseArrayElementsUint32(children []antlr.Tree) []uint32 {
	elements := make([]uint32, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint32(parseSmallUint(text)))
	})
	return elements
}

func parseArrayElementsUint32X(children []antlr.Tree) []uint32 {
	elements := make([]uint32, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint32(parseSmallUintX(text)))
	})
	return elements
}

func parseArrayElementsUint64(children []antlr.Tree) []uint64 {
	elements := make([]uint64, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint64(parseSmallUint(text)))
	})
	return elements
}

func parseArrayElementsUint64X(children []antlr.Tree) []uint64 {
	elements := make([]uint64, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, uint64(parseSmallUintX(text)))
	})
	return elements
}

func parseArrayElementsFloat16(children []antlr.Tree) []float32 {
	elements := make([]float32, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, parseFloat32(text))
	})
	return elements
}

func parseArrayElementsFloat32(children []antlr.Tree) []float32 {
	elements := make([]float32, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, parseFloat32(text))
	})
	return elements
}

func parseArrayElementsFloat64(children []antlr.Tree) []float64 {
	elements := make([]float64, 0, len(children)-1)
	iterateChildren(children, func(text string) {
		elements = append(elements, float64(parseFloat64(text)))
	})
	return elements
}

type eventListener struct {
	*parser.BaseCEEventParserListener
	Events test.Events
}

func (_this *eventListener) setEvents(events ...test.Event) {
	_this.Events = test.Events(events)
}

func (_this *eventListener) ExitEventArrayBits(ctx *parser.EventArrayBitsContext) {
	_this.setEvents(test.AB(parseArrayElementsBit(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayChunkLast(ctx *parser.EventArrayChunkLastContext) {
	_this.setEvents(test.ACL(parseSmallUint(getTokenText(ctx.GetChild(1)))))
}
func (_this *eventListener) ExitEventArrayChunkMore(ctx *parser.EventArrayChunkMoreContext) {
	_this.setEvents(test.ACM(parseSmallUint(getTokenText(ctx.GetChild(1)))))
}
func (_this *eventListener) ExitEventArrayDataBits(ctx *parser.EventArrayDataBitsContext) {
	_this.setEvents(test.ADB(parseArrayElementsBit(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataFloat16(ctx *parser.EventArrayDataFloat16Context) {
	_this.setEvents(test.ADF16(parseArrayElementsFloat16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataFloat32(ctx *parser.EventArrayDataFloat32Context) {
	_this.setEvents(test.ADF32(parseArrayElementsFloat32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataFloat64(ctx *parser.EventArrayDataFloat64Context) {
	_this.setEvents(test.ADF64(parseArrayElementsFloat64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataInt16(ctx *parser.EventArrayDataInt16Context) {
	_this.setEvents(test.ADI16(parseArrayElementsInt16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataInt32(ctx *parser.EventArrayDataInt32Context) {
	_this.setEvents(test.ADI32(parseArrayElementsInt32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataInt64(ctx *parser.EventArrayDataInt64Context) {
	_this.setEvents(test.ADI64(parseArrayElementsInt64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataInt8(ctx *parser.EventArrayDataInt8Context) {
	_this.setEvents(test.ADI8(parseArrayElementsInt8(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataText(ctx *parser.EventArrayDataTextContext) {
	_this.setEvents(test.ADT(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUID(ctx *parser.EventArrayDataUIDContext) {
	_this.setEvents(test.ADU(parseArrayElementsUUID(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint16(ctx *parser.EventArrayDataUint16Context) {
	_this.setEvents(test.ADU16(parseArrayElementsUint16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint16X(ctx *parser.EventArrayDataUint16XContext) {
	_this.setEvents(test.ADU16(parseArrayElementsUint16X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint32(ctx *parser.EventArrayDataUint32Context) {
	_this.setEvents(test.ADU32(parseArrayElementsUint32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint32X(ctx *parser.EventArrayDataUint32XContext) {
	_this.setEvents(test.ADU32(parseArrayElementsUint32X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint64(ctx *parser.EventArrayDataUint64Context) {
	_this.setEvents(test.ADU64(parseArrayElementsUint64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint64X(ctx *parser.EventArrayDataUint64XContext) {
	_this.setEvents(test.ADU64(parseArrayElementsUint64X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint8(ctx *parser.EventArrayDataUint8Context) {
	_this.setEvents(test.ADU8(parseArrayElementsUint8(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayDataUint8X(ctx *parser.EventArrayDataUint8XContext) {
	_this.setEvents(test.ADU8(parseArrayElementsUint8X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayFloat16(ctx *parser.EventArrayFloat16Context) {
	_this.setEvents(test.AF16(parseArrayElementsFloat16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayFloat32(ctx *parser.EventArrayFloat32Context) {
	_this.setEvents(test.AF32(parseArrayElementsFloat32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayFloat64(ctx *parser.EventArrayFloat64Context) {
	_this.setEvents(test.AF64(parseArrayElementsFloat64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayInt16(ctx *parser.EventArrayInt16Context) {
	_this.setEvents(test.AI16(parseArrayElementsInt16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayInt32(ctx *parser.EventArrayInt32Context) {
	_this.setEvents(test.AI32(parseArrayElementsInt32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayInt64(ctx *parser.EventArrayInt64Context) {
	_this.setEvents(test.AI64(parseArrayElementsInt64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayInt8(ctx *parser.EventArrayInt8Context) {
	_this.setEvents(test.AI8(parseArrayElementsInt8(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUID(ctx *parser.EventArrayUIDContext) {
	_this.setEvents(test.AU(parseArrayElementsUUID(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint16(ctx *parser.EventArrayUint16Context) {
	_this.setEvents(test.AU16(parseArrayElementsUint16(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint16X(ctx *parser.EventArrayUint16XContext) {
	_this.setEvents(test.AU16(parseArrayElementsUint16X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint32(ctx *parser.EventArrayUint32Context) {
	_this.setEvents(test.AU32(parseArrayElementsUint32(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint32X(ctx *parser.EventArrayUint32XContext) {
	_this.setEvents(test.AU32(parseArrayElementsUint32X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint64(ctx *parser.EventArrayUint64Context) {
	_this.setEvents(test.AU64(parseArrayElementsUint64(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint64X(ctx *parser.EventArrayUint64XContext) {
	_this.setEvents(test.AU64(parseArrayElementsUint64X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint8(ctx *parser.EventArrayUint8Context) {
	_this.setEvents(test.AU8(parseArrayElementsUint8(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventArrayUint8X(ctx *parser.EventArrayUint8XContext) {
	_this.setEvents(test.AU8(parseArrayElementsUint8X(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventBeginArrayBits(ctx *parser.EventBeginArrayBitsContext) {
	_this.setEvents(test.BAB())
}
func (_this *eventListener) ExitEventBeginArrayFloat16(ctx *parser.EventBeginArrayFloat16Context) {
	_this.setEvents(test.BAF16())
}
func (_this *eventListener) ExitEventBeginArrayFloat32(ctx *parser.EventBeginArrayFloat32Context) {
	_this.setEvents(test.BAF32())
}
func (_this *eventListener) ExitEventBeginArrayFloat64(ctx *parser.EventBeginArrayFloat64Context) {
	_this.setEvents(test.BAF64())
}
func (_this *eventListener) ExitEventBeginArrayInt16(ctx *parser.EventBeginArrayInt16Context) {
	_this.setEvents(test.BAI16())
}
func (_this *eventListener) ExitEventBeginArrayInt32(ctx *parser.EventBeginArrayInt32Context) {
	_this.setEvents(test.BAI32())
}
func (_this *eventListener) ExitEventBeginArrayInt64(ctx *parser.EventBeginArrayInt64Context) {
	_this.setEvents(test.BAI64())
}
func (_this *eventListener) ExitEventBeginArrayInt8(ctx *parser.EventBeginArrayInt8Context) {
	_this.setEvents(test.BAI8())
}
func (_this *eventListener) ExitEventBeginArrayUID(ctx *parser.EventBeginArrayUIDContext) {
	_this.setEvents(test.BAU())
}
func (_this *eventListener) ExitEventBeginArrayUint16(ctx *parser.EventBeginArrayUint16Context) {
	_this.setEvents(test.BAU16())
}
func (_this *eventListener) ExitEventBeginArrayUint32(ctx *parser.EventBeginArrayUint32Context) {
	_this.setEvents(test.BAU32())
}
func (_this *eventListener) ExitEventBeginArrayUint64(ctx *parser.EventBeginArrayUint64Context) {
	_this.setEvents(test.BAU64())
}
func (_this *eventListener) ExitEventBeginArrayUint8(ctx *parser.EventBeginArrayUint8Context) {
	_this.setEvents(test.BAU8())
}
func (_this *eventListener) ExitEventBeginCustomBinary(ctx *parser.EventBeginCustomBinaryContext) {
	customType := parseSmallUint(getTokenText(ctx.GetChild(1)))
	_this.setEvents(test.BCB(customType))
}
func (_this *eventListener) ExitEventBeginCustomText(ctx *parser.EventBeginCustomTextContext) {
	customType := parseSmallUint(getTokenText(ctx.GetChild(1)))
	_this.setEvents(test.BCT(customType))
}
func (_this *eventListener) ExitEventBeginMedia(ctx *parser.EventBeginMediaContext) {
	_this.setEvents(test.BMEDIA(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventBeginResourceId(ctx *parser.EventBeginResourceIdContext) {
	_this.setEvents(test.BRID())
}
func (_this *eventListener) ExitEventBeginString(ctx *parser.EventBeginStringContext) {
	_this.setEvents(test.BS())
}
func (_this *eventListener) ExitEventBeginRemoteReference(ctx *parser.EventBeginRemoteReferenceContext) {
	_this.setEvents(test.BREFR())
}
func (_this *eventListener) ExitEventBoolean(ctx *parser.EventBooleanContext) {
	ttype, _ := getTokenInfo(ctx.GetChild(1))
	_this.setEvents(test.B(ttype == parser.CEEventLexerTRUE))
}
func (_this *eventListener) ExitEventCommentMultiline(ctx *parser.EventCommentMultilineContext) {
	_this.setEvents(test.CM(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventCommentSingleLine(ctx *parser.EventCommentSingleLineContext) {
	_this.setEvents(test.CS(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventCustomBinary(ctx *parser.EventCustomBinaryContext) {
	customType := parseSmallUint(getTokenText(ctx.GetChild(1)))
	children := ctx.GetChildren()[1:]
	_this.setEvents(test.CB(customType, parseArrayElementsUint8X(children)))
}
func (_this *eventListener) ExitEventCustomText(ctx *parser.EventCustomTextContext) {
	customType := parseSmallUint(getTokenText(ctx.GetChild(1)))
	children := ctx.GetChildren()[1:]
	_this.setEvents(test.CT(customType, getStringArg(children)))
}
func (_this *eventListener) ExitEventEdge(ctx *parser.EventEdgeContext) {
	_this.setEvents(test.EDGE())
}
func (_this *eventListener) ExitEventEndContainer(ctx *parser.EventEndContainerContext) {
	_this.setEvents(test.E())
}
func (_this *eventListener) ExitEventList(ctx *parser.EventListContext) {
	_this.setEvents(test.L())
}
func (_this *eventListener) ExitEventMap(ctx *parser.EventMapContext) {
	_this.setEvents(test.M())
}
func (_this *eventListener) ExitEventMarker(ctx *parser.EventMarkerContext) {
	_this.setEvents(test.MARK(getTokenText(ctx.GetChild(1))))
}
func (_this *eventListener) ExitEventMedia(ctx *parser.EventMediaContext) {
	mediaType := getTokenText(ctx.GetChild(1))
	children := ctx.GetChildren()
	elements := make([]uint8, 0, len(children)-2)
	for j := 2; j < len(children); j++ {
		child := children[j]
		if c2, ok := child.(antlr.TerminalNode); ok {
			elements = append(elements, uint8(parseSmallUintX(c2.GetText())))
		}
	}
	_this.setEvents(
		test.MEDIA(mediaType, elements),
	)
}
func (_this *eventListener) ExitEventNode(ctx *parser.EventNodeContext) {
	_this.setEvents(test.NODE())
}
func (_this *eventListener) ExitEventNull(ctx *parser.EventNullContext) {
	_this.setEvents(test.NULL())
}
func (_this *eventListener) ExitEventNumber(ctx *parser.EventNumberContext) {
	ttype, ttext := getTokenInfo(ctx.GetChild(1))
	switch ttype {
	case parser.CEEventParserFLOAT_NAN:
		_this.setEvents(test.N(compact_float.QuietNaN()))
	case parser.CEEventParserFLOAT_SNAN:
		_this.setEvents(test.N(compact_float.SignalingNaN()))
	case parser.CEEventParserFLOAT_INF:
		sign := 1
		if ttext[0] == '-' {
			sign = -1
		}
		_this.setEvents(test.N(math.Inf(sign)))
	case parser.CEEventParserFLOAT_DEC:
		sf, bf := parseDecimalFloat(ttext)
		if bf != nil {
			_this.setEvents(test.N(bf))
		} else {
			_this.setEvents(test.N(sf))
		}
	case parser.CEEventParserFLOAT_HEX:
		sf, bf := parseBinaryFloat(ttext)
		if bf != nil {
			_this.setEvents(test.N(bf))
		} else {
			_this.setEvents(test.N(sf))
		}
	case parser.CEEventParserINT_BIN, parser.CEEventParserINT_OCT, parser.CEEventParserINT_DEC, parser.CEEventParserINT_HEX:
		si, bi := parseInt(ttext)
		if bi != nil {
			_this.setEvents(test.N(bi))
		} else {
			if ttext[0] == '-' && si == 0 {
				_this.setEvents(test.N(compact_float.NegativeZero()))
			} else {
				_this.setEvents(test.N(si))
			}
		}
	default:
		panic(fmt.Errorf("BUG: Unexpected token type %v decoding \"%v\"", ttype, ttext))
	}
}
func (_this *eventListener) ExitEventPad(ctx *parser.EventPadContext) {
	_this.setEvents(test.PAD())
}
func (_this *eventListener) ExitEventLocalReference(ctx *parser.EventLocalReferenceContext) {
	_this.setEvents(test.REFL(getTokenText(ctx.GetChild(1))))
}
func (_this *eventListener) ExitEventRemoteReference(ctx *parser.EventRemoteReferenceContext) {
	_this.setEvents(test.REFR(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventResourceId(ctx *parser.EventResourceIdContext) {
	_this.setEvents(test.RID(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventString(ctx *parser.EventStringContext) {
	_this.setEvents(test.S(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventStructInstance(ctx *parser.EventStructInstanceContext) {
	_this.setEvents(test.SI(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventStructTemplate(ctx *parser.EventStructTemplateContext) {
	_this.setEvents(test.ST(getStringArg(ctx.GetChildren())))
}
func (_this *eventListener) ExitEventTime(ctx *parser.EventTimeContext) {
	text := getTokenText(ctx.GetChild(1))
	if child, ok := ctx.GetChildren()[1].(antlr.TerminalNode); ok {
		ttype := child.GetSymbol().GetTokenType()
		switch ttype {
		case parser.CEEventParserDATE:
			_this.setEvents(test.T(parseDate(text)))
		case parser.CEEventParserTIME:
			_this.setEvents(test.T(parseTime(text)))
		case parser.CEEventParserDATETIME:
			_this.setEvents(test.T(parseDateTime(text)))
		default:
			panic(fmt.Errorf("BUG: Unexpected terminal node type %v", ttype))
		}
	}
}
func (_this *eventListener) ExitEventUID(ctx *parser.EventUIDContext) {
	_this.setEvents(test.UID(parseUUID(getTokenText(ctx.GetChild(1)))))
}
func (_this *eventListener) ExitEventVersion(ctx *parser.EventVersionContext) {
	_this.setEvents(test.V(parseSmallUint(getTokenText(ctx.GetChild(1)))))
}
