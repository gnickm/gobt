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
