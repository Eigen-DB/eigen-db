package health_check

import (
	"eigen_db/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestAuth(c *gin.Context) {
	// if a request makes it here, then it has been authenticated by the middleware
	utils.SendResponse(
		c,
		http.StatusOK,
		"Authenticated.",
		nil,
		nil,
	)
}
