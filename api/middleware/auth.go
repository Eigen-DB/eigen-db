package middleware

import (
	"context"
	"eigen_db/api/utils"
	"eigen_db/constants"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// simple API key authentication middleware
func AuthMiddleware(ctx context.Context, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get(constants.MIDDLEWARE_API_KEY_HEADER)
		if apiKey == "" {
			utils.SendResponse(
				c,
				http.StatusUnauthorized,
				"No API key provided.",
				nil,
				utils.CreateError("NO_API_KEY_PROVIDED", "A valid API key is required to access this endpoint."),
			)
			c.Abort()
			return
		}

		val, err := redisClient.Get(ctx, "apiKey").Result()
		if err != nil {
			utils.SendResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error.",
				nil,
				nil,
			)
			c.Abort()
			fmt.Println(err.Error())
			return
		}

		if val != apiKey {
			utils.SendResponse(
				c,
				http.StatusUnauthorized,
				"Invalid API key.",
				nil,
				utils.CreateError("INVALID_API_KEY", "The API key you provided is invalid."),
			)
			c.Abort()
			return
		}
	}
}
