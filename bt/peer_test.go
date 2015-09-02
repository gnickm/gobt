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

func TestNewPeerErrorConditions(t *testing.T) {
	var err error
	badPeerId := PeerId("bad peerId")
	_, err = NewPeer(&badPeerId, "1.2.3.4", 1234)
	if err == nil {
		t.Error("Expected error with bad peerId")
	}

	_, err = NewPeer(nil, "bad IP", 1234)
	if err == nil {
		t.Error("Expected error with bad IP")
	}

	_, err = NewPeer(nil, "1.2.3.4", -123)
	if err == nil {
		t.Error("Expected error with bad port")
	}

	_, err = NewPeer(nil, "1.2.3.4", 99999)
	if err == nil {
		t.Error("Expected error with bad port")
	}
}

func TestPeerAddInfoHash(t *testing.T) {
	hash1 := InfoHash("aaaaaaaaaaaaaaaaaaaa")
	hash2 := InfoHash("bbbbbbbbbbbbbbbbbbbb")
	peer, _ := NewPeer(nil, "1.2.3.4", 1234)

	if len(peer.infoHashList) != 0 {
		t.Errorf("Initial infoHashList length should be 0")
	}

	peer.AddInfoHash(hash1)
	if len(peer.infoHashList) != 1 {
		t.Errorf("infoHashList length should be 1 after adding hash")
	}

	peer.AddInfoHash(hash1)
	if len(peer.infoHashList) != 1 {
		t.Errorf("infoHashList length should still be 1 after adding same hash")
	}

	peer.AddInfoHash(hash2)
	if len(peer.infoHashList) != 2 {
		t.Errorf("infoHashList length should be 2 after adding new hash")
	}
}
