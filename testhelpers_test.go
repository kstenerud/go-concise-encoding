package cbe

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// General

var stringGeneratorChars = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
}

func generateString(length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(stringGeneratorChars[i%len(stringGeneratorChars)])
	}
	return result.String()
}

func generateBytes(length int) []byte {
	return []byte(generateString(length))
}

func testPanics(function func()) (didPanic bool) {
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
		}
	}()
	didPanic = false
	function()
	return didPanic
}

func assertPanics(t *testing.T, function func()) {
	if !testPanics(function) {
		t.Errorf("Should have panicked but didn't")
	}
}

func assertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func assertFailure(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Unexpected success")
	}
}

// Decoder

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
		switch container.(type) {
		case []interface{}:
			callbacks.currentList = container.([]interface{})
		case *[]interface{}:
			callbacks.currentList = *(container.(*[]interface{}))
		case map[interface{}]interface{}:
			callbacks.currentMap = container.(map[interface{}]interface{})
		case *map[interface{}]interface{}:
			callbacks.currentMap = *(container.(*map[interface{}]interface{}))
		default:
			panic(fmt.Errorf("Unknown container type: %v", container))
		}
	}
}

func (callbacks *testCallbacks) containerBegin(container interface{}) {
	callbacks.containerStack = append(callbacks.containerStack, container)
	callbacks.setCurrentContainer()
}

func (callbacks *testCallbacks) listBegin() {
	callbacks.containerBegin(make([]interface{}, 0))
}

func (callbacks *testCallbacks) mapBegin() {
	callbacks.containerBegin(make(map[interface{}]interface{}))
}

func (callbacks *testCallbacks) containerEnd() {
	var item interface{}

	if callbacks.currentList != nil {
		item = callbacks.currentList
		callbacks.currentList = nil
	} else {
		item = callbacks.currentMap
		callbacks.currentMap = nil
	}
	length := len(callbacks.containerStack)
	if length > 0 {
		callbacks.containerStack = callbacks.containerStack[:length-1]
		callbacks.setCurrentContainer()
	}
	callbacks.storeValue(item)
}

func (callbacks *testCallbacks) arrayBegin(newArrayType arrayType, length int) {
	callbacks.currentArray = make([]byte, 0, length)
	callbacks.currentArrayType = newArrayType
	if length == 0 {
		callbacks.arrayEnd()
	}
}

func (callbacks *testCallbacks) arrayData(data []byte) {
	callbacks.currentArray = append(callbacks.currentArray, data...)
	if len(callbacks.currentArray) == cap(callbacks.currentArray) {
		callbacks.arrayEnd()
	}
}

func (callbacks *testCallbacks) arrayEnd() {
	array := callbacks.currentArray
	callbacks.currentArray = nil
	if callbacks.currentArrayType == arrayTypeBytes {
		callbacks.storeValue(array)
	} else {
		callbacks.storeValue(string(array))
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

func (callbacks *testCallbacks) OnBytesBegin(byteCount uint64) error {
	callbacks.arrayBegin(arrayTypeBytes, int(byteCount))
	return nil
}

func (callbacks *testCallbacks) OnBytesData(bytes []byte) error {
	callbacks.arrayData(bytes)
	return nil
}

func decodeDocument(containerType ContainerType, maxDepth int, encoded []byte) (result interface{}, err error) {
	callbacks := new(testCallbacks)
	decoder := NewCbeDecoder(containerType, maxDepth, callbacks)
	if err := decoder.Feed(encoded); err != nil {
		return nil, err
	}
	if err := decoder.End(); err != nil {
		return nil, err
	}
	result = callbacks.getValue()
	return result, err
}

func decodeWithBufferSize(maxDepth int, encoded []byte, bufferSize int) (result interface{}, err error) {
	unmarshaler := new(Unmarshaler)
	decoder := NewCbeDecoder(ContainerTypeNone, maxDepth, unmarshaler)
	for offset := 0; offset < len(encoded); offset += bufferSize {
		end := offset + bufferSize
		if end > len(encoded) {
			end = len(encoded)
		}
		if err := decoder.Feed(encoded[offset:end]); err != nil {
			return nil, err
		}
	}
	if err := decoder.End(); err != nil {
		return nil, err
	}
	result = unmarshaler.Unmarshaled()
	return result, err
}

func tryDecode(maxDepth int, encoded []byte) error {
	_, err := decodeDocument(ContainerTypeNone, maxDepth, encoded)
	return err
}

func assertDecoded(t *testing.T, containerType ContainerType, encoded []byte, expected interface{}) {
	actual, err := decodeDocument(containerType, 100, encoded)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	if !DeepEquivalence(actual, expected) {
		t.Errorf("Expected [%v], actual [%v]", expected, actual)
	}
}

func assertDecodedPiecemeal(t *testing.T, encoded []byte, minBufferSize int, maxBufferSize int, expected interface{}) {
	for i := minBufferSize; i < maxBufferSize; i++ {
		actual, err := decodeWithBufferSize(100, encoded, i)
		if err != nil {
			t.Errorf("Error: %v", err)
			return
		}
		if !DeepEquivalence(actual, expected) {
			t.Errorf("Expected [%v], actual [%v]", expected, actual)
		}
	}
}

// Encoder

func assertEncoded(t *testing.T, containerType ContainerType, function func(*CbeEncoder), expected []byte) {
	encoder := NewCbeEncoder(containerType, 100)
	function(encoder)
	actual := encoder.EncodedBytes()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

// Marshal / Unmarshal

func assertMarshaled(t *testing.T, containerType ContainerType, value interface{}, expected []byte) {
	encoder := NewCbeEncoder(containerType, 100)
	Marshal(encoder, containerType, value)
	actual := encoder.EncodedBytes()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func assertMarshalUnmarshal(t *testing.T, containerType ContainerType, expected interface{}) {
	assertMarshalUnmarshalProduces(t, containerType, expected, expected)
}

func assertMarshalUnmarshalProduces(t *testing.T, containerType ContainerType, input interface{}, expected interface{}) {
	encoder := NewCbeEncoder(containerType, 100)
	if err := Marshal(encoder, containerType, input); err != nil {
		t.Errorf("Unexpected error while marshling: %v", err)
		return
	}
	document := encoder.EncodedBytes()
	unmarshaler := new(Unmarshaler)
	decoder := NewCbeDecoder(containerType, 100, unmarshaler)
	if err := decoder.Decode(document); err != nil {
		t.Errorf("Unexpected error while decoding: %v", err)
		return
	}
	actual := unmarshaler.Unmarshaled()

	if !DeepEquivalence(actual, expected) {
		t.Errorf("Expected %t: <%v>, actual %t: <%v>", expected, expected, actual, actual)
	}
}
