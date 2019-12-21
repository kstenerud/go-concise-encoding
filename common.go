package cbe

import (
	"errors"
)

var BufferExhaustedError = errors.New("Buffer exhausted")

type InlineContainerType int

const (
	InlineContainerTypeNone InlineContainerType = iota
	InlineContainerTypeList
	InlineContainerTypeMap
)

var digitsMax = [...]uint64{
	0,
	9,
	99,
	999,
	9999,
	99999,
	999999,
	9999999,
	99999999,
	999999999,
	9999999999,
	99999999999,
	999999999999,
	9999999999999,
	99999999999999,
	999999999999999,
	9999999999999999,
	99999999999999999,
	999999999999999999,
	9999999999999999999, // 19 digits
	// Max digits for uint64 is 20
}

func CountDigits(value uint64) int {
	// This is MUCH faster than the string method, and 4x faster than int(math.Log10(float64(value))) + 1
	// Subdividing any further yields no performance gains.
	if value <= digitsMax[10] {
		for i := 1; i < 10; i++ {
			if value <= digitsMax[i] {
				return i
			}
		}
		return 10
	}

	for i := 11; i < 20; i++ {
		if value <= digitsMax[i] {
			return i
		}
	}
	return 20
}
