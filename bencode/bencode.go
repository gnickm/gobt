package bencode

import (
	"fmt"
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
