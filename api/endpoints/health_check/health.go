package health_check

import (
	"context"
	"eigen_db/redis_utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	status_code := http.StatusOK
	status := "healthy"

	// check Redis connection
	ctx := context.Background()
	redisClient, err := redis_utils.GetConnection(ctx)
	if err != nil {
		fmt.Println(err.Error())
		status_code = http.StatusInternalServerError
		status = "unhealthy"
	}
	if err = redis_utils.CheckConnection(ctx, redisClient); err != nil {
		fmt.Println(err.Error())
		status_code = http.StatusInternalServerError
		status = "unhealthy"
	}

	c.JSON(status_code, gin.H{
		"status": status,
	})
}
