Concise Encoding
================

A go implementation of [concise encoding](https://github.com/kstenerud/concise-encoding).

Convenience like JSON, in a twin text/binary format with full datatype support.

 * **No schema necessary.** The simplicity of ad-hoc data, like in JSON & XML.
 * **Rich type support.** No more special code for dates, bytes, large values, etc.
 * **Edit text, transmit binary.** Twin binary and text formats are 1:1 compatible.
 * **Plug and play.** No extra compilation phase or descriptor files.
 * **Metadata.** Out-of-band data? Not a problem!
 * **Recursive/cyclic data.** Recursive structures? Also not a problem!
 * **Fully specified.** No more ambiguities in implementations.



Library Usage
-------------

```golang
package main

import (
	"fmt"
	"time"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"
)

func demonstrateCTEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := cte.Marshal(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CTE: %v\n", string(buffer.Bytes()))
	// Prints: Marshaled CTE: c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}
}

func demonstrateCTEUnmarshal() {
	data := bytes.NewBuffer([]byte(`c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}`))

	value, err := cte.Unmarshal(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CTE: %v\n", value)
	// Prints: Unmarshaled CTE: map[a key:2.5 900:2020-01-15/13:41:00.599]
}

func demonstrateCBEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := cbe.Marshal(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CBE: %v\n", buffer.Bytes())
	// Prints: Marshaled CBE: [1 121 133 97 32 107 101 121 112 0 0 32 64 106 132 3 155 188 18 0 32 109 47 80 0 123]
}

func demonstrateCBEUnmarshal() {
	data := bytes.NewBuffer([]byte{0x01, 0x79, 0x85, 0x61, 0x20, 0x6b, 0x65,
		0x79, 0x70, 0x00, 0x00, 0x20, 0x40, 0x6a, 0x84, 0x03, 0x9b, 0xbc, 0x12,
		0x00, 0x20, 0x6d, 0x2f, 0x50, 0x00, 0x7b})

	value, err := cbe.Unmarshal(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CBE: %v\n", value)
	// Prints: Unmarshaled CBE: map[a key:2.5 900:2020-01-15/13:41:00.599]
}

func main() {
	demonstrateCTEMarshal()
	demonstrateCTEUnmarshal()
	demonstrateCBEMarshal()
	demonstrateCBEUnmarshal()
}
```
