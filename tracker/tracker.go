package main

import (
	"errors"
	"github.com/gnickm/gobt/bencode"
	"github.com/gnickm/gobt/bt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type AnnounceRequest struct {
	peerId      bt.PeerId
	infoHash    bt.InfoHash
	ip          net.IP
	port        int
	downloaded  int
	left        int
	uploaded    int
	eventString string
	compactMode bool
	respChan    chan bencode.BEDictionary
}

func makeAnnounceHandler(reqChan chan AnnounceRequest) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var responseDict bencode.BEDictionary
		ar, err := parseAnnounceRequest(r)
		if err != nil {
			responseDict = bencode.BEDictionary{}
			responseDict["failure message"] = bencode.BEString(err.Error())
		} else {
			reqChan <- *ar
			responseDict = <-ar.respChan
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, responseDict.BEncode())
	}
}

func parseAnnounceRequest(r *http.Request) (*AnnounceRequest, error) {
	query := r.URL.Query()

	// --- Handle 'info_hash' parameter
	if query.Get("info_hash") == "" {
		return nil, errors.New("Missing required parameter 'info_hash'")
	}
	infoHash := bt.InfoHash(query.Get("info_hash"))
	if !infoHash.Validate() {
		return nil, errors.New("Invalid value in 'info_hash'")
	}

	// --- Handle 'peer_id' parameter
	if query.Get("peer_id") == "" {
		return nil, errors.New("Missing required parameter 'peer_id'")
	}
	peerId := bt.PeerId(query.Get("peer_id"))
	if !peerId.Validate() {
		return nil, errors.New("Invalid value in 'peer_id'")
	}

	// --- Handle 'port' parameter
	if query.Get("port") == "" {
		return nil, errors.New("Missing required parameter 'port'")
	}
	port, err := strconv.Atoi(query.Get("port"))
	if err != nil || port < 1 || port > 65535 {
		return nil, errors.New("Invalid value in 'port'")
	}

	// --- Handle optional 'ip' parameter and find peer's IP
	var ip net.IP
	if query.Get("ip") == "" {
		// Use the request remote address if not specified
		chunks := strings.Split(":", r.RemoteAddr)
		ip = net.ParseIP(chunks[0])
	} else {
		ip = net.ParseIP(query.Get("ip"))
	}
	if ip == nil {
		return nil, errors.New("Invalid peer IP address")
	}

	downloaded, _ := strconv.Atoi(query.Get("downloaded"))
	left, _ := strconv.Atoi(query.Get("left"))
	uploaded, _ := strconv.Atoi(query.Get("uploaded"))

	ar := AnnounceRequest{
		peerId,
		infoHash,
		ip,
		port,
		downloaded,
		left,
		uploaded,
		query.Get("eventString"),
		query.Get("compact") == "1",
		make(chan bencode.BEDictionary),
	}

	return &ar, nil
}

func main() {
	announceReqChan := make(chan AnnounceRequest)
	http.HandleFunc("/announce", makeAnnounceHandler(announceReqChan))
	http.ListenAndServe(":80", nil)
}
