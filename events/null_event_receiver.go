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

package events

import (
	"math/big"
	"time"

	"github.com/cockroachdb/apd/v2"
	compact_float "github.com/kstenerud/go-compact-float"
	compact_time "github.com/kstenerud/go-compact-time"
)

// NullEventReceiver receives events and does nothing with them.
type NullEventReceiver struct{}

func NewNullEventReceiver() *NullEventReceiver {
	return &NullEventReceiver{}
}
func (_this *NullEventReceiver) OnBeginDocument()                    {}
func (_this *NullEventReceiver) OnVersion(uint64)                    {}
func (_this *NullEventReceiver) OnComment(bool, []byte)              {}
func (_this *NullEventReceiver) OnPadding(int)                       {}
func (_this *NullEventReceiver) OnNull()                             {}
func (_this *NullEventReceiver) OnBool(bool)                         {}
func (_this *NullEventReceiver) OnTrue()                             {}
func (_this *NullEventReceiver) OnFalse()                            {}
func (_this *NullEventReceiver) OnPositiveInt(uint64)                {}
func (_this *NullEventReceiver) OnNegativeInt(uint64)                {}
func (_this *NullEventReceiver) OnInt(int64)                         {}
func (_this *NullEventReceiver) OnBigInt(*big.Int)                   {}
func (_this *NullEventReceiver) OnFloat(float64)                     {}
func (_this *NullEventReceiver) OnBigFloat(*big.Float)               {}
func (_this *NullEventReceiver) OnDecimalFloat(compact_float.DFloat) {}
func (_this *NullEventReceiver) OnBigDecimalFloat(*apd.Decimal)      {}
func (_this *NullEventReceiver) OnNan(bool)                          {}
func (_this *NullEventReceiver) OnUID([]byte)                        {}
func (_this *NullEventReceiver) OnTime(time.Time)                    {}
func (_this *NullEventReceiver) OnCompactTime(compact_time.Time)     {}
func (_this *NullEventReceiver) OnArray(ArrayType, uint64, []byte)   {}
func (_this *NullEventReceiver) OnStringlikeArray(ArrayType, string) {}
func (_this *NullEventReceiver) OnArrayBegin(ArrayType)              {}
func (_this *NullEventReceiver) OnArrayChunk(uint64, bool)           {}
func (_this *NullEventReceiver) OnArrayData([]byte)                  {}
func (_this *NullEventReceiver) OnList()                             {}
func (_this *NullEventReceiver) OnMap()                              {}
func (_this *NullEventReceiver) OnEdge()                             {}
func (_this *NullEventReceiver) OnNode()                             {}
func (_this *NullEventReceiver) OnEnd()                              {}
func (_this *NullEventReceiver) OnMarker([]byte)                     {}
func (_this *NullEventReceiver) OnReference([]byte)                  {}
func (_this *NullEventReceiver) OnEndDocument()                      {}
