package health_check

import (
	"context"
	"eigen_db/api/utils"
	"eigen_db/redis_utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	// check Redis connection
	ctx := context.Background()
	redisClient, err := redis_utils.GetConnection(ctx)
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"unhealthy",
			nil,
			utils.CreateError("ERROR_CONNECTING_TO_REDIS", err.Error()),
		)
		return
	} else {
		if err := redisClient.Close(); err != nil {
			utils.SendResponse(
				c,
				http.StatusInternalServerError,
				"unhealthy",
				nil,
				utils.CreateError("ERROR_CLOSING_REDIS_TEST_CONNECTION", err.Error()),
			)
			return
		}
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"healthy",
		nil,
		nil,
	)
}
