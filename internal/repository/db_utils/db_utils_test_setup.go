package db_utils

import (
	"central-knot/internal/models"
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTest(tb testing.TB) {
	dbFile, err := os.CreateTemp("", "testing-db")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := gorm.Open(sqlite.Open(dbFile.Name()), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DatabaseConn = db

	err = db.AutoMigrate(&models.PeerTorrent{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Torrent{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Peer{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Exec("PRAGMA foreign_keys = ON")

	tb.Cleanup(func() {
		if err := os.Remove(dbFile.Name()); err != nil {
			log.Fatalln(err)
		}
	})
}
