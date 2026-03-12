package handlers

import (
	"controller/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type stopReqBody struct {
	CustomerId string `json:"customerId" binding:"required"`
}

func Stop(c *gin.Context) {
	var body stopReqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	jail, err := utils.JailFactory(body.CustomerId)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := jail.Stop(); err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "Jail stopped.")
}
