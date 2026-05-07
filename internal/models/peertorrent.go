package models

import (
	"errors"

	"gorm.io/gorm"
)

type Events int

const (
	Started Events = iota
	Completed
	Stopped
	Empty
	Downloaded
)

type PeerTorrent struct {
	gorm.Model
	PeerId     uint
	TorrentId  uint
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      Events
}

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
