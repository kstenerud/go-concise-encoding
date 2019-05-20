Concise Binary Encoding
=======================

A go implementation of [concise binary encoding](https://github.com/kstenerud/concise-encoding/blob/master/cbe-specification.md).

Concise Binary Encoding is a general purpose, machine-readable, compact representation of semi-structured hierarchical data.


Goals
-----

  * General purpose encoding for a large number of applications
  * Supports the most common data types
  * Supports hierarchical data structuring
  * Binary format to minimize parsing costs
  * Little endian byte ordering to allow the most common systems to read directly off the wire
  * Balanced space and computation efficiency
  * Minimal complexity
  * Type compatible with [Concise Text Encoding (CTE)](cte-specification.md)



Library Usage
-------------

```golang
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
	encoder := cbe.NewCbeEncoder(100)
	cbe.Marshal(encoder, dict)
	fmt.Printf("Marshaled bytes: %v\n", encoder.EncodedBytes())
}

func demonstrateUnmarshal() {
	document := []byte{0x94, 0x6b, 0x84, 0x3, 0x79, 0x33, 0xda, 0xc3, 0xe3,
		0xbe, 0xb1, 0x58, 0x31, 0x85, 0x61, 0x20, 0x6b, 0x65, 0x79, 0x72, 0x0,
		0x0, 0x20, 0x40, 0x95}
	unmarshaler := new(cbe.Unmarshaler)
	decoder := cbe.NewCbeDecoder(100, unmarshaler)
	if err := decoder.Decode(document); err != nil {
		panic(fmt.Errorf("Unexpected error while decoding: %v", err))
	}
	fmt.Printf("Unmarshaled object: %v\n", unmarshaler.Unmarshaled())
}

func main() {
	demonstrateMarshal()
	demonstrateUnmarshal()
}
```

Output:

```
Marshaled bytes: [148 107 132 3 121 240 85 214 182 13 178 88 49 133 97 32 107 101 121 114 0 0 32 64 149]
Unmarshaled object: map[a key:2.5 900:2019-05-17 12:27:59.600037939 +0000 UTC]
```
