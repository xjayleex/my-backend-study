package service

import "golang.org/x/net/context"

type Client interface {
	TransferImageFile(ctx context.Context, path string) (stats Stats, err error)
	Close()
}
