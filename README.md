Concise Encoding Go Library
===========================

This is the official reference implementation of [concise encoding](https://github.com/kstenerud/concise-encoding).

Concise Encoding is the friendly data format for human and machine. Ad-hoc, secure, with 1:1 compatible twin binary and text formats and rich type support.

 * **Edit text, transmit binary.** Humans love text. Machines love binary. With Concise Encoding, conversion is 1:1 and seamless.
 * **Secure By Design.** Precise decoding behavior means no poisoned data, no privilege escalations, and no DOSing!
 * **Rich type support.** Boolean, integer, float, string, bytes, time, URI, UUID, list, map, markup, and more!



Prerelease
----------

Please note that this library is still officially in the prerelease stage. The most common functionality is in-place and tested, but there are still some parts to do:

* Constants (partially implemented)
* Proper validation of markup and comment strings
* Extremely large values aren't tested well
* Everything marked with a `TODO` in the code

### Prerelease Document Version

The official Concise Encoding version during pre-release is `0`. Please use this version in your documents so as not to cause issues with old, potentially incompatible documents once the final release occurs. After release, document version 0 will be retired permanently, and considered invalid (there will be no backwards compatibility with the pre-release).

The pre-release library will also accept version `1` even though it's technically not correct, but will always emit version `0`.



Example Tool
------------

Click [here](https://github.com/kstenerud/enctool) for a simple data encoding format conversion tool that uses this library.



Library Usage
-------------

https://play.golang.org/p/6_pD6CQVLuN (If play.golang.org times out, just retry)

```golang
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/kstenerud/go-concise-encoding/ce"
)

func demonstrateCTEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCTE(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CTE: %v\n", string(buffer.Bytes()))
	// Prints: Marshaled CTE: c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}
}

func demonstrateCTEUnmarshal() {
	data := bytes.NewBuffer([]byte(`c1 {"a key"=2.5 900=2020-01-15/13:41:00.000599}`))

	value, err := ce.UnmarshalCTE(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CTE: %v\n", value)
	// Prints: Unmarshaled CTE: map[a key:2.5 900:2020-01-15/13:41:00.000599]
}

func demonstrateCBEMarshal() {
	dict := map[interface{}]interface{}{
		"a key": 2.5,
		900:     time.Date(2020, time.Month(1), 15, 13, 41, 0, 599000, time.UTC),
	}
	var buffer bytes.Buffer
	err := ce.MarshalCBE(dict, &buffer, nil)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Printf("Marshaled CBE: %v\n", buffer.Bytes())
	// Prints: Marshaled CBE: [3 1 121 133 97 32 107 101 121 112 0 0 32 64 106 132 3 155 188 18 0 32 109 47 80 0 123]
}

func demonstrateCBEUnmarshal() {
	data := bytes.NewBuffer([]byte{0x03, 0x01, 0x79, 0x85, 0x61, 0x20, 0x6b, 0x65,
		0x79, 0x70, 0x00, 0x00, 0x20, 0x40, 0x6a, 0x84, 0x03, 0x9b, 0xbc, 0x12,
		0x00, 0x20, 0x6d, 0x2f, 0x50, 0x00, 0x7b})

	value, err := ce.UnmarshalCBE(data, nil, nil)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v", err)
		return
	}
	fmt.Printf("Unmarshaled CBE: %v\n", value)
	// Prints: Unmarshaled CBE: map[a key:2.5 900:2020-01-15/13:41:00.000599]
}

func main() {
	demonstrateCTEMarshal()
	demonstrateCTEUnmarshal()
	demonstrateCBEMarshal()
	demonstrateCBEUnmarshal()
}
```


Architecture
------------

The architecture is designed to keep things as simple as possible for this reference implementation.

See [ARCHITECTURE.md](ARCHITECTURE.md)
