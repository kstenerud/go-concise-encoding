package cbe

import (
	"fmt"
	"time"
)

type PrimitiveDecoderCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnInt(value int64) error
	OnUint(value uint64) error
	OnFloat(value float64) error
	OnTime(value time.Time) error
	OnListBegin() error
	OnListEnd() error
	OnMapBegin() error
	OnMapEnd() error
	OnStringBegin(byteCount uint64) error
	OnStringData(bytes []byte) error
	OnBinaryBegin(byteCount uint64) error
	OnBinaryData(bytes []byte) error
}

type Unmarshaler struct {
	nextValue        interface{}
	containerStack   []interface{}
	currentList      []interface{}
	currentMap       map[interface{}]interface{}
	currentArray     []byte
	currentArrayType arrayType
}

func (this *Unmarshaler) setCurrentContainer() error {
	lastEntry := len(this.containerStack) - 1
	this.currentList = nil
	this.currentMap = nil
	if lastEntry >= 0 {
		container := this.containerStack[lastEntry]
		if list, ok := container.([]interface{}); ok {
			this.currentList = list
		} else {
			this.currentMap = container.(map[interface{}]interface{})
		}
	}
	return nil
}

func (this *Unmarshaler) containerBegin(container interface{}) error {
	this.containerStack = append(this.containerStack, container)
	return this.setCurrentContainer()
}

func (this *Unmarshaler) listBegin() error {
	return this.containerBegin(make([]interface{}, 0))
}

func (this *Unmarshaler) mapBegin() error {
	return this.containerBegin(make(map[interface{}]interface{}))
}

func (this *Unmarshaler) containerEnd() error {
	length := len(this.containerStack)
	this.nextValue = this.containerStack[length-1]
	if length > 0 {
		this.containerStack = this.containerStack[:length-1]
		return this.setCurrentContainer()
	}
	return nil
}

func (this *Unmarshaler) arrayBegin(newArrayType arrayType, length int) error {
	this.currentArray = make([]byte, 0, length)
	this.currentArrayType = newArrayType
	if length == 0 {
		if this.currentArrayType == arrayTypeBinary {
			return this.storeValue(this.currentArray)
		} else {
			return this.storeValue(string(this.currentArray))
		}
	}
	return nil
}

func (this *Unmarshaler) arrayData(data []byte) error {
	this.currentArray = append(this.currentArray, data...)
	if len(this.currentArray) == cap(this.currentArray) {
		if this.currentArrayType == arrayTypeBinary {
			return this.storeValue(this.currentArray)
		} else {
			return this.storeValue(string(this.currentArray))
		}
	}
	return nil
}

func (this *Unmarshaler) storeValue(value interface{}) error {
	if this.currentList != nil {
		this.currentList = append(this.currentList, value)
		return nil
	}

	if this.currentMap != nil {
		if this.nextValue == nil {
			this.nextValue = value
		} else {
			this.currentMap[this.nextValue] = value
			this.nextValue = nil
		}
		return nil
	}

	if this.nextValue != nil {
		return fmt.Errorf("Top level object already exists: %v", this.nextValue)
	}
	this.nextValue = value
	return nil
}

func (this *Unmarshaler) OnNil() error {
	// Ignored
	return nil
}

func (this *Unmarshaler) OnBool(value bool) error {
	return this.storeValue(value)
}

func (this *Unmarshaler) OnInt(value int64) error {
	return this.storeValue(value)
}

func (this *Unmarshaler) OnUint(value uint64) error {
	return this.storeValue(value)
}

func (this *Unmarshaler) OnFloat(value float64) error {
	return this.storeValue(value)
}

func (this *Unmarshaler) OnTime(value time.Time) error {
	return this.storeValue(value)
}

func (this *Unmarshaler) OnListBegin() error {
	return this.listBegin()
}

func (this *Unmarshaler) OnListEnd() error {
	return this.containerEnd()
}

func (this *Unmarshaler) OnMapBegin() error {
	return this.mapBegin()
}

func (this *Unmarshaler) OnMapEnd() error {
	return this.containerEnd()
}

func (this *Unmarshaler) OnStringBegin(byteCount uint64) error {
	return this.arrayBegin(arrayTypeString, int(byteCount))
}

func (this *Unmarshaler) OnStringData(bytes []byte) error {
	return this.arrayData(bytes)
}

func (this *Unmarshaler) OnCommentBegin(byteCount uint64) error {
	return this.arrayBegin(arrayTypeComment, int(byteCount))
}

func (this *Unmarshaler) OnCommentData(bytes []byte) error {
	return this.arrayData(bytes)
}

func (this *Unmarshaler) OnBinaryBegin(byteCount uint64) error {
	return this.arrayBegin(arrayTypeBinary, int(byteCount))
}

func (this *Unmarshaler) OnBinaryData(bytes []byte) error {
	return this.arrayData(bytes)
}

func (this *Unmarshaler) Unmarshaled() interface{} {
	if len(this.containerStack) != 0 {
		return this.containerStack[0]
	}
	return this.nextValue
}

func (this *Unmarshaler) UnmarshaledTo(dest interface{}) interface{} {
	return nil
}
