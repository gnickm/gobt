package bencode

import (
	"testing"
)

func TestBEncodeInteger(t *testing.T) {
	beInt := BEncodeInteger(123)
	if beInt != 123 {
		t.Errorf("Expected Value of '123', was '%s'", beInt)
	}
	if beInt.BEString() != "i123e" {
		t.Errorf("Expected BEString() of 'i123e', was '%s'", beInt.BEString())
	}
}

func TestBEncodeString(t *testing.T) {
	beStr := BEncodeString("Hello World!")
	if beStr != "Hello World!" {
		t.Errorf("Expected Value of 'Hello World', was '%s'", beStr)
	}
	if beStr.BEString() != "12:Hello World!" {
		t.Errorf("Expected BEString() of '12:Hello World!', was '%s'", beStr.BEString())
	}
}
