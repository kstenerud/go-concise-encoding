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
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/test"
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
	return parser(eventStr)
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

	events := []*test.TEvent{test.BD()}
	for index, eventStr = range eventStrings {
		events = append(events, ParseEvent(eventStr))
	}
	events = append(events, test.ED())

	return events
}

// Parsers

type parseEvent func(eventStr string) *test.TEvent

func parseNumericEvent(eventStr string) *test.TEvent {
	iparseNumeric := func(data []byte) (test.TEventType, interface{}) {
		str := string(data)
		if strings.TrimSpace(str) == "-0" {
			return test.TEventNInt, uint64(0)
		}

		if value, err := strconv.ParseInt(str, 0, 64); err == nil {
			return test.TEventInt, value
		}

		bi := &big.Int{}
		if _, success := bi.SetString(str, 0); success {
			return test.TEventBigInt, bi
		}

		if value, err := compact_float.DFloatFromString(str); err == nil {
			return test.TEventDecimalFloat, value
		}

		if strings.Contains(str, "0x") {
			digitCount := len(strings.Split(str, "x")[1])
			bf := &big.Float{}
			bf.SetPrec(uint(digitCount) * 4)
			if _, success := bf.SetString(str); success {
				f64, accuracy := bf.Float64()
				if accuracy == big.Exact {
					return test.TEventFloat, f64
				} else {
					return test.TEventBigFloat, bf
				}
			} else {
				panic(fmt.Errorf("could not convert %v to float", str))
			}
		}

		return test.TEventBigDecimalFloat, test.NewBDF(strings.Replace(str, ",", ".", 1))
	}

	eventType, v := iparseNumeric(get1ParamArg(eventStr))
	return test.NewTEvent(eventType, v, nil)
}

type generalEventParser struct {
	eventType    test.TEventType
	paramParsers []eventParamParser
}

func newParser(eventType test.TEventType, paramParsers ...eventParamParser) *generalEventParser {
	return &generalEventParser{
		eventType:    eventType,
		paramParsers: paramParsers,
	}
}

func (_this *generalEventParser) ParseEvent(eventStr string) *test.TEvent {
	switch len(_this.paramParsers) {
	case 0:
		return test.NewTEvent(_this.eventType, nil, nil)
	case 1:
		param := get1ParamArg(eventStr)
		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](param),
			nil)
	case 2:
		param1, param2 := get2ParamArg(eventStr)
		return test.NewTEvent(_this.eventType,
			_this.paramParsers[0](param1),
			_this.paramParsers[1](param2))
	default:
		panic(fmt.Errorf("BUG: Event parser has %v param parsers", len(_this.paramParsers)))
	}
}

func get1ParamArg(eventStr string) []byte {
	asBytes := []byte(eventStr)
	indices := eventNameAndWSMatcher.FindSubmatchIndex(asBytes)
	if len(indices) != 2 {
		panic(fmt.Errorf("Event [%v] requires 1 parameter", eventStr))
	}
	return asBytes[indices[1]:]
}

func get2ParamArg(eventStr string) ([]byte, []byte) {
	asBytes := []byte(eventStr)
	indices := firstParamAndWSMatcher.FindSubmatchIndex(asBytes)
	if len(indices) != 4 {
		panic(fmt.Errorf("Event [%v] requires 2 parameters", eventStr))
	}
	param1 := asBytes[indices[2]:indices[3]]
	param2 := asBytes[indices[3]:]
	if param2[0] == ' ' || param2[0] == '\r' || param2[0] == '\n' || param2[0] == '\t' {
		param2 = param2[1:]
	}
	return param1, param2
}

type eventParamParser func(bytes []byte) interface{}

func parseString(data []byte) interface{} {
	return string(data)
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

func parseHex(bytes []byte) (result uint64) {
	for _, b := range bytes {
		switch b {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			result = (result << 4) | uint64(b-'0')
		case 'a', 'b', 'c', 'd', 'e', 'f':
			result = (result << 4) | uint64(b-'a'+10)
		case 'A', 'B', 'C', 'D', 'E', 'F':
			result = (result << 4) | uint64(b-'A'+10)
		default:
			panic(fmt.Errorf("Error parsing hexadecimal: Invalid char [%c] in [%v]", b, string(bytes)))
		}
	}
	return
}

func parseUintHex(bytes []byte) interface{} {
	if len(bytes) == 0 {
		panic(fmt.Errorf("Error parsing hexadecimal: no data"))
	}
	return parseHex(bytes)
}

func parseIntHex(bytes []byte) interface{} {
	if len(bytes) == 0 {
		panic(fmt.Errorf("Error parsing hexadecimal: no data"))
	}
	sign := int64(0)
	if bytes[0] == '-' {
		sign = -1
		bytes = bytes[1:]
	}
	value := parseHex(bytes)
	if value&0x8000000000000000 != 0 {
		panic(fmt.Errorf("Overflow parsing [%v]", string(bytes)))
	}
	return sign * int64(value)
}

func parseBigInt(bytes []byte) interface{} {
	return test.NewBigInt(string(bytes))
}

func parseFloat(bytes []byte) interface{} {
	value, err := strconv.ParseFloat(string(bytes), 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing float [%v]: %w", string(bytes), err))
	}
	return value
}

func parseBigFloat(bytes []byte) interface{} {
	return test.NewBigFloat(string(bytes))
}

func parseDecimalFloat(bytes []byte) interface{} {
	return test.NewDFloat(string(bytes))
}

func parseBigDecimalFloat(bytes []byte) interface{} {
	return test.NewBDF(string(bytes))
}

var uuidMatcher = regexp.MustCompile(`^([0-9a-fA-F]{8})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{12})`)

func parseUUID(data []byte) interface{} {
	components := uuidMatcher.FindSubmatch(data)
	if len(components) == 0 {
		panic(fmt.Errorf("Error parsing UUID [%v]: not a UUID", string(data)))
	}
	buff := bytes.Buffer{}
	for iComponent := 1; iComponent < len(components); iComponent++ {
		component := components[iComponent]
		for iByte := 0; iByte < len(component); iByte += 2 {
			buff.WriteByte(byte(parseHex(component[iByte : iByte+2])))
		}
	}
	return buff.Bytes()
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
		digit := int(b - '0')
		if digit < 0 || digit > 9 {
			panic(fmt.Errorf("%c: Invalid integer digit", b))
		}
		accum = accum*10 + digit
	}
	return accum * sign
}

var dateMatcher = regexp.MustCompile(`^(-?\d+)-(\d+)-(\d+)`)

func tryParseDate(bytes []byte) (date compact_time.Time, remainingBytes []byte) {
	remainingBytes = bytes
	indices := dateMatcher.FindSubmatchIndex(bytes)
	if len(indices) == 0 {
		return
	}

	year := bytesToInt(bytes[indices[2]:indices[3]])
	month := bytesToInt(bytes[indices[4]:indices[5]])
	day := bytesToInt(bytes[indices[6]:indices[7]])

	remainingBytes = bytes[indices[1]:]
	if len(remainingBytes) > 0 && remainingBytes[0] == '/' {
		remainingBytes = remainingBytes[1:]
	}
	date = compact_time.NewDate(year, month, day)
	if err := date.Validate(); err != nil {
		panic(fmt.Errorf("Error parsing date from [%v]: %w", string(bytes), err))
	}
	return
}

var utcOffsetMatcher = regexp.MustCompile(`^[+-](\d\d)(\d\d)$`)

func parseTZUTCOffset(data []byte) compact_time.Timezone {
	components := utcOffsetMatcher.FindSubmatch(data)
	if len(components) == 0 {
		panic(fmt.Errorf("Could not parse UTC offset from [%v]", string(data)))
	}
	sign := 1
	if data[0] == '-' {
		sign = -1
	}
	hours := bytesToInt(components[1])
	minutes := bytesToInt(components[2])
	return compact_time.TZWithMiutesOffsetFromUTC(sign * (hours*60 + minutes))
}

var latLongMatcher = regexp.MustCompile(`^(-?\d+(\.\d+)?)/(-?\d+(\.\d+)?)$`)

func parseTZLatLong(data []byte) compact_time.Timezone {
	components := latLongMatcher.FindSubmatch(data)
	if len(components) == 0 {
		panic(fmt.Errorf("Could not parse lat/long from [%v]", string(data)))
	}
	lat, err := strconv.ParseFloat(string(components[1]), 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing latitude from [%v]: %w", string(components[1]), err))
	}
	long, err := strconv.ParseFloat(string(components[3]), 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing longitude from [%v]: %w", string(components[3]), err))
	}
	return compact_time.TZAtLatLong(int(lat*100), int(long*100))
}

func parseTZAreaLocation(data []byte) compact_time.Timezone {
	return compact_time.TZAtAreaLocation(string(data))
}

func parseTZAreaLocationOrLatLong(data []byte) compact_time.Timezone {
	if len(data) == 0 {
		panic(fmt.Errorf("TZ data missing"))
	}
	switch data[0] {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return parseTZLatLong(data)
	default:
		return parseTZAreaLocation(data)
	}
}

func parseTimezone(data []byte) (tz compact_time.Timezone) {
	tz = compact_time.TZAtUTC()
	if len(data) == 0 {
		return
	}

	switch data[0] {
	case '+', '-':
		return parseTZUTCOffset(data)
	case '/':
		return parseTZAreaLocationOrLatLong(data[1:])
	default:
		panic(fmt.Errorf("%v: Invalid timezone", string(data)))
	}
}

var timeMatcher = regexp.MustCompile(`^(\d+):(\d+):(\d+)(\.\d+)?`)

func parseTime(data []byte) (t compact_time.Time) {
	indices := timeMatcher.FindSubmatchIndex(data)
	if len(indices) == 0 {
		return
	}
	hour := bytesToInt(data[indices[2]:indices[3]])
	minute := bytesToInt(data[indices[4]:indices[5]])
	second := bytesToInt(data[indices[6]:indices[7]])
	subsecond := 0
	if indices[8] >= 0 {
		begin := indices[8] + 1
		end := indices[9]
		subsecond = bytesToInt(data[begin:end])
		for i := end - begin; i < 9; i++ {
			subsecond *= 10
		}
	}
	data = data[indices[1]:]
	tz := parseTimezone(data)
	return compact_time.NewTime(hour, minute, second, subsecond, tz)
}

func parseTemporal(data []byte) interface{} {
	originalBytes := data
	var datePart compact_time.Time
	datePart, data = tryParseDate(data)
	timePart := parseTime(data)

	if datePart.IsZeroValue() && timePart.IsZeroValue() {
		panic(fmt.Errorf("Could not parse date [%v]: no date data found", string(originalBytes)))
	}

	if !datePart.IsZeroValue() {
		if err := datePart.Validate(); err != nil {
			panic(fmt.Errorf("Error parsing date [%v]: %w", string(originalBytes), err))
		}
		if timePart.IsZeroValue() {
			return datePart
		}
		timePart.Year = datePart.Year
		timePart.Month = datePart.Month
		timePart.Day = datePart.Day
		timePart.Type = compact_time.TimeTypeTimestamp
	}
	if err := timePart.Validate(); err != nil {
		panic(fmt.Errorf("Error parsing time value [%v]: %w", string(originalBytes), err))
	}
	return timePart
}

func parseTextAsBytes(data []byte) interface{} {
	return data
}

var bitArrayMatcher = regexp.MustCompile(`(\s*[01])+`)

func parseBitArrayEvent(eventStr string) *test.TEvent {
	var array []byte
	asBytes := []byte(eventStr[3:])

	iBytes := 0
	generator := func() (next byte, bitCount int) {
		for iBytes < len(asBytes) {
			b := asBytes[iBytes]
			iBytes++
			switch b {
			case '1':
				next |= byte(1 << bitCount)
				bitCount++
			case '0':
				bitCount++
			case ' ', '\r', '\n', '\t':
				// Skip whitespace
			default:
				panic(fmt.Errorf("[%c]: Invalid bit array character", b))
			}
			if bitCount >= 8 {
				return
			}
		}
		return
	}

	// first byte low bit is first bit of array
	totalBits := uint64(0)
	for {
		b, bitCount := generator()
		if bitCount == 0 {
			break
		}
		totalBits += uint64(bitCount)
		array = append(array, b)
	}

	return test.NewTEvent(test.TEventArrayBoolean,
		totalBits,
		array)
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
	case reflect.Slice:
		typeAppropriate = func(src interface{}) reflect.Value {
			value := reflect.New(elemType).Elem()
			value.Set(reflect.ValueOf(src))
			return value
		}
	default:
		panic(fmt.Errorf("No parser defined for array type %v", elemType))
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
var eventNameAndWSMatcher = regexp.MustCompile(`^\w+\s`)
var firstParamAndWSMatcher = regexp.MustCompile(`^\w+\s+(\w+)\s`)

var eventParsersByName = make(map[string]parseEvent)

func init() {
	eventParsersByName["bd"] = newParser(test.TEventBeginDocument).ParseEvent
	eventParsersByName["ed"] = newParser(test.TEventEndDocument).ParseEvent
	eventParsersByName["v"] = newParser(test.TEventVersion, parseUint).ParseEvent
	eventParsersByName["tt"] = newParser(test.TEventTrue).ParseEvent
	eventParsersByName["ff"] = newParser(test.TEventFalse).ParseEvent
	eventParsersByName["n"] = parseNumericEvent
	eventParsersByName["i"] = newParser(test.TEventInt, parseInt).ParseEvent
	eventParsersByName["f"] = newParser(test.TEventFloat, parseFloat).ParseEvent
	eventParsersByName["bf"] = newParser(test.TEventBigFloat, parseBigFloat).ParseEvent
	eventParsersByName["df"] = newParser(test.TEventDecimalFloat, parseDecimalFloat).ParseEvent
	eventParsersByName["bdf"] = newParser(test.TEventBigDecimalFloat, parseBigDecimalFloat).ParseEvent
	eventParsersByName["null"] = newParser(test.TEventNull).ParseEvent
	eventParsersByName["com"] = newParser(test.TEventComment, parseBool, parseString).ParseEvent
	eventParsersByName["b"] = newParser(test.TEventBool, parseBool).ParseEvent
	eventParsersByName["pi"] = newParser(test.TEventPInt, parseUint).ParseEvent
	eventParsersByName["ni"] = newParser(test.TEventNInt, parseUint).ParseEvent
	eventParsersByName["bi"] = newParser(test.TEventBigInt, parseBigInt).ParseEvent
	eventParsersByName["nan"] = newParser(test.TEventNan).ParseEvent
	eventParsersByName["snan"] = newParser(test.TEventSNan).ParseEvent
	eventParsersByName["uid"] = newParser(test.TEventUID, parseUUID).ParseEvent
	eventParsersByName["t"] = newParser(test.TEventCompactTime, parseTemporal).ParseEvent
	eventParsersByName["s"] = newParser(test.TEventString, parseString).ParseEvent
	eventParsersByName["rid"] = newParser(test.TEventResourceID, parseString).ParseEvent
	eventParsersByName["cb"] = newParser(test.TEventCustomBinary, newArrayParser(reflect.TypeOf(uint8(0)), parseUintHex)).ParseEvent
	eventParsersByName["ct"] = newParser(test.TEventCustomText, parseString).ParseEvent
	eventParsersByName["ab"] = parseBitArrayEvent
	eventParsersByName["ai8"] = newParser(test.TEventArrayInt8, newArrayParser(reflect.TypeOf(int8(0)), parseInt)).ParseEvent
	eventParsersByName["ai8x"] = newParser(test.TEventArrayInt8, newArrayParser(reflect.TypeOf(int8(0)), parseIntHex)).ParseEvent
	eventParsersByName["ai16"] = newParser(test.TEventArrayInt16, newArrayParser(reflect.TypeOf(int16(0)), parseInt)).ParseEvent
	eventParsersByName["ai16x"] = newParser(test.TEventArrayInt16, newArrayParser(reflect.TypeOf(int16(0)), parseIntHex)).ParseEvent
	eventParsersByName["ai32"] = newParser(test.TEventArrayInt32, newArrayParser(reflect.TypeOf(int32(0)), parseInt)).ParseEvent
	eventParsersByName["ai32x"] = newParser(test.TEventArrayInt32, newArrayParser(reflect.TypeOf(int32(0)), parseIntHex)).ParseEvent
	eventParsersByName["ai64"] = newParser(test.TEventArrayInt64, newArrayParser(reflect.TypeOf(int64(0)), parseInt)).ParseEvent
	eventParsersByName["ai64x"] = newParser(test.TEventArrayInt64, newArrayParser(reflect.TypeOf(int64(0)), parseIntHex)).ParseEvent
	eventParsersByName["au8"] = newParser(test.TEventArrayUint8, newArrayParser(reflect.TypeOf(uint8(0)), parseUint)).ParseEvent
	eventParsersByName["au8x"] = newParser(test.TEventArrayUint8, newArrayParser(reflect.TypeOf(uint8(0)), parseUintHex)).ParseEvent
	eventParsersByName["au16"] = newParser(test.TEventArrayUint16, newArrayParser(reflect.TypeOf(uint16(0)), parseUint)).ParseEvent
	eventParsersByName["au16x"] = newParser(test.TEventArrayUint16, newArrayParser(reflect.TypeOf(uint16(0)), parseUintHex)).ParseEvent
	eventParsersByName["au32"] = newParser(test.TEventArrayUint32, newArrayParser(reflect.TypeOf(uint32(0)), parseUint)).ParseEvent
	eventParsersByName["au32x"] = newParser(test.TEventArrayUint32, newArrayParser(reflect.TypeOf(uint32(0)), parseUintHex)).ParseEvent
	eventParsersByName["au64"] = newParser(test.TEventArrayUint64, newArrayParser(reflect.TypeOf(uint64(0)), parseUint)).ParseEvent
	eventParsersByName["au64x"] = newParser(test.TEventArrayUint64, newArrayParser(reflect.TypeOf(uint64(0)), parseUintHex)).ParseEvent
	// TODO: eventParsersByName["af16"] = newParser(test.TEventArrayFloat16, newArrayParser(reflect.TypeOf(float32(0)), parseFloat)).ParseEvent
	eventParsersByName["af32"] = newParser(test.TEventArrayFloat32, newArrayParser(reflect.TypeOf(float32(0)), parseFloat)).ParseEvent
	eventParsersByName["af64"] = newParser(test.TEventArrayFloat64, newArrayParser(reflect.TypeOf(float64(0)), parseFloat)).ParseEvent
	eventParsersByName["au"] = newParser(test.TEventArrayUID, newArrayParser(reflect.TypeOf([]byte{}), parseUUID)).ParseEvent
	eventParsersByName["sb"] = newParser(test.TEventStringBegin).ParseEvent
	eventParsersByName["rb"] = newParser(test.TEventResourceIDBegin).ParseEvent
	eventParsersByName["rrb"] = newParser(test.TEventRemoteRefBegin).ParseEvent
	eventParsersByName["cbb"] = newParser(test.TEventCustomBinaryBegin).ParseEvent
	eventParsersByName["ctb"] = newParser(test.TEventCustomTextBegin).ParseEvent
	eventParsersByName["abb"] = newParser(test.TEventArrayBooleanBegin).ParseEvent
	eventParsersByName["ai8b"] = newParser(test.TEventArrayInt8Begin).ParseEvent
	eventParsersByName["ai16b"] = newParser(test.TEventArrayInt16Begin).ParseEvent
	eventParsersByName["ai32b"] = newParser(test.TEventArrayInt32Begin).ParseEvent
	eventParsersByName["ai64b"] = newParser(test.TEventArrayInt64Begin).ParseEvent
	eventParsersByName["au8b"] = newParser(test.TEventArrayUint8Begin).ParseEvent
	eventParsersByName["au16b"] = newParser(test.TEventArrayUint16Begin).ParseEvent
	eventParsersByName["au32b"] = newParser(test.TEventArrayUint32Begin).ParseEvent
	eventParsersByName["au64b"] = newParser(test.TEventArrayUint64Begin).ParseEvent
	eventParsersByName["aub"] = newParser(test.TEventArrayUIDBegin).ParseEvent
	eventParsersByName["mb"] = newParser(test.TEventMediaBegin).ParseEvent
	eventParsersByName["ac"] = newParser(test.TEventArrayChunk, parseUint, parseBool).ParseEvent
	eventParsersByName["ad"] = newParser(test.TEventArrayData, newArrayParser(reflect.TypeOf(uint8(0)), parseUintHex)).ParseEvent
	eventParsersByName["at"] = newParser(test.TEventArrayData, parseTextAsBytes).ParseEvent
	eventParsersByName["l"] = newParser(test.TEventList).ParseEvent
	eventParsersByName["m"] = newParser(test.TEventMap).ParseEvent
	eventParsersByName["mup"] = newParser(test.TEventMarkup, parseString).ParseEvent
	eventParsersByName["node"] = newParser(test.TEventNode).ParseEvent
	eventParsersByName["edge"] = newParser(test.TEventEdge).ParseEvent
	eventParsersByName["e"] = newParser(test.TEventEnd).ParseEvent
	eventParsersByName["mark"] = newParser(test.TEventMarker, parseString).ParseEvent
	eventParsersByName["ref"] = newParser(test.TEventReference, parseString).ParseEvent
	eventParsersByName["rref"] = newParser(test.TEventRemoteRef, parseString).ParseEvent
	eventParsersByName["const"] = newParser(test.TEventConstant, parseString).ParseEvent
}
