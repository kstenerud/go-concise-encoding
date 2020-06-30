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
	"github.com/kstenerud/go-compact-float"
	"github.com/kstenerud/go-compact-time"
)

// DataEventReceiver receives data events (int, string, etc) and performs
// actions based on those events. Generally, this is used to drive complex
// object builders, and also the encoders.
type DataEventReceiver interface {
	OnVersion(version uint64)
	OnPadding(count int)
	OnNil()
	OnBool(value bool)
	OnTrue()
	OnFalse()
	OnPositiveInt(value uint64)
	OnNegativeInt(value uint64)
	OnInt(value int64)
	OnBigInt(value *big.Int)
	OnFloat(value float64)
	OnBigFloat(value *big.Float)
	OnDecimalFloat(value compact_float.DFloat)
	OnBigDecimalFloat(value *apd.Decimal)
	OnComplex(value complex128)
	OnNan(signaling bool)
	OnUUID(value []byte)
	OnTime(value time.Time)
	OnCompactTime(value *compact_time.Time)
	OnBytes(value []byte)
	OnString(value string)
	OnURI(value string)
	OnCustom(value []byte)
	OnBytesBegin()
	OnStringBegin()
	OnURIBegin()
	OnCustomBegin()
	OnArrayChunk(length uint64, isFinalChunk bool)
	OnArrayData(data []byte)
	OnList()
	OnMap()
	OnMarkup()
	OnMetadata()
	OnComment()
	OnEnd()
	OnMarker()
	OnReference()
	OnEndDocument()
}

// NullEventReceiver receives events and does nothing with them.
type NullEventReceiver struct{}

func NewNullEventReceiver() *NullEventReceiver {
	return &NullEventReceiver{}
}
func (_this *NullEventReceiver) OnVersion(version uint64)                      {}
func (_this *NullEventReceiver) OnPadding(count int)                           {}
func (_this *NullEventReceiver) OnNil()                                        {}
func (_this *NullEventReceiver) OnBool(value bool)                             {}
func (_this *NullEventReceiver) OnTrue()                                       {}
func (_this *NullEventReceiver) OnFalse()                                      {}
func (_this *NullEventReceiver) OnPositiveInt(value uint64)                    {}
func (_this *NullEventReceiver) OnNegativeInt(value uint64)                    {}
func (_this *NullEventReceiver) OnInt(value int64)                             {}
func (_this *NullEventReceiver) OnBigInt(value *big.Int)                       {}
func (_this *NullEventReceiver) OnFloat(value float64)                         {}
func (_this *NullEventReceiver) OnBigFloat(value *big.Float)                   {}
func (_this *NullEventReceiver) OnDecimalFloat(value compact_float.DFloat)     {}
func (_this *NullEventReceiver) OnBigDecimalFloat(value *apd.Decimal)          {}
func (_this *NullEventReceiver) OnComplex(value complex128)                    {}
func (_this *NullEventReceiver) OnNan(signaling bool)                          {}
func (_this *NullEventReceiver) OnUUID(value []byte)                           {}
func (_this *NullEventReceiver) OnTime(value time.Time)                        {}
func (_this *NullEventReceiver) OnCompactTime(value *compact_time.Time)        {}
func (_this *NullEventReceiver) OnBytes(value []byte)                          {}
func (_this *NullEventReceiver) OnString(value string)                         {}
func (_this *NullEventReceiver) OnURI(value string)                            {}
func (_this *NullEventReceiver) OnCustom(value []byte)                         {}
func (_this *NullEventReceiver) OnBytesBegin()                                 {}
func (_this *NullEventReceiver) OnStringBegin()                                {}
func (_this *NullEventReceiver) OnURIBegin()                                   {}
func (_this *NullEventReceiver) OnCustomBegin()                                {}
func (_this *NullEventReceiver) OnArrayChunk(length uint64, isFinalChunk bool) {}
func (_this *NullEventReceiver) OnArrayData(data []byte)                       {}
func (_this *NullEventReceiver) OnList()                                       {}
func (_this *NullEventReceiver) OnMap()                                        {}
func (_this *NullEventReceiver) OnMarkup()                                     {}
func (_this *NullEventReceiver) OnMetadata()                                   {}
func (_this *NullEventReceiver) OnComment()                                    {}
func (_this *NullEventReceiver) OnEnd()                                        {}
func (_this *NullEventReceiver) OnMarker()                                     {}
func (_this *NullEventReceiver) OnReference()                                  {}
func (_this *NullEventReceiver) OnEndDocument()                                {}
