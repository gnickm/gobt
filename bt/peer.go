package bt

import (
	"errors"
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

func (peer *Peer) AddInfoHash(hash InfoHash) {
	for _, existingHash := range peer.infoHashList {
		if existingHash == hash {
			// Already exists
			return
		}
	}
	peer.infoHashList = append(peer.infoHashList, hash)
}

func NewPeer(inPeerId *PeerId, ipString string, port int) (*Peer, error) {
	var peerId PeerId
	if inPeerId == nil {
		peerId = RandomPeerId()
	} else {
		peerId = *inPeerId
		if !peerId.Validate() {
			return nil, errors.New(fmt.Sprintf("NewPeer: Invalid Peer ID: %v", peerId))
		}
	}

	if port < 1 || port > 65535 {
		return nil, errors.New(fmt.Sprintf("NewPeer: Invalid port: %d", port))
	}

	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.New(fmt.Sprintf("NewPeer: Invalid IP address: %s", ipString))
	}

	peer := Peer{peerId, ip, port, []InfoHash{}}

	return &peer, nil
}
