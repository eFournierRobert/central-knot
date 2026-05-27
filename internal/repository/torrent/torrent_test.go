package repo_torrent

import (
	"central-knot/internal/repository/db_utils"
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
