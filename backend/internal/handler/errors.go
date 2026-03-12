package handler

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrLinkExpired    = errors.New("link expired")
	ErrInternal       = errors.New("internal server error")
)
