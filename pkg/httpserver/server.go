package httpserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	mid "aviation-service/pkg/middleware"	
)

type Handler interface {
	RegisterRoutes(r chi.Router)
}

type Server struct {
	server *http.Server
}

func NewRouter(handlers ...Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(mid.ZapLogger)
	r.Use(middleware.Recoverer)

	for _, h := range handlers {
		h.RegisterRoutes(r)
	}

	return r
}

func NewServer(router *chi.Mux, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}