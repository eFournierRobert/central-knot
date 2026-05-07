package service

import (
	"central-knot/internal/models"
	repo_peer "central-knot/internal/repository/peer"
	repo_peerTorrent "central-knot/internal/repository/peerTorrent"
	repo_torrent "central-knot/internal/repository/torrent"
	"errors"
	"strconv"
)

const requestInterval = 300 //5 minutes

func Announce(dto *models.FullPeerDto, ip string) (*models.PeerListDto, error) {
	peer, err := repo_peer.GetPeer(dto.PeerId)
	if err != nil {
		if ip[0] == '[' {
			return nil, errors.New("tracker doesn't support ipv6")
		}
		port, err := strconv.Atoi(dto.Port)
		if err != nil {
			return nil, err
		}

		if err := repo_peer.AddPeer(dto.PeerId, ip, port); err != nil {
			return nil, err
		}

		peer, err = repo_peer.GetPeer(dto.PeerId)
		if err != nil {
			return nil, err
		}
	} else if peer.Ip != ip {
		if err := repo_peer.ChangeIp(peer.ClientId, ip); err != nil {
			return nil, err
		}
	}

	infoHashBytes := []byte(dto.InfoHash)
	torrent, err := repo_torrent.GetTorrent(infoHashBytes)
	if err != nil {
		if err := repo_torrent.AddTorrent(infoHashBytes); err != nil {
			return nil, err
		}

		torrent, err = repo_torrent.GetTorrent(infoHashBytes)
		if err != nil {
			return nil, err
		}
	}

	event, err := models.ParseEvent(dto.Event)
	if err != nil {
		event = models.Started
	}

	uploaded, downloaded, left, err := parsePeerTorrentStatusValues(dto.Uploaded, dto.Downloaded, dto.Left)
	if err != nil {
		return nil, err
	}

	peerTorrent, err := repo_peerTorrent.Get(peer.ID, torrent.ID)
	if err != nil {
		if err := repo_peerTorrent.Add(
			peer.ID,
			torrent.ID,
			int64(uploaded),
			int64(downloaded),
			int64(left),
			event); err != nil {
			return nil, err
		}
	} else {
		peerTorrent.Uploaded = int64(uploaded)
		peerTorrent.Downloaded = int64(downloaded)
		peerTorrent.Left = int64(left)
		peerTorrent.Event = event

		repo_peerTorrent.Update(&peerTorrent)
	}

	peerListDto := models.PeerListDto{Interval: requestInterval, Peers: make([]models.PeerDto, 0)}
	peers := repo_torrent.GetPeersFromHash(infoHashBytes, peer.ClientId)

	for _, p := range peers {
		peerListDto.Peers = append(peerListDto.Peers, p.ToDto())
	}

	return &peerListDto, nil
}

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
