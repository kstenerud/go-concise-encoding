package cbe

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Nil int

type testCallbacks struct {
	nextValue        interface{}
	containerStack   []interface{}
	currentList      []interface{}
	currentMap       map[interface{}]interface{}
	currentArray     []byte
	currentArrayType arrayType
}

func (callbacks *testCallbacks) setCurrentContainer() {
	lastEntry := len(callbacks.containerStack) - 1
	callbacks.currentList = nil
	callbacks.currentMap = nil
	if lastEntry >= 0 {
		container := callbacks.containerStack[lastEntry]
		if list, ok := container.([]interface{}); ok {
			callbacks.currentList = list
		} else {
			callbacks.currentMap = container.(map[interface{}]interface{})
		}
	}
}

func (callbacks *testCallbacks) containerBegin(container interface{}) {
	callbacks.containerStack = append(callbacks.containerStack, container)
	callbacks.setCurrentContainer()
}

func (callbacks *testCallbacks) listBegin() {
	callbacks.containerBegin(new([]interface{}))
}

func (callbacks *testCallbacks) mapBegin() {
	callbacks.containerBegin(new(map[interface{}]interface{}))
}

func (callbacks *testCallbacks) containerEnd() {
	length := len(callbacks.containerStack)
	if length > 0 {
		callbacks.containerStack = callbacks.containerStack[:length-1]
		callbacks.setCurrentContainer()
	}
}

func (callbacks *testCallbacks) arrayBegin(newArrayType arrayType, length int) {
	callbacks.currentArray = make([]byte, 0, length)
	callbacks.currentArrayType = newArrayType
	if length == 0 {
		if callbacks.currentArrayType == arrayTypeBinary {
			callbacks.storeValue(callbacks.currentArray)
		} else {
			callbacks.storeValue(string(callbacks.currentArray))
		}
	}
}

func (callbacks *testCallbacks) arrayData(data []byte) {
	callbacks.currentArray = append(callbacks.currentArray, data...)
	if len(callbacks.currentArray) == cap(callbacks.currentArray) {
		if callbacks.currentArrayType == arrayTypeBinary {
			callbacks.storeValue(callbacks.currentArray)
		} else {
			callbacks.storeValue(string(callbacks.currentArray))
		}
	}
}

func (callbacks *testCallbacks) storeValue(value interface{}) {
	if callbacks.currentList != nil {
		callbacks.currentList = append(callbacks.currentList, value)
		return
	}

	if callbacks.currentMap != nil {
		if callbacks.nextValue == nil {
			callbacks.nextValue = value
		} else {
			callbacks.currentMap[callbacks.nextValue] = value
			callbacks.nextValue = nil
		}
		return
	}

	if callbacks.nextValue != nil {
		panic(fmt.Errorf("Top level object already exists: %v", callbacks.nextValue))
	}
	callbacks.nextValue = value
}

func (callbacks *testCallbacks) getValue() interface{} {
	if len(callbacks.containerStack) != 0 {
		return callbacks.containerStack[0]
	}
	return callbacks.nextValue
}

func (callbacks *testCallbacks) OnNil() error {
	callbacks.storeValue(new(Nil))
	return nil
}

func (callbacks *testCallbacks) OnBool(value bool) error {
	callbacks.storeValue(value)
	return nil
}

func (callbacks *testCallbacks) OnInt(value int64) error {
	callbacks.storeValue(value)
	return nil
}

func (callbacks *testCallbacks) OnUint(value uint64) error {
	callbacks.storeValue(value)
	return nil
}

func (callbacks *testCallbacks) OnFloat(value float64) error {
	callbacks.storeValue(value)
	return nil
}

func (callbacks *testCallbacks) OnTime(value time.Time) error {
	callbacks.storeValue(value)
	return nil
}

func (callbacks *testCallbacks) OnListBegin() error {
	callbacks.listBegin()
	return nil
}

func (callbacks *testCallbacks) OnListEnd() error {
	callbacks.containerEnd()
	return nil
}

func (callbacks *testCallbacks) OnMapBegin() error {
	callbacks.mapBegin()
	return nil
}

func (callbacks *testCallbacks) OnMapEnd() error {
	callbacks.containerEnd()
	return nil
}

func (callbacks *testCallbacks) OnStringBegin(byteCount uint64) error {
	callbacks.arrayBegin(arrayTypeString, int(byteCount))
	return nil
}

func (callbacks *testCallbacks) OnStringData(bytes []byte) error {
	callbacks.arrayData(bytes)
	return nil
}

func (callbacks *testCallbacks) OnCommentBegin(byteCount uint64) error {
	callbacks.arrayBegin(arrayTypeComment, int(byteCount))
	return nil
}

func (callbacks *testCallbacks) OnCommentData(bytes []byte) error {
	callbacks.arrayData(bytes)
	return nil
}

func (callbacks *testCallbacks) OnBinaryBegin(byteCount uint64) error {
	callbacks.arrayBegin(arrayTypeBinary, int(byteCount))
	return nil
}

func (callbacks *testCallbacks) OnBinaryData(bytes []byte) error {
	callbacks.arrayData(bytes)
	return nil
}

func assertDecoded(t *testing.T, encoded []byte, expected interface{}) {
	callbacks := new(testCallbacks)
	decoder := NewCbeDecoder(9, callbacks)
	err := decoder.Feed(encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	err = decoder.End()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	actual := callbacks.getValue()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected [%v], actual [%v]", expected, actual)
	}
}
