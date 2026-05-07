package models

import (
	"errors"
	"net/url"
)

// FullPeerDto is the DTO struct where all
// the query parameters of an announcement by a client are parsed into.
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

// NewPeerDto parses the given query parameters into a FullPeerDto.
func NewPeerDto(query string) (*FullPeerDto, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

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

// dtoValidation validates if the given url values contains all
// the necessary parameters mentioned in BEP 3.
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
