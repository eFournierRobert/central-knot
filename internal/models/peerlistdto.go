package models

import (
	"encoding/binary"
	"strconv"
	"strings"
)

type PeerListDto struct {
	Interval int       `bencode:"interval"`
	Peers    []PeerDto `bencode:"peers"`
}

type CompactPeerListDto struct {
	Interval int    `bencode:"interval"`
	Peers    []byte `bencode:"peers"`
}

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
