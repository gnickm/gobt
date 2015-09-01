package bt

import (
	"testing"
)

func TestRandomPeerId(t *testing.T) {
	peerIdMap := map[PeerId]bool{}
	for i := 0; i < 10000; i++ {
		peerId := RandomPeerId()
		if !peerId.Validate() {
			t.Errorf("Unexpected invalid random PeerId: %v", peerId)
		}
		_, exists := peerIdMap[peerId]
		if exists {
			t.Errorf("Unexpected duplicate PeerId: %v", peerId)
		} else {
			peerIdMap[peerId] = true
		}
	}
}

func TestPeerIdValidate(t *testing.T) {
	var peerId PeerId
	peerId = RandomPeerId()
	if !peerId.Validate() {
		t.Errorf("Unexpected invalid random PeerId: %v", peerId)
	}

	peerId = PeerId("This is valid......!")
	if !peerId.Validate() {
		t.Errorf("Unexpected invalid random PeerId: %v", peerId)
	}

	peerId = PeerId("bad short ID")
	if peerId.Validate() {
		t.Errorf("Unexpected valid PeerId: %v", peerId)
	}

	peerId = PeerId("this is a bad ID since it is waaaaay too long")
	if peerId.Validate() {
		t.Errorf("Unexpected valid PeerId: %v", peerId)
	}
}
