package health_check

import (
	"eigen_db/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	// perform a test call to /test-auth
	// ...

	utils.SendResponse(
		c,
		http.StatusOK,
		"healthy",
		nil,
		nil,
	)
}
