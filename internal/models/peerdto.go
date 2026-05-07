package models

import (
	"strings"
)

type FullPeerDto struct {
	InfoHash   string
	PeerId     string
	Port       string
	Uploaded   string
	Downloaded string
	Left       string
	Corrupt    string
	Key        string
	Event      string
	Numwant    string
	Compact    bool
	NoPeerId   bool
}

func NewPeerDto(query string) FullPeerDto {
	splitString := strings.Split(query, "&")
	values := make(map[string]string)

	for _, s := range splitString {
		pair := strings.Split(s, "=")
		values[pair[0]] = pair[1]
	}

	return FullPeerDto{
		InfoHash:   values["info_hash"],
		PeerId:     values["peer_id"],
		Port:       values["port"],
		Uploaded:   values["uploaded"],
		Downloaded: values["downloaded"],
		Left:       values["left"],
		Corrupt:    values["corrupt"],
		Key:        values["key"],
		Event:      values["event"],
		Numwant:    values["numwant"],
		Compact:    values["compact"] == "1",
		NoPeerId:   values["no_peer_id"] == "1",
	}
}
