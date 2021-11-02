Bug Reporting
=============

When this library is misbehaving, the quickest and easiest way to describe the problem is to demonstrate it in action. This is where Concise Encoding [unit tests](TESTING.md) shine, as they allow you to provide an input and specify what the output should look like.

The unit tests have more knobs than are needed for your average bug report, so the [bugreport](bugreport) directory helps you build your bug report repro case with minimal fuss.



Contents
--------

 * [Quick Example](#quick-example)
 * [Templates](#templates)
 * [Building your bug report](#building-your-bug-report)



Quick Example
-------------

Here's a quick example of what a unit test looks like:

```cte
c0
{
    "tests" = [
        {
            "name" = "My Bug Report"
            "cte" = "\.%%%%
c0
{
    "a" = true
}
            %%%%"
            "cbe"    = |u8x 83 00 79 81 61 7d 7b|
            "from"   = ["t"]
            "to"     = ["b"]
        }
    ]
}
```

This unit test attempts to decode a CTE document (in this case `c0 {"a"=true}`), validates that it is a properly formed document, converts it to CBE, and compares it to the expected [`83 00 79 81 61 7d 7b`].

### Fields:

 * **name**: The name of this test (call it whatever you want).
 * **cte**: A CTE document (encoded using a [verbatim sequence](https://github.com/kstenerud/concise-encoding/blob/master/cte-specification.md#verbatim-sequence) with the sentinel `%%%%` - change it if your data contains this character sequence).
 * **cbe**: A CBE document.
 * **from**: Specifies what field(s) to read from as input (in this case "t" meaning text, i.e. the contents of "cte").
 * **to**: Specifies what field(s) to compare to as output (in this case "b" meaning binary, i.e. the contents of "cbe").

### Notes:

 * After trimming leading and trailing whitespace, the expected CTE document is byte compared with the actual output, so it **must** match the pretty printing of the library in order to be considered a match.
 * There are also other fields that may be used, and more extensive `"from"` and `"to"` field combinations that are possible. For more information, please read the [unit test documentation](TESTING.md).



Templates
---------

Although the [unit test documentation](TESTING.md) describes fully how these unit tests work, the impatient can get by just using a template.

The [templates](bugreport/templates) directory contains templates for the kinds of bugs most likely to be encountered:

| Situation                     | "from" field | "to" field | Notes                                                                  |
| ----------------------------- | ------------ | ---------- | ---------------------------------------------------------------------- |
| CBE is output incorrectly     | "t" or "e"   | "b"        | "cbe" contains what the output should look like                        |
| CTE is output incorrectly     | "t" or "e"   | "t"        | "cte" contains what the output should look like                        |
| CBE is parsed incorrectly     | "b"          | "t" or "e" | "cte" or "events" contains what should be generated from the CBE data  |
| CTE is parsed incorrectly     | "t"          | "t" or "e" | "cbe" or "events" contains what should be generated from the CTE data  |
| Document wrongfully rejected  | "t" or "b"   | nothing    | "cbe" or "cte" contains the document that is being wrongfully rejected |
| Document wrongfully allowed   | "t" or "b"   | nothing    | "cbe" or "cte" contains the document that is being wrongfully allowed  |

Simply copy a template and overwrite [bugreport.cte](bugreport/bugreport.cte), then fill in your CTE and CBE data (and possibly events).



Building your bug report
------------------------

- Update to the latest go-concise-encoding to make sure the problem hasn't already been fixed.

- Modify [bugreport.cte](bugreport/bugreport.cte) (or copy over one of the [templates](bugreport/templates)) to include the input and output that demonstrates the issue you're encountering.

- run `go test` inside the [bugreport](bugreport) directory and observe the test failing in the way you expect it to.

- Include the contents of bugreport.cte in your [bug report](https://github.com/kstenerud/go-concise-encoding/issues).
