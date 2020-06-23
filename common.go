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

package concise_encoding

import (
	"math"
	"math/big"
	"net/url"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// Which version of concise encoding this library adheres to
const ConciseEncodingVersion = 1

// The type of top-level container to assume is already opened (for implied
// structure documents). For normal operation, use TLContainerTypeNone.
type TLContainerType int

const (
	// Assume that no top-level container has already been opened when beginning decoding.
	TLContainerTypeNone = iota
	// Assume a list has already been opened when beginning decoding.
	TLContainerTypeList
	// Assume a map has already been opened when beginning decoding.
	TLContainerTypeMap
)

// Settings to help with debugging
type DebugOptionsStruct struct {
	// Setting this to true will cause all panics to bubble up rather than
	// being handled in the library.
	PassThroughPanics bool
}

// Only set these options when debugging a problem with the library.
var DebugOptions DebugOptionsStruct

var (
	typeInterface    = reflect.TypeOf([]interface{}{}).Elem()
	typeString       = reflect.TypeOf("")
	typeBytes        = reflect.TypeOf([]uint8{})
	typeTime         = reflect.TypeOf(time.Time{})
	typeCompactTime  = reflect.TypeOf(compact_time.Time{})
	typePCompactTime = reflect.TypeOf((*compact_time.Time)(nil))
	typeDFloat       = reflect.TypeOf(compact_float.DFloat{})

	typeBigInt  = reflect.TypeOf(big.Int{})
	typePBigInt = reflect.TypeOf((*big.Int)(nil))

	typeBigFloat  = reflect.TypeOf(big.Float{})
	typePBigFloat = reflect.TypeOf((*big.Float)(nil))

	typeBigDecimalFloat  = reflect.TypeOf(apd.Decimal{})
	typePBigDecimalFloat = reflect.TypeOf((*apd.Decimal)(nil))

	typeURL  = reflect.TypeOf(url.URL{})
	typePURL = reflect.TypeOf((*url.URL)(nil))
)

var keyableTypes = []reflect.Type{
	reflect.TypeOf((*bool)(nil)).Elem(),
	reflect.TypeOf((*int)(nil)).Elem(),
	reflect.TypeOf((*int8)(nil)).Elem(),
	reflect.TypeOf((*int16)(nil)).Elem(),
	reflect.TypeOf((*int32)(nil)).Elem(),
	reflect.TypeOf((*int64)(nil)).Elem(),
	reflect.TypeOf((*uint)(nil)).Elem(),
	reflect.TypeOf((*uint8)(nil)).Elem(),
	reflect.TypeOf((*uint16)(nil)).Elem(),
	reflect.TypeOf((*uint32)(nil)).Elem(),
	reflect.TypeOf((*uint64)(nil)).Elem(),
	reflect.TypeOf((*float32)(nil)).Elem(),
	reflect.TypeOf((*float64)(nil)).Elem(),
	reflect.TypeOf((*string)(nil)).Elem(),
	reflect.TypeOf((*url.URL)(nil)).Elem(),
	reflect.TypeOf((*time.Time)(nil)).Elem(),
	reflect.TypeOf((*compact_float.DFloat)(nil)).Elem(),
	reflect.TypeOf((*interface{})(nil)).Elem(),
}

var nonKeyableTypes = []reflect.Type{
	reflect.TypeOf((*big.Int)(nil)).Elem(),
	reflect.TypeOf((*apd.Decimal)(nil)).Elem(),
}

func isFieldExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

const quietNanBit = uint64(1 << 50)

var signalingNan = math.Float64frombits(math.Float64bits(math.NaN()) & ^quietNanBit)
var quietNan = math.Float64frombits(math.Float64bits(math.NaN()) | quietNanBit)

func isSignalingNan(value float64) bool {
	return math.Float64bits(value)&quietNanBit == 0
}

var bigInt0 = big.NewInt(0)
var bigInt10 = big.NewInt(10)
var bigIntN1 = big.NewInt(-1)

func isBigIntNegative(value *big.Int) bool {
	return value.Cmp(bigInt0) < 0
}
