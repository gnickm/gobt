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

func TestBEncodeList(t *testing.T) {
	beList := BEncodeList{BEncodeInteger(123), BEncodeString("Hello World!")}
	if len(beList) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beList))
	}
	if beList.BEString() != "li123e12:Hello World!e" {
		t.Errorf("Expected BEString() of 'li123e12:Hello World!e', was '%s'", beList.BEString())
	}

	// Can append lists to lists
	beList = append(beList, BEncodeList{BEncodeInteger(456), BEncodeString("nested")})

	if len(beList) != 3 {
		t.Errorf("Expected length of 3, was '%d'", len(beList))
	}
	if beList.BEString() != "li123e12:Hello World!li456e6:nestedee" {
		t.Errorf("Expected BEString() of 'li123e12:Hello World!li456e6:nestede', was '%s'", beList.BEString())
	}
}
