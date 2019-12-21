package cbe

import (
	"net/url"
	"testing"
	"time"
)

func newURI(uri string) *url.URL {
	var result *url.URL
	var err error
	if result, err = url.Parse(uri); err != nil {
		panic(err)
	}
	return result
}

func TestMarshalUnmarshalBool(t *testing.T) {
	var v bool
	assertMarshalUnmarshal(t, false, &v)
	assertMarshalUnmarshal(t, true, &v)
}

func TestMarshalUnmarshal0(t *testing.T) {
	var v int = 1
	assertMarshalUnmarshal(t, 0, &v)
}

func TestMarshalUnmarshal1(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 1, &v)
	assertMarshalUnmarshal(t, -1, &v)
}

func TestMarshalUnmarshal200(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 200, &v)
	assertMarshalUnmarshal(t, -200, &v)
}

func TestMarshalUnmarshal2000(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 2000, &v)
	assertMarshalUnmarshal(t, -2000, &v)
}

func TestMarshalUnmarshal10000(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 10000, &v)
	assertMarshalUnmarshal(t, -10000, &v)
}

func TestMarshalUnmarshal100000(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 100000, &v)
	assertMarshalUnmarshal(t, -100000, &v)
}

func TestMarshalUnmarshal10000000000(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 10000000000, &v)
	assertMarshalUnmarshal(t, -10000000000, &v)
}

func TestMarshalUnmarshal100000000000(t *testing.T) {
	var v int
	assertMarshalUnmarshal(t, 100000000000, &v)
	assertMarshalUnmarshal(t, -100000000000, &v)
}

func TestMarshalUnmarshal1_5(t *testing.T) {
	var v float64
	assertMarshalUnmarshal(t, 1.5, &v)
}

func TestMarshalUnmarshal1_0123(t *testing.T) {
	var v float64
	assertMarshalUnmarshal(t, 1.0123, &v)
}

func TestMarshalUnmarshalTime(t *testing.T) {
	src := time.Date(2019, time.Month(11), 4, 18, 22, 45, 42300, time.UTC)
	dst := time.Date(2000, time.Month(1), 1, 1, 1, 1, 1, time.UTC)
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalEmptyString(t *testing.T) {
	src := ""
	dst := ""
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalShortString(t *testing.T) {
	src := "This is a test"
	dst := ""
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalLongString(t *testing.T) {
	src := "This is a longer string test that goes beyond 15 characters."
	dst := ""
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalUtf8String(t *testing.T) {
	src := "Test ö覚𠜎"
	dst := ""
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalURI(t *testing.T) {
	src := newURI("http://example.com/path?query=something")
	dst := newURI("http://example.com")
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalBytesSlice(t *testing.T) {
	src := []byte{1, 2, 3, 4, 5}
	var dst []byte
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalBytesArray(t *testing.T) {
	array := [...]byte{1, 2}
	dst := [...]byte{0, 0}
	assertMarshalUnmarshal(t, array, &dst)
}

func TestMarshalUnmarshalLongerBytesArray(t *testing.T) {
	array := make([]byte, 60000)
	array[103] = 90
	array[59323] = 100
	dst := make([]byte, 60000)
	assertMarshalUnmarshal(t, array, &dst)
}

func TestMarshalUnmarshalAllBasicTypes(t *testing.T) {
	src := []interface{}{
		1, 250, 1000, 100000, 10000000000,
		-1, -250, -1000, -100000000000,
		1.5, 1.9582384465,
		"string", []byte{10, 11, 12},
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 1, 1, 0, 0, 0, 9999999, time.UTC),
	}
	dst := []interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalList(t *testing.T) {
	src := []interface{}{1, 2, "test", 4}
	dst := []interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalListNil(t *testing.T) {
	src := []interface{}{
		1, nil, "test", 4,
	}
	dst := []interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalMap(t *testing.T) {
	src := map[interface{}]interface{}{
		"a": 1,
		2:   "b",
	}
	dst := map[interface{}]interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalMapNil(t *testing.T) {
	src := map[interface{}]interface{}{
		"a": nil,
	}
	dst := map[interface{}]interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalMapMap(t *testing.T) {
	src := map[interface{}]interface{}{
		1: 2,
		"deep-map": map[interface{}]interface{}{
			3: 1000,
		},
	}
	dst := map[interface{}]interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalListList(t *testing.T) {
	src := []interface{}{
		1, []interface{}{
			2, "test", 4,
		},
	}
	dst := []interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalListMap(t *testing.T) {
	src := []interface{}{
		1, map[interface{}]interface{}{
			2: 3, 4: "blah",
		},
	}
	dst := []interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

func TestMarshalUnmarshalMapList(t *testing.T) {
	src := map[interface{}]interface{}{
		1: 2,
		"deep-list": []interface{}{
			2, "some list entry",
		},
	}
	dst := map[interface{}]interface{}{}
	assertMarshalUnmarshal(t, src, &dst)
}

type MyTestStruct struct {
	IntValue    int
	FloatValue  float32
	StringValue string
	ByteValue   []byte
}

func TestMarshalUnmarshalStruct(t *testing.T) {
	src := new(MyTestStruct)
	src.IntValue = 4000
	src.FloatValue = 2.5
	src.StringValue = "test"
	src.ByteValue = []byte{0x00, 0x01, 0x02}
	dst := new(MyTestStruct)
	assertMarshalUnmarshal(t, src, dst)
}
