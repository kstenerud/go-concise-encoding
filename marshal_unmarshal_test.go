package cbe

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

var chars = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
}

func generateString(length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(chars[i%len(chars)])
	}
	return result.String()
}

func generateBinary(length int) []byte {
	return []byte(generateString(length))
}

func assertMarshaled(t *testing.T, value interface{}, expected []byte) {
	encoder := NewCbeEncoder(100)
	Marshal(encoder, value)
	actual := encoder.Encoded()
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func assertMarshalUnmarshal(t *testing.T, expected interface{}) {
	// Coerce to 64 bit values since that's what the unmarshaler will return.
	// reflect.DeepEqual fails if types are not the same
	rv := reflect.ValueOf(expected)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		iv := rv.Int()
		if iv >= 0 {
			expected = uint64(iv)
		} else {
			expected = iv
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		expected = rv.Uint()
	case reflect.Float32:
		expected = rv.Float()
	}
	encoder := NewCbeEncoder(100)
	Marshal(encoder, expected)
	document := encoder.Encoded()
	unmarshaler := new(Unmarshaler)
	decoder := NewCbeDecoder(100, unmarshaler)
	decoder.Decode(document)
	actual := unmarshaler.Unmarshaled()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %t %v, actual %t %v", expected, expected, actual, actual)
	}
}

func TestNil(t *testing.T) {
	assertMarshalUnmarshal(t, nil)
}

func Test0(t *testing.T) {
	assertMarshalUnmarshal(t, 0)
}

func TestN1(t *testing.T) {
	assertMarshalUnmarshal(t, -1)
}

func Test200(t *testing.T) {
	assertMarshalUnmarshal(t, 200)
}

func Test2000(t *testing.T) {
	assertMarshalUnmarshal(t, 2000)
}

func Test1_5(t *testing.T) {
	assertMarshalUnmarshal(t, 1.5)
}

func Test1_0123(t *testing.T) {
	assertMarshalUnmarshal(t, 1.0123)
}

func Test1_now(t *testing.T) {
	assertMarshalUnmarshal(t, time.Now().UTC())
}
