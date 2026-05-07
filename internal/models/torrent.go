package models

import "gorm.io/gorm"

// Torrent is the struct used for database operations
// on torrents.
type Torrent struct {
	gorm.Model
	InfoHash []byte
	Peers    []PeerTorrent
}
