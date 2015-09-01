package bt

import (
	"fmt"
	"math"
	"math/rand"
	"net"
)

// PeerId type ------------------------------------------------------

type PeerId string

func (peerId PeerId) Validate() bool {
	return len(peerId) == 20
}

func RandomPeerId() PeerId {
	return PeerId(fmt.Sprintf("-GO0001-%012d", rand.Int63n(int64(math.Pow10(12)-1))))
}

// Peer type --------------------------------------------------------

type Peer struct {
	id           PeerId
	ip           net.IP
	port         int
	infoHashList []InfoHash
}

func (peer Peer) AddInfoHash(hash InfoHash) {
	for _, existingHash := range peer.infoHashList {
		if existingHash == hash {
			// Already exists
			return
		}
	}
	peer.infoHashList = append(peer.infoHashList, hash)
}
