package hmgo

import (
	"context"

	"github.com/Chu16537/gomodule/mmgo"
)

type Handler struct {
	ctx    context.Context
	cancel context.CancelFunc
	read   *mmgo.Handler
	write  *mmgo.Handler
}

func New(ctx context.Context, cancel context.CancelFunc, read *mmgo.Handler, write *mmgo.Handler) (*Handler, error) {
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
