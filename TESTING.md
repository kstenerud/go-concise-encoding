Concise Encoding Unit Test System
=================================

General unit tests are now defined as test suites in special CTE documents in order to make tests easier to write, and to make them portable across implementations.

To invoke a unit test, simply import the test runner and give it a test suite to run. It can then be run from the command line using the familiar `go test`:

```go
package xyz

import (
    "testing"

    "github.com/kstenerud/go-concise-encoding/test/test_runner"
)

func TestSomething(t *testing.T) {
    test_runner.RunCEUnitTests(t, "something.cte")
}
```


Anatomy of a test suite
-----------------------

A test suite is a CTE document containing a series of unit tests to run.

#### Sample test suite:

```cte
c0
{
    "ceTests" = [

        // First unit test
        {
            "name" = "My test"
            "cte" = "\.%%%%
c0
{
    "a" = true
}
            %%%%"
            "cbe"    = |u8x 81 00 79 81 61 7d 7b|
            "events" = [
                            "bd" "v 0" "m"
                                "s a" "b true"
                            "e" "ed"
                       ]
            "from"     = ["t" "b" "e"]
            "to"       = ["t" "b" "e"]
            "mustFail" = false
            "trace"    = false
            "debug"    = false
            "skip"     = false
        }

        // Add more unit tests here
    ]
}
```


### `name` field

The name that will be printed in messsages associated with this unit test.


### `cbe` field

If present, contains a complete CBE document that will be used as input and/or output.


### `cte` field

If present, contains a complete CTE document that will be used as input and/or output.

This field is best encoded using a [verbatim sequence](https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#verbatim-sequence) (in this case, with the sentinel `%%%%`) so that we don't need to worry about escape sequences. Set the sentinel to something that does not occur inside your CTE document.

When used for output verification, the document will have leading and trailing whitespace stripped (in this case the linefeed before `c0` and the linefeed after the map), after which it will be byte-compared with the output from the library. This means that the CTE document must match the pretty printing of the library!


### `events` field

If present, this field contains string-encoded elements that are a kind of shorthand for the internal events that this library generates. They're specific to this library but in general map pretty closely to the CBE data types and payloads. They're used for catching unlikely cases where both the CTE and CBE codecs cause the same erroneous condition and would therefore agree with each other and wrongfully pass the test.

| Code  | Param 1       | Param 2 | Meaning                         | Example                                    |
| ----- | ------------- | ------- | ------------------------------- | ------------------------------------------ |
| v     | int           |         | Version                         | `v 1`                                      |
| pad   | int           |         | Padding (byte count)            | `pad 3`                                    |
| com   | bool          | string  | Comment (multiline?, contents)  | `com false this is a single line comment`  |
| null  |               |         | Null                            | `null`                                     |
| b     | bool          |         | Boolean                         | `b false`                                  |
| n     | int or float  |         | Number: auto-detect             | i: `n 1`, df: `n -5.3`, bf: `n 0x1.f`      |
| i     | int           |         | Integer                         | `i -5`                                     |
| bf    | binary float  |         | Binary Float                    | `bf 1.5`                                   |
| df    | decimal float |         | Decimal Float                   | `df 1.5`                                   |
| uid   | uuid          |         | Universal ID                    | `uid f1ce4567-e89b-12d3-a456-426655440000` |
| t     | time          |         | Time                            | `t 2010-07-15/13:28:15.241588124/E/Oslo`   |
| s     | string        |         | String                          | `s this is a string`                       |
| rid   | url           |         | Resource ID                     | `rid https://www.example.com`              |
| cb    | bytes         |         | Custom Binary                   | `cb 00 9f 1a bb 54`                        |
| ct    | string        |         | Custom Text                     | `ct cplx(5.0, 2)`                          |
| ab    | bits          |         | Array: binary                   | `ab 1101001011101011100`                   |
| ai8   | int8 elems    |         | Array: int8                     | `ai8 -4 59 100 -36`                        |
| ai16  | int16 elems   |         | Array: int16                    | `ai16 1000 2000 -3000`                     |
| ai32  | int32 elems   |         | Array: int32                    | `ai32 999999 -3455234`                     |
| ai64  | int64 elems   |         | Array: int64                    | `ai64 99999999999999 34534634730456`       |
| au8   | uint8 elems   |         | Array: uint8                    | `au8 100 54 12 0 6`                        |
| au16  | uint16 elems  |         | Array: uint16                   | `au16 1000 4000 1234`                      |
| au32  | uint32 elems  |         | Array: uint32                   | `au32 56872348 34693939 223444`            |
| au64  | uint64 elems  |         | Array: uint64                   | `au64 87993485978345987 34958345729835`    |
| af16  | float16 elems |         | Array: float16                  | `af16 1.5 2.923`                           |
| af32  | float32 elems |         | Array: float32                  | `af32 6.41239 4.42e9`                      |
| af64  | float64 elems |         | Array: float64                  | `af64 4.2944923e-40 2.4e70`                |
| au    | uuid elems    |         | Array: uid                      | `au f1ce4567-e89b-12d3-a456-426655440000 a1424567-e544-1223-a4f6-4266c5440a00` |
| sb    |               |         | Begin string                    | `sb`                                       |
| rb    |               |         | Begin resource ID               | `rb`                                       |
| rrb   |               |         | Begin remote ref                | `rrb`                                      |
| cbb   |               |         | Begin custom binary             | `cbb`                                      |
| ctb   |               |         | Begin custom text               | `ctb`                                      |
| abb   |               |         | Begin array: bit                | `abb`                                      |
| ai8b  |               |         | Begin array: int8               | `ai8b`                                     |
| ai16b |               |         | Begin array: int16              | `ai16b`                                    |
| ai32b |               |         | Begin array: int32              | `ai32b`                                    |
| ai64b |               |         | Begin array: int64              | `ai64b`                                    |
| au8b  |               |         | Begin array: uint8              | `au8b`                                     |
| au16b |               |         | Begin array: uint16             | `au16b`                                    |
| au32b |               |         | Begin array: uint32             | `au32b`                                    |
| au64b |               |         | Begin array: uint64             | `au64b`                                    |
| aub   |               |         | Begin array: uid                | `aub`                                      |
| mb    |               |         | Begin Media                     | `mb`                                       |
| acm   | uint          |         | Array chunk w/more (elem count) | `acm 10`                                   |
| acl   | uint          |         | Last array chunk (elem count)   | `acl 10`                                   |
| ad    | bytes         |         | Array data as binary            | `ad f8 ee 41 9a 13`                        |
| at    | string        |         | Array data as text              | `at some text data`                        |
| l     |               |         | Begin list                      | `l`                                        |
| m     |               |         | Begin map                       | `m`                                        |
| node  |               |         | Begin node                      | `node`                                     |
| edge  |               |         | Begin edge                      | `edge`                                     |
| e     |               |         | End container                   | `e`                                        |
| mark  | string        |         | Marker                          | `mark some-id`                             |
| ref   | string        |         | Reference                       | `ref some-id`                              |
| rref  | url           |         | Remote Reference                | `rref https://www.example.com`             |

**Notes**:
 * Integer and float values can be written using other bases (base 2, 8, 10, 16 for decimal, 10 and 16 for float).
 * Integer array types also have a hex variant, where all elements must be written in hex without the `0x` prefix (e.g. the hex variant of `au16` is `au16x`).


### `from` and `to` fields

These fields determine which of the `cbe`, `cte`, and `events` fields will be used for input and/or output. They are implemented as lists that can contain any combination of three unique string values:

* `"t"`: Use the content in the `cte` field
* `"b"`: Use the content in the `cbe` field
* `"e"`: Use the content in the `events` field

The test runner will attempt to convert each of the `from` types to each of the `to` types (so for example `from` `["t" "b" "e"]` and `to` `["t" "b" "e"]` will result in 9 combinations).

For each from-to combination, it will ingest the "from" data and expect its output to match the "to" data. If the `to` field is empty or missing, it will simply ingest the "from" data and validate that the document is well formed.

The `from` field must always be present, and must contain at least one type.

#### Possible operations from combinations:

| From | To  | Operation         |
| ---- | --- | ----------------- |
| "t"  | -   | Validate CTE      |
| "t"  | "t" | CTE -> CTE        |
| "t"  | "b" | CTE -> CBE        |
| "t"  | "e" | CTE -> Events     |
| "b"  | -   | Validate CBE      |
| "b"  | "t" | CBE -> CTE        |
| "b"  | "b" | CBE -> CBE        |
| "b"  | "e" | CBE -> Events     |
| "e"  | -   | Validate Events   |
| "e"  | "t" | Events -> CTE     |
| "e"  | "b" | Events -> CBE     |
| "e"  | "e" | Events -> Events  |


### `fail` field

If true, this task is expected to fail. If the task completes without error, the test has failed.

Default false.


### `trace` field

If true, the test will print a stack trace to the point that caused the error.

Default false.


### `debug` field

If true, print out extra debug information while running this test.

Default false.


### `skip` field

If true, skip this test when running the test suite.

Default false.
