package concise_encoding

import (
	"testing"
)

func assertCTEDecodeEncode(t *testing.T, expected string) {
	events, err := cteDecode([]byte(expected))
	if err != nil {
		t.Error(err)
		return
	}
	result := string(cteEncode(events...))
	if result != expected {
		t.Errorf("Expected [%v] but got [%v]", expected, result)
	}
}

func TestMapFloatKey(t *testing.T) {
	assertCTEDecodeEncode(t, "c1 {nil=@nil 1.5=1000}")
}
