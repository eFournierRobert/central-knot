package service

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePeerTorrentStatusValues(t *testing.T) {
	uploadStr := "10"
	downloadedStr := "13"
	leftStr := "14"

	upload, down, left, err := parsePeerTorrentStatusValues(uploadStr, downloadedStr, leftStr)
	assert.Nil(t, err)
	assert.Equal(t, 10, upload)
	assert.Equal(t, 13, down)
	assert.Equal(t, 14, left)
}

func TestParsePeerTorrentInvalidStatusValues(t *testing.T) {
	uploadStr := "10"
	downloadedStr := "1a"
	leftStr := "14"

	_, _, _, err := parsePeerTorrentStatusValues(uploadStr, downloadedStr, leftStr)
	assert.NotNil(t, err)
}

func TestEnsurePeerNotExistingAndExisting(t *testing.T) {
	db_utils.SetupTest(t)

	peerId := "id"
	ip := "10.01.01.12"
	port := "2040"

	peer, err := ensurePeer(peerId, ip, port)
	assert.Nil(t, err)
	assert.Equal(t, peerId, peer.ClientId)
	assert.Equal(t, ip, peer.Ip)
	assert.Equal(t, 2040, peer.Port)

	peer2, err := ensurePeer(peerId, ip, port)
	assert.Nil(t, err)
	assert.Equal(t, peer, peer2)
}

func TestEnsurePeerWithInvalidPort(t *testing.T) {
	db_utils.SetupTest(t)

	peerId := "id"
	ip := "10.01.01.12"
	port := "port"

	peer, err := ensurePeer(peerId, ip, port)
	assert.NotNil(t, err)
	assert.Empty(t, peer)
}

func TestEnsureTorrentNotExistingAndExisting(t *testing.T) {
	db_utils.SetupTest(t)

	infoHash := []byte("hash39021-93021")

	torrent, err := ensureTorrent(infoHash)
	assert.Nil(t, err)
	assert.Equal(t, infoHash, torrent.InfoHash)

	torrent2, err := ensureTorrent(infoHash)
	assert.Nil(t, err)
	assert.Equal(t, torrent, torrent2)
}

func TestUpsertPeerTorrentAndGetPeerList(t *testing.T) {
	db_utils.SetupTest(t)

	peerId := "id"
	ip := "10.01.01.12"
	port := "2040"

	peer, err := ensurePeer(peerId, ip, port)
	assert.Nil(t, err)
	assert.NotEmpty(t, peer)

	peerId2 := "i2"

	peer2, err := ensurePeer(peerId2, ip, port)
	assert.Nil(t, err)
	assert.NotEmpty(t, peer)

	infoHash := []byte("hash39021-93021")

	torrent, err := ensureTorrent(infoHash)
	assert.Nil(t, err)
	assert.NotEmpty(t, torrent)

	dto := models.FullPeerDto{
		InfoHash:   string(infoHash),
		PeerId:     peerId,
		Port:       port,
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	err = upsertPeerTorrent(peer.ID, torrent.ID, &dto, models.Started)
	assert.Nil(t, err)
	err = upsertPeerTorrent(peer2.ID, torrent.ID, &dto, models.Started)
	assert.Nil(t, err)

	peerListDto, err := buildPeerListDto(infoHash, peerId)
	assert.Nil(t, err)
	assert.Equal(t, requestInterval, peerListDto.Interval)
	assert.Len(t, peerListDto.Peers, 1)
	assert.Equal(t, peer2.ToDto(), peerListDto.Peers[0])
}

func TestUpsertPeerTorrentWithNonExistingPeer(t *testing.T) {
	db_utils.SetupTest(t)

	infoHash := []byte("hash39021-93021")

	torrent, err := ensureTorrent(infoHash)
	assert.Nil(t, err)
	assert.NotEmpty(t, torrent)

	dto := models.FullPeerDto{
		InfoHash:   string(infoHash),
		PeerId:     "peerId",
		Port:       "22",
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	err = upsertPeerTorrent(1, torrent.ID, &dto, models.Started)
	assert.NotNil(t, err)
}

func TestUpsertPeerTorrentWithNonExistingTorrent(t *testing.T) {
	db_utils.SetupTest(t)

	peerId := "id"
	ip := "10.01.01.12"
	port := "2040"

	peer, err := ensurePeer(peerId, ip, port)
	assert.Nil(t, err)
	assert.NotEmpty(t, peer)

	dto := models.FullPeerDto{
		InfoHash:   "infohash",
		PeerId:     peer.ClientId,
		Port:       port,
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	err = upsertPeerTorrent(peer.ID, 1, &dto, models.Started)
	assert.NotNil(t, err)
}

func TestBuildEmptyPeerListDto(t *testing.T) {
	db_utils.SetupTest(t)

	peerId := "id"
	ip := "10.01.01.12"
	port := "2040"

	peer, err := ensurePeer(peerId, ip, port)
	assert.Nil(t, err)
	assert.NotEmpty(t, peer)

	infoHash := []byte("hash39021-93021")

	torrent, err := ensureTorrent(infoHash)
	assert.Nil(t, err)
	assert.NotEmpty(t, torrent)

	dto := models.FullPeerDto{
		InfoHash:   string(infoHash),
		PeerId:     peerId,
		Port:       port,
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	err = upsertPeerTorrent(peer.ID, torrent.ID, &dto, models.Started)
	assert.Nil(t, err)

	peerListDto, err := buildPeerListDto(infoHash, peerId)
	assert.Nil(t, err)
	assert.Equal(t, requestInterval, peerListDto.Interval)
	assert.Len(t, peerListDto.Peers, 0)
}

func TestAnnounce(t *testing.T) {
	db_utils.SetupTest(t)

	ip := "10.0.0.19"
	dto := models.FullPeerDto{
		InfoHash:   "infohash",
		PeerId:     "peerid",
		Port:       "1020",
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	peerList, err := Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 0)
	assert.Equal(t, requestInterval, peerList.Interval)

	dto.PeerId = "newid"
	peerList, err = Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 1)
	assert.Equal(t, requestInterval, peerList.Interval)
}

func TestAnnounceInvalidEvent(t *testing.T) {
	db_utils.SetupTest(t)

	ip := "10.0.0.19"
	dto := models.FullPeerDto{
		InfoHash:   "infohash",
		PeerId:     "peerid",
		Port:       "1020",
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "event",
		Compact:    false,
	}

	peerList, err := Announce(&dto, ip)
	assert.NotNil(t, err)
	assert.Nil(t, peerList)
}

func TestAnnounceUpdate(t *testing.T) {
	db_utils.SetupTest(t)

	ip := "10.0.0.19"
	dto := models.FullPeerDto{
		InfoHash:   "infohash",
		PeerId:     "peerid",
		Port:       "1020",
		Uploaded:   "10",
		Downloaded: "12",
		Left:       "14",
		Event:      "started",
		Compact:    false,
	}

	peerList, err := Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 0)
	assert.Equal(t, requestInterval, peerList.Interval)

	dto.PeerId = "newid"
	peerList, err = Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 1)
	assert.Equal(t, requestInterval, peerList.Interval)

	dto.Event = "stopped"
	peerList, err = Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)

	dto.PeerId = "peerid"
	peerList, err = Announce(&dto, ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, peerList)
	assert.Len(t, peerList.Peers, 0)
}
