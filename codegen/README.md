Code Generator for go-concise-encoding
======================================

This program performs all code generation for the go-concise-encoding project.

To build everything, you'll need to extract https://www.unicode.org/Public/UCD/latest/ucdxml/ucd.all.flat.zip

### Usage

 * Build all but chars: `go build && ./codegen`
 * Build everything: `go build && ./codegen -unicode /path/to/ucd.all.flat.xml`

Use `--help` for a list of all options.
