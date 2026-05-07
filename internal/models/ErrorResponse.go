package models

type ErrorResponseDto struct {
	FailureReason string `bencode:"failure reason"`
}
