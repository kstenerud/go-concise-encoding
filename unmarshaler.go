package cbe

import (
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
	topObject        interface{}
	containerStack   []interface{}
	mapKeyStack      []interface{}
	currentList      []interface{}
	currentMap       map[interface{}]interface{}
	currentMapKey    interface{}
	currentArray     []byte
	currentArrayType arrayType
}

func (this *Unmarshaler) stackCurrentContainer() {
	if this.currentList != nil {
		this.containerStack = append(this.containerStack, this.currentList)
		this.currentList = nil
	} else {
		this.containerStack = append(this.containerStack, this.currentMap)
		this.currentMap = nil
	}
	this.mapKeyStack = append(this.mapKeyStack, this.currentMapKey)
	this.currentMapKey = nil
}

func (this *Unmarshaler) listBegin() {
	this.stackCurrentContainer()
	this.currentList = make([]interface{}, 0)
}

func (this *Unmarshaler) mapBegin() {
	this.stackCurrentContainer()
	this.currentMap = make(map[interface{}]interface{})
}

func (this *Unmarshaler) containerEnd() {
	var oldContainer interface{}
	if this.currentList != nil {
		oldContainer = this.currentList
	} else {
		oldContainer = this.currentMap
	}

	length := len(this.containerStack)
	this.currentMapKey = this.mapKeyStack[length-1]
	container := this.containerStack[length-1]
	this.containerStack = this.containerStack[:length-1]

	if container != nil {
		if list, ok := container.([]interface{}); ok {
			this.currentList = list
			this.currentMap = nil
		} else {
			this.currentMap = container.(map[interface{}]interface{})
			this.currentList = nil
		}
	}

	this.storeValue(oldContainer)
}

func (this *Unmarshaler) arrayBegin(newArrayType arrayType, length int) {
	this.currentArray = make([]byte, 0, length)
	this.currentArrayType = newArrayType
	if length == 0 {
		if this.currentArrayType == arrayTypeBinary {
			this.storeValue(this.currentArray)
		} else {
			this.storeValue(string(this.currentArray))
		}
	}
}

func (this *Unmarshaler) arrayData(data []byte) {
	this.currentArray = append(this.currentArray, data...)
	if len(this.currentArray) == cap(this.currentArray) {
		if this.currentArrayType == arrayTypeBinary {
			this.storeValue(this.currentArray)
		} else {
			this.storeValue(string(this.currentArray))
		}
	}
}

func (this *Unmarshaler) storeValue(value interface{}) {
	this.topObject = value

	if this.currentList != nil {
		this.currentList = append(this.currentList, value)
	} else if this.currentMap != nil {
		if this.currentMapKey == nil {
			this.currentMapKey = value
		} else {
			this.currentMap[this.currentMapKey] = value
			this.currentMapKey = nil
		}
	}
}

func (this *Unmarshaler) OnNil() error {
	this.storeValue(nil)
	return nil
}

func (this *Unmarshaler) OnBool(value bool) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnInt(value int64) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnUint(value uint64) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnFloat(value float64) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnTime(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnListBegin() error {
	this.listBegin()
	return nil
}

func (this *Unmarshaler) OnListEnd() error {
	this.containerEnd()
	return nil
}

func (this *Unmarshaler) OnMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *Unmarshaler) OnMapEnd() error {
	this.containerEnd()
	return nil
}

func (this *Unmarshaler) OnStringBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeString, int(byteCount))
	return nil
}

func (this *Unmarshaler) OnStringData(bytes []byte) error {
	this.arrayData(bytes)
	return nil
}

func (this *Unmarshaler) OnCommentBegin(byteCount uint64) error {
	// Ignored
	return nil
}

func (this *Unmarshaler) OnCommentData(bytes []byte) error {
	// Ignored
	return nil
}

func (this *Unmarshaler) OnBinaryBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeBinary, int(byteCount))
	return nil
}

func (this *Unmarshaler) OnBinaryData(bytes []byte) error {
	this.arrayData(bytes)
	return nil
}

func (this *Unmarshaler) Unmarshaled() interface{} {
	return this.topObject
}

func (this *Unmarshaler) UnmarshaledTo(dest interface{}) interface{} {
	return nil
}
