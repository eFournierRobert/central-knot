package models

import "gorm.io/gorm"

type Torrent struct {
	gorm.Model
	InfoHash []byte
	Peers    []PeerTorrent
}
