package db_utils

import (
	"central-knot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DatabaseConn *gorm.DB

func OpenDatabaseConnection() error {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DatabaseConn = db

	db.AutoMigrate(&models.PeerTorrent{})
	db.AutoMigrate(&models.Torrent{})
	db.AutoMigrate(&models.Peer{})

	return nil
}
