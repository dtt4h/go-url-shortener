package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		os.Exit(1)
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// Проверка соединения с БД
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	urlRepository := urlRepo.NewURLRepository(pool)

	var eventService service.EventService
	eventService = service.NewNoOpEventService()

	urlService := service.NewURLService(urlRepository, eventService)
	urlHandler := handler.NewURLHandler(urlService, cfg.URL.Base, log)

	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigin))
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.RateLimiter())

	router.GET("/health", handler.HealthCheck(pool))
	router.GET("/ready", handler.ReadyCheck(pool))

	router.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})

	api := router.Group("/api/v1")
	{
		api.POST("/shorten", urlHandler.Create)
	}

	router.GET("/:code", urlHandler.Get)
	router.DELETE("/:code", urlHandler.Delete)
	router.GET("/:code/qr", urlHandler.GetQRCode)

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info(context.Background(), "Server starting", "address", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(context.Background(), "server error", "error", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info(context.Background(), "Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info(context.Background(), "Server exited")
	return nil
}

func newLogger(cfg *config.Config) (*logger.Logger, error) {
	level, err := logger.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}
	return logger.New(nil, level), nil
}
