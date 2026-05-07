package repo_peerTorrent

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"context"

	"gorm.io/gorm"
)

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

func Update(p *models.PeerTorrent) {
	db_utils.DatabaseConn.Save(p)
}

func Get(peerdId, torrentId uint) (models.PeerTorrent, error) {
	ctx := context.Background()
	return gorm.G[models.PeerTorrent](db_utils.DatabaseConn).
		Where("peer_id = ? AND torrent_id = ?", peerdId, torrentId).
		First(ctx)
}
