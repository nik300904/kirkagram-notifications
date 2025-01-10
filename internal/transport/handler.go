package transport

import (
	"github.com/go-chi/chi"
	"kirkagram-notification/internal/transport/ws"
)

type Handler struct {
	wsHandler *ws.WebSocketManager
}

func NewHandler(wsHandler *ws.WebSocketManager) *Handler {
	return &Handler{
		wsHandler: wsHandler,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/ws/{userID}", h.wsHandler.HandleWebSocket)

	return router
}
