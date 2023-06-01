Concise Encoding Codec Test System
==================================

All codec tests are written as portable CTE test files so that the same tests can be run against multiple implementations (different languages, different platforms, etc).

The tests are kept in the [`suites`](suites) directory, including [template.cte](tests/suites/template.cte) (which contains a full description of the test file structure and what everything's for).


Running the Tests
-----------------

[`tests/tests_test.go`](tests/tests_test.go) will run all tests in the [`tests/suites`](tests/suites) directory (and all subdirs). They'll be automatically run if you invoke `go test ./...`.

To run a specific test, import the test runner and give it a test suite to run:

```go
package xyz

import (
    "testing"

    "github.com/kstenerud/go-concise-encoding/test/test_runner"
)

func TestSomething(t *testing.T) {
    test_runner.RunCEUnitTests(t, "path/to/something.cte")
}
```


Running tests on your own implementation
----------------------------------------

You will need to have at least the following types implemented to be able to load a basic test file:

- String
- Integer
- Boolean
- List
- Map
- Uint8 Array

The event shorthand system will require a bit more work to implement. See the code in the [test](test) directory for an example. There's also an Antlr grammar ([lexer](codegen/test/CEEventLexer.g4), [parser](codegen/test/CEEventParser.g4)) to parse the event shorthand.

If you're building a CBE codec and don't yet have a CTE decoder, you could convert the test files from CTE to CBE using [enctool](https://github.com/kstenerud/enctool) and then load them via your CBE decoder.
