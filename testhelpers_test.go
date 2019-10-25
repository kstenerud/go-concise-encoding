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

func getPanicContents(function func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	function()
	return recovered
}

func assertPanics(t *testing.T, function func()) {
	if getPanicContents(function) == nil {
		t.Errorf("Should have panicked but didn't")
	}
}

func assertDoesNotPanic(t *testing.T, function func()) {
	if result := getPanicContents(function); result != nil {
		t.Errorf("Should not have panicked, but did: %v", result)
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

var NilValue Nil = Nil(-1234)

type testCallbacks struct {
	nextValue        interface{}
	containerStack   []interface{}
	currentList      []interface{}
	currentMap       map[interface{}]interface{}
	currentArray     []byte
	currentArrayType arrayType
}

func (this *testCallbacks) setCurrentContainer() {
	lastEntry := len(this.containerStack) - 1
	this.currentList = nil
	this.currentMap = nil
	if lastEntry >= 0 {
		container := this.containerStack[lastEntry]
		switch container.(type) {
		case []interface{}:
			this.currentList = container.([]interface{})
		case *[]interface{}:
			this.currentList = *(container.(*[]interface{}))
		case map[interface{}]interface{}:
			this.currentMap = container.(map[interface{}]interface{})
		case *map[interface{}]interface{}:
			this.currentMap = *(container.(*map[interface{}]interface{}))
		default:
			panic(fmt.Errorf("Unknown container type: %v", container))
		}
	}
}

func (this *testCallbacks) containerBegin(container interface{}) {
	this.containerStack = append(this.containerStack, container)
	this.setCurrentContainer()
}

func (this *testCallbacks) listBegin() {
	this.containerBegin(make([]interface{}, 0))
}

func (this *testCallbacks) mapBegin() {
	this.containerBegin(make(map[interface{}]interface{}))
}

func (this *testCallbacks) containerEnd() {
	var item interface{}

	if this.currentList != nil {
		item = this.currentList
		this.currentList = nil
	} else {
		item = this.currentMap
		this.currentMap = nil
	}
	length := len(this.containerStack)
	if length > 0 {
		this.containerStack = this.containerStack[:length-1]
		this.setCurrentContainer()
	}
	this.storeValue(item)
}

func (this *testCallbacks) arrayBegin(newArrayType arrayType, length int) {
	this.currentArray = make([]byte, 0, length)
	// fmt.Printf("Make array cap %v: %v\n", length, cap(this.currentArray))
	this.currentArrayType = newArrayType
	if length == 0 {
		this.arrayEnd()
	}
}

func (this *testCallbacks) arrayData(data []byte) {
	this.currentArray = append(this.currentArray, data...)
	if len(this.currentArray) == cap(this.currentArray) {
		this.arrayEnd()
	}
}

func (this *testCallbacks) arrayEnd() {
	array := this.currentArray
	this.currentArray = nil
	if this.currentArrayType == arrayTypeBytes {
		this.storeValue(array)
	} else {
		this.storeValue(string(array))
	}
}

func (this *testCallbacks) storeValue(value interface{}) {
	if this.currentList != nil {
		this.currentList = append(this.currentList, value)
		return
	}

	if this.currentMap != nil {
		if this.nextValue == nil {
			this.nextValue = value
		} else {
			this.currentMap[this.nextValue] = value
			this.nextValue = nil
		}
		return
	}

	if this.nextValue != nil {
		panic(fmt.Errorf("Top level object already exists: %v", this.nextValue))
	}
	this.nextValue = value
}

func (this *testCallbacks) getValue() interface{} {
	if len(this.containerStack) != 0 {
		return this.containerStack[0]
	}
	return this.nextValue
}

func (this *testCallbacks) OnNil() error {
	this.storeValue(NilValue)
	return nil
}

func (this *testCallbacks) OnBool(value bool) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnPositiveInt(value uint64) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnNegativeInt(value uint64) error {
	this.storeValue(-int64(value))
	return nil
}

func (this *testCallbacks) OnFloat(value float64) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnDate(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnTime(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnTimestamp(value time.Time) error {
	this.storeValue(value)
	return nil
}

func (this *testCallbacks) OnListBegin() error {
	this.listBegin()
	return nil
}

func (this *testCallbacks) OnOrderedMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *testCallbacks) OnUnorderedMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *testCallbacks) OnMetadataMapBegin() error {
	this.mapBegin()
	return nil
}

func (this *testCallbacks) OnContainerEnd() error {
	this.containerEnd()
	return nil
}

func (this *testCallbacks) OnStringBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeString, int(byteCount))
	return nil
}

func (this *testCallbacks) OnCommentBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeComment, int(byteCount))
	return nil
}

func (this *testCallbacks) OnURIBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeURI, int(byteCount))
	return nil
}

func (this *testCallbacks) OnBytesBegin(byteCount uint64) error {
	this.arrayBegin(arrayTypeBytes, int(byteCount))
	return nil
}

func (this *testCallbacks) OnArrayData(bytes []byte) error {
	this.arrayData(bytes)
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

func assertEncoded(t *testing.T, containerType ContainerType, function func(*CbeEncoder) error, expected []byte) {
	encoder := NewCbeEncoder(containerType, nil, 100)
	err := function(encoder)
	if err != nil {
		t.Fatal(err)
	}
	actual := encoder.EncodedBytes()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

// Marshal / Unmarshal

func assertMarshaled(t *testing.T, containerType ContainerType, value interface{}, expected []byte) {
	encoder := NewCbeEncoder(containerType, nil, 100)
	Marshal(encoder, containerType, value)
	actual := encoder.EncodedBytes()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func assertEncodesToExternalBuffer(t *testing.T, containerType ContainerType, value interface{}, bufferSize int) {
	buffer := make([]byte, bufferSize)
	encoder := NewCbeEncoder(containerType, buffer, 100)
	if err := Marshal(encoder, containerType, value); err != nil {
		t.Errorf("Unexpected error while marshling: %v", err)
		return
	}

	encoder2 := NewCbeEncoder(containerType, nil, 100)
	Marshal(encoder2, containerType, value)
	expected := encoder2.EncodedBytes()
	if !bytes.Equal(buffer, expected) {
		t.Errorf("Expected %v, actual %v", expected, buffer)
	}
}

func assertFailsEncodingToExternalBuffer(t *testing.T, containerType ContainerType, value interface{}, bufferSize int) {
	buffer := make([]byte, bufferSize)
	encoder := NewCbeEncoder(containerType, buffer, 100)
	assertPanics(t, func() {
		Marshal(encoder, containerType, value)
	})
}

func assertMarshaledSize(t *testing.T, containerType ContainerType, value interface{}, expectedSize int) {
	actualSize := EncodedSize(containerType, value)
	if actualSize != expectedSize {
		t.Errorf("Expected size %v but got %v", expectedSize, actualSize)
	}
}

func assertMarshalUnmarshal(t *testing.T, containerType ContainerType, expected interface{}) {
	assertMarshalUnmarshalProduces(t, containerType, expected, expected)
}

func assertMarshalUnmarshalProduces(t *testing.T, containerType ContainerType, input interface{}, expected interface{}) {
	encoder := NewCbeEncoder(containerType, nil, 100)
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

func ShortCircuit(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
