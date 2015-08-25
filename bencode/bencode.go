package bencode

import (
	"fmt"
	"sort"
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
