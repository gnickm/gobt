package bencode

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type BEncodable interface {
	BEString() string
}

// BEncodeInteger type ----------------------------------------------

type BEncodeInteger int

func (bei BEncodeInteger) BEString() string {
	return fmt.Sprintf("i%de", bei)
}

// BEncodeString type -----------------------------------------------

type BEncodeString string

func (bes BEncodeString) BEString() string {
	return fmt.Sprintf("%d:%s", len(bes), bes)
}

// BEncodeList type -------------------------------------------------

type BEncodeList []BEncodable

func (bel BEncodeList) BEString() string {
	bestr := "l"
	for _, bencodable := range bel {
		bestr += bencodable.BEString()
	}
	return bestr + "e"
}

// BEncodeDictionary type -------------------------------------------

type BEncodeDictionary map[string]BEncodable

func (bed BEncodeDictionary) BEString() string {
	bestr := "d"

	// BitTorrent spec says keys should be sorted
	keys := make([]string, 0, len(bed))
	for key := range bed {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, dictKey := range keys {
		bestr += BEncodeString(dictKey).BEString()
		bestr += bed[dictKey].BEString()
	}
	return bestr + "e"
}

// Decode function --------------------------------------------------

func Decode(beString string) (BEncodable, error) {
	be, _, err := doDecode(beString, 0)
	return be, err
}

func doDecode(beString string, startIndex int) (BEncodable, int, error) {
	if len(beString) == 0 {
		return nil, 0, errors.New("bencode: cannot decode empty string")
	}
	if startIndex > len(beString) {
		return nil, startIndex, errors.New("bencode: unexpected end of string to be decoded")
	}
	if beString[startIndex] == 'i' {
		// Handle integer
		endIndex, err := findChar(beString, startIndex, 'e')
		if err != nil {
			return nil, startIndex, err
		}
		i, err := strconv.Atoi(beString[startIndex+1 : endIndex])
		if err != nil {
			return nil, startIndex, err
		} else {
			return BEncodeInteger(i), endIndex + 1, nil
		}
	} else if beString[startIndex] == 'l' {
		// Handle list
		beList := BEncodeList{}
		startIndex++
		for beString[startIndex] != 'e' {
			var beListItem BEncodable
			var err error
			beListItem, startIndex, err = doDecode(beString, startIndex)
			if err != nil {
				return nil, startIndex, err
			} else {
				beList = append(beList, beListItem)
			}
			// Handle falling off the edge of the string here, otherwise
			// our for condition blows up
			if startIndex >= len(beString) {
				return nil, startIndex, errors.New("bencode: unexpected end of string while trying to find list closure")
			}
		}
		return beList, startIndex + 1, nil
	} else if beString[startIndex] == 'd' {
		// Handle dictionary
		beDict := BEncodeDictionary{}
		startIndex++
		for beString[startIndex] != 'e' {
			var beKey, beValue BEncodable
			beKey, startIndex, _ = doDecode(beString, startIndex)
			beValue, startIndex, _ = doDecode(beString, startIndex)
			beKeyStr := string(beKey.(BEncodeString))
			beDict[beKeyStr] = beValue
		}
		return beDict, startIndex + 1, nil
	} else {
		// Handle string
		endIndex, err := findChar(beString, startIndex, ':')
		if err != nil {
			return nil, startIndex, err
		}
		strLength, err := strconv.Atoi(beString[startIndex:endIndex])
		if err != nil {
			return nil, startIndex, err
		}
		startIndex = endIndex + 1
		if startIndex+strLength > len(beString) {
			return nil, startIndex, errors.New("bencode: string length was greater than encoded string size")
		} else {
			return BEncodeString(beString[startIndex : startIndex+strLength]), startIndex + strLength, nil
		}
	}
}

func findChar(beString string, startIndex int, char uint8) (int, error) {
	for i := startIndex; i < len(beString); i++ {
		if beString[i] == char {
			return i, nil
		}
	}
	return startIndex, errors.New("bencode: failed to find expected terminating character")
}
