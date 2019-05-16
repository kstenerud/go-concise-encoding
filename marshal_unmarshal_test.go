package cbe

import (
	"bytes"
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

func generateBytes(length int) []byte {
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
	assertMarshalUnmarshalProduces(t, expected, expected)
}

func assertMarshalUnmarshalProduces(t *testing.T, input interface{}, expected interface{}) {
	encoder := NewCbeEncoder(100)
	if err := Marshal(encoder, input); err != nil {
		t.Errorf("Unexpected error while marshling: %v", err)
		return
	}
	document := encoder.Encoded()
	unmarshaler := new(Unmarshaler)
	decoder := NewCbeDecoder(100, unmarshaler)
	if err := decoder.Decode(document); err != nil {
		t.Errorf("Unexpected error while decoding: %v", err)
		return
	}
	actual := unmarshaler.Unmarshaled()

	if !DeepEquivalence(actual, expected) {
		t.Errorf("Expected %t: <%v>, actual %t: <%v>", expected, expected, actual, actual)
	}
}

func TestMarshalUnmarshalNil(t *testing.T) {
	assertMarshalUnmarshal(t, nil)
}

func TestMarshalUnmarshalBool(t *testing.T) {
	assertMarshalUnmarshal(t, false)
	assertMarshalUnmarshal(t, true)
}

func TestMarshalUnmarshal0(t *testing.T) {
	assertMarshalUnmarshal(t, 0)
}

func TestMarshalUnmarshalN1(t *testing.T) {
	assertMarshalUnmarshal(t, -1)
}

func TestMarshalUnmarshal200(t *testing.T) {
	assertMarshalUnmarshal(t, 200)
}

func TestMarshalUnmarshal2000(t *testing.T) {
	assertMarshalUnmarshal(t, 2000)
}

func TestMarshalUnmarshal1_5(t *testing.T) {
	assertMarshalUnmarshal(t, 1.5)
}

func TestMarshalUnmarshal1_0123(t *testing.T) {
	assertMarshalUnmarshal(t, 1.0123)
}

func TestMarshalUnmarshalTimeNow(t *testing.T) {
	assertMarshalUnmarshal(t, time.Now().UTC())
}

func TestMarshalUnmarshalEmptyString(t *testing.T) {
	assertMarshalUnmarshal(t, "")
}

func TestMarshalUnmarshalShortString(t *testing.T) {
	assertMarshalUnmarshal(t, "This is a test")
}

func TestMarshalUnmarshalLongString(t *testing.T) {
	assertMarshalUnmarshal(t, "This is a longer string test that goes beyond 15 characters.")
}

func TestMarshalUnmarshalBytesSlice(t *testing.T) {
	assertMarshalUnmarshal(t, []byte{1, 2, 3, 4, 5})
}

func TestMarshalUnmarshalBytesArray(t *testing.T) {
	array := [...]byte{1, 2}
	assertMarshalUnmarshal(t, array)
}

func TestMarshalUnmarshalLongBytesArray(t *testing.T) {
	array := [...]byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}
	assertMarshalUnmarshal(t, array)
}

func TestMarshalUnmarshalLongerBytesArray(t *testing.T) {
	array := make([]byte, 60000)
	assertMarshalUnmarshal(t, array)
}

func TestMarshalUnmarshalAllBasicTypes(t *testing.T) {
	assertMarshalUnmarshal(t, []interface{}{
		1, 250, 1000, 100000, 10000000000,
		-1, -250, -1000, -100000000000,
		1.5, 1.9582384465,
		"string", []byte{10, 11, 12},
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 9999999, time.UTC),
	})
}

func TestMarshalUnmarshalList(t *testing.T) {
	assertMarshalUnmarshal(t, []interface{}{
		1, 2, "test", 4,
	})
}

func TestMarshalUnmarshalListNil(t *testing.T) {
	assertMarshalUnmarshal(t, []interface{}{
		1, nil, "test", 4,
	})
}

func TestMarshalUnmarshalMap(t *testing.T) {
	assertMarshalUnmarshal(t, map[interface{}]interface{}{
		"a": 1,
		2:   "b",
	})
}

func TestMarshalUnmarshalMapNil(t *testing.T) {
	assertMarshalUnmarshal(t, map[interface{}]interface{}{
		"a": nil,
	})
}

func TestMarshalUnmarshalMapMap(t *testing.T) {
	assertMarshalUnmarshal(t, map[interface{}]interface{}{
		1: 2,
		"deep-map": map[interface{}]interface{}{
			3: 1000,
		},
	})
}

func TestMarshalUnmarshalListList(t *testing.T) {
	assertMarshalUnmarshal(t, []interface{}{
		1, []interface{}{
			2, "test", 4,
		},
	})
}

func TestMarshalUnmarshalListMap(t *testing.T) {
	assertMarshalUnmarshal(t, []interface{}{
		1, map[interface{}]interface{}{
			2: 3, 4: "blah",
		},
	})
}

func TestMarshalUnmarshalMapList(t *testing.T) {
	assertMarshalUnmarshal(t, map[interface{}]interface{}{
		1: 2,
		"deep-list": []interface{}{
			2, "some list entry",
		},
	})
}

type MyTestStruct struct {
	IntValue    int
	FloatValue  float32
	StringValue string
	ByteValue   []byte
}

func TestMarshalUnmarshalStructZero(t *testing.T) {
	structValue := new(MyTestStruct)
	assertMarshalUnmarshalProduces(t, structValue, map[interface{}]interface{}{
		"IntValue":    structValue.IntValue,
		"FloatValue":  structValue.FloatValue,
		"StringValue": structValue.StringValue,
		"ByteValue":   structValue.ByteValue,
	})
}

func TestMarshalUnmarshalStruct(t *testing.T) {
	structValue := new(MyTestStruct)
	structValue.IntValue = 4000
	structValue.FloatValue = 2.5
	structValue.StringValue = "test"
	structValue.ByteValue = []byte{0x00, 0x01, 0x02}
	assertMarshalUnmarshalProduces(t, structValue, map[interface{}]interface{}{
		"IntValue":    structValue.IntValue,
		"FloatValue":  structValue.FloatValue,
		"StringValue": structValue.StringValue,
		"ByteValue":   structValue.ByteValue,
	})
}
