ANTLR4 Grammar for Concise Text Encoding
========================================

This Antlr4 grammar is used to generate the CTE parser for this library.

[`codegen`](..) builds the parser code from this grammar.

[Parser.go](../../cte/parser.go) integrates this grammar into the library.

### TODO

- Define the proper character ranges for stringlike (see the bottom of [CTELexer.g4](CTELexer.g4)).
