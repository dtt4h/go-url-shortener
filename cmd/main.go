package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/dtt4h/go-url-shortener/internal/config"
	"github.com/dtt4h/go-url-shortener/internal/handler"
	"github.com/dtt4h/go-url-shortener/internal/logger"
	"github.com/dtt4h/go-url-shortener/internal/middleware"
	urlRepo "github.com/dtt4h/go-url-shortener/internal/repository/url"
	"github.com/dtt4h/go-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log, err := newLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		return errors.New("failed to connect to database")
	}
	defer pool.Close()

	urlRepository := urlRepo.NewURLRepository(pool)
	eventService := service.NewEventService(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer eventService.Close()

	urlService := service.NewURLService(urlRepository, eventService)
	urlHandler := handler.NewURLHandler(urlService, cfg.URL.Base, log)

	router := gin.New()

	router.Use(middleware.CORS())
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.RateLimiter())

	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join("web", "index.html"))
	})

	api := router.Group("/api/v1")
	{
		api.POST("/shorten", urlHandler.Create)
	}

	router.GET("/:code", urlHandler.Get)
	router.DELETE("/:code", urlHandler.Delete)
	router.GET("/:code/qr", urlHandler.GetQRCode)

	log.Info(context.Background(), "Server starting", "address", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, router); err != nil {
		log.Error(context.Background(), "server error", "error", err.Error())
		return errors.New("failed to start server")
	}

	return nil
}

func newLogger(cfg *config.Config) (*logger.Logger, error) {
	level, err := logger.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}
	return logger.New(nil, level), nil
}
