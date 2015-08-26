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

func TestDecodeHappyPath(t *testing.T) {
	var err error
	var be BEncodable

	be, err = Decode("i1234e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beInteger := be.(BEncodeInteger)
	if beInteger != 1234 {
		t.Errorf("Expected value of 1234, was '%d'", beInteger)
	}

	be, err = Decode("12:Hello World!")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beString := be.(BEncodeString)
	if beString != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beString)
	}

	be, err = Decode("li123e12:Hello World!e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList := be.(BEncodeList)
	if len(beList) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beList))
	}
	if beList[0].(BEncodeInteger) != 123 {
		t.Errorf("Expected value of 123, was '%d'", beList[0].(BEncodeInteger))
	}
	if beList[1].(BEncodeString) != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beList[1].(BEncodeString))
	}

	be, err = Decode("d4:KeyAi123e4:KeyB12:Hello World!e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beDict := be.(BEncodeDictionary)
	if len(beDict) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beDict))
	}
	if beDict["KeyA"].(BEncodeInteger) != 123 {
		t.Errorf("Expected value of 123, was '%d'", beDict["KeyA"].(BEncodeInteger))
	}
	if beDict["KeyB"].(BEncodeString) != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beDict["KeyB"].(BEncodeString))
	}
}

func TestDecodeDeepEmbedsAndWeirdStuff(t *testing.T) {
	var err error
	var be BEncodable

	be, err = Decode("lllll9:very deepeeeee")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList := be.(BEncodeList)
	if len(beList) != 1 {
		t.Errorf("Expected length of 1, was '%d'", len(beList))
	}
	if beList[0].(BEncodeList)[0].(BEncodeList)[0].(BEncodeList)[0].(BEncodeList)[0].(BEncodeString) != "very deep" {
		t.Error("Expected value of 'very deep'")
	}

	be, err = Decode("i1234eEverything after valid BEncode will be ignored with no error")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beInteger := be.(BEncodeInteger)
	if beInteger != 1234 {
		t.Errorf("Expected value of 1234, was '%d'", beInteger)
	}

	// Lists can be empty
	be, err = Decode("le")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList = be.(BEncodeList)
	if len(beList) != 0 {
		t.Errorf("Expected length of 0, was '%d'", len(beList))
	}

	// Dictionaries can be empty
	be, err = Decode("de")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beDict := be.(BEncodeDictionary)
	if len(beDict) != 0 {
		t.Errorf("Expected length of 0, was '%d'", len(beDict))
	}
}

func TestDecodeGeneralFailureModes(t *testing.T) {
	var err error

	_, err = Decode("")
	if err == nil {
		t.Error("Expected error when decoding empty string")
	}

	_, err = Decode("This is not the BEncode you are looking for")
	if err == nil {
		t.Error("Expected error when decoding bogus string")
	}
}

func TestDecodeIntegerFailureModes(t *testing.T) {
	var err error
	_, err = Decode("i123")
	if err == nil {
		t.Error("Expected error when decoding unfinished integer")
	}

	_, err = Decode("ie")
	if err == nil {
		t.Error("Expected error when decoding integer with no numeric")
	}

	_, err = Decode("i1a2e")
	if err == nil {
		t.Error("Expected error when decoding integer with non numeric")
	}
}

func TestDecodeStringFailureModes(t *testing.T) {
	var err error
	_, err = Decode("4fail")
	if err == nil {
		t.Error("Expected error when decoding string with no separator")
	}

	_, err = Decode(":fail")
	if err == nil {
		t.Error("Expected error when decoding string with no char count")
	}

	_, err = Decode("4:fai")
	if err == nil {
		t.Error("Expected error when decoding string that is too short")
	}
}

func TestDecodeListFailureModes(t *testing.T) {
	var err error
	_, err = Decode("li1ei2ei3e")
	if err == nil {
		t.Error("Expected error when decoding list with no terminator")
	}
}
