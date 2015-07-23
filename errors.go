package upflare

import (
	"errors"
)

var (
	ErrInvalidURL       = errors.New("upflare: invalid url")
	ErrInvalidHash      = errors.New("upflare: invalid hash")
	ErrInvalidUpload    = errors.New("upflare: invalid upload")
	ErrInvalidSignature = errors.New("upflare: invalid signature")
)
