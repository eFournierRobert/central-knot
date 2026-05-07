package service

import (
	"central-knot/internal/models"
	repo_peer "central-knot/internal/repository/peer"
	repo_peerTorrent "central-knot/internal/repository/peerTorrent"
	repo_torrent "central-knot/internal/repository/torrent"
	"strconv"
)

// requestInterval is the suggested interval
// between announcements of Central Knot. It
// is 300 seconds or 5 minutes.
const requestInterval = 300

// Announce takes in the query DTO and the IP address of the client
// that is making the announcement and returns the DTO struct of peers in
// the swarm.
func Announce(dto *models.FullPeerDto, ip string) (*models.PeerListDto, error) {
	peer, err := ensurePeer(dto.PeerId, ip, dto.Port)
	if err != nil {
		return nil, err
	}

	infoHashBytes := []byte(dto.InfoHash)
	torrent, err := ensureTorrent(infoHashBytes)
	if err != nil {
		return nil, err
	}

	event, err := models.ParseEvent(dto.Event)
	if err != nil {
		event = models.Started
	}

	if err := upsertPeerTorrent(peer.ID, torrent.ID, dto, event); err != nil {
		return nil, err
	}

	return buildPeerListDto(infoHashBytes, peer.ClientId)
}

// buildPeerListDto makes the request for peers in a given swarm and builds the
// PeerListDto with them.
func buildPeerListDto(infoHash []byte, clientId string) (*models.PeerListDto, error) {
	peerListDto := models.PeerListDto{Interval: requestInterval, Peers: make([]models.PeerDto, 0)}
	peers := repo_torrent.GetPeersFromHash(infoHash, clientId)

	for _, p := range peers {
		peerListDto.Peers = append(peerListDto.Peers, p.ToDto())
	}

	return &peerListDto, nil
}

// upsertPeerTorrent updates the PeerTorrent row for a given peer ID (as in primary key) and torrent ID
// with the new data in the announcement DTO.
func upsertPeerTorrent(peerId, torrentId uint, dto *models.FullPeerDto, event models.Events) error {
	uploaded, downloaded, left, err := parsePeerTorrentStatusValues(dto.Uploaded, dto.Downloaded, dto.Left)
	if err != nil {
		return err
	}

	peerTorrent, err := repo_peerTorrent.Get(peerId, torrentId)
	if err != nil {
		if err := repo_peerTorrent.Add(
			peerId,
			torrentId,
			int64(uploaded),
			int64(downloaded),
			int64(left),
			event); err != nil {
			return err
		}
	} else {
		peerTorrent.Uploaded = int64(uploaded)
		peerTorrent.Downloaded = int64(downloaded)
		peerTorrent.Left = int64(left)
		peerTorrent.Event = event

		repo_peerTorrent.Update(&peerTorrent)
	}

	return nil
}

// ensureTorrent returns the Torrent that has the given info hash in the
// database. If it isn't found, it will create it and return the created row.
func ensureTorrent(infoHash []byte) (models.Torrent, error) {
	torrent, err := repo_torrent.Get(infoHash)
	if err != nil {
		if err := repo_torrent.Add(infoHash); err != nil {
			return models.Torrent{}, err
		}

		torrent, err = repo_torrent.Get(infoHash)
		if err != nil {
			return models.Torrent{}, err
		}
	}

	return torrent, nil
}

// ensurePeer returns the Peer that has the given peer ID in the
// database. If it isn't found, it will create it and return the created row.
func ensurePeer(peerId, ip, port string) (models.Peer, error) {
	peer, err := repo_peer.GetPeer(peerId)
	if err != nil {
		port, err := strconv.Atoi(port)
		if err != nil {
			return models.Peer{}, err
		}

		if err := repo_peer.AddPeer(peerId, ip, port); err != nil {
			return models.Peer{}, err
		}

		peer, err = repo_peer.GetPeer(peerId)
		if err != nil {
			return models.Peer{}, err
		}
	} else if peer.Ip != ip {
		if err := repo_peer.ChangeIp(peer.ClientId, ip); err != nil {
			return models.Peer{}, err
		}
	}

	return peer, nil
}

// parsePeerTorrentStatusValues parses the strings for uploaded, downloaded and left in the
// announcement DTO and returns a tuple of integers (uploaded, downloaded, left).
func parsePeerTorrentStatusValues(uploadedStr, downloadedStr, leftStr string) (int, int, int, error) {
	uploaded, err := strconv.Atoi(uploadedStr)
	if err != nil {
		return 0, 0, 0, err
	}

	downloaded, err := strconv.Atoi(downloadedStr)
	if err != nil {
		return 0, 0, 0, err
	}

	left, err := strconv.Atoi(leftStr)
	if err != nil {
		return 0, 0, 0, err
	}

	return uploaded, downloaded, left, nil
}
