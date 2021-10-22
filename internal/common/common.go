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

package common

import (
	"math"
	"math/big"
	"net/url"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/kstenerud/go-concise-encoding/types"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

var (
	TypeNone             = reflect.TypeOf(nil)
	TypeInterface        = reflect.TypeOf([]interface{}{}).Elem()
	TypeInterfaceArray   = reflect.TypeOf([1]interface{}{})
	TypeInterfaceSlice   = reflect.TypeOf([]interface{}{})
	TypeInterfaceMap     = reflect.TypeOf(map[interface{}]interface{}{})
	TypeInterfaceEdge    = reflect.TypeOf(types.Edge{})
	TypeInterfaceNode    = reflect.TypeOf(types.Node{})
	TypeString           = reflect.TypeOf("")
	TypeBytes            = reflect.TypeOf([]uint8{})
	TypeTime             = reflect.TypeOf(time.Time{})
	TypePTime            = reflect.TypeOf((*time.Time)(nil))
	TypeCompactTime      = reflect.TypeOf(compact_time.Time{})
	TypePCompactTime     = reflect.TypeOf((*compact_time.Time)(nil))
	TypeDFloat           = reflect.TypeOf(compact_float.DFloat{})
	TypeBigInt           = reflect.TypeOf(big.Int{})
	TypePBigInt          = reflect.TypeOf((*big.Int)(nil))
	TypeBigFloat         = reflect.TypeOf(big.Float{})
	TypePBigFloat        = reflect.TypeOf((*big.Float)(nil))
	TypeBigDecimalFloat  = reflect.TypeOf(apd.Decimal{})
	TypePBigDecimalFloat = reflect.TypeOf((*apd.Decimal)(nil))
	TypeURL              = reflect.TypeOf(url.URL{})
	TypePURL             = reflect.TypeOf((*url.URL)(nil))
	TypeUID              = reflect.TypeOf(types.UID{})
	TypeMedia            = reflect.TypeOf(types.Media{})
	TypeMarkup           = reflect.TypeOf(types.Markup{})
	TypeEdge             = reflect.TypeOf(types.Edge{})
	TypeNode             = reflect.TypeOf(types.Node{})
)

var KeyableTypes = []reflect.Type{
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
	reflect.TypeOf((*compact_time.Time)(nil)).Elem(),
	reflect.TypeOf((*compact_float.DFloat)(nil)).Elem(),
	reflect.TypeOf((*interface{})(nil)).Elem(),

	// Must be pointers
	reflect.TypeOf((*big.Float)(nil)),
	reflect.TypeOf((*big.Int)(nil)),
}

var NonKeyableTypes = []reflect.Type{
	reflect.TypeOf((*big.Float)(nil)).Elem(),
	reflect.TypeOf((*big.Int)(nil)).Elem(),
	reflect.TypeOf((*apd.Decimal)(nil)).Elem(),
}

func IsFieldExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

const QuietNanBit = uint64(1 << 50)

var SignalingNan = math.Float64frombits(math.Float64bits(math.NaN()) & ^QuietNanBit)
var QuietNan = math.Float64frombits(math.Float64bits(math.NaN()) | QuietNanBit)

func IsSignalingNan(value float64) bool {
	return math.Float64bits(value)&QuietNanBit == 0
}

var Bfloat16SignalingNanBytes = []byte{0x81, 0xff}
var Bfloat16QuietNanBytes = []byte{0xc1, 0xff}

var Bfloat16SignalingNanBits = uint16(0xff81)
var Bfloat16QuietNanBits = uint16(0xffc1)
var Float32SignalingNanBits = math.Float32bits(float32(SignalingNan))
var Float32QuietNanBits = math.Float32bits(float32(math.NaN()))
var Float64SignalingNanBits = math.Float64bits(SignalingNan)
var Float64QuietNanBits = math.Float64bits(math.NaN())

var BigInt0 = big.NewInt(0)
var BigInt2 = big.NewInt(2)
var BigInt8 = big.NewInt(8)
var BigInt10 = big.NewInt(10)
var BigInt16 = big.NewInt(16)
var BigIntN1 = big.NewInt(-1)

func IsBigIntNegative(value *big.Int) bool {
	return value.Cmp(BigInt0) < 0
}

type KindProperty byte

const (
	KindPropertyPointer KindProperty = 1 << iota
	KindPropertyNullable
	KindPropertyLengthable
)

var kindProperties = [64]KindProperty{
	reflect.Chan:          KindPropertyPointer | KindPropertyNullable | KindPropertyLengthable,
	reflect.Func:          KindPropertyPointer | KindPropertyNullable,
	reflect.Interface:     KindPropertyNullable,
	reflect.Map:           KindPropertyPointer | KindPropertyNullable | KindPropertyLengthable,
	reflect.Ptr:           KindPropertyPointer | KindPropertyNullable,
	reflect.Slice:         KindPropertyPointer | KindPropertyNullable | KindPropertyLengthable,
	reflect.String:        KindPropertyLengthable,
	reflect.UnsafePointer: KindPropertyPointer,
}

func IsPointer(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyPointer != 0
}

func IsLengthable(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyLengthable != 0
}

func IsNullable(v reflect.Value) bool {
	return kindProperties[v.Kind()]&KindPropertyNullable != 0
}

func IsNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	return IsNullable(v) && v.IsNil()
}

func CloneBytes(bytes []byte) []byte {
	bytesCopy := make([]byte, len(bytes), len(bytes))
	copy(bytesCopy, bytes)
	return bytesCopy
}

const oneIf64Bit = ((uint64(^uintptr(0)) >> 32) & 1)

// The address space on this machine. This is a conservative value based on:
// * cmd/compile/internal/amd64/galign.go:  arch.MAXWIDTH = 1 << 50
// * cmd/compile/internal/mips/galign.go:   arch.MAXWIDTH = (1 << 31) - 1
const AddressSpace = int((((1 << (oneIf64Bit * 50)) - 1) * oneIf64Bit) + (((^uint64(^uintptr(0))) >> 32) & ((1 << 31) - 1)))
const BytesPerInt = int(oneIf64Bit*4 + 4)

func Is64BitArch() bool {
	return oneIf64Bit == 1
}

var requiresLowercaseAdjust [256]bool

func init() {
	for i := 'A'; i <= 'Z'; i++ {
		requiresLowercaseAdjust[i] = true
	}
}

// Convert ASCII characters A-Z to a-z, ignoring locale.
func ASCIIBytesToLower(bytes []byte) (didChange bool) {
	const lowercaseAdjust = byte('a' - 'A')

	for i, b := range bytes {
		if requiresLowercaseAdjust[b] {
			bytes[i] += lowercaseAdjust
			didChange = true
		}
	}
	return
}

func ASCIIToLower(s string) string {
	asBytes := []byte(s)
	if ASCIIBytesToLower(asBytes) {
		return string(asBytes)
	}
	return s
}

func ByteCountToElementCount(elementBitWidth int, byteCount uint64) uint64 {
	return (byteCount * 8) / uint64(elementBitWidth)
}

func ElementCountToByteCount(elementBitWidth int, elementCount uint64) uint64 {
	byteCount := (elementCount * uint64(elementBitWidth)) / 8
	if elementBitWidth == 1 && elementCount&7 != 0 {
		byteCount++
	}
	return byteCount
}

var MaxDayByMonth = []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
