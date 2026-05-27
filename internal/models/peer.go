package models

import "gorm.io/gorm"

// Peer is the struct used for database operations
// on peers.
type Peer struct {
	gorm.Model
	ClientId string `gorm:"unique;not null"`
	Ip       string `gorm:"not null"`
	Port     int    `gorm:"not null"`
	Swarms   []PeerTorrent
}

// PeerDto is the DTO struct for Peer.
type PeerDto struct {
	PeerId string `bencode:"peer id"`
	Ip     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

// ToDto return a corresponding PeerDto for
// the Peer.
func (p *Peer) ToDto() PeerDto {
	return PeerDto{
		PeerId: p.ClientId,
		Ip:     p.Ip,
		Port:   p.Port,
	}
}
