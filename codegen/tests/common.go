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

	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/configuration"
	"github.com/kstenerud/go-concise-encoding/cte"
	"github.com/kstenerud/go-concise-encoding/test"
	"github.com/kstenerud/go-concise-encoding/version"
)

func truncate(data []byte, count int) []byte {
	return data[:len(data)-count]
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
