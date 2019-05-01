package cbe

import (
	"fmt"
	"testing"
)

func TestGetErrorType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Caught type %T, value \"%v\"\n", r, r)
		}
	}()

	data := make([]byte, 1)
	fmt.Println(data[5])
}
