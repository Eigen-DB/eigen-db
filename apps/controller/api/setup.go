package api

import (
	"controller/api/handlers/instance"
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartAPI(devMode bool) error {
	var port int = 8080
	if devMode {
		port = 1337
	}
	router := gin.Default()

	instanceEndpoints := router.Group("/instance")
	instanceEndpoints.POST("/create", instance.CreateInstance)
	instanceEndpoints.DELETE("/terminate", instance.TerminateInstance)

	return router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}
