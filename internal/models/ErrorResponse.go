package models

// ErrorResponseDto is the DTO struct that is used
// for returning an error message.
type ErrorResponseDto struct {
	FailureReason string `bencode:"failure reason"`
}
