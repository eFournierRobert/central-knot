package repo_peer

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"context"

	"gorm.io/gorm"
)

// AddPeer takes in the information for a new peer and
// adds it to the database.
func AddPeer(id, ip string, port int) error {
	p := models.Peer{
		ClientId: id,
		Ip:       ip,
		Port:     port,
	}

	ctx := context.Background()
	return gorm.G[models.Peer](db_utils.DatabaseConn).Create(ctx, &p)
}

// GetPeer takes in the client ID ("peer id" in BEP 3) and returns
// the corresponding peer.
func GetPeer(id string) (models.Peer, error) {
	ctx := context.Background()
	return gorm.G[models.Peer](db_utils.DatabaseConn).Where("client_id = ?", id).First(ctx)
}

// ChangeIp updates the peer IP with the given client ID ("peer id" in BEP 3)
// to the new given IP.
func ChangeIp(clientId, newIp string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Peer](db_utils.DatabaseConn).
		Where("client_id = ?", clientId).
		Update(ctx, "ip", newIp)
	return err
}
