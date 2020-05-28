package concise_encoding

import (
	"testing"
)

func TestNullEventReceiver(t *testing.T) {
	receiver := new(NullEventReceiver)
	invokeEvents(receiver, v(1), l(), pi(1), s("blah"), e())
}
