package repo_peer

import (
	"central-knot/internal/models"
	"central-knot/internal/repository/db_utils"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTest(tb testing.TB) {
	dbFile, err := os.CreateTemp("", "testing-db")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := gorm.Open(sqlite.Open(dbFile.Name()), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db_utils.DatabaseConn = db

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

	tb.Cleanup(func() {
		if err := os.Remove(dbFile.Name()); err != nil {
			log.Fatalln(err)
		}
	})
}

func TestAddNewPeer(t *testing.T) {
	setupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	peer, err := Get(id)
	assert.Nil(t, err)
	assert.Equal(t, id, peer.ClientId)
	assert.Equal(t, ip, peer.Ip)
	assert.Equal(t, port, peer.Port)
}

func TestAddExistingPeer(t *testing.T) {
	setupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	err = Add(id, ip, port)
	assert.NotNil(t, err)
}

func TestGetNonExistingPeer(t *testing.T) {
	setupTest(t)

	id := "newid"

	peer, err := Get(id)
	assert.NotNil(t, err)
	assert.Empty(t, peer)
}

func TestChangeIpOnExistingPeer(t *testing.T) {
	setupTest(t)

	id := "newid"
	ip := "10.0.0.10"
	port := 2425

	err := Add(id, ip, port)
	assert.Nil(t, err)

	newIp := "10.10.12.12"
	err = ChangeIp(id, newIp)
	assert.Nil(t, err)

	peer, err := Get(id)
	assert.Nil(t, err)
	assert.Equal(t, id, peer.ClientId)
	assert.Equal(t, newIp, peer.Ip)
}

func TestChangeIpOnNoneExistingPeer(t *testing.T) {
	setupTest(t)

	id := "newid"
	newIp := "10.10.12.12"
	err := ChangeIp(id, newIp)
	assert.NotNil(t, err)

	_, err = Get(id)
	assert.NotNil(t, err)
}
