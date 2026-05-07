package api

import (
	"central-knot/internal/models"
	"central-knot/internal/service"
	"log"
	"net/http"
	"strings"

	"github.com/jackpal/bencode-go"
)

func SetUpAnnounceEndpoint() {
	http.HandleFunc("/announce", func(writer http.ResponseWriter, req *http.Request) {
		query := req.URL.RawQuery
		log.Println(query)

		ip := strings.Split(req.RemoteAddr, ":")[0]

		if ip == "127.0.0.1" {
			ip = "10.0.0.59"
		}

		if len(query) != 0 {
			dto := models.NewPeerDto(query)

			peerList, err := service.Announce(&dto, ip)
			if err != nil {
				if err := bencode.Marshal(writer, models.ErrorResponseDto{FailureReason: err.Error()}); err != nil {
					log.Println(err)
				}
				return
			}

			if dto.Compact {
				writer.WriteHeader(http.StatusOK)
				compactPeerList := peerList.ToCompact()
				if err := bencode.Marshal(writer, compactPeerList); err != nil {
					log.Println(err)
				}
			} else {
				writer.WriteHeader(http.StatusOK)
				if err := bencode.Marshal(writer, *peerList); err != nil {
					log.Println(err)
				}
			}
		}
	})
}
