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

package concise_encoding

import (
	"bytes"
	"math/big"
	"time"

	"github.com/kstenerud/go-concise-encoding/cbe"
	"github.com/kstenerud/go-concise-encoding/cte"

	"github.com/cockroachdb/apd/v2"
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
	"github.com/kstenerud/go-concise-encoding/test"
)

func cbeDecode(document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
	err = cbe.Decode(document, receiver, nil)
	events = receiver.Events
	return
}

func cbeEncodeDecode(expected ...*test.TEvent) (events []*test.TEvent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	buffer := &bytes.Buffer{}
	encoder := cbe.NewEncoder(buffer, nil)
	test.InvokeEvents(encoder, expected...)
	document := buffer.Bytes()

	return cbeDecode(document)
}

func cteDecode(document []byte) (events []*test.TEvent, err error) {
	receiver := test.NewTER()
	err = cte.Decode(document, receiver, nil)
	events = receiver.Events
	return
}

func cteEncode(events ...*test.TEvent) []byte {
	buffer := &bytes.Buffer{}
	encoder := cte.NewEncoder(buffer, nil)
	test.InvokeEvents(encoder, events...)
	return buffer.Bytes()
}

func cteEncodeDecode(events ...*test.TEvent) (decodedEvents []*test.TEvent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	return cteDecode(cteEncode(events...))
}

func TT() *test.TEvent                       { return test.TT() }
func FF() *test.TEvent                       { return test.FF() }
func I(v int64) *test.TEvent                 { return test.I(v) }
func F(v float64) *test.TEvent               { return test.F(v) }
func BF(v *big.Float) *test.TEvent           { return test.BF(v) }
func DF(v compact_float.DFloat) *test.TEvent { return test.DF(v) }
func BDF(v *apd.Decimal) *test.TEvent        { return test.BDF(v) }
func V(v uint64) *test.TEvent                { return test.V(v) }
func N() *test.TEvent                        { return test.N() }
func PAD(v int) *test.TEvent                 { return test.PAD(v) }
func B(v bool) *test.TEvent                  { return test.B(v) }
func PI(v uint64) *test.TEvent               { return test.PI(v) }
func NI(v uint64) *test.TEvent               { return test.NI(v) }
func BI(v *big.Int) *test.TEvent             { return test.BI(v) }
func CPLX(v complex128) *test.TEvent         { return test.CPLX(v) }
func NAN() *test.TEvent                      { return test.NAN() }
func SNAN() *test.TEvent                     { return test.SNAN() }
func UUID(v []byte) *test.TEvent             { return test.UUID(v) }
func GT(v time.Time) *test.TEvent            { return test.GT(v) }
func CT(v *compact_time.Time) *test.TEvent   { return test.CT(v) }
func BIN(v []byte) *test.TEvent              { return test.BIN(v) }
func S(v string) *test.TEvent                { return test.S(v) }
func URI(v string) *test.TEvent              { return test.URI(v) }
func CUST(v []byte) *test.TEvent             { return test.CUST(v) }
func BB() *test.TEvent                       { return test.BB() }
func SB() *test.TEvent                       { return test.SB() }
func UB() *test.TEvent                       { return test.UB() }
func CB() *test.TEvent                       { return test.CB() }
func AC(l uint64, term bool) *test.TEvent    { return test.AC(l, term) }
func AD(v []byte) *test.TEvent               { return test.AD(v) }
func L() *test.TEvent                        { return test.L() }
func M() *test.TEvent                        { return test.M() }
func MUP() *test.TEvent                      { return test.MUP() }
func META() *test.TEvent                     { return test.META() }
func CMT() *test.TEvent                      { return test.CMT() }
func E() *test.TEvent                        { return test.E() }
func MARK() *test.TEvent                     { return test.MARK() }
func REF() *test.TEvent                      { return test.REF() }
func ED() *test.TEvent                       { return test.ED() }
