package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zenkiet/zen-attendance-server/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router *chi.Mux
	db     *pgxpool.Pool
	rdb    *redis.Client
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
		rdb:    rdb,
		cfg:    cfg,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
}

func (s *Server) setupRoutes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: s.router,
	}

	fmt.Printf("Server starting on port %s\n", s.cfg.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
