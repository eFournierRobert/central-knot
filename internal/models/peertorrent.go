package models

import (
	"errors"

	"gorm.io/gorm"
)

// Events is an enum for all the possible events
// mentioned in BEP 3.
type Events int

const (
	Started Events = iota
	Completed
	Stopped
	Empty
	Downloaded
)

// PeerTorrent is the struct used for database operations
// on rows of peer_torrent.
type PeerTorrent struct {
	gorm.Model
	PeerId     uint    `gorm:"not null"`
	Peer       Peer    `gorm:"foreignKey:PeerId"`
	TorrentId  uint    `gorm:"not null"`
	Torrent    Torrent `gorm:"foreignKey:TorrentId"`
	Uploaded   int64   `gorm:"not null"`
	Downloaded int64   `gorm:"not null"`
	Left       int64   `gorm:"not null"`
	Event      Events  `gorm:"not null"`
}

// ParseEvent takes a string and returns the corresponding
// item in Events.
func ParseEvent(event string) (Events, error) {
	switch event {
	case "started":
		return Started, nil
	case "completed":
		return Completed, nil
	case "stopped":
		return Stopped, nil
	case "empty":
		return Empty, nil
	case "downloaded":
		return Downloaded, nil
	default:
		return -1, errors.New("event doesn't exist")
	}
}
