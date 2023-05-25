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

package tests

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"time"

	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/test/test_runner"
	"github.com/kstenerud/go-concise-encoding/version"
)

func truncate(data []byte, count int) []byte {
	return data[:len(data)-count]
}

func newUnitTest(name string, mustSucceed []*test_runner.MustSucceedTest, mustFail []*test_runner.MustFailTest) *test_runner.UnitTest {
	unitTest := &test_runner.UnitTest{
		Name: name,
	}

	if mustSucceed != nil {
		unitTest.MustSucceed = mustSucceed
	}
	if mustFail != nil {
		unitTest.MustFail = mustFail
	}

	return unitTest
}

func newMustSucceedUnitTest(name string, mustSucceed ...*test_runner.MustSucceedTest) *test_runner.UnitTest {
	return &test_runner.UnitTest{
		Name:        name,
		MustSucceed: mustSucceed,
	}
}

func newMustFailUnitTest(name string, mustFail ...*test_runner.MustFailTest) *test_runner.UnitTest {
	return &test_runner.UnitTest{
		Name:     name,
		MustFail: mustFail,
	}
}

type Directions int

const (
	DirectionsNone   Directions = 0
	DirectionFromCBE Directions = 1 << iota
	DirectionToCBE
	DirectionFromCTE
	DirectionToCTE
	DirectionsAll Directions = math.MaxInt
)

func (_this Directions) includes(direction Directions) bool {
	return (_this & direction) != 0
}

func (_this Directions) except(direction Directions) Directions {
	return _this & ^direction
}

func (_this Directions) and(direction Directions) Directions {
	return _this | direction
}

func newMustSucceedTest(directions Directions,
	config *configuration.CTEEncoderConfiguration,
	events ...test.Event) *test_runner.MustSucceedTest {

	if containsRecords(events) {
		oldEvents := events
		events = []test.Event{test.EvST, test.EvS, test.EvE}
		events = append(events, oldEvents...)
	}

	hasFromCBE := directions.includes(DirectionFromCBE) && canConvertFromCBE(config, events...)
	hasToCBE := directions.includes(DirectionToCBE) && canConvertToCBE(config, events...)
	cbeDocument := generateCBE(events...)
	cbe := []byte{}
	fromCBE := []byte{}
	toCBE := []byte{}
	if hasFromCBE && hasToCBE {
		cbe = cbeDocument
	} else if hasFromCBE {
		fromCBE = cbeDocument
	} else if hasToCBE {
		toCBE = cbeDocument
	}

	hasFromCTE := directions.includes(DirectionFromCTE) && canConvertFromCTE(config, events...)
	hasToCTE := directions.includes(DirectionToCTE) && canConvertToCTE(config, events...)
	cteDocument := generateCTE(config, events...)
	cte := ""
	fromCTE := ""
	toCTE := ""
	if hasFromCTE && hasToCTE {
		cte = cteDocument
	} else if hasFromCTE {
		fromCTE = cteDocument
	} else if hasToCTE {
		toCTE = cteDocument
	}

	return &test_runner.MustSucceedTest{
		BaseTest: test_runner.BaseTest{
			CBE:    cbe,
			CTE:    cte,
			Events: stringifyEvents(events...),
		},
		FromCBE: fromCBE,
		ToCBE:   toCBE,
		FromCTE: fromCTE,
		ToCTE:   toCTE,
	}
}

func newMustFailTest(testType testType, events ...test.Event) *test_runner.MustFailTest {
	if containsRecords(events) {
		oldEvents := events
		events = []test.Event{test.EvST, test.EvS, test.EvE}
		events = append(events, oldEvents...)
	}

	switch testType {
	case testTypeCbe:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CBE: generateCBE(events...)}}
	case testTypeCte:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{CTE: generateCTE(nil, events...)}}
	case testTypeEvents:
		return &test_runner.MustFailTest{BaseTest: test_runner.BaseTest{Events: stringifyEvents(events...)}}
	default:
		panic(fmt.Errorf("%v: unknown mustFail test type", testType))
	}
}

func newCustomMustFailTest(cteContents string) *test_runner.MustFailTest {
	return &test_runner.MustFailTest{
		BaseTest: test_runner.BaseTest{
			CTE:         cteContents,
			RawDocument: true,
		},
	}
}

func containsRecords(events test.Events) bool {
	for _, event := range events {
		if event.IsEquivalentTo(EvSI) {
			return true
		}
	}
	return false
}

func generateCBE(events ...test.Event) []byte {
	if mustNotConvertToCBE(events...) {
		return []byte{}
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("in events [%v]: %w", test.Events(events), v))
			default:
				panic(fmt.Errorf("in events [%v]: %v", test.Events(events), v))
			}
		}
	}()

	buffer := bytes.Buffer{}
	encoder := cbe.NewEncoder(nil)
	encoder.PrepareToEncode(&buffer)
	encoder.OnBeginDocument()
	encoder.OnVersion(version.ConciseEncodingVersion)
	for _, event := range events {
		event.Invoke(encoder)
	}
	encoder.OnEndDocument()
	result := buffer.Bytes()
	return result[2:]
}

func generateCTE(config *configuration.CTEEncoderConfiguration, events ...test.Event) string {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				panic(fmt.Errorf("in events [%v]: %w", test.Events(events), v))
			default:
				panic(fmt.Errorf("in events [%v]: %v", test.Events(events), v))
			}
		}
	}()

	buffer := bytes.Buffer{}
	encoder := cte.NewEncoder(config)
	encoder.PrepareToEncode(&buffer)
	encoder.OnBeginDocument()
	encoder.OnVersion(version.ConciseEncodingVersion)
	for _, event := range events {
		event.Invoke(encoder)
	}
	encoder.OnEndDocument()
	result := buffer.String()
	return result[3:]
}

type testType int

const (
	testTypeCbe testType = 1 << iota
	testTypeCte
	testTypeEvents
)

func stringifyEvents(events ...test.Event) (stringified []string) {
	for _, event := range events {
		stringified = append(stringified, event.String())
	}
	return
}

var allEvents = test.Events{
	EvAB, EvACL, EvACM, EvAF16, EvAF32, EvAF64, EvAI16, EvAI32, EvAI64, EvAI8,
	EvAU, EvAU16, EvAU32, EvAU64, EvAU8, EvB, EvBAB, EvBAF16, EvBAF32, EvBAF64,
	EvBAI16, EvBAI32, EvBAI64, EvBAI8, EvBAU, EvBAU16, EvBAU32, EvBAU64, EvBAU8,
	EvBCB, EvBCT, EvBMEDIA, EvBRID, EvBS, EvCB, EvCM, EvCS, EvCT, EvE,
	EvEDGE, EvINF, EvL, EvM, EvMARK, EvMEDIA, EvN, EvNAN, EvNINF, EvNODE, EvNULL,
	EvPAD, EvREFL, EvREFR, EvRID, EvS, EvSI, EvSNAN, EvST, EvT, EvUID, EvV,
}

var (
	prefixes = map[string]test.Events{
		EvREFL.Name(): {EvMARK, EvN},
	}
	followups = map[string]test.Events{
		EvL.Name():      {EvE},
		EvM.Name():      {EvE},
		EvSI.Name():     {EvS, EvE},
		EvST.Name():     {EvS, EvE, EvN},
		EvNODE.Name():   {EvN, EvE},
		EvEDGE.Name():   {EvRID, EvRID, EvN, EvE},
		EvBAB.Name():    {EvACL, EvADB},
		EvBAF16.Name():  {EvACL, EvADF16},
		EvBAF32.Name():  {EvACL, EvADF32},
		EvBAF64.Name():  {EvACL, EvADF64},
		EvBAI16.Name():  {EvACL, EvADI16},
		EvBAI32.Name():  {EvACL, EvADI32},
		EvBAI64.Name():  {EvACL, EvADI64},
		EvBAI8.Name():   {EvACL, EvADI8},
		EvBAU16.Name():  {EvACL, EvADU16},
		EvBAU32.Name():  {EvACL, EvADU32},
		EvBAU64.Name():  {EvACL, EvADU64},
		EvBAU8.Name():   {EvACL, EvADU8},
		EvBAU.Name():    {EvACL, EvADU},
		EvBCB.Name():    {EvACL, EvADU8},
		EvBCT.Name():    {EvACL, EvADT},
		EvBMEDIA.Name(): {EvACL, EvADU8},
		EvBRID.Name():   {EvACL, EvADT},
		EvBS.Name():     {EvACL, EvADT},
		EvCM.Name():     {EvN},
		EvCS.Name():     {EvN},
		EvMARK.Name():   {EvN},
		EvPAD.Name():    {EvN},
	}

	noFromCTE = map[string]bool{
		// CTE cannot produce chunked arrays
		EvACL.Name():    true,
		EvACM.Name():    true,
		EvBAB.Name():    true,
		EvBAF16.Name():  true,
		EvBAF32.Name():  true,
		EvBAF64.Name():  true,
		EvBAI16.Name():  true,
		EvBAI32.Name():  true,
		EvBAI64.Name():  true,
		EvBAI8.Name():   true,
		EvBAU16.Name():  true,
		EvBAU32.Name():  true,
		EvBAU64.Name():  true,
		EvBAU8.Name():   true,
		EvBAU.Name():    true,
		EvBCB.Name():    true,
		EvBCT.Name():    true,
		EvBMEDIA.Name(): true,
		EvBRID.Name():   true,
		EvBS.Name():     true,

		// CTE doesn't have this
		EvPAD.Name(): true,
	}

	noFromCBE = map[string]bool{
		// CBE cannot produce non-chunked versions of these
		EvAB.Name():    true,
		EvAU8.Name():   true,
		EvCB.Name():    true,
		EvMEDIA.Name(): true,
		EvREFR.Name():  true,
		EvRID.Name():   true,

		// CBE doesn't have these
		EvCM.Name():  true,
		EvCS.Name():  true,
		EvBCT.Name(): true,
		EvCT.Name():  true,
	}

	// Short form arrays are non-chunked
	arrayHasShortForm = map[string]bool{
		EvS.Name():    true,
		EvAU.Name():   true,
		EvAI8.Name():  true,
		EvAI16.Name(): true,
		EvAI32.Name(): true,
		EvAI64.Name(): true,
		EvAU16.Name(): true,
		EvAU32.Name(): true,
		EvAU64.Name(): true,
		EvAF16.Name(): true,
		EvAF32.Name(): true,
		EvAF64.Name(): true,
	}

	noToCBE = map[string]bool{
		// CBE errors out if this is attempted
		EvBCT.Name(): true,
		EvCT.Name():  true,
	}
)

func hasNonstandardCTEEncoding(formats configuration.CTEEncoderDefaultNumericFormats, event test.Event) bool {
	switch event.Name() {
	case EvAI8.Name(), EvBAI8.Name():
		return formats.Array.Int8 != configuration.CTEEncodingFormatDecimal
	case EvAI16.Name(), EvBAI16.Name():
		return formats.Array.Int16 != configuration.CTEEncodingFormatDecimal
	case EvAI32.Name(), EvBAI32.Name():
		return formats.Array.Int32 != configuration.CTEEncodingFormatDecimal
	case EvAI64.Name(), EvBAI64.Name():
		return formats.Array.Int64 != configuration.CTEEncodingFormatDecimal
	case EvAU8.Name(), EvBAU8.Name():
		return formats.Array.Uint8 != configuration.CTEEncodingFormatDecimal
	case EvAU16.Name(), EvBAU16.Name():
		return formats.Array.Uint16 != configuration.CTEEncodingFormatDecimal
	case EvAU32.Name(), EvBAU32.Name():
		return formats.Array.Uint32 != configuration.CTEEncodingFormatDecimal
	case EvAU64.Name(), EvBAU64.Name():
		return formats.Array.Uint64 != configuration.CTEEncodingFormatDecimal
	case EvAF16.Name(), EvBAF16.Name():
		return formats.Array.Float16 != configuration.CTEEncodingFormatHexadecimal
	case EvAF32.Name(), EvBAF32.Name():
		return formats.Array.Float32 != configuration.CTEEncodingFormatHexadecimal
	case EvAF64.Name(), EvBAF64.Name():
		return formats.Array.Float64 != configuration.CTEEncodingFormatHexadecimal
	}

	switch event.Value().(type) {
	case int, int8, int16, int32, int64, *big.Int, test.NegFFFFFFFFFFFFFFFF:
		return formats.Int != configuration.CTEEncodingFormatDecimal
	case uint, uint8, uint16, uint32, uint64:
		return formats.Uint != configuration.CTEEncodingFormatDecimal
	case float32, float64, *big.Float:
		return formats.BinaryFloat != configuration.CTEEncodingFormatHexadecimal
	}

	return false
}

func canConvertFromCTE(config *configuration.CTEEncoderConfiguration, events ...test.Event) bool {
	for _, event := range events {
		if noFromCTE[event.Name()] {
			return false
		}
	}
	return true
}

func canConvertToCTE(config *configuration.CTEEncoderConfiguration, events ...test.Event) bool {
	for _, event := range events {
		if hasNonstandardCTEEncoding(config.DefaultNumericFormats, event) {
			return false
		}
	}
	return true
}

func canConvertFromCBE(config *configuration.CTEEncoderConfiguration, events ...test.Event) bool {
	for _, event := range events {
		if noFromCBE[event.Name()] {
			return false
		}
	}
	return true
}

func canConvertToCBE(config *configuration.CTEEncoderConfiguration, events ...test.Event) bool {
	for _, event := range events {
		if arrayHasShortForm[event.Name()] && event.ArrayElementCount() > 15 {
			return false
		}
		if noToCBE[event.Name()] {
			return false
		}
	}
	return true
}

func mustNotConvertToCBE(events ...test.Event) bool {
	for _, event := range events {
		if noToCBE[event.Name()] {
			return true
		}
	}
	return false
}

func generateEventPrefixesAndFollowups(events ...test.Event) (eventSets []test.Events) {
	for _, event := range events {
		eventSet := []test.Event{}
		if pre, ok := prefixes[event.Name()]; ok {
			eventSet = append(eventSet, pre...)
		}
		eventSet = append(eventSet, event)
		if post, ok := followups[event.Name()]; ok {
			eventSet = append(eventSet, post...)
		}
		eventSets = append(eventSets, eventSet)
	}
	return
}

func complementaryEvents(events test.Events) test.Events {
	complementary := make(test.Events, 0, len(allEvents)/2)
	for _, event := range allEvents {
		for _, compareEvent := range events {
			if event == compareEvent {
				goto Skip
			}
		}
		complementary = append(complementary, event)
	Skip:
	}
	return complementary
}

var (
	EvAB     = test.AB([]bool{true})
	EvACL    = test.ACL(1)
	EvACM    = test.ACM(1)
	EvADB    = test.ADB([]bool{true})
	EvADF16  = test.ADF16([]float32{1})
	EvADF32  = test.ADF32([]float32{1})
	EvADF64  = test.ADF64([]float64{1})
	EvADI16  = test.ADI16([]int16{1})
	EvADI32  = test.ADI32([]int32{1})
	EvADI64  = test.ADI64([]int64{1})
	EvADI8   = test.ADI8([]int8{1})
	EvADT    = test.ADT("a")
	EvADU    = test.ADU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvADU16  = test.ADU16([]uint16{1})
	EvADU32  = test.ADU32([]uint32{1})
	EvADU64  = test.ADU64([]uint64{1})
	EvADU8   = test.ADU8([]uint8{1})
	EvAF16   = test.AF16([]float32{1})
	EvAF32   = test.AF32([]float32{1})
	EvAF64   = test.AF64([]float64{1})
	EvAI16   = test.AI16([]int16{1})
	EvAI32   = test.AI32([]int32{1})
	EvAI64   = test.AI64([]int64{1})
	EvAI8    = test.AI8([]int8{1})
	EvAU     = test.AU([][]byte{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}})
	EvAU16   = test.AU16([]uint16{1})
	EvAU32   = test.AU32([]uint32{1})
	EvAU64   = test.AU64([]uint64{1})
	EvAU8    = test.AU8([]uint8{1})
	EvB      = test.B(true)
	EvBAB    = test.BAB()
	EvBAF16  = test.BAF16()
	EvBAF32  = test.BAF32()
	EvBAF64  = test.BAF64()
	EvBAI16  = test.BAI16()
	EvBAI32  = test.BAI32()
	EvBAI64  = test.BAI64()
	EvBAI8   = test.BAI8()
	EvBAU    = test.BAU()
	EvBAU16  = test.BAU16()
	EvBAU32  = test.BAU32()
	EvBAU64  = test.BAU64()
	EvBAU8   = test.BAU8()
	EvBCB    = test.BCB(0)
	EvBCT    = test.BCT(0)
	EvBMEDIA = test.BMEDIA("a/b")
	EvBRID   = test.BRID()
	EvBS     = test.BS()
	EvCB     = test.CB(0, []byte{1})
	EvCM     = test.CM("a")
	EvCS     = test.CS("a")
	EvCT     = test.CT(0, "a")
	EvE      = test.E()
	EvEDGE   = test.EDGE()
	EvINF    = test.N(math.Inf(1))
	EvL      = test.L()
	EvM      = test.M()
	EvMARK   = test.MARK("a")
	EvMEDIA  = test.MEDIA("a/b", []byte{1})
	EvN      = test.N(-1)
	EvNAN    = test.N(compact_float.QuietNaN())
	EvNINF   = test.N(math.Inf(-1))
	EvNODE   = test.NODE()
	EvNULL   = test.NULL()
	EvPAD    = test.PAD()
	EvREFL   = test.REFL("a")
	EvREFR   = test.REFR("a")
	EvRID    = test.RID("http://z.com")
	EvS      = test.S("a")
	EvSI     = test.SI("a")
	EvSNAN   = test.N(compact_float.SignalingNaN())
	EvST     = test.ST("a")
	EvT      = test.T(compact_time.AsCompactTime(time.Date(2020, time.Month(1), 1, 1, 1, 1, 1, time.UTC)))
	EvUID    = test.UID([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	EvV      = test.V(version.ConciseEncodingVersion)
)

func AB(v []bool) test.Event              { return test.AB(v) }
func ACL(l uint64) test.Event             { return test.ACL(l) }
func ACM(l uint64) test.Event             { return test.ACM(l) }
func ADB(v []bool) test.Event             { return test.ADB(v) }
func ADF16(v []float32) test.Event        { return test.ADF16(v) }
func ADF32(v []float32) test.Event        { return test.ADF32(v) }
func ADF64(v []float64) test.Event        { return test.ADF64(v) }
func ADI16(v []int16) test.Event          { return test.ADI16(v) }
func ADI32(v []int32) test.Event          { return test.ADI32(v) }
func ADI64(v []int64) test.Event          { return test.ADI64(v) }
func ADI8(v []int8) test.Event            { return test.ADI8(v) }
func ADT(v string) test.Event             { return test.ADT(v) }
func ADU(v [][]byte) test.Event           { return test.ADU(v) }
func ADU16(v []uint16) test.Event         { return test.ADU16(v) }
func ADU32(v []uint32) test.Event         { return test.ADU32(v) }
func ADU64(v []uint64) test.Event         { return test.ADU64(v) }
func ADU8(v []uint8) test.Event           { return test.ADU8(v) }
func AF16(v []float32) test.Event         { return test.AF16(v) }
func AF32(v []float32) test.Event         { return test.AF32(v) }
func AF64(v []float64) test.Event         { return test.AF64(v) }
func AI16(v []int16) test.Event           { return test.AI16(v) }
func AI32(v []int32) test.Event           { return test.AI32(v) }
func AI64(v []int64) test.Event           { return test.AI64(v) }
func AI8(v []int8) test.Event             { return test.AI8(v) }
func AU(v [][]byte) test.Event            { return test.AU(v) }
func AU16(v []uint16) test.Event          { return test.AU16(v) }
func AU32(v []uint32) test.Event          { return test.AU32(v) }
func AU64(v []uint64) test.Event          { return test.AU64(v) }
func AU8(v []byte) test.Event             { return test.AU8(v) }
func B(v bool) test.Event                 { return test.B(v) }
func BAB() test.Event                     { return test.BAB() }
func BAF16() test.Event                   { return test.BAF16() }
func BAF32() test.Event                   { return test.BAF32() }
func BAF64() test.Event                   { return test.BAF64() }
func BAI16() test.Event                   { return test.BAI16() }
func BAI32() test.Event                   { return test.BAI32() }
func BAI64() test.Event                   { return test.BAI64() }
func BAI8() test.Event                    { return test.BAI8() }
func BAU() test.Event                     { return test.BAU() }
func BAU16() test.Event                   { return test.BAU16() }
func BAU32() test.Event                   { return test.BAU32() }
func BAU64() test.Event                   { return test.BAU64() }
func BAU8() test.Event                    { return test.BAU8() }
func BCB(v uint64) test.Event             { return test.BCB(v) }
func BCT(v uint64) test.Event             { return test.BCT(v) }
func BMEDIA(v string) test.Event          { return test.BMEDIA(v) }
func BREFR() test.Event                   { return test.BREFR() }
func BRID() test.Event                    { return test.BRID() }
func BS() test.Event                      { return test.BS() }
func CB(t uint64, v []byte) test.Event    { return test.CB(t, v) }
func CM(v string) test.Event              { return test.CM(v) }
func CS(v string) test.Event              { return test.CS(v) }
func CT(t uint64, v string) test.Event    { return test.CT(t, v) }
func E() test.Event                       { return test.E() }
func EDGE() test.Event                    { return test.EDGE() }
func L() test.Event                       { return test.L() }
func M() test.Event                       { return test.M() }
func MARK(id string) test.Event           { return test.MARK(id) }
func MEDIA(t string, v []byte) test.Event { return test.MEDIA(t, v) }
func N(v interface{}) test.Event          { return test.N(v) }
func NAN() test.Event                     { return test.NAN() }
func NODE() test.Event                    { return test.NODE() }
func NULL() test.Event                    { return test.NULL() }
func PAD() test.Event                     { return test.PAD() }
func REFL(id string) test.Event           { return test.REFL(id) }
func REFR(v string) test.Event            { return test.REFR(v) }
func RID(v string) test.Event             { return test.RID(v) }
func S(v string) test.Event               { return test.S(v) }
func SI(id string) test.Event             { return test.SI(id) }
func SNAN() test.Event                    { return test.SNAN() }
func ST(id string) test.Event             { return test.ST(id) }
func T(v compact_time.Time) test.Event    { return test.T(v) }
func UID(v []byte) test.Event             { return test.UID(v) }
func V(v uint64) test.Event               { return test.V(v) }
