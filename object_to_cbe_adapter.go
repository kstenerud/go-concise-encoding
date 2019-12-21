package cbe

import (
	"fmt"
	"net/url"
	"time"
)

type ObjectToCBEAdapter struct {
	encoder *CBEEncoder
}

func NewObjectToCBEAdapter(encoder *CBEEncoder) *ObjectToCBEAdapter {
	this := new(ObjectToCBEAdapter)
	this.Init(encoder)
	return this
}

func (this *ObjectToCBEAdapter) Init(encoder *CBEEncoder) {
	this.encoder = encoder
}

func (this *ObjectToCBEAdapter) OnNil() error {
	return this.encoder.Nil()
}

func (this *ObjectToCBEAdapter) OnBool(value bool) error {
	return this.encoder.Bool(value)
}

func (this *ObjectToCBEAdapter) OnUint(value uint64) error {
	return this.encoder.PositiveInt(value)
}

func (this *ObjectToCBEAdapter) OnInt(value int64) error {
	if value < 0 {
		return this.encoder.NegativeInt(uint64(-value))
	}
	return this.encoder.PositiveInt(uint64(value))
}

func (this *ObjectToCBEAdapter) OnFloat(value float64) error {
	return this.encoder.Float(value)
}

func (this *ObjectToCBEAdapter) writeArray(array []byte) (err error) {
	if err = this.encoder.BeginChunk(uint64(len(array)), true); err != nil {
		return
	}
	if len(array) > 0 {
		var bytesEncoded int
		if bytesEncoded, err = this.encoder.ArrayData([]byte(array)); err != nil {
			return
		}
		if bytesEncoded != len(array) {
			err = fmt.Errorf("Only %v of %v bytes written", bytesEncoded, len(array))
		}
	}
	return
}

func (this *ObjectToCBEAdapter) OnString(value string) (err error) {
	// TODO: Short strings
	if err = this.encoder.BeginString(); err != nil {
		return
	}
	return this.writeArray([]byte(value))
}

func (this *ObjectToCBEAdapter) OnBytes(value []byte) (err error) {
	if err = this.encoder.BeginBytes(); err != nil {
		return
	}
	return this.writeArray(value)
}

func (this *ObjectToCBEAdapter) OnURI(value *url.URL) (err error) {
	if err = this.encoder.BeginURI(); err != nil {
		return
	}
	return this.writeArray([]byte(value.String()))
}

func (this *ObjectToCBEAdapter) OnTime(value time.Time) error {
	return this.encoder.Time(value)
}

func (this *ObjectToCBEAdapter) OnListBegin() error {
	return this.encoder.BeginList()
}

func (this *ObjectToCBEAdapter) OnListEnd() error {
	return this.encoder.EndContainer()
}

func (this *ObjectToCBEAdapter) OnMapBegin() error {
	return this.encoder.BeginMap()
}

func (this *ObjectToCBEAdapter) OnMapEnd() error {
	return this.encoder.EndContainer()
}

func (this *ObjectToCBEAdapter) OnMarker(name uint64) error {
	return this.encoder.BeginMarker()
}

func (this *ObjectToCBEAdapter) OnReference(name uint64) error {
	return this.encoder.BeginReference()
}
