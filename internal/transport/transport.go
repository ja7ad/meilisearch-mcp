package transport

import "context"

type Server interface {
	Start(ctx context.Context)
	Err() <-chan error
	Stop(ctx context.Context) error
}
