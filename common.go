package concise_encoding

import (
	"time"

	"github.com/kstenerud/go-compact-time"
)

type ConciseEncodingEventHandler interface {
	OnVersion(version uint64)
	OnPadding(count int)
	OnNil()
	OnBool(value bool)
	OnTrue()
	OnFalse()
	OnPositiveInt(value uint64)
	OnNegativeInt(value uint64)
	OnInt(value int64)
	OnFloat(value float64)
	OnComplex(value complex128)
	OnNan()
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

type NullEventHandler struct{}

func NewNullEventHandler() *NullEventHandler {
	return &NullEventHandler{}
}
func (this *NullEventHandler) OnVersion(version uint64)                      {}
func (this *NullEventHandler) OnPadding(count int)                           {}
func (this *NullEventHandler) OnNil()                                        {}
func (this *NullEventHandler) OnBool(value bool)                             {}
func (this *NullEventHandler) OnTrue()                                       {}
func (this *NullEventHandler) OnFalse()                                      {}
func (this *NullEventHandler) OnPositiveInt(value uint64)                    {}
func (this *NullEventHandler) OnNegativeInt(value uint64)                    {}
func (this *NullEventHandler) OnInt(value int64)                             {}
func (this *NullEventHandler) OnFloat(value float64)                         {}
func (this *NullEventHandler) OnComplex(value complex128)                    {}
func (this *NullEventHandler) OnNan()                                        {}
func (this *NullEventHandler) OnUUID(value []byte)                           {}
func (this *NullEventHandler) OnTime(value time.Time)                        {}
func (this *NullEventHandler) OnCompactTime(value *compact_time.Time)        {}
func (this *NullEventHandler) OnBytes(value []byte)                          {}
func (this *NullEventHandler) OnString(value string)                         {}
func (this *NullEventHandler) OnURI(value string)                            {}
func (this *NullEventHandler) OnCustom(value []byte)                         {}
func (this *NullEventHandler) OnBytesBegin()                                 {}
func (this *NullEventHandler) OnStringBegin()                                {}
func (this *NullEventHandler) OnURIBegin()                                   {}
func (this *NullEventHandler) OnCustomBegin()                                {}
func (this *NullEventHandler) OnArrayChunk(length uint64, isFinalChunk bool) {}
func (this *NullEventHandler) OnArrayData(data []byte)                       {}
func (this *NullEventHandler) OnList()                                       {}
func (this *NullEventHandler) OnMap()                                        {}
func (this *NullEventHandler) OnMarkup()                                     {}
func (this *NullEventHandler) OnMetadata()                                   {}
func (this *NullEventHandler) OnComment()                                    {}
func (this *NullEventHandler) OnEnd()                                        {}
func (this *NullEventHandler) OnMarker()                                     {}
func (this *NullEventHandler) OnReference()                                  {}
func (this *NullEventHandler) OnEndDocument()                                {}
