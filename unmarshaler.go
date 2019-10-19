package cbe

import (
	"time"
)

type PrimitiveDecoderCallbacks interface {
	OnNil() error
	OnBool(value bool) error
	OnIntPositive(value uint64) error
	OnIntNegative(value uint64) error
	OnFloat(value float64) error
	OnDate(value time.Time) error
	OnTime(value time.Time) error
	OnTimestamp(value time.Time) error
	OnListBegin() error
	OnOrderedMapBegin() error
	OnUnorderedMapBegin() error
	OnMetadataMapBegin() error
	OnContainerEnd() error
	OnStringBegin(byteCount uint64) error
	OnCommentBegin(byteCount uint64) error
	OnURIBegin(byteCount uint64) error
	OnBytesBegin(byteCount uint64) error
	OnArrayData(bytes []byte) error
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
		if this.currentArrayType == arrayTypeBytes {
			this.storeValue(this.currentArray)
		} else {
			this.storeValue(string(this.currentArray))
		}
	}
}

func (this *Unmarshaler) arrayData(data []byte) {
	this.currentArray = append(this.currentArray, data...)
	if len(this.currentArray) == cap(this.currentArray) {
		if this.currentArrayType == arrayTypeBytes {
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

func (this *Unmarshaler) OnIntPositive(value uint64) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnIntNegative(value uint64) error {
	this.storeValue(-int64(value))
	return nil
}

func (this *Unmarshaler) OnFloat(value float64) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnDate(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnTime(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnTimestamp(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *Unmarshaler) OnListBegin() error {
	this.listBegin()
	return nil
}

func (this *Unmarshaler) OnOrderedMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *Unmarshaler) OnUnorderedMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *Unmarshaler) OnMetadataMapBegin() error {
	// Ignored
	return nil
}

func (this *Unmarshaler) OnContainerEnd() error {
	this.containerEnd()
	return nil
}

func (this *Unmarshaler) OnBytesBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeBytes, int(byteCount))
	return nil
}

func (this *Unmarshaler) OnStringBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeString, int(byteCount))
	return nil
}

func (this *Unmarshaler) OnCommentBegin(byteCount uint64) error {
	// Ignored
	return nil
}

func (this *Unmarshaler) OnURIBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeURI, int(byteCount))
	return nil
}

func (this *Unmarshaler) OnArrayData(bytes []byte) error {
	this.arrayData(bytes)
	return nil
}

func (this *Unmarshaler) Unmarshaled() interface{} {
	return this.topObject
}

func (this *Unmarshaler) UnmarshaledTo(dest interface{}) interface{} {
	return nil
}
