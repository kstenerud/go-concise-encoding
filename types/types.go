// Copyright 2019 Karl Stenerud
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Imposes the structural rules that enforce a well-formed concise encoding
// document.

// Concise Encoding types that don't exist in the standard library.
//
// Go is not expressive enough, nor is its type system capable enough to support
// these properly, so it's necessary to use interface{} everywhere, which also
// removes all compile-time type protections.
package types

// 128-bit universal identifier.
type UID [16]byte

func NewUID(contents []byte) UID {
	return UID{
		contents[0],
		contents[1],
		contents[2],
		contents[3],
		contents[4],
		contents[5],
		contents[6],
		contents[7],
		contents[8],
		contents[9],
		contents[10],
		contents[11],
		contents[12],
		contents[13],
		contents[14],
		contents[15],
	}
}

// Media, containing an IANA media type string, and the media data.
type Media struct {
	MediaType string
	Data      []byte
}

// An edge between two vertices in a graph.
type Edge struct {
	Source      interface{}
	Description interface{}
	Destination interface{}
}

const (
	EdgeFieldIndexSource = iota
	EdgeFieldIndexDescription
	EdgeFieldIndexDestination
)

// A node in a tree.
type Node struct {
	Value    interface{}
	Children []interface{}
}

const (
	NodeFieldIndexValue = iota
	NodeFieldIndexChildren
)
