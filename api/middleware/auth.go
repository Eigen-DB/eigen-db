package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// simple API key authentication middleware
func AuthMiddleware(ctx context.Context, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-Eigen-API-Key")
		if apiKey == "" {
			c.String(http.StatusUnauthorized, "No API key provided.")
			c.Abort()
			return
		}

		val, err := redisClient.Get(ctx, "apiKey").Result()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			fmt.Println(err.Error())
			return
		}

		if val != apiKey {
			c.String(http.StatusUnauthorized, "Invalid API key.")
			c.Abort()
			return
		}
	}
}
