package protocol

import "errors"

var (
	ErrMissingHostHeader = errors.New("missing X-Meili-Instance header")
	ErrInvalidTransport  = errors.New("invalid transport")
	ErrNilPool           = errors.New("nil pool for HTTP transport")
)
