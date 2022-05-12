package handlers

import (
	"github.com/gorilla/mux"

	v1 "finstat/internal/handlers/http/v1"
	"finstat/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *mux.Router {
	return h.initHand()
}

func (h *Handler) initHand() *mux.Router {
	initHandlers := v1.NewHandler(h.services)
	return initHandlers.Init()
}
