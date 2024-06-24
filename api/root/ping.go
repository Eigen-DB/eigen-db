package root_endpoints

import "github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
	c.String(200, "Pong!")
}
