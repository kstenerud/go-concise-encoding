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

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/types"
)

// Numeric

const (
	Int8Min   = -0x80
	Int8Max   = 0x7f
	Int16Min  = -0x8000
	Int16Max  = 0x7fff
	Int32Min  = -0x80000000
	Int32Max  = 0x7fffffff
	Uint8Max  = 0xff
	Uint16Max = 0xffff
	Uint32Max = 0xffffffff

	Float32ExponentMin = -126
	Float32ExponentMax = 127

	Float64ExponentMin        = -1022
	Float64ExponentMax        = 1023
	Float64DecimalExponentMin = -308 // Denormalized can technically go to -324
	Float64DecimalExponentMax = 308
	Float64CoefficientMax     = (uint64(1) << 54) - 1

	Bfloat16SignMask         = uint16(0x8000)
	Bfloat16SpecialMask      = uint16(0x7f80)
	Bfloat16FractionMask     = uint16(0x007f)
	Bfloat16QuietNanBit      = uint16(0x0040)
	Bfloat16QuietNanBits     = uint16(0x7fc1)
	Bfloat16SignalingNanBits = uint16(0x7f81)

	Float32SignMask         = uint32(0x80000000)
	Float32SpecialMask      = uint32(0x7f800000)
	Float32FractionMask     = uint32(0x007fffff)
	Float32QuietNanBit      = uint32(0x00400000)
	Float32QuietNanBits     = uint32(0x7fc00001)
	Float32SignalingNanBits = uint32(0x7f800001)

	Float64SignMask         = uint64(0x8000000000000000)
	Float64SpecialMask      = uint64(0x7ff0000000000000)
	Float64FractionMask     = uint64(0x000fffffffffffff)
	Float64QuietNanBit      = uint64(0x0008000000000000)
	Float64QuietNanBits     = uint64(0x7ff8000000000001)
	Float64SignalingNanBits = uint64(0x7ff0000000000001)
)

var (
	Bfloat16SignalingNanBytes = []byte{0x81, 0x7f}
	Bfloat16QuietNanBytes     = []byte{0xc1, 0x7f}

	Float32SignalingNan = math.Float32frombits(Float32SignalingNanBits)
	Float32QuietNan     = math.Float32frombits(Float32QuietNanBits)

	Float64SignalingNan = math.Float64frombits(Float64SignalingNanBits)
	Float64QuietNan     = math.Float64frombits(Float64QuietNanBits)

	BigInt0  = big.NewInt(0)
	BigInt2  = big.NewInt(2)
	BigInt8  = big.NewInt(8)
	BigInt10 = big.NewInt(10)
	BigInt16 = big.NewInt(16)
	BigIntN1 = big.NewInt(-1)
)

// Architecture

const oneIf64Bit = ((uint64(^uintptr(0)) >> 32) & 1)

// The address space on this machine. This is a conservative value based on:
// * cmd/compile/internal/amd64/galign.go:  arch.MAXWIDTH = 1 << 50
// * cmd/compile/internal/mips/galign.go:   arch.MAXWIDTH = (1 << 31) - 1
const AddressSpace = int((((1 << (oneIf64Bit * 50)) - 1) * oneIf64Bit) + (((^uint64(^uintptr(0))) >> 32) & ((1 << 31) - 1)))
const BytesPerInt = int(oneIf64Bit*4 + 4)

// Reflect

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

// Utility

var MaxDayByMonth = []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
