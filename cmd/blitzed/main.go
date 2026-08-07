package main

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mariogarzac/blitzed/config"
	"github.com/mariogarzac/blitzed/internal/db"
	"github.com/mariogarzac/blitzed/internal/match"
	"github.com/mariogarzac/blitzed/internal/player"
	"github.com/mariogarzac/blitzed/internal/stats"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	runMigrations(cfg.DatabaseURL)

	// Repos
	playerRepo := player.NewRepository(pool)
	matchRepo := match.NewRepository(pool)
	statsRepo := stats.NewRepository(pool)

	// Services
	matchSvc := match.NewService(matchRepo)

	// Controllers
	playerC := player.NewController(playerRepo)
	matchC := match.NewController(matchSvc, matchRepo)
	statsC := stats.NewController(statsRepo)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))

	api := e.Group("/api/v1")

	api.GET("/players", playerC.List)
	api.POST("/players", playerC.Create)
	api.DELETE("/players/:id", playerC.Delete)

	api.GET("/matches", matchC.List)
	api.POST("/matches", matchC.Create)
	api.GET("/matches/:id", matchC.GetDetail)
	api.POST("/matches/:id/rounds", matchC.SubmitRound)
	api.DELETE("/matches/:id/rounds/last", matchC.UndoLastRound)
	api.DELETE("/matches/:id", matchC.Delete)
	api.POST("/matches/:id/end", matchC.EndMatch)

	api.GET("/stats/head-to-head", statsC.HeadToHead)

	log.Printf("starting server on :%s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(databaseURL string) {
	migrateURL := "pgx5://" + databaseURL[len("postgres://"):]
	m, err := migrate.New("file://migrations", migrateURL)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")
}
