Snippet Check
=============

This tool checks documents for CTE snippets and validates them, printing out any errors it encounters. This can be useful for verifying that examples in documents are actually valid.

A snippet is anything within a standard markdown snippet marker using a "cte" tag. For example:

```cte
c1
{
    // An example CTE snippet

    "a" = [
        1
        2
        3
    ]
}
```


### Usage:

```
Usage: ./snippet_check [opts] <files>
  -q	quiet
  -v	verbose
```
