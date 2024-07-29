package hmgo

import (
	"context"

	"github.com/Chu16537/gomodule/zmongo"
)

type Handler struct {
	ctx    context.Context
	cancel context.CancelFunc
	read   *zmongo.Handler
	write  *zmongo.Handler
}

func New(ctx context.Context, cancel context.CancelFunc, read *zmongo.Handler, write *zmongo.Handler) (*Handler, error) {
	return &Handler{
		ctx:    ctx,
		cancel: cancel,
		read:   read,
		write:  write,
	}, nil
}

func (h *Handler) trans() *Handler {
	return &Handler{
		read:  h.write,
		write: h.write,
	}
}
