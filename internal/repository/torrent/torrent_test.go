package repo_torrent

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	repo_peer "central-knot/internal/repository/peer"
	repo_peerTorrent "central-knot/internal/repository/peerTorrent"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGetTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	infoHash := []byte("infohash")
	err := Add(infoHash)
	assert.Nil(t, err)

	torrent, err := Get(infoHash)
	assert.Nil(t, err)
	assert.Equal(t, infoHash, torrent.InfoHash)
}

func TestAddExistingTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	infoHash := []byte("infohash")
	err := Add(infoHash)
	assert.Nil(t, err)

	err = Add(infoHash)
	assert.NotNil(t, err)
}

func TestGetNonExistingTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	infoHash := []byte("infohash")

	torrent, err := Get(infoHash)
	assert.NotNil(t, err)
	assert.Empty(t, torrent)
}

func TestGetPeersFromHash(t *testing.T) {
	db_utils.SetupTest(t)

	db_utils.SetupTest(t)

	id1 := "newid1"
	id2 := "newid2"
	ip := "10.0.0.10"
	port := 2425
	infoHash := []byte("infohash")

	err := Add(infoHash)
	assert.Nil(t, err)

	err = repo_peer.Add(id1, ip, port)
	assert.Nil(t, err)
	err = repo_peer.Add(id2, ip, port)
	assert.Nil(t, err)

	err = repo_peerTorrent.Add(1, 1, 0, 0, 0, models.Started)
	assert.Nil(t, err)
	err = repo_peerTorrent.Add(2, 1, 0, 0, 0, models.Started)
	assert.Nil(t, err)

	peers := GetPeersFromHash(infoHash, id1)
	assert.Len(t, peers, 1)
	assert.Equal(t, id2, peers[0].ClientId)
}
