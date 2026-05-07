package repo_peer

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"context"

	"gorm.io/gorm"
)

func AddPeer(id, ip string, port int) error {
	p := models.Peer{
		ClientId: id,
		Ip:       ip,
		Port:     port,
	}

	ctx := context.Background()
	return gorm.G[models.Peer](db_utils.DatabaseConn).Create(ctx, &p)
}

func GetPeer(id string) (models.Peer, error) {
	ctx := context.Background()
	return gorm.G[models.Peer](db_utils.DatabaseConn).Where("client_id = ?", id).First(ctx)
}

func ChangeIp(clientId, newIp string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Peer](db_utils.DatabaseConn).
		Where("client_id = ?", clientId).
		Update(ctx, "ip", newIp)
	return err
}
