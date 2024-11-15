package hmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mmgo"
)

type Handler struct {
	ctx   context.Context
	read  *mmgo.Handler
	write *mmgo.Handler
}

func New(ctx context.Context, read *mmgo.Handler, write *mmgo.Handler) (*Handler, *errorcode.Error) {
	h := &Handler{
		ctx:   ctx,
		read:  read,
		write: write,
	}

	err := h.Init()

	return h, err
}

func (h *Handler) Done() {
	h.read.Done()
	h.write.Done()
}

func (h *Handler) trans() *Handler {
	return &Handler{
		read:  h.write,
		write: h.write,
	}
}
