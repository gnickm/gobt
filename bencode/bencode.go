package BE

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type BEncodable interface {
	Encode() string
}

// BEInteger type ----------------------------------------------

type BEInteger int

func (bei BEInteger) Encode() string {
	return fmt.Sprintf("i%de", bei)
}

// BEString type -----------------------------------------------

type BEString string

func (bes BEString) Encode() string {
	return fmt.Sprintf("%d:%s", len(bes), bes)
}

// BEList type -------------------------------------------------

type BEList []BEncodable

func (bel BEList) Encode() string {
	bestr := "l"
	for _, bencodable := range bel {
		bestr += bencodable.Encode()
	}
	return bestr + "e"
}

// BEDictionary type -------------------------------------------

type BEDictionary map[string]BEncodable

func (bed BEDictionary) Encode() string {
	bestr := "d"

	// BitTorrent spec says keys should be sorted
	keys := make([]string, 0, len(bed))
	for key := range bed {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, dictKey := range keys {
		bestr += BEString(dictKey).Encode()
		bestr += bed[dictKey].Encode()
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
		return nil, 0, errors.New("BE: cannot decode empty string")
	}
	if startIndex > len(beString) {
		return nil, startIndex, errors.New("BE: unexpected end of string to be decoded")
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
			return BEInteger(i), endIndex + 1, nil
		}
	} else if beString[startIndex] == 'l' {
		// Handle list
		beList := BEList{}
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
				return nil, startIndex, errors.New("BE: unexpected end of string while trying to find list closure")
			}
		}
		return beList, startIndex + 1, nil
	} else if beString[startIndex] == 'd' {
		// Handle dictionary
		beDict := BEDictionary{}
		startIndex++
		for beString[startIndex] != 'e' {
			var beKey, beValue BEncodable
			beKey, startIndex, _ = doDecode(beString, startIndex)
			beValue, startIndex, _ = doDecode(beString, startIndex)
			beKeyStr := string(beKey.(BEString))
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
			return nil, startIndex, errors.New("BE: string length was greater than encoded string size")
		} else {
			return BEString(beString[startIndex : startIndex+strLength]), startIndex + strLength, nil
		}
	}
}

func findChar(beString string, startIndex int, char uint8) (int, error) {
	for i := startIndex; i < len(beString); i++ {
		if beString[i] == char {
			return i, nil
		}
	}
	return startIndex, errors.New("BE: failed to find expected terminating character")
}
