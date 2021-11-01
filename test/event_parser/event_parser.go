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

package event_parser

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kstenerud/go-concise-encoding/test"

	"github.com/google/uuid"
	"github.com/kstenerud/go-compact-time"
)

// API

// Parse a string shorthand into an event for testing
func ParseEvent(eventStr string) *test.TEvent {
	components := eventNameMatcher.FindSubmatch([]byte(eventStr))
	if len(components) == 0 {
		panic(fmt.Errorf("Could not extract event name from [%v]", eventStr))
	}
	name := string(components[0])
	parser := eventParsersByName[name]
	if parser == nil {
		panic(fmt.Errorf("%v: Unknown event name", name))
	}
	return parser.ParseEvent(eventStr)
}

// Parse multiple events
func ParseEvents(eventStrings []string) []*test.TEvent {
	var index = 0
	var eventStr string

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("Event index %v [%v]: %w", index, eventStr, v))
			default:
				panic(v)
			}
		}
	}()

	var events []*test.TEvent

	for index, eventStr = range eventStrings {
		events = append(events, ParseEvent(eventStr))
	}
	return events
}

// Parsers

type eventParser struct {
	eventType    test.TEventType
	paramParsers []eventParamParser
}

func newParser(eventType test.TEventType, paramParsers ...eventParamParser) *eventParser {
	return &eventParser{
		eventType:    eventType,
		paramParsers: paramParsers,
	}
}

func (_this *eventParser) ParseEvent(eventStr string) *test.TEvent {
	switch len(_this.paramParsers) {
	case 0:
		return test.NewTEvent(_this.eventType, nil, nil)
	case 1:
		asBytes := []byte(eventStr)
		indices := eventNameAndWSMatcher.FindSubmatchIndex(asBytes)
		if len(indices) != 2 {
			panic(fmt.Errorf("Could not extract 1-param event components from [%v]", eventStr))
		}
		param := asBytes[indices[1]:]
		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](param),
			nil)
	case 2:
		asBytes := []byte(eventStr)
		indices := firstParamAndWSMatcher.FindSubmatchIndex(asBytes)
		if len(indices) != 4 {
			panic(fmt.Errorf("Could not extract 2-param event components from [%v]", eventStr))
		}
		param1 := asBytes[indices[2]:indices[3]]
		param2 := asBytes[indices[3]:]
		if param2[0] == ' ' || param2[0] == '\r' || param2[0] == '\n' || param2[0] == '\t' {
			param2 = param2[1:]
		}

		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](param1),
			_this.paramParsers[1](param2))
	default:
		panic(fmt.Errorf("BUG: Event parser has %v param parsers", len(_this.paramParsers)))
	}
}

type eventParamParser func(bytes []byte) interface{}

func containsEscape(bytes []byte) bool {
	for _, b := range bytes {
		if b == '\\' {
			return true
		}
	}
	return false
}

func parseString(bytes []byte) interface{} {
	return string(bytes)
}

func parseBool(bytes []byte) interface{} {
	asString := string(bytes)
	if asString == "true" || asString == "t" {
		return true
	}
	if asString == "false" || asString == "f" {
		return false
	}
	panic(fmt.Errorf("Error parsing bool [%v]", string(bytes)))
}

func getBase(bytes []byte) int {
	if len(bytes) > 1 && bytes[0] == '0' {
		switch bytes[1] {
		case 'b', 'B':
			return 2
		case 'o', 'O':
			return 8
		case 'x', 'X':
			return 16
		}
	}
	return 10
}

func parseInt(bytes []byte) interface{} {
	value, err := strconv.ParseInt(string(bytes), 0, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing int [%v]: %w", string(bytes), err))
	}
	return value
}

func parseUint(bytes []byte) interface{} {
	value, err := strconv.ParseUint(string(bytes), 0, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing uint [%v]: %w", string(bytes), err))
	}
	return value
}

func parseBigInt(bytes []byte) interface{} {
	return test.NewBigInt(string(bytes), getBase(bytes))
}

func parseFloat(bytes []byte) interface{} {
	value, err := strconv.ParseFloat(string(bytes), 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing float [%v]: %w", string(bytes), err))
	}
	return value
}

func parseBigFloat(bytes []byte) interface{} {
	sigDigits := 0
	for _, b := range bytes {
		if b >= '0' && b <= '9' {
			sigDigits++
		}
	}
	base := 10
	if len(bytes) > 1 && bytes[0] == '0' && (bytes[1] == 'x' || bytes[1] == 'X') {
		base = 16
		sigDigits--
	}

	return test.NewBigFloat(string(bytes), base, sigDigits)
}

func parseDecimalFloat(bytes []byte) interface{} {
	return test.NewDFloat(string(bytes))
}

func parseBigDecimalFloat(bytes []byte) interface{} {
	return test.NewBDF(string(bytes))
}

func parseUUID(bytes []byte) interface{} {
	value, err := uuid.Parse(string(bytes))
	if err != nil {
		panic(fmt.Errorf("Error parsing UUID [%v]: %w", string(bytes), err))
	}
	return value[:]
}

func bytesToInt(bytes []byte) int {
	if len(bytes) == 0 {
		panic(fmt.Errorf("Tried to parse empty byte array as int"))
	}
	sign := 1
	if bytes[0] == '-' {
		sign = -1
		bytes = bytes[1:]
	}
	accum := 0
	for _, b := range bytes {
		accum = accum*10 + int(b-'0')
	}
	return accum * sign
}

var dateMatcher = regexp.MustCompile(`^(-?\d+)-(\d+)-(\d+)`)
var timeMatcher = regexp.MustCompile(`^(\d+):(\d+):(\d+)(.\d+)?`)

func parseCompactTime(bytes []byte) interface{} {
	originalBytes := bytes
	indices := dateMatcher.FindSubmatchIndex(bytes)
	year := 0
	month := 0
	day := 0
	if len(indices) != 0 {
		year = bytesToInt(bytes[indices[2]:indices[3]])
		month = bytesToInt(bytes[indices[4]:indices[5]])
		day = bytesToInt(bytes[indices[6]:indices[7]])
		bytes = bytes[indices[1]:]
		if len(bytes) == 0 {
			date := compact_time.NewDate(year, month, day)
			if err := date.Validate(); err != nil {
				panic(fmt.Errorf("Error parsing date [%v]: %w", string(originalBytes), err))
			}
			return date
		}
		bytes = bytes[1:]
	}

	indices = timeMatcher.FindSubmatchIndex(bytes)
	if len(indices) == 0 {
		panic(fmt.Errorf("Malformed time [%v]", string(originalBytes)))
	}
	hour := bytesToInt(bytes[indices[2]:indices[3]])
	minute := bytesToInt(bytes[indices[4]:indices[5]])
	second := bytesToInt(bytes[indices[6]:indices[7]])
	subsecond := 0
	if indices[8] >= 0 {
		subsecond := bytesToInt(bytes[indices[8]:indices[9]])
		for i := indices[9] - indices[8]; i <= 9; i++ {
			subsecond *= 10
		}
	}
	bytes = bytes[indices[1]:]

	var tz compact_time.Timezone = compact_time.TZAtUTC()
	if len(bytes) != 0 {
		tz = compact_time.TZAtAreaLocation(string(bytes))
	}

	var time compact_time.Time
	if year == 0 && month == 0 && day == 0 {
		time = compact_time.NewTime(hour, minute, second, subsecond, tz)
	} else {
		time = compact_time.NewTimestamp(year, month, day, hour, minute, second, subsecond, tz)
	}
	if err := time.Validate(); err != nil {
		panic(fmt.Errorf("Error parsing time [%v]: %w", string(originalBytes), err))
	}
	return time

}

func parseBytes(data []byte) interface{} {
	buff := bytes.NewBuffer(make([]byte, 0, len(data)/2))
	var nextByte byte
	isFullByte := true

	for i := 0; i < len(data); i++ {
		b := data[i]
		switch b {
		case ' ', '\r', '\n', '\t':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			nextByte = (nextByte << 4) | (b - '0')
		case 'a', 'b', 'c', 'd', 'e', 'f':
			nextByte = (nextByte << 4) | (b - 'a' + 10)
		case 'A', 'B', 'C', 'D', 'E', 'F':
			nextByte = (nextByte << 4) | (b - 'A' + 10)
		default:
			panic(fmt.Errorf("Error parsing bytes: Invalid char in [%v]", string(data)))
		}
		isFullByte = !isFullByte
		if isFullByte {
			buff.WriteByte(nextByte)
		}
	}
	return buff.Bytes()
}

func parseTextAsBytes(data []byte) interface{} {
	return data
}

func newArrayParser(elemType reflect.Type, elementParser eventParamParser) eventParamParser {
	var typeAppropriate func(src interface{}) reflect.Value
	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		typeAppropriate = func(src interface{}) reflect.Value {
			value := reflect.New(elemType).Elem()
			value.SetInt(reflect.ValueOf(src).Int())
			return value
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeAppropriate = func(src interface{}) reflect.Value {
			value := reflect.New(elemType).Elem()
			value.SetUint(reflect.ValueOf(src).Uint())
			return value
		}
	case reflect.Float32, reflect.Float64:
		typeAppropriate = func(src interface{}) reflect.Value {
			value := reflect.New(elemType).Elem()
			value.SetFloat(reflect.ValueOf(src).Float())
			return value
		}
	}
	return func(data []byte) interface{} {
		fields := strings.Fields(string(data))
		slice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(fields))
		for _, field := range fields {
			elem := elementParser([]byte(field))
			slice = reflect.Append(slice, typeAppropriate(elem))
		}
		return slice.Interface()
	}
}

var eventNameMatcher = regexp.MustCompile(`^(\w+)`)
var eventNameAndWSMatcher = regexp.MustCompile(`^\w+\s+`)
var firstParamAndWSMatcher = regexp.MustCompile(`^\w+\s+(\w+)\s+`)

var eventParsersByName = make(map[string]*eventParser)

func init() {
	eventParsersByName["bd"] = newParser(test.TEventBeginDocument)
	eventParsersByName["ed"] = newParser(test.TEventEndDocument)
	eventParsersByName["v"] = newParser(test.TEventVersion, parseUint)
	eventParsersByName["tt"] = newParser(test.TEventTrue)
	eventParsersByName["ff"] = newParser(test.TEventFalse)
	eventParsersByName["i"] = newParser(test.TEventInt, parseInt)
	eventParsersByName["f"] = newParser(test.TEventFloat, parseFloat)
	eventParsersByName["bf"] = newParser(test.TEventBigFloat, parseBigFloat)
	eventParsersByName["df"] = newParser(test.TEventDecimalFloat, parseDecimalFloat)
	eventParsersByName["bdf"] = newParser(test.TEventBigDecimalFloat, parseBigDecimalFloat)
	eventParsersByName["n"] = newParser(test.TEventNil)
	eventParsersByName["com"] = newParser(test.TEventComment, parseBool, parseString)
	eventParsersByName["b"] = newParser(test.TEventBool, parseBool)
	eventParsersByName["pi"] = newParser(test.TEventPInt, parseUint)
	eventParsersByName["ni"] = newParser(test.TEventNInt, parseUint)
	eventParsersByName["bi"] = newParser(test.TEventBigInt, parseBigInt)
	eventParsersByName["nan"] = newParser(test.TEventNan)
	eventParsersByName["snan"] = newParser(test.TEventSNan)
	eventParsersByName["uid"] = newParser(test.TEventUID, parseUUID)
	// func GT(v time.Time) *test.TEvent            { return test.NewTEvent(test.TEventTime, v, nil) }
	eventParsersByName["ct"] = newParser(test.TEventCompactTime, parseCompactTime)
	eventParsersByName["s"] = newParser(test.TEventString, parseString)
	eventParsersByName["rid"] = newParser(test.TEventResourceID, parseString)
	eventParsersByName["cub"] = newParser(test.TEventCustomBinary, parseBytes)
	eventParsersByName["cut"] = newParser(test.TEventCustomText, parseString)
	// func AB(l uint64, v []byte) *test.TEvent     { return test.NewTEvent(test.TEventArrayBoolean, l, v) }
	eventParsersByName["ai8"] = newParser(test.TEventArrayInt8, newArrayParser(reflect.TypeOf(int8(0)), parseInt))
	eventParsersByName["ai16"] = newParser(test.TEventArrayInt16, newArrayParser(reflect.TypeOf(int16(0)), parseInt))
	eventParsersByName["ai32"] = newParser(test.TEventArrayInt32, newArrayParser(reflect.TypeOf(int32(0)), parseInt))
	eventParsersByName["ai64"] = newParser(test.TEventArrayInt64, newArrayParser(reflect.TypeOf(int64(0)), parseInt))
	eventParsersByName["au8"] = newParser(test.TEventArrayUint8, parseBytes)
	eventParsersByName["au16"] = newParser(test.TEventArrayUint16, newArrayParser(reflect.TypeOf(uint16(0)), parseUint))
	eventParsersByName["au32"] = newParser(test.TEventArrayUint32, newArrayParser(reflect.TypeOf(uint32(0)), parseUint))
	eventParsersByName["au64"] = newParser(test.TEventArrayUint64, newArrayParser(reflect.TypeOf(uint64(0)), parseUint))
	// func AF16(v []byte) *test.TEvent             { return test.NewTEvent(test.TEventArrayFloat16, v, nil) }
	eventParsersByName["af32"] = newParser(test.TEventArrayFloat32, newArrayParser(reflect.TypeOf(float32(0)), parseFloat))
	eventParsersByName["af64"] = newParser(test.TEventArrayFloat64, newArrayParser(reflect.TypeOf(float64(0)), parseFloat))
	// func AUU(v []byte) *test.TEvent              { return test.NewTEvent(test.TEventArrayUID, v, nil) }
	eventParsersByName["sb"] = newParser(test.TEventStringBegin)
	eventParsersByName["rb"] = newParser(test.TEventResourceIDBegin)
	eventParsersByName["rrb"] = newParser(test.TEventRemoteRefBegin)
	eventParsersByName["rbcat"] = newParser(test.TEventResourceIDCatBegin)
	eventParsersByName["cbb"] = newParser(test.TEventCustomBinaryBegin)
	eventParsersByName["ctb"] = newParser(test.TEventCustomTextBegin)
	eventParsersByName["abb"] = newParser(test.TEventArrayBooleanBegin)
	eventParsersByName["ai8b"] = newParser(test.TEventArrayInt8Begin)
	eventParsersByName["ai16b"] = newParser(test.TEventArrayInt16Begin)
	eventParsersByName["ai32b"] = newParser(test.TEventArrayInt32Begin)
	eventParsersByName["ai64b"] = newParser(test.TEventArrayInt64Begin)
	eventParsersByName["au8b"] = newParser(test.TEventArrayUint8Begin)
	eventParsersByName["au16b"] = newParser(test.TEventArrayUint16Begin)
	eventParsersByName["au32b"] = newParser(test.TEventArrayUint32Begin)
	eventParsersByName["au64b"] = newParser(test.TEventArrayUint64Begin)
	eventParsersByName["auub"] = newParser(test.TEventArrayUIDBegin)
	eventParsersByName["mb"] = newParser(test.TEventMediaBegin)
	eventParsersByName["ac"] = newParser(test.TEventArrayChunk, parseUint, parseBool)
	eventParsersByName["ad"] = newParser(test.TEventArrayData, parseBytes)
	eventParsersByName["adt"] = newParser(test.TEventArrayData, parseTextAsBytes)
	eventParsersByName["l"] = newParser(test.TEventList)
	eventParsersByName["m"] = newParser(test.TEventMap)
	eventParsersByName["mup"] = newParser(test.TEventMarkup, parseString)
	eventParsersByName["node"] = newParser(test.TEventNode)
	eventParsersByName["edge"] = newParser(test.TEventEdge)
	eventParsersByName["e"] = newParser(test.TEventEnd)
	eventParsersByName["mark"] = newParser(test.TEventMarker, parseString)
	eventParsersByName["ref"] = newParser(test.TEventReference, parseString)
	eventParsersByName["rref"] = newParser(test.TEventRemoteRef, parseString)
	eventParsersByName["const"] = newParser(test.TEventConstant, parseString)
}
