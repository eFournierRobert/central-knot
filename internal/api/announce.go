package api

import (
	"central-knot/internal/models"
	"central-knot/internal/service"
	"log"
	"net"
	"net/http"

	"github.com/jackpal/bencode-go"
)

func SetUpAnnounceEndpoint() {
	http.HandleFunc("/announce", func(writer http.ResponseWriter, req *http.Request) {
		query := req.URL.RawQuery
		log.Println(query)

		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		if ip[0] == '[' {
			writer.WriteHeader(http.StatusBadRequest)
			if err := bencode.Marshal(writer, models.ErrorResponseDto{FailureReason: "ipv6 not supported"}); err != nil {
				log.Println(err)
			}
			return
		}

		dto, err := models.NewPeerDto(query)

		writer.Header().Set("Content-Type", "text/plain")
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			if err := bencode.Marshal(writer, models.ErrorResponseDto{FailureReason: "empty request"}); err != nil {
				log.Println(err)
			}
			return
		}

		peerList, err := service.Announce(dto, ip)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if err := bencode.Marshal(writer, models.ErrorResponseDto{FailureReason: err.Error()}); err != nil {
				log.Println(err)
			}
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
	})
}
