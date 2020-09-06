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

// Various functions to convert between numeric types.
package conversions

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/kstenerud/go-concise-encoding/internal/common"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
)

// apd.Decimal to other

func BigDecimalFloatToBigFloat(value *apd.Decimal) (*big.Float, error) {
	return StringToBigFloat(value.Text('g'), int(value.NumDigits()))
}

func BigDecimalFloatToBigInt(value *apd.Decimal, maxBase10Exponent int) (*big.Int, error) {
	switch value.Form {
	case apd.NaN, apd.NaNSignaling, apd.Infinite:
		return nil, fmt.Errorf("%v cannot fit into a big.Int", value)
	}
	if value.Exponent < 0 {
		return nil, fmt.Errorf("%v cannot fit into a big.Int", value)
	}
	if value.Exponent > int32(maxBase10Exponent) {
		return nil, fmt.Errorf("%v has a decimal exponential component (%v) that is too large (max %v)", value, value.Exponent, maxBase10Exponent)
	}
	exp := big.NewInt(int64(value.Exponent))
	exp.Exp(common.BigInt10, exp, nil)
	return exp.Mul(exp, &value.Coeff), nil
}

func BigDecimalFloatToUint(value *apd.Decimal) (uint64, error) {
	if i, err := value.Int64(); err == nil {
		return uint64(i), nil
	}

	bf, err := BigDecimalFloatToBigFloat(value)
	if err != nil {
		return 0, err
	}
	return BigFloatToUint(bf)
}

// big.Float to other

func BigFloatToPBigDecimalFloat(value *big.Float) (*apd.Decimal, error) {
	d, _, err := apd.NewFromString(BigFloatToString(value))
	return d, err
}

func BigFloatToBigInt(value *big.Float, maxBase2Exponent int) (*big.Int, error) {
	if value.MantExp(nil) > maxBase2Exponent {
		return nil, fmt.Errorf("%v has a binary exponential component (%v) that is too large for a big int (max %v)", value, value.MantExp(nil), maxBase2Exponent)
	}
	bi, accuracy := value.Int(new(big.Int))
	if accuracy != big.Exact {
		return nil, fmt.Errorf("%v cannot fit into a big.Int", value)
	}
	return bi, nil
}

func BigFloatToFloat(value *big.Float) (float64, error) {
	exp := value.MantExp(nil)
	if exp < -1029 {
		return 0, fmt.Errorf("%v is too small to fit into a float64", value)
	}
	if exp > 1024 {
		return 0, fmt.Errorf("%v is too big to fit into a float64", value)
	}
	f, accuracy := value.Float64()
	if accuracy != big.Exact {
		if f == 0 {
			return 0, fmt.Errorf("%v is too small to fit into a float64", value)
		} else if math.IsInf(f, 0) {
			return 0, fmt.Errorf("%v is too big to fit into a float64", value)
		}
	}

	return f, nil
}

func BigFloatToInt(value *big.Float) (int64, error) {
	i, accuracy := value.Int64()
	if accuracy != big.Exact {
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
	if big.NewFloat(float64(i)).Cmp(value) != 0 {
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
	return i, nil
}

func BigFloatToUint(value *big.Float) (uint64, error) {
	u, accuracy := value.Uint64()
	if accuracy != big.Exact {
		return 0, fmt.Errorf("cannot convert %v to uint", value)
	}
	if big.NewFloat(float64(u)).Cmp(value) != 0 {
		return 0, fmt.Errorf("cannot convert %v to uint", value)
	}
	return u, nil
}

func BigFloatToString(value *big.Float) string {
	return value.Text('g', BitsToDecimalDigits(int(value.Prec())))
}

// Decimal float to other

func DecimalFloatToBigInt(value compact_float.DFloat, maxBase10Exponent int) (*big.Int, error) {
	if value.Exponent > int32(maxBase10Exponent) {
		return nil, fmt.Errorf("%v has a decimal exponential component (%v) that is too large for a big int (max %v)", value, value.Exponent, maxBase10Exponent)
	}
	return value.BigInt()
}

// big.Int to other

func BigIntToBigDecimalFloat(value *big.Int) apd.Decimal {
	return apd.Decimal{
		Coeff: *value,
	}
}

func BigIntToInt(value *big.Int) (int64, error) {
	if !value.IsInt64() {
		return 0, fmt.Errorf("%v is too big to fit into type int64", value)
	}
	return value.Int64(), nil
}

func BigIntToFloat(value *big.Int) (float64, error) {
	asText := value.Text(10)
	f, err := strconv.ParseFloat(asText, 64)
	if err != nil {
		return 0, err
	}
	asBigInt, accuracy := big.NewFloat(f).Int(nil)
	if accuracy != big.Exact {
		return 0, fmt.Errorf("cannot convert %v to float", value)
	}
	if asBigInt.Cmp(value) != 0 {
		return 0, fmt.Errorf("cannot convert %v to float", value)
	}
	return f, nil
}

func BigIntToUint(value *big.Int) (uint64, error) {
	if !value.IsUint64() {
		return 0, fmt.Errorf("%v cannot fit into type uint64", value)
	}
	return value.Uint64(), nil
}

// float to other

func FloatToBigDecimalFloat(value float64) (apd.Decimal, error) {
	var d apd.Decimal
	_, _, err := apd.BaseContext.SetString(&d, FloatToString(value))
	return d, err
}

func FloatToPBigDecimalFloat(value float64) (*apd.Decimal, error) {
	d, _, err := apd.NewFromString(FloatToString(value))
	return d, err
}

func FloatToBigInt(value float64, maxBase2Exponent int) (*big.Int, error) {
	return BigFloatToBigInt(big.NewFloat(value), maxBase2Exponent)
}

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

// int to other

func IntToBigDecimalFloat(value int64) apd.Decimal {
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

func IntToUint(value int64) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%v is negative, and cannot be represented by an unsigned int", value)
	}
	return uint64(value), nil
}

// uint to other

func UintToBigDecimalFloat(value uint64) apd.Decimal {
	return apd.Decimal{
		Coeff: *UintToBigInt(value),
	}
}

func UintToBigInt(value uint64) *big.Int {
	if value <= 0x7fffffffffffffff {
		return big.NewInt(int64(value))
	}

	bi := big.NewInt(int64(value >> 1))
	return bi.Lsh(bi, 1)
}

func UintToInt(value uint64) (int64, error) {
	if value > 0x7fffffffffffffff {
		return 0, fmt.Errorf("%v is too big to fit into type int64", value)
	}
	return int64(value), nil
}

// string to other

func StringToBigFloat(value string, significantDigits int) (*big.Float, error) {
	f, _, err := big.ParseFloat(value, 10, uint(DecimalDigitsToBits(significantDigits)), big.ToNearestEven)
	return f, err
}

// Bits to digits

var bitsToHexDigitsTable = []int{0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4}
var bitsToDecimalDigitsTable = []int{0, 1, 1, 1, 1, 2, 2, 2, 3, 3}
var decimalDigitsToBitsTable = []int{0, 4, 7}

func BitsToDecimalDigits(bitCount int) int {
	return (bitCount/10)*3 + bitsToDecimalDigitsTable[bitCount%10]
}

func DecimalDigitsToBits(digitCount int) int {
	triadCount := digitCount / 3
	remainder := digitCount % 3
	return triadCount*10 + decimalDigitsToBitsTable[remainder]
}

func BitsToHexDigits(bitCount int) int {
	return (bitCount/16)*4 + bitsToHexDigitsTable[bitCount&15]
}

func HexDigitsToBits(digitCount int) int {
	return digitCount * 4
}
