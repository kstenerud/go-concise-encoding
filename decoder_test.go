package cbe

import (
	"fmt"
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
	decoder := NewDecoder(9, callbacks)
	err := decoder.Feed(encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	actual := callbacks.getValue()
	if actual != expected {
		t.Errorf("Expected [%v], actual [%v]", expected, actual)
	}
}

func TestDecodeSmallInt(t *testing.T) {
	assertDecoded(t, []byte{0x00}, uint64(0))
	assertDecoded(t, []byte{0x01}, uint64(1))
	assertDecoded(t, []byte{105}, uint64(105))
	assertDecoded(t, []byte{0xff}, int64(-1))
	assertDecoded(t, []byte{0x96}, int64(-106))
}

func TestDecodeInt8(t *testing.T) {
	assertDecoded(t, []byte{0x6a, 0x6a}, uint64(106))
	assertDecoded(t, []byte{0x6a, 0xff}, uint64(0xff))
	assertDecoded(t, []byte{0x7a, 0xff}, int64(-0xff))
}

func TestDecodeInt16(t *testing.T) {
	assertDecoded(t, []byte{0x6b, 0x00, 0x01}, uint64(0x0100))
	assertDecoded(t, []byte{0x6b, 0xff, 0xff}, uint64(0xffff))
	assertDecoded(t, []byte{0x7b, 0x00, 0x01}, int64(-0x0100))
	assertDecoded(t, []byte{0x7b, 0xff, 0xff}, int64(-0xffff))
}

func TestDecodeInt32(t *testing.T) {
	assertDecoded(t, []byte{0x6c, 0x00, 0x00, 0x00, 0x01}, uint64(0x01000000))
	assertDecoded(t, []byte{0x6c, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffff))
	assertDecoded(t, []byte{0x7c, 0x00, 0x00, 0x00, 0x01}, int64(-0x01000000))
	assertDecoded(t, []byte{0x7c, 0xff, 0xff, 0xff, 0xff}, int64(-0xffffffff))
}

func TestDecodeInt64(t *testing.T) {
	assertDecoded(t, []byte{0x6d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, uint64(0x0100000000000000))
	assertDecoded(t, []byte{0x6d, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, uint64(0xffffffffffffffff))
	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, int64(-0x0100000000000000))
	assertDecoded(t, []byte{0x7d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, int64(-0x8000000000000000))
}

func TestDecodeString(t *testing.T) {
	assertDecoded(t, []byte{0x80}, "")
}
