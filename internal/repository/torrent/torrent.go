package repo_torrent

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"context"

	"gorm.io/gorm"
)

// Add creates a Torrent with the given info hash and
// adds it to the database.
func Add(infoHash []byte) error {
	t := models.Torrent{InfoHash: infoHash}

	ctx := context.Background()

	return gorm.G[models.Torrent](db_utils.DatabaseConn).Create(ctx, &t)
}

// Get returns the Torrent with the given info hash.
func Get(infoHash []byte) (models.Torrent, error) {
	ctx := context.Background()
	return gorm.G[models.Torrent](db_utils.DatabaseConn).Where("info_hash = ?", infoHash).First(ctx)
}

// GetPeersFromHash returns a list of peer that are in the swarm of the
// Torrent with the given info hash.
//
// askingClientId is the client ID ("peer id" in BEP 3) of the client making
// the request. It will be excluded from the returned list.
func GetPeersFromHash(infoHash []byte, askingClientId string) []models.Peer {
	var peers []models.Peer
	db_utils.DatabaseConn.Model(&models.Peer{}).
		Joins("left join peer_torrents ON peers.id = peer_torrents.peer_id").
		Joins("left join torrents ON torrents.id = peer_torrents.torrent_id").
		Where(
			"torrents.info_hash = ? AND peer_torrents.event != ? AND peers.client_id != ?",
			infoHash,
			models.Stopped,
			askingClientId,
		).Scan(&peers)

	return peers
}
