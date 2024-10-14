package health_check

import (
	"eigen_db/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
TODO:
- Check for unusual resource usage
- Check that all endpoints function correctly
- Return DB uptime
*/
func Health(c *gin.Context) {
	utils.SendResponse(
		c,
		http.StatusOK,
		"healthy",
		nil,
		nil,
	)
}
