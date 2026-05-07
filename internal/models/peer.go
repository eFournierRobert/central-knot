package models

import "gorm.io/gorm"

type Peer struct {
	gorm.Model
	ClientId string
	Ip       string
	Port     int
	Swarms   []PeerTorrent
}

type PeerDto struct {
	PeerId string `bencode:"peer id"`
	Ip     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

func (p *Peer) ToDto() PeerDto {
	return PeerDto{
		PeerId: p.ClientId,
		Ip:     p.Ip,
		Port:   p.Port,
	}
}
