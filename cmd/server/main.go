package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"auth-server/internal/config"
	httpdelivery "auth-server/internal/delivery/http"
	"auth-server/internal/repository"
	postgresRepo "auth-server/internal/repository/postgres"
	sqliteRepo "auth-server/internal/repository/sqlite"
	"auth-server/internal/service"
	"auth-server/internal/utils"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, cfg.Database.Type); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	var userRepo repository.UserRepository
	switch cfg.Database.Type {
	case "sqlite":
		userRepo = sqliteRepo.NewUserRepositorySQLite(db)
	case "postgres":
		userRepo = postgresRepo.NewUserRepositoryPostgres(db)
	default:
		log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
	}

	jwtService := utils.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	authService := service.NewAuthService(userRepo, jwtService)
	authHandlers := httpdelivery.NewAuthHandlers(authService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	httpdelivery.SetupRoutes(router, authHandlers, jwtService)

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	switch cfg.Database.Type {
	case "sqlite":
		return sql.Open("sqlite3", cfg.Database.SQLite.Path)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.DBName,
			cfg.Database.Postgres.SSLMode,
		)
		return sql.Open("pgx", dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
}

func runMigrations(db *sql.DB, dbType string) error {
	var query string
	if dbType == "sqlite" {
		query = `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            created_at DATETIME NOT NULL
        );
        CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
        `
	} else {
		query = `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL
        );
        CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
        `
	}
	_, err := db.Exec(query)
	return err
}
