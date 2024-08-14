package health_check

import (
	"context"
	"eigen_db/redis_utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var status_code int
var status string

func markAsUnhealthy(err error) {
	status_code = http.StatusInternalServerError
	status = "unhealthy"
	fmt.Println(err.Error())
}

func Health(c *gin.Context) {
	status_code = http.StatusOK
	status = "healthy"

	// check Redis connection
	ctx := context.Background()
	redisClient, err := redis_utils.GetConnection(ctx)
	defer func() {
		if redisClient != nil { // if connection fails, redisClient = nil -> redisClient.Close() causes nil pointer dereference
			if err := redisClient.Close(); err != nil {
				markAsUnhealthy(err)
			}
		}
	}()
	if err != nil {
		markAsUnhealthy(err)
	}

	c.JSON(status_code, gin.H{
		"status": status,
	})
}
