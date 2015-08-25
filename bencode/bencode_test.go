package bencode

import (
	"testing"
)

func TestBEncodeIntegerBasics(t *testing.T) {
	beInt := BEncodeInteger(123)
	if beInt != 123 {
		t.Errorf("Expected Value of '123', was '%s'", beInt)
	}
	if beInt.BEString() != "i123e" {
		t.Errorf("Expected BEString() of 'i123e', was '%s'", beInt.BEString())
	}
}

func TestBEncodeStringBasics(t *testing.T) {
	beStr := BEncodeString("Hello World!")
	if beStr != "Hello World!" {
		t.Errorf("Expected Value of 'Hello World', was '%s'", beStr)
	}
	if beStr.BEString() != "12:Hello World!" {
		t.Errorf("Expected BEString() of '12:Hello World!', was '%s'", beStr.BEString())
	}
}

func TestBEncodeListBasics(t *testing.T) {
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

func TestBEncodeDictionaryBasics(t *testing.T) {
	beDict := BEncodeDictionary{
		"KeyA": BEncodeInteger(123),
		"KeyB": BEncodeString("Hello World!"),
	}
	if len(beDict) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beDict))
	}
	if beDict.BEString() != "d4:KeyAi123e4:KeyB12:Hello World!e" {
		t.Errorf("Expected BEString() of 'd4:KeyAi123e4:KeyB12:Hello World!e', was '%s'", beDict.BEString())
	}

	beDict["KeyC"] = BEncodeList{BEncodeInteger(456), BEncodeString("nested")}

	if len(beDict) != 3 {
		t.Errorf("Expected length of 3, was '%d'", len(beDict))
	}
	if beDict.BEString() != "d4:KeyAi123e4:KeyB12:Hello World!4:KeyCli456e6:nestedee" {
		t.Errorf("Expected BEString() of 'd4:KeyAi123e4:KeyB12:Hello World!4:KeyCli456e6:nestedee', was '%s'", beDict.BEString())
	}
}

func TestBEncodeDictionaryKeySorting(t *testing.T) {
	// BitTorrent spec says the keys should be sorted in string order...
	beDict := BEncodeDictionary{
		"ZZZ": BEncodeString("zzz"),
		"AAA": BEncodeString("aaa"),
		"MMM": BEncodeString("mmm"),
	}

	if beDict.BEString() != "d3:AAA3:aaa3:MMM3:mmm3:ZZZ3:zzze" {
		t.Errorf("Expected BEString() of 'd3:AAA3:aaa3:MMM3:mmm3:ZZZ3:zzze', was '%s'", beDict.BEString())
	}

}
