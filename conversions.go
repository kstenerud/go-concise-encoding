package concise_encoding

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/cockroachdb/apd/v2"
)

// apd.Decimal to other

func bigDecimalFloatToBigFloat(value *apd.Decimal) (*big.Float, error) {
	return stringToBigFloat(value.Text('g'), int(value.NumDigits()))
}

func bigDecimalFloatToBigInt(value *apd.Decimal) *big.Int {
	// TODO: Check for DOSable values
	exp := big.NewInt(int64(value.Exponent))
	exp.Exp(bigInt10, exp, nil)
	return exp.Mul(exp, &value.Coeff)
}

func bigDecimalFloatToUint(value *apd.Decimal) (uint64, error) {
	if i, err := value.Int64(); err == nil {
		return uint64(i), nil
	}

	bf, err := bigDecimalFloatToBigFloat(value)
	if err != nil {
		return 0, err
	}
	u, accuracy := bf.Uint64()
	if accuracy != big.Exact {
		return 0, fmt.Errorf("%v cannot fit into a uint64", value)
	}
	return u, nil
}

// big.Float to other

func bigFloatToPBigDecimalFloat(value *big.Float) (*apd.Decimal, error) {
	d, _, err := apd.NewFromString(bigFloatToString(value))
	return d, err
}

func bigFloatToBigInt(value *big.Float) (*big.Int, error) {
	// TODO: Check for DOSable values
	bi, accuracy := value.Int(new(big.Int))
	if accuracy != big.Exact {
		return nil, fmt.Errorf("%v cannot fit into a big.Int", value)
	}
	return bi, nil
}

func bigFloatToFloat(value *big.Float) (float64, error) {
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

func bigFloatToString(value *big.Float) string {
	return value.Text('g', bitsToDecimalDigits(int(value.Prec())))
}

// big.Int to other

func bigIntToBigDecimalFloat(value *big.Int) apd.Decimal {
	return apd.Decimal{
		Coeff: *value,
	}
}

func bigIntToInt(value *big.Int) (int64, error) {
	if !value.IsInt64() {
		return 0, fmt.Errorf("%v is too big to fit into type int64", value)
	}
	return value.Int64(), nil
}

func bigIntToFloat(value *big.Int) (float64, error) {
	return stringToFloat(value.Text(10))
}

func bigIntToUint(value *big.Int) (uint64, error) {
	if !value.IsUint64() {
		return 0, fmt.Errorf("%v cannot fit into type uint64", value)
	}
	return value.Uint64(), nil
}

// float to other

func floatToBigDecimalFloat(value float64) (apd.Decimal, error) {
	var d apd.Decimal
	_, _, err := apd.BaseContext.SetString(&d, floatToString(value))
	return d, err
}

func floatToPBigDecimalFloat(value float64) (*apd.Decimal, error) {
	d, _, err := apd.NewFromString(floatToString(value))
	return d, err
}

func floatToBigInt(value float64) *big.Int {
	bi, _ := big.NewFloat(value).Int(nil)
	return bi
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func floatToUint(value float64) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%v is negative, and cannot be represented by an unsigned int", value)
	}
	return uint64(value), nil
}

// int to other

func intToBigDecimalFloat(value int64) apd.Decimal {
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

func intToUint(value int64) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%v is negative, and cannot be represented by an unsigned int", value)
	}
	return uint64(value), nil
}

// uint to other

func uintToBigDecimalFloat(value uint64) apd.Decimal {
	return apd.Decimal{
		Coeff: *uintToBigInt(value),
	}
}

func uintToBigInt(value uint64) *big.Int {
	if value <= 0x7fffffffffffffff {
		return big.NewInt(int64(value))
	}

	bi := big.NewInt(int64(value >> 1))
	return bi.Lsh(bi, 1)
}

func uintToInt(value uint64) (int64, error) {
	if value > 0x7fffffffffffffff {
		return 0, fmt.Errorf("%v is too big to fit into type int64", value)
	}
	return int64(value), nil
}

// string to other

func stringToBigFloat(value string, significantDigits int) (*big.Float, error) {
	f, _, err := big.ParseFloat(value, 10, uint(decimalDigitsToBits(significantDigits)), big.ToNearestEven)
	return f, err
}

func stringToFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// Bits to digits

var bitsToDecimalDigitsTable = []int{0, 1, 1, 1, 1, 2, 2, 2, 3, 3}
var decimalDigitsToBitsTable = []int{0, 4, 7}

func bitsToDecimalDigits(bitCount int) int {
	return (bitCount/10)*3 + bitsToDecimalDigitsTable[bitCount%10]
}

func decimalDigitsToBits(digitCount int) int {
	triadCount := digitCount / 3
	remainder := digitCount % 3
	return triadCount*10 + decimalDigitsToBitsTable[remainder]
}
