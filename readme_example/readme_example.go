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
	encoder := cbe.NewCbeEncoder(cbe.InlineContainerTypeNone, nil, 100)
	cbe.Marshal(encoder, cbe.InlineContainerTypeNone, dict)
	fmt.Printf("Marshaled bytes: %v\n", encoder.EncodedBytes())
}

func demonstrateUnmarshal() {
	document := []byte{0x94, 0x6b, 0x84, 0x3, 0x79, 0x33, 0xda, 0xc3, 0xe3,
		0xbe, 0xb1, 0x58, 0x31, 0x85, 0x61, 0x20, 0x6b, 0x65, 0x79, 0x72, 0x0,
		0x0, 0x20, 0x40, 0x95}
	unmarshaler := new(cbe.Unmarshaler)
	decoder := cbe.NewCbeDecoder(cbe.InlineContainerTypeNone, 100, unmarshaler)
	if err := decoder.Decode(document); err != nil {
		panic(fmt.Errorf("Unexpected error while decoding: %v", err))
	}
	fmt.Printf("Unmarshaled object: %v\n", unmarshaler.Unmarshaled())
}

func main() {
	demonstrateMarshal()
	demonstrateUnmarshal()
}
