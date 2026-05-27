package repo_peerTorrent

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	repo_peer "central-knot/internal/repository/peer"
	repo_torrent "central-knot/internal/repository/torrent"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGetPeerTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425
	infoHash := []byte("infohash")

	err := repo_torrent.Add(infoHash)
	assert.Nil(t, err)

	err = repo_peer.Add(id, ip, port)
	assert.Nil(t, err)

	err = Add(1, 1, 0, 0, 0, models.Started)
	assert.Nil(t, err)

	pt, err := Get(1, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, pt)
}

func TestGetNonExistingPeerTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	pt, err := Get(0, 0)
	assert.NotNil(t, err)
	assert.Empty(t, pt)
}

func TestAddAndUpdatePeerTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425
	infoHash := []byte("infohash")

	err := repo_torrent.Add(infoHash)
	assert.Nil(t, err)

	err = repo_peer.Add(id, ip, port)
	assert.Nil(t, err)

	err = Add(1, 1, 0, 0, 0, models.Started)
	assert.Nil(t, err)

	pt, err := Get(1, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, pt)

	pt.Event = models.Stopped
	Update(&pt)

	pt, err = Get(1, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, pt)
	assert.Equal(t, models.Stopped, pt.Event)
}

func TestAddPeerTorrentWithInvalidPeerAndTorrentId(t *testing.T) {
	db_utils.SetupTest(t)

	err := Add(1, 1, 0, 0, 0, models.Started)
	assert.NotNil(t, err)
}
