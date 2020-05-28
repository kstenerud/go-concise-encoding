package concise_encoding

import (
	"math"
	"math/big"
	"reflect"
	"strconv"

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
	if value > 0x7fffffffffffffff {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetInt(int64(value))
	if uint64(dst.Int()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setIntFromBigInt(value *big.Int, dst reflect.Value) {
	if !value.IsInt64() {
		builderPanicCannotConvert(value, dst.Type())
	}
	i := value.Int64()
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
	if value < 0 {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(uint64(value))
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
	if !value.IsUint64() {
		builderPanicCannotConvert(value, dst.Type())
	}
	u := value.Uint64()
	dst.SetUint(u)
	if dst.Uint() != u {
		builderPanicCannotConvert(value, dst.Type())
	}
}

func setUintFromFloat(value float64, dst reflect.Value) {
	if value < 0 {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(uint64(value))
	if float64(dst.Uint()) != value {
		builderPanicCannotConvert(value, dst.Type())
	}
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
	if i, err := value.Int64(); err == nil {
		setUintFromInt(i, dst)
		return
	}

	if value.Negative {
		builderPanicCannotConvert(value, dst.Type())
	}

	if !value.Coeff.IsUint64() {
		builderPanicCannotConvert(value, dst.Type())
	}

	u := value.Coeff.Uint64()
	exp := uint64(math.Pow10(int(value.Exponent)))
	result := u * exp
	if result/exp != u {
		builderPanicCannotConvert(value, dst.Type())
	}
	dst.SetUint(result)
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
	v, err := strconv.ParseFloat(value.String(), 64)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.SetFloat(v)
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
	dst.Set(reflect.ValueOf(*uintToPBigInt(value)))
}

func setBigIntFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*floatToPBigInt(value)))
}

func setBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := value.BigInt()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(*bi))
}

func setBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*bigFloatToPBigInt(value)))
}

// pBigInt

func setPBigIntFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(big.NewInt(value)))
}

func setPBigIntFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(uintToPBigInt(value)))
}

func setPBigIntFromFloat(value float64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(floatToPBigInt(value)))
}

func setPBigIntFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	bi, err := value.BigInt()
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(bi))
}

func setPBigIntFromBigDecimalFloat(value *apd.Decimal, dst reflect.Value) {
	dst.Set(reflect.ValueOf(bigFloatToPBigInt(value)))
}

func uintToPBigInt(value uint64) *big.Int {
	if value <= 0x7fffffffffffffff {
		return big.NewInt(int64(value))
	}

	bi := big.NewInt(int64(value >> 1))
	return bi.Lsh(bi, 1)
}

func floatToPBigInt(value float64) *big.Int {
	bi, _ := big.NewFloat(value).Int(nil)
	return bi
}

func bigFloatToPBigInt(value *apd.Decimal) *big.Int {
	exp := big.NewInt(int64(value.Exponent))
	exp.Exp(bigInt10, exp, nil)
	return exp.Mul(exp, &value.Coeff)
}

// BigFloat

func setBigFloatFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(intToBigFloat(value)))
}

func setBigFloatFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(uintToBigFloat(value)))
}

func setBigFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(pBigIntToBigFloat(value)))
}

func setBigFloatFromFloat(value float64, dst reflect.Value) {
	v, err := floatToBigFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(v))
}

func setBigFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value.APD()))
}

func intToBigFloat(value int64) apd.Decimal {
	if value < 0 {
		return apd.Decimal{
			Negative: true,
			Coeff:    *big.NewInt(-value),
		}
	}
	return apd.Decimal{
		Coeff: *big.NewInt(value),
	}
}

func uintToBigFloat(value uint64) apd.Decimal {
	return apd.Decimal{
		Coeff: *uintToPBigInt(value),
	}
}

func floatToBigFloat(value float64) (apd.Decimal, error) {
	var d apd.Decimal
	_, _, err := apd.BaseContext.SetString(&d, strconv.FormatFloat(value, 'g', -1, 64))
	return d, err
}

func pBigIntToBigFloat(value *big.Int) apd.Decimal {
	return apd.Decimal{
		Coeff: *value,
	}
}

// PBigFloat

func setPBigFloatFromInt(value int64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(big.NewInt(value), 0)))
}

func setPBigFloatFromUint(value uint64, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(uintToPBigInt(value), 0)))
}

func setPBigFloatFromBigInt(value *big.Int, dst reflect.Value) {
	dst.Set(reflect.ValueOf(apd.NewWithBigInt(value, 0)))
}

func setPBigFloatFromFloat(value float64, dst reflect.Value) {
	v, err := floatToPBigFloat(value)
	if err != nil {
		builderPanicErrorConverting(value, dst.Type(), err)
	}
	dst.Set(reflect.ValueOf(v))
}

func setPBigFloatFromDecimalFloat(value compact_float.DFloat, dst reflect.Value) {
	dst.Set(reflect.ValueOf(*value.APD()))
}

func floatToPBigFloat(value float64) (*apd.Decimal, error) {
	d, _, err := apd.NewFromString(strconv.FormatFloat(value, 'g', -1, 64))
	return d, err
}
