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
)

// API

// Parse a string shorthand into an event for testing
func ParseEvent(eventStr string) *test.TEvent {
	components := parseEventName.FindSubmatch([]byte(eventStr))
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
		components := parseOneParamEvent.FindSubmatch([]byte(eventStr))
		if len(components) == 0 {
			panic(fmt.Errorf("Could not extract event components from [%v]", eventStr))
		}
		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](components[1]),
			nil)
	case 2:
		components := parseOneParamEvent.FindSubmatch([]byte(eventStr))
		if len(components) == 0 {
			panic(fmt.Errorf("Could not extract event components from [%v]", eventStr))
		}
		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](components[1]),
			_this.paramParsers[1](components[2]))
	default:
		panic(fmt.Errorf("BUG: Event parser has %v param parsers", len(_this.paramParsers)))
	}
}

type eventParamParser func(bytes []byte) interface{}

func passthroughString(bytes []byte) interface{} {
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

func parseInt(bytes []byte) interface{} {
	base := 10
	if len(bytes) > 1 && bytes[0] == '0' && (bytes[1] == 'x' || bytes[1] == 'X') {
		base = 16
	}
	value, err := strconv.ParseInt(string(bytes), base, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing int [%v]: %w", string(bytes), err))
	}
	return value
}

func parseUint(bytes []byte) interface{} {
	base := 10
	if len(bytes) > 1 && bytes[0] == '0' && (bytes[1] == 'x' || bytes[1] == 'X') {
		base = 16
	}
	value, err := strconv.ParseUint(string(bytes), base, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing uint [%v]: %w", string(bytes), err))
	}
	return value
}

func parseBigInt(bytes []byte) interface{} {
	base := 10
	if len(bytes) > 1 && bytes[0] == '0' && (bytes[1] == 'x' || bytes[1] == 'X') {
		base = 16
	}
	return test.NewBigInt(string(bytes), base)
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
	return value
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

var parseEventName = regexp.MustCompile("^([a-z0-9]+)")
var parseOneParamEvent = regexp.MustCompile("^[a-z0-9]+ (.+)$")
var parseTwoParamEvent = regexp.MustCompile("^[a-z0-9]+ (\\w+) (.+)$")

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
	eventParsersByName["com"] = newParser(test.TEventComment, parseBool, passthroughString)
	eventParsersByName["b"] = newParser(test.TEventBool, parseBool)
	eventParsersByName["pi"] = newParser(test.TEventPInt, parseUint)
	eventParsersByName["ni"] = newParser(test.TEventNInt, parseUint)
	eventParsersByName["bi"] = newParser(test.TEventBigInt, parseBigInt)
	eventParsersByName["nan"] = newParser(test.TEventNan)
	eventParsersByName["snan"] = newParser(test.TEventSNan)
	eventParsersByName["uid"] = newParser(test.TEventUID, parseUUID)
	// func GT(v time.Time) *test.TEvent            { return test.NewTEvent(test.TEventTime, v, nil) }
	// func CT(v compact_time.Time) *test.TEvent    { return EventOrNil(test.TEventCompactTime, v) }
	eventParsersByName["s"] = newParser(test.TEventString, passthroughString)
	eventParsersByName["rid"] = newParser(test.TEventResourceID, passthroughString)
	eventParsersByName["ridref"] = newParser(test.TEventResourceIDRef, passthroughString)
	eventParsersByName["cub"] = newParser(test.TEventCustomBinary, parseBytes)
	eventParsersByName["cut"] = newParser(test.TEventCustomText, passthroughString)
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
	eventParsersByName["rrb"] = newParser(test.TEventResourceIDRefBegin)
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
	eventParsersByName["l"] = newParser(test.TEventList)
	eventParsersByName["m"] = newParser(test.TEventMap)
	eventParsersByName["mup"] = newParser(test.TEventMarkup, passthroughString)
	eventParsersByName["node"] = newParser(test.TEventNode)
	eventParsersByName["edge"] = newParser(test.TEventEdge)
	eventParsersByName["e"] = newParser(test.TEventEnd)
	eventParsersByName["mark"] = newParser(test.TEventMarker, passthroughString)
	eventParsersByName["ref"] = newParser(test.TEventReference, passthroughString)
	eventParsersByName["const"] = newParser(test.TEventConstant, passthroughString)
}
