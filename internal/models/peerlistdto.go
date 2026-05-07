package models

import (
	"encoding/binary"
	"strconv"
	"strings"
)

// PeerListDto is the DTO struct for the Peer list
// being returned to a client after an announcement.
type PeerListDto struct {
	Interval int       `bencode:"interval"`
	Peers    []PeerDto `bencode:"peers"`
}

// CompactPeerListDto is the compact version of a
// PeerListDto as per BEP 23.
type CompactPeerListDto struct {
	Interval int    `bencode:"interval"`
	Peers    []byte `bencode:"peers"`
}

// ToCompact returns a corresponding CompactPeerListDto
// for a PeerListDto.
func (p *PeerListDto) ToCompact() CompactPeerListDto {
	var buf []byte
	for _, p := range p.Peers {
		ipv4Str := strings.Split(p.Ip, ".")
		for _, subInt := range ipv4Str {
			i, _ := strconv.Atoi(subInt)
			buf, _ = binary.Append(buf, binary.BigEndian, uint8(i))
		}
		buf, _ = binary.Append(buf, binary.BigEndian, uint16(p.Port))
	}

	return CompactPeerListDto{
		Interval: p.Interval,
		Peers:    buf,
	}
}
