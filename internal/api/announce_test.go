package api

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
)

type TestCompactPeerListDto struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func TestAnnounce(t *testing.T) {
	db_utils.SetupTest(t)
	req := httptest.NewRequest(http.MethodGet, "/announce?info_hash=B%a1%3f%0b%db%ce%5b%2bO7%e4%c8%b6%90%90%25%cf%9b%1d%93&peer_id=-qB5210-!mp3edyJ(tiE&port=50789&uploaded=0&downloaded=0&left=407714073&corrupt=0&key=999739FD&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0", nil)
	w := httptest.NewRecorder()

	AnnounceHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	var peerList models.CompactPeerListDto
	bencode.Unmarshal(w.Body, &peerList)
	assert.NotEmpty(t, peerList)
}

func TestAnnounceNoParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/announce", nil)
	w := httptest.NewRecorder()

	AnnounceHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
}

func TestAnnounceWithTwoPeers(t *testing.T) {
	db_utils.SetupTest(t)
	req := httptest.NewRequest(http.MethodGet, "/announce?info_hash=B%a1%3f%0b%db%ce%5b%2bO7%e4%c8%b6%90%90%25%cf%9b%1d%93&peer_id=-qB5210-!mp3edyJ(tiE&port=50789&uploaded=0&downloaded=0&left=407714073&corrupt=0&key=999739FD&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0", nil)
	w := httptest.NewRecorder()

	AnnounceHandler(w, req)

	res := w.Result()

	assert.Equal(t, 200, res.StatusCode)
	res.Body.Close()

	req = httptest.NewRequest(http.MethodGet, "/announce?info_hash=B%a1%3f%0b%db%ce%5b%2bO7%e4%c8%b6%90%90%25%cf%9b%1d%93&peer_id=peer2&port=50789&uploaded=0&downloaded=0&left=407714073&corrupt=0&key=999739FD&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0", nil)
	w = httptest.NewRecorder()

	AnnounceHandler(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	var peerList TestCompactPeerListDto
	bencode.Unmarshal(res.Body, &peerList)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 6)
}
