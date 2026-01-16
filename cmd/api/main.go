package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zenkiet/zen-attendance-server/config"
	"zenkiet/zen-attendance-server/internal/server"
	"zenkiet/zen-attendance-server/pkg/database"
)

// @contact.name	API Support
// @contact.url	http://www.zenkiet.com/support
// @contact.email
// @title           Zen Attendance API
// @version         1.0
// @description     This is the API documentation for the Zen Attendance system.s
// @termsOfService  https://github.com/zenkiet/zen-attendace

// @contact.name   API Support
// @contact.url    https://github.com/zenkiet
// @contact.email  kietgolx65234@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB)
	db, err := database.NewPostgres(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrationDir := "./migrations"
	if err := database.MigrateDB(connString, migrationDir); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb, err := database.NewRedis(ctx, redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("Failed to close redis client: %v", err)
		}
	}()

	srv := server.New(cfg, db, rdb)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
