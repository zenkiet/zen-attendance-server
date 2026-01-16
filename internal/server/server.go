package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"zenkiet/zen-attendance-server/config"
	_ "zenkiet/zen-attendance-server/docs"
	"zenkiet/zen-attendance-server/handler"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router *chi.Mux
	doc    huma.API
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

	humaCfg := huma.DefaultConfig("Zen Attendance API", "1.0.0")
	s.doc = humachi.New(s.router, humaCfg)

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
	h := handler.New(s.db, s.rdb)

	// Health Check
	huma.Register(s.doc, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Check Status Server",
		Tags:        []string{"System"},
	}, h.HealthCheck)

	// API v1 routes
	// s.router.Route("/api/v1", func(r chi.Router) {

	// })
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Server starting on port %s\n", s.cfg.Port)
	fmt.Printf("OpenAPI Doc: http://localhost:%s/docs\n", s.cfg.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
