package db_utils

import (
	"central-knot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DatabaseConn is the opened connection to the database.
// Can be used to query the DB.
var DatabaseConn *gorm.DB

// OpenDatabaseConnection auto migrates all the models to the database
// and stores an opened connection in DatabaseConn.
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
