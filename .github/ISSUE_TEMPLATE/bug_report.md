---
name: Bug report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**Test Case**

<!--
The easiest way to describe the problem is to build a quick test case.
This link will get you set up quickly: https://github.com/kstenerud/go-concise-encoding/blob/master/REPORTING-BUGS.md
-->

```cte
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
            "events" = [
                            "bd" "v 0" "m"
                                "s a" "tt"
                            "e" "ed"
                       ]
            "from"   = ["t" "b" "e"]
            "to"     = ["t" "b" "e"]
        }
    ]
}
```
