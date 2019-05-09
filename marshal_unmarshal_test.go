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
	encoder := NewCbeEncoder(100)
	Marshal(encoder, expected)
	document := encoder.Encoded()
	unmarshaler := new(Unmarshaler)
	decoder := NewCbeDecoder(100, unmarshaler)
	decoder.Decode(document)
	actual := unmarshaler.Unmarshaled()

	if !DeepEquivalence(actual, expected) {
		t.Errorf("Expected %t %v, actual %t %v", expected, expected, actual, actual)
	}
}

func TestNil(t *testing.T) {
	assertMarshalUnmarshal(t, nil)
}

func TestBool(t *testing.T) {
	assertMarshalUnmarshal(t, false)
	assertMarshalUnmarshal(t, true)
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

func TestTimeNow(t *testing.T) {
	assertMarshalUnmarshal(t, time.Now().UTC())
}

func TestEmptyString(t *testing.T) {
	assertMarshalUnmarshal(t, "")
}

func TestShortString(t *testing.T) {
	assertMarshalUnmarshal(t, "This is a test")
}

func TestLongString(t *testing.T) {
	assertMarshalUnmarshal(t, "This is a longer string test that goes beyond 15 characters.")
}

func TestBytesSlice(t *testing.T) {
	assertMarshalUnmarshal(t, []byte{1, 2, 3, 4, 5})
}

func TestBytesArray(t *testing.T) {
	array := [...]byte{1, 2}
	assertMarshalUnmarshal(t, array)
}
