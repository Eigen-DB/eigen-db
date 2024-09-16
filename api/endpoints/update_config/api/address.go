package api

import (
	"eigen_db/cfg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateAddressBody struct {
	UpdatedAddress string `json:"updatedAddress" binding:"required"`
}

func UpdateAddress(config cfg.IConfig) func(*gin.Context) {
	return func(c *gin.Context) {
		var body updateAddressBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		_ = config.SetAPIAddress(body.UpdatedAddress) // handle error
		c.String(http.StatusOK, "API address updated. Please restart the database for it to take effect.")
	}
}
