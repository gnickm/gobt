package main

import (
	"io"
	"net/http"
	"github.com/gnickm/gobt/bencode"
)

func handleAnnounce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	responseDict := bencode.BEDictionary{}

	query := r.URL.Query()
	if query.Get("info_hash") == "" {
		responseDict["failure reason"] = bencode.BEString("Missing required parameter 'info_hash'")
	} else if query.Get("peer_id") == "" {
		responseDict["failure reason"] = bencode.BEString("Missing required parameter 'peer_id'")
	}

	io.WriteString(w, responseDict.BEncode())
}

func main() {
	http.HandleFunc("/announce", handleAnnounce)
	http.ListenAndServe(":80", nil)
}