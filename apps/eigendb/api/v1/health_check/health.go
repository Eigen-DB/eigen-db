package health_check

import (
	"eigen_db/api/utils"
	"eigen_db/metrics"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TODO:
- Check for unusual resource usage
- Check that all endpoints function correctly
*/
func Health(c *gin.Context) {
	uptime := metrics.GetUptime()
	memUsage, err := metrics.GetMemUsage()
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"unhealthy",
			nil,
			utils.CreateError(
				"ERROR_GETTING_MEM_USAGE",
				err.Error(),
			),
		)
	}
	cpuUsage, err := metrics.GetCpuUsage()
	if err != nil {
		utils.SendResponse(
			c,
			http.StatusInternalServerError,
			"unhealthy",
			nil,
			utils.CreateError(
				"ERROR_GETTING_CPU_USAGE",
				err.Error(),
			),
		)
	}

	utils.SendResponse(
		c,
		http.StatusOK,
		"healthy",
		map[string]any{
			"uptime":            uptime.String(),
			"cpu_usage_percent": cpuUsage,
			"mem_usage_percent": memUsage,
		},
		nil,
	)
}
