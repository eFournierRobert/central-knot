package repo_peerTorrent

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"context"

	"gorm.io/gorm"
)

// Add creates a new PeerTorrent and adds it to the database.
func Add(peerId, torrentId uint, uploaded, downloaded, left int64, event models.Events) error {
	p := models.PeerTorrent{
		PeerId:     peerId,
		TorrentId:  torrentId,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		Event:      event,
	}

	ctx := context.Background()
	return gorm.G[models.PeerTorrent](db_utils.DatabaseConn).Create(ctx, &p)
}

// Update takes the updated PeerTorrent and saves that modification
// in the database.
func Update(p *models.PeerTorrent) {
	db_utils.DatabaseConn.Save(p)
}

// Get returns the PeerTorrent that has the given peer ID (as in primary key)
// and torrent ID.
func Get(peerId, torrentId uint) (models.PeerTorrent, error) {
	ctx := context.Background()
	return gorm.G[models.PeerTorrent](db_utils.DatabaseConn).
		Where("peer_id = ? AND torrent_id = ?", peerId, torrentId).
		First(ctx)
}
