package middleware

import (
	"eigen_db/api/utils"
	"eigen_db/constants"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// simple API key authentication middleware
func AuthMiddleware() gin.HandlerFunc {
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

		validKey := os.Getenv(constants.ENV_VAR_API_KEY_NAME)
		if validKey != apiKey {
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
