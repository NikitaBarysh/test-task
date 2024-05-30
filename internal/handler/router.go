package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
	"test-task/internal/config"
	"test-task/internal/service"
)

type Handler struct {
	service *service.Service
	cache   *cache.Cache
}

func NewHandler(newService *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: newService,
		cache:   cache.New(cfg.CacheDefExp, cfg.CacheCleanInterval),
	}
}

func (h *Handler) Register(router *chi.Mux) {
	router.Use(h.Middleware)
	router.Get("/get/{image}", h.handle)
}
