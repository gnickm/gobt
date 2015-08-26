package bencode

import (
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

func Decode(beString string) BEncodable {
	be, _ := doDecode(beString, 0)
	return be
}

func doDecode(beString string, startIndex int) (BEncodable, int) {
	if len(beString) == 0 {
		return nil, 0
	}
	if beString[startIndex] == 'i' {
		// Handle integer
		endIndex := findChar(beString, startIndex, 'e')
		i, _ := strconv.Atoi(beString[startIndex+1 : endIndex])
		return BEncodeInteger(i), endIndex + 1
	} else if beString[startIndex] == 'l' {
		// Handle list
		beList := BEncodeList{}
		startIndex++
		for beString[startIndex] != 'e' {
			var beListItem BEncodable
			beListItem, startIndex = doDecode(beString, startIndex)
			beList = append(beList, beListItem)
		}
		return beList, startIndex + 1
	} else if beString[startIndex] == 'd' {
		// Handle dictionary
		beDict := BEncodeDictionary{}
		startIndex++
		for beString[startIndex] != 'e' {
			var beKey, beValue BEncodable
			beKey, startIndex = doDecode(beString, startIndex)
			beValue, startIndex = doDecode(beString, startIndex)
			beKeyStr := string(beKey.(BEncodeString))
			beDict[beKeyStr] = beValue
		}
		return beDict, startIndex + 1
	} else {
		// Handle string
		endIndex := findChar(beString, startIndex, ':')
		strLength, _ := strconv.Atoi(beString[startIndex:endIndex])
		startIndex = endIndex + 1
		return BEncodeString(beString[startIndex : startIndex+strLength]), startIndex + strLength
	}
}

func findChar(beString string, startIndex int, char uint8) int {
	for i := startIndex; i < len(beString); i++ {
		if beString[i] == char {
			return i
		}
	}
	return startIndex
}
