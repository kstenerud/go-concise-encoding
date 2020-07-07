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

package builder

import (
	"math/big"
	"reflect"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/kstenerud/go-concise-encoding/conversions"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

// Int

func setIntFromInt(value int64, dst reflect.Value) {
	dst.SetInt(value)
	if dst.Int() != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromUint(value uint64, dst reflect.Value) {
	i, err := conversions.UintToInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)

	}
	dst.SetInt(i)
	if uint64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromBigInt(value *big.Int, dst reflect.Value) {
	i, err := conversions.BigIntToInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)

	}
	dst.SetInt(i)
	if dst.Int() != i {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromFloat(value float64, dst reflect.Value) {
	dst.SetInt(int64(value))
	if float64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromBigFloat(value *big.Float, dst reflect.Value) {
	i, err := conversions.BigFloatToInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetInt(i)
	if dst.Int() != i {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	i, err := value.Int()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetInt(i)
	if dst.Int() != i {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	i, err := value.Int64()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetInt(i)
	if dst.Int() != i {
		builderPanicCannotConvert(value, dst.Type())
	}
}

// UInt

func setUintFromInt(value int64, dst reflect.Value) {
	u, err := conversions.IntToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	setUintFromUint(u, dst)
}

func setUintFromUint(value uint64, dst reflect.Value) {
	dst.SetUint(value)
	if dst.Uint() != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromBigInt(value *big.Int, dst reflect.Value) {
	u, err := conversions.BigIntToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	setUintFromUint(u, dst)
}

func setUintFromFloat(value float64, dst reflect.Value) {
	u := uint64(value)
	if float64(u) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
	setUintFromUint(u, dst)
}

func setUintFromBigFloat(value *big.Float, dst reflect.Value) {
	u, err := conversions.BigFloatToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	setUintFromUint(u, dst)
}

func setUintFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	u, err := value.Uint()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	setUintFromUint(u, dst)
}

func setUintFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	u, err := conversions.BigDecimalFloatToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	setUintFromUint(u, dst)
}

// Float

// TODO: When to allow lossy conversions?

func setFloatFromInt(value int64, dst reflect.Value) {
	dst.SetFloat(float64(value))
	if int64(dst.Float()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setFloatFromUint(value uint64, dst reflect.Value) {
	dst.SetFloat(float64(value))
	if uint64(dst.Float()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setFloatFromFloat(value float64, dst reflect.Value) {
	dst.SetFloat(value)
}

func setFloatFromBigInt(value *big.Int, dst reflect.Value) {
	v, err := conversions.BigIntToFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetFloat(v)
}

func setFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	f, err := conversions.BigFloatToFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetFloat(f)
}

func setFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.SetFloat(value.Float())
}

func setFloatFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	v, err := value.Float64()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetFloat(v)
}

// BigInt

func setBigIntFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*big.NewInt(value)))
}

func setBigIntFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*conversions.UintToBigInt(value)))
}

func setBigIntFromFloat(value float64, dst reflect.Value) {
	bi, err := conversions.FloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromBigFloat(value *big.Float, dst reflect.Value) {
	bi, err := conversions.BigFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := conversions.DecimalFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	bi, err := conversions.BigDecimalFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

// pBigInt

func setPBigIntFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(big.NewInt(value)))
}

func setPBigIntFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(conversions.UintToBigInt(value)))
}

func setPBigIntFromFloat(value float64, dst reflect.Value) {
	bi, err := conversions.FloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromBigFloat(value *big.Float, dst reflect.Value) {
	bi, err := conversions.BigFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := conversions.DecimalFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	bi, err := conversions.BigDecimalFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

// BigFloat

func setBigFloatFromInt(value int64, dst reflect.Value) {
	bf := new(big.Float)
	bf.SetInt64(value)
	dst.Set(reflect.ValueOf(*bf))
}

func setBigFloatFromUint(value uint64, dst reflect.Value) {
	bf := new(big.Float)
	bf.SetUint64(value)
	dst.Set(reflect.ValueOf(*bf))
}

func setBigFloatFromFloat(value float64, dst reflect.Value) {
	bf := big.NewFloat(value)
	dst.Set(reflect.ValueOf(*bf))
}

func setBigFloatFromBigInt(value *big.Int, dst reflect.Value) {
	bf := new(big.Float)
	bf.SetInt(value)
	dst.Set(reflect.ValueOf(*bf))
}

func setBigFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bf := value.BigFloat()
	dst.Set(reflect.ValueOf(*bf))
}

func setBigFloatFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	bf, err := conversions.BigDecimalFloatToBigFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bf))
}

// pBigFloat

func setPBigFloatFromInt(value int64, dst reflect.Value) {
	f := new(big.Float)
	f.SetInt64(value)
	dst.Set(reflect.ValueOf(f))
}

func setPBigFloatFromUint(value uint64, dst reflect.Value) {
	bf := new(big.Float)
	bf.SetUint64(value)
	dst.Set(reflect.ValueOf(bf))
}

func setPBigFloatFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(big.NewFloat(value)))
}

func setPBigFloatFromBigInt(value *big.Int, dst reflect.Value) {
	bf := new(big.Float)
	bf.SetInt(value)
	dst.Set(reflect.ValueOf(bf))
}

func setPBigFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(value.BigFloat()))
}

func setPBigFloatFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	bf, err := conversions.BigDecimalFloatToBigFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bf))
}

// BigDecimalFloat

func setBigDecimalFloatFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(conversions.IntToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(conversions.UintToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(conversions.BigIntToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromFloat(value float64, dst reflect.Value) {
	bdf, err := conversions.FloatToBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bdf))
}

func setBigDecimalFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	pbdf, err := conversions.BigFloatToPBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*pbdf))
}

func setBigDecimalFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value.APD()))
}

// PBigDecimalFloat

func setPBigDecimalFloatFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(big.NewInt(value), 0)))
}

func setPBigDecimalFloatFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(conversions.UintToBigInt(value), 0)))
}

func setPBigDecimalFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(value, 0)))
}

func setPBigDecimalFloatFromFloat(value float64, dst reflect.Value) {
	v, err := conversions.FloatToPBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(v))
}

func setPBigDecimalFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	pbdf, err := conversions.BigFloatToPBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(pbdf))
}

func setPBigDecimalFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value.APD()))
}

func setUintFromAnything(src reflect.Value, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setUintFromInt(src.Int(), dst)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setUintFromUint(src.Uint(), dst)
		return
	case reflect.Float32, reflect.Float64:
		setUintFromFloat(src.Float(), dst)
		return
	case reflect.Interface:
		setUintFromAnything(src.Elem(), dst)
		return
	case reflect.Struct:
		switch src.Type() {
		case common.TypeDFloat:
			setUintFromDecimalFloat(src.Interface().(compact_float.DFloat), dst)
			return
		}
	case reflect.Ptr:
		switch src.Type() {
		case common.TypePBigInt:
			setUintFromBigInt(src.Interface().(*big.Int), dst)
			return
		case common.TypePBigFloat:
			setUintFromBigFloat(src.Interface().(*big.Float), dst)
			return
		case common.TypePBigDecimalFloat:
			setUintFromBigDecimalFloat(src.Interface().(*apd.Decimal), dst)
			return
		}
		setUintFromAnything(src.Elem(), dst)
		return
	}
	builderPanicCannotConvertRV(src, dst.Type())
}

func setIntFromAnything(src reflect.Value, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setIntFromInt(src.Int(), dst)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setIntFromUint(src.Uint(), dst)
		return
	case reflect.Float32, reflect.Float64:
		setIntFromFloat(src.Float(), dst)
		return
	case reflect.Interface:
		setIntFromAnything(src.Elem(), dst)
		return
	case reflect.Struct:
		switch src.Type() {
		case common.TypeDFloat:
			setIntFromDecimalFloat(src.Interface().(compact_float.DFloat), dst)
			return
		}
	case reflect.Ptr:
		switch src.Type() {
		case common.TypePBigInt:
			setIntFromBigInt(src.Interface().(*big.Int), dst)
			return
		case common.TypePBigFloat:
			setIntFromBigFloat(src.Interface().(*big.Float), dst)
			return
		case common.TypePBigDecimalFloat:
			setIntFromBigDecimalFloat(src.Interface().(*apd.Decimal), dst)
			return
		}
		setIntFromAnything(src.Elem(), dst)
		return
	}
	builderPanicCannotConvertRV(src, dst.Type())
}

func setFloatFromAnything(src reflect.Value, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setFloatFromInt(src.Int(), dst)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setFloatFromUint(src.Uint(), dst)
		return
	case reflect.Float32, reflect.Float64:
		setFloatFromFloat(src.Float(), dst)
		return
	case reflect.Interface:
		setFloatFromAnything(src.Elem(), dst)
		return
	case reflect.Struct:
		switch src.Type() {
		case common.TypeDFloat:
			setFloatFromDecimalFloat(src.Interface().(compact_float.DFloat), dst)
			return
		}
	case reflect.Ptr:
		switch src.Type() {
		case common.TypePBigInt:
			setFloatFromBigInt(src.Interface().(*big.Int), dst)
			return
		case common.TypePBigFloat:
			setFloatFromBigFloat(src.Interface().(*big.Float), dst)
			return
		case common.TypePBigDecimalFloat:
			setFloatFromBigDecimalFloat(src.Interface().(*apd.Decimal), dst)
			return
		}
		setFloatFromAnything(src.Elem(), dst)
		return
	}
	builderPanicCannotConvert(src, dst.Type())
}

func setAnythingFromAnything(src reflect.Value, dst reflect.Value) {
	if src.Type() == dst.Type() {
		dst.Set(src)
		return
	}

	switch dst.Kind() {
	case reflect.Bool:
		if src.Kind() == reflect.Bool {
			dst.SetBool(src.Bool())
			return
		}
	case reflect.String:
		if src.Kind() == reflect.String {
			dst.SetString(src.String())
			return
		}
	case reflect.Interface:
		dst.Set(src)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setUintFromAnything(src, dst)
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setIntFromAnything(src, dst)
		return
	case reflect.Float32, reflect.Float64:
		setFloatFromAnything(src, dst)
		return
	case reflect.Array:
		panic("TODO: setAnythingFromAnything: Array")
	case reflect.Slice:
		panic("TODO: setAnythingFromAnything: Slice")
	case reflect.Map:
		panic("TODO: setAnythingFromAnything: Map")
	case reflect.Struct:
		panic("TODO: setAnythingFromAnything: Struct")
	case reflect.Ptr:
		panic("TODO: setAnythingFromAnything: Ptr")
	}
	builderPanicCannotConvert(src, dst.Type())
}
