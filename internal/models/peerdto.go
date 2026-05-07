package models

import (
	"errors"
	"net/url"
)

type FullPeerDto struct {
	InfoHash   string
	PeerId     string
	Port       string
	Uploaded   string
	Downloaded string
	Left       string
	Event      string
	Compact    bool
}

func NewPeerDto(query string) (*FullPeerDto, error) {
	values, _ := url.ParseQuery(query)

	if err := dtoValidation(values); err != nil {
		return nil, err
	}

	return &FullPeerDto{
		InfoHash:   values.Get("info_hash"),
		PeerId:     values.Get("peer_id"),
		Port:       values.Get("port"),
		Uploaded:   values.Get("uploaded"),
		Downloaded: values.Get("downloaded"),
		Left:       values.Get("left"),
		Event:      values.Get("event"),
		Compact:    values.Get("compact") == "1",
	}, nil
}

func dtoValidation(values url.Values) error {
	if len(values) == 0 {
		return errors.New("empty params")
	}

	if !values.Has("info_hash") {
		return errors.New("invalid info hash")
	}
	if !values.Has("peer_id") {
		return errors.New("invalid peer id")
	}
	if !values.Has("port") {
		return errors.New("invalid port")
	}
	if !values.Has("uploaded") || !values.Has("downloaded") || !values.Has("left") {
		return errors.New("invalid transfer status")
	}

	return nil
}
