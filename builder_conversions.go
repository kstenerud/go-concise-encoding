package concise_encoding

import (
	"math/big"
	"reflect"

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
	i, err := uintToInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)

	}
	dst.SetInt(i)
	if uint64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromBigInt(value *big.Int, dst reflect.Value) {
	i, err := bigIntToInt(value)
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
	i, accuracy := value.Int64()
	if accuracy != big.Exact {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetInt(i)
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
	u, err := intToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetUint(u)
	if int64(dst.Uint()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromUint(value uint64, dst reflect.Value) {
	dst.SetUint(value)
	if dst.Uint() != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromBigInt(value *big.Int, dst reflect.Value) {
	u, err := bigIntToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetUint(u)
	if dst.Uint() != u {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromFloat(value float64, dst reflect.Value) {
	u, err := floatToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetUint(u)
	if float64(dst.Uint()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromBigFloat(value *big.Float, dst reflect.Value) {
	u, accuracy := value.Uint64()
	if accuracy != big.Exact {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(u)
}

func setUintFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	u, err := value.Uint()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetUint(u)
	if dst.Uint() != u {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	u, err := bigDecimalFloatToUint(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetUint(u)
	if dst.Uint() != u {
		builderPanicCannotConvert(value, dst.Type())
	}
}

// Float

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

func setFloatFromBigInt(value *big.Int, dst reflect.Value) {
	v, err := bigIntToFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetFloat(v)
}

func setFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	f, err := bigFloatToFloat(value)
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
	dst.Set(reflect.ValueOf(*uintToBigInt(value)))
}

func setBigIntFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*floatToBigInt(value)))
}

func setBigIntFromBigFloat(value *big.Float, dst reflect.Value) {
	bi, err := bigFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := value.BigInt()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*bigDecimalFloatToBigInt(value)))
}

// pBigInt

func setPBigIntFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(big.NewInt(value)))
}

func setPBigIntFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(uintToBigInt(value)))
}

func setPBigIntFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(floatToBigInt(value)))
}

func setPBigIntFromBigFloat(value *big.Float, dst reflect.Value) {
	bi, err := bigFloatToBigInt(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := value.BigInt()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(bigDecimalFloatToBigInt(value)))
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
	bf, err := bigDecimalFloatToBigFloat(value)
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
	bf, err := bigDecimalFloatToBigFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bf))
}

// BigDecimalFloat

func setBigDecimalFloatFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(intToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(uintToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(bigIntToBigDecimalFloat(value)))
}

func setBigDecimalFloatFromFloat(value float64, dst reflect.Value) {
	bdf, err := floatToBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bdf))
}

func setBigDecimalFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	pbdf, err := bigFloatToPBigDecimalFloat(value)
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
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(uintToBigInt(value), 0)))
}

func setPBigDecimalFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(value, 0)))
}

func setPBigDecimalFloatFromFloat(value float64, dst reflect.Value) {
	v, err := floatToPBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(v))
}

func setPBigDecimalFloatFromBigFloat(value *big.Float, dst reflect.Value) {
	pbdf, err := bigFloatToPBigDecimalFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(pbdf))
}

func setPBigDecimalFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value.APD()))
}
