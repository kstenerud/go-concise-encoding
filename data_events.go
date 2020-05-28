package concise_encoding

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
func (this *NullEventReceiver) OnVersion(version uint64)                      {}
func (this *NullEventReceiver) OnPadding(count int)                           {}
func (this *NullEventReceiver) OnNil()                                        {}
func (this *NullEventReceiver) OnBool(value bool)                             {}
func (this *NullEventReceiver) OnTrue()                                       {}
func (this *NullEventReceiver) OnFalse()                                      {}
func (this *NullEventReceiver) OnPositiveInt(value uint64)                    {}
func (this *NullEventReceiver) OnNegativeInt(value uint64)                    {}
func (this *NullEventReceiver) OnInt(value int64)                             {}
func (this *NullEventReceiver) OnBigInt(value *big.Int)                       {}
func (this *NullEventReceiver) OnFloat(value float64)                         {}
func (this *NullEventReceiver) OnDecimalFloat(value compact_float.DFloat)     {}
func (this *NullEventReceiver) OnBigDecimalFloat(value *apd.Decimal)          {}
func (this *NullEventReceiver) OnComplex(value complex128)                    {}
func (this *NullEventReceiver) OnNan(signaling bool)                          {}
func (this *NullEventReceiver) OnUUID(value []byte)                           {}
func (this *NullEventReceiver) OnTime(value time.Time)                        {}
func (this *NullEventReceiver) OnCompactTime(value *compact_time.Time)        {}
func (this *NullEventReceiver) OnBytes(value []byte)                          {}
func (this *NullEventReceiver) OnString(value string)                         {}
func (this *NullEventReceiver) OnURI(value string)                            {}
func (this *NullEventReceiver) OnCustom(value []byte)                         {}
func (this *NullEventReceiver) OnBytesBegin()                                 {}
func (this *NullEventReceiver) OnStringBegin()                                {}
func (this *NullEventReceiver) OnURIBegin()                                   {}
func (this *NullEventReceiver) OnCustomBegin()                                {}
func (this *NullEventReceiver) OnArrayChunk(length uint64, isFinalChunk bool) {}
func (this *NullEventReceiver) OnArrayData(data []byte)                       {}
func (this *NullEventReceiver) OnList()                                       {}
func (this *NullEventReceiver) OnMap()                                        {}
func (this *NullEventReceiver) OnMarkup()                                     {}
func (this *NullEventReceiver) OnMetadata()                                   {}
func (this *NullEventReceiver) OnComment()                                    {}
func (this *NullEventReceiver) OnEnd()                                        {}
func (this *NullEventReceiver) OnMarker()                                     {}
func (this *NullEventReceiver) OnReference()                                  {}
func (this *NullEventReceiver) OnEndDocument()                                {}
