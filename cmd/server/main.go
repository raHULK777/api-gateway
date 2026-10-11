package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/repository"
	"api-gateway/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Initialize Redis
	redisClient, err := repository.NewRedisClient(cfg.RedisURL)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize PostgreSQL
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()

	// Initialize services
	keyStore := repository.NewKeyStore(db)
	tokenService := service.NewTokenService(cfg.JWTSecret)
	rateLimiter := middleware.NewRateLimiter(redisClient, cfg.RateLimitRPS)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().Unix()})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register(keyStore, tokenService))
			auth.POST("/login", handler.Login(keyStore, tokenService))
		}

		// Protected routes with rate limiting and auth
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(tokenService))
		protected.Use(rateLimiter.Middleware())
		{
			protected.GET("/users/me", handler.GetMe(keyStore))
			protected.GET("/users/:id", handler.GetUser(keyStore))
			protected.PUT("/users/:id", func(c *gin.Context) {})
			protected.DELETE("/users/:id", func(c *gin.Context) {})

			// Example resource routes
			protected.GET("/items", handler.ListItems(keyStore))
			protected.POST("/items", handler.CreateItem(keyStore, tokenService))
			protected.GET("/items/:id", handler.GetItem(keyStore))
			protected.PUT("/items/:id", handler.UpdateItem(keyStore, tokenService))
			protected.DELETE("/items/:id", handler.DeleteItem(keyStore))
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting API Gateway", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
}
