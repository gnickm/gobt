package BE

import (
	"testing"
)

func TestBEIntegerBasics(t *testing.T) {
	beInt := BEInteger(123)
	if beInt != 123 {
		t.Errorf("Expected Value of '123', was '%s'", beInt)
	}
	if beInt.Encode() != "i123e" {
		t.Errorf("Expected Encode() of 'i123e', was '%s'", beInt.Encode())
	}
}

func TestBEStringBasics(t *testing.T) {
	beStr := BEString("Hello World!")
	if beStr != "Hello World!" {
		t.Errorf("Expected Value of 'Hello World', was '%s'", beStr)
	}
	if beStr.Encode() != "12:Hello World!" {
		t.Errorf("Expected Encode() of '12:Hello World!', was '%s'", beStr.Encode())
	}
}

func TestBEListBasics(t *testing.T) {
	beList := BEList{BEInteger(123), BEString("Hello World!")}
	if len(beList) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beList))
	}
	if beList.Encode() != "li123e12:Hello World!e" {
		t.Errorf("Expected Encode() of 'li123e12:Hello World!e', was '%s'", beList.Encode())
	}

	// Can append lists to lists
	beList = append(beList, BEList{BEInteger(456), BEString("nested")})

	if len(beList) != 3 {
		t.Errorf("Expected length of 3, was '%d'", len(beList))
	}
	if beList.Encode() != "li123e12:Hello World!li456e6:nestedee" {
		t.Errorf("Expected Encode() of 'li123e12:Hello World!li456e6:nestede', was '%s'", beList.Encode())
	}
}

func TestBEDictionaryBasics(t *testing.T) {
	beDict := BEDictionary{
		"KeyA": BEInteger(123),
		"KeyB": BEString("Hello World!"),
	}
	if len(beDict) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beDict))
	}
	if beDict.Encode() != "d4:KeyAi123e4:KeyB12:Hello World!e" {
		t.Errorf("Expected Encode() of 'd4:KeyAi123e4:KeyB12:Hello World!e', was '%s'", beDict.Encode())
	}

	beDict["KeyC"] = BEList{BEInteger(456), BEString("nested")}

	if len(beDict) != 3 {
		t.Errorf("Expected length of 3, was '%d'", len(beDict))
	}
	if beDict.Encode() != "d4:KeyAi123e4:KeyB12:Hello World!4:KeyCli456e6:nestedee" {
		t.Errorf("Expected Encode() of 'd4:KeyAi123e4:KeyB12:Hello World!4:KeyCli456e6:nestedee', was '%s'", beDict.Encode())
	}
}

func TestBEDictionaryKeySorting(t *testing.T) {
	// BitTorrent spec says the keys should be sorted in string order...
	beDict := BEDictionary{
		"ZZZ": BEString("zzz"),
		"AAA": BEString("aaa"),
		"MMM": BEString("mmm"),
	}

	if beDict.Encode() != "d3:AAA3:aaa3:MMM3:mmm3:ZZZ3:zzze" {
		t.Errorf("Expected Encode() of 'd3:AAA3:aaa3:MMM3:mmm3:ZZZ3:zzze', was '%s'", beDict.Encode())
	}

}

func TestDecodeHappyPath(t *testing.T) {
	var err error
	var be BEncodable

	be, err = Decode("i1234e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beInteger := be.(BEInteger)
	if beInteger != 1234 {
		t.Errorf("Expected value of 1234, was '%d'", beInteger)
	}

	be, err = Decode("12:Hello World!")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beString := be.(BEString)
	if beString != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beString)
	}

	be, err = Decode("li123e12:Hello World!e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList := be.(BEList)
	if len(beList) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beList))
	}
	if beList[0].(BEInteger) != 123 {
		t.Errorf("Expected value of 123, was '%d'", beList[0].(BEInteger))
	}
	if beList[1].(BEString) != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beList[1].(BEString))
	}

	be, err = Decode("d4:KeyAi123e4:KeyB12:Hello World!e")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beDict := be.(BEDictionary)
	if len(beDict) != 2 {
		t.Errorf("Expected length of 2, was '%d'", len(beDict))
	}
	if beDict["KeyA"].(BEInteger) != 123 {
		t.Errorf("Expected value of 123, was '%d'", beDict["KeyA"].(BEInteger))
	}
	if beDict["KeyB"].(BEString) != "Hello World!" {
		t.Errorf("Expected value of 'Hello World!', was '%s'", beDict["KeyB"].(BEString))
	}
}

func TestDecodeDeepEmbedsAndWeirdStuff(t *testing.T) {
	var err error
	var be BEncodable

	be, err = Decode("lllll9:very deepeeeee")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList := be.(BEList)
	if len(beList) != 1 {
		t.Errorf("Expected length of 1, was '%d'", len(beList))
	}
	if beList[0].(BEList)[0].(BEList)[0].(BEList)[0].(BEList)[0].(BEString) != "very deep" {
		t.Error("Expected value of 'very deep'")
	}

	be, err = Decode("i1234eEverything after valid BE will be ignored with no error")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beInteger := be.(BEInteger)
	if beInteger != 1234 {
		t.Errorf("Expected value of 1234, was '%d'", beInteger)
	}

	// Lists can be empty
	be, err = Decode("le")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beList = be.(BEList)
	if len(beList) != 0 {
		t.Errorf("Expected length of 0, was '%d'", len(beList))
	}

	// Dictionaries can be empty
	be, err = Decode("de")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	beDict := be.(BEDictionary)
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

	_, err = Decode("This is not the BE you are looking for")
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
