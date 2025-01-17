package serverless

import (
	"net/http"

	"github.com/linden/rpc"
)

type Handler struct {
	srv *rpc.Server
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.srv.ServeCodec(rpc.NewGobServerCodec(r.Body, w))
}

func NewHandler(srv *rpc.Server) *Handler {
	return &Handler{
		srv: srv,
	}
}
