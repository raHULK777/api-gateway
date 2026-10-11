package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"api-gateway/internal/service"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type RateLimiter struct {
	client *redis.Client
	rps    float64
	window time.Duration
}

func NewRateLimiter(client *redis.Client, rps float64) *RateLimiter {
	return &RateLimiter{
		client: client,
		rps:    rps,
		window: time.Second,
	}
}

func (rl *RateLimiter) IsAllowed(c context.Context, key string) (bool, error) {
	redisKey := "rate_limit:" + key
	counter, err := rl.client.Incr(c, redisKey).Result()
	if err != nil {
		return false, err
	}

	if counter == 1 {
		rl.client.Expire(c, redisKey, rl.window)
	}

	return float64(counter) <= rl.rps, nil
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		allowed, err := rl.IsAllowed(ctx, clientIP)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "rate limit check failed"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}

func AuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("user_id", uint(userID))
		c.Set("email", (*claims)["email"])
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
