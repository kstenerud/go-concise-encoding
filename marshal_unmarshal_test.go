package cbe

import (
	"net/url"
	"testing"
	"time"
)

func TestMarshalUnmarshalNil(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, nil)
}

func TestMarshalUnmarshalBool(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, false)
	assertMarshalUnmarshal(t, ContainerTypeNone, true)
}

func TestMarshalUnmarshal0(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 0)
}

func TestMarshalUnmarshal1(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 1)
	assertMarshalUnmarshal(t, ContainerTypeNone, -1)
}

func TestMarshalUnmarshal200(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 200)
	assertMarshalUnmarshal(t, ContainerTypeNone, -200)
}

func TestMarshalUnmarshal2000(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 2000)
	assertMarshalUnmarshal(t, ContainerTypeNone, -2000)
}

func TestMarshalUnmarshal10000(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 10000)
	assertMarshalUnmarshal(t, ContainerTypeNone, -10000)
}

func TestMarshalUnmarshal100000(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 100000)
	assertMarshalUnmarshal(t, ContainerTypeNone, -100000)
}

func TestMarshalUnmarshal10000000000(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 10000000000)
	assertMarshalUnmarshal(t, ContainerTypeNone, -10000000000)
}

func TestMarshalUnmarshal100000000000(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 100000000000)
	assertMarshalUnmarshal(t, ContainerTypeNone, -100000000000)
}

func TestMarshalUnmarshal1_5(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 1.5)
}

func TestMarshalUnmarshal1_0123(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, 1.0123)
}

func TestMarshalUnmarshalTimeNow(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, time.Now().UTC())
}

func TestMarshalUnmarshalEmptyString(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, "")
}

func TestMarshalUnmarshalShortString(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, "This is a test")
}

func TestMarshalUnmarshalLongString(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, "This is a longer string test that goes beyond 15 characters.")
}

func TestMarshalUnmarshalUtf8String(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, "Test ö覚𠜎")
}

func TestMarshalUnmarshalURI(t *testing.T) {
	testURL, err := url.Parse("http://example.com/path?query=something")
	if err != nil {
		t.Fatal(err)
	}
	assertMarshalUnmarshal(t, ContainerTypeNone, testURL)
}

func TestMarshalUnmarshalBytesSlice(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []byte{1, 2, 3, 4, 5})
}

func TestMarshalUnmarshalBytesArray(t *testing.T) {
	array := [...]byte{1, 2}
	assertMarshalUnmarshal(t, ContainerTypeNone, array)
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
	assertMarshalUnmarshal(t, ContainerTypeNone, array)
}

func TestMarshalUnmarshalLongerBytesArray(t *testing.T) {
	array := make([]byte, 60000)
	assertMarshalUnmarshal(t, ContainerTypeNone, array)
}

func TestMarshalUnmarshalAllBasicTypes(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []interface{}{
		1, 250, 1000, 100000, 10000000000,
		-1, -250, -1000, -100000000000,
		1.5, 1.9582384465,
		"string", []byte{10, 11, 12},
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 9999999, time.UTC),
	})
}

func TestMarshalUnmarshalList(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []interface{}{
		1, 2, "test", 4,
	})
}

func TestMarshalUnmarshalListNil(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []interface{}{
		1, nil, "test", 4,
	})
}

func TestMarshalUnmarshalInlineList(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeList, []interface{}{
		1, 2, "test", 4,
	})
}

func TestMarshalUnmarshalMap(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, map[interface{}]interface{}{
		"a": 1,
		2:   "b",
	})
}

func TestMarshalUnmarshalMapNil(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, map[interface{}]interface{}{
		"a": nil,
	})
}

func TestMarshalUnmarshalMapMap(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, map[interface{}]interface{}{
		1: 2,
		"deep-map": map[interface{}]interface{}{
			3: 1000,
		},
	})
}

func TestMarshalUnmarshalInlineMap(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeUnorderedMap, map[interface{}]interface{}{
		"a": 1,
		2:   "b",
	})
}

func TestMarshalUnmarshalListList(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []interface{}{
		1, []interface{}{
			2, "test", 4,
		},
	})
}

func TestMarshalUnmarshalListMap(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, []interface{}{
		1, map[interface{}]interface{}{
			2: 3, 4: "blah",
		},
	})
}

func TestMarshalUnmarshalMapList(t *testing.T) {
	assertMarshalUnmarshal(t, ContainerTypeNone, map[interface{}]interface{}{
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
	assertMarshalUnmarshalProduces(t, ContainerTypeNone, structValue, map[interface{}]interface{}{
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
	assertMarshalUnmarshalProduces(t, ContainerTypeNone, structValue, map[interface{}]interface{}{
		"IntValue":    structValue.IntValue,
		"FloatValue":  structValue.FloatValue,
		"StringValue": structValue.StringValue,
		"ByteValue":   structValue.ByteValue,
	})
}
