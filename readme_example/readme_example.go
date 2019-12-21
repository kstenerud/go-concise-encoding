package main

import (
	"fmt"
	"time"

	"github.com/kstenerud/go-cbe"
)

func demonstrateMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Now(),
	}
	bytes, err := cbe.MarshalCBE(dict)
	if err != nil {
		fmt.Printf("Error marshaling object: %v\n", err)
	}
	fmt.Printf("Marshaled bytes: %v\n", bytes)
}

func demonstrateUnmarshal() {
	document := []byte{0x94, 0x6b, 0x84, 0x3, 0x79, 0x33, 0xda, 0xc3, 0xe3,
		0xbe, 0xb1, 0x58, 0x31, 0x85, 0x61, 0x20, 0x6b, 0x65, 0x79, 0x72, 0x0,
		0x0, 0x20, 0x40, 0x95}

	dict := map[interface{}]interface{}{}
	if err := cbe.UnmarshalCBE(document, &dict); err != nil {
		panic(fmt.Errorf("Error unmarshaling object: %v", err))
	}
	fmt.Printf("Unmarshaled object: %v\n", dict)
}

func main() {
	demonstrateMarshal()
	demonstrateUnmarshal()
}
