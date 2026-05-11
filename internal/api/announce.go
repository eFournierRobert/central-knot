package api

import (
	"central-knot/internal/models"
	"central-knot/internal/service"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/jackpal/bencode-go"
)

// AnnounceHandler handles incoming announcement from clients.
func AnnounceHandler(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.RawQuery
	log.Println(query)

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	if ip[0] == '[' {
		writeError(writer, http.StatusBadRequest, errors.New("IPv6 not supported"))
		return
	}

	dto, err := models.NewPeerDto(query)

	if err != nil {
		writeError(writer, http.StatusBadRequest, errors.New("empty request"))
		return
	}

	peerList, err := service.Announce(dto, ip)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	if dto.Compact {
		compactPeerList := peerList.ToCompact()
		if err := bencode.Marshal(writer, compactPeerList); err != nil {
			log.Println(err)
		}
	} else {
		if err := bencode.Marshal(writer, *peerList); err != nil {
			log.Println(err)
		}
	}
}

func writeError(writer http.ResponseWriter, httpCode int, err error) {
	writer.WriteHeader(httpCode)
	if err := bencode.Marshal(writer, models.ErrorResponseDto{FailureReason: err.Error()}); err != nil {
		log.Println(err)
	}
}
